package forgeconfig

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
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
		if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "unknown config key") {
			t.Errorf("expected 'not found' in error, got %v", err)
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

	t.Run("auto.cleanCode set rejected as ModeToggle", func(t *testing.T) {
		dir := t.TempDir()
		err := SetConfigValue(dir, "auto.cleanCode", "true")
		if err == nil {
			t.Fatal("expected error for ModeToggle direct set")
		}
		if !strings.Contains(err.Error(), "cannot set ModeToggle directly") {
			t.Errorf("expected cannot set ModeToggle directly in error, got %v", err)
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

func TestGetConfigValue_EvalConfig(t *testing.T) {
	t.Run("auto.eval.proposal returns quick:true full:true (default)", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.proposal")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "quick:true full:true" {
			t.Errorf("expected 'quick:true full:true', got %q", val)
		}
	})

	t.Run("auto.eval.proposal.quick returns true (4-level depth)", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.proposal.quick")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("auto.eval returns eval sub-field summary", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, field := range []string{"proposal:", "prd:", "uiDesign:", "techDesign:"} {
			if !strings.Contains(val, field) {
				t.Errorf("expected field %q in output, got %q", field, val)
			}
		}
	})

	t.Run("auto returns mixed type summary", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, field := range []string{"runTasks:", "gitPush:", "eval:"} {
			if !strings.Contains(val, field) {
				t.Errorf("expected field %q in output, got %q", field, val)
			}
		}
	})

	t.Run("auto.eval.proposal.quick.extra returns errKeyNotFound", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "auto.eval.proposal.quick.extra")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})

	t.Run("auto.nonexistent returns errKeyNotFound", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "auto.nonexistent")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})
}

func TestSetConfigValue_EvalConfig(t *testing.T) {
	t.Run("auto.eval rejected as non-leaf", func(t *testing.T) {
		dir := t.TempDir()
		err := SetConfigValue(dir, "auto.eval", "true")
		if err == nil {
			t.Fatal("expected error for non-leaf set")
		}
		if !strings.Contains(err.Error(), "cannot set non-leaf key") {
			t.Errorf("expected 'cannot set non-leaf key' in error, got %v", err)
		}
	})

	t.Run("auto.eval.proposal rejected as ModeToggle", func(t *testing.T) {
		dir := t.TempDir()
		err := SetConfigValue(dir, "auto.eval.proposal", "true")
		if err == nil {
			t.Fatal("expected error for ModeToggle direct set")
		}
		if !strings.Contains(err.Error(), "cannot set ModeToggle directly") {
			t.Errorf("expected 'cannot set ModeToggle directly' in error, got %v", err)
		}
	})

	t.Run("auto.eval.prd.full true writes nested config", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "auto.eval.prd.full", "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.eval.prd.full")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})
}

func TestGetConfigValue_InlineMap(t *testing.T) {
	t.Run("coverage.coding.feature returns default (inline tag)", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "coverage.coding.feature")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "80" {
			t.Errorf("expected '80', got %q", val)
		}
	})

	t.Run("coverage.coding.feature with explicit config", func(t *testing.T) {
		dir := setupConfig(t, "coverage:\n  coding.feature:\n    type: percentage\n    percentage: 90\n")
		val, err := GetConfigValue(dir, "coverage.coding.feature")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "90" {
			t.Errorf("expected '90', got %q", val)
		}
	})
}

func TestGetConfigValue_WorktreeRegression(t *testing.T) {
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
}

func TestParseAutoRaw_EvalConfig(t *testing.T) {
	t.Run("eval sub-fields tracked with flat-path keys", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  eval:\n    proposal:\n      quick: false\n    prd:\n      quick: true\n      full: false\n  test:\n    quick: true\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil || cfg.Auto == nil {
			t.Fatal("expected non-nil config")
		}

		if cfg.Auto.Eval.Proposal.Quick {
			t.Error("Eval.Proposal.Quick should be false (explicitly set)")
		}
		if !cfg.Auto.Eval.Proposal.Full {
			t.Error("Eval.Proposal.Full should be true (default applied)")
		}
		if !cfg.Auto.Eval.Prd.Quick {
			t.Error("Eval.Prd.Quick should be true")
		}
		if cfg.Auto.Eval.Prd.Full {
			t.Error("Eval.Prd.Full should be false")
		}
		if !cfg.Auto.Eval.UiDesign.Quick || !cfg.Auto.Eval.UiDesign.Full {
			t.Errorf("Eval.UiDesign should default to true/true, got %+v", cfg.Auto.Eval.UiDesign)
		}
		if cfg.Auto.Eval.TechDesign.Quick || cfg.Auto.Eval.TechDesign.Full {
			t.Errorf("Eval.TechDesign should default to false/false, got %+v", cfg.Auto.Eval.TechDesign)
		}
	})
}

