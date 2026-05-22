package cmd

import (
	"fmt"

	e2e "forge-cli/pkg/e2e"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var e2eDiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "List all e2e test cases without running",
	Long: `List all discovered e2e test cases for the current feature without running them.
	Uses the project's discovery mechanism (e.g. npx playwright test --list,
	go test -list, python -m pytest --collect-only).`,
	Args: cobra.NoArgs,
	RunE: runE2EDiscover,
}

func runE2EDiscover(_ *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	if err := e2e.Discover(projectRoot); err != nil {
		return fmt.Errorf("e2e discover: %w", err)
	}
	return nil
}
