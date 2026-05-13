package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var e2eRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run e2e tests (profile-aware)",
	Long: `Run end-to-end tests for the current feature using the configured test profile.
Reads the active profile from .forge/config.yaml and dispatches to the
profile-specific test runner (e.g. npx playwright test, go test, pytest).

Requires --feature flag to scope test discovery.`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Fprintln(os.Stderr, "not yet implemented: forge e2e run")
		os.Exit(1)
	},
}
