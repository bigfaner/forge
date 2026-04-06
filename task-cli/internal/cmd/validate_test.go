package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"task-cli/pkg/task"
)

func TestValidator_ValidateTasks(t *testing.T) {
	tests := []struct {
		name            string
		tasks           map[string]task.Task
		wantErrors      int
		wantErrContains []string
	}{
		{
			name: "valid tasks",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "pending", Priority: "P0", File: "task.md"},
			},
			wantErrors: 0,
		},
		{
			name: "missing id",
			tasks: map[string]task.Task{
				"task1": {Title: "Task 1", File: "task.md"},
			},
			wantErrors:      1,
			wantErrContains: []string{"missing 'id'"},
		},
		{
			name: "missing title",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", File: "task.md"},
			},
			wantErrors:      1,
			wantErrContains: []string{"missing 'title'"},
		},
		{
			name: "missing file",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1"},
			},
			wantErrors:      1,
			wantErrContains: []string{"missing 'file'"},
		},
		{
			name: "invalid status",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", File: "task.md", Status: "invalid"},
			},
			wantErrors:      1,
			wantErrContains: []string{"invalid status"},
		},
		{
			name: "invalid priority",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", File: "task.md", Priority: "P5"},
			},
			wantErrors:      1,
			wantErrContains: []string{"invalid priority"},
		},
		{
			name: "multiple errors",
			tasks: map[string]task.Task{
				"task1": {Status: "bad", Priority: "bad"},
			},
			wantErrors: 5, // missing id, title, file, invalid status, invalid priority
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &validator{}
			v.validateTasks(tt.tasks)
			if len(v.errors) != tt.wantErrors {
				t.Errorf("validateTasks() got %d errors, want %d: %v", len(v.errors), tt.wantErrors, v.errors)
			}
			for _, want := range tt.wantErrContains {
				found := false
				for _, err := range v.errors {
					if contains(err, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("validateTasks() missing error containing %q, got %v", want, v.errors)
				}
			}
		})
	}
}

func TestValidator_ValidateDependencies(t *testing.T) {
	tests := []struct {
		name            string
		tasks           map[string]task.Task
		wantErrors      int
		wantErrContains []string
	}{
		{
			name: "no dependencies",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1"},
			},
			wantErrors: 0,
		},
		{
			name: "valid exact dependency",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{}},
				"task2": {ID: "1.2", Dependencies: []string{"1.1"}},
			},
			wantErrors: 0,
		},
		{
			name: "missing dependency",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{"0.1"}},
			},
			wantErrors:      1,
			wantErrContains: []string{"dependency '0.1' not found"},
		},
		{
			name: "valid wildcard dependency",
			tasks: map[string]task.Task{
				"task1": {ID: "0.1", Dependencies: []string{}},
				"task2": {ID: "1.1", Dependencies: []string{"0.x"}},
			},
			wantErrors: 0,
		},
		{
			name: "wildcard matches nothing",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{"9.x"}},
			},
			wantErrors:      1,
			wantErrContains: []string{"wildcard '9.x' matches nothing"},
		},
		{
			name: "multiple missing dependencies",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{"0.1", "0.2"}},
			},
			wantErrors: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &validator{}
			v.validateDependencies(tt.tasks)
			if len(v.errors) != tt.wantErrors {
				t.Errorf("validateDependencies() got %d errors, want %d: %v", len(v.errors), tt.wantErrors, v.errors)
			}
			for _, want := range tt.wantErrContains {
				found := false
				for _, err := range v.errors {
					if contains(err, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("validateDependencies() missing error containing %q, got %v", want, v.errors)
				}
			}
		})
	}
}

func TestValidator_ValidateCircularDeps(t *testing.T) {
	tests := []struct {
		name       string
		tasks      map[string]task.Task
		wantErrors int
	}{
		{
			name: "no circular deps",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{}},
				"task2": {ID: "1.2", Dependencies: []string{"1.1"}},
				"task3": {ID: "1.3", Dependencies: []string{"1.2"}},
			},
			wantErrors: 0,
		},
		{
			name: "simple cycle",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{"1.2"}},
				"task2": {ID: "1.2", Dependencies: []string{"1.1"}},
			},
			wantErrors: 1,
		},
		{
			name: "longer cycle",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{"1.2"}},
				"task2": {ID: "1.2", Dependencies: []string{"1.3"}},
				"task3": {ID: "1.3", Dependencies: []string{"1.1"}},
			},
			wantErrors: 1,
		},
		{
			name: "self cycle",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{"1.1"}},
			},
			wantErrors: 1,
		},
		{
			name: "wildcard ignored for circular check",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{"1.x"}},
			},
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &validator{}
			v.validateCircularDeps(tt.tasks)
			if len(v.errors) != tt.wantErrors {
				t.Errorf("validateCircularDeps() got %d errors, want %d: %v", len(v.errors), tt.wantErrors, v.errors)
			}
		})
	}
}

