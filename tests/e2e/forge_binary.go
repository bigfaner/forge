//go:build e2e

package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// ForgeBinary is the path to a forge CLI binary built from the current source tree.
// Built once at package init and used by all e2e tests to ensure they test the code
// on the current branch, not the system-installed binary.
// Sub-packages import this via `import e2etests "e2e-tests"`.
var ForgeBinary string

var forgeBinaryOnce sync.Once

func init() {
	forgeBinaryOnce.Do(func() {
		tmpDir, err := os.MkdirTemp("", "forge-e2e-binary-*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create temp dir for forge binary: %s\n", err)
			os.Exit(1)
		}

		ForgeBinary = filepath.Join(tmpDir, "forge-test")

		buildCmd := exec.Command("go", "build", "-o", ForgeBinary, "./cmd/forge")
		buildCmd.Dir = findForgeCLIDir()
		if out, err := buildCmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "forge binary build failed: %s\n%s", err, out)
			os.Exit(1)
		}
	})
}

// ForgeCmd returns an exec.Cmd for the forge CLI binary built from source.
func ForgeCmd(args ...string) *exec.Cmd {
	return exec.Command(ForgeBinary, args...)
}

// findForgeCLIDir returns the forge-cli module root directory.
func findForgeCLIDir() string {
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
