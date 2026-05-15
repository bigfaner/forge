package task

import "strings"

// InferType infers the task type from the task ID using pattern matching.
// Returns empty string for unknown IDs (no fallback).
//
// Handles profile-suffixed IDs from the test profile system (D9):
// single profile: T-test-2, T-test-3, T-test-4 (no suffix)
// multi profile:  T-test-2a, T-test-2b, T-test-3a, ... (letter suffix)
func InferType(id string) string {
	switch {
	case strings.HasSuffix(id, ".summary"):
		return TypeDocGenerationSummary
	case strings.HasSuffix(id, ".gate"):
		return TypeGate
	case id == "T-test-1b":
		return TypeTestPipelineEvalCases
	case id == "T-test-1", profileSuffixedID(id, "T-test-1"):
		return TypeTestPipelineGenCases
	case id == "T-test-2", profileSuffixedID(id, "T-test-2"):
		return TypeTestPipelineGenScripts
	case id == "T-test-3", profileSuffixedID(id, "T-test-3"):
		return TypeTestPipelineRun
	case id == "T-test-4", profileSuffixedID(id, "T-test-4"):
		return TypeTestPipelineGraduate
	case id == "T-test-4.5":
		return TypeTestPipelineVerifyRegression
	case id == "T-test-5":
		return TypeDocGenerationConsolidate
	case profileSuffixedID(id, "T-quick-1"):
		return TypeTestPipelineGenCases
	case profileSuffixedID(id, "T-quick-2"):
		return TypeTestPipelineGenScripts
	case profileSuffixedID(id, "T-quick-3"):
		return TypeTestPipelineRun
	case profileSuffixedID(id, "T-quick-4"):
		return TypeTestPipelineGraduate
	case profileSuffixedID(id, "T-quick-5"):
		return TypeTestPipelineVerifyRegression
	case id == "T-quick-6", profileSuffixedID(id, "T-quick-6"):
		return TypeDocGenerationDrift
	case strings.HasPrefix(id, "fix-") || strings.HasPrefix(id, "disc-"):
		return TypeFix
	case id == "T-eval-doc":
		return TypeDocEvaluation
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
