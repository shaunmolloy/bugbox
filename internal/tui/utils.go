package tui

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
