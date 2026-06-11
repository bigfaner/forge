package worktree

import (
	"errors"
	"fmt"
	"strings"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/git"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [<slug>]",
	Short: "Show worktree status",
	Long: `Display the status of a forge-managed worktree.

For each worktree, shows:
  WORKTREE     — feature slug
  BRANCH       — current branch name (or "(detached)")
  COMMIT       — latest commit (git log -1 --oneline)
  UNCOMMITTED  — list of uncommitted files (or "(none)")
  UNPUSHED     — count of commits not yet pushed to remote, or "no remote"

When no slug is provided, shows status for all forge-managed worktrees.
This command is strictly read-only — it never modifies any filesystem state.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runWorktreeStatus,
}

func runWorktreeStatus(cmd *cobra.Command, args []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return base.ErrProjectNotFound()
	}

	if !git.IsGitRepository(projectRoot) {
		return base.ErrNotGitRepository(projectRoot)
	}

	entries, err := listWorktreesFunc(projectRoot)
	if err != nil {
		return base.NewAIError(base.ErrNotFound, fmt.Sprintf("Failed to list worktrees: %v", err), "Could not enumerate git worktrees", "Check git status", "git worktree list")
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
		return base.NewAIError(base.ErrNotFound, fmt.Sprintf("Worktree not found: %s", slug), "No worktree exists with that slug", "Verify the slug is correct", "forge worktree list")
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

	// Unpushed commits
	unpushedStr := formatUnpushed(worktreePath)
	_, _ = fmt.Fprintf(w, "UNPUSHED: %s\n", unpushedStr)

	_, _ = fmt.Fprintln(w, "---")
}

// formatUnpushed returns a human-readable string for the unpushed commit count.
func formatUnpushed(worktreePath string) string {
	count, err := countUnpushedCommitsFunc(worktreePath)
	if errors.Is(err, git.ErrNoUpstream) {
		return "no remote"
	}
	if err != nil {
		return "(unknown)"
	}
	if count == 0 {
		return "(none)"
	}
	if count == 1 {
		return "1 commit"
	}
	return fmt.Sprintf("%d commits", count)
}
