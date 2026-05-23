package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"
)

func TestListCmd_Metadata(t *testing.T) {
	t.Run("command metadata", func(t *testing.T) {
		if listCmd.Use != "list" {
			t.Errorf("Use = %q, want %q", listCmd.Use, "list")
		}
		if listCmd.Short == "" {
			t.Error("Short is empty")
		}
		if listCmd.Args == nil {
			t.Error("Args is nil, expected cobra.NoArgs")
		}
	})

	t.Run("no args rejected by cobra.NoArgs", func(t *testing.T) {
		err := listCmd.Args(listCmd, []string{"extra"})
		if err == nil {
			t.Error("expected error for extra args, got nil")
		}
	})

	t.Run("registered under task parent", func(t *testing.T) {
		found := false
		for _, sub := range Cmd.Commands() {
			if sub.Name() == "list" {
				found = true
				break
			}
		}
		if !found {
			t.Error("list not registered as subcommand of task")
		}
	})
}

func TestListCmd_Output(t *testing.T) {
	t.Run("displays tasks in table format", func(t *testing.T) {
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: "Add feature set subcommand", Type: "coding.feature", Status: "completed"},
			"2": {ID: "2", Title: "Add feature list subcommand", Type: "coding.feature", Status: "in_progress"},
		}
		_ = setupFullProject(t, SetupOpts{
			Tasks:       tasks,
			FeatureName: "my-feature",
		})

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		// Check header with count and feature slug
		if !strings.Contains(output, "2 found") {
			t.Errorf("output should contain '2 found', got:\n%s", output)
		}
		if !strings.Contains(output, "feature: my-feature") {
			t.Errorf("output should contain 'feature: my-feature', got:\n%s", output)
		}

		// Check table header row
		if !strings.Contains(output, "ID") || !strings.Contains(output, "TYPE") ||
			!strings.Contains(output, "TITLE") || !strings.Contains(output, "STATUS") {
			t.Errorf("output should contain column headers ID TYPE TITLE STATUS, got:\n%s", output)
		}

		// Check task data present
		if !strings.Contains(output, "Add feature set subcommand") {
			t.Errorf("output should contain task title, got:\n%s", output)
		}
		if !strings.Contains(output, "completed") {
			t.Errorf("output should contain 'completed', got:\n%s", output)
		}
		if !strings.Contains(output, "in_progress") {
			t.Errorf("output should contain 'in_progress', got:\n%s", output)
		}
	})

	t.Run("shows no tasks message for empty feature", func(t *testing.T) {
		_ = setupFullProject(t, SetupOpts{
			Tasks:       map[string]task.Task{},
			FeatureName: "empty-feat",
		})

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		if !strings.Contains(output, "no tasks found") {
			t.Errorf("output should contain 'no tasks found' for empty feature, got:\n%s", output)
		}
	})

	t.Run("shows no tasks message when index.json missing", func(t *testing.T) {
		// Use setupFullProject with empty tasks to create a proper feature context,
		// then delete the index.json to simulate missing file.
		// Use a feature name that matches the test's temp directory context.
		featureName := "noindex"
		dir := setupFullProject(t, SetupOpts{
			Tasks:       map[string]task.Task{},
			FeatureName: featureName,
		})

		indexPath := filepath.Join(dir, feature.GetFeatureIndexFile(featureName))
		if err := os.Remove(indexPath); err != nil {
			t.Fatalf("failed to remove index.json: %v", err)
		}

		// Recreate the feature directory (EnsureFeatureDir already created it but
		// the feature resolution needs a single feature dir with index.json or state).
		// Write a minimal state.json so RequireFeature can resolve the feature.
		stateDir := feature.GetProcessDir(dir, featureName)
		stateData := `{"feature":"noindex"}`
		if err := os.WriteFile(filepath.Join(stateDir, "state.json"), []byte(stateData), 0644); err != nil {
			t.Fatalf("failed to write state.json: %v", err)
		}

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		if !strings.Contains(output, "no tasks found") {
			t.Errorf("output should contain 'no tasks found' for missing index, got:\n%s", output)
		}
	})
}

