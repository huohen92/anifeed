package config

import "strings"

// DefaultAutoPriority is the order used by the "apply best rules" action.
// IDs match the rule library categories returned by /api/rules.
func DefaultAutoPriority() (language, subtitle, container, resolution []string) {
	return []string{"chs", "bilingual", "cht"},
		[]string{"softsub", "hardsub", "external"},
		[]string{"mkv", "mp4"},
		[]string{"r1080", "r2160", "r720"}
}

func normalizePriority(values, fallback []string, allowed map[string]bool) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] || !allowed[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	for _, value := range fallback {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}
