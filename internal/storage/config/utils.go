package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/shaunmolloy/bugbox/internal/types"
)

// IsExist returns true if path exists
func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	return !os.IsNotExist(err), err
}

// LoadFromFile loads config from a file
func LoadFromFile(path string, conf any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(conf)
}

// SaveToFile saves config to a file
func SaveToFile(path string, data any) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// FlattenIssues converts the hierarchical issue structure to a flat slice
func FlattenIssues(issues Issues) []types.Issue {
	var flat []types.Issue

	for _, repoMap := range issues {
		for _, issueMap := range repoMap {
			for _, issue := range issueMap {
				flat = append(flat, issue)
			}
		}
	}

	return flat
}
