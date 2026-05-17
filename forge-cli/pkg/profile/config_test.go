package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadConfig(t *testing.T) {
	t.Run("file not exists", func(t *testing.T) {
		dir := t.TempDir()
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg != nil {
			t.Fatalf("expected nil, got %v", cfg)
		}
	})

	t.Run("full config", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		content := "project-type: backend\nlanguages:\n  - go\ninterfaces:\n  - tui\n  - api\n"
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.ProjectType != "backend" {
			t.Errorf("expected project-type backend, got %q", cfg.ProjectType)
		}
		if len(cfg.Languages) != 1 || cfg.Languages[0] != "go" {
			t.Errorf("expected [go], got %v", cfg.Languages)
		}
		if len(cfg.Interfaces) != 2 || cfg.Interfaces[0] != "tui" || cfg.Interfaces[1] != "api" {
			t.Errorf("expected [tui api], got %v", cfg.Interfaces)
		}
	})

	t.Run("empty config", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte("{}"), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.ProjectType != "" {
			t.Errorf("expected empty project-type, got %q", cfg.ProjectType)
		}
	})

	t.Run("config with test-framework and test-command", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		content := "languages:\n  - go\ntest-framework: pytest\ntest-command: pytest tests/\n"
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.TestFramework != "pytest" {
			t.Errorf("expected test-framework pytest, got %q", cfg.TestFramework)
		}
		if cfg.TestCommand != "pytest tests/" {
			t.Errorf("expected test-command 'pytest tests/', got %q", cfg.TestCommand)
		}
	})

	t.Run("config without test-framework defaults empty", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		content := "languages:\n  - go\n"
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.TestFramework != "" {
			t.Errorf("expected empty test-framework, got %q", cfg.TestFramework)
		}
		if cfg.TestCommand != "" {
			t.Errorf("expected empty test-command, got %q", cfg.TestCommand)
		}
	})
}

func TestGetConfigValue(t *testing.T) {
	setupConfig := func(t *testing.T, content string) string {
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

	t.Run("project-type scalar", func(t *testing.T) {
		dir := setupConfig(t, "project-type: frontend\nlanguages:\n  - go\n")
		val, err := GetConfigValue(dir, "project-type")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "frontend" {
			t.Errorf("expected 'frontend', got %q", val)
		}
	})

	t.Run("interfaces array", func(t *testing.T) {
		dir := setupConfig(t, "interfaces:\n  - tui\n  - api\n  - cli\n")
		val, err := GetConfigValue(dir, "interfaces")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "tui\napi\ncli"
		if val != expected {
			t.Errorf("expected %q, got %q", expected, val)
		}
	})

	t.Run("languages array", func(t *testing.T) {
		dir := setupConfig(t, "languages:\n  - go\n  - python\n")
		val, err := GetConfigValue(dir, "languages")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "go\npython"
		if val != expected {
			t.Errorf("expected %q, got %q", expected, val)
		}
	})

	t.Run("auto.gitPush true", func(t *testing.T) {
		dir := setupConfig(t, "languages:\n  - go\nauto:\n  gitPush: true\n")
		val, err := GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("auto.gitPush false", func(t *testing.T) {
		dir := setupConfig(t, "languages:\n  - go\nauto:\n  gitPush: false\n")
		val, err := GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false', got %q", val)
		}
	})

	t.Run("auto.gitPush absent returns false (default)", func(t *testing.T) {
		dir := setupConfig(t, "languages:\n  - go\n")
		val, err := GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false' default, got %q", val)
		}
	})

	t.Run("auto block absent returns false (default)", func(t *testing.T) {
		dir := setupConfig(t, "languages:\n  - go\n")
		val, err := GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false' default, got %q", val)
		}
	})

	t.Run("unknown key returns error", func(t *testing.T) {
		dir := setupConfig(t, "project-type: backend\n")
		_, err := GetConfigValue(dir, "nonexistent")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})

	t.Run("missing file returns error", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "project-type")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})

	t.Run("key exists but empty value returns error", func(t *testing.T) {
		dir := setupConfig(t, "project-type: ''\n")
		_, err := GetConfigValue(dir, "project-type")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound for empty string, got %v", err)
		}
	})

	t.Run("test-framework returns value", func(t *testing.T) {
		dir := setupConfig(t, "languages:\n  - go\ntest-framework: pytest\n")
		val, err := GetConfigValue(dir, "test-framework")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "pytest" {
			t.Errorf("expected 'pytest', got %q", val)
		}
	})

	t.Run("test-framework absent returns error", func(t *testing.T) {
		dir := setupConfig(t, "languages:\n  - go\n")
		_, err := GetConfigValue(dir, "test-framework")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})

	t.Run("test-command returns value", func(t *testing.T) {
		dir := setupConfig(t, "languages:\n  - go\ntest-command: go test ./...\n")
		val, err := GetConfigValue(dir, "test-command")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "go test ./..." {
			t.Errorf("expected 'go test ./...', got %q", val)
		}
	})

	t.Run("test-command absent returns error", func(t *testing.T) {
		dir := setupConfig(t, "languages:\n  - go\n")
		_, err := GetConfigValue(dir, "test-command")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})
}

