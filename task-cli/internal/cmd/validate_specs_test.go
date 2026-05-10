package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"task-cli/pkg/feature"
	"task-cli/pkg/task"
)

// setupValidateSpecsProject creates a test project with feature context and e2e spec files.
func setupValidateSpecsProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Create go.mod as project marker
	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Ensure feature directory structure
	if err := feature.EnsureFeatureDir(dir, "test-feature"); err != nil {
		t.Fatal(err)
	}

	// Create index with at least one task
	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test-feature"))
	index := &task.TaskIndex{
		Feature:      "test-feature",
		StatusEnum:   []string{"pending", "in_progress", "completed"},
		PriorityEnum: []string{"P0", "P1", "P2"},
	}
	index.SetTasks(map[string]task.Task{
		"task1": {ID: "1", Title: "Task 1", Priority: "P0", Status: "pending", File: "1-task.md"},
	})
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// Create task file
	taskFile := filepath.Join(dir, "docs", "features", "test-feature", "tasks", "1-task.md")
	if err := os.WriteFile(taskFile, []byte("task content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create e2e staging directory with a sample spec
	e2eDir := filepath.Join(dir, "tests", "e2e", "features", "test-feature")
	if err := os.MkdirAll(e2eDir, 0755); err != nil {
		t.Fatal(err)
	}

	return dir
}

// TestDiscoverSpecFiles verifies that spec files are discovered in the e2e directory.
func TestDiscoverSpecFiles(t *testing.T) {
	dir := setupValidateSpecsProject(t)

	// Create some spec files
	e2eDir := filepath.Join(dir, "tests", "e2e", "features", "test-feature")
	specFiles := []string{"login.spec.ts", "dashboard.spec.ts", "helpers.ts"}
	for _, f := range specFiles {
		if err := os.WriteFile(filepath.Join(e2eDir, f), []byte("// spec"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	files, err := discoverSpecFiles(filepath.Join(dir, "tests", "e2e", "features", "test-feature"))
	if err != nil {
		t.Fatalf("discoverSpecFiles failed: %v", err)
	}

	// Should only find .spec.ts files
	if len(files) != 2 {
		t.Errorf("expected 2 spec files, got %d: %v", len(files), files)
	}
	for _, f := range files {
		if !strings.HasSuffix(f, ".spec.ts") {
			t.Errorf("expected .spec.ts file, got %s", f)
		}
	}
}

// TestDiscoverSpecFiles_EmptyDir verifies empty directory returns error.
func TestDiscoverSpecFiles_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	emptyDir := filepath.Join(dir, "empty")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatal(err)
	}

	_, err := discoverSpecFiles(emptyDir)
	if err == nil {
		t.Error("expected error for empty directory")
	}
}

// TestDiscoverSpecFiles_NonexistentDir verifies nonexistent directory returns error.
func TestDiscoverSpecFiles_NonexistentDir(t *testing.T) {
	_, err := discoverSpecFiles("/nonexistent/path/to/specs")
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

// TestParseValidationOutput verifies parsing of JSON output from validate-specs.mjs.
func TestParseValidationOutput(t *testing.T) {
	tests := []struct {
		name         string
		jsonOutput   string
		wantErrors   int
		wantWarnings int
		wantErr      bool
	}{
		{
			name: "valid clean output",
			jsonOutput: `{
				"errors": [],
				"warnings": []
			}`,
			wantErrors:   0,
			wantWarnings: 0,
		},
		{
			name: "output with errors",
			jsonOutput: `{
				"errors": [
					{"rule": "E1", "file": "login.spec.ts", "line": 10, "message": "Forbidden waitForTimeout call"}
				],
				"warnings": []
			}`,
			wantErrors:   1,
			wantWarnings: 0,
		},
		{
			name: "output with warnings",
			jsonOutput: `{
				"errors": [],
				"warnings": [
					{"rule": "W1", "file": "login.spec.ts", "line": 5, "message": "Serial suite has 20 test() calls"}
				]
			}`,
			wantErrors:   0,
			wantWarnings: 1,
		},
		{
			name: "output with both errors and warnings",
			jsonOutput: `{
				"errors": [
					{"rule": "E1", "file": "a.spec.ts", "line": 1, "message": "E1 error"},
					{"rule": "E3", "file": "b.spec.ts", "line": 2, "message": "E3 error"}
				],
				"warnings": [
					{"rule": "W1", "file": "a.spec.ts", "line": 3, "message": "W1 warning"},
					{"rule": "W4", "file": "b.spec.ts", "line": 4, "message": "W4 warning"}
				]
			}`,
			wantErrors:   2,
			wantWarnings: 2,
		},
		{
			name:       "invalid JSON",
			jsonOutput: "not json",
			wantErr:    true,
		},
		{
			name: "extra fields ignored",
			jsonOutput: `{
				"errors": [{"rule": "E1", "file": "a.ts", "line": 1, "message": "err", "extra": "ignored"}],
				"warnings": [],
				"extra_field": true
			}`,
			wantErrors:   1,
			wantWarnings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseValidationOutput(tt.jsonOutput)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result.Errors) != tt.wantErrors {
				t.Errorf("got %d errors, want %d", len(result.Errors), tt.wantErrors)
			}
			if len(result.Warnings) != tt.wantWarnings {
				t.Errorf("got %d warnings, want %d", len(result.Warnings), tt.wantWarnings)
			}
		})
	}
}

// TestBuildValidateCommand verifies the command arguments are built correctly.
func TestBuildValidateCommand(t *testing.T) {
	scriptPath := "/path/to/validate-specs.mjs"
	specDir := "/project/tests/e2e/features/my-feature"
	testCasesPath := "/project/docs/features/my-feature/testing/test-cases.md"

	cmd := buildValidateCommand(scriptPath, specDir, testCasesPath)

	// exec.Command resolves "node" to the full path on some systems
	if !strings.HasSuffix(cmd.Path, "node") && !strings.HasSuffix(cmd.Path, "node.exe") {
		t.Errorf("expected node executable, got %s", cmd.Path)
	}

	// Args[0] is the executable name, then the script args follow
	// Expected: [node, scriptPath, specDir, --test-cases, testCasesPath]
	expectedArgs := []string{scriptPath, specDir, "--test-cases", testCasesPath}
	if len(cmd.Args) != len(expectedArgs)+1 {
		t.Fatalf("expected %d args, got %d: %v", len(expectedArgs)+1, len(cmd.Args), cmd.Args)
	}
	for i, arg := range expectedArgs {
		if cmd.Args[i+1] != arg {
			t.Errorf("arg[%d]: got %q, want %q", i+1, cmd.Args[i+1], arg)
		}
	}
}

// TestBuildValidateCommand_NoTestCases verifies command without test-cases flag.
func TestBuildValidateCommand_NoTestCases(t *testing.T) {
	scriptPath := "/path/to/validate-specs.mjs"
	specDir := "/project/tests/e2e/features/my-feature"

	cmd := buildValidateCommand(scriptPath, specDir, "")

	// Args[0] is executable name, then scriptPath and specDir
	expectedArgs := []string{scriptPath, specDir}
	if len(cmd.Args) != len(expectedArgs)+1 {
		t.Fatalf("expected %d args, got %d: %v", len(expectedArgs)+1, len(cmd.Args), cmd.Args)
	}
}

// TestFindValidateScript verifies the script path resolution.
func TestFindValidateScript(t *testing.T) {
	t.Run("script found in plugin templates", func(t *testing.T) {
		dir := t.TempDir()

		// Simulate the expected script location at plugins/forge/skills/gen-test-scripts/templates/validate-specs.mjs
		scriptDir := filepath.Join(dir, "plugins", "forge", "skills", "gen-test-scripts", "templates")
		if err := os.MkdirAll(scriptDir, 0755); err != nil {
			t.Fatal(err)
		}
		scriptPath := filepath.Join(scriptDir, "validate-specs.mjs")
		if err := os.WriteFile(scriptPath, []byte("// script"), 0644); err != nil {
			t.Fatal(err)
		}

		found, err := findValidateScript(dir)
		if err != nil {
			t.Fatalf("findValidateScript failed: %v", err)
		}
		if found != scriptPath {
			t.Errorf("expected %s, got %s", scriptPath, found)
		}
	})

	t.Run("script not found", func(t *testing.T) {
		dir := t.TempDir()
		_, err := findValidateScript(dir)
		if err == nil {
			t.Error("expected error when script not found")
		}
	})
}

// TestValidateSpecsOutput_PrintResults verifies the output printing logic.
func TestValidateSpecsOutput_PrintResults(t *testing.T) {
	tests := []struct {
		name           string
		result         *validationResult
		expectedInOutput []string
		notExpected     []string
	}{
		{
			name: "clean output",
			result: &validationResult{
				Errors:   []validationEntry{},
				Warnings: []validationEntry{},
			},
			expectedInOutput: []string{"RESULT: PASS"},
			notExpected:      []string{"[ERRORS]", "[WARNINGS]"},
		},
		{
			name: "only errors",
			result: &validationResult{
				Errors: []validationEntry{
					{Rule: "E1", File: "a.spec.ts", Line: 10, Message: "waitForTimeout"},
				},
				Warnings: []validationEntry{},
			},
			expectedInOutput: []string{"[ERRORS]", "E1", "waitForTimeout", "RESULT: FAIL"},
			notExpected:      []string{"[WARNINGS]"},
		},
		{
			name: "only warnings",
			result: &validationResult{
				Errors: []validationEntry{},
				Warnings: []validationEntry{
					{Rule: "W1", File: "b.spec.ts", Line: 5, Message: "serial suite too large"},
				},
			},
			expectedInOutput: []string{"[WARNINGS]", "W1", "serial suite too large", "RESULT: PASS"},
			notExpected:      []string{"[ERRORS]"},
		},
		{
			name: "errors and warnings",
			result: &validationResult{
				Errors: []validationEntry{
					{Rule: "E1", File: "a.spec.ts", Line: 1, Message: "err1"},
				},
				Warnings: []validationEntry{
					{Rule: "W1", File: "b.spec.ts", Line: 2, Message: "warn1"},
				},
			},
			expectedInOutput: []string{"[ERRORS]", "[WARNINGS]", "RESULT: FAIL"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureValidationOutput(tt.result)

			for _, expected := range tt.expectedInOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("expected output to contain %q, got:\n%s", expected, output)
				}
			}
			for _, notExpected := range tt.notExpected {
				if strings.Contains(output, notExpected) {
					t.Errorf("expected output NOT to contain %q, got:\n%s", notExpected, output)
				}
			}
		})
	}
}

