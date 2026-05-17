package cmd

import (
	"fmt"
	"os"
	"strings"

	"forge-cli/pkg/profile"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Resolve testing strategies based on project language detection",
	Long: `Resolve testing strategies based on project language detection.

Supports auto-detection from project files (go.mod, package.json, Cargo.toml, etc.)
or explicit override via .forge/config.yaml languages field.

Subcommands:
  detect              — output detected language(s)
  get generate        — output generate.md strategy
  get run             — output run.md strategy
  get justfile        — output justfile-recipes
  get template <file> — output specified template file
  interfaces          — output interface types for the project
  promote <journey>   — promote a journey's @feature tags to @regression
  run-journey <name>  — run a single journey in isolated temp directory
  verify              — detect contract breakage against current code`,
	Args: cobra.NoArgs,
	Run:  runTestResolve,
}

var testDetectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect languages from project structure (ignores config overrides)",
	Args:  cobra.NoArgs,
	Run:   runTestDetect,
}

var testGetLanguage string

var testGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get testing strategy file content",
	Long: `Output a testing strategy file for the detected (or specified) language.

Auto-detects the project language when --language flag is not specified.
For multi-language projects, use --language to select a specific language;
without the flag, the first detected language is used.

Examples:
  forge test get generate
  forge test get run --language javascript
  forge test get template test-file.go`,
}

var testGetGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Output generate.md strategy",
	Args:  cobra.NoArgs,
	Run:   runTestGetStrategy("generate"),
}

var testGetRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Output run.md strategy",
	Args:  cobra.NoArgs,
	Run:   runTestGetStrategy("run"),
}

var testGetJustfileCmd = &cobra.Command{
	Use:   "justfile",
	Short: "Output justfile-recipes",
	Args:  cobra.NoArgs,
	Run:   runTestGetJustfile,
}

var testGetTemplateCmd = &cobra.Command{
	Use:   "template <file>",
	Short: "Output a specific template file",
	Args:  cobra.ExactArgs(1),
	Run:   runTestGetTemplate,
}

var testInterfacesCmd = &cobra.Command{
	Use:   "interfaces",
	Short: "Output interface types for the project",
	Long: `Output interface types for the project.

Returns config.Interfaces if set in .forge/config.yaml,
otherwise returns the union of all detected languages' default interfaces.`,
	Args: cobra.NoArgs,
	Run:  runTestInterfaces,
}

var testPromoteCmd = &cobra.Command{
	Use:   "promote <journey-name>",
	Short: "Promote a journey's @feature tags to @regression",
	Long: `Promote a journey by replacing all @feature tags with @regression tags.

Before promoting, runs all tests for the journey. If any test fails,
the promotion is refused and a failure report is printed.

Tag lifecycle:
  @feature (newly generated, under validation) -> @regression (verified, regression)`,
	Args: cobra.ExactArgs(1),
	Run:  runTestPromote,
}

var testRunJourneyCmd = &cobra.Command{
	Use:   "run-journey <journey-name>",
	Short: "Run a single journey in isolated temp directory",
	Long: `Run a single journey's test command in an isolated temporary directory.

Reads the test-command from .forge/config.yaml and executes it in the journey's
isolated work directory. The temp directory is cleaned up after execution,
regardless of success or failure.

The journey name is used as part of the temp directory path for traceability.

Output is a structured block with journey name, result, duration, and any failures.`,
	Args: cobra.ExactArgs(1),
	Run:  runTestRunJourney,
}

var testFrameworkCmd = &cobra.Command{
	Use:   "framework",
	Short: "Resolve the test framework for the project",
	Long: `Resolve the test framework for the project.

Priority:
  1. test-framework field in .forge/config.yaml (explicit override)
  2. Default framework for the first resolved language
  3. No framework resolved

Output is structured fields:
  FRAMEWORK   — framework name (e.g. "go-testing", "pytest")
  PATTERN     — test function pattern (e.g. "func Test*", "def test_*")
  FILES       — file naming pattern (e.g. "*_test.go", "test_*.py")
  SOURCE      — how it was resolved ("config", "language-default", "none")`,
	Args: cobra.NoArgs,
	Run:  runTestFramework,
}

func init() {
	testCmd.AddCommand(testDetectCmd)
	testCmd.AddCommand(testGetCmd)
	testCmd.AddCommand(testInterfacesCmd)
	testCmd.AddCommand(testFrameworkCmd)
	testCmd.AddCommand(testPromoteCmd)
	testCmd.AddCommand(testRunJourneyCmd)

	testGetCmd.PersistentFlags().StringVar(&testGetLanguage, "language", "", "language key (auto-detected if omitted)")

	testGetCmd.AddCommand(testGetGenerateCmd)
	testGetCmd.AddCommand(testGetRunCmd)
	testGetCmd.AddCommand(testGetJustfileCmd)
	testGetCmd.AddCommand(testGetTemplateCmd)
}

func runTestResolve(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	languages, err := profile.ReadLanguages(projectRoot)
	if err != nil {
		Exit(NewAIError(ErrValidation, "Failed to read languages", err.Error(), "Check .forge/config.yaml format", "forge test detect"))
	}

	if len(languages) > 0 {
		printLanguages(languages, "resolved")
		return
	}

	printLanguages(nil, "")
}

