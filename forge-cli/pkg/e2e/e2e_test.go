package e2e

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"forge-cli/pkg/feature"
)

func TestResolveProfile(t *testing.T) {
	t.Run("valid profile from config", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		configContent := "languages:\n  - go\n"
		if err := os.WriteFile(filepath.Join(forgeDir, feature.ForgeConfigFileName), []byte(configContent), 0o644); err != nil {
			t.Fatal(err)
		}

		got, err := ResolveProfile(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "go" {
			t.Fatalf("expected go, got %q", got)
		}
	})

	t.Run("multiple profiles returns first", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		configContent := "languages:\n  - javascript\n  - go\n"
		if err := os.WriteFile(filepath.Join(forgeDir, feature.ForgeConfigFileName), []byte(configContent), 0o644); err != nil {
			t.Fatal(err)
		}

		got, err := ResolveProfile(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "javascript" {
			t.Fatalf("expected javascript, got %q", got)
		}
	})

	t.Run("no config file returns ErrNoProfile", func(t *testing.T) {
		dir := t.TempDir()

		_, err := ResolveProfile(dir)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, ErrNoProfile) {
			t.Fatalf("expected ErrNoProfile, got %v", err)
		}
	})

	t.Run("empty languages returns ErrNoProfile", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		configContent := "{}"
		if err := os.WriteFile(filepath.Join(forgeDir, feature.ForgeConfigFileName), []byte(configContent), 0o644); err != nil {
			t.Fatal(err)
		}

		_, err := ResolveProfile(dir)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, ErrNoProfile) {
			t.Fatalf("expected ErrNoProfile, got %v", err)
		}
	})

	t.Run("unknown profile returns ErrBadProfile", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		configContent := "languages:\n  - fake-profile\n"
		if err := os.WriteFile(filepath.Join(forgeDir, feature.ForgeConfigFileName), []byte(configContent), 0o644); err != nil {
			t.Fatal(err)
		}

		_, err := ResolveProfile(dir)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, ErrBadProfile) {
			t.Fatalf("expected ErrBadProfile, got %v", err)
		}
		// Check error message contains the profile name
		expected := "unknown profile: fake-profile"
		if err.Error() != expected {
			t.Fatalf("expected error message %q, got %q", expected, err.Error())
		}
	})

	t.Run("malformed config returns read error", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		configContent := "languages: [not closed"
		if err := os.WriteFile(filepath.Join(forgeDir, feature.ForgeConfigFileName), []byte(configContent), 0o644); err != nil {
			t.Fatal(err)
		}

		_, err := ResolveProfile(dir)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		// Should be wrapped "read languages: ..." error, not ErrNoProfile or ErrBadProfile
		if errors.Is(err, ErrNoProfile) || errors.Is(err, ErrBadProfile) {
			t.Fatalf("unexpected sentinel error for malformed config: %v", err)
		}
	})
}

func TestRunOpts(t *testing.T) {
	opts := RunOpts{
		ProjectRoot: "/tmp/project",
		Feature:     "my-feature",
		Force:       true,
	}
	if opts.ProjectRoot != "/tmp/project" {
		t.Fatalf("unexpected ProjectRoot: %q", opts.ProjectRoot)
	}
	if opts.Feature != "my-feature" {
		t.Fatalf("unexpected Feature: %q", opts.Feature)
	}
	if !opts.Force {
		t.Fatal("expected Force to be true")
	}
}
