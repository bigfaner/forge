package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

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

// errFixLimitExceeded is returned by appendFixTask when the fix task limit is reached.
var errFixLimitExceeded = errors.New("fix-e2e task limit exceeded")

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

	// Parse test failures and write result files
	var failures []TestFailure
	if e2eOutput != "" {
		failures = parseTestFailures(e2eOutput)

		// Match test names to test case IDs
		testCasesPath := filepath.Join(result.ProjectRoot, feature.GetFeatureTestCasesFile(result.FeatureSlug))
		for i := range failures {
			failures[i].TestCaseID = matchTestCaseID(failures[i].TestName, testCasesPath)
		}

		// Calculate stats
		stats := TestStats{
			Framework: "unknown",
			Total:     len(failures) + countPassingTests(e2eOutput),
			Fail:      len(failures),
		}
		stats.Pass = stats.Total - stats.Fail

		// Write latest.md and failure files
		if err := writeLatestMd(result.ProjectRoot, result.FeatureSlug, stats, failures); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to write latest.md: %v\n", err)
		}
		if err := writeFailureFiles(result.ProjectRoot, result.FeatureSlug, failures); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to write failure files: %v\n", err)
		}
	}

	if !e2eSuccess {
		fmt.Fprintln(os.Stderr, "ERROR: e2e tests failed")
		pendingCount, appendErr := appendFixTask(result.ProjectRoot, result.FeatureSlug, failures)
		if appendErr == errFixLimitExceeded {
			fmt.Fprintf(os.Stderr, "e2e tests failed — fix-e2e task limit (3) reached. Manual intervention required.\n")
			// Stop the agent entirely — no more retries
			printHookJSON(map[string]any{
				"continue":   false,
				"stopReason": "e2e tests failed 3 times. Manual intervention required.",
			})
			os.Exit(0)
		} else if appendErr != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to append fix task: %v\n", appendErr)
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "e2e tests failed — %d fix-e2e task(s) pending.\n", pendingCount)
			// Block Stop: tell Claude to claim and fix the e2e task
			printHookJSON(map[string]any{
				"decision": "block",
				"reason":   fmt.Sprintf("e2e tests failed. %d fix-e2e task(s) added. Run `task claim` to claim the fix-e2e task and fix the failures.", pendingCount),
			})
			os.Exit(0)
		}
		return
	}

	// e2e succeeded — attempt graduation
	if err := graduateTestScripts(result.ProjectRoot, result.FeatureSlug); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: graduation failed: %v\n", err)
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

// appendFixTask appends one fix-e2e task per failure to index.json when e2e tests fail.
// IDs use the format fix-e2e-{round}-{index}. Returns the number of tasks appended
// and errFixLimitExceeded if 3 rounds have already been attempted.
func appendFixTask(projectRoot, featureSlug string, failures []TestFailure) (int, error) {
	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return 0, fmt.Errorf("load index: %w", err)
	}

	// Check for pending/in-progress fix-e2e tasks from the current round
	pendingCount := 0
	for _, t := range index.Tasks {
		if strings.HasPrefix(t.ID, "fix-e2e-") &&
			(t.Status == feature.StatusPending || t.Status == feature.StatusInProgress) {
			pendingCount++
		}
	}

	if pendingCount > 0 {
		fmt.Fprintf(os.Stderr, "INFO: %d pending fix-e2e task(s) already exist — skipping append\n", pendingCount)
		return pendingCount, nil
	}

	if index.E2ERound >= 3 {
		return 0, errFixLimitExceeded
	}

	index.E2ERound++
	round := index.E2ERound
	tasksDir := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug))

	// Ensure at least one task even if no failures were parsed
	if len(failures) == 0 {
		failures = []TestFailure{{TestName: "unknown", TestCaseID: "unknown"}}
	}

	added := 0
	for i, f := range failures {
		idx := i + 1
		id := fmt.Sprintf("fix-e2e-%d-%d", round, idx)
		taskFile := fmt.Sprintf("%s.md", id)
		taskFilePath := filepath.Join(tasksDir, taskFile)

		if err := createFixTaskFile(taskFilePath, round, idx, f); err != nil {
			return added, fmt.Errorf("create task file %s: %w", id, err)
		}

		index.Tasks[id] = task.Task{
			ID:       id,
			Title:    fmt.Sprintf("修复 e2e 测试失败: %s", f.TestName),
			Priority: feature.PriorityP0,
			Status:   feature.StatusPending,
			File:     taskFile,
			Record:   fmt.Sprintf("records/%s.md", id),
		}
		added++
	}

	if err := saveIndexAtomic(indexPath, index); err != nil {
		return 0, err
	}
	return added, nil
}

