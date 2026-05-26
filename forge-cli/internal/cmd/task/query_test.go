package task

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"forge-cli/pkg/task"
)

// captureQueryOutput runs f while capturing os.Stdout output.
func captureQueryOutput(t *testing.T, f func()) string {
	t.Helper()
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

// resetQueryFlags resets the queryVerbose flag to avoid state leakage between tests.
func resetQueryFlags() {
	queryVerbose = false
	_ = queryCmd.Flags().Set("verbose", "false")
}

func TestQuery_DefaultOutput(t *testing.T) {
	resetQueryFlags()
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"phase1-1-task": {ID: "1", Title: "Test task", Status: "pending", Priority: "P1", File: "1.md", Record: "records/1.md", Type: "coding.feature"},
	}})

	Cmd.SetArgs([]string{"query", "1"})
	output := captureQueryOutput(t, func() {
		_ = Cmd.Execute()
	})

	if !strings.Contains(output, "TASK_ID: 1") {
		t.Errorf("expected TASK_ID in output, got:\n%s", output)
	}
	if !strings.Contains(output, "STATUS: pending") {
		t.Errorf("expected STATUS in output, got:\n%s", output)
	}
	// Verbose-only fields must NOT appear
	if strings.Contains(output, "TITLE:") {
		t.Errorf("TITLE should not appear in default mode, got:\n%s", output)
	}
	if strings.Contains(output, "PRIORITY:") {
		t.Errorf("PRIORITY should not appear in default mode, got:\n%s", output)
	}
	if strings.Contains(output, "KEY:") {
		t.Errorf("KEY should not appear in default mode, got:\n%s", output)
	}
	if strings.Contains(output, "RELATED_FIXES:") {
		t.Errorf("RELATED_FIXES should not appear in default mode, got:\n%s", output)
	}
}

func TestQuery_VerboseOutput(t *testing.T) {
	resetQueryFlags()
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"phase1-1-task": {ID: "1", Title: "Test task", Status: "pending", Priority: "P1", Type: "coding.feature", File: "1.md", Record: "records/1.md", Dependencies: []string{"2", "3"}},
	}})

	Cmd.SetArgs([]string{"query", "1", "--verbose"})
	output := captureQueryOutput(t, func() {
		_ = Cmd.Execute()
	})

	expected := []string{
		"KEY: phase1-1-task",
		"TASK_ID: 1",
		"TITLE: Test task",
		"STATUS: pending",
		"PRIORITY: P1",
		"TYPE: coding.feature",
		"DEPENDENCIES:",
		"  2",
		"  3",
		"TASK_FILE:",
		"RECORD_FILE:",
	}
	for _, field := range expected {
		if !strings.Contains(output, field) {
			t.Errorf("expected %q in verbose output, got:\n%s", field, output)
		}
	}
	// RELATED_FIXES should NOT appear when no fix tasks exist
	if strings.Contains(output, "RELATED_FIXES:") {
		t.Errorf("RELATED_FIXES should not appear when no fix tasks exist, got:\n%s", output)
	}
}

func TestQuery_VerboseShorthand(t *testing.T) {
	resetQueryFlags()
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1", Title: "Short flag test", Status: "completed", Priority: "P0", Type: "coding.fix", File: "1.md", Record: "records/1.md"},
	}})

	Cmd.SetArgs([]string{"query", "1", "-v"})
	output := captureQueryOutput(t, func() {
		_ = Cmd.Execute()
	})

	if !strings.Contains(output, "KEY: t1") {
		t.Errorf("expected KEY in -v output, got:\n%s", output)
	}
}

