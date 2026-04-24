package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"task-cli/pkg/feature"
)

// TestFailure represents a single test failure.
type TestFailure struct {
	TestName     string // Human-readable test name
	TestCaseID   string // Test case ID (e.g., "ui/login/login-with-valid-credentials")
	File         string // Source file (if available)
	Line         int    // Line number (if available)
	ErrorMessage string
	Output       string // Full relevant output
	StackTrace   string
}

// TestStats holds aggregated test statistics.
type TestStats struct {
	Total     int
	Pass      int
	Fail      int
	Skip      int
	Framework string
}

// parseTestFailures extracts test failures from output (framework-agnostic).
func parseTestFailures(output string) []TestFailure {
	var failures []TestFailure
	seen := make(map[string]bool) // Deduplicate by test name

	// Common failure patterns across frameworks
	patterns := []struct {
		regex   *regexp.Regexp
		extract func([]string) (name, file string, line int)
	}{
		// npm/jest: FAIL testing/scripts/ui.spec.ts (line in stack trace)
		{
			regex: regexp.MustCompile(`✗\s+(.+?)\s*$`),
			extract: func(match []string) (string, string, int) {
				return strings.TrimSpace(match[1]), "", 0
			},
		},
		// Go: --- FAIL: TestName (0.00s)
		{
			regex: regexp.MustCompile(`---\s+FAIL:\s+(.+?)\s+\(`),
			extract: func(match []string) (string, string, int) {
				return strings.TrimSpace(match[1]), "", 0
			},
		},
		// pytest: FAILED test_file.py::test_name
		{
			regex: regexp.MustCompile(`FAILED\s+(.+?)::(.+)`),
			extract: func(match []string) (string, string, int) {
				return strings.TrimSpace(match[2]), strings.TrimSpace(match[1]), 0
			},
		},
		// Generic: FAIL test_name
		{
			regex: regexp.MustCompile(`FAIL[^:]*:\s*(.+?)\s*$`),
			extract: func(match []string) (string, string, int) {
				return strings.TrimSpace(match[1]), "", 0
			},
		},
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		for _, p := range patterns {
			if matches := p.regex.FindStringSubmatch(line); matches != nil {
				name, file, line := p.extract(matches)

				// Deduplicate
				if seen[name] {
					continue
				}
				seen[name] = true

				// Extract error message and stack trace from surrounding context
				errorMsg, stackTrace := extractErrorContext(output, lineNum)

				failures = append(failures, TestFailure{
					TestName:     name,
					File:         file,
					Line:         line,
					ErrorMessage: errorMsg,
					StackTrace:   stackTrace,
					Output:       extractRelevantOutput(output, lineNum, 50),
				})
				break
			}
		}
	}

	return failures
}

// extractErrorContext extracts error message and stack trace from surrounding lines.
func extractErrorContext(output string, failLineNum int) (errorMsg, stackTrace string) {
	lines := strings.Split(output, "\n")
	if failLineNum >= len(lines) {
		return "", ""
	}

	var errorLines []string
	var stackLines []string
	inError := false
	inStack := false

	// Look ahead for error message and stack trace
	for i := failLineNum; i < len(lines) && i < failLineNum+30; i++ {
		line := lines[i]

		// Error message patterns
		if strings.Contains(line, "Error:") ||
			strings.Contains(line, "AssertionError") ||
			strings.Contains(line, "expected") ||
			strings.Contains(line, "Expected") {
			inError = true
			errorLines = append(errorLines, line)
			continue
		}

		// Stack trace patterns
		if strings.Contains(line, "at ") ||
			strings.Contains(line, "File \"") ||
			strings.HasPrefix(line, "  ") {
			inStack = true
			stackLines = append(stackLines, line)
			continue
		}

		// End of error/stack section
		if inError && line == "" {
			inError = false
		}
		if inStack && line == "" {
			break
		}
	}

	return strings.Join(errorLines, "\n"), strings.Join(stackLines, "\n")
}

// extractRelevantOutput extracts N lines around the failure.
func extractRelevantOutput(output string, failLineNum int, maxLines int) string {
	lines := strings.Split(output, "\n")
	start := failLineNum - 5
	if start < 0 {
		start = 0
	}
	end := failLineNum + maxLines
	if end > len(lines) {
		end = len(lines)
	}

	relevant := lines[start:end]
	return strings.Join(relevant, "\n")
}

