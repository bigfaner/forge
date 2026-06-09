//go:build cli_functional

// Package idempotentstart tests the idempotent-start Journey:
// verifying that forge worktree start is idempotent -- creating a new worktree
// on first invocation and entering the existing worktree on subsequent invocations.
//
// @feature worktree-start-idempotent
// @cli-functional
package idempotentstart

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	testkit "forge-tests/testkit"
)

// ForgeBinary is shared across all tests in this package.
var forgeBinary string

func TestMain(m *testing.M) {
	// Ensure forge binary is built via testkit init
	forgeBinary = testkit.ForgeBinary
	_ = forgeBinary
	m.Run()
}

// --------------------------------------------------------------------------
// Shared helpers for idempotent-start journey tests
// --------------------------------------------------------------------------

// setupGitRepo creates a temporary git repository with forge initialized.
// Returns the repo directory path.
func setupGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Initialize git repo
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@test.com")
	runGit(t, dir, "config", "user.name", "Test User")

	// Create initial commit so HEAD exists
	readmePath := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readmePath, []byte("# test project\n"), 0644); err != nil {
		t.Fatalf("failed to write README.md: %v", err)
	}
	runGit(t, dir, "add", "README.md")
	runGit(t, dir, "commit", "-m", "initial commit")

	return dir
}

// setupGitRepoWithForge creates a git repo with forge config (no includes).
func setupGitRepoWithForge(t *testing.T) string {
	t.Helper()
	dir := setupGitRepo(t)
	writeForgeConfig(t, dir, "")
	return dir
}

// setupGitRepoWithIncludes creates a git repo with forge config that lists includes.
func setupGitRepoWithIncludes(t *testing.T) string {
	t.Helper()
	dir := setupGitRepo(t)

	// Create include files
	secretFile := filepath.Join(dir, "secret.txt")
	if err := os.WriteFile(secretFile, []byte("secret-value\n"), 0644); err != nil {
		t.Fatalf("failed to write secret.txt: %v", err)
	}
	runGit(t, dir, "add", "secret.txt")
	runGit(t, dir, "commit", "-m", "add secret file")

	writeForgeConfig(t, dir, "secret.txt")
	return dir
}

// writeForgeConfig writes a .forge/config.yaml with optional includes.
func writeForgeConfig(t *testing.T, projectRoot, includes string) {
	t.Helper()
	forgeDir := filepath.Join(projectRoot, ".forge")
	if err := os.MkdirAll(forgeDir, 0755); err != nil {
		t.Fatalf("failed to create .forge dir: %v", err)
	}

	configContent := `version: "1"
surfaces: cli
`
	if includes != "" {
		configContent += `worktree:
  includes:
    - ` + includes + "\n"
	}

	configPath := filepath.Join(forgeDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config.yaml: %v", err)
	}
}

// runGit executes a git command in the given directory.
func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_DATE=2026-01-01T00:00:00Z", "GIT_COMMITTER_DATE=2026-01-01T00:00:00Z")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

// forgeStart runs "forge worktree start <slug>" with CLAUDE_PROJECT_DIR set.
// Returns stdout, stderr (separated), and exit code.
func forgeStart(t *testing.T, projectRoot, slug string, extraArgs ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	args := []string{"worktree", "start", slug}
	args = append(args, extraArgs...)
	cmd := exec.Command(testkit.ForgeBinary, args...)
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)

	var stdoutBuf, stderrBuf []byte
	var stdoutErr, stderrErr error

	// Use separate pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("failed to create stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start command: %v", err)
	}

	// Read stdout and stderr concurrently
	done := make(chan struct{})
	go func() {
		stdoutBuf, stdoutErr = readAll(stdoutPipe)
		stderrBuf, stderrErr = readAll(stderrPipe)
		close(done)
	}()
	<-done

	if stdoutErr != nil {
		t.Fatalf("failed to read stdout: %v", stdoutErr)
	}
	if stderrErr != nil {
		t.Fatalf("failed to read stderr: %v", stderrErr)
	}

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return string(stdoutBuf), string(stderrBuf), exitCode
}

// forgeStartNoLaunch runs "forge worktree start [slug] --no-launch" with CLAUDE_PROJECT_DIR set.
// This avoids needing claude binary in PATH.
func forgeStartNoLaunch(t *testing.T, projectRoot, slug string, extraArgs ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	args := []string{"worktree", "start"}
	if slug != "" {
		args = append(args, slug)
	}
	args = append(args, "--no-launch")
	args = append(args, extraArgs...)
	cmd := exec.Command(testkit.ForgeBinary, args...)
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)

	stdoutBuf, stderrBuf, exitCode := runCommandWithSeparateOutputs(cmd)
	return stdoutBuf, stderrBuf, exitCode
}

// forgeRemove runs "forge worktree remove <slug>" with CLAUDE_PROJECT_DIR set.
func forgeRemove(t *testing.T, projectRoot, slug string, extraArgs ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	args := []string{"worktree", "remove", slug}
	args = append(args, extraArgs...)
	cmd := exec.Command(testkit.ForgeBinary, args...)
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)

	stdoutBuf, stderrBuf, exitCode := runCommandWithSeparateOutputs(cmd)
	return stdoutBuf, stderrBuf, exitCode
}

// runCommandWithSeparateOutputs runs a command and returns separated stdout, stderr, and exit code.
func runCommandWithSeparateOutputs(cmd *exec.Cmd) (string, string, int) {
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	_ = cmd.Start()

	var stdoutBuf, stderrBuf []byte
	done := make(chan struct{})
	go func() {
		stdoutBuf, _ = readAll(stdoutPipe)
		stderrBuf, _ = readAll(stderrPipe)
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

// readAll reads all data from an io.Reader.
func readAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

// worktreeDir returns the expected worktree directory path.
func worktreeDir(projectRoot, slug string) string {
	return filepath.Join(projectRoot, ".forge", "worktrees", slug)
}

// assertValidGitFile checks that .git file exists in the worktree dir.
func assertValidGitFile(t *testing.T, worktreePath string) {
	t.Helper()
	gitFile := filepath.Join(worktreePath, ".git")
	info, err := os.Stat(gitFile)
	if err != nil {
		t.Fatalf(".git file missing in worktree %s: %v", worktreePath, err)
	}
	if info.IsDir() {
		t.Fatalf(".git should be a file (not a directory) in worktree %s", worktreePath)
	}
}

// assertFileContentsMatch checks that two files have identical contents.
func assertFileContentsMatch(t *testing.T, file1, file2 string) {
	t.Helper()
	content1, err := os.ReadFile(file1)
	if err != nil {
		t.Fatalf("failed to read %s: %v", file1, err)
	}
	content2, err := os.ReadFile(file2)
	if err != nil {
		t.Fatalf("failed to read %s: %v", file2, err)
	}
	if string(content1) != string(content2) {
		t.Fatalf("file contents mismatch:\n%s (%d bytes)\n%s (%d bytes)", file1, len(content1), file2, len(content2))
	}
}

// createForgeCommand creates an exec.Cmd for forge with given args.
func createForgeCommand(args ...string) *exec.Cmd {
	return exec.Command(testkit.ForgeBinary, args...)
}

// execGit runs a git command and returns its combined output.
func execGit(dir string, args ...string) []byte {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, _ := cmd.CombinedOutput()
	return out
}
