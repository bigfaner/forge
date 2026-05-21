package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"forge-cli/pkg/index"
)

// LoadState loads the task state from the given file path.
// Returns nil if file doesn't exist.
func LoadState(path string) (*TaskState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read state: %w", err)
	}
	var state TaskState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state: %w", err)
	}
	return &state, nil
}

// SaveState saves the task state to the given file path using atomic write (temp+rename).
func SaveState(path string, state *TaskState) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}
	data = append(data, '\n')
	return index.AtomicWrite(path, data, 0644)
}

// DeleteState removes the state file.
func DeleteState(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