// createFixTaskFile creates a fix-e2e task file for a single failure.
func createFixTaskFile(filePath string, round, idx int, f TestFailure) error {
	id := fmt.Sprintf("fix-e2e-%d-%d", round, idx)
	failureRef := fmt.Sprintf("- `testing/results/failures/failure-%s.md` — %s\n", f.TestCaseID, f.TestName)

	content := fmt.Sprintf(`---
id: "%s"
title: "修复 e2e 测试失败: %s"
priority: "P0"
estimated_time: "30min-2h"
dependencies: []
status: pending
---

# %s: 修复 e2e 测试失败

## Description

这是第 %d 轮修复尝试。修复步骤：

1. 读取 `+"`testing/results/latest.md`"+` 查看失败概览
2. 读取 `+"`testing/results/failures/failure-%s.md`"+` 了解具体失败详情
3. 定位根本原因（代码逻辑 / 测试脚本 / 环境配置）
4. 修复并验证

## Reference Files

- `+"`testing/results/latest.md`"+` — 测试结果概览
%s
- `+"`testing/test-cases.md`"+` — 测试用例文档
- `+"`testing/scripts/`"+` — 测试脚本目录

## Acceptance Criteria

- [ ] 已定位失败的根本原因
- [ ] 已修复代码或测试脚本
- [ ] 单元测试全部通过
`, id, f.TestName, id, round, f.TestCaseID, failureRef)

	return os.WriteFile(filePath, []byte(content), 0644)
}

// graduateTestScripts migrates test scripts to tests/e2e/<target>/ on first success.
func graduateTestScripts(projectRoot, featureSlug string) error {
	markerPath := feature.GetE2EGraduatedMarker(projectRoot, featureSlug)
	if fileExists(markerPath) {
		return nil // already graduated
	}

	scriptsDir := filepath.Join(projectRoot, feature.GetFeatureTestingScriptsDir(featureSlug))
	if !fileExists(scriptsDir) {
		return nil // no scripts to graduate
	}

	testCasesPath := filepath.Join(projectRoot, feature.GetFeatureTestCasesFile(featureSlug))
	targets, err := parseTargetsFromTestCases(testCasesPath)
	if err != nil {
		return fmt.Errorf("parse test cases: %w", err)
	}

	if len(targets) == 0 {
		return nil // no targets defined, skip graduation
	}

	// Copy shared infrastructure to tests/e2e/ (helpers, config, types)
	e2eBaseDir := filepath.Join(projectRoot, feature.E2ETestsBaseDir)
	if err := os.MkdirAll(e2eBaseDir, 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", e2eBaseDir, err)
	}
	for _, shared := range []string{"helpers.ts", "package.json", "tsconfig.json"} {
		srcPath := filepath.Join(scriptsDir, shared)
		destPath := filepath.Join(e2eBaseDir, shared)
		if fileExists(srcPath) && !fileExists(destPath) {
			if err := copyFile(srcPath, destPath); err != nil {
				return fmt.Errorf("copy shared %s: %w", shared, err)
			}
			fmt.Printf("INFO: graduated shared %s → %s\n", shared, destPath)
		}
	}

	// Run npm install in tests/e2e/ if node_modules doesn't exist
	nodeModules := filepath.Join(e2eBaseDir, "node_modules")
	if !fileExists(nodeModules) && fileExists(filepath.Join(e2eBaseDir, "package.json")) {
		fmt.Fprintln(os.Stderr, "INFO: running npm install in tests/e2e/ ...")
		c := exec.Command("npm", "install")
		c.Dir = e2eBaseDir
		c.Stdout = os.Stderr
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("npm install in tests/e2e/: %w", err)
		}
	}

	// Copy spec files to tests/e2e/<target>/ and rewrite imports
	typeToSpec := map[string]string{
		"ui":  "ui.spec.ts",
		"api": "api.spec.ts",
		"cli": "cli.spec.ts",
	}

	for testType, specFile := range typeToSpec {
		srcPath := filepath.Join(scriptsDir, specFile)
		if !fileExists(srcPath) {
			continue
		}
		for _, target := range targets[testType] {
			destDir := feature.GetE2ETargetDir(projectRoot, target)
			if err := os.MkdirAll(destDir, 0755); err != nil {
				return fmt.Errorf("mkdir %s: %w", destDir, err)
			}
			destPath := filepath.Join(destDir, specFile)
			if err := copyAndRewriteImports(srcPath, destPath); err != nil {
				return fmt.Errorf("copy %s → %s: %w", srcPath, destPath, err)
			}
			fmt.Printf("INFO: graduated %s → %s\n", specFile, destPath)
		}
	}

	// Create graduation marker
	markerDir := filepath.Dir(markerPath)
	if err := os.MkdirAll(markerDir, 0755); err != nil {
		return fmt.Errorf("mkdir graduated dir: %w", err)
	}
	timestamp := time.Now().Format(time.RFC3339)
	return os.WriteFile(markerPath, []byte(timestamp+"\n"), 0644)
}

