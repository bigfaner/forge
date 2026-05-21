//go:build e2e

package qualitygate

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
// Quality gate fix-task loop breaker tests — feature: quality-gate-fix-task-loop-breaker
// Tests verify that fix-task creation is capped cumulatively per step,
// SourceTaskID uses step-scoped sentinel, and cross-step independence holds.
// ==============================================================================

// qgTaskEntry represents a single task in the index.json fixture.
type qgTaskEntry struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Priority      string   `json:"priority,omitempty"`
	EstimatedTime string   `json:"estimatedTime,omitempty"`
	Dependencies  []string `json:"dependencies,omitempty"`
	Status        string   `json:"status"`
	File          string   `json:"file"`
	Record        string   `json:"record,omitempty"`
	Type          string   `json:"type,omitempty"`
	SourceTaskID  string   `json:"sourceTaskID,omitempty"`
	Breaking      bool     `json:"breaking,omitempty"`
}

// qgIndexFixture represents the top-level index.json structure.
type qgIndexFixture struct {
	Feature      string                 `json:"feature"`
	StatusEnum   []string               `json:"statusEnum"`
	PriorityEnum []string               `json:"priorityEnum"`
	Tasks        map[string]qgTaskEntry `json:"tasks"`
}

// qgSetupProject creates a temp project dir with go.mod and feature structure.
// It creates docs/features/<slug>/tasks/ with index.json containing the given tasks,
// plus .forge/state.json with allCompleted=true.
// Returns the temp dir path.
func qgSetupProject(t *testing.T, slug string, tasks map[string]qgTaskEntry) string {
	t.Helper()
	dir := t.TempDir()

	// Create go.mod (project root marker)
	err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n\ngo 1.26\n"), 0644)
	assert.NoError(t, err, "failed to create go.mod")

	// Create feature task directory structure
	tasksDir := filepath.Join(dir, "docs", "features", slug, "tasks")
	processDir := filepath.Join(tasksDir, "process")
	err = os.MkdirAll(processDir, 0755)
	assert.NoError(t, err, "failed to create process dir")

	// Create task markdown files for each entry
	for _, task := range tasks {
		taskFile := filepath.Join(tasksDir, task.File)
		err := os.MkdirAll(filepath.Dir(taskFile), 0755)
		assert.NoError(t, err, "failed to create task file dir for %s", task.ID)
		err = os.WriteFile(taskFile, []byte("# Task "+task.ID+"\n"), 0644)
		assert.NoError(t, err, "failed to write task file %s", task.File)
	}

	// Write index.json
	idx := qgIndexFixture{
		Feature:      slug,
		StatusEnum:   []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"},
		PriorityEnum: []string{"P0", "P1", "P2"},
		Tasks:        tasks,
	}
	idxData, err := json.MarshalIndent(idx, "", "  ")
	assert.NoError(t, err, "failed to marshal index.json")
	err = os.WriteFile(filepath.Join(tasksDir, "index.json"), idxData, 0644)
	assert.NoError(t, err, "failed to write index.json")

	// Write .forge/state.json with allCompleted=true
	forgeDir := filepath.Join(dir, ".forge")
	err = os.MkdirAll(forgeDir, 0755)
	assert.NoError(t, err, "failed to create .forge dir")
	state := map[string]any{
		"feature":      slug,
		"allCompleted": true,
	}
	stateData, err := json.MarshalIndent(state, "", "  ")
	assert.NoError(t, err, "failed to marshal state.json")
	err = os.WriteFile(filepath.Join(forgeDir, "state.json"), stateData, 0644)
	assert.NoError(t, err, "failed to write state.json")

	return dir
}

// qgWriteJustfile creates a justfile in the project root with a failing compile recipe.
func qgWriteJustfile(t *testing.T, projectRoot string) {
	t.Helper()
	content := `compile:
    @echo "compile error: main.go:10: undefined: foo" && exit 1

fmt:
    @exit 0

lint:
    @exit 0

test:
    @exit 0
`
	err := os.WriteFile(filepath.Join(projectRoot, "justfile"), []byte(content), 0644)
	assert.NoError(t, err, "failed to write justfile")
}

// qgWritePassingJustfile creates a justfile where all recipes pass.
func qgWritePassingJustfile(t *testing.T, projectRoot string) {
	t.Helper()
	content := `compile:
    @exit 0

fmt:
    @exit 0

lint:
    @exit 0

test:
    @exit 0
`
	err := os.WriteFile(filepath.Join(projectRoot, "justfile"), []byte(content), 0644)
	assert.NoError(t, err, "failed to write justfile")
}

