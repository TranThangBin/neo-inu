package internal

import (
	"log"
	"neo-inu/pkg"

	"github.com/bwmarrin/discordgo"
)

type NeoInu struct {
	token              string
	rmcmd              bool
	guildId            string
	session            *discordgo.Session
	commands           map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	registeredCommands []*discordgo.ApplicationCommand
}

func (n *NeoInu) Init() error {
	var err error
	n.session, err = discordgo.New("Bot " + n.token)
	n.registeredCommands = make([]*discordgo.ApplicationCommand, 0, 20)
	n.commands = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	return err
}

func (n *NeoInu) Open() error {
	n.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Neo inu ready: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})
	n.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := n.commands[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	return n.session.Open()
}

func (n *NeoInu) Close() error {
	if n.rmcmd {
		for _, cmd := range n.registeredCommands {
			n.session.ApplicationCommandDelete(n.session.State.User.ID, n.guildId, cmd.ApplicationID)
		}
	}
	return n.session.Close()
}

func (n *NeoInu) AddSlashCommand(cmds ...pkg.Command) error {
	for _, cmd := range cmds {
		appCmd := cmd.ApplicationCommand()
		c, err := n.session.ApplicationCommandCreate(n.session.State.User.ID, n.guildId, appCmd)
		if err != nil {
			return err
		}
		if n.rmcmd {
			n.registeredCommands = append(n.registeredCommands, c)
		}
		n.commands[appCmd.Name] = cmd.Execute
	}
	return nil
}

func NewNeoInu(token string, rmcmd bool, guildId string) *NeoInu {
	return &NeoInu{
		token: token,
	}
}
