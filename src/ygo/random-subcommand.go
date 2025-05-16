package ygo

import (
	"github.com/bwmarrin/discordgo"
)

type YgoRandomSubcommand struct{ ygoprodeckClient ygoprodeckClient }

func (yg YgoRandomSubcommand) NewApplicationCommandOption() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "random",
		Description: "Give you a random Yu-gi-oh card",
	}
}

func (yg YgoRandomSubcommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return err
	}

	resp, err := yg.ygoprodeckClient.SearchRandomCard()

	if err != nil {
		_, err = s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Content: "Something went wrong when retriving ygo card",
		})
		return err
	}

	_, err = s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Type: discordgo.EmbedTypeImage,
				Image: &discordgo.MessageEmbedImage{
					URL: resp.Data[0].CardImages[0].ImageUrl,
				},
			},
		},
	})
	return err
}
