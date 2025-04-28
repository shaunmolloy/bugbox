package github

import (
	"strings"
)

// parseRepo extracts the repository name from html_url.
func parseRepo(url string) string {
	// https://github.com/{org}/{repo}/issues/{id}
	url = strings.Replace(url, "https://github.com", "", 1)
	parts := strings.Split(url, "/")

	if len(parts) >= 2 {
		return parts[2]
	}
	return ""
}
