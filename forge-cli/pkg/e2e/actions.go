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

// Run executes e2e tests using the configured profile.
func Run(opts RunOpts) error {
	profile, err := ResolveProfile(opts.ProjectRoot)
	if err != nil {
		return err
	}

	switch profile {
	case "go-test":
		testPath := "./tests/e2e/..."
		if opts.Feature != "" {
			testPath = fmt.Sprintf("./tests/e2e/features/%s/...", opts.Feature)
		}
		cmdName := "go test"
		args := []string{"test", testPath}
		out, runErr := runner.Run("go", args...)
		if runErr != nil {
			return formatToolError(cmdName, out)
		}
		fmt.Println(string(out))
		return nil

	case "web-playwright":
		cmdName := "npx playwright test"
		args := []string{"playwright", "test"}
		out, runErr := runner.Run("npx", args...)
		if runErr != nil {
			return formatToolError(cmdName, out)
		}
		fmt.Println(string(out))
		return nil

	default:
		return fmt.Errorf("unsupported profile for run: %s", profile)
	}
}

// Setup installs e2e dependencies for the configured profile.
func Setup(opts RunOpts) error {
	profile, err := ResolveProfile(opts.ProjectRoot)
	if err != nil {
		return err
	}

	switch profile {
	case "go-test":
		cmdName := "go install"
		_, runErr := runner.Run("go", "install")
		if runErr != nil {
			return formatToolError(cmdName, nil)
		}
		return nil

	case "web-playwright":
		cmdName := "npx playwright install"
		out, runErr := runner.Run("npx", "playwright", "install")
		if runErr != nil {
			return formatToolError(cmdName, out)
		}
		return nil

	case "pytest":
		cmdName := "python -m pip install pytest"
		out, runErr := runner.Run("python", "-m", "pip", "install", "pytest")
		if runErr != nil {
			return formatToolError(cmdName, out)
		}
		return nil

	default:
		return fmt.Errorf("unsupported profile for setup: %s", profile)
	}
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
	profile, err := ResolveProfile(projectRoot)
	if err != nil {
		return err
	}

	switch profile {
	case "go-test":
		cmdName := "go build"
		out, runErr := runner.Run("go", "build", "./tests/e2e/...")
		if runErr != nil {
			return formatToolError(cmdName, out)
		}
		return nil

	case "web-playwright":
		cmdName := "npx tsc --noEmit"
		out, runErr := runner.Run("npx", "tsc", "--noEmit")
		if runErr != nil {
			return formatToolError(cmdName, out)
		}
		return nil

	case "pytest":
		cmdName := "python -m compileall tests/e2e/ -q"
		out, runErr := runner.Run("python", "-m", "compileall", "tests/e2e/", "-q")
		if runErr != nil {
			return formatToolError(cmdName, out)
		}
		return nil

	default:
		return fmt.Errorf("unsupported profile for compile: %s", profile)
	}
}

// Discover lists all e2e test cases without running them.
func Discover(projectRoot string) error {
	profile, err := ResolveProfile(projectRoot)
	if err != nil {
		return err
	}

	switch profile {
	case "go-test":
		cmdName := "go test -list"
		out, runErr := runner.Run("go", "test", "./tests/e2e/...", "-list", ".*", "-tags=e2e")
		if runErr != nil {
			return formatToolError(cmdName, out)
		}
		fmt.Println(string(out))
		return nil

	case "web-playwright":
		cmdName := "npx playwright test --list"
		out, runErr := runner.Run("npx", "playwright", "test", "--list")
		if runErr != nil {
			return formatToolError(cmdName, out)
		}
		fmt.Println(string(out))
		return nil

	case "pytest":
		cmdName := "python -m pytest --collect-only"
		out, runErr := runner.Run("python", "-m", "pytest", "tests/e2e/", "--collect-only", "-q")
		if runErr != nil {
			return formatToolError(cmdName, out)
		}
		fmt.Println(string(out))
		return nil

	default:
		return fmt.Errorf("unsupported profile for discover: %s", profile)
	}
}