// qgRunQualityGate runs "forge quality-gate" with the project root as CLAUDE_PROJECT_DIR.
// Returns combined output and exit code. Does NOT fatalf on failure.
func qgRunQualityGate(t *testing.T, projectRoot string) (string, int) {
	t.Helper()
	cmd := exec.Command(testkit.ForgeBinary, "quality-gate")
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

// qgLoadIndex reads and parses the index.json from the given project.
func qgLoadIndex(t *testing.T, projectRoot, slug string) qgIndexFixture {
	t.Helper()
	indexPath := filepath.Join(projectRoot, "docs", "features", slug, "tasks", "index.json")
	data, err := os.ReadFile(indexPath)
	assert.NoError(t, err, "failed to read index.json")
	var idx qgIndexFixture
	err = json.Unmarshal(data, &idx)
	assert.NoError(t, err, "failed to unmarshal index.json")
	return idx
}

// qgCountFixTasksForStep counts fix tasks for a given step in the index.
// A fix task is identified by having SourceTaskID "quality-gate:<step>" AND
// a title with the prefix "fix <step>:".
func qgCountFixTasksForStep(idx qgIndexFixture, step string) int {
	count := 0
	prefix := "fix " + step + ":"
	sentinel := "quality-gate:" + step
	for _, task := range idx.Tasks {
		if task.SourceTaskID == sentinel && strings.HasPrefix(task.Title, prefix) {
			count++
		}
	}
	return count
}

// qgCountAllFixTasks counts all fix tasks across all steps in the index.
func qgCountAllFixTasks(idx qgIndexFixture) int {
	count := 0
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.SourceTaskID, "quality-gate:") &&
			strings.HasPrefix(task.Title, "fix ") {
			count++
		}
	}
	return count
}

// Traceability: TC-001 -> Proposal SC #1
func TestTC_001_AddFixTaskCreatesStepScopedSourceTaskID(t *testing.T) {
	slug := "test-qg-tc001"
	projectRoot := qgSetupProject(t, slug, map[string]qgTaskEntry{
		"t1": {ID: "1.1", Status: "completed", File: "1.1.md", Type: "coding.feature"},
	})
	qgWriteJustfile(t, projectRoot)

	output, _ := qgRunQualityGate(t, projectRoot)

	// Verify stderr mentions the compile step failed
	assert.Contains(t, output, "compile check failed",
		"quality gate should report compile failure")

	// Load the updated index and verify fix task has step-scoped SourceTaskID
	idx := qgLoadIndex(t, projectRoot, slug)
	fixCount := qgCountFixTasksForStep(idx, "compile")
	assert.Equal(t, 1, fixCount,
		"expected exactly 1 fix task for compile step")

	// Find the fix task and verify its SourceTaskID sentinel
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.Title, "fix compile:") {
			assert.Equal(t, "quality-gate:compile", task.SourceTaskID,
				"fix task SourceTaskID should use step-scoped sentinel format")
			assert.Equal(t, "P0", task.Priority,
				"fix task should be P0 priority")
			assert.True(t, task.Breaking,
				"fix task should be marked as breaking")
		}
	}

	// Verify fix task markdown was created
	tasksDir := filepath.Join(projectRoot, "docs", "features", slug, "tasks")
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.Title, "fix compile:") {
			mdPath := filepath.Join(tasksDir, task.File)
			data, err := os.ReadFile(mdPath)
			assert.NoError(t, err, "fix task markdown file should exist")
			content := string(data)
			assert.Contains(t, content, "N/A (project-wide gate)",
				"task markdown should contain Vars SOURCE_TASK_ID placeholder")
			assert.Contains(t, content, "just compile",
				"task markdown should reference the test script")
		}
	}
}

