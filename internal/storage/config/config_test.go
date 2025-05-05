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

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		content := `{}`

		tmpFile := createTmpFile(t, content)
		defer os.Remove(tmpFile.Name())

		ConfigPath = tmpFile.Name()
		if err := Validate(); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns error for invalid github_token in JSON", func(t *testing.T) {
		content := `{
			"github_token": "",
			"orgs": ["example"]
		}`

		tmpFile := createTmpFile(t, content)
		defer os.Remove(tmpFile.Name())

		ConfigPath = tmpFile.Name()
		if err := Validate(); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns error for invalid orgs in JSON", func(t *testing.T) {
		content := `{
			"github_token": "example",
			"orgs": []
		}`

		tmpFile := createTmpFile(t, content)
		defer os.Remove(tmpFile.Name())

		ConfigPath = tmpFile.Name()
		if err := Validate(); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns nil for valid config", func(t *testing.T) {
		content := `{
			"github_token": "example",
			"orgs": ["example"]
		}`

		tmpFile := createTmpFile(t, content)
		defer os.Remove(tmpFile.Name())

		ConfigPath = tmpFile.Name()
		if err := Validate(); err != nil {
			t.Fatal("expected nil, got error")
		}
	})
}

func createTmpFile(t *testing.T, content string) *os.File {
	tmpFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	return tmpFile
}
