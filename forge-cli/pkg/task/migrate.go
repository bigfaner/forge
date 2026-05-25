package task

import (
	"fmt"
)

// LegacyScopeError is returned by CheckLegacyScope when tasks with a legacy
// 'scope' field but no 'surface-key' are detected. It implements the error
// interface and carries the count of affected tasks for formatted messages.
type LegacyScopeError struct {
	Count int
}

func (e *LegacyScopeError) Error() string {
	return fmt.Sprintf(
		"migration required: found %d tasks with legacy 'scope' field but no 'surface-key' — run 'forge task migrate' or 'forge breakdown-tasks' to regenerate tasks",
		e.Count,
	)
}

// CheckLegacyScope scans tasks for the legacy 'scope' field without a
// corresponding 'surface-key'. If such tasks are found, it returns a
// *LegacyScopeError with an actionable message (BIZ-error-reporting-002).
// Returns nil when no legacy tasks are detected.
func CheckLegacyScope(tasks []Task) error {
	var count int
	for _, t := range tasks {
		if t.Scope != "" && t.SurfaceKey == "" {
			count++
		}
	}
	if count > 0 {
		return &LegacyScopeError{Count: count}
	}
	return nil
}
