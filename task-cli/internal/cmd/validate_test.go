package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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
			wantErrContains: []string{"wildcard '9.x' matches no business tasks"},
		},
		{
			name: "multiple missing dependencies",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{"0.1", "0.2"}},
			},
			wantErrors: 2,
		},
		{
			name: "wildcard skips .gate and .summary tasks",
			tasks: map[string]task.Task{
				"task1":   {ID: "1.1"},
				"gate":    {ID: "1.gate", Breaking: true},
				"summary": {ID: "1.summary"},
				"task2":   {ID: "2.1", Dependencies: []string{"1.x"}},
			},
			wantErrors: 0,
		},
		{
			name: "wildcard matches nothing when only .gate and .summary exist",
			tasks: map[string]task.Task{
				"gate":    {ID: "1.gate", Breaking: true},
				"summary": {ID: "1.summary"},
				"task2":   {ID: "2.1", Dependencies: []string{"1.x"}},
			},
			wantErrors:      1,
			wantErrContains: []string{"wildcard '1.x' matches no business tasks"},
		},
		{
			name: "wildcard matches only business tasks even with gate/summary present",
			tasks: map[string]task.Task{
				"task1":   {ID: "1.1"},
				"task1b":  {ID: "1.2"},
				"gate":    {ID: "1.gate", Breaking: true},
				"summary": {ID: "1.summary"},
				"task2":   {ID: "2.1", Dependencies: []string{"1.x"}},
			},
			wantErrors: 0,
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

	t.Run("T-test-1 with unresolved placeholder", func(t *testing.T) {
		dir := t.TempDir()
		featureSlug := "test-feature"

		// Create correct directory structure
		tasksDir := filepath.Join(dir, "docs", "features", featureSlug, "tasks")
		if err := os.MkdirAll(tasksDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create T-test-1.md with unresolved placeholder
		taskFile := filepath.Join(tasksDir, "T-test-1.md")
		content := `---
id: "T-test-1"
dependencies: [{{LAST_BUSINESS_TASK_ID}}]
---
`
		if err := os.WriteFile(taskFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		v := &validator{filePath: filepath.Join(dir, "docs", "features", featureSlug, "tasks", "index.json")}
		v.validateFilesExist(featureSlug, map[string]task.Task{
			"t-test-1": {ID: "T-test-1", File: "T-test-1.md"},
		})

		if len(v.errors) != 1 {
			t.Errorf("expected 1 error for unresolved placeholder, got %d: %v", len(v.errors), v.errors)
		}
		if len(v.errors) > 0 && !contains(v.errors[0], "{{LAST_BUSINESS_TASK_ID}}") {
			t.Errorf("error should mention placeholder, got: %s", v.errors[0])
		}
	})

	t.Run("T-test-1 with resolved placeholder", func(t *testing.T) {
		dir := t.TempDir()
		featureSlug := "test-feature"

		// Create correct directory structure
		tasksDir := filepath.Join(dir, "docs", "features", featureSlug, "tasks")
		if err := os.MkdirAll(tasksDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create T-test-1.md with resolved placeholder
		taskFile := filepath.Join(tasksDir, "T-test-1.md")
		content := `---
id: "T-test-1"
dependencies: ["1.5"]
---
`
		if err := os.WriteFile(taskFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		v := &validator{filePath: filepath.Join(dir, "docs", "features", featureSlug, "tasks", "index.json")}
		v.validateFilesExist(featureSlug, map[string]task.Task{
			"t-test-1": {ID: "T-test-1", File: "T-test-1.md"},
		})

		if len(v.errors) != 0 {
			t.Errorf("expected no errors for resolved placeholder, got: %v", v.errors)
		}
	})
}


func TestValidator_Run(t *testing.T) {
	t.Run("valid index", func(t *testing.T) {
		dir := t.TempDir()

		index := &task.TaskIndex{
			Feature:      "test-feature",
			PRD:          "prd/prd-spec.md",
			Design:       "design/tech-design.md",
			StatusEnum:   []string{"pending", "completed"},
			PriorityEnum: []string{"P0", "P1", "P2"},
		}
			index.SetTasks(map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "pending", Priority: "P0", File: "task.md"},
			})

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
		}
			index.SetTasks(map[string]task.Task{
				"task1": {ID: "", Title: "", File: "", Dependencies: []string{"missing"}},
			})

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

	t.Run("missing statusEnum warns", func(t *testing.T) {
		dir := t.TempDir()
		indexFile := filepath.Join(dir, "index.json")
		data, _ := json.Marshal(map[string]interface{}{
			"feature": "test-feature",
			"tasks":   map[string]interface{}{},
		})
		if err := os.WriteFile(indexFile, data, 0644); err != nil {
			t.Fatal(err)
		}

		v := &validator{filePath: indexFile}
		_ = v.run()
		found := false
		for _, w := range v.warnings {
			if contains(w, "statusEnum") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected warning about missing statusEnum, got warnings: %v", v.warnings)
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

func TestIsBusinessTask(t *testing.T) {
	tests := []struct {
		id   string
		want bool
	}{
		{"1.1", true},
		{"2.3.1", true},
		{"1.gate", false},
		{"1.summary", false},
		{"10.gate", false},
		{"10.summary", false},
	}
	for _, tt := range tests {
		if got := isBusinessTask(tt.id); got != tt.want {
			t.Errorf("isBusinessTask(%q) = %v, want %v", tt.id, got, tt.want)
		}
	}
}

// --- V1: Wildcard self-dependency ---

func TestValidator_ValidateWildcardSelfDeps(t *testing.T) {
	tests := []struct {
		name             string
		tasks            map[string]task.Task
		wantErrors       int
		wantWarnings     int
		wantErrContains  []string
		wantWarnContains []string
	}{
		{
			name:         "no deps at all",
			tasks:        map[string]task.Task{
				"task1": {ID: "1.1"},
			},
			wantErrors:   0,
			wantWarnings: 0,
		},
		{
			name:         "only non-wildcard deps",
			tasks:        map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{"1.0"}},
				"task0": {ID: "1.0"},
			},
			wantErrors:   0,
			wantWarnings: 0,
		},
		{
			name:         "wildcard on different phase -> skip",
			tasks:        map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{"2.x"}},
				"task2": {ID: "2.1"},
			},
			wantErrors:   0,
			wantWarnings: 0,
		},
		{
			name:            "single business task with self-only wildcard -> ERROR",
			tasks:           map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{"1.x"}},
			},
			wantErrors:      1,
			wantWarnings:    0,
			wantErrContains: []string{"only matches itself", "1.1", "1.x"},
		},
		{
			name:             "wildcard matches self + 1 other -> WARNING",
			tasks:            map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{"1.x"}},
				"task2": {ID: "1.2"},
			},
			wantErrors:       0,
			wantWarnings:     1,
			wantWarnContains: []string{"matches itself plus 1 others"},
		},
		{
			name:             "wildcard matches self + many others -> WARNING",
			tasks:            map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{"1.x"}},
				"task2": {ID: "1.2"},
				"task3": {ID: "1.3"},
				"task4": {ID: "1.4"},
			},
			wantErrors:       0,
			wantWarnings:     1,
			wantWarnContains: []string{"matches itself plus 3 others"},
		},
		{
			name: "multiple tasks with wildcard self-match -> multiple errors",
			tasks: map[string]task.Task{
				"s1": {ID: "1.1", Dependencies: []string{"1.x"}},
				"s2": {ID: "2.1", Dependencies: []string{"2.x"}},
			},
			wantErrors:   2,
			wantWarnings: 0,
		},
		{
			name: "mixed: exact + wildcard deps, wildcard self-matches",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1"},
				"task2": {ID: "1.2", Dependencies: []string{"1.1", "1.x"}},
			},
			wantErrors:   0,
			wantWarnings: 1,
		},
		// Non-business tasks using wildcard: excluded from self-dep check
		{
			name: "summary with wildcard -> no self-match (non-business)",
			tasks: map[string]task.Task{
				"task1":   {ID: "1.1"},
				"summary": {ID: "1.summary", Dependencies: []string{"1.x"}},
			},
			wantErrors:   0,
			wantWarnings: 0,
		},
		{
			name: "gate with wildcard -> no self-match (non-business)",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1"},
				"gate":  {ID: "1.gate", Breaking: true, Dependencies: []string{"1.x"}},
			},
			wantErrors:   0,
			wantWarnings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &validator{}
			v.validateWildcardSelfDeps(tt.tasks)
			if len(v.errors) != tt.wantErrors {
				t.Errorf("errors: got %d, want %d\n%v", len(v.errors), tt.wantErrors, v.errors)
			}
			if len(v.warnings) != tt.wantWarnings {
				t.Errorf("warnings: got %d, want %d\n%v", len(v.warnings), tt.wantWarnings, v.warnings)
			}
			for _, want := range tt.wantErrContains {
				found := false
				for _, e := range v.errors {
					if contains(e, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing error containing %q in %v", want, v.errors)
				}
			}
			for _, want := range tt.wantWarnContains {
				found := false
				for _, w := range v.warnings {
					if contains(w, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing warning containing %q in %v", want, v.warnings)
				}
			}
		})
	}
}

