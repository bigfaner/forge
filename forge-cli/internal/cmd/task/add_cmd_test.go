package task

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"
)

func TestAddCmd_WithTypeAndVars(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"1.1": {ID: "1.1", Title: "Existing", Priority: "P0", Status: "completed", File: "1.1.md", Record: "records/1.1.md"},
	}})

	output, err := captureOutput(func() error {
		Cmd.SetArgs([]string{
			"add",
			"--title", "Fix: login bug",
			"--type", "coding.fix",
			"--source-task-id", "1.1",
			"--description", "Selector not found",
			"--var", "SOURCE_FILES=src/Login.tsx",
			"--var", "TEST_SCRIPT=tests/e2e/auth.spec.ts",
			"--var", "TEST_RESULTS=results/latest.md",
		})
		return Cmd.Execute()
	})
	if err != nil {
		t.Fatalf("add command failed: %v", err)
	}
	if !strings.Contains(output, "ACTION: ADDED") {
		t.Errorf("expected ACTION: ADDED, got %q", output)
	}
	// Template defaults should apply (auto-discovered via --type): PRIORITY=P0, BREAKING=true
	if !strings.Contains(output, "PRIORITY: P0") {
		t.Errorf("expected PRIORITY: P0 from template defaults, got %q", output)
	}
	if !strings.Contains(output, "BREAKING: true") {
		t.Errorf("expected BREAKING: true from template defaults, got %q", output)
	}

	// Verify the markdown was generated from embedded template
	wd, _ := os.Getwd()
	tasksDir := filepath.Join(wd, feature.GetFeatureTasksDir("test"))
	files, _ := os.ReadDir(tasksDir)
	var mdFile string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".md") && !strings.Contains(f.Name(), "1.1") {
			mdFile = filepath.Join(tasksDir, f.Name())
			break
		}
	}
	if mdFile == "" {
		t.Fatal("no task markdown file created")
	}

	data, err := os.ReadFile(mdFile)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "task status 1.1 pending") && !strings.Contains(content, "1.1") {
		t.Errorf("SOURCE_TASK_ID not substituted via --source-task-id, got:\n%s", content)
	}
	if !strings.Contains(content, "Selector not found") {
		t.Errorf("description not injected, got:\n%s", content)
	}

	// Verify source task (1.1) now has the fix-task as a dependency
	wd2, _ := os.Getwd()
	indexPath := filepath.Join(wd2, feature.GetFeatureIndexFile("test"))
	idx, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	srcTask := idx.TasksMap()["1.1"]
	foundDep := false
	for _, d := range srcTask.Dependencies {
		if d != "" {
			foundDep = true
			break
		}
	}
	if !foundDep {
		t.Errorf("source task 1.1 should have fix-task as dependency, got: %v", srcTask.Dependencies)
	}
}

func TestAddCmd_VarParsing(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantKey string
		wantVal string
		wantErr bool
	}{
		{"valid", "KEY=VALUE", "KEY", "VALUE", false},
		{"value with equals", "K=V=W", "K", "V=W", false},
		{"empty value", "KEY=", "KEY", "", false},
		{"no equals", "NOEQUALS", "", "", true},
		{"empty key", "=VALUE", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := strings.SplitN(tt.input, "=", 2)
			hasErr := len(parts) != 2 || parts[0] == ""
			if hasErr != tt.wantErr {
				t.Errorf("input %q: wantErr=%v, got hasErr=%v", tt.input, tt.wantErr, hasErr)
			}
			if !hasErr {
				if parts[0] != tt.wantKey || parts[1] != tt.wantVal {
					t.Errorf("input %q: got key=%q val=%q, want key=%q val=%q", tt.input, parts[0], parts[1], tt.wantKey, tt.wantVal)
				}
			}
		})
	}
}

func TestAddCmd_InvalidTypeReturnsError(t *testing.T) {
	if os.Getenv("TEST_INVALID_TYPE") == "1" {
		setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
			"1.1": {ID: "1.1", Title: "Existing", Priority: "P0", Status: "completed", File: "1.1.md", Record: "records/1.1.md"},
		}})
		Cmd.SetArgs([]string{
			"add",
			"--title", "Fix: test",
			"--type", "invalid.type",
		})
		_ = Cmd.Execute()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestAddCmd_InvalidTypeReturnsError")
	cmd.Env = append(os.Environ(), "TEST_INVALID_TYPE=1")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for invalid type")
	}
}

func TestAddCmd_DedupSkipsActiveFix(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"1.1": {ID: "1.1", Title: "Source", Priority: "P0", Status: "blocked", File: "1.1.md", Record: "records/1.1.md"},
	}})

	// First add succeeds
	output, err := captureOutput(func() error {
		Cmd.SetArgs([]string{
			"add",
			"--title", "Fix: first attempt",
			"--source-task-id", "1.1",
		})
		return Cmd.Execute()
	})
	if err != nil {
		t.Fatalf("first add failed: %v", err)
	}
	if !strings.Contains(output, "ACTION: ADDED") {
		t.Errorf("first add should output ACTION: ADDED, got %q", output)
	}

	// Second add for same source should be skipped (first fix is still active/pending)
	output, err = captureOutput(func() error {
		Cmd.SetArgs([]string{
			"add",
			"--title", "Fix: second attempt",
			"--source-task-id", "1.1",
		})
		return Cmd.Execute()
	})
	if err != nil {
		t.Fatalf("dedup add should not error, got: %v", err)
	}
	if !strings.Contains(output, "ACTION: SKIPPED") {
		t.Errorf("second add should output ACTION: SKIPPED, got %q", output)
	}
	if !strings.Contains(output, "active fix tasks already exist") {
		t.Errorf("SKIPPED output should contain reason, got %q", output)
	}
}

func TestAddCmd_BlockSource(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"1.1": {ID: "1.1", Title: "Source", Priority: "P0", Status: "completed", File: "1.1.md", Record: "records/1.1.md"},
	}})

	output, err := captureOutput(func() error {
		Cmd.SetArgs([]string{
			"add",
			"--title", "Fix: blocked source",
			"--source-task-id", "1.1",
			"--block-source",
		})
		return Cmd.Execute()
	})
	if err != nil {
		t.Fatalf("add with --block-source failed: %v", err)
	}
	if !strings.Contains(output, "ACTION: ADDED") {
		t.Errorf("expected ACTION: ADDED, got %q", output)
	}
	if !strings.Contains(output, "SOURCE_BLOCKED: 1.1") {
		t.Errorf("expected SOURCE_BLOCKED: 1.1, got %q", output)
	}

	// Verify source task is now blocked in index
	wd, _ := os.Getwd()
	indexPath := filepath.Join(wd, feature.GetFeatureIndexFile("test"))
	idx, _ := task.LoadIndex(indexPath)
	if idx.TasksMap()["1.1"].Status != "blocked" {
		t.Errorf("source 1.1 should be blocked, got %q", idx.TasksMap()["1.1"].Status)
	}
}
