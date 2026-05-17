package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupTestProject creates a temp directory with go.mod for language detection.
func setupTestProject(t *testing.T) string {
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

// resetTestGetLanguage resets the --language flag between tests.
func resetTestGetLanguage() {
	testGetLanguage = ""
}

func TestTestDetect_OutputsDetectedLanguage(t *testing.T) {
	resetTestGetLanguage()
	_ = setupTestProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "detect"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test detect failed: %v", err)
	}

	if !strings.Contains(output, "go") {
		t.Errorf("expected 'go' in output, got: %q", output)
	}
}

func TestTestDetect_NoLanguage(t *testing.T) {
	resetTestGetLanguage()
	_ = setupEmptyProject(t)

	output, _ := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "detect"})
		return rootCmd.Execute()
	})

	if !strings.Contains(output, "(none)") {
		t.Errorf("expected '(none)' in output when no language detected, got: %q", output)
	}
}

func TestTestDetect_OutputFormat(t *testing.T) {
	resetTestGetLanguage()
	_ = setupTestProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "detect"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test detect failed: %v", err)
	}

	// Output should use structured block format with separators
	if !strings.Contains(output, "---") {
		t.Errorf("expected block separator '---' in output, got: %q", output)
	}
	if !strings.Contains(output, "LANGUAGE:") {
		t.Errorf("expected 'LANGUAGE:' in output, got: %q", output)
	}
}

func TestTestGetGenerate_AutoDetect(t *testing.T) {
	resetTestGetLanguage()
	_ = setupTestProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "get", "generate"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test get generate failed: %v", err)
	}

	// Should output the Go generate.md strategy content
	if len(output) == 0 {
		t.Error("expected non-empty output for get generate")
	}
}

func TestTestGetGenerate_WithLanguageFlag(t *testing.T) {
	resetTestGetLanguage()
	_ = setupTestProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "get", "generate", "--language", "go"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test get generate --language go failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for get generate --language go")
	}
}

func TestTestGetRun(t *testing.T) {
	resetTestGetLanguage()
	_ = setupTestProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "get", "run"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test get run failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for get run")
	}
}

func TestTestGetJustfile(t *testing.T) {
	resetTestGetLanguage()
	_ = setupTestProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "get", "justfile"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test get justfile failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for get justfile")
	}
}

func TestTestGetTemplate(t *testing.T) {
	resetTestGetLanguage()
	_ = setupTestProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "get", "template", "test-file.go"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test get template test-file.go failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for get template")
	}
}

func TestTestInterfaces_AutoDetect(t *testing.T) {
	resetTestGetLanguage()
	_ = setupTestProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "interfaces"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test interfaces failed: %v", err)
	}

	if !strings.Contains(output, "api") || !strings.Contains(output, "cli") {
		t.Errorf("expected 'api' and 'cli' in interfaces output, got: %q", output)
	}
}

func TestTestInterfaces_WithConfigOverride(t *testing.T) {
	resetTestGetLanguage()
	dir := setupTestProject(t)

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
		rootCmd.SetArgs([]string{"test", "interfaces"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test interfaces failed: %v", err)
	}

	if !strings.Contains(output, "api") {
		t.Errorf("expected 'api' in interfaces output, got: %q", output)
	}
	// Config overrides: only 'api' should appear, not 'cli'
	if strings.Contains(output, "cli") {
		t.Errorf("expected 'cli' NOT in interfaces output when config overrides, got: %q", output)
	}
}

func TestResolveLanguage_NoLanguageDetected_NoConfig(t *testing.T) {
	resetTestGetLanguage()
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

func TestTestGet_SpecificLanguage(t *testing.T) {
	resetTestGetLanguage()
	_ = setupTestProject(t)

	// Request javascript strategy specifically via --language flag
	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "get", "generate", "--language", "javascript"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test get generate --language javascript failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for javascript generate strategy")
	}
}

func TestProfileCommand_Removed(t *testing.T) {
	// The 'profile' command should not exist on rootCmd
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "profile" {
			t.Error("forge profile command should not exist -- it should be replaced by forge test")
		}
	}
}

func TestTestingCommand_Removed(t *testing.T) {
	// The old 'testing' command should not exist on rootCmd
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "testing" {
			t.Error("forge testing command should not exist -- it is renamed to forge test")
		}
	}
}

func TestTestCommand_Registered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "test" {
			found = true
			break
		}
	}
	if !found {
		t.Error("forge test command should be registered on rootCmd")
	}
}

func TestTestCommand_Subcommands(t *testing.T) {
	subNames := make(map[string]bool)
	for _, cmd := range testCmd.Commands() {
		subNames[cmd.Name()] = true
	}

	expected := []string{"detect", "get", "interfaces", "framework", "run-journey", "verify", "promote"}
	for _, name := range expected {
		if !subNames[name] {
			t.Errorf("test group missing subcommand: %s (have: %v)", name, subNames)
		}
	}
}

func TestTestGetSubcommands(t *testing.T) {
	subNames := make(map[string]bool)
	for _, cmd := range testGetCmd.Commands() {
		subNames[cmd.Name()] = true
	}

	// graduate subcommand removed, replaced by promote at top level
	expected := []string{"generate", "run", "justfile", "template"}
	for _, name := range expected {
		if !subNames[name] {
			t.Errorf("test get missing subcommand: %s (have: %v)", name, subNames)
		}
	}

	// graduate should NOT exist
	if subNames["graduate"] {
		t.Error("test get should NOT have 'graduate' subcommand -- replaced by 'test promote'")
	}
}

func TestTestGet_JavaLanguage(t *testing.T) {
	resetTestGetLanguage()
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Create pom.xml for Java detection
	if err := os.WriteFile(filepath.Join(dir, "pom.xml"), []byte("<project></project>"), 0644); err != nil {
		t.Fatal(err)
	}

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "get", "generate"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test get generate (java auto-detect) failed: %v", err)
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

func TestTestDetect_MultiLanguage(t *testing.T) {
	resetTestGetLanguage()
	setupMultiLanguageProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "detect"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test detect failed: %v", err)
	}

	if !strings.Contains(output, "go") {
		t.Errorf("expected 'go' in multi-language detect output, got: %q", output)
	}
	if !strings.Contains(output, "javascript") {
		t.Errorf("expected 'javascript' in multi-language detect output, got: %q", output)
	}
}

func TestTestGet_MultiLanguage_SelectSpecific(t *testing.T) {
	resetTestGetLanguage()
	setupMultiLanguageProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "get", "generate", "--language", "javascript"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test get generate --language javascript failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for javascript strategy in multi-language project")
	}
}

func TestTestGet_MultiLanguage_DefaultFirst(t *testing.T) {
	resetTestGetLanguage()
	setupMultiLanguageProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "get", "generate"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test get generate (multi-language default) failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output for default language in multi-language project")
	}
}

func TestTestFramework_AutoDetectGo(t *testing.T) {
	resetTestGetLanguage()
	_ = setupTestProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "framework"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test framework failed: %v", err)
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

func TestTestFramework_ConfigOverride(t *testing.T) {
	resetTestGetLanguage()
	dir := setupTestProject(t)

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
		rootCmd.SetArgs([]string{"test", "framework"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test framework failed: %v", err)
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

func TestTestFramework_NoLanguage(t *testing.T) {
	resetTestGetLanguage()
	_ = setupEmptyProject(t)

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "framework"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test framework failed: %v", err)
	}

	if !strings.Contains(output, "(none)") {
		t.Errorf("expected '(none)' in output when no language, got: %q", output)
	}
}
