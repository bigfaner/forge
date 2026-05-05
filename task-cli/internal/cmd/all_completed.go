package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"task-cli/pkg/e2eprobe"
	"task-cli/pkg/feature"
	"task-cli/pkg/just"
	"task-cli/pkg/project"
	"task-cli/pkg/task"
	"task-cli/pkg/testrunner"

	"github.com/spf13/cobra"
)

var allCompletedVerbose bool

var allCompletedCmd = &cobra.Command{
	Use:   "all-completed",
	Short: "Check if all tasks are done, then run tests",
	Long: `Checks if every task in the current feature is completed or skipped.
	Exits 0 silently if any task is still pending, in_progress, or blocked (no-op).
	If all done: runs project-wide unit/integration tests, then e2e regression.

	Feature e2e tests are run by T-test-3 (run-e2e-tests task), not this hook.
	This hook is the project health gate: unit tests + regression suite.

	Use -v to see why the command exits early (useful for debugging).`,
	Run: runAllCompleted,
}

func init() {
	allCompletedCmd.Flags().BoolVarP(&allCompletedVerbose, "verbose", "v", false, "print debug info when exiting early")
}

// AllCompletedResult holds context for running tests after all tasks complete.
type AllCompletedResult struct {
	FeatureSlug string
	ProjectRoot string
	TestCommand string // empty if not set in index.json
}

// checkAllCompleted verifies all tasks are done and returns test context.
// Returns nil (no error) when tasks are not all done — caller should exit 1.
func checkAllCompleted(verbose bool) (*AllCompletedResult, error) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Debugf(verbose, "project root not found: %v", err)
		return nil, nil //nolint:nilerr
	}
	Debugf(verbose, "project root: %s", projectRoot)

	featureSlug, err := feature.GetCurrentFeature(projectRoot)
	if err != nil {
		Debugf(verbose, "feature not found: %v", err)
		return nil, nil //nolint:nilerr
	}
	Debugf(verbose, "feature: %s", featureSlug)

	// Guard: only proceed if .forge/state.json signals allCompleted.
	forgeState := feature.ReadForgeState(projectRoot)
	if forgeState == nil || !forgeState.AllCompleted {
		Debugf(verbose, "no forge state with allCompleted — skipping")
		return nil, nil
	}
	Debugf(verbose, "forge state: feature=%s allCompleted=true", forgeState.Feature)

	// Consume the state — clear it before proceeding
	feature.ClearForgeState(projectRoot)

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		Debugf(verbose, "index.json not found: %s (%v)", indexPath, err)
		return nil, nil //nolint:nilerr
	}
	Debugf(verbose, "loaded index: %d tasks", len(index.Tasks))

	// All tasks must be completed or skipped
	for _, t := range index.Tasks {
		if t.Status != feature.StatusCompleted && t.Status != feature.StatusSkipped {
			Debugf(verbose, "task %s is %s — not all done", t.ID, t.Status)
			return nil, nil
		}
	}

	return &AllCompletedResult{
		FeatureSlug: featureSlug,
		ProjectRoot: projectRoot,
		TestCommand: index.TestCommand,
	}, nil
}

