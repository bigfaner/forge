// Package just provides utilities for interacting with justfile recipes,
// including quality gate execution and scope resolution.
package just

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"forge-cli/pkg/profile"
)

// GateRecipe defines one step in the quality gate sequence.
type GateRecipe struct {
	Name     string // just recipe name
	Optional bool   // if true, skip when recipe not found
	Blocking bool   // if true, failure halts the sequence
}

// DefaultGateSequence returns the standard quality gate: compile → fmt → lint → test.
func DefaultGateSequence() []GateRecipe {
	return []GateRecipe{
		{Name: "compile", Optional: false, Blocking: true},
		{Name: "fmt", Optional: true, Blocking: false},
		{Name: "lint", Optional: true, Blocking: true},
		{Name: "test", Optional: false, Blocking: true},
	}
}

// LintGateSequence returns compile → fmt → lint (without test).
// Used by all-completed hook where test runs independently.
func LintGateSequence() []GateRecipe {
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

// RunCapture runs a command, streams output to stderr, and returns
// the combined output along with whether the command succeeded.
func RunCapture(dir string, name string, args ...string) (string, bool) {
	c := exec.Command(name, args...)
	c.Dir = dir
	output, err := c.CombinedOutput()
	fmt.Fprint(os.Stderr, string(output))
	return string(output), err == nil
}

// ResolveScope applies scope resolution: only pass scope to just if project is mixed.
// Reads project-type from .forge/config.yaml directly — no subprocess call.
func ResolveScope(projectRoot, scope string) string {
	if scope == "" || scope == "all" {
		return ""
	}
	cfg, err := profile.ReadConfig(projectRoot)
	if err != nil || cfg == nil {
		return ""
	}
	projectType := strings.TrimSpace(cfg.ProjectType)
	switch projectType {
	case "mixed":
		return scope
	case "frontend", "backend":
		return ""
	default:
		if projectType != "" {
			fmt.Fprintf(os.Stderr, "WARNING: unexpected project-type %q, expected frontend/backend/mixed; skipping scope\n", projectType)
		}
		return ""
	}
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
		if !HasRecipe(projectRoot, step.Name) {
			if step.Optional {
				continue
			}
			fmt.Fprintf(os.Stderr, "WARNING: required recipe %q not found in justfile; skipping quality gate step\n", step.Name)
			continue
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
			fmt.Fprintf(os.Stderr, "WARNING: non-blocking gate step %q failed\n", step.Name)
		}
	}
	return true
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
