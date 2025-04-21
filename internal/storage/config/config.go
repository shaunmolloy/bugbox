package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var ConfigPath = filepath.Join(os.Getenv("HOME"), ".config", "bugbox", "config.json")

// Validate checks config.json exists with expected structure
func Validate() error {
	file, err := os.Open(ConfigPath)
	if err != nil {
		return fmt.Errorf("File not found")
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return fmt.Errorf("Invalid JSON format")
	}

	if config.GitHubToken == "" {
		return fmt.Errorf("Missing github_token")
	}

	if len(config.Orgs) == 0 {
		return fmt.Errorf("Missing orgs")
	}

	return nil
}

// SaveConfig saves the config to config.json
func SaveConfig(cfg Config) error {
	return SaveToFile(ConfigPath, cfg)
}

// LoadFromFile loads the config from config.json
func LoadConfig() (Config, error) {
	var cfg Config
	err := LoadFromFile(ConfigPath, &cfg)
	return cfg, err
}
