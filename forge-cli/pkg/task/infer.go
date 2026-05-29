package task

import "strings"

// InferType infers the task type from the task ID using registry pattern matching
// with prefix/suffix fallback for runtime tasks.
// Returns empty string for unknown IDs (no fallback).
//
// Resolution order:
// 1. Stage-gate/summary suffixes (*.gate, *.summary)
// 2. Registry pattern matching (exact, per-surface-type, per-surface-key)
// 3. Runtime task prefixes (doc-fix-*, fix-*, disc-*)
//
// The surfaces map is used to validate per-surface-key ID suffixes.
// A nil surfaces map disables surface-key prefix matching entirely.
func InferType(id string, surfaces map[string]string) string {
	// Phase 1: Stage-gate and summary suffixes (not in registry, always prefix-free)
	if strings.HasSuffix(id, IDSuffixSummary) {
		return TypeDocSummary
	}
	if strings.HasSuffix(id, IDSuffixGate) {
		return TypeGate
	}

	// Phase 2: Registry iteration — match against all registry ID patterns
	if typ := matchRegistryID(id, surfaces); typ != "" {
		return typ
	}

	// Phase 3: Runtime task prefix fallback (not in registry)
	if strings.HasPrefix(id, "doc-fix-") {
		return TypeDocFix
	}
	if strings.HasPrefix(id, "fix-") || strings.HasPrefix(id, "disc-") {
		return TypeCodingFix
	}

	return ""
}

// typeSuffixedID checks if id matches the pattern "base" + "-" + type.
// e.g., typeSuffixedID("T-test-gen-scripts-api", "T-test-gen-scripts") → true
// e.g., typeSuffixedID("T-test-gen-scripts", "T-test-gen-scripts") → false (exact match handled separately)
func typeSuffixedID(id, base string) bool {
	if !strings.HasPrefix(id, base) {
		return false
	}
	rem := id[len(base):]
	if len(rem) == 0 || rem[0] != '-' {
		return false
	}
	rem = rem[1:]
	if len(rem) == 0 {
		return false
	}
	for _, c := range rem {
		if (c < 'a' || c > 'z') && c != '-' {
			return false
		}
	}
	return true
}

// ExtractTypeSuffix extracts the type suffix from a type-suffixed task ID.
// Returns empty string if no type suffix is present.
// e.g., ExtractTypeSuffix("T-test-gen-scripts-api", "T-test-gen-scripts") → "api"
// e.g., ExtractTypeSuffix("T-test-gen-scripts", "T-test-gen-scripts") → ""
func ExtractTypeSuffix(id, base string) string {
	if !typeSuffixedID(id, base) {
		return ""
	}
	rem := id[len(base):]
	return rem[1:]
}
