//go:build e2e

package commandregression

import (
	"os/exec"
	"regexp"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Removed forge test commands — CLI e2e tests for feature: test-knowledge-convention-driven
// Tests cover removed commands (forge test detect/get/interfaces/framework) that
// should return errors after Profile removal.
// ==============================================================================

// runForgeTestCommand runs a forge test subcommand and returns output + exit code.
func runForgeTestCommand(t *testing.T, subcommand string, extraArgs ...string) (string, int) {
	t.Helper()
	args := []string{"test", subcommand}
	args = append(args, extraArgs...)
	cmd := exec.Command(testkit.ForgeBinary, args...)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return string(out), exitCode
}

// --- TC-021: Removed Command forge test detect Returns Error ---
// Traceability: TC-021 -> PRD Scope / Remove CLI commands

func TestTC_021_RemovedCommandForgeTestDetectErrors(t *testing.T) {
	output, exitCode := runForgeTestCommand(t, "detect")
	t.Logf("forge test detect output: %s", output)

	// Expect error exit code
	assert.NotEqual(t, 0, exitCode,
		"forge test detect should return non-zero exit code")

	// Expect error message about unknown/removed command
	errorPattern := regexp.MustCompile(`(?i)unknown command|command not found|removed`)
	assert.True(t, errorPattern.MatchString(output),
		"Expected error about removed command, got: %s", output)
}

// --- TC-022: Removed Command forge test get Returns Error ---
// Traceability: TC-022 -> PRD Scope / Remove CLI commands

func TestTC_022_RemovedCommandForgeTestGetErrors(t *testing.T) {
	output, exitCode := runForgeTestCommand(t, "get")
	t.Logf("forge test get output: %s", output)

	assert.NotEqual(t, 0, exitCode,
		"forge test get should return non-zero exit code")

	errorPattern := regexp.MustCompile(`(?i)unknown command|command not found|removed`)
	assert.True(t, errorPattern.MatchString(output),
		"Expected error about removed command, got: %s", output)
}

// --- TC-023: Removed Command forge test interfaces Returns Error ---
// Traceability: TC-023 -> PRD Scope / Remove CLI commands

func TestTC_023_RemovedCommandForgeTestInterfacesErrors(t *testing.T) {
	output, exitCode := runForgeTestCommand(t, "interfaces")
	t.Logf("forge test interfaces output: %s", output)

	assert.NotEqual(t, 0, exitCode,
		"forge test interfaces should return non-zero exit code")

	errorPattern := regexp.MustCompile(`(?i)unknown command|command not found|removed`)
	assert.True(t, errorPattern.MatchString(output),
		"Expected error about removed command, got: %s", output)
}

// --- TC-024: Removed Command forge test framework Returns Error ---
// Traceability: TC-024 -> PRD Scope / Remove CLI commands

func TestTC_024_RemovedCommandForgeTestFrameworkErrors(t *testing.T) {
	output, exitCode := runForgeTestCommand(t, "framework")
	t.Logf("forge test framework output: %s", output)

	assert.NotEqual(t, 0, exitCode,
		"forge test framework should return non-zero exit code")

	errorPattern := regexp.MustCompile(`(?i)unknown command|command not found|removed`)
	assert.True(t, errorPattern.MatchString(output),
		"Expected error about removed command, got: %s", output)
}
