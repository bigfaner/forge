package forgeconfig

import (
	"os"
	"path/filepath"
	"testing"
)

// helper to create a config file in a temp dir
func setupConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	forgeDir := filepath.Join(dir, ".forge")
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestReadConfig(t *testing.T) {
	t.Run("file not exists returns nil nil", func(t *testing.T) {
		dir := t.TempDir()
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg != nil {
			t.Fatalf("expected nil, got %v", cfg)
		}
	})

	t.Run("empty config returns zero struct", func(t *testing.T) {
		dir := setupConfig(t, "{}")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil {
			t.Fatal("expected non-nil config")
		}
		if cfg.Auto != nil {
			t.Errorf("expected Auto nil, got %v", cfg.Auto)
		}
		if cfg.Worktree != nil {
			t.Errorf("expected Worktree nil, got %v", cfg.Worktree)
		}
	})

	t.Run("unknown fields silently ignored", func(t *testing.T) {
		dir := setupConfig(t, "project-type: backend\nlanguages:\n  - go\ntest-framework: pytest\nunknown-field: value\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil {
			t.Fatal("expected non-nil config")
		}
		// Old fields silently ignored — only auto and worktree parsed
		if cfg.Auto != nil {
			t.Errorf("expected Auto nil when not in yaml, got %v", cfg.Auto)
		}
	})

	t.Run("auto block parsed with defaults", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  e2eTest:\n    quick: false\n  gitPush: true\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Auto == nil {
			t.Fatal("expected Auto non-nil")
		}
		if cfg.Auto.E2eTest.Quick {
			t.Error("E2eTest.Quick should be false (explicitly set)")
		}
		if !cfg.Auto.E2eTest.Full {
			t.Error("E2eTest.Full should be true (default applied)")
		}
		if !cfg.Auto.GitPush {
			t.Error("GitPush should be true")
		}
	})

	t.Run("worktree block parsed", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  source-branch: develop\n  copy-files:\n    - .env\n    - .env.local\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Worktree == nil {
			t.Fatal("expected Worktree non-nil")
		}
		if cfg.Worktree.SourceBranch != "develop" {
			t.Errorf("expected source-branch 'develop', got %q", cfg.Worktree.SourceBranch)
		}
		if len(cfg.Worktree.CopyFiles) != 2 || cfg.Worktree.CopyFiles[0] != ".env" || cfg.Worktree.CopyFiles[1] != ".env.local" {
			t.Errorf("expected [.env .env.local], got %v", cfg.Worktree.CopyFiles)
		}
	})

	t.Run("worktree absent is nil", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Worktree != nil {
			t.Error("expected Worktree nil when not configured")
		}
	})
}

func TestReadAutoConfig(t *testing.T) {
	t.Run("missing config returns defaults", func(t *testing.T) {
		dir := t.TempDir()
		auto, err := ReadAutoConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Defaults: e2eTest quick=false/full=true, consolidateSpecs quick=true/full=true,
		// cleanCode false/false, validation false/false, gitPush false
		if auto.E2eTest.Quick || !auto.E2eTest.Full {
			t.Errorf("E2eTest defaults = %+v, want {Quick:false Full:true}", auto.E2eTest)
		}
		if !auto.ConsolidateSpecs.Quick || !auto.ConsolidateSpecs.Full {
			t.Errorf("ConsolidateSpecs defaults = %+v, want {Quick:true Full:true}", auto.ConsolidateSpecs)
		}
		if auto.CleanCode.Quick || auto.CleanCode.Full {
			t.Errorf("CleanCode defaults = %+v, want {Quick:false Full:false}", auto.CleanCode)
		}
		if auto.Validation.Quick || auto.Validation.Full {
			t.Errorf("Validation defaults = %+v, want {Quick:false Full:false}", auto.Validation)
		}
		if auto.GitPush {
			t.Errorf("GitPush default = %v, want false", auto.GitPush)
		}
	})

	t.Run("full auto block", func(t *testing.T) {
		dir := setupConfig(t, `auto:
  e2eTest:
    quick: false
    full: true
  consolidateSpecs:
    quick: false
    full: false
  cleanCode:
    quick: true
    full: true
  validation:
    quick: true
    full: false
  gitPush: true
`)
		auto, err := ReadAutoConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if auto.E2eTest.Quick {
			t.Error("E2eTest.Quick should be false")
		}
		if !auto.E2eTest.Full {
			t.Error("E2eTest.Full should be true")
		}
		if auto.ConsolidateSpecs.Quick {
			t.Error("ConsolidateSpecs.Quick should be false")
		}
		if auto.ConsolidateSpecs.Full {
			t.Error("ConsolidateSpecs.Full should be false")
		}
		if !auto.CleanCode.Quick {
			t.Error("CleanCode.Quick should be true")
		}
		if !auto.CleanCode.Full {
			t.Error("CleanCode.Full should be true")
		}
		if !auto.Validation.Quick {
			t.Error("Validation.Quick should be true")
		}
		if auto.Validation.Full {
			t.Error("Validation.Full should be false")
		}
		if !auto.GitPush {
			t.Error("GitPush should be true")
		}
	})

	t.Run("partial auto block applies defaults", func(t *testing.T) {
		dir := setupConfig(t, `auto:
  e2eTest:
    quick: false
`)
		auto, err := ReadAutoConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if auto.E2eTest.Quick {
			t.Error("E2eTest.Quick should be false (explicitly set)")
		}
		if !auto.E2eTest.Full {
			t.Error("E2eTest.Full should be true (default)")
		}
		if !auto.ConsolidateSpecs.Quick || !auto.ConsolidateSpecs.Full {
			t.Errorf("ConsolidateSpecs should default to true/true, got %+v", auto.ConsolidateSpecs)
		}
		if auto.CleanCode.Quick || auto.CleanCode.Full {
			t.Errorf("CleanCode should default to false/false, got %+v", auto.CleanCode)
		}
		if auto.Validation.Quick || auto.Validation.Full {
			t.Errorf("Validation should default to false/false, got %+v", auto.Validation)
		}
	})

	t.Run("no auto block returns defaults", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  source-branch: main\n")
		auto, err := ReadAutoConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if auto.E2eTest.Quick || !auto.E2eTest.Full {
			t.Errorf("E2eTest defaults = %+v, want {Quick:false Full:true}", auto.E2eTest)
		}
		if !auto.ConsolidateSpecs.Quick || !auto.ConsolidateSpecs.Full {
			t.Errorf("ConsolidateSpecs defaults = %+v, want {Quick:true Full:true}", auto.ConsolidateSpecs)
		}
	})

	t.Run("returns value type not pointer", func(t *testing.T) {
		dir := t.TempDir()
		auto, err := ReadAutoConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Verify it's a value type — this compiles only if AutoConfig is returned as value
		_ = auto.E2eTest
	})
}

