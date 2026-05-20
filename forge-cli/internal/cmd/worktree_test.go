package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/feature"

	gitPkg "forge-cli/pkg/git"

	"github.com/spf13/cobra"
)

// ---------------------------------------------------------------------------
// worktree command group registration
// ---------------------------------------------------------------------------

func TestWorktreeCmd_RegisteredAsGroup(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "worktree" {
			found = true
			// Parent command with no Run — only subcommands
			if cmd.Run != nil && cmd.RunE != nil {
				t.Error("worktreeCmd should have no Run/RunE (group parent only)")
			}
			break
		}
	}
	if !found {
		t.Error("worktree command should be registered as top-level command")
	}
}

func TestWorktreeCmd_HasStartSubcommand(t *testing.T) {
	subcommands := worktreeCmd.Commands()
	found := false
	for _, cmd := range subcommands {
		if cmd.Name() == "start" {
			found = true
			break
		}
	}
	if !found {
		t.Error("worktree group should have 'start' subcommand")
	}
}

func TestWorktreeStartCmd_RequiresSlugArg(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when slug argument is missing")
	}
}

// ---------------------------------------------------------------------------
// worktree start: pre-flight claude check
// ---------------------------------------------------------------------------

func TestWorktreeStart_ErrorWhenClaudeNotInPath(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(_ string) (string, error) {
		return "", &exec.Error{Name: "claude", Err: exec.ErrNotFound}
	}
	defer func() { lookPathFunc = origLookPath }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", "test-slug"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when claude binary not found")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "claude") {
		t.Errorf("error should mention 'claude', got: %s", stderr)
	}
}

// ---------------------------------------------------------------------------
// worktree start: directory conflict checks
// ---------------------------------------------------------------------------

func TestWorktreeStart_ErrorWhenTargetDirExists(t *testing.T) {
	resetSourceBranchFlag(t)

	// Make claude available
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	// Don't actually launch claude
	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create the target directory ahead of time at .forge/worktrees/<slug>
	targetDir := filepath.Join(dir, ".forge", "worktrees", "test-slug")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatalf("create target dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Dir(filepath.Dir(targetDir))) })

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", "test-slug"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when target directory already exists")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "already exists") {
		t.Errorf("error should mention 'already exists', got: %s", stderr)
	}
}

// ---------------------------------------------------------------------------
// worktree start: happy path
// ---------------------------------------------------------------------------

func TestWorktreeStart_CreatesWorktreeAndLaunchesClaude(t *testing.T) {
	resetSourceBranchFlag(t)

	// Make claude available
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	// Capture claude launch args and just succeed
	var capturedArgs []string
	origRunClaude := runClaudeFunc
	runClaudeFunc = func(args []string) error {
		capturedArgs = args
		return nil
	}
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "my-feature"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		// Clean up the worktree
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Errorf("worktree directory %s should exist", targetDir)
	}

	// Verify claude was launched with --dangerously-skip-permissions
	if len(capturedArgs) == 0 {
		t.Fatal("claude should have been launched")
	}
	if capturedArgs[0] != "--dangerously-skip-permissions" {
		t.Errorf("first arg should be --dangerously-skip-permissions, got %q", capturedArgs[0])
	}
}

// ---------------------------------------------------------------------------
// worktree start: resume from existing branch
// ---------------------------------------------------------------------------

func TestWorktreeStart_ResumesFromExistingBranch(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	var capturedArgs []string
	origRunClaude := runClaudeFunc
	runClaudeFunc = func(args []string) error {
		capturedArgs = args
		return nil
	}
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "existing-branch"
	// Create the branch ahead of time
	if err := exec.Command("git", "-C", dir, "branch", slug).Run(); err != nil {
		t.Fatalf("git branch %s: %v", slug, err)
	}

	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Errorf("worktree directory %s should exist", targetDir)
	}

	// Verify claude was launched
	if len(capturedArgs) == 0 {
		t.Fatal("claude should have been launched")
	}
}

// ---------------------------------------------------------------------------
// worktree start: remote branch resolution
// ---------------------------------------------------------------------------

