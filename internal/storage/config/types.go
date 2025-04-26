package config

import "github.com/shaunmolloy/bugbox/internal/types"

type Config struct {
	GitHubToken string   `json:"github_token"`
	Orgs        []string `json:"orgs"`
}

// Issues represents the hierarchical structure of issues organized by org, repo, and issue number
type Issues map[string]map[string]map[int]types.Issue
