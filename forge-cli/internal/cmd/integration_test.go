package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"
)

// SetupOpts configures the test project created by setupFullProject.
type SetupOpts struct {
	// Tasks is the task map to write into index.json (required).
	Tasks map[string]task.Task
	// State, if non-nil, creates a task-state.json in the process directory.
	State *task.TaskState
	// UseEnvVar, when true, sets CLAUDE_PROJECT_DIR instead of using go.mod+chdir+SetFeature.
	UseEnvVar bool
	// FeatureName defaults to "test" if empty.
	FeatureName string
}

// setupFullProject creates a fully configured test project.
//
// By default (UseEnvVar=false) it creates go.mod, feature dirs, index.json,
// task markdown files, records dir, then chdirs and calls feature.SetFeature.
//
// When UseEnvVar=true it instead sets CLAUDE_PROJECT_DIR (no go.mod, no chdir),
// suitable for tests that call project.FindProjectRoot via env-var path.
func setupFullProject(t *testing.T, opts SetupOpts) (dir string) {
	t.Helper()
	dir = t.TempDir()

	featureName := opts.FeatureName
	if featureName == "" {
		featureName = "test"
	}

	if opts.UseEnvVar {
		t.Setenv("CLAUDE_PROJECT_DIR", dir)
	} else {
		// go.mod marks project root
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	if err := feature.EnsureFeatureDir(dir, featureName); err != nil {
		t.Fatal(err)
	}

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile(featureName))
	index := &task.TaskIndex{
		Feature:      featureName,
		PRD:          "prd/prd-spec.md",
		Design:       "design/tech-design.md",
		StatusEnum:   []string{"pending", "in_progress", "completed", "blocked", "skipped"},
		PriorityEnum: []string{"P0", "P1", "P2"},
	}
	if len(opts.Tasks) > 0 {
		index.SetTasks(opts.Tasks)
	}
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Create task markdown files
	tasksDir := filepath.Join(dir, feature.GetFeatureTasksDir(featureName))
	for _, t2 := range opts.Tasks {
		if t2.File != "" {
			if err := os.WriteFile(filepath.Join(tasksDir, t2.File), []byte("# "+t2.Title), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	// Create records dir
	if err := os.MkdirAll(filepath.Join(tasksDir, "records"), 0755); err != nil {
		t.Fatal(err)
	}

	// Optionally create task state file
	if opts.State != nil {
		statePath := feature.GetTaskStatePath(dir, featureName)
		if err := task.SaveState(statePath, opts.State); err != nil {
			t.Fatalf("SaveState failed: %v", err)
		}
	}

	if !opts.UseEnvVar {
		// Set working dir
		origWd, _ := os.Getwd()
		t.Cleanup(func() { _ = os.Chdir(origWd) })
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}

		// Set feature
		if err := feature.SetFeature(dir, featureName); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// ---------- verifyTaskCompletion ----------

func TestVerifyTaskCompletion_HappyPath(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed", File: "1.1.md", Record: "records/1.1.md"},
	}})

	// Write record file
	dir, _ := os.Getwd()
	_ = os.MkdirAll(filepath.Join(dir, "docs", "features", "test", "tasks", "records"), 0755)
	_ = os.WriteFile(filepath.Join(dir, "docs", "features", "test", "tasks", "records", "1.1.md"), []byte("record"), 0644)

	// Save task state
	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

	err := verifyTaskCompletion()
	if err != nil {
		t.Errorf("expected nil, got: %v", err)
	}
}

func TestVerifyTaskCompletion_TaskNotCompleted(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "in_progress", File: "1.1.md", Record: "records/1.1.md"},
	}})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

	err := verifyTaskCompletion()
	if err == nil {
		t.Error("expected error for non-completed task")
	}
	if !strings.Contains(err.Error(), "not completed") {
		t.Errorf("error should mention not completed: %v", err)
	}
}

func TestVerifyTaskCompletion_RecordFileMissing(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed", File: "1.1.md", Record: "records/1.1.md"},
	}})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

	// Don't create the record file
	err := verifyTaskCompletion()
	if err == nil {
		t.Error("expected error for missing record file")
	}
	if !strings.Contains(err.Error(), "record file missing") {
		t.Errorf("error should mention missing record: %v", err)
	}
}

func TestVerifyTaskCompletion_NoProject(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(tmpDir)

	err := verifyTaskCompletion()
	if err != nil {
		t.Errorf("no project should return nil, got: %v", err)
	}
}

func TestVerifyTaskCompletion_NoFeature(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	_ = os.MkdirAll(filepath.Join(dir, "docs", "features"), 0755)

	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	err := verifyTaskCompletion()
	if err != nil {
		t.Errorf("no feature should return nil, got: %v", err)
	}
}

func TestVerifyTaskCompletion_NoState(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed"},
	}})

	err := verifyTaskCompletion()
	if err != nil {
		t.Errorf("no state should return nil, got: %v", err)
	}
}

func TestVerifyTaskCompletion_TaskNotFoundInIndex(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed"},
	}})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "9.9", Key: "missing"})

	err := verifyTaskCompletion()
	if err == nil {
		t.Error("expected error for task not in index")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should mention not found: %v", err)
	}
}

func TestVerifyTaskCompletion_NoRecordField(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed", File: "1.1.md"},
	}})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

	// Task has empty Record — should pass (no file to check)
	err := verifyTaskCompletion()
	if err != nil {
		t.Errorf("empty record field should pass, got: %v", err)
	}
}

// ---------- cleanupCompletedTaskState ----------

func TestCleanupCompletedTaskState_Completed(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed", File: "1.1.md", Record: "records/1.1.md"},
	}})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

	// Also create record.json
	recordPath := feature.GetProcessRecordPath(dir, "test")
	_ = os.MkdirAll(filepath.Dir(recordPath), 0755)
	_ = os.WriteFile(recordPath, []byte("{}"), 0644)

	cleanupCompletedTaskState()

	if _, err := os.Stat(statePath); !os.IsNotExist(err) {
		t.Error("state.json should be deleted for completed task")
	}
	if _, err := os.Stat(recordPath); !os.IsNotExist(err) {
		t.Error("record.json should be deleted for completed task")
	}
}

func TestCleanupCompletedTaskState_InProgress(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "in_progress", File: "1.1.md"},
	}})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

	cleanupCompletedTaskState()

	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Error("state.json should NOT be deleted for in_progress task")
	}
}

func TestCleanupCompletedTaskState_NoState(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed"},
	}})

	// No state file — should not panic
	cleanupCompletedTaskState()
}

func TestCleanupCompletedTaskState_NoProject(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(tmpDir)

	// Should not panic
	cleanupCompletedTaskState()
}

func TestCleanupCompletedTaskState_TaskKeyNotFound(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed"},
	}})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	// State references a key that doesn't exist in index
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "9.9", Key: "nonexistent"})

	cleanupCompletedTaskState()

	// Should not delete state when key doesn't match
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Error("state.json should NOT be deleted when task key not found")
	}
}

// ---------- runSubmit integration ----------

func TestRunRecord_HappyPath(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "Task One", Status: "in_progress", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	}})

	dir, _ := os.Getwd()

	// Create a record data file
	rd := task.RecordData{
		Status:       "completed",
		Summary:      "Did the thing",
		TestsPassed:  5,
		TestsFailed:  0,
		Coverage:     90.0,
		KeyDecisions: []string{"used approach X"},
		AcceptanceCriteria: []task.AcceptanceCriterion{
			{Criterion: "It works", Met: true},
		},
	}
	rdJSON, _ := json.Marshal(rd)
	dataPath := filepath.Join(dir, "record.json")
	_ = os.WriteFile(dataPath, rdJSON, 0644)

	// Save state for startedTime
	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1", StartedTime: "2026-01-01 10:00"})

	submitDataPath = dataPath
	submitJSON = false
	submitQuiet = false
	submitForce = false

	out := captureStdout(func() {
		runSubmit(nil, []string{"1.1"})
	})

	if !strings.Contains(out, "STATUS: completed") {
		t.Errorf("expected status in output, got: %s", out)
	}
	// TASK_ID removed from non-JSON output
	if strings.Contains(out, "TASK_ID:") {
		t.Errorf("TASK_ID should not appear in non-JSON submit output, got: %s", out)
	}

	// Verify record file was created
	recordFile := filepath.Join(dir, "docs", "features", "test", "tasks", "records", "1.1.md")
	if _, err := os.Stat(recordFile); os.IsNotExist(err) {
		t.Error("record file should exist")
	}
}

