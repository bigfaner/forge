package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/just"
	"forge-cli/pkg/project"
	"forge-cli/pkg/serverprobe"
	"forge-cli/pkg/task"
	"forge-cli/pkg/testrunner"

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

var qualityGateCmd = &cobra.Command{
	Use:   "quality-gate",
	Short: "Check if all tasks are done, then run tests",
	Long: `Checks if every task in the current feature is completed or skipped.
Exits 0 silently if any task is still pending, in_progress, or blocked (no-op).
If all done: runs project-wide unit tests, then test regression.

Feature test scripts are run by T-test-run (run-test task), not this hook.
This hook is the project health gate: unit tests + regression suite.

Use -v to see why the command exits early (useful for debugging).`,
	Args: cobra.NoArgs,
	RunE: runQualityGate,
}

func init() {
	qualityGateCmd.Flags().BoolVarP(&qualityGateVerbose, "verbose", "v", false, "print debug info when exiting early")
}

// AllCompletedResult holds context for running tests after all tasks complete.
type AllCompletedResult struct {
	FeatureSlug string
	ProjectRoot string
	DocsOnly    bool // true if no implementation or fix tasks exist
}

// checkAllCompleted verifies all tasks are done and returns test context.
// Returns (nil, nil) when tasks are not all done — caller should exit silently.
// Returns (nil, error) for infrastructure failures (no project, no feature).
func checkAllCompleted(verbose bool) (*AllCompletedResult, error) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Debugf(verbose, "project root not found: %v", err)
		return nil, base.ErrProjectNotFound()
	}
	Debugf(verbose, "project root: %s", projectRoot)

	featureSlug, err := feature.GetCurrentFeature(projectRoot)
	if err != nil {
		Debugf(verbose, "feature not found: %v", err)
		return nil, base.ErrFeatureNotSet()
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
	_ = feature.ClearForgeState(projectRoot)

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		Debugf(verbose, "index.json not found: %s (%v)", indexPath, err)
		return nil, nil
	}
	Debugf(verbose, "loaded index: %d tasks", index.TaskCount())

	// All tasks must be completed or skipped (rejected does not count as done)
	for _, t := range index.TasksMap() {
		if t.Status != feature.StatusCompleted && t.Status != feature.StatusSkipped {
			Debugf(verbose, "task %s is %s — not all done", t.ID, t.Status)
			return nil, nil
		}
	}

	return &AllCompletedResult{
		FeatureSlug: featureSlug,
		ProjectRoot: projectRoot,
		DocsOnly:    isDocsOnly(index),
	}, nil
}

// isDocsOnly returns true if no task has a testable runtime behavior type.
// Docs-only features change only markdown files — no compile/test needed.
// Unlike needsTestPipeline in pkg/task, this checks ALL tasks including auto-generated ones.
func isDocsOnly(index *task.TaskIndex) bool {
	for _, t := range index.TasksMap() {
		if task.IsTestableType(t.Type) {
			return false
		}
	}
	return true
}

func runQualityGate(_ *cobra.Command, _ []string) error {
	result, err := checkAllCompleted(qualityGateVerbose)
	if err != nil {
		base.Exit(err)
	}
	if result == nil {
		return nil // not all done is normal, exit silently
	}

	fmt.Fprintf(os.Stderr, "=== All tasks completed for feature: %s ===\n", result.FeatureSlug)

	// Docs-only features have no code changes — skip compile/test gates.
	if result.DocsOnly {
		fmt.Fprintln(os.Stderr, "Feature is docs-only — skipping quality gate (no implementation or fix tasks)")
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
		fmt.Fprintf(os.Stderr, "ERROR: %s check failed\n", step)
		errorDocPath := "tests/results/unit-raw-output.txt"
		if output != "" {
			if err := testrunner.WriteUnitTestRawOutput(result.ProjectRoot, "=== "+step+" failure ===\n"+output); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: failed to write %s output: %v\n", step, err)
			}
		}
		fixID, fixErr := addFixTask(result.ProjectRoot, result.FeatureSlug, step, output, errorDocPath)
		if fixErr != nil {
			fmt.Fprintf(os.Stderr, "WARNING: %v\n", fixErr)
		}
		gateBlockErr = handleGateFailure(step, errorDocPath, fixID, just.ExtractConciseError(output, 5), fixTypeFromStep(step) == task.TypeCodingFix)
	})
	if gateBlockErr != nil {
		os.Exit(0)
	}

	// Step 2: Project-wide unit tests (with retry-once policy)
	fmt.Fprintln(os.Stderr, "--- Running project-wide tests ---")
	unitPassed, unitFixID, unitErr := runUnitTestStep(
		result.ProjectRoot, result.FeatureSlug,
		testrunner.RunProjectTests,
	)
	if unitErr != nil {
		fmt.Fprintf(os.Stderr, "WARNING: %v\n", unitErr)
	}
	if !unitPassed {
		unitOutput := "" // output already written by runUnitTestStep
		errorDocPath := "tests/results/unit-raw-output.txt"
		if err := handleGateFailure("unit-test", errorDocPath, unitFixID, just.ExtractConciseError(unitOutput, 5), true); err != nil {
			os.Exit(0)
		}
	}

	// Step 3: Full test regression (test scripts in tests/)
	if err := runTestRegression(result.ProjectRoot, result.FeatureSlug); err != nil {
		os.Exit(0)
	}
	return nil
}

