package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/feature"
)

func TestConfigGetCommand(t *testing.T) {
	setupConfig := func(t *testing.T, content string) string {
		t.Helper()
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(forgeDir, feature.ForgeConfigFileName), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		return dir
	}

	t.Run("auto.gitPush returns true", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "get", "auto.gitPush", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "true" {
			t.Errorf("expected 'true', got %q", output)
		}
	})

	t.Run("auto.gitPush returns false when absent", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  e2eTest:\n    quick: true\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "get", "auto.gitPush", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "false" {
			t.Errorf("expected 'false', got %q", output)
		}
	})

	t.Run("test-framework returns value", func(t *testing.T) {
		dir := setupConfig(t, "test-framework: pytest\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "get", "test-framework", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "pytest" {
			t.Errorf("expected 'pytest', got %q", output)
		}
	})

	t.Run("test-command returns value", func(t *testing.T) {
		dir := setupConfig(t, "test-command: go test ./...\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "get", "test-command", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "go test ./..." {
			t.Errorf("expected 'go test ./...', got %q", output)
		}
	})

	t.Run("worktree.source-branch returns value", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  source-branch: develop\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "get", "worktree.source-branch", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "develop" {
			t.Errorf("expected 'develop', got %q", output)
		}
	})

	t.Run("worktree.copy-files returns one per line", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  copy-files:\n    - .env\n    - .env.local\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "get", "worktree.copy-files", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != ".env\n.env.local" {
			t.Errorf("expected '.env\\n.env.local', got %q", output)
		}
	})

	t.Run("nonexistent key exits with error", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"config", "get", "nonexistent", "--project-root", dir})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for nonexistent key")
		}
	})

	t.Run("missing config file exits with error", func(t *testing.T) {
		dir := t.TempDir()

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"config", "get", "test-framework", "--project-root", dir})

		// Silence usage on error so only the error message is captured
		configGetCmd.SilenceUsage = true
		defer func() { configGetCmd.SilenceUsage = false }()

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for missing config")
		}

		// stdout should not contain a config value (usage is suppressed)
		output := strings.TrimSpace(stdout.String())
		if output != "" && !strings.Contains(output, "Usage") {
			t.Errorf("expected no config value output, got %q", output)
		}
	})

	t.Run("output is plain text no formatting blocks", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  source-branch: main\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "get", "worktree.source-branch", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := stdout.String()
		if strings.Contains(output, "```") {
			t.Errorf("output should not contain formatting blocks, got %q", output)
		}
		if strings.HasPrefix(output, "> ") {
			t.Errorf("output should not have block markers, got %q", output)
		}
	})
}

func TestConfigInitCommand(t *testing.T) {
	t.Run("writes config with auto and worktree", func(t *testing.T) {
		dir := t.TempDir()

		// Simulate user input: y (e2e quick), y (e2e full), y (consolidate quick), y (consolidate full), n (clean quick), n (clean full), y (git push), main (source branch), .env (copy files)
		input := "y\ny\ny\ny\nn\nn\ny\nmain\n.env\n"
		var stdin bytes.Buffer
		stdin.WriteString(input)
		var stdout bytes.Buffer

		rootCmd.SetIn(&stdin)
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"config", "init", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify config file was created
		configFile := filepath.Join(dir, feature.ForgeDir, feature.ForgeConfigFileName)
		data, err := os.ReadFile(configFile)
		if err != nil {
			t.Fatalf("config file not created: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "auto:") {
			t.Errorf("expected auto in config, got %q", content)
		}
		if !strings.Contains(content, "worktree:") {
			t.Errorf("expected worktree in config, got %q", content)
		}
	})

	t.Run("prompts when config already exists", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		existingConfig := "auto:\n  gitPush: true\n"
		if err := os.WriteFile(filepath.Join(forgeDir, feature.ForgeConfigFileName), []byte(existingConfig), 0o644); err != nil {
			t.Fatal(err)
		}

		// Answer 'n' to reconfigure prompt
		input := "n\n"
		var stdin bytes.Buffer
		stdin.WriteString(input)
		var stdout bytes.Buffer

		rootCmd.SetIn(&stdin)
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"config", "init", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Config should be unchanged
		data, _ := os.ReadFile(filepath.Join(forgeDir, feature.ForgeConfigFileName))
		if !strings.Contains(string(data), "gitPush: true") {
			t.Errorf("config should not have been modified, got %q", string(data))
		}

		if !strings.Contains(stdout.String(), "already exists") {
			t.Errorf("expected 'already exists' in output, got %q", stdout.String())
		}
	})

	t.Run("reconfigures when user answers yes", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		existingConfig := "auto:\n  gitPush: false\n"
		if err := os.WriteFile(filepath.Join(forgeDir, feature.ForgeConfigFileName), []byte(existingConfig), 0o644); err != nil {
			t.Fatal(err)
		}

		// Answer 'y' to reconfigure, then select defaults for all auto settings, no worktree
		input := "y\ny\ny\ny\ny\nn\nn\nn\n\n\n"
		var stdin bytes.Buffer
		stdin.WriteString(input)
		var stdout bytes.Buffer

		rootCmd.SetIn(&stdin)
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"config", "init", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Config should be updated
		data, _ := os.ReadFile(filepath.Join(forgeDir, feature.ForgeConfigFileName))
		content := string(data)
		if !strings.Contains(content, "auto:") {
			t.Errorf("config should contain auto, got %q", content)
		}
	})
}

func TestParseMultiSelect(t *testing.T) {
	options := []string{"alpha", "beta", "gamma"}

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"single selection", "1", []string{"alpha"}},
		{"multiple selections", "1 3", []string{"alpha", "gamma"}},
		{"empty input", "", nil},
		{"whitespace only", "   ", nil},
		{"out of range high", "5", nil},
		{"out of range low", "0", nil},
		{"mixed valid and invalid", "1 99 2", []string{"alpha", "beta"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseMultiSelect(tt.input, options)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}
