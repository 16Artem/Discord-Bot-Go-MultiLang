package commands

import (
	"fmt"

	"example.com/bot/internal/services"
	"github.com/bwmarrin/discordgo"
)

type SettingsCommand struct {
	svc services.Service
}

func NewSettingsCommand(s services.Service) *SettingsCommand {
	return &SettingsCommand{svc: s}
}

func (c *SettingsCommand) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate, _ interface{}) {
	data := ic.ApplicationCommandData()
	guildID := ic.GuildID
	if len(data.Options) > 0 && data.Options[0].Name == "lang" {
		sub := data.Options[0]
		if len(sub.Options) > 0 {
			val := sub.Options[0].StringValue()
			if err := c.svc.SetGuildLang(guildID, val); err != nil {
				_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{Content: "Failed to save setting."},
				})
				return
			}
			msg := fmt.Sprintf(c.svc.Locale().T(val, "settings_desc_lang_changed"), val)
			_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: msg},
			})
			return
		}
	}
	cur := c.svc.GetGuildLang(guildID)
	msg := fmt.Sprintf(c.svc.Locale().T(cur, "settings_desc_current_lang"), cur)
	_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: msg},
	})
}
