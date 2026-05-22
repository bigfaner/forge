package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	taskpkg "forge-cli/internal/cmd/task"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"

	"github.com/stretchr/testify/assert"
)

func TestRunQuery(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	featureDir := filepath.Join(dir, "docs", "features", "test-feature")
	tasksDir := filepath.Join(featureDir, "tasks")
	indexPath := filepath.Join(tasksDir, "index.json")

	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatal(err)
	}

	index := &task.TaskIndex{
		Feature:      "test-feature",
		PRD:          "prd/prd-spec.md",
		Design:       "design/tech-design.md",
		StatusEnum:   []string{"pending", "in_progress", "completed"},
		PriorityEnum: []string{"P0", "P1", "P2"},
	}
	index.SetTasks(map[string]task.Task{
		"task1": {ID: "1.1", Title: "Task 1", Status: "pending", Priority: "P0", File: "1.1.md", Record: "1.1.md", EstimatedTime: "30m", Dependencies: []string{"1.0"}},
	})

	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Ensure feature directory structure exists
	if err := feature.EnsureFeatureDir(dir, "test-feature"); err != nil {
		t.Fatal(err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"task", "query", "1.1"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("query command failed: %v", err)
	}

	if !strings.Contains(output, "TASK_ID: 1.1") {
		t.Errorf("expected output to contain 'TASK_ID: 1.1', got %q", output)
	}
	if !strings.Contains(output, "STATUS: pending") {
		t.Errorf("expected output to contain 'STATUS: pending', got %q", output)
	}
	// Removed fields must NOT appear
	if strings.Contains(output, "KEY:") {
		t.Errorf("KEY should not appear in query output, got %q", output)
	}
	if strings.Contains(output, "TITLE:") {
		t.Errorf("TITLE should not appear in query output, got %q", output)
	}
	if strings.Contains(output, "DEPENDENCIES:") {
		t.Errorf("DEPENDENCIES should not appear in query output, got %q", output)
	}
}

func TestRunStatus(t *testing.T) {
	// Status command uses ExactArgs(1), cobra rejects 2-arg mutation calls.
	err := taskpkg.StatusCmd.Args(taskpkg.StatusCmd, []string{"1.1", "blocked"})
	if err == nil {
		t.Error("expected ExactArgs(1) to reject 2 arguments")
	}
}

func TestRunCheck(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	featureDir := filepath.Join(dir, "docs", "features", "test-feature")
	tasksDir := filepath.Join(featureDir, "tasks")
	indexPath := filepath.Join(tasksDir, "index.json")

	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatal(err)
	}

	index := &task.TaskIndex{
		Feature:      "test-feature",
		PRD:          "prd/prd-spec.md",
		Design:       "design/tech-design.md",
		StatusEnum:   []string{"pending", "in_progress", "completed"},
		PriorityEnum: []string{"P0", "P1", "P2"},
	}
	index.SetTasks(map[string]task.Task{
		"task1": {ID: "1.1", Title: "Task 1", Status: "pending", Priority: "P0", File: "1.1.md", Record: "1.1.md"},
		"task2": {ID: "1.2", Title: "Task 2", Status: "pending", Priority: "P1", File: "1.2.md", Record: "1.2.md", Dependencies: []string{"1.1"}},
	})

	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Ensure feature directory structure exists
	if err := feature.EnsureFeatureDir(dir, "test-feature"); err != nil {
		t.Fatal(err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"task", "check-deps"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("check command failed: %v", err)
	}

	if !strings.Contains(output, "[TASKS]") {
		t.Errorf("expected output to contain '[TASKS]', got %q", output)
	}
	if !strings.Contains(output, "[DEPENDENCIES]") {
		t.Errorf("expected output to contain '[DEPENDENCIES]', got %q", output)
	}
	if !strings.Contains(output, "RESULT: PASS") {
		t.Errorf("expected output to contain 'RESULT: PASS', got %q", output)
	}
}

// Suppress unused imports.
var _ = json.Marshal
var _ assert.Assertions
