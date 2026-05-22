package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var validateSpecsCmd = &cobra.Command{
	Use:   "validate-specs",
	Short: "Validate generated spec files against structural rules",
	Long: `Run ts-morph-based AST validation against generated Playwright spec files.

Discovers spec files in tests/e2e/features/<slug>/ and runs validate-specs.mjs
against them. Reports errors (blocking) and warnings (non-blocking).

Exit codes:
  0 - No errors (warnings OK)
  1 - Validation errors found
  2 - Script failed to run or prerequisites missing`,
	Args: cobra.NoArgs,
	RunE: runValidateSpecs,
}

// validationResult represents the JSON output from validate-specs.mjs.
type validationResult struct {
	Errors   []validationEntry `json:"errors"`
	Warnings []validationEntry `json:"warnings"`
}

// validationEntry represents a single error or warning from the validation script.
type validationEntry struct {
	Rule    string `json:"rule"`
	File    string `json:"file"`
	Line    int    `json:"line"`
	Message string `json:"message"`
}

// exitFunc allows tests to override os.Exit.
var exitFunc = os.Exit

// validateSpecsScriptPath is the relative path to the validation script from project root.
const validateSpecsScriptRelPath = "plugins/forge/skills/gen-test-scripts/templates/validate-specs.mjs"

func runValidateSpecs(_ *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	slug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		return ErrFeatureNotSet()
	}

	// Find the validation script
	scriptPath, err := findValidateScript(projectRoot)
	if err != nil {
		fmt.Println("WARNING: validate-specs.mjs not found — skipping spec validation")
		fmt.Printf("  HINT: Ensure %s exists in the project\n", validateSpecsScriptRelPath)
		return nil
	}

	// Discover spec files
	specDir := filepath.Join(projectRoot, feature.E2EStagingDir, slug)
	specFiles, err := discoverSpecFiles(specDir)
	if err != nil {
		return NewAIError(
			ErrNotFound,
			"No spec files found",
			fmt.Sprintf("Could not find spec files in %s: %v", specDir, err),
			"Generate spec files first using gen-test-scripts",
			"forge task check-deps",
		)
	}

	fmt.Printf("Validating %d spec file(s) in %s\n", len(specFiles), specDir)

	// Auto-detect test-cases.md path
	testCasesPath := getTestCasesPath(projectRoot, slug)
	if _, err := os.Stat(testCasesPath); err != nil {
		testCasesPath = "" // Not available, skip E2 check
	}

	exitCode := runValidateSpecsInternal(scriptPath, specDir, testCasesPath)
	if exitCode != 0 {
		return fmt.Errorf("spec validation failed with exit code %d", exitCode)
	}
	return nil
}

// runValidateSpecsInternal executes the validation script and prints results.
// Returns exit code: 0 if no errors, 1 if errors, 2 if script fails.
func runValidateSpecsInternal(scriptPath, specDir, testCasesPath string) int {
	nodeCmd := buildValidateCommand(scriptPath, specDir, testCasesPath)

	output, err := nodeCmd.CombinedOutput()
	if err != nil {
		// Check if node is not found
		if isNodeNotFound(err) {
			fmt.Println("WARNING: Node.js not found in PATH — skipping spec validation")
			fmt.Println("  HINT: Install Node.js to enable spec validation")
			return 0 // Graceful degradation
		}

		// Script exited with non-zero — may still have valid JSON output
		if len(output) == 0 {
			fmt.Printf("ERROR: validate-specs.mjs failed: %v\n", err)
			return 2
		}
	}

	// Parse JSON output
	result, parseErr := parseValidationOutput(string(output))
	if parseErr != nil {
		// Script produced non-JSON output — might be a ts-morph import error
		errMsg := string(output)
		if isTSMorphError(errMsg) {
			fmt.Println("WARNING: ts-morph not available — skipping spec validation")
			fmt.Println("  HINT: Run 'npm install' in tests/e2e/ to install ts-morph")
			return 0 // Graceful degradation per proposal risk mitigation
		}
		fmt.Printf("ERROR: Failed to parse validation output: %v\n", parseErr)
		fmt.Printf("  Output: %s\n", errMsg)
		return 2
	}

	// Print results
	printValidationResults(result, specDir)

	if len(result.Errors) > 0 {
		return 1
	}
	return 0
}

