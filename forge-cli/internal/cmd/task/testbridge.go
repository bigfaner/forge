// Package task exports internal symbols for use by cmd package tests.
//
// This file exposes package internals for cross-package testing.
// Named export_for_test.go (not export_test.go) so it's compiled into
// the regular package (not just the test binary). This allows cmd tests
// to import these symbols via the task package.
//
// Functions with substantial logic have been migrated to pkg/task/.
// Where possible, exports point directly to pkg/task implementations.
// Remaining exports reference internal functions that depend on the
// cmd layer (cobra flags, base error types, etc.) and cannot be moved.
package task

import (
	"forge-cli/pkg/task"
)

// Suppress unused import — re-exported functions below keep the import alive,
// but this ensures compile safety even if all re-exports are temporarily removed.
var _ = task.ValidTypes

// --- Functions migrated to pkg/task/ ---
// These export aliases point directly to pkg/task implementations.

var (
	// ExportParseSegment delegates to pkg/task.ParseSegment for ID parsing.
	ExportParseSegment = task.ParseSegment

	// ExportFillRecordTemplate delegates to pkg/task.RenderRecord for category-based record rendering.
	ExportFillRecordTemplate = task.RenderRecord

	// ExportReadSubmitData delegates to pkg/task.ReadSubmitData for file/stdin IO.
	ExportReadSubmitData = task.ReadSubmitData
)

// --- Internal-only functions (cannot migrate due to cmd-layer dependencies) ---
// These reference functions defined in other files within this package.

// Exported function aliases for cross-package testing (cmd/integration_test.go).
var (
	// Command RunE functions
	ExportRunSubmit        = runSubmit
	ExportExecuteClaim     = executeClaim
	ExportClaimNextTask    = claimNextTask
	ExportRunCheckDeps     = runCheckDeps
	ExportRunValidateIndex = runValidateIndex
	ExportRunClaim         = runClaim
	ExportRunAdd           = runAdd
	ExportDoReopen         = doReopen

	// Utility functions (internal dependencies)
	ExportSaveIndexAndSignalCompletion = saveIndexAndSignalCompletion
	ExportValidateRecordData           = validateRecordData
	ExportPrintTaskDetails             = printTaskDetails
	ExportValidateQualityGate          = validateQualityGate
)

// ExportValidator exposes the validator type for cross-package testing.
type ExportValidator = validator

// NewExportValidator creates a validator for cross-package testing.
func NewExportValidator(filePath string) *ExportValidator {
	return &validator{filePath: filePath}
}

// ExportValidatorRun calls the unexported run method.
func (v *ExportValidator) ExportValidatorRun() error {
	return v.run()
}

// ExportValidateFirstTestTaskTemplate calls the unexported method.
func (v *ExportValidator) ExportValidateFirstTestTaskTemplate(taskFile, taskID string, placeholders []string) {
	v.validateFirstTestTaskTemplate(taskFile, taskID, placeholders)
}

// ExportErrors returns the validator errors.
func (v *ExportValidator) ExportErrors() []string {
	return v.errors
}

// ExportWarnings returns the validator warnings.
func (v *ExportValidator) ExportWarnings() []string {
	return v.warnings
}

// StatusCmd exposes statusCmd for cross-package testing.
var StatusCmd = statusCmd

// ExportSubmitCmd exposes submitCmd for cross-package testing.
var ExportSubmitCmd = submitCmd

// Flag variable pointers for cross-package testing.
var (
	ExportSubmitDataPath = &submitDataPath
	ExportSubmitJSON     = &submitJSON
	ExportSubmitQuiet    = &submitQuiet
)

// ExportExecuteAdd exposes executeAdd for cross-package testing.
var ExportExecuteAdd = executeAdd

// ExportPrintNewTask exposes printNewTask for cross-package testing.
var ExportPrintNewTask = printNewTask

// ExportPrintContinueTask exposes printContinueTask for cross-package testing.
var ExportPrintContinueTask = printContinueTask

// Add flag variable pointers for cross-package testing.
var (
	ExportAddTitle         = &addTitle
	ExportAddID            = &addID
	ExportAddPriority      = &addPriority
	ExportAddDependsOn     = &addDependsOn
	ExportAddEstimatedTime = &addEstimatedTime
	ExportAddBreaking      = &addBreaking
	ExportAddDescription   = &addDescription
	ExportAddVars          = &addVars
	ExportAddSourceTaskID  = &addSourceTaskID
	ExportAddBlockSource   = &addBlockSource
	ExportAddType          = &addType
)