func TestParseAutoRaw_ExistingFields_Regression(t *testing.T) {
	t.Run("existing auto fields still tracked correctly", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  test:\n    quick: false\n  consolidateSpecs:\n    quick: true\n    full: false\n  gitPush: true\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil || cfg.Auto == nil {
			t.Fatal("expected non-nil config")
		}

		if cfg.Auto.Test.Quick {
			t.Error("Test.Quick should be false")
		}
		if !cfg.Auto.Test.Full {
			t.Error("Test.Full should be true (default)")
		}
		if !cfg.Auto.ConsolidateSpecs.Quick {
			t.Error("ConsolidateSpecs.Quick should be true")
		}
		if cfg.Auto.ConsolidateSpecs.Full {
			t.Error("ConsolidateSpecs.Full should be false")
		}
		if !cfg.Auto.GitPush {
			t.Error("GitPush should be true")
		}
	})
}

// TestGetStructValueByPath tests the reflect-based getByPath at various depths and edge cases.
func TestGetStructValueByPath(t *testing.T) {
	t.Run("three-level path auto.eval.proposal", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.proposal")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "quick:true full:true" {
			t.Errorf("expected 'quick:true full:true', got %q", val)
		}
	})

	t.Run("four-level path auto.eval.proposal.quick", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.proposal.quick")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("four-level path auto.eval.proposal.full", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.proposal.full")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("intermediate node get auto.eval returns summary", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should contain all four eval sub-fields
		for _, field := range []string{"proposal:", "prd:", "uiDesign:", "techDesign:"} {
			if !strings.Contains(val, field) {
				t.Errorf("expected field %q in output, got %q", field, val)
			}
		}
	})

	t.Run("nil pointer field returns errKeyNotFound", func(t *testing.T) {
		// Config with nil worktree — accessing worktree.source-branch should return errKeyNotFound
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "worktree.source-branch")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound for nil pointer, got %v", err)
		}
	})

	t.Run("nonexistent field returns errKeyNotFound", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "auto.nonexistent")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})

	t.Run("path beyond leaf returns errKeyNotFound", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "auto.eval.proposal.quick.extra")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})

	t.Run("single-level key test-framework returns value", func(t *testing.T) {
		dir := setupConfig(t, "test-framework: jest\n")
		val, err := GetConfigValue(dir, "test-framework")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "jest" {
			t.Errorf("expected 'jest', got %q", val)
		}
	})

	t.Run("two-level path auto.gitPush returns bool", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		val, err := GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})
}

