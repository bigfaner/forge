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

// TestConfigGetCommand_EvalSettings tests CLI "forge config get eval.*" for eval settings.
func TestConfigGetCommand_EvalSettings(t *testing.T) {
	setupEvalConfig := func(t *testing.T, content string) string {
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

	t.Run("eval.proposal.target returns value", func(t *testing.T) {
		dir := setupEvalConfig(t, `eval:
  proposal:
    target: 850
    iterations: 3
`)
		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "get", "eval.proposal.target", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "850" {
			t.Errorf("expected '850', got %q", output)
		}
	})

	t.Run("eval.proposal.iterations returns value", func(t *testing.T) {
		dir := setupEvalConfig(t, `eval:
  proposal:
    target: 850
    iterations: 5
`)
		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "get", "eval.proposal.iterations", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "5" {
			t.Errorf("expected '5', got %q", output)
		}
	})

	t.Run("eval not configured exits with error", func(t *testing.T) {
		dir := setupEvalConfig(t, "auto:\n  gitPush: true\n")

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		configGetCmd.SilenceUsage = true
		defer func() { configGetCmd.SilenceUsage = false }()
		rootCmd.SetArgs([]string{"config", "get", "eval.proposal.target", "--project-root", dir})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for unconfigured eval")
		}
	})

	t.Run("eval returns summary of all types", func(t *testing.T) {
		dir := setupEvalConfig(t, `eval:
  proposal:
    target: 900
    iterations: 3
  ui:
    target: 950
    iterations: 3
`)
		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "get", "eval", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "proposal:") {
			t.Errorf("expected 'proposal:' in eval summary, got %q", output)
		}
		if !strings.Contains(output, "ui:") {
			t.Errorf("expected 'ui:' in eval summary, got %q", output)
		}
	})
}

// TestConfigSetCommand_EvalSettings tests CLI "forge config set eval.*" for eval settings.
func TestConfigSetCommand_EvalSettings(t *testing.T) {
	t.Run("set eval.proposal.target and verify with get", func(t *testing.T) {
		dir := t.TempDir()

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "set", "eval.proposal.target", "850", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		val, err := forgeconfig.GetConfigValue(dir, "eval.proposal.target")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "850" {
			t.Errorf("expected '850', got %q", val)
		}
	})

	t.Run("set eval.journey.iterations and verify with get", func(t *testing.T) {
		dir := t.TempDir()

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"config", "set", "eval.journey.iterations", "5", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		val, err := forgeconfig.GetConfigValue(dir, "eval.journey.iterations")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "5" {
			t.Errorf("expected '5', got %q", val)
		}
	})

	t.Run("set eval non-leaf returns error", func(t *testing.T) {
		dir := t.TempDir()

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"config", "set", "eval", "900", "--project-root", dir})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for non-leaf eval set")
		}
	})
}

// TestConfigInitCommand_EvalBlock tests that forge config init generates an eval block.
func TestConfigInitCommand_EvalBlock(t *testing.T) {
	t.Run("generated config contains eval block with all 7 types", func(t *testing.T) {
		dir := t.TempDir()

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

		configFile := filepath.Join(dir, feature.ForgeDir, feature.ForgeConfigFileName)
		data, err := os.ReadFile(configFile)
		if err != nil {
			t.Fatalf("config file not created: %v", err)
		}

		content := string(data)

		// Verify eval block exists
		if !strings.Contains(content, "eval:") {
			t.Errorf("expected 'eval:' in generated config, got:\n%s", content)
		}

		// Verify all 7 types are present
		for _, evalType := range []string{"proposal:", "prd:", "design:", "ui:", "journey:", "contract:", "consistency:"} {
			if !strings.Contains(content, evalType) {
				t.Errorf("expected %q in generated config, got:\n%s", evalType, content)
			}
		}
	})

	t.Run("generated eval block contains rubric-default values", func(t *testing.T) {
		dir := t.TempDir()

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

		// Read back via GetConfigValue to verify values match rubric defaults
		expectedDefaults := map[string]string{
			"eval.proposal.target":        "900",
			"eval.proposal.iterations":    "3",
			"eval.prd.target":             "900",
			"eval.prd.iterations":         "3",
			"eval.design.target":          "900",
			"eval.design.iterations":      "3",
			"eval.ui.target":              "950",
			"eval.ui.iterations":          "3",
			"eval.journey.target":         "850",
			"eval.journey.iterations":     "3",
			"eval.contract.target":        "850",
			"eval.contract.iterations":    "3",
			"eval.consistency.target":     "900",
			"eval.consistency.iterations": "3",
		}

		for key, expected := range expectedDefaults {
			val, err := forgeconfig.GetConfigValue(dir, key)
			if err != nil {
				t.Errorf("GetConfigValue(%q) returned error: %v", key, err)
				continue
			}
			if val != expected {
				t.Errorf("key %q: expected %q, got %q", key, expected, val)
			}
		}
	})
}
