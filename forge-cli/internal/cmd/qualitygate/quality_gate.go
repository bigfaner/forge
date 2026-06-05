package qualitygate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/forgelog"
	"forge-cli/pkg/just"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"
	"forge-cli/pkg/testrunner"
	"forge-cli/pkg/types"

	"github.com/spf13/cobra"
)

// testRunFunc is the signature for running project tests.
// Returns (output, success).
type testRunFunc func(projectRoot string) (string, bool)

var qualityGateVerbose bool

// maxFixTasksPerStep caps the number of active fix-tasks per quality-gate step.
// When this limit is reached, no new fix-tasks are created for that step.
const maxFixTasksPerStep = 3

// ErrMaxFixTasks is returned when the maximum number of active fix-tasks for a
// quality-gate step has been reached.
var ErrMaxFixTasks = errors.New("max fix-tasks reached")

// QualityGateCmd is the quality gate command.
// QualityGateCmd is the cobra command for running quality gate checks.
var QualityGateCmd = &cobra.Command{
	Use:   "quality-gate",
	Short: "Check if all tasks are done, then run tests",
	Long: `Checks if every task in the current feature is completed or skipped.
Exits 0 silently if any task is still pending, in_progress, or blocked (no-op).
If all done: runs project-wide unit tests, then test regression.

Feature test scripts are run by T-test-run (run-test task), not this hook.
This hook is the project health gate: unit tests + regression suite.

Use -v to see why the command exits early (useful for debugging).`,
	Args: cobra.NoArgs,
	RunE: RunQualityGate,
}

func init() {
	QualityGateCmd.Flags().BoolVarP(&qualityGateVerbose, "verbose", "v", false, "print debug info when exiting early")
}

// AllCompletedResult holds context for running tests after all tasks complete.
type AllCompletedResult struct {
	FeatureSlug string
	ProjectRoot string
	DocsOnly    bool // true if no implementation or fix tasks exist
}

// findProjectRoot is the function used to locate the project root.
// Defaults to project.FindProjectRoot; replaced in tests for reliable isolation.
var findProjectRoot = project.FindProjectRoot

// CheckAllCompleted verifies all tasks are done and returns test context.
// Returns (nil, nil) when tasks are not all done or no feature is set — caller should exit silently.
// Returns (nil, error) for infrastructure failures (no project).
func CheckAllCompleted(verbose bool) (*AllCompletedResult, error) {
	projectRoot, err := findProjectRoot()
	if err != nil {
		base.Debugf(verbose, "project root not found: %v", err)
		return nil, base.ErrProjectNotFound()
	}
	base.Debugf(verbose, "project root: %s", projectRoot)

	featureSlug, err := feature.GetCurrentFeature(projectRoot)
	if err != nil {
		base.Debugf(verbose, "feature not found: %v", err)
		return nil, nil
	}
	base.Debugf(verbose, "feature: %s", featureSlug)

	// Guard: only proceed if .forge/state.json signals allCompleted.
	forgeState := feature.ReadForgeState(projectRoot)
	if forgeState == nil || !forgeState.AllCompleted {
		base.Debugf(verbose, "no forge state with allCompleted — skipping")
		return nil, nil
	}
	base.Debugf(verbose, "forge state: feature=%s allCompleted=true", forgeState.Feature)

	// Consume the state — clear it before proceeding
	_ = feature.ClearForgeState(projectRoot)

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		base.Debugf(verbose, "index.json not found: %s (%v)", indexPath, err)
		return nil, nil
	}
	base.Debugf(verbose, "loaded index: %d tasks", index.TaskCount())

	// All tasks must be completed or skipped (rejected does not count as done)
	for _, t := range index.TasksMap() {
		if t.Status != types.StatusCompleted && t.Status != types.StatusSkipped {
			base.Debugf(verbose, "task %s is %s — not all done", t.ID, t.Status)
			return nil, nil
		}
	}

	return &AllCompletedResult{
		FeatureSlug: featureSlug,
		ProjectRoot: projectRoot,
		DocsOnly:    IsDocsOnly(index),
	}, nil
}

// IsDocsOnly returns true if no task has a testable runtime behavior type.
// Docs-only features change only markdown files — no compile/test needed.
// Unlike needsTestPipeline in pkg/task, this checks ALL tasks including auto-generated ones.
func IsDocsOnly(index *task.TaskIndex) bool {
	for _, t := range index.TasksMap() {
		if task.IsTestableType(t.Type) {
			return false
		}
	}
	return true
}

