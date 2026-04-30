package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"task-cli/pkg/feature"
	"task-cli/pkg/project"
	"task-cli/pkg/task"

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
	// Feature e2e execution is owned by T-test-3 (run-e2e-tests task).
	e2eScriptsDir := feature.GetE2ETargetDir(result.ProjectRoot, result.FeatureSlug)
	markerPath := feature.GetE2EGraduatedMarker(result.ProjectRoot, result.FeatureSlug)
	if fileExists(e2eScriptsDir) && !fileExists(markerPath) {
		fmt.Fprintln(os.Stderr,
			"WARNING: feature e2e scripts exist but haven't been run or graduated.\n"+
				"  Add T-test-3 (run-e2e-tests) and T-test-4 (graduate-tests) to your task index,\n"+
				"  or run /run-e2e-tests and /graduate-tests manually.")
	}

	// Step 1: Compile check (type errors before running tests)
	if hasJustfile(result.ProjectRoot) && hasJustRecipe(result.ProjectRoot, "compile") {
		fmt.Fprintln(os.Stderr, "--- Running compile check (just compile) ---")
		compileOutput, compileSuccess := runCmdCapture(result.ProjectRoot, "just", "compile")
		if !compileSuccess {
			fmt.Fprintln(os.Stderr, "ERROR: compile check failed")
			if compileOutput != "" {
				if err := writeUnitTestRawOutput(result.ProjectRoot, "=== compile failure ===\n"+compileOutput); err != nil {
					fmt.Fprintf(os.Stderr, "WARNING: failed to write compile output: %v\n", err)
				}
			}
			printHookJSON(map[string]any{
				"decision": "block",
				"reason":   "Compile check failed. Read tests/results/unit-raw-output.txt, fix type errors, then re-run.",
			})
			os.Exit(0)
		}
	}

	// Step 2: Project-wide unit/integration tests
	fmt.Fprintln(os.Stderr, "--- Running project-wide tests ---")
	unitOutput, unitSuccess := runProjectTests(result.ProjectRoot, result.TestCommand)
	if !unitSuccess {
		fmt.Fprintln(os.Stderr, "ERROR: unit tests failed")
		if unitOutput != "" {
			if err := writeUnitTestRawOutput(result.ProjectRoot, unitOutput); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: failed to write unit test output: %v\n", err)
			}
		}
		printHookJSON(map[string]any{
			"decision": "block",
			"reason":   "Unit tests failed. Read tests/results/unit-raw-output.txt, fix failures, then re-run.",
		})
		os.Exit(0)
	}

	// Step 3: Full e2e regression (graduated scripts in tests/e2e/)
	if hasJustfile(result.ProjectRoot) && hasJustRecipe(result.ProjectRoot, "test-e2e") {
		// Ensure e2e dependencies are installed; skip regression if setup fails (environment issue)
		e2eReady := true
		if hasJustRecipe(result.ProjectRoot, "e2e-setup") {
			fmt.Fprintln(os.Stderr, "--- Ensuring e2e dependencies (just e2e-setup) ---")
			setupOutput, setupSuccess := runCmdCapture(result.ProjectRoot, "just", "e2e-setup")
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
			fmt.Fprintln(os.Stderr, "--- Running full e2e regression (just test-e2e) ---")
			regressionOutput, regSuccess := runCmdCapture(result.ProjectRoot, "just", "test-e2e")
			if !regSuccess {
				fmt.Fprintln(os.Stderr, "ERROR: e2e regression failed")
				if regressionOutput != "" {
					if err := writeRegressionRawOutput(result.ProjectRoot, regressionOutput); err != nil {
						fmt.Fprintf(os.Stderr, "WARNING: failed to write raw-output.txt: %v\n", err)
					}
				}
				printHookJSON(map[string]any{
					"decision": "block",
					"reason":   "e2e regression failed. Read tests/e2e/results/raw-output.txt, analyze failures, then use `task add --title \"Fix: ...\" --priority P0 --breaking` to create fix tasks.",
				})
				os.Exit(0)
			}
		}
	}
}


func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func hasNpmTestScript(projectRoot string) bool {
	data, err := os.ReadFile(filepath.Join(projectRoot, "package.json"))
	if err != nil {
		return false
	}
	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false
	}
	_, ok := pkg.Scripts["test"]
	return ok
}

func hasMakeTarget(projectRoot, target string) bool {
	c := exec.Command("make", "-n", target)
	c.Dir = projectRoot
	return c.Run() == nil
}

func hasJustfile(dir string) bool {
	return fileExists(filepath.Join(dir, "justfile")) ||
		fileExists(filepath.Join(dir, "Justfile"))
}

func hasJustRecipe(dir, recipe string) bool {
	c := exec.Command("just", "--dry-run", recipe)
	c.Dir = dir
	return c.Run() == nil
}

// runCmdCapture runs a command, streams output to stderr, and returns
// the combined output as a string along with whether the command succeeded.
func runCmdCapture(dir string, name string, args ...string) (string, bool) {
	c := exec.Command(name, args...)
	c.Dir = dir
	output, err := c.CombinedOutput()
	fmt.Fprint(os.Stderr, string(output))
	return string(output), err == nil
}

// printHookJSON writes a Claude Code hook decision as JSON to stdout.
func printHookJSON(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to marshal hook JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

func runCmd(dir string, name string, args ...string) {
	c := exec.Command(name, args...)
	c.Dir = dir
	c.Stdout = os.Stderr
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s failed: %v\n", name, err)
	}
}

func runShell(dir, command string) {
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/C", command)
	} else {
		c = exec.Command("sh", "-c", command)
	}
	c.Dir = dir
	c.Stdout = os.Stderr
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: command failed: %v\n", err)
	}
}

// runShellCapture runs a shell command, streams output to stderr, and returns combined output + success.
func runShellCapture(dir, command string) (string, bool) {
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/C", command)
	} else {
		c = exec.Command("sh", "-c", command)
	}
	c.Dir = dir
	output, err := c.CombinedOutput()
	fmt.Fprint(os.Stderr, string(output))
	return string(output), err == nil
}


func runProjectTests(projectRoot, testCommand string) (string, bool) {
	if testCommand != "" {
		output, ok := runShellCapture(projectRoot, testCommand)
		return output, ok
	}

	switch {
	case hasJustfile(projectRoot) && hasJustRecipe(projectRoot, "test"):
		return runCmdCapture(projectRoot, "just", "test")
	case fileExists(filepath.Join(projectRoot, "Makefile")) && hasMakeTarget(projectRoot, "test"):
		return runCmdCapture(projectRoot, "make", "test")
	case fileExists(filepath.Join(projectRoot, "go.mod")):
		return runCmdCapture(projectRoot, "go", "test", "./...")
	case fileExists(filepath.Join(projectRoot, "package.json")) && hasNpmTestScript(projectRoot):
		return runCmdCapture(projectRoot, "npm", "test")
	case fileExists(filepath.Join(projectRoot, "pytest.ini")) || fileExists(filepath.Join(projectRoot, "pyproject.toml")):
		return runCmdCapture(projectRoot, "pytest")
	default:
		fmt.Println("WARNING: No test command found. Set testCommand in index.json.")
		return "", true
	}
}

