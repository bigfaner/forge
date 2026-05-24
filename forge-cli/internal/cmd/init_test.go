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
	"forge-cli/pkg/forgeconfig"
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
	orig := configInitFunc
	origSurf := surfaceConfigFunc
	configInitFunc = testConfigInit
	surfaceConfigFunc = testSurfaceConfig
	t.Cleanup(func() {
		configInitFunc = orig
		surfaceConfigFunc = origSurf
	})
	return &initTestEnv{
		dir: t.TempDir(),
	}
}

// testSurfaceConfig replaces surfaceConfigFunc for testing.
// Simulates surface detection without requiring a real TTY.
func testSurfaceConfig(_ string) initAction {
	return initAction{status: "SKIPPED", target: "surfaces", detail: "test override"}
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

// testConfigInit replaces configInitFunc for testing.
// Simulates the interactive config flow without requiring a real TTY.
func testConfigInit(projectRoot string) initAction {
	configFile := filepath.Join(projectRoot, feature.ForgeDir, feature.ForgeConfigFileName)

	// Write a sensible default config for testing
	auto := autoConfigDefaults()
	auto.GitPush = true // Explicitly set to true to differentiate from empty/zero configs
	auto.Validation = forgeconfig.ModeToggle{Quick: true, Full: true}
	cfg := forgeconfig.Config{
		Auto: auto,
	}

	if err := writeConfigFile(configFile, &cfg); err != nil {
		return initAction{status: "FAILED", target: ".forge/config.yaml", detail: err.Error()}
	}

	detail := "test override"
	if _, err := os.Stat(configFile); err == nil {
		// Config existed but we overwrote it (reconfigure path)
		detail = "reconfigured"
	}

	return initAction{status: "CREATED", target: ".forge/config.yaml", detail: detail}
}

// autoConfigDefaults returns a default AutoConfig for tests.
func autoConfigDefaults() *forgeconfig.AutoConfig {
	d := forgeconfig.AutoConfigDefaults()
	return &d
}

func TestInitCommand(t *testing.T) {
	t.Run("creates .forge directory", func(t *testing.T) {
		env := newInitTestEnv(t)

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
			"# Forge",
			".forge/state.json",
			".forge/test-state.json",
			".forge/worktrees/",
			"docs/features/*/tasks/process/",
			"docs/features/*/tasks/index.json.lock",
			"docs/features/*/testing/results/",
			"tests/results/",
			"tests/e2e/results/",
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

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		configFile := env.path(feature.ForgeDir, feature.ForgeConfigFileName)
		data, err := os.ReadFile(configFile)
		if err != nil {
			t.Fatalf("config.yaml not created: %v", err)
		}

		if !strings.Contains(string(data), "auto:") {
			t.Errorf("expected 'auto:' in config, got %q", string(data))
		}
	})

	t.Run("overwrites existing config (reconfigure)", func(t *testing.T) {
		env := newInitTestEnv(t)
		forgeDir := env.path(feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		existingConfig := "auto:\n  gitPush: false\n"
		if err := os.WriteFile(env.path(feature.ForgeDir, feature.ForgeConfigFileName), []byte(existingConfig), 0o644); err != nil {
			t.Fatal(err)
		}

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := os.ReadFile(env.path(feature.ForgeDir, feature.ForgeConfigFileName))
		if strings.Contains(string(data), "gitPush: false") {
			t.Error("existing config should have been overwritten")
		}
		if !strings.Contains(string(data), "auto:") {
			t.Error("expected reconfigured config to contain auto")
		}
	})

	t.Run("prints summary report", func(t *testing.T) {
		env := newInitTestEnv(t)

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

func TestInitConfigWithValidation(t *testing.T) {
	t.Run("config includes validation field", func(t *testing.T) {
		env := newInitTestEnv(t)

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		configFile := env.path(feature.ForgeDir, feature.ForgeConfigFileName)
		data, err := os.ReadFile(configFile)
		if err != nil {
			t.Fatalf("config.yaml not created: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "validation:") {
			t.Errorf("expected 'validation:' in config, got %q", content)
		}
		if !strings.Contains(content, "quick: true") {
			t.Errorf("expected 'quick: true' for validation in config, got %q", content)
		}
		if !strings.Contains(content, "full: true") {
			t.Errorf("expected 'full: true' for validation in config, got %q", content)
		}
	})
}

func TestInitConfigWithWorktree(t *testing.T) {
	t.Run("config includes worktree when provided", func(t *testing.T) {
		orig := configInitFunc
		origSurf := surfaceConfigFunc
		configInitFunc = func(projectRoot string) initAction {
			configFile := filepath.Join(projectRoot, feature.ForgeDir, feature.ForgeConfigFileName)
			auto := autoConfigDefaults()
			cfg := forgeconfig.Config{
				Auto: auto,
				Worktree: &forgeconfig.WorktreeConfig{
					SourceBranch: "main",
					CopyFiles:    []string{".env", ".env.local"},
				},
			}
			if err := writeConfigFile(configFile, &cfg); err != nil {
				return initAction{status: "FAILED", target: ".forge/config.yaml", detail: err.Error()}
			}
			return initAction{status: "CREATED", target: ".forge/config.yaml", detail: "with worktree"}
		}
		surfaceConfigFunc = testSurfaceConfig
		defer func() {
			configInitFunc = orig
			surfaceConfigFunc = origSurf
		}()

		env := &initTestEnv{dir: t.TempDir()}
		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		configFile := env.path(feature.ForgeDir, feature.ForgeConfigFileName)
		data, err := os.ReadFile(configFile)
		if err != nil {
			t.Fatalf("config.yaml not created: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "worktree:") {
			t.Errorf("expected 'worktree:' in config, got %q", content)
		}
		if !strings.Contains(content, "source-branch: main") {
			t.Errorf("expected 'source-branch: main' in config, got %q", content)
		}
		if !strings.Contains(content, ".env") {
			t.Errorf("expected '.env' in copy-files, got %q", content)
		}
	})

	t.Run("config omits worktree when both fields empty", func(t *testing.T) {
		orig := configInitFunc
		origSurf := surfaceConfigFunc
		configInitFunc = func(projectRoot string) initAction {
			configFile := filepath.Join(projectRoot, feature.ForgeDir, feature.ForgeConfigFileName)
			auto := autoConfigDefaults()
			cfg := forgeconfig.Config{
				Auto:     auto,
				Worktree: nil, // No worktree when skipped
			}
			if err := writeConfigFile(configFile, &cfg); err != nil {
				return initAction{status: "FAILED", target: ".forge/config.yaml", detail: err.Error()}
			}
			return initAction{status: "CREATED", target: ".forge/config.yaml", detail: "no worktree"}
		}
		surfaceConfigFunc = testSurfaceConfig
		defer func() {
			configInitFunc = orig
			surfaceConfigFunc = origSurf
		}()

		env := &initTestEnv{dir: t.TempDir()}
		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		configFile := env.path(feature.ForgeDir, feature.ForgeConfigFileName)
		data, err := os.ReadFile(configFile)
		if err != nil {
			t.Fatalf("config.yaml not created: %v", err)
		}

		content := string(data)
		if strings.Contains(content, "worktree:") {
			t.Errorf("worktree block should not be present when skipped, got %q", content)
		}
	})
}

func TestAutoBehaviorPrompts_NoE2EReferences(t *testing.T) {
	t.Run("prompt titles do not reference e2e tests", func(t *testing.T) {
		defaults := forgeconfig.AutoConfigDefaults()
		prompts := autoBehaviorPrompts(defaults)

		for _, p := range prompts {
			lower := strings.ToLower(p.title)
			if strings.Contains(lower, "e2e") || strings.Contains(lower, "end-to-end") {
				t.Errorf("prompt title should not reference e2e/end-to-end, got: %q", p.title)
			}
			lowerDesc := strings.ToLower(p.desc)
			if strings.Contains(lowerDesc, "e2e") || strings.Contains(lowerDesc, "end-to-end") {
				t.Errorf("prompt desc should not reference e2e/end-to-end, got: %q", p.desc)
			}
		}
	})

	t.Run("test prompts reference 'test' not 'e2e-test'", func(t *testing.T) {
		defaults := forgeconfig.AutoConfigDefaults()
		prompts := autoBehaviorPrompts(defaults)

		// Find the test-related prompts (first two)
		testPrompts := prompts[:2]
		for _, p := range testPrompts {
			lower := strings.ToLower(p.title)
			if !strings.Contains(lower, "test") {
				t.Errorf("test prompt title should contain 'test', got: %q", p.title)
			}
		}
	})
}

func TestWorktreeConfigRoundTrip(t *testing.T) {
	t.Run("worktree config round-trips through YAML", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &forgeconfig.Config{
			Auto: func() *forgeconfig.AutoConfig {
				d := forgeconfig.AutoConfigDefaults()
				return &d
			}(),
			Worktree: &forgeconfig.WorktreeConfig{
				SourceBranch: "develop",
				CopyFiles:    []string{".env", ".env.local"},
			},
		}

		configFile := filepath.Join(dir, feature.ForgeDir, feature.ForgeConfigFileName)
		if err := writeConfigFile(configFile, cfg); err != nil {
			t.Fatalf("writeConfigFile failed: %v", err)
		}

		read, err := forgeconfig.ReadConfig(dir)
		if err != nil {
			t.Fatalf("ReadConfig failed: %v", err)
		}

		if read.Worktree == nil {
			t.Fatal("worktree should not be nil")
		}
		if read.Worktree.SourceBranch != "develop" {
			t.Errorf("expected source-branch 'develop', got %q", read.Worktree.SourceBranch)
		}
		if len(read.Worktree.CopyFiles) != 2 {
			t.Errorf("expected 2 copy-files, got %d", len(read.Worktree.CopyFiles))
		}
		if read.Worktree.CopyFiles[0] != ".env" || read.Worktree.CopyFiles[1] != ".env.local" {
			t.Errorf("copy-files mismatch: %v", read.Worktree.CopyFiles)
		}
	})
}
