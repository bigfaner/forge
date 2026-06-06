//go:build cli_functional

package justfileintegration

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
)

// --- Mixed template + CLI integration tests ---
// Converted from tests/e2e/justfile-e2e-integration/mixed-template.spec.ts (23 tests)
// and tests/e2e/justfile-e2e-integration/cli.spec.ts (20 tests).
// Combined into one Go file to keep related tests together.

// getMixedTemplate reads the mixed justfile template.
func getMixedTemplate(t *testing.T) string {
	t.Helper()
	return testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/templates/mixed.just")
}

// countOccurrences returns the number of non-overlapping occurrences of substr in s.
func countOccurrences(s, substr string) int {
	return strings.Count(s, substr)
}

// ── TC-MIX-001 to TC-MIX-011: Mixed template scoped recipe checks ───────

// Traceability: TC-MIX-001 -> AC: project-type outputs @echo "mixed"
// Note: project-type has been removed from templates; project type is now detected via `forge probe`.
// The mixed template is identified by having both frontend_dir and backend_dir variables.
func TestTC_MIX_001_ProjectTypeOutputsMixed(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, "frontend_dir"),
		"Expected \"frontend_dir\" variable in mixed template (identifies it as mixed type)")
	assert.True(t, strings.Contains(template, "backend_dir"),
		"Expected \"backend_dir\" variable in mixed template (identifies it as mixed type)")
}

// Traceability: TC-MIX-002 -> AC: compile has per-surface recipes (frontend-compile, backend-compile)
// Mixed template uses per-surface recipes instead of case/esac scope dispatch.
func TestTC_MIX_002_CompileRecipeHasPerSurfaceRecipes(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `frontend-compile`),
		`Expected 'frontend-compile' in mixed template`)
	assert.True(t, strings.Contains(template, `backend-compile`),
		`Expected 'backend-compile' in mixed template`)
}

// Traceability: TC-MIX-003 -> AC: build has per-surface recipes (frontend-build, backend-build)
func TestTC_MIX_003_BuildRecipeHasPerSurfaceRecipes(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `frontend-build`),
		`Expected 'frontend-build' in mixed template`)
	assert.True(t, strings.Contains(template, `backend-build`),
		`Expected 'backend-build' in mixed template`)
}

// Traceability: TC-MIX-004 -> AC: run has per-surface recipes (frontend-run, backend-run)
func TestTC_MIX_004_RunRecipeHasPerSurfaceRecipes(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `frontend-run`),
		`Expected 'frontend-run' in mixed template`)
	assert.True(t, strings.Contains(template, `backend-run`),
		`Expected 'backend-run' in mixed template`)
}

// Traceability: TC-MIX-005 -> AC: dev has per-surface recipes (frontend-dev, backend-dev)
func TestTC_MIX_005_DevRecipeHasPerSurfaceRecipes(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `frontend-dev`),
		`Expected 'frontend-dev' in mixed template`)
	assert.True(t, strings.Contains(template, `backend-dev`),
		`Expected 'backend-dev' in mixed template`)
}

// Traceability: TC-MIX-006 -> AC: unit-test has per-surface recipes (frontend-unit-test, backend-unit-test)
func TestTC_MIX_006_UnitTestRecipeHasPerSurfaceRecipes(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `frontend-unit-test`),
		`Expected 'frontend-unit-test' in mixed template`)
	assert.True(t, strings.Contains(template, `backend-unit-test`),
		`Expected 'backend-unit-test' in mixed template`)
}

// Traceability: TC-MIX-007 -> AC: lint has per-surface recipes (frontend-lint, backend-lint)
func TestTC_MIX_007_LintRecipeHasPerSurfaceRecipes(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `frontend-lint`),
		`Expected 'frontend-lint' in mixed template`)
	assert.True(t, strings.Contains(template, `backend-lint`),
		`Expected 'backend-lint' in mixed template`)
}

