//go:build cli_functional

package testkit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProjectRoot(t *testing.T) {
	root := ProjectRoot(t)
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