func TestRunRecord_JSONOutput(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "Task One", Status: "in_progress", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	}})

	dir, _ := os.Getwd()

	rd := task.RecordData{
		Status:      "completed",
		Summary:     "JSON test",
		TestsPassed: 1,
		Coverage:    80.0,
		AcceptanceCriteria: []task.AcceptanceCriterion{
			{Criterion: "Works", Met: true},
		},
	}
	rdJSON, _ := json.Marshal(rd)
	dataPath := filepath.Join(dir, "record.json")
	_ = os.WriteFile(dataPath, rdJSON, 0644)

	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1", StartedTime: "2026-01-01 10:00"})

	submitDataPath = dataPath
	submitJSON = true
	submitQuiet = false
	submitForce = false

	out := captureStdout(func() {
		runSubmit(nil, []string{"1.1"})
	})

	if !strings.Contains(out, `"recordFile"`) {
		t.Errorf("expected JSON output with recordFile, got: %s", out)
	}
	if !strings.Contains(out, `"taskId"`) || !strings.Contains(out, `"1.1"`) {
		t.Errorf("expected JSON output with taskId 1.1, got: %s", out)
	}
}

func TestRunRecord_QuietOutput(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "Task One", Status: "in_progress", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	}})

	dir, _ := os.Getwd()

	rd := task.RecordData{
		Status:      "completed",
		Summary:     "Quiet test",
		TestsPassed: 1,
		Coverage:    75.0,
		AcceptanceCriteria: []task.AcceptanceCriterion{
			{Criterion: "Works", Met: true},
		},
	}
	rdJSON, _ := json.Marshal(rd)
	dataPath := filepath.Join(dir, "record.json")
	_ = os.WriteFile(dataPath, rdJSON, 0644)

	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1", StartedTime: "2026-01-01 10:00"})

	submitDataPath = dataPath
	submitJSON = false
	submitQuiet = true
	submitForce = false

	out := captureStdout(func() {
		runSubmit(nil, []string{"1.1"})
	})

	if strings.Contains(out, "TASK_ID") {
		t.Errorf("quiet mode should not print block output, got: %s", out)
	}
}

// ---------- executeClaim error paths ----------

func TestExecuteClaim_DataIntegrityError(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	}})

	dir, _ := os.Getwd()
	// Create state pointing to a key that doesn't exist in index
	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "9.9", Key: "nonexistent", StartedTime: "2026-01-01 10:00"})

	_, err := executeClaim()
	if err == nil {
		t.Error("expected data integrity error")
	}
	if !strings.Contains(err.Error(), "integrity") {
		t.Errorf("error should mention integrity: %v", err)
	}
}

func TestExecuteClaim_CompletedStateClaimNew(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed", Priority: "P0", File: "1.1.md"},
		"t2": {ID: "1.2", Title: "T2", Status: "pending", Priority: "P0", File: "1.2.md", Record: "records/1.2.md"},
	}})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1", StartedTime: "2026-01-01 10:00"})

	result, err := executeClaim()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Action != "CLAIMED" {
		t.Errorf("expected CLAIMED after completed state, got %q", result.Action)
	}
	if result.Key != "t2" {
		t.Errorf("expected t2, got %q", result.Key)
	}
}

func TestExecuteClaim_BlockedTaskClearsStateAndProceeds(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "blocked", Priority: "P0", File: "1.1.md"},
	}})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1", StartedTime: "2026-01-01 10:00"})

	_, err := executeClaim()
	// Blocked task clears state, but no pending tasks to claim
	if err == nil {
		t.Error("expected error (no pending tasks)")
	}
	if strings.Contains(err.Error(), "integrity") {
		t.Errorf("blocked status should not trigger integrity error: %v", err)
	}
}

// ---------- runValidateIndex direct validator ----------

func TestValidatorRun_WithFileArg(t *testing.T) {
	dir := t.TempDir()

	// Create a valid index.json
	index := &task.TaskIndex{
		Feature:    "test-feature",
		StatusEnum: []string{"pending", "in_progress", "completed"},
	}
	data, _ := json.Marshal(index)
	indexPath := filepath.Join(dir, "index.json")
	_ = os.WriteFile(indexPath, data, 0644)

	out := captureStdout(func() {
		v := &validator{filePath: indexPath}
		if err := v.run(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS, got: %s", out)
	}
}

func TestValidatorRun_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "index.json")
	_ = os.WriteFile(indexPath, []byte("not json"), 0644)

	v := &validator{filePath: indexPath}
	err := v.run()
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

// ---------- fillRecordTemplate ----------

func TestFillRecordTemplate_NonCompletedStatus(t *testing.T) {
	t2 := &task.Task{ID: "1.1", Title: "Test Task"}
	rd := &task.RecordData{
		Status:       "blocked",
		Summary:      "Blocked due to X",
		KeyDecisions: []string{"Decision A"},
		TestsPassed:  0,
		TestsFailed:  0,
		Coverage:     -1.0,
	}

	content := fillRecordTemplate(t2, rd, "2026-01-01 10:00")
	if !strings.Contains(content, `status: "blocked"`) {
		t.Errorf("expected blocked status, got: %s", content)
	}
	if !strings.Contains(content, "N/A") {
		t.Errorf("expected N/A for non-completed completion time, got: %s", content)
	}
	if !strings.Contains(content, "N/A (task has no tests)") {
		t.Errorf("expected N/A for negative coverage, got: %s", content)
	}
}

func TestFillRecordTemplate_WithNotes(t *testing.T) {
	t2 := &task.Task{ID: "1.1", Title: "Test Task"}
	rd := &task.RecordData{
		Status:      "completed",
		Summary:     "Done",
		Notes:       "Custom notes here",
		TestsPassed: 1,
		Coverage:    50.0,
		AcceptanceCriteria: []task.AcceptanceCriterion{
			{Criterion: "Works", Met: true},
		},
	}

	content := fillRecordTemplate(t2, rd, "")
	if !strings.Contains(content, "Custom notes here") {
		t.Errorf("expected custom notes, got: %s", content)
	}
	// Verify notes section has the custom notes, not default 无
	notesSection := content[strings.LastIndex(content, "## Notes"):]
	if !strings.Contains(notesSection, "Custom notes here") {
		t.Errorf("notes section should have custom notes: %s", notesSection)
	}
	if strings.Contains(notesSection, "无") {
		t.Errorf("notes section should not have default value: %s", notesSection)
	}
}

// ---------- saveIndexAndSignalCompletion ----------

func TestSaveIndexAndSignalCompletion_IncompleteTasks(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	_ = feature.EnsureFeatureDir(dir, "test")

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "in_progress", "completed"},
	}
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1.1", Status: "completed"},
		"t2": {ID: "1.2", Status: "pending"},
	})
	_ = task.SaveIndex(indexPath, index)

	// Should NOT write forge state since not all tasks are done
	saveIndexAndSignalCompletion(indexPath, dir, "test", index)

	forgeState := feature.ReadForgeState(dir)
	if forgeState != nil {
		t.Error("forge state should NOT be written when tasks are incomplete")
	}
}

// ---------- validateRecordData ----------

func TestValidateRecordData_ForceOverride(t *testing.T) {
	rd := &task.RecordData{
		Status:      "completed",
		Summary:     "Done",
		TestsPassed: 0,
		TestsFailed: 0,
		Coverage:    50.0,
		AcceptanceCriteria: []task.AcceptanceCriterion{
			{Criterion: "Works", Met: false},
		},
	}

	// Should not exit when force=true
	out := captureStderr2(func() {
		validateRecordData(rd, true)
	})
	if strings.Contains(out, "ERROR") {
		t.Errorf("force should suppress validation errors, got: %s", out)
	}
}

// ---------- validateRecordData no-test task ----------