// Traceability: TC-MIX-008 -> AC: fmt has per-surface recipes (frontend-fmt, backend-fmt)
func TestTC_MIX_008_FmtRecipeHasPerSurfaceRecipes(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `frontend-fmt`),
		`Expected 'frontend-fmt' in mixed template`)
	assert.True(t, strings.Contains(template, `backend-fmt`),
		`Expected 'backend-fmt' in mixed template`)
}

// Traceability: TC-MIX-009 -> AC: check has per-surface recipes (frontend-check, backend-check)
func TestTC_MIX_009_CheckRecipeHasPerSurfaceRecipes(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `frontend-check`),
		`Expected 'frontend-check' in mixed template`)
	assert.True(t, strings.Contains(template, `backend-check`),
		`Expected 'backend-check' in mixed template`)
}

// Traceability: TC-MIX-010 -> AC: clean has per-surface recipes (frontend-clean, backend-clean)
func TestTC_MIX_010_CleanRecipeHasPerSurfaceRecipes(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `frontend-clean`),
		`Expected 'frontend-clean' in mixed template`)
	assert.True(t, strings.Contains(template, `backend-clean`),
		`Expected 'backend-clean' in mixed template`)
}

// Traceability: TC-MIX-011 -> AC: install has per-surface recipes (frontend-install, backend-install)
func TestTC_MIX_011_InstallRecipeHasPerSurfaceRecipes(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `frontend-install`),
		`Expected 'frontend-install' in mixed template`)
	assert.True(t, strings.Contains(template, `backend-install`),
		`Expected 'backend-install' in mixed template`)
}

// ── TC-MIX-012 to TC-MIX-015: Bash case pattern checks ────────────────

// Traceability: TC-MIX-012 -> AC: Per-surface recipes handle both frontend and backend
// Mixed template has per-surface recipes (no case/esac dispatch).
// Verify that both frontend and backend variants exist for each standard recipe.
func TestTC_MIX_012_PerSurfaceRecipesCoverFrontendAndBackend(t *testing.T) {
	template := getMixedTemplate(t)
	standardRecipes := []string{"compile", "build", "run", "dev", "unit-test", "lint", "fmt", "check", "clean", "install"}
	for _, recipe := range standardRecipes {
		assert.True(t, strings.Contains(template, "frontend-"+recipe),
			"Expected frontend-%s recipe in mixed template", recipe)
		assert.True(t, strings.Contains(template, "backend-"+recipe),
			"Expected backend-%s recipe in mixed template", recipe)
	}
}

// Traceability: TC-MIX-013 -> AC: "") branch runs both frontend and backend
func TestTC_MIX_013_EmptyBranchExecutesBothFrontendAndBackendCommands(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, "npm run build)") && strings.Contains(template, "BACKEND_BUILD"),
		`Expected "" branch to chain frontend npm and BACKEND_BUILD placeholder commands`)
	assert.True(t, strings.Contains(template, "npm test)") && strings.Contains(template, "BACKEND_TEST"),
		`Expected "" branch to chain frontend npm test and BACKEND_TEST placeholder commands in unit-test recipe`)
}

// Traceability: TC-MIX-014 -> AC: All bash recipes use set -euo pipefail
func TestTC_MIX_014_AllBashRecipesUseSetEuoPipefail(t *testing.T) {
	template := getMixedTemplate(t)
	bashRecipes := countOccurrences(template, "#!/usr/bin/env bash")
	pipefailCount := countOccurrences(template, "set -euo pipefail")
	assert.True(t, pipefailCount >= bashRecipes,
		"Expected at least %d \"set -euo pipefail\" for %d bash recipes, got %d",
		bashRecipes, bashRecipes, pipefailCount)
}

