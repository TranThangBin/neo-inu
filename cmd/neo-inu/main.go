package main

import (
	"flag"
	"log"
	"neo-inu/internal"
	"neo-inu/pkg"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	token := flag.String("token", "", "your discord bot token. $TOKEN is prioritized.")
	rmcmd := flag.Bool("rmcmd", true, "remove all command after shutdown. $RMCMD is prioritized")
	guildId := flag.String("guild", "", "test guild ID. $GUILD is prioritized")
	flag.Parse()

	if os.Getenv("TOKEN") != "" {
		*token = os.Getenv("TOKEN")
	}
	if os.Getenv("RMCMD") != "" {
		*rmcmd = os.Getenv("RMCMD") != "false"
	}
	if os.Getenv("GUILD") != "" {
		*guildId = os.Getenv("GUILD")
	}

	var neoinu pkg.App = internal.NewNeoInu(*token, *rmcmd, *guildId, []pkg.Command{
		internal.NewPingCommand(),
		internal.NewYgoCommand(),
	})
	if err := neoinu.Open(); err != nil {
		log.Fatalln(err, "something went wrong when opening connection")
	}
	defer func() {
		if err := neoinu.Close(); err != nil {
			log.Println(err, "something went wrong when closing connection")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Neo inu... Peace out!")
}
