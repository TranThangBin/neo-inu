package internal

import "github.com/bwmarrin/discordgo"

type PingCommand struct{}

func (p *PingCommand) NewApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Say pong",
	}
}

func (p *PingCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func NewPingCommand() *PingCommand {
	return &PingCommand{}
}
