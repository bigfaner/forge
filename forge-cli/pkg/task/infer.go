package task

import "strings"

// InferType infers the task type from the task ID using pattern matching.
// Returns empty string for unknown IDs (no fallback).
//
// Handles profile-suffixed IDs from the test profile system:
// single profile: T-test-gen-scripts, T-test-run, T-test-graduate (no suffix)
// multi profile:  T-test-gen-scriptsa, T-test-gen-scriptsb, T-test-runa, ... (letter suffix)
//
// Handles type-suffixed IDs from per-type test pipeline split:
// no profile:   T-test-gen-scripts-api, T-test-gen-scripts-tui, T-test-gen-scripts-cli (type suffix)
// with profile: T-test-gen-scriptsa-api, T-test-gen-scriptsb-tui (profile letter + type suffix)
func InferType(id string) string {
	switch {
	case strings.HasSuffix(id, ".summary"):
		return TypeDocSummary
	case strings.HasSuffix(id, ".gate"):
		return TypeGate
	case id == "T-test-eval-cases":
		return TypeTestEvalCases
	case id == "T-test-gen-cases", profileSuffixedID(id, "T-test-gen-cases"):
		return TypeTestGenCases
	case id == "T-test-gen-scripts", profileSuffixedID(id, "T-test-gen-scripts"), typeSuffixedID(id, "T-test-gen-scripts"):
		return TypeTestGenScripts
	case id == "T-test-run", profileSuffixedID(id, "T-test-run"), typeSuffixedID(id, "T-test-run"):
		return TypeTestRun
	case id == "T-test-graduate", profileSuffixedID(id, "T-test-graduate"), typeSuffixedID(id, "T-test-graduate"):
		return TypeTestGraduate
	case id == "T-test-verify-regression":
		return TypeTestVerifyRegression
	case id == "T-specs-consolidate":
		return TypeDocConsolidate
	case id == "T-validate-code":
		return TypeValidationCode
	case id == "T-validate-ux":
		return TypeValidationUx
	case profileSuffixedID(id, "T-quick-gen-cases"):
		return TypeTestGenCases
	case profileSuffixedID(id, "T-quick-gen-and-run"), typeSuffixedID(id, "T-quick-gen-and-run"):
		return TypeTestGenAndRun
	case profileSuffixedID(id, "T-quick-graduate"), typeSuffixedID(id, "T-quick-graduate"):
		return TypeTestGraduate
	case profileSuffixedID(id, "T-quick-verify-regression"):
		return TypeTestVerifyRegression
	case id == "T-quick-doc-drift", profileSuffixedID(id, "T-quick-doc-drift"):
		return TypeDocDrift
	case id == "T-clean-code":
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
// e.g., profileSuffixedID("T-test-gen-scriptsa", "T-test-gen-scripts") → true
// e.g., profileSuffixedID("T-test-gen-scripts", "T-test-gen-scripts") → false (exact match handled separately)
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
// e.g., typeSuffixedID("T-test-gen-scripts-api", "T-test-gen-scripts") → true
// e.g., typeSuffixedID("T-test-gen-scriptsa-api", "T-test-gen-scripts") → true
// e.g., typeSuffixedID("T-test-gen-scripts-web-ui", "T-test-gen-scripts") → true
// e.g., typeSuffixedID("T-test-gen-scriptsa", "T-test-gen-scripts") → false (profile suffix only, no type)
// e.g., typeSuffixedID("T-test-gen-scripts", "T-test-gen-scripts") → false (exact match handled separately)
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
// e.g., ExtractTypeSuffix("T-test-gen-scripts-api", "T-test-gen-scripts") → "api"
// e.g., ExtractTypeSuffix("T-test-gen-scriptsa-tui", "T-test-gen-scripts") → "tui"
// e.g., ExtractTypeSuffix("T-test-gen-scripts", "T-test-gen-scripts") → ""
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
