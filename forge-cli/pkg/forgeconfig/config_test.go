package forgeconfig

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	t.Run("version field parsed correctly", func(t *testing.T) {
		dir := setupConfig(t, "version: \"2\"\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil {
			t.Fatal("expected non-nil config")
		}
		if cfg.Version != "2" {
			t.Errorf("expected Version '2', got %q", cfg.Version)
		}
	})

	t.Run("config without version defaults to 1", func(t *testing.T) {
		dir := setupConfig(t, "test-framework: pytest\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil {
			t.Fatal("expected non-nil config")
		}
		if cfg.Version != "1" {
			t.Errorf("expected Version '1' (default), got %q", cfg.Version)
		}
	})

	t.Run("empty config defaults version to 1", func(t *testing.T) {
		dir := setupConfig(t, "{}\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil {
			t.Fatal("expected non-nil config")
		}
		if cfg.Version != "1" {
			t.Errorf("expected Version '1' (default), got %q", cfg.Version)
		}
	})

	t.Run("project-type field parsed correctly", func(t *testing.T) {
		dir := setupConfig(t, "project-type: fullstack\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil {
			t.Fatal("expected non-nil config")
		}
		if cfg.ProjectType != "fullstack" {
			t.Errorf("expected ProjectType 'fullstack', got %q", cfg.ProjectType)
		}
	})

	t.Run("project-type valid values accepted", func(t *testing.T) {
		for _, pt := range []string{"fullstack", "mobile", "library", "mixed"} {
			t.Run(pt, func(t *testing.T) {
				dir := setupConfig(t, "project-type: "+pt+"\n")
				cfg, err := ReadConfig(dir)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if cfg.ProjectType != pt {
					t.Errorf("expected ProjectType %q, got %q", pt, cfg.ProjectType)
				}
			})
		}
	})

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

	t.Run("known fields parsed while unknown silently ignored", func(t *testing.T) {
		dir := setupConfig(t, "project-type: fullstack\nlanguages:\n  - go\ntest-framework: pytest\nunknown-field: value\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil {
			t.Fatal("expected non-nil config")
		}
		if cfg.ProjectType != "fullstack" {
			t.Errorf("expected ProjectType 'fullstack', got %q", cfg.ProjectType)
		}
		if cfg.TestFramework != "pytest" {
			t.Errorf("expected TestFramework 'pytest', got %q", cfg.TestFramework)
		}
		if len(cfg.Languages) != 1 || cfg.Languages[0] != "go" {
			t.Errorf("expected Languages [go], got %v", cfg.Languages)
		}
	})

	t.Run("auto block parsed with defaults", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  test:\n    quick: false\n  gitPush: true\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Auto == nil {
			t.Fatal("expected Auto non-nil")
		}
		if cfg.Auto.Test.Quick {
			t.Error("Test.Quick should be false (explicitly set)")
		}
		if !cfg.Auto.Test.Full {
			t.Error("Test.Full should be true (default applied)")
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
		// Defaults: test quick=false/full=true, consolidateSpecs quick=true/full=true,
		// cleanCode false/false, validation false/false, gitPush false
		if auto.Test.Quick || !auto.Test.Full {
			t.Errorf("Test defaults = %+v, want {Quick:false Full:true}", auto.Test)
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
		if !auto.RunTasks.Quick || auto.RunTasks.Full {
			t.Errorf("RunTasks defaults = %+v, want {Quick:true Full:false}", auto.RunTasks)
		}
		if !auto.KnowledgeSave.Quick || auto.KnowledgeSave.Full {
			t.Errorf("KnowledgeSave defaults = %+v, want {Quick:true Full:false}", auto.KnowledgeSave)
		}
	})

	t.Run("full auto block", func(t *testing.T) {
		dir := setupConfig(t, `auto:
  test:
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
  runTasks:
    quick: false
    full: true
  gitPush: true
`)
		auto, err := ReadAutoConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if auto.Test.Quick {
			t.Error("Test.Quick should be false")
		}
		if !auto.Test.Full {
			t.Error("Test.Full should be true")
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
		if auto.RunTasks.Quick {
			t.Error("RunTasks.Quick should be false")
		}
		if !auto.RunTasks.Full {
			t.Error("RunTasks.Full should be true")
		}
	})

	t.Run("partial auto block applies defaults", func(t *testing.T) {
		dir := setupConfig(t, `auto:
  test:
    quick: false
`)
		auto, err := ReadAutoConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if auto.Test.Quick {
			t.Error("Test.Quick should be false (explicitly set)")
		}
		if !auto.Test.Full {
			t.Error("Test.Full should be true (default)")
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
		if auto.Test.Quick || !auto.Test.Full {
			t.Errorf("Test defaults = %+v, want {Quick:false Full:true}", auto.Test)
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
		_ = auto.Test
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
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})

	t.Run("worktree.copy-files absent returns error", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		_, err := GetConfigValue(dir, "worktree.copy-files")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})

	t.Run("unknown key returns error", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		_, err := GetConfigValue(dir, "nonexistent")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
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
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})

	t.Run("unknown key returns error with no file", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "something.weird")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})

	t.Run("worktree present but source-branch empty returns error", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  copy-files:\n    - .env\n")
		_, err := GetConfigValue(dir, "worktree.source-branch")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})

	t.Run("worktree present but copy-files empty returns error", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  source-branch: develop\n")
		_, err := GetConfigValue(dir, "worktree.copy-files")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
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
		if err := writeConfig(dir, cfg); err != nil {
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
		if err := writeConfig(dir, cfg1); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cfg2 := &Config{
			Worktree: &WorktreeConfig{
				SourceBranch: "main",
			},
		}
		if err := writeConfig(dir, cfg2); err != nil {
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

func TestWriteConfig_VersionRoundtrip(t *testing.T) {
	t.Run("write and read back preserves version", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Version: "2",
		}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		readback, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if readback.Version != "2" {
			t.Errorf("expected Version '2', got %q", readback.Version)
		}
	})
}

func TestWriteConfig_ProjectTypeRoundtrip(t *testing.T) {
	t.Run("write and read back preserves project-type", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			ProjectType: "mobile",
		}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		readback, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if readback.ProjectType != "mobile" {
			t.Errorf("expected ProjectType 'mobile', got %q", readback.ProjectType)
		}
	})
}

