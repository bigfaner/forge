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

func TestCheckAllCompleted(t *testing.T) {
	tests := []struct {
		name        string
		tasks       map[string]task.Task
		testCommand string
		forgeState  bool
		wantNil     bool
		wantTestCmd string
	}{
		{
			name: "all completed with forge state returns result",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
				"t2": {ID: "1.2", Status: "completed"},
			},
			forgeState: true,
			wantNil:    false,
		},
		{
			name: "all skipped with forge state returns result",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "skipped"},
			},
			forgeState: true,
			wantNil:    false,
		},
		{
			name: "mixed completed and skipped with forge state returns result",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
				"t2": {ID: "1.2", Status: "skipped"},
			},
			forgeState: true,
			wantNil:    false,
		},
		{
			name: "one pending task returns nil even with forge state",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
				"t2": {ID: "1.2", Status: "pending"},
			},
			forgeState: true,
			wantNil:    true,
		},
		{
			name: "in_progress task returns nil",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "in_progress"},
			},
			wantNil: true,
		},
		{
			name: "blocked task returns nil",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "blocked"},
			},
			wantNil: true,
		},
		{
			name:       "empty task list with forge state returns result",
			tasks:      map[string]task.Task{},
			forgeState: true,
			wantNil:    false,
		},
		{
			name: "testCommand from index.json is propagated",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
			},
			testCommand: "make test",
			forgeState:  true,
			wantNil:     false,
			wantTestCmd: "make test",
		},
		{
			name: "no forge state returns nil even if all tasks completed",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
				"t2": {ID: "1.2", Status: "completed"},
			},
			forgeState: false,
			wantNil:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv("CLAUDE_PROJECT_DIR", dir)

			if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
				t.Fatal(err)
			}

			indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
			index := &task.TaskIndex{
				Feature:     "test",
				StatusEnum:  []string{"pending", "in_progress", "completed", "blocked", "skipped"},
				Tasks:       tc.tasks,
				TestCommand: tc.testCommand,
			}
			if err := task.SaveIndex(indexPath, index); err != nil {
				t.Fatal(err)
			}

			if tc.forgeState {
				if err := feature.WriteForgeState(dir, "test"); err != nil {
					t.Fatal(err)
				}
			}

			result, err := checkAllCompleted(false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantNil {
				if result != nil {
					t.Errorf("expected nil result, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("expected non-nil result, got nil")
			}

			if result.FeatureSlug != "test" {
				t.Errorf("FeatureSlug = %q, want %q", result.FeatureSlug, "test")
			}
			if result.ProjectRoot == "" {
				t.Error("ProjectRoot should not be empty")
			}

			if result.TestCommand != tc.wantTestCmd {
				t.Errorf("TestCommand = %q, want %q", result.TestCommand, tc.wantTestCmd)
			}
		})
	}
}

func TestCheckAllCompleted_NoFeature(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	if err := os.MkdirAll(filepath.Join(dir, feature.FeaturesDir), 0755); err != nil {
		t.Fatal(err)
	}

	result, err := checkAllCompleted(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result when no feature set, got %+v", result)
	}
}

func TestCheckAllCompleted_NoProject(t *testing.T) {
	t.Setenv("CLAUDE_PROJECT_DIR", "")

	result, err := checkAllCompleted(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result when no project root, got %+v", result)
	}
}

func TestHasJustfile(t *testing.T) {
	tests := []struct {
		name  string
		files []string
		want  bool
	}{
		{name: "no justfile", files: []string{}, want: false},
		{name: "lowercase justfile", files: []string{"justfile"}, want: true},
		{name: "capitalized Justfile", files: []string{"Justfile"}, want: true},
		{name: "both present", files: []string{"justfile", "Justfile"}, want: true},
		{name: "unrelated files only", files: []string{"Makefile", "go.mod"}, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, f := range tc.files {
				if err := os.WriteFile(filepath.Join(dir, f), []byte("test:\n    echo ok\n"), 0644); err != nil {
					t.Fatal(err)
				}
			}
			if got := hasJustfile(dir); got != tc.want {
				t.Errorf("hasJustfile() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHasJustRecipe(t *testing.T) {
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed, skipping")
	}

	t.Run("recipe exists", func(t *testing.T) {
		dir := t.TempDir()
		content := "test:\n    echo ok\n"
		if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		if !hasJustRecipe(dir, "test") {
			t.Error("hasJustRecipe() = false, want true for existing recipe")
		}
	})

	t.Run("recipe does not exist", func(t *testing.T) {
		dir := t.TempDir()
		content := "build:\n    echo build\n"
		if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		if hasJustRecipe(dir, "test") {
			t.Error("hasJustRecipe() = true, want false for missing recipe")
		}
	})

	t.Run("no justfile", func(t *testing.T) {
		dir := t.TempDir()
		if hasJustRecipe(dir, "test") {
			t.Error("hasJustRecipe() = true, want false when no justfile")
		}
	})
}

func TestWriteLatestMd(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
		t.Fatal(err)
	}

	t.Run("fail status", func(t *testing.T) {
		stats := TestStats{Fail: 1}

		err := writeLatestMd(dir, "test", stats)
		if err != nil {
			t.Fatalf("writeLatestMd() error = %v", err)
		}

		resultsDir := filepath.Join(dir, feature.GetFeatureTestingResultsDir("test"))
		data, err := os.ReadFile(filepath.Join(resultsDir, "latest.md"))
		if err != nil {
			t.Fatalf("failed to read latest.md: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "FAIL") {
			t.Error("latest.md should show FAIL status")
		}
		if !strings.Contains(content, "raw-output.txt") {
			t.Error("latest.md should reference raw-output.txt")
		}
	})

	t.Run("pass status", func(t *testing.T) {
		stats := TestStats{Total: 5, Pass: 5}

		err := writeLatestMd(dir, "test", stats)
		if err != nil {
			t.Fatalf("writeLatestMd() error = %v", err)
		}

		resultsDir := filepath.Join(dir, feature.GetFeatureTestingResultsDir("test"))
		data, err := os.ReadFile(filepath.Join(resultsDir, "latest.md"))
		if err != nil {
			t.Fatalf("failed to read latest.md: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "PASS") {
			t.Error("latest.md should show PASS status")
		}
	})
}

func TestWriteRawOutput(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
		t.Fatal(err)
	}

	output := "not ok 1 - Login test\n  Error: expected 200, got 404"
	err := writeRawOutput(dir, "test", output)
	if err != nil {
		t.Fatalf("writeRawOutput() error = %v", err)
	}

	resultsDir := filepath.Join(dir, feature.GetFeatureTestingResultsDir("test"))
	data, err := os.ReadFile(filepath.Join(resultsDir, "raw-output.txt"))
	if err != nil {
		t.Fatalf("failed to read raw-output.txt: %v", err)
	}

	if string(data) != output {
		t.Errorf("raw output mismatch: got %q, want %q", string(data), output)
	}
}

func TestRunCmdCapture(t *testing.T) {
	dir := t.TempDir()

	output, success := runCmdCapture(dir, "echo", "hello")
	if !success {
		t.Error("runCmdCapture() success = false, want true")
	}
	if !strings.Contains(output, "hello") {
		t.Errorf("runCmdCapture() output = %q, want contain hello", output)
	}
}

func TestRunCmdCapture_Failure(t *testing.T) {
	dir := t.TempDir()
	_, success := runCmdCapture(dir, "false")
	if success {
		t.Error("runCmdCapture() success = true, want false for failing command")
	}
}

