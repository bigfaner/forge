package task

import "strings"

// Category constants for task type classification.
const (
	CategoryCoding     = "coding"
	CategoryDoc        = "doc"
	CategoryTest       = "test"
	CategoryValidation = "validation"
	CategoryGate       = "gate"
)

// CategoryForType maps a task type string to its category.
// Returns CategoryCoding for empty or unknown types.
func CategoryForType(typ string) string {
	switch {
	case typ == TypeGate:
		return CategoryGate
	case typ == TypeCleanCode:
		return CategoryCoding
	case strings.HasPrefix(typ, "coding."):
		return CategoryCoding
	case strings.HasPrefix(typ, "doc"):
		return CategoryDoc
	case strings.HasPrefix(typ, "test."):
		return CategoryTest
	case strings.HasPrefix(typ, "validation."):
		return CategoryValidation
	default:
		return CategoryCoding
	}
}