// Traceability: TC-006 -> Proposal SC #6
func TestTC_006_CumulativeCapStopsFixTaskAfter3(t *testing.T) {
	slug := "test-qg-tc006"
	// Pre-populate 3 fix tasks for "compile" step.
	// ALL tasks must be completed/skipped for checkAllCompleted to pass.
	// countFixTasks counts ALL statuses cumulatively, so these 3 completed/skipped
	// tasks still count toward the cap of 3.
	// Map keys must match the fix-task ID pattern so generateAutoID skips over them.
	projectRoot := qgSetupProject(t, slug, map[string]qgTaskEntry{
		"t1":    {ID: "1.1", Status: "completed", File: "1.1.md", Type: "coding.feature"},
		"fix-1": {ID: "fix-1", Title: "fix compile: first error", SourceTaskID: "quality-gate:compile", Status: "completed", File: "fix-1.md", Breaking: true, Priority: "P0", Type: "coding.fix"},
		"fix-2": {ID: "fix-2", Title: "fix compile: second error", SourceTaskID: "quality-gate:compile", Status: "skipped", File: "fix-2.md", Breaking: true, Priority: "P0", Type: "coding.fix"},
		"fix-3": {ID: "fix-3", Title: "fix compile: third error", SourceTaskID: "quality-gate:compile", Status: "completed", File: "fix-3.md", Breaking: true, Priority: "P0", Type: "coding.fix"},
	})
	qgWriteJustfile(t, projectRoot)

	output, _ := qgRunQualityGate(t, projectRoot)

	// Verify stderr contains cap warning
	assert.Contains(t, output, "max fix-tasks reached for compile",
		"should report max fix-tasks reached when 3 cumulative fix tasks exist")

	// Verify no 4th fix task was created
	idx := qgLoadIndex(t, projectRoot, slug)
	fixCount := qgCountFixTasksForStep(idx, "compile")
	assert.Equal(t, 3, fixCount,
		"should still have exactly 3 fix tasks for compile (cap enforced)")
}

// Traceability: TC-007 -> Proposal SC #7
func TestTC_007_CrossStepIndependenceFixADoesNotBlockB(t *testing.T) {
	slug := "test-qg-tc007"
	// Pre-populate 3 fix tasks for "compile" step (at cap).
	// ALL tasks must be completed/skipped for checkAllCompleted to pass.
	// Map keys match fix-task ID pattern for correct auto-ID generation.
	projectRoot := qgSetupProject(t, slug, map[string]qgTaskEntry{
		"t1":    {ID: "1.1", Status: "completed", File: "1.1.md", Type: "coding.feature"},
		"fix-1": {ID: "fix-1", Title: "fix compile: first", SourceTaskID: "quality-gate:compile", Status: "completed", File: "fix-1.md", Breaking: true, Priority: "P0", Type: "coding.fix"},
		"fix-2": {ID: "fix-2", Title: "fix compile: second", SourceTaskID: "quality-gate:compile", Status: "skipped", File: "fix-2.md", Breaking: true, Priority: "P0", Type: "coding.fix"},
		"fix-3": {ID: "fix-3", Title: "fix compile: third", SourceTaskID: "quality-gate:compile", Status: "completed", File: "fix-3.md", Breaking: true, Priority: "P0", Type: "coding.fix"},
	})

	// Create a justfile where compile passes but lint fails
	// This allows us to trigger a fix task for lint while compile is capped
	content := `compile:
    @exit 0

fmt:
    @exit 0

lint:
    @echo "lint error: handler.go:5: unused variable" && exit 1

test:
    @exit 0
`
	err := os.WriteFile(filepath.Join(projectRoot, "justfile"), []byte(content), 0644)
	assert.NoError(t, err, "failed to write justfile")

	output, _ := qgRunQualityGate(t, projectRoot)

	// Compile should pass, lint should fail -> fix task for lint created
	assert.Contains(t, output, "Lint check failed",
		"lint step should fail and be reported")

	// Verify lint fix task was created despite compile being at cap
	idx := qgLoadIndex(t, projectRoot, slug)
	lintFixCount := qgCountFixTasksForStep(idx, "lint")
	assert.Equal(t, 1, lintFixCount,
		"lint fix task should be created independently of compile cap")

	// Verify compile is still at 3
	compileFixCount := qgCountFixTasksForStep(idx, "compile")
	assert.Equal(t, 3, compileFixCount,
		"compile fix tasks should remain at 3")

	// Verify the new fix task has lint-scoped sentinel
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.Title, "fix lint:") {
			assert.Equal(t, "quality-gate:lint", task.SourceTaskID,
				"lint fix task should have lint-scoped SourceTaskID")
		}
	}
}

