//go:build e2e

package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Task lifecycle hardening tests — feature: task-lifecycle-hardening
// Tests verify self-block (SourceTaskID==selfID), lazy unblock scan,
// block-source lifecycle, and auto-downgrade unblock scenarios during
// forge task claim.
// ==============================================================================

// tlhTaskEntry represents a single task in the index.json fixture.
type tlhTaskEntry struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Priority      string   `json:"priority"`
	EstimatedTime string   `json:"estimatedTime,omitempty"`
	Dependencies  []string `json:"dependencies,omitempty"`
	Status        string   `json:"status"`
	File          string   `json:"file"`
	Record        string   `json:"record"`
	Type          string   `json:"type,omitempty"`
	SourceTaskID  string   `json:"sourceTaskID,omitempty"`
	BlockedReason string   `json:"blockedReason,omitempty"`
}

// tlhIndexFixture represents the top-level index.json structure.
type tlhIndexFixture struct {
	Feature      string                  `json:"feature"`
	StatusEnum   []string                `json:"statusEnum"`
	PriorityEnum []string                `json:"priorityEnum"`
	Tasks        map[string]tlhTaskEntry `json:"tasks"`
}

// tlhSetupFeatureFixture creates a temp project root with the given tasks in
// index.json under docs/features/task-lifecycle-hardening/tasks/.
// Returns the temp dir path to use as CLAUDE_PROJECT_DIR.
func tlhSetupFeatureFixture(t *testing.T, tasks map[string]tlhTaskEntry) string {
	t.Helper()
	dir := t.TempDir()

	tasksDir := filepath.Join(dir, "docs", "features", "task-lifecycle-hardening", "tasks")
	processDir := filepath.Join(tasksDir, "process")
	if err := os.MkdirAll(processDir, 0755); err != nil {
		t.Fatalf("failed to create process dir: %v", err)
	}

	// Create task files referenced by each entry
	for _, task := range tasks {
		taskFile := filepath.Join(tasksDir, task.File)
		if err := os.MkdirAll(filepath.Dir(taskFile), 0755); err != nil {
			t.Fatalf("failed to create task file dir: %v", err)
		}
		if err := os.WriteFile(taskFile, []byte("# Task "+task.ID), 0644); err != nil {
			t.Fatalf("failed to write task file: %v", err)
		}
	}

	// Create go.mod so project root detection works
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	idx := tlhIndexFixture{
		Feature:      "task-lifecycle-hardening",
		StatusEnum:   []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"},
		PriorityEnum: []string{"P0", "P1", "P2"},
		Tasks:        tasks,
	}
	idxData, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal index.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tasksDir, "index.json"), idxData, 0644); err != nil {
		t.Fatalf("failed to write index.json: %v", err)
	}

	return dir
}

// tlhForgeClaim runs "forge task claim" with CLAUDE_PROJECT_DIR set.
// Returns combined output and exit code. Does NOT fatalf on failure.
func tlhForgeClaim(t *testing.T, projectRoot string) (string, int) {
	t.Helper()
	cmd := exec.Command("forge", "task", "claim")
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return string(out), exitCode
}

// tlhParseBlock extracts key:value lines between "---" separators from raw CLI output.
// Strips any auto-unblock log lines ("Auto-unblocked task ...") that appear before the block.
func tlhParseBlock(t *testing.T, raw string) []string {
	t.Helper()
	// Strip any lines before the first "---" separator (e.g. auto-unblock logs)
	idx := strings.Index(raw, "\n---")
	if idx >= 0 && !strings.HasPrefix(strings.TrimSpace(raw), "---") {
		raw = raw[idx+1:] // keep the "---" line
	}
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) < 2 || strings.TrimSpace(lines[0]) != "---" || strings.TrimSpace(lines[len(lines)-1]) != "---" {
		t.Fatalf("output must be wrapped in --- separators, got:\n%s", raw)
	}
	result := make([]string, 0, len(lines)-2)
	for _, l := range lines[1 : len(lines)-1] {
		result = append(result, strings.TrimSpace(l))
	}
	return result
}

// tlhGetFieldValue returns the value for the given key from parsed block lines.
func tlhGetFieldValue(lines []string, key string) string {
	prefix := key + ": "
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			return strings.TrimPrefix(l, prefix)
		}
	}
	return ""
}

// tlhHasFieldWithValue checks that a parsed block contains "key: value".
func tlhHasFieldWithValue(lines []string, key, value string) bool {
	prefix := key + ": "
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			if value == "" {
				return true
			}
			return l == prefix+value
		}
	}
	return false
}

