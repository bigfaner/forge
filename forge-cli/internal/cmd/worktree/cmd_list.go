package worktree

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/git"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all git worktrees",
	Long: `List all git worktrees with their name, branch, and path.

Worktrees whose name matches a feature slug in docs/features/ are marked as
forge-managed. The main worktree (current project) is distinguished from
feature worktrees.`,
	Args: cobra.NoArgs,
	RunE: runWorktreeList,
}

func runWorktreeList(cmd *cobra.Command, _ []string) error {
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
