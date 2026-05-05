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

	// E2E graduation paths
	E2ETestsBaseDir   = "tests/e2e"
	E2EStagingDir     = "tests/e2e/features"
	E2EGraduatedDir   = "tests/e2e/.graduated"

	// Forge runtime directory (project-level)
	ForgeDir            = ".forge"
	ForgeStateFileName  = "state.json"

	// Template file
	TemplateFileName = "template.md"

	// Proposals directory
	ProposalBaseDir  = "docs/proposals"
	ProposalFileName = "proposal.md"
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
