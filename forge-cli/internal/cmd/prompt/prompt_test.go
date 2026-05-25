package prompt

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"
)

// captureStdout captures stdout during a function execution.
// Reads from pipe in a goroutine to prevent deadlock when output
// exceeds the OS pipe buffer size.
func captureStdout(f func()) string {
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	done := make(chan struct{})
	go func() {
		_, _ = buf.ReadFrom(r)
		close(done)
	}()

	f()
	_ = w.Close()
	os.Stdout = old
	<-done
	return buf.String()
}

// setupFullProject creates a fully configured test project.
func setupFullProject(t *testing.T, tasks map[string]task.Task) {
	t.Helper()
	dir := t.TempDir()

	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	if err := os.WriteFile(dir+"/go.mod", []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
		t.Fatal(err)
	}

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:      "test",
		PRD:          "prd/prd-spec.md",
		Design:       "design/tech-design.md",
		StatusEnum:   []string{"pending", "in_progress", "completed", "blocked", "skipped"},
		PriorityEnum: []string{"P0", "P1", "P2"},
	}
	if len(tasks) > 0 {
		index.SetTasks(tasks)
	}
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Create task markdown files
	tasksDir := filepath.Join(dir, feature.GetFeatureTasksDir("test"))
	for _, t2 := range tasks {
		if t2.File != "" {
			content := buildTestTaskMD(t2)
			if err := os.WriteFile(tasksDir+"/"+t2.File, []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	// Create records dir
	if err := os.MkdirAll(tasksDir+"/records", 0755); err != nil {
		t.Fatal(err)
	}

	// Set working dir
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	// Set feature
	if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
		t.Fatal(err)
	}
}

// buildTestTaskMD generates markdown with YAML frontmatter for a test task.
func buildTestTaskMD(t task.Task) string {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("id: " + `"` + t.ID + `"` + "\n")
	b.WriteString("title: " + `"` + t.Title + `"` + "\n")
	if t.Priority != "" {
		b.WriteString("priority: " + `"` + t.Priority + `"` + "\n")
	}
	if t.Status != "" {
		b.WriteString("status: " + `"` + t.Status + `"` + "\n")
	}
	if t.Type != "" {
		b.WriteString("type: " + `"` + t.Type + `"` + "\n")
	}
	b.WriteString("---\n\n")
	b.WriteString("# " + t.Title + "\n")
	return b.String()
}

func TestRunPrompt_Success(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {
			ID:     "1.1",
			Title:  "Test task",
			Status: "pending",
			File:   "1.1.md",
			Record: "records/1.1.md",
			Type:   task.TypeCodingFeature,
		},
	})

	out := captureStdout(func() {
		promptFixRecordMissed = false
		_ = runPrompt(nil, []string{"1.1"})
	})

	if out == "" {
		t.Error("expected non-empty prompt output")
	}
	if strings.Contains(out, "{{") {
		t.Errorf("output contains unreplaced placeholder: %s", out)
	}
}

func TestRunPrompt_FixRecordMissed(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {
			ID:     "1.1",
			Title:  "Test task",
			Status: "pending",
			File:   "1.1.md",
			Record: "records/1.1.md",
			Type:   task.TypeCodingFeature,
		},
	})

	out := captureStdout(func() {
		promptFixRecordMissed = true
		_ = runPrompt(nil, []string{"1.1"})
	})
	promptFixRecordMissed = false // reset

	if out == "" {
		t.Error("expected non-empty prompt output for fix-record-missed")
	}
	if strings.Contains(out, "{{") {
		t.Errorf("output contains unreplaced placeholder: %s", out)
	}
}

