package ygo

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type YgoSearchSubcommand struct{ ygoprodeckClient ygoprodeckClient }

func (yg YgoSearchSubcommand) NewApplicationCommandOption() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
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
	}
}

func (yg YgoSearchSubcommand) parseInput(data discordgo.ApplicationCommandInteractionData) map[string]string {
	searchOption := make(map[string]string)
	for _, opt := range data.Options {
		switch opt.Type {
		case discordgo.ApplicationCommandOptionString:
			searchOption[opt.Name] = opt.StringValue()
		case discordgo.ApplicationCommandOptionInteger:
			searchOption[opt.Name] = strconv.FormatInt(opt.IntValue(), 10)
		case discordgo.ApplicationCommandOptionBoolean:
			searchOption[opt.Name] = strconv.FormatBool(opt.BoolValue())
		}
	}
	return searchOption
}

func (yg YgoSearchSubcommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return err
	}

	queries := yg.parseInput(i.ApplicationCommandData())

	resp, err := yg.ygoprodeckClient.SearchCard(queries)

	if err != nil {
		_, err = s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Content: "Something went wrong when retriving ygo card",
		})
		return err
	}

	if len(resp.Data) <= 0 {
		notfoundImage, err := os.Open(filepath.Join("assets", "notfound.jpg"))
		if err != nil {
			return nil
		}
		defer notfoundImage.Close()
		_, err = s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Files: []*discordgo.File{
				{
					Name:        "notfound.jpg",
					ContentType: "image/jpg",
					Reader:      notfoundImage,
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