func TestWorktreeStart_CreatesFromRemoteBranch(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a "remote" repo, push a branch to it, then the slug branch
	// only exists on the remote (origin), not locally.
	remoteDir := t.TempDir()
	if err := exec.Command("git", "init", "--bare", remoteDir).Run(); err != nil {
		t.Fatalf("git init --bare remote: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "remote", "add", "origin", remoteDir).Run(); err != nil {
		t.Fatalf("git remote add: %v", err)
	}

	// Create a branch with a distinct commit, push it to origin
	slug := "remote-branch-test"
	if err := exec.Command("git", "-C", dir, "checkout", "-b", slug).Run(); err != nil {
		t.Fatalf("git checkout -b %s: %v", slug, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "remote.txt"), []byte("from remote"), 0o644); err != nil {
		t.Fatalf("write remote.txt: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "remote commit").Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}
	// Push to origin
	if err := exec.Command("git", "-C", dir, "push", "origin", slug).Run(); err != nil {
		t.Fatalf("git push origin %s: %v", slug, err)
	}

	// Delete the LOCAL branch (go back to master first)
	if err := exec.Command("git", "-C", dir, "checkout", "master").Run(); err != nil {
		t.Fatalf("git checkout master: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "branch", "-D", slug).Run(); err != nil {
		t.Fatalf("git branch -D %s: %v", slug, err)
	}

	// Now: origin/<slug> exists, local <slug> does not.
	// forge worktree start <slug> should detect the remote branch.
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Errorf("worktree directory %s should exist", targetDir)
	}

	// Verify the worktree has the remote.txt file (proving it was based on the remote branch)
	if _, err := os.Stat(filepath.Join(targetDir, "remote.txt")); os.IsNotExist(err) {
		t.Errorf("worktree should have remote.txt (created from remote branch origin/%s)", slug)
	}

	// Verify stdout mentions remote branch
	stdout := buf.String()
	if !strings.Contains(stdout, "origin/"+slug) {
		t.Errorf("stdout should mention remote branch, got: %s", stdout)
	}
}

func TestWorktreeStart_RemoteBranchIgnoresSourceBranch(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Set up remote
	remoteDir := t.TempDir()
	if err := exec.Command("git", "init", "--bare", remoteDir).Run(); err != nil {
		t.Fatalf("git init --bare remote: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "remote", "add", "origin", remoteDir).Run(); err != nil {
		t.Fatalf("git remote add: %v", err)
	}

	// Create "develop" branch with a file
	if err := exec.Command("git", "-C", dir, "checkout", "-b", "develop").Run(); err != nil {
		t.Fatalf("git checkout -b develop: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "develop.txt"), []byte("develop"), 0o644); err != nil {
		t.Fatalf("write develop.txt: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "develop commit").Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}

	// Create the slug branch from develop (has develop.txt + remote-marker.txt)
	slug := "remote-slug"
	if err := exec.Command("git", "-C", dir, "checkout", "-b", slug).Run(); err != nil {
		t.Fatalf("git checkout -b %s: %v", slug, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "remote-marker.txt"), []byte("from remote"), 0o644); err != nil {
		t.Fatalf("write remote-marker.txt: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "slug commit").Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}

	// Push slug branch to origin
	if err := exec.Command("git", "-C", dir, "push", "origin", slug).Run(); err != nil {
		t.Fatalf("git push origin %s: %v", slug, err)
	}

	// Push develop to origin too (so --source-branch is valid)
	if err := exec.Command("git", "-C", dir, "push", "origin", "develop").Run(); err != nil {
		t.Fatalf("git push origin develop: %v", err)
	}

	// Delete local slug branch only (keep develop local for source-branch validation)
	if err := exec.Command("git", "-C", dir, "checkout", "develop").Run(); err != nil {
		t.Fatalf("git checkout develop: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "branch", "-D", slug).Run(); err != nil {
		t.Fatalf("git branch -D %s: %v", slug, err)
	}

	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", "develop").Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	// Use --source-branch develop, but remote branch should take priority
	rootCmd.SetArgs([]string{"worktree", "start", slug, "--source-branch", "develop"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created from remote branch (has remote-marker.txt)
	if _, err := os.Stat(filepath.Join(targetDir, "remote-marker.txt")); os.IsNotExist(err) {
		t.Errorf("worktree should have remote-marker.txt (created from remote branch, ignoring --source-branch)")
	}
}

func TestWorktreeStart_FetchFailureDoesNotBlockWorktree(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Add a fake remote that will fail to fetch (nonexistent path)
	if err := exec.Command("git", "-C", dir, "remote", "add", "origin", "/nonexistent/path/to/remote").Run(); err != nil {
		t.Fatalf("git remote add: %v", err)
	}

	slug := "fetch-fail-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("fetch failure should not block worktree creation, got error: %v", err)
	}

	// Verify worktree was created (fallback to HEAD)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Errorf("worktree directory %s should exist even after fetch failure", targetDir)
	}
}

func TestWorktreeStart_LocalBranchTakesPriorityOverRemote(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Set up remote with a branch
	remoteDir := t.TempDir()
	if err := exec.Command("git", "init", "--bare", remoteDir).Run(); err != nil {
		t.Fatalf("git init --bare remote: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "remote", "add", "origin", remoteDir).Run(); err != nil {
		t.Fatalf("git remote add: %v", err)
	}

	slug := "priority-test"

	// Create remote branch with "remote-marker.txt"
	if err := exec.Command("git", "-C", dir, "checkout", "-b", slug).Run(); err != nil {
		t.Fatalf("git checkout -b %s: %v", slug, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "remote-marker.txt"), []byte("remote"), 0o644); err != nil {
		t.Fatalf("write remote-marker.txt: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "remote version").Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "push", "origin", slug).Run(); err != nil {
		t.Fatalf("git push origin %s: %v", slug, err)
	}

	// Now modify local branch to have "local-marker.txt" instead
	if err := os.WriteFile(filepath.Join(dir, "local-marker.txt"), []byte("local"), 0o644); err != nil {
		t.Fatalf("write local-marker.txt: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "local version").Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}

	// Go back to master; local branch still exists
	if err := exec.Command("git", "-C", dir, "checkout", "master").Run(); err != nil {
		t.Fatalf("git checkout master: %v", err)
	}

	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Local branch takes priority: worktree should have local-marker.txt
	if _, err := os.Stat(filepath.Join(targetDir, "local-marker.txt")); os.IsNotExist(err) {
		t.Errorf("worktree should have local-marker.txt (local branch takes priority over remote)")
	}

	// Should NOT mention remote branch in stdout (since local was used)
	stdout := buf.String()
	if strings.Contains(stdout, "origin/") {
		t.Errorf("stdout should NOT mention remote branch when local exists, got: %s", stdout)
	}
}

// ---------------------------------------------------------------------------
// worktree start: not in a git repo
// ---------------------------------------------------------------------------

func TestWorktreeStart_ErrorWhenNotGitRepo(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", "test-slug"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when not in a git repository")
	}
}

// ---------------------------------------------------------------------------
// worktree start: --source-branch flag
// ---------------------------------------------------------------------------

func TestWorktreeStart_SourceBranchFlag_CreatesFromSpecifiedBranch(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a "develop" branch with a distinct commit
	if err := exec.Command("git", "-C", dir, "checkout", "-b", "develop").Run(); err != nil {
		t.Fatalf("git checkout -b develop: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "develop.txt"), []byte("develop content"), 0o644); err != nil {
		t.Fatalf("write develop.txt: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "develop commit").Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}
	// Go back to main (master)
	if err := exec.Command("git", "-C", dir, "checkout", "master").Run(); err != nil {
		t.Fatalf("git checkout master: %v", err)
	}

	slug := "source-branch-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", "develop").Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug, "--source-branch", "develop"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Errorf("worktree directory %s should exist", targetDir)
	}

	// Verify the worktree has the develop.txt file (proving it was based on develop)
	if _, err := os.Stat(filepath.Join(targetDir, "develop.txt")); os.IsNotExist(err) {
		t.Errorf("worktree should have develop.txt (created from develop branch)")
	}
}

func TestWorktreeStart_SourceBranchShortFlag(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a "release" branch
	if err := exec.Command("git", "-C", dir, "branch", "release").Run(); err != nil {
		t.Fatalf("git branch release: %v", err)
	}

	slug := "short-flag-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", "release").Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug, "-b", "release"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Errorf("worktree directory %s should exist", targetDir)
	}
}

func TestWorktreeStart_SourceBranchErrorWhenBranchNotFound(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", "test-slug", "--source-branch", "nonexistent-branch"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when source branch does not exist")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "nonexistent-branch") {
		t.Errorf("error should mention the branch name, got: %s", stderr)
	}
}

func TestWorktreeStart_SourceBranchFromConfig(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a "staging" branch with a distinct file
	if err := exec.Command("git", "-C", dir, "checkout", "-b", "staging").Run(); err != nil {
		t.Fatalf("git checkout -b staging: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "staging.txt"), []byte("staging content"), 0o644); err != nil {
		t.Fatalf("write staging.txt: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "staging commit").Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "checkout", "master").Run(); err != nil {
		t.Fatalf("git checkout master: %v", err)
	}

	// Create .forge/config.yaml with worktree.source-branch
	forgeDir := filepath.Join(dir, ".forge")
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		t.Fatalf("mkdir .forge: %v", err)
	}
	configContent := "worktree:\n  source-branch: staging\n"
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
		t.Fatalf("write config.yaml: %v", err)
	}

	slug := "config-source-branch"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", "staging").Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created from staging (has staging.txt)
	if _, err := os.Stat(filepath.Join(targetDir, "staging.txt")); os.IsNotExist(err) {
		t.Errorf("worktree should have staging.txt (created from staging branch via config)")
	}
}

func TestWorktreeStart_FlagOverridesConfigSourceBranch(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create "develop" branch with develop.txt
	if err := exec.Command("git", "-C", dir, "checkout", "-b", "develop").Run(); err != nil {
		t.Fatalf("git checkout -b develop: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "develop.txt"), []byte("develop content"), 0o644); err != nil {
		t.Fatalf("write develop.txt: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "develop commit").Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}

	// Create "v3" branch with v3.txt
	if err := exec.Command("git", "-C", dir, "checkout", "-b", "v3").Run(); err != nil {
		t.Fatalf("git checkout -b v3: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "v3.txt"), []byte("v3 content"), 0o644); err != nil {
		t.Fatalf("write v3.txt: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "v3 commit").Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}

	if err := exec.Command("git", "-C", dir, "checkout", "master").Run(); err != nil {
		t.Fatalf("git checkout master: %v", err)
	}

	// Config says "develop"
	forgeDir := filepath.Join(dir, ".forge")
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		t.Fatalf("mkdir .forge: %v", err)
	}
	configContent := "worktree:\n  source-branch: develop\n"
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
		t.Fatalf("write config.yaml: %v", err)
	}

	slug := "flag-override-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", "develop").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", "v3").Run()
	})

	// Flag says "v3" — should override config
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug, "--source-branch", "v3"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created from v3 (has v3.txt, not develop.txt)
	if _, err := os.Stat(filepath.Join(targetDir, "v3.txt")); os.IsNotExist(err) {
		t.Errorf("worktree should have v3.txt (created from v3 branch via flag override)")
	}
}

func TestWorktreeStart_SourceBranchNotUsedForExistingBranch(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a "develop" branch
	if err := exec.Command("git", "-C", dir, "checkout", "-b", "develop").Run(); err != nil {
		t.Fatalf("git checkout -b develop: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "develop.txt"), []byte("develop content"), 0o644); err != nil {
		t.Fatalf("write develop.txt: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "develop commit").Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}

	// Create the "existing-slug" branch from develop (NOT from master)
	slug := "existing-slug"
	if err := exec.Command("git", "-C", dir, "branch", slug).Run(); err != nil {
		t.Fatalf("git branch %s: %v", slug, err)
	}
	if err := exec.Command("git", "-C", dir, "checkout", "master").Run(); err != nil {
		t.Fatalf("git checkout master: %v", err)
	}

	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", "develop").Run()
	})

	// Start with --source-branch develop, but branch already exists
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug, "--source-branch", "develop"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Worktree should still be created (existing branch path ignores source-branch)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Errorf("worktree directory %s should exist", targetDir)
	}
}

func TestWorktreeStart_NoSourceBranch_DefaultsToHEAD(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "default-head-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created from HEAD (has README.md from initial commit)
	if _, err := os.Stat(filepath.Join(targetDir, "README.md")); os.IsNotExist(err) {
		t.Errorf("worktree should have README.md (created from HEAD)")
	}
}

func TestWorktreeStart_SourceBranchFlagRegistered(t *testing.T) {
	flag := worktreeStartCmd.Flags().Lookup("source-branch")
	if flag == nil {
		t.Fatal("worktree start command should have --source-branch flag")
	}
	if flag.Shorthand != "b" {
		t.Errorf("source-branch shorthand should be 'b', got %q", flag.Shorthand)
	}
}

func TestWorktreeStart_ConfigSourceBranchErrorWhenBranchNotFound(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create .forge/config.yaml with a nonexistent source branch
	forgeDir := filepath.Join(dir, ".forge")
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		t.Fatalf("mkdir .forge: %v", err)
	}
	configContent := "worktree:\n  source-branch: nonexistent-config-branch\n"
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
		t.Fatalf("write config.yaml: %v", err)
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", "test-slug"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when config source branch does not exist")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "nonexistent-config-branch") {
		t.Errorf("error should mention the branch name, got: %s", stderr)
	}
}

// ---------------------------------------------------------------------------
// worktree start: GetWorktreeName auto-detects feature
// ---------------------------------------------------------------------------

func TestWorktreeStart_WorktreeNameAutoDetection(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "auto-detect-feature"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify GetWorktreeName returns the slug from the created worktree
	detectedName := gitPkg.GetWorktreeName(targetDir)
	if detectedName != slug {
		t.Errorf("GetWorktreeName() = %q, want %q", detectedName, slug)
	}
}

// ---------------------------------------------------------------------------
// worktree start: uses filepath.Join for path construction
// ---------------------------------------------------------------------------

func TestWorktreeStart_CreatesWorktreeInsideDotForgeWorktrees(t *testing.T) {
	resetSourceBranchFlag(t)

	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "inside-dot-forge"
	expectedDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", expectedDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created at .forge/worktrees/<slug>
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("worktree should be at %s, not as a sibling of project root", expectedDir)
	}

	// Verify it's NOT at the old sibling location
	oldPath := filepath.Join(filepath.Dir(dir), slug)
	if _, err := os.Stat(oldPath); err == nil {
		t.Errorf("worktree should NOT exist at old sibling path %s", oldPath)
		_ = exec.Command("git", "worktree", "remove", oldPath, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	}
}

func TestWorktreeStart_SlugValidation(t *testing.T) {
	// Empty slug should fail
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", ""})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when slug is empty")
	}
}

// ---------------------------------------------------------------------------
// worktree list: subcommand registration
// ---------------------------------------------------------------------------

func TestWorktreeCmd_HasListSubcommand(t *testing.T) {
	subcommands := worktreeCmd.Commands()
	found := false
	for _, cmd := range subcommands {
		if cmd.Name() == "list" {
			found = true
			break
		}
	}
	if !found {
		t.Error("worktree group should have 'list' subcommand")
	}
}

// ---------------------------------------------------------------------------
// worktree list: displays worktrees
// ---------------------------------------------------------------------------

func TestWorktreeList_ShowsMainWorktree(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "list"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "[main]") {
		t.Errorf("output should mark main worktree, got:\n%s", output)
	}
}

func TestWorktreeList_ShowsMultipleWorktrees(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create an additional worktree
	slug := "list-test-feature"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cmd := exec.Command("git", "worktree", "add", "-b", slug, targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "list"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, slug) {
		t.Errorf("output should contain worktree name %q, got:\n%s", slug, output)
	}
	if !strings.Contains(output, "[main]") {
		t.Errorf("output should mark main worktree, got:\n%s", output)
	}
}

func TestWorktreeList_MarksForgeManaged(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a feature directory matching a worktree name
	slug := "forge-managed-feat"
	featureDir := filepath.Join(dir, feature.FeaturesDir, slug, feature.TasksDirName)
	if err := os.MkdirAll(featureDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Create index.json so feature scanner recognizes it
	if err := os.WriteFile(filepath.Join(featureDir, feature.IndexFileName), []byte("{}"), 0o644); err != nil {
		t.Fatalf("write index.json: %v", err)
	}

	// Create the worktree
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cmd := exec.Command("git", "worktree", "add", "-b", slug, targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "list"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "[forge]") {
		t.Errorf("output should mark forge-managed worktree, got:\n%s", output)
	}
}

func TestWorktreeList_NoWorktrees(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Override ListWorktrees to return empty list
	origListWorktrees := listWorktreesFunc
	listWorktreesFunc = func(_ string) ([]gitPkg.WorktreeEntry, error) {
		return nil, nil
	}
	defer func() { listWorktreesFunc = origListWorktrees }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "list"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	if output != "No worktrees found" {
		t.Errorf("expected 'No worktrees found', got:\n%s", output)
	}
}

func TestWorktreeList_NotGitRepo(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "list"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when not a git repo")
	}
}

// ---------------------------------------------------------------------------
// worktree remove: subcommand registration
// ---------------------------------------------------------------------------

func TestWorktreeCmd_HasRemoveSubcommand(t *testing.T) {
	subcommands := worktreeCmd.Commands()
	found := false
	for _, cmd := range subcommands {
		if cmd.Name() == "remove" {
			found = true
			break
		}
	}
	if !found {
		t.Error("worktree group should have 'remove' subcommand")
	}
}

func TestWorktreeRemoveCmd_RequiresSlugArg(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "remove"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when slug argument is missing")
	}
}

// ---------------------------------------------------------------------------
// worktree remove: happy path
// ---------------------------------------------------------------------------

func TestWorktreeRemove_ResolvesDotForgeWorktreesPath(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create worktree at .forge/worktrees/<slug>
	slug := "remove-dot-forge-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		t.Fatalf("mkdir .forge/worktrees: %v", err)
	}
	cmd := exec.Command("git", "worktree", "add", "-b", slug, targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "remove", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Worktree directory should no longer exist
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		t.Errorf("worktree directory %s should have been removed", targetDir)
	}

	// Branch should still exist
	output, err := exec.Command("git", "-C", dir, "branch", "--list", slug).Output()
	if err != nil {
		t.Fatalf("git branch --list: %v", err)
	}
	if !strings.Contains(string(output), slug) {
		t.Errorf("branch %q should still exist after worktree removal", slug)
	}
}

func TestWorktreeRemove_RemovesWorktreeAndKeepsBranch(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a worktree
	slug := "remove-test-feature"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cmd := exec.Command("git", "worktree", "add", "-b", slug, targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "remove", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Worktree directory should no longer exist
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		t.Errorf("worktree directory %s should have been removed", targetDir)
	}

	// Branch should still exist
	output, err := exec.Command("git", "-C", dir, "branch", "--list", slug).Output()
	if err != nil {
		t.Fatalf("git branch --list: %v", err)
	}
	if !strings.Contains(string(output), slug) {
		t.Errorf("branch %q should still exist after worktree removal", slug)
	}

	// Confirmation output should include branch name
	outputStr := buf.String()
	if !strings.Contains(outputStr, slug) {
		t.Errorf("output should mention branch name %q, got:\n%s", slug, outputStr)
	}
}

// ---------------------------------------------------------------------------
// worktree remove: error cases
// ---------------------------------------------------------------------------

func TestWorktreeRemove_ErrorWhenWorktreeNotFound(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "remove", "nonexistent-worktree"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when worktree does not exist")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "not found") {
		t.Errorf("error should mention 'not found', got: %s", stderr)
	}
}

func TestWorktreeRemove_ErrorWhenUncommittedChanges(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a worktree
	slug := "dirty-worktree"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cmd := exec.Command("git", "worktree", "add", "-b", slug, targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	// Create uncommitted changes in the worktree
	dirtyFile := filepath.Join(targetDir, "dirty.txt")
	if err := os.WriteFile(dirtyFile, []byte("uncommitted"), 0o644); err != nil {
		t.Fatalf("write dirty file: %v", err)
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "remove", slug})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when worktree has uncommitted changes")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "commit") && !strings.Contains(stderr, "stash") {
		t.Errorf("error should hint to commit or stash, got: %s", stderr)
	}
}

