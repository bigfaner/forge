package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"task-cli/pkg/profile"
	"task-cli/pkg/project"

	"github.com/spf13/cobra"
)

var profileJSON bool

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Resolve or set the active test profile",
	Long: `Resolve or set the active test profile.

Without arguments: resolves the active profile(s) from .forge/config.yaml,
falling back to project structure detection.
Use subcommands to set or detect profiles.`,
	Args: cobra.NoArgs,
	Run:  runProfileResolve,
}

var profileSetCmd = &cobra.Command{
	Use:   "set <name>",
	Short: "Set the active test profile in .forge/config.yaml",
	Args:  cobra.ExactArgs(1),
	Run:   runProfileSet,
}

var profileDetectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect test profiles from project structure (ignores config)",
	Args:  cobra.NoArgs,
	Run:   runProfileDetect,
}

func init() {
	profileCmd.Flags().BoolVar(&profileJSON, "json", false, "output as JSON")
	profileCmd.AddCommand(profileSetCmd)
	profileCmd.AddCommand(profileDetectCmd)
}

// profileResult holds the resolved profile info.
type profileResult struct {
	Profiles []string `json:"profiles"`
	Source   string   `json:"source"`
}

func runProfileResolve(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	// 1. Try config
	configured, err := profile.ReadTestProfiles(projectRoot)
	if err != nil {
		Exit(NewAIError(ErrValidation, "Failed to read config", err.Error(), "Check .forge/config.yaml format", "task profile detect"))
	}
	if len(configured) > 0 {
		printProfileResult(profileResult{Profiles: configured, Source: "config"})
		return
	}

	// 2. Try detection
	detected, err := profile.DetectProfiles(projectRoot)
	if err != nil {
		Exit(NewAIError(ErrValidation, "Detection failed", err.Error(), "Run task profile set <name> manually", "task profile set web-playwright"))
	}
	if len(detected) > 0 {
		printProfileResult(profileResult{Profiles: detected, Source: "detected"})
		return
	}

	// 3. No match — signal AI to ask user
	printProfileResult(profileResult{Profiles: nil, Source: ""})
}

func runProfileSet(_ *cobra.Command, args []string) {
	name := args[0]

	if !profile.IsKnownProfile(name) {
		Exit(NewAIError(
			ErrInvalidInput,
			fmt.Sprintf("Unknown profile: %s", name),
			"Profile name is not in the known profiles list",
			fmt.Sprintf("Choose from: %s", strings.Join(profile.KnownProfiles, ", ")),
			fmt.Sprintf("task profile set %s", profile.KnownProfiles[0]),
		))
	}

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	if err := profile.WriteTestProfiles(projectRoot, []string{name}); err != nil {
		Exit(NewAIError(ErrValidation, "Failed to write config", err.Error(), "Check .forge/ directory permissions", "ls -la .forge/"))
	}

	PrintBlockStart()
	PrintField("PROFILE", name)
	PrintField("SOURCE", "config")
	PrintField("ACTION", "written to .forge/config.yaml")
	PrintBlockEnd()
}

func runProfileDetect(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	detected, err := profile.DetectProfiles(projectRoot)
	if err != nil {
		Exit(NewAIError(ErrValidation, "Detection failed", err.Error(), "Run task profile set <name> manually", "task profile set web-playwright"))
	}

	printProfileResult(profileResult{Profiles: detected, Source: "detected"})
}

func printProfileResult(r profileResult) {
	if profileJSON {
		data, _ := json.Marshal(r)
		fmt.Println(string(data))
		return
	}

	PrintBlockStart()
	if len(r.Profiles) == 0 {
		PrintField("PROFILE", "(none)")
		fmt.Fprintln(os.Stderr, "HINT: No profile detected. Ask user to choose and run: task profile set <name>")
	} else {
		for _, p := range r.Profiles {
			PrintField("PROFILE", p)
		}
		if r.Source != "" {
			PrintField("SOURCE", r.Source)
		}
	}
	PrintBlockEnd()
}
