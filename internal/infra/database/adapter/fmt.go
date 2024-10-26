package adapter

import (
	"strings"
)

func ExtractColumnFromContraint(constraint string) string {
	parts := strings.Split(constraint, "_")
	switch len(parts) {
	case 3:
		return parts[1]
	default:
		return ""
	}
}