func TestListCmd_Sorting(t *testing.T) {
	t.Run("sorts numeric IDs before test/gate IDs", func(t *testing.T) {
		tasks := map[string]task.Task{
			"T-1": {ID: "T-1", Title: "Compile check", Type: "gate", Status: "pending"},
			"3":   {ID: "3", Title: "Third task", Type: "coding.feature", Status: "pending"},
			"1":   {ID: "1", Title: "First task", Type: "coding.feature", Status: "completed"},
			"T-2": {ID: "T-2", Title: "Test run", Type: "test.run", Status: "pending"},
			"2":   {ID: "2", Title: "Second task", Type: "coding.enhancement", Status: "in_progress"},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		// Extract lines that look like task rows (contain status keywords)
		lines := strings.Split(output, "\n")
		var taskLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			// Skip header/separator lines
			if strings.HasPrefix(trimmed, "---") || strings.HasPrefix(trimmed, "ID ") ||
				strings.HasPrefix(trimmed, "──") || strings.Contains(trimmed, "found") {
				continue
			}
			// Task data rows start with a number or T-
			if len(trimmed) > 0 && (trimmed[0] >= '0' && trimmed[0] <= '9' || strings.HasPrefix(trimmed, "T-")) {
				taskLines = append(taskLines, trimmed)
			}
		}

		if len(taskLines) < 5 {
			t.Fatalf("expected at least 5 task lines, got %d\noutput:\n%s", len(taskLines), output)
		}

		// Verify order: 1, 2, 3, T-1, T-2
		expectedOrder := []string{"1", "2", "3", "T-1", "T-2"}
		for i, expectedID := range expectedOrder {
			if !strings.HasPrefix(taskLines[i], expectedID) {
				t.Errorf("task line %d: expected to start with %q, got %q\nfull output:\n%s", i, expectedID, taskLines[i], output)
			}
		}
	})
}

func TestListCmd_ColumnAlignment(t *testing.T) {
	t.Run("bug: columns misalign when IDs exceed fixed width", func(t *testing.T) {
		tasks := map[string]task.Task{
			"1":         {ID: "1", Title: "First task", Type: "coding.feature", Status: "completed"},
			"1.summary": {ID: "1.summary", Title: "Summary", Type: "doc.summary", Status: "pending"},
			"1.gate":    {ID: "1.gate", Title: "Phase Gate", Type: "gate", Status: "pending"},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		// Find the separator line to determine column widths
		lines := strings.Split(output, "\n")
		var sepLine string
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "---") {
				sepLine = line
				break
			}
		}
		if sepLine == "" {
			t.Fatalf("separator line not found in output:\n%s", output)
		}

		segments := strings.Split(sepLine, "  ")
		if len(segments) != 4 {
			t.Fatalf("expected 4 segments in separator, got %d: %q", len(segments), segments)
		}

		// ID column must accommodate the longest ID ("1.summary" = 9 chars)
		idDashCount := len(strings.TrimRight(segments[0], " "))
		if idDashCount < 9 {
			t.Errorf("bug: ID column is %d chars wide but '1.summary' needs 9 — columns misalign\nseparator: %q\noutput:\n%s",
				idDashCount, sepLine, output)
		}

		// TYPE column must accommodate the longest type ("doc.summary" = 11 chars)
		typeDashCount := len(strings.TrimRight(segments[1], " "))
		if typeDashCount < 11 {
			t.Errorf("bug: TYPE column is %d chars wide but 'doc.summary' needs 11\nseparator: %q\noutput:\n%s",
				typeDashCount, sepLine, output)
		}
	})
}

func TestListCmd_TitleTruncation(t *testing.T) {
	t.Run("truncates long title", func(t *testing.T) {
		longTitle := "This is a very long task title that should be truncated because it exceeds the maximum column width"
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: longTitle, Type: "coding.feature", Status: "pending"},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		if strings.Contains(output, longTitle) {
			t.Errorf("long title should have been truncated, got:\n%s", output)
		}
		if !strings.Contains(output, "...") {
			t.Errorf("truncated title should end with '...', got:\n%s", output)
		}
	})
}

func TestNaturalSortTaskIDs(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"3", "1", "2"},
			expected: []string{"1", "2", "3"},
		},
		{
			input:    []string{"T-2", "T-1", "1"},
			expected: []string{"1", "T-1", "T-2"},
		},
		{
			input:    []string{"10", "2", "1"},
			expected: []string{"1", "2", "10"},
		},
		{
			input:    []string{"1.gate", "1", "2"},
			expected: []string{"1", "1.gate", "2"},
		},
		{
			input:    []string{"T-10", "T-2", "3", "1"},
			expected: []string{"1", "3", "T-2", "T-10"},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("sort %v", tt.input), func(t *testing.T) {
			result := naturalSortTaskIDs(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, result)
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("position %d: expected %q, got %q", i, tt.expected[i], result[i])
				}
			}
		})
	}
}
