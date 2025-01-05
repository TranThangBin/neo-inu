package internal

import (
	"log"
	"neo-inu/internal/ygo"

	"github.com/bwmarrin/discordgo"
)

type YgoCommand struct{}

type YgoCommandOptionType int

const (
	YgoCommandOptionTypeRandom YgoCommandOptionType = iota
)

type YgoCommandParams struct {
	Option YgoCommandOptionType
}

func (yg *YgoCommand) NewApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ygo",
		Description: "Use ygoprodeck api to do cool stuff",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "random",
				Required:    false,
				Description: "Give you a random Yu-gi-oh card",
			},
		},
	}
}

func (yg *YgoCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	params := YgoCommandParams{}

	if i.ApplicationCommandData().Options[0].Name == "random" {
		params.Option = YgoCommandOptionTypeRandom
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return err
	}

	resp := yg.NewResponse(params)
	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds:  resp.Data.Embeds,
		Content: resp.Data.Content,
	})

	return err
}

func (yg *YgoCommand) NewResponse(params interface{}) *discordgo.InteractionResponse {
	cmdParams := params.(YgoCommandParams)

	if cmdParams.Option == YgoCommandOptionTypeRandom {
		return yg.NewRandomCardResponse()
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Something went wrong when fetching from ygoprodeck api",
		},
	}
}

func (yg YgoCommand) NewRandomCardResponse() *discordgo.InteractionResponse {
	resp, err := ygo.RandomCard()

	if err != nil {
		log.Printf("An error happened when getting random ygo card: {%v}", err)
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something went wrong when retriving ygo card",
			},
		}
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Type: discordgo.EmbedTypeImage,
					Image: &discordgo.MessageEmbedImage{
						URL: resp.Data[0].CardImages[0].ImageUrl,
					},
				},
			},
		},
	}
}

func NewYgoCommand() *YgoCommand {
	return &YgoCommand{}
}
