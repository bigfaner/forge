// Package worktree contains all forge worktree subcommand implementations.
//
// Commands are registered into the CLI tree via Register(), called from
// the parent cmd package during initialization.
package worktree

import "github.com/spf13/cobra"

// Cmd is the parent worktree command, exported for use by the cmd package.
var Cmd = &cobra.Command{
	Use:   "worktree",
	Short: "Manage git worktrees for feature development",
	Long: `Manage git worktrees for parallel feature development.

Each worktree is created inside the project at .forge/worktrees/<slug> with a
branch named <slug>. Forge's feature auto-detection resolves the correct
feature from the worktree name.`,
	Args: cobra.NoArgs,
}

// Register adds all worktree subcommands to Cmd.
func Register() {
	Cmd.AddCommand(
		startCmd,
		listCmd,
		removeCmd,
		resumeCmd,
		pushCmd,
		statusCmd,
	)
}