func TestValidateRecordData_NoTestTask(t *testing.T) {
	rd := &task.RecordData{
		Status:       "completed",
		Summary:      "Docs only",
		Coverage:     -1.0,
		KeyDecisions: []string{"doc-only"},
		AcceptanceCriteria: []task.AcceptanceCriterion{
			{Criterion: "Docs written", Met: true},
		},
	}

	out := captureStderr2(func() {
		validateRecordData(rd, false)
	})
	if strings.Contains(out, "ERROR") {
		t.Errorf("coverage=-1.0 should pass for no-test tasks, got: %s", out)
	}
}

// ---------- runValidateIndex no file arg, feature-based ----------

func TestValidatorRun_FeatureBased(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	_ = feature.EnsureFeatureDir(dir, "test")

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "completed"},
	}
	_ = task.SaveIndex(indexPath, index)

	_ = feature.SetFeature(dir, "test")

	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	out := captureStdout(func() {
		v := &validator{filePath: indexPath}
		if err := v.run(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS: %s", out)
	}
}

// ---------- validateTTest1Template ----------

func TestValidateTTest1Template_UnresolvedPlaceholder(t *testing.T) {
	dir := t.TempDir()
	taskFile := filepath.Join(dir, "T-test-1.md")
	_ = os.WriteFile(taskFile, []byte("# Task\nReplace {{LAST_BUSINESS_TASK_ID}} with actual ID\n"), 0644)

	v := &validator{}
	v.validateFirstTestTaskTemplate(taskFile, "T-test-1", []string{"{{LAST_BUSINESS_TASK_ID}}"})
	if len(v.errors) == 0 {
		t.Error("expected error for unresolved placeholder")
	}
	if !strings.Contains(v.errors[0], "{{LAST_BUSINESS_TASK_ID}}") {
		t.Errorf("error should mention placeholder: %s", v.errors[0])
	}
}

func TestValidateTTest1Template_ResolvedPlaceholder(t *testing.T) {
	dir := t.TempDir()
	taskFile := filepath.Join(dir, "T-test-1.md")
	_ = os.WriteFile(taskFile, []byte("# Task\nDepends on 1.5\n"), 0644)

	v := &validator{}
	v.validateFirstTestTaskTemplate(taskFile, "T-test-1", []string{"{{LAST_BUSINESS_TASK_ID}}"})
	if len(v.errors) != 0 {
		t.Errorf("expected no errors, got: %v", v.errors)
	}
}

// ---------- findTask ----------

func TestFindTaskByKey(t *testing.T) {
	index := &task.TaskIndex{}
	index.SetTasks(map[string]task.Task{
		"task1": {ID: "1.1", Title: "Task One"},
	})

	key, t2, err := task.FindTask(index, "task1")
	if err != nil {
		t.Fatal(err)
	}
	if key != "task1" {
		t.Errorf("key = %q, want task1", key)
	}
	_ = t2
}

// ---------- readSubmitData ----------

func TestReadRecordData_FromFile(t *testing.T) {
	dir := t.TempDir()
	rd := task.RecordData{Summary: "test summary", TestsPassed: 1, Coverage: 50.0}
	data, _ := json.Marshal(rd)
	path := filepath.Join(dir, "record.json")
	_ = os.WriteFile(path, data, 0644)

	result, err := readSubmitData(path)
	if err != nil {
		t.Fatal(err)
	}
	if result.Summary != "test summary" {
		t.Errorf("Summary = %q, want 'test summary'", result.Summary)
	}
}

func TestReadRecordData_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "record.json")
	_ = os.WriteFile(path, []byte("not json"), 0644)

	_, err := readSubmitData(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

// ---------- WarnMissingFields ----------

func TestWarnMissingFields(t *testing.T) {
	out := captureStderr2(func() {
		WarnMissingFields([]string{"keyDecisions", "acceptanceCriteria"})
	})
	if !strings.Contains(out, "WARNING") {
		t.Errorf("expected warning output, got: %s", out)
	}
	if !strings.Contains(out, "keyDecisions") {
		t.Errorf("expected field name in warning, got: %s", out)
	}
}

// ---------- parseSegment default branch ----------

func TestParseSegment_DefaultAlphabetic(t *testing.T) {
	// Unknown alphabetic segment should return 0, false
	val, isNum := parseSegment([]string{"unknown"}, 0)
	if isNum {
		t.Error("expected non-numeric for 'unknown'")
	}
	if val != 0 {
		t.Errorf("val = %d, want 0", val)
	}
}

func TestParseSegment_OutOfRange(t *testing.T) {
	val, isNum := parseSegment([]string{"1"}, 1)
	if !isNum {
		t.Error("expected numeric for missing segment")
	}
	if val != -1 {
		t.Errorf("val = %d, want -1", val)
	}
}

// ---------- printTaskDetails with Breaking ----------

func TestPrintTaskDetails_Breaking(t *testing.T) {
	t2 := &task.Task{
		ID:            "2.gate",
		Title:         "Gate Task",
		Priority:      "P0",
		Status:        "pending",
		Breaking:      true,
		File:          "2.gate.md",
		Record:        "records/2.gate.md",
		EstimatedTime: "30min",
		Dependencies:  []string{"1.summary"},
	}

	out := captureStdout(func() {
		printTaskDetails("gate-2", t2, "/project", "test")
	})
	if !strings.Contains(out, "BREAKING: true") {
		t.Errorf("expected BREAKING field, got: %s", out)
	}
	if !strings.Contains(out, "FEATURE: test") {
		t.Errorf("expected FEATURE: test, got: %s", out)
	}
	if !strings.Contains(out, "TASK_ID: 2.gate") {
		t.Errorf("expected TASK_ID: 2.gate, got: %s", out)
	}
}

// ---------- runStatus update mode ----------

func TestRunStatus_Update(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "Status Task", Status: "pending", Priority: "P0", File: "1.1.md", Record: "records/1.1.md", Dependencies: []string{}},
	}})

	out := captureStdout(func() {
		runStatus(nil, []string{"1.1", "blocked"})
	})
	if !strings.Contains(out, "STATUS: blocked") {
		t.Errorf("expected updated status, got: %s", out)
	}

	// Verify index was updated
	dir, _ := os.Getwd()
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index, _ := task.LoadIndex(indexPath)
	if index.TasksMap()["t1"].Status != "blocked" {
		t.Errorf("index status = %q, want blocked", index.TasksMap()["t1"].Status)
	}
}

// ---------- executeClaim error: no project ----------

func TestExecuteClaim_NoProject(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(tmpDir)

	_, err := executeClaim()
	if err == nil {
		t.Error("expected error for no project root")
	}
}

// ---------- executeClaim: save index error ----------

func TestExecuteClaim_SaveIndexError(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	_ = feature.EnsureFeatureDir(dir, "test")

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:      "test",
		StatusEnum:   []string{"pending", "in_progress", "completed"},
		PriorityEnum: []string{"P0", "P1", "P2"},
	}
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "pending", Priority: "P0", File: "1.1.md", Record: "1.1.md"},
	})
	_ = task.SaveIndex(indexPath, index)

	_ = feature.SetFeature(dir, "test")
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Make index.json read-only so task.SaveIndex (os.WriteFile) fails
	_ = os.Chmod(indexPath, 0444)
	defer func() { _ = os.Chmod(indexPath, 0644) }()

	_, err := executeClaim()
	if err == nil {
		t.Error("expected error when save index fails")
	}
}

// ---------- runClaim output paths ----------

func TestRunClaim_Output(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "Claim Task", Status: "pending", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	}})

	out := captureStdout(func() {
		runClaim(nil, []string{})
	})
	if !strings.Contains(out, "ACTION: CLAIMED") {
		t.Errorf("expected CLAIMED output, got: %s", out)
	}
	if !strings.Contains(out, "ID: 1.1") {
		t.Errorf("expected task ID in output, got: %s", out)
	}
}

// ---------- runCheckDeps integration (valid deps, exits 0 via PrintResult) ----------

func TestRunCheck_AllValid(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t0": {ID: "1.0", Title: "T0", Status: "completed", Dependencies: []string{}},
		"t1": {ID: "1.1", Title: "T1", Status: "pending", Dependencies: []string{"1.0"}},
	}})

	out := captureStdout(func() {
		captureStderr2(func() {
			runCheckDeps(nil, []string{})
		})
	})
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS for valid deps, got: %s", out)
	}
	if !strings.Contains(out, "TASKS") {
		t.Errorf("expected TASKS section, got: %s", out)
	}
	if !strings.Contains(out, "DEPENDENCIES") {
		t.Errorf("expected DEPENDENCIES section, got: %s", out)
	}
}