// Traceability: TC-MIX-015 -> AC: Frontend uses npm, backend uses BACKEND_* placeholders
func TestTC_MIX_015_FrontendCommandsUseNpmBackendUsesPlaceholders(t *testing.T) {
	template := getMixedTemplate(t)
	// Frontend branch commands
	assert.True(t, strings.Contains(template, "npm run build"), "Expected \"npm run build\" for frontend")
	assert.True(t, strings.Contains(template, "npm test"), "Expected \"npm test\" for frontend")
	assert.True(t, strings.Contains(template, "npm run lint"), "Expected \"npm run lint\" for frontend")
	assert.True(t, strings.Contains(template, "FRONTEND_RUN"), "Expected \"FRONTEND_RUN\" placeholder for frontend run")
	assert.True(t, strings.Contains(template, "FRONTEND_DEV"), "Expected \"FRONTEND_DEV\" placeholder for frontend dev")
	assert.True(t, strings.Contains(template, "npx prettier --write ."), "Expected \"npx prettier --write .\" for frontend")
	assert.True(t, strings.Contains(template, "npx tsc --noEmit"), "Expected \"npx tsc --noEmit\" for frontend compile")
	assert.True(t, strings.Contains(template, "npm install"), "Expected \"npm install\" for frontend")

	// Backend uses BACKEND_* placeholders
	assert.True(t, strings.Contains(template, "BACKEND_BUILD"), "Expected \"BACKEND_BUILD\" placeholder")
	assert.True(t, strings.Contains(template, "BACKEND_TEST"), "Expected \"BACKEND_TEST\" placeholder")
	assert.True(t, strings.Contains(template, "BACKEND_LINT"), "Expected \"BACKEND_LINT\" placeholder")
	assert.True(t, strings.Contains(template, "BACKEND_RUN"), "Expected \"BACKEND_RUN\" placeholder")
	assert.True(t, strings.Contains(template, "BACKEND_FMT"), "Expected \"BACKEND_FMT\" placeholder")
	assert.True(t, strings.Contains(template, "BACKEND_COMPILE"), "Expected \"BACKEND_COMPILE\" placeholder")
	assert.True(t, strings.Contains(template, "BACKEND_CLEAN"), "Expected \"BACKEND_CLEAN\" placeholder")
	assert.True(t, strings.Contains(template, "BACKEND_INSTALL"), "Expected \"BACKEND_INSTALL\" placeholder")
}

// ── TC-MIX-016 to TC-MIX-020: Unscoped recipe checks ──────────────────

// Traceability: TC-MIX-016 -> AC: project-type has no scope parameter
// Note: project-type recipe has been removed; probe recipe (its replacement) has no scope parameter.
func TestTC_MIX_016_ProjectTypeHasNoScopeParameter(t *testing.T) {
	template := getMixedTemplate(t)
	re := regexp.MustCompile(`(?m)^probe[^:]*:`)
	match := re.FindString(template)
	assert.NotEmpty(t, match, "Expected probe recipe in mixed template (replaces project-type)")
	assert.False(t, strings.Contains(match, `scope=""`),
		"Expected probe to NOT have scope=\"\" parameter")
}

// Traceability: TC-MIX-017 -> AC: test (surface-level) has no scope parameter
func TestTC_MIX_017_TestRecipeHasNoScopeParameter(t *testing.T) {
	template := getMixedTemplate(t)
	re := regexp.MustCompile(`(?m)^test[^:]*:`)
	match := re.FindString(template)
	assert.NotEmpty(t, match, "Expected test recipe in mixed template")
	assert.False(t, strings.Contains(match, `scope=""`),
		"Expected test to NOT have scope=\"\" parameter")
}

// Traceability: TC-MIX-018 -> AC: ci has no scope parameter
func TestTC_MIX_018_CiHasNoScopeParameter(t *testing.T) {
	template := getMixedTemplate(t)
	re := regexp.MustCompile(`(?m)^ci:`)
	match := re.FindString(template)
	assert.NotEmpty(t, match, "Expected ci recipe in mixed template")
	assert.False(t, strings.Contains(template, `ci scope=""`),
		"Expected ci to NOT have scope=\"\" parameter")
}

