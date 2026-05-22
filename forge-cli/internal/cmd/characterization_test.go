package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	taskpkg "forge-cli/internal/cmd/task"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"
)

// ============================================================================
// Characterization tests for current (possibly buggy) behavior.
// These tests lock in the CURRENT behavior exactly, including gaps.
// Phase 2 will update expectations to reflect behavior changes.
// DO NOT fix bugs in these tests — they are the regression safety net.
// ============================================================================

// ---------------------------------------------------------------------------
// TestSubmit_RejectsCompletedResubmit
//
// Submit on a completed task is rejected by the state machine.
// ValidateTransition(current="completed", target="completed", role=RoleSubmit)
// returns an error because completed is a terminal state.
// ---------------------------------------------------------------------------
func TestSubmit_RejectsCompletedResubmit(t *testing.T) {
	if os.Getenv("TEST_CHAR_SUBMIT_COMPLETED_RESUBMIT") == "1" {
		setupFullProject(t, SetupOpts{
			Tasks: map[string]task.Task{
				"t1": {ID: "1", Title: "Already Done", Status: "completed", Type: task.TypeDoc, File: "1.md", Record: "records/1.md"},
			},
		})

		dir, _ := os.Getwd()

		dataPath := filepath.Join(dir, "record.json")
		rd := task.RecordData{
			Status:      "completed",
			Summary:     "Resubmit on completed task",
			Coverage:    -1.0,
			TestsPassed: 0,
			TestsFailed: 0,
		}
		rdJSON, _ := json.Marshal(rd)
		_ = os.WriteFile(dataPath, rdJSON, 0644)

		*taskpkg.ExportSubmitDataPath = dataPath
		*taskpkg.ExportSubmitJSON = false
		*taskpkg.ExportSubmitQuiet = false
		if err := taskpkg.ExportRunSubmit(taskpkg.ExportSubmitCmd, []string{"1"}); err != nil {
			Exit(err)
		}
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestSubmit_RejectsCompletedResubmit")
	cmd.Env = append(os.Environ(), "TEST_CHAR_SUBMIT_COMPLETED_RESUBMIT=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected error: submit on completed task should be rejected by state machine")
	}
	out := string(output)
	if !strings.Contains(out, "INVALID_TRANSITION") {
		t.Errorf("expected INVALID_TRANSITION error, got: %s", out)
	}
}

// ---------------------------------------------------------------------------
// TestAdd_BlockSource_CurrentBehavior_AllowsCompletedToBlocked
//
// Current behavior: --block-source on a completed task succeeds. The add
// command sets the source task to blocked status even though completed is a
// terminal state. The status command would reject this transition, but add
// bypasses state machine validation.
//
// Desired future behavior (Phase 2): --block-source should validate source
// task state and reject terminal-state transitions.
// ---------------------------------------------------------------------------
func TestAdd_BlockSource_CurrentBehavior_AllowsCompletedToBlocked(t *testing.T) {
	// NOTE: This tests the executeAdd path indirectly through the
	// block-source flag. The current implementation in task.AddTask
	// does not validate that the source task is non-terminal before
	// setting it to blocked. This test documents that gap.

	// We verify the current behavior: block-source on a completed task
	// succeeds by checking the returned SourceBlocked field.
	t.Run("block-source reports blocked even for completed source", func(t *testing.T) {
		dir := setupFullProject(t, SetupOpts{
			Tasks: map[string]task.Task{
				"source": {ID: "1", Title: "Source Task", Status: "completed", Priority: "P0", Type: task.TypeCodingFeature, File: "1.md", Record: "records/1.md"},
			},
		})

		// Reset flag defaults
		*taskpkg.ExportAddTitle = "Fix task"
		*taskpkg.ExportAddPriority = "P0"
		*taskpkg.ExportAddBreaking = true
		*taskpkg.ExportAddID = ""
		*taskpkg.ExportAddDependsOn = ""
		*taskpkg.ExportAddEstimatedTime = ""
		*taskpkg.ExportAddDescription = "Fix something"
		*taskpkg.ExportAddTemplate = ""
		*taskpkg.ExportAddVars = nil
		*taskpkg.ExportAddSourceTaskID = "1"
		*taskpkg.ExportAddBlockSource = true
		*taskpkg.ExportAddType = task.TypeCodingFix

		// executeAdd should succeed — current behavior does not validate
		// that the source task is non-terminal
		result, err := taskpkg.ExportExecuteAdd(nil)
		if err != nil {
			t.Fatalf("CURRENT BEHAVIOR: add --block-source should succeed on completed source. Got error: %v", err)
		}

		// Verify SourceBlocked is reported (current behavior documents the gap)
		if result.SourceBlocked != "1" {
			t.Errorf("expected SourceBlocked='1', got %q", result.SourceBlocked)
		}

		// Check that the index was updated — source task status
		indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
		index, _ := task.LoadIndex(indexPath)

		// After Phase 2, BuildIndex re-keys tasks by filename.
		// The source task has File "1.md" so its key is "1", not "source".
		srcTask, exists := index.ByID("1")
		if !exists {
			t.Fatal("source task not found in index")
		}
		// CURRENT BEHAVIOR: AddTask DOES set source to blocked even though
		// it was completed (a terminal state). This bypasses the state machine
		// guards that the status command enforces. Phase 2 should reject this.
		if srcTask.Status != "blocked" {
			t.Errorf("CURRENT BEHAVIOR: source should be set to blocked by AddTask (--block-source bypasses terminal state guards). Got status: %s", srcTask.Status)
		}
	})
}

