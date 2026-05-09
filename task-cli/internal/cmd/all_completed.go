package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"task-cli/pkg/e2eprobe"
	"task-cli/pkg/feature"
	"task-cli/pkg/just"
	"task-cli/pkg/project"
	"task-cli/pkg/task"
	tmpl "task-cli/pkg/template"
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
		addFixTask(result.ProjectRoot, result.FeatureSlug, step, output, errorDocPath)
		handleGateFailure(step, errorDocPath, just.ExtractConciseError(output, 5))
	})

	// Step 2: Project-wide unit/integration tests
	fmt.Fprintln(os.Stderr, "--- Running project-wide tests ---")
	unitOutput, unitSuccess := testrunner.RunProjectTests(result.ProjectRoot, result.TestCommand)
	if !unitSuccess {
		fmt.Fprintln(os.Stderr, "ERROR: unit tests failed")
		errorDocPath := "tests/results/unit-raw-output.txt"
		if unitOutput != "" {
			if err := writeUnitTestRawOutput(result.ProjectRoot, unitOutput); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: failed to write unit test output: %v\n", err)
			}
		}
		addFixTask(result.ProjectRoot, result.FeatureSlug, "unit-test", unitOutput, errorDocPath)
		handleGateFailure("unit-test", errorDocPath, just.ExtractConciseError(unitOutput, 5))
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
				errorDocPath := "tests/e2e/results/raw-output.txt"
				if regressionOutput != "" {
					if err := writeRegressionRawOutput(result.ProjectRoot, regressionOutput); err != nil {
						fmt.Fprintf(os.Stderr, "WARNING: failed to write raw-output.txt: %v\n", err)
					}
				}
				addFixTask(result.ProjectRoot, result.FeatureSlug, "test-e2e", regressionOutput, errorDocPath)
				handleGateFailure("test-e2e", errorDocPath, just.ExtractConciseError(regressionOutput, 5))
			}
		}
	}
}

// handleGateFailure prints the hook JSON block reason and exits.
// Each stage has a distinct reason with stage-specific guidance.
func handleGateFailure(step, errorDocPath, concise string) {
	var reason string
	switch step {
	case "compile":
		reason = fmt.Sprintf(
			"Project compilation failed in all-completed hook. A fix task has been added (P0, breaking) — run `task claim` to pick it up and fix compilation errors.\nError output: %s\n%s",
			errorDocPath, concise)
	case "lint":
		reason = fmt.Sprintf(
			"Lint check failed in all-completed hook. A fix task has been added (P0, breaking) — run `task claim` to pick it up and fix lint errors.\nError output: %s\n%s",
			errorDocPath, concise)
	case "unit-test":
		reason = fmt.Sprintf(
			"Unit tests failed in all-completed hook. A fix task has been added (P0, breaking) — run `task claim` to pick it up and fix failing tests.\nError output: %s\n%s",
			errorDocPath, concise)
	case "test-e2e":
		reason = fmt.Sprintf(
			"E2e regression tests failed in all-completed hook. A fix task has been added (P0, breaking) — run `task claim` to pick it up and fix failing e2e tests.\nError output: %s\n%s",
			errorDocPath, concise)
	default:
		reason = fmt.Sprintf(
			"%s check failed in all-completed hook. A fix task has been added (P0, breaking) — run `task claim` to pick it up and fix the issue.\nError output: %s\n%s",
			testrunner.Capitalize(step), errorDocPath, concise)
	}

	testrunner.PrintHookJSON(map[string]any{
		"decision": "block",
		"reason":   reason,
	})
	os.Exit(0)
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

// addFixTask creates a fix task using the same internal API as `task add`.
// Mirrors executeAdd() from add.go: template defaults → AddTask → CreateTaskMarkdown → EnsureForgeState.
func addFixTask(projectRoot, featureSlug, step, output, errorDocPath string) string {
	sourceFiles := extractSourceFiles(output)

	testScript := "just " + step
	if step == "unit-test" {
		testScript = "just test"
	}

	title := fmt.Sprintf("Fix: %s failure in all-completed quality gate", step)
	description := fmt.Sprintf(
		"Quality gate step `%s` failed during all-completed hook.\n\n"+
			"Error output saved to: `%s`\n\n"+
			"Concise error:\n```\n%s\n```",
		testScript, errorDocPath, just.ExtractConciseError(output, 10),
	)

	// Build opts — same as task add --template fix-task
	opts := task.AddTaskOpts{
		Title:         title,
		Priority:      "P0",
		EstimatedTime: "30min",
		Breaking:      true,
		Description:   description,
		Template:      "fix-task",
		Vars: map[string]string{
			"SOURCE_FILES":   sourceFiles,
			"TEST_SCRIPT":    testScript,
			"TEST_RESULTS":   errorDocPath,
			"SOURCE_TASK_ID": "N/A (project-wide gate)",
		},
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	tasksDir := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug))

	if _, err := tmpl.Get(opts.Template); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: template %q not found: %v\n", opts.Template, err)
		return ""
	}

	id, err := task.AddTask(indexPath, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to add fix task: %v\n", err)
		return ""
	}

	opts.ID = id

	if err := task.CreateTaskMarkdown(tasksDir, id+".md", opts); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to create fix task file: %v\n", err)
		return ""
	}

	if err := feature.EnsureForgeState(projectRoot, featureSlug); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to update .forge/state.json: %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "Fix task %s added (P0, breaking)\n", id)
	return id
}