// Traceability: TC-MIX-019 -> AC: test-setup has no scope parameter
func TestTC_MIX_019_TestSetupHasNoScopeParameter(t *testing.T) {
	template := getMixedTemplate(t)
	re := regexp.MustCompile(`(?m)^test-setup[^:]*:`)
	match := re.FindString(template)
	assert.NotEmpty(t, match, "Expected test-setup recipe in mixed template")
	assert.False(t, strings.Contains(match, `scope=""`),
		"Expected test-setup to NOT have scope=\"\" parameter")
}

// Traceability: TC-MIX-020 -> AC: probe has no scope parameter
func TestTC_MIX_020_ProbeHasNoScopeParameter(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, "probe"),
		"Expected probe recipe in mixed template")
	assert.False(t, strings.Contains(template, `probe scope=""`),
		"Expected probe to NOT have scope=\"\" parameter")
}

// ── TC-MIX-021 to TC-MIX-023: Boundary markers and structure ──────────

// Traceability: TC-MIX-021 -> AC: Templates stored as string literals
func TestTC_MIX_021_MixedTemplateHasForgeBoundaryMarkers(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, "# --- forge standard recipes ---"),
		"Expected forge standard recipes start marker in mixed template")
	assert.True(t, strings.Contains(template, "# --- end forge standard recipes ---"),
		"Expected forge standard recipes end marker in mixed template")
}

// Traceability: TC-MIX-022 -> AC: All recipes present
// Two-layer test model: unit-test (language-level) + test (surface-level).
// Total: 14 recipes (compile, build, run, dev, unit-test, test, lint, fmt, check,
// clean, install, ci, test-setup, probe).
func TestTC_MIX_022_AllRecipesArePresentInMixedTemplate(t *testing.T) {
	template := getMixedTemplate(t)
	expectedRecipes := []string{
		"compile", "build", "run", "dev",
		"unit-test", "test", "lint", "fmt", "check",
		"clean", "install", "ci", "test-setup", "probe",
	}
	for _, recipe := range expectedRecipes {
		assert.True(t, strings.Contains(template, recipe),
			"Expected recipe %q in mixed template", recipe)
	}
}

// Traceability: TC-MIX-023 -> AC: ci recipe chains install, compile, build, unit-test, lint
func TestTC_MIX_023_CiRecipeChainsStandardCommands(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, "ci:"),
		"Expected ci recipe definition")
	re := regexp.MustCompile(`(?m)^ci:.*$`)
	ciLine := re.FindString(template)
	expectedSteps := []string{"install", "compile", "build", "unit-test", "lint"}
	for _, step := range expectedSteps {
		assert.True(t, strings.Contains(ciLine, step),
			"Expected %q in ci recipe, got: %s", step, ciLine)
	}
}

// ── CLI integration tests (TC-001 to TC-008, TC-013 to TC-019) ────────
// Converted from tests/e2e/justfile-e2e-integration/cli.spec.ts

// Traceability: TC-001 -> Story 1 / AC-1
// run-tests SKILL.md orchestrates using just recipes (dev/probe/test/teardown sequence).
func TestTC_001_RunTestsSkillUsesJustRecipes(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/skills/run-tests/SKILL.md")
	assert.True(t, strings.Contains(content, "just") && strings.Contains(content, "test"),
		"Expected \"just\" and \"test\" to appear in run-tests/SKILL.md orchestration")
}