func TestGetConfigValue(t *testing.T) {
	t.Run("auto.gitPush true", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		val, err := GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("auto.gitPush false", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: false\n")
		val, err := GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false', got %q", val)
		}
	})

	t.Run("auto.gitPush absent returns false (default)", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  source-branch: main\n")
		val, err := GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false' default, got %q", val)
		}
	})

	t.Run("worktree.source-branch returns value", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  source-branch: develop\n")
		val, err := GetConfigValue(dir, "worktree.source-branch")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "develop" {
			t.Errorf("expected 'develop', got %q", val)
		}
	})

	t.Run("worktree.copy-files returns joined list", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  copy-files:\n    - .env\n    - .env.local\n")
		val, err := GetConfigValue(dir, "worktree.copy-files")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := ".env\n.env.local"
		if val != expected {
			t.Errorf("expected %q, got %q", expected, val)
		}
	})

	t.Run("worktree.source-branch absent returns error", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		_, err := GetConfigValue(dir, "worktree.source-branch")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})

	t.Run("worktree.copy-files absent returns error", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		_, err := GetConfigValue(dir, "worktree.copy-files")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})

	t.Run("unknown key returns error", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		_, err := GetConfigValue(dir, "nonexistent")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})

	t.Run("missing file returns default for auto.gitPush", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false' default, got %q", val)
		}
	})

	t.Run("missing file returns error for worktree key", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "worktree.source-branch")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})

	t.Run("unknown key returns error with no file", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "something.weird")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})

	t.Run("worktree present but source-branch empty returns error", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  copy-files:\n    - .env\n")
		_, err := GetConfigValue(dir, "worktree.source-branch")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})

	t.Run("worktree present but copy-files empty returns error", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  source-branch: develop\n")
		_, err := GetConfigValue(dir, "worktree.copy-files")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})
}

func TestWriteConfig(t *testing.T) {
	t.Run("create new file", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Worktree: &WorktreeConfig{
				SourceBranch: "main",
			},
		}
		if err := WriteConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		readback, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if readback.Worktree == nil {
			t.Fatal("expected Worktree non-nil")
		}
		if readback.Worktree.SourceBranch != "main" {
			t.Errorf("expected 'main', got %q", readback.Worktree.SourceBranch)
		}
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		dir := t.TempDir()
		cfg1 := &Config{
			Worktree: &WorktreeConfig{
				SourceBranch: "develop",
			},
		}
		if err := WriteConfig(dir, cfg1); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cfg2 := &Config{
			Worktree: &WorktreeConfig{
				SourceBranch: "main",
			},
		}
		if err := WriteConfig(dir, cfg2); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		readback, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if readback.Worktree.SourceBranch != "main" {
			t.Errorf("expected 'main' after overwrite, got %q", readback.Worktree.SourceBranch)
		}
	})
}