// ---------------------------------------------------------------------------
// TestClaim_AutoUnblock_CurrentBehavior
//
// Current behavior: claimNextTask auto-unblocks blocked tasks whose
// dependencies are all met. The auto-unblock sets status to "pending" but
// does NOT clear BlockedReason. The BlockedReason field persists even after
// the task transitions out of blocked state.
//
// Desired future behavior (Phase 2): auto-unblock should clear BlockedReason
// when transitioning to pending.
// ---------------------------------------------------------------------------
func TestClaim_AutoUnblock_CurrentBehavior(t *testing.T) {
	t.Run("auto-unblock preserves BlockedReason (current gap)", func(t *testing.T) {
		index := &task.TaskIndex{
			StatusEnum:   []string{"pending", "in_progress", "completed", "blocked"},
			PriorityEnum: []string{"P0", "P1", "P2"},
		}
		index.SetTasks(map[string]task.Task{
			"task1": {ID: "1", Title: "Dep", Priority: "P0", Status: "completed", Dependencies: []string{}},
			"task2": {
				ID:            "2",
				Title:         "Blocked task",
				Priority:      "P0",
				Status:        "blocked",
				Dependencies:  []string{"1"},
				BlockedReason: "auto-downgrade: testsFailed=2",
			},
		})

		_, _, err := taskpkg.ExportClaimNextTask(index)
		if err != nil {
			t.Fatalf("taskpkg.ExportClaimNextTask() error = %v", err)
		}

		// After auto-unblock, task2 should be in_progress (claimed)
		// CURRENT BEHAVIOR: BlockedReason is NOT cleared
		task2 := index.TasksMap()["task2"]
		if task2.Status != "in_progress" {
			t.Errorf("expected in_progress, got %s", task2.Status)
		}

		// This is the current behavior gap: BlockedReason persists
		if task2.BlockedReason == "" {
			t.Error("CURRENT BEHAVIOR: BlockedReason should persist after auto-unblock (this is a gap). Phase 2 should clear it.")
		}
	})

	t.Run("auto-unblock without BlockedReason works normally", func(t *testing.T) {
		index := &task.TaskIndex{
			StatusEnum:   []string{"pending", "in_progress", "completed", "blocked"},
			PriorityEnum: []string{"P0", "P1", "P2"},
		}
		index.SetTasks(map[string]task.Task{
			"task1": {ID: "1", Title: "Dep", Priority: "P0", Status: "completed", Dependencies: []string{}},
			"task2": {
				ID:           "2",
				Title:        "Blocked task",
				Priority:     "P0",
				Status:       "blocked",
				Dependencies: []string{"1"},
			},
		})

		key, _, err := taskpkg.ExportClaimNextTask(index)
		if err != nil {
			t.Fatalf("taskpkg.ExportClaimNextTask() error = %v", err)
		}
		if key != "task2" {
			t.Errorf("expected key 'task2', got %q", key)
		}
		task2 := index.TasksMap()["task2"]
		if task2.Status != "in_progress" {
			t.Errorf("expected in_progress, got %s", task2.Status)
		}
	})
}

