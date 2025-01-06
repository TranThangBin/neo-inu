package internal

import (
	"log"
	"neo-inu/internal/ygo"
	"os"

	"github.com/bwmarrin/discordgo"
)

type YgoCommand struct{}

type YgoCommandOptionType int

const (
	YgoCommandOptionTypeRandom YgoCommandOptionType = iota
	YgoCommandOptionTypeSearch
)

type YgoCommandParams struct {
	Option       YgoCommandOptionType
	SearchOption map[string]string
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
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "search",
				Description: "Search a Yu-gi-oh card",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "fname",
						Description: "Fuzzy search a card with name",
						Required:    false,
					},
				},
			},
		},
	}
}

func (yg *YgoCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return err
	}

	resp := yg.NewResponse(yg.GetParams(i.ApplicationCommandData()))
	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds:  resp.Data.Embeds,
		Content: resp.Data.Content,
		Files:   resp.Data.Files,
	})

	for _, file := range resp.Data.Files {
		file.Reader.(*os.File).Close()
	}

	return err
}

func (yg *YgoCommand) NewResponse(params interface{}) *discordgo.InteractionResponse {
	cmdParams := params.(YgoCommandParams)

	switch cmdParams.Option {
	case YgoCommandOptionTypeRandom:
		return yg.NewRandomCardResponse()

	case YgoCommandOptionTypeSearch:
		return yg.NewSearchCardResponse(cmdParams.SearchOption)
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Something went wrong when fetching from ygoprodeck api",
		},
	}
}

func (yg *YgoCommand) GetParams(data discordgo.ApplicationCommandInteractionData) YgoCommandParams {
	params := YgoCommandParams{}
	subCmd := data.Options[0]
	switch subCmd.Name {
	case "random":
		params.Option = YgoCommandOptionTypeRandom

	case "search":
		params.Option = YgoCommandOptionTypeSearch
		params.SearchOption = make(map[string]string)
		for _, opt := range subCmd.Options {
			switch opt.Type {
			case discordgo.ApplicationCommandOptionString:
				params.SearchOption[opt.Name] = opt.StringValue()
			}
		}
	}
	return params
}

func (yg YgoCommand) NewRandomCardResponse() *discordgo.InteractionResponse {
	resp, err := ygo.SearchRandomCard()

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

func (yg *YgoCommand) NewSearchCardResponse(opt map[string]string) *discordgo.InteractionResponse {
	resp, err := ygo.SearchCard(opt)

	if err != nil {
		log.Printf("An error happened when searching ygo card: {%v}", err)
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something went wrong when retriving ygo card",
			},
		}
	}

	if len(resp.Data) < 1 {
		file, err := os.Open("assets/notfound.jpg")
		if err != nil {
			log.Printf("Something went wrong when opening the notfound image")
			return nil
		}
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Files: []*discordgo.File{
					{
						Name:        "notfound.jpg",
						ContentType: "image/jpg",
						Reader:      file,
					},
				},
				Embeds: []*discordgo.MessageEmbed{
					{
						Type: discordgo.EmbedTypeImage,
						Image: &discordgo.MessageEmbedImage{
							URL: "attachment://notfound.jpg",
						},
					},
				},
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