// discoverSpecFiles finds all *.spec.ts files in the given directory.
func discoverSpecFiles(dir string) ([]string, error) {
	pattern := filepath.Join(dir, "*.spec.ts")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob spec files: %w", err)
	}

	if len(files) == 0 {
		// Check if directory exists
		if _, statErr := os.Stat(dir); os.IsNotExist(statErr) {
			return nil, fmt.Errorf("spec directory not found: %s", dir)
		}
		return nil, fmt.Errorf("no .spec.ts files found in %s", dir)
	}

	return files, nil
}

// parseValidationOutput parses the JSON output from validate-specs.mjs.
func parseValidationOutput(raw string) (*validationResult, error) {
	var result validationResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return &result, nil
}

// buildValidateCommand constructs the exec.Cmd to run the validation script.
func buildValidateCommand(scriptPath, specDir, testCasesPath string) *exec.Cmd {
	args := []string{scriptPath, specDir}
	if testCasesPath != "" {
		args = append(args, "--test-cases", testCasesPath)
	}
	return exec.Command("node", args...)
}

// findValidateScript locates validate-specs.mjs in the project.
func findValidateScript(projectRoot string) (string, error) {
	scriptPath := filepath.Join(projectRoot, validateSpecsScriptRelPath)
	if _, err := os.Stat(scriptPath); err != nil {
		return "", fmt.Errorf("validate-specs.mjs not found at %s", scriptPath)
	}
	return scriptPath, nil
}

// getSpecDir returns the e2e spec directory path for a feature.
func getSpecDir(baseDir, slug string) string {
	return filepath.Join(baseDir, slug)
}

// getTestCasesPath returns the path to test-cases.md for a feature.
func getTestCasesPath(projectRoot, slug string) string {
	return filepath.Join(projectRoot, feature.GetFeatureTestCasesFile(slug))
}

// printValidationResults prints human-readable validation results.
func printValidationResults(result *validationResult, specDir string) {
	if len(result.Warnings) > 0 {
		PrintSection("WARNINGS")
		for _, w := range result.Warnings {
			PrintListItem(fmt.Sprintf("[%s] %s:%d %s", w.Rule, w.File, w.Line, w.Message))
		}
	}

	if len(result.Errors) > 0 {
		PrintSection("ERRORS")
		for _, e := range result.Errors {
			PrintListItem(fmt.Sprintf("[%s] %s:%d %s", e.Rule, e.File, e.Line, e.Message))
		}
	}

	if len(result.Errors) == 0 {
		PrintResult("PASS", fmt.Sprintf("%s (0 errors, %d warnings)", specDir, len(result.Warnings)))
	} else {
		PrintResult("FAIL", fmt.Sprintf("%s (%d errors, %d warnings)", specDir, len(result.Errors), len(result.Warnings)))
	}
}

// captureValidationOutput captures the output of printValidationResults.
func captureValidationOutput(result *validationResult) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printValidationResults(result, "test-dir")

	_ = w.Close()
	os.Stdout = old

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	return string(buf[:n])
}

// isNodeNotFound checks if the error indicates node is not in PATH.
func isNodeNotFound(err error) bool {
	if exitErr, ok := err.(*exec.Error); ok {
		return exitErr.Err == exec.ErrNotFound
	}
	return false
}

// isTSMorphError checks if the output indicates a ts-morph import failure.
func isTSMorphError(output string) bool {
	return containsAny(output, "Cannot find module 'ts-morph'", "ERR_MODULE_NOT_FOUND", "ts-morph")
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}
