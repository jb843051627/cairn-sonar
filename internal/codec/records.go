package codec

import "strings"

func RecordFields(values []string) string {
	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			cleaned = append(cleaned, value)
		}
	}
	return strings.Join(cleaned, "|")
}

func RecordCount(values []string) int { return len(RecordFields(values)) }
