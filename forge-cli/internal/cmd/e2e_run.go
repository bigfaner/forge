package cmd

import (
	"fmt"

	e2e "forge-cli/pkg/e2e"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var e2eRunFeature string

var e2eRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run e2e tests",
	Long: `Run end-to-end tests for the current feature.
	Dispatches to the appropriate test runner (e.g. npx playwright test, go test, pytest).`,
	Args: cobra.NoArgs,
	RunE: runE2ERun,
}

func init() {
	e2eRunCmd.Flags().StringVar(&e2eRunFeature, "feature", "", "Run tests for a specific feature (empty = all)")
}

func runE2ERun(_ *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	opts := e2e.RunOpts{
		ProjectRoot: projectRoot,
		Feature:     e2eRunFeature,
	}

	if err := e2e.Run(opts); err != nil {
		return fmt.Errorf("e2e run: %w", err)
	}
	return nil
}
