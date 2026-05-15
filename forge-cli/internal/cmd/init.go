package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/internal/embedded"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/just"
	"forge-cli/pkg/profile"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize forge project environment",
	Long: `One-stop initialization for forge project.

Creates .forge/ directory, generates CLAUDE.md from embedded template,
appends runtime entries to .gitignore, appends claude/claude-c recipes
to justfile, and runs interactive config if .forge/config.yaml doesn't exist.`,
	Args: cobra.NoArgs,
	RunE: runInit,
}

func init() {
	initCmd.Flags().String("project-root", "", "project root directory (defaults to current directory)")
	initCmd.Flags().Bool("skip-just", false, "skip just installation check")
}

// gitignoreEntries are the lines to append to .gitignore.
var gitignoreEntries = []string{
	"# Forge runtime",
	"docs/features/*/tasks/process/",
	".forge/state.json",
	"tests/results/.last-run.json",
	"tests/e2e/results/.last-run.json",
	"tests/e2e/results/*/error-context.md",
}

// justfileRecipes are the recipes to append to justfile.
var justfileRecipes = []struct {
	name string
	body string
}{
	{"claude", "    claude --dangerously-skip-permissions"},
	{"claude-c", "    claude --dangerously-skip-permissions -c"},
}

// initAction records a single action taken during init.
type initAction struct {
	status string // CREATED, APPENDED, INSTALLED, SKIPPED, FAILED
	target string // file or directory path
	detail string // extra info (e.g., "5 entries", "from template")
}

func runInit(cmd *cobra.Command, _ []string) error {
	projectRoot, _ := cmd.Flags().GetString("project-root")
	if projectRoot == "" {
		projectRoot = "."
	}

	skipJust, _ := cmd.Flags().GetBool("skip-just")

	out := cmd.OutOrStdout()
	var actions []initAction

	// Step 1: Create .forge/ directory
	action := createForgeDir(projectRoot)
	actions = append(actions, action)

	// Step 2: Generate CLAUDE.md
	action = createCLAUDEmd(projectRoot)
	actions = append(actions, action)

	// Step 3: Update .gitignore
	action = updateGitignore(projectRoot)
	actions = append(actions, action)

	// Step 3.5: Ensure just is installed (before justfile update)
	action = ensureJustStep(skipJust, cmd.InOrStdin(), out)
	actions = append(actions, action)

	// Step 4: Update justfile
	action = updateJustfile(projectRoot)
	actions = append(actions, action)

	// Step 5: Interactive config (only if config doesn't exist)
	action = runConfigInitIfNeeded(projectRoot, cmd.InOrStdin(), out)
	actions = append(actions, action)

	// Print summary report
	printInitSummary(out, actions)

	return nil
}

func createForgeDir(projectRoot string) initAction {
	forgeDir := filepath.Join(projectRoot, feature.ForgeDir)
	if _, err := os.Stat(forgeDir); err == nil {
		return initAction{status: "SKIPPED", target: feature.ForgeDir, detail: "already exists"}
	}
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to create %s: %v\n", forgeDir, err)
		return initAction{status: "FAILED", target: feature.ForgeDir, detail: err.Error()}
	}
	return initAction{status: "CREATED", target: feature.ForgeDir}
}

func createCLAUDEmd(projectRoot string) initAction {
	claudePath := filepath.Join(projectRoot, "CLAUDE.md")
	if _, err := os.Stat(claudePath); err == nil {
		return initAction{status: "SKIPPED", target: "CLAUDE.md", detail: "already exists"}
	}
	if err := os.WriteFile(claudePath, []byte(embedded.CLAUDEmdTemplate), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to create CLAUDE.md: %v\n", err)
		return initAction{status: "FAILED", target: "CLAUDE.md", detail: err.Error()}
	}
	return initAction{status: "CREATED", target: "CLAUDE.md", detail: "from template"}
}

