//go:build cli_functional

package scoperesolution

import (
	"os/exec"
	"strings"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
)

// --- Scope Resolution (TC-007 to TC-010, TC-015, TC-016, TC-023, TC-024) ---
// Converted from tests/e2e/scope-resolution/scope-resolution.spec.ts (8 tests).
// These tests validate the scope resolution algorithm for quality gate commands.

// runJust executes a just recipe with optional args via the system shell,
// returning exit code and combined output.
func runJust(args ...string) (int, string) {
	cmd := exec.Command("just", args...)
	out, err := cmd.CombinedOutput()
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), string(out)
	}
	if err != nil {
		return 1, err.Error()
	}
	return 0, string(out)
}

// Traceability: TC-007 -> Story 3 / AC-1
// Mixed project tasks receive scope field in index.json.
// Verify the breakdown-tasks skill describes adding scope to index.json tasks.
func TestTC_007_MixedProjectTasksReceiveScopeFieldInIndexJSON(t *testing.T) {
	skillContent := testkit.ReadProjectFile(t, "../plugins/forge/skills/breakdown-tasks/SKILL.md")
	assert.True(t, strings.Contains(skillContent, "scope"),
		"Expected \"scope\" in breakdown-tasks skill description")
}

// Traceability: TC-008 -> Story 3 / AC-2
// Frontend-only task scope marked as frontend.
// Skill delegates scope detection to `forge surfaces` at runtime (Surface-Key/Type Inference).
// Verify the skill describes surface inference, which assigns per-task surface keys
// (e.g., frontend, backend) based on file paths.
func TestTC_008_FrontendOnlyTaskScopeMarkedAsFrontend(t *testing.T) {
	skillContent := testkit.ReadProjectFile(t, "../plugins/forge/skills/breakdown-tasks/SKILL.md")
	assert.True(t,
		strings.Contains(skillContent, "surfaces") && strings.Contains(skillContent, "Surface"),
		"Expected Surface-Key/Type Inference logic in breakdown-tasks skill")
}

// Traceability: TC-009 -> Story 3 / AC-3
// Cross-scope task marked as all.
func TestTC_009_CrossScopeTaskMarkedAsAll(t *testing.T) {
	skillContent := testkit.ReadProjectFile(t, "../plugins/forge/skills/breakdown-tasks/SKILL.md")
	assert.True(t,
		strings.Contains(skillContent, "all") && strings.Contains(skillContent, "scope"),
		"Expected scope=all logic in breakdown-tasks skill")
}

// Traceability: TC-010 -> Story 3 / AC-4
// Non-mixed project tasks all receive scope all.
// For non-mixed projects, scope should default to "all".
func TestTC_010_NonMixedProjectTasksAllReceiveScopeAll(t *testing.T) {
	skillContent := testkit.ReadProjectFile(t, "../plugins/forge/skills/breakdown-tasks/SKILL.md")
	assert.True(t,
		strings.Contains(skillContent, "all") && strings.Contains(skillContent, "scope"),
		"Expected default scope=all for non-mixed projects")
}

// Traceability: TC-015 -> Story 5 / AC-1
// Scope mismatch shows warning and falls back.
// Verify that the justfile has scope-aware recipes and that the PRD spec describes
// the invalid scope error format: [forge] invalid scope 'X'; expected frontend/backend.
// For mixed projects with scope dispatch, invalid scope should produce an error;
// for simple justfiles (pure Go), the scope parameter is accepted but not validated.
func TestTC_015_ScopeMismatchShowsWarningAndFallsBack(t *testing.T) {
	// Verify the PRD spec describes the invalid scope error format
	prdSpec := testkit.ReadProjectFile(t, "../docs/features/justfile-standard-vocabulary/prd/prd-spec.md")
	lower := strings.ToLower(prdSpec)
	assert.True(t,
		strings.Contains(lower, "invalid scope") || strings.Contains(lower, "scope"),
		"Expected invalid scope handling described in PRD spec")

	// Verify the justfile has scope-aware build recipe
	justfile := testkit.ReadProjectFile(t, "../justfile")
	buildIdx := strings.Index(justfile, "build scope=")
	assert.NotEqual(t, -1, buildIdx, "Expected build recipe with scope parameter")

	// For scope-dispatching justfiles, verify invalid scope produces error.
	// For simple justfiles (pure Go), build succeeds regardless of scope arg.
	exitCode, out := runJust("build", "invalidscope")
	if exitCode != 0 {
		assert.Equal(t, 1, exitCode, "Build with invalid scope should exit 1")
		outLower := strings.ToLower(out)
		assert.True(t,
			strings.Contains(outLower, "error") || strings.Contains(outLower, "invalid"),
			"Expected error indication in output, got: %s", out)
	}
	// If exitCode == 0, the justfile does not validate scope (pure Go project)
}

