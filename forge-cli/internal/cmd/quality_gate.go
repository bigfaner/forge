package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"forge-cli/pkg/e2eprobe"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/just"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"
	tmpl "forge-cli/pkg/template"
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
			If all done: runs project-wide unit/integration tests, then e2e regression.

			Feature e2e tests are run by T-test-run (run-e2e-tests task), not this hook.
			This hook is the project health gate: unit tests + regression suite.

			Use -v to see why the command exits early (useful for debugging).`,
	Run: runQualityGate,
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
// Returns nil when tasks are not all done — caller should exit silently.
func checkAllCompleted(verbose bool) *AllCompletedResult {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Debugf(verbose, "project root not found: %v", err)
		return nil
	}
	Debugf(verbose, "project root: %s", projectRoot)

	featureSlug, err := feature.GetCurrentFeature(projectRoot)
	if err != nil {
		Debugf(verbose, "feature not found: %v", err)
		return nil
	}
	Debugf(verbose, "feature: %s", featureSlug)

	// Guard: only proceed if .forge/state.json signals allCompleted.
	forgeState := feature.ReadForgeState(projectRoot)
	if forgeState == nil || !forgeState.AllCompleted {
		Debugf(verbose, "no forge state with allCompleted — skipping")
		return nil
	}
	Debugf(verbose, "forge state: feature=%s allCompleted=true", forgeState.Feature)

	// Consume the state — clear it before proceeding
	_ = feature.ClearForgeState(projectRoot)

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		Debugf(verbose, "index.json not found: %s (%v)", indexPath, err)
		return nil
	}
	Debugf(verbose, "loaded index: %d tasks", index.TaskCount())

	// All tasks must be completed or skipped (rejected does not count as done)
	for _, t := range index.TasksMap() {
		if t.Status != feature.StatusCompleted && t.Status != feature.StatusSkipped {
			Debugf(verbose, "task %s is %s — not all done", t.ID, t.Status)
			return nil
		}
	}

	return &AllCompletedResult{
		FeatureSlug: featureSlug,
		ProjectRoot: projectRoot,
		DocsOnly:    isDocsOnly(index),
	}
}

// isDocsOnly returns true if no task has a testable runtime behavior type.
// Docs-only features change only markdown files — no compile/test/e2e needed.
// Unlike needsTestPipeline in pkg/task, this checks ALL tasks including auto-generated ones.
func isDocsOnly(index *task.TaskIndex) bool {
	for _, t := range index.TasksMap() {
		if task.IsTestableType(t.Type) {
			return false
		}
	}
	return true
}

func runQualityGate(_ *cobra.Command, _ []string) {
	result := checkAllCompleted(qualityGateVerbose)
	if result == nil {
		os.Exit(0) // not all done is normal, exit silently
	}

	fmt.Fprintf(os.Stderr, "=== All tasks completed for feature: %s ===\n", result.FeatureSlug)

	// Docs-only features have no code changes — skip compile/test/e2e gates.
	if result.DocsOnly {
		fmt.Fprintln(os.Stderr, "Feature is docs-only — skipping quality gate (no implementation or fix tasks)")
		os.Exit(0)
	}

	// Warn if feature e2e scripts exist but haven't been promoted.
	e2eScriptsDir := feature.GetE2EStagingDir(result.ProjectRoot, result.FeatureSlug)
	markerPath := feature.GetE2EGraduatedMarker(result.ProjectRoot, result.FeatureSlug)
	if just.FileExists(e2eScriptsDir) && !just.FileExists(markerPath) {
		fmt.Fprintln(os.Stderr,
			"WARNING: feature e2e scripts exist but haven't been run or promoted.\n"+
				"  Add T-test-run (run-e2e-tests) and T-test-graduate (promote) to your task index,\n"+
				"  or run /run-e2e-tests and forge test promote <journey> manually.")
	}

	// Step 1: Quality gate (compile -> fmt -> lint)
	// Stops at first blocking failure.
	gateSteps := just.LintGateSequence()
	just.RunGate(result.ProjectRoot, "", gateSteps, func(step, output string) {
		fmt.Fprintf(os.Stderr, "ERROR: %s check failed\n", step)
		errorDocPath := "tests/results/unit-raw-output.txt"
		if output != "" {
			if err := writeUnitTestRawOutput(result.ProjectRoot, "=== "+step+" failure ===\n"+output); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: failed to write %s output: %v\n", step, err)
			}
		}
		fixID, fixErr := addFixTask(result.ProjectRoot, result.FeatureSlug, step, output, errorDocPath)
		if fixErr != nil {
			fmt.Fprintf(os.Stderr, "WARNING: %v\n", fixErr)
		}
		handleGateFailure(step, errorDocPath, fixID, just.ExtractConciseError(output, 5))
	})

	// Step 2: Project-wide unit/integration tests (with retry-once policy)
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
		handleGateFailure("unit-test", errorDocPath, unitFixID, just.ExtractConciseError(unitOutput, 5))
	}

	// Step 3: Full e2e regression (promoted scripts in tests/e2e/)
	if just.HasJustfile(result.ProjectRoot) && just.HasRecipe(result.ProjectRoot, "e2e-test") {
		e2eReady := true
		if just.HasRecipe(result.ProjectRoot, "e2e-setup") {
			fmt.Fprintln(os.Stderr, "--- Ensuring e2e dependencies (just e2e-setup) ---")
			setupOutput, setupSuccess := just.RunCapture(result.ProjectRoot, "just", "e2e-setup")
			if !setupSuccess {
				fmt.Fprintln(os.Stderr, "WARNING: e2e-setup failed; skipping e2e regression")
				fmt.Fprintln(os.Stderr, "  To retry manually: just e2e-setup && just e2e-test")
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
			if !e2eprobe.ProbeServers(result.ProjectRoot, "") {
				fmt.Fprintln(os.Stderr, "WARNING: e2e server health check failed; skipping e2e regression")
				fmt.Fprintln(os.Stderr, "  Start dev server and retry: just dev && just e2e-test")
				e2eReady = false
			}
		}
		if e2eReady {
			fmt.Fprintln(os.Stderr, "--- Running full e2e regression (just e2e-test) ---")
			regressionOutput, regSuccess := just.RunCapture(result.ProjectRoot, "just", "e2e-test")
			if !regSuccess {
				fmt.Fprintln(os.Stderr, "ERROR: e2e regression failed")
				errorDocPath := "tests/e2e/results/raw-output.txt"
				if regressionOutput != "" {
					if err := writeRegressionRawOutput(result.ProjectRoot, regressionOutput); err != nil {
						fmt.Fprintf(os.Stderr, "WARNING: failed to write raw-output.txt: %v\n", err)
					}
				}
				fixID, fixErr := addFixTask(result.ProjectRoot, result.FeatureSlug, "e2e-test", regressionOutput, errorDocPath)
				if fixErr != nil {
					fmt.Fprintf(os.Stderr, "WARNING: %v\n", fixErr)
				}
				handleGateFailure("e2e-test", errorDocPath, fixID, just.ExtractConciseError(regressionOutput, 5))
			}
		}
	}
}

// handleGateFailure prints the hook JSON block reason and exits.
// fixID is the ID returned by addFixTask; empty means task creation failed.
func handleGateFailure(step, errorDocPath, fixID, concise string) {
	action := "run `forge task add --template fix-task` to create one manually, then `forge task claim`"
	if fixID != "" {
		action = "run `forge task claim` to pick it up"
	}

	guide := map[string]string{
		"compile":   "fix compilation errors",
		"lint":      "fix lint errors",
		"unit-test": "fix failing tests",
		"e2e-test":  "fix failing e2e tests",
	}
	label := map[string]string{
		"compile":   "Project compilation",
		"lint":      "Lint check",
		"unit-test": "Unit tests",
		"e2e-test":  "E2e regression tests",
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
		fixMsg = fmt.Sprintf("Fix task %s added (P0, breaking)", fixID)
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
	os.Exit(0)
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
		if err := writeUnitTestRawOutput(projectRoot, combinedOutput); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to write unit test output: %v\n", err)
		}
	}

	fixID, fixErr := addFixTask(projectRoot, featureSlug, "unit-test", combinedOutput, errorDocPath)
	return false, fixID, fixErr
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

// countFixTasks counts ALL fix-tasks for a step regardless of status (completed +
// active + blocked + skipped). A fix-task is identified by having a non-empty
// SourceTaskID AND a title with the prefix "fix <step>:".
func countFixTasks(index *task.TaskIndex, step string) int {
	count := 0
	prefix := "fix " + step + ":"
	for _, t := range index.TasksMap() {
		if t.SourceTaskID != "" &&
			strings.HasPrefix(t.Title, prefix) {
			count++
		}
	}
	return count
}

// fixTypeFromStep returns the deterministic task type for a quality gate failure step.
// compile/test failures → TypeCodingFix, fmt/lint failures → TypeCodingCleanup.
func fixTypeFromStep(step string) string {
	switch step {
	case "compile", "unit-test", "e2e-test":
		return task.TypeCodingFix
	case "fmt", "lint":
		return task.TypeCodingCleanup
	default:
		return task.TypeCodingFix
	}
}

// addFixTask creates a fix task using the same internal API as `forge task add`.
// Mirrors executeAdd() from add.go: template defaults -> AddTask -> CreateTaskMarkdown -> EnsureForgeState.
// Returns (taskID, nil) on success.
// Returns ("", error) on failure: template not found, task add failure, markdown creation failure, or cap exceeded.
func addFixTask(projectRoot, featureSlug, step, output, errorDocPath string) (string, error) {
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

	sourceFiles := extractSourceFiles(output)

	testScript := "just " + step
	if step == "unit-test" {
		testScript = "just test"
	}

	title := fmt.Sprintf("fix %s: %s failure in quality gate", step, testScript)
	description := fmt.Sprintf(
		"Quality gate step `%s` failed during quality-gate hook.\n\n"+
			"Error output saved to: `%s`\n\n"+
			"Concise error:\n```\n%s\n```",
		testScript, errorDocPath, just.ExtractConciseError(output, 10),
	)

	// Build opts — Priority/Breaking/EstimatedTime intentionally hardcoded
	// (not read from template defaults) since this is a programmatic caller.
	// SourceTaskID uses step-scoped sentinel for cumulative counting;
	// Vars["SOURCE_TASK_ID"] diverges intentionally for template rendering.
	taskType := fixTypeFromStep(step)
	tmplName := "fix-task"
	if taskType == task.TypeCodingCleanup {
		tmplName = "cleanup-task"
	}

	opts := task.AddTaskOpts{
		Title:         title,
		Priority:      "P0",
		EstimatedTime: "30min",
		Breaking:      true,
		Description:   description,
		SourceTaskID:  "quality-gate:" + step,
		Template:      tmplName,
		Type:          taskType,
		Vars: map[string]string{
			"SOURCE_FILES":   sourceFiles,
			"TEST_SCRIPT":    testScript,
			"TEST_RESULTS":   errorDocPath,
			"SOURCE_TASK_ID": "N/A (project-wide gate)",
		},
	}

	tasksDir := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug))

	if _, err := tmpl.Get(opts.Template); err != nil {
		return "", fmt.Errorf("template %q not found: %w", opts.Template, err)
	}
	if defs, err := tmpl.GetDefaults(opts.Template); err == nil && defs.IDPrefix != "" {
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

	fmt.Fprintf(os.Stderr, "Fix task %s added (P0, breaking)\n", id)
	return id, nil
}
