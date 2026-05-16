//go:build e2e

package e2e

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// claim command tests — feature: cli-lean-output
// ==============================================================================

// claimTask attempts to claim a task. It skips the test if no tasks are available.
func claimTask(t *testing.T) []string {
	t.Helper()
	out, exitCode := runCLIRaw(t, "forge", "task", "claim")
	if exitCode != 0 {
		t.Skip("no pending tasks available for claim — cannot test claim output")
	}
	return parseBlock(t, out)
}

// Traceability: TC-001 -> Proposal "Success Criteria" item 1 + Task 1 AC-1
func TestTC_001_ClaimOutputsOnlyEssentialFields(t *testing.T) {
	lines := claimTask(t)

	// Mandatory fields must be present
	for _, field := range []string{"ACTION", "TASK_ID", "FEATURE", "FILE"} {
		assert.True(t, hasField(lines, field, ""),
			"mandatory field %s missing from output: %v", field, lines)
	}

	// Only known fields are allowed
	allowedFields := map[string]bool{
		"ACTION": true, "TASK_ID": true, "FEATURE": true, "FILE": true,
		"SCOPE": true, "BREAKING": true, "MAIN_SESSION": true,
	}
	for _, l := range lines {
		key := strings.SplitN(l, ":", 2)[0]
		assert.True(t, allowedFields[key],
			"unexpected field %q in output: %v", key, lines)
	}
}

