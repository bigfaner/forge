package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var e2eSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Install e2e dependencies (idempotent)",
	Long: `Install external dependencies for the configured e2e test profile.
Idempotent: safe to run multiple times. Uses the active profile from
.forge/config.yaml to determine what to install (e.g. Playwright browsers,
Go test tools, pytest packages).`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Fprintln(os.Stderr, "not yet implemented: forge e2e setup")
		os.Exit(1)
	},
}
