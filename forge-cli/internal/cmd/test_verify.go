package cmd

import (
	"fmt"

	"forge-cli/pkg/contract"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var testVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Detect contract breakage by comparing Contract specs against current code",
	Long: `Scan all Contract spec files (tests/<journey>/_contracts/*.md), re-collect
	the Fact Table from the current codebase, and compare each Output/State assertion
	against actual values. Reports broken contracts with dimension-level detail.

	Hard Rules:
	  - verify does not modify any files, only reads and reports
	  - Fact Table is freshly collected on each run (no cached snapshots)
	  - Zero false positives on unchanged contracts`,
	Args: cobra.NoArgs,
	RunE: runTestVerify,
}

func init() {
	testCmd.AddCommand(testVerifyCmd)
}

func runTestVerify(_ *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	collector := contract.RealFactCollector{}
	summary, err := contract.Verify(projectRoot, collector)
	if err != nil {
		return NewErrEvalParseFailure(err.Error())
	}

	if summary.Total == 0 {
		return NewErrContractUnverifiable("no contracts to verify")
	}

	fmt.Print(summary.FormatReport())

	if summary.Broken > 0 {
		return NewAIError(ErrValidation,
			"Broken contracts detected",
			fmt.Sprintf("%d contract(s) are broken", summary.Broken),
			"Fix the broken contracts to match implementation",
			"forge test verify")
	}
	return nil
}
