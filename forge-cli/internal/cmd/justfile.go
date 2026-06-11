package cmd

import (
	scaffoldpkg "forge-cli/internal/cmd/scaffold"

	"github.com/spf13/cobra"
)

var justfileCmd = &cobra.Command{
	Use:   "justfile",
	Short: "Justfile management commands",
	Args:  cobra.NoArgs,
}

func init() {
	justfileCmd.AddCommand(scaffoldpkg.Cmd)
	rootCmd.AddCommand(justfileCmd)
}
