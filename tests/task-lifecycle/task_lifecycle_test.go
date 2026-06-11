//go:build cli_functional

package tasklifecycle

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Task lifecycle + fix-task claim priority tests — Journey: task-lifecycle
// Tests verify self-block (SourceTaskID==selfID), lazy unblock scan,
// block-source lifecycle, auto-downgrade unblock, and fix-task claim priority
// scenarios during forge task claim.
// ==============================================================================

// taskEntry represents a single task in the index.json fixture.
type taskEntry struct {
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

// indexFixture represents the top-level index.json structure.
type indexFixture struct {
	Feature      string               `json:"feature"`
	StatusEnum   []string             `json:"statusEnum"`
	PriorityEnum []string             `json:"priorityEnum"`
	Tasks        map[string]taskEntry `json:"tasks"`
}

// setupFeatureFixture creates a temp project root with the given tasks in
// index.json under docs/features/<featureSlug>/tasks/. Returns the temp dir
// path to use as CLAUDE_PROJECT_DIR.
func setupFeatureFixture(t *testing.T, featureSlug string, tasks map[string]taskEntry) string {
	t.Helper()
	dir := t.TempDir()

	tasksDir := filepath.Join(dir, "docs", "features", featureSlug, "tasks")
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

	idx := indexFixture{
		Feature:      featureSlug,
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

// forgeClaim runs "forge task claim" with CLAUDE_PROJECT_DIR set.
// Returns combined output and exit code. Does NOT fatalf on failure.
func forgeClaim(t *testing.T, projectRoot string) (string, int) {
	t.Helper()
	cmd := exec.Command(testkit.ForgeBinary, "task", "claim")
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

// completeTask updates a task's status to "completed" and removes state.json.
func completeTask(t *testing.T, projectRoot, featureSlug, taskKey string) {
	t.Helper()
	indexPath := filepath.Join(projectRoot, "docs", "features", featureSlug, "tasks", "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index.json: %v", err)
	}
	var idx indexFixture
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
	_ = os.Remove(filepath.Join(projectRoot, "docs", "features", featureSlug, "tasks", "process", "state.json"))
}

// readIndex reads and returns the current index.json from the project root.
func readIndex(t *testing.T, projectRoot, featureSlug string) indexFixture {
	t.Helper()
	indexPath := filepath.Join(projectRoot, "docs", "features", featureSlug, "tasks", "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index.json: %v", err)
	}
	var idx indexFixture
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("failed to unmarshal index.json: %v", err)
	}
	return idx
}

// parseClaimOutput strips any auto-unblock log lines and extracts the block between --- separators.
func parseClaimOutput(t *testing.T, raw string) []string {
	t.Helper()
	// Strip any lines before the first "---" separator (e.g. auto-unblock logs)
	idx := strings.Index(raw, "\n---")
	if idx >= 0 && !strings.HasPrefix(strings.TrimSpace(raw), "---") {
		raw = raw[idx+1:] // keep the "---" line
	}
	return testkit.ParseBlock(t, raw)
}

// ==============================================================================
// Self-block tests (SourceTaskID == selfID)
// ==============================================================================

// Traceability: TC-001 -> Task 1 / AC-1 ("checkDependenciesMet returns false when an active fix-task has SourceTaskID == selfID")
func TestTC_001_ActiveFixTaskWithSourceTaskIDEqSelfBlocksClaim(t *testing.T) {
	// Setup: pending task "3" (no deps), pending fix-task "fix-1" (SourceTaskID="3", Type="coding.fix")
	projectRoot := setupFeatureFixture(t, "task-lifecycle-hardening", map[string]taskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "pending", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "pending", File: "fix-1.md", Record: "",
			Type: "coding.fix", SourceTaskID: "3",
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: fix-1 should be claimed (task 3 is blocked by self-targeting fix)
	assert.True(t, testkit.HasField(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, testkit.HasField(lines, "TASK_ID", "fix-1"),
		"expected TASK_ID: fix-1 (task 3 blocked by self-targeting fix), got: %v", lines)
}

// Traceability: TC-002 -> Task 1 / AC-1 (extended -- active includes in_progress status)
func TestTC_002_InProgressFixTaskWithSourceTaskIDEqSelfBlocksClaim(t *testing.T) {
	// Setup: pending task "3" (no deps), in_progress fix-task "fix-1" (SourceTaskID="3", Type="coding.fix")
	projectRoot := setupFeatureFixture(t, "task-lifecycle-hardening", map[string]taskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "pending", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "in_progress", File: "fix-1.md", Record: "",
			Type: "coding.fix", SourceTaskID: "3",
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)
	// No pending tasks are claimable (task 3 blocked, fix-1 in_progress)
	assert.NotEqual(t, 0, exitCode,
		"forge task claim should fail when no eligible tasks, output: %s", out)
}

// Traceability: TC-003 -> Task 1 / AC-2 ("checkDependenciesMet returns true when fix-task targeting self is completed")
func TestTC_003_CompletedFixTaskTargetingSelfDoesNotBlock(t *testing.T) {
	// Setup: pending task "3" (no deps), completed fix-task "fix-1" (SourceTaskID="3", Type="coding.fix")
	projectRoot := setupFeatureFixture(t, "task-lifecycle-hardening", map[string]taskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "pending", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "completed", File: "fix-1.md", Record: "",
			Type: "coding.fix", SourceTaskID: "3",
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: task 3 should be claimed (completed fix does not block)
	assert.True(t, testkit.HasField(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, testkit.HasField(lines, "TASK_ID", "3"),
		"expected TASK_ID: 3 (completed fix does not block), got: %v", lines)
}

// Traceability: TC-004 -> Task 1 / AC-1 + Proposal flowchart (node A4)
func TestTC_004_SelfBlockTakesPrecedenceOverMetRegularDependencies(t *testing.T) {
	// Setup: completed task "2", pending task "3" (depends on "2"),
	// pending fix-task "fix-1" (SourceTaskID="3", Type="coding.fix")
	projectRoot := setupFeatureFixture(t, "task-lifecycle-hardening", map[string]taskEntry{
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
			Type: "coding.fix", SourceTaskID: "3",
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: fix-1 should be claimed (task 3 blocked despite regular deps being met)
	assert.True(t, testkit.HasField(lines, "TASK_ID", "fix-1"),
		"expected TASK_ID: fix-1 (self-block takes precedence), got: %v", lines)
}

// Traceability: TC-005 -> Task 1 / AC-3 ("Existing behavior unchanged for tasks without active fix-tasks targeting them")
func TestTC_005_FixTaskTargetingOtherTaskDoesNotCauseSelfBlock(t *testing.T) {
	// Setup: completed task "2", pending task "3" (no deps),
	// pending fix-task "fix-1" (SourceTaskID="2", Type="coding.fix" -- targets task 2, not task 3)
	projectRoot := setupFeatureFixture(t, "task-lifecycle-hardening", map[string]taskEntry{
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
			Type: "coding.fix", SourceTaskID: "2",
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: either fix-1 or task 3 can be claimed (both eligible, neither blocked)
	taskID := testkit.FieldValue(lines, "TASK_ID")
	assert.True(t, taskID == "3" || taskID == "fix-1",
		"expected TASK_ID to be '3' or 'fix-1' (fix targeting other task does not self-block), got: %s, lines: %v", taskID, lines)
}

// Traceability: TC-006 -> Task 1 / AC-1 (extended -- multiple fix-tasks)
func TestTC_006_MultipleFixTasksTargetingSelfMustAllComplete(t *testing.T) {
	// Setup: pending task "3" (no deps), completed fix-task "fix-1" (SourceTaskID="3"),
	// pending fix-task "fix-2" (SourceTaskID="3")
	projectRoot := setupFeatureFixture(t, "task-lifecycle-hardening", map[string]taskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "pending", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "completed", File: "fix-1.md", Record: "",
			Type: "coding.fix", SourceTaskID: "3",
		},
		"fix-2": {
			ID: "fix-2", Title: "Fix Task 2", Priority: "P0",
			Status: "pending", File: "fix-2.md", Record: "",
			Type: "coding.fix", SourceTaskID: "3",
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: fix-2 should be claimed (task 3 still blocked by pending fix-2)
	assert.True(t, testkit.HasField(lines, "TASK_ID", "fix-2"),
		"expected TASK_ID: fix-2 (task 3 blocked while fix-2 pending), got: %v", lines)
}

// ==============================================================================
// Lazy unblock scan tests
// ==============================================================================

// Traceability: TC-007 -> Task 2 / AC-4 + Proposal SC-1 ("Blocked task auto-transitions to pending when checkDependenciesMet returns true")
func TestTC_007_BlockedTaskAutoUnblockedWhenDependenciesMet(t *testing.T) {
	// Setup: completed task "1", blocked task "2" (depends on "1")
	projectRoot := setupFeatureFixture(t, "task-lifecycle-hardening", map[string]taskEntry{
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
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: task 2 should be claimed (auto-unblocked from blocked -> pending -> in_progress)
	assert.True(t, testkit.HasField(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, testkit.HasField(lines, "TASK_ID", "2"),
		"expected TASK_ID: 2 (auto-unblocked and claimed), got: %v", lines)
}

// Traceability: TC-008 -> Task 2 / AC-4 (negative case)
func TestTC_008_BlockedTaskStaysBlockedWhenDependenciesNotMet(t *testing.T) {
	// Setup: pending task "1", blocked task "2" (depends on "1")
	projectRoot := setupFeatureFixture(t, "task-lifecycle-hardening", map[string]taskEntry{
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
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: task 1 should be claimed (task 2 stays blocked because dep is pending)
	assert.True(t, testkit.HasField(lines, "TASK_ID", "1"),
		"expected TASK_ID: 1 (task 2 stays blocked), got: %v", lines)

	// Verify: task 2 should still be blocked in index
	idx := readIndex(t, projectRoot, "task-lifecycle-hardening")
	assert.Equal(t, "blocked", idx.Tasks["2"].Status,
		"task 2 should remain blocked after claim")
}

// Traceability: TC-009 -> Task 2 / AC-5 ("Auto-unblocked tasks are logged")
func TestTC_009_AutoUnblockLoggedToStdout(t *testing.T) {
	// Setup: completed task "1", blocked task "2" (depends on "1")
	projectRoot := setupFeatureFixture(t, "task-lifecycle-hardening", map[string]taskEntry{
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
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)

	// Verify: stdout contains auto-unblock log message
	assert.True(t, strings.Contains(out, "Auto-unblocked task 2"),
		"expected auto-unblock log for task 2, got: %s", out)
}

// Traceability: TC-010 -> Task 2 / AC-4 (extended -- multiple tasks) + AC-6 ("scan runs before hasPending check")
func TestTC_010_MultipleBlockedTasksUnblockedSimultaneously(t *testing.T) {
	// Setup: completed task "1", blocked task "2" (depends on "1", P1),
	// blocked task "3" (depends on "1", P0)
	projectRoot := setupFeatureFixture(t, "task-lifecycle-hardening", map[string]taskEntry{
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
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: task 3 (P0) should be claimed, task 2 auto-unblocked to pending
	assert.True(t, testkit.HasField(lines, "TASK_ID", "3"),
		"expected TASK_ID: 3 (P0 beats P1), got: %v", lines)

	// Verify: task 2 should be auto-unblocked to pending
	idx := readIndex(t, projectRoot, "task-lifecycle-hardening")
	assert.Equal(t, "pending", idx.Tasks["2"].Status,
		"task 2 should be auto-unblocked to pending")
}

// Traceability: TC-011 -> Task 2 / AC-7 + Proposal SC-2 ("Fix-task in progress targeting the task keeps it blocked")
func TestTC_011_BlockedTaskWithActiveFixTargetingItStaysBlocked(t *testing.T) {
	// Setup: completed task "1", blocked task "2" (depends on "1"),
	// pending fix-task "fix-1" (SourceTaskID="2", Type="coding.fix")
	projectRoot := setupFeatureFixture(t, "task-lifecycle-hardening", map[string]taskEntry{
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
			Type: "coding.fix", SourceTaskID: "2",
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: fix-1 should be claimed (task 2 stays blocked by active fix targeting it)
	assert.True(t, testkit.HasField(lines, "TASK_ID", "fix-1"),
		"expected TASK_ID: fix-1 (task 2 blocked by fix), got: %v", lines)

	// Verify: task 2 should remain blocked
	idx := readIndex(t, projectRoot, "task-lifecycle-hardening")
	assert.Equal(t, "blocked", idx.Tasks["2"].Status,
		"task 2 should remain blocked due to active fix targeting it")
}

// ==============================================================================
// Block-source lifecycle tests
// ==============================================================================

// Traceability: TC-012 -> Proposal SC-3 + Task 2 / AC-4
func TestTC_012_FixCompletedAutoUnblocksBlockedSourceTask(t *testing.T) {
	// Setup: blocked source task "1" (no deps), completed fix-task "fix-1" (SourceTaskID="1", Type="coding.fix")
	projectRoot := setupFeatureFixture(t, "task-lifecycle-hardening", map[string]taskEntry{
		"1": {
			ID: "1", Title: "Source task", Priority: "P0",
			Status: "blocked", File: "task-1.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "completed", File: "fix-1.md", Record: "",
			Type: "coding.fix", SourceTaskID: "1",
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: source task "1" should be auto-unblocked and claimed
	assert.True(t, testkit.HasField(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, testkit.HasField(lines, "TASK_ID", "1"),
		"expected TASK_ID: 1 (source auto-unblocked after fix completed), got: %v", lines)
}

// Traceability: TC-013 -> Proposal SC-3 (negative case -- fix still active)
func TestTC_013_SourceStaysBlockedWhenFixIsStillInProgress(t *testing.T) {
	// Setup: blocked source task "1" (no deps), in_progress fix-task "fix-1" (SourceTaskID="1", Type="coding.fix")
	projectRoot := setupFeatureFixture(t, "task-lifecycle-hardening", map[string]taskEntry{
		"1": {
			ID: "1", Title: "Source task", Priority: "P0",
			Status: "blocked", File: "task-1.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "in_progress", File: "fix-1.md", Record: "",
			Type: "coding.fix", SourceTaskID: "1",
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)

	// Verify: no eligible tasks (fix still active, source stays blocked)
	assert.NotEqual(t, 0, exitCode,
		"forge task claim should fail (no eligible tasks), output: %s", out)

	// Verify: source task remains blocked in index
	idx := readIndex(t, projectRoot, "task-lifecycle-hardening")
	assert.Equal(t, "blocked", idx.Tasks["1"].Status,
		"source task should remain blocked while fix is in_progress")
}

// ==============================================================================
// Auto-downgrade unblock tests
// ==============================================================================

// Traceability: TC-014 -> Proposal SC-4 ("Auto-downgrade scenario: task blocked -> dep completed -> claim auto-unblocks")
func TestTC_014_AutoDowngradedTaskAutoUnblockedWhenDepCompletes(t *testing.T) {
	// Setup: completed task "1", blocked task "2" (depends on "1", BlockedReason="auto-downgrade: testsFailed=2")
	projectRoot := setupFeatureFixture(t, "task-lifecycle-hardening", map[string]taskEntry{
		"1": {
			ID: "1", Title: "Task 1", Priority: "P0",
			Status: "completed", File: "task-1.md", Record: "",
		},
		"2": {
			ID: "2", Title: "Task 2", Priority: "P0",
			Status: "blocked", File: "task-2.md", Record: "",
			Dependencies:  []string{"1"},
			BlockedReason: "auto-downgrade: testsFailed=2",
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: auto-downgraded task should be auto-unblocked and claimed
	assert.True(t, testkit.HasField(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, testkit.HasField(lines, "TASK_ID", "2"),
		"expected TASK_ID: 2 (auto-downgraded task auto-unblocked), got: %v", lines)
}

// ==============================================================================
// Fix-task claim priority tests (from fix-task-claim-priority feature)
// ==============================================================================

// Traceability: TC-015 -> Proposal Key Scenario "Primary" + Success Criterion 1
func TestTC_015_PendingFixTaskBlocksDependentBusinessTask(t *testing.T) {
	// Setup: task 3 completed, fix-1 pending (sourceTaskID "3"), task 4 pending (depends on task 3)
	projectRoot := setupFeatureFixture(t, "fix-task-claim-priority", map[string]taskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "completed", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "pending", File: "fix-1.md", Record: "",
			Type: "coding.fix", SourceTaskID: "3",
		},
		"4": {
			ID: "4", Title: "Task 4", Priority: "P0",
			Status: "pending", File: "task-4.md", Record: "",
			Dependencies: []string{"3"},
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: fix-1 should be claimed, not task 4
	assert.True(t, testkit.HasField(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, testkit.HasField(lines, "TASK_ID", "fix-1"),
		"expected TASK_ID: fix-1 (fix task should be claimed before blocked business task), got: %v", lines)
}

// Traceability: TC-016 -> Proposal Key Scenario "Fix completed" + Success Criterion 2
func TestTC_016_CompletedFixTaskAllowsDependentBusinessTask(t *testing.T) {
	// Setup: task 3 completed, fix-1 completed (sourceTaskID "3"), task 4 pending (depends on task 3)
	projectRoot := setupFeatureFixture(t, "fix-task-claim-priority", map[string]taskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "completed", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "completed", File: "fix-1.md", Record: "",
			Type: "coding.fix", SourceTaskID: "3",
		},
		"4": {
			ID: "4", Title: "Task 4", Priority: "P0",
			Status: "pending", File: "task-4.md", Record: "",
			Dependencies: []string{"3"},
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: task 4 should be claimable since fix-1 is completed
	assert.True(t, testkit.HasField(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, testkit.HasField(lines, "TASK_ID", "4"),
		"expected TASK_ID: 4 (business task should be claimable when fix is completed), got: %v", lines)
}

// Traceability: TC-017 -> Success Criterion 3
func TestTC_017_FixTaskClaimedBeforeBusinessTaskWhenBothEligible(t *testing.T) {
	// Setup: task 3 completed, fix-1 pending (sourceTaskID "3", depends on task 3),
	// task 4 pending (depends on task 3). Both fix-1 and task 4 are P0 with met deps.
	projectRoot := setupFeatureFixture(t, "fix-task-claim-priority", map[string]taskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "completed", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "pending", File: "fix-1.md", Record: "",
			Dependencies: []string{"3"},
			Type:         "coding.fix", SourceTaskID: "3",
		},
		"4": {
			ID: "4", Title: "Task 4", Priority: "P0",
			Status: "pending", File: "task-4.md", Record: "",
			Dependencies: []string{"3"},
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: fix-1 should be returned because it has met dependencies
	assert.True(t, testkit.HasField(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, testkit.HasField(lines, "TASK_ID", "fix-1"),
		"expected TASK_ID: fix-1 (fix task should be claimed before business task), got: %v", lines)
}

// Traceability: TC-018 -> Proposal Key Scenario "Fix chain" + Success Criterion 5
func TestTC_018_FixChainBlocksDependentTaskUntilAllFixTasksComplete(t *testing.T) {
	// Phase 1: task 3 completed, fix-1 completed (sourceTaskID "3"),
	// fix-2 pending (sourceTaskID "3"), task 4 pending (depends on task 3)
	projectRoot := setupFeatureFixture(t, "fix-task-claim-priority", map[string]taskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "completed", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "completed", File: "fix-1.md", Record: "",
			Type: "coding.fix", SourceTaskID: "3",
		},
		"fix-2": {
			ID: "fix-2", Title: "Fix Task 2", Priority: "P0",
			Status: "pending", File: "fix-2.md", Record: "",
			Type: "coding.fix", SourceTaskID: "3",
		},
		"4": {
			ID: "4", Title: "Task 4", Priority: "P0",
			Status: "pending", File: "task-4.md", Record: "",
			Dependencies: []string{"3"},
		},
	})

	// Phase 1: forge task claim should pick fix-2 (task 4 blocked by pending fix)
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "Phase 1: forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)
	assert.True(t, testkit.HasField(lines, "TASK_ID", "fix-2"),
		"Phase 1: expected TASK_ID: fix-2 (pending fix blocks business task), got: %v", lines)

	// Phase 2: Complete fix-2 and claim again
	completeTask(t, projectRoot, "fix-task-claim-priority", "fix-2")
	out, exitCode = forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "Phase 2: forge task claim should succeed, output: %s", out)
	lines = parseClaimOutput(t, out)

	// Verify: task 4 should now be eligible
	assert.True(t, testkit.HasField(lines, "TASK_ID", "4"),
		"Phase 2: expected TASK_ID: 4 (all fix tasks completed), got: %v", lines)
}

// Traceability: TC-019 -> Proposal Key Scenario "Fix for unrelated task" + Success Criterion 6
func TestTC_019_UnrelatedFixTaskDoesNotBlockTaskWithDifferentDependency(t *testing.T) {
	// Setup: task 2 completed, task 3 completed, fix-1 pending (sourceTaskID "2"),
	// task 4 pending (depends on task 3 only)
	projectRoot := setupFeatureFixture(t, "fix-task-claim-priority", map[string]taskEntry{
		"2": {
			ID: "2", Title: "Task 2", Priority: "P0",
			Status: "completed", File: "task-2.md", Record: "",
		},
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "completed", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "pending", File: "fix-1.md", Record: "",
			Type: "coding.fix", SourceTaskID: "2",
		},
		"4": {
			ID: "4", Title: "Task 4", Priority: "P0",
			Status: "pending", File: "task-4.md", Record: "",
			Dependencies: []string{"3"},
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: task 4 must not be blocked by the unrelated fix-1.
	// Either fix-1 or task 4 may be claimed (both have met dependencies),
	// but task 4 MUST be eligible (not blocked).
	taskID := testkit.FieldValue(lines, "TASK_ID")
	assert.True(t, taskID == "4" || taskID == "fix-1",
		"expected TASK_ID to be either '4' or 'fix-1' (unrelated fix should not block task 4), got: %s, lines: %v", taskID, lines)
}

// Traceability: TC-020 -> Proposal Key Scenario "No fix tasks" + Success Criterion 4
func TestTC_020_NoFixTasksPreservesExistingClaimBehavior(t *testing.T) {
	// Setup: task 3 completed, task 4 pending (depends on task 3), no fix tasks
	projectRoot := setupFeatureFixture(t, "fix-task-claim-priority", map[string]taskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "completed", File: "task-3.md", Record: "",
		},
		"4": {
			ID: "4", Title: "Task 4", Priority: "P0",
			Status: "pending", File: "task-4.md", Record: "",
			Dependencies: []string{"3"},
		},
	})

	// Execute: forge task claim
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge task claim should succeed, output: %s", out)
	lines := parseClaimOutput(t, out)

	// Verify: task 4 should be claimed as before (no fix tasks to interfere)
	assert.True(t, testkit.HasField(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, testkit.HasField(lines, "TASK_ID", "4"),
		"expected TASK_ID: 4 (existing claim behavior unchanged), got: %v", lines)
}
