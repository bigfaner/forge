package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"
)

func TestRunClaim(t *testing.T) {
	dir := setupClaimTestProject(t)
	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"claim"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("claim command failed: %v", err)
	}

	if !strings.Contains(output, "ACTION:") {
		t.Errorf("expected output to contain 'ACTION', got %q", output)
	}
	if !strings.Contains(output, "KEY:") {
		t.Errorf("expected output to contain 'KEY:', got %q", output)
	}
}

func TestRunClaim_Continue(t *testing.T) {
	dir := setupClaimTestProject(t)

	// Create task state in new location
	statePath := feature.GetTaskStatePath(dir, "test-feature")
	state := &task.TaskState{
		TaskID:      "1.1",
		Key:         "task1",
		Title:       "Task 1",
		Priority:    "P0",
		StartedTime: time.Now().Format("2006-01-02 15:04"),
	}
	if err := task.SaveState(statePath, state); err != nil {
		t.Fatal(err)
	}

	// Mark task as in_progress in index
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test-feature"))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	t1 := index.TasksMap()["task1"]
	t1.Status = "in_progress"
	index.SetTask("task1", t1)
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"claim"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("claim command failed: %v", err)
	}

	if !strings.Contains(output, "ACTION: CONTINUE") {
		t.Errorf("expected output to contain 'ACTION: CONTINUE', got %q", output)
	}
}