// runTestRegression runs the full test regression suite when a justfile with
// a test recipe is present. When surfaces are configured in .forge/config.yaml,
// it orchestrates per-surface lifecycle (dev→probe→test→teardown for web/api/mobile;
// test→teardown for cli/tui). Falls back to legacy behavior when no surfaces configured.
// Returns an error when a gate failure is detected, nil otherwise.
func runTestRegression(projectRoot, featureSlug string) error {
	if !just.HasJustfile(projectRoot) || !just.HasRecipe(projectRoot, "test") {
		return nil
	}

	// Detect surface types from config.
	surfaces, _ := forgeconfig.ReadSurfaces(projectRoot)
	surfaceTypes := forgeconfig.SurfaceTypes(surfaces)

	if len(surfaceTypes) == 0 {
		// No surfaces configured — fall back to legacy behavior.
		return runTestRegressionLegacy(projectRoot, featureSlug)
	}

	// Surface-aware orchestration: run lifecycle per surface type.
	return runTestRegressionSurface(projectRoot, featureSlug, surfaceTypes)
}

// runTestRegressionLegacy is the pre-surface-aware test regression logic.
// Runs test-setup (optional), serverprobe health check, then just test.
func runTestRegressionLegacy(projectRoot, featureSlug string) error {
	// Optional setup step — skip regression on failure.
	if just.HasRecipe(projectRoot, "test-setup") {
		fmt.Fprintln(os.Stderr, "--- Ensuring test dependencies (just test-setup) ---")
		setupOutput, setupSuccess := just.RunCapture(projectRoot, "just", "test-setup")
		if !setupSuccess {
			fmt.Fprintln(os.Stderr, "WARNING: test-setup failed; skipping test regression")
			fmt.Fprintln(os.Stderr, "  To retry manually: just test-setup && just test")
			if setupOutput != "" {
				if err := testrunner.WriteRegressionRawOutput(projectRoot, "=== test-setup failure ===\n"+setupOutput); err != nil {
					fmt.Fprintf(os.Stderr, "WARNING: failed to write setup output: %v\n", err)
				} else {
					fmt.Fprintln(os.Stderr, "  Setup output saved to tests/results/raw-output.txt")
				}
			}
			return nil
		}
	}

	// Health check — skip regression if servers aren't ready.
	if !serverprobe.ProbeServers(projectRoot, "") {
		fmt.Fprintln(os.Stderr, "WARNING: server health check failed; skipping test regression")
		fmt.Fprintln(os.Stderr, "  Start dev server and retry: just dev && just test")
		return nil
	}

	// Run the regression suite.
	fmt.Fprintln(os.Stderr, "--- Running full test regression (just test) ---")
	regressionOutput, regSuccess := just.RunCapture(projectRoot, "just", "test")
	if !regSuccess {
		fmt.Fprintln(os.Stderr, "ERROR: test regression failed")
		errorDocPath := "tests/results/raw-output.txt"
		if regressionOutput != "" {
			if err := testrunner.WriteRegressionRawOutput(projectRoot, regressionOutput); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: failed to write raw-output.txt: %v\n", err)
			}
		}
		fixID, fixErr := addFixTask(projectRoot, featureSlug, "test", regressionOutput, errorDocPath)
		if fixErr != nil {
			fmt.Fprintf(os.Stderr, "WARNING: %v\n", fixErr)
		}
		return handleGateFailure("test", errorDocPath, fixID, just.ExtractConciseError(regressionOutput, 5), true)
	}
	return nil
}