func TestQuery_VerboseWithRelatedFixes(t *testing.T) {
	resetQueryFlags()
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"source-task": {ID: "5", Title: "Source task", Status: "completed", Priority: "P0", Type: "coding.feature", File: "5.md", Record: "records/5.md"},
		"fix-task-1":  {ID: "5.1", Title: "Fix issue A", Status: "pending", Priority: "P1", Type: "coding.fix", File: "5.1.md", Record: "records/5.1.md", SourceTaskID: "5"},
		"fix-task-2":  {ID: "5.2", Title: "Fix issue B", Status: "completed", Priority: "P2", Type: "coding.fix", File: "5.2.md", Record: "records/5.2.md", SourceTaskID: "5"},
	}})

	Cmd.SetArgs([]string{"query", "5", "-v"})
	output := captureQueryOutput(t, func() {
		_ = Cmd.Execute()
	})

	if !strings.Contains(output, "RELATED_FIXES:") {
		t.Errorf("expected RELATED_FIXES header, got:\n%s", output)
	}
	if !strings.Contains(output, "5.1 [pending] Fix issue A") {
		t.Errorf("expected fix task 5.1, got:\n%s", output)
	}
	if !strings.Contains(output, "5.2 [completed] Fix issue B") {
		t.Errorf("expected fix task 5.2, got:\n%s", output)
	}
}

func TestQuery_VerboseNoRelatedFixes(t *testing.T) {
	resetQueryFlags()
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1", Title: "No fixes", Status: "pending", Priority: "P2", Type: "coding.cleanup", File: "1.md", Record: "records/1.md"},
	}})

	Cmd.SetArgs([]string{"query", "1", "-v"})
	output := captureQueryOutput(t, func() {
		_ = Cmd.Execute()
	})

	if strings.Contains(output, "RELATED_FIXES:") {
		t.Errorf("RELATED_FIXES should not appear when no fixes exist, got:\n%s", output)
	}
}

func TestQuery_VerboseWithScope(t *testing.T) {
	resetQueryFlags()
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1", Title: "Scoped task", Status: "pending", Priority: "P1", Type: "coding.feature", File: "1.md", Record: "records/1.md", SurfaceKey: "backend"},
	}})

	Cmd.SetArgs([]string{"query", "1", "-v"})
	output := captureQueryOutput(t, func() {
		_ = Cmd.Execute()
	})

	if !strings.Contains(output, "SURFACE_KEY: backend") {
		t.Errorf("expected SURFACE_KEY in verbose output, got:\n%s", output)
	}
}

func TestQuery_VerboseEmptyScope(t *testing.T) {
	resetQueryFlags()
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1", Title: "No scope", Status: "pending", Priority: "P1", Type: "coding.feature", File: "1.md", Record: "records/1.md"},
	}})

	Cmd.SetArgs([]string{"query", "1", "-v"})
	output := captureQueryOutput(t, func() {
		_ = Cmd.Execute()
	})

	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "SURFACE_KEY:") {
			t.Errorf("SURFACE_KEY should not appear when empty, got:\n%s", output)
			break
		}
	}
}

func TestQuery_DefaultOutputWithScope(t *testing.T) {
	resetQueryFlags()
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1", Title: "Scoped", Status: "pending", File: "1.md", Record: "records/1.md", SurfaceKey: "backend"},
	}})

	Cmd.SetArgs([]string{"query", "1"})
	output := captureQueryOutput(t, func() {
		_ = Cmd.Execute()
	})

	if !strings.Contains(output, "SURFACE_KEY: backend") {
		t.Errorf("expected SURFACE_KEY in default output when set, got:\n%s", output)
	}
}

func TestQuery_DefaultBreakingTrue(t *testing.T) {
	resetQueryFlags()
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1", Title: "Breaking", Status: "pending", File: "1.md", Record: "records/1.md", Breaking: true},
	}})

	Cmd.SetArgs([]string{"query", "1"})
	output := captureQueryOutput(t, func() {
		_ = Cmd.Execute()
	})

	if !strings.Contains(output, "BREAKING: true") {
		t.Errorf("expected BREAKING in default output when true, got:\n%s", output)
	}
}

func TestQuery_VerboseWithNoDependencies(t *testing.T) {
	resetQueryFlags()
	setupFullProject(t, SetupOpts{Tasks: map[string]task.Task{
		"t1": {ID: "1", Title: "No deps", Status: "pending", Priority: "P1", Type: "coding.feature", File: "1.md", Record: "records/1.md"},
	}})

	Cmd.SetArgs([]string{"query", "1", "-v"})
	output := captureQueryOutput(t, func() {
		_ = Cmd.Execute()
	})

	if strings.Contains(output, "DEPENDENCIES:") {
		t.Errorf("DEPENDENCIES should not appear when empty, got:\n%s", output)
	}
}
