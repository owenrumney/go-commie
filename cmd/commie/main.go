package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/owenrumney/go-commie/internal/git"
	"github.com/owenrumney/go-commie/internal/logger"
)

func main() {
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()

	// Create a new logger with the debug option
	log := logger.New(logger.WithDebug(*debug))
	log.Debug("Starting commie")

	gClient, err := git.New(log)
	if err != nil {
		log.Fatal(err)
	}

	notifyChan := make(chan os.Signal, 1)

	signal.Notify(notifyChan, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-notifyChan
		log.Debug("Received signal to stop")
		os.Exit(0)
	}()

	if err := gClient.Commit(); err != nil {
		log.Fatal(err)
	}
}
