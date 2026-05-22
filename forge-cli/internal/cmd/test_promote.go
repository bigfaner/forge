package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/pkg/project"
	"forge-cli/pkg/testrunner"

	"github.com/spf13/cobra"
)

// runTestPromote promotes a journey's @feature tags to @regression.
// It first runs the journey's tests, then on success replaces @feature with @regression
// in all test files under the journey directory.
func runTestPromote(_ *cobra.Command, args []string) error {
	journeyName := args[0]

	if err := validateJourneyName(journeyName); err != nil {
		return err
	}

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	// Validate the journey directory exists
	journeyDir := filepath.Join(projectRoot, "tests", journeyName)
	if _, err := os.Stat(journeyDir); os.IsNotExist(err) {
		Exit(NewAIError(ErrValidation,
			fmt.Sprintf("Journey %q not found", journeyName),
			fmt.Sprintf("Expected directory: %s", journeyDir),
			"List available journeys with: ls tests/",
			"forge test promote <journey-name>"))
	}

	// Run the journey's tests first
	cfg := testrunner.ResolveJourneyExecutionConfig(projectRoot)

	workDir, cleanup, err := testrunner.CreateJourneyWorkDir(projectRoot, journeyName)
	if err != nil {
		return NewAIError(ErrValidation, "Failed to create journey work directory", err.Error(),
			"Check temp directory permissions", "forge test promote "+journeyName)
	}

	result := testrunner.ExecuteJourneyInIsolation(cfg, workDir, journeyName)

	if !result.Passed {
		cleanup()
		Exit(NewAIError(ErrValidation,
			"Journey tests failed, promotion refused",
			"One or more tests in the journey did not pass",
			"Fix the failing tests before promoting",
			"forge test run-journey "+journeyName))
	}

	cleanup()

	// Find all test files under the journey directory and replace @feature with @regression
	filesModified, err := promoteJourneyTags(journeyDir)
	if err != nil {
		return NewAIError(ErrValidation, "Failed to promote tags", err.Error(),
			"Check file permissions in journey directory",
			"forge test promote "+journeyName)
	}

	PrintBlockStart()
	PrintField("JOURNEY", journeyName)
	PrintField("RESULT", "PROMOTED")
	PrintField("FILES_MODIFIED", fmt.Sprintf("%d", filesModified))
	PrintField("TAG_CHANGE", "@feature -> @regression")
	PrintBlockEnd()
	return nil
}

// validateJourneyName checks that the journey name does not contain path traversal.
// Uses filepath.Base() and rejects ".." components per Hard Rules.
func validateJourneyName(name string) *AIError {
	if filepath.Base(name) != name || strings.Contains(name, "..") {
		return NewErrInvalidPath(name)
	}
	return nil
}

// promoteJourneyTags walks the journey directory, finds test files, and replaces
// @feature tags with @regression tags. Returns the number of files modified.
func promoteJourneyTags(journeyDir string) (int, error) {
	filesModified := 0

	err := filepath.Walk(journeyDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip _contracts directory
			if filepath.Base(path) == "_contracts" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process test files (common patterns)
		if !isTestFile(path) {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file %s: %w", path, err)
		}

		content := string(data)
		if !strings.Contains(content, "@feature") {
			return nil
		}

		// Replace @feature with @regression
		newContent := replaceFeatureTag(content, filepath.Ext(path))

		if newContent != content {
			if err := os.WriteFile(path, []byte(newContent), info.Mode()); err != nil {
				return fmt.Errorf("write file %s: %w", path, err)
			}
			filesModified++
		}

		return nil
	})

	return filesModified, err
}

// isTestFile checks if a file is a test file based on common naming patterns.
func isTestFile(path string) bool {
	base := filepath.Base(path)
	ext := filepath.Ext(path)

	switch ext {
	case ".go":
		return strings.HasSuffix(base, "_test.go")
	case ".py":
		return strings.HasPrefix(base, "test_") || strings.HasSuffix(base, "_test.py")
	case ".ts", ".js":
		return strings.Contains(base, ".test.") || strings.Contains(base, ".spec.")
	case ".java":
		return strings.HasSuffix(base, "Test.java") || strings.HasSuffix(base, "Tests.java")
	case ".rs":
		return strings.Contains(base, "_test") || strings.HasSuffix(base, "_tests.rs")
	}
	return false
}

// replaceFeatureTag replaces @feature with @regression using the appropriate
// syntax for the file's language.
func replaceFeatureTag(content, ext string) string {
	switch ext {
	case ".go":
		// Go: replace in build tags and comments
		// //go:build feature -> //go:build regression
		// // +build feature -> // +build regression
		result := content
		result = strings.ReplaceAll(result, "//go:build feature", "//go:build regression")
		result = strings.ReplaceAll(result, "// +build feature", "// +build regression")
		result = strings.ReplaceAll(result, "@feature", "@regression")
		return result
	case ".py":
		// Python: @pytest.mark.feature -> @pytest.mark.regression
		result := content
		result = strings.ReplaceAll(result, "@pytest.mark.feature", "@pytest.mark.regression")
		result = strings.ReplaceAll(result, "@feature", "@regression")
		return result
	case ".ts", ".js":
		// Playwright/JS: @feature -> @regression in test annotations
		return strings.ReplaceAll(content, "@feature", "@regression")
	case ".java":
		// Java: @Tag("feature") -> @Tag("regression")
		result := content
		result = strings.ReplaceAll(result, `@Tag("feature")`, `@Tag("regression")`)
		result = strings.ReplaceAll(result, "@feature", "@regression")
		return result
	case ".rs":
		// Rust: #[cfg(feature = "feature")] -> #[cfg(feature = "regression")]
		result := content
		result = strings.ReplaceAll(result, `#[cfg(feature = "feature")]`, `#[cfg(feature = "regression")]`)
		result = strings.ReplaceAll(result, "@feature", "@regression")
		return result
	default:
		// Generic: just replace @feature with @regression
		return strings.ReplaceAll(content, "@feature", "@regression")
	}
}

// PromoteDiffSummary generates a summary of what would change in a promote operation
// without actually modifying files. Used for dry-run verification.
func PromoteDiffSummary(journeyDir string) (bytes.Buffer, error) {
	var buf bytes.Buffer

	err := filepath.Walk(journeyDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if filepath.Base(path) == "_contracts" {
				return filepath.SkipDir
			}
			return nil
		}

		if !isTestFile(path) {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content := string(data)
		if strings.Contains(content, "@feature") {
			relPath, _ := filepath.Rel(journeyDir, path)
			newContent := replaceFeatureTag(content, filepath.Ext(path))

			// Show diff-style output
			oldLines := strings.Split(content, "\n")
			newLines := strings.Split(newContent, "\n")

			for i, oldLine := range oldLines {
				if i < len(newLines) && oldLine != newLines[i] {
					fmt.Fprintf(&buf, "--- %s\n", relPath)
					fmt.Fprintf(&buf, "+++ %s\n", relPath)
					fmt.Fprintf(&buf, "-%s\n", oldLine)
					fmt.Fprintf(&buf, "+%s\n", newLines[i])
				}
			}
		}

		return nil
	})

	return buf, err
}
