package src

import (
	"neo-inu/src/ygo"

	"github.com/bwmarrin/discordgo"
)

type YgoCommand struct {
	subcommands       map[string]Subcommand
	subcommandOptions []*discordgo.ApplicationCommandOption
}

func (yg *YgoCommand) NewApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ygo",
		Description: "Use ygoprodeck api to do cool stuff",
		Options:     yg.subcommandOptions,
	}
}

func (yg *YgoCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opt := i.ApplicationCommandData().Options[0]
	subCmdInteraction := *i.Interaction
	subCmdInteraction.Data = discordgo.ApplicationCommandInteractionData{
		Name:    opt.Name,
		Options: opt.Options,
	}
	return yg.subcommands[opt.Name].Execute(s, &discordgo.InteractionCreate{
		Interaction: &subCmdInteraction,
	})
}

func NewYgoCommand() *YgoCommand {
	subCommandList := []Subcommand{
		&ygo.YgoRandomSubcommand{},
		&ygo.YgoSearchSubcommand{},
	}

	subcommands := make(map[string]Subcommand, len(subCommandList))
	subcommandOptions := make([]*discordgo.ApplicationCommandOption, len(subCommandList))

	for i, subCommand := range subCommandList {
		subcommandOptions[i] = subCommand.NewApplicationCommandOption()
		subcommands[subcommandOptions[i].Name] = subCommand
	}

	return &YgoCommand{
		subcommands:       subcommands,
		subcommandOptions: subcommandOptions,
	}
}
