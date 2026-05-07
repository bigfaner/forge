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

func TestSaveState_ErrorPaths(t *testing.T) {
	t.Run("write permission denied", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "task-state.json")
		state := &TaskState{TaskID: "1.1"}

		if err := SaveState(path, state); err != nil {
			t.Fatalf("initial SaveState() error = %v", err)
		}

		// Make the file read-only
		os.Chmod(path, 0444)
		defer os.Chmod(path, 0644)

		state2 := &TaskState{TaskID: "2.1"}
		err := SaveState(path, state2)
		if err == nil {
			t.Error("expected error for read-only file")
		}
	})

	t.Run("trailing newline preserved", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "task-state.json")
		state := &TaskState{TaskID: "1.1"}

		if err := SaveState(path, state); err != nil {
			t.Fatalf("SaveState() error = %v", err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile() error = %v", err)
		}
		if len(data) == 0 || data[len(data)-1] != '\n' {
			t.Errorf("saved file should end with newline, got %q", data)
		}
	})
}

func TestLoadState_ErrorPaths(t *testing.T) {
	t.Run("read permission denied", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "task-state.json")
		os.WriteFile(path, []byte(`{"task_id":"1.1"}`), 0644)

		os.Chmod(path, 0000)
		defer os.Chmod(path, 0644)

		_, err := LoadState(path)
		if err == nil {
			t.Error("expected error for unreadable file")
		}
	})
}

func TestDeleteState_ErrorPaths(t *testing.T) {
	t.Run("permission denied on directory", func(t *testing.T) {
		dir := t.TempDir()
		subDir := filepath.Join(dir, "sub")
		os.MkdirAll(subDir, 0755)
		path := filepath.Join(subDir, "task-state.json")
		os.WriteFile(path, []byte("{}"), 0644)

		// Make parent directory read-only
		os.Chmod(subDir, 0555)
		defer os.Chmod(subDir, 0755)

		err := DeleteState(path)
		if err == nil {
			t.Error("expected error for read-only directory")
		}
	})
}

func TestStateLifecycle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task-state.json")

	// Save -> Load -> modify -> Save -> Load -> Delete -> Load
	state1 := &TaskState{TaskID: "1.1", Key: "task1", Title: "First", StartedTime: "2024-01-01"}
	if err := SaveState(path, state1); err != nil {
		t.Fatalf("save1: %v", err)
	}

	loaded, err := LoadState(path)
	if err != nil {
		t.Fatalf("load1: %v", err)
	}
	if loaded.TaskID != "1.1" {
		t.Errorf("load1 TaskID = %q, want 1.1", loaded.TaskID)
	}

	// Modify and save again
	loaded.TaskID = "2.1"
	if err := SaveState(path, loaded); err != nil {
		t.Fatalf("save2: %v", err)
	}

	loaded2, err := LoadState(path)
	if err != nil {
		t.Fatalf("load2: %v", err)
	}
	if loaded2.TaskID != "2.1" {
		t.Errorf("load2 TaskID = %q, want 2.1", loaded2.TaskID)
	}

	// Delete
	if err := DeleteState(path); err != nil {
		t.Fatalf("delete: %v", err)
	}

	loaded3, err := LoadState(path)
	if err != nil {
		t.Fatalf("load3: %v", err)
	}
	if loaded3 != nil {
		t.Errorf("load3 should be nil after delete, got %v", loaded3)
	}
}
