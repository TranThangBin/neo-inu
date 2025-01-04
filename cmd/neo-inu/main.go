package main

import (
	"flag"
	"log"
	"neo-inu/internal"
	"neo-inu/pkg"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	token := flag.String("token", "", "Your discord bot token. $TOKEN is prioritized.")
	rmcmd := flag.Bool("rmcmd", true, "Remove all command after shutdown. $RMCMD is prioritized")
	guildId := flag.String("guild", "", "Test guild ID. $GUILD is prioritized")
	flag.Parse()

	godotenv.Load()
	if os.Getenv("TOKEN") != "" {
		*token = os.Getenv("TOKEN")
	}
	if os.Getenv("RMCMD") != "" {
		*rmcmd = os.Getenv("RMCMD") != "false"
	}
	if os.Getenv("GUILD") != "" {
		*guildId = os.Getenv("GUILD")
	}

	var neoinu pkg.App = internal.NewNeoInu(*token, *rmcmd, *guildId,
		internal.NewPingCommand(),
		internal.NewYgoCommand())
	neoinu.Init()
	if err := neoinu.Open(); err != nil {
		log.Fatalln("Something went wrong when starting the bot: ", err.Error())
	}
	defer neoinu.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Neo inu... Peace out!")
}
