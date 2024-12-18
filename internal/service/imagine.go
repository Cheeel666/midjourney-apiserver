package service

import (
	"context"
	"fmt"
	"github.com/Cheeel666/midjourney-go/midjourney"
	"log"
	"time"

	"github.com/Cheeel666/midjourney-apiserver/pkg/api"
	"github.com/Cheeel666/midjourney-apiserver/pkg/store"
	"github.com/google/uuid"
)

/*
flow:
1. create mesasge id: 1
2. update message id: 1
3. create message id: 2 -> contains attachments
4. delete message id: 1
*/
func (s *Service) Imagine(ctx context.Context, in *api.ImagineRequest) (*api.ImagineResponse, error) {
	if in.RequestId == "" {
		in.RequestId = uuid.NewString()
	}

	if err := s.Store.CheckPrompt(ctx, in.Prompt); err != nil {
		e := err.(store.Error)
		return &api.ImagineResponse{
			RequestId: in.RequestId,
			Code:      e.Code,
			Msg:       e.Msg,
		}, nil
	}

	key := store.GetKey(in.Prompt)

	log.Printf("Imagine, key: %s, len: %d", key, len(key))

	if !KeyChan.Init(key) {
		return &api.ImagineResponse{
			RequestId: in.RequestId,
			Code:      api.Codes_CODES_INVALID_PARAMETER_ERROR,
			Msg:       fmt.Sprintf("The same prompt is being processed, please try again later."),
		}, nil
	}

	defer KeyChan.Del(key)

	if err := s.MJClient.Imagine(ctx, &midjourney.ImagineRequest{
		GuildID:   s.Config.Midjourney.GuildID,
		ChannelID: s.Config.Midjourney.ChannelID,
		Prompt:    in.Prompt,
		SessionID: s.discordSessionID,
	}); err != nil {
		return &api.ImagineResponse{
			RequestId: in.RequestId,
			Code:      api.Codes_CODES_SERVER_INTERNAL_ERROR,
			Msg:       fmt.Sprint(err),
		}, nil
	}

	select {
	case <-time.After(10 * time.Second):
		return &api.ImagineResponse{
			RequestId: in.RequestId,
			Code:      api.Codes_CODES_PROCESSING_TIMEOUT,
			Msg:       "timeout",
		}, nil
	case msgInfo := <-KeyChan.Get(key):
		if msgInfo.Error != nil {
			code := api.Codes_CODES_SERVER_INTERNAL_ERROR

			switch msgInfo.Error.Title {
			case "Invalid parameter":
				code = api.Codes_CODES_INVALID_PARAMETER_ERROR
			}

			return &api.ImagineResponse{
				RequestId: in.RequestId,
				Code:      code,
				Msg:       msgInfo.Error.Description,
			}, nil
		}

		if err := s.Store.SaveWebhook(ctx, msgInfo.ID, in.Webhook); err != nil {
			e := err.(store.Error)
			return &api.ImagineResponse{
				RequestId: in.RequestId,
				Code:      e.Code,
				Msg:       e.Msg,
			}, nil
		}

		return &api.ImagineResponse{
			RequestId: in.RequestId,
			Code:      api.Codes_CODES_SUCCESS,
			Msg:       "success",
			Data: &api.ImagineResponseData{
				TaskId:    msgInfo.ID,
				StartTime: msgInfo.StartTime,
			},
		}, nil
	}
}