// RunQualityGate executes the quality gate check.
// RunQualityGate executes the quality gate check.
func RunQualityGate(_ *cobra.Command, _ []string) error {
	result, err := CheckAllCompleted(qualityGateVerbose)
	if err != nil {
		base.Exit(err)
	}
	if result == nil {
		return nil // not all done is normal, exit silently
	}

	forgelog.Info("=== All tasks completed for feature: %s ===\n", result.FeatureSlug)

	// Docs-only features have no code changes — skip compile/test gates.
	if result.DocsOnly {
		forgelog.Info("Feature is docs-only — skipping quality gate (no implementation or fix tasks)\n")
		os.Exit(0)
	}

	// NOTE: The legacy promotion model was removed in v3.0.0 in favor of
	// tag-based test promotion. No pre-run validation is needed here —
	// runTestRegression handles the full lifecycle per surface type.
	_ = result.ProjectRoot
	_ = result.FeatureSlug

	// Step 1: Quality gate (compile -> fmt -> lint)
	// Stops at first blocking failure.
	gateSteps := just.NonBreakingGateSequence()
	var gateBlockErr error
	just.RunGate(result.ProjectRoot, "", gateSteps, func(step, output string) {
		forgelog.Error("ERROR: %s check failed\n", step)
		errorDocPath := feature.TestResultsDir + "/" + feature.UnitTestOutputFileName
		if output != "" {
			if err := testrunner.WriteUnitTestRawOutput(result.ProjectRoot, "=== "+step+" failure ===\n"+output); err != nil {
				forgelog.Warn("WARNING: failed to write %s output: %v\n", step, err)
			}
		}
		fixID, fixErr := AddFixTask(result.ProjectRoot, result.FeatureSlug, step, output, errorDocPath)
		if fixErr != nil {
			forgelog.Warn("WARNING: %v\n", fixErr)
		}
		gateBlockErr = HandleGateFailure(step, errorDocPath, fixID, just.ExtractConciseError(output, conciseErrorMaxLines), fixTypeFromStep(step) == task.TypeCodingFix)
	})
	if gateBlockErr != nil {
		os.Exit(0)
	}

	// Step 2: Project-wide unit tests (with retry-once policy)
	forgelog.Info("--- Running project-wide tests ---\n")
	unitPassed, unitFixID, unitErr := runUnitTestStep(
		result.ProjectRoot, result.FeatureSlug,
		testrunner.RunProjectTests,
	)
	if unitErr != nil {
		forgelog.Warn("WARNING: %v\n", unitErr)
	}
	if !unitPassed {
		unitOutput := "" // output already written by runUnitTestStep
		errorDocPath := feature.TestResultsDir + "/" + feature.UnitTestOutputFileName
		if err := HandleGateFailure("unit-test", errorDocPath, unitFixID, just.ExtractConciseError(unitOutput, conciseErrorMaxLines), true); err != nil {
			os.Exit(0)
		}
	}

	// Step 3: Full test regression (test scripts in tests/)
	if err := runTestRegression(result.ProjectRoot, result.FeatureSlug); err != nil {
		os.Exit(0)
	}
	return nil
}

// HandleGateFailure prints the hook JSON block reason and returns an error
// signalling that the gate blocked. The caller (RunE handler) decides exit behavior.
// fixID is the ID returned by addFixTask; empty means task creation failed.
// breaking indicates whether the fix task blocks downstream tasks.
func HandleGateFailure(step, errorDocPath, fixID, concise string, breaking bool) error {
	action := "run `forge task add --type coding.fix` to create one manually, then `forge task claim`"
	if fixID != "" {
		action = "run `forge task claim` to pick it up"
	}

	guide := map[string]string{
		"compile":   "fix compilation errors",
		"lint":      "fix lint errors",
		"unit-test": "fix failing unit tests",
		"test":      "fix failing tests",
	}
	label := map[string]string{
		"compile":   "Project compilation",
		"lint":      "Lint check",
		"unit-test": "Unit tests",
		"test":      "Advanced tests",
	}

	g := guide[step]
	if g == "" {
		g = "fix the issue"
	}
	l := label[step]
	if l == "" {
		l = testrunner.Capitalize(step) + " check"
	}

	var fixMsg string
	if fixID != "" {
		fixMsg = fmt.Sprintf("Fix task %s added (P0, breaking=%v)", fixID, breaking)
	} else {
		fixMsg = "Failed to add fix task automatically"
	}

	reason := fmt.Sprintf(
		"%s failed in quality-gate hook. %s — %s and %s.\nError output: %s\n%s",
		l, fixMsg, action, g, errorDocPath, concise)

	testrunner.PrintHookJSON(map[string]any{
		"decision": "block",
		"reason":   reason,
	})
	return fmt.Errorf("quality gate blocked: %s", step)
}

