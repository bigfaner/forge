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
		dir := setupConfig(t, "worktree:\n  source-branch: develop\n  includes:\n    - .env\n    - .env.local\n")
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
		if len(cfg.Worktree.Includes) != 2 || cfg.Worktree.Includes[0] != ".env" || cfg.Worktree.Includes[1] != ".env.local" {
			t.Errorf("expected [.env .env.local], got %v", cfg.Worktree.Includes)
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

	t.Run("worktree.includes returns joined list", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  includes:\n    - .env\n    - .env.local\n")
		val, err := GetConfigValue(dir, "worktree.includes")
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

	t.Run("worktree.includes absent returns error", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		_, err := GetConfigValue(dir, "worktree.includes")
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
		dir := setupConfig(t, "worktree:\n  includes:\n    - .env\n")
		_, err := GetConfigValue(dir, "worktree.source-branch")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound, got %v", err)
		}
	})

	t.Run("worktree present but includes empty returns error", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  source-branch: develop\n")
		_, err := GetConfigValue(dir, "worktree.includes")
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
	// Eval defaults (bool)
	if !defaults.Eval.Proposal {
		t.Errorf("Eval.Proposal = %v, want true", defaults.Eval.Proposal)
	}
	if defaults.Eval.Prd {
		t.Errorf("Eval.Prd = %v, want false", defaults.Eval.Prd)
	}
	if !defaults.Eval.UiDesign {
		t.Errorf("Eval.UiDesign = %v, want true", defaults.Eval.UiDesign)
	}
	if defaults.Eval.TechDesign {
		t.Errorf("Eval.TechDesign = %v, want false", defaults.Eval.TechDesign)
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

	t.Run("only Eval.Proposal set is not zero", func(t *testing.T) {
		a := AutoConfig{Eval: EvalConfig{Proposal: true}}
		if a.IsZero() {
			t.Error("expected AutoConfig with Eval.Proposal set to be non-zero")
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
		if cfg.Worktree.Includes != nil {
			t.Errorf("expected Includes nil, got %v", cfg.Worktree.Includes)
		}
	})

	t.Run("only includes", func(t *testing.T) {
		dir := setupConfig(t, "worktree:\n  includes:\n    - .env\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Worktree.SourceBranch != "" {
			t.Errorf("expected empty source-branch, got %q", cfg.Worktree.SourceBranch)
		}
		if len(cfg.Worktree.Includes) != 1 || cfg.Worktree.Includes[0] != ".env" {
			t.Errorf("expected [.env], got %v", cfg.Worktree.Includes)
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
	t.Run("auto.eval.proposal returns true (default)", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.proposal")
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

	t.Run("auto.eval.proposal.extra returns errKeyNotFound (bool is leaf)", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "auto.eval.proposal.extra")
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

	t.Run("auto.eval.proposal set to true (bool field)", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "auto.eval.proposal", "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.eval.proposal")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("auto.eval.prd set to true writes and persists", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "auto.eval.prd", "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.eval.prd")
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
		dir := setupConfig(t, "auto:\n  eval:\n    proposal: false\n    prd: true\n  test:\n    quick: true\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil || cfg.Auto == nil {
			t.Fatal("expected non-nil config")
		}

		if cfg.Auto.Eval.Proposal {
			t.Error("Eval.Proposal should be false (explicitly set)")
		}
		if !cfg.Auto.Eval.Prd {
			t.Error("Eval.Prd should be true (explicitly set)")
		}
		if !cfg.Auto.Eval.UiDesign {
			t.Errorf("Eval.UiDesign should default to true, got %v", cfg.Auto.Eval.UiDesign)
		}
		if cfg.Auto.Eval.TechDesign {
			t.Errorf("Eval.TechDesign should default to false, got %v", cfg.Auto.Eval.TechDesign)
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
	t.Run("three-level path auto.eval.proposal returns bool", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.proposal")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("three-level path auto.eval.prd returns bool", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.prd")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false', got %q", val)
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
		_, err := GetConfigValue(dir, "auto.eval.proposal.extra")
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
	t.Run("multi-level path set auto.eval.prd true", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "auto.eval.prd", "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.eval.prd")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("bool field direct set succeeds", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "auto.eval.proposal", "false"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.eval.proposal")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false', got %q", val)
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

	t.Run("set auto.eval.prd to true (flip default)", func(t *testing.T) {
		dir := t.TempDir()
		// Default for prd is false, flip to true
		if err := SetConfigValue(dir, "auto.eval.prd", "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.eval.prd")
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
	// Here we test that the config system supports bool-based eval keys.

	t.Run("auto.eval.proposal default is true", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.proposal")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true' (enabled by default), got %q", val)
		}
	})

	t.Run("auto.eval.prd default is false", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.prd")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false' (disabled by default), got %q", val)
		}
	})

	t.Run("auto.eval.uiDesign default is true", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.uiDesign")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("auto.eval.techDesign default is false", func(t *testing.T) {
		dir := t.TempDir()
		val, err := GetConfigValue(dir, "auto.eval.techDesign")
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
		yaml := []byte("auto:\n  eval:\n    proposal: true\n")
		raw, err := parseAutoRaw(yaml)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_, exists := raw["eval.proposal"]
		if !exists {
			t.Fatal("expected 'eval.proposal' key in raw map")
		}
	})

	t.Run("raw map tracks multiple eval sub-fields", func(t *testing.T) {
		yaml := []byte("auto:\n  eval:\n    proposal: false\n    prd: true\n")
		raw, err := parseAutoRaw(yaml)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_, exists := raw["eval.proposal"]
		if !exists {
			t.Fatal("expected 'eval.proposal' in raw map")
		}

		_, exists = raw["eval.prd"]
		if !exists {
			t.Fatal("expected 'eval.prd' in raw map")
		}
	})

	t.Run("applyDefaults preserves explicit values", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  eval:\n    proposal: false\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// proposal was explicitly set to false — should remain false
		if cfg.Auto.Eval.Proposal {
			t.Error("Eval.Proposal should be false (explicitly set)")
		}
		// prd was not set — should get default (false)
		if cfg.Auto.Eval.Prd {
			t.Error("Eval.Prd should be false (default)")
		}
		// uiDesign was not set — should get default (true)
		if !cfg.Auto.Eval.UiDesign {
			t.Error("Eval.UiDesign should be true (default)")
		}
	})

	t.Run("raw map tracks old ModeToggle format for eval", func(t *testing.T) {
		yaml := []byte("auto:\n  eval:\n    proposal:\n      quick: false\n      full: true\n")
		raw, err := parseAutoRaw(yaml)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_, exists := raw["eval.proposal"]
		if !exists {
			t.Fatal("expected 'eval.proposal' in raw map (old ModeToggle format)")
		}
	})
}

// TestGetConfigValue_EvalExplicitConfig tests eval config with explicit YAML values.
func TestGetConfigValue_EvalExplicitConfig(t *testing.T) {
	t.Run("explicit eval config overrides defaults", func(t *testing.T) {
		dir := setupConfig(t, `auto:
  eval:
    proposal: false
    prd: true
    uiDesign: false
    techDesign: true
`)
		tests := []struct {
			key      string
			expected string
		}{
			{"auto.eval.proposal", "false"},
			{"auto.eval.prd", "true"},
			{"auto.eval.uiDesign", "false"},
			{"auto.eval.techDesign", "true"},
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

		// Check eval sub-fields appear (bool format)
		for _, field := range []string{"proposal:", "prd:", "uiDesign:", "techDesign:"} {
			if !strings.Contains(val, field) {
				t.Errorf("expected eval sub-field %q in auto summary, got:\n%s", field, val)
			}
		}
	})
}

func TestEvalConfigDefaults(t *testing.T) {
	t.Run("eval defaults match spec", func(t *testing.T) {
		defaults := AutoConfigDefaults()

		// proposal: true
		if !defaults.Eval.Proposal {
			t.Errorf("Eval.Proposal = %v, want true", defaults.Eval.Proposal)
		}
		// prd: false
		if defaults.Eval.Prd {
			t.Errorf("Eval.Prd = %v, want false", defaults.Eval.Prd)
		}
		// uiDesign: true
		if !defaults.Eval.UiDesign {
			t.Errorf("Eval.UiDesign = %v, want true", defaults.Eval.UiDesign)
		}
		// techDesign: false
		if defaults.Eval.TechDesign {
			t.Errorf("Eval.TechDesign = %v, want false", defaults.Eval.TechDesign)
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
		a := AutoConfig{Eval: EvalConfig{Proposal: true}}
		if a.IsZero() {
			t.Error("expected AutoConfig with Eval set to be non-zero")
		}
	})
}

func TestSetConfigValue_EvalPersistence(t *testing.T) {
	t.Run("set auto.eval.techDesign persists across read", func(t *testing.T) {
		dir := t.TempDir()
		if err := SetConfigValue(dir, "auto.eval.techDesign", "true"); err != nil {
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
		if !cfg.Auto.Eval.TechDesign {
			t.Error("Eval.TechDesign should be true after set")
		}
	})

	t.Run("set auto.eval.prd true after false default", func(t *testing.T) {
		dir := t.TempDir()
		// Default for prd is false, set to true
		if err := SetConfigValue(dir, "auto.eval.prd", "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "auto.eval.prd")
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
    proposal: false
    prd: true
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

		// Eval fields (bool)
		if cfg.Auto.Eval.Proposal {
			t.Error("Eval.Proposal should be false (explicit)")
		}
		if !cfg.Auto.Eval.Prd {
			t.Error("Eval.Prd should be true (explicit)")
		}
		if !cfg.Auto.Eval.UiDesign {
			t.Errorf("Eval.UiDesign should default to true, got %v", cfg.Auto.Eval.UiDesign)
		}
		if cfg.Auto.Eval.TechDesign {
			t.Errorf("Eval.TechDesign should default to false, got %v", cfg.Auto.Eval.TechDesign)
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
		{"auto.eval.proposal", "true"},
		{"auto.eval.prd", "false"},
		{"auto.eval.uiDesign", "true"},
		{"auto.eval.techDesign", "false"},
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

// TestEvalSettingsDefaults tests the default eval settings matching rubric frontmatter.
func TestEvalSettingsDefaults(t *testing.T) {
	t.Run("defaults match rubric frontmatter", func(t *testing.T) {
		defaults := EvalSettingsDefaults()

		tests := []struct {
			name       string
			target     int
			iterations int
		}{
			{"proposal", 900, 3},
			{"prd", 900, 3},
			{"design", 900, 3},
			{"ui", 950, 3},
			{"journey", 850, 3},
			{"contract", 850, 3},
			{"consistency", 900, 3},
		}

		fields := map[string]EvalTypeSettings{
			"proposal":    defaults.Proposal,
			"prd":         defaults.Prd,
			"design":      defaults.Design,
			"ui":          defaults.Ui,
			"journey":     defaults.Journey,
			"contract":    defaults.Contract,
			"consistency": defaults.Consistency,
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				field, ok := fields[tt.name]
				if !ok {
					t.Fatalf("field %q not found in EvalSettings", tt.name)
				}
				if field.Target == nil {
					t.Fatalf("Target is nil for %q", tt.name)
				}
				if *field.Target != tt.target {
					t.Errorf("%s.Target = %d, want %d", tt.name, *field.Target, tt.target)
				}
				if field.Iterations == nil {
					t.Fatalf("Iterations is nil for %q", tt.name)
				}
				if *field.Iterations != tt.iterations {
					t.Errorf("%s.Iterations = %d, want %d", tt.name, *field.Iterations, tt.iterations)
				}
			})
		}
	})

	t.Run("returns fresh instances (immutable)", func(t *testing.T) {
		d1 := EvalSettingsDefaults()
		d2 := EvalSettingsDefaults()
		*d1.Proposal.Target = 0
		if *d2.Proposal.Target != 900 {
			t.Error("mutating one default affected the other")
		}
	})
}

// TestEvalSettingsWriteReadRoundtrip tests that EvalSettings round-trips through YAML.
func TestEvalSettingsWriteReadRoundtrip(t *testing.T) {
	t.Run("write and read eval block", func(t *testing.T) {
		dir := t.TempDir()
		proposalTarget := 900
		proposalIter := 3
		uiTarget := 950
		uiIter := 3
		cfg := &Config{
			Eval: &EvalSettings{
				Proposal: EvalTypeSettings{Target: &proposalTarget, Iterations: &proposalIter},
				Ui:       EvalTypeSettings{Target: &uiTarget, Iterations: &uiIter},
			},
		}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		readback, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if readback.Eval == nil {
			t.Fatal("expected Eval non-nil")
		}
		if readback.Eval.Proposal.Target == nil || *readback.Eval.Proposal.Target != 900 {
			t.Errorf("Proposal.Target = %v, want 900", readback.Eval.Proposal.Target)
		}
		if readback.Eval.Proposal.Iterations == nil || *readback.Eval.Proposal.Iterations != 3 {
			t.Errorf("Proposal.Iterations = %v, want 3", readback.Eval.Proposal.Iterations)
		}
		if readback.Eval.Ui.Target == nil || *readback.Eval.Ui.Target != 950 {
			t.Errorf("Ui.Target = %v, want 950", readback.Eval.Ui.Target)
		}
	})

	t.Run("eval block absent is nil", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Eval != nil {
			t.Error("expected Eval nil when not configured")
		}
	})

	t.Run("full defaults roundtrip", func(t *testing.T) {
		dir := t.TempDir()
		defaults := EvalSettingsDefaults()
		cfg := &Config{Eval: &defaults}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		readback, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if readback.Eval == nil {
			t.Fatal("expected Eval non-nil")
		}
		if readback.Eval.Proposal.Target == nil || *readback.Eval.Proposal.Target != 900 {
			t.Errorf("Proposal.Target = %v, want 900", readback.Eval.Proposal.Target)
		}
		if readback.Eval.Consistency.Target == nil || *readback.Eval.Consistency.Target != 900 {
			t.Errorf("Consistency.Target = %v, want 900", readback.Eval.Consistency.Target)
		}
	})
}

// TestGetConfigValue_EvalSettingsKeys tests config get for eval type settings.
func TestGetConfigValue_EvalSettingsKeys(t *testing.T) {
	t.Run("eval.proposal.target returns value", func(t *testing.T) {
		dir := t.TempDir()
		defaults := EvalSettingsDefaults()
		cfg := &Config{Eval: &defaults}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "eval.proposal.target")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "900" {
			t.Errorf("expected '900', got %q", val)
		}
	})

	t.Run("eval.ui.target returns 950", func(t *testing.T) {
		dir := t.TempDir()
		defaults := EvalSettingsDefaults()
		cfg := &Config{Eval: &defaults}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "eval.ui.target")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "950" {
			t.Errorf("expected '950', got %q", val)
		}
	})

	t.Run("eval.proposal.iterations returns 3", func(t *testing.T) {
		dir := t.TempDir()
		defaults := EvalSettingsDefaults()
		cfg := &Config{Eval: &defaults}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "eval.proposal.iterations")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "3" {
			t.Errorf("expected '3', got %q", val)
		}
	})

	t.Run("eval not configured returns errKeyNotFound", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "eval.proposal.target")
		if err != errKeyNotFound {
			t.Errorf("expected errKeyNotFound for unconfigured eval, got %v", err)
		}
	})

	t.Run("eval returns summary of all types", func(t *testing.T) {
		dir := t.TempDir()
		defaults := EvalSettingsDefaults()
		cfg := &Config{Eval: &defaults}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := GetConfigValue(dir, "eval")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, field := range []string{"proposal:", "prd:", "design:", "ui:", "journey:", "contract:", "consistency:"} {
			if !strings.Contains(val, field) {
				t.Errorf("expected field %q in eval summary, got:\n%s", field, val)
			}
		}
	})
}

// TestEvalConfig_OldModeToggleCompat tests backward compatibility with old ModeToggle format.
func TestEvalConfig_OldModeToggleCompat(t *testing.T) {
	t.Run("old ModeToggle map format reads 'full' sub-key", func(t *testing.T) {
		dir := setupConfig(t, `auto:
  eval:
    proposal:
      quick: false
      full: true
    prd:
      quick: true
      full: false
    uiDesign:
      quick: false
      full: false
    techDesign:
      quick: true
      full: true
`)
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil || cfg.Auto == nil {
			t.Fatal("expected non-nil config")
		}

		// proposal: full=true → true
		if !cfg.Auto.Eval.Proposal {
			t.Errorf("Eval.Proposal should be true (full=true from ModeToggle), got %v", cfg.Auto.Eval.Proposal)
		}
		// prd: full=false → false
		if cfg.Auto.Eval.Prd {
			t.Errorf("Eval.Prd should be false (full=false from ModeToggle), got %v", cfg.Auto.Eval.Prd)
		}
		// uiDesign: full=false → false
		if cfg.Auto.Eval.UiDesign {
			t.Errorf("Eval.UiDesign should be false (full=false from ModeToggle), got %v", cfg.Auto.Eval.UiDesign)
		}
		// techDesign: full=true → true
		if !cfg.Auto.Eval.TechDesign {
			t.Errorf("Eval.TechDesign should be true (full=true from ModeToggle), got %v", cfg.Auto.Eval.TechDesign)
		}
	})

	t.Run("old ModeToggle map format with only quick", func(t *testing.T) {
		dir := setupConfig(t, `auto:
  eval:
    proposal:
      quick: true
`)
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// No 'full' sub-key → defaults to false (map without 'full' key)
		if cfg.Auto.Eval.Proposal {
			t.Errorf("Eval.Proposal should be false (no 'full' sub-key in ModeToggle map), got %v", cfg.Auto.Eval.Proposal)
		}
	})
}

// --- LogsConfig tests ---

// boolPtr is a test helper that returns a pointer to the given bool value.
func boolPtr(v bool) *bool { return &v }

func TestResolveLogsConfig(t *testing.T) {
	t.Run("nil config returns defaults", func(t *testing.T) {
		resolved := ResolveLogsConfig(nil)
		if !*resolved.Enabled {
			t.Error("Enabled should default to true")
		}
		if resolved.Level != "info" {
			t.Errorf("Level = %q, want %q", resolved.Level, "info")
		}
		if resolved.RetentionDays != 7 {
			t.Errorf("RetentionDays = %d, want 7", resolved.RetentionDays)
		}
	})

	t.Run("empty config returns defaults", func(t *testing.T) {
		resolved := ResolveLogsConfig(&LogsConfig{})
		if !*resolved.Enabled {
			t.Error("Enabled should default to true (nil *bool -> true)")
		}
		if resolved.Level != "info" {
			t.Errorf("Level = %q, want %q (empty falls back to default)", resolved.Level, "info")
		}
		if resolved.RetentionDays != 7 {
			t.Errorf("RetentionDays = %d, want 7 (zero falls back to default)", resolved.RetentionDays)
		}
	})

	t.Run("valid config preserved", func(t *testing.T) {
		cfg := &LogsConfig{
			Enabled:       boolPtr(false),
			Level:         "warn",
			RetentionDays: 14,
		}
		resolved := ResolveLogsConfig(cfg)
		if *resolved.Enabled {
			t.Error("Enabled should be false (explicitly set)")
		}
		if resolved.Level != "warn" {
			t.Errorf("Level = %q, want %q", resolved.Level, "warn")
		}
		if resolved.RetentionDays != 14 {
			t.Errorf("RetentionDays = %d, want 14", resolved.RetentionDays)
		}
	})

	t.Run("invalid retentionDays falls back to 7", func(t *testing.T) {
		cfg := &LogsConfig{
			RetentionDays: 0,
		}
		resolved := ResolveLogsConfig(cfg)
		if resolved.RetentionDays != 7 {
			t.Errorf("RetentionDays = %d, want 7 (0 falls back)", resolved.RetentionDays)
		}
	})

	t.Run("negative retentionDays falls back to 7", func(t *testing.T) {
		cfg := &LogsConfig{
			RetentionDays: -5,
		}
		resolved := ResolveLogsConfig(cfg)
		if resolved.RetentionDays != 7 {
			t.Errorf("RetentionDays = %d, want 7 (negative falls back)", resolved.RetentionDays)
		}
	})

	t.Run("retentionDays 1 is minimum valid value", func(t *testing.T) {
		cfg := &LogsConfig{
			RetentionDays: 1,
		}
		resolved := ResolveLogsConfig(cfg)
		if resolved.RetentionDays != 1 {
			t.Errorf("RetentionDays = %d, want 1", resolved.RetentionDays)
		}
	})

	t.Run("empty level falls back to info", func(t *testing.T) {
		cfg := &LogsConfig{
			Level: "",
		}
		resolved := ResolveLogsConfig(cfg)
		if resolved.Level != "info" {
			t.Errorf("Level = %q, want %q (empty falls back)", resolved.Level, "info")
		}
	})

	t.Run("bogus level preserved for parseLogLevel to handle", func(t *testing.T) {
		// ResolveLogsConfig does not validate level values -- forgelog.parseLogLevel handles that
		cfg := &LogsConfig{
			Level: "bogus",
		}
		resolved := ResolveLogsConfig(cfg)
		if resolved.Level != "bogus" {
			t.Errorf("Level = %q, want %q (ResolveLogsConfig preserves for downstream)", resolved.Level, "bogus")
		}
	})

	t.Run("explicit enabled true preserved", func(t *testing.T) {
		cfg := &LogsConfig{
			Enabled: boolPtr(true),
		}
		resolved := ResolveLogsConfig(cfg)
		if !*resolved.Enabled {
			t.Error("Enabled should be true (explicitly set)")
		}
	})

	t.Run("explicit enabled false preserved", func(t *testing.T) {
		cfg := &LogsConfig{
			Enabled: boolPtr(false),
		}
		resolved := ResolveLogsConfig(cfg)
		if *resolved.Enabled {
			t.Error("Enabled should be false (explicitly set)")
		}
	})
}

func TestReadConfig_LogsBlock(t *testing.T) {
	t.Run("logs block parsed correctly", func(t *testing.T) {
		dir := setupConfig(t, "logs:\n  enabled: false\n  level: warn\n  retentionDays: 14\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Logs == nil {
			t.Fatal("expected Logs non-nil")
		}
		if cfg.Logs.Enabled == nil || *cfg.Logs.Enabled {
			t.Error("Logs.Enabled should be false")
		}
		if cfg.Logs.Level != "warn" {
			t.Errorf("Logs.Level = %q, want %q", cfg.Logs.Level, "warn")
		}
		if cfg.Logs.RetentionDays != 14 {
			t.Errorf("Logs.RetentionDays = %d, want 14", cfg.Logs.RetentionDays)
		}
	})

	t.Run("logs block with defaults", func(t *testing.T) {
		dir := setupConfig(t, "logs:\n  enabled: true\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Logs == nil {
			t.Fatal("expected Logs non-nil")
		}
		if cfg.Logs.Enabled == nil || !*cfg.Logs.Enabled {
			t.Error("Logs.Enabled should be true")
		}
		if cfg.Logs.Level != "" {
			t.Errorf("Logs.Level should be empty (default applied by ResolveLogsConfig), got %q", cfg.Logs.Level)
		}
		if cfg.Logs.RetentionDays != 0 {
			t.Errorf("Logs.RetentionDays should be 0 (default applied by ResolveLogsConfig), got %d", cfg.Logs.RetentionDays)
		}
	})

	t.Run("logs absent is nil", func(t *testing.T) {
		dir := setupConfig(t, "auto:\n  gitPush: true\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Logs != nil {
			t.Error("expected Logs nil when not configured")
		}
	})

	t.Run("existing config without logs deserializes cleanly", func(t *testing.T) {
		// Hard rule: omitempty on Logs field -- existing configs deserialize cleanly
		dir := setupConfig(t, "auto:\n  gitPush: true\nworktree:\n  source-branch: main\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Logs != nil {
			t.Error("expected Logs nil for config without logs section")
		}
		if cfg.Auto == nil || !cfg.Auto.GitPush {
			t.Error("existing fields should still parse correctly")
		}
	})

	t.Run("logs section without enabled key gets nil *bool", func(t *testing.T) {
		dir := setupConfig(t, "logs:\n  level: warn\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Logs == nil {
			t.Fatal("expected Logs non-nil")
		}
		if cfg.Logs.Enabled != nil {
			t.Errorf("Logs.Enabled should be nil (absent from YAML), got %v", *cfg.Logs.Enabled)
		}
		if cfg.Logs.Level != "warn" {
			t.Errorf("Logs.Level = %q, want %q", cfg.Logs.Level, "warn")
		}
	})
}

func TestWriteConfig_LogsBlock(t *testing.T) {
	t.Run("write and read logs block roundtrip", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Logs: &LogsConfig{
				Enabled:       boolPtr(false),
				Level:         "debug",
				RetentionDays: 3,
			},
		}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		readback, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if readback.Logs == nil {
			t.Fatal("expected Logs non-nil")
		}
		if readback.Logs.Enabled == nil || *readback.Logs.Enabled {
			t.Error("Logs.Enabled should be false")
		}
		if readback.Logs.Level != "debug" {
			t.Errorf("Logs.Level = %q, want %q", readback.Logs.Level, "debug")
		}
		if readback.Logs.RetentionDays != 3 {
			t.Errorf("Logs.RetentionDays = %d, want 3", readback.Logs.RetentionDays)
		}
	})

	t.Run("write config without logs omits logs key", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Auto: &AutoConfig{GitPush: true},
		}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Read raw file and check no "logs:" key present
		data, err := os.ReadFile(filepath.Join(dir, ".forge", "config.yaml"))
		if err != nil {
			t.Fatal(err)
		}
		content := string(data)
		if strings.Contains(content, "logs:") {
			t.Errorf("logs key should be omitted when nil, got:\n%s", content)
		}
	})
}
