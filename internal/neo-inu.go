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
	commands           []pkg.Command
	commandHandlers    map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) error
	registeredCommands []*discordgo.ApplicationCommand
}

func (n *NeoInu) Init() {
	var err error
	n.session, err = discordgo.New("Bot " + n.token)

	if err != nil {
		log.Fatalln("Something went wrong when initializing Neo Inu: ", err.Error())
	}

	n.registeredCommands = make([]*discordgo.ApplicationCommand, 0, 20)
	n.commandHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) error)

}

func (n *NeoInu) Open() error {
	n.session.AddHandler(n.onReady)
	n.session.AddHandler(n.onCommand)
	return n.session.Open()
}

func (n *NeoInu) Close() error {
	if n.rmcmd {
		for _, cmd := range n.registeredCommands {
			if err := n.session.ApplicationCommandDelete(
				n.session.State.User.ID, n.guildId, cmd.ID,
			); err != nil {
				log.Printf("Cannot delete command %s because of {%v}\n", cmd.Name, err)
			} else {
				log.Println("Successfully deleted command: ", cmd.Name)
			}
		}
	}
	return n.session.Close()
}

func (n *NeoInu) onReady(s *discordgo.Session, _ *discordgo.Ready) {
	log.Printf("Neo inu ready: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	log.Println("Initializing commands...")
	for _, cmd := range n.commands {
		c, err := n.addSlashCommand(s, cmd)
		if err != nil {
			log.Printf("Cannot add command %s because of {%v}\n", cmd.NewApplicationCommand().Name, err)
		} else {
			log.Println("Successfully added command: ", c.Name)
		}
	}
	log.Println("Finish initializing command!")
}

func (n *NeoInu) onCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := n.commandHandlers[i.ApplicationCommandData().Name]; ok {
		err := h(s, i)
		if err != nil {
			log.Printf("An error happened when executing command %s: {%v}",
				i.ApplicationCommandData().Name, err)
		}
	}
}

func (n *NeoInu) addSlashCommand(s *discordgo.Session, cmd pkg.Command) (
	*discordgo.ApplicationCommand, error,
) {
	c, err := s.ApplicationCommandCreate(n.session.State.User.ID, n.guildId, cmd.NewApplicationCommand())

	if err != nil {
		return nil, err
	}
	if n.rmcmd {
		n.registeredCommands = append(n.registeredCommands, c)
	}

	n.commandHandlers[c.Name] = cmd.Execute
	return c, nil
}

func NewNeoInu(token string, rmcmd bool, guildId string, commands []pkg.Command) *NeoInu {
	return &NeoInu{
		token:    token,
		rmcmd:    rmcmd,
		guildId:  guildId,
		commands: commands,
	}
}
