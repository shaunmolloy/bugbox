package config

import (
	"os"
	"path/filepath"
)

var IssuesPath = filepath.Join(os.Getenv("HOME"), ".config", "bugbox", "issues.json")

// SaveIssues saves issues to a file
func SaveIssues(cfg Issues) error {
	return SaveToFile(IssuesPath, cfg)
}

// LoadIssues loads issues from a file
func LoadIssues() (Issues, error) {
	var cfg Issues
	err := LoadFromFile(IssuesPath, &cfg)
	return cfg, err
}
