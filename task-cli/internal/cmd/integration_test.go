package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"task-cli/pkg/feature"
	"task-cli/pkg/task"
)

// setupFullProject creates a project with go.mod, feature dir, index, task files, and chdirs.
func setupFullProject(t *testing.T, tasks map[string]task.Task) (dir string) {
	t.Helper()
	dir = t.TempDir()

	// go.mod marks project root
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
		t.Fatal(err)
	}

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:      "test",
		PRD:          "prd/prd-spec.md",
		Design:       "design/tech-design.md",
		StatusEnum:   []string{"pending", "in_progress", "completed", "blocked", "skipped"},
		PriorityEnum: []string{"P0", "P1", "P2"},
		Tasks:        tasks,
	}
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Create task markdown files
	tasksDir := filepath.Join(dir, feature.GetFeatureTasksDir("test"))
	for _, t2 := range tasks {
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

	// Set working dir
	origWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(origWd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	// Set feature
	if err := feature.SetFeature(dir, "test"); err != nil {
		t.Fatal(err)
	}
	return dir
}

// ---------- verifyTaskCompletion ----------

func TestVerifyTaskCompletion_HappyPath(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed", File: "1.1.md", Record: "records/1.1.md"},
	})

	// Write record file
	dir, _ := os.Getwd()
	_ = os.MkdirAll(filepath.Join(dir, "docs", "features", "test", "tasks", "records"), 0755)
	os.WriteFile(filepath.Join(dir, "docs", "features", "test", "tasks", "records", "1.1.md"), []byte("record"), 0644)

	// Save task state
	statePath := feature.GetTaskStatePath(dir, "test")
	task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

	err := verifyTaskCompletion()
	if err != nil {
		t.Errorf("expected nil, got: %v", err)
	}
}

func TestVerifyTaskCompletion_TaskNotCompleted(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "in_progress", File: "1.1.md", Record: "records/1.1.md"},
	})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

	err := verifyTaskCompletion()
	if err == nil {
		t.Error("expected error for non-completed task")
	}
	if !strings.Contains(err.Error(), "not completed") {
		t.Errorf("error should mention not completed: %v", err)
	}
}

func TestVerifyTaskCompletion_RecordFileMissing(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed", File: "1.1.md", Record: "records/1.1.md"},
	})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

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
	t.Cleanup(func() { os.Chdir(origWd) })
	os.Chdir(tmpDir)

	err := verifyTaskCompletion()
	if err != nil {
		t.Errorf("no project should return nil, got: %v", err)
	}
}

func TestVerifyTaskCompletion_NoFeature(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	os.MkdirAll(filepath.Join(dir, "docs", "features"), 0755)

	origWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(origWd) })
	os.Chdir(dir)

	err := verifyTaskCompletion()
	if err != nil {
		t.Errorf("no feature should return nil, got: %v", err)
	}
}

func TestVerifyTaskCompletion_NoState(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed"},
	})

	err := verifyTaskCompletion()
	if err != nil {
		t.Errorf("no state should return nil, got: %v", err)
	}
}

func TestVerifyTaskCompletion_TaskNotFoundInIndex(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed"},
	})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	task.SaveState(statePath, &task.TaskState{TaskID: "9.9", Key: "missing"})

	err := verifyTaskCompletion()
	if err == nil {
		t.Error("expected error for task not in index")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should mention not found: %v", err)
	}
}

func TestVerifyTaskCompletion_NoRecordField(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed", File: "1.1.md"},
	})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

	// Task has empty Record — should pass (no file to check)
	err := verifyTaskCompletion()
	if err != nil {
		t.Errorf("empty record field should pass, got: %v", err)
	}
}

// ---------- cleanupCompletedTaskState ----------

func TestCleanupCompletedTaskState_Completed(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed", File: "1.1.md", Record: "records/1.1.md"},
	})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

	// Also create record.json
	recordPath := feature.GetProcessRecordPath(dir, "test")
	os.MkdirAll(filepath.Dir(recordPath), 0755)
	os.WriteFile(recordPath, []byte("{}"), 0644)

	cleanupCompletedTaskState()

	if _, err := os.Stat(statePath); !os.IsNotExist(err) {
		t.Error("state.json should be deleted for completed task")
	}
	if _, err := os.Stat(recordPath); !os.IsNotExist(err) {
		t.Error("record.json should be deleted for completed task")
	}
}