func TestRunPrompt_TypeMissing_ExitsWithError(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "No type", Status: "pending", File: "1.1.md", Record: "records/1.1.md", Type: ""},
	})

	if os.Getenv("TEST_PROMPT_TYPE_MISSING") == "1" {
		promptFixRecordMissed = false
		_ = runPrompt(nil, []string{"1.1"})
		return
	}

	dir, _ := os.Getwd()
	cmd := exec.Command(os.Args[0], "-test.run=TestRunPrompt_TypeMissing_ExitsWithError")
	cmd.Env = append(os.Environ(), "TEST_PROMPT_TYPE_MISSING=1")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for missing type")
	}
	if !strings.Contains(string(output), "type field missing") {
		t.Errorf("expected 'type field missing' in stderr, got: %s", string(output))
	}
}

func TestRunPrompt_UnknownType_ExitsWithError(t *testing.T) {
	setupFullProject(t, map[string]task.Task{
		"t1": {ID: "1.1", Title: "Unknown type", Status: "pending", File: "1.1.md", Record: "records/1.1.md", Type: "unknown-xyz"},
	})

	if os.Getenv("TEST_PROMPT_UNKNOWN_TYPE") == "1" {
		promptFixRecordMissed = false
		_ = runPrompt(nil, []string{"1.1"})
		return
	}

	dir, _ := os.Getwd()
	cmd := exec.Command(os.Args[0], "-test.run=TestRunPrompt_UnknownType_ExitsWithError")
	cmd.Env = append(os.Environ(), "TEST_PROMPT_UNKNOWN_TYPE=1")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for unknown type")
	}
	if !strings.Contains(string(output), "unknown type") {
		t.Errorf("expected 'unknown type' in stderr, got: %s", string(output))
	}
}

func TestRunPrompt_NoProject_ExitsWithError(t *testing.T) {
	if os.Getenv("TEST_PROMPT_NO_PROJECT") == "1" {
		promptFixRecordMissed = false
		_ = runPrompt(nil, []string{"1.1"})
		return
	}

	tmpDir := t.TempDir()
	cmd := exec.Command(os.Args[0], "-test.run=TestRunPrompt_NoProject_ExitsWithError")
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "CLAUDE_PROJECT_DIR=") || strings.HasPrefix(e, "PROJECT_ROOT=") {
			continue
		}
		env = append(env, e)
	}
	env = append(env, "TEST_PROMPT_NO_PROJECT=1", "CLAUDE_PROJECT_DIR=")
	cmd.Env = env
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for no project")
	}
	out := string(output)
	if !strings.Contains(out, "NO_PROJECT") && !strings.Contains(out, "NO_FEATURE") {
		t.Errorf("expected NO_PROJECT or NO_FEATURE error, got: %s", out)
	}
}

func TestRunPrompt_NoFeature_ExitsWithError(t *testing.T) {
	if os.Getenv("TEST_PROMPT_NO_FEATURE") == "1" {
		promptFixRecordMissed = false
		_ = runPrompt(nil, []string{"1.1"})
		return
	}

	dir := t.TempDir()
	if err := os.WriteFile(dir+"/go.mod", []byte("module test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir+"/docs/features", 0755); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRunPrompt_NoFeature_ExitsWithError")
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "CLAUDE_PROJECT_DIR=") || strings.HasPrefix(e, "PROJECT_ROOT=") {
			continue
		}
		env = append(env, e)
	}
	env = append(env, "TEST_PROMPT_NO_FEATURE=1", "CLAUDE_PROJECT_DIR="+dir)
	cmd.Env = env
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit for no feature")
	}
	if !strings.Contains(string(output), "NO_FEATURE") {
		t.Errorf("expected NO_FEATURE error, got: %s", string(output))
	}
}

func TestExit_NonAIError(t *testing.T) {
	if os.Getenv("TEST_EXIT_PLAIN_ERR") == "1" {
		base.Exit(fmt.Errorf("plain error"))
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExit_NonAIError")
	cmd.Env = append(os.Environ(), "TEST_EXIT_PLAIN_ERR=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit")
	}
	if !strings.Contains(string(output), "ERROR: plain error") {
		t.Errorf("expected plain error message, got: %s", string(output))
	}
}
