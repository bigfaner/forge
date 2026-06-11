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

var resumeCmd = &cobra.Command{
	Use:   "resume <slug>",
	Short: "Re-launch Claude in an existing worktree",
	Long: `Launch claude with session restore (-c) and --dangerously-skip-permissions in an
existing worktree directory. If the -c flag is not supported by the installed
Claude CLI, falls back to launching without session restore.

Verifies that the worktree exists and is a valid git worktree before launching.`,
	Args: cobra.ExactArgs(1),
	RunE: runWorktreeResume,
}

func runWorktreeResume(_ *cobra.Command, args []string) error {
	slug := args[0]

	if slug == "" {
		return base.ErrSlugRequired()
	}

	// Pre-flight: verify claude binary exists in PATH
	if _, err := lookPathFunc("claude"); err != nil {
		return base.NewAIError(base.ErrNotFound, "Claude binary not found in PATH", "The claude CLI is required but not installed", "Install Claude Code or check your PATH", "pip install claude-code")
	}

	// Find project root
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return base.ErrProjectNotFound()
	}

	// Verify we're in a git repository
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

	// Evaluate symlinks so the path matches os.Getwd() on macOS (/var → /private/var).
	targetDir, err = filepath.EvalSymlinks(targetDir)
	if err != nil {
		return base.NewAIError(base.ErrInvalidInput, fmt.Sprintf("Unable to resolve target path: %v", err), "Failed to resolve the worktree directory path", "Check that the path is valid", "forge worktree list")
	}

	// Verify it's a git worktree (.git file or directory must exist)
	gitFile := filepath.Join(targetDir, ".git")
	if _, err := os.Stat(gitFile); os.IsNotExist(err) {
		return base.NewAIError(base.ErrInvalidInput, fmt.Sprintf("Not a git worktree: %s", targetDir), "The directory is not a git worktree", "Ensure the slug corresponds to a valid worktree", "forge worktree list")
	}

	// Launch claude in the worktree directory
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	if err := os.Chdir(targetDir); err != nil {
		return base.NewAIError(base.ErrInvalidInput, fmt.Sprintf("Failed to change to worktree directory: %v", err), "Could not change directory", "Check that the worktree path is accessible", "ls .forge/worktrees/")
	}

	allArgs := []string{"--dangerously-skip-permissions"}

	// Add -c for session restore if supported (no positional arg — slug would be sent as message)
	if claudeSupportsContinueFlagFunc() {
		allArgs = append([]string{"-c"}, allArgs...)
	}

	return runClaudeFunc(allArgs)
}
