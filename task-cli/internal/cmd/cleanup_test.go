package cmd

import (
	"os"
	"testing"

	"task-cli/pkg/feature"
	"task-cli/pkg/task"
)

func TestCleanupCompletedTaskState(t *testing.T) {
	t.Run("no state file returns no error", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("CLAUDE_PROJECT_DIR", dir)
		cleanupCompletedTaskState()
	})

	t.Run("no project root env returns no error", func(t *testing.T) {
		cleanupCompletedTaskState()
	})

	t.Run("deletes state file when task is completed", func(t *testing.T) {
		dir := setupFullProject(t, SetupOpts{
			UseEnvVar: true,
			Tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "completed"},
			},
			State: &task.TaskState{TaskID: "1.1", Key: "task1"},
		})

		statePath := feature.GetTaskStatePath(dir, "test")
		cleanupCompletedTaskState()

		// Verify state file is deleted
		if _, err := os.Stat(statePath); !os.IsNotExist(err) {
			t.Errorf("state file should be deleted")
		}
	})

	t.Run("keeps state file when task is not completed", func(t *testing.T) {
		dir := setupFullProject(t, SetupOpts{
			UseEnvVar: true,
			Tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "in_progress"},
			},
			State: &task.TaskState{TaskID: "1.1", Key: "task1"},
		})

		statePath := feature.GetTaskStatePath(dir, "test")
		cleanupCompletedTaskState()

		// Verify state file still exists
		if _, err := os.Stat(statePath); os.IsNotExist(err) {
			t.Errorf("state file should NOT be deleted for in_progress task")
		}
	})

	t.Run("deletes record.json when task is completed", func(t *testing.T) {
		dir := setupFullProject(t, SetupOpts{
			UseEnvVar: true,
			Tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "completed"},
			},
			State: &task.TaskState{TaskID: "1.1", Key: "task1"},
		})

		statePath := feature.GetTaskStatePath(dir, "test")
		recordPath := feature.GetProcessRecordPath(dir, "test")

		// Create record.json
		os.WriteFile(recordPath, []byte("{}"), 0644)

		cleanupCompletedTaskState()

		// Verify both files are deleted
		if _, err := os.Stat(statePath); !os.IsNotExist(err) {
			t.Errorf("state file should be deleted")
		}
		if _, err := os.Stat(recordPath); !os.IsNotExist(err) {
			t.Errorf("record.json should be deleted")
		}
	})

	t.Run("preserves .forge/state.json during cleanup", func(t *testing.T) {
		dir := setupFullProject(t, SetupOpts{
			UseEnvVar: true,
			Tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "completed"},
			},
			State: &task.TaskState{TaskID: "1.1", Key: "task1"},
		})

		// Simulate .forge/state.json created by task claim
		feature.EnsureForgeState(dir, "test")

		cleanupCompletedTaskState()

		// .forge/state.json should NOT be deleted
		forgeState := feature.ReadForgeState(dir)
		if forgeState == nil {
			t.Error(".forge/state.json should NOT be deleted by cleanup")
		}
	})

	t.Run("preserves .forge/state.json with in_progress task", func(t *testing.T) {
		dir := setupFullProject(t, SetupOpts{
			UseEnvVar: true,
			Tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "in_progress"},
			},
			State: &task.TaskState{TaskID: "1.1", Key: "task1"},
		})

		// Simulate .forge/state.json created by task claim
		feature.EnsureForgeState(dir, "test")

		cleanupCompletedTaskState()

		forgeState := feature.ReadForgeState(dir)
		if forgeState == nil {
			t.Error(".forge/state.json should NOT be deleted by cleanup")
		}
	})

	t.Run("deletes state file when task is blocked", func(t *testing.T) {
		dir := setupFullProject(t, SetupOpts{
			UseEnvVar: true,
			Tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "blocked"},
			},
			State: &task.TaskState{TaskID: "1.1", Key: "task1"},
		})

		statePath := feature.GetTaskStatePath(dir, "test")
		recordPath := feature.GetProcessRecordPath(dir, "test")

		// Create record.json
		os.WriteFile(recordPath, []byte("{}"), 0644)

		cleanupCompletedTaskState()

		// Verify state file is deleted
		if _, err := os.Stat(statePath); !os.IsNotExist(err) {
			t.Errorf("state file should be deleted for blocked task")
		}
		if _, err := os.Stat(recordPath); !os.IsNotExist(err) {
			t.Errorf("record.json should be deleted for blocked task")
		}
	})
}
