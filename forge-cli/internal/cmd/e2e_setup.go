package cmd

import (
	"fmt"
	"os"

	e2e "forge-cli/pkg/e2e"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var e2eSetupForce bool

var e2eSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Install e2e dependencies (idempotent)",
	Long: `Install external dependencies for the configured e2e test profile.
Idempotent: safe to run multiple times. Uses the active profile from
.forge/config.yaml to determine what to install (e.g. Playwright browsers,
Go test tools, pytest packages).`,
	Run: runE2ESetup,
}

func init() {
	e2eSetupCmd.Flags().BoolVar(&e2eSetupForce, "force", false, "Force reinstall dependencies")
}

func runE2ESetup(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	opts := e2e.RunOpts{
		ProjectRoot: projectRoot,
		Force:       e2eSetupForce,
	}

	if err := e2e.Setup(opts); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
