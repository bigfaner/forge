package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/feature"
)

// surfacesTestHelper creates a temp dir with .forge/config.yaml containing the given content.
func surfacesTestHelper(t *testing.T, content string) string {
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

// resetSurfacesFlag resets the --types flag to avoid state leaking between tests.
func resetSurfacesFlag(t *testing.T) {
	t.Helper()
	surfacesTypesFlag = false
}

// TestSurfacesScalarForm tests `forge surfaces` with scalar form config.
func TestSurfacesScalarForm(t *testing.T) {
	resetSurfacesFlag(t)

	t.Run("outputs single type with exit 0", func(t *testing.T) {
		dir := surfacesTestHelper(t, "surfaces: api\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "api" {
			t.Errorf("expected 'api', got %q", output)
		}
	})

	t.Run("any path returns single value with exit 0", func(t *testing.T) {
		resetSurfacesFlag(t)
		dir := surfacesTestHelper(t, "surfaces: api\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "src/main.go", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "api" {
			t.Errorf("expected 'api', got %q", output)
		}
	})
}

// TestSurfacesMapForm tests `forge surfaces` with map form config.
func TestSurfacesMapForm(t *testing.T) {
	resetSurfacesFlag(t)

	t.Run("outputs path=surface per line", func(t *testing.T) {
		dir := surfacesTestHelper(t, "surfaces:\n  frontend: web\n  backend: api\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		lines := strings.Split(output, "\n")
		if len(lines) != 2 {
			t.Fatalf("expected 2 lines, got %d: %q", len(lines), output)
		}
		// Verify each line contains path=surface format
		seen := map[string]bool{}
		for _, line := range lines {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				t.Errorf("expected path=surface format, got %q", line)
				continue
			}
			seen[parts[0]+"="+parts[1]] = true
		}
		if !seen["frontend=web"] {
			t.Errorf("missing 'frontend=web', got: %s", output)
		}
		if !seen["backend=api"] {
			t.Errorf("missing 'backend=api', got: %s", output)
		}
	})

	t.Run("path query returns surface type", func(t *testing.T) {
		resetSurfacesFlag(t)
		dir := surfacesTestHelper(t, "surfaces:\n  frontend: web\n  backend: api\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "frontend/src", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "web" {
			t.Errorf("expected 'web', got %q", output)
		}
	})

	t.Run("path query unmatched returns exit 1 and stderr", func(t *testing.T) {
		resetSurfacesFlag(t)
		dir := surfacesTestHelper(t, "surfaces:\n  frontend: web\n")

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"surfaces", "unknown-dir", "--project-root", dir})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for unmatched path")
		}

		// stderr should contain the error message with config hint
		stderrContent := stderr.String()
		if !strings.Contains(stderrContent, "forge init") {
			t.Errorf("error should mention 'forge init', got stderr: %s", stderrContent)
		}
		// stdout should be empty
		if strings.TrimSpace(stdout.String()) != "" {
			t.Errorf("stdout should be empty on error, got: %q", stdout.String())
		}
	})
}

// TestSurfacesTypes tests `forge surfaces --types`.
func TestSurfacesTypes(t *testing.T) {
	resetSurfacesFlag(t)

	t.Run("returns space-separated deduplicated types", func(t *testing.T) {
		dir := surfacesTestHelper(t, "surfaces:\n  frontend: web\n  backend: api\n  cli: cli\n  shared: web\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "--types", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		types := strings.Fields(output)
		// Should be deduplicated: web, api, cli (order not guaranteed but must have 3)
		if len(types) != 3 {
			t.Errorf("expected 3 types, got %d: %q", len(types), output)
		}
		seen := map[string]bool{}
		for _, typ := range types {
			if seen[typ] {
				t.Errorf("duplicate type: %q", typ)
			}
			seen[typ] = true
		}
		if !seen["web"] || !seen["api"] || !seen["cli"] {
			t.Errorf("expected web, api, cli; got: %s", output)
		}
	})

	t.Run("scalar form returns single type", func(t *testing.T) {
		resetSurfacesFlag(t)
		dir := surfacesTestHelper(t, "surfaces: api\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "--types", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "api" {
			t.Errorf("expected 'api', got %q", output)
		}
	})
}

// TestSurfacesInterfacesIgnored tests that `interfaces` field is completely ignored.
func TestSurfacesInterfacesIgnored(t *testing.T) {
	resetSurfacesFlag(t)

	t.Run("interfaces field completely ignored", func(t *testing.T) {
		dir := surfacesTestHelper(t, "surfaces:\n  frontend: web\ninterfaces:\n  - api\n  - cli\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "--types", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		// Should only contain web, not api or cli (from interfaces)
		if output != "web" {
			t.Errorf("expected only 'web' (interfaces ignored), got %q", output)
		}
	})
}

// TestSurfacesEdgeCases tests edge cases.
func TestSurfacesEdgeCases(t *testing.T) {
	resetSurfacesFlag(t)

	t.Run("empty surfaces outputs nothing with exit 0", func(t *testing.T) {
		resetSurfacesFlag(t)
		dir := surfacesTestHelper(t, "surfaces: {}\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "" {
			t.Errorf("expected empty output for empty surfaces, got %q", output)
		}
	})

	t.Run("no config file outputs nothing with exit 0", func(t *testing.T) {
		resetSurfacesFlag(t)
		dir := t.TempDir()

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "" {
			t.Errorf("expected empty output for missing config, got %q", output)
		}
	})

	t.Run("path query on empty surfaces returns exit 1", func(t *testing.T) {
		resetSurfacesFlag(t)
		dir := surfacesTestHelper(t, "surfaces: {}\n")

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"surfaces", "frontend", "--project-root", dir})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for path query on empty surfaces")
		}

		if !strings.Contains(stderr.String(), "forge init") {
			t.Errorf("error should mention 'forge init', got: %s", stderr.String())
		}
	})

	t.Run("output is raw text parseable by scripts", func(t *testing.T) {
		resetSurfacesFlag(t)
		dir := surfacesTestHelper(t, "surfaces: api\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := stdout.String()
		// No extra formatting, no block markers
		if strings.Contains(output, "```") {
			t.Errorf("output should not contain formatting blocks, got %q", output)
		}
		if strings.HasPrefix(output, "> ") {
			t.Errorf("output should not have block markers, got %q", output)
		}
	})

	t.Run("--types on empty surfaces outputs nothing", func(t *testing.T) {
		resetSurfacesFlag(t)
		dir := surfacesTestHelper(t, "surfaces: {}\n")

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs([]string{"surfaces", "--types", "--project-root", dir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "" {
			t.Errorf("expected empty output for --types on empty surfaces, got %q", output)
		}
	})

	t.Run("dotdot path returns exit 1", func(t *testing.T) {
		resetSurfacesFlag(t)
		dir := surfacesTestHelper(t, "surfaces:\n  frontend: web\n")

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"surfaces", "../etc", "--project-root", dir})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for path with '..'")
		}

		if !strings.Contains(stderr.String(), "'..'") {
			t.Errorf("error should mention '..', got: %s", stderr.String())
		}
	})
}
