package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// helper: initGitRepo creates a real git repo in a temp dir with an initial commit.
// Returns the directory path.
func initGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// git init
	if err := runGit(dir, "init"); err != nil {
		t.Fatalf("git init: %v", err)
	}

	// Configure user so commits work
	if err := runGit(dir, "config", "user.email", "test@test.com"); err != nil {
		t.Fatalf("git config email: %v", err)
	}
	if err := runGit(dir, "config", "user.name", "Test"); err != nil {
		t.Fatalf("git config name: %v", err)
	}

	// Create an initial file and commit so HEAD is not empty
	f := filepath.Join(dir, "README.md")
	if err := os.WriteFile(f, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := runGit(dir, "add", "."); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := runGit(dir, "commit", "-m", "initial"); err != nil {
		t.Fatalf("git commit: %v", err)
	}

	return dir
}

// helper: runGit executes a git command in the given directory.
func runGit(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	return cmd.Run()
}

// ---------------------------------------------------------------------------
// ExtractFeatureFromBranch
// ---------------------------------------------------------------------------

func TestExtractFeatureFromBranch(t *testing.T) {
	tests := []struct {
		branch string
		want   string
	}{
		{"feature/auth-login", "auth-login"},
		{"feat/user-registration", "user-registration"},
		{"fix/null-pointer", "null-pointer"},
		{"bugfix/memory-leak", "memory-leak"},
		{"hotfix/security-issue", "security-issue"},
		{"chore/update-deps", "update-deps"},
		{"main", "main"},
		{"master", "master"},
		{"custom-branch", "custom-branch"},
		{"feature/nested/path", "nested/path"},
		{"nested/branch/name", "nested-branch-name"},
	}

	for _, tt := range tests {
		t.Run(tt.branch, func(t *testing.T) {
			got := ExtractFeatureFromBranch(tt.branch)
			if got != tt.want {
				t.Errorf("ExtractFeatureFromBranch(%q) = %q, want %q", tt.branch, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// GetCurrentBranch
// ---------------------------------------------------------------------------

func TestGetCurrentBranch_NotGitRepo(t *testing.T) {
	dir := t.TempDir()
	branch := GetCurrentBranch(dir)
	if branch != "" {
		t.Errorf("expected empty branch in non-git dir, got %q", branch)
	}
}

func TestGetCurrentBranch_RealRepo(t *testing.T) {
	dir := initGitRepo(t)

	// After init + commit on the default branch, the branch name should be
	// either "main" or "master" depending on the git version.
	branch := GetCurrentBranch(dir)
	if branch == "" {
		t.Fatal("expected non-empty branch name, got empty")
	}
	if branch != "main" && branch != "master" {
		t.Errorf("expected main or master, got %q", branch)
	}
}

func TestGetCurrentBranch_FeatureBranch(t *testing.T) {
	dir := initGitRepo(t)

	// Create and checkout a feature branch
	if err := runGit(dir, "checkout", "-b", "feature/my-feature"); err != nil {
		t.Fatalf("git checkout -b: %v", err)
	}

	branch := GetCurrentBranch(dir)
	if branch != "feature/my-feature" {
		t.Errorf("expected %q, got %q", "feature/my-feature", branch)
	}
}

// ---------------------------------------------------------------------------
// IsGitRepository
// ---------------------------------------------------------------------------

func TestIsGitRepository(t *testing.T) {
	t.Run("not a git repository", func(t *testing.T) {
		dir := t.TempDir()
		if IsGitRepository(dir) {
			t.Error("expected false for non-git directory")
		}
	})

	t.Run("real git repository", func(t *testing.T) {
		dir := initGitRepo(t)
		if !IsGitRepository(dir) {
			t.Error("expected true for git directory")
		}
	})
}

// ---------------------------------------------------------------------------
// GetWorktreeName
// ---------------------------------------------------------------------------

func TestGetWorktreeName_NotInWorktree(t *testing.T) {
	dir := t.TempDir()
	name := GetWorktreeName(dir)
	if name != "" {
		t.Errorf("expected empty worktree name in non-worktree dir, got %q", name)
	}
}

func TestGetWorktreeName_GitFileWorktree(t *testing.T) {
	// Simulate a worktree by creating a .git file (not directory) with the
	// standard "gitdir: ..." content.
	dir := t.TempDir()

	// Also create a fake main repo worktree directory so the path is plausible.
	mainRepo := t.TempDir()
	worktreesDir := filepath.Join(mainRepo, ".git", "worktrees", "my-worktree")
	if err := os.MkdirAll(worktreesDir, 0o755); err != nil {
		t.Fatalf("mkdir all: %v", err)
	}

	gitFile := filepath.Join(dir, ".git")
	content := "gitdir: " + strings.ReplaceAll(filepath.Join(mainRepo, ".git", "worktrees", "my-worktree"), "\\", "/") + "\n"
	if err := os.WriteFile(gitFile, []byte(content), 0o644); err != nil {
		t.Fatalf("write .git file: %v", err)
	}

	name := GetWorktreeName(dir)
	if name != "my-worktree" {
		t.Errorf("expected %q, got %q", "my-worktree", name)
	}
}

func TestGetWorktreeName_GitFileWithoutGitdirPrefix(t *testing.T) {
	dir := t.TempDir()
	gitFile := filepath.Join(dir, ".git")
	if err := os.WriteFile(gitFile, []byte("garbage content"), 0o644); err != nil {
		t.Fatalf("write .git file: %v", err)
	}

	name := GetWorktreeName(dir)
	if name != "" {
		t.Errorf("expected empty name for malformed .git file, got %q", name)
	}
}

func TestGetWorktreeName_GitFileWithoutWorktreesPath(t *testing.T) {
	dir := t.TempDir()
	gitFile := filepath.Join(dir, ".git")
	content := "gitdir: /some/random/path\n"
	if err := os.WriteFile(gitFile, []byte(content), 0o644); err != nil {
		t.Fatalf("write .git file: %v", err)
	}

	name := GetWorktreeName(dir)
	if name != "" {
		t.Errorf("expected empty name when no /worktrees/ in path, got %q", name)
	}
}

func TestGetWorktreeName_GitFileWithTrailingPath(t *testing.T) {
	dir := t.TempDir()
	mainRepo := t.TempDir()

	gitFile := filepath.Join(dir, ".git")
	// Include trailing path after worktree name, e.g. ".../worktrees/wt-name/something"
	worktreePath := filepath.Join(mainRepo, ".git", "worktrees", "wt-name", "something")
	content := "gitdir: " + strings.ReplaceAll(worktreePath, "\\", "/") + "\n"
	if err := os.WriteFile(gitFile, []byte(content), 0o644); err != nil {
		t.Fatalf("write .git file: %v", err)
	}

	name := GetWorktreeName(dir)
	if name != "wt-name" {
		t.Errorf("expected %q, got %q", "wt-name", name)
	}
}

func TestGetWorktreeName_RegularGitDir(t *testing.T) {
	// When .git is a regular directory (not a file), it falls through to the
	// git worktree list code path. In a normal repo with a single worktree,
	// the current worktree entry should match the project root.
	dir := initGitRepo(t)

	name := GetWorktreeName(dir)
	// In a non-worktree situation, the worktree list includes the main worktree.
	// The code extracts the branch name from the [branch] field and runs
	// ExtractFeatureFromBranch. For "main"/"master" this returns the name as-is
	// (no prefix match), but the code path depends on absRoot matching.
	// We just ensure it doesn't panic and returns a string (may be empty or a slug).
	_ = name
}

// ---------------------------------------------------------------------------
// GetFeatureFromGit
// ---------------------------------------------------------------------------

func TestGetFeatureFromGit_NotGitRepo(t *testing.T) {
	dir := t.TempDir()
	feature := GetFeatureFromGit(dir)
	if feature != "" {
		t.Errorf("expected empty feature in non-git dir, got %q", feature)
	}
}

func TestGetFeatureFromGit_MainBranch(t *testing.T) {
	dir := initGitRepo(t)

	// Ensure we're on the default branch
	branch := GetCurrentBranch(dir)
	if branch != "main" {
		if err := runGit(dir, "checkout", "-b", "main"); err != nil {
			t.Fatalf("git checkout -b main: %v", err)
		}
	}

	feature := GetFeatureFromGit(dir)
	if feature != "" {
		t.Errorf("expected empty feature on main branch, got %q", feature)
	}
}

func TestGetFeatureFromGit_FeatureBranch(t *testing.T) {
	dir := initGitRepo(t)

	if err := runGit(dir, "checkout", "-b", "feature/add-auth"); err != nil {
		t.Fatalf("git checkout -b: %v", err)
	}

	feature := GetFeatureFromGit(dir)
	if feature != "add-auth" {
		t.Errorf("expected %q, got %q", "add-auth", feature)
	}
}

func TestGetFeatureFromGit_FixBranch(t *testing.T) {
	dir := initGitRepo(t)

	if err := runGit(dir, "checkout", "-b", "fix/timeout-bug"); err != nil {
		t.Fatalf("git checkout -b: %v", err)
	}

	feature := GetFeatureFromGit(dir)
	if feature != "timeout-bug" {
		t.Errorf("expected %q, got %q", "timeout-bug", feature)
	}
}

func TestGetFeatureFromGit_WorktreePriority(t *testing.T) {
	// When .git is a file (worktree), GetFeatureFromGit should return the
	// worktree name, not the branch name.
	dir := t.TempDir()
	mainRepo := t.TempDir()
	worktreesDir := filepath.Join(mainRepo, ".git", "worktrees", "wt-slug")
	if err := os.MkdirAll(worktreesDir, 0o755); err != nil {
		t.Fatalf("mkdir all: %v", err)
	}

	gitFile := filepath.Join(dir, ".git")
	content := "gitdir: " + strings.ReplaceAll(filepath.Join(mainRepo, ".git", "worktrees", "wt-slug"), "\\", "/") + "\n"
	if err := os.WriteFile(gitFile, []byte(content), 0o644); err != nil {
		t.Fatalf("write .git file: %v", err)
	}

	feature := GetFeatureFromGit(dir)
	if feature != "wt-slug" {
		t.Errorf("expected worktree name %q, got %q", "wt-slug", feature)
	}
}

// ---------------------------------------------------------------------------
// Run
// ---------------------------------------------------------------------------

func TestRun_Success(t *testing.T) {
	dir := initGitRepo(t)

	out, err := Run(dir, "rev-parse", "--git-dir")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if out == "" {
		t.Error("expected non-empty output from git rev-parse --git-dir")
	}
}

func TestRun_Failure(t *testing.T) {
	dir := t.TempDir()

	_, err := Run(dir, "status")
	if err == nil {
		t.Error("expected error running git status in non-git dir")
	}
}

func TestRun_OutputTrimmed(t *testing.T) {
	dir := initGitRepo(t)

	out, err := Run(dir, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	// Verify no trailing newline
	if strings.Contains(out, "\n") {
		t.Errorf("output should be trimmed, got %q", out)
	}
}

// ---------------------------------------------------------------------------
// IsInsideWorktree
// ---------------------------------------------------------------------------

func TestIsInsideWorktree_NotGitDir(t *testing.T) {
	dir := t.TempDir()
	if IsInsideWorktree(dir) {
		t.Error("expected false for non-git directory")
	}
}

func TestIsInsideWorktree_RegularGitRepo(t *testing.T) {
	dir := initGitRepo(t)
	if IsInsideWorktree(dir) {
		t.Error("expected false for regular git repo (.git is a directory)")
	}
}

func TestIsInsideWorktree_LinkedWorktree(t *testing.T) {
	// Simulate a linked worktree: .git is a file, not a directory
	dir := t.TempDir()
	mainRepo := t.TempDir()
	gitFile := filepath.Join(dir, ".git")
	content := "gitdir: " + strings.ReplaceAll(filepath.Join(mainRepo, ".git", "worktrees", "my-wt"), "\\", "/") + "\n"
	if err := os.WriteFile(gitFile, []byte(content), 0o644); err != nil {
		t.Fatalf("write .git file: %v", err)
	}
	if !IsInsideWorktree(dir) {
		t.Error("expected true for linked worktree (.git is a file)")
	}
}

// ---------------------------------------------------------------------------
// Push
// ---------------------------------------------------------------------------

func TestPush_NotGitRepo(t *testing.T) {
	dir := t.TempDir()
	_, err := Push(dir)
	if err == nil {
		t.Error("expected error when pushing from non-git directory")
	}
}

// ---------------------------------------------------------------------------
// ErrNoUpstream
// ---------------------------------------------------------------------------

func TestErrNoUpstream_IsSentinel(t *testing.T) {
	err := ErrNoUpstream
	if err == nil {
		t.Error("ErrNoUpstream should be non-nil")
	}
	if err.Error() != "no upstream tracking branch" {
		t.Errorf("unexpected error message: %q", err.Error())
	}
}

// ---------------------------------------------------------------------------
// CountUnpushedCommits
// ---------------------------------------------------------------------------

func TestCountUnpushedCommits_NoUpstream(t *testing.T) {
	dir := initGitRepo(t)

	count, err := CountUnpushedCommits(dir)
	if count != 0 {
		t.Errorf("expected count 0, got %d", count)
	}
	if err != ErrNoUpstream {
		t.Errorf("expected ErrNoUpstream, got %v", err)
	}
}

func TestCountUnpushedCommits_NoUpstream_OnFeatureBranch(t *testing.T) {
	dir := initGitRepo(t)

	// Create a feature branch (still no remote/upstream)
	if err := runGit(dir, "checkout", "-b", "feature/test"); err != nil {
		t.Fatalf("git checkout -b: %v", err)
	}

	count, err := CountUnpushedCommits(dir)
	if count != 0 {
		t.Errorf("expected count 0, got %d", count)
	}
	if err != ErrNoUpstream {
		t.Errorf("expected ErrNoUpstream, got %v", err)
	}
}

func TestCountUnpushedCommits_WithUpstream_ZeroUnpushed(t *testing.T) {
	dir := initGitRepo(t)

	// Create a "remote" repo and push to it to establish upstream tracking
	remoteDir := t.TempDir()
	if err := runGit(remoteDir, "init", "--bare"); err != nil {
		t.Fatalf("git init --bare: %v", err)
	}
	if err := runGit(dir, "remote", "add", "origin", remoteDir); err != nil {
		t.Fatalf("git remote add: %v", err)
	}
	// Push and set upstream
	cmd := exec.Command("git", "push", "-u", "origin", "HEAD")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git push -u origin HEAD: %v", err)
	}

	count, err := CountUnpushedCommits(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 unpushed commits, got %d", count)
	}
}

func TestCountUnpushedCommits_WithUpstream_SomeUnpushed(t *testing.T) {
	dir := initGitRepo(t)

	// Create a bare remote and push to set upstream
	remoteDir := t.TempDir()
	if err := runGit(remoteDir, "init", "--bare"); err != nil {
		t.Fatalf("git init --bare: %v", err)
	}
	if err := runGit(dir, "remote", "add", "origin", remoteDir); err != nil {
		t.Fatalf("git remote add: %v", err)
	}
	cmd := exec.Command("git", "push", "-u", "origin", "HEAD")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git push -u origin HEAD: %v", err)
	}

	// Create 3 additional commits
	for i := 0; i < 3; i++ {
		f := filepath.Join(dir, fmt.Sprintf("file%d.txt", i))
		if err := os.WriteFile(f, []byte(fmt.Sprintf("content %d", i)), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := runGit(dir, "add", "."); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := runGit(dir, "commit", "-m", fmt.Sprintf("commit %d", i)); err != nil {
			t.Fatalf("git commit: %v", err)
		}
	}

	count, err := CountUnpushedCommits(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 unpushed commits, got %d", count)
	}
}

func TestCountUnpushedCommits_NotGitRepo(t *testing.T) {
	dir := t.TempDir()

	_, err := CountUnpushedCommits(dir)
	if err == nil {
		t.Error("expected error for non-git directory")
	}
}
