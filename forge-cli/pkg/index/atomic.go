package index

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveIndexAtomic writes data to a temp file then renames it to path.
// This ensures the target file is never in a partial state.
func SaveIndexAtomic(path string, data any) (retErr error) {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}
	content = append(content, '\n')

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	tmp, err := os.CreateTemp(dir, ".index.json.tmp.*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmp.Name()

	defer func() {
		if retErr != nil {
			_ = tmp.Close()
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmp.Write(content); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}
	return nil
}