func TestValidProjectType(t *testing.T) {
	t.Run("valid types return true", func(t *testing.T) {
		for _, pt := range []string{"fullstack", "mobile", "library", "mixed"} {
			if !ValidProjectType(pt) {
				t.Errorf("expected %q to be valid", pt)
			}
		}
	})

	t.Run("invalid types return false", func(t *testing.T) {
		for _, pt := range []string{"frontend", "backend", "", "unknown", "full stack"} {
			if ValidProjectType(pt) {
				t.Errorf("expected %q to be invalid", pt)
			}
		}
	})
}

func TestAutoConfigDefaults(t *testing.T) {
	defaults := AutoConfigDefaults()
	if defaults.Test.Quick || !defaults.Test.Full {
		t.Errorf("Test = %+v, want {Quick:false Full:true}", defaults.Test)
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
	if !defaults.RunTasks.Quick || defaults.RunTasks.Full {
		t.Errorf("RunTasks = %+v, want {Quick:true Full:false}", defaults.RunTasks)
	}
	if !defaults.KnowledgeSave.Quick || defaults.KnowledgeSave.Full {
		t.Errorf("KnowledgeSave = %+v, want {Quick:true Full:false}", defaults.KnowledgeSave)
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

	t.Run("only RunTasks set is not zero", func(t *testing.T) {
		a := AutoConfig{RunTasks: ModeToggle{Quick: true, Full: false}}
		if a.IsZero() {
			t.Error("expected AutoConfig with RunTasks set to be non-zero")
		}
	})

	t.Run("only KnowledgeSave set is not zero", func(t *testing.T) {
		a := AutoConfig{KnowledgeSave: ModeToggle{Quick: true, Full: false}}
		if a.IsZero() {
			t.Error("expected AutoConfig with KnowledgeSave set to be non-zero")
		}
	})
}

func TestAutoConfigWithDefaults(t *testing.T) {
	t.Run("zero returns full defaults", func(t *testing.T) {
		a := AutoConfig{}.WithDefaults()
		if a.Test != (ModeToggle{Quick: false, Full: true}) {
			t.Errorf("Test = %+v, want {Quick:false Full:true}", a.Test)
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
		if result.Test != (ModeToggle{}) {
			t.Errorf("Test should be returned unchanged for partial config, got %+v", result.Test)
		}
	})

	t.Run("RunTasks defaults to quick:true full:false", func(t *testing.T) {
		a := AutoConfig{}.WithDefaults()
		if a.RunTasks.Quick != true || a.RunTasks.Full != false {
			t.Errorf("RunTasks = %+v, want {Quick:true Full:false}", a.RunTasks)
		}
	})

	t.Run("RunTasks preserved when set", func(t *testing.T) {
		a := AutoConfig{
			RunTasks: ModeToggle{Quick: false, Full: true},
		}
		result := a.WithDefaults()
		if result.RunTasks.Quick != false || result.RunTasks.Full != true {
			t.Errorf("RunTasks should be preserved as false/true, got %+v", result.RunTasks)
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
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})

	t.Run("test-framework empty returns error", func(t *testing.T) {
		dir := setupConfig(t, "test-framework: \"\"\n")
		_, err := GetConfigValue(dir, "test-framework")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})

	t.Run("test-command returns error (removed field)", func(t *testing.T) {
		dir := setupConfig(t, "test-command: npm test\n")
		_, err := GetConfigValue(dir, "test-command")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound for removed test-command key, got %v", err)
		}
	})

	t.Run("missing file returns error for test-framework", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "test-framework")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})
}

func TestReadConfig_CoverageBlock(t *testing.T) {
	t.Run("coverage block parsed with defaults", func(t *testing.T) {
		dir := setupConfig(t, `coverage:
  coding.feature:
    type: percentage
    percentage: 90
  coding.refactor:
    type: maintain
`)
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Coverage == nil {
			t.Fatal("expected Coverage non-nil")
		}
		strategy, ok := cfg.Coverage.ByType["coding.feature"]
		if !ok {
			t.Fatal("expected coding.feature strategy")
		}
		if strategy.Type != "percentage" {
			t.Errorf("coding.feature type = %q, want percentage", strategy.Type)
		}
		if strategy.Percentage == nil || *strategy.Percentage != 90 {
			t.Errorf("coding.feature percentage = %v, want 90", strategy.Percentage)
		}

		refactor, ok := cfg.Coverage.ByType["coding.refactor"]
		if !ok {
			t.Fatal("expected coding.refactor strategy")
		}
		if refactor.Type != "maintain" {
			t.Errorf("coding.refactor type = %q, want maintain", refactor.Type)
		}
		if refactor.Percentage != nil {
			t.Errorf("coding.refactor percentage = %v, want nil", refactor.Percentage)
		}
	})

	t.Run("coverage absent is nil", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Coverage != nil {
			t.Error("expected Coverage nil when not configured")
		}
	})

	t.Run("coverage with unknown fields silently ignored", func(t *testing.T) {
		dir := setupConfig(t, `coverage:
  coding.feature:
    type: percentage
    percentage: 80
    unknown-extra: value
`)
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Coverage == nil {
			t.Fatal("expected Coverage non-nil")
		}
		strategy := cfg.Coverage.ByType["coding.feature"]
		if strategy.Type != "percentage" {
			t.Errorf("type = %q, want percentage", strategy.Type)
		}
	})
}

