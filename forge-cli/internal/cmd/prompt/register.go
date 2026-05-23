// Package prompt contains all forge prompt subcommand implementations.
//
// Commands are registered into the CLI tree via Register(), called from
// the parent cmd package during initialization.
package prompt

import "github.com/spf13/cobra"

// Cmd is the parent prompt command, exported for use by the cmd package.
var Cmd = &cobra.Command{
	Use:   "prompt",
	Short: "Manage agent execution prompts",
	Long: `Synthesize and print agent execution prompts for tasks.

Subcommands:
  get-by-task-id   Synthesize the agent prompt for a task`,
	Args: cobra.NoArgs,
}

// Register adds all prompt subcommands to Cmd.
func Register() {
	Cmd.AddCommand(promptGetCmd)
}
