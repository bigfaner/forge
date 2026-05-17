package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"forge-cli/pkg/git"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var worktreeCmd = &cobra.Command{
	Use:   "worktree",
	Short: "Manage git worktrees for feature development",
	Long: `Manage git worktrees for parallel feature development.

Each worktree is created as a sibling directory (../<slug>) with a branch
named <slug>. Forge's feature auto-detection resolves the correct feature
from the worktree name.`,
}

var worktreeStartCmd = &cobra.Command{
	Use:   "start <slug>",
	Short: "Create a worktree and launch Claude in it",
	Long: `Create a git worktree at ../<slug> with branch <slug> from HEAD,
then launch claude --dangerously-skip-permissions in the worktree directory.

If branch <slug> already exists, creates the worktree from that branch
(resume context).`,
	Args: cobra.ExactArgs(1),
	RunE: runWorktreeStart,
}

func runWorktreeStart(cmd *cobra.Command, args []string) error {
	slug := args[0]

	if slug == "" {
		return fmt.Errorf("slug must not be empty")
	}

	// Pre-flight: verify claude binary exists in PATH
	if _, err := lookPathFunc("claude"); err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error: claude binary not found in PATH\n")
		return fmt.Errorf("claude: %w", err)
	}

	// Find project root
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return fmt.Errorf("find project root: %w", err)
	}

	// Verify we're in a git repository
	if !git.IsGitRepository(projectRoot) {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error: not a git repository: %s\n", projectRoot)
		return fmt.Errorf("not a git repository: %s", projectRoot)
	}

	// Compute target directory as sibling: filepath.Join(projectRoot, "..", slug)
	targetDir := filepath.Join(projectRoot, "..", slug)
	targetDir, err = filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("resolve target path: %w", err)
	}

	// Check if target directory already exists
	if info, err := os.Stat(targetDir); err == nil && info.IsDir() {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error: target directory already exists: %s\n", targetDir)
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "hint: use 'forge worktree resume %s' to re-open an existing worktree\n", slug)
		return fmt.Errorf("target directory already exists: %s", targetDir)
	}

	// Check if branch already exists
	branchExists := false
	if _, err := git.Run(projectRoot, "rev-parse", "--verify", slug); err == nil {
		branchExists = true
	}

	// Create the worktree
	if branchExists {
		// Resume: create worktree from existing branch
		_, err = git.Run(projectRoot, "worktree", "add", targetDir, slug)
	} else {
		// New: create worktree with new branch from HEAD
		_, err = git.Run(projectRoot, "worktree", "add", "-b", slug, targetDir)
	}
	if err != nil {
		return fmt.Errorf("git worktree add: %w", err)
	}

	// Launch claude in the worktree directory
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	if err := os.Chdir(targetDir); err != nil {
		return fmt.Errorf("change to worktree directory: %w", err)
	}

	allArgs := []string{"--dangerously-skip-permissions"}
	return runClaudeFunc(allArgs)
}