// ---------- runValidateIndex with explicit file arg ----------

func TestRunValidate_Integration(t *testing.T) {
	dir := t.TempDir()

	index := &task.TaskIndex{
		Feature:    "my-feature",
		PRD:        "prd/prd-spec.md",
		Design:     "design/tech-design.md",
		StatusEnum: []string{"pending", "in_progress", "completed"},
	}
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "pending", Priority: "P0", File: "1.1.md", Dependencies: []string{}, Type: "implementation"},
	})
	data, _ := json.Marshal(index)
	indexPath := filepath.Join(dir, "index.json")
	_ = os.WriteFile(indexPath, data, 0644)

	// Create tasks dir and task file
	tasksDir := filepath.Join(dir, "tasks")
	_ = os.MkdirAll(tasksDir, 0755)
	_ = os.WriteFile(filepath.Join(tasksDir, "1.1.md"), []byte("# T1"), 0644)

	out := captureStdout(func() {
		v := &validator{filePath: indexPath}
		if err := v.run(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS, got: %s", out)
	}
	if !strings.Contains(out, "Feature: my-feature") {
		t.Errorf("expected feature info, got: %s", out)
	}
}

// ---------- saveIndexAndSignalCompletion with forge state ----------

func TestSaveIndexAndSignalCompletion_AllDone(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	_ = feature.EnsureFeatureDir(dir, "test")

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "completed", "skipped"},
	}
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1.1", Status: "completed"},
		"t2": {ID: "1.2", Status: "skipped"},
	})
	_ = task.SaveIndex(indexPath, index)

	out := captureStderr2(func() {
		saveIndexAndSignalCompletion(indexPath, dir, "test", index)
	})

	// Forge state should be written
	fs := feature.ReadForgeState(dir)
	if fs == nil || !fs.AllCompleted {
		t.Error("forge state should be written when all tasks done")
	}
	_ = out
}

// ---------- runValidateIndex no-args (feature-based path) ----------

func TestRunValidate_NoArgs(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "pending", Priority: "P0", File: "1.1.md", Dependencies: []string{}, Type: "implementation"},
	}})

	out := captureStdout(func() {
		runValidateIndex(nil, []string{})
	})
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS via feature resolution, got: %s", out)
	}
}

// ---------- runCheckDeps with wildcard ----------

func TestRunCheck_WildcardMatch(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t0": {ID: "1.0", Title: "T0", Status: "completed", Dependencies: []string{}},
		"t1": {ID: "1.1", Title: "T1", Status: "pending", Dependencies: []string{"1.x"}},
	}})

	out := captureStdout(func() {
		captureStderr2(func() {
			runCheckDeps(nil, []string{})
		})
	})
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS for wildcard deps, got: %s", out)
	}
	if !strings.Contains(out, "wildcard") {
		t.Errorf("expected wildcard in output, got: %s", out)
	}
}

// ---------- readSubmitData no data path ----------

func TestReadRecordData_NoData(t *testing.T) {
	// When no --data flag and no stdin pipe
	_, err := readSubmitData("")
	if err == nil {
		t.Error("expected error when no data provided")
	}
	if !strings.Contains(err.Error(), "no input") {
		t.Errorf("error should mention no input: %v", err)
	}
}

// ---------- checkExistingTaskState: failed to load state ----------

func TestCheckExistingTaskState_LoadFail(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	_ = feature.EnsureFeatureDir(dir, "test")

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "in_progress", "completed"},
	}
	_ = task.SaveIndex(indexPath, index)

	// Write invalid JSON to state file to trigger load failure
	statePath := feature.GetTaskStatePath(dir, "test")
	_ = os.MkdirAll(filepath.Dir(statePath), 0755)
	_ = os.WriteFile(statePath, []byte("invalid json"), 0644)

	continueTask, hasIssues, issues := checkExistingTaskState(dir, index, statePath)
	if continueTask {
		t.Error("should not continue with invalid state")
	}
	if hasIssues {
		t.Errorf("load failure should not report issues: %v", issues)
	}
}

// ---------- runSubmit with blocked status ----------

func TestRunRecord_BlockedStatus(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "Task One", Status: "in_progress", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	}})

	dir, _ := os.Getwd()

	rd := task.RecordData{
		Status:  "blocked",
		Summary: "Blocked by dependency",
		Notes:   "Waiting for upstream",
	}
	rdJSON, _ := json.Marshal(rd)
	dataPath := filepath.Join(dir, "record.json")
	_ = os.WriteFile(dataPath, rdJSON, 0644)

	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1", StartedTime: "2026-01-01 10:00"})

	submitDataPath = dataPath
	submitJSON = false
	submitQuiet = false
	submitForce = false

	out := captureStdout(func() {
		runSubmit(nil, []string{"1.1"})
	})
	if !strings.Contains(out, "STATUS: blocked") {
		t.Errorf("expected blocked status, got: %s", out)
	}
}

// ---------- appendFixTask removed (agent handles fix tasks now) ----------

// ---------- writeUnitTestRawOutput ----------

// TestWriteUnitTestRawOutput_CompilePrefix verifies compile failure output is prefixed correctly.
func TestWriteUnitTestRawOutput_CompilePrefix(t *testing.T) {
	dir := t.TempDir()
	compileOutput := "src/main.ts(10,5): error TS2345: Argument of type 'number' is not assignable"
	err := writeUnitTestRawOutput(dir, "=== compile failure ===\n"+compileOutput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "tests", "results", "unit-raw-output.txt"))
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	if !strings.Contains(string(data), "compile failure") {
		t.Errorf("expected compile prefix in output, got: %s", string(data))
	}
	if !strings.Contains(string(data), "TS2345") {
		t.Errorf("expected compile error in output, got: %s", string(data))
	}
}

func TestWriteUnitTestRawOutput(t *testing.T) {
	dir := t.TempDir()
	output := "FAIL\n--- FAIL: TestFoo (0.01s)"

	err := writeUnitTestRawOutput(dir, output)
	if err != nil {
		t.Fatalf("writeUnitTestRawOutput() error = %v", err)
	}

	path := filepath.Join(dir, "tests", "results", "unit-raw-output.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected file at %s, got error: %v", path, err)
	}
	if string(data) != output {
		t.Errorf("content = %q, want %q", string(data), output)
	}
}

func TestWriteUnitTestRawOutput_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	// tests/results/ does not exist yet — function must create it
	err := writeUnitTestRawOutput(dir, "output")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "tests", "results")); os.IsNotExist(err) {
		t.Error("tests/results/ directory should have been created")
	}
}

// ---------- writeRegressionRawOutput ----------

func TestWriteRegressionRawOutput(t *testing.T) {
	dir := t.TempDir()
	output := "not ok 1 - login test\n  Error: expected 200, got 404"

	err := writeRegressionRawOutput(dir, output)
	if err != nil {
		t.Fatalf("writeRegressionRawOutput() error = %v", err)
	}

	path := filepath.Join(dir, "tests", "e2e", "results", "raw-output.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected file at %s, got error: %v", path, err)
	}
	if string(data) != output {
		t.Errorf("content = %q, want %q", string(data), output)
	}
}

func TestWriteRegressionRawOutput_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	err := writeRegressionRawOutput(dir, "output")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "tests", "e2e", "results")); os.IsNotExist(err) {
		t.Error("tests/e2e/results/ directory should have been created")
	}
}

// ---------- runFeature: display no feature ----------

func TestRunFeature_None(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	_ = os.MkdirAll(filepath.Join(dir, "docs", "features"), 0755)

	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	out := captureStdout(func() {
		runFeature(nil, []string{})
	})
	if !strings.Contains(out, "(none)") {
		t.Errorf("expected (none) for no feature, got: %s", out)
	}
}

// ---------- runValidateIndex with invalid file ----------

