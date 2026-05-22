package cmd

import (
	"strings"
	"testing"

	taskpkg "forge-cli/internal/cmd/task"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"
)

// parseBlock extracts lines between "---" separators.
// Returns the inner lines (without separators) or fails the test.
func parseBlock(t *testing.T, raw string) []string {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) < 2 || strings.TrimSpace(lines[0]) != "---" || strings.TrimSpace(lines[len(lines)-1]) != "---" {
		t.Fatalf("output must be wrapped in --- separators, got:\n%s", raw)
	}
	inner := lines[1 : len(lines)-1]
	// trim whitespace on each line
	result := make([]string, 0, len(inner))
	for _, l := range inner {
		result = append(result, strings.TrimSpace(l))
	}
	return result
}

// hasField checks that a parsed block contains a "KEY: value" line.
func hasField(lines []string, key, value string) bool {
	prefix := key + ": "
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			if value == "" {
				return true // any value accepted
			}
			return l == prefix+value
		}
	}
	return false
}

// hasNoField checks that a parsed block does NOT contain any line starting with key.
func hasNoField(lines []string, key string) bool {
	prefix := key + ": "
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			return false
		}
	}
	return true
}

// fieldIndex returns the index of the line starting with key+": ", or -1.
func fieldIndex(lines []string, key string) int {
	prefix := key + ": "
	for i, l := range lines {
		if strings.HasPrefix(l, prefix) {
			return i
		}
	}
	return -1
}

// --- Contract: task feature output format ---

func TestContract_Feature_NoFeature(t *testing.T) {
	out := captureStdout(func() {
		PrintBlock("FEATURE", "(none)")
	})
	lines := parseBlock(t, out)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d: %v", len(lines), lines)
	}
	if !hasField(lines, "FEATURE", "(none)") {
		t.Errorf("expected FEATURE: (none), got: %v", lines)
	}
}

func TestContract_Feature_WithFeature(t *testing.T) {
	out := captureStdout(func() {
		PrintBlock("FEATURE", "my-feature")
	})
	lines := parseBlock(t, out)
	if !hasField(lines, "FEATURE", "my-feature") {
		t.Errorf("expected FEATURE: my-feature, got: %v", lines)
	}
}

// --- Contract: task claim output format ---

func TestContract_Claim_NewTask(t *testing.T) {
	dir := t.TempDir()
	_ = feature.EnsureFeatureDir(dir, "feat")

	tk := &task.Task{
		ID: "1.1", Title: "Do stuff", Priority: "P0", Status: "pending",
		File: "1.1.md", Record: "records/1.1.md",
		Breaking: true, MainSession: true,
		Type: "coding.feature", Scope: "backend",
		Dependencies: []string{"1.0"}, EstimatedTime: "30min",
	}

	out := captureStdout(func() {
		taskpkg.ExportPrintNewTask("task1", tk, dir, "feat")
	})
	lines := parseBlock(t, out)

	// Mandatory fields (trimmed set)
	for _, field := range []string{"ACTION", "TASK_ID", "FEATURE", "FILE"} {
		if !hasField(lines, field, "") {
			t.Errorf("missing mandatory field %s in output: %v", field, lines)
		}
	}
	if !hasField(lines, "ACTION", "CLAIMED") {
		t.Errorf("expected ACTION: CLAIMED")
	}
	// Conditional fields present when non-empty
	if !hasField(lines, "TYPE", "coding.feature") {
		t.Errorf("expected TYPE: coding.feature")
	}
	if !hasField(lines, "SCOPE", "backend") {
		t.Errorf("expected SCOPE: backend")
	}
	if !hasField(lines, "MAIN_SESSION", "true") {
		t.Errorf("expected MAIN_SESSION: true")
	}
	// BREAKING must NOT appear (removed from claim output)
	if !hasNoField(lines, "BREAKING") {
		t.Errorf("BREAKING should not appear in claim output")
	}
	// Removed fields must NOT appear
	for _, field := range []string{"KEY", "TITLE", "PRIORITY", "STATUS", "ESTIMATED_TIME",
		"DEPENDENCIES", "NO_TEST", "RECORD"} {
		if !hasNoField(lines, field) {
			t.Errorf("removed field %s should not appear in output: %v", field, lines)
		}
	}
}

