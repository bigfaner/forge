package task

import (
	"os"
	"path/filepath"
	"testing"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"
	"forge-cli/pkg/types"
)

// ---------------------------------------------------------------------------
// TestTransition_ClearsStateForClaimedTask
//
// When transitioning a task that is currently claimed (has process/state.json),
// doTransition should delete the state file so subsequent claim works.
// Root cause: transition only updated index.json, leaving stale state.json.
// ---------------------------------------------------------------------------
func TestTransition_ClearsStateForClaimedTask(t *testing.T) {
	tests := []struct {
		name         string
		targetStatus string
	}{
		{"pending", "pending"},
		{"blocked", "blocked"},
		{"skipped", "skipped"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := setupFullProject(t, SetupOpts{
				Tasks: map[string]task.Task{
					"t1": {ID: "1.1", Title: "Claimed Task", Status: types.StatusInProgress, Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
					"t2": {ID: "1.2", Title: "Other Task", Status: types.StatusPending, Priority: "P0", File: "1.2.md", Record: "records/1.2.md"},
				},
				State: &task.TaskState{
					TaskID: "1.1",
					Key:    "t1",
					Title:  "Claimed Task",
				},
			})

			indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
			statePath := feature.GetTaskStatePath(dir, "test")

			// Verify state file exists before transition
			if _, err := os.Stat(statePath); os.IsNotExist(err) {
				t.Fatal("state.json should exist before transition")
			}

			transitionReason = "resetting"
			err := doTransition(indexPath, statePath, "1.1", tt.targetStatus)
			if err != nil {
				t.Fatalf("doTransition() error = %v", err)
			}

			// State file should be deleted after transitioning claimed task away from in_progress
			if _, err := os.Stat(statePath); !os.IsNotExist(err) {
				t.Error("state.json should be deleted after transitioning claimed task away from in_progress")
			}

			// Verify index was updated
			index, _ := task.LoadIndex(indexPath)
			if index.TasksMap()["t1"].Status != types.Status(tt.targetStatus) {
				t.Errorf("task status = %q, want %q", index.TasksMap()["t1"].Status, tt.targetStatus)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestTransition_DoesNotClearStateForOtherTask
//
// Transitioning a different task (not the one in state.json) should NOT
// delete the state file.
// ---------------------------------------------------------------------------
func TestTransition_DoesNotClearStateForOtherTask(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "Claimed Task", Status: types.StatusInProgress, Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
			"t2": {ID: "1.2", Title: "Other Task", Status: types.StatusBlocked, Priority: "P0", File: "1.2.md", Record: "records/1.2.md"},
		},
		State: &task.TaskState{
			TaskID: "1.1",
			Key:    "t1",
			Title:  "Claimed Task",
		},
	})

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	statePath := feature.GetTaskStatePath(dir, "test")

	transitionReason = "unblocking"
	err := doTransition(indexPath, statePath, "1.2", "pending")
	if err != nil {
		t.Fatalf("doTransition() error = %v", err)
	}

	// State file should still exist (it's for task 1.1, not 1.2)
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Error("state.json should NOT be deleted when transitioning a different task")
	}
}

// ---------------------------------------------------------------------------
// TestTransition_NoStateFileIsNoop
//
// Transitioning when no state.json exists should work fine (no error).
// ---------------------------------------------------------------------------
func TestTransition_NoStateFileIsNoop(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "Task", Status: types.StatusBlocked, Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
		},
		// No State — simulates no previously claimed task
	})

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	statePath := feature.GetTaskStatePath(dir, "test")

	transitionReason = "unblocking"
	err := doTransition(indexPath, statePath, "1.1", "pending")
	if err != nil {
		t.Fatalf("doTransition() error = %v", err)
	}

	index, _ := task.LoadIndex(indexPath)
	if index.TasksMap()["t1"].Status != types.StatusPending {
		t.Errorf("task status = %q, want pending", index.TasksMap()["t1"].Status)
	}
}

// ---------------------------------------------------------------------------
// TestTransition_ThenClaimSucceeds
//
// End-to-end: claim -> transition to pending -> claim again should work
// without data integrity errors. This is the exact bug scenario.
// ---------------------------------------------------------------------------
func TestTransition_ThenClaimSucceeds(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "Task One", Status: types.StatusInProgress, Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
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

	// Step 1: Transition task 1.1 from in_progress to pending
	transitionReason = "resetting"
	err := doTransition(indexPath, statePath, "1.1", "pending")
	if err != nil {
		t.Fatalf("doTransition() error = %v", err)
	}

	// Step 2: Claim should succeed without data integrity error
	result, err := executeClaim()
	if err != nil {
		t.Fatalf("executeClaim() after transition should succeed, got error: %v", err)
	}

	// Should claim task 1.1 (pending, no deps, P0, lowest topo depth)
	if result.Action != "CLAIMED" {
		t.Errorf("expected CLAIMED, got %q", result.Action)
	}
	if result.Task.ID != "1.1" {
		t.Errorf("expected task 1.1, got %q", result.Task.ID)
	}
}