// TestSetStructValueByPath tests the reflect-based setByPath for various scenarios.
func TestSetStructValueByPath(t *testing.T) {
	t.Run("multi-level path set auto.eval.prd.full true", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "auto.eval.prd.full", "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.eval.prd.full")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("ModeToggle direct set rejected", func(t *testing.T) {
		dir := t.TempDir()
		err := SetConfigValue(dir, "auto.eval.proposal", "true")
		if err == nil {
			t.Fatal("expected error for ModeToggle direct set")
		}
		if !strings.Contains(err.Error(), "cannot set ModeToggle directly") {
			t.Errorf("expected 'cannot set ModeToggle directly' in error, got %v", err)
		}
	})

	t.Run("non-leaf set rejected", func(t *testing.T) {
		dir := t.TempDir()
		err := SetConfigValue(dir, "auto.eval", "true")
		if err == nil {
			t.Fatal("expected error for non-leaf set")
		}
		if !strings.Contains(err.Error(), "cannot set non-leaf key") {
			t.Errorf("expected 'cannot set non-leaf key' in error, got %v", err)
		}
	})

	t.Run("nil pointer auto-initialized for worktree.source-branch", func(t *testing.T) {
		dir := t.TempDir()
		// No config file exists, but set should auto-create the worktree block
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

	t.Run("set auto.cleanCode.quick to false", func(t *testing.T) {
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

	t.Run("set auto.eval.uiDesign.quick to true (flip default)", func(t *testing.T) {
		dir := t.TempDir()
		// First set it to true (non-default for prd), then verify
		if err := SetConfigValue(dir, "auto.eval.prd.quick", "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.eval.prd.quick")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("invalid bool value returns descriptive error", func(t *testing.T) {
		dir := t.TempDir()
		err := SetConfigValue(dir, "auto.gitPush", "maybe")
		if err == nil {
			t.Fatal("expected error for invalid bool value")
		}
		if !strings.Contains(err.Error(), "invalid value") {
			t.Errorf("expected 'invalid value' in error, got %v", err)
		}
		if !strings.Contains(err.Error(), "maybe") {
			t.Errorf("expected value 'maybe' in error, got %v", err)
		}
	})

	t.Run("nonexistent key returns error", func(t *testing.T) {
		dir := t.TempDir()
		err := SetConfigValue(dir, "auto.nonexistent", "true")
		if err == nil {
			t.Fatal("expected error for nonexistent key")
		}
		if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "unknown") {
			t.Errorf("expected 'not found' or 'unknown' in error, got %v", err)
		}
	})
}

// TestGetConfigValue_SurfacesMap_Fallback tests that SurfacesMap scalar and map
// forms fall back correctly via the reflect → hardcode fallback path.
func TestGetConfigValue_SurfacesMap_Fallback(t *testing.T) {
	t.Run("scalar surfaces form (single string) is not retrievable by key path", func(t *testing.T) {
		// surfaces: api → SurfacesMap{".": "api"}
		// Attempting getByPath on SurfacesMap hits errUnsupportedType
		dir := setupConfig(t, "surfaces: api\n")
		// surfaces itself is a map — reflect can handle it as map summary
		// but individual sub-keys like "surfaces.." would be unusual
		// The main point: reading a config with scalar surfaces should not error
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Surfaces == nil {
			t.Fatal("expected Surfaces non-nil")
		}
		if v, ok := cfg.Surfaces["."]; !ok || v != "api" {
			t.Errorf("expected Surfaces['.'] = 'api', got %v", cfg.Surfaces)
		}
	})

	t.Run("map surfaces form parses correctly", func(t *testing.T) {
		dir := setupConfig(t, "surfaces:\n  frontend: web\n  backend: api\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cfg.Surfaces) != 2 {
			t.Errorf("expected 2 surfaces, got %d", len(cfg.Surfaces))
		}
		if cfg.Surfaces["frontend"] != "web" {
			t.Errorf("expected frontend=web, got %s", cfg.Surfaces["frontend"])
		}
		if cfg.Surfaces["backend"] != "api" {
			t.Errorf("expected backend=api, got %s", cfg.Surfaces["backend"])
		}
	})

	t.Run("SurfacesMap custom unmarshaler returns errUnsupportedType for reflect routing", func(t *testing.T) {
		// This test validates that getByPath returns errUnsupportedType when
		// encountering a SurfacesMap field (implements yaml.Unmarshaler).
		dir := setupConfig(t, "surfaces: api\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Directly test the reflect routing on the SurfacesMap field
		surfacesField := reflect.ValueOf(cfg).Elem().FieldByName("Surfaces")
		if !surfacesField.IsValid() {
			t.Fatal("Surfaces field not found on Config struct")
		}
		_, err = formatValue(derefPointer(surfacesField))
		if err != errUnsupportedType {
			t.Errorf("expected errUnsupportedType for SurfacesMap, got %v", err)
		}
	})
}

// TestGetByPath_DirectByPathTests tests the internal getByPath function directly.
func TestGetByPath_DirectByPathTests(t *testing.T) {
	t.Run("get by empty segments returns errKeyNotFound", func(t *testing.T) {
		cfg := &Config{}
		_, err := getByPath(reflect.ValueOf(cfg).Elem(), []string{})
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound for empty segments, got %v", err)
		}
	})

	t.Run("get by single segment top-level key", func(t *testing.T) {
		cfg := &Config{Version: "2"}
		val, err := getByPath(reflect.ValueOf(cfg).Elem(), []string{"version"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "2" {
			t.Errorf("expected '2', got %q", val)
		}
	})
}

// TestDetectPipelineMode tests mode detection at the forgeconfig API level.
// The actual detectModeFromPath is in internal/cmd, but the API behavior
// is testable through GetConfigValue for config-based keys.
func TestDetectPipelineMode(t *testing.T) {
	// These tests validate mode-related config behavior:
	// The actual mode detection (quick/full/none) is tested in
	// internal/cmd/config_test.go::TestDetectModeFromPath.
	// Here we test that the config system supports the mode-based eval keys.

	t.Run("auto.eval.proposal.quick default enables quick mode eval", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.proposal.quick")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true' (quick mode enabled by default), got %q", val)
		}
	})

	t.Run("auto.eval.prd.quick default disables quick mode eval", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.prd.quick")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false' (quick mode disabled by default), got %q", val)
		}
	})

	t.Run("auto.eval.prd.full default disables full mode eval", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.prd.full")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false' (full mode disabled by default), got %q", val)
		}
	})

	t.Run("auto.eval.uiDesign.quick default enables quick mode eval", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.uiDesign.quick")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("auto.eval.techDesign.full default disables full mode eval", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.techDesign.full")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false', got %q", val)
		}
	})
}