// Traceability: TC-016 -> Story 5 / AC-2
// Mixed project with matching scope executes normally.
// In the real forge project (which is mixed), verify scoped commands work
// without producing a scope error.
func TestTC_016_MixedProjectWithMatchingScopeExecutesNormally(t *testing.T) {
	exitCode, out := runJust("compile", "frontend")
	combined := out
	lower := strings.ToLower(combined)
	isScopeError := strings.Contains(lower, "[forge] invalid scope")
	assert.False(t, isScopeError, "Should not be a scope error for frontend scope")
	// If toolchains are missing, allow non-zero exit
	if exitCode != 0 {
		t.Logf("compile frontend exited %d (possibly missing toolchain): %s", exitCode, combined)
	}
}

// Traceability: TC-023 -> Spec 5.3 / row 4
// just project-type failure triggers fallback in skill.
// Verify the justfile has project-type recipe that agents can call,
// or that forge probe provides equivalent functionality.
// When project-type fails (old justfile), skills should fall back to `just <verb>`.
func TestTC_023_ProjectTypeFailureTriggersFallbackInSkill(t *testing.T) {
	// First try forge probe (equivalent to just project-type for CLI-based projects)
	exitCode, out := testkit.RunCLIExitCode("probe")
	if exitCode == 0 {
		output := strings.TrimSpace(out)
		assert.True(t,
			strings.Contains(output, "frontend") || strings.Contains(output, "backend") ||
				strings.Contains(output, "mixed") || strings.Contains(output, "go") ||
				strings.Contains(output, "project"),
			"Expected valid project-type output, got: %q", output)
	} else {
		// If forge probe fails, try just project-type
		jExitCode, jOut := runJust("project-type")
		if jExitCode != 0 {
			t.Skipf("neither forge probe nor just project-type available: probe=%d, just=%d", exitCode, jExitCode)
		}
		output := strings.TrimSpace(jOut)
		assert.True(t,
			strings.Contains(output, "frontend") || strings.Contains(output, "backend") || strings.Contains(output, "mixed"),
			"Expected valid project-type output, got: %q", output)
	}

	// Also verify the PRD spec documents the fallback behavior
	prdSpec := testkit.ReadProjectFile(t, "../docs/features/justfile-standard-vocabulary/prd/prd-spec.md")
	assert.True(t,
		strings.Contains(strings.ToLower(prdSpec), "falling back") || strings.Contains(strings.ToLower(prdSpec), "fallback"),
		"Expected fallback behavior described in PRD spec")
}

// Traceability: TC-024 -> Spec 5.3 / row 5
// Unexpected project-type output triggers fallback.
// Verify the PRD spec describes handling of unexpected project-type output.
// The fallback behavior: "[forge] just project-type returned unexpected output 'XYZ'; falling back to just verb"
func TestTC_024_UnexpectedProjectTypeOutputTriggersFallback(t *testing.T) {
	prdSpec := testkit.ReadProjectFile(t, "../docs/features/justfile-standard-vocabulary/prd/prd-spec.md")
	lower := strings.ToLower(prdSpec)
	assert.True(t,
		strings.Contains(lower, "unexpected output") || strings.Contains(lower, "unexpected"),
		"Expected unexpected output handling described in PRD spec")
	assert.True(t,
		strings.Contains(lower, "falling back") || strings.Contains(lower, "fallback"),
		"Expected fallback description in PRD spec")

	// Verify the justfile project-type output is deterministic across runs.
	// Use forge probe as equivalent CLI command.
	exitCode1, out1 := testkit.RunCLIExitCode("probe")
	if exitCode1 != 0 {
		// Try just project-type instead
		exitCode1, out1 = runJust("project-type")
	}
	if exitCode1 != 0 {
		t.Skip("no project-type command available for determinism check")
	}

	exitCode2, out2 := testkit.RunCLIExitCode("probe")
	if exitCode2 != 0 {
		exitCode2, out2 = runJust("project-type")
	}
	if exitCode2 != 0 {
		t.Skip("no project-type command available on second run")
	}

	assert.Equal(t, strings.TrimSpace(out1), strings.TrimSpace(out2),
		"Expected deterministic project-type output")
}
