package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/internal/embedded"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/just"
)

// setupInitTest creates a temp directory with optional pre-existing files.
type initTestEnv struct {
	dir    string
	stdout bytes.Buffer
	stderr bytes.Buffer
	stdin  bytes.Buffer
}

func newInitTestEnv(t *testing.T) *initTestEnv {
	t.Helper()
	return &initTestEnv{
		dir: t.TempDir(),
	}
}

func (e *initTestEnv) run(extraArgs ...string) error {
	args := []string{"init", "--project-root", e.dir}
	args = append(args, extraArgs...)
	// Reset flags to defaults to avoid state leakage between tests.
	_ = initCmd.Flags().Set("skip-just", "false")
	rootCmd.SetOut(&e.stdout)
	rootCmd.SetErr(&e.stderr)
	rootCmd.SetIn(&e.stdin)
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

func (e *initTestEnv) path(parts ...string) string {
	return filepath.Join(append([]string{e.dir}, parts...)...)
}

func TestInitCommand(t *testing.T) {
	t.Run("creates .forge directory", func(t *testing.T) {
		env := newInitTestEnv(t)
		// Answer 'n' to config init prompt (config doesn't exist yet, but stdin needs input)
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		forgeDir := env.path(feature.ForgeDir)
		info, err := os.Stat(forgeDir)
		if err != nil {
			t.Fatalf(".forge directory not created: %v", err)
		}
		if !info.IsDir() {
			t.Fatal(".forge is not a directory")
		}
	})

	t.Run("skips .forge directory when already exists", func(t *testing.T) {
		env := newInitTestEnv(t)
		forgeDir := env.path(feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(env.stdout.String(), "SKIPPED") || !strings.Contains(env.stdout.String(), ".forge") {
			t.Errorf("expected SKIPPED for .forge in output, got %q", env.stdout.String())
		}
	})

	t.Run("creates CLAUDE.md from template", func(t *testing.T) {
		env := newInitTestEnv(t)
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		claudeFile := env.path("CLAUDE.md")
		data, err := os.ReadFile(claudeFile)
		if err != nil {
			t.Fatalf("CLAUDE.md not created: %v", err)
		}

		if string(data) != embedded.CLAUDEmdTemplate {
			t.Errorf("CLAUDE.md content doesn't match template")
		}
	})

	t.Run("skips CLAUDE.md when already exists", func(t *testing.T) {
		env := newInitTestEnv(t)
		existing := "existing content"
		if err := os.WriteFile(env.path("CLAUDE.md"), []byte(existing), 0o644); err != nil {
			t.Fatal(err)
		}
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := os.ReadFile(env.path("CLAUDE.md"))
		if string(data) != existing {
			t.Error("CLAUDE.md should not be overwritten")
		}

		output := env.stdout.String()
		if !strings.Contains(output, "SKIPPED") {
			t.Errorf("expected SKIPPED status for CLAUDE.md, got %q", output)
		}
	})

	t.Run("appends entries to .gitignore", func(t *testing.T) {
		env := newInitTestEnv(t)
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(env.path(".gitignore"))
		if err != nil {
			t.Fatalf(".gitignore not created: %v", err)
		}

		content := string(data)
		expected := []string{
			"# Forge runtime",
			"docs/features/*/tasks/process/",
			".forge/state.json",
			"tests/results/.last-run.json",
			"tests/e2e/results/.last-run.json",
			"tests/e2e/results/*/error-context.md",
		}
		for _, line := range expected {
			if !strings.Contains(content, line) {
				t.Errorf(".gitignore missing %q", line)
			}
		}
	})

	t.Run("deduplicates .gitignore entries", func(t *testing.T) {
		env := newInitTestEnv(t)
		// Pre-populate with some entries already present
		existing := "# Forge runtime\n.forge/state.json\nsome-other-file\n"
		if err := os.WriteFile(env.path(".gitignore"), []byte(existing), 0o644); err != nil {
			t.Fatal(err)
		}
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := os.ReadFile(env.path(".gitignore"))
		content := string(data)

		// Count occurrences of .forge/state.json
		count := strings.Count(content, ".forge/state.json")
		if count != 1 {
			t.Errorf("expected 1 occurrence of '.forge/state.json', got %d", count)
		}
	})

	t.Run("does not create justfile", func(t *testing.T) {
		env := newInitTestEnv(t)
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// justfile should NOT be created since no recipes to append
		if _, err := os.Stat(env.path("justfile")); !os.IsNotExist(err) {
			data, _ := os.ReadFile(env.path("justfile"))
			t.Errorf("justfile should not be created, but found: %q", string(data))
		}
	})

	t.Run("runs config init when config does not exist", func(t *testing.T) {
		env := newInitTestEnv(t)
		// Input for config init: project-type=backend(2), select go-test(1), done, done
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		configFile := env.path(feature.ForgeDir, feature.ForgeConfigFileName)
		data, err := os.ReadFile(configFile)
		if err != nil {
			t.Fatalf("config.yaml not created: %v", err)
		}

		if !strings.Contains(string(data), "project-type: backend") {
			t.Errorf("expected 'project-type: backend' in config, got %q", string(data))
		}
	})

	t.Run("skips config init when config already exists", func(t *testing.T) {
		env := newInitTestEnv(t)
		forgeDir := env.path(feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		existingConfig := "project-type: frontend\n"
		if err := os.WriteFile(env.path(feature.ForgeDir, feature.ForgeConfigFileName), []byte(existingConfig), 0o644); err != nil {
			t.Fatal(err)
		}

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := os.ReadFile(env.path(feature.ForgeDir, feature.ForgeConfigFileName))
		if string(data) != existingConfig {
			t.Error("existing config should not be modified")
		}
	})

	t.Run("prints summary report", func(t *testing.T) {
		env := newInitTestEnv(t)
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := env.stdout.String()
		// Should contain status markers
		if !strings.Contains(output, ">>>") || !strings.Contains(output, "<<<") {
			t.Errorf("expected summary block markers, got %q", output)
		}
		// Should have CREATED or APPENDED actions
		if !strings.Contains(output, "CREATED") && !strings.Contains(output, "APPENDED") {
			t.Errorf("expected action status in output, got %q", output)
		}
	})

	t.Run("full init on empty directory creates all artifacts", func(t *testing.T) {
		env := newInitTestEnv(t)
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify all artifacts exist
		checks := []struct {
			name string
			path string
		}{
			{".forge directory", feature.ForgeDir},
			{"CLAUDE.md", "CLAUDE.md"},
			{".gitignore", ".gitignore"},
			{"config.yaml", filepath.Join(feature.ForgeDir, feature.ForgeConfigFileName)},
		}

		for _, check := range checks {
			p := env.path(check.path)
			if _, err := os.Stat(p); err != nil {
				t.Errorf("%s not found at %s: %v", check.name, p, err)
			}
		}
	})
}

func TestAppendGitignoreEntries(t *testing.T) {
	t.Run("dedup checks each line individually", func(t *testing.T) {
		existing := "# Forge runtime\n.forge/state.json\nsome-other-file\n"
		result := buildGitignoreAppend(existing, []string{
			"# Forge runtime",
			"docs/features/*/tasks/process/",
			".forge/state.json",
			"tests/results/.last-run.json",
		})

		containsEntry := func(entries []string, target string) bool {
			for _, e := range entries {
				if strings.TrimSpace(e) == strings.TrimSpace(target) {
					return true
				}
			}
			return false
		}

		if containsEntry(result, "# Forge runtime") {
			t.Error("should not include duplicate '# Forge runtime'")
		}
		if containsEntry(result, ".forge/state.json") {
			t.Error("should not include duplicate '.forge/state.json'")
		}
		if !containsEntry(result, "docs/features/*/tasks/process/") {
			t.Error("should include new entry 'docs/features/*/tasks/process/'")
		}
		if !containsEntry(result, "tests/results/.last-run.json") {
			t.Error("should include new entry 'tests/results/.last-run.json'")
		}
	})

	t.Run("returns all entries when file is empty", func(t *testing.T) {
		result := buildGitignoreAppend("", []string{"a", "b"})
		if len(result) != 2 {
			t.Errorf("expected 2 entries, got %d", len(result))
		}
	})
}

func TestInitSkipJustFlag(t *testing.T) {
	t.Run("--skip-just reports SKIPPED for just step", func(t *testing.T) {
		env := newInitTestEnv(t)
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run("--skip-just")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := env.stdout.String()
		if !strings.Contains(output, "SKIPPED") || !strings.Contains(output, "just") {
			t.Errorf("expected SKIPPED status for just in output, got %q", output)
		}
	})

	t.Run("--skip-just still runs all other steps", func(t *testing.T) {
		env := newInitTestEnv(t)
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run("--skip-just")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// All other artifacts should still be created
		checks := []struct {
			name string
			path string
		}{
			{".forge directory", feature.ForgeDir},
			{"CLAUDE.md", "CLAUDE.md"},
			{".gitignore", ".gitignore"},
			{"config.yaml", filepath.Join(feature.ForgeDir, feature.ForgeConfigFileName)},
		}

		for _, check := range checks {
			p := env.path(check.path)
			if _, err := os.Stat(p); err != nil {
				t.Errorf("%s not found at %s: %v", check.name, p, err)
			}
		}
	})

	t.Run("ensureJust step appears in summary", func(t *testing.T) {
		env := newInitTestEnv(t)
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run("--skip-just")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := env.stdout.String()
		if !strings.Contains(output, "just installation") {
			t.Fatalf("expected 'just installation' in output, got %q", output)
		}
	})
}

func TestInitEnsureJustIntegration(t *testing.T) {
	t.Run("just already installed reports SKIPPED", func(t *testing.T) {
		origEnsure := ensureJustFunc
		ensureJustFunc = func(_ io.Reader, _ io.Writer) just.EnsureResult {
			return just.EnsureResult{
				Status:  just.StatusSkipped,
				Version: "1.40.0",
				Detail:  "just 1.40.0 found at /usr/bin/just (meets minimum 1.40.0)",
			}
		}
		defer func() { ensureJustFunc = origEnsure }()

		env := newInitTestEnv(t)
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := env.stdout.String()
		if !strings.Contains(output, "SKIPPED") || !strings.Contains(output, "just") {
			t.Errorf("expected SKIPPED for just (already installed), got %q", output)
		}
	})

	t.Run("just not installed triggers installation attempt", func(t *testing.T) {
		origEnsure := ensureJustFunc
		ensureJustFunc = func(_ io.Reader, _ io.Writer) just.EnsureResult {
			return just.EnsureResult{
				Status:  just.StatusInstalled,
				Version: "1.40.0",
				Method:  "brew",
				Detail:  "installed via brew",
			}
		}
		defer func() { ensureJustFunc = origEnsure }()

		env := newInitTestEnv(t)
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := env.stdout.String()
		if !strings.Contains(output, "INSTALLED") || !strings.Contains(output, "just") {
			t.Errorf("expected INSTALLED for just, got %q", output)
		}
	})

	t.Run("just installation failure is non-blocking", func(t *testing.T) {
		origEnsure := ensureJustFunc
		ensureJustFunc = func(_ io.Reader, _ io.Writer) just.EnsureResult {
			return just.EnsureResult{
				Status: just.StatusFailed,
				Detail: "no supported package manager found",
			}
		}
		defer func() { ensureJustFunc = origEnsure }()

		env := newInitTestEnv(t)
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := env.stdout.String()
		if !strings.Contains(output, "FAILED") {
			t.Errorf("expected FAILED for just installation, got %q", output)
		}

		// Other steps should still complete
		if !strings.Contains(output, "CREATED") {
			t.Errorf("expected other steps to still succeed, got %q", output)
		}
	})
}

func TestEnsureResultToAction(t *testing.T) {
	tests := []struct {
		name   string
		result just.EnsureResult
		status string
		target string
	}{
		{
			name:   "installed status",
			result: just.EnsureResult{Status: just.StatusInstalled, Version: "1.40.0", Method: "brew", Detail: "installed via brew"},
			status: "INSTALLED",
			target: "just installation",
		},
		{
			name:   "skipped status",
			result: just.EnsureResult{Status: just.StatusSkipped, Version: "1.40.0", Detail: "just 1.40.0 found"},
			status: "SKIPPED",
			target: "just installation",
		},
		{
			name:   "failed status",
			result: just.EnsureResult{Status: just.StatusFailed, Detail: "no package manager"},
			status: "FAILED",
			target: "just installation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := ensureResultToAction(tt.result)
			if action.status != tt.status {
				t.Errorf("expected status %q, got %q", tt.status, action.status)
			}
			if action.target != tt.target {
				t.Errorf("expected target %q, got %q", tt.target, action.target)
			}
		})
	}
}