// tlhCompleteTask updates a task's status to "completed" and removes state.json.
func tlhCompleteTask(t *testing.T, projectRoot, taskKey string) {
	t.Helper()
	indexPath := filepath.Join(projectRoot, "docs", "features", "task-lifecycle-hardening", "tasks", "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index.json: %v", err)
	}
	var idx tlhIndexFixture
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("failed to unmarshal index.json: %v", err)
	}
	task := idx.Tasks[taskKey]
	task.Status = "completed"
	idx.Tasks[taskKey] = task
	updated, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal updated index.json: %v", err)
	}
	if err := os.WriteFile(indexPath, updated, 0644); err != nil {
		t.Fatalf("failed to write updated index.json: %v", err)
	}
	// Remove state.json so next claim picks a new task
	_ = os.Remove(filepath.Join(projectRoot, "docs", "features", "task-lifecycle-hardening", "tasks", "process", "state.json"))
}

// tlhReadIndex reads and returns the current index.json from the project root.
func tlhReadIndex(t *testing.T, projectRoot string) tlhIndexFixture {
	t.Helper()
	indexPath := filepath.Join(projectRoot, "docs", "features", "task-lifecycle-hardening", "tasks", "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index.json: %v", err)
	}
	var idx tlhIndexFixture
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("failed to unmarshal index.json: %v", err)
	}
	return idx
}

// ==============================================================================
// Self-block tests (SourceTaskID == selfID)
// ==============================================================================

// Traceability: TC-001 -> Task 1 / AC-1 ("checkDependenciesMet returns false when an active fix-task has SourceTaskID == selfID")
func TestTC_001_ActiveFixTaskWithSourceTaskIDEqSelfBlocksClaim(t *testing.T) {
	// Setup: pending task "3" (no deps), pending fix-task "fix-1" (SourceTaskID="3", Type="fix")
	projectRoot := tlhSetupFeatureFixture(t, map[string]tlhTaskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "pending", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "pending", File: "fix-1.md", Record: "",
			Type: "fix", SourceTaskID: "3",
		},
	})

	// Execute: forge task claim
	out, exitCode := tlhForgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := tlhParseBlock(t, out)

	// Verify: fix-1 should be claimed (task 3 is blocked by self-targeting fix)
	assert.True(t, tlhHasFieldWithValue(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, tlhHasFieldWithValue(lines, "TASK_ID", "fix-1"),
		"expected TASK_ID: fix-1 (task 3 blocked by self-targeting fix), got: %v", lines)
}

// Traceability: TC-002 -> Task 1 / AC-1 (extended -- active includes in_progress status)
func TestTC_002_InProgressFixTaskWithSourceTaskIDEqSelfBlocksClaim(t *testing.T) {
	// Setup: pending task "3" (no deps), in_progress fix-task "fix-1" (SourceTaskID="3", Type="fix")
	projectRoot := tlhSetupFeatureFixture(t, map[string]tlhTaskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "pending", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "in_progress", File: "fix-1.md", Record: "",
			Type: "fix", SourceTaskID: "3",
		},
	})

	// Execute: forge task claim
	out, exitCode := tlhForgeClaim(t, projectRoot)
	// No pending tasks are claimable (task 3 blocked, fix-1 in_progress)
	assert.NotEqual(t, 0, exitCode,
		"forge task claim should fail when no eligible tasks, output: %s", out)
}

// Traceability: TC-003 -> Task 1 / AC-2 ("checkDependenciesMet returns true when fix-task targeting self is completed")
func TestTC_003_CompletedFixTaskTargetingSelfDoesNotBlock(t *testing.T) {
	// Setup: pending task "3" (no deps), completed fix-task "fix-1" (SourceTaskID="3", Type="fix")
	projectRoot := tlhSetupFeatureFixture(t, map[string]tlhTaskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "pending", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "completed", File: "fix-1.md", Record: "",
			Type: "fix", SourceTaskID: "3",
		},
	})

	// Execute: forge task claim
	out, exitCode := tlhForgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := tlhParseBlock(t, out)

	// Verify: task 3 should be claimed (completed fix does not block)
	assert.True(t, tlhHasFieldWithValue(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, tlhHasFieldWithValue(lines, "TASK_ID", "3"),
		"expected TASK_ID: 3 (completed fix does not block), got: %v", lines)
}