func TestWorktreeRemove_ErrorWhenNotGitRepo(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "remove", "test-slug"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when not in a git repository")
	}
}

// ---------------------------------------------------------------------------
// worktree remove --hard: hard removal tests
// ---------------------------------------------------------------------------

// resetHardFlag resets the --hard flag on worktreeRemoveCmd.
func resetHardFlag(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		f := worktreeRemoveCmd.Flags().Lookup("hard")
		if f != nil {
			f.Changed = false
			_ = f.Value.Set("false")
		}
	})
}

// resetForceFlag resets the --force flag on worktreeRemoveCmd.
func resetForceFlag(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		f := worktreeRemoveCmd.Flags().Lookup("force")
		if f != nil {
			f.Changed = false
			_ = f.Value.Set("false")
		}
	})
}

func TestWorktreeRemove_Hard_RemovesWorktreeAndBranch(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create worktree
	slug := "hard-remove-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cmd := exec.Command("git", "worktree", "add", "-b", slug, targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	resetHardFlag(t)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "remove", slug, "--hard"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Worktree directory should no longer exist
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		t.Errorf("worktree directory %s should have been removed", targetDir)
	}

	// Branch should also be deleted
	output, err := exec.Command("git", "-C", dir, "branch", "--list", slug).Output()
	if err != nil {
		t.Fatalf("git branch --list: %v", err)
	}
	if strings.Contains(string(output), slug) {
		t.Errorf("branch %q should have been deleted with --hard", slug)
	}

	// Output should report all three steps
	outputStr := buf.String()
	if !strings.Contains(outputStr, "Removed worktree") {
		t.Errorf("output should mention worktree removal, got:\n%s", outputStr)
	}
	if !strings.Contains(outputStr, "Deleted branch") {
		t.Errorf("output should mention branch deletion, got:\n%s", outputStr)
	}
	if !strings.Contains(outputStr, "Pruned") {
		t.Errorf("output should mention pruning, got:\n%s", outputStr)
	}
}

func TestWorktreeRemove_Hard_UnmergedBranch_WarnsButProceeds(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a worktree and make a commit on the branch (so it's not merged into main)
	slug := "unmerged-branch"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cmd := exec.Command("git", "worktree", "add", "-b", slug, targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}

	// Make a commit on the branch that is NOT on main
	unmergedFile := filepath.Join(targetDir, "uncommitted.txt")
	if err := os.WriteFile(unmergedFile, []byte("unique content"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = targetDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	cmd = exec.Command("git", "commit", "-m", "unmerged commit")
	cmd.Dir = targetDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}

	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	resetHardFlag(t)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "remove", slug, "--hard"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Branch should be deleted (unmerged but --hard without --force still allows it per Hard Rules)
	output, err := exec.Command("git", "-C", dir, "branch", "--list", slug).Output()
	if err != nil {
		t.Fatalf("git branch --list: %v", err)
	}
	if strings.Contains(string(output), slug) {
		t.Errorf("branch %q should have been deleted", slug)
	}

	// Should warn about unmerged branch
	stderr := buf.String()
	if !strings.Contains(stderr, "not fully merged") && !strings.Contains(stderr, "unmerged") {
		t.Errorf("should warn about unmerged branch, got:\n%s", stderr)
	}
}

func TestWorktreeRemove_Hard_UncommittedChanges_Blocks(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "hard-dirty"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cmd := exec.Command("git", "worktree", "add", "-b", slug, targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	// Create uncommitted changes
	dirtyFile := filepath.Join(targetDir, "dirty.txt")
	if err := os.WriteFile(dirtyFile, []byte("uncommitted"), 0o644); err != nil {
		t.Fatalf("write dirty file: %v", err)
	}

	resetHardFlag(t)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "remove", slug, "--hard"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when worktree has uncommitted changes with --hard (no --force)")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "uncommitted changes") {
		t.Errorf("error should mention uncommitted changes, got: %s", stderr)
	}

	// Branch should still exist (removal was blocked)
	output, err := exec.Command("git", "-C", dir, "branch", "--list", slug).Output()
	if err != nil {
		t.Fatalf("git branch --list: %v", err)
	}
	if !strings.Contains(string(output), slug) {
		t.Errorf("branch %q should still exist when --hard is blocked by uncommitted changes", slug)
	}
}

func TestWorktreeRemove_HardForce_UncommittedChanges_Proceeds(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "hard-force-dirty"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cmd := exec.Command("git", "worktree", "add", "-b", slug, targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	// Create uncommitted changes
	dirtyFile := filepath.Join(targetDir, "dirty.txt")
	if err := os.WriteFile(dirtyFile, []byte("uncommitted"), 0o644); err != nil {
		t.Fatalf("write dirty file: %v", err)
	}

	resetHardFlag(t)
	resetForceFlag(t)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "remove", slug, "--hard", "--force"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Worktree should be removed
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		t.Errorf("worktree directory %s should have been removed", targetDir)
	}

	// Branch should be deleted
	output, err := exec.Command("git", "-C", dir, "branch", "--list", slug).Output()
	if err != nil {
		t.Fatalf("git branch --list: %v", err)
	}
	if strings.Contains(string(output), slug) {
		t.Errorf("branch %q should have been deleted with --hard --force", slug)
	}
}

func TestWorktreeRemove_WithoutHard_BehaviorUnchanged(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "no-hard-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cmd := exec.Command("git", "worktree", "add", "-b", slug, targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "remove", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Worktree removed, branch preserved
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		t.Errorf("worktree directory %s should have been removed", targetDir)
	}
	output, err := exec.Command("git", "-C", dir, "branch", "--list", slug).Output()
	if err != nil {
		t.Fatalf("git branch --list: %v", err)
	}
	if !strings.Contains(string(output), slug) {
		t.Errorf("branch %q should still exist without --hard", slug)
	}

	// Output should say "preserved"
	outputStr := buf.String()
	if !strings.Contains(outputStr, "preserved") {
		t.Errorf("output should say 'preserved', got:\n%s", outputStr)
	}
}

func TestWorktreeRemove_Hard_ReportsSkippedBranchIfAlreadyGone(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "hard-branch-gone"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cmd := exec.Command("git", "worktree", "add", "-b", slug, targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}

	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	// Override listWorktreesFunc to return a non-existent branch name.
	// This simulates the case where the branch name cannot be resolved.
	origListWorktrees := listWorktreesFunc
	listWorktreesFunc = func(_ string) ([]gitPkg.WorktreeEntry, error) {
		return []gitPkg.WorktreeEntry{
			{Path: targetDir, Branch: "nonexistent-branch-xyz"},
		}, nil
	}
	defer func() { listWorktreesFunc = origListWorktrees }()

	resetHardFlag(t)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "remove", slug, "--hard"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outputStr := buf.String()
	if !strings.Contains(outputStr, "Skipped branch deletion") {
		t.Errorf("should report skipped branch deletion when branch unknown, got:\n%s", outputStr)
	}
}

// ---------------------------------------------------------------------------
// listForgeFeatures helper
// ---------------------------------------------------------------------------

func TestListForgeFeatures(t *testing.T) {
	dir := t.TempDir()

	// No features dir — should return nil or empty
	f := listForgeFeatures(dir)
	if len(f) > 0 {
		t.Errorf("expected empty for missing features dir, got %d", len(f))
	}

	// Create features
	for _, slug := range []string{"feat-a", "feat-b", "feat-c"} {
		if err := os.MkdirAll(filepath.Join(dir, feature.FeaturesDir, slug), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	f = listForgeFeatures(dir)
	if len(f) != 3 {
		t.Errorf("expected 3 features, got %d", len(f))
	}
	if !f["feat-a"] || !f["feat-b"] || !f["feat-c"] {
		t.Errorf("expected all three features, got: %v", f)
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// resetSourceBranchFlag resets the --source-branch flag on worktreeStartCmd to
// prevent state leakage between tests. Cobra flags persist across Execute calls
// on the same Command instance.
func resetSourceBranchFlag(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		f := worktreeStartCmd.Flags().Lookup("source-branch")
		if f != nil {
			f.Changed = false
			_ = f.Value.Set("")
		}
	})
}

// resetNoLaunchFlag resets the --no-launch flag on worktreeStartCmd to
// prevent state leakage between tests.
func resetNoLaunchFlag(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		f := worktreeStartCmd.Flags().Lookup("no-launch")
		if f != nil {
			f.Changed = false
			_ = f.Value.Set("false")
		}
	})
}

// initGitRepoForWorktree creates a git repo with initial commit for worktree testing.
func initGitRepoForWorktree(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// git init
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git init: %v", err)
	}

	// Configure user
	cmd = exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git config email: %v", err)
	}
	cmd = exec.Command("git", "config", "user.name", "Test")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git config name: %v", err)
	}

	// Initial commit
	f := filepath.Join(dir, "README.md")
	if err := os.WriteFile(f, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}

	return dir
}

// ---------------------------------------------------------------------------
// validateCopyFilePath: path validation
// ---------------------------------------------------------------------------

func TestValidateCopyFilePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{name: "simple relative path", path: ".env", wantErr: false},
		{name: "nested relative path", path: "config/.env", wantErr: false},
		{name: "deep nested relative path", path: "a/b/c/.env", wantErr: false},
		{name: "Windows absolute path rejected", path: "C:\\Windows\\System32", wantErr: true},
		{name: "dot-dot traversal rejected", path: "../../etc/passwd", wantErr: true},
		{name: "dot-dot in middle rejected", path: "foo/../../etc/passwd", wantErr: true},
		{name: "dot-dot at start rejected", path: "../secret", wantErr: true},
		{name: "empty path allowed", path: "", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCopyFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCopyFilePath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// copyFilesToWorktree: file copy logic
// ---------------------------------------------------------------------------

func TestCopyFilesToWorktree_CopiesSingleFile(t *testing.T) {
	projectRoot := t.TempDir()
	worktreeDir := t.TempDir()

	// Create source file in project root
	if err := os.WriteFile(filepath.Join(projectRoot, ".env"), []byte("KEY=VALUE"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	err := copyFilesToWorktree(projectRoot, worktreeDir, []string{".env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file was copied
	data, err := os.ReadFile(filepath.Join(worktreeDir, ".env"))
	if err != nil {
		t.Fatalf("read copied file: %v", err)
	}
	if string(data) != "KEY=VALUE" {
		t.Errorf("copied content = %q, want %q", string(data), "KEY=VALUE")
	}
}

func TestCopyFilesToWorktree_CopiesMultipleFiles(t *testing.T) {
	projectRoot := t.TempDir()
	worktreeDir := t.TempDir()

	// Create source files
	if err := os.WriteFile(filepath.Join(projectRoot, ".env"), []byte("ENV=dev"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectRoot, ".env.local"), []byte("LOCAL=true"), 0o644); err != nil {
		t.Fatalf("write .env.local: %v", err)
	}

	err := copyFilesToWorktree(projectRoot, worktreeDir, []string{".env", ".env.local"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, file := range []string{".env", ".env.local"} {
		if _, err := os.Stat(filepath.Join(worktreeDir, file)); os.IsNotExist(err) {
			t.Errorf("file %s should exist in worktree", file)
		}
	}
}

func TestCopyFilesToWorktree_OverwritesExistingFile(t *testing.T) {
	projectRoot := t.TempDir()
	worktreeDir := t.TempDir()

	// Create source file
	if err := os.WriteFile(filepath.Join(projectRoot, ".env"), []byte("NEW=content"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	// Create existing file in worktree (simulates git checkout having it)
	if err := os.WriteFile(filepath.Join(worktreeDir, ".env"), []byte("OLD=content"), 0o644); err != nil {
		t.Fatalf("write old .env: %v", err)
	}

	err := copyFilesToWorktree(projectRoot, worktreeDir, []string{".env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(worktreeDir, ".env"))
	if err != nil {
		t.Fatalf("read .env: %v", err)
	}
	if string(data) != "NEW=content" {
		t.Errorf("should overwrite with project root version, got %q", string(data))
	}
}

func TestCopyFilesToWorktree_CopiesNestedFile(t *testing.T) {
	projectRoot := t.TempDir()
	worktreeDir := t.TempDir()

	// Create nested source file
	nestedDir := filepath.Join(projectRoot, "config")
	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("mkdir config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nestedDir, "app.conf"), []byte("port=8080"), 0o644); err != nil {
		t.Fatalf("write app.conf: %v", err)
	}

	err := copyFilesToWorktree(projectRoot, worktreeDir, []string{"config/app.conf"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(worktreeDir, "config", "app.conf"))
	if err != nil {
		t.Fatalf("read copied file: %v", err)
	}
	if string(data) != "port=8080" {
		t.Errorf("copied content = %q, want %q", string(data), "port=8080")
	}
}

func TestCopyFilesToWorktree_ErrorOnInvalidPath(t *testing.T) {
	projectRoot := t.TempDir()
	worktreeDir := t.TempDir()

	err := copyFilesToWorktree(projectRoot, worktreeDir, []string{"../../etc/passwd"})
	if err == nil {
		t.Error("expected error for path traversal")
	}
}

// ---------------------------------------------------------------------------
// validateCopyFiles: pre-validation of all copy-files
// ---------------------------------------------------------------------------

func TestValidateCopyFiles_AllFilesExist(t *testing.T) {
	dir := t.TempDir()

	// Create files
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("KEY=VAL"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".env.local"), []byte("LOCAL=true"), 0o644); err != nil {
		t.Fatalf("write .env.local: %v", err)
	}

	err := validateCopyFiles(dir, []string{".env", ".env.local"})
	if err != nil {
		t.Errorf("expected no error when all files exist, got: %v", err)
	}
}

func TestValidateCopyFiles_MissingFile(t *testing.T) {
	dir := t.TempDir()

	// Create only one of two files
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("KEY=VAL"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	err := validateCopyFiles(dir, []string{".env", ".env.missing"})
	if err == nil {
		t.Error("expected error when a copy-file is missing")
	}
	if !strings.Contains(err.Error(), ".env.missing") {
		t.Errorf("error should mention the missing file, got: %v", err)
	}
}

func TestValidateCopyFiles_InvalidPathRejected(t *testing.T) {
	dir := t.TempDir()

	err := validateCopyFiles(dir, []string{"/etc/passwd"})
	if err == nil {
		t.Error("expected error for absolute path")
	}
}

func TestValidateCopyFiles_EmptyList(t *testing.T) {
	dir := t.TempDir()

	err := validateCopyFiles(dir, nil)
	if err != nil {
		t.Errorf("expected no error for empty list, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// worktree start: copy-files integration
// ---------------------------------------------------------------------------

func TestWorktreeStart_CopyFilesFromConfig(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create .env in project root (not committed, not gitignored -- just exists)
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("DB_HOST=localhost"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	// Create .forge/config.yaml with copy-files
	forgeDir := filepath.Join(dir, ".forge")
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		t.Fatalf("mkdir .forge: %v", err)
	}
	configContent := "worktree:\n  copy-files:\n    - .env\n"
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
		t.Fatalf("write config.yaml: %v", err)
	}

	slug := "copy-files-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify .env was copied to worktree
	data, err := os.ReadFile(filepath.Join(targetDir, ".env"))
	if err != nil {
		t.Fatalf("read .env from worktree: %v", err)
	}
	if string(data) != "DB_HOST=localhost" {
		t.Errorf("copied .env content = %q, want %q", string(data), "DB_HOST=localhost")
	}
}

func TestWorktreeStart_AbortsWhenCopyFileMissing(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// .env does NOT exist in project root

	// Create .forge/config.yaml with copy-files
	forgeDir := filepath.Join(dir, ".forge")
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		t.Fatalf("mkdir .forge: %v", err)
	}
	configContent := "worktree:\n  copy-files:\n    - .env\n"
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
		t.Fatalf("write config.yaml: %v", err)
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", "test-slug"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when copy-file is missing from project root")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, ".env") {
		t.Errorf("error should mention the missing file, got: %s", stderr)
	}

	// Verify NO worktree was created (pre-validation)
	slug := "test-slug"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		// Clean up any orphan worktree
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
		t.Error("worktree should NOT have been created when copy-file is missing")
	}
}

func TestWorktreeStart_NoCopyWhenConfigAbsent(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// No .forge/config.yaml at all

	slug := "no-copy-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created normally
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Errorf("worktree directory %s should exist", targetDir)
	}
}

func TestWorktreeStart_NoCopyWhenCopyFilesEmpty(t *testing.T) {
	resetSourceBranchFlag(t)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create .forge/config.yaml with empty copy-files
	forgeDir := filepath.Join(dir, ".forge")
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		t.Fatalf("mkdir .forge: %v", err)
	}
	configContent := "worktree:\n  copy-files: []\n"
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
		t.Fatalf("write config.yaml: %v", err)
	}

	slug := "empty-copy-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created normally
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Errorf("worktree directory %s should exist", targetDir)
	}
}

// ---------------------------------------------------------------------------
// worktree resume: subcommand registration
// ---------------------------------------------------------------------------

func TestWorktreeCmd_HasResumeSubcommand(t *testing.T) {
	subcommands := worktreeCmd.Commands()
	found := false
	for _, cmd := range subcommands {
		if cmd.Name() == "resume" {
			found = true
			break
		}
	}
	if !found {
		t.Error("worktree group should have 'resume' subcommand")
	}
}

func TestWorktreeResumeCmd_RequiresSlugArg(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "resume"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when slug argument is missing")
	}
}

// ---------------------------------------------------------------------------
// worktree resume: pre-flight claude check
// ---------------------------------------------------------------------------

func TestWorktreeResume_ErrorWhenClaudeNotInPath(t *testing.T) {
	origLookPath := lookPathFunc
	lookPathFunc = func(_ string) (string, error) {
		return "", &exec.Error{Name: "claude", Err: exec.ErrNotFound}
	}
	defer func() { lookPathFunc = origLookPath }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "resume", "test-slug"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when claude binary not found")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "claude") {
		t.Errorf("error should mention 'claude', got: %s", stderr)
	}
}

// ---------------------------------------------------------------------------
// worktree resume: error when worktree not found
// ---------------------------------------------------------------------------

func TestWorktreeResume_ErrorWhenWorktreeNotFound(t *testing.T) {
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "resume", "nonexistent-worktree"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when worktree does not exist")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "not found") {
		t.Errorf("error should mention 'not found', got: %s", stderr)
	}
}

// ---------------------------------------------------------------------------
// worktree resume: error when not a git repo
// ---------------------------------------------------------------------------

func TestWorktreeResume_ErrorWhenNotGitRepo(t *testing.T) {
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "resume", "test-slug"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when not in a git repository")
	}
}

// ---------------------------------------------------------------------------
// worktree resume: error when directory exists but not a git worktree
// ---------------------------------------------------------------------------

func TestWorktreeResume_ErrorWhenDirExistsButNotWorktree(t *testing.T) {
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a sibling directory that is NOT a git worktree
	slug := "not-a-worktree"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(targetDir) })

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "resume", slug})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when directory is not a git worktree")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "not a git worktree") {
		t.Errorf("error should mention 'not a git worktree', got: %s", stderr)
	}
}

// ---------------------------------------------------------------------------
// worktree resume: happy path
// ---------------------------------------------------------------------------

