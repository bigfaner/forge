// Package types defines typed constants for enum-like values used throughout
// forge-cli (Status, SurfaceType, Priority). It is a leaf package with zero
// internal dependencies — only the Go standard library is imported.
package types

// Status represents the lifecycle state of a task.
type Status string

// Task status constants.
const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusBlocked    Status = "blocked"
	StatusSuspended  Status = "suspended"
	StatusSkipped    Status = "skipped"
	StatusRejected   Status = "rejected"
)

// AllStatuses returns all defined Status constants.
func AllStatuses() []Status {
	return []Status{
		StatusPending,
		StatusInProgress,
		StatusCompleted,
		StatusBlocked,
		StatusSuspended,
		StatusSkipped,
		StatusRejected,
	}
}

// IsTerminalStatus reports whether the given status is a terminal state
// from which no further transitions are expected. The terminal statuses
// are completed, skipped, and rejected.
func IsTerminalStatus(s Status) bool {
	return s == StatusCompleted || s == StatusSkipped || s == StatusRejected
}
