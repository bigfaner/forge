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

// --- Tests for the forge surfaces detect subcommand ---

// resetDetectFlags resets the --apply flag to avoid state leaking between tests.
func resetDetectFlags(t *testing.T) {
	t.Helper()
	detectApplyFlag = false
}

// TestSurfacesDetectReadOnly runs detection in read-only mode (no --apply).
// AC: shows results with source annotations, exits without writing config.
func TestSurfacesDetectReadOnly(t *testing.T) {
	resetDetectFlags(t)

	t.Run("detect shows results and does not write config", func(t *testing.T) {
		resetDetectFlags(t)
		dir := t.TempDir()
		writeGoModWithCobra(t, dir)

		// Force non-interactive to avoid TTY dependency
		origIsInteractive := isInteractiveTerminalFunc
		isInteractiveTerminalFunc = func() bool { return false }
		defer func() { isInteractiveTerminalFunc = origIsInteractive }()

		var stdout, stderr bytes.Buffer
		rootCmd.SetIn(strings.NewReader(""))
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"surfaces", "detect", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output == "" {
			t.Fatal("expected output for detected surfaces, got empty")
		}
		// Should contain the detected type
		if !strings.Contains(output, "cli") {
			t.Errorf("expected 'cli' in output, got %q", output)
		}
		// Should contain source annotation
		if !strings.Contains(output, "detected:") {
			t.Errorf("expected 'detected:' source annotation in output, got %q", output)
		}

		// Config should NOT be written (read-only mode)
		configFile := filepath.Join(dir, feature.ForgeDir, feature.ForgeConfigFileName)
		if _, err := os.Stat(configFile); !os.IsNotExist(err) {
			t.Error("config file should NOT exist in read-only mode")
		}
	})
}

// TestSurfacesDetectEmptyDetection tests empty detection results.
// AC: prints nothing to stdout, exits with code 1.
func TestSurfacesDetectEmptyDetection(t *testing.T) {
	resetDetectFlags(t)

	t.Run("empty detection prints nothing and exits 1", func(t *testing.T) {
		resetDetectFlags(t)
		dir := t.TempDir()
		// Empty directory, no manifest files → no detection

		var stdout, stderr bytes.Buffer
		rootCmd.SetIn(strings.NewReader(""))
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"surfaces", "detect", "--project-root", dir})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error (exit 1) for empty detection")
		}

		output := strings.TrimSpace(stdout.String())
		if output != "" {
			t.Errorf("expected empty stdout for empty detection, got %q", output)
		}
	})
}

