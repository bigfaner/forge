package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"forge-cli/pkg/forgeconfig"
)

// ContractFailure records a single contract dimension failure.
type ContractFailure struct {
	Dimension    string `json:"dimension"`    // Preconditions, Input, Output, State, Side-effect, Invariants
	ContractPath string `json:"contractPath"` // Path to the contract spec file
	Expected     string `json:"expected"`     // Expected value
	Actual       string `json:"actual"`       // Actual value observed
}

// Format returns a human-readable description of this contract failure.
func (cf ContractFailure) Format() string {
	return fmt.Sprintf("%s dimension FAILED: %s, expected '%s', got '%s'",
		cf.Dimension, cf.ContractPath, cf.Expected, cf.Actual)
}

// JourneyResult holds the execution result of a single journey.
type JourneyResult struct {
	JourneyName string            `json:"journeyName"`
	Passed      bool              `json:"passed"`
	Duration    time.Duration     `json:"duration"`
	ExitCode    int               `json:"exitCode"`
	Output      string            `json:"output,omitempty"`
	Error       string            `json:"error,omitempty"`
	Failures    []ContractFailure `json:"failures,omitempty"`
}

// FormatReport returns a structured text report for the journey result.
func (r JourneyResult) FormatReport() string {
	var buf strings.Builder

	status := "PASS"
	if !r.Passed {
		status = "FAIL"
	}

	buf.WriteString("---\n")
	fmt.Fprintf(&buf, "JOURNEY: %s\n", r.JourneyName)
	fmt.Fprintf(&buf, "RESULT: %s\n", status)
	fmt.Fprintf(&buf, "DURATION: %s\n", r.Duration.Round(time.Millisecond))
	if r.ExitCode != 0 {
		fmt.Fprintf(&buf, "EXIT_CODE: %d\n", r.ExitCode)
	}

	if len(r.Failures) > 0 {
		buf.WriteString("\nFAILURES:\n")
		for _, f := range r.Failures {
			fmt.Fprintf(&buf, "  - %s\n", f.Format())
		}
	}

	if r.Error != "" {
		fmt.Fprintf(&buf, "ERROR: %s\n", r.Error)
	}

	buf.WriteString("---\n")
	return buf.String()
}

// JourneyExecutionConfig holds resolved configuration for journey execution.
type JourneyExecutionConfig struct {
	TestCommand string
	Language    string
}

// resolveJourneyExecutionConfig reads test execution settings from project config.
func resolveJourneyExecutionConfig(projectRoot string) (*JourneyExecutionConfig, error) {
	// Read test-command from config
	testCmd, err := readTestCommand(projectRoot)
	if err != nil {
		return nil, err
	}

	return &JourneyExecutionConfig{
		TestCommand: testCmd,
	}, nil
}

// readTestCommand reads the test-command field from .forge/config.yaml.
func readTestCommand(projectRoot string) (string, error) {
	cfg, err := forgeconfig.ReadConfig(projectRoot)
	if err != nil {
		return "", fmt.Errorf("failed to read config: %w", err)
	}
	if cfg == nil {
		return "", fmt.Errorf("no .forge/config.yaml found; set test-command in config")
	}
	if cfg.TestCommand == "" {
		return "", fmt.Errorf("test-command not set in .forge/config.yaml; add 'test-command: go test ./...' (or your project's test command)")
	}
	return cfg.TestCommand, nil
}

// createJourneyWorkDir creates an isolated temporary directory for a journey.
// The directory name includes the journey name and a random suffix.
// Returns the work dir path, a cleanup function, and any error.
func createJourneyWorkDir(_, journeyName string) (string, func(), error) {
	// Create a temp dir in the system temp location (not inside project)
	prefix := fmt.Sprintf("forge-journey-%s-", journeyName)
	workDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		return "", nil, fmt.Errorf("create temp dir for journey %s: %w", journeyName, err)
	}

	cleanup := func() {
		_ = os.RemoveAll(workDir)
	}

	return workDir, cleanup, nil
}

// copyFileToWorkDir copies a single file from projectRoot to workDir.
func copyFileToWorkDir(projectRoot, workDir, filename string) error {
	src := filepath.Join(projectRoot, filename)
	dst := filepath.Join(workDir, filename)

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("create dest dir: %w", err)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create dest file: %w", err)
	}
	defer func() { _ = dstFile.Close() }()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copy file content: %w", err)
	}

	return nil
}

// executeJourneyInIsolation runs a journey's test command in its isolated work directory.
// The test command is executed with the work dir as cwd.
// Returns the execution result with output and exit code.
func executeJourneyInIsolation(cfg *JourneyExecutionConfig, workDir, journeyName string) JourneyResult {
	start := time.Now()
	result := JourneyResult{
		JourneyName: journeyName,
	}

	// Parse the test command into command name and args
	parts := strings.Fields(cfg.TestCommand)
	if len(parts) == 0 {
		result.Passed = false
		result.Error = "empty test command"
		result.Duration = time.Since(start)
		return result
	}

	cmdName := parts[0]
	cmdArgs := parts[1:]

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+workDir)

	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	result.Duration = time.Since(start)

	if err != nil {
		result.Passed = false
		result.ExitCode = 1
		var exitErr *exec.ExitError
		if ok := isErrorType(err, &exitErr); ok {
			result.ExitCode = exitErr.ExitCode()
		}
		result.Error = fmt.Sprintf("test command failed: %v", err)
	} else {
		result.Passed = true
		result.ExitCode = 0
	}

	return result
}

// isErrorType checks if an error is an exec.ExitError and extracts the exit code.
func isErrorType(err error, exitErr **exec.ExitError) bool {
	*exitErr, _ = err.(*exec.ExitError)
	return *exitErr != nil
}
