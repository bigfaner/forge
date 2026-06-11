package worktree

import (
	"fmt"
	"os"
	"path/filepath"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/git"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push [<slug>]",
	Short: "Push the current worktree branch to remote",
	Long: `Push a worktree's branch to origin with upstream tracking set.

When run inside a worktree (no arguments), pushes the current worktree's branch.
When given a <slug>, pushes the named worktree's branch from any directory.

Uses "git push -u origin HEAD" to push and set the upstream tracking reference.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runWorktreePush,
}

func runWorktreePush(cmd *cobra.Command, args []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return base.ErrProjectNotFound()
	}

	if !git.IsGitRepository(projectRoot) {
		return base.ErrNotGitRepository(projectRoot)
	}

	var workDir string

	if len(args) == 1 {
		// Slug provided: resolve to worktree directory
		slug := args[0]
		workDir, err = resolveWorktreeDir(projectRoot, slug)
		if err != nil {
			return err
		}
	} else {
		// No slug: use current directory (must be inside a worktree)
		if !isInsideWorktreeFunc(projectRoot) {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "error: not inside a worktree")
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "hint: run this command from within a forge worktree directory")
			return base.ErrNotInsideWorktree()
		}
		workDir = projectRoot
	}

	// Hard Rule: refuse to push from main worktree's main branch
	branch := getCurrentBranchFunc(workDir)
	if branch == "main" || branch == "master" {
		return base.ErrRefusingDefaultBranch(branch)
	}

	// Push with upstream tracking
	output, err := gitPushFunc(workDir)
	if err != nil {
		return base.NewAIError(base.ErrConflict, fmt.Sprintf("Push failed: %v", err), "Git push command failed", "Check remote connectivity and branch status", "git push -u origin HEAD")
	}

	if output != "" {
		_, _ = fmt.Fprint(cmd.OutOrStdout(), output)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Pushed branch %q to origin\n", branch)
	return nil
}

// resolveWorktreeDir resolves a slug to its worktree directory path.
// Validates that the directory exists and is a valid linked worktree.
func resolveWorktreeDir(projectRoot, slug string) (string, error) {
	targetDir := filepath.Join(projectRoot, ".forge", "worktrees", slug)
	targetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return "", base.NewAIError(base.ErrInvalidInput, fmt.Sprintf("Unable to resolve worktree path: %v", err), "Failed to resolve the worktree directory path", "Check that the slug is valid", "forge worktree list")
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return "", base.NewAIError(base.ErrNotFound, fmt.Sprintf("Worktree not found: %s", slug), "No worktree exists with that slug", "Verify the slug is correct", "forge worktree list")
	}

	// Evaluate symlinks so the path matches os.Getwd() on macOS (/var → /private/var).
	targetDir, err = filepath.EvalSymlinks(targetDir)
	if err != nil {
		return "", base.NewAIError(base.ErrInvalidInput, fmt.Sprintf("Unable to resolve worktree path: %v", err), "Failed to resolve the worktree directory path", "Check that the path is valid", "forge worktree list")
	}

	return targetDir, nil
}