// Traceability: TC-002 -> Proposal SC #2
func TestTC_002_CountFixTasksCountsCumulativeRegardlessOfStatus(t *testing.T) {
	slug := "test-qg-tc002"
	// Create a project with 3 completed + 1 skipped fix tasks for compile (total = 4).
	// ALL tasks must be completed/skipped for checkAllCompleted to pass.
	// Total = 4 exceeds the cap of 3, so next attempt should be blocked.
	projectRoot := qgSetupProject(t, slug, map[string]qgTaskEntry{
		"t1":    {ID: "1.1", Status: "completed", File: "1.1.md", Type: "coding.feature"},
		"fix-1": {ID: "fix-1", Title: "fix compile: first", SourceTaskID: "quality-gate:compile", Status: "completed", File: "fix-1.md", Breaking: true, Priority: "P0", Type: "coding.fix"},
		"fix-2": {ID: "fix-2", Title: "fix compile: second", SourceTaskID: "quality-gate:compile", Status: "completed", File: "fix-2.md", Breaking: true, Priority: "P0", Type: "coding.fix"},
		"fix-3": {ID: "fix-3", Title: "fix compile: third", SourceTaskID: "quality-gate:compile", Status: "completed", File: "fix-3.md", Breaking: true, Priority: "P0", Type: "coding.fix"},
		"fix-4": {ID: "fix-4", Title: "fix compile: fourth", SourceTaskID: "quality-gate:compile", Status: "skipped", File: "fix-4.md", Breaking: true, Priority: "P0", Type: "coding.fix"},
	})
	qgWriteJustfile(t, projectRoot)

	output, _ := qgRunQualityGate(t, projectRoot)

	// With 4 cumulative fix tasks (already over cap), should report max reached
	assert.Contains(t, output, "max fix-tasks reached for compile",
		"should block when cumulative count (4) exceeds cap (3)")

	// Verify still only 4 fix tasks (no 5th created)
	idx := qgLoadIndex(t, projectRoot, slug)
	fixCount := qgCountFixTasksForStep(idx, "compile")
	assert.Equal(t, 4, fixCount,
		"should have 4 cumulative fix tasks, no new one created")
}

// Traceability: TC-003 -> Proposal SC #3
func TestTC_003_QualityGateExits0OnNotAllCompleted(t *testing.T) {
	slug := "test-qg-tc003"
	// Create a project where tasks are NOT all completed
	// quality-gate should exit 0 silently
	projectRoot := qgSetupProject(t, slug, map[string]qgTaskEntry{
		"t1": {ID: "1.1", Status: "completed", File: "1.1.md", Type: "coding.feature"},
		"t2": {ID: "1.2", Status: "pending", File: "1.2.md", Type: "coding.feature"},
	})

	output, exitCode := qgRunQualityGate(t, projectRoot)

	assert.Equal(t, 0, exitCode,
		"quality-gate should exit 0 when not all tasks completed")
	assert.NotContains(t, output, "All tasks completed",
		"should not print completion message when tasks are still pending")
}

// Traceability: TC-004 -> Proposal SC #4
func TestTC_004_QualityGateSkipsDocsOnlyFeatures(t *testing.T) {
	slug := "test-qg-tc004"
	// Create a docs-only feature (no implementation or fix tasks)
	projectRoot := qgSetupProject(t, slug, map[string]qgTaskEntry{
		"t1": {ID: "1.1", Status: "completed", File: "1.1.md", Type: "documentation"},
		"t2": {ID: "T-eval-1", Status: "completed", File: "T-eval-1.md", Type: "doc-evaluation"},
	})

	output, exitCode := qgRunQualityGate(t, projectRoot)

	assert.Equal(t, 0, exitCode,
		"quality-gate should exit 0 for docs-only features")
	assert.Contains(t, output, "docs-only",
		"should mention docs-only skip reason")
}

// Traceability: TC-005 -> Proposal SC #5
func TestTC_005_FixTaskMarkdownCreatedOnDisk(t *testing.T) {
	slug := "test-qg-tc005"
	projectRoot := qgSetupProject(t, slug, map[string]qgTaskEntry{
		"t1": {ID: "1.1", Status: "completed", File: "1.1.md", Type: "coding.feature"},
	})
	qgWriteJustfile(t, projectRoot)

	qgRunQualityGate(t, projectRoot)

	idx := qgLoadIndex(t, projectRoot, slug)
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.Title, "fix compile:") {
			mdPath := filepath.Join(projectRoot, "docs", "features", slug, "tasks", task.File)
			data, err := os.ReadFile(mdPath)
			assert.NoError(t, err, "fix task markdown file should exist on disk")
			content := string(data)
			// Verify fix task markdown contains expected sections from the fix-task template
			assert.Contains(t, content, "Root Cause",
				"fix task markdown should contain Root Cause section")
			assert.Contains(t, content, "Verification",
				"fix task markdown should contain Verification section")
			assert.Contains(t, content, "Reference Files",
				"fix task markdown should contain Reference Files section")
			// Verify it references the error output file
			assert.Contains(t, content, "tests/results/unit-raw-output.txt",
				"fix task should reference the error output path")
		}
	}
}
