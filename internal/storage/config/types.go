package config

import "github.com/shaunmolloy/bugbox/internal/types"

type Config struct {
	GitHubToken string   `json:"github_token"`
	Orgs        []string `json:"orgs"`
}

// Issues as hierarchical structure of issues organized by org, repo, and id
type Issues map[string]map[string]map[int]types.Issue