// matchTestCaseID attempts to match test name to test case ID from test-cases.md.
func matchTestCaseID(testName, testCasesPath string) string {
	// If test-cases.md doesn't exist, return sanitized test name
	if _, err := os.Stat(testCasesPath); os.IsNotExist(err) {
		return sanitizeTestName(testName)
	}

	data, err := os.ReadFile(testCasesPath)
	if err != nil {
		return sanitizeTestName(testName)
	}

	// Build map of test names to test case IDs
	// Format in test-cases.md:
	// ## TC-001: Login with valid credentials
	// - **Test ID**: ui/login/login-with-valid-credentials
	testNameToID := make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	var currentTestID string

	for scanner.Scan() {
		line := scanner.Text()

		// Extract Test ID
		if strings.Contains(line, "**Test ID**") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				currentTestID = strings.TrimSpace(parts[1])
			}
		}

		// Extract test name from title
		if strings.HasPrefix(line, "## TC-") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				title := strings.TrimSpace(parts[1])
				if currentTestID != "" {
					testNameToID[title] = currentTestID
					testNameToID[strings.ToLower(title)] = currentTestID
				}
			}
		}
	}

	// Try exact match
	if id, ok := testNameToID[testName]; ok {
		return id
	}

	// Try case-insensitive match
	if id, ok := testNameToID[strings.ToLower(testName)]; ok {
		return id
	}

	// Fallback to sanitized test name
	return sanitizeTestName(testName)
}

// sanitizeTestName converts test name to valid file name.
func sanitizeTestName(name string) string {
	// Lowercase
	name = strings.ToLower(name)
	// Replace spaces with hyphens
	name = strings.ReplaceAll(name, " ", "-")
	// Remove special characters
	reg := regexp.MustCompile(`[^a-z0-9\-/]`)
	name = reg.ReplaceAllString(name, "")
	// Replace / with - (for test case IDs)
	name = strings.ReplaceAll(name, "/", "-")
	return name
}

// writeLatestMd writes test result overview to latest.md.
func writeLatestMd(projectRoot, featureSlug string, stats TestStats, failures []TestFailure) error {
	resultsDir := filepath.Join(projectRoot, feature.GetFeatureTestingResultsDir(featureSlug))
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return err
	}

	latestPath := filepath.Join(resultsDir, "latest.md")

	status := "PASS"
	if stats.Fail > 0 {
		status = "FAIL"
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("# Test Results: %s\n\n", featureSlug))
	buf.WriteString(fmt.Sprintf("**Date**: %s\n", time.Now().Format("2006-01-02 15:04")))
	buf.WriteString(fmt.Sprintf("**Status**: %s\n", status))
	buf.WriteString(fmt.Sprintf("**Total**: %d tests, %d passed, %d failed\n\n",
		stats.Total, stats.Pass, stats.Fail))

	if len(failures) > 0 {
		buf.WriteString("## Failures\n\n")
		for _, f := range failures {
			failureFile := fmt.Sprintf("failures/failure-%s.md", f.TestCaseID)
			buf.WriteString(fmt.Sprintf("- [%s](%s) — %s\n", failureFile, failureFile, f.TestName))
		}
	}

	return os.WriteFile(latestPath, []byte(buf.String()), 0644)
}

// writeFailureFiles writes individual failure files.
func writeFailureFiles(projectRoot, featureSlug string, failures []TestFailure) error {
	if len(failures) == 0 {
		return nil
	}

	resultsDir := filepath.Join(projectRoot, feature.GetFeatureTestingResultsDir(featureSlug))
	failuresDir := filepath.Join(resultsDir, "failures")
	if err := os.MkdirAll(failuresDir, 0755); err != nil {
		return err
	}

	for _, f := range failures {
		failurePath := filepath.Join(failuresDir, fmt.Sprintf("failure-%s.md", f.TestCaseID))

		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("# Failure: %s\n\n", f.TestName))
		buf.WriteString(fmt.Sprintf("**Test Case ID**: %s\n", f.TestCaseID))
		if f.File != "" {
			buf.WriteString(fmt.Sprintf("**File**: %s", f.File))
			if f.Line > 0 {
				buf.WriteString(fmt.Sprintf(":%d", f.Line))
			}
			buf.WriteString("\n")
		}
		buf.WriteString(fmt.Sprintf("**Generated**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

		if f.ErrorMessage != "" {
			buf.WriteString("## Error\n\n```\n")
			buf.WriteString(f.ErrorMessage)
			buf.WriteString("\n```\n\n")
		}

		if f.Output != "" {
			buf.WriteString("## Output\n\n```\n")
			buf.WriteString(f.Output)
			buf.WriteString("\n```\n\n")
		}

		if f.StackTrace != "" {
			buf.WriteString("## Stack Trace\n\n```\n")
			buf.WriteString(f.StackTrace)
			buf.WriteString("\n```\n")
		}

		if err := os.WriteFile(failurePath, []byte(buf.String()), 0644); err != nil {
			return err
		}
	}

	return nil
}
