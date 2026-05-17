package cmd

import (
	"bytes"
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

	// Create the target sibling directory ahead of time
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, "test-slug")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatalf("create target dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(targetDir) })

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
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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

	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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

	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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

func TestResolveSourceBranch(t *testing.T) {
	tests := []struct {
		name         string
		flagValue    string
		configBranch string
		want         string
	}{
		{
			name:         "flag overrides everything",
			flagValue:    "develop",
			configBranch: "main",
			want:         "develop",
		},
		{
			name:         "flag overrides empty config",
			flagValue:    "v3.0.0",
			configBranch: "",
			want:         "v3.0.0",
		},
		{
			name:         "config used when no flag",
			flagValue:    "",
			configBranch: "develop",
			want:         "develop",
		},
		{
			name:         "empty when neither set",
			flagValue:    "",
			configBranch: "",
			want:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveSourceBranch(tt.flagValue, tt.configBranch)
			if got != tt.want {
				t.Errorf("resolveSourceBranch(%q, %q) = %q, want %q",
					tt.flagValue, tt.configBranch, got, tt.want)
			}
		})
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
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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

func TestWorktreeRemove_RemovesWorktreeAndKeepsBranch(t *testing.T) {
	dir := initGitRepoForWorktree(t)
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	_ = os.Chdir(dir)

	// Create a worktree
	slug := "remove-test-feature"
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, slug)
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
	if capturedWd != absTarget {
		t.Errorf("claude should have been launched in %s, got %s", absTarget, capturedWd)
	}
}
