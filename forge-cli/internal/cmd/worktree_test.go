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
// worktree start: table-driven remote branch detection (mocked gitRunFunc)
// ---------------------------------------------------------------------------

// mockGitResponse is a simplified mock response for git commands.
type mockGitResponse struct {
	output string
	err    error
}

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
				"worktree add -b remote-only-feature":                   {},
			},
			wantErr:            false,
			wantStdoutContains: "origin/remote-only-feature",
			wantWorktreeArgs:   []string{"worktree", "add", "-b", "remote-only-feature"},
		},
		{
			name:         "fetch fails, no remote branch, falls back to HEAD",
			slug:         "fetch-fail-feature",
			sourceBranch: "",
			mockResponses: map[string]mockGitResponse{
				"rev-parse --verify fetch-fail-feature":                {err: fmt.Errorf("not found")},
				"fetch origin":                                         {err: fmt.Errorf("network error")},
				"rev-parse --verify remotes/origin/fetch-fail-feature": {err: fmt.Errorf("not found")},
				"worktree add -b fetch-fail-feature":                   {},
			},
			wantErr:            false,
			wantStderrContains: "warning: git fetch origin failed",
			wantWorktreeArgs:   []string{"worktree", "add", "-b", "fetch-fail-feature"},
		},
		{
			name:         "both local and remote absent, creates from HEAD",
			slug:         "new-feature",
			sourceBranch: "",
			mockResponses: map[string]mockGitResponse{
				"rev-parse --verify new-feature":                {err: fmt.Errorf("not found")},
				"fetch origin":                                  {},
				"rev-parse --verify remotes/origin/new-feature": {err: fmt.Errorf("not found")},
				"worktree add -b new-feature":                   {},
			},
			wantErr:          false,
			wantWorktreeArgs: []string{"worktree", "add", "-b", "new-feature"},
		},
		{
			name:         "source-branch set but remote branch exists, remote wins",
			slug:         "remote-wins-feature",
			sourceBranch: "develop",
			mockResponses: map[string]mockGitResponse{
				"rev-parse --verify remote-wins-feature":                {err: fmt.Errorf("not found")},
				"fetch origin":                                          {},
				"rev-parse --verify remotes/origin/remote-wins-feature": {},
				"worktree add -b remote-wins-feature":                   {},
			},
			wantErr:            false,
			wantStdoutContains: "origin/remote-wins-feature",
			wantWorktreeArgs:   []string{"worktree", "add", "-b", "remote-wins-feature"},
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
					if len(args) >= 4 && args[0] == "worktree" && args[1] == "add" {
						// worktree add targetDir slug (layer 1) or worktree add -b slug targetDir [ref] (layer 2/3)
						targetIdx := 2
						if args[2] == "-b" {
							targetIdx = 4
						}
						if targetIdx < len(args) {
							_ = os.MkdirAll(args[targetIdx], 0o755)
						}
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