func TestCleanupCompletedTaskState_InProgress(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "in_progress", File: "1.1.md"},
	})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1"})

	cleanupCompletedTaskState()

	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Error("state.json should NOT be deleted for in_progress task")
	}
}

func TestCleanupCompletedTaskState_NoState(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed"},
	})

	// No state file — should not panic
	cleanupCompletedTaskState()
}

func TestCleanupCompletedTaskState_NoProject(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(origWd) })
	os.Chdir(tmpDir)

	// Should not panic
	cleanupCompletedTaskState()
}

func TestCleanupCompletedTaskState_TaskKeyNotFound(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed"},
	})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	// State references a key that doesn't exist in index
	task.SaveState(statePath, &task.TaskState{TaskID: "9.9", Key: "nonexistent"})

	cleanupCompletedTaskState()

	// Should not delete state when key doesn't match
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Error("state.json should NOT be deleted when task key not found")
	}
}

// ---------- runRecord integration ----------

func TestRunRecord_HappyPath(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "Task One", Status: "in_progress", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	})

	dir, _ := os.Getwd()

	// Create a record data file
	rd := RecordData{
		Status:       "completed",
		Summary:      "Did the thing",
		TestsPassed:  5,
		TestsFailed:  0,
		Coverage:     90.0,
		KeyDecisions: []string{"used approach X"},
		AcceptanceCriteria: []AcceptanceCriterion{
			{Criterion: "It works", Met: true},
		},
	}
	rdJSON, _ := json.Marshal(rd)
	dataPath := filepath.Join(dir, "record.json")
	os.WriteFile(dataPath, rdJSON, 0644)

	// Save state for startedTime
	statePath := feature.GetTaskStatePath(dir, "test")
	task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1", StartedTime: "2026-01-01 10:00"})

	recordDataPath = dataPath
	recordJSON = false
	recordQuiet = false
	recordForce = false

	out := captureStdout(func() {
		runRecord(nil, []string{"1.1"})
	})

	if !strings.Contains(out, "TASK_ID: 1.1") {
		t.Errorf("expected task ID in output, got: %s", out)
	}
	if !strings.Contains(out, "STATUS: completed") {
		t.Errorf("expected status in output, got: %s", out)
	}

	// Verify record file was created
	recordFile := filepath.Join(dir, "docs", "features", "test", "tasks", "records", "1.1.md")
	if _, err := os.Stat(recordFile); os.IsNotExist(err) {
		t.Error("record file should exist")
	}
}

func TestRunRecord_JSONOutput(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "Task One", Status: "in_progress", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	})

	dir, _ := os.Getwd()

	rd := RecordData{
		Status:      "completed",
		Summary:     "JSON test",
		TestsPassed: 1,
		Coverage:    80.0,
		AcceptanceCriteria: []AcceptanceCriterion{
			{Criterion: "Works", Met: true},
		},
	}
	rdJSON, _ := json.Marshal(rd)
	dataPath := filepath.Join(dir, "record.json")
	os.WriteFile(dataPath, rdJSON, 0644)

	statePath := feature.GetTaskStatePath(dir, "test")
	task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1", StartedTime: "2026-01-01 10:00"})

	recordDataPath = dataPath
	recordJSON = true
	recordQuiet = false
	recordForce = false

	out := captureStdout(func() {
		runRecord(nil, []string{"1.1"})
	})

	if !strings.Contains(out, `"recordFile"`) {
		t.Errorf("expected JSON output with recordFile, got: %s", out)
	}
	if !strings.Contains(out, `"taskId"`) || !strings.Contains(out, `"1.1"`) {
		t.Errorf("expected JSON output with taskId 1.1, got: %s", out)
	}
}

