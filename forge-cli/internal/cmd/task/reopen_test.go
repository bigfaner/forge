package task

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"
	"forge-cli/pkg/types"
)

// ---------------------------------------------------------------------------
// TestReopen_RejectedToPending
//
// Reopen a rejected task should succeed and set status to pending.
// ---------------------------------------------------------------------------
func TestReopen_RejectedToPending(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "Rejected Task", Status: "rejected", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	}})

	out := captureStdout(func() {
		if err := runReopen(nil, []string{"1.1"}); err != nil {
			base.Exit(err)
		}
	})
	if !strings.Contains(out, "STATUS: pending") {
		t.Errorf("expected status pending, got: %s", out)
	}
	if !strings.Contains(out, "TASK_ID: 1.1") {
		t.Errorf("expected TASK_ID in output, got: %s", out)
	}

	// Verify index was updated
	dir, _ := os.Getwd()
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index, _ := task.LoadIndex(indexPath)
	if index.TasksMap()["t1"].Status != "pending" {
		t.Errorf("index status = %q, want pending", index.TasksMap()["t1"].Status)
	}
}

// ---------------------------------------------------------------------------
// TestReopen_SkippedToPending
//
// Reopen a skipped task should succeed and set status to pending.
// ---------------------------------------------------------------------------
func TestReopen_SkippedToPending(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "Skipped Task", Status: "skipped", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	}})

	out := captureStdout(func() {
		if err := runReopen(nil, []string{"1.1"}); err != nil {
			base.Exit(err)
		}
	})
	if !strings.Contains(out, "STATUS: pending") {
		t.Errorf("expected status pending, got: %s", out)
	}

	dir, _ := os.Getwd()
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index, _ := task.LoadIndex(indexPath)
	if index.TasksMap()["t1"].Status != "pending" {
		t.Errorf("index status = %q, want pending", index.TasksMap()["t1"].Status)
	}
}

// ---------------------------------------------------------------------------
// TestReopen_CompletedBlocked
//
// Completed tasks are NEVER re-openable.
// ---------------------------------------------------------------------------
func TestReopen_CompletedBlocked(t *testing.T) {
	if os.Getenv("TEST_REOPEN_COMPLETED") == "1" {
		setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "Completed Task", Status: "completed", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
		}})
		if err := runReopen(nil, []string{"1.1"}); err != nil {
			base.Exit(err)
		}
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestReopen_CompletedBlocked")
	cmd.Env = append(os.Environ(), "TEST_REOPEN_COMPLETED=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected error: completed tasks should not be re-openable")
	}
	out := string(output)
	if !strings.Contains(out, "INVALID_TRANSITION") {
		t.Errorf("expected INVALID_TRANSITION error, got: %s", out)
	}
}

// ---------------------------------------------------------------------------
// TestReopen_NonTerminalBlocked
//
// Reopen on non-terminal states (pending, in_progress, blocked) should fail.
// ---------------------------------------------------------------------------
func TestReopen_NonTerminalBlocked(t *testing.T) {
	tests := []struct {
		name   string
		status string
	}{
		{"pending", "pending"},
		{"in_progress", "in_progress"},
		{"blocked", "blocked"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envKey := "TEST_REOPEN_NONTERM_" + strings.ToUpper(tt.name)
			if os.Getenv(envKey) == "1" {
				setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
					"t1": {ID: "1.1", Title: "Task", Status: types.Status(tt.status), Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
				}})
				if err := runReopen(nil, []string{"1.1"}); err != nil {
					base.Exit(err)
				}
				return
			}

			cmd := exec.Command(os.Args[0], "-test.run=TestReopen_NonTerminalBlocked/"+tt.name)
			cmd.Env = append(os.Environ(), envKey+"=1")
			output, err := cmd.CombinedOutput()
			if err == nil {
				t.Errorf("expected error: reopen on %s should fail", tt.status)
			}
			out := string(output)
			if !strings.Contains(out, "INVALID_TRANSITION") {
				t.Errorf("expected INVALID_TRANSITION for %s, got: %s", tt.status, out)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestReopen_TaskNotFound
//
// Reopen a nonexistent task should fail.
// ---------------------------------------------------------------------------
func TestReopen_TaskNotFound(t *testing.T) {
	if os.Getenv("TEST_REOPEN_NOT_FOUND") == "1" {
		setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "Task", Status: "rejected", Priority: "P0", File: "1.1.md"},
		}})
		if err := runReopen(nil, []string{"9.9"}); err != nil {
			base.Exit(err)
		}
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestReopen_TaskNotFound")
	cmd.Env = append(os.Environ(), "TEST_REOPEN_NOT_FOUND=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected error for task not found")
	}
	out := string(output)
	if !strings.Contains(out, "NOT_FOUND") {
		t.Errorf("expected NOT_FOUND error, got: %s", out)
	}
}

