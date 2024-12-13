package webhook

import "github.com/Cheeel666/midjourney-apiserver/pkg/store"

type WebhookRequest struct {
	TaskID       string       `json:"task_id"`
	Prompt       string       `json:"prompt"`
	Type         store.Type   `json:"type"`
	Status       store.Status `json:"status"`
	Mode         string       `json:"mode"`
	ImageURL     string       `json:"image_url"`
	StartTime    int64        `json:"start_time"`
	CompleteTime int64        `json:"complete_time"`
}
