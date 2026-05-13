package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"forge-cli/pkg/task"
)

func TestRunPrompt_Success(t *testing.T) {
	setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {
				ID:     "1.1",
				Title:  "Test task",
				Status: "pending",
				File:   "1.1.md",
				Record: "records/1.1.md",
				Type:   task.TypeImplementation,
				Scope:  "backend",
			},
		},
	})

	out := captureStdout(func() {
		promptFixRecordMissed = false
		runPrompt(nil, []string{"1.1"})
	})

	if out == "" {
		t.Error("expected non-empty prompt output")
	}
	if strings.Contains(out, "{{") {
		t.Errorf("output contains unreplaced placeholder: %s", out)
	}
}

func TestRunPrompt_FixRecordMissed(t *testing.T) {
	setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {
				ID:     "1.1",
				Title:  "Test task",
				Status: "pending",
				File:   "1.1.md",
				Record: "records/1.1.md",
				Type:   task.TypeImplementation,
				Scope:  "backend",
			},
		},
	})

	out := captureStdout(func() {
		promptFixRecordMissed = true
		runPrompt(nil, []string{"1.1"})
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
	setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "No type", Status: "pending", File: "1.1.md", Record: "records/1.1.md", Type: ""},
		},
	})

	if os.Getenv("TEST_PROMPT_TYPE_MISSING") == "1" {
		promptFixRecordMissed = false
		runPrompt(nil, []string{"1.1"})
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
	setupFullProject(t, SetupOpts{
		Tasks: map[string]task.Task{
			"t1": {ID: "1.1", Title: "Unknown type", Status: "pending", File: "1.1.md", Record: "records/1.1.md", Type: "unknown-xyz"},
		},
	})

	if os.Getenv("TEST_PROMPT_UNKNOWN_TYPE") == "1" {
		promptFixRecordMissed = false
		runPrompt(nil, []string{"1.1"})
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

func TestPromptCmd_RegisteredInRoot(t *testing.T) {
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "prompt" {
			// Verify it has the get-by-task-id subcommand
			for _, sub := range cmd.Commands() {
				if sub.Name() == "get-by-task-id" {
					return
				}
			}
			t.Error("prompt command does not have get-by-task-id subcommand")
			return
		}
	}
	t.Error("promptCmd not registered in rootCmd")
}

func TestRunPrompt_NoProject_ExitsWithError(t *testing.T) {
	if os.Getenv("TEST_PROMPT_NO_PROJECT") == "1" {
		promptFixRecordMissed = false
		runPrompt(nil, []string{"1.1"})
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
		runPrompt(nil, []string{"1.1"})
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
