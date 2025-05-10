package github

import "testing"

func TestParseRepo(t *testing.T) {
	t.Run("returns empty string for invalid URL", func(t *testing.T) {
		want := ""
		got := parseRepo("invalid-url")

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns repo name from URL", func(t *testing.T) {
		want := "repo"
		got := parseRepo("https://github.com/org/repo/issues/1")

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
