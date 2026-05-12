package profile

import (
	"os"
	"path/filepath"
	"testing"
)

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
