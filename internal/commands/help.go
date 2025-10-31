package commands

import (
    "github.com/bwmarrin/discordgo"
    "example.com/bot/internal/services"
)

type HelpCommand struct {
    svc services.Service
}

func NewHelpCommand(s services.Service) *HelpCommand {
    return &HelpCommand{svc: s}
}

func (c *HelpCommand) Execute(s *discordgo.Session, ic *discordgo.InteractionCreate, _ interface{}) {
    guildID := ic.GuildID
    lang := c.svc.GetGuildLang(guildID)
    title := c.svc.Locale().T(lang, "help_title")
    desc := c.svc.Locale().T(lang, "help_description")

    _ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Embeds: []*discordgo.MessageEmbed{{
                Title:       title,
                Description: desc,
            }},
        },
    })
}