// TestRunValidateSpecs_Integration tests the full command with a mock Node script.
func TestRunValidateSpecs_Integration(t *testing.T) {
	dir := setupValidateSpecsProject(t)

	// Create e2e spec files
	e2eDir := filepath.Join(dir, "tests", "e2e", "features", "test-feature")
	specContent := `import { test } from '@playwright/test';
// Traceability: TC-001
test('example test', async ({ page }) => {
	await page.goto('/');
});
`
	if err := os.WriteFile(filepath.Join(e2eDir, "login.spec.ts"), []byte(specContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create test-cases.md
	testCasesDir := filepath.Join(dir, "docs", "features", "test-feature", "testing")
	if err := os.MkdirAll(testCasesDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testCasesDir, "test-cases.md"), []byte("# Test Cases\n## TC-001 Login\nTest login"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a mock validate-specs.mjs script
	mockScript := filepath.Join(dir, "mock-validate-specs.mjs")
	mockContent := `const result = { errors: [], warnings: [] };
console.log(JSON.stringify(result));
process.exit(0);
`
	if err := os.WriteFile(mockScript, []byte(mockContent), 0755); err != nil {
		t.Fatal(err)
	}

	// Test the execution with mock script
	cmd := buildValidateCommand(mockScript, e2eDir, filepath.Join(testCasesDir, "test-cases.md"))
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("mock script failed: %v, output: %s", err, output)
	}

	var result validationResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("failed to parse output: %v", err)
	}

	if len(result.Errors) != 0 {
		t.Errorf("expected 0 errors, got %d", len(result.Errors))
	}
}

// TestRunValidateSpecs_Integration_WithErrors tests handling of script that reports errors.
func TestRunValidateSpecs_Integration_WithErrors(t *testing.T) {
	dir := setupValidateSpecsProject(t)
	e2eDir := filepath.Join(dir, "tests", "e2e", "features", "test-feature")
	if err := os.WriteFile(filepath.Join(e2eDir, "bad.spec.ts"), []byte("// bad spec"), 0644); err != nil {
		t.Fatal(err)
	}

	// Mock script that reports errors
	mockScript := filepath.Join(dir, "mock-validate-specs.mjs")
	mockContent := `const result = { errors: [{rule: "E1", file: "bad.spec.ts", line: 5, message: "waitForTimeout"}], warnings: [{rule: "W1", file: "bad.spec.ts", line: 1, message: "too many tests"}] };
console.log(JSON.stringify(result));
process.exit(1);
`
	if err := os.WriteFile(mockScript, []byte(mockContent), 0755); err != nil {
		t.Fatal(err)
	}

	cmd := buildValidateCommand(mockScript, e2eDir, "")
	output, _ := cmd.CombinedOutput()

	var result validationResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("failed to parse output: %v", err)
	}

	if len(result.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Rule != "E1" {
		t.Errorf("expected E1, got %s", result.Errors[0].Rule)
	}
	if len(result.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(result.Warnings))
	}
}

