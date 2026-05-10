package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"task-cli/pkg/feature"
	"task-cli/pkg/task"
)

func TestVerifyTaskCompletion(t *testing.T) {
	t.Run("no project root env returns nil", func(t *testing.T) {
		// Use a temp dir with no git root or feature state to fully isolate
		dir := t.TempDir()
		t.Setenv("CLAUDE_PROJECT_DIR", dir)
		// Also clear PROJECT_ROOT to prevent fallback to real project
		oldProjectRoot := os.Getenv("PROJECT_ROOT")
		_ = os.Unsetenv("PROJECT_ROOT")
		defer func() { _ = os.Setenv("PROJECT_ROOT", oldProjectRoot) }()

		err := verifyTaskCompletion()
		if err != nil {
			t.Errorf("expected nil, got error: %v", err)
		}
	})

	t.Run("no state file returns nil", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("CLAUDE_PROJECT_DIR", dir)

		err := verifyTaskCompletion()
		if err != nil {
			t.Errorf("expected nil, got error: %v", err)
		}
	})

	t.Run("no active feature returns nil", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("CLAUDE_PROJECT_DIR", dir)
		// Create features dir but no valid feature
		_ = os.MkdirAll(filepath.Join(dir, feature.FeaturesDir), 0755)

		err := verifyTaskCompletion()
		if err != nil {
			t.Errorf("expected nil, got error: %v", err)
		}
	})

	t.Run("incomplete task returns error", func(t *testing.T) {
		setupFullProject(t, SetupOpts{
			UseEnvVar: true,
			Tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "in_progress"},
			},
			State: &task.TaskState{TaskID: "1.1", Key: "task1"},
		})

		err := verifyTaskCompletion()
		if err == nil {
			t.Error("expected error for incomplete task")
		}
	})

	t.Run("completed task without record returns error", func(t *testing.T) {
		setupFullProject(t, SetupOpts{
			UseEnvVar: true,
			Tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "completed", Record: "records/1.1.md"},
			},
			State: &task.TaskState{TaskID: "1.1", Key: "task1"},
		})

		err := verifyTaskCompletion()
		if err == nil {
			t.Error("expected error for missing record file")
		}
	})

	t.Run("completed task with record passes", func(t *testing.T) {
		dir := setupFullProject(t, SetupOpts{
			UseEnvVar: true,
			Tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "completed", Record: "records/1.1.md"},
			},
			State: &task.TaskState{TaskID: "1.1", Key: "task1"},
		})

		// Create record file
		recordsDir := filepath.Join(dir, feature.FeaturesDir, "test", feature.TasksDirName, feature.RecordsDirName)
		_ = os.MkdirAll(recordsDir, 0755)
		recordFile := filepath.Join(recordsDir, "1.1.md")
		_ = os.WriteFile(recordFile, []byte("record content"), 0644)

		err := verifyTaskCompletion()
		if err != nil {
			t.Errorf("expected nil for completed task with record, got: %v", err)
		}
	})

	t.Run("completed task without record field passes", func(t *testing.T) {
		setupFullProject(t, SetupOpts{
			UseEnvVar: true,
			Tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "completed"},
			},
			State: &task.TaskState{TaskID: "1.1", Key: "task1"},
		})

		err := verifyTaskCompletion()
		if err != nil {
			t.Errorf("expected nil for completed task without record field, got: %v", err)
		}
	})
}
