package cmd

import (
	"fmt"
	"os"
	"strings"

	"forge-cli/pkg/profile"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var testingCmd = &cobra.Command{
	Use:   "testing",
	Short: "Resolve testing strategies based on project language detection",
	Long: `Resolve testing strategies based on project language detection.

Supports auto-detection from project files (go.mod, package.json, Cargo.toml, etc.)
or explicit override via .forge/config.yaml languages field.

Subcommands:
  detect              — output detected language(s)
  get generate        — output generate.md strategy
  get run             — output run.md strategy
  get graduate        — output graduate.md strategy
  get justfile        — output justfile-recipes
  get template <file> — output specified template file
  interfaces          — output interface types for the project`,
	Args: cobra.NoArgs,
	Run:  runTestingResolve,
}

var testingDetectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect languages from project structure (ignores config overrides)",
	Args:  cobra.NoArgs,
	Run:   runTestingDetect,
}

var testingGetLanguage string

var testingGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get testing strategy file content",
	Long: `Output a testing strategy file for the detected (or specified) language.

Auto-detects the project language when --language flag is not specified.
For multi-language projects, use --language to select a specific language;
without the flag, the first detected language is used.

Examples:
  forge testing get generate
  forge testing get run --language javascript
  forge testing get template test-file.go`,
}

var testingGetGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Output generate.md strategy",
	Args:  cobra.NoArgs,
	Run:   runTestingGetStrategy("generate"),
}

var testingGetRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Output run.md strategy",
	Args:  cobra.NoArgs,
	Run:   runTestingGetStrategy("run"),
}

var testingGetGraduateCmd = &cobra.Command{
	Use:   "graduate",
	Short: "Output graduate.md strategy",
	Args:  cobra.NoArgs,
	Run:   runTestingGetStrategy("graduate"),
}

var testingGetJustfileCmd = &cobra.Command{
	Use:   "justfile",
	Short: "Output justfile-recipes",
	Args:  cobra.NoArgs,
	Run:   runTestingGetJustfile,
}

var testingGetTemplateCmd = &cobra.Command{
	Use:   "template <file>",
	Short: "Output a specific template file",
	Args:  cobra.ExactArgs(1),
	Run:   runTestingGetTemplate,
}

var testingInterfacesCmd = &cobra.Command{
	Use:   "interfaces",
	Short: "Output interface types for the project",
	Long: `Output interface types for the project.

Returns config.Interfaces if set in .forge/config.yaml,
otherwise returns the union of all detected languages' default interfaces.`,
	Args: cobra.NoArgs,
	Run:  runTestingInterfaces,
}

func init() {
	testingCmd.AddCommand(testingDetectCmd)
	testingCmd.AddCommand(testingGetCmd)
	testingCmd.AddCommand(testingInterfacesCmd)

	testingGetCmd.PersistentFlags().StringVar(&testingGetLanguage, "language", "", "language key (auto-detected if omitted)")

	testingGetCmd.AddCommand(testingGetGenerateCmd)
	testingGetCmd.AddCommand(testingGetRunCmd)
	testingGetCmd.AddCommand(testingGetGraduateCmd)
	testingGetCmd.AddCommand(testingGetJustfileCmd)
	testingGetCmd.AddCommand(testingGetTemplateCmd)
}

func runTestingResolve(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	languages, err := profile.ReadLanguages(projectRoot)
	if err != nil {
		Exit(NewAIError(ErrValidation, "Failed to read languages", err.Error(), "Check .forge/config.yaml format", "forge testing detect"))
	}

	if len(languages) > 0 {
		printTestingLanguages(languages, "resolved")
		return
	}

	printTestingLanguages(nil, "")
}

func runTestingDetect(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	detected, err := profile.DetectLanguages(projectRoot)
	if err != nil {
		Exit(NewAIError(ErrValidation, "Detection failed", err.Error(), "Set languages in .forge/config.yaml manually", "forge testing detect"))
	}

	var names []string
	for _, l := range detected {
		names = append(names, string(l))
	}
	printTestingLanguages(names, "detected")
}

