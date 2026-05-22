package cmd

import (
	"fmt"

	e2e "forge-cli/pkg/e2e"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var e2eSetupForce bool

var e2eSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Install e2e dependencies (idempotent)",
	Long: `Install external dependencies for e2e tests.
	Idempotent: safe to run multiple times. Determines what to install based
	on the project's test framework (e.g. Playwright browsers, Go test tools, pytest packages).`,
	Args: cobra.NoArgs,
	RunE: runE2ESetup,
}

func init() {
	e2eSetupCmd.Flags().BoolVar(&e2eSetupForce, "force", false, "Force reinstall dependencies")
}

func runE2ESetup(_ *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	opts := e2e.RunOpts{
		ProjectRoot: projectRoot,
		Force:       e2eSetupForce,
	}

	if err := e2e.Setup(opts); err != nil {
		return fmt.Errorf("e2e setup: %w", err)
	}
	return nil
}
