package homie

import "strings"

// IsValidID checks the ID can be used in a topic.
func IsValidID(id string) bool {
	if id == "" || strings.HasPrefix(id, "-") || strings.HasSuffix(id, "-") {
		return false
	}
	for _, char := range id {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < '0' || char > '9') && char != '-' {
			return false
		}
	}
	return true
}
