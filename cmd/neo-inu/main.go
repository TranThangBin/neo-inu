package main

import (
	"flag"
	"log"
	"neo-inu/internal"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	token := flag.String("token", os.Getenv("TOKEN"), "Your discord bot token look for TOKEN variable if not provide")
	rmcmd := flag.Bool("rmcmd", true, "Remove all command after shutdown default: true")
	guildId := flag.String("guild", "", "Test guild ID default: \"\" (mean global)")

	flag.Parse()

	neoinu := internal.NewNeoInu(*token, *rmcmd, *guildId)

	if err := neoinu.Init(); err != nil {
		log.Fatalln("Invalid bot token " + *token + ": " + err.Error())
	}

	if err := neoinu.Open(); err != nil {
		log.Fatalln("Something went wrong when starting the bot: ", err.Error())
	}

	if err := neoinu.AddSlashCommand(internal.NewPingCommand()); err != nil {
		log.Println("Something went wrong when adding ping command")
	}

	defer neoinu.Close()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, os.Kill)
	log.Println("Press Ctrl+C to exit")

	<-stop

	log.Println("Neo inu... Peace out!")
}