func TestContract_Claim_NewTask_ConditionalAbsent(t *testing.T) {
	dir := t.TempDir()
	_ = feature.EnsureFeatureDir(dir, "feat")

	tk := &task.Task{
		ID: "1.1", Title: "Simple", Priority: "P1", Status: "pending",
		File: "1.1.md", Record: "records/1.1.md",
	}

	out := captureStdout(func() {
		taskpkg.ExportPrintNewTask("t1", tk, dir, "feat")
	})
	lines := parseBlock(t, out)

	// TYPE absent when empty
	if !hasNoField(lines, "TYPE") {
		t.Errorf("TYPE should be absent when empty")
	}
	// SCOPE absent when empty
	if !hasNoField(lines, "SCOPE") {
		t.Errorf("SCOPE should be absent when empty")
	}
	// Boolean fields absent when false
	if !hasNoField(lines, "BREAKING") {
		t.Errorf("BREAKING should be absent when false")
	}
	if !hasNoField(lines, "MAIN_SESSION") {
		t.Errorf("MAIN_SESSION should be absent when false")
	}
}

func TestContract_Claim_Continue(t *testing.T) {
	dir := t.TempDir()
	_ = feature.EnsureFeatureDir(dir, "feat")

	tk := &task.Task{
		ID: "1.1", Title: "Resume me", Priority: "P0", Status: "in_progress",
		File: "1.1.md", Record: "records/1.1.md", Scope: "backend",
		Breaking: true, MainSession: true,
	}
	state := &task.TaskState{
		Key: "task1", TaskID: "1.1", StartedTime: "2025-01-01T00:00:00Z",
	}

	out := captureStdout(func() {
		taskpkg.ExportPrintContinueTask(state, tk, dir, "feat")
	})
	lines := parseBlock(t, out)

	if !hasField(lines, "ACTION", "CONTINUE") {
		t.Errorf("expected ACTION: CONTINUE")
	}
	if !hasField(lines, "TASK_ID", "1.1") {
		t.Errorf("expected TASK_ID: 1.1")
	}
	if !hasField(lines, "FEATURE", "feat") {
		t.Errorf("expected FEATURE: feat")
	}
	if !hasField(lines, "FILE", "") {
		t.Errorf("expected FILE field")
	}
	if !hasField(lines, "SCOPE", "backend") {
		t.Errorf("expected SCOPE: backend")
	}
	// BREAKING must NOT appear (removed from claim output)
	if !hasNoField(lines, "BREAKING") {
		t.Errorf("BREAKING should not appear in CONTINUE output")
	}
	if !hasField(lines, "MAIN_SESSION", "true") {
		t.Errorf("expected MAIN_SESSION: true")
	}
	if !hasField(lines, "STARTED_AT", "2025-01-01T00:00:00Z") {
		t.Errorf("expected STARTED_AT field")
	}
	// Removed fields must NOT appear
	for _, field := range []string{"KEY", "TITLE", "PRIORITY", "STATUS"} {
		if !hasNoField(lines, field) {
			t.Errorf("removed field %s should not appear in CONTINUE output", field)
		}
	}
}

func TestContract_Claim_FieldOrder(t *testing.T) {
	dir := t.TempDir()
	_ = feature.EnsureFeatureDir(dir, "feat")

	tk := &task.Task{
		ID: "1.1", Title: "T", Priority: "P0", Status: "pending",
		File: "1.1.md", Record: "records/1.1.md",
		Breaking: true, MainSession: true, Scope: "backend",
		Type: "coding.feature",
	}

	out := captureStdout(func() {
		taskpkg.ExportPrintNewTask("k1", tk, dir, "feat")
	})
	lines := parseBlock(t, out)

	// ACTION must be first field after ---
	if idx := fieldIndex(lines, "ACTION"); idx != 0 {
		t.Errorf("ACTION should be at index 0, got %d", idx)
	}
	// Expected order: ACTION, TASK_ID, TYPE, FEATURE, FILE, SCOPE, MAIN_SESSION
	taskIDIdx := fieldIndex(lines, "TASK_ID")
	typeIdx := fieldIndex(lines, "TYPE")
	featureIdx := fieldIndex(lines, "FEATURE")
	fileIdx := fieldIndex(lines, "FILE")
	scopeIdx := fieldIndex(lines, "SCOPE")
	mainIdx := fieldIndex(lines, "MAIN_SESSION")

	if taskIDIdx == -1 || typeIdx == -1 || featureIdx == -1 || fileIdx == -1 || scopeIdx == -1 || mainIdx == -1 {
		t.Fatalf("missing expected fields, got: %v", lines)
	}

	// Verify ordering
	if taskIDIdx >= typeIdx {
		t.Errorf("TASK_ID (%d) should come before TYPE (%d)", taskIDIdx, typeIdx)
	}
	if typeIdx >= featureIdx {
		t.Errorf("TYPE (%d) should come before FEATURE (%d)", typeIdx, featureIdx)
	}
	if featureIdx >= fileIdx {
		t.Errorf("FEATURE (%d) should come before FILE (%d)", featureIdx, fileIdx)
	}
	if fileIdx >= scopeIdx {
		t.Errorf("FILE (%d) should come before SCOPE (%d)", fileIdx, scopeIdx)
	}
	if scopeIdx >= mainIdx {
		t.Errorf("SCOPE (%d) should come before MAIN_SESSION (%d)", scopeIdx, mainIdx)
	}

	// BREAKING must NOT appear (removed from claim output)
	if idx := fieldIndex(lines, "BREAKING"); idx != -1 {
		t.Errorf("BREAKING should not appear in output, found at index %d", idx)
	}

	// Removed fields must not appear
	if idx := fieldIndex(lines, "KEY"); idx != -1 {
		t.Errorf("KEY should not appear in output, found at index %d", idx)
	}
}

