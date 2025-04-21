package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shaunmolloy/bugbox/cmd/setup"
	"github.com/shaunmolloy/bugbox/internal/logging"
	"github.com/shaunmolloy/bugbox/internal/scheduler"
	"github.com/shaunmolloy/bugbox/internal/tui"
)

func main() {
	if err := logging.SetupLogger(); err != nil {
		log.Printf("Error setting up logger: %v\n", err)
		os.Exit(1)
	}
	logging.Info("Bugbox started")

	if err := setup.Setup(); err != nil {
		logging.Error(fmt.Sprintf("Setup failed: %v\n", err))
		os.Exit(1)
	}

	scheduler.FetchIssues()
	tui.Start()
}

