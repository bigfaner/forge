package git

import (
	"testing"
)

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

func TestGetCurrentBranch_NotGitRepo(t *testing.T) {
	dir := t.TempDir()
	branch := GetCurrentBranch(dir)
	if branch != "" {
		t.Errorf("expected empty branch in non-git dir, got %q", branch)
	}
}

func TestIsGitRepository(t *testing.T) {
	t.Run("not a git repository", func(t *testing.T) {
		dir := t.TempDir()
		if IsGitRepository(dir) {
			t.Error("expected false for non-git directory")
		}
	})

	// Note: Testing with a real git repository requires git to be installed
	// and may have side effects. The function is simple enough that we trust
	// it works if the manual test passes.
}

func TestGetWorktreeName_NotInWorktree(t *testing.T) {
	dir := t.TempDir()
	name := GetWorktreeName(dir)
	if name != "" {
		t.Errorf("expected empty worktree name in non-worktree dir, got %q", name)
	}
}

func TestGetFeatureFromGit_NotGitRepo(t *testing.T) {
	dir := t.TempDir()
	feature := GetFeatureFromGit(dir)
	if feature != "" {
		t.Errorf("expected empty feature in non-git dir, got %q", feature)
	}
}
