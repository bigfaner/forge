package cmd

import (
	"fmt"
	"os"

	e2e "forge-cli/pkg/e2e"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var e2eRunFeature string

var e2eRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run e2e tests (profile-aware)",
	Long: `Run end-to-end tests for the current feature using the configured test profile.
Reads the active profile from .forge/config.yaml and dispatches to the
profile-specific test runner (e.g. npx playwright test, go test, pytest).`,
	Run: runE2ERun,
}

func init() {
	e2eRunCmd.Flags().StringVar(&e2eRunFeature, "feature", "", "Run tests for a specific feature (empty = all)")
}

func runE2ERun(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	opts := e2e.RunOpts{
		ProjectRoot: projectRoot,
		Feature:     e2eRunFeature,
	}

	if err := e2e.Run(opts); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
