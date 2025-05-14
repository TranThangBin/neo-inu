package main

import (
	"log"
	"neo-inu/src"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	var neoinu src.App = src.NewNeoInu(
		os.Getenv("TOKEN"),
		os.Getenv("RMCMD") != "false",
		os.Getenv("GUILD"),
		logger,
	)
	if err := neoinu.Open(); err != nil {
		logger.Fatalln(err, "something went wrong when opening connection")
	}
	defer func() {
		if err := neoinu.Close(); err != nil {
			logger.Fatalln(err, "something went wrong when closing connection")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	logger.Println("Press Ctrl+C to exit")
	<-stop

	logger.Println("Neo inu... Peace out!")
}