func TestValidator_ValidateFilesExist(t *testing.T) {
	t.Run("file exists", func(t *testing.T) {
		dir := t.TempDir()
		featureSlug := "test-feature"

		// Create correct directory structure: .../docs/features/test-feature/tasks/
		tasksDir := filepath.Join(dir, "docs", "features", featureSlug, "tasks")
		if err := os.MkdirAll(tasksDir, 0755); err != nil {
			t.Fatal(err)
		}
		taskFile := filepath.Join(tasksDir, "task.md")
		if err := os.WriteFile(taskFile, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}

		// FilePath must be in format: .../docs/features/<slug>/tasks/index.json
		v := &validator{filePath: filepath.Join(dir, "docs", "features", featureSlug, "tasks", "index.json")}
		v.validateFilesExist(featureSlug, map[string]task.Task{
			"task1": {ID: "1.1", File: "task.md"},
		})
		if len(v.warnings) != 0 {
			t.Errorf("expected no warnings, got %v", v.warnings)
		}
	})

	t.Run("file missing", func(t *testing.T) {
		dir := t.TempDir()
		featureSlug := "test-feature"

		// FilePath must be in format: .../docs/features/<slug>/tasks/index.json
		v := &validator{filePath: filepath.Join(dir, "docs", "features", featureSlug, "tasks", "index.json")}
		v.validateFilesExist(featureSlug, map[string]task.Task{
			"task1": {ID: "1.1", File: "missing.md"},
		})
		if len(v.warnings) != 1 {
			t.Errorf("expected 1 warning, got %v", v.warnings)
		}
	})

	t.Run("empty file field skipped", func(t *testing.T) {
		v := &validator{}
		v.validateFilesExist("test-feature", map[string]task.Task{
			"task1": {ID: "1.1", File: ""},
		})
		if len(v.warnings) != 0 {
			t.Errorf("expected no warnings for empty file, got %v", v.warnings)
		}
	})
}

func TestValidator_Run(t *testing.T) {
	t.Run("valid index", func(t *testing.T) {
		dir := t.TempDir()

		index := &task.TaskIndex{
			Feature:      "test-feature",
			PRD:          "prd.md",
			Design:       "design.md",
			StatusEnum:   []string{"pending", "completed"},
			PriorityEnum: []string{"P0", "P1", "P2"},
			Tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "pending", Priority: "P0", File: "task.md"},
			},
		}

		// Create task file
		taskFile := filepath.Join(dir, "task.md")
		if err := os.WriteFile(taskFile, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}

		// Create index file
		indexFile := filepath.Join(dir, "index.json")
		data, err := json.Marshal(index)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(indexFile, data, 0644); err != nil {
			t.Fatal(err)
		}

		v := &validator{filePath: indexFile}
		err = v.run()
		if err != nil {
			t.Errorf("run() returned error: %v", err)
		}
	})

	t.Run("index with errors", func(t *testing.T) {
		dir := t.TempDir()

		index := &task.TaskIndex{
			Feature: "test-feature",
			Tasks: map[string]task.Task{
				"task1": {ID: "", Title: "", File: "", Dependencies: []string{"missing"}},
			},
		}

		// Create index file
		indexFile := filepath.Join(dir, "index.json")
		data, err := json.Marshal(index)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(indexFile, data, 0644); err != nil {
			t.Fatal(err)
		}

		v := &validator{filePath: indexFile}
		err = v.run()
		if err == nil {
			t.Error("expected error for invalid index")
		}
	})

	t.Run("file not found", func(t *testing.T) {
		v := &validator{filePath: "/nonexistent/path/index.json"}
		err := v.run()
		if err == nil {
			t.Error("expected error for nonexistent file")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		dir := t.TempDir()
		indexFile := filepath.Join(dir, "index.json")
		if err := os.WriteFile(indexFile, []byte("not valid json"), 0644); err != nil {
			t.Fatal(err)
		}

		v := &validator{filePath: indexFile}
		err := v.run()
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("missing feature field", func(t *testing.T) {
		dir := t.TempDir()
		indexFile := filepath.Join(dir, "index.json")
		data, _ := json.Marshal(map[string]interface{}{"tasks": map[string]interface{}{}})
		if err := os.WriteFile(indexFile, data, 0644); err != nil {
			t.Fatal(err)
		}

		v := &validator{filePath: indexFile}
		err := v.run()
		if err == nil {
			t.Error("expected error for missing feature field")
		}
	})
}

func TestValidStatusAndPriority(t *testing.T) {
	// Test validStatus map
	for _, s := range []string{"pending", "in_progress", "completed", "blocked", "skipped"} {
		if !validStatus[s] {
			t.Errorf("validStatus[%q] should be true", s)
		}
	}
	if validStatus["invalid"] {
		t.Error("validStatus['invalid'] should be false")
	}

	// Test validPriority map
	for _, p := range []string{"P0", "P1", "P2"} {
		if !validPriority[p] {
			t.Errorf("validPriority[%q] should be true", p)
		}
	}
	if validPriority["P3"] {
		t.Error("validPriority['P3'] should be false")
	}
}

// helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
