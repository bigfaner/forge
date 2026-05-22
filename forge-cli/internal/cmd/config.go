package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage forge configuration",
	Long:  `Manage .forge/config.yaml for project settings like auto-behavior and worktree.`,
	Args:  cobra.NoArgs,
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

Collects auto-behavior and worktree settings through interactive prompts.`,
	Args: cobra.NoArgs,
	RunE: runConfigInitHuh,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Long: `Set a config value in .forge/config.yaml.

	Supports dot-notation keys for nested config (e.g. "auto.gitPush true").
	Returns an error for unknown keys or invalid values.`,
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runConfigSet,
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]
	projectRoot := resolveProjectRoot(cmd)

	if err := forgeconfig.SetConfigValue(projectRoot, key, value); err != nil {
		return err
	}

	return nil
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
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

// runConfigInitHuh delegates to the huh TUI interactive config init path.
func runConfigInitHuh(cmd *cobra.Command, _ []string) error {
	projectRoot := resolveProjectRoot(cmd)
	action := configInitFunc(projectRoot)
	switch action.status {
	case "CREATED":
		write(cmd.OutOrStdout(), "Config written to .forge/config.yaml (%s)\n", action.detail)
	case "SKIPPED":
		write(cmd.OutOrStdout(), "Config init skipped: %s\n", action.detail)
	case "CANCELLED":
		write(cmd.OutOrStdout(), "Config init cancelled.\n")
	case "FAILED":
		return fmt.Errorf("config init failed: %s", action.detail)
	}
	return nil
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
