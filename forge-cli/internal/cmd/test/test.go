// Package test contains all forge test subcommand implementations.
//
// Commands are registered into the CLI tree via Register(), called from
// the parent cmd package during initialization.
package test

import (
	"fmt"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/project"
	"forge-cli/pkg/testrunner"

	"github.com/spf13/cobra"
)

// Cmd is the parent test command, exported for use by the cmd package.
var Cmd = &cobra.Command{
	Use:   "test",
	Short: "Testing utilities for forge projects",
	Long: `Testing utilities for forge projects.

Subcommands:
  promote <journey>   — promote a journey's @feature tags to @regression
  run-journey <name>  — run a single journey in isolated temp directory
  verify              — detect contract breakage against current code`,
	Args: cobra.NoArgs,
	RunE: runTestHelp,
}

var testRunJourneyCmd = &cobra.Command{
	Use:   "run-journey <journey-name>",
	Short: "Run a single journey in isolated temp directory",
	Long: `Run a single journey's advanced tests in an isolated temporary directory.

Runs just test from the project root with the journey name as filter.
The temp directory is cleaned up after execution, regardless of success or failure.

The journey name is used as part of the temp directory path for traceability.

Output is a structured block with journey name, result, duration, and any failures.`,
	Args: cobra.ExactArgs(1),
	RunE: runTestRunJourney,
}

// Register adds all test subcommands to Cmd.
func Register() {
	Cmd.AddCommand(testPromoteCmd)
	Cmd.AddCommand(testRunJourneyCmd)
	Cmd.AddCommand(testVerifyCmd)
}

func runTestHelp(_ *cobra.Command, _ []string) error {
	base.PrintBlockStart()
	base.PrintField("USAGE", "forge test <subcommand>")
	base.PrintField("SUBCOMMANDS", "promote, run-journey, verify")
	base.PrintField("HINT", "Run 'forge test verify' to check contracts, 'forge test promote <journey>' to graduate features")
	base.PrintBlockEnd()
	return nil
}

func runTestRunJourney(_ *cobra.Command, args []string) error {
	journeyName := args[0]

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		base.Exit(base.ErrProjectNotFound())
	}

	cfg := testrunner.ResolveJourneyExecutionConfig(projectRoot)

	// Create isolated work directory
	workDir, cleanup, err := testrunner.CreateJourneyWorkDir(projectRoot, journeyName)
	if err != nil {
		base.Exit(base.NewAIError(base.ErrValidation, "Failed to create journey work directory", err.Error(),
			"Check temp directory permissions", "forge test run-journey "+journeyName))
	}
	defer cleanup()

	// Execute the test command in isolation
	result := testrunner.ExecuteJourneyInIsolation(cfg, workDir, journeyName)

	// Output the result report
	fmt.Print(result.FormatReport())
	return nil
}
