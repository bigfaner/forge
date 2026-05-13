package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var e2eVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Check for unresolved VERIFY markers",
	Long: `Scan generated e2e test files for unresolved VERIFY markers that indicate
placeholder assertions needing human review.

Requires --feature flag to scope the search.`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Fprintln(os.Stderr, "not yet implemented: forge e2e verify")
		os.Exit(1)
	},
}
