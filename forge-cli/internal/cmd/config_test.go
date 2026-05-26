package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
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
		dir := setupConfig(t, "auto:\n  test:\n    quick: true\n")

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
	t.Run("delegates to configInitFunc and creates config", func(t *testing.T) {
		dir := t.TempDir()

		// Override configInitFunc to simulate interactive config
		origConfigInit := configInitFunc
		configInitFunc = testConfigInit
		t.Cleanup(func() { configInitFunc = origConfigInit })

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
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

		// Verify output acknowledges config creation
		if !strings.Contains(stdout.String(), "Config written") {
			t.Errorf("expected 'Config written' in output, got %q", stdout.String())
		}
	})

	t.Run("prints skipped message when config init is skipped", func(t *testing.T) {
		dir := t.TempDir()

		origConfigInit := configInitFunc
		configInitFunc = func(_ string) initAction {
			return initAction{status: "SKIPPED", target: ".forge/config.yaml", detail: "kept existing"}
		}
		t.Cleanup(func() { configInitFunc = origConfigInit })

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "init", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(stdout.String(), "Config init skipped") {
			t.Errorf("expected 'Config init skipped' in output, got %q", stdout.String())
		}
	})

	t.Run("returns error on FAILED status", func(t *testing.T) {
		dir := t.TempDir()

		origConfigInit := configInitFunc
		configInitFunc = func(_ string) initAction {
			return initAction{status: "FAILED", target: ".forge/config.yaml", detail: "test failure"}
		}
		t.Cleanup(func() { configInitFunc = origConfigInit })

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"config", "init", "--project-root", dir})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for FAILED status")
		}
	})

	t.Run("prints cancelled message on CANCELLED status", func(t *testing.T) {
		dir := t.TempDir()

		origConfigInit := configInitFunc
		configInitFunc = func(_ string) initAction {
			return initAction{status: "CANCELLED", target: ".forge/config.yaml", detail: "Ctrl+C"}
		}
		t.Cleanup(func() { configInitFunc = origConfigInit })

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "init", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(stdout.String(), "Config init cancelled") {
			t.Errorf("expected 'Config init cancelled' in output, got %q", stdout.String())
		}
	})
}

func TestConfigSetCommand(t *testing.T) {
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

	t.Run("set auto.gitPush and verify with get", func(t *testing.T) {
		dir := t.TempDir()

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "set", "auto.gitPush", "true", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify via get
		val, err := forgeconfig.GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("set worktree.source-branch and verify", func(t *testing.T) {
		dir := t.TempDir()

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "set", "worktree.source-branch", "develop", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		val, err := forgeconfig.GetConfigValue(dir, "worktree.source-branch")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "develop" {
			t.Errorf("expected 'develop', got %q", val)
		}
	})

	t.Run("set test-framework and verify with get command", func(t *testing.T) {
		dir := t.TempDir()

		var setStdout bytes.Buffer
		rootCmd.SetOut(&setStdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "set", "test-framework", "pytest", "--project-root", dir})
		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify via get command
		var getStdout bytes.Buffer
		rootCmd.SetOut(&getStdout)
		rootCmd.SetArgs([]string{"config", "get", "test-framework", "--project-root", dir})
		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(getStdout.String())
		if output != "pytest" {
			t.Errorf("expected 'pytest', got %q", output)
		}
	})

	t.Run("unknown key returns error", func(t *testing.T) {
		dir := t.TempDir()

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"config", "set", "nonexistent", "value", "--project-root", dir})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for unknown key")
		}
	})

	t.Run("invalid args count returns error", func(t *testing.T) {
		dir := t.TempDir()

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"config", "set", "auto.gitPush", "--project-root", dir})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for missing value arg")
		}
	})

	t.Run("set auto.cleanCode rejected as ModeToggle", func(t *testing.T) {
		dir := t.TempDir()

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "set", "auto.cleanCode", "true", "--project-root", dir})
		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for ModeToggle direct set")
		}
		if !strings.Contains(err.Error(), "cannot set ModeToggle directly") {
			t.Errorf("expected cannot set ModeToggle directly in error, got %v", err)
		}
	})

	t.Run("set overwrites existing value", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: false\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "set", "auto.gitPush", "true", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		val, err := forgeconfig.GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})
}
