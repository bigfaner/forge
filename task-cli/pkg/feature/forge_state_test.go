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
		os.MkdirAll(filepath.Dir(statePath), 0755)
		content := ForgeState{Feature: "f1", AllCompleted: true, UpdatedAt: "2026-01-01T00:00:00Z"}
		data, _ := json.Marshal(content)
		os.WriteFile(statePath, data, 0644)

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
		os.MkdirAll(filepath.Dir(statePath), 0755)
		os.WriteFile(statePath, []byte("not json"), 0644)

		state := ReadForgeState(dir)
		if state != nil {
			t.Errorf("expected nil for invalid JSON, got %+v", state)
		}
	})
}

func TestClearForgeState(t *testing.T) {
	t.Run("deletes existing file", func(t *testing.T) {
		dir := t.TempDir()
		WriteForgeState(dir, "test")

		err := ClearForgeState(dir)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if _, err := os.Stat(GetForgeStatePath(dir)); !os.IsNotExist(err) {
			t.Error("expected file to be deleted")
		}
	})

	t.Run("no error when file doesn't exist", func(t *testing.T) {
		dir := t.TempDir()
		err := ClearForgeState(dir)
		if err != nil && !os.IsNotExist(err) {
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
		WriteForgeState(dir, "my-feature")
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
