package scheduler

import (
	"fmt"
	"time"

	"github.com/shaunmolloy/bugbox/internal/issues/github"
	"github.com/shaunmolloy/bugbox/internal/logging"
)

// FetchIssues starts a goroutine that polls for GitHub issues every minute
func FetchIssues() {
	// Fetch when app starts
	handleGitHub()

	// Fetch every minute
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			handleGitHub()
		}
	}()
}

func handleGitHub() {
	logging.Info("Fetching GitHub issues...")
	err := github.FetchAllIssues()
	if err != nil {
		logging.Error(fmt.Sprintf("Fetching error: %v", err))
	}
	logging.Info("Fetched GitHub issues successfully")
}
