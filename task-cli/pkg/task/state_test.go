package task

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadState(t *testing.T) {
	t.Run("file does not exist", func(t *testing.T) {
		state, err := LoadState(filepath.Join(t.TempDir(), "nonexistent.json"))
		if err != nil {
			t.Errorf("LoadState() error = %v, want nil", err)
		}
		if state != nil {
			t.Errorf("LoadState() = %v, want nil", state)
		}
	})

	t.Run("valid state file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "task-state.json")
		wantState := &TaskState{
			TaskID:      "1.1",
			Key:         "task1",
			Title:       "Test Task",
			Priority:    "P0",
			StartedTime: "2024-01-01 10:00",
		}
		if err := SaveState(path, wantState); err != nil {
			t.Fatalf("SaveState() error = %v", err)
		}

		gotState, err := LoadState(path)
		if err != nil {
			t.Fatalf("LoadState() error = %v", err)
		}
		if gotState.TaskID != wantState.TaskID {
			t.Errorf("TaskID = %q, want %q", gotState.TaskID, wantState.TaskID)
		}
		if gotState.Key != wantState.Key {
			t.Errorf("Key = %q, want %q", gotState.Key, wantState.Key)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "invalid.json")
		if err := os.WriteFile(path, []byte("invalid json"), 0644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		_, err := LoadState(path)
		if err == nil {
			t.Error("LoadState() error = nil, want error for invalid JSON")
		}
	})
}

func TestSaveState(t *testing.T) {
	t.Run("creates directory if not exists", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "subdir", "nested", "task-state.json")
		state := &TaskState{
			TaskID:      "1.1",
			Key:         "task1",
			Title:       "Test Task",
			StartedTime: "2024-01-01 10:00",
		}

		if err := SaveState(path, state); err != nil {
			t.Fatalf("SaveState() error = %v", err)
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Error("SaveState() did not create file")
		}
	})

	t.Run("overwrites existing file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "task-state.json")

		state1 := &TaskState{TaskID: "1.1"}
		if err := SaveState(path, state1); err != nil {
			t.Fatalf("SaveState() error = %v", err)
		}

		state2 := &TaskState{TaskID: "2.1"}
		if err := SaveState(path, state2); err != nil {
			t.Fatalf("SaveState() error = %v", err)
		}

		got, err := LoadState(path)
		if err != nil {
			t.Fatalf("LoadState() error = %v", err)
		}
		if got.TaskID != "2.1" {
			t.Errorf("TaskID = %q, want %q", got.TaskID, "2.1")
		}
	})
}

func TestDeleteState(t *testing.T) {
	t.Run("deletes existing file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "task-state.json")
		state := &TaskState{TaskID: "test"}
		if err := SaveState(path, state); err != nil {
			t.Fatalf("SaveState() error = %v", err)
		}

		if err := DeleteState(path); err != nil {
			t.Fatalf("DeleteState() error = %v", err)
		}

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Error("DeleteState() did not remove file")
		}
	})

	t.Run("no error if file does not exist", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "nonexistent.json")

		if err := DeleteState(path); err != nil {
			t.Errorf("DeleteState() error = %v, want nil", err)
		}
	})
}
