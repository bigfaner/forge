package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/git"
	"forge-cli/pkg/project"
	"forge-cli/pkg/proposal"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// listWorktreesFunc lists git worktrees for a project root.
// Overridable for testing.
var listWorktreesFunc = git.ListWorktrees

// gitRunFunc executes a git command. Overridable for testing.
var gitRunFunc = git.Run

// gitPushFunc pushes to remote. Overridable for testing.
var gitPushFunc = git.Push

// isInsideWorktreeFunc checks if inside a linked worktree. Overridable for testing.
var isInsideWorktreeFunc = git.IsInsideWorktree

// getCurrentBranchFunc returns the current branch name. Overridable for testing.
var getCurrentBranchFunc = git.GetCurrentBranch

var worktreeCmd = &cobra.Command{
	Use:   "worktree",
	Short: "Manage git worktrees for feature development",
	Long: `Manage git worktrees for parallel feature development.

Each worktree is created inside the project at .forge/worktrees/<slug> with a
branch named <slug>. Forge's feature auto-detection resolves the correct
feature from the worktree name.`,
	Args: cobra.NoArgs,
}

var worktreeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all git worktrees",
	Long: `List all git worktrees with their name, branch, and path.

Worktrees whose name matches a feature slug in docs/features/ are marked as
forge-managed. The main worktree (current project) is distinguished from
feature worktrees.`,
	Args: cobra.NoArgs,
	RunE: runWorktreeList,
}

var worktreeRemoveCmd = &cobra.Command{
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
		return ErrSlugRequired()
	}

	hard, _ := cmd.Flags().GetBool("hard")
	force, _ := cmd.Flags().GetBool("force")

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	if !git.IsGitRepository(projectRoot) {
		return ErrNotGitRepository(projectRoot)
	}

	// Resolve worktree path
	targetDir := filepath.Join(projectRoot, ".forge", "worktrees", slug)
	targetDir, err = filepath.Abs(targetDir)
	if err != nil {
		return NewAIError(ErrInvalidInput, fmt.Sprintf("Unable to resolve target path: %v", err), "Failed to resolve the worktree directory path", "Check that the path is valid", "forge worktree list")
	}

	// Check that the worktree directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return NewAIError(ErrNotFound, fmt.Sprintf("Worktree not found: %s", targetDir), "The worktree directory does not exist", "Verify the slug is correct", "forge worktree list")
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

	// Build git worktree remove args
	removeArgs := []string{"worktree", "remove"}
	if force {
		removeArgs = append(removeArgs, "--force")
	}
	removeArgs = append(removeArgs, targetDir)

	// Use git worktree remove
	_, err = git.Run(projectRoot, removeArgs...)
	if err != nil {
		// Check if error is due to uncommitted changes
		errMsg := err.Error()
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
		return NewAIError(ErrConflict, fmt.Sprintf("Failed to remove worktree: %v", err), "Git worktree remove command failed", "Check git status and retry", "forge worktree list")
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

var worktreeResumeCmd = &cobra.Command{
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
		return ErrSlugRequired()
	}

	// Pre-flight: verify claude binary exists in PATH
	if _, err := lookPathFunc("claude"); err != nil {
		return NewAIError(ErrNotFound, "Claude binary not found in PATH", "The claude CLI is required but not installed", "Install Claude Code or check your PATH", "pip install claude-code")
	}

	// Find project root
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	// Verify we're in a git repository
	if !git.IsGitRepository(projectRoot) {
		return ErrNotGitRepository(projectRoot)
	}

	// Resolve worktree path
	targetDir := filepath.Join(projectRoot, ".forge", "worktrees", slug)
	targetDir, err = filepath.Abs(targetDir)
	if err != nil {
		return NewAIError(ErrInvalidInput, fmt.Sprintf("Unable to resolve target path: %v", err), "Failed to resolve the worktree directory path", "Check that the path is valid", "forge worktree list")
	}

	// Check that the worktree directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return NewAIError(ErrNotFound, fmt.Sprintf("Worktree not found: %s", targetDir), "The worktree directory does not exist", "Verify the slug is correct", "forge worktree list")
	}

	// Evaluate symlinks so the path matches os.Getwd() on macOS (/var → /private/var).
	targetDir, err = filepath.EvalSymlinks(targetDir)
	if err != nil {
		return NewAIError(ErrInvalidInput, fmt.Sprintf("Unable to resolve target path: %v", err), "Failed to resolve the worktree directory path", "Check that the path is valid", "forge worktree list")
	}

	// Verify it's a git worktree (.git file or directory must exist)
	gitFile := filepath.Join(targetDir, ".git")
	if _, err := os.Stat(gitFile); os.IsNotExist(err) {
		return NewAIError(ErrInvalidInput, fmt.Sprintf("Not a git worktree: %s", targetDir), "The directory is not a git worktree", "Ensure the slug corresponds to a valid worktree", "forge worktree list")
	}

	// Launch claude in the worktree directory
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	if err := os.Chdir(targetDir); err != nil {
		return NewAIError(ErrInvalidInput, fmt.Sprintf("Failed to change to worktree directory: %v", err), "Could not change directory", "Check that the worktree path is accessible", "ls .forge/worktrees/")
	}

	allArgs := []string{"--dangerously-skip-permissions"}

	// Add -c <slug> for session restore if supported
	if claudeSupportsContinueFlagFunc() {
		allArgs = append([]string{"-c", slug}, allArgs...)
	}

	return runClaudeFunc(allArgs)
}

