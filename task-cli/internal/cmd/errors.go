// Package cmd provides structured error output utilities for AI-friendly error messages.
package cmd

import (
	"fmt"
	"os"
	"strings"
)

// ErrorCode represents a structured error code for AI-friendly error messages.
type ErrorCode string

const (
	// ErrNoProject indicates that no project root was found.
	ErrNoProject ErrorCode = "NO_PROJECT"
	// ErrNoFeature indicates that no feature is set.
	ErrNoFeature ErrorCode = "NO_FEATURE"
	// ErrInvalidInput indicates invalid input (arguments, flags, etc.).
	ErrInvalidInput ErrorCode = "INVALID_INPUT"
	// ErrNotFound indicates a resource not found (task, file, etc.).
	ErrNotFound ErrorCode = "NOT_FOUND"
	// ErrConflict indicates a conflict (dependencies, state, etc.).
	ErrConflict ErrorCode = "CONFLICT"
	// ErrValidation indicates validation failure.
	ErrValidation ErrorCode = "VALIDATION_ERROR"
)

// AIError represents a structured error with AI-friendly context.
type AIError struct {
	Code    ErrorCode
	Message string
	Cause   string   // What caused the error
	Hint   string   // How to fix
	Action  string   // Suggested next step
}

// Error implements the error interface.
func (e *AIError) Error() string {
	return e.Message
}

// NewAIError creates a new AI-friendly error.
func NewAIError(code ErrorCode, message, cause, hint, action string) *AIError {
	return &AIError{
		Code:    code,
		Message: message,
		Cause:   cause,
		Hint:   hint,
		Action:  action,
	}
}

// Exit prints the AI-friendly error and exits.
func Exit(err error) {
	if aiErr, ok := err.(*AIError); ok {
		printAIError(aiErr)
	} else {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	}
	os.Exit(1)
}

// printAIError prints the error in AI-friendly format.
func printAIError(err *AIError) {
	fmt.Fprintln(os.Stderr, "---")
	fmt.Fprintf(os.Stderr, "ERROR_CODE: %s\n", err.Code)
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Message)
	if err.Cause != "" {
		fmt.Fprintf(os.Stderr, "CAUSE: %s\n", err.Cause)
	}
	if err.Hint != "" {
		fmt.Fprintf(os.Stderr, "HINT: %s\n", err.Hint)
	}
	if err.Action != "" {
		fmt.Fprintf(os.Stderr, "ACTION: %s\n", err.Action)
	}
	fmt.Fprintln(os.Stderr, "---")
}

// --- Helper functions for common errors ---

// ErrProjectNotFound creates a project root not found error.
func ErrProjectNotFound() *AIError {
	return NewAIError(
		ErrNoProject,
		"Project root not found",
		"No .git directory or CLAUDE.md found in current or parent directories",
		"Run from a directory containing CLAUDE.md or .git",
		"cd /path/to/project && task <command>",
	)
}

// ErrFeatureNotSet creates a feature not set error.
func ErrFeatureNotSet() *AIError {
	return NewAIError(
		ErrNoFeature,
		"No feature set",
		"Feature context is required but not configured",
		"Set a feature first using: task feature <slug>",
		"task feature <feature-slug>",
	)
}

// ErrTaskNotFound creates a task not found error.
func ErrTaskNotFound(taskID string) *AIError {
	return NewAIError(
		ErrNotFound,
		fmt.Sprintf("Task not found: %s", taskID),
		"The task ID does not exist in index.json",
		"Verify the task ID is correct. Check available tasks with: task query",
		"task query --status pending",
	)
}

// ErrNoInput creates a no input error.
func ErrNoInput(details string) *AIError {
	return NewAIError(
		ErrInvalidInput,
		"No input provided",
		details,
		"Provide the required input",
		"Check command usage with: task <command> -h",
	)
}

// ErrInvalidJSON creates an invalid JSON error.
func ErrInvalidJSON(path, details string) *AIError {
	return NewAIError(
		ErrValidation,
		fmt.Sprintf("Invalid JSON in %s", path),
		details,
		"Ensure JSON is valid and matches expected schema",
		fmt.Sprintf("Fix JSON syntax in %s and retry", path),
	)
}

// ErrFileNotFound creates a file not found error.
func ErrFileNotFound(path string) *AIError {
	return NewAIError(
		ErrNotFound,
		fmt.Sprintf("File not found: %s", path),
		"The specified file does not exist",
		"Check file path is correct and file exists",
		"ls "+path,
	)
}

// ErrNoPendingTasks creates a no pending tasks error.
func ErrNoPendingTasks() *AIError {
	return NewAIError(
		ErrNotFound,
		"No pending tasks available",
		"All tasks are either in_progress or completed, or no tasks defined",
		"Add new tasks to docs/features/<slug>/tasks/index.json",
		"task query --status pending",
	)
}

// ErrDependenciesNotMet creates dependencies not met error.
func ErrDependenciesNotMet(taskID string, unmetDeps []string) *AIError {
	return NewAIError(
		ErrConflict,
		fmt.Sprintf("Dependencies not met for task %s", taskID),
		fmt.Sprintf("Unmet dependencies: %s", strings.Join(unmetDeps, ", ")),
		"Complete the dependency tasks first",
		fmt.Sprintf("task status %s completed", strings.Join(unmetDeps, " ")),
	)
}

// ErrDataIntegrity creates a data integrity error.
func ErrDataIntegrity(issues []string) *AIError {
	return NewAIError(
		ErrConflict,
		"Task data integrity issues detected",
		strings.Join(issues, "; "),
		"Fix data inconsistency manually or cleanup state",
		"task cleanup",
	)
}

// ErrInvalidStatus creates an invalid status error.
func ErrInvalidStatus(status string, validStatuses []string) *AIError {
	action := "task status <id> <valid-status>"
	if len(validStatuses) > 0 {
		action = fmt.Sprintf("task status <id> %s", validStatuses[0])
	}
	cause := "statusEnum is not defined in index.json"
	if len(validStatuses) > 0 {
		cause = fmt.Sprintf("Valid statuses: %s", strings.Join(validStatuses, ", "))
	}
	return NewAIError(
		ErrValidation,
		fmt.Sprintf("Invalid status: %s", status),
		cause,
		"Use one of the valid status values",
		action,
	)
}

// ErrMissingFields creates an error for missing required fields in record data.
func ErrMissingFields(missing []string) *AIError {
	return NewAIError(
		ErrValidation,
		fmt.Sprintf("Missing required fields: %s", strings.Join(missing, ", ")),
		fmt.Sprintf("The following fields are required but empty: %s", strings.Join(missing, ", ")),
		"Include all required fields in record.json",
		"See record.json schema: taskId, summary, keyDecisions, testsPassed, testsFailed, coverage, acceptanceCriteria",
	)
}

// WarnMissingFields prints a warning for recommended but non-required fields.
func WarnMissingFields(missing []string) {
	fmt.Fprintln(os.Stderr, "---")
	fmt.Fprintf(os.Stderr, "WARNING: Missing recommended fields: %s\n", strings.Join(missing, ", "))
	fmt.Fprintf(os.Stderr, "HINT: Include these fields for complete records. Record will still be saved.\n")
	fmt.Fprintln(os.Stderr, "---")
}

// ErrFeatureNotFound creates a feature not found error.
func ErrFeatureNotFound(slug string) *AIError {
	return NewAIError(
		ErrNotFound,
		fmt.Sprintf("Feature not found: %s", slug),
		"Feature directory does not exist",
		"Check feature slug is correct",
		"ls docs/features/",
	)
}
