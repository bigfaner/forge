package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"task-cli/pkg/feature"
)

// TestStats holds aggregated test statistics.
type TestStats struct {
	Total     int
	Pass      int
	Fail      int
	Skip      int
	Framework string
}

// writeLatestMd writes a test result summary to latest.md.
// When stats.Fail > 0, it reports FAIL with a reference to raw-output.txt.
// When stats.Fail == 0, it reports PASS.
func writeLatestMd(projectRoot, featureSlug string, stats TestStats) error {
	resultsDir := filepath.Join(projectRoot, feature.GetFeatureTestingResultsDir(featureSlug))
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return err
	}

	latestPath := filepath.Join(resultsDir, "latest.md")

	status := "PASS"
	if stats.Fail > 0 {
		status = "FAIL"
	}

	var buf string
	if status == "FAIL" {
		buf = fmt.Sprintf("# Test Results: %s\n\n", featureSlug)
		buf += fmt.Sprintf("**Date**: %s\n", time.Now().Format("2006-01-02 15:04"))
		buf += fmt.Sprintf("**Status**: %s\n\n", status)
		buf += "Tests failed. See `raw-output.txt` for full output.\n\n"
		buf += "## Next Steps\n\n"
		buf += "1. Read `testing/results/raw-output.txt`\n"
		buf += "2. Analyze failures and determine root causes\n"
		buf += "3. Use `task add --title \"Fix: <description>\" --priority P0 --breaking` for each issue\n"
		buf += "4. Run `task claim` to pick up fix tasks\n"
	} else {
		buf = fmt.Sprintf("# Test Results: %s\n\n", featureSlug)
		buf += fmt.Sprintf("**Date**: %s\n", time.Now().Format("2006-01-02 15:04"))
		buf += fmt.Sprintf("**Status**: %s\n", status)
		buf += fmt.Sprintf("**Total**: %d tests, %d passed\n", stats.Total, stats.Pass)
	}

	return os.WriteFile(latestPath, []byte(buf), 0644)
}

// writeRawOutput saves the raw test output to a file for agent analysis.
func writeRawOutput(projectRoot, featureSlug string, output string) error {
	resultsDir := filepath.Join(projectRoot, feature.GetFeatureTestingResultsDir(featureSlug))
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return err
	}
	rawPath := filepath.Join(resultsDir, "raw-output.txt")
	return os.WriteFile(rawPath, []byte(output), 0644)
}

// writeUnitTestRawOutput saves project-wide unit test output to tests/results/unit-raw-output.txt.
func writeUnitTestRawOutput(projectRoot, output string) error {
	resultsDir := filepath.Join(projectRoot, "tests", "results")
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return err
	}
	rawPath := filepath.Join(resultsDir, "unit-raw-output.txt")
	return os.WriteFile(rawPath, []byte(output), 0644)
}

// writeRegressionRawOutput saves project-wide regression output to tests/e2e/results/raw-output.txt.
// Kept separate from writeRawOutput because regression failures span all features, not one.
func writeRegressionRawOutput(projectRoot, output string) error {
	resultsDir := filepath.Join(projectRoot, "tests", "e2e", "results")
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return err
	}
	rawPath := filepath.Join(resultsDir, "raw-output.txt")
	return os.WriteFile(rawPath, []byte(output), 0644)
}
