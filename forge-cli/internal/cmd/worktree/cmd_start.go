package worktree

import (
	"fmt"
	"os"
	"path/filepath"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/git"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [slug]",
	Short: "Create a worktree and launch Claude in it",
	Long: `Create a git worktree at .forge/worktrees/<slug> with branch <slug> from HEAD,
then launch claude --dangerously-skip-permissions in the worktree directory.

If branch <slug> already exists, creates the worktree from that branch
(resume context).

The source branch for new worktrees can be set via --source-branch / -b flag
or worktree.source-branch in .forge/config.yaml. Priority: flag > config > HEAD.

Use -i/--interactive to select a proposal or feature from a list.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runWorktreeStart,
}

func runWorktreeStart(cmd *cobra.Command, args []string) error {
	slug := ""
	if len(args) > 0 {
		slug = args[0]
	}

	// Interactive mode: prompt user to select a proposal or feature
	interactive, _ := cmd.Flags().GetBool("interactive")
	if slug == "" && interactive {
		// Find project root first (needed for scanning)
		projectRoot, err := project.FindProjectRoot()
		if err != nil {
			return base.ErrProjectNotFound()
		}

		// Check TTY
		if !isTerminalFunc() {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "error: interactive mode requires a terminal (TTY)")
			return base.NewAIError(base.ErrInvalidInput, "Interactive mode requires a terminal (TTY)", "The -i flag requires a real terminal", "Run without -i or use a TTY", "forge worktree start <slug>")
		}

		items := listUnfinishedItems(projectRoot)
		if len(items) == 0 {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No unfinished proposals or features found.")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Create one with: forge proposal <slug> or forge feature <slug>")
			return nil
		}

		selected, err := promptSelection(items, cmd.OutOrStdout())
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error: %v\n", err)
			return err
		}
		slug = selected
	}

	if slug == "" {
		return base.ErrSlugRequired()
	}

	// Check --no-launch early: skip claude pre-flight when not launching
	noLaunch, _ := cmd.Flags().GetBool("no-launch")

	// Pre-flight: verify claude binary exists in PATH (skip when --no-launch)
	if !noLaunch {
		if _, err := lookPathFunc("claude"); err != nil {
			return base.NewAIError(base.ErrNotFound, "Claude binary not found in PATH", "The claude CLI is required but not installed", "Install Claude Code or check your PATH", "pip install claude-code")
		}
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

	// Compute target directory inside project: .forge/worktrees/<slug>
	targetDir := filepath.Join(projectRoot, ".forge", "worktrees", slug)
	targetDir, err = filepath.Abs(targetDir)
	if err != nil {
		return base.NewAIError(base.ErrInvalidInput, fmt.Sprintf("Unable to resolve target path: %v", err), "Failed to resolve the worktree directory path", "Check that the path is valid", "forge worktree list")
	}

	// Ensure .forge/worktrees/ parent directory exists
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		return base.NewAIError(base.ErrInvalidInput, fmt.Sprintf("Failed to create worktrees directory: %v", err), "Could not create .forge/worktrees directory", "Check filesystem permissions", "mkdir -p .forge/worktrees")
	}

	// Check if target directory already exists
	if info, err := os.Stat(targetDir); err == nil && info.IsDir() {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error: target directory already exists: %s\n", targetDir)
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "hint: use 'forge worktree resume %s' to re-open an existing worktree\n", slug)
		return fmt.Errorf("target directory already exists: %s", targetDir)
	}

	// Load config for source-branch and copy-files
	cfg, _ := forgeconfig.ReadConfig(projectRoot)

	// Pre-validate copy-files BEFORE git worktree add (to avoid orphan worktrees)
	var copyFiles []string
	if cfg != nil && cfg.Worktree != nil {
		copyFiles = cfg.Worktree.CopyFiles
	}
	if err := validateCopyFiles(projectRoot, copyFiles); err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error: %v\n", err)
		return err
	}

	// Resolve source branch: flag > config > HEAD
	var sourceBranch string
	if cmd.Flags().Changed("source-branch") {
		sourceBranch, _ = cmd.Flags().GetString("source-branch")
	} else if cfg != nil && cfg.Worktree != nil {
		sourceBranch = cfg.Worktree.SourceBranch
	}

	// Layer 1: Check if local branch already exists
	localBranchExists := false
	if _, err := gitRunFunc(projectRoot, "rev-parse", "--verify", slug); err == nil {
		localBranchExists = true
	}

	// Layer 2: Fetch from origin (best-effort) and check remote branch
	remoteBranchExists := false
	if !localBranchExists {
		// Best-effort fetch: failure degrades gracefully (no remote, offline, etc.)
		if _, fetchErr := gitRunFunc(projectRoot, "fetch", "origin"); fetchErr != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "warning: git fetch origin failed: %v\n", fetchErr)
		}
		if _, err := gitRunFunc(projectRoot, "rev-parse", "--verify", "remotes/origin/"+slug); err == nil {
			remoteBranchExists = true
		}
	}

	// Create the worktree using three-layer resolution with branch-first approach
	switch {
	case localBranchExists:
		// Layer 1: Resume from existing local branch
		_, err = gitRunFunc(projectRoot, "worktree", "add", targetDir, slug)
	case remoteBranchExists:
		// Layer 2: Create branch from remote tracking branch, then add worktree
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "creating worktree from remote branch origin/%s\n", slug)
		_, err = gitRunFunc(projectRoot, "branch", slug, "origin/"+slug)
		if err != nil {
			return base.NewAIError(base.ErrConflict, fmt.Sprintf("Failed to create branch from remote: %v", err), "Git branch from remote failed", "Check remote connectivity", "git fetch origin")
		}
		_, err = gitRunFunc(projectRoot, "worktree", "add", targetDir, slug)
		if err != nil {
			// Cleanup: remove the branch we just created
			_, _ = gitRunFunc(projectRoot, "branch", "-D", slug)
		}
	default:
		// Layer 3: Pre-validate source branch if specified
		if sourceBranch != "" {
			if _, err := gitRunFunc(projectRoot, "rev-parse", "--verify", sourceBranch); err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error: source branch %q not found\n", sourceBranch)
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "hint: verify the branch exists locally or fetch from remote\n")
				return fmt.Errorf("source branch not found: %s", sourceBranch)
			}
		}

		// Branch-first: create branch from source, then add worktree
		branchArgs := []string{"branch", slug}
		if sourceBranch != "" {
			branchArgs = append(branchArgs, sourceBranch)
		}
		_, err = gitRunFunc(projectRoot, branchArgs...)
		if err != nil {
			return base.NewAIError(base.ErrConflict, fmt.Sprintf("Failed to create branch: %v", err), "Git branch command failed", "Check that the source branch exists", "git branch")
		}
		_, err = gitRunFunc(projectRoot, "worktree", "add", targetDir, slug)
		if err != nil {
			// Cleanup: remove the branch we just created
			_, _ = gitRunFunc(projectRoot, "branch", "-D", slug)
		}
	}
	if err != nil {
		return base.NewAIError(base.ErrConflict, fmt.Sprintf("Failed to add worktree: %v", err), "Git worktree add command failed", "Check for conflicts with existing branches", "git worktree list")
	}

	// Copy configured files from project root to worktree
	if err := copyFilesToWorktree(projectRoot, targetDir, copyFiles); err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "warning: copy-files failed: %v\n", err)
	}

	// --no-launch: print path and exit without launching claude
	if noLaunch {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "worktree created at %s\n", targetDir)
		return nil
	}

	// Launch claude in the worktree directory
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	if err := os.Chdir(targetDir); err != nil {
		return base.NewAIError(base.ErrInvalidInput, fmt.Sprintf("Failed to change to worktree directory: %v", err), "Could not change directory", "Check that the worktree path is accessible", "ls .forge/worktrees/")
	}

	allArgs := []string{"--dangerously-skip-permissions"}
	return runClaudeFunc(allArgs)
}