func TestCoverageConfigDefaults(t *testing.T) {
	defaults := CoverageConfigDefaults()
	if len(defaults.ByType) == 0 {
		t.Error("expected non-empty default strategies")
	}

	tests := []struct {
		taskType     string
		wantType     string
		wantPct      int
		wantMaintain bool
	}{
		{"coding.feature", "percentage", 80, false},
		{"coding.enhancement", "percentage", 60, false},
		{"coding.fix", "percentage", 60, false},
		{"coding.refactor", "maintain", 0, true},
		{"coding.cleanup", "maintain", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.taskType, func(t *testing.T) {
			s, ok := defaults.ByType[tt.taskType]
			if !ok {
				t.Fatalf("no default for %q", tt.taskType)
			}
			if s.Type != tt.wantType {
				t.Errorf("type = %q, want %q", s.Type, tt.wantType)
			}
			if tt.wantMaintain {
				if s.Percentage != nil {
					t.Errorf("percentage = %v, want nil for maintain", s.Percentage)
				}
			} else {
				if s.Percentage == nil || *s.Percentage != tt.wantPct {
					t.Errorf("percentage = %v, want %d", s.Percentage, tt.wantPct)
				}
			}
		})
	}
}

func TestCoverageConfigDefaults_UnknownType(t *testing.T) {
	defaults := CoverageConfigDefaults()
	_, ok := defaults.ByType["coding.unknown"]
	if ok {
		t.Error("expected no default for unknown type")
	}
}

