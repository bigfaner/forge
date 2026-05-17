package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/internal/embedded"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/just"
	"forge-cli/pkg/profile"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
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
	status string // CREATED, APPENDED, INSTALLED, SKIPPED, FAILED, CANCELLED
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
	action = configInitFunc(projectRoot)
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

func runConfigInitIfNeeded(projectRoot string) initAction {
	configFile := filepath.Join(projectRoot, feature.ForgeDir, feature.ForgeConfigFileName)

	fi, _ := os.Stdin.Stat()
	if fi.Mode()&os.ModeCharDevice == 0 {
		return initAction{status: "SKIPPED", target: ".forge/config.yaml", detail: "non-interactive terminal"}
	}

	// When config exists, ask if user wants to reconfigure
	if _, err := os.Stat(configFile); err == nil {
		reconfigure := false
		if err := huh.NewForm(huh.NewGroup(
			huh.NewConfirm().
				Title("Config already exists. Reconfigure?").
				Description("Select Yes to overwrite .forge/config.yaml with new settings.").
				Affirmative("Yes").
				Negative("No").
				Value(&reconfigure),
		)).Run(); err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return initAction{status: "CANCELLED", target: ".forge/config.yaml", detail: "Ctrl+C"}
			}
			return initAction{status: "FAILED", target: ".forge/config.yaml", detail: err.Error()}
		}
		if !reconfigure {
			return initAction{status: "SKIPPED", target: ".forge/config.yaml", detail: "kept existing"}
		}
	}

	// Step 1: Project type
	projectType := "backend"
	if err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("What type of project is this?").
			Description("Determines scope rules for generated tasks (frontend, backend, or both).").
			Options(
				huh.NewOption("Frontend (UI-focused)", "frontend"),
				huh.NewOption("Backend (server/API)", "backend"),
				huh.NewOption("Mixed (full-stack)", "mixed"),
			).
			Value(&projectType),
	)).Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return initAction{status: "CANCELLED", target: ".forge/config.yaml", detail: "Ctrl+C"}
		}
		return initAction{status: "FAILED", target: ".forge/config.yaml", detail: err.Error()}
	}

	// Step 2: Test profile
	selectedProfile := ""
	profileOpts := make([]huh.Option[string], 0, len(profile.KnownProfiles)+1)
	profileOpts = append(profileOpts, huh.NewOption("None — skip test generation", ""))
	for _, p := range profile.KnownProfiles {
		profileOpts = append(profileOpts, huh.NewOption(p, p))
	}
	if err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Which test framework does this project use?").
			Description("Select a profile matching your test tooling, or skip test task generation.").
			Options(profileOpts...).
			Value(&selectedProfile),
	)).Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return initAction{status: "CANCELLED", target: ".forge/config.yaml", detail: "Ctrl+C"}
		}
		return initAction{status: "FAILED", target: ".forge/config.yaml", detail: err.Error()}
	}

	var selectedProfiles []string
	if selectedProfile != "" {
		selectedProfiles = []string{selectedProfile}
	}

	// Step 3: Capabilities (only if profile selected)
	var selectedCaps []string
	if len(selectedProfiles) > 0 {
		union, err := profile.UnionCapabilities(selectedProfiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: could not resolve capabilities: %v\n", err)
		} else if len(union) > 0 {
			selectedCaps = make([]string, len(union))
			copy(selectedCaps, union)
			capOpts := make([]huh.Option[string], 0, len(union))
			for _, c := range union {
				capOpts = append(capOpts, huh.NewOption(c, c).Selected(true))
			}
			if err := huh.NewForm(huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Which test capabilities should be enabled?").
					Description("Controls which test task types are generated. Space to toggle, Enter to confirm.").
					Options(capOpts...).
					Limit(0).
					Value(&selectedCaps),
			)).Run(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					return initAction{status: "CANCELLED", target: ".forge/config.yaml", detail: "Ctrl+C"}
				}
				return initAction{status: "FAILED", target: ".forge/config.yaml", detail: err.Error()}
			}
		}
	}

	// Step 4: Auto-behavior config
	auto, cancelled := askAutoBehavior(len(selectedProfiles) > 0)
	if cancelled {
		return initAction{status: "CANCELLED", target: ".forge/config.yaml", detail: "Ctrl+C"}
	}

	cfg := profile.ForgeConfig{
		ProjectType:  projectType,
		TestProfiles: selectedProfiles,
		Capabilities: selectedCaps,
		Auto:         auto,
	}

	if err := writeConfigFile(configFile, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to write config: %v\n", err)
		return initAction{status: "FAILED", target: ".forge/config.yaml", detail: err.Error()}
	}

	return initAction{status: "CREATED", target: ".forge/config.yaml", detail: "interactive"}
}

