package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupTestingProject creates a temp directory with go.mod for language detection.
func setupTestingProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	return dir
}

// setupEmptyProject creates a temp directory with no language signals.
func setupEmptyProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Write a CLAUDE.md so FindProjectRoot succeeds but no language files
	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	return dir
}

// resetTestingGetLanguage resets the --language flag between tests.
func resetTestingGetLanguage() {
	testingGetLanguage = ""
}

func TestTestingDetect_OutputsDetectedLanguage(t *testing.T) {
	resetTestingGetLanguage()
	_ = setupTestingProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "detect"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing detect failed: %v", err)
	}

	if !strings.Contains(output, "go") {
		t.Errorf("expected 'go' in output, got: %q", output)
	}
}

func TestTestingDetect_NoLanguage(t *testing.T) {
	resetTestingGetLanguage()
	_ = setupEmptyProject(t)

	output, _ := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "detect"})
		return rootCmd.Execute()
	})

	if !strings.Contains(output, "(none)") {
		t.Errorf("expected '(none)' in output when no language detected, got: %q", output)
	}
}

func TestTestingDetect_OutputFormat(t *testing.T) {
	resetTestingGetLanguage()
	_ = setupTestingProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "detect"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing detect failed: %v", err)
	}

	// Output should use structured block format with separators
	if !strings.Contains(output, "---") {
		t.Errorf("expected block separator '---' in output, got: %q", output)
	}
	if !strings.Contains(output, "LANGUAGE:") {
		t.Errorf("expected 'LANGUAGE:' in output, got: %q", output)
	}
}

func TestTestingGetGenerate_AutoDetect(t *testing.T) {
	resetTestingGetLanguage()
	_ = setupTestingProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "get", "generate"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing get generate failed: %v", err)
	}

	// Should output the Go generate.md strategy content
	if len(output) == 0 {
		t.Error("expected non-empty output for get generate")
	}
}

func TestTestingGetGenerate_WithLanguageFlag(t *testing.T) {
	resetTestingGetLanguage()
	_ = setupTestingProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "get", "generate", "--language", "go"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing get generate --language go failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for get generate --language go")
	}
}

func TestTestingGetRun(t *testing.T) {
	resetTestingGetLanguage()
	_ = setupTestingProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "get", "run"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing get run failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for get run")
	}
}

func TestTestingGetGraduate(t *testing.T) {
	resetTestingGetLanguage()
	_ = setupTestingProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "get", "graduate"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing get graduate failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for get graduate")
	}
}

func TestTestingGetJustfile(t *testing.T) {
	resetTestingGetLanguage()
	_ = setupTestingProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "get", "justfile"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing get justfile failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for get justfile")
	}
}

func TestTestingGetTemplate(t *testing.T) {
	resetTestingGetLanguage()
	_ = setupTestingProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "get", "template", "test-file.go"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing get template test-file.go failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for get template")
	}
}

func TestTestingInterfaces_AutoDetect(t *testing.T) {
	resetTestingGetLanguage()
	_ = setupTestingProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "interfaces"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing interfaces failed: %v", err)
	}

	if !strings.Contains(output, "api") || !strings.Contains(output, "cli") {
		t.Errorf("expected 'api' and 'cli' in interfaces output, got: %q", output)
	}
}

func TestTestingInterfaces_WithConfigOverride(t *testing.T) {
	resetTestingGetLanguage()
	dir := setupTestingProject(t)

	// Write config.yaml with explicit interfaces
	configDir := filepath.Join(dir, ".forge")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	configContent := "project-type: backend\ninterfaces:\n  - api\n"
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "interfaces"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing interfaces failed: %v", err)
	}

	if !strings.Contains(output, "api") {
		t.Errorf("expected 'api' in interfaces output, got: %q", output)
	}
	// Config overrides: only 'api' should appear, not 'cli'
	if strings.Contains(output, "cli") {
		t.Errorf("expected 'cli' NOT in interfaces output when config overrides, got: %q", output)
	}
}

func TestTestingResolveLanguage_NoLanguageDetected_NoConfig(t *testing.T) {
	resetTestingGetLanguage()
	dir := setupEmptyProject(t)

	// Test the resolveLanguageFromFlags helper directly
	_, err := resolveLanguageFromFlags(dir)
	if err == nil {
		t.Error("expected error when no language detected and no config")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "No language detected") {
		t.Errorf("expected error message to contain 'No language detected', got: %s", errMsg)
	}
}

func TestTestingGet_SpecificLanguage(t *testing.T) {
	resetTestingGetLanguage()
	_ = setupTestingProject(t)

	// Request javascript strategy specifically via --language flag
	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "get", "generate", "--language", "javascript"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing get generate --language javascript failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for javascript generate strategy")
	}
}

