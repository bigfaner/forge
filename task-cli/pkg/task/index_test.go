package task

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadIndex(t *testing.T) {
	t.Run("file does not exist", func(t *testing.T) {
		_, err := LoadIndex(filepath.Join(t.TempDir(), "nonexistent.json"))
		if err == nil {
			t.Error("LoadIndex() error = nil, want error for nonexistent file")
		}
	})

	t.Run("valid index file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "index.json")
		wantIndex := &TaskIndex{
			Feature: "test-feature",
			tasks: map[string]Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "pending"},
			},
			StatusEnum:   []string{"pending", "completed"},
			PriorityEnum: []string{"P0", "P1"},
		}
		if err := SaveIndex(path, wantIndex); err != nil {
			t.Fatalf("SaveIndex() error = %v", err)
		}

		gotIndex, err := LoadIndex(path)
		if err != nil {
			t.Fatalf("LoadIndex() error = %v", err)
		}
		if gotIndex.Feature != wantIndex.Feature {
			t.Errorf("Feature = %q, want %q", gotIndex.Feature, wantIndex.Feature)
		}
		if len(gotIndex.tasks) != len(wantIndex.tasks) {
			t.Errorf("Tasks count = %d, want %d", len(gotIndex.tasks), len(wantIndex.tasks))
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "invalid.json")
		if err := os.WriteFile(path, []byte("invalid json"), 0644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		_, err := LoadIndex(path)
		if err == nil {
			t.Error("LoadIndex() error = nil, want error for invalid JSON")
		}
	})
}

func TestSaveIndex(t *testing.T) {
	t.Run("saves valid index", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "index.json")
		index := &TaskIndex{
			Feature: "test-feature",
			tasks: map[string]Task{
				"task1": {ID: "1.1", Title: "Task 1", Status: "pending"},
			},
		}

		if err := SaveIndex(path, index); err != nil {
			t.Fatalf("SaveIndex() error = %v", err)
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Error("SaveIndex() did not create file")
		}
	})

	t.Run("creates parent directories", func(t *testing.T) {
		dir := t.TempDir()
		// EnsureIndexDir should create parent dirs, not SaveIndex
		path := filepath.Join(dir, "subdir", "index.json")
		index := &TaskIndex{Feature: "test"}

		// SaveIndex should fail if parent doesn't exist
		err := SaveIndex(path, index)
		if err == nil {
			t.Error("SaveIndex() should fail when parent directory doesn't exist")
		}
	})
}

func TestByID(t *testing.T) {
	index := &TaskIndex{
		Feature: "test-feature",
		tasks: map[string]Task{
			"task1":    {ID: "1.1", Title: "Task 1", Status: "pending"},
			"run-e2e":  {ID: "T-test-3", Title: "Run e2e", Status: "pending"},
			"disc-1":   {ID: "disc-1", Title: "Discovery", Status: "pending"},
		},
	}

	t.Run("key equals ID", func(t *testing.T) {
		got, ok := index.ByID("disc-1")
		if !ok {
			t.Fatal("ByID() returned false")
		}
		if got.Title != "Discovery" {
			t.Errorf("Title = %q, want %q", got.Title, "Discovery")
		}
	})

	t.Run("key differs from ID, lookup by ID", func(t *testing.T) {
		got, ok := index.ByID("T-test-3")
		if !ok {
			t.Fatal("ByID() returned false")
		}
		if got.Title != "Run e2e" {
			t.Errorf("Title = %q, want %q", got.Title, "Run e2e")
		}
	})

	t.Run("key differs from ID, lookup by key", func(t *testing.T) {
		got, ok := index.ByID("run-e2e")
		if !ok {
			t.Fatal("ByID() returned false")
		}
		if got.ID != "T-test-3" {
			t.Errorf("ID = %q, want %q", got.ID, "T-test-3")
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, ok := index.ByID("nonexistent")
		if ok {
			t.Error("ByID() should return false for nonexistent ID")
		}
	})

	t.Run("empty index", func(t *testing.T) {
		empty := &TaskIndex{tasks: map[string]Task{}}
		_, ok := empty.ByID("anything")
		if ok {
			t.Error("ByID() should return false in empty index")
		}
	})
}

func TestFindTask(t *testing.T) {
	index := &TaskIndex{
		Feature: "test-feature",
		tasks: map[string]Task{
			"task1":   {ID: "1.1", Title: "Task 1", Status: "pending"},
			"task2":   {ID: "1.2", Title: "Task 2", Status: "pending"},
			"run-e2e": {ID: "T-test-3", Title: "Run e2e", Status: "pending"},
		},
	}

	tests := []struct {
		name    string
		idOrKey string
		wantKey string
		wantID  string
		wantErr bool
	}{
		{"find by key (key==ID)", "task1", "task1", "1.1", false},
		{"find by ID", "1.2", "task2", "1.2", false},
		{"find by ID when key differs", "T-test-3", "run-e2e", "T-test-3", false},
		{"find by key when key differs from ID", "run-e2e", "run-e2e", "T-test-3", false},
		{"not found", "nonexistent", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, task, err := FindTask(index, tt.idOrKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if key != tt.wantKey {
					t.Errorf("key = %q, want %q", key, tt.wantKey)
				}
				if task.ID != tt.wantID {
					t.Errorf("task.ID = %q, want %q", task.ID, tt.wantID)
				}
			}
		})
	}
}

func TestIsValidStatus(t *testing.T) {
	index := &TaskIndex{
		StatusEnum: []string{"pending", "in_progress", "completed", "blocked", "skipped"},
	}

	tests := []struct {
		status   string
		expected bool
	}{
		{"pending", true},
		{"in_progress", true},
		{"completed", true},
		{"blocked", true},
		{"skipped", true},
		{"invalid", false},
		{"", false},
		{"PENDING", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			if got := IsValidStatus(index, tt.status); got != tt.expected {
				t.Errorf("IsValidStatus(%q) = %v, want %v", tt.status, got, tt.expected)
			}
		})
	}
}

func TestEnsureIndexDir(t *testing.T) {
	t.Run("creates directory", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "subdir", "nested", "index.json")

		if err := EnsureIndexDir(path); err != nil {
			t.Fatalf("EnsureIndexDir() error = %v", err)
		}

		expectedDir := filepath.Dir(path)
		if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
			t.Error("EnsureIndexDir() did not create directory")
		}
	})

	t.Run("no error if directory exists", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "index.json")

		if err := EnsureIndexDir(path); err != nil {
			t.Fatalf("EnsureIndexDir() error = %v", err)
		}
	})
}
