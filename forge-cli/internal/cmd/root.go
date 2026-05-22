// Package cmd provides the CLI commands for the forge CLI tool.
package cmd

import (
	"os"

	featurepkg "forge-cli/internal/cmd/feature"
	forensicpkg "forge-cli/internal/cmd/forensic"
	taskpkg "forge-cli/internal/cmd/task"
	testpkg "forge-cli/internal/cmd/test"
	worktreepkg "forge-cli/internal/cmd/worktree"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "forge",
	Short: "CLI tool for managing tasks in Claude Code projects",
	Long: `A unified CLI tool for managing tasks in Claude Code projects.

Supports the docs/features/<slug>/ directory structure for task management.`,
	Args: cobra.NoArgs,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Initialize subcommand groups
	featurepkg.Register()
	forensicpkg.Register()
	taskpkg.Register()
	testpkg.Register()
	worktreepkg.Register()

	// Group parents (5)
	rootCmd.AddCommand(taskpkg.Cmd)
	rootCmd.AddCommand(forensicpkg.Cmd)
	rootCmd.AddCommand(testpkg.Cmd)
	rootCmd.AddCommand(promptCmd)
	rootCmd.AddCommand(worktreepkg.Cmd)

	// Top-level commands
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(qualityGateCmd)
	rootCmd.AddCommand(verifyTaskDoneCmd)
	rootCmd.AddCommand(featurepkg.Cmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(proposalCmd)
	rootCmd.AddCommand(lessonCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(claudeCmd)

	// Version is hidden from --help but accessible via `forge version`
	versionCmd.Hidden = true

	// Task group subcommands — registered via taskpkg.Register() above
	// Worktree group subcommands — registered via worktreepkg.Register() above

	// Prompt group subcommands
	promptCmd.AddCommand(promptGetCmd)
}
