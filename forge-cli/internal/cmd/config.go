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
	"forge-cli/pkg/profile"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage forge configuration",
	Long:  `Manage .forge/config.yaml for project settings like project-type, test-profiles, and capabilities.`,
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

Collects project-type, test-profiles, and capabilities through stdin prompts.`,
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

	value, err := profile.GetConfigValue(projectRoot, key)
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
	write(out, "\nSelect test profiles (enter numbers, space-separated, then 'done'):\n")
	for i, p := range profile.KnownProfiles {
		write(out, "  %d. %s\n", i+1, p)
	}
	write(out, "Selections: ")

	input, _ = reader.ReadString('\n')
	selectedProfiles := parseMultiSelect(input, profile.KnownProfiles)

	// Step 3: Capabilities
	var availableCaps []string
	if len(selectedProfiles) > 0 {
		union, err := profile.UnionCapabilities(selectedProfiles)
		if err != nil {
			return fmt.Errorf("resolve capabilities: %w", err)
		}
		availableCaps = union
	}

	var selectedCaps []string
	if len(availableCaps) > 0 {
		write(out, "\nSelect capabilities from detected profiles:\n")
		for i, c := range availableCaps {
			write(out, "  %d. %s\n", i+1, c)
		}
		write(out, "Selections: ")

		input, _ = reader.ReadString('\n')
		selectedCaps = parseMultiSelect(input, availableCaps)
	}

	// Write config
	cfg := profile.ForgeConfig{
		ProjectType:  projectType,
		TestProfiles: selectedProfiles,
		Capabilities: selectedCaps,
	}

	if err := writeConfigFile(configFile, &cfg); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	write(out, "\nConfig written to %s\n", configFile)
	return nil
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

// writeConfigFile writes ForgeConfig to the given path, creating parent dirs.
func writeConfigFile(path string, cfg *profile.ForgeConfig) error {
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
