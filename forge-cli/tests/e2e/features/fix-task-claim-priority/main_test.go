//go:build e2e

package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"forge-cli/tests/e2e/testkit"
)

// forgeBinary is the path to a forge CLI binary built from the current source tree.
// Built once in TestMain and used by all e2e tests to ensure they test the code
// on the current branch, not the system-installed binary.
var forgeBinary string

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "forge-e2e-binary-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp dir for forge binary: %s\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	forgeBinary = filepath.Join(tmpDir, "forge-test")

	buildCmd := exec.Command("go", "build", "-o", forgeBinary, "./cmd/forge")
	buildCmd.Dir = forgeCLIDir()
	if out, err := buildCmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "forge binary build failed: %s\n%s", err, out)
		os.Exit(1)
	}

	// Propagate binary path to testkit so its public helpers also use the built binary.
	testkit.SetForgeBinary(forgeBinary)

	code := m.Run()
	os.Exit(code)
}

// forgeCLIDir returns the forge-cli module root directory by walking up from the
// current working directory to find the directory containing cmd/forge.
func forgeCLIDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "cmd", "forge")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find forge-cli module root (no cmd/forge found)")
		}
		dir = parent
	}
}
