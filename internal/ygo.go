package internal

import (
	"log"
	"neo-inu/internal/ygo"

	"github.com/bwmarrin/discordgo"
)

type YgoCommand struct{}

type ygoCommandParams struct {
	Random bool
}

func (yg *YgoCommand) NewApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ygo",
		Description: "Use ygoprodeck api to do cool stuff",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "random",
				Description: "Give you a random Yu-gi-oh card",
			},
		},
	}
}

func (yg *YgoCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	params := ygoCommandParams{
		Random: false,
	}
	for _, opt := range i.ApplicationCommandData().Options {
		switch opt.Name {
		case "random":
			params.Random = true
		}
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		return err
	}
	resp := yg.NewResponse(params)
	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds:  resp.Data.Embeds,
		Flags:   resp.Data.Flags,
		Content: resp.Data.Content,
	})
	return err
}

func (yg *YgoCommand) NewResponse(params interface{}) *discordgo.InteractionResponse {
	cmdParams := params.(ygoCommandParams)

	if cmdParams.Random {
		resp, err := ygo.RandomCard()
		if err != nil {
			log.Printf("An error happened when getting random ygo card: {%v}", err)
			return &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Something went wrong when retriving ygo card",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}
		}
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{
					{
						Type: discordgo.EmbedTypeImage,
						Image: &discordgo.MessageEmbedImage{
							URL: resp.Data[0].CardImages[0].ImageUrlSmall,
						},
					},
				},
			},
		}
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Why are we here",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
}

func NewYgoCommand() *YgoCommand {
	return &YgoCommand{}
}
