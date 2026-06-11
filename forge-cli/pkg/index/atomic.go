package index

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// AtomicWrite writes data to a temp file in the same directory then renames it
// to path. This ensures the target file is never in a partial state.
// The perm value is applied to the final file via the temp file's mode.
func AtomicWrite(path string, data []byte, perm os.FileMode) (retErr error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp.*")
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

	if err := os.Chmod(tmp.Name(), perm); err != nil {
		return fmt.Errorf("failed to set temp file permissions: %w", err)
	}

	if _, err := tmp.Write(data); err != nil {
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

// SaveIndexAtomic marshals data to JSON and writes it atomically to path.
// This ensures the target file is never in a partial state.
func SaveIndexAtomic(path string, data any) error {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}
	content = append(content, '\n')

	return AtomicWrite(path, content, 0o644)
}
