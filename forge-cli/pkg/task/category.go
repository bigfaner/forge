package task

import (
	"strings"

	"forge-cli/pkg/forgelog"
)

// Category constants for task type classification.
const (
	CategoryCoding     = "coding"
	CategoryDoc        = "doc"
	CategoryTest       = "test"
	CategoryValidation = "validation"
	CategoryGate       = "gate"
	CategoryEval       = "eval"
)

// CategoryForType maps a task type string to its category.
// Returns CategoryCoding for empty or unknown types, with a log warning for unknowns.
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
	case strings.HasPrefix(typ, "eval."):
		return CategoryEval
	default:
		forgelog.Info("CategoryForType: unknown type %q, defaulting to coding\n", typ)
		return CategoryCoding
	}
}
