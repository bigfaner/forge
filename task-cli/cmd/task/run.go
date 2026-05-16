package main

import (
	"os"

	"task-cli/internal/cmd"
	"task-cli/pkg/version"
)

// Run executes the main entry point. This function is testable.
func Run() {
	cmd.Execute()
}

// GetVersion returns the CLI version.
func GetVersion() string {
	return version.GetVersion()
}

// GetName returns the CLI name.
func GetName() string {
	return version.GetName()
}

// IsTestMode checks if running in test mode.
func IsTestMode() bool {
	return os.Getenv("GO_TEST") == "1"
}