// ---------------------------------------------------------------------------
// TestReopen_ValidateTransition_UsesRoleReopen
//
// Verify doReopen uses ValidateTransition with RoleReopen.
// This is a unit-level test that calls doReopen directly.
// ---------------------------------------------------------------------------
func TestReopen_ValidateTransition_UsesRoleReopen(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		wantErr bool
	}{
		{"rejected -> pending succeeds", "rejected", false},
		{"skipped -> pending succeeds", "skipped", false},
		{"completed -> pending fails", "completed", true},
		{"pending -> pending fails", "pending", true},
		{"in_progress -> pending fails", "in_progress", true},
		{"blocked -> pending fails", "blocked", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
				"t1": {ID: "1.1", Title: "Task", Status: types.Status(tt.status), Priority: "P0", File: "1.1.md"},
			}})

			indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
			err := doReopen(indexPath, "", "1.1")

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %s -> pending", tt.status)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %s -> pending: %v", tt.status, err)
				}
				// Verify status was updated
				index, _ := task.LoadIndex(indexPath)
				if index.TasksMap()["t1"].Status != "pending" {
					t.Errorf("status = %q, want pending", index.TasksMap()["t1"].Status)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestReopen_WithLock
//
// Verify reopen uses WithLock for index write. We test this by verifying
// the runReopen path succeeds and the index is properly updated.
// ---------------------------------------------------------------------------
func TestReopen_WithLock(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "Rejected Task", Status: "rejected", Priority: "P0", File: "1.1.md"},
		},
		UseEnvVar: true,
	})

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))

	// Use runReopen which goes through WithLock
	Cmd.SetArgs([]string{"reopen", "1.1"})
	if err := Cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Lock file should exist (WithLock persists it)
	lockPath := indexPath + ".lock"
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Error("lock file should exist after WithLock usage")
	}

	// Task should be pending
	index, _ := task.LoadIndex(indexPath)
	if index.TasksMap()["t1"].Status != "pending" {
		t.Errorf("status = %q, want pending", index.TasksMap()["t1"].Status)
	}
}

// ---------------------------------------------------------------------------
// TestReopen_SlugKeyedTask
//
// Reopen should work with tasks that have slug keys different from their IDs.
// ---------------------------------------------------------------------------
func TestReopen_SlugKeyedTask(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"run-e2e": {ID: "T-test-run", Title: "Run E2E", Status: "rejected", Priority: "P0", File: "T-test-run.md"},
	}})

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))

	// Should find the task by ID and reopen it
	err := doReopen(indexPath, "", "T-test-run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	index, _ := task.LoadIndex(indexPath)
	if index.TasksMap()["run-e2e"].Status != "pending" {
		t.Errorf("status = %q, want pending", index.TasksMap()["run-e2e"].Status)
	}
}

// ---------------------------------------------------------------------------
// TestStatus_ForceFlagRemoved
//
// Verify --force flag is no longer registered on the status command.
// ---------------------------------------------------------------------------
func TestStatus_ForceFlagRemoved(t *testing.T) {
	flag := statusCmd.Flags().Lookup("force")
	if flag != nil {
		t.Error("--force flag should be removed from status command")
	}
}

