package internal

import "github.com/bwmarrin/discordgo"

type PingCommand struct{}

func (p *PingCommand) NewApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Say pong",
	}
}

func (p *PingCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, p.NewResponse(nil))
}

func (p *PingCommand) NewResponse(interface{}) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
}

func NewPingCommand() *PingCommand {
	return &PingCommand{}
}
