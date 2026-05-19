// Package journey provides Journey-Driven test code generation utilities.
// It generates test files with Convention-based assertions and Journey smoke tests.
package journey

import (
	"fmt"
	"path/filepath"
	"strings"
)

// TestGenerationOpts holds options for generating Journey test code.
type TestGenerationOpts struct {
	// Journey is the Journey name (kebab-case, e.g. "task-lifecycle").
	Journey string
	// ContractsDir is the directory containing Contract specifications.
	ContractsDir string
	// Facts are the Fact Table entries as a plain string (Convention content).
	Facts string
	// CustomTemplatePath is the optional path to a custom template directory.
	// Empty means use built-in default templates.
	CustomTemplatePath string
}

// GeneratedTest represents a single generated test file.
type GeneratedTest struct {
	// Filename is the output file name (e.g. "claim_submit_test.go").
	Filename string
	// Content is the generated test code.
	Content string
	// IsSmokeTest is true if this is the Journey smoke test.
	IsSmokeTest bool
}

// TestOutputDir returns the output directory for a Journey's tests.
// Tests go directly into tests/<journey>/ (no staging area).
func TestOutputDir(projectRoot, journey string) string {
	return filepath.Join(projectRoot, "tests", journey)
}

// SmokeTestName generates the smoke test function name for a Journey.
func SmokeTestName(journey string) string {
	parts := strings.Split(journey, "-")
	var name strings.Builder
	for _, p := range parts {
		if len(p) > 0 {
			name.WriteString(strings.ToUpper(p[:1]))
			name.WriteString(p[1:])
		}
	}
	return fmt.Sprintf("TestJourney%sSmoke", name.String())
}

// sanitizeName converts a human-readable name to a safe identifier component.
func sanitizeName(name string) string {
	replacer := strings.NewReplacer(
		" ", "_",
		"-", "_",
		"/", "_",
		".", "_",
		":", "_",
	)
	s := replacer.Replace(name)
	s = strings.ToLower(s)
	// Remove consecutive underscores
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}
	s = strings.Trim(s, "_")
	return s
}