func TestRunRecord_QuietOutput(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "Task One", Status: "in_progress", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	})

	dir, _ := os.Getwd()

	rd := RecordData{
		Status:      "completed",
		Summary:     "Quiet test",
		TestsPassed: 1,
		Coverage:    75.0,
		AcceptanceCriteria: []AcceptanceCriterion{
			{Criterion: "Works", Met: true},
		},
	}
	rdJSON, _ := json.Marshal(rd)
	dataPath := filepath.Join(dir, "record.json")
	os.WriteFile(dataPath, rdJSON, 0644)

	statePath := feature.GetTaskStatePath(dir, "test")
	task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1", StartedTime: "2026-01-01 10:00"})

	recordDataPath = dataPath
	recordJSON = false
	recordQuiet = true
	recordForce = false

	out := captureStdout(func() {
		runRecord(nil, []string{"1.1"})
	})

	if strings.Contains(out, "TASK_ID") {
		t.Errorf("quiet mode should not print block output, got: %s", out)
	}
}

// ---------- printHookJSON ----------

func TestPrintHookJSON(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		out := captureStdout(func() {
			printHookJSON(map[string]any{"decision": "block", "reason": "test"})
		})
		if !strings.Contains(out, `"decision"`) {
			t.Errorf("expected JSON output, got: %s", out)
		}
	})
}

// ---------- hasNpmTestScript ----------

func TestHasNpmTestScript(t *testing.T) {
	t.Run("has test script", func(t *testing.T) {
		dir := t.TempDir()
		pkg := `{"scripts": {"test": "jest"}}`
		os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0644)
		if !hasNpmTestScript(dir) {
			t.Error("expected true for package with test script")
		}
	})

	t.Run("no test script", func(t *testing.T) {
		dir := t.TempDir()
		pkg := `{"scripts": {"build": "tsc"}}`
		os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0644)
		if hasNpmTestScript(dir) {
			t.Error("expected false for package without test script")
		}
	})

	t.Run("no package.json", func(t *testing.T) {
		dir := t.TempDir()
		if hasNpmTestScript(dir) {
			t.Error("expected false when no package.json")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "package.json"), []byte("not json"), 0644)
		if hasNpmTestScript(dir) {
			t.Error("expected false for invalid JSON")
		}
	})
}

// ---------- executeClaim error paths ----------

func TestExecuteClaim_DataIntegrityError(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	})

	dir, _ := os.Getwd()
	// Create state pointing to a key that doesn't exist in index
	statePath := feature.GetTaskStatePath(dir, "test")
	task.SaveState(statePath, &task.TaskState{TaskID: "9.9", Key: "nonexistent", StartedTime: "2026-01-01 10:00"})

	_, err := executeClaim()
	if err == nil {
		t.Error("expected data integrity error")
	}
	if !strings.Contains(err.Error(), "integrity") {
		t.Errorf("error should mention integrity: %v", err)
	}
}

func TestExecuteClaim_CompletedStateClaimNew(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "completed", Priority: "P0", File: "1.1.md"},
		"t2": {ID: "1.2", Title: "T2", Status: "pending", Priority: "P0", File: "1.2.md", Record: "records/1.2.md"},
	})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1", StartedTime: "2026-01-01 10:00"})

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

func TestExecuteClaim_UnexpectedStatus(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "blocked", Priority: "P0", File: "1.1.md"},
	})

	dir, _ := os.Getwd()
	statePath := feature.GetTaskStatePath(dir, "test")
	task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1", StartedTime: "2026-01-01 10:00"})

	_, err := executeClaim()
	if err == nil {
		t.Error("expected error for unexpected status")
	}
	if !strings.Contains(err.Error(), "integrity") {
		t.Errorf("error should mention integrity: %v", err)
	}
}

// ---------- runValidate direct validator ----------

func TestValidatorRun_WithFileArg(t *testing.T) {
	dir := t.TempDir()

	// Create a valid index.json
	index := &task.TaskIndex{
		Feature:    "test-feature",
		StatusEnum: []string{"pending", "in_progress", "completed"},
		Tasks:      map[string]task.Task{},
	}
	data, _ := json.Marshal(index)
	indexPath := filepath.Join(dir, "index.json")
	os.WriteFile(indexPath, data, 0644)

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
	os.WriteFile(indexPath, []byte("not json"), 0644)

	v := &validator{filePath: indexPath}
	err := v.run()
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

// ---------- fillRecordTemplate ----------

func TestFillRecordTemplate_NonCompletedStatus(t *testing.T) {
	t2 := &task.Task{ID: "1.1", Title: "Test Task"}
	rd := &RecordData{
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
	rd := &RecordData{
		Status:      "completed",
		Summary:     "Done",
		Notes:       "Custom notes here",
		TestsPassed: 1,
		Coverage:    50.0,
		AcceptanceCriteria: []AcceptanceCriterion{
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
	feature.EnsureFeatureDir(dir, "test")

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "in_progress", "completed"},
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Status: "completed"},
			"t2": {ID: "1.2", Status: "pending"},
		},
	}
	task.SaveIndex(indexPath, index)

	// Should NOT write forge state since not all tasks are done
	saveIndexAndSignalCompletion(indexPath, dir, "test", index)

	forgeState := feature.ReadForgeState(dir)
	if forgeState != nil {
		t.Error("forge state should NOT be written when tasks are incomplete")
	}
}

