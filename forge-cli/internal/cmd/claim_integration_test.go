package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"
)

func TestCheckExistingTaskState(t *testing.T) {
	t.Run("no state file", func(t *testing.T) {
		dir := t.TempDir()
		indexPath := filepath.Join(dir, "index.json")
		statePath := filepath.Join(dir, "task-state.json")

		index := &task.TaskIndex{}
		if err := task.SaveIndex(indexPath, index); err != nil {
			t.Fatal(err)
		}

		cont, hasIssues, issues := checkExistingTaskState(dir, index, statePath)
		if cont || hasIssues || len(issues) != 0 {
			t.Errorf("expected (false, false, nil), got (%v, %v, %v)", cont, hasIssues, issues)
		}
	})

	t.Run("corrupted state file", func(t *testing.T) {
		dir := t.TempDir()
		indexPath := filepath.Join(dir, "index.json")
		statePath := filepath.Join(dir, "task-state.json")

		index := &task.TaskIndex{}
		if err := task.SaveIndex(indexPath, index); err != nil {
			t.Fatal(err)
		}

		// Create corrupted state file (invalid JSON)
		_ = os.WriteFile(statePath, []byte("not valid json {{{"), 0644)

		cont, hasIssues, issues := checkExistingTaskState(dir, index, statePath)
		if cont || hasIssues || len(issues) != 0 {
			t.Errorf("expected (false, false, nil) for corrupted file, got (%v, %v, %v)", cont, hasIssues, issues)
		}
	})

	t.Run("state with task in_progress", func(t *testing.T) {
		dir := t.TempDir()
		indexPath := filepath.Join(dir, "index.json")
		statePath := filepath.Join(dir, "task-state.json")

		index := &task.TaskIndex{}
		index.SetTasks(map[string]task.Task{
			"task1": {ID: "1.1", Title: "Task 1", Status: "in_progress"},
		})
		if err := task.SaveIndex(indexPath, index); err != nil {
			t.Fatal(err)
		}

		state := &task.TaskState{TaskID: "1.1", Key: "task1"}
		if err := task.SaveState(statePath, state); err != nil {
			t.Fatal(err)
		}

		cont, hasIssues, issues := checkExistingTaskState(dir, index, statePath)
		if !cont || hasIssues || len(issues) != 0 {
			t.Errorf("expected (true, false, nil), got (%v, %v, %v)", cont, hasIssues, issues)
		}
	})

	t.Run("state with task completed", func(t *testing.T) {
		dir := t.TempDir()
		indexPath := filepath.Join(dir, "index.json")
		statePath := filepath.Join(dir, "task-state.json")

		index := &task.TaskIndex{}
		index.SetTasks(map[string]task.Task{
			"task1": {ID: "1.1", Title: "Task 1", Status: "completed"},
		})
		if err := task.SaveIndex(indexPath, index); err != nil {
			t.Fatal(err)
		}

		state := &task.TaskState{TaskID: "1.1", Key: "task1"}
		if err := task.SaveState(statePath, state); err != nil {
			t.Fatal(err)
		}

		cont, hasIssues, issues := checkExistingTaskState(dir, index, statePath)
		if cont || hasIssues || len(issues) != 0 {
			t.Errorf("expected (false, false, nil), got (%v, %v, %v)", cont, hasIssues, issues)
		}

		// State should be deleted
		if _, err := os.Stat(statePath); !os.IsNotExist(err) {
			t.Error("state file should be deleted")
		}
	})

	t.Run("state with task key not in index", func(t *testing.T) {
		dir := t.TempDir()
		indexPath := filepath.Join(dir, "index.json")
		statePath := filepath.Join(dir, "task-state.json")

		index := &task.TaskIndex{}
		if err := task.SaveIndex(indexPath, index); err != nil {
			t.Fatal(err)
		}

		state := &task.TaskState{TaskID: "1.1", Key: "task1"}
		if err := task.SaveState(statePath, state); err != nil {
			t.Fatal(err)
		}

		cont, hasIssues, issues := checkExistingTaskState(dir, index, statePath)
		if cont || !hasIssues || len(issues) != 1 {
			t.Errorf("expected (false, true, 1 issue), got (%v, %v, %v)", cont, hasIssues, issues)
		}
	})

	t.Run("state with blocked task clears state and claims new", func(t *testing.T) {
		dir := t.TempDir()
		indexPath := filepath.Join(dir, "index.json")
		statePath := filepath.Join(dir, "task-state.json")

		index := &task.TaskIndex{}
		index.SetTasks(map[string]task.Task{
			"task1": {ID: "1.1", Title: "Task 1", Status: "blocked"},
		})
		if err := task.SaveIndex(indexPath, index); err != nil {
			t.Fatal(err)
		}

		state := &task.TaskState{TaskID: "1.1", Key: "task1"}
		if err := task.SaveState(statePath, state); err != nil {
			t.Fatal(err)
		}

		cont, hasIssues, issues := checkExistingTaskState(dir, index, statePath)
		if cont || hasIssues || len(issues) != 0 {
			t.Errorf("expected (false, false, nil) for blocked, got (%v, %v, %v)", cont, hasIssues, issues)
		}

		// State should be deleted
		if _, err := os.Stat(statePath); !os.IsNotExist(err) {
			t.Error("state file should be deleted for blocked task")
		}
	})

	t.Run("state with unexpected task status", func(t *testing.T) {
		dir := t.TempDir()
		indexPath := filepath.Join(dir, "index.json")
		statePath := filepath.Join(dir, "task-state.json")

		index := &task.TaskIndex{}
		index.SetTasks(map[string]task.Task{
			"task1": {ID: "1.1", Title: "Task 1", Status: "unknown_status"},
		})
		if err := task.SaveIndex(indexPath, index); err != nil {
			t.Fatal(err)
		}

		state := &task.TaskState{TaskID: "1.1", Key: "task1"}
		if err := task.SaveState(statePath, state); err != nil {
			t.Fatal(err)
		}

		cont, hasIssues, issues := checkExistingTaskState(dir, index, statePath)
		if cont || !hasIssues || len(issues) != 1 {
			t.Errorf("expected (false, true, 1 issue) for truly unexpected status, got (%v, %v, %v)", cont, hasIssues, issues)
		}
	})
}

