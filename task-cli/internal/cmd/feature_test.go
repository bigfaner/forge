package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"task-cli/pkg/feature"
	"task-cli/pkg/task"
)

// captureOutput captures stdout and stderr during a function execution
func captureOutput(f func() error) (string, error) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, err := os.Pipe()
	if err != nil {
		return "", err
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		return "", err
	}

	os.Stdout = wOut
	os.Stderr = wErr

	outCh := make(chan string)
	errCh := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rOut)
		outCh <- buf.String()
	}()

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rErr)
		errCh <- buf.String()
	}()

	runErr := f()

	wOut.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	stdout := <-outCh
	stderr := <-errCh

	return stdout + stderr, runErr
}

func TestRunFeature_Display(t *testing.T) {
	dir := t.TempDir()

	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create feature with proper structure
	if err := feature.EnsureFeatureDir(dir, "test-feature"); err != nil {
		t.Fatal(err)
	}

	// Create index.json
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test-feature"))
	index := &task.TaskIndex{Feature: "test-feature"}
	indexData, _ := json.Marshal(index)
	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		t.Fatal(err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"feature"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("feature command failed: %v", err)
	}

	if !strings.Contains(output, "test-feature") {
		t.Errorf("expected output to contain 'test-feature', got %q", output)
	}
}

func TestRunFeature_Set(t *testing.T) {
	dir := t.TempDir()

	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	_, err = captureOutput(func() error {
		rootCmd.SetArgs([]string{"feature", "my-new-feature"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("feature set command failed: %v", err)
	}

	// Verify feature directory structure was created
	featureProcessDir := filepath.Join(dir, feature.FeaturesDir, "my-new-feature", feature.TasksDirName, feature.ProcessDirName)
	if _, err := os.Stat(featureProcessDir); os.IsNotExist(err) {
		t.Errorf("feature process directory %s was not created", featureProcessDir)
	}
}

func TestRunFeature_NoFeatureSet(t *testing.T) {
	dir := t.TempDir()

	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create features directory but no feature subdirectory
	featuresDir := filepath.Join(dir, feature.FeaturesDir)
	if err := os.MkdirAll(featuresDir, 0755); err != nil {
		t.Fatal(err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	output, _ := captureOutput(func() error {
		rootCmd.SetArgs([]string{"feature"})
		return rootCmd.Execute()
	})

	if !strings.Contains(output, "FEATURE: (none)") {
		t.Errorf("expected output to contain 'FEATURE: (none)', got %q", output)
	}
}

func TestRunQuery(t *testing.T) {
	dir := t.TempDir()

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
		Tasks: map[string]task.Task{
			"task1": {ID: "1.1", Title: "Task 1", Status: "pending", Priority: "P0", File: "1.1.md", Record: "1.1.md", EstimatedTime: "30m", Dependencies: []string{"1.0"}},
		},
	}

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
	defer os.Chdir(origWd)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"query", "1.1"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("query command failed: %v", err)
	}

	if !strings.Contains(output, "KEY: task1") {
		t.Errorf("expected output to contain 'KEY: task1', got %q", output)
	}
	if !strings.Contains(output, "ID: 1.1") {
		t.Errorf("expected output to contain 'ID: 1.1', got %q", output)
	}
	if !strings.Contains(output, "TITLE: Task 1") {
		t.Errorf("expected output to contain 'TITLE: Task 1', got %q", output)
	}
	if !strings.Contains(output, "DEPENDENCIES: 1.0") {
		t.Errorf("expected output to contain 'DEPENDENCIES: 1.0', got %q", output)
	}
}

func TestRunStatus(t *testing.T) {
	dir := t.TempDir()

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
			StatusEnum:   []string{"pending", "in_progress", "completed", "blocked", "skipped"},
		PriorityEnum: []string{"P0", "P1", "P2"},
		Tasks: map[string]task.Task{
			"task1": {ID: "1.1", Title: "Task 1", Status: "pending", Priority: "P0", File: "1.1.md", Record: "1.1.md"},
		},
	}

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
	defer os.Chdir(origWd)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	_, err = captureOutput(func() error {
		rootCmd.SetArgs([]string{"status", "1.1", "blocked"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("status command failed: %v", err)
	}

	updatedIndex, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatalf("failed to load updated index: %v", err)
	}
	if updatedIndex.Tasks["task1"].Status != "blocked" {
		t.Errorf("expected status 'blocked', got %q", updatedIndex.Tasks["task1"].Status)
	}
}

func TestRunCheck(t *testing.T) {
	dir := t.TempDir()

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
		Tasks: map[string]task.Task{
			"task1": {ID: "1.1", Title: "Task 1", Status: "pending", Priority: "P0", File: "1.1.md", Record: "1.1.md"},
			"task2": {ID: "1.2", Title: "Task 2", Status: "pending", Priority: "P1", File: "1.2.md", Record: "1.2.md", Dependencies: []string{"1.1"}},
		},
	}

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
	defer os.Chdir(origWd)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"check"})
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
