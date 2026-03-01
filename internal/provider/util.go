package provider

import (
	"strings"
)

func splitImportID(id string, expectedParts int) []string {
	parts := strings.SplitN(id, "/", expectedParts)
	if len(parts) != expectedParts {
		return nil
	}
	for _, p := range parts {
		if p == "" {
			return nil
		}
	}
	return parts
}
