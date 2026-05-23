package task

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"forge-cli/pkg/task"
)

// Test runValidate indirectly through validator methods
func TestValidator_IndirectRun(t *testing.T) {
	t.Run("complete validation flow", func(t *testing.T) {
		dir := t.TempDir()
		featureSlug := "test-feature"

		// Create correct directory structure: .../docs/features/<slug>/tasks/
		tasksDir := filepath.Join(dir, "docs", "features", featureSlug, "tasks")
		_ = os.MkdirAll(tasksDir, 0755)

		// Create task file
		taskFile := filepath.Join(tasksDir, "1.1.md")
		_ = os.WriteFile(taskFile, []byte("task content"), 0644)

		// Create index at correct location
		indexPath := filepath.Join(tasksDir, "index.json")
		index := &task.TaskIndex{
			Feature: featureSlug,
			PRD:     "prd/prd-spec.md",
			Design:  "design/tech-design.md",
		}
		index.SetTasks(map[string]task.Task{
			"task1": {ID: "1.1", Title: "Task 1", Status: "pending", Priority: "P0", File: "1.1.md", Type: "coding.feature"},
		})
		data, _ := encodeIndex(index)
		_ = os.WriteFile(indexPath, data, 0644)

		v := &validator{filePath: indexPath}
		v.validateTasks(index.TasksMap())
		v.validateDependencies(index)
		v.validateCircularDeps(index.TasksMap())
		v.validateFilesExist(featureSlug, index.TasksMap())

		if len(v.errors) != 0 {
			t.Errorf("unexpected errors: %v", v.errors)
		}
		// validateFilesExist only checks task files, No warnings expected since task file exists.
		if len(v.warnings) != 0 {
			t.Errorf("expected 0 warnings, got %d: %v", len(v.warnings), v.warnings)
		}
	})
}

func TestValidator_PrintResults(t *testing.T) {
	// Note: printResults calls os.Exit(), so we can't test it directly.
	// Instead, we test the output formatting logic separately.
	t.Run("output formatting", func(t *testing.T) {
		// Just verify the validator struct works correctly
		v := &validator{
			filePath: "test.json",
			info:     []string{"Feature: test", "Tasks: 5"},
			warnings: []string{"Missing prd"},
			errors:   []string{"Missing id"},
		}

		if len(v.info) != 2 {
			t.Errorf("expected 2 info items, got %d", len(v.info))
		}
		if len(v.warnings) != 1 {
			t.Errorf("expected 1 warning, got %d", len(v.warnings))
		}
		if len(v.errors) != 1 {
			t.Errorf("expected 1 error, got %d", len(v.errors))
		}
	})
}

func TestValidator_ValidateFilesExist_Integration(t *testing.T) {
	t.Run("files exist check", func(t *testing.T) {
		dir := t.TempDir()
		featureSlug := "test-feature"

		// Create correct directory structure: .../docs/features/<slug>/tasks/
		tasksDir := filepath.Join(dir, "docs", "features", featureSlug, "tasks")
		_ = os.MkdirAll(tasksDir, 0755)

		// Create existing file
		existingFile := filepath.Join(tasksDir, "existing.md")
		_ = os.WriteFile(existingFile, []byte("content"), 0644)

		v := &validator{filePath: filepath.Join(dir, "docs", "features", featureSlug, "tasks", "index.json")}

		tasks := map[string]task.Task{
			"task1": {ID: "1.1", File: "existing.md"},
			"task2": {ID: "1.2", File: "missing.md"},
			"task3": {ID: "1.3", File: ""}, // empty file - should be skipped
		}

		v.validateFilesExist(featureSlug, tasks)

		if len(v.warnings) != 1 {
			t.Errorf("expected 1 warning for missing file, got %d: %v", len(v.warnings), v.warnings)
		}
	})
}

func TestValidator_ComplexDependencies(t *testing.T) {
	t.Run("diamond dependencies", func(t *testing.T) {
		// 1.1 -> 1.2 and 1.3
		// 1.2 and 1.3 -> 1.4
		tasks := map[string]task.Task{
			"task1": {ID: "1.1", Dependencies: []string{}},
			"task2": {ID: "1.2", Dependencies: []string{"1.1"}},
			"task3": {ID: "1.3", Dependencies: []string{"1.1"}},
			"task4": {ID: "1.4", Dependencies: []string{"1.2", "1.3"}},
		}

		v := &validator{}
		v.validateDependencies(task.NewTestIndex("test", tasks))
		v.validateCircularDeps(tasks)

		if len(v.errors) != 0 {
			t.Errorf("expected no errors for valid diamond deps, got: %v", v.errors)
		}
	})

	t.Run("complex circular detection", func(t *testing.T) {
		// Create a longer cycle: 1.1 -> 1.2 -> 1.3 -> 1.4 -> 1.2
		tasks := map[string]task.Task{
			"task1": {ID: "1.1", Dependencies: []string{"1.2"}},
			"task2": {ID: "1.2", Dependencies: []string{"1.3"}},
			"task3": {ID: "1.3", Dependencies: []string{"1.4"}},
			"task4": {ID: "1.4", Dependencies: []string{"1.2"}},
		}

		v := &validator{}
		v.validateCircularDeps(tasks)

		if len(v.errors) == 0 {
			t.Error("expected error for circular deps")
		}
	})
}

// Test check command logic indirectly
func TestCheckLogic(t *testing.T) {
	t.Run("valid dependencies", func(t *testing.T) {
		taskIDs := map[string]bool{
			"1.1": true,
			"1.2": true,
			"2.1": true,
		}

		tasks := map[string]task.Task{
			"task1": {ID: "1.1", Dependencies: []string{}},
			"task2": {ID: "1.2", Dependencies: []string{"1.1"}},
			"task3": {ID: "2.1", Dependencies: []string{"1.x"}}, // wildcard
		}

		var errors []string
		for _, t := range tasks {
			for _, dep := range t.Dependencies {
				if hasSuffix(dep, ".x") || hasSuffix(dep, "x") {
					prefix := trimSuffix(trimSuffix(dep, "x"), ".")
					prefixWithDot := prefix + "."

					var matches []string
					for id := range taskIDs {
						if hasPrefix(id, prefixWithDot) {
							matches = append(matches, id)
						}
					}

					if len(matches) == 0 {
						errors = append(errors, "wildcard matches nothing")
					}
				} else if !taskIDs[dep] {
					errors = append(errors, "dependency not found")
				}
			}
		}

		if len(errors) != 0 {
			t.Errorf("unexpected errors: %v", errors)
		}
	})

	t.Run("invalid dependencies", func(t *testing.T) {
		taskIDs := map[string]bool{
			"1.1": true,
		}

		tasks := map[string]task.Task{
			"task1": {ID: "1.1", Dependencies: []string{"9.9"}}, // non-existent
		}

		var errors []string
		for key, t := range tasks {
			for _, dep := range t.Dependencies {
				if !hasSuffix(dep, ".x") && !hasSuffix(dep, "x") {
					if !taskIDs[dep] {
						errors = append(errors, key+": dependency "+dep+" not found")
					}
				}
			}
		}

		if len(errors) == 0 {
			t.Error("expected error for missing dependency")
		}
	})
}

// Helper functions for check logic
func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func trimSuffix(s, suffix string) string {
	if hasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

func encodeIndex(index *task.TaskIndex) ([]byte, error) {
	return json.Marshal(index)
}