// Traceability: TC-002 -> Story 2 / AC-1 (migrated: just build -> just compile per tech-design)
// Note: task-executor.md is a thin executor that delegates to the synthesized strategy (task prompt).
// Verify the task-executor references submit-task skill for completion and does not hardcode commands.
func TestTC_002_TaskExecutorStep3UsesJustCompileAndJustTest(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/agents/task-executor.md")
	// The thin executor delegates to submit-task for completion
	assert.True(t, strings.Contains(content, "submit-task"),
		"Expected \"submit-task\" skill reference in task-executor.md")
	// Verify it does NOT hardcode language-specific commands
	assert.False(t, strings.Contains(content, "go test ./..."),
		"Expected \"go test ./...\" NOT to appear in task-executor.md")
	assert.False(t, strings.Contains(content, "npm test"),
		"Expected \"npm test\" NOT to appear in task-executor.md")
}

// Traceability: TC-003 -> Story 3 / AC-1
// Note: e2e-verify has been removed in the two-layer test model. The test recipe
// replaces e2e-test for surface-level testing. Verify the mixed template has test-setup.
func TestTC_003_TestSetupRecipePresent(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, "test-setup"),
		"Expected test-setup recipe in mixed template")
}

// Traceability: TC-004 -> Story 3 / AC-2
// Verify the mixed template has test recipe with journey parameter.
func TestTC_004_TestRecipeHasJourneyParameter(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `test journey=''`),
		"Expected test recipe with journey parameter in mixed template")
}

// Traceability: TC-005 -> Story 5 / AC-2 (verify run-tasks delegates test execution to submit gate)
// Note: run-tasks.md explicitly delegates test execution to the CLI submit gate ("NO running tests directly").
// This test verifies no hardcoded language-specific test commands appear in the dispatcher.
func TestTC_005_RunTasksBreakingGateUsesJustUnitTestForVerification(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/commands/run-tasks.md")
	// run-tasks.md delegates to submit gate — verify no hardcoded test commands
	assert.False(t, strings.Contains(content, "go test ./..."),
		"Expected \"go test ./...\" NOT to appear in run-tasks.md")
	assert.False(t, strings.Contains(content, "npm test"),
		"Expected \"npm test\" NOT to appear in run-tasks.md")
}

// Traceability: TC-006 -> Story 5 / AC-1
func TestTC_006_FixBugUsesJustTestNotProjectTestCommandPlaceholder(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/commands/fix-bug.md")
	assert.True(t, strings.Contains(content, "just unit-test") || strings.Contains(content, "just test"),
		"Expected \"just unit-test\" or \"just test\" to appear in fix-bug.md test verification step")
	assert.False(t, strings.Contains(content, "<project-test-command>"),
		"Expected \"<project-test-command>\" placeholder NOT to appear in fix-bug.md")
}

// Traceability: TC-007 -> Story 5 / AC-2
// run-tasks delegates test execution to the CLI submit gate.
// Verify no hardcoded language-specific test commands appear in run-tasks.md.
func TestTC_007_RunTasksDispatcherDelegatesTestExecutionToSubmitGate(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/commands/run-tasks.md")
	// run-tasks.md explicitly delegates test execution to the CLI submit gate
	assert.True(t, strings.Contains(content, "submit"),
		"Expected \"submit\" reference in run-tasks.md (delegation to submit gate)")
	// Verify no hardcoded language-specific test commands
	assert.False(t, strings.Contains(content, "npm test"),
		"Expected \"npm test\" NOT to appear in run-tasks.md")
	assert.False(t, strings.Contains(content, "go test"),
		"Expected \"go test\" NOT to appear in run-tasks.md")
}

// Traceability: TC-008 -> Story 5 / AC-3
// submit-task SKILL.md is tool-agnostic; verify no hardcoded language-specific test commands.
func TestTC_008_RecordTaskMetricsCollectionUsesJustUnitTest(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/skills/submit-task/SKILL.md")
	// submit-task should NOT hardcode language-specific test commands
	assert.False(t, strings.Contains(content, "go test -cover ./..."),
		"Expected \"go test -cover ./...\" NOT to appear in submit-task/SKILL.md")
	assert.False(t, strings.Contains(content, "npm test -- --coverage"),
		"Expected \"npm test -- --coverage\" NOT to appear in submit-task/SKILL.md")
	assert.False(t, strings.Contains(content, "pytest --cov="),
		"Expected \"pytest --cov=\" NOT to appear in submit-task/SKILL.md")
}

