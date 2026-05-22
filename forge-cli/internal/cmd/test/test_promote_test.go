package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/internal/cmd/base"

	"github.com/stretchr/testify/assert"
)

// --- Test: promote command registered ---

func TestTestPromote_CommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range Cmd.Commands() {
		if cmd.Name() == "promote" {
			found = true
			break
		}
	}
	if !found {
		t.Error("test group missing 'promote' subcommand")
	}
}

// --- Test: isTestFile helper ---

func TestIsTestFile_Go(t *testing.T) {
	if !isTestFile("foo_test.go") {
		t.Error("foo_test.go should be a test file")
	}
	if isTestFile("foo.go") {
		t.Error("foo.go should NOT be a test file")
	}
}

func TestIsTestFile_Python(t *testing.T) {
	if !isTestFile("test_something.py") {
		t.Error("test_something.py should be a test file")
	}
	if !isTestFile("something_test.py") {
		t.Error("something_test.py should be a test file")
	}
	if isTestFile("something.py") {
		t.Error("something.py should NOT be a test file")
	}
}

func TestIsTestFile_JavaScript(t *testing.T) {
	if !isTestFile("app.test.ts") {
		t.Error("app.test.ts should be a test file")
	}
	if !isTestFile("app.spec.js") {
		t.Error("app.spec.js should be a test file")
	}
	if isTestFile("app.ts") {
		t.Error("app.ts should NOT be a test file")
	}
}

func TestIsTestFile_Java(t *testing.T) {
	if !isTestFile("FooTest.java") {
		t.Error("FooTest.java should be a test file")
	}
	if !isTestFile("FooTests.java") {
		t.Error("FooTests.java should be a test file")
	}
	if isTestFile("Foo.java") {
		t.Error("Foo.java should NOT be a test file")
	}
}

func TestIsTestFile_Rust(t *testing.T) {
	if !isTestFile("foo_test.rs") {
		t.Error("foo_test.rs should be a test file")
	}
	if isTestFile("foo.rs") {
		t.Error("foo.rs should NOT be a test file")
	}
}

func TestIsTestFile_OtherExtension(t *testing.T) {
	if isTestFile("readme.md") {
		t.Error("readme.md should NOT be a test file")
	}
	if isTestFile("data.json") {
		t.Error("data.json should NOT be a test file")
	}
}

// --- Test: replaceFeatureTag ---

func TestReplaceFeatureTag_Go(t *testing.T) {
	content := `//go:build feature

package main

// @feature test
func TestSomething(t *testing.T) {}`

	result := replaceFeatureTag(content, ".go")

	if !strings.Contains(result, "//go:build regression") {
		t.Error("expected //go:build regression in Go file")
	}
	if strings.Contains(result, "//go:build feature") {
		t.Error("//go:build feature should be replaced")
	}
	if !strings.Contains(result, "@regression") {
		t.Error("expected @regression tag")
	}
}

func TestReplaceFeatureTag_Python(t *testing.T) {
	content := `import pytest

@pytest.mark.feature
def test_something():
    pass`

	result := replaceFeatureTag(content, ".py")

	if !strings.Contains(result, "@pytest.mark.regression") {
		t.Error("expected @pytest.mark.regression in Python file")
	}
	if strings.Contains(result, "@pytest.mark.feature") {
		t.Error("@pytest.mark.feature should be replaced")
	}
}

func TestReplaceFeatureTag_JavaScript(t *testing.T) {
	content := `// @feature journey
test('something', async ({ page }) => {});`

	result := replaceFeatureTag(content, ".ts")

	if !strings.Contains(result, "@regression journey") {
		t.Error("expected @regression in JS/TS file")
	}
	if strings.Contains(result, "@feature") {
		t.Error("@feature should be replaced")
	}
}

func TestReplaceFeatureTag_Java(t *testing.T) {
	content := `@Tag("feature")
public class FooTest {}`

	result := replaceFeatureTag(content, ".java")

	if !strings.Contains(result, `@Tag("regression")`) {
		t.Error("expected @Tag(\"regression\") in Java file")
	}
	if strings.Contains(result, `@Tag("feature")`) {
		t.Error("@Tag(\"feature\") should be replaced")
	}
}

func TestReplaceFeatureTag_Rust(t *testing.T) {
	content := `#[cfg(feature = "feature")]
mod tests {}`

	result := replaceFeatureTag(content, ".rs")

	if !strings.Contains(result, `#[cfg(feature = "regression")]`) {
		t.Error("expected #[cfg(feature = \"regression\")] in Rust file")
	}
}

func TestReplaceFeatureTag_NoFeature(t *testing.T) {
	content := `package main

func TestSomething(t *testing.T) {}`

	result := replaceFeatureTag(content, ".go")

	if result != content {
		t.Error("content without @feature should not change")
	}
}

// --- Test: promoteJourneyTags ---

func TestPromoteJourneyTags_ReplacesFeatureTag(t *testing.T) {
	dir := t.TempDir()
	journeyDir := filepath.Join(dir, "my-journey")
	if err := os.MkdirAll(journeyDir, 0755); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(journeyDir, "main_test.go")
	content := `//go:build feature

package main

// @feature journey
func TestMain(t *testing.T) {}`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	filesModified, err := promoteJourneyTags(journeyDir)
	if err != nil {
		t.Fatalf("promoteJourneyTags failed: %v", err)
	}

	if filesModified != 1 {
		t.Errorf("expected 1 file modified, got %d", filesModified)
	}

	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	result := string(data)
	if strings.Contains(result, "@feature") {
		t.Error("@feature should be replaced with @regression")
	}
	if !strings.Contains(result, "@regression") {
		t.Error("expected @regression in promoted file")
	}
}