// parseTargetsFromTestCases reads test-cases.md and returns unique targets grouped by type.
// Returns map[type][]target, e.g. {"ui": ["ui/login", "ui/dashboard"], "api": ["api/auth"]}.
func parseTargetsFromTestCases(path string) (map[string][]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	seen := map[string]bool{}
	result := map[string][]string{}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Match: - **Target**: ui/login
		if !strings.Contains(line, "**Target**") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		target := strings.TrimSpace(parts[1])
		if target == "" || seen[target] {
			continue
		}
		seen[target] = true

		// Determine type from target prefix (e.g. "ui/login" → "ui")
		typeParts := strings.SplitN(target, "/", 2)
		if len(typeParts) < 2 {
			continue
		}
		testType := typeParts[0]
		result[testType] = append(result[testType], target)
	}
	return result, scanner.Err()
}

// saveIndexAtomic writes index.json atomically (temp file + rename).
func saveIndexAtomic(path string, index *task.TaskIndex) error {
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal index: %w", err)
	}
	data = append(data, '\n')

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	// On Windows, os.Rename cannot overwrite an existing file
	if runtime.GOOS == "windows" {
		os.Remove(path)
	}
	return os.Rename(tmpPath, path)
}

// copyFile copies src to dst, overwriting dst if it exists.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// helpersImportRe matches import paths like './helpers.js', "./helpers", '../helpers'
var helpersImportRe = regexp.MustCompile(`['"]\.\./?helpers(?:\.js)?['"]`)

// copyAndRewriteImports copies src to dst, rewriting "./helpers.js" → "../helpers.js"
// so graduated specs can import shared helpers from tests/e2e/helpers.ts.
func copyAndRewriteImports(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	rewritten := helpersImportRe.ReplaceAllString(string(data), "'../helpers.js'")
	return os.WriteFile(dst, []byte(rewritten), 0644)
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

// runCmdCapture runs a command, streams output to stderr, and returns
// the combined output as a string along with whether the command succeeded.
// Output goes to stderr so stdout remains clean for hook JSON decisions.
func runCmdCapture(dir string, name string, args ...string) (string, bool) {
	c := exec.Command(name, args...)
	c.Dir = dir
	var buf bytes.Buffer
	mw := io.MultiWriter(os.Stderr, &buf)
	c.Stdout = mw
	c.Stderr = mw
	err := c.Run()
	return buf.String(), err == nil
}

// printHookJSON writes a Claude Code hook decision as JSON to stdout.
// Stdout must contain only this JSON for Claude Code to parse it correctly.
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

// countPassingTests counts passing tests from node:test TAP output.
// Matches lines starting with "ok " or "✓" at line start.
func countPassingTests(output string) int {
	count := 0
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "ok ") || strings.HasPrefix(trimmed, "✓") {
			count++
		}
	}
	return count
}

// runSpecsIndividually runs each .spec.ts in dir sequentially, collecting all output.
// Unlike npm run test:all, this runs all specs regardless of individual failures.
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
