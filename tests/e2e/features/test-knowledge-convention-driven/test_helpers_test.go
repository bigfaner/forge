//go:build e2e

package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// forgeBinary is the path to a forge CLI binary built from the current source tree.
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

	os.Exit(m.Run())
}

// forgeCLIDir returns the forge-cli module root directory containing go.mod and cmd/forge.
func forgeCLIDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "forge-cli", "cmd", "forge")); err == nil {
			return filepath.Join(dir, "forge-cli")
		}
		if _, err := os.Stat(filepath.Join(dir, "cmd", "forge")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find forge-cli module root")
		}
		dir = parent
	}
}

// forgeCmd returns an exec.Cmd for the forge CLI binary built from source.
func forgeCmd(args ...string) *exec.Cmd {
	return exec.Command(forgeBinary, args...)
}
