package common

import (
	"github.com/Cheeel666/midjourney-go/midjourney"
	"github.com/bwmarrin/discordgo"
	"midjourney-apiserver/internal/config"
	"midjourney-apiserver/pkg/store"
)

type Base struct {
	*discordgo.Session
	Store    *store.Store
	MJClient *midjourney.Client
	Config   *config.Config
}