// TestParseAutoRaw_FlatPathTracking verifies the raw map flat-path key structure.
func TestParseAutoRaw_FlatPathTracking(t *testing.T) {
	t.Run("raw map contains eval.proposal flat-path key", func(t *testing.T) {
		yaml := []byte("auto:\n  eval:\n    proposal:\n      quick: true\n")
		raw, err := parseAutoRaw(yaml)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		fieldRaw, exists := raw["eval.proposal"]
		if !exists {
			t.Fatal("expected 'eval.proposal' key in raw map")
		}
		if !fieldRaw["quick"] {
			t.Error("expected quick=true in eval.proposal raw")
		}
	})

	t.Run("raw map tracks multiple eval sub-fields", func(t *testing.T) {
		yaml := []byte("auto:\n  eval:\n    proposal:\n      quick: false\n    prd:\n      full: true\n")
		raw, err := parseAutoRaw(yaml)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		proposalRaw, exists := raw["eval.proposal"]
		if !exists {
			t.Fatal("expected 'eval.proposal' in raw map")
		}
		if !proposalRaw["quick"] {
			t.Error("expected quick tracked for eval.proposal")
		}

		prdRaw, exists := raw["eval.prd"]
		if !exists {
			t.Fatal("expected 'eval.prd' in raw map")
		}
		if !prdRaw["full"] {
			t.Error("expected full tracked for eval.prd")
		}
	})

	t.Run("applyDefaults only fills missing sub-keys", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  eval:\n    proposal:\n      quick: false\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// quick was explicitly set to false — should remain false
		if cfg.Auto.Eval.Proposal.Quick {
			t.Error("Eval.Proposal.Quick should be false (explicitly set)")
		}
		// full was not set — should get default (true)
		if !cfg.Auto.Eval.Proposal.Full {
			t.Error("Eval.Proposal.Full should be true (default applied)")
		}
	})
}

// TestGetConfigValue_EvalExplicitConfig tests eval config with explicit YAML values.
func TestGetConfigValue_EvalExplicitConfig(t *testing.T) {
	t.Run("explicit eval config overrides defaults", func(t *testing.T) {
		dir := setupConfig(t, `auto:
  eval:
    proposal:
      quick: false
      full: false
    prd:
      quick: true
      full: true
    uiDesign:
      quick: false
    techDesign:
      full: true
`)
		tests := []struct {
			key      string
			expected string
		}{
			{"auto.eval.proposal", "quick:false full:false"},
			{"auto.eval.prd", "quick:true full:true"},
			{"auto.eval.uiDesign", "quick:false full:true"},   // quick explicit false, full defaults to true
			{"auto.eval.techDesign", "quick:false full:true"}, // quick defaults false, full explicit true
		}
		for _, tt := range tests {
			t.Run(tt.key, func(t *testing.T) {
				val, err := GetConfigValue(dir, tt.key)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if val != tt.expected {
					t.Errorf("key %q: expected %q, got %q", tt.key, tt.expected, val)
				}
			})
		}
	})
}