// runTestRegressionSurface orchestrates per-surface-type lifecycle sequences.
// For each unique surface type, runs the appropriate sequence:
//   - web/api: dev → probe → test → teardown (full lifecycle)
//   - mobile: dev → probe → test-setup → test → teardown (full lifecycle with mobile setup)
//   - cli/tui: test → teardown (simplified)
//
// Surfaces of the same type share a single lifecycle (dev/probe run once per type).
// Teardown is mandatory regardless of prior step success/failure.
func runTestRegressionSurface(projectRoot, featureSlug string, surfaceTypes []string) error {
	var lastErr error
	for _, surfaceType := range surfaceTypes {
		fmt.Fprintf(os.Stderr, "--- Running surface orchestration for %s ---\n", surfaceType)
		result := runSurfaceLifecycle(projectRoot, surfaceType)
		if !result.success {
			errorDocPath := "tests/results/raw-output.txt"
			if result.output != "" {
				if err := testrunner.WriteRegressionRawOutput(projectRoot, result.output); err != nil {
					fmt.Fprintf(os.Stderr, "WARNING: failed to write raw-output.txt: %v\n", err)
				}
			}
			fixID, fixErr := addFixTask(projectRoot, featureSlug, "test", result.output, errorDocPath)
			if fixErr != nil {
				fmt.Fprintf(os.Stderr, "WARNING: %v\n", fixErr)
			}
			lastErr = handleGateFailure("test", errorDocPath, fixID, just.ExtractConciseError(result.output, 5), true)
		}
	}
	return lastErr
}

// lifecycleResult holds the result of a surface lifecycle execution.
type lifecycleResult struct {
	success bool
	output  string
}

// needsFullLifecycle returns true for surface types that require dev→probe→test→teardown.
// cli and tui surfaces use the simplified test→teardown sequence.
func needsFullLifecycle(surfaceType string) bool {
	return surfaceType == "web" || surfaceType == "api" || surfaceType == "mobile"
}

// resolveRecipe attempts to find a surface-specific recipe (e.g., "web-dev"),
// falling back to the generic recipe (e.g., "dev") if not found.
// Returns the recipe name to use, or empty string if neither exists.
func resolveRecipe(projectRoot, surfaceType, genericRecipe string) string {
	specificRecipe := surfaceType + "-" + genericRecipe
	if just.HasRecipe(projectRoot, specificRecipe) {
		return specificRecipe
	}
	if just.HasRecipe(projectRoot, genericRecipe) {
		return genericRecipe
	}
	return ""
}

// runSurfaceLifecycle executes the per-surface lifecycle sequence.
// For web/api: dev → probe → test → teardown
// For mobile: dev → probe → mobile-test-setup → test → teardown
// For cli/tui: test → teardown
// Teardown always executes (via defer-like pattern).
func runSurfaceLifecycle(projectRoot, surfaceType string) lifecycleResult {
	full := needsFullLifecycle(surfaceType)

	// Phase 1: Dev (full lifecycle only)
	if full {
		recipe := resolveRecipe(projectRoot, surfaceType, "dev")
		if recipe != "" {
			fmt.Fprintf(os.Stderr, "  Starting dev server (just %s)...\n", recipe)
			output, success := just.RunCapture(projectRoot, "just", recipe)
			if !success {
				fmt.Fprintf(os.Stderr, "  ERROR: dev failed (just %s)\n", recipe)
				runTeardown(projectRoot, surfaceType)
				return lifecycleResult{success: false, output: output}
			}
		}
	}

	// Phase 2: Probe (full lifecycle only)
	if full {
		probeRecipe := resolveRecipe(projectRoot, surfaceType, "probe")
		if !probeWithRetry(projectRoot, probeRecipe, 3, 5*time.Second) {
			fmt.Fprintln(os.Stderr, "  ERROR: probe failed after retries")
			runTeardown(projectRoot, surfaceType)
			return lifecycleResult{success: false, output: "probe failed: server not responding after 3 retries"}
		}
	}

	// Phase 2b: Mobile test setup (mobile surfaces only)
	if surfaceType == "mobile" {
		setupRecipe := resolveRecipe(projectRoot, surfaceType, "test-setup")
		if setupRecipe != "" {
			fmt.Fprintf(os.Stderr, "  Running mobile test setup (just %s)...\n", setupRecipe)
			output, success := just.RunCapture(projectRoot, "just", setupRecipe)
			if !success {
				fmt.Fprintf(os.Stderr, "  ERROR: mobile-test-setup failed (just %s)\n", setupRecipe)
				runTeardown(projectRoot, surfaceType)
				return lifecycleResult{success: false, output: output}
			}
		}
	}

	// Phase 3: Test
	var result lifecycleResult
	testRecipe := resolveRecipe(projectRoot, surfaceType, "test")
	if testRecipe != "" {
		fmt.Fprintf(os.Stderr, "  Running tests (just %s)...\n", testRecipe)
		output, success := just.RunCapture(projectRoot, "just", testRecipe)
		result = lifecycleResult{success: success, output: output}
		if !success {
			fmt.Fprintln(os.Stderr, "  ERROR: test failed")
		}
	} else {
		result = lifecycleResult{success: true}
	}

	// Phase 4: Teardown (always)
	runTeardown(projectRoot, surfaceType)

	return result
}

