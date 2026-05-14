package cmd

import (
	"fmt"
	"os"

	e2e "forge-cli/pkg/e2e"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var e2eVerifyFeature string

var e2eVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Check for unresolved VERIFY markers",
	Long: `Scan generated e2e test files for unresolved VERIFY markers that indicate
placeholder assertions needing human review.`,
	Run: runE2EVerify,
}

func init() {
	e2eVerifyCmd.Flags().StringVar(&e2eVerifyFeature, "feature", "", "Feature slug to verify (required)")
	_ = e2eVerifyCmd.MarkFlagRequired("feature")
}

func runE2EVerify(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	opts := e2e.RunOpts{
		ProjectRoot: projectRoot,
		Feature:     e2eVerifyFeature,
	}

	if err := e2e.Verify(opts); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