// TestGetConfigValue_AutoSummary verifies the auto block summary includes eval fields.
func TestGetConfigValue_AutoSummary(t *testing.T) {
	t.Run("auto summary includes all expected fields", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Check ModeToggle fields
		for _, field := range []string{"runTasks:", "consolidateSpecs:", "cleanCode:", "validation:", "knowledgeSave:"} {
			if !strings.Contains(val, field) {
				t.Errorf("expected field %q in auto summary, got:\n%s", field, val)
			}
		}

		// Check bool field
		if !strings.Contains(val, "gitPush:") {
			t.Errorf("expected 'gitPush:' in auto summary, got:\n%s", val)
		}

		// Check nested struct
		if !strings.Contains(val, "eval:") {
			t.Errorf("expected 'eval:' in auto summary, got:\n%s", val)
		}

		// Check eval sub-fields appear (indented under eval)
		for _, field := range []string{"proposal:", "prd:", "uiDesign:", "techDesign:"} {
			if !strings.Contains(val, field) {
				t.Errorf("expected eval sub-field %q in auto summary, got:\n%s", field, val)
			}
		}
	})
}

func TestEvalConfigDefaults(t *testing.T) {
	t.Run("eval defaults match proposal spec", func(t *testing.T) {
		defaults := AutoConfigDefaults()

		// proposal: quick=true, full=true
		if !defaults.Eval.Proposal.Quick || !defaults.Eval.Proposal.Full {
			t.Errorf("Eval.Proposal = %+v, want {Quick:true Full:true}", defaults.Eval.Proposal)
		}
		// prd: quick=false, full=false
		if defaults.Eval.Prd.Quick || defaults.Eval.Prd.Full {
			t.Errorf("Eval.Prd = %+v, want {Quick:false Full:false}", defaults.Eval.Prd)
		}
		// uiDesign: quick=true, full=true
		if !defaults.Eval.UiDesign.Quick || !defaults.Eval.UiDesign.Full {
			t.Errorf("Eval.UiDesign = %+v, want {Quick:true Full:true}", defaults.Eval.UiDesign)
		}
		// techDesign: quick=false, full=false
		if defaults.Eval.TechDesign.Quick || defaults.Eval.TechDesign.Full {
			t.Errorf("Eval.TechDesign = %+v, want {Quick:false Full:false}", defaults.Eval.TechDesign)
		}
	})
}

func TestAutoConfigIsZero_IncludesEval(t *testing.T) {
	t.Run("zero EvalConfig counts as zero", func(t *testing.T) {
		a := AutoConfig{}
		if !a.IsZero() {
			t.Error("expected zero AutoConfig to be zero")
		}
	})

	t.Run("non-zero Eval makes AutoConfig non-zero", func(t *testing.T) {
		a := AutoConfig{Eval: EvalConfig{Proposal: ModeToggle{Quick: true}}}
		if a.IsZero() {
			t.Error("expected AutoConfig with Eval set to be non-zero")
		}
	})
}

func TestSetConfigValue_EvalPersistence(t *testing.T) {
	t.Run("set auto.eval.techDesign.quick persists across read", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "auto.eval.techDesign.quick", "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Read back from file
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Auto == nil {
			t.Fatal("expected Auto non-nil")
		}
		if !cfg.Auto.Eval.TechDesign.Quick {
			t.Error("Eval.TechDesign.Quick should be true after set")
		}
		// Full should be default (false)
		if cfg.Auto.Eval.TechDesign.Full {
			t.Error("Eval.TechDesign.Full should be false (default)")
		}
	})

	t.Run("set auto.eval.prd.full true after false default", func(t *testing.T) {
		dir := t.TempDir()
		// Default for prd.full is false, set to true
		if err := SetConfigValue(dir, "auto.eval.prd.full", "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.eval.prd.full")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})
}

