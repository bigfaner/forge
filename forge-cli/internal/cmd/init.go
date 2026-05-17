package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
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
appends runtime entries to .gitignore, ensures just is installed, and
runs interactive config if .forge/config.yaml doesn't exist.`,
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
	"docs/features/*/tasks/index.json.lock",
	".forge/state.json",
	"tests/results/.last-run.json",
	"tests/e2e/results/.last-run.json",
	"tests/e2e/results/*/error-context.md",
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

	// Step 4: Ensure just is installed
	action = ensureJustStep(skipJust, cmd.InOrStdin(), out)
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

func runConfigInitIfNeeded(projectRoot string, in io.Reader, out io.Writer) initAction {
	configFile := filepath.Join(projectRoot, feature.ForgeDir, feature.ForgeConfigFileName)
	if _, err := os.Stat(configFile); err == nil {
		return initAction{status: "SKIPPED", target: ".forge/config.yaml", detail: "already exists"}
	}

	reader := bufio.NewReader(in)

	// Step 1: Project type (single select)
	projectType := askSelect(reader, out, "Project type", []string{"frontend", "backend", "mixed"}, 2)

	// Step 2: Test profile (single select)
	profileOpts := profile.KnownProfiles
	selectedProfile := askSelect(reader, out, "Test profile", profileOpts, 0)
	var selectedProfiles []string
	if selectedProfile != "" {
		selectedProfiles = []string{selectedProfile}
	}

	// Step 3: Capabilities (multi-select)
	var selectedCaps []string
	if len(selectedProfiles) > 0 {
		union, err := profile.UnionCapabilities(selectedProfiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: could not resolve capabilities: %v\n", err)
		} else if len(union) > 0 {
			selectedCaps = askMultiSelect(reader, out, "Capabilities", union, nil)
		}
	}

	// Step 4: Auto-behavior config
	auto := askAutoBehavior(reader, out, len(selectedProfiles) > 0)

	cfg := profile.ForgeConfig{
		ProjectType:  projectType,
		TestProfiles: selectedProfiles,
		Capabilities: selectedCaps,
	}
	if auto != nil {
		cfg.Auto = auto
	}

	if err := writeConfigFile(configFile, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to write config: %v\n", err)
		return initAction{status: "FAILED", target: ".forge/config.yaml", detail: err.Error()}
	}

	write(out, "\n")
	return initAction{status: "CREATED", target: ".forge/config.yaml", detail: "interactive"}
}

// askSelect presents a numbered single-select menu. Returns the selected option.
// defaultIdx is 0-based; -1 means no default.
func askSelect(reader *bufio.Reader, out io.Writer, label string, options []string, defaultIdx int) string {
	write(out, "\n%s:\n", label)
	for i, opt := range options {
		marker := "  "
		if i == defaultIdx {
			marker = "> "
		}
		write(out, "%s%d. %s\n", marker, i+1, opt)
	}
	if defaultIdx >= 0 {
		write(out, "Select [%d]: ", defaultIdx+1)
	} else {
		write(out, "Select (0 to skip): ")
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		if defaultIdx >= 0 && defaultIdx < len(options) {
			return options[defaultIdx]
		}
		return ""
	}

	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > len(options) {
		if defaultIdx >= 0 && defaultIdx < len(options) {
			return options[defaultIdx]
		}
		return ""
	}
	return options[n-1]
}

// askMultiSelect presents a numbered multi-select menu. Returns selected options.
// defaults are option values that start pre-selected.
func askMultiSelect(reader *bufio.Reader, out io.Writer, label string, options []string, defaults []string) []string {
	defSet := make(map[string]bool)
	for _, d := range defaults {
		defSet[d] = true
	}

	write(out, "\n%s (enter numbers, space-separated):\n", label)
	for i, opt := range options {
		marker := "  "
		if defSet[opt] {
			marker = "> "
		}
		write(out, "%s%d. %s\n", marker, i+1, opt)
	}
	if len(defaults) > 0 {
		write(out, "Select [%s]: ", strings.Join(defaults, " "))
	} else {
		write(out, "Select (0 to skip): ")
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaults
	}

	return parseMultiSelect(input, options)
}

// askAutoBehavior runs the auto-behavior config step with select menus.
// Returns nil if user accepts all defaults.
func askAutoBehavior(reader *bufio.Reader, out io.Writer, hasProfiles bool) *profile.AutoConfig {
	write(out, "\nAuto-behavior configuration:\n")

	var auto profile.AutoConfig
	changed := false

	if hasProfiles {
		auto.E2eTest = askToggleSelect(reader, out, "e2eTest", profile.ModeToggle{Quick: true, Full: true}, &changed)
		auto.ConsolidateSpecs = askToggleSelect(reader, out, "consolidateSpecs", profile.ModeToggle{Quick: true, Full: true}, &changed)
	}

	auto.CleanCode = askToggleSelect(reader, out, "cleanCode", profile.ModeToggle{Quick: false, Full: false}, &changed)

	if askYesNo(reader, out, "Auto git push after completion", false) {
		auto.GitPush = true
		changed = true
	}

	if !changed {
		return nil
	}
	return &auto
}

// askToggleSelect presents a numbered menu for a mode-scoped toggle.
func askToggleSelect(reader *bufio.Reader, out io.Writer, name string, defaults profile.ModeToggle, changed *bool) profile.ModeToggle {
	opts := []string{
		"enabled (quick + full)",
		"quick only",
		"full only",
		"disabled",
	}
	defIdx := 0
	if defaults.Quick && defaults.Full {
		defIdx = 0
	} else if defaults.Quick && !defaults.Full {
		defIdx = 1
	} else if !defaults.Quick && defaults.Full {
		defIdx = 2
	} else {
		defIdx = 3
	}

	write(out, "  %s:\n", name)
	for i, opt := range opts {
		marker := "    "
		if i == defIdx {
			marker = "  > "
		}
		write(out, "%s%d. %s\n", marker, i+1, opt)
	}
	write(out, "  Select [%d]: ", defIdx+1)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaults
	}

	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > 4 {
		return defaults
	}

	*changed = true
	switch n {
	case 1:
		return profile.ModeToggle{Quick: true, Full: true}
	case 2:
		return profile.ModeToggle{Quick: true, Full: false}
	case 3:
		return profile.ModeToggle{Quick: false, Full: true}
	default:
		return profile.ModeToggle{Quick: false, Full: false}
	}
}

// askYesNo presents a yes/no confirmation.
func askYesNo(reader *bufio.Reader, out io.Writer, label string, defVal bool) bool {
	prompt := "y/N"
	if defVal {
		prompt = "Y/n"
	}
	write(out, "  %s? [%s]: ", label, prompt)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return defVal
	}
	return input == "y" || input == "yes"
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
