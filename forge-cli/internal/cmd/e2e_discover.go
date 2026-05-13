package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var e2eDiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "List all e2e test cases without running",
	Long: `List all discovered e2e test cases for the current feature without running them.
Uses the active profile's discovery mechanism (e.g. npx playwright test --list,
go test -list, python -m pytest --collect-only).`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Fprintln(os.Stderr, "not yet implemented: forge e2e discover")
		os.Exit(1)
	},
}
