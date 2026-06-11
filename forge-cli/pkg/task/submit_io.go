package task

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ReadSubmitData reads task record data from a file path or stdin.
// If dataPath is non-empty, reads from the file; otherwise reads from stdin.
func ReadSubmitData(dataPath string) (*RecordData, error) {
	var data []byte
	var err error

	if dataPath != "" {
		data, err = os.ReadFile(dataPath)
	} else {
		stat, _ := os.Stdin.Stat()
		if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
			return nil, fmt.Errorf("no input: provide --data flag or pipe JSON to stdin")
		}
		data, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read record data: %w", err)
	}

	var rd RecordData
	if err := json.Unmarshal(data, &rd); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return &rd, nil
}
