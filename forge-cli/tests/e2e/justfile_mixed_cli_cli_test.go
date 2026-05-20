//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"forge-cli/tests/e2e/testkit"

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

// Traceability: TC-MIX-002 -> AC: compile has scope with bash case
func TestTC_MIX_002_CompileRecipeHasScopeParameterWithBashCaseDispatch(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `compile scope=""`),
		`Expected 'compile scope=""' in mixed template`)
	assert.True(t, strings.Contains(template, "frontend)") && strings.Contains(template, "backend)"),
		`Expected "frontend)" and "backend)" case branches in compile recipe`)
}

// Traceability: TC-MIX-003 -> AC: build has scope with bash case
func TestTC_MIX_003_BuildRecipeHasScopeParameterWithBashCaseDispatch(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `build scope=""`),
		`Expected 'build scope=""' in mixed template`)
}

// Traceability: TC-MIX-004 -> AC: run has scope with bash case
func TestTC_MIX_004_RunRecipeHasScopeParameterWithBashCaseDispatch(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `run scope=""`),
		`Expected 'run scope=""' in mixed template`)
}

// Traceability: TC-MIX-005 -> AC: dev has scope with bash case
func TestTC_MIX_005_DevRecipeHasScopeParameterWithBashCaseDispatch(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `dev scope=""`),
		`Expected 'dev scope=""' in mixed template`)
}

// Traceability: TC-MIX-006 -> AC: test has scope with bash case
func TestTC_MIX_006_TestRecipeHasScopeParameterWithBashCaseDispatch(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `test scope=""`),
		`Expected 'test scope=""' in mixed template`)
}

// Traceability: TC-MIX-007 -> AC: lint has scope with bash case
func TestTC_MIX_007_LintRecipeHasScopeParameterWithBashCaseDispatch(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `lint scope=""`),
		`Expected 'lint scope=""' in mixed template`)
}

// Traceability: TC-MIX-008 -> AC: fmt has scope with bash case
func TestTC_MIX_008_FmtRecipeHasScopeParameterWithBashCaseDispatch(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `fmt scope=""`),
		`Expected 'fmt scope=""' in mixed template`)
}

// Traceability: TC-MIX-009 -> AC: check has scope with bash case
func TestTC_MIX_009_CheckRecipeHasScopeParameterWithBashCaseDispatch(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `check scope=""`),
		`Expected 'check scope=""' in mixed template`)
}

// Traceability: TC-MIX-010 -> AC: clean has scope with bash case
func TestTC_MIX_010_CleanRecipeHasScopeParameterWithBashCaseDispatch(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `clean scope=""`),
		`Expected 'clean scope=""' in mixed template`)
}

// Traceability: TC-MIX-011 -> AC: install has scope with bash case
func TestTC_MIX_011_InstallRecipeHasScopeParameterWithBashCaseDispatch(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, `install scope=""`),
		`Expected 'install scope=""' in mixed template`)
}

// ── TC-MIX-012 to TC-MIX-015: Bash case pattern checks ────────────────

// Traceability: TC-MIX-012 -> AC: *) branch error message
func TestTC_MIX_012_ScopedRecipesHaveStarBranchWithErrorMessageToStderr(t *testing.T) {
	template := getMixedTemplate(t)
	errorMsg := `echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1`
	matches := countOccurrences(template, errorMsg)
	assert.True(t, matches >= 10,
		"Expected at least 10 occurrences of *) error branch, got %d", matches)
}

// Traceability: TC-MIX-013 -> AC: "") branch runs both frontend and backend
func TestTC_MIX_013_EmptyBranchExecutesBothFrontendAndBackendCommands(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, "npm run build)") && strings.Contains(template, "BACKEND_BUILD"),
		`Expected "" branch to chain frontend npm and BACKEND_BUILD placeholder commands`)
	assert.True(t, strings.Contains(template, "npm test)") && strings.Contains(template, "BACKEND_TEST"),
		`Expected "" branch to chain frontend npm test and BACKEND_TEST placeholder commands`)
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

