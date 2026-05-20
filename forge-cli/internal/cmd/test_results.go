package cmd

import (
	"os"
	"path/filepath"
)

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
func writeRegressionRawOutput(projectRoot, output string) error {
	resultsDir := filepath.Join(projectRoot, "tests", "e2e", "results")
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return err
	}
	rawPath := filepath.Join(resultsDir, "raw-output.txt")
	return os.WriteFile(rawPath, []byte(output), 0644)
}
