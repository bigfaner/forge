package feature

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureForgeDir(t *testing.T) {
	t.Run("creates .forge/ directory", func(t *testing.T) {
		dir := t.TempDir()

		if err := EnsureForgeDir(dir); err != nil {
			t.Fatalf("EnsureForgeDir() error = %v", err)
		}

		forgeDir := filepath.Join(dir, ForgeDir)
		info, err := os.Stat(forgeDir)
		if err != nil {
			t.Fatalf(".forge/ directory not created: %v", err)
		}
		if !info.IsDir() {
			t.Error(".forge is not a directory")
		}
	})

	t.Run("idempotent when directory exists", func(t *testing.T) {
		dir := t.TempDir()

		if err := EnsureForgeDir(dir); err != nil {
			t.Fatalf("first EnsureForgeDir() error = %v", err)
		}
		if err := EnsureForgeDir(dir); err != nil {
			t.Fatalf("second EnsureForgeDir() error = %v", err)
		}
	})
}

func TestWriteForgeState(t *testing.T) {
	dir := t.TempDir()

	err := WriteForgeState(dir, "my-feature")
	if err != nil {
		t.Fatalf("WriteForgeState() error = %v", err)
	}

	statePath := GetForgeStatePath(dir)
	data, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatalf("failed to read state file: %v", err)
	}

	var state ForgeState
	if err := json.Unmarshal(data, &state); err != nil {
		t.Fatalf("failed to parse state: %v", err)
	}

	if state.Feature != "my-feature" {
		t.Errorf("feature = %q, want %q", state.Feature, "my-feature")
	}
	if !state.AllCompleted {
		t.Error("allCompleted = false, want true")
	}
	if state.UpdatedAt == "" {
		t.Error("updatedAt is empty")
	}
}