func TestAutoConfigDefaults(t *testing.T) {
	defaults := AutoConfigDefaults()
	if defaults.E2eTest.Quick || !defaults.E2eTest.Full {
		t.Errorf("E2eTest = %+v, want {Quick:false Full:true}", defaults.E2eTest)
	}
	if !defaults.ConsolidateSpecs.Quick || !defaults.ConsolidateSpecs.Full {
		t.Errorf("ConsolidateSpecs = %+v, want {Quick:true Full:true}", defaults.ConsolidateSpecs)
	}
	if defaults.CleanCode.Quick || defaults.CleanCode.Full {
		t.Errorf("CleanCode = %+v, want {Quick:false Full:false}", defaults.CleanCode)
	}
	if defaults.Validation.Quick || defaults.Validation.Full {
		t.Errorf("Validation = %+v, want {Quick:false Full:false}", defaults.Validation)
	}
	if defaults.GitPush {
		t.Errorf("GitPush = %v, want false", defaults.GitPush)
	}
}

func TestAutoConfigIsZero(t *testing.T) {
	t.Run("zero value is zero", func(t *testing.T) {
		a := AutoConfig{}
		if !a.IsZero() {
			t.Error("expected zero AutoConfig to be zero")
		}
	})

	t.Run("non-zero is not zero", func(t *testing.T) {
		a := AutoConfigDefaults()
		if a.IsZero() {
			t.Error("expected defaults to be non-zero")
		}
	})
}

func TestAutoConfigWithDefaults(t *testing.T) {
	t.Run("zero returns full defaults", func(t *testing.T) {
		a := AutoConfig{}.WithDefaults()
		if a.E2eTest != (ModeToggle{Quick: false, Full: true}) {
			t.Errorf("E2eTest = %+v, want {Quick:false Full:true}", a.E2eTest)
		}
	})

	t.Run("partial preserves set values", func(t *testing.T) {
		a := AutoConfig{
			CleanCode: ModeToggle{Quick: true, Full: true},
		}
		result := a.WithDefaults()
		if result.CleanCode.Quick != true || result.CleanCode.Full != true {
			t.Errorf("CleanCode should be preserved as true/true, got %+v", result.CleanCode)
		}
		if result.E2eTest != (ModeToggle{Quick: false, Full: true}) {
			t.Errorf("E2eTest should default to {Quick:false Full:true}, got %+v", result.E2eTest)
		}
	})
}

func TestWorktreeConfig(t *testing.T) {
	t.Run("only source-branch", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  source-branch: main\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Worktree.SourceBranch != "main" {
			t.Errorf("expected 'main', got %q", cfg.Worktree.SourceBranch)
		}
		if cfg.Worktree.CopyFiles != nil {
			t.Errorf("expected CopyFiles nil, got %v", cfg.Worktree.CopyFiles)
		}
	})

	t.Run("only copy-files", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  copy-files:\n    - .env\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Worktree.SourceBranch != "" {
			t.Errorf("expected empty source-branch, got %q", cfg.Worktree.SourceBranch)
		}
		if len(cfg.Worktree.CopyFiles) != 1 || cfg.Worktree.CopyFiles[0] != ".env" {
			t.Errorf("expected [.env], got %v", cfg.Worktree.CopyFiles)
		}
	})
}

func TestGetConfigValueLegacyKeys(t *testing.T) {
	t.Run("test-framework returns value", func(t *testing.T) {
		dir := setupConfig(t, "test-framework: playwright\n")
		val, err := GetConfigValue(dir, "test-framework")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "playwright" {
			t.Errorf("expected 'playwright', got %q", val)
		}
	})

	t.Run("test-framework absent returns error", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		_, err := GetConfigValue(dir, "test-framework")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})

	t.Run("test-framework empty returns error", func(t *testing.T) {
		dir := setupConfig(t, "test-framework: \"\"\n")
		_, err := GetConfigValue(dir, "test-framework")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})

	t.Run("test-command returns value", func(t *testing.T) {
		dir := setupConfig(t, "test-command: npm test\n")
		val, err := GetConfigValue(dir, "test-command")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "npm test" {
			t.Errorf("expected 'npm test', got %q", val)
		}
	})

	t.Run("test-command absent returns error", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		_, err := GetConfigValue(dir, "test-command")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})

	t.Run("missing file returns error for test-framework", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "test-framework")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})
}

func TestWriteConfigAutoBlock(t *testing.T) {
	t.Run("write and read auto block", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Auto: &AutoConfig{
				E2eTest:          ModeToggle{Quick: false, Full: true},
				ConsolidateSpecs: ModeToggle{Quick: true, Full: true},
				GitPush:          true,
			},
		}
		if err := WriteConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		readback, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if readback.Auto == nil {
			t.Fatal("expected Auto non-nil")
		}
		if readback.Auto.GitPush != true {
			t.Errorf("expected GitPush true, got %v", readback.Auto.GitPush)
		}
	})
}