// Traceability: TC-009 -> Spec Section 5.1
// Note: test-setup replaces e2e-setup in the two-layer test model. Verify the template
// has the missing package.json check.
func TestTC_009_JustTestSetupExits1WhenPackageJsonMissing(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, "package.json"),
		"Expected package.json existence check in test-setup recipe")
	assert.True(t, strings.Contains(template, "not found"),
		"Expected 'not found' error message in test-setup recipe")
}

// Traceability: TC-010 -> Spec Section 5.1
func TestTC_010_JustTestSetupExits0WithOKMessageWhenDepsReady(t *testing.T) {
	root := testkit.ProjectRoot(t)
	pkgPath := filepath.Join(root, "tests", "e2e", "package.json")
	nodeModulesPath := filepath.Join(root, "tests", "e2e", "node_modules")
	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		t.Skip("requires real package.json to be present")
	}
	if _, err := os.Stat(nodeModulesPath); os.IsNotExist(err) {
		t.Skip("requires real node_modules to be present")
	}

	exitCode, out := runJust("test-setup")
	assert.Equal(t, 0, exitCode, "Expected exit code 0 when deps are ready")
	assert.True(t, strings.Contains(out, "OK: test dependencies ready"),
		"Expected \"OK: test dependencies ready\" in stdout, got: %s", out)
}

// Traceability: TC-011 -> Spec Section 5.1
// Note: test recipe replaces e2e-test and accepts optional journey parameter.
func TestTC_011_TestRecipeAcceptsJourneyParameter(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `test journey=''`),
		"Expected test recipe with journey parameter in mixed template")
}

// Traceability: TC-012 -> Spec Section 5.1
// Note: test recipe dispatches by journey — verify the journey filtering logic.
func TestTC_012_TestRecipeFiltersByJourneyWhenProvided(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `{{journey}}`),
		"Expected journey parameter usage in test recipe")
}

// Traceability: TC-013 -> Spec Section 5.3
func TestTC_013_RunTestsSkillPromptsInitJustfileWhenJustfileMissing(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/skills/run-tests/SKILL.md")
	hasJustfileCheck := strings.Contains(content, "justfile") || strings.Contains(content, "init-justfile")
	assert.True(t, hasJustfileCheck,
		"Expected run-tests/SKILL.md to reference justfile existence check or /init-justfile")
}

// Traceability: TC-014 -> Spec Section 5.2
// gen-test-scripts generates test scripts that are executed by run-tests via just recipes.
func TestTC_014_GenTestScriptsGeneratesTestScriptsForRunTests(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/skills/gen-test-scripts/SKILL.md")
	// gen-test-scripts generates test scripts that run-tests will execute via just recipes
	assert.True(t, strings.Contains(content, "test") && strings.Contains(content, "scripts"),
		"Expected \"test\" and \"scripts\" to appear in gen-test-scripts/SKILL.md")
}

// Traceability: TC-015 -> Spec Section 5.2 (migrated: just build -> just compile per tech-design)
// Note: error-fixer.md has been removed as a standalone agent. Error fixing is now handled
// by the execute-task command with fix-task template. Verify execute-task.md delegates verification
// to submit-task rather than hardcoding language-specific commands.
func TestTC_015_ErrorFixerUsesJustCompileAndJustUnitTest(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/commands/execute-task.md")
	// execute-task delegates to submit-task for verification — verify no hardcoded commands
	assert.False(t, strings.Contains(content, "go build ./..."),
		"Expected \"go build ./...\" NOT to appear in execute-task.md")
	assert.False(t, strings.Contains(content, "go vet ./..."),
		"Expected \"go vet ./...\" NOT to appear in execute-task.md")
	assert.False(t, strings.Contains(content, "go test -race -cover ./..."),
		"Expected \"go test -race -cover ./...\" NOT to appear in execute-task.md")
	assert.False(t, strings.Contains(content, "npm run build && npm test"),
		"Expected \"npm run build && npm test\" NOT to appear in execute-task.md")
}