// TestSurfacesDetectStdoutFormat tests the stdout format for non-interactive mode.
// AC: one line per surface <path>=<type> (<source>), where <source> is detected:<signal> or inferred:<rule-id>
func TestSurfacesDetectStdoutFormat(t *testing.T) {
	resetDetectFlags(t)

	t.Run("scalar form output format", func(t *testing.T) {
		resetDetectFlags(t)
		dir := t.TempDir()
		writeGoModWithCobra(t, dir)

		origIsInteractive := isInteractiveTerminalFunc
		isInteractiveTerminalFunc = func() bool { return false }
		defer func() { isInteractiveTerminalFunc = origIsInteractive }()

		var stdout bytes.Buffer
		rootCmd.SetIn(strings.NewReader(""))
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "detect", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		// Format: .=cli (detected:cobra) or .=cli (inferred:cmd-dir)
		if !strings.Contains(output, "=") {
			t.Errorf("expected path=type format, got %q", output)
		}
		if !strings.Contains(output, "(detected:") {
			t.Errorf("expected (detected:<signal>) source annotation, got %q", output)
		}
	})

	t.Run("inference source format", func(t *testing.T) {
		resetDetectFlags(t)
		dir := t.TempDir()
		// Go project with no framework deps but cmd/ subdirectories
		writeGoModMinimal(t, dir)
		cmdDir := filepath.Join(dir, "cmd", "myapp")
		if err := os.MkdirAll(cmdDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(cmdDir, "main.go"), []byte("package main\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		origIsInteractive := isInteractiveTerminalFunc
		isInteractiveTerminalFunc = func() bool { return false }
		defer func() { isInteractiveTerminalFunc = origIsInteractive }()

		var stdout bytes.Buffer
		rootCmd.SetIn(strings.NewReader(""))
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "detect", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		// Format: .=cli (inferred:cmd-dir)
		if !strings.Contains(output, "(inferred:") {
			t.Errorf("expected (inferred:<rule-id>) source annotation, got %q", output)
		}
	})
}

// TestSurfacesDetectApplyWithMockedTUI tests --apply mode with mocked TUI.
// AC: shows TUI confirmation, writes to config on confirm, exit code 0.
func TestSurfacesDetectApplyWithMockedTUI(t *testing.T) {
	resetDetectFlags(t)

	t.Run("apply with confirm writes to config", func(t *testing.T) {
		resetDetectFlags(t)
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}

		// Create config file
		configFile := filepath.Join(forgeDir, feature.ForgeConfigFileName)
		cfg := &forgeconfig.Config{Auto: &forgeconfig.AutoConfig{}}
		if err := writeConfigFile(configFile, cfg); err != nil {
			t.Fatal(err)
		}

		writeGoModWithCobra(t, dir)

		// Mock askSurfaceConfirmation to simulate confirm
		origAsk := askSurfaceConfirmation
		askSurfaceConfirmation = func(_ string) (forgeconfig.SurfacesMap, forgeconfig.SourcesMap, bool) {
			return forgeconfig.SurfacesMap{".": "cli"},
				forgeconfig.SourcesMap{".": "dependency:cobra"},
				false
		}
		defer func() { askSurfaceConfirmation = origAsk }()

		// Mock isInteractiveTerminal to return true
		origIsInteractive := isInteractiveTerminalFunc
		isInteractiveTerminalFunc = func() bool { return true }
		defer func() { isInteractiveTerminalFunc = origIsInteractive }()

		var stdout bytes.Buffer
		rootCmd.SetIn(strings.NewReader(""))
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "detect", "--apply", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify config file on disk contains the detected surfaces
		updatedCfg, err := forgeconfig.ReadConfig(dir)
		if err != nil {
			t.Fatalf("read config: %v", err)
		}
		if updatedCfg.Surfaces == nil || updatedCfg.Surfaces["."] != "cli" {
			t.Errorf("expected surfaces['.']=cli in config, got %v", updatedCfg.Surfaces)
		}
	})

	t.Run("apply creates config file if missing", func(t *testing.T) {
		resetDetectFlags(t)
		dir := t.TempDir()
		writeGoModWithCobra(t, dir)

		// Mock askSurfaceConfirmation to simulate confirm
		origAsk := askSurfaceConfirmation
		askSurfaceConfirmation = func(_ string) (forgeconfig.SurfacesMap, forgeconfig.SourcesMap, bool) {
			return forgeconfig.SurfacesMap{".": "cli"},
				forgeconfig.SourcesMap{".": "dependency:cobra"},
				false
		}
		defer func() { askSurfaceConfirmation = origAsk }()

		origIsInteractive := isInteractiveTerminalFunc
		isInteractiveTerminalFunc = func() bool { return true }
		defer func() { isInteractiveTerminalFunc = origIsInteractive }()

		var stdout bytes.Buffer
		rootCmd.SetIn(strings.NewReader(""))
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "detect", "--apply", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Config file should have been created
		updatedCfg, err := forgeconfig.ReadConfig(dir)
		if err != nil {
			t.Fatalf("read config: %v", err)
		}
		if updatedCfg.Surfaces == nil || updatedCfg.Surfaces["."] != "cli" {
			t.Errorf("expected surfaces['.']=cli in config, got %v", updatedCfg.Surfaces)
		}
	})

	t.Run("apply with cancel exits without writing", func(t *testing.T) {
		resetDetectFlags(t)
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}

		configFile := filepath.Join(forgeDir, feature.ForgeConfigFileName)
		cfg := &forgeconfig.Config{Auto: &forgeconfig.AutoConfig{}}
		if err := writeConfigFile(configFile, cfg); err != nil {
			t.Fatal(err)
		}

		writeGoModWithCobra(t, dir)

		// Mock askSurfaceConfirmation to simulate cancel
		origAsk := askSurfaceConfirmation
		askSurfaceConfirmation = func(_ string) (forgeconfig.SurfacesMap, forgeconfig.SourcesMap, bool) {
			return nil, nil, true // cancelled
		}
		defer func() { askSurfaceConfirmation = origAsk }()

		origIsInteractive := isInteractiveTerminalFunc
		isInteractiveTerminalFunc = func() bool { return true }
		defer func() { isInteractiveTerminalFunc = origIsInteractive }()

		var stdout bytes.Buffer
		rootCmd.SetIn(strings.NewReader(""))
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "detect", "--apply", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Config should NOT have surfaces written
		updatedCfg, _ := forgeconfig.ReadConfig(dir)
		if updatedCfg != nil && len(updatedCfg.Surfaces) > 0 {
			t.Errorf("surfaces should NOT be written on cancel, got %v", updatedCfg.Surfaces)
		}
	})
}