// Traceability: TC-004 -> Task 1 / AC-1 + Proposal flowchart (node A4)
func TestTC_004_SelfBlockTakesPrecedenceOverMetRegularDependencies(t *testing.T) {
	// Setup: completed task "2", pending task "3" (depends on "2"),
	// pending fix-task "fix-1" (SourceTaskID="3", Type="fix")
	projectRoot := tlhSetupFeatureFixture(t, map[string]tlhTaskEntry{
		"2": {
			ID: "2", Title: "Task 2", Priority: "P0",
			Status: "completed", File: "task-2.md", Record: "",
		},
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "pending", File: "task-3.md", Record: "",
			Dependencies: []string{"2"},
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "pending", File: "fix-1.md", Record: "",
			Type: "fix", SourceTaskID: "3",
		},
	})

	// Execute: forge task claim
	out, exitCode := tlhForgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := tlhParseBlock(t, out)

	// Verify: fix-1 should be claimed (task 3 blocked despite regular deps being met)
	assert.True(t, tlhHasFieldWithValue(lines, "TASK_ID", "fix-1"),
		"expected TASK_ID: fix-1 (self-block takes precedence), got: %v", lines)
}

// Traceability: TC-005 -> Task 1 / AC-3 ("Existing behavior unchanged for tasks without active fix-tasks targeting them")
func TestTC_005_FixTaskTargetingOtherTaskDoesNotCauseSelfBlock(t *testing.T) {
	// Setup: completed task "2", pending task "3" (no deps),
	// pending fix-task "fix-1" (SourceTaskID="2", Type="fix" -- targets task 2, not task 3)
	projectRoot := tlhSetupFeatureFixture(t, map[string]tlhTaskEntry{
		"2": {
			ID: "2", Title: "Task 2", Priority: "P0",
			Status: "completed", File: "task-2.md", Record: "",
		},
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "pending", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "pending", File: "fix-1.md", Record: "",
			Type: "fix", SourceTaskID: "2",
		},
	})

	// Execute: forge task claim
	out, exitCode := tlhForgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := tlhParseBlock(t, out)

	// Verify: either fix-1 or task 3 can be claimed (both eligible, neither blocked)
	taskID := tlhGetFieldValue(lines, "TASK_ID")
	assert.True(t, taskID == "3" || taskID == "fix-1",
		"expected TASK_ID to be '3' or 'fix-1' (fix targeting other task does not self-block), got: %s, lines: %v", taskID, lines)
}

// Traceability: TC-006 -> Task 1 / AC-1 (extended -- multiple fix-tasks)
func TestTC_006_MultipleFixTasksTargetingSelfMustAllComplete(t *testing.T) {
	// Setup: pending task "3" (no deps), completed fix-task "fix-1" (SourceTaskID="3"),
	// pending fix-task "fix-2" (SourceTaskID="3")
	projectRoot := tlhSetupFeatureFixture(t, map[string]tlhTaskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "pending", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "completed", File: "fix-1.md", Record: "",
			Type: "fix", SourceTaskID: "3",
		},
		"fix-2": {
			ID: "fix-2", Title: "Fix Task 2", Priority: "P0",
			Status: "pending", File: "fix-2.md", Record: "",
			Type: "fix", SourceTaskID: "3",
		},
	})

	// Execute: forge task claim
	out, exitCode := tlhForgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := tlhParseBlock(t, out)

	// Verify: fix-2 should be claimed (task 3 still blocked by pending fix-2)
	assert.True(t, tlhHasFieldWithValue(lines, "TASK_ID", "fix-2"),
		"expected TASK_ID: fix-2 (task 3 blocked while fix-2 pending), got: %v", lines)
}

// ==============================================================================
// Lazy unblock scan tests
// ==============================================================================

