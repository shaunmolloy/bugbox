package config

import (
	"testing"

	"github.com/shaunmolloy/bugbox/internal/types"
)

func TestIsExist(t *testing.T) {
	t.Run("returns false for invalid file", func(t *testing.T) {
		if _, err := IsExist("./README.md"); err == nil {
			t.Fatal("expected nil, got error")
		}
	})

	t.Run("returns true for valid file", func(t *testing.T) {
		exists, err := IsExist("../../../README.md")
		if err != nil {
			t.Fatal("expected nil, got error")
		}
		if exists != true {
			t.Fatal("expected true, got false")
		}
	})
}

func TestFlattenIssues(t *testing.T) {
	t.Run("returns flat slice for non-empty input", func(t *testing.T) {
		issues := Issues{
			"shaunmolloy": {
				"repo1": {
					1: {ID: 1, Title: "Issue 1"},
				},
			},
		}
		expected := []types.Issue{
			{ID: 1, Title: "Issue 1"},
		}
		flat := FlattenIssues(issues)
		if len(flat) != len(expected) {
			t.Fatalf("expected %d, got %d", len(expected), len(flat))
		}
	})
}