// TestSurfacesDetectNonInteractive tests non-interactive terminal behavior.
// AC: prints results to stdout, no TUI, no config write, exit 0 on success, 1 if no surfaces found.
// Hard Rule: no config write regardless of --apply in non-interactive mode.
func TestSurfacesDetectNonInteractive(t *testing.T) {
	resetDetectFlags(t)

	t.Run("non-interactive prints to stdout, no config write", func(t *testing.T) {
		resetDetectFlags(t)
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}

		configFile := filepath.Join(forgeDir, feature.ForgeConfigFileName)
		cfg := &forgeconfig.Config{Auto: &forgeconfig.AutoConfig{}}
		if err := writeConfigFile(configFile, cfg); err != nil {
			t.Fatal(err)
		}

		writeGoModWithCobra(t, dir)

		// Force non-interactive
		origIsInteractive := isInteractiveTerminalFunc
		isInteractiveTerminalFunc = func() bool { return false }
		defer func() { isInteractiveTerminalFunc = origIsInteractive }()

		var stdout bytes.Buffer
		rootCmd.SetIn(strings.NewReader(""))
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "detect", "--apply", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output == "" {
			t.Fatal("expected detection output, got empty")
		}

		// Config should NOT be updated despite --apply
		updatedCfg, _ := forgeconfig.ReadConfig(dir)
		if updatedCfg != nil && len(updatedCfg.Surfaces) > 0 {
			t.Errorf("config should NOT be written in non-interactive mode, got %v", updatedCfg.Surfaces)
		}
	})

	t.Run("non-interactive empty detection exits 1", func(t *testing.T) {
		resetDetectFlags(t)
		dir := t.TempDir()

		origIsInteractive := isInteractiveTerminalFunc
		isInteractiveTerminalFunc = func() bool { return false }
		defer func() { isInteractiveTerminalFunc = origIsInteractive }()

		var stdout bytes.Buffer
		rootCmd.SetIn(strings.NewReader(""))
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "detect", "--project-root", dir})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error (exit 1) for empty detection in non-interactive mode")
		}

		output := strings.TrimSpace(stdout.String())
		if output != "" {
			t.Errorf("expected empty stdout for empty detection, got %q", output)
		}
	})
}

// TestSurfacesDetectProjectRoot tests --project-root flag support.
// AC: --project-root flag supported (consistent with existing forge surfaces command).
func TestSurfacesDetectProjectRoot(t *testing.T) {
	resetDetectFlags(t)

	t.Run("project-root flag resolves correctly", func(t *testing.T) {
		resetDetectFlags(t)
		dir := t.TempDir()
		writeGoModWithCobra(t, dir)

		// Force non-interactive to avoid TTY dependency
		origIsInteractive := isInteractiveTerminalFunc
		isInteractiveTerminalFunc = func() bool { return false }
		defer func() { isInteractiveTerminalFunc = origIsInteractive }()

		var stdout bytes.Buffer
		rootCmd.SetIn(strings.NewReader(""))
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "detect", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if !strings.Contains(output, "cli") {
			t.Errorf("expected 'cli' in output, got %q", output)
		}
	})
}

// TestFormatDetectSourceAnnotation tests the source annotation format for detect stdout.
// AC: <source> is detected:<signal> or inferred:<rule-id>
func TestFormatDetectSourceAnnotation(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{"dependency cobra", "dependency:cobra", "detected:cobra"},
		{"dependency react", "dependency:react", "detected:react"},
		{"inference cmd-dir", "inference:cmd-dir", "inferred:cmd-dir"},
		{"inference api-dir", "inference:api-dir", "inferred:api-dir"},
		{"empty source", "", ""},
		{"unknown format", "other:value", "other:value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDetectSourceAnnotation(tt.source)
			if got != tt.want {
				t.Errorf("formatDetectSourceAnnotation(%q) = %q, want %q", tt.source, got, tt.want)
			}
		})
	}
}

// --- Test helpers for detect tests ---

func writeGoModWithCobra(t *testing.T, dir string) {
	t.Helper()
	content := "module example.com/test\n\ngo 1.25\n\nrequire github.com/spf13/cobra v1.0.0\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeGoModMinimal(t *testing.T, dir string) {
	t.Helper()
	content := "module example.com/test\n\ngo 1.25\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