// Traceability: TC-007 -> Task 2 / AC-4 + Proposal SC-1 ("Blocked task auto-transitions to pending when checkDependenciesMet returns true")
func TestTC_007_BlockedTaskAutoUnblockedWhenDependenciesMet(t *testing.T) {
	// Setup: completed task "1", blocked task "2" (depends on "1")
	projectRoot := tlhSetupFeatureFixture(t, map[string]tlhTaskEntry{
		"1": {
			ID: "1", Title: "Task 1", Priority: "P0",
			Status: "completed", File: "task-1.md", Record: "",
		},
		"2": {
			ID: "2", Title: "Task 2", Priority: "P0",
			Status: "blocked", File: "task-2.md", Record: "",
			Dependencies: []string{"1"},
		},
	})

	// Execute: forge task claim
	out, exitCode := tlhForgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := tlhParseBlock(t, out)

	// Verify: task 2 should be claimed (auto-unblocked from blocked -> pending -> in_progress)
	assert.True(t, tlhHasFieldWithValue(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, tlhHasFieldWithValue(lines, "TASK_ID", "2"),
		"expected TASK_ID: 2 (auto-unblocked and claimed), got: %v", lines)
}

// Traceability: TC-008 -> Task 2 / AC-4 (negative case)
func TestTC_008_BlockedTaskStaysBlockedWhenDependenciesNotMet(t *testing.T) {
	// Setup: pending task "1", blocked task "2" (depends on "1")
	projectRoot := tlhSetupFeatureFixture(t, map[string]tlhTaskEntry{
		"1": {
			ID: "1", Title: "Task 1", Priority: "P0",
			Status: "pending", File: "task-1.md", Record: "",
		},
		"2": {
			ID: "2", Title: "Task 2", Priority: "P0",
			Status: "blocked", File: "task-2.md", Record: "",
			Dependencies: []string{"1"},
		},
	})

	// Execute: forge task claim
	out, exitCode := tlhForgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := tlhParseBlock(t, out)

	// Verify: task 1 should be claimed (task 2 stays blocked because dep is pending)
	assert.True(t, tlhHasFieldWithValue(lines, "TASK_ID", "1"),
		"expected TASK_ID: 1 (task 2 stays blocked), got: %v", lines)

	// Verify: task 2 should still be blocked in index
	idx := tlhReadIndex(t, projectRoot)
	assert.Equal(t, "blocked", idx.Tasks["2"].Status,
		"task 2 should remain blocked after claim")
}

// Traceability: TC-009 -> Task 2 / AC-5 ("Auto-unblocked tasks are logged")
func TestTC_009_AutoUnblockLoggedToStdout(t *testing.T) {
	// Setup: completed task "1", blocked task "2" (depends on "1")
	projectRoot := tlhSetupFeatureFixture(t, map[string]tlhTaskEntry{
		"1": {
			ID: "1", Title: "Task 1", Priority: "P0",
			Status: "completed", File: "task-1.md", Record: "",
		},
		"2": {
			ID: "2", Title: "Task 2", Priority: "P0",
			Status: "blocked", File: "task-2.md", Record: "",
			Dependencies: []string{"1"},
		},
	})

	// Execute: forge task claim
	out, exitCode := tlhForgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)

	// Verify: stdout contains auto-unblock log message
	assert.True(t, strings.Contains(out, "Auto-unblocked task 2"),
		"expected auto-unblock log for task 2, got: %s", out)
}

// Traceability: TC-010 -> Task 2 / AC-4 (extended -- multiple tasks) + AC-6 ("scan runs before hasPending check")
func TestTC_010_MultipleBlockedTasksUnblockedSimultaneously(t *testing.T) {
	// Setup: completed task "1", blocked task "2" (depends on "1", P1),
	// blocked task "3" (depends on "1", P0)
	projectRoot := tlhSetupFeatureFixture(t, map[string]tlhTaskEntry{
		"1": {
			ID: "1", Title: "Task 1", Priority: "P0",
			Status: "completed", File: "task-1.md", Record: "",
		},
		"2": {
			ID: "2", Title: "Task 2", Priority: "P1",
			Status: "blocked", File: "task-2.md", Record: "",
			Dependencies: []string{"1"},
		},
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "blocked", File: "task-3.md", Record: "",
			Dependencies: []string{"1"},
		},
	})

	// Execute: forge task claim
	out, exitCode := tlhForgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := tlhParseBlock(t, out)

	// Verify: task 3 (P0) should be claimed, task 2 auto-unblocked to pending
	assert.True(t, tlhHasFieldWithValue(lines, "TASK_ID", "3"),
		"expected TASK_ID: 3 (P0 beats P1), got: %v", lines)

	// Verify: task 2 should be auto-unblocked to pending
	idx := tlhReadIndex(t, projectRoot)
	assert.Equal(t, "pending", idx.Tasks["2"].Status,
		"task 2 should be auto-unblocked to pending")
}