// ---------------------------------------------------------------------------
// TestStatus_ReadOnly_AnyStatusArgument
//
// Status command with 2 args mutates task status via the state machine.
// Transitions blocked by the state machine (terminal, completed) return errors.
// ---------------------------------------------------------------------------
func TestStatus_ReadOnly_AnyStatusArgument(t *testing.T) {
	// Status command uses ExactArgs(1), so cobra rejects 2-arg calls at the framework level.
	// Verify the Args validator rejects extra arguments.
	err := statusCmd.Args(statusCmd, []string{"1.1", "completed"})
	if err == nil {
		t.Error("expected ExactArgs(1) to reject 2 arguments")
	}
}

// ---------------------------------------------------------------------------
// TestReopen_CLI_Integration
//
// Full CLI integration test for forge task reopen.
// ---------------------------------------------------------------------------
func TestReopen_CLI_Integration(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "Rejected Task", Status: "rejected", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
		},
		UseEnvVar: true,
	})

	// Build args for root command
	Cmd.SetArgs([]string{"reopen", "1.1"})
	_ = dir

	// Capture output
	out := captureStdout(func() {
		if err := Cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "STATUS: pending") {
		t.Errorf("expected status pending in output, got: %s", out)
	}
}

// ---------------------------------------------------------------------------
// TestReopen_WithLock_SaveIndexError
//
// Verify error handling when WithLock fails.
// ---------------------------------------------------------------------------
func TestReopen_WithLock_SaveIndexError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix chmod 0555 does not prevent writes on Windows")
	}
	if os.Getenv("TEST_REOPEN_SAVE_ERROR") == "1" {
		dir := setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "Rejected Task", Status: "rejected", Priority: "P0", File: "1.1.md"},
		}})

		indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
		// Make tasks directory read-only so LockFile fails
		tasksDir := filepath.Dir(indexPath)
		_ = os.Chmod(tasksDir, 0555)
		defer func() { _ = os.Chmod(tasksDir, 0755) }()

		if err := runReopen(nil, []string{"1.1"}); err != nil {
			base.Exit(err)
		}
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestReopen_WithLock_SaveIndexError")
	cmd.Env = append(os.Environ(), "TEST_REOPEN_SAVE_ERROR=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected error when lock cannot be acquired")
	}
	out := string(output)
	if !strings.Contains(out, "Failed to acquire lock") {
		t.Errorf("expected lock acquisition error, got: %s", out)
	}
}

// ---------------------------------------------------------------------------
// TestReopen_NoProject
// ---------------------------------------------------------------------------
func TestReopen_NoProject(t *testing.T) {
	if os.Getenv("TEST_REOPEN_NO_PROJECT") == "1" {
		if err := runReopen(nil, []string{"1.1"}); err != nil {
			base.Exit(err)
		}
		return
	}

	tmpDir := t.TempDir()
	cmd := exec.Command(os.Args[0], "-test.run=TestReopen_NoProject")
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "CLAUDE_PROJECT_DIR=") {
			continue
		}
		env = append(env, e)
	}
	cmd.Env = append(slices.Clone(env), "TEST_REOPEN_NO_PROJECT=1", "CLAUDE_PROJECT_DIR=")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected error for no project root")
	}
	if !strings.Contains(string(output), "NO_PROJECT") && !strings.Contains(string(output), "NO_FEATURE") {
		t.Errorf("expected NO_PROJECT or NO_FEATURE error, got: %s", string(output))
	}
}

// ---------------------------------------------------------------------------
// TestReopen_SetsStatusToPending_WithExistingIndex
//
// Verify reopen writes the updated index correctly, preserving other tasks.
// ---------------------------------------------------------------------------
func TestReopen_SetsStatusToPending_WithExistingIndex(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "Rejected Task", Status: "rejected", Priority: "P0", File: "1.1.md"},
		"t2": {ID: "1.2", Title: "Other Task", Status: "completed", Priority: "P1", File: "1.2.md", Record: "records/1.2.md"},
		"t3": {ID: "1.3", Title: "Pending Task", Status: "pending", Priority: "P0", File: "1.3.md"},
	}})

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	err := doReopen(indexPath, "", "1.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	index, _ := task.LoadIndex(indexPath)

	// t1 should be pending
	if index.TasksMap()["t1"].Status != "pending" {
		t.Errorf("t1 status = %q, want pending", index.TasksMap()["t1"].Status)
	}

	// Other tasks should be unchanged
	if index.TasksMap()["t2"].Status != "completed" {
		t.Errorf("t2 status = %q, want completed (unchanged)", index.TasksMap()["t2"].Status)
	}
	if index.TasksMap()["t3"].Status != "pending" {
		t.Errorf("t3 status = %q, want pending (unchanged)", index.TasksMap()["t3"].Status)
	}
}

