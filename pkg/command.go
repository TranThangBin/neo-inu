package pkg

import "github.com/bwmarrin/discordgo"

type Command interface {
	NewApplicationCommand() *discordgo.ApplicationCommand
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error
	NewResponse(params interface{}) *discordgo.InteractionResponse
}

var PublicCommandOption = &discordgo.ApplicationCommandOption{
	Type:        discordgo.ApplicationCommandOptionBoolean,
	Name:        "public",
	Description: "If you want to make this response public on the guild",
}
