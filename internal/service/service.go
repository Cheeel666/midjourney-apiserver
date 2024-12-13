package service

import (
	"github.com/Cheeel666/midjourney-apiserver/internal/common"
	"github.com/Cheeel666/midjourney-apiserver/pkg/api"
)

type Service struct {
	api.UnimplementedAPIServiceServer
	discordSessionID string
	*common.Base
}

func New(base *common.Base, discordSessionID string) *Service {
	return &Service{
		discordSessionID: discordSessionID,
		Base:             base,
	}
}