func TestRunValidate_InvalidFile(t *testing.T) {
	v := &validator{filePath: "/nonexistent/path/index.json"}
	err := v.run()
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

// ---------- validateTTest1Template file read error ----------

func TestValidateTTest1Template_MissingFile(t *testing.T) {
	v := &validator{}
	v.validateFirstTestTaskTemplate("/nonexistent/task.md", "T-test-1", []string{"{{LAST_BUSINESS_TASK_ID}}"})
	if len(v.errors) != 0 {
		t.Errorf("missing file should not add errors, got: %v", v.errors)
	}
}

// ---------- validateQualityGate ----------

func TestValidateQualityGate_PassingGate(t *testing.T) {
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed, skipping")
	}

	dir := t.TempDir()
	// Create a justfile with all required recipes that succeed
	justfile := `
compile:
    echo "compile ok"

fmt:
    echo "fmt ok"

lint:
    echo "lint ok"

test:
    echo "test ok"
`
	if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte(justfile), 0644); err != nil {
		t.Fatal(err)
	}

	// Should not exit -- validateQualityGate only exits on failure
	exited := false
	// validateQualityGate calls Exit on failure which calls os.Exit(1).
	// For success path, it just returns.
	validateQualityGate(dir, "")
	_ = exited
}

func TestValidateQualityGate_FailingCompileGate(t *testing.T) {
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed, skipping")
	}

	dir := t.TempDir()
	// Create a justfile where compile fails (blocking step)
	justfile := `
compile:
    echo "compile fail" && exit 1

fmt:
    echo "fmt ok"

lint:
    echo "lint ok"

test:
    echo "test ok"
`
	if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte(justfile), 0644); err != nil {
		t.Fatal(err)
	}

	if os.Getenv("TEST_QUALITY_GATE_COMPILE_FAIL") == "1" {
		validateQualityGate(dir, "")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestValidateQualityGate_FailingCompileGate")
	cmd.Env = append(os.Environ(), "TEST_QUALITY_GATE_COMPILE_FAIL=1")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit when compile fails")
	}
	out := string(output)
	if !strings.Contains(out, "Quality gate failed") {
		t.Errorf("expected quality gate failure message, got: %s", out)
	}
	if !strings.Contains(out, "VALIDATION_ERROR") {
		t.Errorf("expected VALIDATION_ERROR code, got: %s", out)
	}
}

func TestValidateQualityGate_FailingLintGate(t *testing.T) {
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed, skipping")
	}

	dir := t.TempDir()
	// compile passes, lint fails (blocking step)
	justfile := `
compile:
    echo "compile ok"

fmt:
    echo "fmt ok"

lint:
    echo "lint fail" && exit 1

test:
    echo "test ok"
`
	if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte(justfile), 0644); err != nil {
		t.Fatal(err)
	}

	if os.Getenv("TEST_QUALITY_GATE_LINT_FAIL") == "1" {
		validateQualityGate(dir, "")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestValidateQualityGate_FailingLintGate")
	cmd.Env = append(os.Environ(), "TEST_QUALITY_GATE_LINT_FAIL=1")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit when lint fails")
	}
	out := string(output)
	if !strings.Contains(out, "Quality gate failed") {
		t.Errorf("expected quality gate failure message, got: %s", out)
	}
}

func TestValidateQualityGate_FailingTestGate(t *testing.T) {
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed, skipping")
	}

	dir := t.TempDir()
	// compile and lint pass, test fails (blocking step)
	justfile := `
compile:
    echo "compile ok"

fmt:
    echo "fmt ok"

lint:
    echo "lint ok"

test:
    echo "test fail" && exit 1
`
	if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte(justfile), 0644); err != nil {
		t.Fatal(err)
	}

	if os.Getenv("TEST_QUALITY_GATE_TEST_FAIL") == "1" {
		validateQualityGate(dir, "")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestValidateQualityGate_FailingTestGate")
	cmd.Env = append(os.Environ(), "TEST_QUALITY_GATE_TEST_FAIL=1")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit when test fails")
	}
	out := string(output)
	if !strings.Contains(out, "Quality gate failed") {
		t.Errorf("expected quality gate failure message, got: %s", out)
	}
}

func TestValidateQualityGate_NoJustfile(t *testing.T) {
	dir := t.TempDir()
	// No justfile -- RunGate returns true immediately, no exit
	validateQualityGate(dir, "")
}

func TestValidateQualityGate_FmtNonBlockingFailure(t *testing.T) {
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed, skipping")
	}

	dir := t.TempDir()
	// fmt fails but is non-blocking, should still pass
	justfile := `
compile:
    echo "compile ok"

fmt:
    exit 1

lint:
    echo "lint ok"

test:
    echo "test ok"
`
	if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte(justfile), 0644); err != nil {
		t.Fatal(err)
	}

	// Should not exit -- fmt is non-blocking
	validateQualityGate(dir, "")
}

// ---------- write*Output MkdirAll error paths ----------

func TestWriteRawOutput_MkdirAllError(t *testing.T) {
	dir := t.TempDir()
	// Create a file where the directory should be, so MkdirAll fails
	resultsDir := filepath.Join(dir, feature.GetFeatureTestingResultsDir("test"))
	_ = os.MkdirAll(filepath.Dir(resultsDir), 0755)
	_ = os.WriteFile(resultsDir, []byte("blocker"), 0644)

	err := writeRawOutput(dir, "test", "output")
	if err == nil {
		t.Error("expected error when MkdirAll fails")
	}
}

func TestWriteUnitTestRawOutput_MkdirAllError(t *testing.T) {
	dir := t.TempDir()
	// Create a file where tests/results/ should be, so MkdirAll fails
	testsDir := filepath.Join(dir, "tests")
	_ = os.WriteFile(testsDir, []byte("blocker"), 0644)

	err := writeUnitTestRawOutput(dir, "output")
	if err == nil {
		t.Error("expected error when MkdirAll fails")
	}
}

func TestWriteRegressionRawOutput_MkdirAllError(t *testing.T) {
	dir := t.TempDir()
	// Create a file where tests/e2e/results/ should be, so MkdirAll fails
	testsDir := filepath.Join(dir, "tests")
	_ = os.WriteFile(testsDir, []byte("blocker"), 0644)

	err := writeRegressionRawOutput(dir, "output")
	if err == nil {
		t.Error("expected error when MkdirAll fails")
	}
}

// ---------- runValidateIndex error paths ----------

func TestRunValidate_NoProjectRoot(t *testing.T) {
	if os.Getenv("TEST_RUN_VALIDATE_NO_PROJECT") == "1" {
		runValidateIndex(nil, []string{})
		return
	}

	tmpDir := t.TempDir()
	cmd := exec.Command(os.Args[0], "-test.run=TestRunValidate_NoProjectRoot")
	// Build clean env: clear CLAUDE_PROJECT_DIR so FindProjectRoot walks up
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "CLAUDE_PROJECT_DIR=") || strings.HasPrefix(e, "PROJECT_ROOT=") {
			continue
		}
		env = append(env, e)
	}
	cmd.Env = append(slices.Clone(env), "TEST_RUN_VALIDATE_NO_PROJECT=1", "CLAUDE_PROJECT_DIR=")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for no project root")
	}
	// Accept both NO_PROJECT (no markers found) and NO_FEATURE (markers found in ancestor dirs, but no feature set)
	if !strings.Contains(string(output), "NO_PROJECT") && !strings.Contains(string(output), "NO_FEATURE") {
		t.Errorf("expected NO_PROJECT or NO_FEATURE error, got: %s", string(output))
	}
}

func TestRunValidate_NoFeatureSet(t *testing.T) {
	if os.Getenv("TEST_RUN_VALIDATE_NO_FEATURE") == "1" {
		runValidateIndex(nil, []string{})
		return
	}

	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	_ = os.MkdirAll(filepath.Join(dir, "docs", "features"), 0755)

	cmd := exec.Command(os.Args[0], "-test.run=TestRunValidate_NoFeatureSet")
	// Clear env vars and set CLAUDE_PROJECT_DIR to our temp dir
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "CLAUDE_PROJECT_DIR=") || strings.HasPrefix(e, "PROJECT_ROOT=") {
			continue
		}
		env = append(env, e)
	}
	cmd.Env = append(slices.Clone(env), "TEST_RUN_VALIDATE_NO_FEATURE=1", "CLAUDE_PROJECT_DIR="+dir)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for no feature")
	}
	if !strings.Contains(string(output), "NO_FEATURE") {
		t.Errorf("expected NO_FEATURE error, got: %s", string(output))
	}
}