// modeHighlight styles mode keywords for terminal display.
var modeHighlight = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7DCFFF"))

// hl returns a highlighted version of text using modeHighlight.
func hl(text string) string {
	return modeHighlight.Render(text)
}

// hlMode returns "Quick mode" or "Full mode" with the whole phrase highlighted.
func hlMode(mode string) string {
	return hl(mode + " mode")
}

// askAutoBehavior runs the auto-behavior config steps, one question per screen.
// Returns the config and whether the user cancelled.
func askAutoBehavior(hasProfiles bool) (*profile.AutoConfig, bool) {
	defaults := profile.AutoConfigDefaults()
	auto := &profile.AutoConfig{}

	if hasProfiles {
		val, ok := askConfirm(
			fmt.Sprintf("%s: auto-run e2e tests?", hlMode("Quick")),
			fmt.Sprintf("Automatically run end-to-end tests during %s (lightweight verification after each task).", hl("quick mode")),
			defaults.E2eTest.Quick,
		)
		if !ok {
			return nil, true
		}
		auto.E2eTest.Quick = val

		val, ok = askConfirm(
			fmt.Sprintf("%s: auto-run e2e tests?", hlMode("Full")),
			fmt.Sprintf("Automatically run end-to-end tests during %s (comprehensive coverage).", hl("full mode")),
			defaults.E2eTest.Full,
		)
		if !ok {
			return nil, true
		}
		auto.E2eTest.Full = val

		val, ok = askConfirm(
			fmt.Sprintf("%s: auto-consolidate specs?", hlMode("Quick")),
			fmt.Sprintf("Automatically extract and consolidate specs from code after %s tasks.", hl("quick-mode")),
			defaults.ConsolidateSpecs.Quick,
		)
		if !ok {
			return nil, true
		}
		auto.ConsolidateSpecs.Quick = val

		val, ok = askConfirm(
			fmt.Sprintf("%s: auto-consolidate specs?", hlMode("Full")),
			fmt.Sprintf("Automatically extract and consolidate specs from code after %s tasks.", hl("full-mode")),
			defaults.ConsolidateSpecs.Full,
		)
		if !ok {
			return nil, true
		}
		auto.ConsolidateSpecs.Full = val
	}

	val, ok := askConfirm(
		fmt.Sprintf("%s: auto code cleanup?", hlMode("Quick")),
		fmt.Sprintf("Automatically simplify and clean code during %s.", hl("quick mode")),
		defaults.CleanCode.Quick,
	)
	if !ok {
		return nil, true
	}
	auto.CleanCode.Quick = val

	val, ok = askConfirm(
		fmt.Sprintf("%s: auto code cleanup?", hlMode("Full")),
		fmt.Sprintf("Automatically simplify and clean code during %s.", hl("full mode")),
		defaults.CleanCode.Full,
	)
	if !ok {
		return nil, true
	}
	auto.CleanCode.Full = val

	val, ok = askConfirm(
		"Auto git push after all tasks complete?",
		"Push to remote automatically when every task in a run finishes successfully.",
		defaults.GitPush,
	)
	if !ok {
		return nil, true
	}
	auto.GitPush = val

	return auto, false
}

// askConfirm shows a single confirm prompt. Returns (value, ok).
// ok is false when the user pressed Ctrl+C.
func askConfirm(title, desc string, defaultVal bool) (bool, bool) {
	val := defaultVal
	err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title(title).
			Description(desc).
			Affirmative("Yes").
			Negative("No").
			Value(&val),
	)).Run()
	if err != nil {
		return defaultVal, false
	}
	return val, true
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

// configInitFunc is the function that runs interactive config init.
// Variable for testability — huh requires a real TTY, so tests override this.
var configInitFunc = runConfigInitIfNeeded

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
