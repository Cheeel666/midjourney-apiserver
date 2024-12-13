package common

import (
	"github.com/Cheeel666/midjourney-apiserver/internal/config"
	"github.com/Cheeel666/midjourney-apiserver/pkg/store"
	"github.com/Cheeel666/midjourney-go/midjourney"
	"github.com/bwmarrin/discordgo"
)

type Base struct {
	*discordgo.Session
	Store    *store.Store
	MJClient *midjourney.Client
	Config   *config.Config
}
