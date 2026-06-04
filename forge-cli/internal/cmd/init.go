package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/just"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize forge project environment",
	Long: `One-stop initialization for forge project.

Creates .forge/ directory, appends runtime entries to .gitignore,
ensures just is installed, and runs interactive config if
.forge/config.yaml doesn't exist.`,
	Args: cobra.NoArgs,
	RunE: runInit,
}

func init() {
	initCmd.Flags().String("project-root", "", "project root directory (defaults to current directory)")
	initCmd.Flags().Bool("skip-just", false, "skip just installation check")
}

// gitignoreEntries are the lines to append to .gitignore.
var gitignoreEntries = []string{
	"# Forge",
	".forge/state.json",
	".forge/test-state.json",
	".forge/worktrees/",
	".forge/logs/",
	"docs/features/*/tasks/process/",
	"docs/features/*/tasks/index.json.lock",
	"docs/features/*/testing/results/",
	feature.TestResultsDir + "/",
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

	// Step 2: Update .gitignore
	action = updateGitignore(projectRoot)
	actions = append(actions, action)

	// Step 4: Ensure just is installed
	action = ensureJustStep(skipJust, cmd.InOrStdin(), out)
	actions = append(actions, action)

	// Step 5: Interactive config (only if config doesn't exist)
	action = configInitFunc(projectRoot)
	actions = append(actions, action)

	// Step 6: Surface detection and TUI confirmation
	action = surfaceConfigFunc(projectRoot)
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

// surfaceConfigFunc is the function that runs surface detection and TUI confirmation.
// Variable for testability — huh requires a real TTY, so tests override this.
var surfaceConfigFunc = runSurfaceConfig

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

// runSurfaceConfig runs surface detection and TUI confirmation.
// Skipped in non-interactive mode or when config doesn't exist.
// Re-run behavior: when config.yaml already has surfaces configured, prompts
// the user with Confirm / Re-detect / Edit options.
func runSurfaceConfig(projectRoot string) initAction {
	configFile := filepath.Join(projectRoot, feature.ForgeDir, feature.ForgeConfigFileName)

	// Only run surface config if config file exists (created by step 5)
	if _, err := os.Stat(configFile); err != nil {
		return initAction{status: "SKIPPED", target: "surfaces", detail: "no config file"}
	}

	// Check for interactive terminal
	fi, _ := os.Stdin.Stat()
	if fi.Mode()&os.ModeCharDevice == 0 {
		return initAction{status: "SKIPPED", target: "surfaces", detail: "non-interactive terminal"}
	}

	// Check if surfaces are already configured (re-run behavior)
	existingCfg, err := forgeconfig.ReadConfig(projectRoot)
	if err != nil {
		return initAction{status: "FAILED", target: "surfaces", detail: err.Error()}
	}

	if existingCfg != nil && len(existingCfg.Surfaces) > 0 {
		// Re-run flow: ask user what to do
		return handleRerunSurfaceConfig(projectRoot, configFile, existingCfg)
	}

	// First-run flow: run detection + TUI confirmation
	return runNewSurfaceDetection(projectRoot, configFile)
}

// handleRerunSurfaceConfig handles the re-run flow when surfaces already exist in config.
func handleRerunSurfaceConfig(projectRoot, configFile string, existingCfg *forgeconfig.Config) initAction {
	action, cancelled := askRerunPrompt(existingCfg.Surfaces)
	if cancelled {
		return initAction{status: "CANCELLED", target: "surfaces", detail: "Ctrl+C"}
	}

	switch action {
	case "confirm":
		// Keep existing surfaces
		return initAction{status: "SKIPPED", target: "surfaces", detail: "already configured"}
	case "edit":
		// Hard Rule: Edit calls the same manualSurfaceEntry function as first-run
		surfaces, cancelled := manualSurfaceEntry()
		if cancelled {
			return initAction{status: "CANCELLED", target: "surfaces", detail: "Ctrl+C"}
		}
		if len(surfaces) == 0 {
			return initAction{status: "SKIPPED", target: "surfaces", detail: "no surfaces entered"}
		}
		existingCfg.Surfaces = surfaces
		if err := writeConfigFile(configFile, existingCfg); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: failed to write surfaces: %v\n", err)
			return initAction{status: "FAILED", target: "surfaces", detail: err.Error()}
		}
		return initAction{status: "CREATED", target: "surfaces", detail: formatSurfacesSummary(surfaces, nil)}
	default:
		// "redetect" — run full detection + inference pipeline
		return runNewSurfaceDetection(projectRoot, configFile)
	}
}

// runNewSurfaceDetection runs detection and TUI confirmation, then writes results.
// runNewSurfaceDetection is the function variable for running new surface detection.
// Variable for testability.
var runNewSurfaceDetection = runNewSurfaceDetectionImpl

func runNewSurfaceDetectionImpl(projectRoot, configFile string) initAction {
	// Run TUI confirmation (detection + display + user interaction)
	surfaces, sources, cancelled := askSurfaceConfirmation(projectRoot)
	if cancelled {
		return initAction{status: "CANCELLED", target: "surfaces", detail: "Ctrl+C"}
	}
	if len(surfaces) == 0 {
		return initAction{status: "SKIPPED", target: "surfaces", detail: "no surfaces detected"}
	}

	// Write surfaces to config (source annotations are display-only, not persisted)
	return writeSurfacesToConfig(configFile, surfaces, sources)
}

// writeSurfacesToConfig writes surfaces to the config file.
// Source annotations are NOT persisted — only surface types are written.
// Sources are used only for the detail string in the init summary.
func writeSurfacesToConfig(configFile string, surfaces forgeconfig.SurfacesMap, sources forgeconfig.SourcesMap) initAction {
	// Read config from the directory containing the config file
	projectRoot := filepath.Dir(filepath.Dir(configFile))
	cfg, err := forgeconfig.ReadConfig(projectRoot)
	if err != nil || cfg == nil {
		cfg = &forgeconfig.Config{}
	}

	cfg.Surfaces = surfaces
	if err := writeConfigFile(configFile, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to write surfaces: %v\n", err)
		return initAction{status: "FAILED", target: "surfaces", detail: err.Error()}
	}

	return initAction{status: "CREATED", target: "surfaces", detail: formatSurfacesSummary(surfaces, sources)}
}