// TestRunValidateSpecs_ScriptNotFound tests graceful degradation when script is missing.
func TestRunValidateSpecs_ScriptNotFound(t *testing.T) {
	dir := setupValidateSpecsProject(t)

	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	// Temporarily override exitFunc to prevent os.Exit
	oldExit := exitFunc
	defer func() { exitFunc = oldExit }()

	var exited bool
	var exitCode int
	exitFunc = func(code int) {
		exited = true
		exitCode = code
	}

	// Create a non-existent script path
	_ = runValidateSpecsInternal("/nonexistent/validate-specs.mjs", filepath.Join(dir, "tests", "e2e", "features", "test-feature"), "")

	// Should NOT exit with error — graceful degradation
	if exited {
		t.Errorf("expected graceful degradation (no exit), but exited with code %d", exitCode)
	}
}

// TestGetSpecDir verifies the spec directory path construction.
func TestGetSpecDir(t *testing.T) {
	dir := "tests/e2e/features"
	slug := "my-feature"
	result := getSpecDir(dir, slug)

	expected := filepath.Join(dir, slug)
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

// TestGetTestCasesPath verifies the test-cases.md path construction.
func TestGetTestCasesPath(t *testing.T) {
	tests := []struct {
		name        string
		projectRoot string
		slug        string
		expected    string
	}{
		{
			name:        "standard path",
			projectRoot: "/project",
			slug:        "my-feature",
			expected:    filepath.Join("/project", "docs", "features", "my-feature", "testing", "test-cases.md"),
		},
		{
			name:        "empty slug",
			projectRoot: "/project",
			slug:        "",
			expected:    filepath.Join("/project", "docs", "features", "testing", "test-cases.md"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTestCasesPath(tt.projectRoot, tt.slug)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
