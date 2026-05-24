// Package testrunner provides test execution and output formatting.
package testrunner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/pkg/just"
)

// PrintHookJSON writes a Claude Code hook decision as JSON to stdout.
func PrintHookJSON(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to marshal hook JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

// Capitalize returns s with its first letter uppercased.
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// hasNpmTestScript checks if package.json has a "test" script.
func hasNpmTestScript(projectRoot string) bool {
	data, err := os.ReadFile(filepath.Join(projectRoot, "package.json"))
	if err != nil {
		return false
	}
	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false
	}
	_, ok := pkg.Scripts["test"]
	return ok
}

// RunProjectTests detects and runs the project's test command.
// Probe chain: just unit-test -> just test -> go test -> npm -> pytest -> warning.
// Retains fallback mechanism: first matching recipe wins.
func RunProjectTests(projectRoot string) (string, bool) {
	switch {
	case just.HasJustfile(projectRoot) && just.HasRecipe(projectRoot, "unit-test"):
		return just.RunCapture(projectRoot, "just", "unit-test")
	case just.HasJustfile(projectRoot) && just.HasRecipe(projectRoot, "test"):
		return just.RunCapture(projectRoot, "just", "test")
	case just.FileExists(filepath.Join(projectRoot, "go.mod")):
		return just.RunCapture(projectRoot, "go", "test", "./...")
	case just.FileExists(filepath.Join(projectRoot, "package.json")) && hasNpmTestScript(projectRoot):
		return just.RunCapture(projectRoot, "npm", "test")
	case just.FileExists(filepath.Join(projectRoot, "pytest.ini")) || just.FileExists(filepath.Join(projectRoot, "pyproject.toml")):
		return just.RunCapture(projectRoot, "pytest")
	default:
		fmt.Println("WARNING: No test command found. Ensure justfile has a 'unit-test' recipe, or add go.mod/package.json/pytest.ini.")
		return "", true
	}
}