// ---------------------------------------------------------------------------
// TestReopen_ClearsStateForClaimedTask
//
// When reopening a task that is tracked in process/state.json, doReopen
// should delete state.json so subsequent claim works.
// ---------------------------------------------------------------------------
func TestReopen_ClearsStateForClaimedTask(t *testing.T) {
	tests := []struct {
		name   string
		status string
	}{
		{"rejected", "rejected"},
		{"skipped", "skipped"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := setupFullProject(t, SetupOpts{
				Tasks: map[string]task.Task{
					"t1": {ID: "1.1", Title: "Task", Status: types.Status(tt.status), Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
					"t2": {ID: "1.2", Title: "Other", Status: types.StatusPending, Priority: "P0", File: "1.2.md", Record: "records/1.2.md", Dependencies: []string{}},
				},
				State: &task.TaskState{
					TaskID: "1.1",
					Key:    "t1",
					Title:  "Task",
				},
			})

			indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
			statePath := feature.GetTaskStatePath(dir, "test")

			if _, err := os.Stat(statePath); os.IsNotExist(err) {
				t.Fatal("state.json should exist before reopen")
			}

			err := doReopen(indexPath, statePath, "1.1")
			if err != nil {
				t.Fatalf("doReopen() error = %v", err)
			}

			if _, err := os.Stat(statePath); !os.IsNotExist(err) {
				t.Error("state.json should be deleted after reopening claimed task")
			}

			index, _ := task.LoadIndex(indexPath)
			if index.TasksMap()["t1"].Status != types.StatusPending {
				t.Errorf("task status = %q, want pending", index.TasksMap()["t1"].Status)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestReopen_DoesNotClearStateForOtherTask
// ---------------------------------------------------------------------------
func TestReopen_DoesNotClearStateForOtherTask(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "Claimed", Status: types.StatusInProgress, Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
			"t2": {ID: "1.2", Title: "Rejected", Status: "rejected", Priority: "P0", File: "1.2.md", Record: "records/1.2.md"},
		},
		State: &task.TaskState{
			TaskID: "1.1",
			Key:    "t1",
			Title:  "Claimed",
		},
	})

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	statePath := feature.GetTaskStatePath(dir, "test")

	err := doReopen(indexPath, statePath, "1.2")
	if err != nil {
		t.Fatalf("doReopen() error = %v", err)
	}

	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Error("state.json should NOT be deleted when reopening a different task")
	}
}

// ---------------------------------------------------------------------------
// TestReopen_ThenClaimSucceeds
//
// skipped with stale state -> reopen -> claim should work.
// CheckExistingTaskState does NOT handle "skipped" in its switch,
// so this would fail without the state cleanup in doReopen.
// ---------------------------------------------------------------------------
func TestReopen_ThenClaimSucceeds(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "Task One", Status: "skipped", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
			"t2": {ID: "1.2", Title: "Task Two", Status: types.StatusPending, Priority: "P0", File: "1.2.md", Record: "records/1.2.md", Dependencies: []string{}},
		},
		State: &task.TaskState{
			TaskID: "1.1",
			Key:    "t1",
			Title:  "Task One",
		},
	})

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	statePath := feature.GetTaskStatePath(dir, "test")

	err := doReopen(indexPath, statePath, "1.1")
	if err != nil {
		t.Fatalf("doReopen() error = %v", err)
	}

	result, err := executeClaim()
	if err != nil {
		t.Fatalf("executeClaim() after reopen should succeed, got error: %v", err)
	}

	if result.Action != "CLAIMED" {
		t.Errorf("expected CLAIMED, got %q", result.Action)
	}
	if result.Task.ID != "1.1" {
		t.Errorf("expected task 1.1, got %q", result.Task.ID)
	}
}
