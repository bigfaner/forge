// Package cmd provides the CLI commands for the forge CLI tool.
package cmd

import (
	"os"

	factpkg "forge-cli/internal/cmd/fact"
	featurepkg "forge-cli/internal/cmd/feature"
	forensicpkg "forge-cli/internal/cmd/forensic"
	promptpkg "forge-cli/internal/cmd/prompt"
	qualitygatepkg "forge-cli/internal/cmd/qualitygate"
	taskpkg "forge-cli/internal/cmd/task"
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
	factpkg.Register()
	featurepkg.Register()
	forensicpkg.Register()
	promptpkg.Register()
	taskpkg.Register()
	worktreepkg.Register()

	// Group parents (5)
	rootCmd.AddCommand(taskpkg.Cmd)
	rootCmd.AddCommand(forensicpkg.Cmd)
	rootCmd.AddCommand(promptpkg.Cmd)
	rootCmd.AddCommand(worktreepkg.Cmd)
	rootCmd.AddCommand(factpkg.Cmd)

	// Top-level commands
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(qualitygatepkg.QualityGateCmd)
	rootCmd.AddCommand(verifyTaskDoneCmd)
	rootCmd.AddCommand(featurepkg.Cmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(proposalCmd)
	rootCmd.AddCommand(lessonCmd)
	rootCmd.AddCommand(researchCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(claudeCmd)
	rootCmd.AddCommand(surfacesCmd)

	// Version is hidden from --help but accessible via `forge version`
	versionCmd.Hidden = true
}
