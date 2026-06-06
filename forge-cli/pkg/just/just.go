// Package just provides utilities for interacting with justfile recipes,
// including quality gate execution and scope resolution.
package just

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"forge-cli/pkg/forgelog"
)

// GateRecipe defines one step in the quality gate sequence.
type GateRecipe struct {
	Name     string // just recipe name
	Optional bool   // if true, skip when recipe not found
	Blocking bool   // if true, failure halts the sequence
}

// FullGateSequence returns the full quality gate: compile → fmt → lint → unit-test → test → probe.
// Used by all-completed hook for complete project validation.
func FullGateSequence() []GateRecipe {
	return []GateRecipe{
		{Name: "compile", Optional: false, Blocking: true},
		{Name: "fmt", Optional: true, Blocking: false},
		{Name: "lint", Optional: true, Blocking: true},
		{Name: "unit-test", Optional: false, Blocking: true},
		{Name: "test", Optional: false, Blocking: true},
		{Name: "probe", Optional: true, Blocking: false},
	}
}

// UnitGateSequence returns compile → fmt → lint → unit-test.
// Used by breaking tasks on submit for fast feedback.
func UnitGateSequence() []GateRecipe {
	return []GateRecipe{
		{Name: "compile", Optional: false, Blocking: true},
		{Name: "fmt", Optional: true, Blocking: false},
		{Name: "lint", Optional: true, Blocking: true},
		{Name: "unit-test", Optional: false, Blocking: true},
	}
}

// NonBreakingGateSequence returns compile → fmt → lint (without tests).
// Used by non-breaking tasks on submit.
func NonBreakingGateSequence() []GateRecipe {
	return []GateRecipe{
		{Name: "compile", Optional: false, Blocking: true},
		{Name: "fmt", Optional: true, Blocking: false},
		{Name: "lint", Optional: true, Blocking: true},
	}
}

// HasJustfile checks if a justfile exists in the given directory.
func HasJustfile(dir string) bool {
	return FileExists(filepath.Join(dir, "justfile")) ||
		FileExists(filepath.Join(dir, "Justfile"))
}

// HasRecipe checks if a recipe exists in the justfile using dry-run.
func HasRecipe(dir, recipe string) bool {
	c := exec.Command("just", "--dry-run", recipe)
	c.Dir = dir
	return c.Run() == nil
}

// hasRecipeWithArg checks if a recipe exists with an argument using dry-run.
func hasRecipeWithArg(dir, recipe, arg string) bool {
	c := exec.Command("just", "--dry-run", recipe, arg)
	c.Dir = dir
	return c.Run() == nil
}

// RunCapture runs a command, streams output to stderr, and returns
// the combined output along with whether the command succeeded.
func RunCapture(dir string, name string, args ...string) (string, bool) {
	c := exec.Command(name, args...)
	c.Dir = dir
	output, err := c.CombinedOutput()
	fmt.Fprint(os.Stderr, string(output))
	return string(output), err == nil
}

// ResolveScope applies scope resolution: only pass scope to just if the justfile
// has recipes that accept a scope argument. Probes by running `just --dry-run compile <scope>`.
func ResolveScope(projectRoot, scope string) string {
	if scope == "" || scope == "all" {
		return ""
	}
	if !HasJustfile(projectRoot) {
		return ""
	}
	// Probe whether the compile recipe accepts a scope argument.
	// If `just --dry-run compile <scope>` succeeds, the recipe takes scope.
	c := exec.Command("just", "--dry-run", "compile", scope)
	c.Dir = projectRoot
	if c.Run() == nil {
		return scope
	}
	return ""
}

// RunGate executes the gate sequence in order.
// scope: task scope (frontend/backend/empty). Only passed to just if project is mixed.
// onFail: called when a blocking step fails. Receives step name and output.
// Returns true if all steps passed (or skipped gracefully).
func RunGate(projectRoot, scope string, steps []GateRecipe, onFail func(step, output string)) bool {
	if !HasJustfile(projectRoot) {
		return true
	}

	resolvedScope := ResolveScope(projectRoot, scope)

	for _, step := range steps {
		recipeExists := HasRecipe(projectRoot, step.Name)
		// For scoped recipes, also probe with the resolved scope.
		if !recipeExists && resolvedScope != "" {
			recipeExists = hasRecipeWithArg(projectRoot, step.Name, resolvedScope)
		}
		if !recipeExists {
			if step.Optional {
				continue
			}
			// Required recipe missing — no fallback, report error with init-justfile hint.
			output := fmt.Sprintf("required recipe %q not found in justfile. Run `just init-justfile` to generate standard recipes.", step.Name)
			forgelog.Error("ERROR: %s\n", output)
			if step.Blocking && onFail != nil {
				onFail(step.Name, output)
			}
			return false
		}

		args := []string{step.Name}
		if resolvedScope != "" {
			args = append(args, resolvedScope)
		}

		output, success := RunCapture(projectRoot, "just", args...)
		if !success {
			if step.Blocking {
				if onFail != nil {
					onFail(step.Name, output)
				}
				return false
			}
			forgelog.Warn("WARNING: non-blocking gate step %q failed\n", step.Name)
		}
	}
	return true
}

// ExtractFailLines extracts all lines starting with "--- FAIL:" from output.
// Returns the joined lines separated by newlines.
// Returns empty string if no "--- FAIL:" lines are found.
func ExtractFailLines(output string) string {
	var failLines []string
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "--- FAIL:") {
			failLines = append(failLines, line)
		}
	}
	return strings.Join(failLines, "\n")
}

// ExtractConciseError returns the last N non-empty lines from output.
func ExtractConciseError(output string, maxLines int) string {
	lines := strings.Split(output, "\n")
	var nonEmpty []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			nonEmpty = append(nonEmpty, trimmed)
		}
	}
	if len(nonEmpty) <= maxLines {
		return output
	}
	return "...\n" + strings.Join(nonEmpty[len(nonEmpty)-maxLines:], "\n")
}

// FileExists checks if a file or directory exists at the given path.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
