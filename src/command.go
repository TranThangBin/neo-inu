package src

import "github.com/bwmarrin/discordgo"

type Executable interface {
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

type Command interface {
	Executable
	NewApplicationCommand() *discordgo.ApplicationCommand
}

type Subcommand interface {
	Executable
	NewApplicationCommandOption() *discordgo.ApplicationCommandOption
}
