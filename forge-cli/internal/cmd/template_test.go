package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestListTemplates(t *testing.T) {
	t.Run("lists available templates", func(t *testing.T) {
		out := captureStdout(func() {
			listTemplates()
		})
		if !strings.Contains(out, "Available templates:") {
			t.Errorf("output should contain 'Available templates:', got: %s", out)
		}
		if !strings.Contains(out, "fix-task") {
			t.Errorf("output should contain 'fix-task', got: %s", out)
		}
	})
}

func TestShowTemplate(t *testing.T) {
	t.Run("shows known template", func(t *testing.T) {
		out := captureStdout(func() {
			showTemplate("fix-task")
		})
		if out == "" {
			t.Error("output should not be empty for known template")
		}
	})

	t.Run("unknown template exits", func(t *testing.T) {
		if os.Getenv("TEST_SHOW_UNKNOWN_TEMPLATE") == "1" {
			showTemplate("nonexistent-template")
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestShowTemplate/unknown_template_exits")
		cmd.Env = append(os.Environ(), "TEST_SHOW_UNKNOWN_TEMPLATE=1")
		output, err := cmd.CombinedOutput()
		if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 1 {
			t.Fatalf("expected exit 1, got: %v, output: %s", err, output)
		}
	})
}

func TestRunTemplate(t *testing.T) {
	t.Run("no args delegates to listTemplates", func(t *testing.T) {
		out := captureStdout(func() {
			runTemplate(nil, []string{})
		})
		if !strings.Contains(out, "Available templates:") {
			t.Errorf("no-args should list templates, got: %s", out)
		}
	})

	t.Run("with arg delegates to showTemplate", func(t *testing.T) {
		out := captureStdout(func() {
			runTemplate(nil, []string{"fix-task"})
		})
		if out == "" {
			t.Error("with valid arg should show template content")
		}
	})
}
