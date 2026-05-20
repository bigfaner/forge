package cmd

import (
	"fmt"
	"os"

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
	Run: runE2EDiscover,
}

func runE2EDiscover(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	if err := e2e.Discover(projectRoot); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
