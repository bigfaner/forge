// Package feature provides feature-related utilities.
package feature

// Path constants for the features-based directory structure.
const (
	// FeaturesDir is the base directory for all features
	FeaturesDir = "docs/features"

	// ProcessDirName is the name of the process subdirectory within tasks
	ProcessDirName = "process"

	// StateFileName is the name of the task state file (was task-state.json)
	StateFileName = "state.json"

	// RecordFileName is the name of the in-progress record file
	RecordFileName = "record.json"
)

// File names and directory names within a feature directory
const (
	// IndexFileName is the name of the task index file
	IndexFileName = "index.json"

	// PRDFileName is the name of the PRD file
	PRDFileName = "prd.md"

	// DesignFileName is the name of the design file
	DesignFileName = "design.md"

	// TasksDirName is the name of the tasks subdirectory
	TasksDirName = "tasks"

	// RecordsDirName is the name of the records subdirectory
	RecordsDirName = "records"

	// TemplateFileName is the name of the task template
	TemplateFileName = "template.md"
)

// Task status values
const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusBlocked    = "blocked"
	StatusSkipped    = "skipped"
)

// Priority values
const (
	PriorityP0 = "P0"
	PriorityP1 = "P1"
	PriorityP2 = "P2"
)
