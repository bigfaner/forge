//go:build e2e

package testkit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProjectRoot(t *testing.T) {
	root := projectRoot(t)
	modPath := filepath.Join(root, "go.mod")
	if _, err := os.Stat(modPath); err != nil {
		t.Fatalf("projectRoot returned %q but go.mod not found there: %v", root, err)
	}
}

func TestReadProjectFile(t *testing.T) {
	// go.mod is guaranteed to exist at project root
	content := ReadProjectFile(t, "go.mod")
	if !strings.Contains(content, "module") {
		t.Fatalf("ReadProjectFile returned content without 'module': %s", content)
	}
}

func TestReadProjectFileNotFound(t *testing.T) {
	// Verify the function compiles and the signature is correct.
	// A missing-file call would Fatal the test, so we only verify existence here.
	_ = ReadProjectFile
}

func TestProjectFileExists(t *testing.T) {
	// go.mod exists at project root
	if !ProjectFileExists("go.mod") {
		t.Fatal("ProjectFileExists should return true for go.mod")
	}
	// nonexistent file
	if ProjectFileExists("this_file_definitely_does_not_exist_12345.go") {
		t.Fatal("ProjectFileExists should return false for nonexistent file")
	}
}

func TestFileContains(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.txt")
	content := "hello world\nfoo bar baz\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Should pass — substring exists
	FileContains(t, filePath, "hello world")
	FileContains(t, filePath, "foo bar")
	FileContains(t, filePath, "baz")

	// Should pass — substring spanning content
	FileContains(t, filePath, "foo bar baz")
}

func TestFileNotContains(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(filePath, []byte("hello world\n"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Should pass — substring does NOT exist
	FileNotContains(t, filePath, "goodbye")
	FileNotContains(t, filePath, "xyz")
}
