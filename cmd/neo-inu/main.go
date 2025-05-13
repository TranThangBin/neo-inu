package main

import (
	"log"
	"neo-inu/internal"
	"neo-inu/pkg"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var neoinu pkg.App = internal.NewNeoInu(
		os.Getenv("TOKEN"),
		os.Getenv("RMCMD") != "false",
		os.Getenv("GUILD"),
	)
	if err := neoinu.Open(); err != nil {
		log.Fatalln(err, "something went wrong when opening connection")
	}
	defer func() {
		if err := neoinu.Close(); err != nil {
			log.Fatalln(err, "something went wrong when closing connection")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Neo inu... Peace out!")
}