// Traceability: TC-011 -> Task 2 / AC-7 + Proposal SC-2 ("Fix-task in progress targeting the task keeps it blocked")
func TestTC_011_BlockedTaskWithActiveFixTargetingItStaysBlocked(t *testing.T) {
	// Setup: completed task "1", blocked task "2" (depends on "1"),
	// pending fix-task "fix-1" (SourceTaskID="2", Type="fix")
	projectRoot := tlhSetupFeatureFixture(t, map[string]tlhTaskEntry{
		"1": {
			ID: "1", Title: "Task 1", Priority: "P0",
			Status: "completed", File: "task-1.md", Record: "",
		},
		"2": {
			ID: "2", Title: "Task 2", Priority: "P0",
			Status: "blocked", File: "task-2.md", Record: "",
			Dependencies: []string{"1"},
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "pending", File: "fix-1.md", Record: "",
			Type: "fix", SourceTaskID: "2",
		},
	})

	// Execute: forge task claim
	out, exitCode := tlhForgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := tlhParseBlock(t, out)

	// Verify: fix-1 should be claimed (task 2 stays blocked by active fix targeting it)
	assert.True(t, tlhHasFieldWithValue(lines, "TASK_ID", "fix-1"),
		"expected TASK_ID: fix-1 (task 2 blocked by fix), got: %v", lines)

	// Verify: task 2 should remain blocked
	idx := tlhReadIndex(t, projectRoot)
	assert.Equal(t, "blocked", idx.Tasks["2"].Status,
		"task 2 should remain blocked due to active fix targeting it")
}

// ==============================================================================
// Block-source lifecycle tests
// ==============================================================================

// Traceability: TC-012 -> Proposal SC-3 + Task 2 / AC-4
func TestTC_012_FixCompletedAutoUnblocksBlockedSourceTask(t *testing.T) {
	// Setup: blocked source task "1" (no deps), completed fix-task "fix-1" (SourceTaskID="1", Type="fix")
	projectRoot := tlhSetupFeatureFixture(t, map[string]tlhTaskEntry{
		"1": {
			ID: "1", Title: "Source task", Priority: "P0",
			Status: "blocked", File: "task-1.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "completed", File: "fix-1.md", Record: "",
			Type: "fix", SourceTaskID: "1",
		},
	})

	// Execute: forge task claim
	out, exitCode := tlhForgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := tlhParseBlock(t, out)

	// Verify: source task "1" should be auto-unblocked and claimed
	assert.True(t, tlhHasFieldWithValue(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, tlhHasFieldWithValue(lines, "TASK_ID", "1"),
		"expected TASK_ID: 1 (source auto-unblocked after fix completed), got: %v", lines)
}

// Traceability: TC-013 -> Proposal SC-3 (negative case -- fix still active)
func TestTC_013_SourceStaysBlockedWhenFixIsStillInProgress(t *testing.T) {
	// Setup: blocked source task "1" (no deps), in_progress fix-task "fix-1" (SourceTaskID="1", Type="fix")
	projectRoot := tlhSetupFeatureFixture(t, map[string]tlhTaskEntry{
		"1": {
			ID: "1", Title: "Source task", Priority: "P0",
			Status: "blocked", File: "task-1.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "in_progress", File: "fix-1.md", Record: "",
			Type: "fix", SourceTaskID: "1",
		},
	})

	// Execute: forge task claim
	out, exitCode := tlhForgeClaim(t, projectRoot)

	// Verify: no eligible tasks (fix still active, source stays blocked)
	assert.NotEqual(t, 0, exitCode,
		"forge task claim should fail (no eligible tasks), output: %s", out)

	// Verify: source task remains blocked in index
	idx := tlhReadIndex(t, projectRoot)
	assert.Equal(t, "blocked", idx.Tasks["1"].Status,
		"source task should remain blocked while fix is in_progress")
}

// ==============================================================================
// Auto-downgrade unblock tests
// ==============================================================================

// Traceability: TC-014 -> Proposal SC-4 ("Auto-downgrade scenario: task blocked -> dep completed -> claim auto-unblocks")
func TestTC_014_AutoDowngradedTaskAutoUnblockedWhenDepCompletes(t *testing.T) {
	// Setup: completed task "1", blocked task "2" (depends on "1", BlockedReason="auto-downgrade: testsFailed=2")
	projectRoot := tlhSetupFeatureFixture(t, map[string]tlhTaskEntry{
		"1": {
			ID: "1", Title: "Task 1", Priority: "P0",
			Status: "completed", File: "task-1.md", Record: "",
		},
		"2": {
			ID: "2", Title: "Task 2", Priority: "P0",
			Status: "blocked", File: "task-2.md", Record: "",
			Dependencies: []string{"1"},
			BlockedReason: "auto-downgrade: testsFailed=2",
		},
	})

	// Execute: forge task claim
	out, exitCode := tlhForgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := tlhParseBlock(t, out)

	// Verify: auto-downgraded task should be auto-unblocked and claimed
	assert.True(t, tlhHasFieldWithValue(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, tlhHasFieldWithValue(lines, "TASK_ID", "2"),
		"expected TASK_ID: 2 (auto-downgraded task auto-unblocked), got: %v", lines)
}
