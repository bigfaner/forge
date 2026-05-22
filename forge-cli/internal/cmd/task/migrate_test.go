package task

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/prompt"
	"forge-cli/pkg/task"
)

// TestRunMigrate_HappyPath verifies that all tasks get type fields inferred and
// the success message is printed.
func TestRunMigrate_HappyPath(t *testing.T) {
	setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1":      {ID: "1.1", Title: "Impl task", Status: "pending", File: "1.1.md", Record: "records/1.1.md"},
			"t-gate":  {ID: "1.gate", Title: "Gate", Status: "pending", File: "1-gate.md", Record: "records/1-gate.md"},
			"t-sum":   {ID: "1.summary", Title: "Summary", Status: "pending", File: "1-summary.md", Record: "records/1-summary.md"},
			"t-fix":   {ID: "fix-1", Title: "Fix", Status: "pending", File: "fix-1.md", Record: "records/fix-1.md"},
			"t-test1": {ID: "T-test-gen-cases", Title: "Gen cases", Status: "pending", File: "T-test-gen-cases.md", Record: "records/T-test-gen-cases.md"},
		},
	})

	out := captureStdout(func() {
		_ = runMigrate(nil, []string{})
	})

	if !strings.Contains(out, "Migrated 5 tasks") {
		t.Errorf("expected 'Migrated 5 tasks' in output, got: %s", out)
	}
	if !strings.Contains(out, "Run task validate to verify") {
		t.Errorf("expected 'Run task validate to verify' in output, got: %s", out)
	}

	// Verify types were written to index
	dir, _ := os.Getwd()
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	cases := map[string]string{
		"t1":      task.TypeCodingFeature, // unknown ID → conservative default is "feature"
		"t-gate":  task.TypeGate,
		"t-sum":   task.TypeDocSummary,
		"t-fix":   task.TypeCodingFix,
		"t-test1": task.TypeTestGenCases,
	}
	for key, wantType := range cases {
		got := index.TasksMap()[key].Type
		if got != wantType {
			t.Errorf("task %q: type = %q, want %q", key, got, wantType)
		}
	}
}

// TestRunMigrate_StatusPreserved verifies that task statuses are not changed.
func TestRunMigrate_StatusPreserved(t *testing.T) {
	setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "T1", Status: "completed", File: "1.1.md", Record: "records/1.1.md"},
			"t2": {ID: "1.2", Title: "T2", Status: "blocked", File: "1.2.md", Record: "records/1.2.md"},
			"t3": {ID: "1.3", Title: "T3", Status: "pending", File: "1.3.md", Record: "records/1.3.md"},
		},
	})

	captureStdout(func() {
		_ = runMigrate(nil, []string{})
	})

	dir, _ := os.Getwd()
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	if index.TasksMap()["t1"].Status != "completed" {
		t.Errorf("t1 status should remain completed, got %s", index.TasksMap()["t1"].Status)
	}
	if index.TasksMap()["t2"].Status != "blocked" {
		t.Errorf("t2 status should remain blocked, got %s", index.TasksMap()["t2"].Status)
	}
	if index.TasksMap()["t3"].Status != "pending" {
		t.Errorf("t3 status should remain pending, got %s", index.TasksMap()["t3"].Status)
	}
}

// TestRunMigrate_Idempotent verifies that running migrate twice produces the same result.
func TestRunMigrate_Idempotent(t *testing.T) {
	setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "T1", Status: "pending", File: "1.1.md", Record: "records/1.1.md", Type: task.TypeGate},
		},
	})

	captureStdout(func() { _ = runMigrate(nil, []string{}) })
	captureStdout(func() { _ = runMigrate(nil, []string{}) })

	dir, _ := os.Getwd()
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	// InferType("1.1") → "" (no fallback), migrate defaults to feature (conservative default)
	if index.TasksMap()["t1"].Type != task.TypeCodingFeature {
		t.Errorf("type = %q, want %q", index.TasksMap()["t1"].Type, task.TypeCodingFeature)
	}
}

