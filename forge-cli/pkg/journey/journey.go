// Package journey provides Journey-Driven test code generation utilities.
// It generates test files with @feature tags, Contract-based assertions,
// and Journey smoke tests.
package journey

import (
	"fmt"
	"path/filepath"
	"strings"

	"forge-cli/pkg/contract"
	"forge-cli/pkg/descriptor"
	"forge-cli/pkg/profile"
)

// TestGenerationOpts holds options for generating Journey test code.
type TestGenerationOpts struct {
	// Journey is the Journey name (kebab-case, e.g. "task-lifecycle").
	Journey string
	// Contracts are the Contract specifications for the Journey steps.
	Contracts []contract.Contract
	// Framework is the resolved test framework info.
	Framework profile.FrameworkInfo
	// Facts are the Fact Table entries for code reconnaissance.
	Facts []descriptor.FactEntry
	// CustomTemplatePath is the optional path to a custom template directory.
	// Empty means use built-in default templates.
	CustomTemplatePath string
}

// GeneratedTest represents a single generated test file.
type GeneratedTest struct {
	// Filename is the output file name (e.g. "claim_submit_test.go").
	Filename string
	// Content is the generated test code.
	Content string
	// IsSmokeTest is true if this is the Journey smoke test.
	IsSmokeTest bool
}

// FeatureTag returns the framework-native @feature tag string for the given language.
func FeatureTag(framework profile.FrameworkInfo) string {
	switch framework.Name {
	case "go-testing":
		return "//go:build feature"
	case "pytest":
		return "@pytest.mark.feature"
	case "mocha":
		return `describe("@feature", ...)`
	case "junit5":
		return "@Tag(\"feature\")"
	case "rust-test":
		return "#[cfg(feature = \"feature\")]"
	case "maestro":
		return "# feature"
	default:
		return "// @feature"
	}
}

// RegressionTag returns the framework-native @regression tag string.
func RegressionTag(framework profile.FrameworkInfo) string {
	switch framework.Name {
	case "go-testing":
		return "//go:build regression"
	case "pytest":
		return "@pytest.mark.regression"
	case "mocha":
		return `describe("@regression", ...)`
	case "junit5":
		return "@Tag(\"regression\")"
	case "rust-test":
		return "#[cfg(feature = \"regression\")]"
	case "maestro":
		return "# regression"
	default:
		return "// @regression"
	}
}

// TestOutputDir returns the output directory for a Journey's tests.
// Tests go directly into tests/<journey>/ (no staging area).
func TestOutputDir(projectRoot, journey string) string {
	return filepath.Join(projectRoot, "tests", journey)
}

// SmokeTestName generates the smoke test function name for a Journey.
func SmokeTestName(journey string) string {
	parts := strings.Split(journey, "-")
	var name strings.Builder
	for _, p := range parts {
		if len(p) > 0 {
			name.WriteString(strings.ToUpper(p[:1]))
			name.WriteString(p[1:])
		}
	}
	return fmt.Sprintf("TestJourney%sSmoke", name.String())
}

// GenerateGoTest generates a Go test file from Contracts using the go-testing framework.
// The output includes @feature build tags and Contract-based assertions.
func GenerateGoTest(opts TestGenerationOpts) ([]GeneratedTest, error) {
	if opts.Framework.Name != "go-testing" {
		return nil, fmt.Errorf("expected go-testing framework, got %s", opts.Framework.Name)
	}

	var tests []GeneratedTest

	// Generate individual step tests
	for _, c := range opts.Contracts {
		content := generateGoContractTest(c, opts.Facts)
		filename := fmt.Sprintf("step%d_%s_test.go", c.Step, sanitizeName(c.Action))
		tests = append(tests, GeneratedTest{
			Filename:    filename,
			Content:     content,
			IsSmokeTest: false,
		})
	}

	// Generate Journey smoke test (happy path end-to-end)
	smokeContent := generateGoSmokeTest(opts)
	tests = append(tests, GeneratedTest{
		Filename:    fmt.Sprintf("%s_smoke_test.go", sanitizeName(opts.Journey)),
		Content:     smokeContent,
		IsSmokeTest: true,
	})

	return tests, nil
}

