package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"task-cli/pkg/feature"
	"task-cli/pkg/task"
)

func TestCheckAllCompleted(t *testing.T) {
	tests := []struct {
		name         string
		tasks        map[string]task.Task
		testCommand  string
		createE2EDir bool
		wantNil      bool
		wantE2EDir   bool
		wantTestCmd  string
	}{
		{
			name: "all completed returns result",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
				"t2": {ID: "1.2", Status: "completed"},
			},
			wantNil: false,
		},
		{
			name: "all skipped returns result",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "skipped"},
			},
			wantNil: false,
		},
		{
			name: "mixed completed and skipped returns result",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
				"t2": {ID: "1.2", Status: "skipped"},
			},
			wantNil: false,
		},
		{
			name: "one pending task returns nil",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
				"t2": {ID: "1.2", Status: "pending"},
			},
			wantNil: true,
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
			name:    "empty task list returns result (vacuously all done)",
			tasks:   map[string]task.Task{},
			wantNil: false,
		},
		{
			name: "e2e scripts dir present is reported",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
			},
			createE2EDir: true,
			wantNil:      false,
			wantE2EDir:   true,
		},
		{
			name: "e2e scripts dir absent gives empty field",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
			},
			createE2EDir: false,
			wantNil:      false,
			wantE2EDir:   false,
		},
		{
			name: "testCommand from index.json is propagated",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
			},
			testCommand: "make test",
			wantNil:     false,
			wantTestCmd: "make test",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv("CLAUDE_PROJECT_DIR", dir)

			// Create feature directory structure
			if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
				t.Fatal(err)
			}

			// Write index.json
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

			// Optionally create e2e scripts dir
			if tc.createE2EDir {
				e2eDir := filepath.Join(dir, feature.GetFeatureTestingScriptsDir("test"))
				if err := os.MkdirAll(e2eDir, 0755); err != nil {
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

			if tc.wantE2EDir && result.E2EScriptsDir == "" {
				t.Error("expected E2EScriptsDir to be set")
			}
			if !tc.wantE2EDir && result.E2EScriptsDir != "" {
				t.Errorf("expected E2EScriptsDir to be empty, got %q", result.E2EScriptsDir)
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
	// Create features dir but no feature subdirectory
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
		name     string
		files    []string // files to create in temp dir
		want     bool
	}{
		{
			name:  "no justfile",
			files: []string{},
			want:  false,
		},
		{
			name:  "lowercase justfile",
			files: []string{"justfile"},
			want:  true,
		},
		{
			name:  "capitalized Justfile",
			files: []string{"Justfile"},
			want:  true,
		},
		{
			name:  "both present",
			files: []string{"justfile", "Justfile"},
			want:  true,
		},
		{
			name:  "unrelated files only",
			files: []string{"Makefile", "go.mod"},
			want:  false,
		},
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
	// Skip if just is not installed
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
