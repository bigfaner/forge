package task

import "strings"

// InferType infers the task type from the task ID using pattern matching.
// Returns empty string for unknown IDs (no fallback).
//
// Handles type-suffixed IDs from per-type test pipeline split:
// T-test-gen-scripts-api, T-test-gen-scripts-cli (type suffix)
//
// The surfaces map is used to resolve T-test-run-{surface-key} IDs:
// the suffix after "T-test-run-" is looked up as a surface-key in the map.
// If the key exists, the ID is recognized as TypeTestRun.
// If the key does not exist, the ID falls through to exact-match logic.
// A nil surfaces map disables surface-key prefix matching entirely.
func InferType(id string, surfaces map[string]string) string {
	switch {
	case strings.HasSuffix(id, IDSuffixSummary):
		return TypeDocSummary
	case strings.HasSuffix(id, IDSuffixGate):
		return TypeGate
	case id == "T-test-gen-contracts":
		return TypeTestGenContracts
	case id == "T-test-gen-journeys", typeSuffixedID(id, "T-test-gen-journeys"):
		return TypeTestGenJourneys
	case id == "T-test-gen-scripts", typeSuffixedID(id, "T-test-gen-scripts"):
		return TypeTestGenScripts
	case id == "T-test-run":
		return TypeTestRun
	case testRunSurfaceKeyMatch(id, surfaces):
		return TypeTestRun
	case id == "T-test-verify-regression":
		return "test.verify-regression"
	case id == "T-specs-consolidate":
		return TypeDocConsolidate
	case id == "T-validate-code":
		return TypeValidationCode
	case id == "T-validate-ux":
		return TypeValidationUx
	case id == "T-quick-doc-drift":
		return TypeDocDrift
	case id == "T-clean-code":
		return TypeCleanCode
	case id == "T-eval-journey":
		return TypeEvalJourney
	case id == "T-eval-contract":
		return TypeEvalContract
	case strings.HasPrefix(id, "fix-") || strings.HasPrefix(id, "disc-"):
		return TypeCodingFix
	case id == "T-review-doc":
		return TypeDocReview
	default:
		return ""
	}
}

// testRunSurfaceKeyMatch checks if id matches "T-test-run-" + a known surface-key.
// Returns true if the suffix after "T-test-run-" is a key in the surfaces map.
func testRunSurfaceKeyMatch(id string, surfaces map[string]string) bool {
	const prefix = "T-test-run-"
	if !strings.HasPrefix(id, prefix) {
		return false
	}
	suffix := id[len(prefix):]
	if suffix == "" {
		return false
	}
	_, ok := surfaces[suffix]
	return ok
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