func updateGitignore(projectRoot string) initAction {
	gitignorePath := filepath.Join(projectRoot, ".gitignore")

	var existingContent string
	data, err := os.ReadFile(gitignorePath)
	if err == nil {
		existingContent = string(data)
	}

	toAppend := buildGitignoreAppend(existingContent, gitignoreEntries)
	if len(toAppend) == 0 {
		return initAction{status: "SKIPPED", target: ".gitignore", detail: "all entries already present"}
	}

	// Build content to append
	var buf strings.Builder
	if existingContent != "" && !strings.HasSuffix(existingContent, "\n") {
		buf.WriteByte('\n')
	}
	for _, line := range toAppend {
		buf.WriteString(line)
		buf.WriteByte('\n')
	}

	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to update .gitignore: %v\n", err)
		return initAction{status: "FAILED", target: ".gitignore", detail: err.Error()}
	}
	defer func() { _ = f.Close() }()

	if _, err := f.WriteString(buf.String()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to write .gitignore: %v\n", err)
		return initAction{status: "FAILED", target: ".gitignore", detail: err.Error()}
	}

	return initAction{status: "APPENDED", target: ".gitignore", detail: fmt.Sprintf("%d entries", len(toAppend))}
}

// buildGitignoreAppend returns only the lines that are not already present.
func buildGitignoreAppend(existing string, entries []string) []string {
	existingLines := make(map[string]bool)
	for _, line := range strings.Split(existing, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			existingLines[trimmed] = true
		}
	}

	var toAppend []string
	for _, entry := range entries {
		if !existingLines[strings.TrimSpace(entry)] {
			toAppend = append(toAppend, entry)
		}
	}
	return toAppend
}

func updateJustfile(projectRoot string) initAction {
	justfilePath := filepath.Join(projectRoot, "justfile")

	var existingContent string
	data, err := os.ReadFile(justfilePath)
	if err == nil {
		existingContent = string(data)
	}

	toAppend := buildJustfileAppend(existingContent)
	if len(toAppend) == 0 {
		return initAction{status: "SKIPPED", target: "justfile", detail: "all recipes already present"}
	}

	// Build content to append
	var buf strings.Builder
	if existingContent != "" && !strings.HasSuffix(existingContent, "\n") {
		buf.WriteByte('\n')
	}
	buf.WriteString(toAppend)

	f, err := os.OpenFile(justfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to update justfile: %v\n", err)
		return initAction{status: "FAILED", target: "justfile", detail: err.Error()}
	}
	defer func() { _ = f.Close() }()

	if _, err := f.WriteString(buf.String()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to write justfile: %v\n", err)
		return initAction{status: "FAILED", target: "justfile", detail: err.Error()}
	}

	// Count recipes added
	recipeNames := collectAddedRecipeNames(existingContent)
	return initAction{status: "APPENDED", target: "justfile", detail: fmt.Sprintf("%d recipes: %s", len(recipeNames), strings.Join(recipeNames, ", "))}
}

// buildJustfileAppend returns the text to append, excluding recipes that already exist.
func buildJustfileAppend(existing string) string {
	var buf strings.Builder
	for _, recipe := range justfileRecipes {
		// Check if recipe name already exists as a line prefix
		if recipeExists(existing, recipe.name) {
			continue
		}
		buf.WriteString(recipe.name + ":\n")
		buf.WriteString(recipe.body + "\n\n")
	}
	if buf.Len() == 0 {
		return ""
	}
	return strings.TrimRight(buf.String(), "\n") + "\n"
}

// recipeExists checks if a recipe name already exists in justfile content.
func recipeExists(content, recipeName string) bool {
	lines := strings.Split(content, "\n")
	prefix := recipeName + ":"
	for _, line := range lines {
		// Recipe names appear at the start of a line (no leading whitespace)
		trimmed := strings.TrimLeft(line, " \t")
		// Match exact recipe name followed by ':' or ':' + space (to avoid
		// matching 'claude:' when checking for 'claude-c:').
		if trimmed == prefix || strings.HasPrefix(trimmed, prefix+" ") {
			return true
		}
	}
	return false
}

