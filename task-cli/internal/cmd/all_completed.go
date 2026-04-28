package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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
If all done: runs feature e2e tests, then project-wide tests.

Use -v to see why the command exits early (useful for debugging).`,
	Run: runAllCompleted,
}

func init() {
	allCompletedCmd.Flags().BoolVarP(&allCompletedVerbose, "verbose", "v", false, "print debug info when exiting early")
}

// AllCompletedResult holds context for running tests after all tasks complete.
type AllCompletedResult struct {
	FeatureSlug   string
	ProjectRoot   string
	E2EScriptsDir string // empty if dir doesn't exist
	TestCommand   string // empty if not set in index.json
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

	// Resolve e2e scripts dir
	e2eRelDir := feature.GetFeatureTestingScriptsDir(featureSlug)
	e2eAbsDir := filepath.Join(projectRoot, e2eRelDir)
	if _, err := os.Stat(e2eAbsDir); err != nil {
		Debugf(verbose, "e2e scripts dir not found: %s", e2eAbsDir)
		e2eAbsDir = ""
	} else {
		Debugf(verbose, "e2e scripts dir: %s", e2eAbsDir)
	}

	return &AllCompletedResult{
		FeatureSlug:   featureSlug,
		ProjectRoot:   projectRoot,
		E2EScriptsDir: e2eAbsDir,
		TestCommand:   index.TestCommand,
	}, nil
}

func runAllCompleted(cmd *cobra.Command, args []string) {
	result, err := checkAllCompleted(allCompletedVerbose)
	if err != nil || result == nil {
		os.Exit(0) // not all done is normal, exit silently
	}

	fmt.Fprintf(os.Stderr, "=== All tasks completed for feature: %s ===\n", result.FeatureSlug)

	// Step 1: Feature e2e tests
	var e2eOutput string
	var e2eSuccess bool

	switch {
	case hasJustfile(result.ProjectRoot) && hasJustRecipe(result.ProjectRoot, "test-e2e"):
		fmt.Fprintf(os.Stderr, "--- Running feature e2e tests (just test-e2e --feature %s) ---\n", result.FeatureSlug)
		e2eOutput, e2eSuccess = runCmdCapture(result.ProjectRoot, "just", "test-e2e", "--feature", result.FeatureSlug)
	case fileExists(filepath.Join(result.ProjectRoot, "Makefile")) && hasMakeTarget(result.ProjectRoot, "test-e2e"):
		fmt.Fprintf(os.Stderr, "--- Running feature e2e tests (make test-e2e FEATURE=%s) ---\n", result.FeatureSlug)
		e2eOutput, e2eSuccess = runCmdCapture(result.ProjectRoot, "make", "test-e2e", "FEATURE="+result.FeatureSlug)
	case result.E2EScriptsDir != "":
		pkgJSON := filepath.Join(result.E2EScriptsDir, "package.json")
		if _, err := os.Stat(pkgJSON); err == nil {
			fmt.Fprintln(os.Stderr, "--- Running feature e2e tests (individual specs) ---")
			e2eOutput, e2eSuccess = runSpecsIndividually(result.E2EScriptsDir)
		} else {
			fmt.Fprintf(os.Stderr, "WARNING: %s has no package.json — skipping e2e tests\n", result.E2EScriptsDir)
			e2eSuccess = true // no scripts to run, treat as success
		}
	default:
		e2eSuccess = true // no e2e configured, skip
	}

	if !e2eSuccess {
		fmt.Fprintln(os.Stderr, "ERROR: e2e tests failed")

		// Save raw output for agent analysis
		if e2eOutput != "" {
			if err := writeRawOutput(result.ProjectRoot, result.FeatureSlug, e2eOutput); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: failed to write raw-output.txt: %v\n", err)
			}
			if err := writeLatestMd(result.ProjectRoot, result.FeatureSlug, TestStats{Fail: 1}); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: failed to write latest.md: %v\n", err)
			}
		}

		// Block Stop: tell Claude to analyze failures and add fix tasks
		printHookJSON(map[string]any{
			"decision": "block",
			"reason":   "e2e tests failed. Read testing/results/raw-output.txt, analyze failures, then use `task add --title \"Fix: ...\" --priority P0 --breaking` to create fix tasks for each failure.",
		})
		os.Exit(0)
		return
	}

	// Write passing results
	if e2eOutput != "" {
		stats := TestStats{
			Total: countTestLines(e2eOutput),
			Pass:  countTestLines(e2eOutput), // all passed since exit code was 0
		}
		if err := writeLatestMd(result.ProjectRoot, result.FeatureSlug, stats); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to write latest.md: %v\n", err)
		}
	}

	// Step 2: Project-wide unit tests
	fmt.Fprintln(os.Stderr, "--- Running project-wide tests ---")
	runProjectTests(result.ProjectRoot, result.TestCommand)

	// Step 3: Full e2e regression (graduated scripts in tests/e2e/)
	if hasJustfile(result.ProjectRoot) && hasJustRecipe(result.ProjectRoot, "test-e2e") {
		fmt.Fprintln(os.Stderr, "--- Running full e2e regression (just test-e2e) ---")
		_, regSuccess := runCmdCapture(result.ProjectRoot, "just", "test-e2e")
		if !regSuccess {
			fmt.Fprintln(os.Stderr, "ERROR: e2e regression failed: see output above")
			os.Exit(1)
		}
	}
}

// countTestLines provides a rough count of test results from output.
// Not used for accurate stats — only for the passing summary.
func countTestLines(output string) int {
	count := 0
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "ok ") || strings.HasPrefix(trimmed, "✓") ||
			strings.HasPrefix(trimmed, "not ok ") || strings.HasPrefix(trimmed, "✗") {
			count++
		}
	}
	return count
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

// runSpecsIndividually runs each .spec.ts in dir sequentially, collecting all output.
func runSpecsIndividually(dir string) (string, bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Sprintf("ERROR: read dir %s: %v", dir, err), false
	}

	var allOutput strings.Builder
	allSuccess := true

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".spec.ts") {
			continue
		}
		specPath := filepath.Join(dir, entry.Name())
		output, success := runCmdCapture(dir, "npx", "tsx", specPath)
		allOutput.WriteString(output)
		if !success {
			allSuccess = false
		}
	}

	return allOutput.String(), allSuccess
}

func runProjectTests(projectRoot, testCommand string) {
	if testCommand != "" {
		runShell(projectRoot, testCommand)
		return
	}

	switch {
	case hasJustfile(projectRoot) && hasJustRecipe(projectRoot, "test"):
		runCmd(projectRoot, "just", "test")
	case fileExists(filepath.Join(projectRoot, "Makefile")) && hasMakeTarget(projectRoot, "test"):
		runCmd(projectRoot, "make", "test")
	case fileExists(filepath.Join(projectRoot, "go.mod")):
		runCmd(projectRoot, "go", "test", "./...")
	case fileExists(filepath.Join(projectRoot, "package.json")) && hasNpmTestScript(projectRoot):
		runCmd(projectRoot, "npm", "test")
	case fileExists(filepath.Join(projectRoot, "pytest.ini")) || fileExists(filepath.Join(projectRoot, "pyproject.toml")):
		runCmd(projectRoot, "pytest")
	default:
		fmt.Println("WARNING: No test command found. Set testCommand in index.json.")
	}
}

