// Package profile provides test profile resolution utilities.
package profile

import (
	"fmt"
	"slices"
)

// FrameworkInfo describes a test framework's code generation conventions.
// This is the data the skill uses to generate correct test code structure.
type FrameworkInfo struct {
	// Name is the canonical framework key (e.g. "go-testing", "pytest", "mocha").
	Name string `yaml:"name"`

	// TestFunctionPattern describes the function/test declaration pattern.
	// e.g. "func Test*", "def test_*", "describe/it"
	TestFunctionPattern string `yaml:"test-function-pattern"`

	// FilePattern describes the file naming convention.
	// e.g. "*_test.go", "test_*.py", "*.spec.ts"
	FilePattern string `yaml:"file-pattern"`

	// LanguageHint is the language this framework is associated with.
	// Empty for custom/user-defined frameworks.
	LanguageHint string `yaml:"language-hint,omitempty"`
}

// builtinFrameworks is the registry of built-in framework definitions.
// Projects can override these by declaring test-framework in config.
// The key is the framework name used in config.
var builtinFrameworks = map[string]FrameworkInfo{
	"go-testing": {
		Name:                "go-testing",
		TestFunctionPattern: "func Test*",
		FilePattern:         "*_test.go",
		LanguageHint:        "go",
	},
	"pytest": {
		Name:                "pytest",
		TestFunctionPattern: "def test_*",
		FilePattern:         "test_*.py",
		LanguageHint:        "python",
	},
	"mocha": {
		Name:                "mocha",
		TestFunctionPattern: "describe/it",
		FilePattern:         "*.spec.ts",
		LanguageHint:        "javascript",
	},
	"junit5": {
		Name:                "junit5",
		TestFunctionPattern: "@Test void *",
		FilePattern:         "*Test.java",
		LanguageHint:        "java",
	},
	"rust-test": {
		Name:                "rust-test",
		TestFunctionPattern: "#[test] fn *",
		FilePattern:         "*_test.rs",
		LanguageHint:        "rust",
	},
	"maestro": {
		Name:                "maestro",
		TestFunctionPattern: "YAML flow",
		FilePattern:         "*.yaml",
		LanguageHint:        "mobile",
	},
}

// KnownFrameworkNames returns the sorted list of built-in framework names.
func KnownFrameworkNames() []string {
	names := make([]string, 0, len(builtinFrameworks))
	for k := range builtinFrameworks {
		names = append(names, k)
	}
	slices.Sort(names)
	return names
}

// GetBuiltinFramework returns the FrameworkInfo for a built-in framework name.
// Returns an error if the framework name is not in the built-in registry.
func GetBuiltinFramework(name string) (FrameworkInfo, error) {
	fw, ok := builtinFrameworks[name]
	if !ok {
		return FrameworkInfo{}, fmt.Errorf("unknown test framework: %s (known: %v)", name, KnownFrameworkNames())
	}
	return fw, nil
}

// IsKnownFramework checks whether a framework name is a built-in framework.
func IsKnownFramework(name string) bool {
	_, ok := builtinFrameworks[name]
	return ok
}

// DefaultFrameworkForLanguage returns the default framework for a given language.
// Returns ("", false) if the language has no default framework mapping.
func DefaultFrameworkForLanguage(lang string) (string, bool) {
	defaults := map[string]string{
		"go":         "go-testing",
		"python":     "pytest",
		"javascript": "mocha",
		"java":       "junit5",
		"rust":       "rust-test",
		"mobile":     "maestro",
	}
	fw, ok := defaults[lang]
	return fw, ok
}

// ResolveTestFramework resolves the effective test framework for a project.
// Priority:
//  1. config.TestFramework (explicit override in .forge/config.yaml)
//  2. Default framework for the first resolved language
//  3. Empty string (no framework resolved)
//
// Returns the FrameworkInfo for the resolved framework.
// If no framework is resolved, returns a zero FrameworkInfo and nil error.
func ResolveTestFramework(projectRoot string) (FrameworkInfo, error) {
	cfg, err := ReadConfig(projectRoot)
	if err != nil {
		return FrameworkInfo{}, fmt.Errorf("resolve test framework: %w", err)
	}

	// Priority 1: explicit config override
	if cfg != nil && cfg.TestFramework != "" {
		fw, err := GetBuiltinFramework(cfg.TestFramework)
		if err != nil {
			// Not a built-in — return a custom FrameworkInfo with just the name.
			// Custom frameworks are fully user-defined; only the name comes from config.
			return FrameworkInfo{Name: cfg.TestFramework}, nil
		}
		return fw, nil
	}

	// Priority 2: default for detected language
	langs, err := ReadLanguages(projectRoot)
	if err != nil {
		return FrameworkInfo{}, fmt.Errorf("resolve test framework: %w", err)
	}
	if len(langs) == 0 {
		return FrameworkInfo{}, nil
	}

	if fwName, ok := DefaultFrameworkForLanguage(langs[0]); ok {
		fw, _ := GetBuiltinFramework(fwName)
		return fw, nil
	}

	return FrameworkInfo{}, nil
}

// RegisterCustomFramework adds a framework to the registry at runtime.
// This allows projects to define custom frameworks that can be resolved
// by name, just like built-in ones. Overwrites if name already exists.
func RegisterCustomFramework(fw FrameworkInfo) {
	builtinFrameworks[fw.Name] = fw
}