// ---------------------------------------------------------------------------
// TestQualityGate_SourceTaskID_IsEmpty
//
// Phase 2: addFixTask no longer uses SourceTaskID "quality-gate:<step>" sentinel.
// SourceTaskID is now empty. countFixTasks identifies fix-tasks by title prefix
// "fix <step>:" only, and counts active (non-terminal) tasks only.
// ---------------------------------------------------------------------------
func TestQualityGate_SourceTaskID_IsEmpty(t *testing.T) {
	t.Run("countFixTasks matches by title prefix only", func(t *testing.T) {
		index := &task.TaskIndex{Feature: "test"}
		index.SetTasks(map[string]task.Task{
			"fix-1": {
				ID:           "fix-1",
				Title:        "fix compile: compile failure in quality gate",
				SourceTaskID: "",
				Status:       "pending",
				Type:         task.TypeCodingFix,
			},
			"fix-2": {
				ID:           "fix-2",
				Title:        "fix compile: compile failure in quality gate",
				SourceTaskID: "",
				Status:       "in_progress",
				Type:         task.TypeCodingFix,
			},
			"fix-3": {
				ID:           "fix-3",
				Title:        "fix lint: lint failure in quality gate",
				SourceTaskID: "",
				Status:       "pending",
				Type:         task.TypeCodingCleanup,
			},
		})

		// countFixTasks matches by title prefix, not SourceTaskID
		compileCount := countFixTasks(index, "compile")
		if compileCount != 2 {
			t.Errorf("expected 2 compile fix-tasks, got %d", compileCount)
		}

		lintCount := countFixTasks(index, "lint")
		if lintCount != 1 {
			t.Errorf("expected 1 lint fix-task, got %d", lintCount)
		}

		testCount := countFixTasks(index, "unit-test")
		if testCount != 0 {
			t.Errorf("expected 0 unit-test fix-tasks, got %d", testCount)
		}
	})

	t.Run("countFixTasks excludes completed fix-tasks (active-only)", func(t *testing.T) {
		index := &task.TaskIndex{Feature: "test"}
		index.SetTasks(map[string]task.Task{
			"fix-1": {
				ID:           "fix-1",
				Title:        "fix compile: compile failure in quality gate",
				SourceTaskID: "",
				Status:       "completed",
				Type:         task.TypeCodingFix,
			},
			"fix-2": {
				ID:           "fix-2",
				Title:        "fix compile: compile failure in quality gate",
				SourceTaskID: "",
				Status:       "pending",
				Type:         task.TypeCodingFix,
			},
		})

		count := countFixTasks(index, "compile")
		// Active-only: only pending/in_progress/blocked count
		if count != 1 {
			t.Errorf("countFixTasks should count only active fix-tasks (not completed). Got %d", count)
		}
	})
}

