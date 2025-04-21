package scheduler

import (
	"fmt"
	"time"

	"github.com/shaunmolloy/bugbox/internal/issues/github"
	"github.com/shaunmolloy/bugbox/internal/logging"
	"github.com/shaunmolloy/bugbox/internal/tui"
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
		return
	}
	logging.Info("Fetched GitHub issues successfully")

	// Signal the TUI to refresh with the updated issues
	// Use non-blocking send to prevent hanging if channel is full
	select {
	case tui.RefreshChan <- struct{}{}:
		logging.Info("Sent refresh signal to TUI")
	default:
		// Channel is full, no need to send another signal
	}
}