// runTeardown executes the teardown recipe for a surface type.
// Errors are logged but never fail the lifecycle — teardown is best-effort cleanup.
func runTeardown(projectRoot, surfaceType string) {
	recipe := resolveRecipe(projectRoot, surfaceType, "teardown")
	if recipe != "" {
		fmt.Fprintf(os.Stderr, "  Running teardown (just %s)...\n", recipe)
		output, success := just.RunCapture(projectRoot, "just", recipe)
		if !success {
			fmt.Fprintf(os.Stderr, "  WARNING: teardown failed (just %s)\n", recipe)
			if output != "" {
				fmt.Fprintf(os.Stderr, "  %s\n", just.ExtractConciseError(output, 3))
			}
		}
	}
}

// probeWithRetry runs the probe recipe with the specified number of retries.
// Returns true if the probe succeeds within the retry count.
// Returns true (skip) if the probe recipe doesn't exist.
// interval is the delay between retries (0 for no delay, useful in tests).
func probeWithRetry(projectRoot, probeRecipe string, maxRetries int, interval time.Duration) bool {
	if probeRecipe == "" {
		return true // no probe recipe — skip
	}

	// Verify the recipe actually exists before retrying.
	if !just.HasRecipe(projectRoot, probeRecipe) {
		return true // recipe not found — skip
	}

	for attempt := range maxRetries {
		if attempt > 0 && interval > 0 {
			fmt.Fprintf(os.Stderr, "  Probe retry %d/%d (waiting %v)...\n", attempt+1, maxRetries, interval)
			time.Sleep(interval)
		}
		fmt.Fprintf(os.Stderr, "  Probing (just %s) attempt %d/%d...\n", probeRecipe, attempt+1, maxRetries)
		_, success := just.RunCapture(projectRoot, "just", probeRecipe)
		if success {
			fmt.Fprintln(os.Stderr, "  Probe succeeded")
			return true
		}
	}
	return false
}

// handleGateFailure prints the hook JSON block reason and returns an error
// signalling that the gate blocked. The caller (RunE handler) decides exit behavior.
// fixID is the ID returned by addFixTask; empty means task creation failed.
// breaking indicates whether the fix task blocks downstream tasks.
func handleGateFailure(step, errorDocPath, fixID, concise string, breaking bool) error {
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
	fmt.Fprintln(os.Stderr, "WARNING: unit tests failed on first attempt, retrying once...")
	retryOutput, retrySuccess := runTest(projectRoot)
	if retrySuccess {
		fmt.Fprintln(os.Stderr, "WARNING: unit tests passed on retry (transient failure)")
		return true, "", nil
	}

	// Both attempts failed — write combined output and create fix task.
	fmt.Fprintln(os.Stderr, "ERROR: unit tests failed (retried once, both attempts failed)")
	errorDocPath := "tests/results/unit-raw-output.txt"
	combinedOutput := fmt.Sprintf(
		"retried once, both attempts failed\n\n=== First attempt ===\n%s\n\n=== Retry attempt ===\n%s",
		unitOutput, retryOutput,
	)
	if combinedOutput != "" {
		if err := testrunner.WriteUnitTestRawOutput(projectRoot, combinedOutput); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to write unit test output: %v\n", err)
		}
	}

	fixID, fixErr := addFixTask(projectRoot, featureSlug, "unit-test", combinedOutput, errorDocPath)
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

