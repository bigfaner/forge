package cmd

import (
	"task-cli/pkg/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the CLI version",
	Long:  `Print the version number of the task CLI.`,
	Run:   runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	PrintBlock("VERSION", version.GetVersion())
}
