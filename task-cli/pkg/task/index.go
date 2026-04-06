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

// FindTask finds a task by ID or key.
func FindTask(index *TaskIndex, idOrKey string) (string, *Task, error) {
	// Try exact key match first
	if task, ok := index.Tasks[idOrKey]; ok {
		return idOrKey, &task, nil
	}
	// Try ID match
	for key, task := range index.Tasks {
		if task.ID == idOrKey {
			return key, &task, nil
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