// --- V2: Gate integrity ---

func TestValidator_ValidateGateIntegrity(t *testing.T) {
	tests := []struct {
		name            string
		tasks           map[string]task.Task
		wantErrors      int
		wantErrContains []string
	}{
		{
			name:  "no gates at all -> PASS",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1"},
			},
			wantErrors: 0,
		},
		{
			name: "gate without ID suffix is ignored",
			tasks: map[string]task.Task{
				"notgate": {ID: "2.something", Breaking: true},
			},
			wantErrors: 0,
		},
		{
			name: "gate not breaking is ignored",
			tasks: map[string]task.Task{
				"gate":  {ID: "2.gate", Breaking: false},
				"task2": {ID: "2.1", Dependencies: []string{}},
			},
			wantErrors: 0,
		},
		{
			name: "gate at phase 1 depends on own summary -> PASS",
			tasks: map[string]task.Task{
				"task1":    {ID: "1.1"},
				"summary1": {ID: "1.summary"},
				"gate":     {ID: "1.gate", Breaking: true, Dependencies: []string{"1.summary"}},
				"task2":    {ID: "2.1", Dependencies: []string{"1.gate"}},
			},
			wantErrors: 0,
		},
		{
			name: "gate missing own summary dep -> ERROR",
			tasks: map[string]task.Task{
				"summary1": {ID: "1.summary"},
				"summary2": {ID: "2.summary"},
				"gate":     {ID: "2.gate", Breaking: true, Dependencies: []string{}},
			},
			wantErrors:      1,
			wantErrContains: []string{"must depend on own phase summary '2.summary'"},
		},
		{
			name: "gate depends on own summary -> OK, but next phase task missing gate -> ERROR",
			tasks: map[string]task.Task{
				"summary1": {ID: "1.summary"},
				"summary2": {ID: "2.summary"},
				"gate":     {ID: "2.gate", Breaking: true, Dependencies: []string{"2.summary"}},
				"task3a":   {ID: "3.1", Dependencies: []string{}},
				"task3b":   {ID: "3.2", Dependencies: []string{"2.gate"}},
			},
			wantErrors:      1,
			wantErrContains: []string{"must depend on gate '2.gate'", "3.1"},
		},
		{
			name: "fully correct exit-gate chain -> PASS",
			tasks: map[string]task.Task{
				"task1":    {ID: "1.1"},
				"summary1": {ID: "1.summary"},
				"gate1":    {ID: "1.gate", Breaking: true, Dependencies: []string{"1.summary"}},
				"task2":    {ID: "2.1", Dependencies: []string{"1.gate"}},
				"summary2": {ID: "2.summary"},
				"gate2":    {ID: "2.gate", Breaking: true, Dependencies: []string{"2.summary"}},
			},
			wantErrors: 0,
		},
		{
			name: "next phase task wildcard does NOT satisfy gate dep -> ERROR",
			tasks: map[string]task.Task{
				"summary1": {ID: "1.summary"},
				"summary2": {ID: "2.summary"},
				"gate":     {ID: "2.gate", Breaking: true, Dependencies: []string{"2.summary"}},
				"task3":    {ID: "3.1", Dependencies: []string{"2.x"}},
			},
			wantErrors:      1,
			wantErrContains: []string{"must depend on gate '2.gate'"},
		},
		{
			name: "multiple gates at different phases",
			tasks: map[string]task.Task{
				"task1":    {ID: "1.1"},
				"summary1": {ID: "1.summary"},
				"gate1":    {ID: "1.gate", Breaking: true, Dependencies: []string{"1.summary"}},
				"task2":    {ID: "2.1", Dependencies: []string{"1.gate"}},
				"summary2": {ID: "2.summary"},
				"gate2":    {ID: "2.gate", Breaking: true, Dependencies: []string{"2.summary"}},
				"task3":    {ID: "3.1", Dependencies: []string{"2.gate"}},
			},
			wantErrors: 0,
		},
		{
			name: "multiple gates, second gate broken",
			tasks: map[string]task.Task{
				"task1":    {ID: "1.1"},
				"summary1": {ID: "1.summary"},
				"gate1":    {ID: "1.gate", Breaking: true, Dependencies: []string{"1.summary"}},
				"task2":    {ID: "2.1", Dependencies: []string{"1.gate"}},
				"summary2": {ID: "2.summary"},
				"gate2":    {ID: "2.gate", Breaking: true, Dependencies: []string{}},
				"task3":    {ID: "3.1", Dependencies: []string{}},
			},
			wantErrors:      2,
			wantErrContains: []string{"must depend on own phase summary '2.summary'", "must depend on gate '2.gate'"},
		},
		{
			name: "gate phase 0 -> skipped (non-numeric prefix)",
			tasks: map[string]task.Task{
				"gate":  {ID: "abc.gate", Breaking: true, Dependencies: []string{}},
				"task1": {ID: "abc.1", Dependencies: []string{}},
			},
			wantErrors: 0,
		},
		// Wildcard dep on own phase summary
		{
			name: "gate wildcard does NOT satisfy own summary dep -> ERROR",
			tasks: map[string]task.Task{
				"task1":    {ID: "1.1"},
				"summary1": {ID: "1.summary"},
				"gate":     {ID: "1.gate", Breaking: true, Dependencies: []string{"1.x"}},
				"task2":    {ID: "2.1", Dependencies: []string{"1.gate"}},
			},
			wantErrors:      1,
			wantErrContains: []string{"must depend on own phase summary '1.summary'"},
		},
		{
			name: "gate depends on unrelated wildcard but not own summary -> ERROR",
			tasks: map[string]task.Task{
				"task1":    {ID: "1.1"},
				"summary1": {ID: "1.summary"},
				"gate":     {ID: "1.gate", Breaking: true, Dependencies: []string{"3.x"}},
				"task2":    {ID: "2.1", Dependencies: []string{"1.gate"}},
			},
			wantErrors:      1,
			wantErrContains: []string{"must depend on own phase summary '1.summary'"},
		},
		{
			name: "gate with explicit summary dep + wildcard on same phase -> PASS",
			tasks: map[string]task.Task{
				"task1":    {ID: "1.1"},
				"summary1": {ID: "1.summary"},
				"gate":     {ID: "1.gate", Breaking: true, Dependencies: []string{"1.summary", "1.x"}},
				"task2":    {ID: "2.1", Dependencies: []string{"1.gate"}},
			},
			wantErrors: 0,
		},
		{
			name: "next phase task must have explicit gate dep, wildcard alone fails",
			tasks: map[string]task.Task{
				"task1":    {ID: "1.1"},
				"summary1": {ID: "1.summary"},
				"gate":     {ID: "1.gate", Breaking: true, Dependencies: []string{"1.summary"}},
				"task2":    {ID: "2.1", Dependencies: []string{"1.x"}},
			},
			wantErrors:      1,
			wantErrContains: []string{"must depend on gate '1.gate'"},
		},
		{
			name: "next phase task with explicit gate + wildcard -> PASS",
			tasks: map[string]task.Task{
				"task1":    {ID: "1.1"},
				"summary1": {ID: "1.summary"},
				"gate":     {ID: "1.gate", Breaking: true, Dependencies: []string{"1.summary"}},
				"task2":    {ID: "2.1", Dependencies: []string{"1.gate", "1.x"}},
			},
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &validator{}
			v.validateGateIntegrity(tt.tasks)
			if len(v.errors) != tt.wantErrors {
				t.Errorf("errors: got %d, want %d\n%v", len(v.errors), tt.wantErrors, v.errors)
			}
			for _, want := range tt.wantErrContains {
				found := false
				for _, e := range v.errors {
					if contains(e, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing error containing %q in %v", want, v.errors)
				}
			}
		})
	}
}