// runUnitTestStep runs unit tests with a retry-once policy for transient failures.
// On first failure, tests are re-run once. If retry passes, a warning is logged and
// no fix task is created. If retry also fails, a fix task is created with both outputs.
// Returns (passed, fixTaskID, error).
func runUnitTestStep(projectRoot, featureSlug string, runTest testRunFunc) (bool, string, error) {
	unitOutput, unitSuccess := runTest(projectRoot)
	if unitSuccess {
		return true, "", nil
	}

	// First attempt failed — retry once.
	forgelog.Warn("WARNING: unit tests failed on first attempt, retrying once...\n")
	retryOutput, retrySuccess := runTest(projectRoot)
	if retrySuccess {
		forgelog.Warn("WARNING: unit tests passed on retry (transient failure)\n")
		return true, "", nil
	}

	// Both attempts failed — write combined output and create fix task.
	forgelog.Error("ERROR: unit tests failed (retried once, both attempts failed)\n")
	errorDocPath := feature.TestResultsDir + "/" + feature.UnitTestOutputFileName
	combinedOutput := fmt.Sprintf(
		"retried once, both attempts failed\n\n=== First attempt ===\n%s\n\n=== Retry attempt ===\n%s",
		unitOutput, retryOutput,
	)
	if combinedOutput != "" {
		if err := testrunner.WriteUnitTestRawOutput(projectRoot, combinedOutput); err != nil {
			forgelog.Warn("WARNING: failed to write unit test output: %v\n", err)
		}
	}

	fixID, fixErr := AddFixTask(projectRoot, featureSlug, "unit-test", combinedOutput, errorDocPath)
	return false, fixID, fixErr
}

// requireSurfaceInference wraps inferSurface with a hard-failure policy.
// When surface inference fails (returns empty key+type), it returns an error
// with guidance to run `forge surfaces detect`.
// This is the hard-constraint entry point used by addSingleFixTask.
// To revert to soft behavior, replace the call with inferSurface and use empty strings.
func requireSurfaceInference(projectRoot, sourceFiles string) (surfaceKey, surfaceType string, err error) {
	key, typ := inferSurface(projectRoot, sourceFiles)
	if key == "" && typ == "" {
		return "", "", fmt.Errorf("surface inference failed: no surfaces configured or no match for source files %q. Run 'forge surfaces detect' to configure surfaces", sourceFiles)
	}
	return key, typ, nil
}

// inferSurface attempts to determine the surface-key and surface-type for a
// fix-task by querying forge surfaces with all extracted source file paths.
// Returns ("", "") on any failure (no surfaces configured, no match, parse error)
// — the caller falls back to empty values and fix-task creation proceeds unblocked.
// Uses all source files (not just the first) to correctly handle multi-surface projects
// where different files may belong to different surfaces.
func inferSurface(projectRoot, sourceFiles string) (surfaceKey, surfaceType string) {
	surfaces, err := forgeconfig.ReadSurfaces(projectRoot)
	if err != nil || len(surfaces) == 0 {
		return "", ""
	}

	// Parse all source file paths from the comma-separated list.
	// sourceFiles may be "See error output for affected files" when no files found.
	if sourceFiles == "" || strings.HasPrefix(sourceFiles, "See error") {
		return "", ""
	}

	var files []string
	for _, part := range strings.Split(sourceFiles, ",") {
		f := strings.TrimSpace(part)
		if f != "" {
			files = append(files, f)
		}
	}
	if len(files) == 0 {
		return "", ""
	}

	// Try each file until one matches a configured surface.
	for _, file := range files {
		match, err := forgeconfig.MatchSurface(surfaces, file)
		if err == nil {
			return match.Key, match.Type
		}
	}
	return "", ""
}
