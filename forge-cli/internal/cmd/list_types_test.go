package cmd

import (
	"os"
	"strings"
	"testing"

	"forge-cli/pkg/task"
)

func TestListTypesCmd(t *testing.T) {
	t.Run("outputs all registered types with descriptions", func(t *testing.T) {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		runListTypes(nil, []string{})

		_ = w.Close()
		os.Stdout = old

		buf := make([]byte, 4096)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

		if len(lines) != len(task.TaskTypeRegistry) {
			t.Fatalf("expected %d lines, got %d\noutput:\n%s", len(task.TaskTypeRegistry), len(lines), output)
		}

		// Verify each line has format "<name>  <description>"
		for _, line := range lines {
			if !strings.Contains(line, "  ") {
				t.Errorf("line does not contain two-space separator: %q", line)
			}
			parts := strings.SplitN(line, "  ", 2)
			if len(parts) != 2 {
				t.Errorf("line does not split into name and description: %q", line)
				continue
			}
			name := parts[0]
			desc := parts[1]
			if !task.ValidTypes[name] {
				t.Errorf("unknown type name in output: %q", name)
			}
			if desc == "" {
				t.Errorf("empty description for type %q", name)
			}
			if len(desc) > 60 {
				t.Errorf("description too long (%d chars) for type %q: %q", len(desc), name, desc)
			}
		}
	})

	t.Run("no args rejected by cobra.NoArgs", func(t *testing.T) {
		err := listTypesCmd.Args(listTypesCmd, []string{"extra"})
		if err == nil {
			t.Error("expected error for extra args, got nil")
		}
	})

	t.Run("command metadata", func(t *testing.T) {
		if listTypesCmd.Use != "list-types" {
			t.Errorf("Use = %q, want %q", listTypesCmd.Use, "list-types")
		}
		if listTypesCmd.Short != "List all supported task types" {
			t.Errorf("Short = %q, want %q", listTypesCmd.Short, "List all supported task types")
		}
		if listTypesCmd.Args == nil {
			t.Error("Args is nil, expected cobra.NoArgs")
		}
	})

	t.Run("registered under task parent", func(t *testing.T) {
		found := false
		for _, sub := range taskCmd.Commands() {
			if sub.Name() == "list-types" {
				found = true
				break
			}
		}
		if !found {
			t.Error("list-types not registered as subcommand of task")
		}
	})
}
