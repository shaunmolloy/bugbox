package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var ConfigPath = filepath.Join(os.Getenv("HOME"), ".config", "bugbox", "config.json")

type Config struct {
	GitHubToken string   `json:"github_token"`
	Orgs        []string `json:"orgs"`
}

func IsExist() bool {
	_, err := os.Stat(ConfigPath)
	return !os.IsNotExist(err)
}

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

func LoadConfig() (Config, error) {
	var conf Config
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return conf, nil // return empty config
		}
		return conf, err
	}

	if err := json.Unmarshal(data, &conf); err != nil {
		return conf, err
	}
	return conf, nil
}

func SaveConfig(config Config) error {
	dir := filepath.Dir(ConfigPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(ConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}