// resolveLanguageFromFlags resolves the language from --language flag, config, or auto-detect.
// Returns an error if no language can be resolved.
func resolveLanguageFromFlags(projectRoot string) (string, error) {
	// Explicit --language flag takes highest priority
	if testingGetLanguage != "" {
		if !profile.IsKnownLanguage(testingGetLanguage) {
			return "", NewAIError(
				ErrInvalidInput,
				fmt.Sprintf("Unknown language: %s", testingGetLanguage),
				"Language key is not in the known languages list",
				fmt.Sprintf("Choose from: %s", strings.Join(profile.KnownLanguages, ", ")),
				fmt.Sprintf("forge testing get generate --language %s", profile.KnownLanguages[0]),
			)
		}
		return testingGetLanguage, nil
	}

	// Try config override, then auto-detect
	languages, err := profile.ReadLanguages(projectRoot)
	if err != nil {
		return "", NewAIError(ErrValidation, "Failed to read languages", err.Error(), "Check .forge/config.yaml format", "forge testing detect")
	}

	if len(languages) > 0 {
		return languages[0], nil
	}

	// No language detected and no config override
	return "", NewAIError(
		ErrValidation,
		"No language detected",
		"No language signals found in project and no languages override in config",
		"Add a languages field to .forge/config.yaml (e.g., languages: [go])",
		"echo 'languages: [go]' >> .forge/config.yaml",
	)
}

// runTestingGetStrategy returns a Run function for the given strategy kind.
func runTestingGetStrategy(kind string) func(*cobra.Command, []string) {
	return func(_ *cobra.Command, _ []string) {
		projectRoot, err := project.FindProjectRoot()
		if err != nil {
			Exit(ErrProjectNotFound())
		}

		language, err := resolveLanguageFromFlags(projectRoot)
		if err != nil {
			Exit(err)
		}

		data, err := profile.GetStrategy(language, kind)
		if err != nil {
			Exit(NewAIError(ErrInvalidInput, fmt.Sprintf("Failed to get %s strategy", kind), err.Error(), fmt.Sprintf("Check that %q is a valid language with a %s strategy", language, kind), "forge testing detect"))
		}

		fmt.Print(string(data))
	}
}

func runTestingGetJustfile(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	language, err := resolveLanguageFromFlags(projectRoot)
	if err != nil {
		Exit(err)
	}

	data, err := profile.GetJustfileRecipes(language)
	if err != nil {
		Exit(NewAIError(ErrInvalidInput, "Failed to get justfile", err.Error(), fmt.Sprintf("Check that %q is a valid language with justfile-recipes", language), "forge testing detect"))
	}

	fmt.Print(string(data))
}

func runTestingGetTemplate(_ *cobra.Command, args []string) {
	filename := args[0]

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	language, err := resolveLanguageFromFlags(projectRoot)
	if err != nil {
		Exit(err)
	}

	data, err := profile.GetTemplate(language, filename)
	if err != nil {
		Exit(NewAIError(ErrInvalidInput, "Failed to get template", err.Error(), fmt.Sprintf("Check that %q is a valid template for language %q", filename, language), "forge testing detect"))
	}

	fmt.Print(string(data))
}

func runTestingInterfaces(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	interfaces, err := profile.ReadInterfaces(projectRoot)
	if err != nil {
		Exit(NewAIError(ErrValidation, "Failed to read interfaces", err.Error(), "Check .forge/config.yaml format", "forge testing detect"))
	}

	PrintBlockStart()
	if len(interfaces) == 0 {
		PrintField("INTERFACES", "(none)")
	} else {
		for _, iface := range interfaces {
			PrintField("INTERFACE", iface)
		}
		PrintField("SOURCE", "resolved")
	}
	PrintBlockEnd()
}

// printTestingLanguages outputs languages in the structured block format.
func printTestingLanguages(languages []string, source string) {
	PrintBlockStart()
	if len(languages) == 0 {
		PrintField("LANGUAGE", "(none)")
		fmt.Fprintln(os.Stderr, "HINT: No language detected. Add languages to .forge/config.yaml or run: forge testing detect")
	} else {
		for _, l := range languages {
			PrintField("LANGUAGE", l)
		}
		if source != "" {
			PrintField("SOURCE", source)
		}
	}
	PrintBlockEnd()
}