func TestCoverageConfigDefaults_Immutable(t *testing.T) {
	d1 := CoverageConfigDefaults()
	d2 := CoverageConfigDefaults()
	// Mutating one should not affect the other
	delete(d1.ByType, "coding.feature")
	if _, ok := d2.ByType["coding.feature"]; !ok {
		t.Error("mutating one default affected the other")
	}
}

func TestReadCoverageConfig(t *testing.T) {
	t.Run("missing config returns defaults", func(t *testing.T) {
		dir := t.TempDir()
		coverage, err := ReadCoverageConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defaults := CoverageConfigDefaults()
		if len(coverage.ByType) != len(defaults.ByType) {
			t.Errorf("coverage types = %d, want %d", len(coverage.ByType), len(defaults.ByType))
		}
		for k, v := range defaults.ByType {
			got, ok := coverage.ByType[k]
			if !ok {
				t.Errorf("missing default for %q", k)
				continue
			}
			if got.Type != v.Type {
				t.Errorf("type for %q = %q, want %q", k, got.Type, v.Type)
			}
		}
	})

	t.Run("partial config merges with defaults", func(t *testing.T) {
		dir := setupConfig(t, `coverage:
  coding.feature:
    type: percentage
    percentage: 95
`)
		coverage, err := ReadCoverageConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// coding.feature should be overridden
		s := coverage.ByType["coding.feature"]
		if s.Percentage == nil || *s.Percentage != 95 {
			t.Errorf("coding.feature percentage = %v, want 95", s.Percentage)
		}
		// Other defaults should still be present
		if _, ok := coverage.ByType["coding.refactor"]; !ok {
			t.Error("expected coding.refactor default to be present")
		}
	})
}

func TestGetConfigValue_CoverageKeys(t *testing.T) {
	t.Run("coverage.coding.feature returns percentage", func(t *testing.T) {
		dir := setupConfig(t, `coverage:
  coding.feature:
    type: percentage
    percentage: 80
  coding.refactor:
    type: maintain
`)
		val, err := GetConfigValue(dir, "coverage.coding.feature")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "80" {
			t.Errorf("expected '80', got %q", val)
		}
	})

	t.Run("coverage.coding.refactor returns maintain", func(t *testing.T) {
		dir := setupConfig(t, `coverage:
  coding.refactor:
    type: maintain
`)
		val, err := GetConfigValue(dir, "coverage.coding.refactor")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "maintain" {
			t.Errorf("expected 'maintain', got %q", val)
		}
	})

	t.Run("coverage.unknown-type returns error", func(t *testing.T) {
		dir := setupConfig(t, `coverage:
  coding.feature:
    type: percentage
    percentage: 80
`)
		_, err := GetConfigValue(dir, "coverage.unknown.type")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})

	t.Run("coverage key with no config returns default", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "coverage.coding.feature")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "80" {
			t.Errorf("expected '80' (default), got %q", val)
		}
	})
}

func TestWriteConfigAutoBlock(t *testing.T) {
	t.Run("write and read auto block", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Auto: &AutoConfig{
				Test:             ModeToggle{Quick: false, Full: true},
				ConsolidateSpecs: ModeToggle{Quick: true, Full: true},
				GitPush:          true,
			},
		}
		if err := writeConfig(dir, cfg); err != nil {
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

func TestSetConfigValue(t *testing.T) {
	t.Run("auto.gitPush set to true", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "auto.gitPush", "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("auto.cleanCode.quick set to false", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "auto.cleanCode.quick", "false"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.cleanCode.quick")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false', got %q", val)
		}
	})

	t.Run("auto.test.full set to true", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "auto.test.full", "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.test.full")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("worktree.source-branch set to develop", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "worktree.source-branch", "develop"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "worktree.source-branch")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "develop" {
			t.Errorf("expected 'develop', got %q", val)
		}
	})

	t.Run("test-framework set to pytest", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "test-framework", "pytest"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "test-framework")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "pytest" {
			t.Errorf("expected 'pytest', got %q", val)
		}
	})

	t.Run("unknown key returns meaningful error", func(t *testing.T) {
		dir := t.TempDir()
		err := SetConfigValue(dir, "nonexistent", "value")
		if err == nil {
			t.Fatal("expected error for unknown key")
		}
		if !strings.Contains(err.Error(), "unknown config key") {
			t.Errorf("expected 'unknown config key' in error, got %v", err)
		}
	})

	t.Run("invalid bool value returns error", func(t *testing.T) {
		dir := t.TempDir()
		err := SetConfigValue(dir, "auto.gitPush", "notabool")
		if err == nil {
			t.Fatal("expected error for invalid bool")
		}
	})

	t.Run("set and verify persistence with existing config", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: false\n")
		if err := SetConfigValue(dir, "auto.gitPush", "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("auto.cleanCode set both quick and full", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "auto.cleanCode", "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.cleanCode")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "quick:true full:true" {
			t.Errorf("expected 'quick:true full:true', got %q", val)
		}
	})

	t.Run("coverage set", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "coverage.coding.feature", "90"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "coverage.coding.feature")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "90" {
			t.Errorf("expected '90', got %q", val)
		}
	})
}

