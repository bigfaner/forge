package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage forge configuration",
	Long:  `Manage .forge/config.yaml for project settings like auto-behavior and worktree.`,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value (plain text output)",
	Long: `Get a config value from .forge/config.yaml.

Output is plain text: scalars print the raw value, arrays print one item per line.
Exits with code 1 if the key doesn't exist or config file is missing.`,
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runConfigGet,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactively initialize .forge/config.yaml",
	Long: `Interactively create or reconfigure .forge/config.yaml.

Collects auto-behavior and worktree settings through stdin prompts.`,
	RunE: runConfigInit,
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configInitCmd)

	configCmd.PersistentFlags().String("project-root", "", "project root directory (defaults to auto-detection)")
}

// write prints to w, ignoring write errors (interactive output is best-effort).
func write(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}

func resolveProjectRoot(cmd *cobra.Command) string {
	root, _ := cmd.Flags().GetString("project-root")
	if root != "" {
		return root
	}
	// Auto-detect project root
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return "."
	}
	return projectRoot
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]
	projectRoot := resolveProjectRoot(cmd)

	value, err := forgeconfig.GetConfigValue(projectRoot, key)
	if err != nil {
		return err
	}

	write(cmd.OutOrStdout(), "%s", value)
	return nil
}

func runConfigInit(cmd *cobra.Command, _ []string) error {
	projectRoot := resolveProjectRoot(cmd)

	configFile := filepath.Join(projectRoot, feature.ForgeDir, feature.ForgeConfigFileName)
	reader := bufio.NewReader(cmd.InOrStdin())
	out := cmd.OutOrStdout()

	// Check if config already exists
	if _, err := os.Stat(configFile); err == nil {
		write(out, "Config already exists. Reconfigure? [y/N] ")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			write(out, "Keeping existing config.\n")
			return nil
		}
	}

	// Collect auto-behavior settings
	write(out, "\nAuto-behavior settings:\n")

	write(out, "  Auto-run e2e tests in quick mode? [y/N] ")
	autoE2eQuick := readBool(reader, false)

	write(out, "  Auto-run e2e tests in full mode? [Y/n] ")
	autoE2eFull := readBool(reader, true)

	write(out, "  Auto-consolidate specs in quick mode? [Y/n] ")
	autoConsolidateQuick := readBool(reader, true)

	write(out, "  Auto-consolidate specs in full mode? [Y/n] ")
	autoConsolidateFull := readBool(reader, true)

	write(out, "  Auto code cleanup in quick mode? [y/N] ")
	autoCleanQuick := readBool(reader, false)

	write(out, "  Auto code cleanup in full mode? [y/N] ")
	autoCleanFull := readBool(reader, false)

	write(out, "  Auto git push after all tasks complete? [y/N] ")
	autoGitPush := readBool(reader, false)

	// Collect worktree settings
	write(out, "\nWorktree settings (press Enter to skip):\n")

	write(out, "  Source branch for worktrees [default: HEAD]: ")
	sourceBranch, _ := reader.ReadString('\n')
	sourceBranch = strings.TrimSpace(sourceBranch)

	write(out, "  Copy files (space-separated, e.g. '.env .env.local'): ")
	copyInput, _ := reader.ReadString('\n')
	copyFiles := parseSpaceSeparated(copyInput)

	// Build config
	cfg := forgeconfig.Config{
		Auto: &forgeconfig.AutoConfig{
			E2eTest:          forgeconfig.ModeToggle{Quick: autoE2eQuick, Full: autoE2eFull},
			ConsolidateSpecs: forgeconfig.ModeToggle{Quick: autoConsolidateQuick, Full: autoConsolidateFull},
			CleanCode:        forgeconfig.ModeToggle{Quick: autoCleanQuick, Full: autoCleanFull},
			GitPush:          autoGitPush,
		},
	}

	if sourceBranch != "" || len(copyFiles) > 0 {
		cfg.Worktree = &forgeconfig.WorktreeConfig{
			SourceBranch: sourceBranch,
			CopyFiles:    copyFiles,
		}
	}

	if err := writeConfigFile(configFile, &cfg); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	write(out, "\nConfig written to %s\n", configFile)
	return nil
}

// readBool reads a y/n answer from the reader, returning the default on empty input.
func readBool(reader *bufio.Reader, defaultVal bool) bool {
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return defaultVal
	}
	return input == "y" || input == "yes"
}

// parseSpaceSeparated splits input by whitespace, trimming and filtering empty strings.
func parseSpaceSeparated(input string) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}
	return strings.Fields(input)
}

// parseMultiSelect parses space-separated numbers from user input
// and returns the corresponding items from the options list.
func parseMultiSelect(input string, options []string) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	var selected []string
	parts := strings.Fields(input)
	for _, p := range parts {
		idx, err := strconv.Atoi(p)
		if err != nil || idx < 1 || idx > len(options) {
			continue
		}
		selected = append(selected, options[idx-1])
	}
	return selected
}

// writeConfigFile writes forgeconfig.Config to the given path, creating parent dirs.
func writeConfigFile(path string, cfg *forgeconfig.Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create .forge dir: %w", err)
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(cfg); err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	return os.WriteFile(path, buf.Bytes(), 0o644)
}