func TestProfileCommand_Removed(t *testing.T) {
	// The 'profile' command should not exist on rootCmd
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "profile" {
			t.Error("forge profile command should not exist -- it should be replaced by forge testing")
		}
	}
}

func TestTestingCommand_Registered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "testing" {
			found = true
			break
		}
	}
	if !found {
		t.Error("forge testing command should be registered on rootCmd")
	}
}

func TestTestingCommand_Subcommands(t *testing.T) {
	subNames := make(map[string]bool)
	for _, cmd := range testingCmd.Commands() {
		subNames[cmd.Name()] = true
	}

	expected := []string{"detect", "get", "interfaces", "framework", "run-journey"}
	for _, name := range expected {
		if !subNames[name] {
			t.Errorf("testing group missing subcommand: %s (have: %v)", name, subNames)
		}
	}
}

func TestTestingGetGetSubcommands(t *testing.T) {
	subNames := make(map[string]bool)
	for _, cmd := range testingGetCmd.Commands() {
		subNames[cmd.Name()] = true
	}

	expected := []string{"generate", "run", "graduate", "justfile", "template"}
	for _, name := range expected {
		if !subNames[name] {
			t.Errorf("testing get missing subcommand: %s (have: %v)", name, subNames)
		}
	}
}

func TestTestingGet_JavaLanguage(t *testing.T) {
	resetTestingGetLanguage()
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Create pom.xml for Java detection
	if err := os.WriteFile(filepath.Join(dir, "pom.xml"), []byte("<project></project>"), 0644); err != nil {
		t.Fatal(err)
	}

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "get", "generate"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing get generate (java auto-detect) failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for java generate strategy")
	}
}

// setupMultiLanguageProject creates a project with both go.mod and package.json (playwright).
func setupMultiLanguageProject(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Go signal
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// JavaScript/Playwright signal
	pkgJSON := `{"devDependencies": {"@playwright/test": "^1.0.0"}}`
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestTestingDetect_MultiLanguage(t *testing.T) {
	resetTestingGetLanguage()
	setupMultiLanguageProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "detect"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing detect failed: %v", err)
	}

	if !strings.Contains(output, "go") {
		t.Errorf("expected 'go' in multi-language detect output, got: %q", output)
	}
	if !strings.Contains(output, "javascript") {
		t.Errorf("expected 'javascript' in multi-language detect output, got: %q", output)
	}
}

func TestTestingGet_MultiLanguage_SelectSpecific(t *testing.T) {
	resetTestingGetLanguage()
	setupMultiLanguageProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "get", "generate", "--language", "javascript"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing get generate --language javascript failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for javascript strategy in multi-language project")
	}
}

func TestTestingGet_MultiLanguage_DefaultFirst(t *testing.T) {
	resetTestingGetLanguage()
	setupMultiLanguageProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "get", "generate"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing get generate (multi-language default) failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for default language in multi-language project")
	}
}

func TestTestingFramework_AutoDetectGo(t *testing.T) {
	resetTestingGetLanguage()
	_ = setupTestingProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "framework"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing framework failed: %v", err)
	}

	if !strings.Contains(output, "FRAMEWORK: go-testing") {
		t.Errorf("expected 'FRAMEWORK: go-testing' in output, got: %q", output)
	}
	if !strings.Contains(output, "PATTERN: func Test*") {
		t.Errorf("expected 'PATTERN: func Test*' in output, got: %q", output)
	}
	if !strings.Contains(output, "FILES: *_test.go") {
		t.Errorf("expected 'FILES: *_test.go' in output, got: %q", output)
	}
	if !strings.Contains(output, "SOURCE: language-default") {
		t.Errorf("expected 'SOURCE: language-default' in output, got: %q", output)
	}
}

func TestTestingFramework_ConfigOverride(t *testing.T) {
	resetTestingGetLanguage()
	dir := setupTestingProject(t)

	// Write config.yaml with test-framework override
	configDir := filepath.Join(dir, ".forge")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	configContent := "languages:\n  - go\ntest-framework: pytest\n"
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "framework"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing framework failed: %v", err)
	}

	if !strings.Contains(output, "FRAMEWORK: pytest") {
		t.Errorf("expected 'FRAMEWORK: pytest' in output, got: %q", output)
	}
	if !strings.Contains(output, "PATTERN: def test_*") {
		t.Errorf("expected 'PATTERN: def test_*' in output, got: %q", output)
	}
	if !strings.Contains(output, "SOURCE: config") {
		t.Errorf("expected 'SOURCE: config' in output, got: %q", output)
	}
}

func TestTestingFramework_NoLanguage(t *testing.T) {
	resetTestingGetLanguage()
	_ = setupEmptyProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"testing", "framework"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("testing framework failed: %v", err)
	}

	if !strings.Contains(output, "(none)") {
		t.Errorf("expected '(none)' in output when no language, got: %q", output)
	}
}
