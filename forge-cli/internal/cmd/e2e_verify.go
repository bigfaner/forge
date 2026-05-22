package cmd

import (
	"fmt"

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
	Args: cobra.NoArgs,
	RunE: runE2EVerify,
}

func init() {
	e2eVerifyCmd.Flags().StringVar(&e2eVerifyFeature, "feature", "", "Feature slug to verify (required)")
	_ = e2eVerifyCmd.MarkFlagRequired("feature")
}

func runE2EVerify(_ *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	opts := e2e.RunOpts{
		ProjectRoot: projectRoot,
		Feature:     e2eVerifyFeature,
	}

	if err := e2e.Verify(opts); err != nil {
		return fmt.Errorf("e2e verify: %w", err)
	}
	return nil
}
