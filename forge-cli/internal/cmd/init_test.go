package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/internal/embedded"
	"forge-cli/pkg/feature"
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

func (e *initTestEnv) run() error {
	rootCmd.SetOut(&e.stdout)
	rootCmd.SetErr(&e.stderr)
	rootCmd.SetIn(&e.stdin)
	rootCmd.SetArgs([]string{"init", "--project-root", e.dir})
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

	t.Run("appends recipes to justfile", func(t *testing.T) {
		env := newInitTestEnv(t)
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(env.path("justfile"))
		if err != nil {
			t.Fatalf("justfile not created: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "claude:") {
			t.Error("justfile missing 'claude:' recipe")
		}
		if !strings.Contains(content, "claude --dangerously-skip-permissions") {
			t.Error("justfile missing claude recipe content")
		}
		if !strings.Contains(content, "claude-c:") {
			t.Error("justfile missing 'claude-c:' recipe")
		}
	})

	t.Run("deduplicates justfile recipes", func(t *testing.T) {
		env := newInitTestEnv(t)
		existing := "build:\n    go build ./...\n\nclaude:\n    claude --dangerously-skip-permissions\n"
		if err := os.WriteFile(env.path("justfile"), []byte(existing), 0o644); err != nil {
			t.Fatal(err)
		}
		env.stdin.WriteString("2\n1\n\n\n")

		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, _ := os.ReadFile(env.path("justfile"))
		content := string(data)

		// Count occurrences of "claude:" as line prefix (recipe name)
		count := strings.Count(content, "claude:\n")
		if count != 1 {
			t.Errorf("expected 1 occurrence of 'claude:' recipe, got %d", count)
		}

		// claude-c should be added
		if !strings.Contains(content, "claude-c:") {
			t.Error("justfile should have claude-c recipe added")
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
			{"justfile", "justfile"},
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

func TestBuildJustfileAppend(t *testing.T) {
	t.Run("skips recipe that already exists", func(t *testing.T) {
		existing := "build:\n    go build\n\nclaude:\n    claude --dangerously-skip-permissions\n"
		result := buildJustfileAppend(existing)

		if strings.Contains(result, "claude:") {
			t.Error("should skip 'claude' recipe when it already exists")
		}
		if !strings.Contains(result, "claude-c:") {
			t.Error("should include 'claude-c' recipe when it doesn't exist")
		}
	})

	t.Run("includes both recipes when neither exists", func(t *testing.T) {
		existing := "build:\n    go build\n"
		result := buildJustfileAppend(existing)

		if !strings.Contains(result, "claude:") {
			t.Error("should include 'claude' recipe")
		}
		if !strings.Contains(result, "claude-c:") {
			t.Error("should include 'claude-c' recipe")
		}
	})

	t.Run("skips both recipes when both exist", func(t *testing.T) {
		existing := "claude:\n    claude --dangerously-skip-permissions\n\nclaude-c:\n    claude --dangerously-skip-permissions -c\n"
		result := buildJustfileAppend(existing)

		if len(result) != 0 {
			t.Errorf("expected no recipes to append, got %d", len(result))
		}
	})
}
