package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LoadIndex loads the task index from the given file path.
func LoadIndex(path string) (*TaskIndex, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read index: %w", err)
	}
	var index TaskIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("failed to parse index: %w", err)
	}
	return &index, nil
}

// SaveIndex saves the task index to the given file path.
func SaveIndex(path string, index *TaskIndex) error {
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}
	return nil
}

// ByID returns the task matching the given ID or map key.
// Returns (Task, false) if not found.
func (ti *TaskIndex) ByID(id string) (Task, bool) {
	if t, ok := ti.Tasks[id]; ok {
		return t, true
	}
	for _, t := range ti.Tasks {
		if t.ID == id {
			return t, true
		}
	}
	return Task{}, false
}

// FindTask finds a task by ID or key and returns the map key for write-back.
func FindTask(index *TaskIndex, idOrKey string) (string, *Task, error) {
	if t, ok := index.ByID(idOrKey); ok {
		if _, direct := index.Tasks[idOrKey]; direct {
			return idOrKey, &t, nil
		}
		for key, task := range index.Tasks {
			if task.ID == idOrKey {
				return key, &t, nil
			}
		}
	}
	return "", nil, fmt.Errorf("task not found: %s", idOrKey)
}

// IsValidStatus checks if the status is valid.
func IsValidStatus(index *TaskIndex, status string) bool {
	for _, s := range index.StatusEnum {
		if s == status {
			return true
		}
	}
	return false
}

// EnsureIndexDir ensures the directory for the index file exists.
func EnsureIndexDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}
