package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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
Exits 1 silently if any task is still pending, in_progress, or blocked.
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
	debugf := func(format string, args ...any) {
		if verbose {
			fmt.Fprintf(os.Stderr, "[debug] "+format+"\n", args...)
		}
	}

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		debugf("project root not found: %v", err)
		return nil, nil //nolint:nilerr
	}
	debugf("project root: %s", projectRoot)

	featureSlug, err := feature.GetCurrentFeature(projectRoot)
	if err != nil {
		debugf("feature not found: %v", err)
		return nil, nil //nolint:nilerr
	}
	debugf("feature: %s", featureSlug)

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		debugf("index.json not found: %s (%v)", indexPath, err)
		return nil, nil //nolint:nilerr
	}
	debugf("loaded index: %d tasks", len(index.Tasks))

	// All tasks must be completed or skipped
	for _, t := range index.Tasks {
		if t.Status != feature.StatusCompleted && t.Status != feature.StatusSkipped {
			debugf("task %s is %s — not all done", t.ID, t.Status)
			return nil, nil
		}
	}

	// Resolve e2e scripts dir
	e2eRelDir := feature.GetFeatureTestingScriptsDir(featureSlug)
	e2eAbsDir := filepath.Join(projectRoot, e2eRelDir)
	if _, err := os.Stat(e2eAbsDir); err != nil {
		debugf("e2e scripts dir not found: %s", e2eAbsDir)
		e2eAbsDir = ""
	} else {
		debugf("e2e scripts dir: %s", e2eAbsDir)
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
		os.Exit(1)
	}

	fmt.Printf("=== All tasks completed for feature: %s ===\n", result.FeatureSlug)

	// Step 1: Feature e2e tests
	if result.E2EScriptsDir != "" {
		pkgJSON := filepath.Join(result.E2EScriptsDir, "package.json")
		if _, err := os.Stat(pkgJSON); err == nil {
			fmt.Println("--- Running feature e2e tests ---")
			runCmd(result.E2EScriptsDir, "npm", "run", "test:all", "--if-present")
		} else {
			fmt.Printf("WARNING: %s has no package.json — skipping e2e tests\n", result.E2EScriptsDir)
		}
	}

	// Step 2: Project-wide tests
	fmt.Println("--- Running project-wide tests ---")
	runProjectTests(result.ProjectRoot, result.TestCommand)
}

func runProjectTests(projectRoot, testCommand string) {
	if testCommand != "" {
		runShell(projectRoot, testCommand)
		return
	}

	switch {
	case fileExists(filepath.Join(projectRoot, "go.mod")):
		runCmd(projectRoot, "go", "test", "./...")
	case fileExists(filepath.Join(projectRoot, "package.json")) && hasNpmTestScript(projectRoot):
		runCmd(projectRoot, "npm", "test")
	case fileExists(filepath.Join(projectRoot, "Makefile")) && hasMakeTarget(projectRoot, "test"):
		runCmd(projectRoot, "make", "test")
	case fileExists(filepath.Join(projectRoot, "pytest.ini")) || fileExists(filepath.Join(projectRoot, "pyproject.toml")):
		runCmd(projectRoot, "pytest")
	default:
		fmt.Println("WARNING: No test command found. Set testCommand in index.json.")
	}
}

func runCmd(dir string, name string, args ...string) {
	c := exec.Command(name, args...)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s failed: %v\n", name, err)
	}
}

func runShell(dir, command string) {
	c := exec.Command("sh", "-c", command)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: command failed: %v\n", err)
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