func TestRunValidate_IndexFileNotFound(t *testing.T) {
	if os.Getenv("TEST_RUN_VALIDATE_NO_INDEX") == "1" {
		runValidateIndex(nil, []string{})
		return
	}

	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	_ = feature.EnsureFeatureDir(dir, "testf")
	_ = feature.SetFeature(dir, "testf")

	cmd := exec.Command(os.Args[0], "-test.run=TestRunValidate_IndexFileNotFound")
	// Clear env and set CLAUDE_PROJECT_DIR
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "CLAUDE_PROJECT_DIR=") || strings.HasPrefix(e, "PROJECT_ROOT=") {
			continue
		}
		env = append(env, e)
	}
	cmd.Env = append(slices.Clone(env), "TEST_RUN_VALIDATE_NO_INDEX=1", "CLAUDE_PROJECT_DIR="+dir)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for missing index file")
	}
	if !strings.Contains(string(output), "NO_FEATURE") {
		t.Errorf("expected NO_FEATURE error, got: %s", string(output))
	}
}

// ---------- runCheckDeps error paths ----------

func TestRunCheck_NoProjectRoot(t *testing.T) {
	if os.Getenv("TEST_RUN_CHECK_NO_PROJECT") == "1" {
		runCheckDeps(nil, []string{})
		return
	}

	tmpDir := t.TempDir()
	cmd := exec.Command(os.Args[0], "-test.run=TestRunCheck_NoProjectRoot")
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "CLAUDE_PROJECT_DIR=") || strings.HasPrefix(e, "PROJECT_ROOT=") {
			continue
		}
		env = append(env, e)
	}
	cmd.Env = append(slices.Clone(env), "TEST_RUN_CHECK_NO_PROJECT=1", "CLAUDE_PROJECT_DIR=")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for no project root")
	}
	// Accept both NO_PROJECT (no markers found) and NO_FEATURE (markers found in ancestor dirs, but no feature set)
	if !strings.Contains(string(output), "NO_PROJECT") && !strings.Contains(string(output), "NO_FEATURE") {
		t.Errorf("expected NO_PROJECT or NO_FEATURE error, got: %s", string(output))
	}
}

func TestRunCheck_NoFeatureSet(t *testing.T) {
	if os.Getenv("TEST_RUN_CHECK_NO_FEATURE") == "1" {
		runCheckDeps(nil, []string{})
		return
	}

	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	_ = os.MkdirAll(filepath.Join(dir, "docs", "features"), 0755)

	cmd := exec.Command(os.Args[0], "-test.run=TestRunCheck_NoFeatureSet")
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "CLAUDE_PROJECT_DIR=") || strings.HasPrefix(e, "PROJECT_ROOT=") {
			continue
		}
		env = append(env, e)
	}
	cmd.Env = append(slices.Clone(env), "TEST_RUN_CHECK_NO_FEATURE=1", "CLAUDE_PROJECT_DIR="+dir)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for no feature")
	}
	if !strings.Contains(string(output), "NO_FEATURE") {
		t.Errorf("expected NO_FEATURE error, got: %s", string(output))
	}
}

func TestRunCheck_IndexFileNotFound(t *testing.T) {
	if os.Getenv("TEST_RUN_CHECK_NO_INDEX") == "1" {
		runCheckDeps(nil, []string{})
		return
	}

	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	_ = feature.EnsureFeatureDir(dir, "testf")
	_ = feature.SetFeature(dir, "testf")

	cmd := exec.Command(os.Args[0], "-test.run=TestRunCheck_IndexFileNotFound")
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "CLAUDE_PROJECT_DIR=") || strings.HasPrefix(e, "PROJECT_ROOT=") {
			continue
		}
		env = append(env, e)
	}
	cmd.Env = append(slices.Clone(env), "TEST_RUN_CHECK_NO_INDEX=1", "CLAUDE_PROJECT_DIR="+dir)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for missing index file")
	}
	if !strings.Contains(string(output), "NO_FEATURE") {
		t.Errorf("expected NO_FEATURE error, got: %s", string(output))
	}
}

// ---------- runCheckDeps with invalid deps ----------

func TestRunCheck_InvalidDeps(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "pending", Dependencies: []string{"9.9"}},
	}})

	if os.Getenv("TEST_RUN_CHECK_INVALID_DEPS") == "1" {
		runCheckDeps(nil, []string{})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRunCheck_InvalidDeps")
	cmd.Env = append(os.Environ(), "TEST_RUN_CHECK_INVALID_DEPS=1")
	dir, _ := os.Getwd()
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for invalid deps")
	}
	if !strings.Contains(string(output), "FAIL") {
		t.Errorf("expected FAIL output, got: %s", string(output))
	}
	if !strings.Contains(string(output), "does NOT exist") {
		t.Errorf("expected 'does NOT exist' error message, got: %s", string(output))
	}
}

// ---------- saveIndexAndSignalCompletion error paths ----------

func TestSaveIndexAndSignalCompletion_SaveIndexError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod on directories has no effect on Windows")
	}
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	_ = feature.EnsureFeatureDir(dir, "test")

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "completed"},
	}
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1.1", Status: "completed"},
	})
	_ = task.SaveIndex(indexPath, index)

	// Make the parent directory read-only so SaveIndexAtomic (temp+rename) fails
	indexDir := filepath.Dir(indexPath)
	_ = os.Chmod(indexDir, 0555)
	defer func() { _ = os.Chmod(indexDir, 0755) }()

	if os.Getenv("TEST_SAVE_INDEX_ERROR") == "1" {
		saveIndexAndSignalCompletion(indexPath, dir, "test", index)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestSaveIndexAndSignalCompletion_SaveIndexError")
	cmd.Env = append(os.Environ(), "TEST_SAVE_INDEX_ERROR=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for save index failure")
	}
	if !strings.Contains(string(output), "Failed to update task index") {
		t.Errorf("expected save index error message, got: %s", string(output))
	}
}

func TestSaveIndexAndSignalCompletion_WriteForgeStateWarning(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Create the feature directory structure manually
	featureDir := filepath.Join(dir, "docs", "features", "test", "tasks")
	_ = os.MkdirAll(featureDir, 0755)

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "completed"},
	}
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1.1", Status: "completed"},
	})
	_ = task.SaveIndex(indexPath, index)

	// Make .forge directory read-only to trigger WriteForgeState warning
	forgeDir := filepath.Join(dir, ".forge")
	_ = os.MkdirAll(forgeDir, 0755)
	// Create a file named state.json that is a directory (causes write to fail)
	_ = os.MkdirAll(filepath.Join(forgeDir, "state.json"), 0755)

	out := captureStderr2(func() {
		saveIndexAndSignalCompletion(indexPath, dir, "test", index)
	})
	if !strings.Contains(out, "WARNING") {
		t.Errorf("expected warning about failed forge state write, got: %s", out)
	}
}

