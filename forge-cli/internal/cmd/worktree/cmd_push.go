package worktree

import (
	"fmt"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/git"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push the current worktree branch to remote",
	Long: `Push the current worktree's branch to origin with upstream tracking set.

Must be run inside a linked worktree (not the main worktree on its default branch).
Uses "git push -u origin HEAD" to push and set the upstream tracking reference.`,
	Args: cobra.NoArgs,
	RunE: runWorktreePush,
}

func runWorktreePush(cmd *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return base.ErrProjectNotFound()
	}

	if !git.IsGitRepository(projectRoot) {
		return base.ErrNotGitRepository(projectRoot)
	}

	// Hard Rule: must detect worktree context before pushing
	if !isInsideWorktreeFunc(projectRoot) {
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "error: not inside a worktree")
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "hint: run this command from within a forge worktree directory")
		return base.ErrNotInsideWorktree()
	}

	// Hard Rule: refuse to push from main worktree's main branch
	branch := getCurrentBranchFunc(projectRoot)
	if branch == "main" || branch == "master" {
		return base.ErrRefusingDefaultBranch(branch)
	}

	// Push with upstream tracking
	output, err := gitPushFunc(projectRoot)
	if err != nil {
		return base.NewAIError(base.ErrConflict, fmt.Sprintf("Push failed: %v", err), "Git push command failed", "Check remote connectivity and branch status", "git push -u origin HEAD")
	}

	if output != "" {
		_, _ = fmt.Fprint(cmd.OutOrStdout(), output)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Pushed branch %q to origin\n", branch)
	return nil
}
