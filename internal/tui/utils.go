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

func sortByCreated(issues []types.Issue) {
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].CreatedAt.After(issues[j].CreatedAt)
	})
}
