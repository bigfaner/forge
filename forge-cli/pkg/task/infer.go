package task

import "strings"

// InferType infers the task type from the task ID using pattern matching.
// Returns empty string for unknown IDs (no fallback).
//
// Handles profile-suffixed IDs from the test profile system (D9):
// single profile: T-test-2, T-test-3, T-test-4 (no suffix)
// multi profile:  T-test-2a, T-test-2b, T-test-3a, ... (letter suffix)
//
// Handles type-suffixed IDs from per-type test pipeline split:
// no profile:   T-test-2-api, T-test-2-tui, T-test-2-cli (type suffix after number)
// with profile: T-test-2a-api, T-test-2b-tui (profile letter + type suffix)
func InferType(id string) string {
	switch {
	case strings.HasSuffix(id, ".summary"):
		return TypeDocSummary
	case strings.HasSuffix(id, ".gate"):
		return TypeGate
	case id == "T-test-1b":
		return TypeTestEvalCases
	case id == "T-test-1", profileSuffixedID(id, "T-test-1"):
		return TypeTestGenCases
	case id == "T-test-2", profileSuffixedID(id, "T-test-2"), typeSuffixedID(id, "T-test-2"):
		return TypeTestGenScripts
	case id == "T-test-3", profileSuffixedID(id, "T-test-3"), typeSuffixedID(id, "T-test-3"):
		return TypeTestRun
	case id == "T-test-4", profileSuffixedID(id, "T-test-4"), typeSuffixedID(id, "T-test-4"):
		return TypeTestGraduate
	case id == "T-test-4.5":
		return TypeTestVerifyRegression
	case id == "T-test-5", id == "T-specs-1":
		return TypeDocConsolidate
	case profileSuffixedID(id, "T-quick-1"):
		return TypeTestGenCases
	case profileSuffixedID(id, "T-quick-2"), typeSuffixedID(id, "T-quick-2"):
		return TypeTestGenAndRun
	case profileSuffixedID(id, "T-quick-3"), typeSuffixedID(id, "T-quick-3"):
		return TypeTestGraduate
	case profileSuffixedID(id, "T-quick-4"):
		return TypeTestVerifyRegression
	case id == "T-quick-5", id == "T-quick-specs-1", profileSuffixedID(id, "T-quick-5"), profileSuffixedID(id, "T-quick-specs-1"):
		return TypeDocDrift
	case id == "T-clean-code-1":
		return TypeCleanCode
	case strings.HasPrefix(id, "fix-") || strings.HasPrefix(id, "disc-"):
		return TypeCodingFix
	case id == "T-eval-doc":
		return TypeDocEval
	default:
		return ""
	}
}

// profileSuffixedID checks if id matches the pattern "base" + single lowercase letter.
// e.g., profileSuffixedID("T-test-2a", "T-test-2") → true
// e.g., profileSuffixedID("T-test-2", "T-test-2") → false (exact match handled separately)
func profileSuffixedID(id, base string) bool {
	if !strings.HasPrefix(id, base) {
		return false
	}
	suffix := id[len(base):]
	if len(suffix) != 1 {
		return false
	}
	return suffix[0] >= 'a' && suffix[0] <= 'z'
}

// typeSuffixedID checks if id matches the pattern "base" + optional profile letter + "-" + capability.
// e.g., typeSuffixedID("T-test-2-api", "T-test-2") → true
// e.g., typeSuffixedID("T-test-2a-api", "T-test-2") → true
// e.g., typeSuffixedID("T-test-2-web-ui", "T-test-2") → true
// e.g., typeSuffixedID("T-test-2a", "T-test-2") → false (profile suffix only, no type)
// e.g., typeSuffixedID("T-test-2", "T-test-2") → false (exact match handled separately)
func typeSuffixedID(id, base string) bool {
	if !strings.HasPrefix(id, base) {
		return false
	}
	rem := id[len(base):]
	if len(rem) == 0 {
		return false
	}

	// Skip optional single profile letter.
	if rem[0] >= 'a' && rem[0] <= 'z' {
		rem = rem[1:]
	}

	// Must start with hyphen followed by at least one lowercase letter or digit.
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

// ExtractTypeSuffix extracts the type/capability suffix from a type-suffixed task ID.
// Returns empty string if no type suffix is present.
// e.g., ExtractTypeSuffix("T-test-2-api", "T-test-2") → "api"
// e.g., ExtractTypeSuffix("T-test-2a-tui", "T-test-2") → "tui"
// e.g., ExtractTypeSuffix("T-test-2", "T-test-2") → ""
func ExtractTypeSuffix(id, base string) string {
	if !typeSuffixedID(id, base) {
		return ""
	}
	rem := id[len(base):]
	if rem[0] >= 'a' && rem[0] <= 'z' {
		rem = rem[1:]
	}
	// rem starts with "-", strip it
	return rem[1:]
}
