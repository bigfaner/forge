// Package git provides git-related utilities.
package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetCurrentBranch returns the current git branch name.
// Returns empty string if not in a git repository.
func GetCurrentBranch(projectRoot string) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = projectRoot
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// GetWorktreeName returns the current worktree name.
// Returns empty string if not in a worktree or on error.
func GetWorktreeName(projectRoot string) string {
	// Check if we're in a worktree by comparing .git path
	gitDir := filepath.Join(projectRoot, ".git")
	info, err := os.Lstat(gitDir)
	if err != nil {
		return ""
	}

	// If .git is a file (not directory), we're in a worktree
	if !info.IsDir() {
		// Read the .git file to get the worktree path
		data, err := os.ReadFile(gitDir)
		if err != nil {
			return ""
		}
		// Format: gitdir: /path/to/main/.git/worktrees/<name>
		content := string(data)
		_, gitdirPath, found := strings.Cut(content, "gitdir: ")
		if !found {
			return ""
		}
		gitdirPath = strings.TrimSpace(gitdirPath)

		// Extract worktree name from path like .../worktrees/<name>
		_, afterWorktrees, found := strings.Cut(gitdirPath, "/worktrees/")
		if !found {
			return ""
		}
		// Remove trailing path components
		if slashIdx := strings.Index(afterWorktrees, "/"); slashIdx != -1 {
			return afterWorktrees[:slashIdx]
		}
		return afterWorktrees
	}

	// Alternative: use git worktree list
	cmd := exec.Command("git", "worktree", "list")
	cmd.Dir = projectRoot
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// Parse worktree list to find current worktree name
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			// Format: /path/to/worktree  commit  [branch]
			// Check if this is the current worktree
			absRoot, _ := filepath.Abs(projectRoot)
			if fields[0] == absRoot {
				// Extract branch name and convert to feature slug
				branch := strings.Trim(fields[2], "[]")
				return ExtractFeatureFromBranch(branch)
			}
		}
	}

	return ""
}

// ExtractFeatureFromBranch extracts feature slug from branch name.
// E.g., "feature/auth-login" -> "auth-login"
func ExtractFeatureFromBranch(branch string) string {
	// Common branch prefixes
	prefixes := []string{"feature/", "feat/", "fix/", "bugfix/", "hotfix/", "chore/"}

	for _, prefix := range prefixes {
		if _, after, found := strings.Cut(branch, prefix); found {
			return after
		}
	}

	// If no prefix, use the branch name as-is
	// Replace slashes with dashes for feature slug
	return strings.ReplaceAll(branch, "/", "-")
}

// IsGitRepository checks if the directory is a git repository.
func IsGitRepository(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = dir
	return cmd.Run() == nil
}

// GetFeatureFromGit attempts to get feature from git context.
// Priority: worktree name -> branch name
func GetFeatureFromGit(projectRoot string) string {
	// Try worktree name first (more specific)
	if worktree := GetWorktreeName(projectRoot); worktree != "" {
		return worktree
	}

	// Try branch name
	branch := GetCurrentBranch(projectRoot)
	if branch == "" || branch == "main" || branch == "master" || branch == "HEAD" {
		return ""
	}

	return ExtractFeatureFromBranch(branch)
}

// IsInsideWorktree returns true if the directory is a linked git worktree
// (i.e., .git is a file pointing to the main repository's worktrees directory).
func IsInsideWorktree(projectRoot string) bool {
	gitPath := filepath.Join(projectRoot, ".git")
	info, err := os.Lstat(gitPath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Push pushes the current branch to origin with upstream tracking set.
// Uses "git push -u origin HEAD" pattern.
func Push(projectRoot string) (string, error) {
	out, err := Run(projectRoot, "push", "-u", "origin", "HEAD")
	if err != nil {
		return "", fmt.Errorf("git push: %w", err)
	}
	return out, nil
}

// Run executes a git command and returns the output.
func Run(projectRoot string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = projectRoot
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}