// ---------- validateRecordData ----------

func TestValidateRecordData_ForceOverride(t *testing.T) {
	rd := &RecordData{
		Status:       "completed",
		Summary:      "Done",
		TestsPassed:  0,
		TestsFailed:  0,
		Coverage:     50.0,
		AcceptanceCriteria: []AcceptanceCriterion{
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
	rd := &RecordData{
		Status:      "completed",
		Summary:     "Docs only",
		Coverage:    -1.0,
		KeyDecisions: []string{"doc-only"},
		AcceptanceCriteria: []AcceptanceCriterion{
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

// ---------- runValidate no file arg, feature-based ----------

func TestValidatorRun_FeatureBased(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	feature.EnsureFeatureDir(dir, "test")

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "completed"},
		Tasks:      map[string]task.Task{},
	}
	task.SaveIndex(indexPath, index)

	feature.SetFeature(dir, "test")

	origWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(origWd) })
	os.Chdir(dir)

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
	os.WriteFile(taskFile, []byte("# Task\nReplace {{LAST_BUSINESS_TASK_ID}} with actual ID\n"), 0644)

	v := &validator{}
	v.validateTTest1Template(taskFile)
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
	os.WriteFile(taskFile, []byte("# Task\nDepends on 1.5\n"), 0644)

	v := &validator{}
	v.validateTTest1Template(taskFile)
	if len(v.errors) != 0 {
		t.Errorf("expected no errors, got: %v", v.errors)
	}
}

// ---------- runCmd / runShell ----------

func TestRunCmd_Success(t *testing.T) {
	out := captureStderr2(func() {
		runCmd(t.TempDir(), "echo", "hello from runcmd")
	})
	if !strings.Contains(out, "hello from runcmd") {
		t.Errorf("expected echo output on stderr, got: %s", out)
	}
}

func TestRunCmd_Failure(t *testing.T) {
	out := captureStderr2(func() {
		runCmd(t.TempDir(), "false")
	})
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected error output for failing command, got: %s", out)
	}
}

func TestRunShell_Success(t *testing.T) {
	out := captureStderr2(func() {
		runShell(t.TempDir(), "echo shell-output")
	})
	if !strings.Contains(out, "shell-output") {
		t.Errorf("expected shell output, got: %s", out)
	}
}

func TestRunShell_Failure(t *testing.T) {
	out := captureStderr2(func() {
		runShell(t.TempDir(), "exit 1")
	})
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected error for failing shell command, got: %s", out)
	}
}

// ---------- runProjectTests ----------

func TestRunProjectTests_GoMod(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)

	out := captureStderr2(func() {
		runProjectTests(dir, "")
	})
	// go test ./... on empty module should still run (even if no tests)
	// Just verify it doesn't panic
	_ = out
}

func TestRunProjectTests_CustomCommand(t *testing.T) {
	dir := t.TempDir()
	out := captureStderr2(func() {
		runProjectTests(dir, "echo custom-test-command")
	})
	if !strings.Contains(out, "custom-test-command") {
		t.Errorf("expected custom command output, got: %s", out)
	}
}

func TestRunProjectTests_NpmTest(t *testing.T) {
	dir := t.TempDir()
	pkg := `{"scripts": {"test": "echo npm-test-pass"}}`
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0644)

	out := captureStderr2(func() {
		runProjectTests(dir, "")
	})
	if !strings.Contains(out, "npm-test-pass") {
		t.Errorf("expected npm test output, got: %s", out)
	}
}

// ---------- findTask ----------

