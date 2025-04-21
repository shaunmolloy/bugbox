package config

import "github.com/shaunmolloy/bugbox/internal/types"

type Config struct {
	GitHubToken string   `json:"github_token"`
	Orgs        []string `json:"orgs"`
}

type Issues []types.Issue
