package cmd

import (
	"strings"
	"testing"

	"task-cli/pkg/feature"
	"task-cli/pkg/task"
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
		Breaking: true, MainSession: true, NoTest: false,
		Type: "implementation", Scope: "backend",
		Dependencies: []string{"1.0"}, EstimatedTime: "30min",
	}

	out := captureStdout(func() {
		printNewTask("task1", tk, dir, "feat")
	})
	lines := parseBlock(t, out)

	// Mandatory fields
	for _, field := range []string{"ACTION", "KEY", "TASK_ID", "TITLE", "PRIORITY",
		"STATUS", "BREAKING", "MAIN_SESSION", "TYPE", "NO_TEST", "FEATURE", "FILE", "RECORD"} {
		if !hasField(lines, field, "") {
			t.Errorf("missing mandatory field %s in output: %v", field, lines)
		}
	}
	if !hasField(lines, "ACTION", "CLAIMED") {
		t.Errorf("expected ACTION: CLAIMED")
	}
	if !hasField(lines, "KEY", "task1") {
		t.Errorf("expected KEY: task1")
	}
	// Conditional fields present when non-empty
	if !hasField(lines, "ESTIMATED_TIME", "30min") {
		t.Errorf("expected ESTIMATED_TIME: 30min")
	}
	if !hasField(lines, "DEPENDENCIES", "1.0") {
		t.Errorf("expected DEPENDENCIES: 1.0")
	}
	if !hasField(lines, "SCOPE", "backend") {
		t.Errorf("expected SCOPE: backend")
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
		printNewTask("t1", tk, dir, "feat")
	})
	lines := parseBlock(t, out)

	// Conditional fields absent when empty
	if !hasNoField(lines, "ESTIMATED_TIME") {
		t.Errorf("ESTIMATED_TIME should be absent when empty")
	}
	if !hasNoField(lines, "DEPENDENCIES") {
		t.Errorf("DEPENDENCIES should be absent when empty")
	}
	if !hasNoField(lines, "SCOPE") {
		t.Errorf("SCOPE should be absent when empty")
	}
	if !hasNoField(lines, "PROFILE") {
		t.Errorf("PROFILE should be absent when empty")
	}
}

func TestContract_Claim_ProfilePresent(t *testing.T) {
	dir := t.TempDir()
	_ = feature.EnsureFeatureDir(dir, "feat")

	tk := &task.Task{
		ID: "1.1", Title: "Test task", Priority: "P0", Status: "pending",
		File: "1.1.md", Record: "records/1.1.md",
		Profile: "go-test",
	}

	out := captureStdout(func() {
		printNewTask("t1", tk, dir, "feat")
	})
	lines := parseBlock(t, out)

	if !hasField(lines, "PROFILE", "go-test") {
		t.Errorf("expected PROFILE: go-test, got: %v", lines)
	}
}

func TestContract_Claim_ProfileAbsent(t *testing.T) {
	dir := t.TempDir()
	_ = feature.EnsureFeatureDir(dir, "feat")

	tk := &task.Task{
		ID: "1.1", Title: "Biz task", Priority: "P0", Status: "pending",
		File: "1.1.md", Record: "records/1.1.md",
	}

	out := captureStdout(func() {
		printNewTask("t1", tk, dir, "feat")
	})
	lines := parseBlock(t, out)

	if !hasNoField(lines, "PROFILE") {
		t.Errorf("PROFILE should be absent when empty, got: %v", lines)
	}
}

func TestContract_Claim_Continue(t *testing.T) {
	dir := t.TempDir()
	_ = feature.EnsureFeatureDir(dir, "feat")

	tk := &task.Task{
		ID: "1.1", Title: "Resume me", Priority: "P0", Status: "in_progress",
		File: "1.1.md", Record: "records/1.1.md",
	}
	state := &task.TaskState{
		Key: "task1", TaskID: "1.1", StartedTime: "2025-01-01T00:00:00Z",
	}

	out := captureStdout(func() {
		printContinueTask(state, tk, dir, "feat")
	})
	lines := parseBlock(t, out)

	if !hasField(lines, "ACTION", "CONTINUE") {
		t.Errorf("expected ACTION: CONTINUE")
	}
	if !hasField(lines, "STARTED_AT", "2025-01-01T00:00:00Z") {
		t.Errorf("expected STARTED_AT field")
	}
}

func TestContract_Claim_FieldOrder(t *testing.T) {
	dir := t.TempDir()
	_ = feature.EnsureFeatureDir(dir, "feat")

	tk := &task.Task{
		ID: "1.1", Title: "T", Priority: "P0", Status: "pending",
		File: "1.1.md", Record: "records/1.1.md", Type: "fix",
	}

	out := captureStdout(func() {
		printNewTask("k1", tk, dir, "feat")
	})
	lines := parseBlock(t, out)

	// ACTION must be first field after ---
	if idx := fieldIndex(lines, "ACTION"); idx != 0 {
		t.Errorf("ACTION should be at index 0, got %d", idx)
	}
	// KEY before TASK_ID
	keyIdx := fieldIndex(lines, "KEY")
	idIdx := fieldIndex(lines, "TASK_ID")
	if keyIdx == -1 || idIdx == -1 || keyIdx >= idIdx {
		t.Errorf("KEY (%d) should come before TASK_ID (%d)", keyIdx, idIdx)
	}
	// BREAKING before MAIN_SESSION before TYPE
	brkIdx := fieldIndex(lines, "BREAKING")
	mainIdx := fieldIndex(lines, "MAIN_SESSION")
	typeIdx := fieldIndex(lines, "TYPE")
	if brkIdx == -1 || mainIdx == -1 || typeIdx == -1 {
		t.Fatalf("missing BREAKING/MAIN_SESSION/TYPE fields")
	}
	if brkIdx >= mainIdx {
		t.Errorf("BREAKING (%d) should come before MAIN_SESSION (%d)", brkIdx, mainIdx)
	}
	if mainIdx >= typeIdx {
		t.Errorf("MAIN_SESSION (%d) should come before TYPE (%d)", mainIdx, typeIdx)
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
		PrintField("TASK_ID", "1.1")
		PrintField("RECORD_FILE", "/path/to/records/1.1.md")
		PrintField("STATUS", "completed")
		PrintBlockEnd()
	})
	lines := parseBlock(t, out)

	if !hasField(lines, "TASK_ID", "1.1") {
		t.Errorf("expected TASK_ID: 1.1")
	}
	if !hasField(lines, "RECORD_FILE", "") {
		t.Errorf("expected RECORD_FILE field")
	}
	if !hasField(lines, "STATUS", "completed") {
		t.Errorf("expected STATUS: completed")
	}
}

func TestContract_Record_Blocked(t *testing.T) {
	out := captureStdout(func() {
		PrintBlockStart()
		PrintField("TASK_ID", "1.1")
		PrintField("RECORD_FILE", "/path/to/records/1.1.md")
		PrintField("STATUS", "blocked")
		PrintBlockEnd()
	})
	lines := parseBlock(t, out)

	if !hasField(lines, "STATUS", "blocked") {
		t.Errorf("expected STATUS: blocked")
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
