//go:build e2e

package justfile_canonical_e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "forge-e2e-justfile-binary-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp dir for forge binary: %s\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	forgeBinaryPath = filepath.Join(tmpDir, "forge-test")

	buildCmd := exec.Command("go", "build", "-o", forgeBinaryPath, "./cmd/forge")
	buildCmd.Dir = forgeCLIDir()
	if out, err := buildCmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "forge binary build failed: %s\n%s", err, out)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

// forgeCLIDir returns the forge-cli module root directory.
// It walks up from the test module directory to find forge-cli/cmd/forge.
func forgeCLIDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "forge-cli", "cmd", "forge")); err == nil {
			return filepath.Join(dir, "forge-cli")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find forge-cli module root")
		}
		dir = parent
	}
}
