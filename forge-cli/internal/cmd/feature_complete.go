package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"
	"forge-cli/pkg/testrunner"

	"github.com/spf13/cobra"
)

// completionResult holds context for post-completion actions.
type completionResult struct {
	FeatureSlug string
	ProjectRoot string
	QuickMode   bool // true when proposal.md exists in feature directory
	ManifestRel string
	ProposalRel string // empty when full pipeline (no proposal.md)
}

var featureCompleteCmd = &cobra.Command{
	Use:   "complete",
	Short: "Complete feature lifecycle after all tasks are done",
	Long: `Post-completion hook that updates status files and commits them.

Run with --if-done to guard: exits 0 silently when not all tasks
are completed or skipped. Updates manifest.md (and proposal.md in
quick mode), commits, and optionally pushes to remote.

When uncommitted post-loop artifacts are detected (feature workspace
files, decisions, lessons, conventions, business-rules), outputs a
block decision so the agent can commit them via /git-commit. The
CompletedAt guard ensures artifact detection runs at most once.

This command is designed as a Stop hook — it always exits 0 so it
never blocks the agent stop flow.`,
	Args: cobra.NoArgs,
	RunE: runFeatureCompleteCmd,
}

var ifDone bool

func init() {
	featureCompleteCmd.Flags().BoolVar(&ifDone, "if-done", false, "only act when all tasks are completed or skipped")
	featureCmd.AddCommand(featureCompleteCmd)
}

// runFeatureCompleteCmd is the cobra Run function for the complete subcommand.
// Always exits 0 (hook protocol: non-blocking).
func runFeatureCompleteCmd(_ *cobra.Command, _ []string) error {
	if !ifDone {
		// Without --if-done flag, this is a no-op (safety guard).
		return nil
	}

	result := checkFeatureCompletion()
	if result == nil {
		return nil
	}

	if err := completeFeature(result); err != nil {
		fmt.Fprintf(os.Stderr, "[feature:complete] Error: %v\n", err)
	}
	return nil
}

// checkFeatureCompletion verifies all tasks are done and returns context for completion.
// Returns nil when conditions are not met — caller should exit silently.
func checkFeatureCompletion() *completionResult {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return nil
	}

	featureSlug, err := feature.GetCurrentFeature(projectRoot)
	if err != nil {
		return nil
	}

	// Already completed by a previous hook run — skip entirely.
	if state := feature.ReadForgeState(projectRoot); state != nil && state.CompletedAt != "" {
		return nil
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return nil
	}

	// All tasks must be completed or skipped (rejected does not count as done).
	for _, t := range index.TasksMap() {
		if t.Status != feature.StatusCompleted && t.Status != feature.StatusSkipped {
			return nil
		}
	}

	// Detect pipeline mode: proposal.md in feature directory = quick mode.
	featureDir := filepath.Join(projectRoot, feature.FeaturesDir, featureSlug)
	proposalPath := filepath.Join(featureDir, feature.ProposalFileName)
	_, proposalErr := os.Stat(proposalPath)
	quickMode := proposalErr == nil

	manifestRel := filepath.Join(feature.FeaturesDir, featureSlug, feature.ManifestFileName)
	proposalRel := ""
	if quickMode {
		proposalRel = filepath.Join(feature.FeaturesDir, featureSlug, feature.ProposalFileName)
	}

	return &completionResult{
		FeatureSlug: featureSlug,
		ProjectRoot: projectRoot,
		QuickMode:   quickMode,
		ManifestRel: manifestRel,
		ProposalRel: proposalRel,
	}
}

// completeFeature performs the status update, commit, and optional push.
func completeFeature(result *completionResult) error {
	featureDir := filepath.Join(result.ProjectRoot, feature.FeaturesDir, result.FeatureSlug)

	manifestPath := filepath.Join(featureDir, feature.ManifestFileName)

	// 1. Update manifest.md status to completed
	if err := updateFileStatus(manifestPath, "completed"); err != nil {
		return fmt.Errorf("update manifest: %w", err)
	}

	// 2. In quick mode, also update proposal.md status
	if result.QuickMode {
		proposalPath := filepath.Join(featureDir, feature.ProposalFileName)
		if err := updateFileStatus(proposalPath, "Completed"); err != nil {
			return fmt.Errorf("update proposal: %w", err)
		}
	}

	// 3. Build list of files to commit (only specific paths — never broad staging)
	var filesToCommit []string
	filesToCommit = append(filesToCommit, filepath.Join(feature.FeaturesDir, result.FeatureSlug, feature.ManifestFileName))
	if result.QuickMode {
		filesToCommit = append(filesToCommit, filepath.Join(feature.FeaturesDir, result.FeatureSlug, feature.ProposalFileName))
	}

	// 4. Git commit
	commitMsg := fmt.Sprintf("feat(%s): mark feature completed", result.FeatureSlug)
	if err := gitCommitFiles(result.ProjectRoot, filesToCommit, commitMsg); err != nil {
		return fmt.Errorf("git commit: %w", err)
	}

	commaFiles := "manifest.md"
	if result.QuickMode {
		commaFiles = "manifest.md, proposal.md"
	}
	fmt.Fprintf(os.Stderr, "[feature:complete] Status committed: %s\n", commaFiles)

	// 5. Mark state.json so future hook runs skip this feature
	if err := feature.MarkFeatureCompleted(result.ProjectRoot); err != nil {
		fmt.Fprintf(os.Stderr, "[feature:complete] Warning: failed to mark state: %v\n", err)
	}

	// 6. Optional push (before artifact detection — if artifacts block the hook,
	// the next run skips due to CompletedAt guard, so push must happen here).
	if enabled, _ := isGitPushEnabled(result.ProjectRoot); enabled {
		if err := gitPush(result.ProjectRoot); err != nil {
			_ = err
		} else {
			fmt.Fprintln(os.Stderr, "[feature:complete] Pushed to remote")
		}
	}

	// 7. Detect uncommitted post-loop artifacts and block so agent can commit them.
	// Runs at most once: CompletedAt guard in checkFeatureCompletion prevents re-entry.
	artifacts := detectUncommittedArtifacts(result.ProjectRoot, result.FeatureSlug)
	if len(artifacts) > 0 {
		var fileList strings.Builder
		for _, a := range artifacts {
			fmt.Fprintf(&fileList, "\n  - %s", a)
		}
		reason := fmt.Sprintf(
			"Post-loop artifacts detected for feature %s.%s\n\nUse /git-commit to commit these files.",
			result.FeatureSlug, fileList.String())
		testrunner.PrintHookJSON(map[string]any{
			"decision": "block",
			"reason":   reason,
		})
		return nil
	}

	return nil
}