var worktreePushCmd = &cobra.Command{
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
		return ErrProjectNotFound()
	}

	if !git.IsGitRepository(projectRoot) {
		return ErrNotGitRepository(projectRoot)
	}

	// Hard Rule: must detect worktree context before pushing
	if !isInsideWorktreeFunc(projectRoot) {
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "error: not inside a worktree")
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "hint: run this command from within a forge worktree directory")
		return ErrNotInsideWorktree()
	}

	// Hard Rule: refuse to push from main worktree's main branch
	branch := getCurrentBranchFunc(projectRoot)
	if branch == "main" || branch == "master" {
		return ErrRefusingDefaultBranch(branch)
	}

	// Push with upstream tracking
	output, err := gitPushFunc(projectRoot)
	if err != nil {
		return NewAIError(ErrConflict, fmt.Sprintf("Push failed: %v", err), "Git push command failed", "Check remote connectivity and branch status", "git push -u origin HEAD")
	}

	if output != "" {
		_, _ = fmt.Fprint(cmd.OutOrStdout(), output)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Pushed branch %q to origin\n", branch)
	return nil
}

var worktreeStatusCmd = &cobra.Command{
	Use:   "status [<slug>]",
	Short: "Show worktree status",
	Long: `Display the status of a forge-managed worktree: branch name, latest commit,
and uncommitted files list.

When no slug is provided, shows status for all forge-managed worktrees.
This command is strictly read-only — it never modifies any filesystem state.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runWorktreeStatus,
}

func runWorktreeStatus(cmd *cobra.Command, args []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	if !git.IsGitRepository(projectRoot) {
		return ErrNotGitRepository(projectRoot)
	}

	entries, err := listWorktreesFunc(projectRoot)
	if err != nil {
		return NewAIError(ErrNotFound, fmt.Sprintf("Failed to list worktrees: %v", err), "Could not enumerate git worktrees", "Check git status", "git worktree list")
	}

	// Build set of forge-managed feature slugs
	forgeFeatures := listForgeFeatures(projectRoot)

	if len(args) == 1 {
		// Specific slug requested
		return showWorktreeStatus(cmd, projectRoot, entries, forgeFeatures, args[0])
	}

	// No slug: show status for all forge-managed worktrees
	return showAllWorktreeStatus(cmd, projectRoot, entries, forgeFeatures)
}

// showWorktreeStatus displays status for a specific worktree slug.
func showWorktreeStatus(cmd *cobra.Command, projectRoot string, entries []git.WorktreeEntry, _ map[string]bool, slug string) error {
	// Find the worktree by slug
	var found *git.WorktreeEntry
	for i := range entries {
		if entries[i].Name() == slug {
			found = &entries[i]
			break
		}
	}

	if found == nil {
		return NewAIError(ErrNotFound, fmt.Sprintf("Worktree not found: %s", slug), "No worktree exists with that slug", "Verify the slug is correct", "forge worktree list")
	}

	printWorktreeStatus(cmd, projectRoot, found)
	return nil
}

// showAllWorktreeStatus displays status for all forge-managed worktrees.
func showAllWorktreeStatus(cmd *cobra.Command, projectRoot string, entries []git.WorktreeEntry, forgeFeatures map[string]bool) error {
	// Filter to forge-managed worktrees (non-main, with matching feature slug)
	var forgeWorktrees []git.WorktreeEntry
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsMain && forgeFeatures[name] {
			forgeWorktrees = append(forgeWorktrees, entry)
		}
	}

	if len(forgeWorktrees) == 0 {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No forge-managed worktrees found")
		return nil
	}

	for i := range forgeWorktrees {
		printWorktreeStatus(cmd, projectRoot, &forgeWorktrees[i])
	}
	return nil
}

// printWorktreeStatus prints the status of a single worktree using structured output format.
// This is a read-only operation — it only reads git state, never modifies it.
func printWorktreeStatus(cmd *cobra.Command, _ string, entry *git.WorktreeEntry) {
	worktreePath := entry.Path

	// Branch name
	branch := entry.Branch
	if branch == "" {
		branch = "(detached)"
	}

	// Latest commit: git log -1 --oneline
	commitInfo := ""
	if output, err := gitRunFunc(worktreePath, "log", "-1", "--oneline"); err == nil {
		commitInfo = output
	}

	// Uncommitted files: git status --porcelain
	var uncommittedFiles []string
	if output, err := gitRunFunc(worktreePath, "status", "--porcelain"); err == nil && output != "" {
		uncommittedFiles = strings.Split(output, "\n")
	}

	// Print structured output using PrintBlockStart/PrintField/PrintBlockEnd pattern
	w := cmd.OutOrStdout()
	_, _ = fmt.Fprintln(w, "---")
	_, _ = fmt.Fprintf(w, "WORKTREE: %s\n", entry.Name())
	_, _ = fmt.Fprintf(w, "BRANCH: %s\n", branch)
	_, _ = fmt.Fprintf(w, "COMMIT: %s\n", commitInfo)
	if len(uncommittedFiles) > 0 {
		_, _ = fmt.Fprintf(w, "UNCOMMITTED: %s\n", strings.Join(uncommittedFiles, ", "))
	} else {
		_, _ = fmt.Fprintln(w, "UNCOMMITTED: (none)")
	}
	_, _ = fmt.Fprintln(w, "---")
}

var worktreeStartCmd = &cobra.Command{
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

func init() {
	worktreeStartCmd.Flags().StringP("source-branch", "b", "", "source branch for the new worktree (default: HEAD)")
	worktreeStartCmd.Flags().Bool("no-launch", false, "create worktree without launching claude")
	worktreeStartCmd.Flags().BoolP("interactive", "i", false, "interactively select a proposal or feature")

	worktreeRemoveCmd.Flags().Bool("hard", false, "delete worktree, local branch, and prune stale administrative files")
	worktreeRemoveCmd.Flags().Bool("force", false, "force removal even with uncommitted changes (use with --hard)")

	// Shell completion functions
	worktreeStartCmd.ValidArgsFunction = worktreeStartCompletion
	worktreeRemoveCmd.ValidArgsFunction = worktreeRemoveCompletion
	worktreeResumeCmd.ValidArgsFunction = worktreeResumeCompletion

	worktreeCmd.AddCommand(worktreePushCmd)
	worktreeCmd.AddCommand(worktreeStatusCmd)
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
			return ErrProjectNotFound()
		}

		// Check TTY
		if !isTerminalFunc() {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "error: interactive mode requires a terminal (TTY)")
			return NewAIError(ErrInvalidInput, "Interactive mode requires a terminal (TTY)", "The -i flag requires a real terminal", "Run without -i or use a TTY", "forge worktree start <slug>")
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
		return ErrSlugRequired()
	}

	// Check --no-launch early: skip claude pre-flight when not launching
	noLaunch, _ := cmd.Flags().GetBool("no-launch")

	// Pre-flight: verify claude binary exists in PATH (skip when --no-launch)
	if !noLaunch {
		if _, err := lookPathFunc("claude"); err != nil {
			return NewAIError(ErrNotFound, "Claude binary not found in PATH", "The claude CLI is required but not installed", "Install Claude Code or check your PATH", "pip install claude-code")
		}
	}

	// Find project root
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	// Verify we're in a git repository
	if !git.IsGitRepository(projectRoot) {
		return ErrNotGitRepository(projectRoot)
	}

	// Compute target directory inside project: .forge/worktrees/<slug>
	targetDir := filepath.Join(projectRoot, ".forge", "worktrees", slug)
	targetDir, err = filepath.Abs(targetDir)
	if err != nil {
		return NewAIError(ErrInvalidInput, fmt.Sprintf("Unable to resolve target path: %v", err), "Failed to resolve the worktree directory path", "Check that the path is valid", "forge worktree list")
	}

	// Ensure .forge/worktrees/ parent directory exists
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		return NewAIError(ErrInvalidInput, fmt.Sprintf("Failed to create worktrees directory: %v", err), "Could not create .forge/worktrees directory", "Check filesystem permissions", "mkdir -p .forge/worktrees")
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
			return NewAIError(ErrConflict, fmt.Sprintf("Failed to create branch from remote: %v", err), "Git branch from remote failed", "Check remote connectivity", "git fetch origin")
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
			return NewAIError(ErrConflict, fmt.Sprintf("Failed to create branch: %v", err), "Git branch command failed", "Check that the source branch exists", "git branch")
		}
		_, err = gitRunFunc(projectRoot, "worktree", "add", targetDir, slug)
		if err != nil {
			// Cleanup: remove the branch we just created
			_, _ = gitRunFunc(projectRoot, "branch", "-D", slug)
		}
	}
	if err != nil {
		return NewAIError(ErrConflict, fmt.Sprintf("Failed to add worktree: %v", err), "Git worktree add command failed", "Check for conflicts with existing branches", "git worktree list")
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
		return NewAIError(ErrInvalidInput, fmt.Sprintf("Failed to change to worktree directory: %v", err), "Could not change directory", "Check that the worktree path is accessible", "ls .forge/worktrees/")
	}

	allArgs := []string{"--dangerously-skip-permissions"}
	return runClaudeFunc(allArgs)
}

func runWorktreeList(cmd *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	if !git.IsGitRepository(projectRoot) {
		return ErrNotGitRepository(projectRoot)
	}

	entries, err := listWorktreesFunc(projectRoot)
	if err != nil {
		return NewAIError(ErrNotFound, fmt.Sprintf("Failed to list worktrees: %v", err), "Could not enumerate git worktrees", "Check git status", "git worktree list")
	}

	if len(entries) == 0 {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No worktrees found")
		return nil
	}

	// Build set of forge-managed feature slugs
	forgeFeatures := listForgeFeatures(projectRoot)

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	for _, entry := range entries {
		name := entry.Name()
		branch := entry.Branch
		if branch == "" {
			branch = "(detached)"
		}

		suffix := ""
		if entry.IsMain {
			suffix = "  [main]"
		} else if forgeFeatures[name] {
			suffix = "  [forge]"
		}

		_, _ = fmt.Fprintf(w, "%s\t%s\t%s%s\n", name, branch, entry.Path, suffix)
	}
	return w.Flush()
}

// listForgeFeatures returns a set of feature slugs that exist under docs/features/.
func listForgeFeatures(projectRoot string) map[string]bool {
	featuresDir := filepath.Join(projectRoot, feature.FeaturesDir)
	entries, err := os.ReadDir(featuresDir)
	if err != nil {
		return nil
	}

	result := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() {
			result[entry.Name()] = true
		}
	}
	return result
}

// validateCopyFilePath checks that a single copy-file path is safe.
// Rejects absolute paths and paths containing ".." traversals.
func validateCopyFilePath(relPath string) error {
	if filepath.IsAbs(relPath) {
		return NewAIError(ErrInvalidPath, fmt.Sprintf("Copy-file path must be relative: %s", relPath), "Absolute paths are not allowed for copy-files", "Use a relative path", "forge config set worktree.copy-files [path]")
	}
	// Reject Windows-style absolute paths (e.g. C:\Windows) on all platforms.
	if len(relPath) >= 2 && relPath[1] == ':' && (relPath[2] == '\\' || relPath[2] == '/') {
		return NewAIError(ErrInvalidPath, fmt.Sprintf("Copy-file path must be relative: %s", relPath), "Absolute paths are not allowed for copy-files", "Use a relative path", "forge config set worktree.copy-files [path]")
	}
	if strings.Contains(relPath, "..") {
		return NewAIError(ErrInvalidPath, fmt.Sprintf("Copy-file path must not contain \"..\": %s", relPath), "Path traversal is not allowed", "Use a simple relative path without ..", relPath)
	}
	return nil
}

// validateCopyFiles pre-validates that all copy-files exist in the project root
// and have safe paths. Returns an error describing the first problem found.
// Returns nil if copyFiles is empty or nil.
func validateCopyFiles(projectRoot string, copyFiles []string) error {
	for _, relPath := range copyFiles {
		if err := validateCopyFilePath(relPath); err != nil {
			return err
		}
		fullPath := filepath.Join(projectRoot, relPath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return NewAIError(ErrNotFound, fmt.Sprintf("Copy-file not found: %s", relPath), "The specified copy-file does not exist in the project root", "Verify the file exists", "ls "+relPath)
		}
	}
	return nil
}

// copyFilesToWorktree copies the listed files from projectRoot to worktreeDir.
// Creates parent directories as needed. Overwrites existing files.
// Returns nil if copyFiles is empty or nil.
func copyFilesToWorktree(projectRoot, worktreeDir string, copyFiles []string) error {
	for _, relPath := range copyFiles {
		src := filepath.Join(projectRoot, relPath)
		dst := filepath.Join(worktreeDir, relPath)

		data, err := os.ReadFile(src)
		if err != nil {
			return NewAIError(ErrNotFound, fmt.Sprintf("Failed to read %s: %v", relPath, err), "Could not read the source file", "Check file exists and is readable", "cat "+relPath)
		}

		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return NewAIError(ErrInvalidInput, fmt.Sprintf("Failed to create directory for %s: %v", relPath, err), "Could not create parent directory", "Check destination permissions", "mkdir -p "+relPath)
		}

		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return NewAIError(ErrInvalidInput, fmt.Sprintf("Failed to write %s: %v", relPath, err), "Could not write the destination file", "Check destination permissions", "ls -la "+relPath)
		}
	}
	return nil
}

// selectableItem represents a proposal or feature that can be selected interactively.
type selectableItem struct {
	Slug   string // directory name (used as slug)
	Type   string // "proposal" or "feature"
	Status string // frontmatter status, e.g. "Draft", "in_progress", "completed"
}

// listUnfinishedItems scans docs/proposals/ and docs/features/ for unfinished work.
//
// Proposals are considered unfinished if their status is not "completed".
// Features are considered unfinished if their manifest status is not "completed".
// Proposals appear before features in the returned list.
func listUnfinishedItems(projectRoot string) []selectableItem {
	var items []selectableItem

	// Scan proposals: any proposal directory with status != "completed"
	proposals, err := proposal.Discover(projectRoot)
	if err == nil {
		for _, p := range proposals {
			if strings.EqualFold(p.Status, "completed") {
				continue
			}
			status := p.Status
			if status == "" {
				status = "Draft"
			}
			items = append(items, selectableItem{
				Slug:   p.Slug,
				Type:   "proposal",
				Status: status,
			})
		}
	}

	// Scan features: any feature directory with manifest status != "completed"
	featuresDir := filepath.Join(projectRoot, feature.FeaturesDir)
	entries, err := os.ReadDir(featuresDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			slug := entry.Name()

			// Skip if already listed as a proposal
			if containsSlug(items, slug) {
				continue
			}

			manifestPath := filepath.Join(featuresDir, slug, feature.ManifestFileName)
			status := readManifestStatus(manifestPath)

			if strings.EqualFold(status, "completed") {
				continue
			}
			if status == "" {
				status = "active"
			}

			items = append(items, selectableItem{
				Slug:   slug,
				Type:   "feature",
				Status: status,
			})
		}
	}

	return items
}

// containsSlug checks if a selectableItem with the given slug exists.
func containsSlug(items []selectableItem, slug string) bool {
	for _, item := range items {
		if item.Slug == slug {
			return true
		}
	}
	return false
}

// readManifestStatus reads the status field from a manifest.md frontmatter.
// Returns empty string if file doesn't exist or frontmatter can't be parsed.
func readManifestStatus(manifestPath string) string {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return ""
	}

	var meta struct {
		Status string `yaml:"status"`
	}
	if err := parseFrontmatter(data, &meta); err != nil {
		return ""
	}
	return meta.Status
}

// parseFrontmatter extracts YAML frontmatter from markdown content.
func parseFrontmatter(content []byte, target any) error {
	text := string(content)
	if !strings.HasPrefix(text, "---") {
		return nil
	}
	text = text[3:]
	closeIdx := strings.Index(text, "\n---")
	if closeIdx < 0 {
		return nil
	}
	yamlContent := text[:closeIdx]
	return yaml.Unmarshal([]byte(yamlContent), target)
}

// stdinFunc is the function used to read from stdin. Overridable for testing.
var stdinFunc = defaultStdinRead

// defaultStdinRead reads a line from stdin.
func defaultStdinRead() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// isTerminalFunc checks if stdin is connected to a terminal (TTY).
// Overridable for testing.
var isTerminalFunc = defaultIsTerminal

func defaultIsTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// promptSelection presents a numbered list and reads the user's choice.
// Returns the selected slug, or empty string if selection is invalid.
func promptSelection(items []selectableItem, stdout io.Writer) (string, error) {
	_, _ = fmt.Fprintln(stdout, "Select a proposal or feature:")
	_, _ = fmt.Fprintln(stdout)
	for i, item := range items {
		_, _ = fmt.Fprintf(stdout, "  %d. [%s] %s (%s)\n", i+1, item.Type, item.Slug, item.Status)
	}
	_, _ = fmt.Fprintln(stdout)
	_, _ = fmt.Fprint(stdout, "Enter number: ")

	line, err := stdinFunc()
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}

	num, err := strconv.Atoi(line)
	if err != nil {
		return "", fmt.Errorf("invalid selection: %q", line)
	}

	if num < 1 || num > len(items) {
		return "", fmt.Errorf("selection %d out of range (1-%d)", num, len(items))
	}

	return items[num-1].Slug, nil
}

// ---------------------------------------------------------------------------
// Shell completion functions
// ---------------------------------------------------------------------------

// worktreeStartCompletion returns unfinished proposal/feature slugs for shell completion.
// Hard Rule: return empty list on error (never return error to shell).
func worktreeStartCompletion(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// If a slug arg is already provided, no more completion needed
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	items := listUnfinishedItems(projectRoot)

	var completions []string
	for _, item := range items {
		if strings.HasPrefix(item.Slug, toComplete) {
			completions = append(completions, item.Slug+"\t"+item.Type)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// worktreeRemoveCompletion returns existing non-main worktree slugs for shell completion.
// Hard Rule: return empty list on error (never return error to shell).
func worktreeRemoveCompletion(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return worktreeSlugCompletion(args, toComplete)
}

// worktreeResumeCompletion returns existing non-main worktree slugs for shell completion.
// Hard Rule: return empty list on error (never return error to shell).
func worktreeResumeCompletion(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return worktreeSlugCompletion(args, toComplete)
}

// worktreeSlugCompletion returns non-main worktree slugs filtered by toComplete prefix.
func worktreeSlugCompletion(args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	entries, err := listWorktreesFunc(projectRoot)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	for _, entry := range entries {
		if entry.IsMain {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, toComplete) {
			completions = append(completions, name)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}
