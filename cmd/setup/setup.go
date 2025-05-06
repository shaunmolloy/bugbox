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

	conf, err := config.LoadConfig()
	if err != nil {
		conf = config.Config{} // fallback to default
	}

	handleGitHubToken(&conf)
	handleGitHubOrgs(&conf)

	logging.Info("Saving config...")
	if err := config.SaveConfig(conf); err != nil {
		logging.Error(fmt.Sprintf("Error saving config: %v\n", err))
		return err
	}

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
	if exists, _ := config.IsExist(config.ConfigPath); exists {
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

func handleGitHubToken(conf *config.Config) error {
	fmt.Printf("\nEnter your GitHub personal access token [%s]: ", conf.GitHubToken)
	input := parseInput()
	if input != "" {
		conf.GitHubToken = input
	}
	return nil
}

func handleGitHubOrgs(conf *config.Config) error {
	fmt.Printf("\nEnter GitHub org(s) (space-separated) [%s]: ", strings.Join(conf.Orgs, " "))
	input := strings.Split(parseInput(), " ")
	if len(input) > 0 && input[0] != "" {
		conf.Orgs = input
	}
	return nil
}

func parseInput() string {
	reader := bufio.NewReader(os.Stdin)
	value, err := reader.ReadString('\n')
	value = strings.TrimSpace(value)
	if err != nil {
		logging.Error(fmt.Sprintf("Error reading value: %v\n", err))
	}
	return value
}