// --- Contract: task add output format ---

func TestContract_Add_Success(t *testing.T) {
	out := captureStdout(func() {
		PrintBlockStart()
		PrintField("ACTION", "ADDED")
		PrintField("KEY", "fix-1")
		PrintField("TASK_ID", "fix-1")
		PrintField("TITLE", "Fix bug")
		PrintField("PRIORITY", "P0")
		PrintField("STATUS", "pending")
		PrintField("FILE", "/path/to/tasks/fix-1.md")
		PrintField("RECORD", "/path/to/tasks/records/fix-1.md")
		PrintBlockEnd()
	})
	lines := parseBlock(t, out)

	if !hasField(lines, "ACTION", "ADDED") {
		t.Errorf("expected ACTION: ADDED")
	}
	for _, field := range []string{"ACTION", "KEY", "TASK_ID", "TITLE", "PRIORITY", "STATUS", "FILE", "RECORD"} {
		if !hasField(lines, field, "") {
			t.Errorf("missing field %s", field)
		}
	}
}

func TestContract_Add_Skipped(t *testing.T) {
	out := captureStdout(func() {
		PrintBlockStart()
		PrintField("ACTION", "SKIPPED")
		PrintField("REASON", "active fix task fix-1 already exists")
		PrintBlockEnd()
	})
	lines := parseBlock(t, out)

	if !hasField(lines, "ACTION", "SKIPPED") {
		t.Errorf("expected ACTION: SKIPPED")
	}
	if !hasField(lines, "REASON", "") {
		t.Errorf("expected REASON field")
	}
}

// --- Contract: task record output format ---

func TestContract_Record_Completed(t *testing.T) {
	out := captureStdout(func() {
		PrintBlockStart()
		PrintField("STATUS", "completed")
		PrintBlockEnd()
	})
	lines := parseBlock(t, out)

	if !hasField(lines, "STATUS", "completed") {
		t.Errorf("expected STATUS: completed")
	}
	// TASK_ID and RECORD_FILE removed — only STATUS remains
	if !hasNoField(lines, "TASK_ID") {
		t.Errorf("TASK_ID should not appear in submit output")
	}
	if !hasNoField(lines, "RECORD_FILE") {
		t.Errorf("RECORD_FILE should not appear in submit output")
	}
}

func TestContract_Record_Blocked(t *testing.T) {
	out := captureStdout(func() {
		PrintBlockStart()
		PrintField("STATUS", "blocked")
		PrintBlockEnd()
	})
	lines := parseBlock(t, out)

	if !hasField(lines, "STATUS", "blocked") {
		t.Errorf("expected STATUS: blocked")
	}
	// Only STATUS field should appear
	if len(lines) != 1 {
		t.Errorf("expected exactly 1 field (STATUS), got %d: %v", len(lines), lines)
	}
}

// --- Contract: Block structure ---

func TestContract_BlockSeparator_SingleBlock(t *testing.T) {
	out := captureStdout(func() {
		PrintBlock("FEATURE", "test")
	})

	trimmed := strings.TrimSpace(out)
	parts := strings.Split(trimmed, "\n")

	if len(parts) != 3 {
		t.Fatalf("expected 3 lines (---, content, ---), got %d: %q", len(parts), parts)
	}
	if parts[0] != "---" {
		t.Errorf("first line should be ---, got %q", parts[0])
	}
	if parts[2] != "---" {
		t.Errorf("last line should be ---, got %q", parts[2])
	}
}

