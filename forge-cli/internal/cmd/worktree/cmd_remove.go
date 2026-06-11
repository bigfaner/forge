package worktree

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/git"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <slug>",
	Short: "Remove a git worktree while preserving its branch",
	Long: `Remove the git worktree at .forge/worktrees/<slug>.

The branch is preserved after removal so you can merge it later with
'git merge <slug>'. Fails if the worktree has uncommitted changes —
commit or stash first.

Use --hard to also delete the local branch and prune stale administrative
files. Without --hard, only the worktree directory is removed.

Use --force with --hard to force deletion even when the worktree has
uncommitted changes or the branch is not fully merged.`,
	Args: cobra.ExactArgs(1),
	RunE: runWorktreeRemove,
}

func runWorktreeRemove(cmd *cobra.Command, args []string) error {
	slug := args[0]

	if slug == "" {
		return base.ErrSlugRequired()
	}

	hard, _ := cmd.Flags().GetBool("hard")
	force, _ := cmd.Flags().GetBool("force")

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return base.ErrProjectNotFound()
	}

	if !git.IsGitRepository(projectRoot) {
		return base.ErrNotGitRepository(projectRoot)
	}

	// Resolve worktree path
	targetDir := filepath.Join(projectRoot, ".forge", "worktrees", slug)
	targetDir, err = filepath.Abs(targetDir)
	if err != nil {
		return base.NewAIError(base.ErrInvalidInput, fmt.Sprintf("Unable to resolve target path: %v", err), "Failed to resolve the worktree directory path", "Check that the path is valid", "forge worktree list")
	}

	// Check that the worktree directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return base.NewAIError(base.ErrNotFound, fmt.Sprintf("Worktree not found: %s", targetDir), "The worktree directory does not exist", "Verify the slug is correct", "forge worktree list")
	}

	// Look up the branch name before removal
	branchName := slug
	entries, err := listWorktreesFunc(projectRoot)
	if err == nil {
		for _, entry := range entries {
			if entry.Name() == slug && entry.Branch != "" {
				branchName = entry.Branch
				break
			}
		}
	}

	// Check for unpushed commits before removal (unless --force)
	unpushedCount, unpushedErr := countUnpushedCommitsFunc(targetDir)
	if !errors.Is(unpushedErr, git.ErrNoUpstream) && unpushedErr != nil {
		// Unexpected error — report but don't block removal
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "warning: could not check unpushed commits: %v\n", unpushedErr)
	}
	if unpushedErr == nil && unpushedCount > 0 && !force {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error: branch has %d unpushed commit(s) — push first, or use --force to discard\n", unpushedCount)
		return fmt.Errorf("branch has %d unpushed commit(s)", unpushedCount)
	}

	// Build git worktree remove args
	removeArgs := []string{"worktree", "remove"}
	if force {
		removeArgs = append(removeArgs, "--force")
	}
	removeArgs = append(removeArgs, targetDir)

	// Use git worktree remove
	_, err = git.Run(projectRoot, removeArgs...)
	if err != nil {
		errMsg := err.Error()

		// Check if error is due to uncommitted changes
		if strings.Contains(errMsg, "dirty") || strings.Contains(errMsg, "modified") ||
			strings.Contains(errMsg, "local changes") {
			if hard && !force {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error: uncommitted changes in worktree — use --force to discard\n")
				return fmt.Errorf("uncommitted changes in worktree: %w", err)
			}
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error: worktree has uncommitted changes\n")
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "hint: commit or stash your changes before removing, or use --force to discard\n")
			return fmt.Errorf("uncommitted changes in worktree: %w", err)
		}

		// Fallback: handle corrupted worktrees (e.g. missing .git file)
		// git worktree remove fails with ".git does not exist" or "not a valid working tree"
		// In this case, manually remove the directory and prune stale administrative files.
		if strings.Contains(errMsg, ".git") || strings.Contains(errMsg, "validation failed") ||
			strings.Contains(errMsg, "not a valid") || strings.Contains(errMsg, "could not identify") {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "warning: git worktree remove failed, falling back to manual cleanup: %v\n", err)
			if removeErr := os.RemoveAll(targetDir); removeErr != nil {
				return base.NewAIError(base.ErrConflict, fmt.Sprintf("Failed to remove worktree directory: %v", removeErr), "Manual worktree directory removal failed", "Check directory permissions and retry", "forge worktree list")
			}
			// Prune stale worktree administrative files so git no longer tracks it
			_, _ = git.Run(projectRoot, "worktree", "prune")
		} else {
			return base.NewAIError(base.ErrConflict, fmt.Sprintf("Failed to remove worktree: %v", err), "Git worktree remove command failed", "Check git status and retry", "forge worktree list")
		}
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Removed worktree %q\n", slug)

	// --hard: also delete branch and prune
	if !hard {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Branch %s preserved\n", branchName)
		return nil
	}

	return runHardCleanup(cmd, projectRoot, branchName)
}

// runHardCleanup performs branch deletion and worktree pruning after worktree removal.
func runHardCleanup(cmd *cobra.Command, projectRoot, branchName string) error {
	// Delete local branch (only local — never remote)
	branchDeleted := false
	if branchName != "" {
		// Try safe delete first (git branch -d)
		_, err := git.Run(projectRoot, "branch", "-d", branchName)
		if err != nil {
			// Check if the error is about unmerged changes
			errMsg := err.Error()
			if strings.Contains(errMsg, "not fully merged") || strings.Contains(errMsg, "unmerged") {
				// --hard without --force: warn but still allow deletion per Hard Rules
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Warning: branch %q is not fully merged\n", branchName)
				_, err = git.Run(projectRoot, "branch", "-D", branchName)
				if err != nil {
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Skipped branch deletion: %v\n", err)
				} else {
					branchDeleted = true
				}
			} else {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Skipped branch deletion: %v\n", err)
			}
		} else {
			branchDeleted = true
		}
	} else {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Skipped branch deletion: branch name unknown\n")
	}

	if branchDeleted {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Deleted branch %q\n", branchName)
	}

	// Prune stale worktree administrative files
	_, _ = git.Run(projectRoot, "worktree", "prune")
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Pruned stale worktree administrative files\n")

	return nil
}