func TestPrintTaskDetails(t *testing.T) {
	testRoot := filepath.Join("/", "test-root")
	testFeature := "test-feature"

	tests := []struct {
		name         string
		key          string
		task         *task.Task
		wantContains []string
	}{
		{
			name: "basic task",
			key:  "task1",
			task: &task.Task{
				ID:       "1.1",
				Title:    "Test Task",
				Priority: "P0",
				Status:   "in_progress",
				File:     "1.1.md",
				Record:   "records/1.1.md",
			},
			wantContains: []string{
				"TASK_ID: 1.1",
				"FILE: " + filepath.Join(testRoot, feature.GetTaskFile(testFeature, "1.1.md")),
				"FEATURE: test-feature",
			},
		},
		{
			name: "task with scope and booleans",
			key:  "task2",
			task: &task.Task{
				ID:          "2.1",
				Title:       "Full Task",
				Priority:    "P1",
				Status:      "pending",
				File:        "2.1.md",
				Record:      "2.1.md",
				Scope:       "backend",
				Breaking:    true,
				MainSession: true,
			},
			wantContains: []string{
				"SCOPE: backend",
				"BREAKING: true",
				"MAIN_SESSION: true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			// Redirect stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			printTaskDetails(tt.key, tt.task, testRoot, testFeature)

			_ = w.Close()
			os.Stdout = old
			_, _ = buf.ReadFrom(r)

			output := buf.String()
			for _, want := range tt.wantContains {
				if !bytes.Contains([]byte(output), []byte(want)) {
					t.Errorf("printTaskDetails() output missing %q", want)
				}
			}
		})
	}
}

func TestPrintContinueTask(t *testing.T) {
	const testRoot = "/test-root"
	const testFeature = "test-feature"

	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	state := &task.TaskState{

		TaskID:      "1.1",
		Key:         "task1",
		Title:       "Test Task",
		StartedTime: "2026-04-06 10:00",
	}
	tt := &task.Task{
		ID:       "1.1",
		Title:    "Test Task",
		Priority: "P0",
		Status:   "in_progress",
	}

	printContinueTask(state, tt, testRoot, testFeature)

	_ = w.Close()
	os.Stdout = old
	_, _ = buf.ReadFrom(r)

	output := buf.String()
	wantStrings := []string{"ACTION: CONTINUE", "STARTED_AT: 2026-04-06 10:00"}
	for _, want := range wantStrings {
		if !bytes.Contains([]byte(output), []byte(want)) {
			t.Errorf("printContinueTask() output missing %q", want)
		}
	}
}

func TestPrintNewTask(t *testing.T) {
	const testRoot = "/test-root"
	const testFeature = "test-feature"

	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printNewTask("task1", &task.Task{
		ID:       "1.1",
		Title:    "New Task",
		Priority: "P0",
		Status:   "in_progress",
	}, testRoot, testFeature)

	_ = w.Close()
	os.Stdout = old
	_, _ = buf.ReadFrom(r)

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("ACTION: CLAIMED")) {
		t.Errorf("printNewTask() output missing 'ACTION: CLAIMED'")
	}
}

// Integration test for claim command setup
func TestClaimCommand_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("setup and validate task state", func(t *testing.T) {
		dir := t.TempDir()

		// Force FindProjectRoot to resolve to temp dir
		t.Setenv("CLAUDE_PROJECT_DIR", dir)

		// Create feature structure
		featureDir := filepath.Join(dir, feature.FeaturesDir, "test-feature")
		tasksDir := filepath.Join(featureDir, feature.TasksDirName)
		indexPath := filepath.Join(featureDir, feature.IndexFileName)

		if err := os.MkdirAll(tasksDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create index
		index := &task.TaskIndex{

			StatusEnum:   []string{"pending", "in_progress", "completed"},
			PriorityEnum: []string{"P0", "P1", "P2"},
		}
		index.SetTasks(map[string]task.Task{
			"task1": {ID: "1.1", Title: "Task 1", Status: "pending", Priority: "P0", File: "1.1.md"},
			"task2": {ID: "1.2", Title: "Task 2", Status: "pending", Priority: "P1", File: "1.2.md"},
		})
		if err := task.SaveIndex(indexPath, index); err != nil {
			t.Fatal(err)
		}

		// Test findTask
		key, gotTask, err := task.FindTask(index, "1.1")
		if err != nil {
			t.Fatalf("findTask error: %v", err)
		}
		if key != "task1" {
			t.Errorf("expected key 'task1', got %q", key)
		}
		if gotTask.Title != "Task 1" {
			t.Errorf("expected title 'Task 1', got %q", gotTask.Title)
		}

		// Test claimNextTask
		claimedKey, claimedTask, err := claimNextTask(index)
		if err != nil {
			t.Fatalf("claimNextTask error: %v", err)
		}
		if claimedKey != "task1" {
			t.Errorf("expected key 'task1', got %q", claimedKey)
		}
		if claimedTask.Status != "in_progress" {
			t.Errorf("expected status 'in_progress', got %q", claimedTask.Status)
		}
	})
}
