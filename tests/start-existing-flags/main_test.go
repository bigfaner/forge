//go:build cli_functional

// Package startexistingflags tests the start-existing-flags Journey:
// verifying flag behavior (--source-branch, --no-launch, --interactive) on existing worktrees.
//
// @feature worktree-start-idempotent
// @cli-functional
package startexistingflags

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	testkit "forge-tests/testkit"
)

func TestMain(m *testing.M) {
	_ = testkit.ForgeBinary
	m.Run()
}

// setupGitRepoWithForge creates a temp git repo with forge config.
func setupGitRepoWithForge(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %v\n%s", err, out)
	}

	cmd = exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = dir
	_ = cmd.Run()
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = dir
	_ = cmd.Run()

	readmePath := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readmePath, []byte("# test\n"), 0644); err != nil {
		t.Fatalf("failed to write README: %v", err)
	}

	cmd = exec.Command("git", "add", "README.md")
	cmd.Dir = dir
	_ = cmd.Run()
	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = dir
	_ = cmd.Run()

	forgeDir := filepath.Join(dir, ".forge")
	if err := os.MkdirAll(forgeDir, 0755); err != nil {
		t.Fatalf("failed to create .forge: %v", err)
	}
	configContent := `version: "1"
surfaces: cli
`
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	return dir
}

// setupGitRepoWithIncludes creates a temp git repo with includes configured.
func setupGitRepoWithIncludes(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	_, _ = cmd.CombinedOutput()
	cmd = exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = dir
	_ = cmd.Run()
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = dir
	_ = cmd.Run()

	readmePath := filepath.Join(dir, "README.md")
	_ = os.WriteFile(readmePath, []byte("# test\n"), 0644)
	cmd = exec.Command("git", "add", "README.md")
	cmd.Dir = dir
	_ = cmd.Run()
	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = dir
	_ = cmd.Run()

	// Create include file
	secretFile := filepath.Join(dir, "secret.txt")
	_ = os.WriteFile(secretFile, []byte("secret\n"), 0644)
	cmd = exec.Command("git", "add", "secret.txt")
	cmd.Dir = dir
	_ = cmd.Run()
	cmd = exec.Command("git", "commit", "-m", "add secret")
	cmd.Dir = dir
	_ = cmd.Run()

	forgeDir := filepath.Join(dir, ".forge")
	_ = os.MkdirAll(forgeDir, 0755)
	configContent := `version: "1"
surfaces: cli
worktree:
  includes:
    - secret.txt
`
	_ = os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0644)
	return dir
}

// runForgeStartNoLaunch runs forge worktree start with --no-launch.
func runForgeStartNoLaunch(t *testing.T, projectRoot, slug string, extraArgs ...string) (string, string, int) {
	t.Helper()
	args := []string{"worktree", "start"}
	if slug != "" {
		args = append(args, slug)
	}
	args = append(args, "--no-launch")
	args = append(args, extraArgs...)
	cmd := exec.Command(testkit.ForgeBinary, args...)
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)

	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()
	_ = cmd.Start()

	var stdoutBuf, stderrBuf []byte
	done := make(chan struct{})
	go func() {
		stdoutBuf, _ = io.ReadAll(stdoutPipe)
		stderrBuf, _ = io.ReadAll(stderrPipe)
		close(done)
	}()
	<-done

	exitCode := 0
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return string(stdoutBuf), string(stderrBuf), exitCode
}

// runForgeStartInteractive runs forge worktree start with --interactive.
// Inherits parent stdin (TTY-aware): works correctly for tests that expect TTY behavior.
func runForgeStartInteractive(t *testing.T, projectRoot string, extraArgs ...string) (string, string, int) {
	t.Helper()
	args := []string{"worktree", "start", "--interactive", "--no-launch"}
	args = append(args, extraArgs...)
	cmd := exec.Command(testkit.ForgeBinary, args...)
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)

	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()
	_ = cmd.Start()

	var stdoutBuf, stderrBuf []byte
	done := make(chan struct{})
	go func() {
		stdoutBuf, _ = io.ReadAll(stdoutPipe)
		stderrBuf, _ = io.ReadAll(stderrPipe)
		close(done)
	}()
	<-done

	exitCode := 0
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return string(stdoutBuf), string(stderrBuf), exitCode
}

// runForgeStartInteractiveNonTTY runs forge worktree start with --interactive
// with stdin piped to force non-TTY detection.
func runForgeStartInteractiveNonTTY(t *testing.T, projectRoot string, extraArgs ...string) (string, string, int) {
	t.Helper()
	args := []string{"worktree", "start", "--interactive", "--no-launch"}
	args = append(args, extraArgs...)
	cmd := exec.Command(testkit.ForgeBinary, args...)
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	// Use an os.Pipe (not /dev/null) to force non-TTY detection.
	// On macOS, /dev/null is a char device and would pass the TTY check.
	// An os.Pipe creates a regular pipe (not char device), so os.Stdin.Stat()
	// in the subprocess will NOT have ModeCharDevice set.
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	defer pr.Close()
	defer pw.Close()
	cmd.Stdin = pr
	// Close write end in the parent so the subprocess reads EOF immediately
	pw.Close()

	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()
	_ = cmd.Start()

	var stdoutBuf, stderrBuf []byte
	done := make(chan struct{})
	go func() {
		stdoutBuf, _ = io.ReadAll(stdoutPipe)
		stderrBuf, _ = io.ReadAll(stderrPipe)
		close(done)
	}()
	<-done

	exitCode := 0
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return string(stdoutBuf), string(stderrBuf), exitCode
}

// worktreeDir returns the expected worktree path.
func worktreeDir(projectRoot, slug string) string {
	return filepath.Join(projectRoot, ".forge", "worktrees", slug)
}

// assertValidGitFile checks .git file exists in worktree.
func assertValidGitFile(t *testing.T, wtDir string) {
	t.Helper()
	gitFile := filepath.Join(wtDir, ".git")
	info, err := os.Stat(gitFile)
	if err != nil {
		t.Fatalf(".git file missing: %v", err)
	}
	if info.IsDir() {
		t.Fatalf(".git should be a file, not a directory")
	}
}

// runGitIn runs a git command in the given directory.
func runGitIn(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}
