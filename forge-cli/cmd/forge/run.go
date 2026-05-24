package main

import (
	"fmt"
	"os"

	"forge-cli/internal/cmd"
	"forge-cli/pkg/prompt"
	"forge-cli/pkg/task"
)

// Run executes the main entry point. This function is testable.
func Run() {
	if err := prompt.ValidatePromptTemplates(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := task.ValidateAutogenTemplates(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cmd.Execute()
}