func runAllCompleted(cmd *cobra.Command, args []string) {
	result, err := checkAllCompleted(allCompletedVerbose)
	if err != nil || result == nil {
		os.Exit(0) // not all done is normal, exit silently
	}

	fmt.Fprintf(os.Stderr, "=== All tasks completed for feature: %s ===\n", result.FeatureSlug)

	// Warn if feature e2e scripts exist but haven't been graduated.
	e2eScriptsDir := feature.GetE2EStagingDir(result.ProjectRoot, result.FeatureSlug)
	markerPath := feature.GetE2EGraduatedMarker(result.ProjectRoot, result.FeatureSlug)
	if just.FileExists(e2eScriptsDir) && !just.FileExists(markerPath) {
		fmt.Fprintln(os.Stderr,
			"WARNING: feature e2e scripts exist but haven't been run or graduated.\n"+
				"  Add T-test-3 (run-e2e-tests) and T-test-4 (graduate-tests) to your task index,\n"+
				"  or run /run-e2e-tests and /graduate-tests manually.")
	}

	// Step 1: Quality gate (compile → fmt → lint)
	gateSteps := just.LintGateSequence()
	just.RunGate(result.ProjectRoot, "", gateSteps, func(step, output string) {
		fmt.Fprintf(os.Stderr, "ERROR: %s check failed\n", step)
		if output != "" {
			if err := writeUnitTestRawOutput(result.ProjectRoot, "=== "+step+" failure ===\n"+output); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: failed to write %s output: %v\n", step, err)
			}
		}
		concise := just.ExtractConciseError(output, 5)
		testrunner.PrintHookJSON(map[string]any{
			"decision": "block",
			"reason":   fmt.Sprintf("%s check failed. Read tests/results/unit-raw-output.txt, fix errors, then re-run.\n%s", testrunner.Capitalize(step), concise),
		})
		os.Exit(0)
	})

	// Step 2: Project-wide unit/integration tests
	fmt.Fprintln(os.Stderr, "--- Running project-wide tests ---")
	unitOutput, unitSuccess := testrunner.RunProjectTests(result.ProjectRoot, result.TestCommand)
	if !unitSuccess {
		fmt.Fprintln(os.Stderr, "ERROR: unit tests failed")
		if unitOutput != "" {
			if err := writeUnitTestRawOutput(result.ProjectRoot, unitOutput); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: failed to write unit test output: %v\n", err)
			}
		}
		testrunner.PrintHookJSON(map[string]any{
			"decision": "block",
			"reason":   "Unit tests failed. Read tests/results/unit-raw-output.txt, fix failures, then re-run.",
		})
		os.Exit(0)
	}

	// Step 3: Full e2e regression (graduated scripts in tests/e2e/)
	if just.HasJustfile(result.ProjectRoot) && just.HasRecipe(result.ProjectRoot, "test-e2e") {
		e2eReady := true
		if just.HasRecipe(result.ProjectRoot, "e2e-setup") {
			fmt.Fprintln(os.Stderr, "--- Ensuring e2e dependencies (just e2e-setup) ---")
			setupOutput, setupSuccess := just.RunCapture(result.ProjectRoot, "just", "e2e-setup")
			if !setupSuccess {
				fmt.Fprintln(os.Stderr, "WARNING: e2e-setup failed; skipping e2e regression")
				fmt.Fprintln(os.Stderr, "  To retry manually: just e2e-setup && just test-e2e")
				if setupOutput != "" {
					if err := writeRegressionRawOutput(result.ProjectRoot, "=== e2e-setup failure ===\n"+setupOutput); err != nil {
						fmt.Fprintf(os.Stderr, "WARNING: failed to write setup output: %v\n", err)
					} else {
						fmt.Fprintln(os.Stderr, "  Setup output saved to tests/e2e/results/raw-output.txt")
					}
				}
				e2eReady = false
			}
		}
		if e2eReady {
			if !e2eprobe.ProbeServers(result.ProjectRoot) {
				fmt.Fprintln(os.Stderr, "WARNING: e2e server health check failed; skipping e2e regression")
				fmt.Fprintln(os.Stderr, "  Start dev server and retry: just dev && just test-e2e")
				e2eReady = false
			}
		}
		if e2eReady {
			fmt.Fprintln(os.Stderr, "--- Running full e2e regression (just test-e2e) ---")
			regressionOutput, regSuccess := just.RunCapture(result.ProjectRoot, "just", "test-e2e")
			if !regSuccess {
				fmt.Fprintln(os.Stderr, "ERROR: e2e regression failed")
				if regressionOutput != "" {
					if err := writeRegressionRawOutput(result.ProjectRoot, regressionOutput); err != nil {
						fmt.Fprintf(os.Stderr, "WARNING: failed to write raw-output.txt: %v\n", err)
					}
				}
				testrunner.PrintHookJSON(map[string]any{
					"decision": "block",
					"reason":   "e2e regression failed. Read tests/e2e/results/raw-output.txt, analyze failures, then use `task add --title \"Fix: ...\" --priority P0 --breaking` to create fix tasks.",
				})
				os.Exit(0)
			}
		}
	}
}
