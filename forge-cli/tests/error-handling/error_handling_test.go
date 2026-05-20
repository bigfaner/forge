//go:build e2e

package errorhandling

import (
	"strings"
	"testing"

	"forge-cli/tests/e2e/testkit"

	"github.com/stretchr/testify/assert"
)

// --- Error Handling (TC-031 to TC-041) ---

// Traceability: TC-031 -> Spec Error Handling Table
func TestTC_031_TaskClaimNoAvailableTasksReturnsError(t *testing.T) {
	t.Skip("requires manual setup: feature with no available tasks to claim")
}

// Traceability: TC-032 -> Spec Error Handling Table
func TestTC_032_TaskClaimCorruptedIndexReturnsError(t *testing.T) {
	t.Skip("requires manual setup: corrupted or missing index.json")
}

// Traceability: TC-033 -> Spec Error Handling Table
func TestTC_033_TaskCheckDepsUnmetDependencyReturnsError(t *testing.T) {
	t.Skip("requires manual setup: task with unmet dependency in index.json")
}

// Traceability: TC-034 -> Spec Error Handling Table
func TestTC_034_TaskValidateIndexInvalidSchemaReturnsError(t *testing.T) {
	t.Skip("requires manual setup: index.json with schema validation errors")
}

// Traceability: TC-035 -> Spec Error Handling Table
func TestTC_035_TaskStatusNonexistentIDReturnsError(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("task", "status", "NONEXISTENT-999")

	assert.Equal(t, 1, exitCode, "status for nonexistent task should exit 1")
	lower := strings.ToLower(out)
	assert.True(t,
		strings.Contains(lower, "not found"),
		"output should contain 'not found': %s", out)
}

// Traceability: TC-036 -> Spec Error Handling Table
func TestTC_036_ForensicSearchNoResultsReturnsEmptyOutput(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("forensic", "search",
		"--keyword", "ZZZ_NO_MATCH_POSSIBLE_XYZ_12345",
		"--last", "1")

	assert.Equal(t, 0, exitCode, "forensic search with no matches should exit 0")
	assert.Equal(t, "[]", strings.TrimSpace(out),
		"forensic search with no matches should return empty array")
}

// Traceability: TC-037 -> Spec Error Handling Table
func TestTC_037_ForensicSearchMissingRecordsDirReturnsError(t *testing.T) {
	t.Skip("requires manual setup: missing ~/.claude/history.jsonl")
}

// Traceability: TC-038 -> Spec Error Handling Table
func TestTC_038_VerifyTaskDoneIncompleteTasksReturnsError(t *testing.T) {
	t.Skip("requires manual setup: feature with incomplete tasks and active state.json")
}

// Traceability: TC-039 -> Spec Error Handling Table
func TestTC_039_TaskSubmitConcurrentWriteConflictReturnsRetryError(t *testing.T) {
	t.Skip("requires manual setup: concurrent lock contention scenario")
}

// Traceability: TC-040 -> Spec Error Handling Table
func TestTC_040_TaskSubmitMissingIndexReturnsError(t *testing.T) {
	t.Skip("requires manual setup: feature directory without index.json")
}

// assertTaskNotFound verifies the standard task-not-found error pattern.
func assertTaskNotFound(t *testing.T, output, taskID string) {
	t.Helper()
	lower := strings.ToLower(output)
	assert.True(t,
		strings.Contains(lower, "not found") && strings.Contains(lower, strings.ToLower(taskID)),
		"output should mention task not found with ID %s: %s", taskID, output)
}

// assertIndexLoadError verifies the standard index-load error pattern.
func assertIndexLoadError(t *testing.T, output string) {
	t.Helper()
	lower := strings.ToLower(output)
	assert.True(t,
		strings.Contains(lower, "failed") || strings.Contains(lower, "not found"),
		"output should mention index load failure: %s", output)
}