func TestRunValidate(t *testing.T) {
	dir := setupClaimTestProject(t)

	// Add more tasks
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test-feature"))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	index.SetTask("task2", task.Task{ID: "1.2", Title: "Task 2", Priority: "P1", Status: "pending", File: "1.2.md", Record: "1.2.md", Dependencies: []string{"1.1"}, Type: "implementation"})
	index.SetTask("task3", task.Task{ID: "1.3", Title: "Task 3", Priority: "P2", Status: "pending", File: "1.3.md", Record: "1.3.md", Dependencies: []string{"1.1", "1.2"}, Type: "implementation"})
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Create task files
	tasksDir := filepath.Join(dir, "docs", "features", "test-feature", "tasks")
	for i := 2; i <= 3; i++ {
		taskFile := filepath.Join(tasksDir, fmt.Sprintf("1.%d.md", i))
		if err := os.WriteFile(taskFile, []byte("task content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"validate", indexPath})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("validate command failed: %v", err)
	}

	if !strings.Contains(output, "Feature: test-feature") {
		t.Errorf("expected output to contain 'Feature: test-feature', got %q", output)
	}
	if !strings.Contains(output, "PASS") {
		t.Errorf("expected output to contain 'PASS' got %q", output)
	}
}

func TestRunRecord(t *testing.T) {
	dir := setupClaimTestProject(t)

	// Mark task as in_progress
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test-feature"))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	t1 := index.TasksMap()["task1"]
	t1.Status = "in_progress"
	index.SetTask("task1", t1)
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Create task state in new location
	statePath := feature.GetTaskStatePath(dir, "test-feature")
	state := &task.TaskState{
		TaskID:      "1.1",
		Key:         "task1",
		Title:       "Task 1",
		Priority:    "P0",
		StartedTime: time.Now().Format("2006-01-02 15:04"),
	}
	if err := task.SaveState(statePath, state); err != nil {
		t.Fatal(err)
	}

	// Create record data file
	recordData := map[string]any{
		"status":        "completed",
		"summary":       "Test summary",
		"filesCreated":  []string{"file1.go"},
		"filesModified": []string{"file2.go"},
		"keyDecisions":  []string{"decision1"},
		"testsPassed":   1,
		"testsFailed":   0,
		"coverage":      85.5,
		"acceptanceCriteria": []map[string]any{
			{"criterion": "Test criterion", "met": true},
		},
		"notes": "test notes",
	}
	recordDataBytes, err := json.Marshal(recordData)
	if err != nil {
		t.Fatal(err)
	}
	recordDataFile := filepath.Join(dir, "record-data.json")
	if err := os.WriteFile(recordDataFile, recordDataBytes, 0644); err != nil {
		t.Fatal(err)
	}

	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	output, _ := captureOutput(func() error {
		rootCmd.SetArgs([]string{"record", "1.1", "--data", recordDataFile})
		return rootCmd.Execute()
	})

	// Verify record was created
	recordPath := filepath.Join(dir, "docs", "features", "test-feature", "tasks", "records", "1.1.md")
	if _, statErr := os.Stat(recordPath); os.IsNotExist(statErr) {
		t.Errorf("expected record file to be created at %s, output=%q", recordPath, output)
	}
	_ = output // avoid unused variable warning
}

func TestRunHookCleanup(t *testing.T) {
	dir := setupClaimTestProject(t)

	// Mark task as completed in index (cleanup only deletes state for completed tasks)
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test-feature"))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	t1 := index.TasksMap()["task1"]
	t1.Status = "completed"
	index.SetTask("task1", t1)
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Create task state in new location
	statePath := feature.GetTaskStatePath(dir, "test-feature")
	state := &task.TaskState{
		TaskID:      "1.1",
		Key:         "task1",
		Title:       "Task 1",
		Priority:    "P0",
		StartedTime: time.Now().Format("2006-01-02 15:04"),
	}
	if err := task.SaveState(statePath, state); err != nil {
		t.Fatal(err)
	}

	// Verify state exists before cleanup
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Fatal("state.json should exist before cleanup")
	}

	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	// Run cleanup directly - cleanupCompletedTaskState() doesn't call os.Exit
	cleanupCompletedTaskState()

	// Verify state was deleted
	if _, err := os.Stat(statePath); !os.IsNotExist(err) {
		t.Error("state.json should be deleted after cleanup")
	}
}

func TestRunHookPreCommit_Success(t *testing.T) {
	dir := setupClaimTestProject(t)

	// Create task state in new location
	statePath := feature.GetTaskStatePath(dir, "test-feature")
	state := &task.TaskState{
		TaskID:      "1.1",
		Key:         "task1",
		Title:       "Task 1",
		Priority:    "P0",
		StartedTime: time.Now().Format("2006-01-02 15:04"),
	}
	if err := task.SaveState(statePath, state); err != nil {
		t.Fatal(err)
	}

	// Mark task as completed in index
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test-feature"))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	t1 := index.TasksMap()["task1"]
	t1.Status = "completed"
	index.SetTask("task1", t1)
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Create record file
	recordsDir := filepath.Join(dir, "docs", "features", "test-feature", "tasks", "records")
	if err := os.MkdirAll(recordsDir, 0755); err != nil {
		t.Fatal(err)
	}
	recordPath := filepath.Join(recordsDir, "1.1.md")
	if err := os.WriteFile(recordPath, []byte("record content"), 0644); err != nil {
		t.Fatal(err)
	}

	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	// Call preCommit directly - it returns error instead of calling os.Exit
	err = verifyTaskCompletion()
	if err != nil {
		t.Errorf("hook pre-commit should succeed, got error: %v", err)
	}
}

// setupClaimTestProject creates a minimal test project for claim tests
func setupClaimTestProject(t *testing.T) string {
	dir := t.TempDir()

	// Create go.mod
	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Ensure feature directory structure (this creates tasks/index.json and process directory)
	if err := feature.EnsureFeatureDir(dir, "test-feature"); err != nil {
		t.Fatal(err)
	}

	// Create index with task in tasks/index.json
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test-feature"))
	index := &task.TaskIndex{
		Feature:      "test-feature",
		PRD:          "prd/prd-spec.md",
		Design:       "design/tech-design.md",
		StatusEnum:   []string{"pending", "in_progress", "completed"},
		PriorityEnum: []string{"P0", "P1", "P2"},
	}
	index.SetTasks(map[string]task.Task{
		"task1": {ID: "1.1", Title: "Task 1", Priority: "P0", Status: "pending", File: "1.1.md", Record: "records/1.1.md", Type: "implementation"},
	})
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Create task file
	taskFile := filepath.Join(dir, "docs", "features", "test-feature", "tasks", "1.1.md")
	if err := os.WriteFile(taskFile, []byte("task content"), 0644); err != nil {
		t.Fatal(err)
	}

	return dir
}

func TestRunStatus_QueryMode(t *testing.T) {
	dir := setupClaimTestProject(t)
	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	// First claim to set up state
	rootCmd.SetArgs([]string{"claim"})
	_ = rootCmd.Execute()

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"status", "1.1"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("status command failed: %v", err)
	}
	if !strings.Contains(output, "STATUS:") {
		t.Errorf("expected STATUS in output, got %q", output)
	}
}

func TestRunVersion(t *testing.T) {
	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"version"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}
	if !strings.Contains(output, "VERSION:") {
		t.Errorf("expected VERSION in output, got %q", output)
	}
}
