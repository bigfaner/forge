package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"forge-cli/pkg/forgelog"
	"forge-cli/pkg/index"
	"forge-cli/pkg/types"
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
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}
	data = append(data, '\n')
	return index.AtomicWrite(path, data, 0o644)
}

// DeleteState removes the state file.
func DeleteState(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// CheckExistingTaskState checks for an existing task state and determines
// whether to continue an in-progress task or claim a new one.
// Parameters: _ is unused (reserved for projectRoot), index is the task index,
// statePath is the path to the state file.
// Returns: (shouldContinue, hasIssues, issues).
func CheckExistingTaskState(_ string, index *TaskIndex, statePath string) (bool, bool, []string) {
	state, err := LoadState(statePath)
	if err != nil {
		forgelog.Warn("Warning: failed to load task state: %v\n", err)
		return false, false, nil
	}
	if state == nil {
		return false, false, nil
	}

	t, exists := index.ByID(state.Key)
	if !exists {
		return false, true, []string{fmt.Sprintf("Task key '%s' not found in index.json", state.Key)}
	}

	switch t.Status {
	case types.StatusInProgress:
		return true, false, nil
	case types.StatusCompleted:
		fmt.Printf("Previous task '%s' is completed. Claiming new task...\n", t.Title)
		_ = DeleteState(statePath)
		return false, false, nil
	case types.StatusBlocked:
		fmt.Printf("Previous task '%s' is blocked. Claiming new task...\n", t.Title)
		_ = DeleteState(statePath)
		return false, false, nil
	case types.StatusSuspended:
		fmt.Printf("Previous task '%s' is suspended. Claiming new task...\n", t.Title)
		_ = DeleteState(statePath)
		return false, false, nil
	case types.StatusRejected:
		fmt.Printf("Previous task '%s' was rejected. Claiming new task...\n", t.Title)
		_ = DeleteState(statePath)
		return false, false, nil
	default:
		return false, true, []string{fmt.Sprintf("Task '%s' has unexpected status: %s", t.Title, t.Status)}
	}
}