// sourceFileRe matches source file paths followed by :line or :line:col patterns.
var sourceFileRe = regexp.MustCompile(`([\w][\w./-]*\.\w{1,10})(?::\d+){1,2}`)

// sourceExts is a whitelist of source code extensions for file extraction.
var sourceExts = map[string]bool{
	".go": true, ".ts": true, ".js": true, ".tsx": true, ".jsx": true,
	".py": true, ".rs": true, ".java": true, ".rb": true,
	".c": true, ".cpp": true, ".h": true, ".hpp": true,
	".css": true, ".scss": true, ".html": true, ".sql": true,
	".vue": true, ".svelte": true,
}

// extractSourceFiles parses error output and returns comma-separated file paths.
func extractSourceFiles(output string) string {
	seen := make(map[string]bool)
	var files []string
	for _, match := range sourceFileRe.FindAllStringSubmatch(output, -1) {
		path := strings.TrimPrefix(match[1], "./")
		if path == "" || seen[path] {
			continue
		}
		if !sourceExts[filepath.Ext(path)] {
			continue
		}
		seen[path] = true
		files = append(files, path)
	}

	if len(files) > 10 {
		files = files[:10]
	}
	if len(files) == 0 {
		return "See error output for affected files"
	}
	return strings.Join(files, ", ")
}

// countFixTasks counts active (non-terminal) fix-tasks for a step.
// A fix-task is identified by having a title with the prefix "fix <step>:".
// Terminal statuses (completed, rejected, skipped) are excluded from the count.
// This ensures the fix-task cap reflects work-in-progress only.
func countFixTasks(index *task.TaskIndex, step string) int {
	count := 0
	prefix := "fix " + step + ":"
	for _, t := range index.TasksMap() {
		if !strings.HasPrefix(t.Title, prefix) {
			continue
		}
		// Exclude terminal statuses
		if t.Status == "completed" || t.Status == "rejected" || t.Status == "skipped" {
			continue
		}
		count++
	}
	return count
}

// fixTypeFromStep returns the deterministic task type for a quality gate failure step.
// compile/test failures → TypeCodingFix, fmt/lint failures → TypeCodingCleanup.
func fixTypeFromStep(step string) string {
	switch step {
	case "compile", "unit-test", "test":
		return task.TypeCodingFix
	case "fmt", "lint":
		return task.TypeCodingCleanup
	default:
		return task.TypeCodingFix
	}
}

// addFixTask creates fix tasks grouped by test suite (directory) using the same
// internal API as `forge task add`. Source files are extracted from the output and
// grouped by directory. Each directory group becomes a separate fix-task, enabling
// parallel execution and bounded scope. Returns the first task ID on success.
// Returns ("", error) on failure: template not found, task add failure, markdown creation failure, or cap exceeded.
func addFixTask(projectRoot, featureSlug, step, output, errorDocPath string) (string, error) {
	sourceFiles := extractSourceFiles(output)

	// Group source files by directory for parallel execution.
	groups := groupFilesByDir(sourceFiles)

	// If no meaningful groups (e.g. "See error output" fallback), create a single task.
	if len(groups) == 0 {
		return addSingleFixTask(projectRoot, featureSlug, step, sourceFiles, output, errorDocPath)
	}

	// One task per directory group.
	var firstID string
	for _, group := range groups {
		id, err := addSingleFixTask(projectRoot, featureSlug, step, group, output, errorDocPath)
		if err != nil {
			// If a group fails (e.g. cap reached), return the error.
			// Already-created tasks remain in the index.
			return firstID, err
		}
		if firstID == "" {
			firstID = id
		}
	}
	return firstID, nil
}

// groupFilesByDir splits comma-separated source files into groups by directory.
// Files in the same directory stay in one group. Returns nil if files is empty
// or the fallback message.
func groupFilesByDir(files string) []string {
	if files == "" || strings.HasPrefix(files, "See error") {
		return nil
	}

	dirMap := make(map[string][]string)
	for _, f := range strings.Split(files, ",") {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		dir := filepath.Dir(f)
		dirMap[dir] = append(dirMap[dir], f)
	}

	if len(dirMap) <= 1 {
		// All files in one directory (or no files) — single group.
		return nil
	}

	var groups []string
	for _, dirFiles := range dirMap {
		groups = append(groups, strings.Join(dirFiles, ", "))
	}
	return groups
}