func TestFindTaskByKey(t *testing.T) {
	index := &task.TaskIndex{
		Tasks: map[string]task.Task{
			"task1": {ID: "1.1", Title: "Task One"},
		},
	}

	key, t2, err := findTask(index, "task1")
	if err != nil {
		t.Fatal(err)
	}
	if key != "task1" {
		t.Errorf("key = %q, want task1", key)
	}
	_ = t2
}

// ---------- readRecordData ----------

func TestReadRecordData_FromFile(t *testing.T) {
	dir := t.TempDir()
	rd := RecordData{Summary: "test summary", TestsPassed: 1, Coverage: 50.0}
	data, _ := json.Marshal(rd)
	path := filepath.Join(dir, "record.json")
	os.WriteFile(path, data, 0644)

	result, err := readRecordData(path)
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
	os.WriteFile(path, []byte("not json"), 0644)

	_, err := readRecordData(path)
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
		ID:           "2.gate",
		Title:        "Gate Task",
		Priority:     "P0",
		Status:       "pending",
		Breaking:     true,
		File:         "2.gate.md",
		Record:       "records/2.gate.md",
		EstimatedTime: "30min",
		Dependencies: []string{"1.summary"},
	}

	out := captureStdout(func() {
		printTaskDetails("gate-2", t2, "/project", "test")
	})
	if !strings.Contains(out, "BREAKING: true") {
		t.Errorf("expected BREAKING field, got: %s", out)
	}
	if !strings.Contains(out, "ESTIMATED_TIME: 30min") {
		t.Errorf("expected ESTIMATED_TIME, got: %s", out)
	}
	if !strings.Contains(out, "DEPENDENCIES: 1.summary") {
		t.Errorf("expected DEPENDENCIES, got: %s", out)
	}
}

// ---------- runProjectTests no test command ----------

func TestRunProjectTests_NoTestCommand(t *testing.T) {
	dir := t.TempDir()
	out := captureStdout(func() {
		runProjectTests(dir, "")
	})
	if !strings.Contains(out, "WARNING") {
		t.Errorf("expected warning when no test command found, got: %s", out)
	}
}

// ---------- runStatus update mode ----------

func TestRunStatus_Update(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "Status Task", Status: "pending", Priority: "P0", File: "1.1.md", Record: "records/1.1.md", Dependencies: []string{}},
	})

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
	if index.Tasks["t1"].Status != "blocked" {
		t.Errorf("index status = %q, want blocked", index.Tasks["t1"].Status)
	}
}

// ---------- executeClaim error: no project ----------

func TestExecuteClaim_NoProject(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(origWd) })
	os.Chdir(tmpDir)

	_, err := executeClaim()
	if err == nil {
		t.Error("expected error for no project root")
	}
}

// ---------- executeClaim: save index error ----------

func TestExecuteClaim_SaveIndexError(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	feature.EnsureFeatureDir(dir, "test")

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:      "test",
		StatusEnum:   []string{"pending", "in_progress", "completed"},
		PriorityEnum: []string{"P0", "P1", "P2"},
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "T1", Status: "pending", Priority: "P0", File: "1.1.md", Record: "1.1.md"},
		},
	}
	task.SaveIndex(indexPath, index)

	feature.SetFeature(dir, "test")
	origWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(origWd) })
	os.Chdir(dir)

	// Make index.json read-only so SaveIndex fails
	os.Chmod(indexPath, 0444)
	defer os.Chmod(indexPath, 0644)

	_, err := executeClaim()
	if err == nil {
		t.Error("expected error when save index fails")
	}
}

// ---------- runClaim output paths ----------

func TestRunClaim_Output(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "Claim Task", Status: "pending", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	})

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

// ---------- fileExists ----------

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, "exists.txt")
	os.WriteFile(existing, []byte("x"), 0644)

	if !fileExists(existing) {
		t.Error("expected true for existing file")
	}
	if fileExists(filepath.Join(dir, "nope.txt")) {
		t.Error("expected false for non-existing file")
	}
}

// ---------- runCheck integration (valid deps, exits 0 via PrintResult) ----------

