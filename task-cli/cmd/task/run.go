package main

import (
	"os"

	"task-cli/internal/cmd"
)

// Version returns the version of the CLI.
var Version = "dev"

// Name returns the name of the CLI.
var Name = "task"

// Run executes the main entry point. This function is testable.
func Run() {
	cmd.Execute()
}

// GetVersion returns the CLI version.
func GetVersion() string {
	return Version
}

// GetName returns the CLI name.
func GetName() string {
	return Name
}

// IsTestMode checks if running in test mode.
func IsTestMode() bool {
	return os.Getenv("GO_TEST") == "1"
}