func TestWorktreeResume_LaunchesClaudeInExistingWorktree(t *testing.T) {
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	// Simulate -c not supported to test basic launch behavior
	origResumeSupport := claudeSupportsContinueFlagFunc
	claudeSupportsContinueFlagFunc = func() bool { return false }
	defer func() { claudeSupportsContinueFlagFunc = origResumeSupport }()

	// Capture claude launch args and working directory
	var capturedArgs []string
	var capturedWd string
	origRunClaude := runClaudeFunc
	runClaudeFunc = func(args []string) error {
		capturedArgs = args
		wd, _ := os.Getwd()
		capturedWd = wd
		return nil
	}
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a real worktree
	slug := "resume-feature"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cmd := exec.Command("git", "worktree", "add", "-b", slug, targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "resume", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify claude was launched with --dangerously-skip-permissions
	if len(capturedArgs) == 0 {
		t.Fatal("claude should have been launched")
	}
	if capturedArgs[0] != "--dangerously-skip-permissions" {
		t.Errorf("first arg should be --dangerously-skip-permissions, got %q", capturedArgs[0])
	}

	// Verify claude was launched in the worktree directory
	absTarget, _ := filepath.Abs(targetDir)
	resolvedTarget, _ := filepath.EvalSymlinks(absTarget)
	if capturedWd != resolvedTarget {
		t.Errorf("claude should have been launched in %s, got %s", resolvedTarget, capturedWd)
	}
}

// ---------------------------------------------------------------------------
// worktree resume: -c session restore
// ---------------------------------------------------------------------------

func TestWorktreeResume_UsesContinueFlagWhenSupported(t *testing.T) {
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	// Simulate claude -c support
	origResumeSupport := claudeSupportsContinueFlagFunc
	claudeSupportsContinueFlagFunc = func() bool { return true }
	defer func() { claudeSupportsContinueFlagFunc = origResumeSupport }()

	// Capture claude launch args
	var capturedArgs []string
	origRunClaude := runClaudeFunc
	runClaudeFunc = func(args []string) error {
		capturedArgs = args
		return nil
	}
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a real worktree
	slug := "resume-with-continue"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cmd := exec.Command("git", "worktree", "add", "-b", slug, targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "resume", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify claude was launched with -c <slug> --dangerously-skip-permissions
	if len(capturedArgs) < 3 {
		t.Fatalf("expected at least 3 args, got %v", capturedArgs)
	}
	if capturedArgs[0] != "-c" {
		t.Errorf("first arg should be '-c', got %q", capturedArgs[0])
	}
	if capturedArgs[1] != slug {
		t.Errorf("second arg should be slug %q, got %q", slug, capturedArgs[1])
	}
	if capturedArgs[2] != "--dangerously-skip-permissions" {
		t.Errorf("third arg should be '--dangerously-skip-permissions', got %q", capturedArgs[2])
	}
}

func TestWorktreeResume_FallsBackWhenContinueNotSupported(t *testing.T) {
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	// Simulate claude -c NOT supported
	origResumeSupport := claudeSupportsContinueFlagFunc
	claudeSupportsContinueFlagFunc = func() bool { return false }
	defer func() { claudeSupportsContinueFlagFunc = origResumeSupport }()

	// Capture claude launch args
	var capturedArgs []string
	origRunClaude := runClaudeFunc
	runClaudeFunc = func(args []string) error {
		capturedArgs = args
		return nil
	}
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a real worktree
	slug := "resume-no-continue"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cmd := exec.Command("git", "worktree", "add", "-b", slug, targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "resume", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify claude was launched WITHOUT -c, just --dangerously-skip-permissions
	if len(capturedArgs) != 1 {
		t.Fatalf("expected 1 arg, got %v", capturedArgs)
	}
	if capturedArgs[0] != "--dangerously-skip-permissions" {
		t.Errorf("arg should be '--dangerously-skip-permissions', got %q", capturedArgs[0])
	}
}

// ---------------------------------------------------------------------------
// worktree start: table-driven remote branch detection (mocked gitRunFunc)
// ---------------------------------------------------------------------------

// mockGitResponse is a simplified mock response for git commands.
type mockGitResponse struct {
	output string
	err    error
}

// ---------------------------------------------------------------------------
// worktree start: --no-launch flag
// ---------------------------------------------------------------------------

func TestWorktreeStart_NoLaunchFlagRegistered(t *testing.T) {
	flag := worktreeStartCmd.Flags().Lookup("no-launch")
	if flag == nil {
		t.Fatal("worktree start command should have --no-launch flag")
	}
	if flag.DefValue != "false" {
		t.Errorf("no-launch default should be false, got %q", flag.DefValue)
	}
}

func TestWorktreeStart_NoLaunch_CreatesWorktreeWithoutLaunchingClaude(t *testing.T) {
	resetSourceBranchFlag(t)
	resetNoLaunchFlag(t)

	// Make claude available (but should NOT be called)
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	// Track whether claude was called
	claudeCalled := false
	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error {
		claudeCalled = true
		return nil
	}
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "no-launch-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug, "--no-launch"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Errorf("worktree directory %s should exist", targetDir)
	}

	// Verify claude was NOT called
	if claudeCalled {
		t.Error("claude should NOT have been launched with --no-launch")
	}

	// Verify stdout contains the worktree path
	stdout := buf.String()
	if !strings.Contains(stdout, targetDir) {
		t.Errorf("stdout should contain worktree path %q, got: %s", targetDir, stdout)
	}
}

func TestWorktreeStart_NoLaunch_WithSourceBranch(t *testing.T) {
	resetSourceBranchFlag(t)
	resetNoLaunchFlag(t)

	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a "develop" branch with a distinct file
	if err := exec.Command("git", "-C", dir, "checkout", "-b", "develop").Run(); err != nil {
		t.Fatalf("git checkout -b develop: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "develop.txt"), []byte("develop"), 0o644); err != nil {
		t.Fatalf("write develop.txt: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "develop commit").Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "checkout", "master").Run(); err != nil {
		t.Fatalf("git checkout master: %v", err)
	}

	slug := "no-launch-source"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", "develop").Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug, "--source-branch", "develop", "--no-launch"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created from develop (has develop.txt)
	if _, err := os.Stat(filepath.Join(targetDir, "develop.txt")); os.IsNotExist(err) {
		t.Errorf("worktree should have develop.txt (created from develop branch)")
	}
}

func TestWorktreeStart_WithoutNoLaunch_LaunchesClaude(t *testing.T) {
	resetSourceBranchFlag(t)
	resetNoLaunchFlag(t)

	// Make claude available
	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	// Capture claude launch
	claudeCalled := false
	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error {
		claudeCalled = true
		return nil
	}
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "with-launch-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify claude WAS called (default behavior)
	if !claudeCalled {
		t.Error("claude should have been launched without --no-launch (default behavior)")
	}
}

// ---------------------------------------------------------------------------
// worktree start: branch-first creation
// ---------------------------------------------------------------------------

func TestWorktreeStart_BranchFirstCreation_NewBranchFromHead(t *testing.T) {
	resetSourceBranchFlag(t)
	resetNoLaunchFlag(t)

	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "branch-first-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Errorf("worktree directory %s should exist", targetDir)
	}

	// Verify the branch exists independently (not just as part of worktree)
	output, err := exec.Command("git", "-C", dir, "branch", "--list", slug).Output()
	if err != nil {
		t.Fatalf("git branch --list: %v", err)
	}
	if !strings.Contains(string(output), slug) {
		t.Errorf("branch %q should exist", slug)
	}
}

func TestWorktreeStart_BranchFirstCreation_SkipsCheckoutWhenBranchExists(t *testing.T) {
	resetSourceBranchFlag(t)
	resetNoLaunchFlag(t)

	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "existing-branch-first"
	// Pre-create the branch
	if err := exec.Command("git", "-C", dir, "branch", slug).Run(); err != nil {
		t.Fatalf("git branch %s: %v", slug, err)
	}

	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created (existing branch should just be used)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Errorf("worktree directory %s should exist", targetDir)
	}
}

func TestWorktreeStart_BranchFirstCreation_WithSourceBranch(t *testing.T) {
	resetSourceBranchFlag(t)
	resetNoLaunchFlag(t)

	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create "staging" branch with unique file
	if err := exec.Command("git", "-C", dir, "checkout", "-b", "staging").Run(); err != nil {
		t.Fatalf("git checkout -b staging: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "staging.txt"), []byte("staging"), 0o644); err != nil {
		t.Fatalf("write staging.txt: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "staging commit").Run(); err != nil {
		t.Fatalf("git commit: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "checkout", "master").Run(); err != nil {
		t.Fatalf("git checkout master: %v", err)
	}

	slug := "branch-first-source"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", "staging").Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug, "--source-branch", "staging"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created from staging (has staging.txt)
	if _, err := os.Stat(filepath.Join(targetDir, "staging.txt")); os.IsNotExist(err) {
		t.Errorf("worktree should have staging.txt (created from staging branch)")
	}

	// Verify the branch was created from staging (not from HEAD)
	output, err := exec.Command("git", "-C", dir, "log", "--oneline", "-1", slug).Output()
	if err != nil {
		t.Fatalf("git log: %v", err)
	}
	if !strings.Contains(string(output), "staging commit") {
		t.Errorf("branch should be based on staging, got: %s", string(output))
	}
}

func TestWorktreeStart_BranchFirstCreation_CleansUpBranchOnWorktreeFailure(t *testing.T) {
	resetSourceBranchFlag(t)
	resetNoLaunchFlag(t)

	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "cleanup-test"

	// Make worktree add fail by mocking git
	var branchCreated bool
	origGitRun := gitRunFunc
	gitRunFunc = func(_ string, args ...string) (string, error) {
		// Track if branch was created via `git branch <slug>`
		if len(args) >= 2 && args[0] == "branch" && args[1] == slug {
			branchCreated = true
			return "", nil
		}
		// Make worktree add fail
		if len(args) >= 2 && args[0] == "worktree" && args[1] == "add" {
			return "", fmt.Errorf("worktree add failed")
		}
		// Allow other git commands
		return origGitRun(dir, args...)
	}
	defer func() { gitRunFunc = origGitRun }()

	// Also need to mock branch deletion cleanup
	var branchDeleted bool
	cleanupOrigGitRun := gitRunFunc
	gitRunFunc = func(root string, args ...string) (string, error) {
		if len(args) >= 3 && args[0] == "branch" && args[1] == "-D" && args[2] == slug {
			branchDeleted = true
			return "", nil
		}
		return cleanupOrigGitRun(root, args...)
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", slug})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when worktree add fails")
	}

	// Verify branch was created
	if !branchCreated {
		t.Error("branch should have been created before worktree add")
	}

	// Verify branch was cleaned up after failure
	if !branchDeleted {
		t.Error("branch should have been cleaned up after worktree add failure")
	}
}

// ---------------------------------------------------------------------------
// worktree start: table-driven remote branch detection (mocked gitRunFunc)
// ---------------------------------------------------------------------------

func TestWorktreeStart_RemoteBranchResolution(t *testing.T) {
	tests := []struct {
		name               string
		slug               string
		sourceBranch       string // --source-branch flag (empty = not set)
		mockResponses      map[string]mockGitResponse
		wantErr            bool
		wantStdoutContains string
		wantStderrContains string
		wantWorktreeArgs   []string // the git worktree add args prefix that should be used
	}{
		{
			name:         "remote branch exists, local does not, no source-branch",
			slug:         "remote-only-feature",
			sourceBranch: "",
			mockResponses: map[string]mockGitResponse{
				"rev-parse --verify remote-only-feature":                {err: fmt.Errorf("not found")},
				"fetch origin":                                          {},
				"rev-parse --verify remotes/origin/remote-only-feature": {},
				"branch remote-only-feature":                            {},
				"worktree add":                                          {},
			},
			wantErr:            false,
			wantStdoutContains: "origin/remote-only-feature",
			wantWorktreeArgs:   []string{"worktree", "add"},
		},
		{
			name:         "fetch fails, no remote branch, falls back to HEAD",
			slug:         "fetch-fail-feature",
			sourceBranch: "",
			mockResponses: map[string]mockGitResponse{
				"rev-parse --verify fetch-fail-feature":                {err: fmt.Errorf("not found")},
				"fetch origin":                                         {err: fmt.Errorf("network error")},
				"rev-parse --verify remotes/origin/fetch-fail-feature": {err: fmt.Errorf("not found")},
				"branch fetch-fail-feature":                            {},
				"worktree add":                                         {},
			},
			wantErr:            false,
			wantStderrContains: "warning: git fetch origin failed",
			wantWorktreeArgs:   []string{"worktree", "add"},
		},
		{
			name:         "both local and remote absent, creates from HEAD",
			slug:         "new-feature",
			sourceBranch: "",
			mockResponses: map[string]mockGitResponse{
				"rev-parse --verify new-feature":                {err: fmt.Errorf("not found")},
				"fetch origin":                                  {},
				"rev-parse --verify remotes/origin/new-feature": {err: fmt.Errorf("not found")},
				"branch new-feature":                            {},
				"worktree add":                                  {},
			},
			wantErr:          false,
			wantWorktreeArgs: []string{"worktree", "add"},
		},
		{
			name:         "source-branch set but remote branch exists, remote wins",
			slug:         "remote-wins-feature",
			sourceBranch: "develop",
			mockResponses: map[string]mockGitResponse{
				"rev-parse --verify remote-wins-feature":                {err: fmt.Errorf("not found")},
				"fetch origin":                                          {},
				"rev-parse --verify remotes/origin/remote-wins-feature": {},
				"branch remote-wins-feature":                            {},
				"worktree add":                                          {},
			},
			wantErr:            false,
			wantStdoutContains: "origin/remote-wins-feature",
			wantWorktreeArgs:   []string{"worktree", "add"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSourceBranchFlag(t)

			// Mock claude in PATH
			origLookPath := lookPathFunc
			lookPathFunc = func(name string) (string, error) {
				if name == "claude" {
					return "/usr/bin/claude", nil
				}
				return exec.LookPath(name)
			}
			defer func() { lookPathFunc = origLookPath }()

			// Mock claude launch
			origRunClaude := runClaudeFunc
			runClaudeFunc = func(_ []string) error { return nil }
			defer func() { runClaudeFunc = origRunClaude }()

			// Track which git worktree add args were used
			var capturedWorktreeArgs []string
			origGitRun := gitRunFunc
			gitRunFunc = func(_ string, args ...string) (string, error) {
				key := strings.Join(args, " ")

				// Helper: extract target dir from worktree add args and create it
				// Capture worktree add args specifically
				if len(args) >= 3 && args[0] == "worktree" && args[1] == "add" {
					capturedWorktreeArgs = args
				}
				maybeCreateTarget := func() {
					if len(args) >= 3 && args[0] == "worktree" && args[1] == "add" {
						// worktree add targetDir [slug] (branch-first: no -b flag)
						// The target dir is always args[2]
						_ = os.MkdirAll(args[2], 0o755)
					}
				}
				// Check for exact match first
				if resp, ok := tt.mockResponses[key]; ok {
					if resp.err == nil {
						maybeCreateTarget()
					}
					return resp.output, resp.err
				}

				// Check for prefix match (worktree add args include dynamic target path)
				for pattern, resp := range tt.mockResponses {
					if strings.HasPrefix(key, pattern) {
						if resp.err == nil {
							maybeCreateTarget()
						}
						return resp.output, resp.err
					}
				}

				// Default: no error for unmocked commands
				return "", nil
			}
			defer func() { gitRunFunc = origGitRun }()

			// Set up isolated git repo (real git init for IsGitRepository check)
			dir := initGitRepoForWorktree(t)
			t.Setenv("CLAUDE_PROJECT_DIR", dir)
			origWd, _ := os.Getwd()
			t.Cleanup(func() { _ = os.Chdir(origWd) })
			_ = os.Chdir(dir)

			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)

			cmdArgs := []string{"worktree", "start", tt.slug}
			if tt.sourceBranch != "" {
				cmdArgs = append(cmdArgs, "--source-branch", tt.sourceBranch)
			}
			rootCmd.SetArgs(cmdArgs)

			err := rootCmd.Execute()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check stdout
			if tt.wantStdoutContains != "" {
				stdout := buf.String()
				if !strings.Contains(stdout, tt.wantStdoutContains) {
					t.Errorf("stdout should contain %q, got: %s", tt.wantStdoutContains, stdout)
				}
			}

			// Check stderr
			if tt.wantStderrContains != "" {
				stderr := buf.String()
				if !strings.Contains(stderr, tt.wantStderrContains) {
					t.Errorf("stderr should contain %q, got: %s", tt.wantStderrContains, stderr)
				}
			}

			// Check worktree args prefix
			if tt.wantWorktreeArgs != nil {
				if capturedWorktreeArgs == nil {
					t.Fatal("expected worktree add command to be called, but it wasn't")
				}
				for i, want := range tt.wantWorktreeArgs {
					if i >= len(capturedWorktreeArgs) {
						t.Errorf("worktree args too short: want %q at index %d, got %v", want, i, capturedWorktreeArgs)
						break
					}
					if want != capturedWorktreeArgs[i] {
						t.Errorf("worktree args[%d] = %q, want %q, full args: %v", i, capturedWorktreeArgs[i], want, capturedWorktreeArgs)
					}
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// interactive mode: listUnfinishedItems
// ---------------------------------------------------------------------------

func TestListUnfinishedItems_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	items := listUnfinishedItems(dir)
	if len(items) != 0 {
		t.Errorf("expected no items for empty dir, got %d", len(items))
	}
}

func TestListUnfinishedItems_NoDir(t *testing.T) {
	dir := t.TempDir()
	// Neither docs/proposals/ nor docs/features/ exist
	items := listUnfinishedItems(dir)
	if len(items) != 0 {
		t.Errorf("expected no items when dirs don't exist, got %d", len(items))
	}
}

func TestListUnfinishedItems_SkipsCompletedProposal(t *testing.T) {
	dir := t.TempDir()

	// Create a completed proposal
	proposalDir := filepath.Join(dir, "docs", "proposals", "done-proposal")
	if err := os.MkdirAll(proposalDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\nstatus: completed\ncreated: 2026-01-01\n---\n# Done"
	if err := os.WriteFile(filepath.Join(proposalDir, "proposal.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	items := listUnfinishedItems(dir)
	if len(items) != 0 {
		t.Errorf("expected no items for completed proposal, got %d", len(items))
	}
}

func TestListUnfinishedItems_IncludesDraftProposal(t *testing.T) {
	dir := t.TempDir()

	// Create a draft proposal
	proposalDir := filepath.Join(dir, "docs", "proposals", "my-proposal")
	if err := os.MkdirAll(proposalDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\nstatus: Draft\ncreated: 2026-01-01\n---\n# Draft"
	if err := os.WriteFile(filepath.Join(proposalDir, "proposal.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	items := listUnfinishedItems(dir)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Slug != "my-proposal" {
		t.Errorf("expected slug 'my-proposal', got %q", items[0].Slug)
	}
	if items[0].Type != "proposal" {
		t.Errorf("expected type 'proposal', got %q", items[0].Type)
	}
	if items[0].Status != "Draft" {
		t.Errorf("expected status 'Draft', got %q", items[0].Status)
	}
}

func TestListUnfinishedItems_IncludesProposalWithNoStatus(t *testing.T) {
	dir := t.TempDir()

	// Create a proposal with no status field
	proposalDir := filepath.Join(dir, "docs", "proposals", "no-status-proposal")
	if err := os.MkdirAll(proposalDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\ncreated: 2026-01-01\n---\n# No Status"
	if err := os.WriteFile(filepath.Join(proposalDir, "proposal.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	items := listUnfinishedItems(dir)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Status != "Draft" {
		t.Errorf("expected status 'Draft' for proposal without status, got %q", items[0].Status)
	}
}

func TestListUnfinishedItems_IncludesFeatureWithoutManifest(t *testing.T) {
	dir := t.TempDir()

	// Create a feature directory without manifest.md
	featureDir := filepath.Join(dir, "docs", "features", "my-feature", "tasks")
	if err := os.MkdirAll(featureDir, 0o755); err != nil {
		t.Fatal(err)
	}

	items := listUnfinishedItems(dir)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Slug != "my-feature" {
		t.Errorf("expected slug 'my-feature', got %q", items[0].Slug)
	}
	if items[0].Type != "feature" {
		t.Errorf("expected type 'feature', got %q", items[0].Type)
	}
	if items[0].Status != "active" {
		t.Errorf("expected status 'active', got %q", items[0].Status)
	}
}

func TestListUnfinishedItems_SkipsCompletedFeature(t *testing.T) {
	dir := t.TempDir()

	// Create a feature directory with completed manifest
	featureDir := filepath.Join(dir, "docs", "features", "completed-feat")
	if err := os.MkdirAll(featureDir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := "---\nstatus: completed\n---\n# Completed"
	if err := os.WriteFile(filepath.Join(featureDir, "manifest.md"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	items := listUnfinishedItems(dir)
	if len(items) != 0 {
		t.Errorf("expected no items for completed feature, got %d", len(items))
	}
}

func TestListUnfinishedItems_IncludesFeatureWithStatus(t *testing.T) {
	dir := t.TempDir()

	// Create a feature directory with in_progress manifest
	featureDir := filepath.Join(dir, "docs", "features", "active-feat")
	if err := os.MkdirAll(featureDir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := "---\nstatus: in_progress\n---\n# Active"
	if err := os.WriteFile(filepath.Join(featureDir, "manifest.md"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	items := listUnfinishedItems(dir)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Status != "in_progress" {
		t.Errorf("expected status 'in_progress', got %q", items[0].Status)
	}
}

func TestListUnfinishedItems_DoesNotDuplicateProposalSlug(t *testing.T) {
	dir := t.TempDir()

	// Create both a proposal and a feature with the same slug
	proposalDir := filepath.Join(dir, "docs", "proposals", "shared-slug")
	if err := os.MkdirAll(proposalDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\nstatus: Draft\ncreated: 2026-01-01\n---\n# Proposal"
	if err := os.WriteFile(filepath.Join(proposalDir, "proposal.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	featureDir := filepath.Join(dir, "docs", "features", "shared-slug")
	if err := os.MkdirAll(featureDir, 0o755); err != nil {
		t.Fatal(err)
	}

	items := listUnfinishedItems(dir)
	if len(items) != 1 {
		t.Fatalf("expected 1 item (no duplicate), got %d", len(items))
	}
	if items[0].Type != "proposal" {
		t.Errorf("expected type 'proposal' (proposal takes priority), got %q", items[0].Type)
	}
}

func TestListUnfinishedItems_MixedProposalsAndFeatures(t *testing.T) {
	dir := t.TempDir()

	// Proposal 1: Draft
	p1Dir := filepath.Join(dir, "docs", "proposals", "draft-proposal")
	if err := os.MkdirAll(p1Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(p1Dir, "proposal.md"), []byte("---\nstatus: Draft\ncreated: 2026-01-01\n---\n# P1"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Proposal 2: completed (should be skipped)
	p2Dir := filepath.Join(dir, "docs", "proposals", "done-proposal")
	if err := os.MkdirAll(p2Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(p2Dir, "proposal.md"), []byte("---\nstatus: completed\ncreated: 2026-01-01\n---\n# P2"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Feature 1: active
	f1Dir := filepath.Join(dir, "docs", "features", "active-feat")
	if err := os.MkdirAll(f1Dir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Feature 2: completed (should be skipped)
	f2Dir := filepath.Join(dir, "docs", "features", "done-feat")
	if err := os.MkdirAll(f2Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f2Dir, "manifest.md"), []byte("---\nstatus: completed\n---\n# Done"), 0o644); err != nil {
		t.Fatal(err)
	}

	items := listUnfinishedItems(dir)
	if len(items) != 2 {
		t.Fatalf("expected 2 items (1 draft proposal + 1 active feature), got %d: %+v", len(items), items)
	}
}

// ---------------------------------------------------------------------------
// interactive mode: promptSelection
// ---------------------------------------------------------------------------

func TestPromptSelection_ValidNumber(t *testing.T) {
	items := []selectableItem{
		{Slug: "foo", Type: "proposal", Status: "Draft"},
		{Slug: "bar", Type: "feature", Status: "active"},
	}

	// Mock stdin to return "1"
	origStdin := stdinFunc
	stdinFunc = func() (string, error) { return "1", nil }
	defer func() { stdinFunc = origStdin }()

	var buf bytes.Buffer
	slug, err := promptSelection(items, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if slug != "foo" {
		t.Errorf("expected slug 'foo', got %q", slug)
	}

	// Verify output contains the numbered list
	output := buf.String()
	if !strings.Contains(output, "1.") || !strings.Contains(output, "foo") {
		t.Errorf("output should contain numbered list, got:\n%s", output)
	}
}

func TestPromptSelection_SecondItem(t *testing.T) {
	items := []selectableItem{
		{Slug: "foo", Type: "proposal", Status: "Draft"},
		{Slug: "bar", Type: "feature", Status: "active"},
	}

	origStdin := stdinFunc
	stdinFunc = func() (string, error) { return "2", nil }
	defer func() { stdinFunc = origStdin }()

	var buf bytes.Buffer
	slug, err := promptSelection(items, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if slug != "bar" {
		t.Errorf("expected slug 'bar', got %q", slug)
	}
}

func TestPromptSelection_InvalidNumber(t *testing.T) {
	items := []selectableItem{
		{Slug: "foo", Type: "proposal", Status: "Draft"},
	}

	origStdin := stdinFunc
	stdinFunc = func() (string, error) { return "5", nil }
	defer func() { stdinFunc = origStdin }()

	var buf bytes.Buffer
	_, err := promptSelection(items, &buf)
	if err == nil {
		t.Error("expected error for out-of-range selection")
	}
}

func TestPromptSelection_NonNumericInput(t *testing.T) {
	items := []selectableItem{
		{Slug: "foo", Type: "proposal", Status: "Draft"},
	}

	origStdin := stdinFunc
	stdinFunc = func() (string, error) { return "abc", nil }
	defer func() { stdinFunc = origStdin }()

	var buf bytes.Buffer
	_, err := promptSelection(items, &buf)
	if err == nil {
		t.Error("expected error for non-numeric input")
	}
}

func TestPromptSelection_StdinError(t *testing.T) {
	items := []selectableItem{
		{Slug: "foo", Type: "proposal", Status: "Draft"},
	}

	origStdin := stdinFunc
	stdinFunc = func() (string, error) { return "", fmt.Errorf("stdin error") }
	defer func() { stdinFunc = origStdin }()

	var buf bytes.Buffer
	_, err := promptSelection(items, &buf)
	if err == nil {
		t.Error("expected error for stdin failure")
	}
}

// ---------------------------------------------------------------------------
// interactive mode: -i flag registration
// ---------------------------------------------------------------------------

func TestWorktreeStart_InteractiveFlagRegistered(t *testing.T) {
	flag := worktreeStartCmd.Flags().Lookup("interactive")
	if flag == nil {
		t.Fatal("worktree start command should have --interactive flag")
	}
	if flag.Shorthand != "i" {
		t.Errorf("interactive shorthand should be 'i', got %q", flag.Shorthand)
	}
	if flag.DefValue != "false" {
		t.Errorf("interactive default should be false, got %q", flag.DefValue)
	}
}

func TestWorktreeStart_SlugArgIsOptional(t *testing.T) {
	// The command should accept 0 args when using -i
	// MaximumNArgs(1) should allow 0 or 1 args
	if worktreeStartCmd.Args == nil {
		t.Fatal("worktreeStartCmd.Args should not be nil")
	}
	// Test that 0 args is accepted by Args validator (but will fail at runtime without -i)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when no slug and no -i flag")
	}
}

// ---------------------------------------------------------------------------
// interactive mode: slug takes precedence over -i
// ---------------------------------------------------------------------------

func TestWorktreeStart_SlugTakesPrecedenceOverInteractive(t *testing.T) {
	resetSourceBranchFlag(t)
	resetInteractiveFlag(t)

	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "slug-precedence-test"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	// Both -i and slug provided: slug should take precedence
	rootCmd.SetArgs([]string{"worktree", "start", slug, "-i"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created with the provided slug
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Errorf("worktree directory %s should exist", targetDir)
	}
}

// ---------------------------------------------------------------------------
// interactive mode: empty list handling
// ---------------------------------------------------------------------------

func TestWorktreeStart_InteractiveEmptyList(t *testing.T) {
	resetSourceBranchFlag(t)
	resetInteractiveFlag(t)

	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	// Override isTerminal to return true
	origIsTerminal := isTerminalFunc
	isTerminalFunc = func() bool { return true }
	defer func() { isTerminalFunc = origIsTerminal }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// No proposals or features in this git repo

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", "-i"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stdout := buf.String()
	if !strings.Contains(stdout, "No unfinished proposals or features found") {
		t.Errorf("output should mention no items found, got:\n%s", stdout)
	}
}

// ---------------------------------------------------------------------------
// interactive mode: non-TTY detection
// ---------------------------------------------------------------------------

func TestWorktreeStart_InteractiveNonTTY(t *testing.T) {
	resetSourceBranchFlag(t)
	resetInteractiveFlag(t)

	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	// Override isTerminal to return false (non-TTY)
	origIsTerminal := isTerminalFunc
	isTerminalFunc = func() bool { return false }
	defer func() { isTerminalFunc = origIsTerminal }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", "-i"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when interactive mode used in non-TTY")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "terminal") {
		t.Errorf("error should mention terminal, got: %s", stderr)
	}
}

// ---------------------------------------------------------------------------
// interactive mode: successful selection
// ---------------------------------------------------------------------------

func TestWorktreeStart_InteractiveSelectsProposal(t *testing.T) {
	resetSourceBranchFlag(t)
	resetInteractiveFlag(t)

	origLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "claude" {
			return "/usr/bin/claude", nil
		}
		return exec.LookPath(name)
	}
	defer func() { lookPathFunc = origLookPath }()

	origRunClaude := runClaudeFunc
	runClaudeFunc = func(_ []string) error { return nil }
	defer func() { runClaudeFunc = origRunClaude }()

	// Override isTerminal to return true
	origIsTerminal := isTerminalFunc
	isTerminalFunc = func() bool { return true }
	defer func() { isTerminalFunc = origIsTerminal }()

	// Mock stdin to return "1"
	origStdin := stdinFunc
	stdinFunc = func() (string, error) { return "1", nil }
	defer func() { stdinFunc = origStdin }()

	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a proposal
	proposalDir := filepath.Join(dir, "docs", "proposals", "interactive-test")
	if err := os.MkdirAll(proposalDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\nstatus: Draft\ncreated: 2026-01-01\n---\n# Test"
	if err := os.WriteFile(filepath.Join(proposalDir, "proposal.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	targetDir := filepath.Join(dir, ".forge", "worktrees", "interactive-test")
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", "interactive-test").Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start", "-i"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify worktree was created with the selected slug
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Errorf("worktree directory %s should exist", targetDir)
	}
}

// ---------------------------------------------------------------------------
// interactive mode: no slug and no -i
// ---------------------------------------------------------------------------

func TestWorktreeStart_NoSlugNoInteractive(t *testing.T) {
	resetSourceBranchFlag(t)
	resetInteractiveFlag(t)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "start"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when no slug and no -i flag")
	}
	if !strings.Contains(err.Error(), "slug") {
		t.Errorf("error should mention slug requirement, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// resetInteractiveFlag resets the --interactive flag on worktreeStartCmd.
// ---------------------------------------------------------------------------

func resetInteractiveFlag(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		f := worktreeStartCmd.Flags().Lookup("interactive")
		if f != nil {
			f.Changed = false
			_ = f.Value.Set("false")
		}
	})
}

// ---------------------------------------------------------------------------
// worktree push: subcommand registration
// ---------------------------------------------------------------------------

func TestWorktreeCmd_HasPushSubcommand(t *testing.T) {
	subcommands := worktreeCmd.Commands()
	found := false
	for _, cmd := range subcommands {
		if cmd.Name() == "push" {
			found = true
			break
		}
	}
	if !found {
		t.Error("worktree group should have 'push' subcommand")
	}
}

// ---------------------------------------------------------------------------
// worktree push: error cases
// ---------------------------------------------------------------------------

func TestWorktreePush_ErrorWhenNotInWorktree(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Ensure isInsideWorktreeFunc returns false (default for regular repo)
	origInsideWt := isInsideWorktreeFunc
	defer func() { isInsideWorktreeFunc = origInsideWt }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "push"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when not inside a worktree")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "not inside a worktree") {
		t.Errorf("error should mention worktree context, got: %s", stderr)
	}
}

func TestWorktreePush_ErrorOnDefaultBranch(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Override to simulate being inside a worktree on default branch
	origInsideWt := isInsideWorktreeFunc
	isInsideWorktreeFunc = func(_ string) bool { return true }
	defer func() { isInsideWorktreeFunc = origInsideWt }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "push"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when on default branch in worktree")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "refusing to push default branch") {
		t.Errorf("error should mention refusing default branch, got: %s", stderr)
	}
}

func TestWorktreePush_ErrorOnPushFailure(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Override to simulate worktree on a feature branch
	origInsideWt := isInsideWorktreeFunc
	isInsideWorktreeFunc = func(_ string) bool { return true }
	defer func() { isInsideWorktreeFunc = origInsideWt }()

	origBranchFunc := getCurrentBranchFunc
	getCurrentBranchFunc = func(_ string) string { return "feature/my-branch" }
	defer func() { getCurrentBranchFunc = origBranchFunc }()

	// Override gitPushFunc to simulate failure
	origPushFunc := gitPushFunc
	gitPushFunc = func(_ string) (string, error) {
		return "", fmt.Errorf("network error: connection refused")
	}
	defer func() { gitPushFunc = origPushFunc }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "push"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when push fails")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "push failed") {
		t.Errorf("error should mention push failure, got: %s", stderr)
	}
}

// ---------------------------------------------------------------------------
// worktree push: happy path
// ---------------------------------------------------------------------------

func TestWorktreePush_Success(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Override to simulate worktree on a feature branch
	origInsideWt := isInsideWorktreeFunc
	isInsideWorktreeFunc = func(_ string) bool { return true }
	defer func() { isInsideWorktreeFunc = origInsideWt }()

	origBranchFunc := getCurrentBranchFunc
	getCurrentBranchFunc = func(_ string) string { return "my-feature" }
	defer func() { getCurrentBranchFunc = origBranchFunc }()

	// Override gitPushFunc to simulate success
	origPushFunc := gitPushFunc
	gitPushFunc = func(_ string) (string, error) {
		return "", nil
	}
	defer func() { gitPushFunc = origPushFunc }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "push"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stdout := buf.String()
	if !strings.Contains(stdout, "Pushed branch") {
		t.Errorf("output should confirm push, got: %s", stdout)
	}
	if !strings.Contains(stdout, "my-feature") {
		t.Errorf("output should mention branch name, got: %s", stdout)
	}
}

func TestWorktreePush_PrintsPushOutput(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Override to simulate worktree on a feature branch
	origInsideWt := isInsideWorktreeFunc
	isInsideWorktreeFunc = func(_ string) bool { return true }
	defer func() { isInsideWorktreeFunc = origInsideWt }()

	origBranchFunc := getCurrentBranchFunc
	getCurrentBranchFunc = func(_ string) string { return "feature/test" }
	defer func() { getCurrentBranchFunc = origBranchFunc }()

	// Override gitPushFunc to return push output
	origPushFunc := gitPushFunc
	gitPushFunc = func(_ string) (string, error) {
		return "remote: Enumerating objects: 5, done.\nTo github.com:user/repo.git\n", nil
	}
	defer func() { gitPushFunc = origPushFunc }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "push"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stdout := buf.String()
	if !strings.Contains(stdout, "Enumerating objects") {
		t.Errorf("output should include push output, got: %s", stdout)
	}
}

// ---------------------------------------------------------------------------
// worktree status: command registration
// ---------------------------------------------------------------------------

func TestWorktreeCmd_HasStatusSubcommand(t *testing.T) {
	subcommands := worktreeCmd.Commands()
	found := false
	for _, cmd := range subcommands {
		if cmd.Name() == "status" {
			found = true
			break
		}
	}
	if !found {
		t.Error("worktree group should have 'status' subcommand")
	}
}

func TestWorktreeStatusCmd_AcceptsOptionalSlug(t *testing.T) {
	// status accepts 0 or 1 args
	if worktreeStatusCmd.Args != nil {
		// cobra.MaximumNArgs returns a PositionalArgs function
		// Verify it allows 0 and 1 args
		if err := worktreeStatusCmd.Args(worktreeStatusCmd, []string{}); err != nil {
			t.Errorf("status should accept 0 args: %v", err)
		}
		if err := worktreeStatusCmd.Args(worktreeStatusCmd, []string{"my-slug"}); err != nil {
			t.Errorf("status should accept 1 arg: %v", err)
		}
	}
}

// ---------------------------------------------------------------------------
// worktree status: error cases
// ---------------------------------------------------------------------------

func TestWorktreeStatus_ErrorWhenNotGitRepo(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "status", "my-slug"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when not a git repo")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "not a git repository") {
		t.Errorf("error should mention 'not a git repository', got: %s", stderr)
	}
}

func TestWorktreeStatus_ErrorOnNonExistentSlug(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "status", "non-existent-slug"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for non-existent slug")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "not found") {
		t.Errorf("error should mention 'not found', got: %s", stderr)
	}
}

// ---------------------------------------------------------------------------
// worktree status: specific slug — shows branch, commit, uncommitted files
// ---------------------------------------------------------------------------

func TestWorktreeStatus_ShowsBranchAndCommit(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a real worktree
	slug := "status-test-wt"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	cmd := exec.Command("git", "worktree", "add", targetDir, "-b", slug)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "status", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stdout := buf.String()
	if !strings.Contains(stdout, "BRANCH:") {
		t.Errorf("output should contain BRANCH:, got: %s", stdout)
	}
	if !strings.Contains(stdout, slug) {
		t.Errorf("output should mention slug/branch %s, got: %s", slug, stdout)
	}
	if !strings.Contains(stdout, "COMMIT:") {
		t.Errorf("output should contain COMMIT:, got: %s", stdout)
	}
	if !strings.Contains(stdout, "initial") {
		t.Errorf("output should contain latest commit message, got: %s", stdout)
	}
	if !strings.Contains(stdout, "---") {
		t.Errorf("output should use structured block format (---), got: %s", stdout)
	}
}

func TestWorktreeStatus_ShowsUncommittedFiles(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a real worktree
	slug := "status-dirty-wt"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	cmd := exec.Command("git", "worktree", "add", targetDir, "-b", slug)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	// Create an uncommitted file in the worktree
	uncommittedFile := filepath.Join(targetDir, "dirty-file.txt")
	if err := os.WriteFile(uncommittedFile, []byte("dirty"), 0o644); err != nil {
		t.Fatalf("write dirty file: %v", err)
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "status", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stdout := buf.String()
	if !strings.Contains(stdout, "UNCOMMITTED:") {
		t.Errorf("output should contain UNCOMMITTED:, got: %s", stdout)
	}
	if !strings.Contains(stdout, "dirty-file.txt") {
		t.Errorf("output should list dirty-file.txt, got: %s", stdout)
	}
}

func TestWorktreeStatus_CleanWorktree(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a real worktree
	slug := "status-clean-wt"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	cmd := exec.Command("git", "worktree", "add", targetDir, "-b", slug)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "status", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stdout := buf.String()
	if !strings.Contains(stdout, "UNCOMMITTED:") {
		t.Errorf("output should contain UNCOMMITTED: (none), got: %s", stdout)
	}
	if strings.Contains(stdout, "dirty-file") {
		t.Errorf("clean worktree should not list dirty files, got: %s", stdout)
	}
}

// ---------------------------------------------------------------------------
// worktree status: no slug — shows all forge-managed worktrees
// ---------------------------------------------------------------------------

func TestWorktreeStatus_NoSlug_ShowsAllForgeWorktrees(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create two worktrees
	slug1 := "status-all-wt1"
	slug2 := "status-all-wt2"
	for _, slug := range []string{slug1, slug2} {
		targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
		cmd := exec.Command("git", "worktree", "add", targetDir, "-b", slug)
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			t.Fatalf("git worktree add %s: %v", slug, err)
		}
		t.Cleanup(func() {
			td := filepath.Join(dir, ".forge", "worktrees", slug)
			_ = exec.Command("git", "worktree", "remove", td, "--force").Run()
			_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
		})
	}

	// Create feature dir so worktrees are forge-managed
	featureDir := filepath.Join(dir, "docs", "features", "status-all-wt1")
	if err := os.MkdirAll(featureDir, 0o755); err != nil {
		t.Fatalf("create feature dir: %v", err)
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "status"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stdout := buf.String()
	// Should show the forge-managed worktree
	if !strings.Contains(stdout, "BRANCH:") {
		t.Errorf("output should contain BRANCH:, got: %s", stdout)
	}
	if !strings.Contains(stdout, "COMMIT:") {
		t.Errorf("output should contain COMMIT:, got: %s", stdout)
	}
}

func TestWorktreeStatus_NoSlug_NoWorktrees(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "status"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stdout := buf.String()
	if !strings.Contains(stdout, "No forge-managed worktrees found") {
		t.Errorf("output should indicate no forge-managed worktrees, got: %s", stdout)
	}
}

// ---------------------------------------------------------------------------
// worktree status: read-only guarantee
// ---------------------------------------------------------------------------

func TestWorktreeStatus_IsReadOnly(t *testing.T) {
	// Verify the status command does not modify filesystem state.
	// We check by running status on a clean worktree and verifying
	// the directory mtime doesn't change and git status stays clean.
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	slug := "readonly-wt"
	targetDir := filepath.Join(dir, ".forge", "worktrees", slug)
	cmd := exec.Command("git", "worktree", "add", targetDir, "-b", slug)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", slug).Run()
	})

	// Get git status before
	cmd = exec.Command("git", "status", "--porcelain")
	cmd.Dir = targetDir
	beforeStatus, _ := cmd.Output()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"worktree", "status", slug})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Get git status after
	cmd = exec.Command("git", "status", "--porcelain")
	cmd.Dir = targetDir
	afterStatus, _ := cmd.Output()

	if string(beforeStatus) != string(afterStatus) {
		t.Errorf("status command modified filesystem state: before=%q after=%q", beforeStatus, afterStatus)
	}
}

// ---------------------------------------------------------------------------
// Shell completion: ValidArgsFunction registration
// ---------------------------------------------------------------------------

func TestWorktreeStartCmd_HasValidArgsFunction(t *testing.T) {
	if worktreeStartCmd.ValidArgsFunction == nil {
		t.Error("worktreeStartCmd should have a ValidArgsFunction for shell completion")
	}
}

func TestWorktreeRemoveCmd_HasValidArgsFunction(t *testing.T) {
	if worktreeRemoveCmd.ValidArgsFunction == nil {
		t.Error("worktreeRemoveCmd should have a ValidArgsFunction for shell completion")
	}
}

func TestWorktreeResumeCmd_HasValidArgsFunction(t *testing.T) {
	if worktreeResumeCmd.ValidArgsFunction == nil {
		t.Error("worktreeResumeCmd should have a ValidArgsFunction for shell completion")
	}
}

func TestWorktreeListCmd_NoValidArgsFunction(t *testing.T) {
	if worktreeListCmd.ValidArgsFunction != nil {
		t.Error("worktreeListCmd should NOT have a ValidArgsFunction (list takes no slug arg)")
	}
}

func TestWorktreePushCmd_NoValidArgsFunction(t *testing.T) {
	if worktreePushCmd.ValidArgsFunction != nil {
		t.Error("worktreePushCmd should NOT have a ValidArgsFunction (push takes no slug arg)")
	}
}

func TestWorktreeStatusCmd_NoValidArgsFunction(t *testing.T) {
	if worktreeStatusCmd.ValidArgsFunction != nil {
		t.Error("worktreeStatusCmd should NOT have a ValidArgsFunction (status uses optional slug)")
	}
}

// ---------------------------------------------------------------------------
// Shell completion: start — unfinished proposal/feature slugs
// ---------------------------------------------------------------------------

func TestWorktreeStartCompletion_UnfinishedItems(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create proposal and feature directories
	propDir := filepath.Join(dir, "docs", "proposals", "my-proposal")
	if err := os.MkdirAll(propDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(propDir, "proposal.md"), []byte("---\nstatus: Draft\n---\n# Test"), 0o644); err != nil {
		t.Fatal(err)
	}

	featDir := filepath.Join(dir, "docs", "features", "my-feature")
	if err := os.MkdirAll(featDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(featDir, "manifest.md"), []byte("---\nstatus: in_progress\n---\n# Test"), 0o644); err != nil {
		t.Fatal(err)
	}

	completions, directive := worktreeStartCmd.ValidArgsFunction(worktreeStartCmd, []string{}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected ShellCompDirectiveNoFileComp, got %v", directive)
	}

	var slugs []string
	for _, c := range completions {
		parts := strings.SplitN(c, "\t", 2)
		slugs = append(slugs, parts[0])
	}

	if !containsStr(slugs, "my-proposal") {
		t.Errorf("expected completions to contain 'my-proposal', got %v", slugs)
	}
	if !containsStr(slugs, "my-feature") {
		t.Errorf("expected completions to contain 'my-feature', got %v", slugs)
	}
}

func TestWorktreeStartCompletion_FilterByPrefix(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create two proposals
	for _, slug := range []string{"alpha-proposal", "beta-proposal"} {
		propDir := filepath.Join(dir, "docs", "proposals", slug)
		if err := os.MkdirAll(propDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(propDir, "proposal.md"), []byte("---\nstatus: Draft\n---\n# Test"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	completions, _ := worktreeStartCmd.ValidArgsFunction(worktreeStartCmd, []string{}, "alpha")

	var slugs []string
	for _, c := range completions {
		parts := strings.SplitN(c, "\t", 2)
		slugs = append(slugs, parts[0])
	}

	if !containsStr(slugs, "alpha-proposal") {
		t.Errorf("expected completions to contain 'alpha-proposal', got %v", slugs)
	}
	if containsStr(slugs, "beta-proposal") {
		t.Errorf("expected completions NOT to contain 'beta-proposal' when filtering by 'alpha', got %v", slugs)
	}
}

func TestWorktreeStartCompletion_SkipsCompleted(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a completed proposal — should not appear in completions
	propDir := filepath.Join(dir, "docs", "proposals", "done-proposal")
	if err := os.MkdirAll(propDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(propDir, "proposal.md"), []byte("---\nstatus: completed\n---\n# Done"), 0o644); err != nil {
		t.Fatal(err)
	}

	completions, _ := worktreeStartCmd.ValidArgsFunction(worktreeStartCmd, []string{}, "")

	var slugs []string
	for _, c := range completions {
		parts := strings.SplitN(c, "\t", 2)
		slugs = append(slugs, parts[0])
	}

	if containsStr(slugs, "done-proposal") {
		t.Errorf("completed proposals should not appear in completions, got %v", slugs)
	}
}

func TestWorktreeStartCompletion_ErrorReturnsEmpty(t *testing.T) {
	// No project root → FindProjectRoot will fail → should return empty list
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	completions, directive := worktreeStartCmd.ValidArgsFunction(worktreeStartCmd, []string{}, "")
	if len(completions) != 0 {
		t.Errorf("error should return empty completions, got %v", completions)
	}
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected ShellCompDirectiveNoFileComp on error, got %v", directive)
	}
}

func TestWorktreeStartCompletion_AlreadyHasArg(t *testing.T) {
	// When args already has a slug (cobra.ExactArgs or arg already provided), return empty
	completions, directive := worktreeStartCmd.ValidArgsFunction(worktreeStartCmd, []string{"existing-slug"}, "")
	if len(completions) != 0 {
		t.Errorf("should return empty when arg already provided, got %v", completions)
	}
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected ShellCompDirectiveNoFileComp, got %v", directive)
	}
}

// ---------------------------------------------------------------------------
// Shell completion: remove/resume — existing worktree slugs
// ---------------------------------------------------------------------------

func TestWorktreeRemoveCompletion_ExistingSlugs(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Override listWorktreesFunc to return predictable results
	origList := listWorktreesFunc
	listWorktreesFunc = func(_ string) ([]gitPkg.WorktreeEntry, error) {
		return []gitPkg.WorktreeEntry{
			{Path: dir, Branch: "main", IsMain: true},
			{Path: filepath.Join(dir, ".forge", "worktrees", "feature-a"), Branch: "feature-a"},
			{Path: filepath.Join(dir, ".forge", "worktrees", "feature-b"), Branch: "feature-b"},
		}, nil
	}
	defer func() { listWorktreesFunc = origList }()

	completions, directive := worktreeRemoveCmd.ValidArgsFunction(worktreeRemoveCmd, []string{}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected ShellCompDirectiveNoFileComp, got %v", directive)
	}

	var slugs []string
	for _, c := range completions {
		parts := strings.SplitN(c, "\t", 2)
		slugs = append(slugs, parts[0])
	}

	// Main worktree (basename of dir) should not appear
	mainName := filepath.Base(dir)
	if containsStr(slugs, mainName) {
		t.Errorf("main worktree %q should not appear, got %v", mainName, slugs)
	}
	if !containsStr(slugs, "feature-a") {
		t.Errorf("expected 'feature-a' in completions, got %v", slugs)
	}
	if !containsStr(slugs, "feature-b") {
		t.Errorf("expected 'feature-b' in completions, got %v", slugs)
	}
}

func TestWorktreeRemoveCompletion_FilterByPrefix(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	origList := listWorktreesFunc
	listWorktreesFunc = func(_ string) ([]gitPkg.WorktreeEntry, error) {
		return []gitPkg.WorktreeEntry{
			{Path: dir, Branch: "main", IsMain: true},
			{Path: filepath.Join(dir, ".forge", "worktrees", "alpha"), Branch: "alpha"},
			{Path: filepath.Join(dir, ".forge", "worktrees", "beta"), Branch: "beta"},
		}, nil
	}
	defer func() { listWorktreesFunc = origList }()

	completions, _ := worktreeRemoveCmd.ValidArgsFunction(worktreeRemoveCmd, []string{}, "alp")

	var slugs []string
	for _, c := range completions {
		parts := strings.SplitN(c, "\t", 2)
		slugs = append(slugs, parts[0])
	}

	if !containsStr(slugs, "alpha") {
		t.Errorf("expected 'alpha' when filtering by 'alp', got %v", slugs)
	}
	if containsStr(slugs, "beta") {
		t.Errorf("expected 'beta' to be filtered out when prefix is 'alp', got %v", slugs)
	}
}

func TestWorktreeRemoveCompletion_ErrorReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	origList := listWorktreesFunc
	listWorktreesFunc = func(_ string) ([]gitPkg.WorktreeEntry, error) {
		return nil, fmt.Errorf("simulated error")
	}
	defer func() { listWorktreesFunc = origList }()

	completions, directive := worktreeRemoveCmd.ValidArgsFunction(worktreeRemoveCmd, []string{}, "")
	if len(completions) != 0 {
		t.Errorf("error should return empty completions, got %v", completions)
	}
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected ShellCompDirectiveNoFileComp on error, got %v", directive)
	}
}

func TestWorktreeResumeCompletion_ExistingSlugs(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	origList := listWorktreesFunc
	listWorktreesFunc = func(_ string) ([]gitPkg.WorktreeEntry, error) {
		return []gitPkg.WorktreeEntry{
			{Path: dir, Branch: "main", IsMain: true},
			{Path: filepath.Join(dir, ".forge", "worktrees", "resume-me"), Branch: "resume-me"},
		}, nil
	}
	defer func() { listWorktreesFunc = origList }()

	completions, directive := worktreeResumeCmd.ValidArgsFunction(worktreeResumeCmd, []string{}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected ShellCompDirectiveNoFileComp, got %v", directive)
	}

	var slugs []string
	for _, c := range completions {
		parts := strings.SplitN(c, "\t", 2)
		slugs = append(slugs, parts[0])
	}

	if !containsStr(slugs, "resume-me") {
		t.Errorf("expected 'resume-me' in completions, got %v", slugs)
	}
}

func TestWorktreeResumeCompletion_ErrorReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	origList := listWorktreesFunc
	listWorktreesFunc = func(_ string) ([]gitPkg.WorktreeEntry, error) {
		return nil, fmt.Errorf("simulated error")
	}
	defer func() { listWorktreesFunc = origList }()

	completions, directive := worktreeResumeCmd.ValidArgsFunction(worktreeResumeCmd, []string{}, "")
	if len(completions) != 0 {
		t.Errorf("error should return empty completions, got %v", completions)
	}
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected ShellCompDirectiveNoFileComp on error, got %v", directive)
	}
}

func TestWorktreeRemoveCompletion_ExcludesMainWorktree(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	origList := listWorktreesFunc
	listWorktreesFunc = func(_ string) ([]gitPkg.WorktreeEntry, error) {
		return []gitPkg.WorktreeEntry{
			{Path: dir, Branch: "main", IsMain: true},
		}, nil
	}
	defer func() { listWorktreesFunc = origList }()

	completions, _ := worktreeRemoveCmd.ValidArgsFunction(worktreeRemoveCmd, []string{}, "")
	if len(completions) != 0 {
		t.Errorf("only main worktree — should return empty, got %v", completions)
	}
}

// containsStr checks if a string slice contains the given value.
func containsStr(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
