package tui

import (
	"sort"

	"github.com/shaunmolloy/bugbox/internal/types"
)

func indexOf(list []string, value string) int {
	for i, v := range list {
		if v == value {
			return i
		}
	}
	return -1 // Not found
}

func fallback(primary string, alt string) string {
	if primary != "" {
		return primary
	}
	return alt
}

// sortIssues sorts by created, org, and repo
func sortIssues(issues []types.Issue) {
	sort.Slice(issues, func(i, j int) bool {
		if !issues[i].CreatedAt.Equal(issues[j].CreatedAt) {
			return issues[i].CreatedAt.After(issues[j].CreatedAt)
		}
		if issues[i].Org != issues[j].Org {
			return issues[i].Org < issues[j].Org
		}
		return issues[i].Repo < issues[j].Repo
	})
}
