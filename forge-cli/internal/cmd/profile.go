package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"forge-cli/pkg/profile"
	"forge-cli/pkg/project"

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

var (
	profileGetManifest bool
	profileGetGenerate bool
	profileGetRun      bool
	profileGetGraduate bool
	profileGetJustfile bool
	profileGetTemplate string
)

var profileGetCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Get profile strategy file content",
	Long: `Output a profile's strategy file or template content.
Used by skills to retrieve profile data from embedded storage.

Examples:
  forge profile get go-test --manifest
  forge profile get go-test --generate
  forge profile get web-playwright --template helpers.ts`,
	Args: cobra.ExactArgs(1),
	Run:  runProfileGet,
}

func init() {
	profileCmd.Flags().BoolVar(&profileJSON, "json", false, "output as JSON")
	profileCmd.AddCommand(profileSetCmd)
	profileCmd.AddCommand(profileDetectCmd)
	profileCmd.AddCommand(profileGetCmd)

	profileGetCmd.Flags().BoolVar(&profileGetManifest, "manifest", false, "output manifest.yaml")
	profileGetCmd.Flags().BoolVar(&profileGetGenerate, "generate", false, "output generate.md strategy")
	profileGetCmd.Flags().BoolVar(&profileGetRun, "run", false, "output run.md strategy")
	profileGetCmd.Flags().BoolVar(&profileGetGraduate, "graduate", false, "output graduate.md strategy")
	profileGetCmd.Flags().BoolVar(&profileGetJustfile, "justfile", false, "output justfile-recipes")
	profileGetCmd.Flags().StringVar(&profileGetTemplate, "template", "", "output a specific template file")
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
		Exit(NewAIError(ErrValidation, "Failed to read config", err.Error(), "Check .forge/config.yaml format", "forge profile detect"))
	}
	if len(configured) > 0 {
		printProfileResult(profileResult{Profiles: configured, Source: "config"})
		return
	}

	// 2. Try detection
	detected, err := profile.DetectProfiles(projectRoot)
	if err != nil {
		Exit(NewAIError(ErrValidation, "Detection failed", err.Error(), "Run forge profile set <name> manually", "forge profile set web-playwright"))
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
			fmt.Sprintf("forge profile set %s", profile.KnownProfiles[0]),
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
		Exit(NewAIError(ErrValidation, "Detection failed", err.Error(), "Run forge profile set <name> manually", "forge profile set web-playwright"))
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
		fmt.Fprintln(os.Stderr, "HINT: No profile detected. Ask user to choose and run: forge profile set <name>")
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

func runProfileGet(_ *cobra.Command, args []string) {
	name := args[0]

	var data []byte
	var err error
	var label string

	switch {
	case profileGetManifest:
		data, err = profile.GetManifest(name)
		label = "manifest"
	case profileGetGenerate:
		data, err = profile.GetStrategy(name, "generate")
		label = "generate"
	case profileGetRun:
		data, err = profile.GetStrategy(name, "run")
		label = "run"
	case profileGetGraduate:
		data, err = profile.GetStrategy(name, "graduate")
		label = "graduate"
	case profileGetJustfile:
		data, err = profile.GetJustfileRecipes(name)
		label = "justfile"
	case profileGetTemplate != "":
		data, err = profile.GetTemplate(name, profileGetTemplate)
		label = "template:" + profileGetTemplate
	default:
		Exit(NewAIError(ErrInvalidInput, "No flag specified", "Choose one: --manifest, --generate, --run, --graduate, --justfile, --template <file>", "forge profile get go-test --generate", ""))
	}

	if err != nil {
		Exit(NewAIError(ErrInvalidInput, "Failed to get profile data", err.Error(), fmt.Sprintf("Check that %q is a valid profile name", name), "forge profile detect"))
	}

	fmt.Print(string(data))
	_ = label
}
