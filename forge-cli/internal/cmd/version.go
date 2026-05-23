package cmd

import (
	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the CLI version",
	Long:  `Print the version number of the task CLI.`,
	Args:  cobra.NoArgs,
	RunE:  runVersion,
}

func runVersion(_ *cobra.Command, _ []string) error {
	base.PrintBlock("VERSION", version.GetVersion())
	return nil
}
