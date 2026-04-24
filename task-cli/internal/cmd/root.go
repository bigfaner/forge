// Package cmd provides the CLI commands for the task management tool.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "Task management CLI for Claude Code projects",
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
	rootCmd.AddCommand(claimCmd)
	rootCmd.AddCommand(recordCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(featureCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(verifyCompletionCmd)
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(allCompletedCmd)
	rootCmd.AddCommand(versionCmd)
}