func TestReadForgeState(t *testing.T) {
	t.Run("file exists and valid", func(t *testing.T) {
		dir := t.TempDir()
		statePath := GetForgeStatePath(dir)
		_ = os.MkdirAll(filepath.Dir(statePath), 0755)
		content := ForgeState{Feature: "f1", AllCompleted: true, UpdatedAt: "2026-01-01T00:00:00Z"}
		data, _ := json.Marshal(content)
		_ = os.WriteFile(statePath, data, 0644)

		state := ReadForgeState(dir)
		if state == nil {
			t.Fatal("expected non-nil state")
		}
		if state.Feature != "f1" {
			t.Errorf("feature = %q, want %q", state.Feature, "f1")
		}
		if !state.AllCompleted {
			t.Error("allCompleted = false, want true")
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		dir := t.TempDir()
		state := ReadForgeState(dir)
		if state != nil {
			t.Errorf("expected nil, got %+v", state)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		dir := t.TempDir()
		statePath := GetForgeStatePath(dir)
		_ = os.MkdirAll(filepath.Dir(statePath), 0755)
		_ = os.WriteFile(statePath, []byte("not json"), 0644)

		state := ReadForgeState(dir)
		if state != nil {
			t.Errorf("expected nil for invalid JSON, got %+v", state)
		}
	})
}

func TestClearForgeState(t *testing.T) {
	t.Run("writes allCompleted=false instead of deleting", func(t *testing.T) {
		dir := t.TempDir()
		_ = WriteForgeState(dir, "test")

		err := ClearForgeState(dir)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// File should still exist
		statePath := GetForgeStatePath(dir)
		if _, err := os.Stat(statePath); os.IsNotExist(err) {
			t.Fatal("state file should still exist after ClearForgeState")
		}

		// allCompleted should be false
		state := ReadForgeState(dir)
		if state == nil {
			t.Fatal("expected non-nil state after ClearForgeState")
		}
		if state.AllCompleted {
			t.Error("allCompleted should be false after ClearForgeState")
		}
		if state.Feature != "test" {
			t.Errorf("feature = %q, want %q", state.Feature, "test")
		}
	})

	t.Run("no error when file doesn't exist", func(t *testing.T) {
		dir := t.TempDir()
		err := ClearForgeState(dir)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestEnsureForgeState(t *testing.T) {
	t.Run("creates state.json with allCompleted=false", func(t *testing.T) {
		dir := t.TempDir()

		if err := EnsureForgeState(dir, "my-feature"); err != nil {
			t.Fatalf("EnsureForgeState() error = %v", err)
		}

		statePath := GetForgeStatePath(dir)
		data, err := os.ReadFile(statePath)
		if err != nil {
			t.Fatalf("state.json not created: %v", err)
		}

		var state ForgeState
		if err := json.Unmarshal(data, &state); err != nil {
			t.Fatalf("failed to parse state: %v", err)
		}

		if state.Feature != "my-feature" {
			t.Errorf("feature = %q, want %q", state.Feature, "my-feature")
		}
		if state.AllCompleted {
			t.Error("allCompleted = true, want false")
		}
		if state.UpdatedAt == "" {
			t.Error("updatedAt is empty")
		}
	})

	t.Run("creates .forge/ directory if missing", func(t *testing.T) {
		dir := t.TempDir()

		if err := EnsureForgeState(dir, "test"); err != nil {
			t.Fatalf("EnsureForgeState() error = %v", err)
		}

		forgeDir := filepath.Join(dir, ForgeDir)
		info, err := os.Stat(forgeDir)
		if err != nil {
			t.Fatalf(".forge/ directory not created: %v", err)
		}
		if !info.IsDir() {
			t.Error(".forge is not a directory")
		}
	})

	t.Run("overwrites existing allCompleted=true", func(t *testing.T) {
		dir := t.TempDir()

		// Simulate: all tasks were done, then fix-e2e tasks added
		_ = WriteForgeState(dir, "my-feature")
		state := ReadForgeState(dir)
		if !state.AllCompleted {
			t.Fatal("setup: WriteForgeState should set allCompleted=true")
		}

		// Claim overwrites with false
		if err := EnsureForgeState(dir, "my-feature"); err != nil {
			t.Fatalf("EnsureForgeState() error = %v", err)
		}

		state = ReadForgeState(dir)
		if state.AllCompleted {
			t.Error("allCompleted should be false after EnsureForgeState overwrite")
		}
	})
}

func TestMarkFeatureCompleted(t *testing.T) {
	t.Run("sets completedAt on existing state", func(t *testing.T) {
		dir := t.TempDir()
		_ = EnsureForgeState(dir, "test-feature")

		err := MarkFeatureCompleted(dir)
		if err != nil {
			t.Fatalf("MarkFeatureCompleted() error = %v", err)
		}

		state := ReadForgeState(dir)
		if state == nil {
			t.Fatal("expected non-nil state")
		}
		if state.CompletedAt == "" {
			t.Error("completedAt should be set")
		}
		if state.Feature != "test-feature" {
			t.Errorf("feature should be preserved, got %q", state.Feature)
		}
	})

	t.Run("no-op when state file does not exist", func(t *testing.T) {
		dir := t.TempDir()
		err := MarkFeatureCompleted(dir)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("no-op when state file is malformed", func(t *testing.T) {
		dir := t.TempDir()
		statePath := GetForgeStatePath(dir)
		_ = os.MkdirAll(filepath.Dir(statePath), 0755)
		_ = os.WriteFile(statePath, []byte("not json"), 0644)

		err := MarkFeatureCompleted(dir)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestForgeState_AtomicWriteNoTempFiles(t *testing.T) {
	dir := t.TempDir()

	// Write multiple times and verify no temp files remain
	for i := 0; i < 5; i++ {
		_ = WriteForgeState(dir, "test")
		_ = EnsureForgeState(dir, "test")
	}

	forgeDir := filepath.Join(dir, ForgeDir)
	entries, err := os.ReadDir(forgeDir)
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}
	for _, e := range entries {
		if matched, _ := filepath.Match(".state.json.tmp.*", e.Name()); matched {
			t.Errorf("temp file should not remain: %s", e.Name())
		}
	}
}