// generateGoContractTest generates a Go test file for a single Contract step.
func generateGoContractTest(c contract.Contract, facts []descriptor.FactEntry) string {
	var sb strings.Builder

	// @feature build tag
	sb.WriteString("//go:build feature\n\n")
	sb.WriteString("package journey_test\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"os/exec\"\n")
	sb.WriteString("\t\"testing\"\n\n")
	sb.WriteString("\t\"github.com/stretchr/testify/assert\"\n")
	sb.WriteString(")\n\n")

	// Generate test for each Outcome
	for _, o := range c.Outcomes {
		funcName := fmt.Sprintf("TestStep%d_%s_%s", c.Step, sanitizeName(c.Action), sanitizeName(o.Name))
		fmt.Fprintf(&sb, "func %s(t *testing.T) {\n", funcName)
		fmt.Fprintf(&sb, "\t// Traceability: Journey %s, Step %d, Outcome %q\n", c.Journey, c.Step, o.Name)
		fmt.Fprintf(&sb, "\t// Action: %s\n\n", c.Action)

		// Setup: isolated temp directory
		sb.WriteString("\tdir := t.TempDir()\n")
		sb.WriteString("\t_ = dir // VERIFY: setup project structure as needed\n\n")

		// Execute command
		cmdParts := strings.Fields(c.Action)
		if len(cmdParts) > 0 {
			fmt.Fprintf(&sb, "\tcmd := exec.Command(%q", cmdParts[0])
			for _, arg := range cmdParts[1:] {
				fmt.Fprintf(&sb, ", %q", arg)
			}
			sb.WriteString(")\n")
		}
		sb.WriteString("\tout, err := cmd.CombinedOutput()\n\n")

		// Outcome-specific assertions
		if o.Name == "success" || !strings.Contains(strings.ToLower(o.Name), "error") {
			sb.WriteString("\tif err != nil {\n")
			sb.WriteString("\t\tt.Fatalf(\"command failed: %s\", err)\n")
			sb.WriteString("\t}\n")
		} else {
			sb.WriteString("\t// Expected non-zero exit for this Outcome\n")
			sb.WriteString("\tif err == nil {\n")
			fmt.Fprintf(&sb, "\t\tt.Fatal(\"expected command to fail for Outcome %q\")\n", o.Name)
			sb.WriteString("\t}\n")
		}

		// Output assertion using Fact Table
		outputPattern := descriptor.BuildAssertionRegex(o.Output, facts)
		fmt.Fprintf(&sb, "\tassert.Regexp(t, %q, string(out), \"output should match Outcome %q assertion\")\n", outputPattern, o.Name)

		sb.WriteString("}\n\n")
	}

	return sb.String()
}

// generateGoSmokeTest generates the Journey smoke test that runs the happy path end-to-end.
func generateGoSmokeTest(opts TestGenerationOpts) string {
	var sb strings.Builder

	sb.WriteString("//go:build feature\n\n")
	sb.WriteString("package journey_test\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"os/exec\"\n")
	sb.WriteString("\t\"path/filepath\"\n")
	sb.WriteString("\t\"testing\"\n\n")
	sb.WriteString("\t\"github.com/stretchr/testify/assert\"\n")
	sb.WriteString(")\n\n")

	fmt.Fprintf(&sb, "// %s verifies the %s Journey happy path end-to-end.\n",
		SmokeTestName(opts.Journey), opts.Journey)
	sb.WriteString("// Each step output matches the Contract 'success' Outcome.\n")
	fmt.Fprintf(&sb, "func %s(t *testing.T) {\n", SmokeTestName(opts.Journey))

	// Isolated temp directory for the Journey
	sb.WriteString("\tdir := t.TempDir()\n")
	sb.WriteString("\t_ = filepath.Join(dir, \"dummy\") // VERIFY: setup project structure\n\n")

	// Sequential step execution with state passing
	for _, c := range opts.Contracts {
		// Find success Outcome
		var successOutcome *contract.Outcome
		for i := range c.Outcomes {
			if c.Outcomes[i].Name == "success" {
				successOutcome = &c.Outcomes[i]
				break
			}
		}
		if successOutcome == nil {
			continue
		}

		fmt.Fprintf(&sb, "\t// Step %d: %s\n", c.Step, c.Action)
		cmdParts := strings.Fields(c.Action)
		if len(cmdParts) > 0 {
			fmt.Fprintf(&sb, "\tstep%dCmd := exec.Command(%q", c.Step, cmdParts[0])
			for _, arg := range cmdParts[1:] {
				fmt.Fprintf(&sb, ", %q", arg)
			}
			sb.WriteString(")\n")
		}
		fmt.Fprintf(&sb, "\tstep%dOut, err := step%dCmd.CombinedOutput()\n", c.Step, c.Step)
		sb.WriteString("\tif err != nil {\n")
		fmt.Fprintf(&sb, "\t\tt.Fatalf(\"Step %d failed: %%s\\nOutput: %%s\", err, step%dOut)\n", c.Step, c.Step)
		sb.WriteString("\t}\n")

		// Assert output matches success Outcome
		outputPattern := descriptor.BuildAssertionRegex(successOutcome.Output, opts.Facts)
		fmt.Fprintf(&sb, "\tassert.Regexp(t, %q, string(step%dOut), \"Step %d output should match success Outcome\")\n", outputPattern, c.Step, c.Step)
		sb.WriteString("\n")
	}

	// Verify Journey Invariants
	if len(opts.Contracts) > 0 && len(opts.Contracts[0].Invariants) > 0 {
		sb.WriteString("\t// Journey Invariants verification\n")
		for _, inv := range opts.Contracts[0].Invariants {
			fmt.Fprintf(&sb, "\t// Invariant: %s\n", inv)
			sb.WriteString("\t// VERIFY: check invariant holds across all steps\n")
		}
	}

	sb.WriteString("}\n")

	return sb.String()
}

