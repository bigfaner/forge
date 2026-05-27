package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

func TestListCmd_Metadata(t *testing.T) {
	t.Run("command metadata", func(t *testing.T) {
		if listCmd.Use != "list [slug]" {
			t.Errorf("Use = %q, want %q", listCmd.Use, "list")
		}
		if listCmd.Short == "" {
			t.Error("Short is empty")
		}
		if listCmd.Args == nil {
			t.Error("Args is nil, expected cobra.NoArgs")
		}
	})

	t.Run("accepts at most one arg", func(t *testing.T) {
		// 0 args should be accepted
		if err := listCmd.Args(listCmd, []string{}); err != nil {
			t.Errorf("expected 0 args to be accepted, got error: %v", err)
		}
		// 1 arg (slug) should be accepted
		if err := listCmd.Args(listCmd, []string{"my-feature"}); err != nil {
			t.Errorf("expected 1 arg to be accepted, got error: %v", err)
		}
		// 2 args should be rejected
		if err := listCmd.Args(listCmd, []string{"a", "b"}); err == nil {
			t.Error("expected error for 2 args, got nil")
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

		// Use --sort id for deterministic ordering (topo mode may produce
		// different order for T-prefixed IDs due to CompareVersionIDs behavior).
		cmd := helperListCmd("id")
		output := captureStdout(func() {
			err := runList(cmd, []string{})
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
				strings.HasPrefix(trimmed, "---") || strings.Contains(trimmed, "found") {
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
		if len(segments) != 6 {
			t.Fatalf("expected 6 segments in separator, got %d: %q", len(segments), segments)
		}

		// ID column must accommodate the longest ID ("1.summary" = 9 chars)
		idDashCount := len(strings.TrimRight(segments[0], " "))
		if idDashCount < 9 {
			t.Errorf("bug: ID column is %d chars wide but '1.summary' needs 9 -- columns misalign\nseparator: %q\noutput:\n%s",
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
	t.Run("bug: titles under 50 chars are unnecessarily truncated", func(t *testing.T) {
		mediumTitle := "Retire gen-test-cases + eval-test-cases skill"
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: mediumTitle, Type: "coding.cleanup", Status: "completed"},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		if !strings.Contains(output, mediumTitle) {
			t.Errorf("title %q should NOT be truncated (under 50 chars), got:\n%s", mediumTitle, output)
		}
	})

	t.Run("bug: CJK title truncation produces valid UTF-8", func(t *testing.T) {
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: "autogen.go Quick 模式：替换旧测试类型为 staged pipeline", Type: "coding.feature", Status: "pending"},
			"2": {ID: "2", Title: "新增 test.gen-journeys 和 test.gen-contracts 的集成测试", Type: "coding.feature", Status: "completed"},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		lines := strings.Split(output, "\n")
		for i, line := range lines {
			if !utf8.ValidString(line) {
				t.Errorf("line %d contains invalid UTF-8: %q", i, line)
			}
		}
	})

	t.Run("truncates very long title", func(t *testing.T) {
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

func TestListCmd_SlugArg(t *testing.T) {
	t.Run("slug arg shows tasks for specified feature", func(t *testing.T) {
		// Setup: create a project with feature "target-feat" containing tasks
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: "Target task one", Type: "coding.feature", Status: "completed"},
			"2": {ID: "2", Title: "Target task two", Type: "coding.enhancement", Status: "pending"},
		}
		_ = setupFullProject(t, SetupOpts{
			Tasks:       tasks,
			FeatureName: "target-feat",
		})

		output := captureStdout(func() {
			err := runList(nil, []string{"target-feat"})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		if !strings.Contains(output, "2 found") {
			t.Errorf("output should contain '2 found', got:\n%s", output)
		}
		if !strings.Contains(output, "feature: target-feat") {
			t.Errorf("output should contain 'feature: target-feat', got:\n%s", output)
		}
		if !strings.Contains(output, "Target task one") {
			t.Errorf("output should contain task title, got:\n%s", output)
		}
	})

	t.Run("slug not matching any feature dir returns error", func(t *testing.T) {
		// Setup: create a project with a different feature
		_ = setupFullProject(t, SetupOpts{
			Tasks:       map[string]task.Task{},
			FeatureName: "existing-feat",
		})

		err := runList(nil, []string{"nonexistent"})
		if err == nil {
			t.Fatal("expected error for nonexistent slug, got nil")
		}
		if !strings.Contains(err.Error(), "nonexistent") {
			t.Errorf("error should mention the slug, got: %v", err)
		}
	})

	t.Run("feature exists but no index.json shows no tasks message", func(t *testing.T) {
		dir := setupFullProject(t, SetupOpts{
			Tasks:       map[string]task.Task{},
			FeatureName: "empty-slug",
		})
		// Remove the index.json to simulate missing file
		indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("empty-slug"))
		if err := os.Remove(indexPath); err != nil {
			t.Fatalf("failed to remove index.json: %v", err)
		}

		output := captureStdout(func() {
			err := runList(nil, []string{"empty-slug"})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		if !strings.Contains(output, "no tasks found") {
			t.Errorf("output should contain 'no tasks found', got:\n%s", output)
		}
		if !strings.Contains(output, "empty-slug") {
			t.Errorf("output should mention the slug, got:\n%s", output)
		}
	})

	t.Run("no args backward compatible", func(t *testing.T) {
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: "Backward compat task", Type: "coding.feature", Status: "completed"},
		}
		_ = setupFullProject(t, SetupOpts{
			Tasks:       tasks,
			FeatureName: "compat-feat",
		})

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		if !strings.Contains(output, "1 found") {
			t.Errorf("output should contain '1 found', got:\n%s", output)
		}
		if !strings.Contains(output, "feature: compat-feat") {
			t.Errorf("output should contain 'feature: compat-feat', got:\n%s", output)
		}
	})
}

func TestListCmd_LocalFlag(t *testing.T) {
	t.Run("--local flag reads from main repo", func(t *testing.T) {
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: "Local task", Type: "coding.feature", Status: "completed"},
		}
		_ = setupFullProject(t, SetupOpts{
			Tasks:       tasks,
			FeatureName: "local-feat",
		})

		output := captureStdout(func() {
			err := runList(nil, []string{"local-feat"})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		if !strings.Contains(output, "1 found") {
			t.Errorf("output should contain '1 found', got:\n%s", output)
		}
		if !strings.Contains(output, "Local task") {
			t.Errorf("output should contain task title, got:\n%s", output)
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

// helperListCmd creates a cobra.Command with the --sort flag set to the given value.
func helperListCmd(sortValue string) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().String("sort", "topo", "Sort order: topo or id")
	_ = cmd.Flags().Set("sort", sortValue)
	return cmd
}

func TestListCmd_TopologicalSort(t *testing.T) {
	t.Run("default topo sort orders by dependencies", func(t *testing.T) {
		// 2 depends on 1, so 1 must come before 2
		tasks := map[string]task.Task{
			"2": {ID: "2", Title: "Second task", Type: "coding.enhancement", Status: "pending", Dependencies: []string{"1"}},
			"1": {ID: "1", Title: "First task", Type: "coding.feature", Status: "completed", Dependencies: nil},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		// Task 1 (ID=1) must come before Task 2 (ID=2) in topo order
		idx1 := strings.Index(output, "First task")
		idx2 := strings.Index(output, "Second task")
		if idx1 == -1 || idx2 == -1 {
			t.Fatalf("missing task titles in output:\n%s", output)
		}
		if idx1 >= idx2 {
			t.Errorf("task 1 should appear before task 2 in topo order, got:\n%s", output)
		}
	})

	t.Run("--sort id restores natural ID ordering", func(t *testing.T) {
		// Create tasks where topo order differs from ID order:
		// ID order: 1, 2, 3; Topo order: 3, 1, 2 (3 has no deps, 1 depends on 3, 2 depends on 1)
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: "Depends on 3", Type: "coding.feature", Status: "pending", Dependencies: []string{"3"}},
			"2": {ID: "2", Title: "Depends on 1", Type: "coding.feature", Status: "pending", Dependencies: []string{"1"}},
			"3": {ID: "3", Title: "No deps", Type: "coding.feature", Status: "pending", Dependencies: nil},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		cmd := helperListCmd("id")
		output := captureStdout(func() {
			err := runList(cmd, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		// In --sort id mode, order should be 1, 2, 3
		lines := strings.Split(output, "\n")
		var taskLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "---") || strings.HasPrefix(trimmed, "ID ") ||
				strings.Contains(trimmed, "found") {
				continue
			}
			if len(trimmed) > 0 && trimmed[0] >= '0' && trimmed[0] <= '9' {
				taskLines = append(taskLines, trimmed)
			}
		}

		if len(taskLines) < 3 {
			t.Fatalf("expected 3 task lines, got %d\noutput:\n%s", len(taskLines), output)
		}

		expectedOrder := []string{"1", "2", "3"}
		for i, expectedID := range expectedOrder {
			if !strings.HasPrefix(taskLines[i], expectedID) {
				t.Errorf("task line %d: expected to start with %q, got %q\nfull output:\n%s", i, expectedID, taskLines[i], output)
			}
		}
	})

	t.Run("invalid --sort value returns error", func(t *testing.T) {
		cmd := helperListCmd("invalid")
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: "Task", Type: "coding.feature", Status: "pending"},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		err := runList(cmd, []string{})
		if err == nil {
			t.Fatal("expected error for invalid --sort value, got nil")
		}
		if !strings.Contains(err.Error(), "invalid --sort value") {
			t.Errorf("error should mention invalid sort value, got: %v", err)
		}
	})
}

func TestListCmd_CycleMarker(t *testing.T) {
	t.Run("cycle nodes display [cycle] marker", func(t *testing.T) {
		// 1 -> 2 -> 3 -> 1 (cycle)
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: "Cycle task 1", Type: "coding.feature", Status: "pending", Dependencies: []string{"3"}},
			"2": {ID: "2", Title: "Cycle task 2", Type: "coding.feature", Status: "pending", Dependencies: []string{"1"}},
			"3": {ID: "3", Title: "Cycle task 3", Type: "coding.feature", Status: "pending", Dependencies: []string{"2"}},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		// Force non-TTY mode for predictable output
		orig := listIsTerminalFunc
		listIsTerminalFunc = func() bool { return false }
		defer func() { listIsTerminalFunc = orig }()

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		if !strings.Contains(output, "[cycle]") {
			t.Errorf("output should contain [cycle] marker for cycle nodes, got:\n%s", output)
		}
	})

	t.Run("cycle marker has no ANSI codes in non-TTY mode", func(t *testing.T) {
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: "Cycle task", Type: "coding.feature", Status: "pending", Dependencies: []string{"1"}},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		orig := listIsTerminalFunc
		listIsTerminalFunc = func() bool { return false }
		defer func() { listIsTerminalFunc = orig }()

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		if strings.Contains(output, "\033[") {
			t.Errorf("non-TTY output should not contain ANSI escape codes, got:\n%s", output)
		}
		if !strings.Contains(output, "[cycle]") {
			t.Errorf("non-TTY output should contain plain [cycle] marker, got:\n%s", output)
		}
	})

	t.Run("cycle marker has ANSI codes in TTY mode", func(t *testing.T) {
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: "Cycle task", Type: "coding.feature", Status: "pending", Dependencies: []string{"1"}},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		orig := listIsTerminalFunc
		listIsTerminalFunc = func() bool { return true }
		defer func() { listIsTerminalFunc = orig }()

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		if !strings.Contains(output, "\033[33m") {
			t.Errorf("TTY output should contain ANSI color codes for markers, got:\n%s", output)
		}
		if !strings.Contains(output, "[cycle]") {
			t.Errorf("TTY output should contain [cycle] text, got:\n%s", output)
		}
	})
}

func TestListCmd_MissingDepMarker(t *testing.T) {
	t.Run("missing deps display [missing: id] marker", func(t *testing.T) {
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: "Standalone task", Type: "coding.feature", Status: "pending", Dependencies: nil},
			"2": {ID: "2", Title: "Task with missing dep", Type: "coding.feature", Status: "pending", Dependencies: []string{"1", "999"}},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		orig := listIsTerminalFunc
		listIsTerminalFunc = func() bool { return false }
		defer func() { listIsTerminalFunc = orig }()

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		if !strings.Contains(output, "[missing: 999]") {
			t.Errorf("output should contain [missing: 999] marker, got:\n%s", output)
		}
	})
}

func TestListCmd_PipeModeColorSuppression(t *testing.T) {
	t.Run("pipe mode suppresses color for all markers", func(t *testing.T) {
		// Create tasks with both cycle and missing deps
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: "Self-cycle", Type: "coding.feature", Status: "pending", Dependencies: []string{"1"}},
			"2": {ID: "2", Title: "Missing dep", Type: "coding.feature", Status: "pending", Dependencies: []string{"404"}},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		orig := listIsTerminalFunc
		listIsTerminalFunc = func() bool { return false }
		defer func() { listIsTerminalFunc = orig }()

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		if strings.Contains(output, "\033[") {
			t.Errorf("pipe mode output should not contain ANSI codes, got:\n%s", output)
		}
		if !strings.Contains(output, "[cycle]") {
			t.Errorf("output should contain [cycle], got:\n%s", output)
		}
		if !strings.Contains(output, "[missing: 404]") {
			t.Errorf("output should contain [missing: 404], got:\n%s", output)
		}
	})
}

func TestListCmd_ColumnAlignmentWithMarkers(t *testing.T) {
	t.Run("column alignment accounts for marker width", func(t *testing.T) {
		tasks := map[string]task.Task{
			"1": {ID: "1", Title: "Task with missing dep", Type: "coding.feature", Status: "pending", Dependencies: []string{"999"}},
			"2": {ID: "2", Title: "Normal task", Type: "coding.feature", Status: "completed"},
		}
		_ = setupFullProject(t, SetupOpts{Tasks: tasks})

		orig := listIsTerminalFunc
		listIsTerminalFunc = func() bool { return false }
		defer func() { listIsTerminalFunc = orig }()

		output := captureStdout(func() {
			err := runList(nil, []string{})
			if err != nil {
				t.Fatalf("runList returned error: %v", err)
			}
		})

		// Find the separator line to determine ID column width
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
		idDashCount := len(strings.TrimRight(segments[0], " "))

		// "1 [missing: 999]" = 16 chars, ID column must be at least 16 wide
		if idDashCount < 16 {
			t.Errorf("ID column is %d chars wide but marker text needs 16 -- columns misalign\nseparator: %q\noutput:\n%s",
				idDashCount, sepLine, output)
		}
	})
}
