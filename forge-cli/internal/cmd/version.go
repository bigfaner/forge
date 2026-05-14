package cmd

import (
	"forge-cli/pkg/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the CLI version",
	Long:  `Print the version number of the task CLI.`,
	Run:   runVersion,
}

func runVersion(_ *cobra.Command, _ []string) {
	PrintBlock("VERSION", version.GetVersion())
}