func TestPromoteJourneyTags_SkipsContractsDir(t *testing.T) {
	dir := t.TempDir()
	journeyDir := filepath.Join(dir, "my-journey")
	contractsDir := filepath.Join(journeyDir, "_contracts")
	if err := os.MkdirAll(contractsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Put a file with @feature in _contracts - should NOT be modified
	contractFile := filepath.Join(contractsDir, "step-1.md")
	content := `# @feature contract`
	if err := os.WriteFile(contractFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	filesModified, err := promoteJourneyTags(journeyDir)
	if err != nil {
		t.Fatalf("promoteJourneyTags failed: %v", err)
	}

	if filesModified != 0 {
		t.Errorf("_contracts directory should be skipped, got %d files modified", filesModified)
	}

	// Verify the contract file was not modified
	data, err := os.ReadFile(contractFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "@feature") {
		t.Error("_contracts file should NOT have @feature replaced")
	}
}

func TestPromoteJourneyTags_SkipsNonTestFiles(t *testing.T) {
	dir := t.TempDir()
	journeyDir := filepath.Join(dir, "my-journey")
	if err := os.MkdirAll(journeyDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Non-test file with @feature - should NOT be modified
	readme := filepath.Join(journeyDir, "README.md")
	if err := os.WriteFile(readme, []byte("# @feature test\n"), 0644); err != nil {
		t.Fatal(err)
	}

	filesModified, err := promoteJourneyTags(journeyDir)
	if err != nil {
		t.Fatalf("promoteJourneyTags failed: %v", err)
	}

	if filesModified != 0 {
		t.Errorf("non-test files should be skipped, got %d files modified", filesModified)
	}
}

func TestPromoteJourneyTags_MultipleLanguages(t *testing.T) {
	dir := t.TempDir()
	journeyDir := filepath.Join(dir, "multi-journey")
	if err := os.MkdirAll(journeyDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Go test file with @feature tag
	goFile := filepath.Join(journeyDir, "main_test.go")
	goContent := "//go:build feature\n\npackage main\n\n// @feature journey\nfunc TestMain(t *testing.T) {}\n"
	if err := os.WriteFile(goFile, []byte(goContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Python test file with @feature tag (generic annotation, not pytest.mark)
	pyFile := filepath.Join(journeyDir, "test_main.py")
	pyContent := "# @feature journey\ndef test_main(): pass\n"
	if err := os.WriteFile(pyFile, []byte(pyContent), 0644); err != nil {
		t.Fatal(err)
	}

	filesModified, err := promoteJourneyTags(journeyDir)
	if err != nil {
		t.Fatalf("promoteJourneyTags failed: %v", err)
	}

	if filesModified != 2 {
		t.Errorf("expected 2 files modified, got %d", filesModified)
	}

	// Verify Go file got regression build tag
	goData, _ := os.ReadFile(goFile)
	if !strings.Contains(string(goData), "//go:build regression") {
		t.Error("Go file should have //go:build regression after promotion")
	}

	// Verify Python file got @regression tag
	pyData, _ := os.ReadFile(pyFile)
	if !strings.Contains(string(pyData), "@regression") {
		t.Error("Python file should have @regression after promotion")
	}
}

// --- Test: PromoteDiffSummary ---

// --- Test: validateJourneyName path traversal ---

func TestValidateJourneyName_RejectsPathTraversal(t *testing.T) {
	err := validateJourneyName("../other-journey")
	assert.Equal(t, base.ErrInvalidPath, err.Code)
	assert.Contains(t, err.Cause, "..")
}

func TestValidateJourneyName_RejectsDoubleDot(t *testing.T) {
	err := validateJourneyName("foo/../../bar")
	assert.Equal(t, base.ErrInvalidPath, err.Code)
	assert.Contains(t, err.Cause, "..")
}

func TestValidateJourneyName_RejectsAbsolutePath(t *testing.T) {
	err := validateJourneyName("/etc/passwd")
	assert.Equal(t, base.ErrInvalidPath, err.Code)
}

func TestValidateJourneyName_AcceptsSimpleName(t *testing.T) {
	err := validateJourneyName("my-journey")
	assert.Nil(t, err)
}

func TestValidateJourneyName_AcceptsHyphenatedName(t *testing.T) {
	err := validateJourneyName("task-lifecycle")
	assert.Nil(t, err)
}

func TestValidateJourneyName_RejectsJustDotDot(t *testing.T) {
	err := validateJourneyName("..")
	assert.Equal(t, base.ErrInvalidPath, err.Code)
}

func TestPromoteDiffSummary_ShowsChanges(t *testing.T) {
	dir := t.TempDir()
	journeyDir := filepath.Join(dir, "diff-journey")
	if err := os.MkdirAll(journeyDir, 0755); err != nil {
		t.Fatal(err)
	}

	goFile := filepath.Join(journeyDir, "main_test.go")
	content := `// @feature journey
package main
func TestMain(t *testing.T) {}`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	buf, err := PromoteDiffSummary(journeyDir)
	if err != nil {
		t.Fatalf("PromoteDiffSummary failed: %v", err)
	}

	diff := buf.String()
	if !strings.Contains(diff, "--- main_test.go") {
		t.Errorf("expected diff to contain file name, got: %s", diff)
	}
	if !strings.Contains(diff, "-// @feature") {
		t.Errorf("expected diff to show removed @feature, got: %s", diff)
	}
	if !strings.Contains(diff, "+// @regression") {
		t.Errorf("expected diff to show added @regression, got: %s", diff)
	}
}
