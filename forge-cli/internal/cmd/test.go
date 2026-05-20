package cmd

import (
	"fmt"

	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Testing utilities for forge projects",
	Long: `Testing utilities for forge projects.

Subcommands:
  promote <journey>   — promote a journey's @feature tags to @regression
  run-journey <name>  — run a single journey in isolated temp directory
  verify              — detect contract breakage against current code`,
	Args: cobra.NoArgs,
	Run:  runTestHelp,
}

var testPromoteCmd = &cobra.Command{
	Use:   "promote <journey-name>",
	Short: "Promote a journey's @feature tags to @regression",
	Long: `Promote a journey by replacing all @feature tags with @regression tags.

Before promoting, runs all tests for the journey. If any test fails,
the promotion is refused and a failure report is printed.

Tag lifecycle:
  @feature (newly generated, under validation) -> @regression (verified, regression)`,
	Args: cobra.ExactArgs(1),
	Run:  runTestPromote,
}

var testRunJourneyCmd = &cobra.Command{
	Use:   "run-journey <journey-name>",
	Short: "Run a single journey in isolated temp directory",
	Long: `Run a single journey's e2e tests in an isolated temporary directory.

Runs just e2e-test from the project root with the journey name as filter.
The temp directory is cleaned up after execution, regardless of success or failure.

The journey name is used as part of the temp directory path for traceability.

Output is a structured block with journey name, result, duration, and any failures.`,
	Args: cobra.ExactArgs(1),
	Run:  runTestRunJourney,
}

func init() {
	testCmd.AddCommand(testPromoteCmd)
	testCmd.AddCommand(testRunJourneyCmd)
}

func runTestHelp(_ *cobra.Command, _ []string) {
	PrintBlockStart()
	PrintField("USAGE", "forge test <subcommand>")
	PrintField("SUBCOMMANDS", "promote, run-journey, verify")
	PrintField("HINT", "Run 'forge test verify' to check contracts, 'forge test promote <journey>' to graduate features")
	PrintBlockEnd()
}

func runTestRunJourney(_ *cobra.Command, args []string) {
	journeyName := args[0]

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	cfg := resolveJourneyExecutionConfig(projectRoot)

	// Create isolated work directory
	workDir, cleanup, err := createJourneyWorkDir(projectRoot, journeyName)
	if err != nil {
		Exit(NewAIError(ErrValidation, "Failed to create journey work directory", err.Error(),
			"Check temp directory permissions", "forge test run-journey "+journeyName))
	}
	defer cleanup()

	// Execute the test command in isolation
	result := executeJourneyInIsolation(cfg, workDir, journeyName)

	// Output the result report
	fmt.Print(result.FormatReport())
}
