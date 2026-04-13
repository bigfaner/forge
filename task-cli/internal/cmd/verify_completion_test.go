package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"task-cli/pkg/feature"
	"task-cli/pkg/task"
)

// setupTestProject creates a minimal test project for tests.
func setupTestProject(t *testing.T) string {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Create feature directory structure with index.json
	featureDir := filepath.Join(dir, feature.FeaturesDir, "test")
	tasksDir := filepath.Join(featureDir, feature.TasksDirName)
	processDir := filepath.Join(tasksDir, feature.ProcessDirName)
	os.MkdirAll(processDir, 0755)

	// Create index.json
	indexData := &task.TaskIndex{Feature: "test", Tasks: make(map[string]task.Task)}
	task.SaveIndex(filepath.Join(tasksDir, feature.IndexFileName), indexData)

	return dir
}

// setupTestProjectWithTask creates a test project with a task in the given status.
func setupTestProjectWithTask(t *testing.T, status string) string {
	dir := setupTestProject(t)

	// Create index with task
	indexPath := filepath.Join(dir, feature.FeaturesDir, "test", feature.TasksDirName, feature.IndexFileName)
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "in_progress", "completed"},
		Tasks: map[string]task.Task{
			"task1": {ID: "1.1", Title: "Task 1", Status: status},
		},
	}
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	// Create state
	statePath := feature.GetTaskStatePath(dir, "test")
	state := &task.TaskState{TaskID: "1.1", Key: "task1"}
	if err := task.SaveState(statePath, state); err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	return dir
}

func TestVerifyTaskCompletion(t *testing.T) {
	t.Run("no project root env returns nil", func(t *testing.T) {
		oldEnv := os.Getenv("CLAUDE_PROJECT_DIR")
		os.Unsetenv("CLAUDE_PROJECT_DIR")
		defer os.Setenv("CLAUDE_PROJECT_DIR", oldEnv)

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
		os.MkdirAll(filepath.Join(dir, feature.FeaturesDir), 0755)

		err := verifyTaskCompletion()
		if err != nil {
			t.Errorf("expected nil, got error: %v", err)
		}
	})

	t.Run("incomplete task returns error", func(t *testing.T) {
		setupTestProjectWithTask(t, "in_progress")

		err := verifyTaskCompletion()
		if err == nil {
			t.Error("expected error for incomplete task")
		}
	})

	t.Run("completed task without record returns error", func(t *testing.T) {
		dir := setupTestProject(t)

		// Create index with completed task that has a record file
		indexPath := filepath.Join(dir, feature.FeaturesDir, "test", feature.TasksDirName, feature.IndexFileName)
		index := &task.TaskIndex{
			Feature:    "test",
			StatusEnum: []string{"pending", "in_progress", "completed"},
			Tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "completed", Record: "records/1.1.md"},
			},
		}
		task.SaveIndex(indexPath, index)

		// Create state
		statePath := feature.GetTaskStatePath(dir, "test")
		state := &task.TaskState{TaskID: "1.1", Key: "task1"}
		task.SaveState(statePath, state)

		err := verifyTaskCompletion()
		if err == nil {
			t.Error("expected error for missing record file")
		}
	})

	t.Run("completed task with record passes", func(t *testing.T) {
		dir := setupTestProject(t)

		// Create index with completed task
		indexPath := filepath.Join(dir, feature.FeaturesDir, "test", feature.TasksDirName, feature.IndexFileName)
		index := &task.TaskIndex{
			Feature:    "test",
			StatusEnum: []string{"pending", "in_progress", "completed"},
			Tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "completed", Record: "records/1.1.md"},
			},
		}
		task.SaveIndex(indexPath, index)

		// Create record file
		recordsDir := filepath.Join(dir, feature.FeaturesDir, "test", feature.TasksDirName, feature.RecordsDirName)
		os.MkdirAll(recordsDir, 0755)
		recordFile := filepath.Join(recordsDir, "1.1.md")
		os.WriteFile(recordFile, []byte("record content"), 0644)

		// Create state
		statePath := feature.GetTaskStatePath(dir, "test")
		state := &task.TaskState{TaskID: "1.1", Key: "task1"}
		task.SaveState(statePath, state)

		err := verifyTaskCompletion()
		if err != nil {
			t.Errorf("expected nil for completed task with record, got: %v", err)
		}
	})

	t.Run("completed task without record field passes", func(t *testing.T) {
		setupTestProjectWithTask(t, "completed")

		err := verifyTaskCompletion()
		if err != nil {
			t.Errorf("expected nil for completed task without record field, got: %v", err)
		}
	})
}
