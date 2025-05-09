package utils

import (
	"testing"
	"time"
)

func TestRelativeTime(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "just now"},
		{10, "10 minutes ago"},
		{60, "1 hour ago"},
		{1440, "1 day ago"},
		{43200, "1 month ago"},
	}

	for _, test := range tests {
		now := time.Now().Add(-time.Duration(test.input) * time.Minute)
		result := RelativeTime(now)
		expected := test.expected
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	}
}

func TestPlural(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "s"},
		{1, ""},
		{2, "s"},
	}

	for _, test := range tests {
		result := plural(test.input)
		if result != test.expected {
			t.Errorf("expected %q, got %q", test.expected, result)
		}
	}
}
