package internal

import (
	"log"
	"neo-inu/pkg"

	"github.com/bwmarrin/discordgo"
)

type NeoInu struct {
	session            *discordgo.Session
	rmcmd              bool
	guildId            string
	commands           []pkg.Command
	commandHandlers    map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) error
	registeredCommands []*discordgo.ApplicationCommand
}

func (n *NeoInu) Open() error {
	n.session.AddHandler(n.onReady)
	n.session.AddHandler(n.onCommand)
	return n.session.Open()
}

func (n *NeoInu) Close() error {
	if n.rmcmd {
		for _, cmd := range n.registeredCommands {
			err := n.session.ApplicationCommandDelete(n.session.State.User.ID, n.guildId, cmd.ID)
			if err != nil {
				log.Printf("cannot delete command %s because of {%v}\n", cmd.Name, err)
				continue
			}
			log.Println("successfully deleted command: ", cmd.Name)
		}
	}
	return n.session.Close()
}

func (n *NeoInu) onReady(s *discordgo.Session, _ *discordgo.Ready) {
	log.Printf("Neo inu ready: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	log.Println("Initializing commands...")

	for i, cmd := range n.commands {
		_cmd := cmd.NewApplicationCommand()
		c, err := s.ApplicationCommandCreate(n.session.State.User.ID, n.guildId, _cmd)
		if err != nil {
			log.Printf("cannot add command %s because of {%v}\n", _cmd.Name, err)
			continue
		}
		log.Println("successfully added command: ", c.Name)
		n.registeredCommands[i] = c
		n.commandHandlers[c.Name] = cmd.Execute
	}

	log.Println("finish initializing command!")
}

func (n *NeoInu) onCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := n.commandHandlers[i.ApplicationCommandData().Name]; ok {
		err := h(s, i)
		if err != nil {
			log.Println(
				err, "an error happened when executing command",
				i.ApplicationCommandData().Name,
			)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Error",
							Description: "Something went wrong when handling the command",
							Color:       16711680,
						},
					},
				},
			})
		}
	}
}

func NewNeoInu(token string, rmcmd bool, guildId string) *NeoInu {
	bot := &NeoInu{}

	session, err := discordgo.New("Bot " + token)

	if err != nil {
		log.Fatalln(err, "something went wrong when initializing Neo Inu")
	}

	bot.session = session
	bot.rmcmd = rmcmd
	bot.guildId = guildId

	bot.commands = []pkg.Command{
		NewPingCommand(),
		NewYgoCommand(),
		NewJoinVoiceCommand(),
	}
	bot.registeredCommands = make([]*discordgo.ApplicationCommand, len(bot.commands))
	bot.commandHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) error, len(bot.commands))

	return bot
}
