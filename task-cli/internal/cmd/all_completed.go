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
		os.Exit(1)
	}

	fmt.Printf("=== All tasks completed for feature: %s ===\n", result.FeatureSlug)

	// Step 1: Feature e2e tests
	var e2eOutput string
	var e2eSuccess bool

	switch {
	case hasJustfile(result.ProjectRoot) && hasJustRecipe(result.ProjectRoot, "test-e2e"):
		fmt.Println("--- Running feature e2e tests (just test-e2e) ---")
		e2eOutput, e2eSuccess = runCmdCapture(result.ProjectRoot, "just", "test-e2e")
	case fileExists(filepath.Join(result.ProjectRoot, "Makefile")) && hasMakeTarget(result.ProjectRoot, "test-e2e"):
		fmt.Println("--- Running feature e2e tests (make test-e2e) ---")
		e2eOutput, e2eSuccess = runCmdCapture(result.ProjectRoot, "make", "test-e2e")
	case result.E2EScriptsDir != "":
		pkgJSON := filepath.Join(result.E2EScriptsDir, "package.json")
		if _, err := os.Stat(pkgJSON); err == nil {
			fmt.Println("--- Running feature e2e tests ---")
			e2eOutput, e2eSuccess = runCmdCapture(result.E2EScriptsDir, "npm", "run", "test:all", "--if-present")
		} else {
			fmt.Printf("WARNING: %s has no package.json — skipping e2e tests\n", result.E2EScriptsDir)
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
			Total:     len(failures) + strings.Count(e2eOutput, "✓") + strings.Count(e2eOutput, "ok"),
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
		appendErr := appendFixTask(result.ProjectRoot, result.FeatureSlug, failures)
		if appendErr == errFixLimitExceeded {
			fmt.Println("WARNING: fix-e2e task limit (3) reached — skipping fix task append")
			// exit 0 so the Stop hook doesn't loop forever
		} else if appendErr != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to append fix task: %v\n", appendErr)
			os.Exit(1)
		} else {
			os.Exit(1)
		}
		return
	}

	// e2e succeeded — attempt graduation
	if err := graduateTestScripts(result.ProjectRoot, result.FeatureSlug); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: graduation failed: %v\n", err)
	}

	// Step 2: Project-wide tests
	fmt.Println("--- Running project-wide tests ---")
	runProjectTests(result.ProjectRoot, result.TestCommand)
}

// appendFixTask appends a fix-e2e-N task to index.json when e2e tests fail.
// Returns errFixLimitExceeded if the limit (3) is reached.
func appendFixTask(projectRoot, featureSlug string, failures []TestFailure) error {
	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return fmt.Errorf("load index: %w", err)
	}

	// Count existing fix-e2e tasks and check for pending ones
	fixCount := 0
	for _, t := range index.Tasks {
		if strings.HasPrefix(t.ID, "fix-e2e-") {
			fixCount++
			if t.Status == feature.StatusPending {
				fmt.Println("INFO: pending fix-e2e task already exists — skipping append")
				return nil
			}
		}
	}

	if fixCount >= 3 {
		return errFixLimitExceeded
	}

	n := fixCount + 1
	id := fmt.Sprintf("fix-e2e-%d", n)
	key := id

	// Create task file from template
	tasksDir := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug))
	taskFile := fmt.Sprintf("%s.md", id)
	taskFilePath := filepath.Join(tasksDir, taskFile)

	if err := createFixTaskFile(taskFilePath, n, failures); err != nil {
		return fmt.Errorf("create task file: %w", err)
	}

	newTask := task.Task{
		ID:       id,
		Title:    "修复 e2e 测试失败",
		Priority: feature.PriorityP0,
		Status:   feature.StatusPending,
		File:     taskFile,
		Record:   fmt.Sprintf("records/%s.md", id),
	}
	index.Tasks[key] = newTask

	return saveIndexAtomic(indexPath, index)
}

// createFixTaskFile creates a fix-e2e task file from template.
func createFixTaskFile(filePath string, n int, failures []TestFailure) error {
	// Build failure references section
	var failureRefs strings.Builder
	for _, f := range failures {
		failureRefs.WriteString(fmt.Sprintf("- `testing/results/failures/failure-%s.md` — %s\n",
			f.TestCaseID, f.TestName))
	}

	template := fmt.Sprintf(`---
id: "fix-e2e-%d"
title: "修复 e2e 测试失败"
priority: "P0"
estimated_time: "30min-2h"
dependencies: []
status: pending
---

# fix-e2e-%d: 修复 e2e 测试失败

## Description

e2e 测试失败，需要分析失败原因并修复代码。

## Reference Files

- `+"`testing/results/latest.md`"+` — 测试结果概览
%s
- `+"`testing/test-cases.md`"+` — 测试用例文档
- `+"`testing/scripts/`"+` — 测试脚本目录

## Acceptance Criteria

- [ ] 已读取 `+"`testing/results/latest.md`"+` 了解失败概览
- [ ] 已读取 failure 文件了解具体失败详情
- [ ] 已定位失败的根本原因
- [ ] 已修复代码或测试脚本
- [ ] 本地验证测试通过（可选）
- [ ] `+"`task all-completed`"+` 再次运行时测试通过

## User Stories

No direct user story mapping. This is a test fix task.

## Implementation Notes

1. 读取 `+"`testing/results/latest.md`"+` 查看失败概览
2. 读取对应的 failure-{test-case-id}.md 文件了解具体失败详情
3. 分析失败原因：
   - 代码逻辑错误？
   - 测试脚本问题？
   - 环境配置问题？
4. 修复问题
5. 如果需要，可运行 `+"`npm run test:all`"+` 本地验证
6. 完成后执行 `+"`task record`"+` 记录修复内容

## Context

这是第 %d 次尝试修复 e2e 测试失败。如果修复后测试仍失败，会创建 fix-e2e-%d 任务。
最多允许 3 次修复尝试。
`, n, n, failureRefs.String(), n, n+1)

	return os.WriteFile(filePath, []byte(template), 0644)
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

	// Copy spec files to tests/e2e/<target>/ for each target
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
			if err := copyFile(srcPath, destPath); err != nil {
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

// runCmdCapture runs a command, streams output to stdout/stderr, and returns
// the combined output as a string along with whether the command succeeded.
func runCmdCapture(dir string, name string, args ...string) (string, bool) {
	c := exec.Command(name, args...)
	c.Dir = dir
	var buf bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &buf)
	c.Stdout = mw
	c.Stderr = mw
	err := c.Run()
	return buf.String(), err == nil
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

func hasJustfile(dir string) bool {
	return fileExists(filepath.Join(dir, "justfile")) ||
		fileExists(filepath.Join(dir, "Justfile"))
}

func hasJustRecipe(dir, recipe string) bool {
	c := exec.Command("just", "--dry-run", recipe)
	c.Dir = dir
	return c.Run() == nil
}
