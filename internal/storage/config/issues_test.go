package config

import (
	"os"
	"testing"
)

func TestSaveIssues(t *testing.T) {
	t.Run("returns nil when saving issues config", func(t *testing.T) {
		tmpFile := createTmpFile(t, "{}")
		defer os.Remove(tmpFile.Name())

		IssuesPath = tmpFile.Name()
		issues := Issues{
			"example": {
				"issue": {
					1: {
						ID:    1,
						Title: "example",
					},
				},
			},
		}

		if err := SaveIssues(issues); err != nil {
			t.Fatal("expected nil, got error")
		}
	})
}

func TestLoadIssues(t *testing.T) {
	t.Run("returns error for missing issues config", func(t *testing.T) {
		IssuesPath = "./issues-load.json"
		if _, err := LoadIssues(); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns nil for existing issues config", func(t *testing.T) {
		tmpFile := createTmpFile(t, "{}")
		defer os.Remove(tmpFile.Name())

		IssuesPath = tmpFile.Name()
		if _, err := LoadIssues(); err != nil {
			t.Fatal("expected error, got nil")
		}
	})
}