// updateFileStatus updates the "status" field in a markdown file's YAML frontmatter.
// The file must have valid YAML frontmatter (--- delimited).
func updateFileStatus(filePath, value string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file %s: %w", filePath, err)
	}

	text := string(data)

	// Must have YAML frontmatter
	if !strings.HasPrefix(text, "---") {
		return fmt.Errorf("file %s has no YAML frontmatter", filePath)
	}

	// Find closing ---
	rest := text[3:]
	closeIdx := strings.Index(rest, "\n---")
	if closeIdx < 0 {
		return fmt.Errorf("file %s has malformed frontmatter (no closing ---)", filePath)
	}

	frontmatter := rest[:closeIdx]
	body := rest[closeIdx+4:] // after \n---

	// Find and update the status line in frontmatter
	lines := strings.Split(strings.TrimPrefix(frontmatter, "\n"), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, "status:") {
			// Idempotent: skip if already set to target value
			if strings.TrimSpace(line) == "status: "+value {
				return nil
			}
			lines[i] = fmt.Sprintf("status: %s", value)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("file %s frontmatter missing key %q", filePath, "status")
	}

	newContent := "---\n" + strings.Join(lines, "\n") + "\n---" + body
	return os.WriteFile(filePath, []byte(newContent), 0644)
}

// gitCommitFiles stages only the specified files and commits them.
// Never uses git add -A or git add . — only explicit file paths.
func gitCommitFiles(projectRoot string, files []string, message string) error {
	// Stage specific files only
	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add failed: %w\n%s", err, string(output))
	}

	// Commit
	commitArgs := []string{"commit", "-m", message}
	cmd = exec.Command("git", commitArgs...)
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit failed: %w\n%s", err, string(output))
	}

	return nil
}

// gitPush pushes to the remote. Errors are logged to stderr but not returned
// (non-blocking: exit code is always 0 per hook protocol).
func gitPush(projectRoot string) error {
	cmd := exec.Command("git", "push", "-u", "origin", "HEAD")
	cmd.Dir = projectRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[feature:complete] Push failed: %v\n%s\n", err, string(output))
		return err
	}
	return nil
}

// isGitPushEnabled reads auto.gitPush from .forge/config.yaml.
func isGitPushEnabled(projectRoot string) (bool, error) {
	auto, err := forgeconfig.ReadAutoConfig(projectRoot)
	if err != nil {
		return false, err
	}
	return auto.GitPush, nil
}

// artifactScopePaths lists directories whose uncommitted files qualify as post-loop artifacts.
var artifactScopePaths = []string{
	"docs/decisions/",
	"docs/lessons/",
	"docs/conventions/",
	"docs/business-rules/",
}

// detectUncommittedArtifacts runs git status --porcelain and filters results to
// feature-scope paths: the feature workspace and knowledge directories.
func detectUncommittedArtifacts(projectRoot, featureSlug string) []string {
	cmd := exec.Command("git", "status", "--porcelain", "-uall")
	cmd.Dir = projectRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}

	if len(strings.TrimSpace(string(output))) == 0 {
		return nil
	}

	featurePrefix := feature.FeaturesDir + "/" + featureSlug + "/"
	prefixes := append([]string{featurePrefix}, artifactScopePaths...)

	var matched []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// porcelain format: "XY path" — skip status chars (3 bytes)
		if len(line) < 4 {
			continue
		}
		path := line[3:]
		// Renames: "old -> new"
		if idx := strings.Index(path, " -> "); idx >= 0 {
			path = path[idx+4:]
		}
		// Normalize separators
		path = filepath.ToSlash(path)
		for _, prefix := range prefixes {
			if strings.HasPrefix(path, prefix) {
				matched = append(matched, path)
				break
			}
		}
	}
	return matched
}
