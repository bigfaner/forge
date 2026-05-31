// Package feature provides feature-related utilities.
package feature

import "forge-cli/pkg/types"

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

	// Subdirectory names within a feature directory
	PRDDirName    = "prd"
	DesignDirName = "design"
	UIDirName     = "ui"

	// PRD file names (source-prefixed to prevent naming collisions)
	PRDSpecFile        = "prd-spec.md"
	PRDUserStoriesFile = "prd-user-stories.md"
	PRDUIFunctionsFile = "prd-ui-functions.md"

	// Design file names
	TechDesignFile  = "tech-design.md"
	APIHandbookFile = "api-handbook.md"

	// UI design file names
	UIDesignFile = "ui-design.md"

	// Manifest file
	ManifestFileName = "manifest.md"

	// Tasks subdirectory names
	TasksDirName   = "tasks"
	RecordsDirName = "records"

	// Testing subdirectory names
	TestingDirName        = "testing"
	TestingResultsDirName = "testing/results"
	TestCasesFileName     = "testing/test-cases.md"

	// Test paths (flat structure)
	TestBaseDir    = "tests"
	TestResultsDir = "tests/results"
	TestConfigFile = "tests/config.yaml"

	// Test output file names (relative to TestResultsDir)
	TestOutputFileName     = "raw-output.txt"
	UnitTestOutputFileName = "unit-raw-output.txt"

	// Forge runtime directory (project-level)
	ForgeDir            = ".forge"
	ForgeStateFileName  = "state.json"
	ForgeConfigFileName = "config.yaml"

	// Template file
	TemplateFileName = "template.md"

	// Proposals directory
	ProposalBaseDir  = "docs/proposals"
	ProposalFileName = "proposal.md"
)

// Status is a type alias for types.Status, ensuring feature.Status and
// types.Status are the same type.
type Status = types.Status

// Priority is a type alias for types.Priority, ensuring feature.Priority and
// types.Priority are the same type.
type Priority = types.Priority

// Task status constants — re-exported from pkg/types.
const (
	StatusPending    Status = types.StatusPending
	StatusInProgress Status = types.StatusInProgress
	StatusCompleted  Status = types.StatusCompleted
	StatusBlocked    Status = types.StatusBlocked
	StatusSuspended  Status = types.StatusSuspended
	StatusSkipped    Status = types.StatusSkipped
	StatusRejected   Status = types.StatusRejected
)

// Priority values — re-exported from pkg/types.
const (
	PriorityP0 Priority = types.PriorityP0
	PriorityP1 Priority = types.PriorityP1
	PriorityP2 Priority = types.PriorityP2
)