// TestForgeStateLifecycle verifies the full .forge/state.json lifecycle:
// claim (creates allCompleted=false) → record (overwrites to true) → all-completed (deletes)
func TestForgeStateLifecycle(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	_ = feature.EnsureFeatureDir(dir, "lf")

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("lf"))
	index := &task.TaskIndex{
		Feature:      "lf",
		StatusEnum:   []string{"pending", "in_progress", "completed"},
		PriorityEnum: []string{"P0"},
	}
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "pending", Priority: "P0", File: "1.1.md", Record: "1.1.md"},
	})
	_ = task.SaveIndex(indexPath, index)
	_ = os.WriteFile(filepath.Join(dir, "docs", "features", "lf", "tasks", "1.1.md"), []byte("# T1"), 0644)
	_ = os.MkdirAll(filepath.Join(dir, "docs", "features", "lf", "tasks", "records"), 0755)

	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Phase 1: claim creates state.json with allCompleted=false
	claimResult, err := executeClaim()
	if err != nil {
		t.Fatalf("claim failed: %v", err)
	}
	state := feature.ReadForgeState(dir)
	if state == nil {
		t.Fatal("state.json should exist after claim")
	}
	if state.AllCompleted {
		t.Error("allCompleted should be false after claim")
	}

	// Phase 2: record overwrites state.json with allCompleted=true
	submitDataPath := filepath.Join(dir, "docs", "features", "lf", "tasks", "process", "record.json")
	rd := map[string]any{
		"taskId":      "1.1",
		"status":      "completed",
		"summary":     "done",
		"coverage":    -1.0,
		"testsPassed": 0,
		"testsFailed": 0,
	}
	rdJSON, _ := json.Marshal(rd)
	_ = os.WriteFile(submitDataPath, rdJSON, 0644)

	rootCmd.SetArgs([]string{"task", "submit", claimResult.Task.ID, "--data", submitDataPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("record failed: %v", err)
	}

	state = feature.ReadForgeState(dir)
	if state == nil {
		t.Fatal("state.json should exist after record")
	}
	if !state.AllCompleted {
		t.Error("allCompleted should be true after all tasks recorded")
	}

	// Phase 3: all-completed reads and deletes state.json
	result := checkAllCompleted(false)
	if result == nil {
		t.Fatal("checkAllCompleted should return result when all done with state")
	}

	state = feature.ReadForgeState(dir)
	if state != nil {
		t.Error("state.json should be deleted after all-completed consumes it")
	}
}

// ---------- error constructors ----------

func TestErrTaskIDConflict(t *testing.T) {
	err := ErrTaskIDConflict("1.1")
	if err.Code != ErrConflict {
		t.Errorf("Code = %q, want %q", err.Code, ErrConflict)
	}
	if !strings.Contains(err.Message, "1.1") {
		t.Errorf("Message should contain '1.1', got %q", err.Message)
	}
}

func TestErrInvalidDependency(t *testing.T) {
	err := ErrInvalidDependency([]string{"2.1", "2.2"})
	if err.Code != ErrValidation {
		t.Errorf("Code = %q, want %q", err.Code, ErrValidation)
	}
	if !strings.Contains(err.Message, "2.1") {
		t.Errorf("Message should contain '2.1', got %q", err.Message)
	}
}

// ---------- Exit with non-AIError ----------

func TestExit_NonAIError(t *testing.T) {
	if os.Getenv("TEST_EXIT_PLAIN_ERR") == "1" {
		Exit(fmt.Errorf("plain error"))
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExit_NonAIError")
	cmd.Env = append(os.Environ(), "TEST_EXIT_PLAIN_ERR=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit")
	}
	if !strings.Contains(string(output), "ERROR: plain error") {
		t.Errorf("expected plain error message, got: %s", string(output))
	}
}

// ---------- runAdd / executeAdd ----------

func TestRunAdd_NoProject(t *testing.T) {
	if os.Getenv("TEST_RUN_ADD_NO_PROJECT") == "1" {
		addTitle = "Test"
		runAdd(nil, []string{})
		return
	}

	tmpDir := t.TempDir()
	cmd := exec.Command(os.Args[0], "-test.run=TestRunAdd_NoProject")
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "CLAUDE_PROJECT_DIR=") {
			continue
		}
		env = append(env, e)
	}
	cmd.Env = append(slices.Clone(env), "TEST_RUN_ADD_NO_PROJECT=1", "CLAUDE_PROJECT_DIR=")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit")
	}
	out := string(output)
	if !strings.Contains(out, "NO_PROJECT") && !strings.Contains(out, "NO_FEATURE") {
		t.Errorf("expected project or feature error, got: %s", out)
	}
}

func TestRunAdd_Success(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{}})

	if os.Getenv("TEST_RUN_ADD_SUCCESS") == "1" {
		addTitle = "New Task"
		addPriority = "P1"
		runAdd(nil, []string{})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRunAdd_Success")
	cmd.Env = append(os.Environ(), "TEST_RUN_ADD_SUCCESS=1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("expected success, got error: %v\n%s", err, string(output))
	}
	if !strings.Contains(string(output), "ADDED") {
		t.Errorf("expected ADDED in output, got: %s", string(output))
	}
}

// ---------- runCleanup ----------

func TestRunCleanup_Success(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed", File: "1.1.md", Record: "records/1.1.md"},
	}})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

	if os.Getenv("TEST_RUN_CLEANUP") == "1" {
		runCleanup(nil, []string{})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRunCleanup_Success")
	cmd.Env = append(os.Environ(), "TEST_RUN_CLEANUP=1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("expected success, got error: %v\n%s", err, string(output))
	}
	_ = output
}

// ---------- runVerifyTaskDone ----------

func TestRunVerifyCompletion_Success(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed", File: "1.1.md", Record: "records/1.1.md"},
	}})

	dir, _ := os.Getwd()
	_ = os.MkdirAll(filepath.Join(dir, "docs", "features", "test", "tasks", "records"), 0755)
	_ = os.WriteFile(filepath.Join(dir, "docs", "features", "test", "tasks", "records", "1.1.md"), []byte("record"), 0644)

	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

	if os.Getenv("TEST_RUN_VERIFY_OK") == "1" {
		runVerifyTaskDone(nil, []string{})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRunVerifyCompletion_Success")
	cmd.Env = append(os.Environ(), "TEST_RUN_VERIFY_OK=1")
	err := cmd.Run()
	if err != nil {
		t.Errorf("expected success (exit 0), got: %v", err)
	}
}

func TestRunVerifyCompletion_Fail(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "in_progress", File: "1.1.md"},
	}})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

	if os.Getenv("TEST_RUN_VERIFY_FAIL") == "1" {
		runVerifyTaskDone(nil, []string{})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRunVerifyCompletion_Fail")
	cmd.Env = append(os.Environ(), "TEST_RUN_VERIFY_FAIL=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit")
	}
	if !strings.Contains(string(output), "not completed") {
		t.Errorf("expected 'not completed' in output, got: %s", string(output))
	}
}

// ---------- runQualityGate ----------

func TestRunAllCompleted_NotAllDone(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "pending", File: "1.1.md"},
	}})

	if os.Getenv("TEST_RUN_ALL_COMPLETED_NOT_DONE") == "1" {
		runQualityGate(nil, []string{})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRunAllCompleted_NotAllDone")
	cmd.Env = append(os.Environ(), "TEST_RUN_ALL_COMPLETED_NOT_DONE=1")
	err := cmd.Run()
	if err != nil {
		t.Errorf("expected exit 0 when not all done, got: %v", err)
	}
}

func TestRunRecord_AutoRestore_SlugKeyedSource(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"run-e2e":  {ID: "T-test-3", Title: "Run e2e tests", Status: "blocked", Priority: "P0", File: "T-test-3.md", Record: "records/T-test-3.md", Dependencies: []string{"fix-auth"}},
		"fix-auth": {ID: "fix-auth", Title: "Fix auth", Status: "in_progress", Priority: "P0", File: "fix-auth.md", Record: "records/fix-auth.md", SourceTaskID: "T-test-3"},
	}})

	rd := task.RecordData{
		Status:             "completed",
		Summary:            "Fixed auth",
		TestsPassed:        3,
		Coverage:           85.0,
		KeyDecisions:       []string{"added retry logic"},
		AcceptanceCriteria: []task.AcceptanceCriterion{{Criterion: "Tests pass", Met: true}},
	}
	rdJSON, _ := json.Marshal(rd)
	dataPath := filepath.Join(dir, "record.json")
	_ = os.WriteFile(dataPath, rdJSON, 0644)

	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "fix-auth", Key: "fix-auth", StartedTime: "2026-01-01 10:00"})

	submitDataPath = dataPath
	submitJSON = false
	submitQuiet = false
	submitForce = false

	_ = captureStdout(func() {
		runSubmit(nil, []string{"fix-auth"})
	})

	// Verify source task was auto-restored
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	// Source restored to pending
	if index.TasksMap()["run-e2e"].Status != "pending" {
		t.Errorf("source task should be restored to pending, got %s", index.TasksMap()["run-e2e"].Status)
	}

	// Fix task completed
	if index.TasksMap()["fix-auth"].Status != "completed" {
		t.Errorf("fix task should be completed, got %s", index.TasksMap()["fix-auth"].Status)
	}

	// No duplicate key created under task ID
	if _, hasDup := index.TasksMap()["T-test-3"]; hasDup {
		t.Error("should not create duplicate entry under ID key 'T-test-3'")
	}
}

