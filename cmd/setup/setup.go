package setup

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shaunmolloy/bugbox/internal/logging"
	"github.com/shaunmolloy/bugbox/internal/storage/config"
)

func Setup() error {
	if !isSetup() {
		return nil
	}

	clearTerminal()

	logging.Info("Starting setup")

	handleGitHubToken()
	handleGitHubOrgs()

	if err := config.Validate(); err != nil {
		logging.Error(fmt.Sprintf("Config format is invalid: %v\n", err))
		return err
	}

	logging.Info("Setup complete!")
	return nil
}

func isSetup() bool {
	if len(os.Args) > 1 && os.Args[1] == "setup" {
		return true
	}
	if !config.IsExist(config.ConfigPath) {
		return true
	}
	if err := config.Validate(); err != nil {
		return true
	}
	return false
}

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func handleGitHubToken() error {
	fmt.Print("\nEnter your GitHub personal access token: ")

	reader := bufio.NewReader(os.Stdin)
	value, err := reader.ReadString('\n')
	token := strings.TrimSpace(value)
	if err != nil {
		logging.Error(fmt.Sprintf("Error reading token: %v\n", err))
		return err
	}

	conf, err := config.LoadConfig()
	if err != nil {
		conf = config.Config{} // fallback to default
	}

	conf.GitHubToken = token
	if err := config.SaveConfig(conf); err != nil {
		logging.Error(fmt.Sprintf("Error saving config: %v\n", err))
		return err
	}
	return nil
}

func handleGitHubOrgs() error {
	fmt.Print("\nEnter GitHub org(s) to find issues from: ")

	reader := bufio.NewReader(os.Stdin)
	value, err := reader.ReadString('\n')
	orgs := strings.Split(strings.TrimSpace(value), ",")
	if err != nil {
		logging.Error(fmt.Sprintf("Error reading orgs: %v\n", err))
		return err
	}

	conf, err := config.LoadConfig()
	if err != nil {
		conf = config.Config{} // fallback to default
	}

	conf.Orgs = orgs
	if err := config.SaveConfig(conf); err != nil {
		logging.Error(fmt.Sprintf("Error saving config: %v\n", err))
		return err
	}
	return nil
}
