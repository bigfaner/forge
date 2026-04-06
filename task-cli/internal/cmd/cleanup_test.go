package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
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
		dir := setupCompletedTaskProject(t)

		statePath := feature.GetTaskStatePath(dir, "test")
		cleanupCompletedTaskState()

		// Verify state file is deleted
		if _, err := os.Stat(statePath); !os.IsNotExist(err) {
			t.Errorf("state file should be deleted")
		}
	})

	t.Run("keeps state file when task is not completed", func(t *testing.T) {
		dir := setupInProgressTaskProject(t)

		statePath := feature.GetTaskStatePath(dir, "test")
		cleanupCompletedTaskState()

		// Verify state file still exists
		if _, err := os.Stat(statePath); os.IsNotExist(err) {
			t.Errorf("state file should NOT be deleted for in_progress task")
		}
	})

	t.Run("deletes record.json when task is completed", func(t *testing.T) {
		dir := setupCompletedTaskProject(t)

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
}

// setupCompletedTaskProject creates a test project with a completed task.
func setupCompletedTaskProject(t *testing.T) string {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Create feature directory structure
	featureDir := filepath.Join(dir, feature.FeaturesDir, "test", feature.TasksDirName)
	processDir := filepath.Join(featureDir, feature.ProcessDirName)
	os.MkdirAll(processDir, 0755)

	// Create index.json with completed task
	indexPath := filepath.Join(featureDir, feature.IndexFileName)
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "in_progress", "completed"},
		Tasks: map[string]task.Task{
			"task1": {ID: "1.1", Title: "Task 1", Status: "completed"},
		},
	}
	indexData, _ := json.Marshal(index)
	os.WriteFile(indexPath, indexData, 0644)

	// Create state.json
	statePath := feature.GetTaskStatePath(dir, "test")
	state := &task.TaskState{TaskID: "1.1", Key: "task1"}
	stateData, _ := json.Marshal(state)
	os.WriteFile(statePath, stateData, 0644)

	return dir
}

// setupInProgressTaskProject creates a test project with an in_progress task.
func setupInProgressTaskProject(t *testing.T) string {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Create feature directory structure
	featureDir := filepath.Join(dir, feature.FeaturesDir, "test", feature.TasksDirName)
	processDir := filepath.Join(featureDir, feature.ProcessDirName)
	os.MkdirAll(processDir, 0755)

	// Create index.json with in_progress task
	indexPath := filepath.Join(featureDir, feature.IndexFileName)
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "in_progress", "completed"},
		Tasks: map[string]task.Task{
			"task1": {ID: "1.1", Title: "Task 1", Status: "in_progress"},
		},
	}
	indexData, _ := json.Marshal(index)
	os.WriteFile(indexPath, indexData, 0644)

	// Create state.json
	statePath := feature.GetTaskStatePath(dir, "test")
	state := &task.TaskState{TaskID: "1.1", Key: "task1"}
	stateData, _ := json.Marshal(state)
	os.WriteFile(statePath, stateData, 0644)

	return dir
}
