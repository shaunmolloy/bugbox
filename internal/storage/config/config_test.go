package config

import (
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	t.Run("returns error for missing config", func(t *testing.T) {
		ConfigPath = "./config.json"
		if err := Validate(); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns nil for valid config", func(t *testing.T) {
		content := `{
			"github_token": "example",
			"orgs": ["example"]
		}`

		tmpFile, err := os.CreateTemp("", "config-*.json")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write to temp file: %v", err)
		}
		tmpFile.Close()

		ConfigPath = tmpFile.Name()
		if err := Validate(); err != nil {
			t.Fatal("expected nil, got error")
		}
	})
}
