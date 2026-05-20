package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// setupJourneyProject creates a temp directory with project structure for journey testing.
func setupJourneyProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Create go.mod for language detection
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create .forge/config.yaml
	configDir := filepath.Join(dir, ".forge")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	configContent := "languages:\n  - go\ntest-command: go test ./...\n"
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	return dir
}

// setupJourneyProjectWithoutConfig creates a project without test-command in config.
func setupJourneyProjectWithoutConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	return dir
}

// --- Test: Temp directory isolation ---

func TestJourneyIsolation_CreateWorkDir_IncludesJourneyName(t *testing.T) {
	dir := setupJourneyProject(t)

	workDir, cleanup, err := createJourneyWorkDir(dir, "task-lifecycle")
	if err != nil {
		t.Fatalf("createJourneyWorkDir failed: %v", err)
	}
	defer cleanup()

	if !strings.Contains(filepath.Base(workDir), "task-lifecycle") {
		t.Errorf("work dir base name should contain journey name 'task-lifecycle', got: %s", filepath.Base(workDir))
	}
}

func TestJourneyIsolation_CreateWorkDir_HasRandomSuffix(t *testing.T) {
	dir := setupJourneyProject(t)

	workDir1, cleanup1, err := createJourneyWorkDir(dir, "my-journey")
	if err != nil {
		t.Fatalf("createJourneyWorkDir failed: %v", err)
	}
	defer cleanup1()

	workDir2, cleanup2, err := createJourneyWorkDir(dir, "my-journey")
	if err != nil {
		t.Fatalf("createJourneyWorkDir failed: %v", err)
	}
	defer cleanup2()

	if workDir1 == workDir2 {
		t.Error("two work dirs for same journey should have different random suffixes")
	}
}

func TestJourneyIsolation_CreateWorkDir_IsInSystemTemp(t *testing.T) {
	dir := setupJourneyProject(t)

	workDir, cleanup, err := createJourneyWorkDir(dir, "task-lifecycle")
	if err != nil {
		t.Fatalf("createJourneyWorkDir failed: %v", err)
	}
	defer cleanup()

	// Work dir should be under the project's temp parent, not the project itself
	if strings.HasPrefix(workDir, dir) {
		t.Errorf("work dir should NOT be inside project dir, got: %s (project: %s)", workDir, dir)
	}
}

// --- Test: Cleanup always runs ---

