package services

import (
	"example.com/bot/internal/config"
	"example.com/bot/internal/locale"
)

type Service interface {
	Config() *config.Config
	Locale() *locale.Loader
	GetGuildLang(guildID string) string
	SetGuildLang(guildID, lang string) error
}
