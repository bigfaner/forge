package e2e

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// firstStderrLine extracts the first line from stderr output.
func firstStderrLine(output []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(output))
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

// formatToolError creates a formatted error for external tool failures.
// Format: "<command> failed: <first line of child stderr>"
func formatToolError(cmd string, output []byte) error {
	firstLine := firstStderrLine(output)
	if firstLine == "" {
		return fmt.Errorf("%s failed", cmd)
	}
	return fmt.Errorf("%s failed: %s", cmd, firstLine)
}

const justNotFoundMsg = "error: 'just' is required but not found on PATH. Install: https://github.com/casey/just"

// isJustNotFound checks if the error indicates just is not on PATH.
func isJustNotFound(err error) bool {
	return strings.Contains(err.Error(), "just") &&
		(strings.Contains(err.Error(), "executable file not found") ||
			strings.Contains(err.Error(), "file not found"))
}

// runJust executes a just recipe and returns its output and error.
// It wraps subprocess errors with a not-found check for just itself.
func runJust(recipe string, args ...string) ([]byte, error) {
	cmdArgs := append([]string{recipe}, args...)
	out, err := runner.Run("just", cmdArgs...)
	if err != nil {
		if isJustNotFound(err) {
			return nil, fmt.Errorf("%s", justNotFoundMsg)
		}
		return nil, formatToolError("just "+recipe, out)
	}
	return out, nil
}

// Run executes e2e tests using the configured profile.
func Run(opts RunOpts) error {
	if _, err := ResolveProfile(opts.ProjectRoot); err != nil {
		return err
	}

	args := []string{}
	if opts.Feature != "" {
		args = append(args, "feature="+opts.Feature)
	}

	out, err := runJust("test-e2e", args...)
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

// Setup installs e2e dependencies for the configured profile.
func Setup(opts RunOpts) error {
	if _, err := ResolveProfile(opts.ProjectRoot); err != nil {
		return err
	}

	_, err := runJust("e2e-setup")
	return err
}

// Verify scans e2e test files for unresolved VERIFY markers.
func Verify(opts RunOpts) error {
	// Validate profile is configured (all profiles use the same scan logic)
	if _, err := ResolveProfile(opts.ProjectRoot); err != nil {
		return err
	}

	// Determine scan directory
	scanDir := filepath.Join(opts.ProjectRoot, "tests", "e2e")
	if opts.Feature != "" {
		featureDir := filepath.Join(scanDir, "features", opts.Feature)
		if _, statErr := os.Stat(featureDir); os.IsNotExist(statErr) {
			return fmt.Errorf("%w: %s", ErrFeatureNotFound, opts.Feature)
		}
		scanDir = featureDir
	}

	// Walk files and look for VERIFY markers
	var filesWithMarkers []string
	walkErr := filepath.WalkDir(scanDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		// Only scan source files
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".go" && ext != ".ts" && ext != ".py" && ext != ".js" {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		if strings.Contains(string(data), "VERIFY") {
			rel, _ := filepath.Rel(opts.ProjectRoot, path)
			filesWithMarkers = append(filesWithMarkers, rel)
		}
		return nil
	})
	if walkErr != nil {
		return fmt.Errorf("scan directory: %w", walkErr)
	}

	if len(filesWithMarkers) > 0 {
		return fmt.Errorf("VERIFY markers found in: %s", strings.Join(filesWithMarkers, ", "))
	}

	return nil
}

// Compile runs a compile-only check on e2e test files.
func Compile(projectRoot string) error {
	if _, err := ResolveProfile(projectRoot); err != nil {
		return err
	}

	_, err := runJust("e2e-compile")
	return err
}

// Discover lists all e2e test cases without running them.
func Discover(projectRoot string) error {
	if _, err := ResolveProfile(projectRoot); err != nil {
		return err
	}

	out, err := runJust("e2e-discover")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
