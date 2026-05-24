package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"

	"github.com/spf13/cobra"
)

// isInteractiveTerminalFunc is the function variable for terminal interactivity detection.
// Variable for testability — tests override this to simulate non-interactive mode.
var isInteractiveTerminalFunc = defaultIsInteractiveTerminal

// defaultIsInteractiveTerminal checks if stdin is an interactive terminal (character device).
func defaultIsInteractiveTerminal() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice != 0
}

var detectApplyFlag bool

// detectCmd is the `forge surfaces detect` subcommand.
var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect project surfaces without modifying config",
	Long: `Run surface detection and structural inference, showing results with source
annotations. Default mode is read-only: surfaces are displayed but NOT written to config.

Use --apply to enable config writing via interactive TUI confirmation.

Output format (non-interactive): one line per surface as "path=type (source)".
The non-interactive output format is UNSTABLE and may change between minor versions.
Scripted consumers should pin the forge version to avoid breakage.`,
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runDetect,
}

func init() {
	detectCmd.Flags().BoolVar(&detectApplyFlag, "apply", false, "enable TUI confirmation and config writing")
	detectCmd.Flags().String("project-root", "", "project root directory (defaults to auto-detection)")

	surfacesCmd.AddCommand(detectCmd)
}

// runDetect implements the forge surfaces detect command.
func runDetect(cmd *cobra.Command, _ []string) error {
	projectRoot := resolveProjectRoot(cmd)

	// Run detection + inference pipeline
	result, err := forgeconfig.DetectSurfacesWithConflicts(projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: surface detection failed: %v\n", err)
		return err
	}

	// Empty detection: print nothing, exit 1
	if len(result.Surfaces) == 0 {
		return &detectEmptyError{}
	}

	// Non-interactive terminal: print to stdout, no TUI, no config write (Hard Rule)
	if !isInteractiveTerminalFunc() {
		printDetectResult(cmd.OutOrStdout(), result)
		return nil
	}

	// Interactive mode without --apply: print results and exit (read-only)
	if !detectApplyFlag {
		printDetectResult(cmd.OutOrStdout(), result)
		return nil
	}

	// Interactive mode with --apply: show TUI confirmation, write on confirm
	return runDetectApply(projectRoot, result, cmd)
}

// runDetectApply handles the --apply flow: TUI confirmation + config write.
// Hard Rule: no config write without explicit --apply flag.
func runDetectApply(projectRoot string, _ *forgeconfig.DetectResult, cmd *cobra.Command) error {
	// Check if config file exists; if not, create it
	configFile := filepath.Join(projectRoot, feature.ForgeDir, feature.ForgeConfigFileName)

	// Ensure config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create .forge directory and initial config
		if err := os.MkdirAll(filepath.Dir(configFile), 0o755); err != nil {
			return fmt.Errorf("create .forge dir: %w", err)
		}
		cfg := &forgeconfig.Config{Auto: &forgeconfig.AutoConfig{}}
		if err := writeConfigFile(configFile, cfg); err != nil {
			return fmt.Errorf("create config file: %w", err)
		}
	}

	// Run TUI confirmation using askSurfaceConfirmation
	// Re-run detection through the TUI flow to get user confirmation
	surfaces, sources, cancelled := askSurfaceConfirmation(projectRoot)
	if cancelled {
		return nil
	}
	if len(surfaces) == 0 {
		return &detectEmptyError{}
	}

	// Write surfaces to config
	action := writeSurfacesToConfig(configFile, surfaces, sources)
	if action.status == "FAILED" {
		fmt.Fprintf(os.Stderr, "ERROR: failed to write surfaces: %s\n", action.detail)
		return fmt.Errorf("failed to write surfaces: %s", action.detail)
	}

	write(cmd.OutOrStdout(), "%s\n", formatSurfacesSummary(surfaces, sources))
	return nil
}

// printDetectResult prints detection results in the non-interactive stdout format.
// Format: one line per surface: <path>=<type> (<source>)
// where <source> is detected:<signal> or inferred:<rule-id>
func printDetectResult(out interface{ Write([]byte) (int, error) }, result *forgeconfig.DetectResult) {
	paths := make([]string, 0, len(result.Surfaces))
	for p := range result.Surfaces {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	for _, p := range paths {
		surfaceType := result.Surfaces[p]
		line := surfaceType
		if p != "." {
			line = fmt.Sprintf("%s=%s", p, surfaceType)
		}

		if source, ok := result.Sources[p]; ok && source != "" {
			line += fmt.Sprintf(" (%s)", formatDetectSourceAnnotation(source))
		}

		write(out, "%s\n", line)
	}
}

// formatDetectSourceAnnotation converts internal source annotation to detect stdout format.
// Internal: "dependency:cobra" -> "detected:cobra"
// Internal: "inference:cmd-dir" -> "inferred:cmd-dir"
func formatDetectSourceAnnotation(source string) string {
	if source == "" {
		return ""
	}

	parts := strings.SplitN(source, ":", 2)
	if len(parts) != 2 {
		return source
	}

	category := parts[0]
	detail := parts[1]

	switch category {
	case "dependency":
		return "detected:" + detail
	case "inference":
		return "inferred:" + detail
	default:
		return source
	}
}

// detectEmptyError is a sentinel error for empty detection results.
// Cobra SilenceErrors=true suppresses the "Error:" prefix, and we
// use this to signal exit code 1 without printing an error message.
type detectEmptyError struct{}

func (e *detectEmptyError) Error() string {
	return "no surfaces detected"
}
