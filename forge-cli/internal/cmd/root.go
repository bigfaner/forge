// Package cmd provides the CLI commands for the forge CLI tool.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "forge",
	Short: "CLI tool for managing tasks in Claude Code projects",
	Long: `A unified CLI tool for managing tasks in Claude Code projects.

Supports the docs/features/<slug>/ directory structure for task management.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Group parents (5)
	rootCmd.AddCommand(taskCmd)
	rootCmd.AddCommand(e2eCmd)
	rootCmd.AddCommand(forensicCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(promptCmd)

	// Top-level commands
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(probeCmd)
	rootCmd.AddCommand(qualityGateCmd)
	rootCmd.AddCommand(verifyTaskDoneCmd)
	rootCmd.AddCommand(featureCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(proposalCmd)
	rootCmd.AddCommand(lessonCmd)
	rootCmd.AddCommand(initCmd)

	// Version is hidden from --help but accessible via `forge version`
	versionCmd.Hidden = true

	// Task group subcommands
	taskCmd.AddCommand(claimCmd)
	taskCmd.AddCommand(submitCmd)
	taskCmd.AddCommand(statusCmd)
	taskCmd.AddCommand(queryCmd)
	taskCmd.AddCommand(checkDepsCmd)
	taskCmd.AddCommand(validateIndexCmd)
	taskCmd.AddCommand(addCmd)
	taskCmd.AddCommand(indexCmd)
	taskCmd.AddCommand(migrateCmd)
	taskCmd.AddCommand(listTypesCmd)

	// E2E group subcommands (6 total)
	e2eCmd.AddCommand(validateSpecsCmd)
	e2eCmd.AddCommand(e2eRunCmd)
	e2eCmd.AddCommand(e2eSetupCmd)
	e2eCmd.AddCommand(e2eVerifyCmd)
	e2eCmd.AddCommand(e2eCompileCmd)
	e2eCmd.AddCommand(e2eDiscoverCmd)

	// Prompt group subcommands
	promptCmd.AddCommand(promptGetCmd)
}
