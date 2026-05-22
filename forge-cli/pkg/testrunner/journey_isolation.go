package testrunner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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
	ProjectRoot string
}

// ResolveJourneyExecutionConfig creates the journey execution config.
func ResolveJourneyExecutionConfig(projectRoot string) *JourneyExecutionConfig {
	return &JourneyExecutionConfig{
		ProjectRoot: projectRoot,
	}
}

// CreateJourneyWorkDir creates an isolated temporary directory for a journey.
// The directory name includes the journey name and a random suffix.
// Returns the work dir path, a cleanup function, and any error.
func CreateJourneyWorkDir(_, journeyName string) (string, func(), error) {
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

// CopyFileToWorkDir copies a single file from projectRoot to workDir.
func CopyFileToWorkDir(projectRoot, workDir, filename string) error {
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

// ExecuteJourneyInIsolation runs a journey's e2e tests using `just e2e-test`
// from the project root with the journey filter.
// Returns the execution result with output and exit code.
func ExecuteJourneyInIsolation(cfg *JourneyExecutionConfig, _, journeyName string) JourneyResult {
	start := time.Now()
	result := JourneyResult{
		JourneyName: journeyName,
	}

	cmd := exec.Command("just", "e2e-test", journeyName)
	cmd.Dir = cfg.ProjectRoot
	cmd.Env = append(os.Environ(), "FORGE_JOURNEY="+journeyName)

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
		result.Error = fmt.Sprintf("just e2e-test %s failed: %v", journeyName, err)
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
