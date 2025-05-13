package internal

import "github.com/bwmarrin/discordgo"

type VCCCommand struct{}

func (jv *VCCCommand) NewApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "vcc",
		Description: "Voice channel control control bot behavior to channel",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "leave",
				Description: "Tell the bot to leave the voice channel",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "join",
				Description: "Allow the bot to join a voice channel",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "channel",
						Description: "A channel you want the bot to join",
						Required:    true,
						ChannelTypes: []discordgo.ChannelType{
							discordgo.ChannelTypeGuildVoice,
						},
					},
				},
			},
		},
	}
}

func (vc *VCCCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	subCmd := i.ApplicationCommandData().Options[0]
	switch subCmd.Name {
	case "leave":
		v, ok := s.VoiceConnections[i.GuildID]
		if !ok {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "The bot is not in any voice channel in this guild",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
		if err := v.Disconnect(); err != nil {
			return err
		}
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Successfully leave the voice channel",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	case "join":
		targetedChannel := subCmd.Options[0].ChannelValue(s)
		_, err := s.ChannelVoiceJoin(targetedChannel.GuildID, targetedChannel.ID, false, true)
		if err != nil {
			return err
		}
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Successfully join the voice channel",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
	return nil
}

func NewJoinVoiceCommand() *VCCCommand {
	return &VCCCommand{}
}