// --- V3: Phase order sanity ---

func TestValidator_ValidatePhaseOrder(t *testing.T) {
	tests := []struct {
		name             string
		tasks            map[string]task.Task
		wantWarnings     int
		wantWarnContains []string
	}{
		{
			name:         "phase 1 with empty deps -> no warning",
			tasks:        map[string]task.Task{
				"task1": {ID: "1.1", Dependencies: []string{}},
			},
			wantWarnings: 0,
		},
		{
			name:         "phase 1 with deps -> no warning",
			tasks:        map[string]task.Task{
				"task0": {ID: "1.0"},
				"task1": {ID: "1.1", Dependencies: []string{"1.0"}},
			},
			wantWarnings: 0,
		},
		{
			name:         "phase 2 with cross-phase exact dep -> no warning",
			tasks:        map[string]task.Task{
				"task1": {ID: "1.1"},
				"task2": {ID: "2.1", Dependencies: []string{"1.1"}},
			},
			wantWarnings: 0,
		},
		{
			name:         "phase 2 with cross-phase wildcard dep -> no warning",
			tasks:        map[string]task.Task{
				"task1": {ID: "1.1"},
				"task2": {ID: "2.1", Dependencies: []string{"1.x"}},
			},
			wantWarnings: 0,
		},
		{
			name:         "phase 3 with wildcard on phase 1 -> no warning",
			tasks:        map[string]task.Task{
				"task1": {ID: "1.1"},
				"task3": {ID: "3.1", Dependencies: []string{"1.x"}},
			},
			wantWarnings: 0,
		},
		{
			name:         "phase 2 with only same-phase dep -> WARNING for both",
			tasks:        map[string]task.Task{
				"task2a": {ID: "2.1"},
				"task2b": {ID: "2.2", Dependencies: []string{"2.1"}},
			},
			wantWarnings:     2,
		},
		{
			name:         "phase 2 with no deps at all -> WARNING",
			tasks:        map[string]task.Task{
				"task2": {ID: "2.1", Dependencies: []string{}},
			},
			wantWarnings:     1,
			wantWarnContains: []string{"2.1"},
		},
		{
			name:         "multiple phases, all missing cross-phase deps -> multiple warnings",
			tasks:        map[string]task.Task{
				"task2": {ID: "2.1", Dependencies: []string{}},
				"task3": {ID: "3.1", Dependencies: []string{}},
			},
			wantWarnings: 2,
		},
		{
			name:         "gate task skipped",
			tasks:        map[string]task.Task{
				"gate": {ID: "2.gate", Breaking: true, Dependencies: []string{}},
			},
			wantWarnings: 0,
		},
		{
			name:         "summary task skipped",
			tasks:        map[string]task.Task{
				"summary": {ID: "2.summary", Dependencies: []string{}},
			},
			wantWarnings: 0,
		},
		{
			name:         "mixed: same-phase + cross-phase dep -> no warning",
			tasks:        map[string]task.Task{
				"task1":  {ID: "1.1"},
				"task2a": {ID: "2.1", Dependencies: []string{"1.1"}},
				"task2b": {ID: "2.2", Dependencies: []string{"2.1", "1.1"}},
			},
			wantWarnings: 0,
		},
		{
			name:         "non-numeric phase ID -> skipped (phase=0)",
			tasks:        map[string]task.Task{
				"task": {ID: "abc.1", Dependencies: []string{}},
			},
			wantWarnings: 0,
		},
		{
			name:         "non-numeric wildcard dep -> fallback branch",
			tasks:        map[string]task.Task{
				"task1": {ID: "1.1"},
				"task2": {ID: "2.1", Dependencies: []string{"abc.x"}},
			},
			wantWarnings:     1,
			wantWarnContains: []string{"no dependency on previous phase"},
		},
		{
			name:         "task with nil deps -> WARNING",
			tasks:        map[string]task.Task{
				"task2": {ID: "2.1", Dependencies: nil},
			},
			wantWarnings: 1,
		},
		// Exit-gate convention: gate N.gate in phase N, next phase tasks depend on it
		{
			name:         "business depends on exit gate from prev phase -> no warning",
			tasks:        map[string]task.Task{
				"task1":    {ID: "1.1"},
				"summary1": {ID: "1.summary", Dependencies: []string{"1.x"}},
				"gate1":    {ID: "1.gate", Breaking: true, Dependencies: []string{"1.summary"}},
				"task2":    {ID: "2.1", Dependencies: []string{"1.gate"}},
			},
			wantWarnings: 0,
		},
		{
			name:         "business depends on exit gate via wildcard -> no warning",
			tasks:        map[string]task.Task{
				"task1":    {ID: "1.1"},
				"summary1": {ID: "1.summary"},
				"gate1":    {ID: "1.gate", Breaking: true, Dependencies: []string{"1.summary"}},
				"task2":    {ID: "2.1", Dependencies: []string{"1.x"}},
			},
			wantWarnings: 0,
		},
		{
			name:         "business depends on exit gate only, no other cross-phase -> no warning",
			tasks:        map[string]task.Task{
				"summary1": {ID: "1.summary"},
				"gate1":    {ID: "1.gate", Breaking: true, Dependencies: []string{"1.summary"}},
				"task2":    {ID: "2.1", Dependencies: []string{"1.gate"}},
				"task2b":   {ID: "2.2", Dependencies: []string{"1.gate"}},
			},
			wantWarnings: 0,
		},
		{
			name:         "phase 2 wildcard skips gate and summary in matching",
			tasks:        map[string]task.Task{
				"task1":    {ID: "1.1"},
				"summary1": {ID: "1.summary"},
				"gate1":    {ID: "1.gate", Breaking: true, Dependencies: []string{"1.summary"}},
				"task2":    {ID: "2.1", Dependencies: []string{"1.x"}},
			},
			wantWarnings: 0,
		},
		{
			name:         "phase 2 with same-phase wildcard only -> WARNING",
			tasks:        map[string]task.Task{
				"task2a": {ID: "2.1"},
				"task2b": {ID: "2.2", Dependencies: []string{"2.x"}},
			},
			wantWarnings:     2,
			wantWarnContains: []string{"no dependency on previous phase"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &validator{}
			v.validatePhaseOrder(tt.tasks)
			if len(v.warnings) != tt.wantWarnings {
				t.Errorf("warnings: got %d, want %d\n%v", len(v.warnings), tt.wantWarnings, v.warnings)
			}
			for _, want := range tt.wantWarnContains {
				found := false
				for _, w := range v.warnings {
					if contains(w, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing warning containing %q in %v", want, v.warnings)
				}
			}
		})
	}
}

// --- V4: Phase summary existence ---

func TestValidator_ValidatePhaseSummaries(t *testing.T) {
	tests := []struct {
		name             string
		tasks            map[string]task.Task
		wantWarnings     int
		wantWarnContains []string
	}{
		{
			name:         "empty tasks -> no warning",
			tasks:        map[string]task.Task{},
			wantWarnings: 0,
		},
		{
			name: "phase with summary -> no warning",
			tasks: map[string]task.Task{
				"task1":   {ID: "1.1"},
				"summary": {ID: "1.summary"},
			},
			wantWarnings: 0,
		},
		{
			name: "phase without summary -> WARNING",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1"},
			},
			wantWarnings:     1,
			wantWarnContains: []string{"no '1.summary' task"},
		},
		{
			name: "only summary and gate, no business tasks -> no warning",
			tasks: map[string]task.Task{
				"summary": {ID: "1.summary"},
				"gate":    {ID: "1.gate", Breaking: true},
			},
			wantWarnings: 0,
		},
		{
			name: "multiple phases, one missing summary",
			tasks: map[string]task.Task{
				"task1":    {ID: "1.1"},
				"summary1": {ID: "1.summary"},
				"task2":    {ID: "2.1"},
			},
			wantWarnings:     1,
			wantWarnContains: []string{"no '2.summary' task"},
		},
		{
			name: "multiple phases, all missing summaries",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1"},
				"task2": {ID: "2.1"},
			},
			wantWarnings: 2,
		},
		{
			name: "multiple phases, all have summaries -> no warning",
			tasks: map[string]task.Task{
				"task1":    {ID: "1.1"},
				"summary1": {ID: "1.summary"},
				"task2":    {ID: "2.1"},
				"summary2": {ID: "2.summary"},
			},
			wantWarnings: 0,
		},
		// Gate/summary compat: gate doesn't count as business task
		{
			name: "phase with only gate and summary, no business -> no warning",
			tasks: map[string]task.Task{
				"gate":    {ID: "1.gate", Breaking: true},
				"summary": {ID: "1.summary"},
			},
			wantWarnings: 0,
		},
		{
			name: "phase has gate+summary+business but no summary for gate phase",
			tasks: map[string]task.Task{
				"summary1": {ID: "1.summary"},
				"gate2":    {ID: "2.gate", Breaking: true},
				"task2":    {ID: "2.1"},
			},
			wantWarnings:     1,
			wantWarnContains: []string{"no '2.summary' task"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &validator{}
			v.validatePhaseSummaries(tt.tasks)
			if len(v.warnings) != tt.wantWarnings {
				t.Errorf("warnings: got %d, want %d\n%v", len(v.warnings), tt.wantWarnings, v.warnings)
			}
			for _, want := range tt.wantWarnContains {
				found := false
				for _, w := range v.warnings {
					if contains(w, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing warning containing %q in %v", want, v.warnings)
				}
			}
		})
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

func TestRunValidate_WithFileArg(t *testing.T) {
	t.Run("valid file via arg", func(t *testing.T) {
		dir := t.TempDir()
		indexFile := filepath.Join(dir, "index.json")
		data, _ := json.Marshal(map[string]interface{}{
			"feature": "test",
			"tasks":   map[string]interface{}{},
		})
		os.WriteFile(indexFile, data, 0644)

		// Should not exit (would kill test process)
		runValidate(nil, []string{indexFile})
	})
}

func TestValidator_ValidateLiveness(t *testing.T) {
	tests := []struct {
		name            string
		tasks           map[string]task.Task
		wantWarnings    int
		wantErrors      int
		wantWarnContains []string
		wantErrContains  []string
	}{
		{
			name: "non-blocked task produces no warnings",
			tasks: map[string]task.Task{
				"a": {ID: "a", Status: "pending"},
			},
			wantWarnings: 0,
			wantErrors:   0,
		},
		{
			name: "blocked with no deps warns orphaned",
			tasks: map[string]task.Task{
				"a": {ID: "a", Status: "blocked"},
			},
			wantWarnings:    1,
			wantWarnContains: []string{"orphaned"},
		},
		{
			name: "blocked with all deps completed warns stale",
			tasks: map[string]task.Task{
				"a": {ID: "a", Status: "blocked", Dependencies: []string{"b"}},
				"b": {ID: "b", Status: "completed"},
			},
			wantWarnings:    1,
			wantWarnContains: []string{"stale"},
		},
		{
			name: "blocked with all deps skipped warns stale",
			tasks: map[string]task.Task{
				"a": {ID: "a", Status: "blocked", Dependencies: []string{"b"}},
				"b": {ID: "b", Status: "skipped"},
			},
			wantWarnings:    1,
			wantWarnContains: []string{"stale"},
		},
		{
			name: "blocked with active dep produces no warning",
			tasks: map[string]task.Task{
				"a": {ID: "a", Status: "blocked", Dependencies: []string{"b"}},
				"b": {ID: "b", Status: "in_progress"},
			},
			wantWarnings: 0,
			wantErrors:   0,
		},
		{
			name: "blocked with pending dep produces no warning",
			tasks: map[string]task.Task{
				"a": {ID: "a", Status: "blocked", Dependencies: []string{"b"}},
				"b": {ID: "b", Status: "pending"},
			},
			wantWarnings: 0,
			wantErrors:   0,
		},
		{
			name: "blocked on missing dep errors and warns no path",
			tasks: map[string]task.Task{
				"a": {ID: "a", Status: "blocked", Dependencies: []string{"missing"}},
			},
			wantErrors:      1,
			wantWarnings:    1,
			wantErrContains: []string{"missing dependency"},
			wantWarnContains: []string{"no path to resolution"},
		},
		{
			name: "blocked chain with no path to resolution warns",
			tasks: map[string]task.Task{
				"a": {ID: "a", Status: "blocked", Dependencies: []string{"b"}},
				"b": {ID: "b", Status: "blocked"},
			},
			wantWarnings:    2, // a: no path, b: orphaned (no deps)
			wantWarnContains: []string{"no path to resolution", "orphaned"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &validator{}
			v.validateLiveness(tt.tasks)
			if len(v.warnings) != tt.wantWarnings {
				t.Errorf("got %d warnings, want %d: %v", len(v.warnings), tt.wantWarnings, v.warnings)
			}
			if len(v.errors) != tt.wantErrors {
				t.Errorf("got %d errors, want %d: %v", len(v.errors), tt.wantErrors, v.errors)
			}
			for _, want := range tt.wantWarnContains {
				found := false
				for _, w := range v.warnings {
					if strings.Contains(w, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected warning containing %q, got %v", want, v.warnings)
				}
			}
			for _, want := range tt.wantErrContains {
				found := false
				for _, e := range v.errors {
					if strings.Contains(e, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error containing %q, got %v", want, v.errors)
				}
			}
		})
	}
}

func TestValidator_ValidateLiveness_Wildcard(t *testing.T) {
	tests := []struct {
		name             string
		tasks            map[string]task.Task
		wantWarnings     int
		wantErrors       int
		wantWarnContains []string
	}{
		{
			name: "wildcard deps all completed warns stale",
			tasks: map[string]task.Task{
				"a":   {ID: "a", Status: "blocked", Dependencies: []string{"1.x"}},
				"1.1": {ID: "1.1", Status: "completed"},
				"1.2": {ID: "1.2", Status: "completed"},
			},
			wantWarnings:     1,
			wantWarnContains: []string{"stale"},
		},
		{
			name: "wildcard deps with pending task no warning",
			tasks: map[string]task.Task{
				"a":   {ID: "a", Status: "blocked", Dependencies: []string{"1.x"}},
				"1.1": {ID: "1.1", Status: "completed"},
				"1.2": {ID: "1.2", Status: "pending"},
			},
			wantWarnings: 0,
			wantErrors:   0,
		},
		{
			name: "wildcard no matches no error",
			tasks: map[string]task.Task{
				"a": {ID: "a", Status: "blocked", Dependencies: []string{"9.x"}},
			},
			wantWarnings:     1,
			wantWarnContains: []string{"stale"},
		},
		{
			name: "mixed wildcard and exact deps",
			tasks: map[string]task.Task{
				"a":     {ID: "a", Status: "blocked", Dependencies: []string{"1.x", "fix-1"}},
				"1.1":   {ID: "1.1", Status: "completed"},
				"fix-1": {ID: "fix-1", Status: "completed"},
			},
			wantWarnings:     1,
			wantWarnContains: []string{"stale"},
		},
		{
			name: "wildcard self-exclusion: blocked task matching own wildcard",
			tasks: map[string]task.Task{
				"1.3": {ID: "1.3", Status: "blocked", Dependencies: []string{"1.x"}},
				"1.1": {ID: "1.1", Status: "completed"},
				"1.2": {ID: "1.2", Status: "completed"},
			},
			wantWarnings:     1,
			wantWarnContains: []string{"stale"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &validator{}
			v.validateLiveness(tt.tasks)
			if len(v.warnings) != tt.wantWarnings {
				t.Errorf("got %d warnings, want %d: %v", len(v.warnings), tt.wantWarnings, v.warnings)
			}
			if len(v.errors) != tt.wantErrors {
				t.Errorf("got %d errors, want %d: %v", len(v.errors), tt.wantErrors, v.errors)
			}
			for _, want := range tt.wantWarnContains {
				found := false
				for _, w := range v.warnings {
					if strings.Contains(w, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected warning containing %q, got %v", want, v.warnings)
				}
			}
		})
	}
}

// --- Quick mode validation ---

func TestValidator_QuickMode(t *testing.T) {
	t.Run("quick mode with proposal field suppresses prd/design/summary warnings", func(t *testing.T) {
		dir := t.TempDir()

		index := &task.TaskIndex{
			Feature:      "test-quick",
			Proposal:     "docs/proposals/test-quick/proposal.md",
			StatusEnum:   []string{"pending", "in_progress", "completed", "blocked", "skipped"},
			PriorityEnum: []string{"P0", "P1", "P2"},
		}
		index.SetTasks(map[string]task.Task{
			"task1": {ID: "1", Title: "Task 1", Status: "pending", Priority: "P0", File: "1-task.md"},
			"task2": {ID: "2", Title: "Task 2", Status: "pending", Priority: "P0", Dependencies: []string{"1"}, File: "2-task.md"},
			"quick-test-cases": {ID: "T-quick-1", Title: "Test Cases", Status: "pending", Priority: "P1", Dependencies: []string{"2"}, File: "quick-test-cases.md"},
		})

		for _, fname := range []string{"1-task.md", "2-task.md", "quick-test-cases.md"} {
			if err := os.WriteFile(filepath.Join(dir, fname), []byte("content"), 0644); err != nil {
				t.Fatal(err)
			}
		}

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

		for _, w := range v.warnings {
			if contains(w, "prd") {
				t.Errorf("quick mode should not warn about missing prd: %s", w)
			}
			if contains(w, "design") {
				t.Errorf("quick mode should not warn about missing design: %s", w)
			}
			if contains(w, "summary") {
				t.Errorf("quick mode should not warn about missing summary: %s", w)
			}
			if contains(w, "previous phase") {
				t.Errorf("quick mode should not warn about phase order: %s", w)
			}
		}
	})

	t.Run("quick mode deserializes proposal field", func(t *testing.T) {
		dir := t.TempDir()

		rawIndex := map[string]interface{}{
			"feature":      "test-quick",
			"proposal":     "docs/proposals/test-quick/proposal.md",
			"statusEnum":   []string{"pending", "completed"},
			"priorityEnum": []string{"P0", "P1", "P2"},
			"tasks": map[string]interface{}{
				"task1": map[string]interface{}{
					"id": "1", "title": "Task", "status": "pending", "priority": "P0", "file": "task.md", "scope": "all",
				},
			},
		}
		data, _ := json.Marshal(rawIndex)
		indexFile := filepath.Join(dir, "index.json")
		if err := os.WriteFile(indexFile, data, 0644); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(filepath.Join(dir, "task.md"), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}

		v := &validator{filePath: indexFile}
		err := v.run()
		if err != nil {
			t.Errorf("run() returned error: %v", err)
		}
		if !v.quickMode {
			t.Error("expected quickMode to be true when proposal is set and prd/design are empty")
		}
	})

	t.Run("quick mode round-trips proposal through marshal/unmarshal", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature:      "test-quick",
			Proposal:     "docs/proposals/test-quick/proposal.md",
			StatusEnum:   []string{"pending"},
			PriorityEnum: []string{"P0"},
		}
		index.SetTasks(map[string]task.Task{
			"t1": {ID: "1", Title: "T", File: "t.md"},
		})

		data, err := json.Marshal(index)
		if err != nil {
			t.Fatal(err)
		}

		var restored task.TaskIndex
		if err := json.Unmarshal(data, &restored); err != nil {
			t.Fatal(err)
		}
		if restored.Proposal != "docs/proposals/test-quick/proposal.md" {
			t.Errorf("proposal not preserved: got %q", restored.Proposal)
		}
		if restored.PRD != "" || restored.Design != "" {
			t.Errorf("prd/design should be empty: prd=%q design=%q", restored.PRD, restored.Design)
		}
	})

	t.Run("full mode with prd+design+proposal has no quick mode suppression", func(t *testing.T) {
		dir := t.TempDir()

		index := &task.TaskIndex{
			Feature:      "test-feature",
			PRD:          "prd/prd-spec.md",
			Design:       "design/tech-design.md",
			Proposal:     "docs/proposals/test/proposal.md",
			StatusEnum:   []string{"pending", "completed"},
			PriorityEnum: []string{"P0", "P1", "P2"},
		}
		index.SetTasks(map[string]task.Task{
			"task1": {ID: "1.1", Title: "Task 1", Status: "pending", Priority: "P0", File: "task.md"},
		})

		taskFile := filepath.Join(dir, "task.md")
		if err := os.WriteFile(taskFile, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}

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
		if v.quickMode {
			t.Error("should not be quick mode when prd and design are both set")
		}
	})
}

func TestValidator_QuickMode_FirstTestTaskPlaceholder(t *testing.T) {
	t.Run("T-quick-1 with unresolved placeholder", func(t *testing.T) {
		dir := t.TempDir()
		featureSlug := "test-feature"

		tasksDir := filepath.Join(dir, "docs", "features", featureSlug, "tasks")
		if err := os.MkdirAll(tasksDir, 0755); err != nil {
			t.Fatal(err)
		}

		taskFile := filepath.Join(tasksDir, "quick-test-cases.md")
		content := `---
id: "T-quick-1"
dependencies: [{{T_QUICK_1_DEP}}]
---
`
		if err := os.WriteFile(taskFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		v := &validator{filePath: filepath.Join(dir, "docs", "features", featureSlug, "tasks", "index.json")}
		v.validateFilesExist(featureSlug, map[string]task.Task{
			"quick-test-cases": {ID: "T-quick-1", File: "quick-test-cases.md"},
		})

		if len(v.errors) != 1 {
			t.Errorf("expected 1 error for unresolved placeholder, got %d: %v", len(v.errors), v.errors)
		}
		if len(v.errors) > 0 && !contains(v.errors[0], "{{T_QUICK_1_DEP}}") {
			t.Errorf("error should mention placeholder, got: %s", v.errors[0])
		}
	})

	t.Run("T-quick-1 with resolved placeholder", func(t *testing.T) {
		dir := t.TempDir()
		featureSlug := "test-feature"

		tasksDir := filepath.Join(dir, "docs", "features", featureSlug, "tasks")
		if err := os.MkdirAll(tasksDir, 0755); err != nil {
			t.Fatal(err)
		}

		taskFile := filepath.Join(tasksDir, "quick-test-cases.md")
		content := `---
id: "T-quick-1"
dependencies: ["2"]
---
`
		if err := os.WriteFile(taskFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		v := &validator{filePath: filepath.Join(dir, "docs", "features", featureSlug, "tasks", "index.json")}
		v.validateFilesExist(featureSlug, map[string]task.Task{
			"quick-test-cases": {ID: "T-quick-1", File: "quick-test-cases.md"},
		})

		if len(v.errors) != 0 {
			t.Errorf("expected no errors for resolved placeholder, got: %v", v.errors)
		}
	})
}

func TestValidator_ValidateTasks_RejectedStatusValid(t *testing.T) {
	v := &validator{filePath: "test.json"}
	tasks := map[string]task.Task{
		"a": {ID: "1.1", Title: "Task", File: "1.1.md", Status: "rejected"},
	}
	v.validateTasks(tasks)
	if len(v.errors) != 0 {
		t.Errorf("rejected should be a valid status, got errors: %v", v.errors)
	}
}

func TestValidator_ValidateLiveness_BlockedWithRejectedDep(t *testing.T) {
	v := &validator{filePath: "test.json"}
	tasks := map[string]task.Task{
		"blocked-task": {ID: "1.1", Status: "blocked", Dependencies: []string{"1.2"}},
		"rejected-dep": {ID: "1.2", Status: "rejected"},
	}
	v.validateLiveness(tasks)
	found := false
	for _, w := range v.warnings {
		if strings.Contains(w, "no path to resolution") {
			found = true
		}
	}
	if !found {
		t.Errorf("blocked on rejected dep should warn 'no path to resolution', got: %v", v.warnings)
	}
}

func TestValidator_ValidateLiveness_RejectedTaskNotFlagged(t *testing.T) {
	v := &validator{filePath: "test.json"}
	tasks := map[string]task.Task{
		"rejected-task": {ID: "1.1", Status: "rejected", Dependencies: []string{"1.2"}},
		"completed-dep": {ID: "1.2", Status: "completed"},
	}
	v.validateLiveness(tasks)
	for _, w := range v.warnings {
		if strings.Contains(w, "1.1") {
			t.Errorf("rejected task should not be flagged by liveness checks, got: %s", w)
		}
	}
}

func TestValidator_ValidateLiveness_BlockedOnRejectedViaWildcard(t *testing.T) {
	v := &validator{filePath: "test.json"}
	tasks := map[string]task.Task{
		"blocked-task": {ID: "2.1", Status: "blocked", Dependencies: []string{"1.x"}},
		"1.1":          {ID: "1.1", Status: "completed"},
		"1.2":          {ID: "1.2", Status: "rejected"},
	}
	v.validateLiveness(tasks)
	found := false
	for _, w := range v.warnings {
		if strings.Contains(w, "no path to resolution") {
			found = true
		}
	}
	if !found {
		t.Errorf("blocked on wildcard with rejected task should warn, got: %v", v.warnings)
	}
}