// GeneratePythonTest generates a Python test file from Contracts using pytest.
func GeneratePythonTest(opts TestGenerationOpts) ([]GeneratedTest, error) {
	if opts.Framework.Name != "pytest" {
		return nil, fmt.Errorf("expected pytest framework, got %s", opts.Framework.Name)
	}

	var tests []GeneratedTest

	for _, c := range opts.Contracts {
		content := generatePythonContractTest(c, opts.Facts)
		filename := fmt.Sprintf("test_step%d_%s.py", c.Step, sanitizeName(c.Action))
		tests = append(tests, GeneratedTest{
			Filename:    filename,
			Content:     content,
			IsSmokeTest: false,
		})
	}

	// Smoke test
	smokeContent := generatePythonSmokeTest(opts)
	tests = append(tests, GeneratedTest{
		Filename:    fmt.Sprintf("test_%s_smoke.py", sanitizeName(opts.Journey)),
		Content:     smokeContent,
		IsSmokeTest: true,
	})

	return tests, nil
}

func generatePythonContractTest(c contract.Contract, facts []descriptor.FactEntry) string {
	var sb strings.Builder

	sb.WriteString("import pytest\n\n")
	sb.WriteString("pytestmark = pytest.mark.feature\n\n")

	for _, o := range c.Outcomes {
		funcName := fmt.Sprintf("test_step%d_%s_%s", c.Step, sanitizeName(c.Action), sanitizeName(o.Name))
		fmt.Fprintf(&sb, "def %s(tmp_path):\n", funcName)
		fmt.Fprintf(&sb, "    \"\"\"Journey %s, Step %d, Outcome %q: %s\"\"\"\n",
			c.Journey, c.Step, o.Name, c.Action)
		sb.WriteString("    # VERIFY: setup project structure in tmp_path\n")
		sb.WriteString("    import subprocess\n")
		cmdParts := strings.Fields(c.Action)
		fmt.Fprintf(&sb, "    result = subprocess.run(%q", cmdParts)
		if len(cmdParts) > 1 {
			fmt.Fprintf(&sb, " + %q", cmdParts[1:])
		}
		sb.WriteString(", capture_output=True, text=True)\n\n")

		outputPattern := descriptor.BuildAssertionRegex(o.Output, facts)
		sb.WriteString("    import re\n")
		fmt.Fprintf(&sb, "    assert re.search(%q, result.stdout), \"output should match Outcome %q\"\n", outputPattern, o.Name)
		sb.WriteString("\n")
	}

	return sb.String()
}

func generatePythonSmokeTest(opts TestGenerationOpts) string {
	var sb strings.Builder

	sb.WriteString("import pytest\n\n")
	sb.WriteString("pytestmark = pytest.mark.feature\n\n")

	funcName := fmt.Sprintf("test_%s_smoke", sanitizeName(opts.Journey))
	fmt.Fprintf(&sb, "def %s(tmp_path):\n", funcName)
	fmt.Fprintf(&sb, "    \"\"\"Smoke test for Journey %s: happy path end-to-end.\"\"\"\n", opts.Journey)
	sb.WriteString("    import subprocess\n\n")

	for _, c := range opts.Contracts {
		var successOutcome *contract.Outcome
		for i := range c.Outcomes {
			if c.Outcomes[i].Name == "success" {
				successOutcome = &c.Outcomes[i]
				break
			}
		}
		if successOutcome == nil {
			continue
		}

		fmt.Fprintf(&sb, "    # Step %d: %s\n", c.Step, c.Action)
		cmdParts := strings.Fields(c.Action)
		fmt.Fprintf(&sb, "    step%d = subprocess.run(%q, capture_output=True, text=True)\n", c.Step, cmdParts)
		fmt.Fprintf(&sb, "    assert step%d.returncode == 0, f\"Step %d failed: {step%d.stderr}\"\n", c.Step, c.Step, c.Step)
		outputPattern := descriptor.BuildAssertionRegex(successOutcome.Output, opts.Facts)
		sb.WriteString("    import re\n")
		fmt.Fprintf(&sb, "    assert re.search(%q, step%d.stdout), \"Step %d output mismatch\"\n", outputPattern, c.Step, c.Step)
		sb.WriteString("\n")
	}

	return sb.String()
}

