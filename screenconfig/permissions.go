package screenconfig

import "strings"

// ExtractResourceKeys parses unique resource keys from a slice of permissions in "resource:action" format.
func ExtractResourceKeys(permissions []string) []string {
	keySet := make(map[string]bool)
	for _, perm := range permissions {
		parts := strings.SplitN(perm, ":", 2)
		if len(parts) >= 2 {
			keySet[parts[0]] = true
		}
	}

	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	return keys
}

// HasPermission checks if a specific permission exists in a permissions slice.
func HasPermission(permissions []string, required string) bool {
	for _, p := range permissions {
		if p == required {
			return true
		}
	}
	return false
}