// collectAddedRecipeNames returns the names of recipes that were added.
func collectAddedRecipeNames(existing string) []string {
	var names []string
	for _, recipe := range justfileRecipes {
		if !recipeExists(existing, recipe.name) {
			names = append(names, recipe.name)
		}
	}
	return names
}

func runConfigInitIfNeeded(projectRoot string, in io.Reader, out io.Writer) initAction {
	configFile := filepath.Join(projectRoot, feature.ForgeDir, feature.ForgeConfigFileName)
	if _, err := os.Stat(configFile); err == nil {
		return initAction{status: "SKIPPED", target: ".forge/config.yaml", detail: "already exists"}
	}

	reader := bufio.NewReader(in)

	// Step 1: Project type
	write(out, "\nSelect project type:\n")
	write(out, "  1. frontend\n")
	write(out, "  2. backend\n")
	write(out, "  3. mixed\n")
	write(out, "Enter number [2]: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	projectType := "backend"
	switch input {
	case "1":
		projectType = "frontend"
	case "2", "":
		projectType = "backend"
	case "3":
		projectType = "mixed"
	}

	// Step 2: Test profiles
	write(out, "\nSelect test profiles (enter numbers, space-separated):\n")
	for i, p := range profile.KnownProfiles {
		write(out, "  %d. %s\n", i+1, p)
	}
	write(out, "Selections: ")

	input, _ = reader.ReadString('\n')
	selectedProfiles := parseMultiSelect(input, profile.KnownProfiles)

	// Step 3: Capabilities
	var selectedCaps []string
	if len(selectedProfiles) > 0 {
		union, err := profile.UnionCapabilities(selectedProfiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: could not resolve capabilities: %v\n", err)
		} else if len(union) > 0 {
			write(out, "\nSelect capabilities from detected profiles:\n")
			for i, c := range union {
				write(out, "  %d. %s\n", i+1, c)
			}
			write(out, "Selections: ")

			input, _ = reader.ReadString('\n')
			selectedCaps = parseMultiSelect(input, union)
		}
	}

	cfg := profile.ForgeConfig{
		ProjectType:  projectType,
		TestProfiles: selectedProfiles,
		Capabilities: selectedCaps,
	}

	if err := writeConfigFile(configFile, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to write config: %v\n", err)
		return initAction{status: "FAILED", target: ".forge/config.yaml", detail: err.Error()}
	}

	write(out, "\n")
	return initAction{status: "CREATED", target: ".forge/config.yaml", detail: "interactive"}
}

func printInitSummary(out io.Writer, actions []initAction) {
	write(out, ">>>\n")
	for _, a := range actions {
		detail := ""
		if a.detail != "" {
			detail = fmt.Sprintf(" (%s)", a.detail)
		}
		write(out, "%-10s %s%s\n", a.status, a.target, detail)
	}
	write(out, "<<<\n")
}

// ensureJustFunc is the function that runs the ensure-just flow.
// Variable for testability.
var ensureJustFunc = just.EnsureJust

// ensureJustStep runs the ensure-just flow or skips it based on the flag.
// Installation failure is non-blocking — init continues with a WARNING.
func ensureJustStep(skipJust bool, in io.Reader, out io.Writer) initAction {
	if skipJust {
		return initAction{status: "SKIPPED", target: "just installation", detail: "skipped via --skip-just flag"}
	}

	result := ensureJustFunc(in, out)

	if result.Status == just.StatusFailed {
		fmt.Fprintf(os.Stderr, "WARNING: just installation failed: %s\n", result.Detail)
	}

	return ensureResultToAction(result)
}

// ensureResultToAction converts an EnsureResult to an initAction.
func ensureResultToAction(r just.EnsureResult) initAction {
	detail := r.Detail
	if r.Version != "" && r.Status == just.StatusSkipped {
		detail = fmt.Sprintf("just %s already available", r.Version)
	}
	if r.Method != "" && r.Status == just.StatusInstalled {
		detail = fmt.Sprintf("installed via %s (%s)", r.Method, r.Version)
	}
	return initAction{
		status: string(r.Status),
		target: "just installation",
		detail: detail,
	}
}