func TestRunCheck_AllValid(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t0": {ID: "1.0", Title: "T0", Status: "completed", Dependencies: []string{}},
		"t1": {ID: "1.1", Title: "T1", Status: "pending", Dependencies: []string{"1.0"}},
	})

	out := captureStdout(func() {
		captureStderr2(func() {
			runCheck(nil, []string{})
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

// ---------- runValidate with explicit file arg ----------

func TestRunValidate_Integration(t *testing.T) {
	dir := t.TempDir()

	index := &task.TaskIndex{
		Feature:    "my-feature",
		PRD:        "prd/prd-spec.md",
		Design:     "design/tech-design.md",
		StatusEnum: []string{"pending", "in_progress", "completed"},
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "T1", Status: "pending", Priority: "P0", File: "1.1.md", Dependencies: []string{}},
		},
	}
	data, _ := json.Marshal(index)
	indexPath := filepath.Join(dir, "index.json")
	os.WriteFile(indexPath, data, 0644)

	// Create tasks dir and task file
	tasksDir := filepath.Join(dir, "tasks")
	os.MkdirAll(tasksDir, 0755)
	os.WriteFile(filepath.Join(tasksDir, "1.1.md"), []byte("# T1"), 0644)

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

// ---------- printHookJSON error path ----------

func TestPrintHookJSON_MarshalError(t *testing.T) {
	// Chan cannot be marshaled to JSON
	out := captureStderr2(func() {
		printHookJSON(map[string]any{"ch": make(chan int)})
	})
	if !strings.Contains(out, "WARNING") {
		t.Errorf("expected warning for marshal error, got: %s", out)
	}
}

// ---------- saveIndexAndSignalCompletion with forge state ----------

func TestSaveIndexAndSignalCompletion_AllDone(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	feature.EnsureFeatureDir(dir, "test")

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "completed", "skipped"},
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Status: "completed"},
			"t2": {ID: "1.2", Status: "skipped"},
		},
	}
	task.SaveIndex(indexPath, index)

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

// ---------- runProjectTests: pytest branch ----------

func TestRunProjectTests_Pytest(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "pytest.ini"), []byte("[pytest]\n"), 0644)

	out := captureStderr2(func() {
		runProjectTests(dir, "")
	})
	_ = out
	// Just verify it doesn't panic — pytest may or may not be installed
}

// ---------- runValidate no-args (feature-based path) ----------

func TestRunValidate_NoArgs(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "T1", Status: "pending", Priority: "P0", File: "1.1.md", Dependencies: []string{}},
	})

	out := captureStdout(func() {
		runValidate(nil, []string{})
	})
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS via feature resolution, got: %s", out)
	}
}

// ---------- runCheck with wildcard ----------

func TestRunCheck_WildcardMatch(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t0": {ID: "1.0", Title: "T0", Status: "completed", Dependencies: []string{}},
		"t1": {ID: "1.1", Title: "T1", Status: "pending", Dependencies: []string{"1.x"}},
	})

	out := captureStdout(func() {
		captureStderr2(func() {
			runCheck(nil, []string{})
		})
	})
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS for wildcard deps, got: %s", out)
	}
	if !strings.Contains(out, "wildcard") {
		t.Errorf("expected wildcard in output, got: %s", out)
	}
}

// ---------- readRecordData no data path ----------

func TestReadRecordData_NoData(t *testing.T) {
	// When no --data flag and no stdin pipe
	_, err := readRecordData("")
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
	feature.EnsureFeatureDir(dir, "test")

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "in_progress", "completed"},
		Tasks:      map[string]task.Task{},
	}
	task.SaveIndex(indexPath, index)

	// Write invalid JSON to state file to trigger load failure
	statePath := feature.GetTaskStatePath(dir, "test")
	os.MkdirAll(filepath.Dir(statePath), 0755)
	os.WriteFile(statePath, []byte("invalid json"), 0644)

	continueTask, hasIssues, issues := checkExistingTaskState(dir, index, statePath)
	if continueTask {
		t.Error("should not continue with invalid state")
	}
	if hasIssues {
		t.Errorf("load failure should not report issues: %v", issues)
	}
}

// ---------- runRecord with blocked status ----------

