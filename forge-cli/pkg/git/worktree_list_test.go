package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// ParsePorcelainWorktrees
// ---------------------------------------------------------------------------

func TestParsePorcelainWorktrees_SingleWorktree(t *testing.T) {
	input := `worktree /home/user/project
HEAD abc123def456
branch refs/heads/main`

	entries, err := ParsePorcelainWorktrees(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Path != "/home/user/project" {
		t.Errorf("Path = %q, want %q", entries[0].Path, "/home/user/project")
	}
	if entries[0].Branch != "main" {
		t.Errorf("Branch = %q, want %q", entries[0].Branch, "main")
	}
	if entries[0].HEAD != "abc123def456" {
		t.Errorf("HEAD = %q, want %q", entries[0].HEAD, "abc123def456")
	}
	if !entries[0].IsMain {
		t.Errorf("IsMain should be true for the first (main) worktree")
	}
}

func TestParsePorcelainWorktrees_MultipleWorktrees(t *testing.T) {
	input := `worktree /home/user/project
HEAD abc123def456
branch refs/heads/main

worktree /home/user/auth-login
HEAD def456abc789
branch refs/heads/auth-login`

	entries, err := ParsePorcelainWorktrees(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	// First entry is the main worktree
	if !entries[0].IsMain {
		t.Error("first entry should be main worktree")
	}
	if entries[0].Branch != "main" {
		t.Errorf("first Branch = %q, want %q", entries[0].Branch, "main")
	}

	// Second entry
	if entries[1].IsMain {
		t.Error("second entry should not be main worktree")
	}
	if entries[1].Path != "/home/user/auth-login" {
		t.Errorf("second Path = %q, want %q", entries[1].Path, "/home/user/auth-login")
	}
	if entries[1].Branch != "auth-login" {
		t.Errorf("second Branch = %q, want %q", entries[1].Branch, "auth-login")
	}
}

func TestParsePorcelainWorktrees_DetachedHead(t *testing.T) {
	input := `worktree /home/user/project
HEAD abc123def456
branch refs/heads/main

worktree /home/user/detached
HEAD def456abc789`

	entries, err := ParsePorcelainWorktrees(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[1].Branch != "" {
		t.Errorf("detached HEAD should have empty branch, got %q", entries[1].Branch)
	}
}

func TestParsePorcelainWorktrees_EmptyInput(t *testing.T) {
	entries, err := ParsePorcelainWorktrees("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for empty input, got %d", len(entries))
	}
}

func TestParsePorcelainWorktrees_BareRepo(t *testing.T) {
	input := `bare
HEAD abc123def456`

	entries, err := ParsePorcelainWorktrees(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for bare repo, got %d", len(entries))
	}
}

// ---------------------------------------------------------------------------
// ListWorktrees (integration with git binary)
// ---------------------------------------------------------------------------

func TestListWorktrees_RealGitRepo(t *testing.T) {
	dir := initTestGitRepo(t)

	entries, err := ListWorktrees(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry (main), got %d", len(entries))
	}
	if !entries[0].IsMain {
		t.Error("only entry should be main worktree")
	}
	if entries[0].Branch == "" {
		t.Error("branch should not be empty")
	}
}

func TestListWorktrees_WithAdditionalWorktree(t *testing.T) {
	dir := initTestGitRepo(t)
	parentDir := filepath.Dir(dir)
	targetDir := filepath.Join(parentDir, "test-feature")

	// Create a worktree
	cmd := exec.Command("git", "worktree", "add", "-b", "test-feature", targetDir)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add: %v", err)
	}
	t.Cleanup(func() {
		_ = exec.Command("git", "worktree", "remove", targetDir, "--force").Run()
		_ = exec.Command("git", "-C", dir, "branch", "-D", "test-feature").Run()
	})

	entries, err := ListWorktrees(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	// Verify the feature worktree
	found := false
	for _, e := range entries {
		if e.Branch == "test-feature" {
			found = true
			if e.Name() != "test-feature" {
				t.Errorf("Name() = %q, want %q", e.Name(), "test-feature")
			}
			break
		}
	}
	if !found {
		t.Error("expected to find worktree with branch 'test-feature'")
	}
}

func TestListWorktrees_NotGitRepo(t *testing.T) {
	dir := t.TempDir()

	_, err := ListWorktrees(dir)
	if err == nil {
		t.Error("expected error when not a git repo")
	}
}

// ---------------------------------------------------------------------------
// WorktreeEntry.Name()
// ---------------------------------------------------------------------------

func TestWorktreeEntry_Name_FromPath(t *testing.T) {
	e := WorktreeEntry{Path: "/home/user/my-feature"}
	if got := e.Name(); got != "my-feature" {
		t.Errorf("Name() = %q, want %q", got, "my-feature")
	}
}

func TestWorktreeEntry_Name_TrailingSlash(t *testing.T) {
	e := WorktreeEntry{Path: "/home/user/my-feature/"}
	if got := e.Name(); got != "my-feature" {
		t.Errorf("Name() = %q, want %q", got, "my-feature")
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func initTestGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			t.Fatalf("%s: %v", strings.Join(args, " "), err)
		}
	}

	// Initial commit
	f := filepath.Join(dir, "README.md")
	if err := os.WriteFile(f, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	cmd := exec.Command("git", "add", ".")
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
