package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

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
// TestSubmit_CurrentBehavior_AllowsCompletedResubmit
//
// Current behavior: submit on a completed task succeeds without error.
// The submit command does not check whether the task is already in a terminal
// state — it overwrites status unconditionally (BIZ-task-lifecycle-002 gap).
//
// Desired future behavior (Phase 2): submit on a terminal task should fail
// unless --force is used.
// ---------------------------------------------------------------------------
func TestSubmit_CurrentBehavior_AllowsCompletedResubmit(t *testing.T) {
	if os.Getenv("TEST_CHAR_SUBMIT_COMPLETED_RESUBMIT") == "1" {
		setupFullProject(t, SetupOpts{
			Tasks: map[string]task.Task{
				"t1": {ID: "1", Title: "Already Done", Status: "completed", Type: task.TypeDoc, File: "1.md", Record: "records/1.md"},
			},
		})

		dir, _ := os.Getwd()

		// Write existing record so write-once protection triggers without --force
		recordPath := filepath.Join(dir, "docs", "features", "test", "tasks", "records", "1.md")
		_ = os.WriteFile(recordPath, []byte("existing record"), 0644)

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

		submitDataPath = dataPath
		submitForce = true // need --force to overwrite the existing record file
		submitJSON = false
		submitQuiet = false
		runSubmit(submitCmd, []string{"1"})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestSubmit_CurrentBehavior_AllowsCompletedResubmit")
	cmd.Env = append(os.Environ(), "TEST_CHAR_SUBMIT_COMPLETED_RESUBMIT=1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("CURRENT BEHAVIOR: submit on completed task should succeed (this is a known gap). Got error: %v\n%s", err, string(output))
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
				"source": {ID: "1", Title: "Source Task", Status: "completed", Priority: "P0", File: "1.md", Record: "records/1.md"},
			},
		})

		// Reset flag defaults
		addTitle = "Fix task"
		addPriority = "P0"
		addBreaking = true
		addID = ""
		addDependsOn = ""
		addEstimatedTime = ""
		addDescription = "Fix something"
		addTemplate = ""
		addVars = nil
		addSourceTaskID = "1"
		addBlockSource = true
		addType = task.TypeCodingFix

		// executeAdd should succeed — current behavior does not validate
		// that the source task is non-terminal
		result, err := executeAdd(nil)
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

		srcTask, exists := index.ByID("source")
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

		_, _, err := claimNextTask(index)
		if err != nil {
			t.Fatalf("claimNextTask() error = %v", err)
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

		key, _, err := claimNextTask(index)
		if err != nil {
			t.Fatalf("claimNextTask() error = %v", err)
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
// TestQualityGate_SourceTaskID_IsSentinel
//
// Current behavior: addFixTask uses SourceTaskID "quality-gate:<step>" as a
// sentinel value. This is NOT a real task ID — it's a synthetic identifier
// used for cumulative counting via countFixTasks. The Vars["SOURCE_TASK_ID"]
// is set to "N/A (project-wide gate)" for template rendering.
//
// This test verifies that the sentinel format is "quality-gate:<step>" and
// that countFixTasks correctly counts tasks with this sentinel.
// ---------------------------------------------------------------------------
func TestQualityGate_SourceTaskID_IsSentinel(t *testing.T) {
	t.Run("sentinel format is quality-gate:STEP", func(t *testing.T) {
		index := &task.TaskIndex{Feature: "test"}
		index.SetTasks(map[string]task.Task{
			"fix-1": {
				ID:           "fix-1",
				Title:        "fix compile: compile failure in quality gate",
				SourceTaskID: "quality-gate:compile",
				Type:         task.TypeCodingFix,
			},
			"fix-2": {
				ID:           "fix-2",
				Title:        "fix compile: compile failure in quality gate",
				SourceTaskID: "quality-gate:compile",
				Type:         task.TypeCodingFix,
			},
			"fix-3": {
				ID:           "fix-3",
				Title:        "fix lint: lint failure in quality gate",
				SourceTaskID: "quality-gate:lint",
				Type:         task.TypeCodingCleanup,
			},
		})

		// countFixTasks counts ALL fix-tasks (completed + active) for a step
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

	t.Run("countFixTasks counts completed fix-tasks too (not active-only)", func(t *testing.T) {
		index := &task.TaskIndex{Feature: "test"}
		index.SetTasks(map[string]task.Task{
			"fix-1": {
				ID:           "fix-1",
				Title:        "fix compile: compile failure in quality gate",
				SourceTaskID: "quality-gate:compile",
				Status:       "completed",
				Type:         task.TypeCodingFix,
			},
			"fix-2": {
				ID:           "fix-2",
				Title:        "fix compile: compile failure in quality gate",
				SourceTaskID: "quality-gate:compile",
				Status:       "pending",
				Type:         task.TypeCodingFix,
			},
		})

		count := countFixTasks(index, "compile")
		// CURRENT BEHAVIOR: counts ALL fix-tasks regardless of status
		if count != 2 {
			t.Errorf("CURRENT BEHAVIOR: countFixTasks should count all fix-tasks (completed + pending). Got %d", count)
		}
	})
}

// ---------------------------------------------------------------------------
// TestQualityGate_CountFixTasks_CountsAll
//
// Current behavior: countFixTasks counts ALL fix-tasks for a step regardless
// of status (completed + active + blocked + skipped). It identifies fix-tasks
// by having a non-empty SourceTaskID AND a title with prefix "fix <step>:".
//
// Desired future behavior (Phase 2): may want to count only active fix-tasks
// for the cap check.
// ---------------------------------------------------------------------------
func TestQualityGate_CountFixTasks_CountsAll(t *testing.T) {
	tests := []struct {
		name  string
		step  string
		tasks map[string]task.Task
		want  int
	}{
		{
			"counts completed fix-tasks",
			"compile",
			map[string]task.Task{
				"fix-1": {ID: "fix-1", Title: "fix compile: error", SourceTaskID: "quality-gate:compile", Status: "completed", Type: task.TypeCodingFix},
			},
			1,
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
			"counts skipped fix-tasks",
			"compile",
			map[string]task.Task{
				"fix-1": {ID: "fix-1", Title: "fix compile: error", SourceTaskID: "quality-gate:compile", Status: "skipped", Type: task.TypeCodingFix},
			},
			1,
		},
		{
			"mix of statuses",
			"compile",
			map[string]task.Task{
				"fix-1": {ID: "fix-1", Title: "fix compile: error", SourceTaskID: "quality-gate:compile", Status: "completed", Type: task.TypeCodingFix},
				"fix-2": {ID: "fix-2", Title: "fix compile: error", SourceTaskID: "quality-gate:compile", Status: "pending", Type: task.TypeCodingFix},
				"fix-3": {ID: "fix-3", Title: "fix compile: error", SourceTaskID: "quality-gate:compile", Status: "blocked", Type: task.TypeCodingFix},
			},
			3,
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
			"excludes tasks without SourceTaskID",
			"compile",
			map[string]task.Task{
				"fix-1": {ID: "fix-1", Title: "fix compile: error", SourceTaskID: "", Status: "pending", Type: task.TypeCodingFix},
			},
			0,
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
		t.Errorf("CURRENT BEHAVIOR: expected orphan warning, got %v", result.Warnings)
	}

	if result.PreservedCount != 1 {
		t.Errorf("CURRENT BEHAVIOR: expected PreservedCount=1 (orphan preserved), got %d", result.PreservedCount)
	}

	// Verify the orphan is still in the index (not cleaned)
	index, _ := task.LoadIndex(indexPath)
	_, exists := index.ByID("1")
	if !exists {
		t.Error("CURRENT BEHAVIOR: orphan task should be preserved in index, not removed")
	}
}

// ---------------------------------------------------------------------------
// TestStatus_AllowsMutation
//
// Current behavior: the status command allows state changes via the 2-arg
// form for non-terminal transitions. The isTransitionAllowed function blocks:
//   - completed -> anything (without --force)
//   - rejected -> anything (without --force)
//   - anything -> completed (must use submit)
//
// But it allows all other transitions including:
//   - pending -> blocked (skipping in_progress)
//   - in_progress -> pending (regression)
//   - in_progress -> skipped (skipping submit)
//   - blocked -> in_progress (skipping pending)
//
// This test verifies the permissive current behavior for non-terminal transitions.
// ---------------------------------------------------------------------------
func TestStatus_AllowsMutation(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		to      string
		force   bool
		allowed bool
	}{
		// Allowed non-terminal transitions (current permissive behavior)
		{"pending -> blocked", "pending", "blocked", false, true},
		{"pending -> in_progress", "pending", "in_progress", false, true},
		{"pending -> skipped", "pending", "skipped", false, true},
		{"pending -> rejected", "pending", "rejected", false, true},
		{"in_progress -> pending", "in_progress", "pending", false, true},
		{"in_progress -> blocked", "in_progress", "blocked", false, true},
		{"in_progress -> skipped", "in_progress", "skipped", false, true},
		{"in_progress -> rejected", "in_progress", "rejected", false, true},
		{"blocked -> pending", "blocked", "pending", false, true},
		{"blocked -> in_progress", "blocked", "in_progress", false, true},
		{"blocked -> skipped", "blocked", "skipped", false, true},
		{"blocked -> rejected", "blocked", "rejected", false, true},
		{"skipped -> pending", "skipped", "pending", false, true},
		{"skipped -> in_progress", "skipped", "in_progress", false, true},
		{"skipped -> blocked", "skipped", "blocked", false, true},

		// Blocked transitions (terminal guards)
		{"completed -> pending", "completed", "pending", false, false},
		{"completed -> in_progress", "completed", "in_progress", false, false},
		{"rejected -> pending", "rejected", "pending", false, false},
		{"anything -> completed", "in_progress", "completed", false, false},
		{"pending -> completed", "pending", "completed", false, false},

		// Same state (allowed)
		{"pending -> pending", "pending", "pending", false, true},
		{"completed -> completed", "completed", "completed", false, true},

		// Force overrides
		{"completed -> pending with --force", "completed", "pending", true, true},
		{"rejected -> pending with --force", "rejected", "pending", true, true},
		{"in_progress -> completed with --force", "in_progress", "completed", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// isTransitionAllowed(from, to) returns true if the transition is valid.
			// When force=true, the caller skips the check entirely (not handled here).
			// We test isTransitionAllowed directly for the non-force path.
			var got bool
			if tt.force {
				// With --force, the caller bypasses isTransitionAllowed entirely
				got = true
			} else {
				got = isTransitionAllowed(tt.from, tt.to)
			}
			if got != tt.allowed {
				t.Errorf("isTransitionAllowed(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.allowed)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestSubmit_AutoDowngrade_NoBlockedReason
//
// Current behavior (SM-6): auto-downgrade sets status to "blocked" but does
// NOT set BlockedReason on the task. The auto-downgrade happens in
// validateRecordData which mutates rd.Status but has no access to the Task
// struct's BlockedReason field. The index update only sets t.Status = rd.Status.
//
// Desired future behavior (Phase 2): auto-downgrade should set BlockedReason
// to something like "auto-downgrade: N test failures".
// ---------------------------------------------------------------------------
func TestSubmit_AutoDowngrade_NoBlockedReason(t *testing.T) {
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

	submitDataPath = dataPath
	submitJSON = false
	submitQuiet = false
	submitForce = false

	_ = captureStdout(func() {
		runSubmit(nil, []string{"1"})
	})

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index, _ := task.LoadIndex(indexPath)

	task1 := index.TasksMap()["t1"]
	if task1.Status != "blocked" {
		t.Fatalf("expected auto-downgrade to blocked, got %s", task1.Status)
	}

	// CURRENT BEHAVIOR (SM-6): BlockedReason is NOT set after auto-downgrade
	if task1.BlockedReason != "" {
		t.Errorf("CURRENT BEHAVIOR: BlockedReason should be empty after auto-downgrade (SM-6 gap). Got: %q", task1.BlockedReason)
	}
}
