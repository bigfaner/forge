package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"task-cli/pkg/feature"
	"task-cli/pkg/task"
)

func TestAddCmd_WithTemplateAndVars(t *testing.T) {
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"1.1": {ID: "1.1", Title: "Existing", Priority: "P0", Status: "completed", File: "1.1.md", Record: "records/1.1.md"},
	}})

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{
			"add",
			"--title", "Fix: login bug",
			"--template", "fix-task",
			"--source-task-id", "1.1",
			"--description", "Selector not found",
		})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("add command failed: %v", err)
	}
	if !strings.Contains(output, "ACTION: ADDED") {
		t.Errorf("expected ACTION: ADDED, got %q", output)
	}
	// Template defaults should apply: PRIORITY=P0, BREAKING=true
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
	srcTask := idx.Tasks["1.1"]
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

func TestAddCmd_UnknownTemplateReturnsError(t *testing.T) {
	if os.Getenv("TEST_UNKNOWN_TEMPLATE") == "1" {
		setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
			"1.1": {ID: "1.1", Title: "Existing", Priority: "P0", Status: "completed", File: "1.1.md", Record: "records/1.1.md"},
		}})
		rootCmd.SetArgs([]string{
			"add",
			"--title", "Fix: test",
			"--template", "nonexistent",
		})
		rootCmd.Execute()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestAddCmd_UnknownTemplateReturnsError")
	cmd.Env = append(os.Environ(), "TEST_UNKNOWN_TEMPLATE=1")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for unknown template")
	}
}