func TestRunRecord_BlockedStatus(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "Task One", Status: "in_progress", Priority: "P0", File: "1.1.md", Record: "records/1.1.md"},
	})

	dir, _ := os.Getwd()

	rd := RecordData{
		Status:  "blocked",
		Summary: "Blocked by dependency",
		Notes:   "Waiting for upstream",
	}
	rdJSON, _ := json.Marshal(rd)
	dataPath := filepath.Join(dir, "record.json")
	os.WriteFile(dataPath, rdJSON, 0644)

	statePath := feature.GetTaskStatePath(dir, "test")
	task.SaveState(statePath, &task.TaskState{TaskID: "1.1", Key: "t1", StartedTime: "2026-01-01 10:00"})

	recordDataPath = dataPath
	recordJSON = false
	recordQuiet = false
	recordForce = false

	out := captureStdout(func() {
		runRecord(nil, []string{"1.1"})
	})
	if !strings.Contains(out, "STATUS: blocked") {
		t.Errorf("expected blocked status, got: %s", out)
	}
}

// ---------- appendFixTask removed (agent handles fix tasks now) ----------

// ---------- runProjectTests: justfile branch ----------

func TestRunProjectTests_Justfile(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "justfile"), []byte("test:\n    echo just-test-output\n"), 0644)

	out := captureStderr2(func() {
		runProjectTests(dir, "")
	})
	_ = out
	// Just verify no panic
}

// ---------- runProjectTests: Makefile branch ----------

func TestRunProjectTests_Makefile(t *testing.T) {
	dir := t.TempDir()

	// Skip if make not available
	if _, err := exec.LookPath("make"); err != nil {
		t.Skip("make not installed")
	}

	os.WriteFile(filepath.Join(dir, "Makefile"), []byte("test:\n\t@echo make-test-output\n"), 0644)

	out := captureStderr2(func() {
		runProjectTests(dir, "")
	})
	if !strings.Contains(out, "make-test-output") {
		t.Errorf("expected make test output, got: %s", out)
	}
}

// ---------- runFeature: display no feature ----------

func TestRunFeature_None(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	os.MkdirAll(filepath.Join(dir, "docs", "features"), 0755)

	origWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(origWd) })
	os.Chdir(dir)

	out := captureStdout(func() {
		runFeature(nil, []string{})
	})
	if !strings.Contains(out, "(none)") {
		t.Errorf("expected (none) for no feature, got: %s", out)
	}
}

// ---------- runValidate with invalid file ----------

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
	v.validateTTest1Template("/nonexistent/task.md")
	if len(v.errors) != 0 {
		t.Errorf("missing file should not add errors, got: %v", v.errors)
	}
}

// TestForgeStateLifecycle verifies the full .forge/state.json lifecycle:
// claim (creates allCompleted=false) → record (overwrites to true) → all-completed (deletes)
func TestForgeStateLifecycle(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	feature.EnsureFeatureDir(dir, "lf")

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("lf"))
	index := &task.TaskIndex{
		Feature:      "lf",
		StatusEnum:   []string{"pending", "in_progress", "completed"},
		PriorityEnum: []string{"P0"},
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "T1", Status: "pending", Priority: "P0", File: "1.1.md", Record: "1.1.md"},
		},
	}
	task.SaveIndex(indexPath, index)
	os.WriteFile(filepath.Join(dir, "docs", "features", "lf", "tasks", "1.1.md"), []byte("# T1"), 0644)
	os.MkdirAll(filepath.Join(dir, "docs", "features", "lf", "tasks", "records"), 0755)

	origWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(origWd) })
	os.Chdir(dir)

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
	recordDataPath := filepath.Join(dir, "docs", "features", "lf", "tasks", "process", "record.json")
	rd := map[string]any{
		"taskId":      "1.1",
		"status":      "completed",
		"summary":     "done",
		"coverage":    -1.0,
		"testsPassed": 0,
		"testsFailed": 0,
	}
	rdJSON, _ := json.Marshal(rd)
	os.WriteFile(recordDataPath, rdJSON, 0644)

	rootCmd.SetArgs([]string{"record", claimResult.Task.ID, "--data", recordDataPath})
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
	result, err := checkAllCompleted(false)
	if err != nil || result == nil {
		t.Fatal("checkAllCompleted should return result when all done with state")
	}

	state = feature.ReadForgeState(dir)
	if state != nil {
		t.Error("state.json should be deleted after all-completed consumes it")
	}
}
