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
// worktree start: GetWorktreeName auto-detects feature
// ---------------------------------------------------------------------------

func TestWorktreeStart_WorktreeNameAutoDetection(t *testing.T) {
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
