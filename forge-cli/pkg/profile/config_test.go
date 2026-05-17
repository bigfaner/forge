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
		content := "project-type: backend\ntest-profiles:\n  - go-test\ncapabilities:\n  - tui\n  - api\n"
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
		if len(cfg.TestProfiles) != 1 || cfg.TestProfiles[0] != "go-test" {
			t.Errorf("expected [go-test], got %v", cfg.TestProfiles)
		}
		if len(cfg.Capabilities) != 2 || cfg.Capabilities[0] != "tui" || cfg.Capabilities[1] != "api" {
			t.Errorf("expected [tui api], got %v", cfg.Capabilities)
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
		dir := setupConfig(t, "project-type: frontend\ntest-profiles:\n  - go-test\n")
		val, err := GetConfigValue(dir, "project-type")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "frontend" {
			t.Errorf("expected 'frontend', got %q", val)
		}
	})

	t.Run("capabilities array", func(t *testing.T) {
		dir := setupConfig(t, "capabilities:\n  - tui\n  - api\n  - cli\n")
		val, err := GetConfigValue(dir, "capabilities")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "tui\napi\ncli"
		if val != expected {
			t.Errorf("expected %q, got %q", expected, val)
		}
	})

	t.Run("test-profiles array", func(t *testing.T) {
		dir := setupConfig(t, "test-profiles:\n  - go-test\n  - pytest\n")
		val, err := GetConfigValue(dir, "test-profiles")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "go-test\npytest"
		if val != expected {
			t.Errorf("expected %q, got %q", expected, val)
		}
	})

	t.Run("auto.gitPush true", func(t *testing.T) {
		dir := setupConfig(t, "test-profiles:\n  - go-test\nauto:\n  gitPush: true\n")
		val, err := GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "true" {
			t.Errorf("expected 'true', got %q", val)
		}
	})

	t.Run("auto.gitPush false", func(t *testing.T) {
		dir := setupConfig(t, "test-profiles:\n  - go-test\nauto:\n  gitPush: false\n")
		val, err := GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false', got %q", val)
		}
	})

	t.Run("auto.gitPush absent returns false (default)", func(t *testing.T) {
		dir := setupConfig(t, "test-profiles:\n  - go-test\n")
		val, err := GetConfigValue(dir, "auto.gitPush")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "false" {
			t.Errorf("expected 'false' default, got %q", val)
		}
	})

	t.Run("auto block absent returns false (default)", func(t *testing.T) {
		dir := setupConfig(t, "test-profiles:\n  - go-test\n")
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
}

func TestReadTestProfiles(t *testing.T) {
	t.Run("file not exists", func(t *testing.T) {
		dir := t.TempDir()
		profiles, err := ReadTestProfiles(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if profiles != nil {
			t.Fatalf("expected nil, got %v", profiles)
		}
	})

	t.Run("valid config", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		configContent := "test-profiles:\n  - web-playwright\n  - go-test\n"
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
			t.Fatal(err)
		}

		profiles, err := ReadTestProfiles(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(profiles) != 2 || profiles[0] != "web-playwright" || profiles[1] != "go-test" {
			t.Fatalf("expected [web-playwright go-test], got %v", profiles)
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

		profiles, err := ReadTestProfiles(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if profiles != nil {
			t.Fatalf("expected nil for empty config, got %v", profiles)
		}
	})
}

func TestWriteTestProfiles(t *testing.T) {
	t.Run("create new file", func(t *testing.T) {
		dir := t.TempDir()
		if err := WriteTestProfiles(dir, []string{"go-test"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify written
		profiles, err := ReadTestProfiles(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(profiles) != 1 || profiles[0] != "go-test" {
			t.Fatalf("expected [go-test], got %v", profiles)
		}
	})

	t.Run("overwrite existing", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte("test-profiles:\n  - old\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		if err := WriteTestProfiles(dir, []string{"rust-test"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		profiles, err := ReadTestProfiles(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(profiles) != 1 || profiles[0] != "rust-test" {
			t.Fatalf("expected [rust-test], got %v", profiles)
		}
	})
}

func TestIsKnownProfile(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"web-playwright", true},
		{"go-test", true},
		{"maestro", true},
		{"java-junit", true},
		{"rust-test", true},
		{"pytest", true},
		{"unknown", false},
		{"", false},
		{"Web-Playwright", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsKnownProfile(tt.name); got != tt.expected {
				t.Errorf("IsKnownProfile(%q) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}