// ---------------------------------------------------------------------------
// TestQualityGate_CountFixTasks_CountsActiveOnly
//
// Phase 2: countFixTasks counts only active (non-terminal) fix-tasks for a step.
// It identifies fix-tasks by title prefix "fix <step>:" (no longer requires
// SourceTaskID). Terminal statuses (completed, rejected, skipped) are excluded.
// ---------------------------------------------------------------------------
func TestQualityGate_CountFixTasks_CountsActiveOnly(t *testing.T) {
	tests := []struct {
		name  string
		step  string
		tasks map[string]task.Task
		want  int
	}{
		{
			"excludes completed fix-tasks",
			"compile",
			map[string]task.Task{
				"fix-1": {ID: "fix-1", Title: "fix compile: error", SourceTaskID: "quality-gate:compile", Status: "completed", Type: task.TypeCodingFix},
			},
			0,
		},
		{
			"counts pending fix-tasks",
			"compile",
			map[string]task.Task{
				"fix-1": {ID: "fix-1", Title: "fix compile: error", SourceTaskID: "quality-gate:compile", Status: "pending", Type: task.TypeCodingFix},
			},
			1,
		},
		{
			"counts blocked fix-tasks",
			"compile",
			map[string]task.Task{
				"fix-1": {ID: "fix-1", Title: "fix compile: error", SourceTaskID: "quality-gate:compile", Status: "blocked", Type: task.TypeCodingFix},
			},
			1,
		},
		{
			"excludes skipped fix-tasks",
			"compile",
			map[string]task.Task{
				"fix-1": {ID: "fix-1", Title: "fix compile: error", SourceTaskID: "quality-gate:compile", Status: "skipped", Type: task.TypeCodingFix},
			},
			0,
		},
		{
			"mix of statuses — active only",
			"compile",
			map[string]task.Task{
				"fix-1": {ID: "fix-1", Title: "fix compile: error", SourceTaskID: "quality-gate:compile", Status: "completed", Type: task.TypeCodingFix},
				"fix-2": {ID: "fix-2", Title: "fix compile: error", SourceTaskID: "quality-gate:compile", Status: "pending", Type: task.TypeCodingFix},
				"fix-3": {ID: "fix-3", Title: "fix compile: error", SourceTaskID: "quality-gate:compile", Status: "blocked", Type: task.TypeCodingFix},
			},
			2,
		},
		{
			"excludes different step",
			"compile",
			map[string]task.Task{
				"fix-1": {ID: "fix-1", Title: "fix lint: error", SourceTaskID: "quality-gate:lint", Status: "pending", Type: task.TypeCodingCleanup},
			},
			0,
		},
		{
			"counts tasks without SourceTaskID (matched by title)",
			"compile",
			map[string]task.Task{
				"fix-1": {ID: "fix-1", Title: "fix compile: error", SourceTaskID: "", Status: "pending", Type: task.TypeCodingFix},
			},
			1,
		},
		{
			"excludes tasks without title prefix",
			"compile",
			map[string]task.Task{
				"fix-1": {ID: "fix-1", Title: "some other title", SourceTaskID: "quality-gate:compile", Status: "pending", Type: task.TypeCodingFix},
			},
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := &task.TaskIndex{Feature: "test"}
			index.SetTasks(tt.tasks)
			got := countFixTasks(index, tt.step)
			if got != tt.want {
				t.Errorf("countFixTasks(%q) = %d, want %d", tt.step, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestBuildIndex_Orphan_WarningOnly
//
// Current behavior: BuildIndex detects orphan tasks (in index but no .md file)
// and emits a WARNING but does NOT remove them. The orphan task is preserved
// in the index with PreservedCount incremented.
//
// Desired future behavior (Phase 2): may want to clean up orphans or fail.
// ---------------------------------------------------------------------------
func TestBuildIndex_Orphan_WarningOnly(t *testing.T) {
	projectRoot := t.TempDir()
	featureSlug := "test-feature"
	tasksDir := filepath.Join(projectRoot, "docs", "features", featureSlug, "tasks")
	indexPath := filepath.Join(tasksDir, "index.json")
	_ = os.MkdirAll(tasksDir, 0755)

	// Create initial task .md file and build index
	taskContent := "---\nid: \"1\"\ntitle: \"Exists\"\npriority: P0\ntype: coding.feature\n---\n\n# Task 1\n"
	_ = os.WriteFile(filepath.Join(tasksDir, "1.md"), []byte(taskContent), 0644)

	// Create go.mod for project root detection
	_ = os.WriteFile(filepath.Join(projectRoot, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0644)

	opts := task.BuildIndexOpts{
		FeatureSlug: featureSlug,
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}
	_, err := task.BuildIndex(opts)
	if err != nil {
		t.Fatalf("first build: %v", err)
	}

	// Now remove the .md file to create an orphan
	_ = os.Remove(filepath.Join(tasksDir, "1.md"))

	// Add a different task so the dir isn't empty
	taskContent2 := "---\nid: \"2\"\ntitle: \"New\"\npriority: P0\ntype: coding.feature\n---\n\n# Task 2\n"
	_ = os.WriteFile(filepath.Join(tasksDir, "2.md"), []byte(taskContent2), 0644)

	result, err := task.BuildIndex(opts)
	if err != nil {
		t.Fatalf("rebuild: %v", err)
	}

	// CURRENT BEHAVIOR: orphan is only warned, not cleaned
	found := false
	for _, w := range result.Warnings {
		if strings.HasPrefix(w, "orphan") {
			found = true
		}
	}
	if !found {
		t.Errorf("PHASE 2: expected orphan warning for removed orphan, got %v", result.Warnings)
	}

	if result.PreservedCount != 1 {
		t.Errorf("PHASE 2: expected PreservedCount=1 (orphan removed), got %d", result.PreservedCount)
	}

	// PHASE 2 BEHAVIOR: orphan is removed from the index (not preserved)
	index, _ := task.LoadIndex(indexPath)
	_, exists := index.ByID("1")
	if exists {
		t.Error("PHASE 2: orphan task should be removed from index, not preserved")
	}
}

// ---------------------------------------------------------------------------
// TestStatus_AllowsMutation
//
// Status command is now READ-ONLY. 2-arg mutation calls return an error.
// Non-terminal transitions that were previously allowed are now rejected.
// Use "forge task submit" to complete a task or "forge task reopen" to
// re-activate rejected/skipped tasks.
// ---------------------------------------------------------------------------
func TestStatus_AllowsMutation(t *testing.T) {
	// Status command uses ExactArgs(1), so cobra rejects 2-arg mutation calls.
	// Verify the Args validator rejects extra arguments for various target statuses.
	targets := []string{"blocked", "in_progress", "skipped", "rejected", "completed", "pending"}
	for _, target := range targets {
		t.Run("rejects_"+target, func(t *testing.T) {
			err := taskpkg.StatusCmd.Args(taskpkg.StatusCmd, []string{"1.1", target})
			if err == nil {
				t.Errorf("expected ExactArgs(1) to reject 2nd arg %q", target)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestSubmit_AutoDowngrade_SetsBlockedReason
//
// Auto-downgrade sets status to "blocked" and sets BlockedReason on the task.
// The BlockedReason format is "auto-downgrade: testsFailed=N".
// ---------------------------------------------------------------------------
func TestSubmit_AutoDowngrade_SetsBlockedReason(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {ID: "1", Title: "Task 1", Status: "in_progress", Type: task.TypeCodingFeature, File: "1.md", Record: "records/1.md"},
		},
	})

	rd := task.RecordData{
		Status:       "completed",
		Summary:      "Tests partially pass",
		TestsPassed:  3,
		TestsFailed:  2,
		Coverage:     60.0,
		KeyDecisions: []string{"partial fix"},
	}
	rdJSON, _ := json.Marshal(rd)
	dataPath := filepath.Join(dir, "record.json")
	_ = os.WriteFile(dataPath, rdJSON, 0644)

	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1", Key: "t1", StartedTime: "2026-01-01 10:00"})

	*taskpkg.ExportSubmitDataPath = dataPath
	*taskpkg.ExportSubmitJSON = false
	*taskpkg.ExportSubmitQuiet = false

	_ = captureStdout(func() {
		_ = taskpkg.ExportRunSubmit(nil, []string{"1"})
	})

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index, _ := task.LoadIndex(indexPath)

	task1 := index.TasksMap()["t1"]
	if task1.Status != "blocked" {
		t.Fatalf("expected auto-downgrade to blocked, got %s", task1.Status)
	}

	// BlockedReason should be set after auto-downgrade
	if task1.BlockedReason != "auto-downgrade: testsFailed=2" {
		t.Errorf("expected BlockedReason 'auto-downgrade: testsFailed=2', got %q", task1.BlockedReason)
	}
}
