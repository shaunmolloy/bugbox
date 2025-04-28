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

// PruneInvalidOrgs removes orgs from issues config no longer in main config
func PruneInvalidOrgs() error {
	// Load the current issues configuration
	issuesConf, err := LoadIssues()
	if err != nil {
		return err
	}

	// Load the current main configuration
	mainConf, err := LoadConfig()
	if err != nil {
		return err
	}

	validOrgs := make(map[string]struct{}, len(mainConf.Orgs))
	for _, org := range mainConf.Orgs {
		validOrgs[org] = struct{}{}
	}

	// Iterate over the issues configuration and remove invalid orgs
	for org := range issuesConf {
		if _, exists := validOrgs[org]; !exists {
			delete(issuesConf, org)
		}
	}

	return SaveIssues(issuesConf)
}