// addSingleFixTask creates a single fix task using the same internal API as `forge task add`.
// Mirrors executeAdd() from add.go: template defaults -> AddTask -> CreateTaskMarkdown -> EnsureForgeState.
// Returns (taskID, nil) on success.
// Returns ("", error) on failure: template not found, task add failure, markdown creation failure, or cap exceeded.
func addSingleFixTask(projectRoot, featureSlug, step, sourceFiles, output, errorDocPath string) (string, error) {
	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))

	// Check cap before creating a new fix-task.
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to load index for cap check: %v\n", err)
		// Proceed without cap check if index can't be loaded.
	} else {
		active := countFixTasks(index, step)
		if active >= maxFixTasksPerStep {
			fmt.Fprintf(os.Stderr, "max fix-tasks reached for %s, manual intervention required\n", step)
			return "", ErrMaxFixTasks
		}
	}

	// Surface inference with soft-failure policy.
	// When surfaces are not configured or no match is found, fix-task creation
	// proceeds with empty surface key/type.
	surfaceKey, surfaceType := inferSurface(projectRoot, sourceFiles)
	if surfaceKey == "" && surfaceType == "" {
		fmt.Fprintln(os.Stderr, "WARNING: surface inference failed: no surfaces configured or no match for source files. Run 'forge surfaces detect' to configure surfaces")
	}

	testScript := "just " + step

	title := fmt.Sprintf("fix %s: %s failure in quality gate", step, testScript)
	description := fmt.Sprintf(
		"Quality gate step `%s` failed during quality-gate hook.\n\n"+
			"Error output saved to: `%s`\n\n"+
			"Concise error:\n```\n%s\n```",
		testScript, errorDocPath, just.ExtractConciseError(output, 10),
	)

	// Build opts — Priority/Breaking/EstimatedTime sourced from template defaults
	// when available (dual-source truth: opts is authoritative, template provides defaults).
	// SourceTaskID is deliberately empty (project-wide gate has no source task).
	// Vars["SOURCE_TASK_ID"] is "N/A (project-wide gate)" for template rendering.
	taskType := fixTypeFromStep(step)

	// Derive Breaking and EstimatedTime from template defaults (dual-source truth).
	// For coding.cleanup (fmt/lint failures): Breaking=false, EstimatedTime="15min".
	// For coding.fix (compile/test failures): Breaking=true, EstimatedTime="30min".
	breaking := true
	estimatedTime := "30min"
	if defs, err := task.GetTaskTemplateDefaults(taskType); err == nil {
		breaking = defs.Breaking
		estimatedTime = defs.EstimatedTime
	}

	opts := task.AddTaskOpts{
		Title:         title,
		Priority:      "P0",
		EstimatedTime: estimatedTime,
		Breaking:      breaking,
		Description:   description,
		Template:      taskType,
		Type:          taskType,
		SurfaceKey:    surfaceKey,
		SurfaceType:   surfaceType,
		Vars: map[string]string{
			"SOURCE_FILES":   sourceFiles,
			"TEST_SCRIPT":    testScript,
			"TEST_RESULTS":   errorDocPath,
			"SOURCE_TASK_ID": "N/A (project-wide gate)",
		},
	}

	tasksDir := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug))

	if _, err := task.GetTaskTemplate(opts.Template); err != nil {
		return "", fmt.Errorf("template %q not found: %w", opts.Template, err)
	}
	if defs, err := task.GetTaskTemplateDefaults(opts.Template); err == nil && defs.IDPrefix != "" {
		opts.IDPrefix = defs.IDPrefix
	}

	id, err := task.AddTask(indexPath, opts)
	if err != nil {
		return "", fmt.Errorf("failed to add fix task: %w", err)
	}

	opts.ID = id

	if err := task.CreateTaskMarkdown(tasksDir, id+".md", opts); err != nil {
		return "", fmt.Errorf("failed to create fix task file %s: %w", id+".md", err)
	}

	if err := feature.EnsureForgeState(projectRoot, featureSlug); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to update .forge/state.json: %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "Fix task %s added (P0, breaking=%v)\n", id, breaking)
	return id, nil
}
