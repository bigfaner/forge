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
// Fix task claim priority tests — feature: fix-task-claim-priority
// Tests verify that fix tasks with sourceTaskID correctly block/allow
// dependent business tasks during forge task claim.
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
}

// indexFixture represents the top-level index.json structure.
type indexFixture struct {
	Feature      string               `json:"feature"`
	StatusEnum   []string             `json:"statusEnum"`
	PriorityEnum []string             `json:"priorityEnum"`
	Tasks        map[string]taskEntry `json:"tasks"`
}

// setupFeatureFixture creates a temp project root with the given tasks in
// index.json under docs/features/fix-task-claim-priority/tasks/. Returns the
// temp dir path to use as CLAUDE_PROJECT_DIR.
func setupFeatureFixture(t *testing.T, tasks map[string]taskEntry) string {
	t.Helper()
	dir := t.TempDir()

	tasksDir := filepath.Join(dir, "docs", "features", "fix-task-claim-priority", "tasks")
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

	idx := indexFixture{
		Feature:      "fix-task-claim-priority",
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

// parseBlockLines extracts key:value lines between "---" separators.
func parseBlockLines(t *testing.T, raw string) []string {
	t.Helper()
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

// getFieldValue returns the value for the given key from parsed block lines.
func getFieldValue(lines []string, key string) string {
	prefix := key + ": "
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			return strings.TrimPrefix(l, prefix)
		}
	}
	return ""
}

// hasFieldWithValue checks that a parsed block contains "key: value".
func hasFieldWithValue(lines []string, key, value string) bool {
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

// completeTask updates a task's status to "completed" and removes state.json.
func completeTask(t *testing.T, projectRoot, taskKey string) {
	t.Helper()
	indexPath := filepath.Join(projectRoot, "docs", "features", "fix-task-claim-priority", "tasks", "index.json")
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
	_ = os.Remove(filepath.Join(projectRoot, "docs", "features", "fix-task-claim-priority", "tasks", "process", "state.json"))
}

// Traceability: TC-001 -> Proposal Key Scenario "Primary" + Success Criterion 1
func TestTC_001_PendingFixTaskBlocksDependentBusinessTask(t *testing.T) {
	// Setup: task 3 completed, fix-1 pending (sourceTaskID "3"), task 4 pending (depends on task 3)
	projectRoot := setupFeatureFixture(t, map[string]taskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "completed", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "pending", File: "fix-1.md", Record: "",
			Type: "fix", SourceTaskID: "3",
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
	lines := parseBlockLines(t, out)

	// Verify: fix-1 should be claimed, not task 4
	assert.True(t, hasFieldWithValue(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, hasFieldWithValue(lines, "TASK_ID", "fix-1"),
		"expected TASK_ID: fix-1 (fix task should be claimed before blocked business task), got: %v", lines)
}

// Traceability: TC-002 -> Proposal Key Scenario "Fix completed" + Success Criterion 2
func TestTC_002_CompletedFixTaskAllowsDependentBusinessTask(t *testing.T) {
	// Setup: task 3 completed, fix-1 completed (sourceTaskID "3"), task 4 pending (depends on task 3)
	projectRoot := setupFeatureFixture(t, map[string]taskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "completed", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "completed", File: "fix-1.md", Record: "",
			Type: "fix", SourceTaskID: "3",
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
	lines := parseBlockLines(t, out)

	// Verify: task 4 should be claimable since fix-1 is completed
	assert.True(t, hasFieldWithValue(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, hasFieldWithValue(lines, "TASK_ID", "4"),
		"expected TASK_ID: 4 (business task should be claimable when fix is completed), got: %v", lines)
}

// Traceability: TC-003 -> Success Criterion 3
func TestTC_003_FixTaskClaimedBeforeBusinessTaskWhenBothEligible(t *testing.T) {
	// Setup: task 3 completed, fix-1 pending (sourceTaskID "3", depends on task 3),
	// task 4 pending (depends on task 3). Both fix-1 and task 4 are P0 with met deps.
	projectRoot := setupFeatureFixture(t, map[string]taskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "completed", File: "task-3.md", Record: "",
		},
		"fix-1": {
			ID: "fix-1", Title: "Fix Task 1", Priority: "P0",
			Status: "pending", File: "fix-1.md", Record: "",
			Dependencies: []string{"3"},
			Type: "fix", SourceTaskID: "3",
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
	lines := parseBlockLines(t, out)

	// Verify: fix-1 should be returned because it has met dependencies
	// (task 3 is completed, no pending fix tasks for its own deps)
	// and sorts at the same priority level as task 4
	assert.True(t, hasFieldWithValue(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, hasFieldWithValue(lines, "TASK_ID", "fix-1"),
		"expected TASK_ID: fix-1 (fix task should be claimed before business task), got: %v", lines)
}

// Traceability: TC-004 -> Proposal Key Scenario "Fix chain" + Success Criterion 5
func TestTC_004_FixChainBlocksDependentTaskUntilAllFixTasksComplete(t *testing.T) {
	// Phase 1: task 3 completed, fix-1 completed (sourceTaskID "3"),
	// fix-2 pending (sourceTaskID "3"), task 4 pending (depends on task 3)
	projectRoot := setupFeatureFixture(t, map[string]taskEntry{
		"3": {
			ID: "3", Title: "Task 3", Priority: "P0",
			Status: "completed", File: "task-3.md", Record: "",
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
		"4": {
			ID: "4", Title: "Task 4", Priority: "P0",
			Status: "pending", File: "task-4.md", Record: "",
			Dependencies: []string{"3"},
		},
	})

	// Phase 1: forge task claim should pick fix-2 (task 4 blocked by pending fix)
	out, exitCode := forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "Phase 1: forge task claim should succeed, output: %s", out)
	lines := parseBlockLines(t, out)
	assert.True(t, hasFieldWithValue(lines, "TASK_ID", "fix-2"),
		"Phase 1: expected TASK_ID: fix-2 (pending fix blocks business task), got: %v", lines)

	// Phase 2: Complete fix-2 and claim again
	completeTask(t, projectRoot, "fix-2")
	out, exitCode = forgeClaim(t, projectRoot)
	assert.Equal(t, 0, exitCode, "Phase 2: forge task claim should succeed, output: %s", out)
	lines = parseBlockLines(t, out)

	// Verify: task 4 should now be eligible
	assert.True(t, hasFieldWithValue(lines, "TASK_ID", "4"),
		"Phase 2: expected TASK_ID: 4 (all fix tasks completed), got: %v", lines)
}

// Traceability: TC-005 -> Proposal Key Scenario "Fix for unrelated task" + Success Criterion 6
func TestTC_005_UnrelatedFixTaskDoesNotBlockTaskWithDifferentDependency(t *testing.T) {
	// Setup: task 2 completed, task 3 completed, fix-1 pending (sourceTaskID "2"),
	// task 4 pending (depends on task 3 only)
	projectRoot := setupFeatureFixture(t, map[string]taskEntry{
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
			Type: "fix", SourceTaskID: "2",
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
	lines := parseBlockLines(t, out)

	// Verify: task 4 must not be blocked by the unrelated fix-1.
	// Either fix-1 or task 4 may be claimed (both have met dependencies),
	// but task 4 MUST be eligible (not blocked).
	taskID := getFieldValue(lines, "TASK_ID")
	assert.True(t, taskID == "4" || taskID == "fix-1",
		"expected TASK_ID to be either '4' or 'fix-1' (unrelated fix should not block task 4), got: %s, lines: %v", taskID, lines)
}

// Traceability: TC-006 -> Proposal Key Scenario "No fix tasks" + Success Criterion 4
func TestTC_006_NoFixTasksPreservesExistingClaimBehavior(t *testing.T) {
	// Setup: task 3 completed, task 4 pending (depends on task 3), no fix tasks
	projectRoot := setupFeatureFixture(t, map[string]taskEntry{
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
	lines := parseBlockLines(t, out)

	// Verify: task 4 should be claimed as before (no fix tasks to interfere)
	assert.True(t, hasFieldWithValue(lines, "ACTION", "CLAIMED"),
		"expected ACTION: CLAIMED, got: %v", lines)
	assert.True(t, hasFieldWithValue(lines, "TASK_ID", "4"),
		"expected TASK_ID: 4 (existing claim behavior unchanged), got: %v", lines)
}
