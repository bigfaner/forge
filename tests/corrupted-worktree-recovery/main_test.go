//go:build cli_functional

// Package corruptedworktreerecovery tests the corrupted-worktree-recovery Journey:
// verifying that corrupted worktrees are detected, removed, and can be re-created.
//
// @feature worktree-start-idempotent
// @cli-functional
package corruptedworktreerecovery

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

	forgeDir := filepath.Join(dir, ".forge")
	_ = os.MkdirAll(forgeDir, 0755)
	configContent := `version: "1"
surfaces: cli
`
	_ = os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0644)
	return dir
}

// runForgeStartNoLaunch runs forge worktree start with --no-launch.
func runForgeStartNoLaunch(t *testing.T, projectRoot, slug string, extraArgs ...string) (string, string, int) {
	t.Helper()
	args := []string{"worktree", "start", slug, "--no-launch"}
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

// runForgeRemove runs forge worktree remove with given slug and extra args.
func runForgeRemove(t *testing.T, projectRoot, slug string, extraArgs ...string) (string, string, int) {
	t.Helper()
	args := []string{"worktree", "remove", slug}
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

// createCorruptedWorktree creates a directory at the worktree path without .git file.
func createCorruptedWorktree(t *testing.T, projectRoot, slug string) string {
	t.Helper()
	wtDir := worktreeDir(projectRoot, slug)
	if err := os.MkdirAll(wtDir, 0755); err != nil {
		t.Fatalf("failed to create corrupted worktree dir: %v", err)
	}
	return wtDir
}

// createOrphanWorktree creates a directory that git doesn't recognize as a worktree.
func createOrphanWorktree(t *testing.T, projectRoot, slug string) string {
	t.Helper()
	wtDir := createCorruptedWorktree(t, projectRoot, slug)
	// Write a random file to make it non-empty
	_ = os.WriteFile(filepath.Join(wtDir, "orphan.txt"), []byte("orphan content\n"), 0644)
	return wtDir
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