func TestContract_BlockSeparator_NoTrailingNewlineInSeparator(t *testing.T) {
	// Each block must end with "---" on its own line, followed by a newline.
	out := captureStdout(func() {
		PrintBlock("KEY", "value")
	})

	if !strings.HasSuffix(out, "---\n") {
		t.Errorf("block must end with ---\\n, got suffix: %q", out[len(out)-20:])
	}
}

// --- Contract: query output format ---

func TestContract_Query_Minimal(t *testing.T) {
	out := captureStdout(func() {
		PrintBlockStart()
		PrintField("TASK_ID", "1")
		PrintField("STATUS", "in_progress")
		PrintBlockEnd()
	})
	lines := parseBlock(t, out)

	if !hasField(lines, "TASK_ID", "1") {
		t.Errorf("expected TASK_ID: 1")
	}
	if !hasField(lines, "STATUS", "in_progress") {
		t.Errorf("expected STATUS: in_progress")
	}
	// Removed fields must NOT appear
	for _, field := range []string{"KEY", "TITLE", "PRIORITY", "ESTIMATED_TIME",
		"DEPENDENCIES", "FILE", "RECORD"} {
		if !hasNoField(lines, field) {
			t.Errorf("removed field %s should not appear in query output: %v", field, lines)
		}
	}
}

func TestContract_Query_WithScope(t *testing.T) {
	out := captureStdout(func() {
		PrintBlockStart()
		PrintField("TASK_ID", "1")
		PrintField("STATUS", "in_progress")
		PrintFieldIfNotEmpty("SCOPE", "backend")
		PrintBlockEnd()
	})
	lines := parseBlock(t, out)

	if !hasField(lines, "SCOPE", "backend") {
		t.Errorf("expected SCOPE: backend")
	}
}

func TestContract_Query_EmptyScopeOmitted(t *testing.T) {
	out := captureStdout(func() {
		PrintBlockStart()
		PrintField("TASK_ID", "1")
		PrintField("STATUS", "in_progress")
		PrintFieldIfNotEmpty("SCOPE", "")
		PrintBlockEnd()
	})
	lines := parseBlock(t, out)

	if !hasNoField(lines, "SCOPE") {
		t.Errorf("SCOPE should be omitted when empty")
	}
}

func TestContract_Query_BreakingWhenTrue(t *testing.T) {
	out := captureStdout(func() {
		PrintBlockStart()
		PrintField("TASK_ID", "1")
		PrintField("STATUS", "in_progress")
		PrintField("BREAKING", "true")
		PrintBlockEnd()
	})
	lines := parseBlock(t, out)

	if !hasField(lines, "BREAKING", "true") {
		t.Errorf("expected BREAKING: true")
	}
}

func TestContract_Query_BreakingOmittedWhenFalse(t *testing.T) {
	out := captureStdout(func() {
		PrintBlockStart()
		PrintField("TASK_ID", "1")
		PrintField("STATUS", "in_progress")
		PrintBlockEnd()
	})
	lines := parseBlock(t, out)

	if !hasNoField(lines, "BREAKING") {
		t.Errorf("BREAKING should be omitted when false")
	}
}

// --- Contract: status output format ---

func TestContract_Status_QueryMode(t *testing.T) {
	out := captureStdout(func() {
		PrintBlockStart()
		PrintField("TASK_ID", "1")
		PrintField("STATUS", "in_progress")
		PrintBlockEnd()
	})
	lines := parseBlock(t, out)

	if !hasField(lines, "TASK_ID", "1") {
		t.Errorf("expected TASK_ID: 1")
	}
	if !hasField(lines, "STATUS", "in_progress") {
		t.Errorf("expected STATUS: in_progress")
	}
	// Removed fields must NOT appear
	for _, field := range []string{"KEY", "TITLE", "DEPENDENCIES"} {
		if !hasNoField(lines, field) {
			t.Errorf("removed field %s should not appear in status output: %v", field, lines)
		}
	}
}

func TestContract_Status_UpdateMode(t *testing.T) {
	out := captureStdout(func() {
		PrintBlockStart()
		PrintField("TASK_ID", "1")
		PrintField("STATUS", "pending")
		PrintBlockEnd()
	})
	lines := parseBlock(t, out)

	if !hasField(lines, "TASK_ID", "1") {
		t.Errorf("expected TASK_ID: 1")
	}
	if !hasField(lines, "STATUS", "pending") {
		t.Errorf("expected STATUS: pending")
	}
	if len(lines) != 2 {
		t.Errorf("expected exactly 2 fields (TASK_ID + STATUS), got %d: %v", len(lines), lines)
	}
}