// TestRunMigrate_InProgress_ExitsWithError verifies that migrate aborts when any
// task is in_progress, leaving index.json unchanged.
func TestRunMigrate_InProgress_ExitsWithError(t *testing.T) {
	setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "T1", Status: "in_progress", File: "1.1.md", Record: "records/1.1.md"},
			"t2": {ID: "1.2", Title: "T2", Status: "pending", File: "1.2.md", Record: "records/1.2.md"},
		},
	})

	if os.Getenv("TEST_MIGRATE_IN_PROGRESS") == "1" {
		_ = runMigrate(nil, []string{})
		return
	}

	dir, _ := os.Getwd()
	cmd := exec.Command(os.Args[0], "-test.run=TestRunMigrate_InProgress_ExitsWithError")
	cmd.Env = append(os.Environ(), "TEST_MIGRATE_IN_PROGRESS=1")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit when in_progress task exists")
	}
	if !strings.Contains(string(output), "in_progress") {
		t.Errorf("expected 'in_progress' in stderr, got: %s", string(output))
	}

	// index.json must not be modified: t1 should still have no type
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index, loadErr := task.LoadIndex(indexPath)
	if loadErr != nil {
		t.Fatal(loadErr)
	}
	if index.TasksMap()["t1"].Type != "" {
		t.Errorf("index should not be modified when in_progress task exists, got type=%q", index.TasksMap()["t1"].Type)
	}
}

// TestRunMigrate_AllKnownIDPatterns verifies InferType covers all documented patterns.
func TestRunMigrate_AllKnownIDPatterns(t *testing.T) {
	cases := []struct {
		id       string
		wantType string
	}{
		{"1.1", ""},
		{"2.3", ""},
		{"1.gate", task.TypeGate},
		{"2.gate", task.TypeGate},
		{"1.summary", task.TypeDocSummary},
		{"2.summary", task.TypeDocSummary},
		{"fix-1", task.TypeCodingFix},
		{"disc-1", task.TypeCodingFix},
		{"T-test-gen-cases", task.TypeTestGenCases},
		{"T-test-eval-cases", task.TypeTestEvalCases},
		{"T-test-gen-scripts", task.TypeTestGenScripts},
		{"T-test-run", task.TypeTestRun},
		{"T-test-graduate", task.TypeTestGraduate},
		{"T-test-verify-regression", task.TypeTestVerifyRegression},
		{"T-specs-consolidate", task.TypeDocConsolidate},
	}

	for _, tc := range cases {
		got := prompt.InferType(tc.id)
		if got != tc.wantType {
			t.Errorf("InferType(%q) = %q, want %q", tc.id, got, tc.wantType)
		}
	}
}

// TestMigrateCmd_RegisteredInRoot verifies migrateCmd is registered.
func TestMigrateCmd_RegisteredInRoot(t *testing.T) {
	for _, cmd := range Cmd.Commands() {
		if cmd.Use == "migrate" || cmd.Name() == "migrate" {
			return
		}
	}
	t.Error("migrateCmd not registered in Cmd")
}

// TestRunMigrate_NoProject_ExitsWithError verifies error when no project root.
func TestRunMigrate_NoProject_ExitsWithError(t *testing.T) {
	if os.Getenv("TEST_MIGRATE_NO_PROJECT") == "1" {
		_ = runMigrate(nil, []string{})
		return
	}

	tmpDir := t.TempDir()
	cmd := exec.Command(os.Args[0], "-test.run=TestRunMigrate_NoProject_ExitsWithError")
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "CLAUDE_PROJECT_DIR=") || strings.HasPrefix(e, "PROJECT_ROOT=") {
			continue
		}
		env = append(env, e)
	}
	env = append(env, "TEST_MIGRATE_NO_PROJECT=1", "CLAUDE_PROJECT_DIR=")
	cmd.Env = env
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for no project")
	}
	out := string(output)
	if !strings.Contains(out, "NO_PROJECT") && !strings.Contains(out, "NO_FEATURE") {
		t.Errorf("expected NO_PROJECT or NO_FEATURE error, got: %s", out)
	}
}

// TestRunMigrate_NoFeature_ExitsWithError verifies error when no feature is set.
func TestRunMigrate_NoFeature_ExitsWithError(t *testing.T) {
	if os.Getenv("TEST_MIGRATE_NO_FEATURE") == "1" {
		_ = runMigrate(nil, []string{})
		return
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "docs", "features"), 0755); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRunMigrate_NoFeature_ExitsWithError")
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "CLAUDE_PROJECT_DIR=") || strings.HasPrefix(e, "PROJECT_ROOT=") {
			continue
		}
		env = append(env, e)
	}
	env = append(env, "TEST_MIGRATE_NO_FEATURE=1", "CLAUDE_PROJECT_DIR="+dir)
	cmd.Env = env
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for no feature")
	}
	if !strings.Contains(string(output), "NO_FEATURE") {
		t.Errorf("expected NO_FEATURE error, got: %s", string(output))
	}
}