func TestRunRecord_FixTaskAutoDowngrade_NoRestore(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"source": {ID: "T-test-3", Title: "Run e2e tests", Status: "blocked", Priority: "P0", File: "T-test-3.md", Record: "records/T-test-3.md", Dependencies: []string{"fix-1"}},
		"fix-1":  {ID: "fix-1", Title: "Fix auth", Status: "in_progress", Priority: "P0", File: "fix-1.md", Record: "records/fix-1.md", SourceTaskID: "T-test-3"},
	}})

	rd := task.RecordData{
		Status:       "completed",
		Summary:      "Attempted fix, some tests still fail",
		TestsPassed:  2,
		TestsFailed:  1,
		Coverage:     60.0,
		KeyDecisions: []string{"partial fix"},
	}
	rdJSON, _ := json.Marshal(rd)
	dataPath := filepath.Join(dir, "record.json")
	_ = os.WriteFile(dataPath, rdJSON, 0644)

	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "fix-1", Key: "fix-1", StartedTime: "2026-01-01 10:00"})

	submitDataPath = dataPath
	submitJSON = false
	submitQuiet = false
	submitForce = false

	_ = captureStdout(func() {
		runSubmit(nil, []string{"fix-1"})
	})

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	if index.TasksMap()["fix-1"].Status != "blocked" {
		t.Errorf("fix-task should be auto-downgraded to blocked, got %s", index.TasksMap()["fix-1"].Status)
	}
	if index.TasksMap()["source"].Status != "blocked" {
		t.Errorf("source should stay blocked when fix-task also fails, got %s", index.TasksMap()["source"].Status)
	}
}

func TestRunRecord_AutoDowngrade_ThenCleanup(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"task1": {ID: "1.1", Title: "Task 1", Status: "in_progress", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	}})

	rd := task.RecordData{
		Status:       "completed",
		Summary:      "Some tests failed",
		TestsPassed:  3,
		TestsFailed:  2,
		Coverage:     60.0,
		KeyDecisions: []string{"best effort"},
	}
	rdJSON, _ := json.Marshal(rd)
	dataPath := filepath.Join(dir, "record.json")
	_ = os.WriteFile(dataPath, rdJSON, 0644)

	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "task1", StartedTime: "2026-01-01 10:00"})

	submitDataPath = dataPath
	submitJSON = false
	submitQuiet = false
	submitForce = false

	_ = captureStdout(func() {
		runSubmit(nil, []string{"1.1"})
	})

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index, _ := task.LoadIndex(indexPath)
	if index.TasksMap()["task1"].Status != "blocked" {
		t.Fatalf("expected blocked, got %s", index.TasksMap()["task1"].Status)
	}

	// state.json persists after record (record.go doesn't delete it)
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Fatal("state.json should exist after record")
	}

	// Cleanup should delete state.json for blocked tasks
	cleanupCompletedTaskState()

	if _, err := os.Stat(statePath); !os.IsNotExist(err) {
		t.Error("state.json should be deleted by cleanup for blocked task")
	}
}

func TestRunRecord_AutoDowngrade_ThenClaim(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"task1": {ID: "1.1", Title: "Task 1", Status: "in_progress", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
		"task2": {ID: "1.2", Title: "Task 2", Status: "pending", Priority: "P1", File: "1.2.md", Record: "records/1.2.md"},
	}})

	rd := task.RecordData{
		Status:       "completed",
		Summary:      "Tests failed",
		TestsPassed:  3,
		TestsFailed:  2,
		Coverage:     60.0,
		KeyDecisions: []string{"partial"},
	}
	rdJSON, _ := json.Marshal(rd)
	dataPath := filepath.Join(dir, "record.json")
	_ = os.WriteFile(dataPath, rdJSON, 0644)

	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "task1", StartedTime: "2026-01-01 10:00"})

	submitDataPath = dataPath
	submitJSON = false
	submitQuiet = false
	submitForce = false

	_ = captureStdout(func() {
		runSubmit(nil, []string{"1.1"})
	})

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index, _ := task.LoadIndex(indexPath)
	if index.TasksMap()["task1"].Status != "blocked" {
		t.Fatalf("expected blocked, got %s", index.TasksMap()["task1"].Status)
	}

	// Claim should clear blocked state and claim task2
	result, err := executeClaim()
	if err != nil {
		t.Fatalf("claim should succeed after blocked task cleanup: %v", err)
	}
	if result.Task.ID != "1.2" {
		t.Errorf("expected to claim task2, got %s", result.Task.ID)
	}

	// New state should be for task2
	newStatePath := feature.GetTaskStatePath(dir, "test")
	newState, _ := task.LoadState(newStatePath)
	if newState.TaskID != "1.2" {
		t.Errorf("new state should be for task2, got %s", newState.TaskID)
	}
}

func TestRunRecord_MultiFixTask_PartialDowngrade(t *testing.T) {
	dir := setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"source": {ID: "T-test-3", Title: "Run e2e tests", Status: "blocked", Priority: "P0", File: "T-test-3.md", Record: "records/T-test-3.md", Dependencies: []string{"fix-1", "fix-2"}},
		"fix-1":  {ID: "fix-1", Title: "Fix auth", Status: "in_progress", Priority: "P0", File: "fix-1.md", Record: "records/fix-1.md", SourceTaskID: "T-test-3"},
		"fix-2":  {ID: "fix-2", Title: "Fix timeout", Status: "pending", Priority: "P0", File: "fix-2.md", Record: "records/fix-2.md", SourceTaskID: "T-test-3"},
	}})

	// Step 1: Record fix-1 as completed
	rd1 := task.RecordData{
		Status:             "completed",
		Summary:            "Fixed auth",
		TestsPassed:        3,
		Coverage:           85.0,
		KeyDecisions:       []string{"added retry"},
		AcceptanceCriteria: []task.AcceptanceCriterion{{Criterion: "Tests pass", Met: true}},
	}
	rd1JSON, _ := json.Marshal(rd1)
	dataPath1 := filepath.Join(dir, "record1.json")
	_ = os.WriteFile(dataPath1, rd1JSON, 0644)

	statePath := feature.GetTaskStatePath(dir, "test")
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "fix-1", Key: "fix-1", StartedTime: "2026-01-01 10:00"})

	submitDataPath = dataPath1
	submitJSON = false
	submitQuiet = false
	submitForce = false

	_ = captureStdout(func() {
		runSubmit(nil, []string{"fix-1"})
	})

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index, _ := task.LoadIndex(indexPath)

	if index.TasksMap()["fix-1"].Status != "completed" {
		t.Fatalf("fix-1 should be completed, got %s", index.TasksMap()["fix-1"].Status)
	}
	if index.TasksMap()["source"].Status != "blocked" {
		t.Errorf("source should stay blocked with fix-2 pending, got %s", index.TasksMap()["source"].Status)
	}

	// Step 2: Record fix-2 with test failures (auto-downgrade)
	_ = task.SaveState(statePath, &task.TaskState{TaskID: "fix-2", Key: "fix-2", StartedTime: "2026-01-01 11:00"})

	rd2 := task.RecordData{
		Status:       "completed",
		Summary:      "Attempted fix, tests still fail",
		TestsPassed:  2,
		TestsFailed:  1,
		Coverage:     50.0,
		KeyDecisions: []string{"partial"},
	}
	rd2JSON, _ := json.Marshal(rd2)
	dataPath2 := filepath.Join(dir, "record2.json")
	_ = os.WriteFile(dataPath2, rd2JSON, 0644)

	submitDataPath = dataPath2

	_ = captureStdout(func() {
		runSubmit(nil, []string{"fix-2"})
	})

	index, _ = task.LoadIndex(indexPath)

	if index.TasksMap()["fix-2"].Status != "blocked" {
		t.Errorf("fix-2 should be auto-downgraded to blocked, got %s", index.TasksMap()["fix-2"].Status)
	}
	if index.TasksMap()["source"].Status != "blocked" {
		t.Errorf("source should stay blocked when fix-2 is blocked, got %s", index.TasksMap()["source"].Status)
	}
}