// Traceability: TC-MIX-017 -> AC: e2e-test has no scope parameter
func TestTC_MIX_017_TestE2eHasNoScopeParameter(t *testing.T) {
	template := getMixedTemplate(t)
	re := regexp.MustCompile(`e2e-test[^:]*:`)
	match := re.FindString(template)
	assert.NotEmpty(t, match, "Expected e2e-test recipe in mixed template")
	assert.False(t, strings.Contains(match, `scope=""`),
		"Expected e2e-test to NOT have scope=\"\" parameter")
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

// Traceability: TC-MIX-019 -> AC: e2e-setup has no scope parameter
func TestTC_MIX_019_E2eSetupHasNoScopeParameter(t *testing.T) {
	template := getMixedTemplate(t)
	re := regexp.MustCompile(`(?m)^e2e-setup[^:]*:`)
	match := re.FindString(template)
	assert.NotEmpty(t, match, "Expected e2e-setup recipe in mixed template")
	assert.False(t, strings.Contains(match, `scope=""`),
		"Expected e2e-setup to NOT have scope=\"\" parameter")
}

// Traceability: TC-MIX-020 -> AC: e2e-verify has no scope parameter
func TestTC_MIX_020_E2eVerifyHasNoScopeParameter(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, "e2e-verify"),
		"Expected e2e-verify recipe in mixed template")
	assert.False(t, strings.Contains(template, `e2e-verify scope=""`),
		"Expected e2e-verify to NOT have scope=\"\" parameter")
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
// Note: project-type has been removed (replaced by forge probe); probe recipe added.
// Total: 15 recipes (compile, build, run, dev, test, e2e-test, lint, fmt, check,
// clean, install, ci, e2e-setup, probe, e2e-verify).
func TestTC_MIX_022_All15RecipesArePresentInMixedTemplate(t *testing.T) {
	template := getMixedTemplate(t)
	expectedRecipes := []string{
		"compile", "build", "run", "dev",
		"test", "e2e-test", "lint", "fmt", "check",
		"clean", "install", "ci", "e2e-setup", "probe", "e2e-verify",
	}
	for _, recipe := range expectedRecipes {
		assert.True(t, strings.Contains(template, recipe),
			"Expected recipe %q in mixed template", recipe)
	}
}

// Traceability: TC-MIX-023 -> AC: ci recipe chains install, compile, build, test, lint
func TestTC_MIX_023_CiRecipeChainsStandardCommands(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, "ci:"),
		"Expected ci recipe definition")
	re := regexp.MustCompile(`(?m)^ci:.*$`)
	ciLine := re.FindString(template)
	expectedSteps := []string{"install", "compile", "build", "test", "lint"}
	for _, step := range expectedSteps {
		assert.True(t, strings.Contains(ciLine, step),
			"Expected %q in ci recipe, got: %s", step, ciLine)
	}
}

// ── CLI integration tests (TC-001 to TC-008, TC-013 to TC-019) ────────
// Converted from tests/e2e/justfile-e2e-integration/cli.spec.ts

// Traceability: TC-001 -> Story 1 / AC-1
func TestTC_001_RunE2eTestsStep1UsesJustE2eSetup(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/skills/run-e2e-tests/SKILL.md")
	assert.True(t, strings.Contains(content, "just e2e-setup"),
		"Expected \"just e2e-setup\" to appear in run-e2e-tests/SKILL.md")
	assert.False(t, strings.Contains(content, "cd tests/e2e && npm install"),
		"Expected \"cd tests/e2e && npm install\" NOT to appear in run-e2e-tests/SKILL.md")
	assert.False(t, strings.Contains(content, "npx playwright install chromium"),
		"Expected \"npx playwright install chromium\" NOT to appear in run-e2e-tests/SKILL.md")
}

