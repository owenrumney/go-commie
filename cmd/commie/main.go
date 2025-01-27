package main

import (
	"flag"

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

	if err := gClient.Commit(); err != nil {
		log.Fatal(err)
	}
}
