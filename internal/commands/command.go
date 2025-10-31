package commands

import "github.com/bwmarrin/discordgo"

type Command interface {
    Execute(s *discordgo.Session, ic *discordgo.InteractionCreate, svc interface{})
}