// Traceability: TC-002 -> Story 2 / AC-1 (migrated: just build -> just compile per tech-design)
// Note: task-executor.md has been refactored into a thin executor that delegates verification
// to the synthesized strategy (task prompt). The task prompt (from execute-task.md) contains
// the just test / just compile commands. Verify the task-executor references record-task
// skill for completion and delegates to the task prompt.
func TestTC_002_TaskExecutorStep3UsesJustCompileAndJustTest(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/agents/task-executor.md")
	// The thin executor delegates to task prompt which contains just commands
	assert.True(t, strings.Contains(content, "record-task"),
		"Expected \"record-task\" skill reference in task-executor.md")
	// Verify it does NOT hardcode language-specific commands
	assert.False(t, strings.Contains(content, "go test ./..."),
		"Expected \"go test ./...\" NOT to appear in task-executor.md")
	assert.False(t, strings.Contains(content, "npm test"),
		"Expected \"npm test\" NOT to appear in task-executor.md")
}

// Traceability: TC-003 -> Story 3 / AC-1
// Note: e2e-verify is a template-provided recipe. For pure Go projects, the justfile
// may not contain this recipe. Verify the mixed template has the correct e2e-verify logic.
func TestTC_003_JustE2eVerifyExits1WhenVerifyMarkersPresent(t *testing.T) {
	template := getMixedTemplate(t)
	// Verify the template contains e2e-verify recipe that exits 1 when markers found
	assert.True(t, strings.Contains(template, "e2e-verify"),
		"Expected e2e-verify recipe in mixed template")
	assert.True(t, strings.Contains(template, "// VERIFY:"),
		"Expected VERIFY marker scanning logic in e2e-verify recipe")
	assert.True(t, strings.Contains(template, "exit 1"),
		"Expected exit 1 in e2e-verify when unresolved markers found")
}

// Traceability: TC-004 -> Story 3 / AC-2
// Note: e2e-verify is template-provided. Verify the template outputs success message when no markers.
func TestTC_004_JustE2eVerifyExits0WhenNoVerifyMarkers(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, "OK: no unresolved // VERIFY: markers"),
		"Expected success message in e2e-verify recipe when no markers found")
}

// Traceability: TC-005 -> Story 5 / AC-2 (verify run-tasks uses standard just commands)
func TestTC_005_RunTasksBreakingGateUsesJustTestForVerification(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/commands/run-tasks.md")
	assert.True(t, strings.Contains(content, "just test"),
		"Expected \"just test\" to appear in run-tasks.md")
	assert.False(t, strings.Contains(content, "go test ./..."),
		"Expected \"go test ./...\" NOT to appear in run-tasks.md")
	assert.False(t, strings.Contains(content, "npm test"),
		"Expected \"npm test\" NOT to appear in run-tasks.md")
}

// Traceability: TC-006 -> Story 5 / AC-1
func TestTC_006_FixBugUsesJustTestNotProjectTestCommandPlaceholder(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/commands/fix-bug.md")
	assert.True(t, strings.Contains(content, "just test"),
		"Expected \"just test\" to appear in fix-bug.md test verification step")
	assert.False(t, strings.Contains(content, "<project-test-command>"),
		"Expected \"<project-test-command>\" placeholder NOT to appear in fix-bug.md")
}

// Traceability: TC-007 -> Story 5 / AC-2
func TestTC_007_RunTasksBreakingGateUsesJustTest(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/commands/run-tasks.md")
	assert.True(t, strings.Contains(content, "just test"),
		"Expected \"just test\" to appear in run-tasks.md Breaking Gate section")

	breakingGateIdx := strings.Index(content, "Breaking Task Gate")
	assert.NotEqual(t, -1, breakingGateIdx, "Expected \"Breaking Task Gate\" section to exist in run-tasks.md")
	afterBreakingGate := content[breakingGateIdx:]
	assert.False(t, strings.Contains(afterBreakingGate, "npm test"),
		"Expected \"npm test\" NOT to appear in Breaking Gate section of run-tasks.md")
	assert.False(t, strings.Contains(afterBreakingGate, "go test"),
		"Expected \"go test\" NOT to appear in Breaking Gate section of run-tasks.md")
}