func TestGetConfigValue_KnowledgeSave(t *testing.T) {
	t.Run("auto.knowledgeSave default returns quick:true full:false", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.knowledgeSave")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "quick:true full:false" {
			t.Errorf("expected 'quick:true full:false', got %q", val)
		}
	})

	t.Run("auto.knowledgeSave explicit", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  knowledgeSave:\n    quick: false\n    full: true\n")
		val, err := GetConfigValue(dir, "auto.knowledgeSave")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "quick:false full:true" {
			t.Errorf("expected 'quick:false full:true', got %q", val)
		}
	})

	t.Run("auto.knowledgeSave partial applies default", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  knowledgeSave:\n    full: true\n")
		val, err := GetConfigValue(dir, "auto.knowledgeSave")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "quick:true full:true" {
			t.Errorf("expected 'quick:true full:true' (quick defaulted), got %q", val)
		}
	})
}

func TestGetConfigValue_RunTasks(t *testing.T) {
	t.Run("auto.runTasks default returns quick:true full:false", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.runTasks")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "quick:true full:false" {
			t.Errorf("expected 'quick:true full:false', got %q", val)
		}
	})

	t.Run("auto.runTasks explicit", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  runTasks:\n    quick: false\n    full: true\n")
		val, err := GetConfigValue(dir, "auto.runTasks")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "quick:false full:true" {
			t.Errorf("expected 'quick:false full:true', got %q", val)
		}
	})

	t.Run("auto.runTasks partial applies default", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  runTasks:\n    full: true\n")
		val, err := GetConfigValue(dir, "auto.runTasks")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "quick:true full:true" {
			t.Errorf("expected 'quick:true full:true' (quick defaulted), got %q", val)
		}
	})

	t.Run("auto.runTasks.quick returns bool", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  runTasks:\n    quick: true\n    full: false\n")
		val, err := GetConfigValue(dir, "auto.runTasks.quick")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("auto.runTasks.full returns bool", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  runTasks:\n    quick: false\n    full: true\n")
		val, err := GetConfigValue(dir, "auto.runTasks.full")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("auto.runTasks.quick returns default when absent", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.runTasks.quick")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true' (default), got %q", val)
		}
	})
}

func TestReadConfig_OldKeyMigration(t *testing.T) {
	t.Run("old key e2eTest maps to Test field", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  e2eTest:\n    quick: true\n    full: false\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Auto == nil {
			t.Fatal("expected Auto non-nil")
		}
		if !cfg.Auto.Test.Quick {
			t.Error("Test.Quick should be true (mapped from old e2eTest key)")
		}
		if cfg.Auto.Test.Full {
			t.Error("Test.Full should be false (mapped from old e2eTest key)")
		}
	})

	t.Run("only new key test no migration hint", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  test:\n    quick: true\n    full: false\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !cfg.Auto.Test.Quick {
			t.Error("Test.Quick should be true")
		}
		if cfg.Auto.Test.Full {
			t.Error("Test.Full should be false")
		}
	})

	t.Run("old key e2eTest outputs migration hint to stderr", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  e2eTest:\n    quick: false\n")
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		cfg, err := ReadConfig(dir)
		_ = w.Close()
		os.Stderr = old
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		stderr := buf.String()
		if !strings.Contains(stderr, "config key 'auto.e2eTest' is renamed to 'auto.test' in v3.0.0") {
			t.Errorf("expected migration hint in stderr, got: %s", stderr)
		}
		if cfg.Auto.Test.Quick {
			t.Error("Test.Quick should be false (mapped from old key)")
		}
		if !cfg.Auto.Test.Full {
			t.Error("Test.Full should be true (default applied)")
		}
		_ = cfg
	})
}
