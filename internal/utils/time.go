package utils

import (
	"fmt"
	"time"
)

// RelativeTime returns a human-readable relative time string like "2 hours ago".
func RelativeTime(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%d minute%s ago", minutes, plural(minutes))
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		return fmt.Sprintf("%d hour%s ago", hours, plural(hours))
	case duration < 30*24*time.Hour:
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d day%s ago", days, plural(days))
	default:
		months := int(duration.Hours() / (24 * 30))
		return fmt.Sprintf("%d month%s ago", months, plural(months))
	}
}

func plural(n int) string {
	if n != 1 {
		return "s"
	}
	return ""
}