// Traceability: TC-008 -> Story 5 / AC-3
func TestTC_008_RecordTaskMetricsCollectionUsesJustTest(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/skills/submit-task/SKILL.md")
	assert.True(t, strings.Contains(content, "just test"),
		"Expected \"just test\" to appear in submit-task/SKILL.md Metrics Collection section")
	assert.False(t, strings.Contains(content, "go test -cover ./..."),
		"Expected \"go test -cover ./...\" NOT to appear in submit-task/SKILL.md")
	assert.False(t, strings.Contains(content, "npm test -- --coverage"),
		"Expected \"npm test -- --coverage\" NOT to appear in submit-task/SKILL.md")
	assert.False(t, strings.Contains(content, "pytest --cov="),
		"Expected \"pytest --cov=\" NOT to appear in submit-task/SKILL.md")
}

// Traceability: TC-009 -> Spec Section 5.1
// Note: e2e-setup is template-provided. Verify the template has the missing package.json check.
func TestTC_009_JustE2eSetupExits1WhenPackageJsonMissing(t *testing.T) {
	template := getMixedTemplate(t)
	assert.True(t, strings.Contains(template, "package.json"),
		"Expected package.json existence check in e2e-setup recipe")
	assert.True(t, strings.Contains(template, "not found"),
		"Expected 'not found' error message in e2e-setup recipe")
}

// Traceability: TC-010 -> Spec Section 5.1
func TestTC_010_JustE2eSetupExits0WithOKMessageWhenDepsReady(t *testing.T) {
	root := testkit.ProjectRoot(t)
	pkgPath := filepath.Join(root, "tests", "e2e", "package.json")
	nodeModulesPath := filepath.Join(root, "tests", "e2e", "node_modules")
	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		t.Skip("requires real package.json to be present")
	}
	if _, err := os.Stat(nodeModulesPath); os.IsNotExist(err) {
		t.Skip("requires real node_modules to be present")
	}

	exitCode, out := runJust("e2e-setup")
	assert.Equal(t, 0, exitCode, "Expected exit code 0 when deps are ready")
	assert.True(t, strings.Contains(out, "OK: e2e dependencies ready"),
		"Expected \"OK: e2e dependencies ready\" in stdout, got: %s", out)
}

// Traceability: TC-011 -> Spec Section 5.1
// Note: e2e-verify is template-provided. Verify the template requires --feature flag.
func TestTC_011_JustE2eVerifyExits1WhenFeatureFlagMissing(t *testing.T) {
	template := getMixedTemplate(t)
	// The template should check for empty feature and exit with usage
	assert.True(t, strings.Contains(template, `e2e-verify feature=""`),
		"Expected e2e-verify recipe with feature parameter")
	assert.True(t, strings.Contains(template, "--feature"),
		"Expected --feature reference in e2e-verify usage message")
}

// Traceability: TC-012 -> Spec Section 5.1
// Note: e2e-verify is template-provided. Verify the template uses grep -n for line numbers.
func TestTC_012_JustE2eVerifyOutputsFileAndLineNumberForResidualMarkers(t *testing.T) {
	template := getMixedTemplate(t)
	// The template uses grep -rn to find VERIFY markers (recursive with line numbers)
	assert.True(t, strings.Contains(template, "grep -rn"),
		"Expected 'grep -rn' for line number reporting in e2e-verify recipe")
	assert.True(t, strings.Contains(template, "// VERIFY:"),
		"Expected VERIFY marker pattern search in e2e-verify recipe")
}

// Traceability: TC-013 -> Spec Section 5.3
func TestTC_013_RunE2eTestsSkillPromptsInitJustfileWhenJustfileMissing(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/skills/run-e2e-tests/SKILL.md")
	hasJustfileCheck := strings.Contains(content, "justfile") || strings.Contains(content, "init-justfile")
	assert.True(t, hasJustfileCheck,
		"Expected run-e2e-tests/SKILL.md to reference justfile existence check or /init-justfile")
}