func runTestDetect(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	detected, err := profile.DetectLanguages(projectRoot)
	if err != nil {
		Exit(NewAIError(ErrValidation, "Detection failed", err.Error(), "Set languages in .forge/config.yaml manually", "forge test detect"))
	}

	var names []string
	for _, l := range detected {
		names = append(names, string(l))
	}
	printLanguages(names, "detected")
}

// resolveLanguageFromFlags resolves the language from --language flag, config, or auto-detect.
// Returns an error if no language can be resolved.
func resolveLanguageFromFlags(projectRoot string) (string, error) {
	// Explicit --language flag takes highest priority
	if testGetLanguage != "" {
		if !profile.IsKnownLanguage(testGetLanguage) {
			return "", NewAIError(
				ErrInvalidInput,
				fmt.Sprintf("Unknown language: %s", testGetLanguage),
				"Language key is not in the known languages list",
				fmt.Sprintf("Choose from: %s", strings.Join(profile.KnownLanguages, ", ")),
				fmt.Sprintf("forge test get generate --language %s", profile.KnownLanguages[0]),
			)
		}
		return testGetLanguage, nil
	}

	// Try config override, then auto-detect
	languages, err := profile.ReadLanguages(projectRoot)
	if err != nil {
		return "", NewAIError(ErrValidation, "Failed to read languages", err.Error(), "Check .forge/config.yaml format", "forge test detect")
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

// runTestGetStrategy returns a Run function for the given strategy kind.
func runTestGetStrategy(kind string) func(*cobra.Command, []string) {
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
			Exit(NewAIError(ErrInvalidInput, fmt.Sprintf("Failed to get %s strategy", kind), err.Error(), fmt.Sprintf("Check that %q is a valid language with a %s strategy", language, kind), "forge test detect"))
		}

		fmt.Print(string(data))
	}
}

func runTestGetJustfile(_ *cobra.Command, _ []string) {
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
		Exit(NewAIError(ErrInvalidInput, "Failed to get justfile", err.Error(), fmt.Sprintf("Check that %q is a valid language with justfile-recipes", language), "forge test detect"))
	}

	fmt.Print(string(data))
}

func runTestGetTemplate(_ *cobra.Command, args []string) {
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
		Exit(NewAIError(ErrInvalidInput, "Failed to get template", err.Error(), fmt.Sprintf("Check that %q is a valid template for language %q", filename, language), "forge test detect"))
	}

	fmt.Print(string(data))
}

func runTestInterfaces(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	interfaces, err := profile.ReadInterfaces(projectRoot)
	if err != nil {
		Exit(NewAIError(ErrValidation, "Failed to read interfaces", err.Error(), "Check .forge/config.yaml format", "forge test detect"))
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

// printLanguages outputs languages in the structured block format.
func printLanguages(languages []string, source string) {
	PrintBlockStart()
	if len(languages) == 0 {
		PrintField("LANGUAGE", "(none)")
		fmt.Fprintln(os.Stderr, "HINT: No language detected. Add languages to .forge/config.yaml or run: forge test detect")
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

func runTestRunJourney(_ *cobra.Command, args []string) {
	journeyName := args[0]

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	cfg, err := resolveJourneyExecutionConfig(projectRoot)
	if err != nil {
		Exit(NewAIError(ErrValidation, "Cannot resolve journey execution config", err.Error(),
			"Set test-command in .forge/config.yaml", "echo 'test-command: go test ./...' >> .forge/config.yaml"))
	}

	// Create isolated work directory
	workDir, cleanup, err := createJourneyWorkDir(projectRoot, journeyName)
	if err != nil {
		Exit(NewAIError(ErrValidation, "Failed to create journey work directory", err.Error(),
			"Check temp directory permissions", "forge test run-journey "+journeyName))
	}
	defer cleanup()

	// Execute the test command in isolation
	result := executeJourneyInIsolation(cfg, workDir, journeyName)

	// Output the result report
	fmt.Print(result.FormatReport())
}

func runTestFramework(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	// Determine source
	cfg, _ := profile.ReadConfig(projectRoot)
	source := "none"
	if cfg != nil && cfg.TestFramework != "" {
		source = "config"
	}

	fw, err := profile.ResolveTestFramework(projectRoot)
	if err != nil {
		Exit(NewAIError(ErrValidation, "Failed to resolve test framework", err.Error(), "Set test-framework in .forge/config.yaml", "forge test framework"))
	}

	PrintBlockStart()
	if fw.Name == "" {
		PrintField("FRAMEWORK", "(none)")
		PrintField("SOURCE", "none")
	} else {
		PrintField("FRAMEWORK", fw.Name)
		if fw.TestFunctionPattern != "" {
			PrintField("PATTERN", fw.TestFunctionPattern)
		}
		if fw.FilePattern != "" {
			PrintField("FILES", fw.FilePattern)
		}
		if source == "config" {
			PrintField("SOURCE", "config")
		} else {
			PrintField("SOURCE", "language-default")
		}
	}
	PrintBlockEnd()
}
