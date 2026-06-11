package task

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"
)

// SetupOpts configures the test project created by setupFullProject.
type SetupOpts struct {
	// Tasks is the task map to write into index.json (required).
	Tasks map[string]task.Task
	// State, if non-nil, creates a task-state.json in the process directory.
	State *task.TaskState
	// UseEnvVar, when true, sets CLAUDE_PROJECT_DIR instead of using go.mod+chdir+EnsureFeatureDir.
	UseEnvVar bool
	// FeatureName defaults to "test" if empty.
	FeatureName string
}

// setupFullProject creates a fully configured test project.
func setupFullProject(t *testing.T, opts SetupOpts) (dir string) {
	t.Helper()
	dir = t.TempDir()

	featureName := opts.FeatureName
	if featureName == "" {
		featureName = "test"
	}

	// Always set env var for isolation
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	if !opts.UseEnvVar {
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	if err := feature.EnsureFeatureDir(dir, featureName); err != nil {
		t.Fatal(err)
	}

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile(featureName))
	index := &task.TaskIndex{
		Feature:      featureName,
		PRD:          "prd/prd-spec.md",
		Design:       "design/tech-design.md",
		StatusEnum:   []string{"pending", "in_progress", "completed", "blocked", "skipped"},
		PriorityEnum: []string{"P0", "P1", "P2"},
	}
	if len(opts.Tasks) > 0 {
		index.SetTasks(opts.Tasks)
	}
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Create task markdown files
	tasksDir := filepath.Join(dir, feature.GetFeatureTasksDir(featureName))
	for _, t2 := range opts.Tasks {
		if t2.File != "" {
			content := buildTestTaskMD(t2)
			if err := os.WriteFile(filepath.Join(tasksDir, t2.File), []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	// Create records dir
	if err := os.MkdirAll(filepath.Join(tasksDir, "records"), 0755); err != nil {
		t.Fatal(err)
	}

	// Optionally create task state file
	if opts.State != nil {
		statePath := feature.GetTaskStatePath(dir, featureName)
		if err := task.SaveState(statePath, opts.State); err != nil {
			t.Fatalf("SaveState failed: %v", err)
		}
	}

	if !opts.UseEnvVar {
		origWd, _ := os.Getwd()
		t.Cleanup(func() { _ = os.Chdir(origWd) })
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		if err := feature.EnsureFeatureDir(dir, featureName); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// buildTestTaskMD generates markdown with YAML frontmatter for a test task.
func buildTestTaskMD(t task.Task) string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "id: %q\n", t.ID)
	fmt.Fprintf(&b, "title: %q\n", t.Title)
	if t.Priority != "" {
		fmt.Fprintf(&b, "priority: %q\n", t.Priority)
	}
	if t.Status != "" {
		fmt.Fprintf(&b, "status: %q\n", t.Status)
	}
	if t.Type != "" {
		fmt.Fprintf(&b, "type: %q\n", t.Type)
	}
	if len(t.Dependencies) > 0 {
		b.WriteString("dependencies:\n")
		for _, d := range t.Dependencies {
			fmt.Fprintf(&b, "  - %q\n", d)
		}
	}
	if t.Breaking {
		b.WriteString("breaking: true\n")
	}
	if t.EstimatedTime != "" {
		fmt.Fprintf(&b, "estimated_time: %q\n", t.EstimatedTime)
	}
	b.WriteString("---\n\n")
	fmt.Fprintf(&b, "# %s\n", t.Title)
	return b.String()
}

// captureOutput captures stdout and stderr during a function execution.
func captureOutput(f func() error) (string, error) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, err := os.Pipe()
	if err != nil {
		return "", err
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		return "", err
	}

	os.Stdout = wOut
	os.Stderr = wErr

	outCh := make(chan string)
	errCh := make(chan string)

	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, rOut)
		outCh <- buf.String()
	}()

	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, rErr)
		errCh <- buf.String()
	}()

	runErr := f()

	_ = wOut.Close()
	_ = wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	stdout := <-outCh
	stderr := <-errCh

	return stdout + stderr, runErr
}

// captureStdout captures stdout during a function execution.
func captureStdout(f func()) string {
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
