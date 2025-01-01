package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	token   *string
	rmcmd   *bool
	guildId *string
)

func main() {
	godotenv.Load()

	token = flag.String("token", os.Getenv("TOKEN"), "Your discord bot token look for TOKEN variable if not provide")
	rmcmd = flag.Bool("rmcmd", true, "Remove all command after shutdown default: true")
	guildId = flag.String("guild", "", "Test guild ID default: \"\" (mean global)")

	flag.Parse()

	s, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatalln("Invalid bot token " + *token + ": " + err.Error())
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Neo inu ready: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})

	if err = s.Open(); err != nil {
		log.Fatalln("Something went wrong when starting the bot: ", err.Error())
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Neo inu... Peace out!")
}