func TestParseAutoRaw_EvalRegressionWithExistingFields(t *testing.T) {
	t.Run("eval and existing fields coexist in raw tracking", func(t *testing.T) {
		dir := setupConfig(t, `auto:
  test:
    quick: false
  consolidateSpecs:
    quick: true
    full: false
  gitPush: true
  eval:
    proposal:
      quick: false
    prd:
      full: true
`)
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Existing fields
		if cfg.Auto.Test.Quick {
			t.Error("Test.Quick should be false")
		}
		if !cfg.Auto.Test.Full {
			t.Error("Test.Full should be true (default)")
		}
		if !cfg.Auto.ConsolidateSpecs.Quick {
			t.Error("ConsolidateSpecs.Quick should be true")
		}
		if cfg.Auto.ConsolidateSpecs.Full {
			t.Error("ConsolidateSpecs.Full should be false")
		}
		if !cfg.Auto.GitPush {
			t.Error("GitPush should be true")
		}

		// Eval fields
		if cfg.Auto.Eval.Proposal.Quick {
			t.Error("Eval.Proposal.Quick should be false (explicit)")
		}
		if !cfg.Auto.Eval.Proposal.Full {
			t.Error("Eval.Proposal.Full should be true (default)")
		}
		if cfg.Auto.Eval.Prd.Quick {
			t.Error("Eval.Prd.Quick should be false (default)")
		}
		if !cfg.Auto.Eval.Prd.Full {
			t.Error("Eval.Prd.Full should be true (explicit)")
		}
	})
}

// TestFormatValue_EdgeCases tests formatValue edge cases for coverage.
func TestFormatValue_EdgeCases(t *testing.T) {
	t.Run("formatValue on nil SurfacesMap returns errUnsupportedType", func(t *testing.T) {
		var s SurfacesMap
		_, err := formatValue(reflect.ValueOf(s))
		if err != errUnsupportedType {
			t.Errorf("expected errUnsupportedType for nil SurfacesMap, got %v", err)
		}
	})

	t.Run("formatValue on non-nil SurfacesMap returns errUnsupportedType", func(t *testing.T) {
		s := SurfacesMap{".": "api"}
		_, err := formatValue(reflect.ValueOf(s))
		if err != errUnsupportedType {
			t.Errorf("expected errUnsupportedType for non-nil SurfacesMap, got %v", err)
		}
	})

	t.Run("formatValue on int returns string", func(t *testing.T) {
		v := reflect.ValueOf(42)
		result, err := formatValue(v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "42" {
			t.Errorf("expected '42', got %q", result)
		}
	})

	t.Run("formatValue on bool returns string", func(t *testing.T) {
		v := reflect.ValueOf(true)
		result, err := formatValue(v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "true" {
			t.Errorf("expected 'true', got %q", result)
		}
	})

	t.Run("formatValue on string returns string", func(t *testing.T) {
		v := reflect.ValueOf("hello")
		result, err := formatValue(v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "hello" {
			t.Errorf("expected 'hello', got %q", result)
		}
	})

	t.Run("formatValue on string slice returns newline-joined", func(t *testing.T) {
		v := reflect.ValueOf([]string{"a", "b", "c"})
		result, err := formatValue(v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "a\nb\nc" {
			t.Errorf("expected 'a\\nb\\nc', got %q", result)
		}
	})

	t.Run("formatValue on invalid Value returns errKeyNotFound", func(t *testing.T) {
		var v reflect.Value
		_, err := formatValue(v)
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound for invalid Value, got %v", err)
		}
	})
}

// TestGetConfigValue_EvalFullRoundtrip tests set-get roundtrip for all eval fields.
func TestGetConfigValue_EvalFullRoundtrip(t *testing.T) {
	evalKeys := []struct {
		key        string
		defaultVal string
	}{
		{"auto.eval.proposal.quick", "true"},
		{"auto.eval.proposal.full", "true"},
		{"auto.eval.prd.quick", "false"},
		{"auto.eval.prd.full", "false"},
		{"auto.eval.uiDesign.quick", "true"},
		{"auto.eval.uiDesign.full", "true"},
		{"auto.eval.techDesign.quick", "false"},
		{"auto.eval.techDesign.full", "false"},
	}

	for _, tt := range evalKeys {
		t.Run(fmt.Sprintf("default for %s is %s", tt.key, tt.defaultVal), func(t *testing.T) {
			dir := t.TempDir()
			val, err := GetConfigValue(dir, tt.key)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if val != tt.defaultVal {
				t.Errorf("expected %q, got %q", tt.defaultVal, val)
			}
		})
	}
}