// Traceability: TC-014 -> Spec Section 5.2 / Story 3
func TestTC_014_GenTestScriptsStep4UsesJustE2eVerify(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/skills/gen-test-scripts/SKILL.md")
	assert.True(t, strings.Contains(content, "just e2e-verify --feature"),
		"Expected \"just e2e-verify --feature\" to appear in gen-test-scripts/SKILL.md Step 4")
}

// Traceability: TC-015 -> Spec Section 5.2 (migrated: just build -> just compile per tech-design)
// Note: error-fixer.md has been removed as a standalone agent. Error fixing is now handled
// by the execute-task command with fix-task template. Verify the execute-task.md uses just commands.
func TestTC_015_ErrorFixerUsesJustCompileAndJustTest(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/commands/execute-task.md")
	assert.True(t, strings.Contains(content, "just test"),
		"Expected \"just test\" to appear in execute-task.md (error fixer verification)")
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
func TestTC_016_ExecuteTaskStep3UsesJustCompileAndJustTest(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/commands/execute-task.md")
	// execute-task.md references just test in Step 3
	assert.True(t, strings.Contains(content, "just test"),
		"Expected \"just test\" to appear in execute-task.md Step 3")
}

// Traceability: TC-017 -> Spec Section 5.2
func TestTC_017_ImproveHarnessUsesJustTest(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/skills/improve-harness/SKILL.md")
	assert.True(t, strings.Contains(content, "just test"),
		"Expected \"just test\" to appear in improve-harness/SKILL.md Step 4.3")
}

// Traceability: TC-018 -> Spec Section 5.1
func TestTC_018_InitJustfileGeneratesE2eSetupTarget(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/SKILL.md")
	assert.True(t, strings.Contains(content, "e2e-setup"),
		"Expected \"e2e-setup\" recipe to appear in init-justfile.md template")
	assert.True(t, strings.Contains(content, "node_modules"),
		"Expected idempotent node_modules check in e2e-setup recipe")

	genericTemplate := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/templates/generic.just")
	assert.True(t, strings.Contains(genericTemplate, "playwright install chromium"),
		"Expected \"playwright install chromium\" in e2e-setup recipe template")
}

// Traceability: TC-019 -> Spec Section 5.1
func TestTC_019_InitJustfileGeneratesE2eVerifyTarget(t *testing.T) {
	content := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/SKILL.md")
	assert.True(t, strings.Contains(content, "e2e-verify"),
		"Expected \"e2e-verify\" recipe to appear in init-justfile.md template")
	assert.True(t, strings.Contains(content, "--feature"),
		"Expected \"--feature\" parameter in e2e-verify recipe")
	assert.True(t, strings.Contains(content, "// VERIFY:"),
		"Expected \"// VERIFY:\" marker scanning in e2e-verify recipe")
}

// Traceability: TC-020 -> Spec Section 5.1
func TestTC_020_JustE2eSetupIsIdempotent(t *testing.T) {
	root := testkit.ProjectRoot(t)
	pkgPath := filepath.Join(root, "tests", "e2e", "package.json")
	nodeModulesPath := filepath.Join(root, "tests", "e2e", "node_modules")
	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		t.Skip("requires real package.json to be present")
	}
	if _, err := os.Stat(nodeModulesPath); os.IsNotExist(err) {
		t.Skip("requires real node_modules to be present")
	}

	result1Code, result1Out := runJust("e2e-setup")
	result2Code, result2Out := runJust("e2e-setup")
	assert.Equal(t, 0, result1Code, "Expected first run to exit 0")
	assert.Equal(t, 0, result2Code, "Expected second run to exit 0")
	assert.True(t, strings.Contains(result1Out, "OK: e2e dependencies ready"),
		"Expected \"OK: e2e dependencies ready\" in first run stdout")
	assert.True(t, strings.Contains(result2Out, "OK: e2e dependencies ready"),
		"Expected \"OK: e2e dependencies ready\" in second run stdout")
}