func TestReadLanguages(t *testing.T) {
	t.Run("file not exists falls back to detect", func(t *testing.T) {
		dir := t.TempDir()
		// No config file, no project files — should return empty
		profiles, err := ReadLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if profiles != nil {
			t.Fatalf("expected nil when nothing detected, got %v", profiles)
		}
	})

	t.Run("languages from config", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		configContent := "languages:\n  - javascript\n  - go\n"
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
			t.Fatal(err)
		}

		profiles, err := ReadLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(profiles) != 2 || profiles[0] != "javascript" || profiles[1] != "go" {
			t.Fatalf("expected [javascript go], got %v", profiles)
		}
	})

	t.Run("empty config falls back to detect", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte("{}"), 0o644); err != nil {
			t.Fatal(err)
		}

		// No project files to detect, should return nil
		profiles, err := ReadLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if profiles != nil {
			t.Fatalf("expected nil for empty config with no detectable files, got %v", profiles)
		}
	})
}

func TestWriteLanguages(t *testing.T) {
	t.Run("create new file", func(t *testing.T) {
		dir := t.TempDir()
		if err := WriteLanguages(dir, []string{"go"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify written
		profiles, err := ReadLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(profiles) != 1 || profiles[0] != "go" {
			t.Fatalf("expected [go], got %v", profiles)
		}
	})

	t.Run("overwrite existing", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte("languages:\n  - old\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		if err := WriteLanguages(dir, []string{"rust"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		profiles, err := ReadLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(profiles) != 1 || profiles[0] != "rust" {
			t.Fatalf("expected [rust], got %v", profiles)
		}
	})
}

func TestWorktreeConfig(t *testing.T) {
	t.Run("worktree config parsed from yaml", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		content := "worktree:\n  source-branch: develop\n  copy-files:\n    - .env\n    - .env.local\n"
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Worktree == nil {
			t.Fatal("expected Worktree to be non-nil")
		}
		if cfg.Worktree.SourceBranch != "develop" {
			t.Errorf("expected source-branch 'develop', got %q", cfg.Worktree.SourceBranch)
		}
		if len(cfg.Worktree.CopyFiles) != 2 || cfg.Worktree.CopyFiles[0] != ".env" || cfg.Worktree.CopyFiles[1] != ".env.local" {
			t.Errorf("expected [.env .env.local], got %v", cfg.Worktree.CopyFiles)
		}
	})

	t.Run("worktree absent is backward compatible", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		content := "project-type: backend\nlanguages:\n  - go\n"
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Worktree != nil {
			t.Error("expected Worktree to be nil when not configured")
		}
	})

	t.Run("worktree with only source-branch", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		content := "worktree:\n  source-branch: main\n"
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Worktree == nil {
			t.Fatal("expected Worktree to be non-nil")
		}
		if cfg.Worktree.SourceBranch != "main" {
			t.Errorf("expected source-branch 'main', got %q", cfg.Worktree.SourceBranch)
		}
		if cfg.Worktree.CopyFiles != nil {
			t.Errorf("expected CopyFiles nil, got %v", cfg.Worktree.CopyFiles)
		}
	})

	t.Run("worktree with only copy-files", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		content := "worktree:\n  copy-files:\n    - .env\n"
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Worktree == nil {
			t.Fatal("expected Worktree to be non-nil")
		}
		if cfg.Worktree.SourceBranch != "" {
			t.Errorf("expected empty source-branch, got %q", cfg.Worktree.SourceBranch)
		}
		if len(cfg.Worktree.CopyFiles) != 1 || cfg.Worktree.CopyFiles[0] != ".env" {
			t.Errorf("expected [.env], got %v", cfg.Worktree.CopyFiles)
		}
	})
}

func TestGetConfigValue_Worktree(t *testing.T) {
	setupConfig := func(t *testing.T, content string) string {
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
		dir := setupConfig(t, "project-type: backend\n")
		_, err := GetConfigValue(dir, "worktree.source-branch")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})

	t.Run("worktree.copy-files absent returns error", func(t *testing.T) {
		dir := setupConfig(t, "project-type: backend\n")
		_, err := GetConfigValue(dir, "worktree.copy-files")
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

	t.Run("worktree key with no config file returns error", func(t *testing.T) {
		dir := t.TempDir()
		_, err := GetConfigValue(dir, "worktree.source-branch")
		if err != ErrKeyNotFound {
			t.Errorf("expected ErrKeyNotFound, got %v", err)
		}
	})
}

func TestIsKnownLanguage(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"go", true},
		{"javascript", true},
		{"mobile", true},
		{"java", true},
		{"rust", true},
		{"python", true},
		{"unknown", false},
		{"", false},
		{"JavaScript", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsKnownLanguage(tt.name); got != tt.expected {
				t.Errorf("IsKnownLanguage(%q) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}
