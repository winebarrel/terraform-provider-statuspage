package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func stringValueOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

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