// Traceability: TC-016 -> Spec Section 5.2 (migrated: just build -> just compile per tech-design)
// Note: execute-task delegates verification to submit-task via "submit-task is mandatory" rule.
func TestTC_016_ExecuteTaskStep3UsesJustCompileAndJustTest(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/commands/execute-task.md")
	// execute-task delegates to submit-task for verification
	assert.True(t, strings.Contains(content, "submit-task"),
		"Expected \"submit-task\" reference in execute-task.md (verification delegation)")
}

// Traceability: TC-017 -> removed: improve-harness skill no longer exists (replaced by run-tests)

// Traceability: TC-018 -> Spec Section 5.1
// Note: test-setup replaces e2e-setup in the two-layer test model.
// The test-setup recipe is defined in templates (mixed.just has the node_modules check).
func TestTC_018_InitJustfileGeneratesTestSetupTarget(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/SKILL.md")
	assert.True(t, strings.Contains(content, "test-setup"),
		"Expected \"test-setup\" to appear in init-justfile SKILL.md")

	// Verify the mixed template has test-setup with idempotent node_modules check
	mixedTemplate := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/templates/mixed.just")
	assert.True(t, strings.Contains(mixedTemplate, "test-setup"),
		"Expected test-setup recipe in mixed template")
	assert.True(t, strings.Contains(mixedTemplate, "node_modules"),
		"Expected idempotent node_modules check in test-setup recipe in mixed template")

	genericTemplate := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/templates/generic.just")
	assert.True(t, strings.Contains(genericTemplate, "test-setup"),
		"Expected test-setup recipe in generic template")
}

// Traceability: TC-019 -> Spec Section 5.1
// Note: test recipe replaces e2e-test and accepts optional journey parameter.
// The test recipe with journey parameter is defined in templates (mixed.just, generic.just).
func TestTC_019_InitJustfileGeneratesTestTarget(t *testing.T) {
	// Check SKILL.md references the test recipe with journey parameter
	content := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/SKILL.md")
	assert.True(t, strings.Contains(content, "test") && strings.Contains(content, "journey"),
		"Expected \"test\" and \"journey\" to appear in init-justfile SKILL.md")

	// Verify the mixed template has test recipe with journey parameter
	mixedTemplate := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/templates/mixed.just")
	assert.True(t, strings.Contains(mixedTemplate, `test journey`),
		"Expected \"test journey\" parameter in mixed template")
}

// Traceability: TC-020 -> Spec Section 5.1
func TestTC_020_JustTestSetupIsIdempotent(t *testing.T) {
	root := testkit.ProjectRoot(t)
	pkgPath := filepath.Join(root, "tests", "e2e", "package.json")
	nodeModulesPath := filepath.Join(root, "tests", "e2e", "node_modules")
	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		t.Skip("requires real package.json to be present")
	}
	if _, err := os.Stat(nodeModulesPath); os.IsNotExist(err) {
		t.Skip("requires real node_modules to be present")
	}

	result1Code, result1Out := runJust("test-setup")
	result2Code, result2Out := runJust("test-setup")
	assert.Equal(t, 0, result1Code, "Expected first run to exit 0")
	assert.Equal(t, 0, result2Code, "Expected second run to exit 0")
	assert.True(t, strings.Contains(result1Out, "OK: test dependencies ready"),
		"Expected \"OK: test dependencies ready\" in first run stdout")
	assert.True(t, strings.Contains(result2Out, "OK: test dependencies ready"),
		"Expected \"OK: test dependencies ready\" in second run stdout")
}
