package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"forge-cli/pkg/types"
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
	index.ensureStatusEnumHas(string(types.StatusSuspended))
	return &index, nil
}

// ensureStatusEnumHas adds the given status to StatusEnum if not already present.
func (ti *TaskIndex) ensureStatusEnumHas(status string) {
	for _, s := range ti.StatusEnum {
		if s == status {
			return
		}
	}
	ti.StatusEnum = append(ti.StatusEnum, status)
}

// SaveIndex saves the task index to the given file path.
func SaveIndex(path string, index *TaskIndex) error {
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}
	return nil
}

// ByID returns the task matching the given ID or map key.
// Returns (Task, false) if not found.
func (ti *TaskIndex) ByID(id string) (Task, bool) {
	if t, ok := ti.tasks[id]; ok {
		return t, true
	}
	for _, t := range ti.tasks {
		if t.ID == id {
			return t, true
		}
	}
	return Task{}, false
}

// FindTask finds a task by ID or key and returns the map key for write-back.
func FindTask(index *TaskIndex, idOrKey string) (string, *Task, error) {
	if t, ok := index.ByID(idOrKey); ok {
		if _, direct := index.tasks[idOrKey]; direct {
			return idOrKey, &t, nil
		}
		for key, task := range index.tasks {
			if task.ID == idOrKey {
				return key, &t, nil
			}
		}
	}
	return "", nil, fmt.Errorf("task not found: %s", idOrKey)
}

// SetTask inserts or updates a task under the given map key.
func (ti *TaskIndex) SetTask(key string, t Task) {
	if ti.tasks == nil {
		ti.tasks = make(map[string]Task)
	}
	ti.tasks[key] = t
}

// SetTasks bulk-inserts tasks, replacing the internal map.
func (ti *TaskIndex) SetTasks(tasks map[string]Task) {
	if ti.tasks == nil {
		ti.tasks = make(map[string]Task, len(tasks))
	}
	for k, v := range tasks {
		ti.tasks[k] = v
	}
}

// TasksMap returns the internal tasks map for iteration.
func (ti *TaskIndex) TasksMap() map[string]Task {
	return ti.tasks
}

// TaskCount returns the number of tasks.
func (ti *TaskIndex) TaskCount() int {
	return len(ti.tasks)
}

// NewTestIndex creates a TaskIndex for testing with the given tasks.
func NewTestIndex(feature string, tasks map[string]Task, statusEnum ...[]string) *TaskIndex {
	ti := NewTaskIndex(feature)
	if len(statusEnum) > 0 {
		ti.StatusEnum = statusEnum[0]
	}
	ti.SetTasks(tasks)
	return ti
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
	return os.MkdirAll(dir, 0o755)
}