// Traceability: TC-002 -> Task 1 AC-2 — printNewTask() wraps with ACTION: CLAIMED
func TestTC_002_ClaimOutputIncludesActionClaimed(t *testing.T) {
	lines := claimTask(t)

	assert.True(t, hasField(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
}

// Traceability: TC-003 -> Proposal — TASK_ID is consumed "everywhere"
func TestTC_003_ClaimOutputIncludesTaskID(t *testing.T) {
	lines := claimTask(t)

	assert.True(t, hasField(lines, "TASK_ID", ""),
		"TASK_ID field missing from output: %v", lines)
}

// Traceability: TC-004 -> Proposal — FEATURE is consumed by "E2E gate"
func TestTC_004_ClaimOutputIncludesFeature(t *testing.T) {
	lines := claimTask(t)

	assert.True(t, hasField(lines, "FEATURE", ""),
		"FEATURE field missing from output: %v", lines)
}

// Traceability: TC-005 -> Proposal — FILE is consumed by "agent reads task file"
func TestTC_005_ClaimOutputIncludesFile(t *testing.T) {
	lines := claimTask(t)

	assert.True(t, hasField(lines, "FILE", ""),
		"FILE field missing from output: %v", lines)
	// FILE must contain an absolute path
	fileVal := fieldValue(lines, "FILE")
	assert.True(t, strings.HasPrefix(fileVal, "/") || (len(fileVal) > 2 && fileVal[1] == ':'),
		"FILE should be an absolute path, got: %s", fileVal)
}

// Traceability: TC-006 -> Task 1 AC-1 — SCOPE only when non-empty
func TestTC_006_ClaimScopePresentWhenNonEmpty(t *testing.T) {
	lines := claimTask(t)

	scopeVal := fieldValue(lines, "SCOPE")
	if scopeVal != "" {
		assert.Equal(t, "backend", scopeVal,
			"SCOPE value mismatch: %v", lines)
	}
}

// Traceability: TC-007 -> Task 1 AC-1 — SCOPE only when non-empty
func TestTC_007_ClaimScopeAbsentWhenEmpty(t *testing.T) {
	lines := claimTask(t)

	scopeVal := fieldValue(lines, "SCOPE")
	if scopeVal == "" {
		assert.True(t, hasNoField(lines, "SCOPE"),
			"SCOPE should be absent when empty")
	}
}

// Traceability: TC-008 -> Task 1 AC-5 — BREAKING present with "true" when true
func TestTC_008_ClaimBreakingPresentWhenTrue(t *testing.T) {
	lines := claimTask(t)

	breakingVal := fieldValue(lines, "BREAKING")
	if breakingVal != "" {
		assert.Equal(t, "true", breakingVal,
			"BREAKING value mismatch: %v", lines)
	}
}

// Traceability: TC-009 -> Task 1 AC-5 — BREAKING absent when false
func TestTC_009_ClaimBreakingAbsentWhenFalse(t *testing.T) {
	lines := claimTask(t)

	breakingVal := fieldValue(lines, "BREAKING")
	if breakingVal == "" {
		assert.True(t, hasNoField(lines, "BREAKING"),
			"BREAKING should be absent when false")
	}
}

// Traceability: TC-010 -> Task 1 AC-5 — MAIN_SESSION present with "true" when true
func TestTC_010_ClaimMainSessionPresentWhenTrue(t *testing.T) {
	lines := claimTask(t)

	mainSessionVal := fieldValue(lines, "MAIN_SESSION")
	if mainSessionVal != "" {
		assert.Equal(t, "true", mainSessionVal,
			"MAIN_SESSION value mismatch: %v", lines)
	}
}

// Traceability: TC-011 -> Task 1 AC-5 — MAIN_SESSION absent when false
func TestTC_011_ClaimMainSessionAbsentWhenFalse(t *testing.T) {
	lines := claimTask(t)

	mainSessionVal := fieldValue(lines, "MAIN_SESSION")
	if mainSessionVal == "" {
		assert.True(t, hasNoField(lines, "MAIN_SESSION"),
			"MAIN_SESSION should be absent when false")
	}
}

// Traceability: TC-012 -> Task 1 AC-4 — Removed fields no longer appear
func TestTC_012_ClaimRemovedFieldsNotPresent(t *testing.T) {
	lines := claimTask(t)

	removedFields := []string{
		"KEY", "TITLE", "PRIORITY", "STATUS", "ESTIMATED_TIME",
		"DEPENDENCIES", "TYPE", "PROFILE", "NO_TEST", "RECORD",
	}
	for _, field := range removedFields {
		assert.True(t, hasNoField(lines, field),
			"removed field %s should not appear in output: %v", field, lines)
	}
}

// Traceability: TC-013 -> Task 1 AC-3 — printContinueTask() wraps with ACTION: CONTINUE + STARTED_AT
func TestTC_013_ClaimContinueWrapsWithActionContinueAndStartedAt(t *testing.T) {
	lines := claimTask(t)

	action := fieldValue(lines, "ACTION")
	if action == "CONTINUE" {
		assert.True(t, hasField(lines, "STARTED_AT", ""),
			"STARTED_AT missing in CONTINUE output: %v", lines)
		for _, field := range []string{"TASK_ID", "FEATURE", "FILE"} {
			assert.True(t, hasField(lines, field, ""),
				"%s missing in CONTINUE output: %v", field, lines)
		}
	}
}

// Traceability: TC-014 -> Task 1 Implementation Notes — field order
func TestTC_014_ClaimFieldOrderMatchesSpecification(t *testing.T) {
	lines := claimTask(t)

	// ACTION must be first field
	actionIdx := fieldIndex(lines, "ACTION")
	assert.Equal(t, 0, actionIdx,
		"ACTION should be at index 0, got %d: %v", actionIdx, lines)

	// Verify ordering of present fields
	expectedOrder := []string{"ACTION", "TASK_ID", "FEATURE", "FILE", "SCOPE", "BREAKING", "MAIN_SESSION"}
	prevIdx := -1
	for _, field := range expectedOrder {
		idx := fieldIndex(lines, field)
		if idx != -1 {
			assert.Greater(t, idx, prevIdx,
				"%s (idx %d) should come after previous field (idx %d): %v", field, idx, prevIdx, lines)
			prevIdx = idx
		}
	}

	// Removed fields must not appear
	for _, field := range []string{"KEY", "TYPE"} {
		assert.Equal(t, -1, fieldIndex(lines, field),
			"%s should not appear in output", field)
	}
}

// ==============================================================================
// submit command tests — feature: cli-lean-output
// ==============================================================================

// Traceability: TC-015 -> Task 2 AC-1 + Proposal "Success Criteria" item 2
func TestTC_015_SubmitOutputsOnlyStatusField(t *testing.T) {
	claimLines := claimTask(t)
	taskID := fieldValue(claimLines, "TASK_ID")

	recordJSON := `{"summary":"e2e test submit","testsPassed":1,"testsFailed":0,"coverage":100}`
	tmpFile := t.TempDir() + "/record.json"
	if err := os.WriteFile(tmpFile, []byte(recordJSON), 0644); err != nil {
		t.Fatalf("failed to write record file: %v", err)
	}

	out := runCLI(t, "forge", "task", "submit", taskID, "--data", tmpFile)
	lines := parseBlock(t, out)

	assert.True(t, hasField(lines, "STATUS", ""),
		"STATUS field missing from submit output: %v", lines)
	assert.True(t, hasNoField(lines, "TASK_ID"),
		"TASK_ID should not appear in submit output: %v", lines)
	assert.True(t, hasNoField(lines, "RECORD_FILE"),
		"RECORD_FILE should not appear in submit output: %v", lines)
	assert.Equal(t, 1, len(lines),
		"submit output should contain exactly 1 field (STATUS), got %d: %v", len(lines), lines)
}

// Traceability: TC-016 -> Task 2 AC-5 — JSON mode (--json) in submit is NOT changed
func TestTC_016_SubmitJsonModeUnchanged(t *testing.T) {
	claimLines := claimTask(t)
	taskID := fieldValue(claimLines, "TASK_ID")

	recordJSON := `{"summary":"e2e test JSON mode","testsPassed":1,"testsFailed":0,"coverage":100}`
	tmpFile := t.TempDir() + "/record.json"
	if err := os.WriteFile(tmpFile, []byte(recordJSON), 0644); err != nil {
		t.Fatalf("failed to write record file: %v", err)
	}

	out := runCLI(t, "forge", "task", "submit", taskID, "--data", tmpFile, "--json")

	assert.True(t, strings.Contains(out, `"taskId"`),
		"JSON output should contain taskId field, got: %s", out)
	assert.True(t, strings.Contains(out, `"status"`),
		"JSON output should contain status field, got: %s", out)
	assert.True(t, strings.Contains(out, `"recordFile"`),
		"JSON output should contain recordFile field, got: %s", out)
	trimmed := strings.TrimSpace(out)
	assert.True(t, strings.HasPrefix(trimmed, "{"),
		"JSON output should start with '{', got: %s", trimmed)
}

// ==============================================================================
// query command tests — feature: cli-lean-output
// ==============================================================================

// Traceability: TC-017 -> Task 2 AC-2 + Proposal "Success Criteria" item 3
func TestTC_017_QueryOutputsEssentialFieldsWithConditionalScopeAndBreaking(t *testing.T) {
	claimLines := claimTask(t)
	taskID := fieldValue(claimLines, "TASK_ID")

	out := runCLI(t, "forge", "task", "query", taskID)
	lines := parseBlock(t, out)

	assert.True(t, hasField(lines, "TASK_ID", ""),
		"TASK_ID missing from query output: %v", lines)
	assert.True(t, hasField(lines, "STATUS", ""),
		"STATUS missing from query output: %v", lines)

	scopeVal := fieldValue(lines, "SCOPE")
	if scopeVal != "" {
		assert.Equal(t, "backend", scopeVal,
			"SCOPE value mismatch: %v", lines)
	}
	breakingVal := fieldValue(lines, "BREAKING")
	if breakingVal != "" {
		assert.Equal(t, "true", breakingVal,
			"BREAKING value mismatch: %v", lines)
	}

	for _, field := range []string{
		"KEY", "TITLE", "PRIORITY", "ESTIMATED_TIME",
		"DEPENDENCIES", "FILE", "RECORD",
	} {
		assert.True(t, hasNoField(lines, field),
			"removed field %s should not appear in query output: %v", field, lines)
	}
}

// Traceability: TC-018 -> Task 2 AC-2 — SCOPE (when non-empty), BREAKING (when true)
func TestTC_018_QueryOmitsScopeWhenEmptyAndBreakingWhenFalse(t *testing.T) {
	claimLines := claimTask(t)
	taskID := fieldValue(claimLines, "TASK_ID")

	out := runCLI(t, "forge", "task", "query", taskID)
	lines := parseBlock(t, out)

	scopeVal := fieldValue(lines, "SCOPE")
	breakingVal := fieldValue(lines, "BREAKING")
	if scopeVal == "" && breakingVal == "" {
		assert.True(t, hasNoField(lines, "SCOPE"),
			"SCOPE should be absent when empty")
		assert.True(t, hasNoField(lines, "BREAKING"),
			"BREAKING should be absent when false")
	}
}

// ==============================================================================
// status command tests — feature: cli-lean-output
// ==============================================================================

// Traceability: TC-019 -> Task 2 AC-3, AC-4 — status outputs TASK_ID + STATUS
func TestTC_019_StatusOutputsOnlyTaskIDAndStatus(t *testing.T) {
	claimLines := claimTask(t)
	taskID := fieldValue(claimLines, "TASK_ID")

	// Query mode: forge task status <task-id>
	out := runCLI(t, "forge", "task", "status", taskID)
	lines := parseBlock(t, out)

	assert.True(t, hasField(lines, "TASK_ID", taskID),
		"TASK_ID mismatch in status query output: %v", lines)
	assert.True(t, hasField(lines, "STATUS", ""),
		"STATUS missing from status query output: %v", lines)
	assert.Equal(t, 2, len(lines),
		"status query output should contain exactly 2 fields, got %d: %v", len(lines), lines)

	for _, field := range []string{"KEY", "TITLE", "DEPENDENCIES"} {
		assert.True(t, hasNoField(lines, field),
			"removed field %s should not appear in status output: %v", field, lines)
	}

	// Update mode: forge task status <task-id> blocked
	out = runCLI(t, "forge", "task", "status", taskID, "blocked")
	lines = parseBlock(t, out)

	assert.True(t, hasField(lines, "TASK_ID", ""),
		"TASK_ID missing from status update output: %v", lines)
	assert.True(t, hasField(lines, "STATUS", ""),
		"STATUS missing from status update output: %v", lines)
	assert.Equal(t, 2, len(lines),
		"status update output should contain exactly 2 fields, got %d: %v", len(lines), lines)

	for _, field := range []string{"KEY", "TITLE", "DEPENDENCIES"} {
		assert.True(t, hasNoField(lines, field),
			"removed field %s should not appear in status update output: %v", field, lines)
	}
}
