package worktree

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/git"
	"forge-cli/pkg/project"
	"forge-cli/pkg/proposal"
	"forge-cli/pkg/types"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// ---------------------------------------------------------------------------
// Overridable function variables (for testing)
// ---------------------------------------------------------------------------

// listWorktreesFunc lists git worktrees for a project root.
// Overridable for testing.
var listWorktreesFunc = git.ListWorktrees

// gitRunFunc executes a git command. Overridable for testing.
var gitRunFunc = git.Run

// countUnpushedCommitsFunc counts unpushed commits. Overridable for testing.
var countUnpushedCommitsFunc = git.CountUnpushedCommits

// gitPushFunc pushes to remote. Overridable for testing.
var gitPushFunc = git.Push

// isInsideWorktreeFunc checks if inside a linked worktree. Overridable for testing.
var isInsideWorktreeFunc = git.IsInsideWorktree

// getCurrentBranchFunc returns the current branch name. Overridable for testing.
var getCurrentBranchFunc = git.GetCurrentBranch

// lookPathFunc resolves a binary name to its full path.
// Overridable for testing.
var lookPathFunc = exec.LookPath

// runClaudeFunc executes claude with the given args.
// Overridable for testing.
var runClaudeFunc = base.RunClaude

// claudeSupportsContinueFlagFunc checks whether the installed claude CLI
// supports the -c / --continue flag. Overridable for testing.
var claudeSupportsContinueFlagFunc = defaultClaudeSupportsContinueFlag

// stdinFunc is the function used to read from stdin. Overridable for testing.
var stdinFunc = defaultStdinRead

// isTerminalFunc checks if stdin is connected to a terminal (TTY).
// Overridable for testing.
var isTerminalFunc = defaultIsTerminal

// ---------------------------------------------------------------------------
// Default implementations
// ---------------------------------------------------------------------------

func defaultClaudeSupportsContinueFlag() bool {
	cmd := exec.Command("claude", "--help")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	helpText := string(output)
	return strings.Contains(helpText, "-c,") || strings.Contains(helpText, "--continue")
}

func defaultStdinRead() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func defaultIsTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// ---------------------------------------------------------------------------
// Flag registration (init)
// ---------------------------------------------------------------------------

func init() {
	startCmd.Flags().StringP("source-branch", "b", "", "source branch for the new worktree (default: HEAD)")
	startCmd.Flags().Bool("no-launch", false, "create worktree without launching claude")
	startCmd.Flags().BoolP("interactive", "i", false, "interactively select a proposal or feature")

	removeCmd.Flags().Bool("hard", false, "delete worktree, local branch, and prune stale administrative files")
	removeCmd.Flags().Bool("force", false, "force removal even with uncommitted changes (use with --hard)")

	// Shell completion functions
	startCmd.ValidArgsFunction = worktreeStartCompletion
	removeCmd.ValidArgsFunction = worktreeRemoveCompletion
	resumeCmd.ValidArgsFunction = worktreeResumeCompletion
	pushCmd.ValidArgsFunction = worktreePushCompletion
}

// ---------------------------------------------------------------------------
// File operation helpers (copy-files for worktree start)
// ---------------------------------------------------------------------------

// validateCopyFilePath checks that a single copy-file path is safe.
// Rejects absolute paths and paths containing ".." traversals.
func validateCopyFilePath(relPath string) error {
	if filepath.IsAbs(relPath) {
		return base.NewAIError(base.ErrInvalidPath, fmt.Sprintf("Copy-file path must be relative: %s", relPath), "Absolute paths are not allowed for copy-files", "Use a relative path", "forge config set worktree.copy-files [path]")
	}
	// Reject Windows-style absolute paths (e.g. C:\Windows) on all platforms.
	if len(relPath) >= 2 && relPath[1] == ':' && (relPath[2] == '\\' || relPath[2] == '/') {
		return base.NewAIError(base.ErrInvalidPath, fmt.Sprintf("Copy-file path must be relative: %s", relPath), "Absolute paths are not allowed for copy-files", "Use a relative path", "forge config set worktree.copy-files [path]")
	}
	if strings.Contains(relPath, "..") {
		return base.NewAIError(base.ErrInvalidPath, fmt.Sprintf("Copy-file path must not contain \"..\": %s", relPath), "Path traversal is not allowed", "Use a simple relative path without ..", relPath)
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
			return base.NewAIError(base.ErrNotFound, fmt.Sprintf("Copy-file not found: %s", relPath), "The specified copy-file does not exist in the project root", "Verify the file exists", "ls "+relPath)
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
			return base.NewAIError(base.ErrNotFound, fmt.Sprintf("Failed to read %s: %v", relPath, err), "Could not read the source file", "Check file exists and is readable", "cat "+relPath)
		}

		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return base.NewAIError(base.ErrInvalidInput, fmt.Sprintf("Failed to create directory for %s: %v", relPath, err), "Could not create parent directory", "Check destination permissions", "mkdir -p "+relPath)
		}

		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return base.NewAIError(base.ErrInvalidInput, fmt.Sprintf("Failed to write %s: %v", relPath, err), "Could not write the destination file", "Check destination permissions", "ls -la "+relPath)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// TUI / interactive selection helpers
// ---------------------------------------------------------------------------

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
			if strings.EqualFold(p.Status, string(types.StatusCompleted)) {
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

			if strings.EqualFold(status, string(types.StatusCompleted)) {
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

// worktreePushCompletion returns existing non-main worktree slugs for shell completion.
// Hard Rule: return empty list on error (never return error to shell).
func worktreePushCompletion(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
