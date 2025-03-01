package internal

import (
	"io"
	"log"
	"neo-inu/internal/ygo"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type YgoCommand struct{}

type YgoCommandOptionType int

const (
	YgoCommandOptionTypeRandom YgoCommandOptionType = iota
	YgoCommandOptionTypeSearch
	YgoCommandOptionTypeBanlist
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
						Name:        "name",
						Description: "The exact name of the card.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "fname",
						Description: "A fuzzy search using a string.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "id",
						Description: "The 8-digit passcode of the card. You cannot pass this alongside name.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "konami_id",
						Description: "The Konami ID of the card. This is not the passcode.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "type",
						Description: "The type of card you want to filter by.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "atk",
						Description: "Filter by atk value.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "def",
						Description: "Filter by def value.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "level",
						Description: "Filter by card level/RANK.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "race",
						Description: "Filter by the card race which is officially called type. This is also used for Spell/Trap cards",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "attribute",
						Description: "Filter by the card attribute. You can pass multiple comma separated Attributes to this parameter.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "link",
						Description: "Filter the cards by Link value.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "linkmarker",
						Description: "Filter the cards by Link Marker value. You can pass multiple ',' separated values to this parameter.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "scale",
						Description: "Filter the cards by Pendulum Scale value.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "cardset",
						Description: "Filter the cards by card set.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "archetype",
						Description: "Filter the cards by archetype.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "banlist",
						Description: "Filter the cards by banlist.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "sort",
						Description: "Sort the order of the cards.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "format",
						Description: "Sort the format of the cards.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "misc",
						Description: "Pass yes to show additional response info",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "staple",
						Description: "Pass yes to check if card is a staple.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "has_effect",
						Description: "Check if a card actually has an effect or not by passing a boolean true/false.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "startdate",
						Description: "Filter based on cards' release date. Format dates as YYYY-mm-dd.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "enddate",
						Description: "Filter based on cards' release date. Format dates as YYYY-mm-dd.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "dateregion",
						Description: "Pass dateregion as tcg (default) or ocg.",
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
	defer func() {
		for _, file := range resp.Data.Files {
			if f, ok := file.Reader.(io.Closer); ok {
				f.Close()
			}
		}
	}()

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds:  resp.Data.Embeds,
		Content: resp.Data.Content,
		Files:   resp.Data.Files,
	})

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
			case discordgo.ApplicationCommandOptionInteger:
				params.SearchOption[opt.Name] = strconv.FormatInt(opt.IntValue(), 10)
			case discordgo.ApplicationCommandOptionBoolean:
				params.SearchOption[opt.Name] = strconv.FormatBool(opt.BoolValue())
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