// GenerateJSTest generates a JavaScript test file from Contracts using mocha/describe-it.
func GenerateJSTest(opts TestGenerationOpts) ([]GeneratedTest, error) {
	if opts.Framework.Name != "mocha" {
		return nil, fmt.Errorf("expected mocha framework, got %s", opts.Framework.Name)
	}

	var tests []GeneratedTest

	for _, c := range opts.Contracts {
		content := generateJSContractTest(c, opts.Facts)
		filename := fmt.Sprintf("step%d_%s.spec.ts", c.Step, sanitizeName(c.Action))
		tests = append(tests, GeneratedTest{
			Filename:    filename,
			Content:     content,
			IsSmokeTest: false,
		})
	}

	smokeContent := generateJSSmokeTest(opts)
	tests = append(tests, GeneratedTest{
		Filename:    fmt.Sprintf("%s_smoke.spec.ts", sanitizeName(opts.Journey)),
		Content:     smokeContent,
		IsSmokeTest: true,
	})

	return tests, nil
}

func generateJSContractTest(c contract.Contract, facts []descriptor.FactEntry) string {
	var sb strings.Builder

	sb.WriteString("describe(\"@feature\", () => {\n")

	for _, o := range c.Outcomes {
		fmt.Fprintf(&sb, "  it(\"Step %d %s - %s\", async () => {\n", c.Step, c.Action, o.Name)
		fmt.Fprintf(&sb, "    // Traceability: Journey %s, Step %d, Outcome %q\n", c.Journey, c.Step, o.Name)
		sb.WriteString("    // VERIFY: setup test environment\n")
		sb.WriteString("    const { execSync } = require('child_process');\n")
		cmdParts := strings.Fields(c.Action)
		fmt.Fprintf(&sb, "    const output = execSync('%s').toString();\n", strings.Join(cmdParts, " "))
		outputPattern := descriptor.BuildAssertionRegex(o.Output, facts)
		fmt.Fprintf(&sb, "    expect(output).toMatch(/%s/);\n", outputPattern)
		sb.WriteString("  });\n\n")
	}

	sb.WriteString("});\n")

	return sb.String()
}

func generateJSSmokeTest(opts TestGenerationOpts) string {
	var sb strings.Builder

	sb.WriteString("describe(\"@feature\", () => {\n")
	fmt.Fprintf(&sb, "  it(\"Journey %s smoke test - happy path\", async () => {\n", opts.Journey)
	sb.WriteString("    const { execSync } = require('child_process');\n\n")

	for _, c := range opts.Contracts {
		var successOutcome *contract.Outcome
		for i := range c.Outcomes {
			if c.Outcomes[i].Name == "success" {
				successOutcome = &c.Outcomes[i]
				break
			}
		}
		if successOutcome == nil {
			continue
		}

		cmdParts := strings.Fields(c.Action)
		fmt.Fprintf(&sb, "    // Step %d: %s\n", c.Step, c.Action)
		fmt.Fprintf(&sb, "    const step%d = execSync('%s').toString();\n", c.Step, strings.Join(cmdParts, " "))
		outputPattern := descriptor.BuildAssertionRegex(successOutcome.Output, opts.Facts)
		fmt.Fprintf(&sb, "    expect(step%d).toMatch(/%s/);\n", c.Step, outputPattern)
		sb.WriteString("\n")
	}

	sb.WriteString("  });\n")
	sb.WriteString("});\n")

	return sb.String()
}

// GenerateDispatched dispatches test generation to the correct framework handler.
func GenerateDispatched(opts TestGenerationOpts) ([]GeneratedTest, error) {
	switch opts.Framework.Name {
	case "go-testing":
		return GenerateGoTest(opts)
	case "pytest":
		return GeneratePythonTest(opts)
	case "mocha":
		return GenerateJSTest(opts)
	default:
		return nil, fmt.Errorf("unsupported framework for Journey test generation: %s", opts.Framework.Name)
	}
}

// sanitizeName converts a human-readable name to a safe identifier component.
func sanitizeName(name string) string {
	replacer := strings.NewReplacer(
		" ", "_",
		"-", "_",
		"/", "_",
		".", "_",
		":", "_",
	)
	s := replacer.Replace(name)
	s = strings.ToLower(s)
	// Remove consecutive underscores
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}
	s = strings.Trim(s, "_")
	return s
}