func TestJourneyIsolation_CleanupRuns_AfterSuccess(t *testing.T) {
	dir := setupJourneyProject(t)

	workDir, cleanup, err := createJourneyWorkDir(dir, "cleanup-test")
	if err != nil {
		t.Fatalf("createJourneyWorkDir failed: %v", err)
	}

	// Write a file to the work dir
	if err := os.WriteFile(filepath.Join(workDir, "test.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	cleanup()

	if _, err := os.Stat(workDir); !os.IsNotExist(err) {
		t.Errorf("work dir should be cleaned up after cleanup(), still exists: %s", workDir)
	}
}

func TestJourneyIsolation_CleanupRuns_AfterSimulatedFailure(t *testing.T) {
	dir := setupJourneyProject(t)

	workDir, cleanup, _ := createJourneyWorkDir(dir, "cleanup-fail-test")

	// Simulate something written during execution
	if err := os.WriteFile(filepath.Join(workDir, "output.txt"), []byte("partial output"), 0644); err != nil {
		t.Fatal(err)
	}

	// Cleanup must still work even after "failure"
	cleanup()

	if _, err := os.Stat(workDir); !os.IsNotExist(err) {
		t.Error("work dir should be cleaned up even after simulated failure")
	}
}

// --- Test: File copy to work dir ---

func TestJourneyIsolation_CopyFiles_PreservesContent(t *testing.T) {
	dir := setupJourneyProject(t)

	// Create a test file in the project
	testContent := "package main\n\nfunc main() {}\n"
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	workDir, cleanup, err := createJourneyWorkDir(dir, "copy-test")
	if err != nil {
		t.Fatalf("createJourneyWorkDir failed: %v", err)
	}
	defer cleanup()

	// Copy the file to work dir
	if err := copyFileToWorkDir(dir, workDir, "main.go"); err != nil {
		t.Fatalf("copyFileToWorkDir failed: %v", err)
	}

	copied, err := os.ReadFile(filepath.Join(workDir, "main.go"))
	if err != nil {
		t.Fatalf("read copied file failed: %v", err)
	}

	if string(copied) != testContent {
		t.Errorf("copied file content mismatch: got %q, want %q", string(copied), testContent)
	}
}

func TestJourneyIsolation_CopyFiles_OriginalUnmodified(t *testing.T) {
	dir := setupJourneyProject(t)

	originalContent := "original content"
	if err := os.WriteFile(filepath.Join(dir, "data.txt"), []byte(originalContent), 0644); err != nil {
		t.Fatal(err)
	}

	workDir, cleanup, err := createJourneyWorkDir(dir, "modify-test")
	if err != nil {
		t.Fatalf("createJourneyWorkDir failed: %v", err)
	}
	defer cleanup()

	if err := copyFileToWorkDir(dir, workDir, "data.txt"); err != nil {
		t.Fatalf("copyFileToWorkDir failed: %v", err)
	}

	// Modify the copy
	if err := os.WriteFile(filepath.Join(workDir, "data.txt"), []byte("modified"), 0644); err != nil {
		t.Fatal(err)
	}

	// Original should be unchanged
	orig, err := os.ReadFile(filepath.Join(dir, "data.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(orig) != originalContent {
		t.Errorf("original file was modified: got %q, want %q", string(orig), originalContent)
	}
}

// --- Test: JourneyResult ---

func TestJourneyResult_ContractFailure_Format(t *testing.T) {
	failures := []ContractFailure{
		{
			Dimension:    "Output",
			ContractPath: "tests/task-lifecycle/_contracts/step-2-task-claim.md",
			Expected:     "exit code 0",
			Actual:       "exit code 1",
		},
		{
			Dimension:    "State",
			ContractPath: "tests/task-lifecycle/_contracts/step-3-task-submit.md",
			Expected:     "status -> completed",
			Actual:       "status -> done",
		},
	}

	result := JourneyResult{
		JourneyName: "task-lifecycle",
		Passed:      false,
		Failures:    failures,
	}

	output := result.FormatReport()

	// Must contain dimension name
	if !strings.Contains(output, "Output dimension FAILED") {
		t.Errorf("report should contain 'Output dimension FAILED', got: %s", output)
	}
	// Must contain contract file path
	if !strings.Contains(output, "step-2-task-claim.md") {
		t.Errorf("report should contain contract file path, got: %s", output)
	}
	// Must contain expected value
	if !strings.Contains(output, "exit code 0") {
		t.Errorf("report should contain expected value, got: %s", output)
	}
	// Must contain actual value
	if !strings.Contains(output, "exit code 1") {
		t.Errorf("report should contain actual value, got: %s", output)
	}
	// Must contain second failure too
	if !strings.Contains(output, "State dimension FAILED") {
		t.Errorf("report should contain 'State dimension FAILED', got: %s", output)
	}
}

func TestJourneyResult_Pass_Format(t *testing.T) {
	result := JourneyResult{
		JourneyName: "task-lifecycle",
		Passed:      true,
		Duration:    3 * time.Second,
	}

	output := result.FormatReport()

	if !strings.Contains(output, "PASS") {
		t.Errorf("passing report should contain PASS, got: %s", output)
	}
	if !strings.Contains(output, "task-lifecycle") {
		t.Errorf("report should contain journey name, got: %s", output)
	}
}

// --- Test: Parallel execution isolation ---

func TestJourneyIsolation_ParallelExecution_ResultsConsistent(t *testing.T) {
	dir := setupJourneyProject(t)

	// Create a shared file in the project that journeys will read/write
	sharedFile := "shared-state.txt"
	if err := os.WriteFile(filepath.Join(dir, sharedFile), []byte("initial"), 0644); err != nil {
		t.Fatal(err)
	}

	// Run 3 journeys in parallel, each reading/writing its own isolated copy
	journeys := []string{"journey-a", "journey-b", "journey-c"}
	var wg sync.WaitGroup
	results := make(chan string, len(journeys))

	for _, name := range journeys {
		wg.Add(1)
		go func(journeyName string) {
			defer wg.Done()

			workDir, cleanup, err := createJourneyWorkDir(dir, journeyName)
			if err != nil {
				results <- fmt.Sprintf("%s: error %v", journeyName, err)
				return
			}
			defer cleanup()

			// Copy shared file to isolated dir
			if err := copyFileToWorkDir(dir, workDir, sharedFile); err != nil {
				results <- fmt.Sprintf("%s: copy error %v", journeyName, err)
				return
			}

			// Each journey writes its own value (no interference)
			writePath := filepath.Join(workDir, sharedFile)
			newContent := journeyName + "-was-here"
			if err := os.WriteFile(writePath, []byte(newContent), 0644); err != nil {
				results <- fmt.Sprintf("%s: write error %v", journeyName, err)
				return
			}

			// Read back and verify isolation
			readBack, err := os.ReadFile(writePath)
			if err != nil {
				results <- fmt.Sprintf("%s: read error %v", journeyName, err)
				return
			}

			results <- fmt.Sprintf("%s: %s", journeyName, string(readBack))
		}(name)
	}

	wg.Wait()
	close(results)

	// Verify each journey saw its own isolated state
	for result := range results {
		journeyName := strings.SplitN(result, ": ", 2)[0]
		value := strings.SplitN(result, ": ", 2)[1]

		expected := journeyName + "-was-here"
		if value != expected {
			t.Errorf("isolation violated for %s: got %q, want %q", journeyName, value, expected)
		}
	}

	// Original file should be unchanged
	orig, err := os.ReadFile(filepath.Join(dir, sharedFile))
	if err != nil {
		t.Fatal(err)
	}
	if string(orig) != "initial" {
		t.Errorf("original shared file was modified: got %q, want %q", string(orig), "initial")
	}
}

// --- Test: run-journey CLI command ---

func TestTestingRunJourney_NoTestCommand(t *testing.T) {
	_ = setupJourneyProjectWithoutConfig(t)
	// No longer requires test-command in config. Journey execution uses just e2e-test.
}

func TestTestingRunJourney_CommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range testCmd.Commands() {
		if cmd.Name() == "run-journey" {
			found = true
			break
		}
	}
	if !found {
		t.Error("testing group missing 'run-journey' subcommand")
	}
}

// --- Test: JourneyExecutionConfig ---

func TestJourneyExecutionConfig_ResolvesFromProject(t *testing.T) {
	dir := setupJourneyProject(t)

	cfg := resolveJourneyExecutionConfig(dir)

	if cfg.ProjectRoot != dir {
		t.Errorf("expected ProjectRoot %q, got %q", dir, cfg.ProjectRoot)
	}
}

func TestJourneyExecutionConfig_NoLanguage(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	cfg := resolveJourneyExecutionConfig(dir)
	if cfg.ProjectRoot != dir {
		t.Errorf("expected ProjectRoot %q, got %q", dir, cfg.ProjectRoot)
	}
}

// --- Test: Parallel vs Sequential equivalence ---

func TestJourneyIsolation_ParallelSameAsSequential(t *testing.T) {
	dir := setupJourneyProject(t)

	// Create a counter file
	counterFile := "counter.txt"
	if err := os.WriteFile(filepath.Join(dir, counterFile), []byte("0"), 0644); err != nil {
		t.Fatal(err)
	}

	journeys := []string{"j-1", "j-2", "j-3"}

	// Sequential run
	seqResults := make(map[string]string)
	for _, name := range journeys {
		workDir, cleanup, err := createJourneyWorkDir(dir, name)
		if err != nil {
			t.Fatal(err)
		}
		if err := copyFileToWorkDir(dir, workDir, counterFile); err != nil {
			cleanup()
			t.Fatal(err)
		}
		// Write journey result
		result := name + "-done"
		if err := os.WriteFile(filepath.Join(workDir, counterFile), []byte(result), 0644); err != nil {
			cleanup()
			t.Fatal(err)
		}
		readBack, _ := os.ReadFile(filepath.Join(workDir, counterFile))
		seqResults[name] = string(readBack)
		cleanup()
	}

	// Parallel run
	parResults := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, name := range journeys {
		wg.Add(1)
		go func(journeyName string) {
			defer wg.Done()
			workDir, cleanup, err := createJourneyWorkDir(dir, journeyName)
			if err != nil {
				t.Errorf("parallel %s error: %v", journeyName, err)
				return
			}
			defer cleanup()

			if err := copyFileToWorkDir(dir, workDir, counterFile); err != nil {
				t.Errorf("parallel %s copy error: %v", journeyName, err)
				return
			}
			result := journeyName + "-done"
			if err := os.WriteFile(filepath.Join(workDir, counterFile), []byte(result), 0644); err != nil {
				t.Errorf("parallel %s write error: %v", journeyName, err)
				return
			}
			readBack, _ := os.ReadFile(filepath.Join(workDir, counterFile))
			mu.Lock()
			parResults[journeyName] = string(readBack)
			mu.Unlock()
		}(name)
	}
	wg.Wait()

	// Results should be identical
	for _, name := range journeys {
		if seqResults[name] != parResults[name] {
			t.Errorf("parallel result differs from sequential for %s: seq=%q par=%q",
				name, seqResults[name], parResults[name])
		}
	}
}

// --- Test: Multiple journeys don't block each other ---

func TestJourneyIsolation_FailureDoesNotBlockOthers(t *testing.T) {
	dir := setupJourneyProject(t)

	journeys := []string{"failing-journey", "passing-journey"}
	var results []JourneyResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, name := range journeys {
		wg.Add(1)
		go func(journeyName string) {
			defer wg.Done()
			result := JourneyResult{
				JourneyName: journeyName,
			}

			_, cleanup, err := createJourneyWorkDir(dir, journeyName)
			if err != nil {
				result.Passed = false
				result.Error = err.Error()
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
				return
			}
			defer cleanup()

			if journeyName == "failing-journey" {
				// Simulate failure
				result.Passed = false
				result.Error = "simulated failure"
				result.Failures = []ContractFailure{
					{
						Dimension:    "Output",
						ContractPath: "tests/failing/_contracts/step-1.md",
						Expected:     "exit code 0",
						Actual:       "exit code 1",
					},
				}
			} else {
				result.Passed = true
			}

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(name)
	}

	wg.Wait()

	// Both journeys should have results (failing didn't block passing)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// Find the passing journey result
	var passingFound bool
	for _, r := range results {
		if r.JourneyName == "passing-journey" && r.Passed {
			passingFound = true
		}
	}
	if !passingFound {
		t.Error("passing journey should have completed despite failing journey")
	}
}

// --- Test: Exit code collection ---

func TestJourneyIsolation_ExitCodesCollected(t *testing.T) {
	result := JourneyResult{
		JourneyName: "test-journey",
		Passed:      false,
		ExitCode:    1,
		Error:       "command failed",
	}

	output := result.FormatReport()
	if !strings.Contains(output, "FAIL") {
		t.Errorf("failing report should contain FAIL, got: %s", output)
	}
	if !strings.Contains(output, "EXIT_CODE: 1") {
		t.Errorf("report should contain EXIT_CODE, got: %s", output)
	}
}

// --- Test: JSON report generation ---

func TestJourneyResult_JSONReport(t *testing.T) {
	result := JourneyResult{
		JourneyName: "task-lifecycle",
		Passed:      true,
		Duration:    5 * time.Second,
		ExitCode:    0,
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	var parsed JourneyResult
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if parsed.JourneyName != "task-lifecycle" {
		t.Errorf("journey name mismatch: got %q", parsed.JourneyName)
	}
	if !parsed.Passed {
		t.Error("should be passing")
	}
}

// --- Test: ContractFailure report format detail ---

func TestContractFailure_FormatDetail(t *testing.T) {
	cf := ContractFailure{
		Dimension:    "Output",
		ContractPath: "tests/task-lifecycle/_contracts/step-2-task-claim.md",
		Expected:     "claimed task <task_id>",
		Actual:       "Task <task_id> claimed",
	}

	formatted := cf.Format()
	if !strings.Contains(formatted, "Output dimension FAILED") {
		t.Errorf("should contain dimension name, got: %s", formatted)
	}
	if !strings.Contains(formatted, "step-2-task-claim.md") {
		t.Errorf("should contain contract path, got: %s", formatted)
	}
	if !strings.Contains(formatted, "claimed task <task_id>") {
		t.Errorf("should contain expected value, got: %s", formatted)
	}
	if !strings.Contains(formatted, "Task <task_id> claimed") {
		t.Errorf("should contain actual value, got: %s", formatted)
	}
}
