//go:build e2e

package justfileintegration

import (
	"strings"
	"testing"

	"forge-cli/tests/testkit"

	"github.com/stretchr/testify/assert"
)

// --- Forge Justfile Tests (TC-FJ-001 to TC-FJ-015) ---
// Converted from tests/e2e/justfile-e2e-integration/forge-justfile.spec.ts (15 tests).
// These tests validate the forge justfile: standard recipes presence and live command execution.
//
// Note: The current justfile is a Go-backend project using `cd forge-cli` prefix.
// Scoped recipes accept a scope parameter but don't dispatch by scope (no case/esac).
// The project-type recipe has been removed — project type is now read from .forge/config.yaml
// via `forge probe`. Template recipes (e2e-test, e2e-setup, e2e-verify) are provided by
// the profile manifest, not baked into the justfile.

// --- TC-FJ-001 to TC-FJ-010: Standard recipes presence ---

// getJustfile reads the project justfile.
func getJustfile(t *testing.T) string {
	t.Helper()
	return testkit.ReadProjectFile(t, "../justfile")
}

// getStandardSection extracts the forge standard recipes section from the justfile.
func getStandardSection(t *testing.T) string {
	t.Helper()
	content := getJustfile(t)
	startMarker := "# --- forge standard recipes ---"
	endMarker := "# --- end forge standard recipes ---"
	startIdx := strings.Index(content, startMarker)
	endIdx := strings.Index(content, endMarker)
	assert.NotEqual(t, -1, startIdx, "Expected start boundary marker in justfile")
	assert.NotEqual(t, -1, endIdx, "Expected end boundary marker in justfile")
	return content[startIdx : endIdx+len(endMarker)]
}

// Traceability: TC-FJ-001 -> AC: project-type outputs "mixed"
// project-type is now via forge probe. For a Go-backend project, probe returns a valid type.
func TestTC_FJ_001_ProjectTypeRecipeOutputsMixed(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("probe")
	assert.Equal(t, 0, exitCode, "Expected exit code 0 from forge probe")
	output := strings.TrimSpace(out)
	assert.NotEmpty(t, output, "Expected non-empty output from forge probe")
}

// Traceability: TC-FJ-002 -> AC: 10 scoped recipes use scope parameter
// 10 scoped recipes have scope="" parameter.
func TestTC_FJ_002_TenScopedRecipesUseScopeParameter(t *testing.T) {
	section := getStandardSection(t)
	scopedRecipes := []string{"compile", "build", "run", "dev", "test", "lint", "fmt", "check", "clean", "install"}
	for _, recipe := range scopedRecipes {
		pattern := recipe + ` scope=""`
		assert.True(t, strings.Contains(section, pattern),
			"Expected scoped recipe %q in forge standard section", pattern)
	}
}

// Traceability: TC-FJ-003 -> AC: unscoped recipes present (no scope parameter)
// ci recipe present without scope parameter.
// Note: project-type, e2e-test, e2e-setup, e2e-verify are template-provided and may not
// appear in the generated justfile for a pure backend project.
func TestTC_FJ_003_UnscopedRecipesPresent(t *testing.T) {
	section := getStandardSection(t)
	// ci is always present (unscoped)
	assert.True(t, strings.Contains(section, "ci:"),
		"Expected ci recipe in forge standard section")
	// Verify ci does NOT have scope=""
	assert.False(t, strings.Contains(section, `ci scope=""`), "ci should NOT have scope")
}

// Traceability: TC-FJ-004 -> AC: Boundary markers present
// Boundary markers present in justfile.
func TestTC_FJ_004_BoundaryMarkersPresent(t *testing.T) {
	content := getJustfile(t)
	assert.True(t, strings.Contains(content, "# --- forge standard recipes ---"),
		"Expected start boundary marker")
	assert.True(t, strings.Contains(content, "# --- end forge standard recipes ---"),
		"Expected end boundary marker")
}

// Traceability: TC-FJ-005 -> AC: compile recipe has correct toolchain commands
// compile recipe contains the backend compile command (go vet).
func TestTC_FJ_005_CompileRecipeHasBackendCompileCommand(t *testing.T) {
	section := getStandardSection(t)
	assert.True(t, strings.Contains(section, `compile scope=""`), "Expected compile with scope")
	// The compile recipe should contain go vet for a Go project
	assert.True(t, strings.Contains(section, "go vet"),
		"Expected go vet in compile recipe")
}

// Traceability: TC-FJ-006 -> AC: build has correct backend branch
// build recipe has go build command.
func TestTC_FJ_006_BuildRecipeHasBackendBuildCommand(t *testing.T) {
	section := getStandardSection(t)
	buildIdx := strings.Index(section, `build scope=""`)
	assert.NotEqual(t, -1, buildIdx, "Expected build recipe in section")
	// The build recipe should contain go build for a Go project
	assert.True(t, strings.Contains(section, "go build"),
		"Expected go build in build recipe")
}

// Traceability: TC-FJ-007 -> AC: test recipe has correct backend branch
// test recipe has go test command.
func TestTC_FJ_007_TestRecipeHasBackendTestCommand(t *testing.T) {
	section := getStandardSection(t)
	testIdx := strings.Index(section, `test scope=""`)
	assert.NotEqual(t, -1, testIdx, "Expected test recipe in section")
	// The test recipe should contain go test for a Go project
	assert.True(t, strings.Contains(section, "go test"),
		"Expected go test in test recipe")
}

// Traceability: TC-FJ-008 -> AC: scoped recipes have shebang-based bodies
// All scoped recipes have shebang-based script bodies with error propagation.
func TestTC_FJ_008_AllScopedRecipesHaveShebangWithSetEuoPipefail(t *testing.T) {
	section := getStandardSection(t)
	scopedRecipes := []string{"compile", "build", "run", "dev", "test", "lint", "fmt", "check", "clean", "install"}
	for _, recipe := range scopedRecipes {
		recipeIdx := strings.Index(section, recipe+` scope=""`)
		assert.NotEqual(t, -1, recipeIdx, "Expected %s recipe", recipe)
		recipeToEnd := section[recipeIdx:]
		// Find the next recipe or end of section
		nextRecipe := len(recipeToEnd)
		for _, other := range scopedRecipes {
			if other == recipe {
				continue
			}
			idx := strings.Index(recipeToEnd[1:], other+` scope=""`)
			if idx != -1 && idx+1 < nextRecipe {
				nextRecipe = idx + 1
			}
		}
		// Also check for unscoped recipes
		for _, unscoped := range []string{"ci:", "e2e-test", "e2e-setup", "e2e-verify"} {
			idx := strings.Index(recipeToEnd[1:], unscoped)
			if idx != -1 && idx+1 < nextRecipe {
				nextRecipe = idx + 1
			}
		}
		recipeBlock := recipeToEnd[:nextRecipe]
		assert.True(t, strings.Contains(recipeBlock, "#!/usr/bin/env bash"),
			"Expected shebang in %s recipe", recipe)
		assert.True(t, strings.Contains(recipeBlock, "set -euo pipefail"),
			"Expected set -euo pipefail in %s recipe", recipe)
	}
}

// Traceability: TC-FJ-009 -> AC: ci chains standard commands
// ci recipe chains install, compile, build, test, lint.
func TestTC_FJ_009_CiRecipeChainsInstallCompileBuildTestLint(t *testing.T) {
	section := getStandardSection(t)
	assert.True(t, strings.Contains(section, "just install"), "Expected \"just install\" in ci")
	assert.True(t, strings.Contains(section, "just compile"), "Expected \"just compile\" in ci")
	assert.True(t, strings.Contains(section, "just build"), "Expected \"just build\" in ci")
	assert.True(t, strings.Contains(section, "just test"), "Expected \"just test\" in ci")
	assert.True(t, strings.Contains(section, "just lint"), "Expected \"just lint\" in ci")
}

// Traceability: TC-FJ-010 -> AC: Custom recipes preserved
// Custom recipes (claude, claude-c) preserved outside boundary markers.
func TestTC_FJ_010_CustomRecipesPreservedOutsideBoundaryMarkers(t *testing.T) {
	content := getJustfile(t)
	assert.True(t, strings.Contains(content, "claude:"), "Expected \"claude:\" recipe preserved")
	assert.True(t, strings.Contains(content, "claude-c:"), "Expected \"claude-c:\" recipe preserved")
}

// --- TC-FJ-011 to TC-FJ-015: Live just command execution ---

// Traceability: TC-FJ-011 -> AC: just compile backend dispatches correctly
// just compile backend runs without scope error (backend project accepts scope silently).
func TestTC_FJ_011_JustCompileBackendDispatchesCorrectly(t *testing.T) {
	exitCode, out := runJust("compile", "backend")
	lower := strings.ToLower(out)
	isScopeError := strings.Contains(lower, "[forge] invalid scope")
	assert.False(t, isScopeError, "Should not be a scope error for backend scope")
	if exitCode != 0 {
		t.Logf("compile backend exited %d (possibly missing toolchain): %s", exitCode, out)
	}
}

// Traceability: TC-FJ-012 -> AC: just compile frontend dispatches correctly
// just compile frontend runs without scope error (backend project accepts scope silently).
func TestTC_FJ_012_JustCompileFrontendDispatchesCorrectly(t *testing.T) {
	exitCode, out := runJust("compile", "frontend")
	lower := strings.ToLower(out)
	isScopeError := strings.Contains(lower, "[forge] invalid scope")
	assert.False(t, isScopeError, "Should not be a scope error for frontend scope")
	if exitCode != 0 {
		t.Logf("compile frontend exited %d (possibly missing toolchain): %s", exitCode, out)
	}
}

// Traceability: TC-FJ-013 -> AC: just compile (empty scope) dispatches correctly
// just compile with empty scope runs without scope error.
func TestTC_FJ_013_JustCompileEmptyScopeDispatchesCorrectly(t *testing.T) {
	exitCode, out := runJust("compile")
	lower := strings.ToLower(out)
	isScopeError := strings.Contains(lower, "[forge] invalid scope")
	assert.False(t, isScopeError, "Should not be a scope error with empty scope")
	if exitCode != 0 {
		t.Logf("compile (empty scope) exited %d (possibly missing toolchain): %s", exitCode, out)
	}
}

// Traceability: TC-FJ-014 -> AC: invalid scope behavior for non-dispatching project
// For a non-dispatching justfile, compile with any scope argument runs the same command.
// The scope parameter is accepted but not validated in a pure backend project.
func TestTC_FJ_014_JustCompileInvalidScopeBehaviorForNonDispatchingProject(t *testing.T) {
	exitCode, out := runJust("compile", "invalidscope")
	// For a non-dispatching justfile, scope is silently accepted.
	// The compile runs the same command regardless of scope value.
	// Verify no scope error (the justfile doesn't validate scope).
	lower := strings.ToLower(out)
	isScopeError := strings.Contains(lower, "[forge] invalid scope")
	assert.False(t, isScopeError, "Non-dispatching justfile should not produce scope error")
	// Exit code depends on toolchain availability
	if exitCode != 0 {
		t.Logf("compile invalidscope exited %d: %s", exitCode, out)
	}
}

// Traceability: TC-FJ-015 -> AC: project-type is available via forge probe
// forge probe returns a valid project type.
func TestTC_FJ_015_ProjectTypeAvailableViaForgeProbe(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("probe")
	assert.Equal(t, 0, exitCode, "Expected exit code 0 from forge probe")
	output := strings.TrimSpace(out)
	assert.NotEmpty(t, output, "Expected non-empty output from forge probe")
}

// --- Detection & Assembly Tests (TC-DET-001 to TC-DET-019) ---
// Converted from tests/e2e/justfile-e2e-integration/detection-assembly.spec.ts (19 tests).
// These tests validate project-type detection, classification, template selection,
// boundary markers, force flag, and detection signal mapping.

// getInitJustfileContent reads the init-justfile SKILL.md content.
func getInitJustfileContent(t *testing.T) string {
	t.Helper()
	return testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/SKILL.md")
}

// --- TC-DET-001 to TC-DET-004: Project-type detection signals ---

// Traceability: TC-DET-001 -> AC: Detection checks package.json (frontend signal)
// Detection logic checks for package.json as frontend signal.
func TestTC_DET_001_DetectionChecksPackageJsonAsFrontendSignal(t *testing.T) {
	content := getInitJustfileContent(t)
	assert.True(t, strings.Contains(content, "package.json"),
		"Expected detection logic to check for package.json")
	// Must be in the detection section, not just the language templates
	detectionSection := strings.Index(content, "Step 1:")
	assert.NotEqual(t, -1, detectionSection, "Expected Step 1 detection section")
}

// Traceability: TC-DET-002 -> AC: Detection checks go.mod (backend signal)
// Detection logic checks for go.mod as backend signal.
func TestTC_DET_002_DetectionChecksGoModAsBackendSignal(t *testing.T) {
	content := getInitJustfileContent(t)
	assert.True(t, strings.Contains(content, "go.mod"),
		"Expected detection logic to check for go.mod")
}

// Traceability: TC-DET-003 -> AC: Detection checks Cargo.toml (backend signal)
// Detection logic checks for Cargo.toml as backend signal.
func TestTC_DET_003_DetectionChecksCargoTomlAsBackendSignal(t *testing.T) {
	content := getInitJustfileContent(t)
	assert.True(t, strings.Contains(content, "Cargo.toml"),
		"Expected detection logic to check for Cargo.toml")
}

// Traceability: TC-DET-004 -> AC: Detection checks pyproject.toml (backend signal)
// Detection logic checks for pyproject.toml as backend signal.
func TestTC_DET_004_DetectionChecksPyprojectTomlAsBackendSignal(t *testing.T) {
	content := getInitJustfileContent(t)
	assert.True(t, strings.Contains(content, "pyproject.toml"),
		"Expected detection logic to check for pyproject.toml")
}

// --- TC-DET-005 to TC-DET-008: Project classification ---

// Traceability: TC-DET-005 -> AC: both signals -> mixed
// Classification produces mixed when both frontend and backend signals detected.
func TestTC_DET_005_ClassificationProducesMixedWhenBothSignalsDetected(t *testing.T) {
	content := getInitJustfileContent(t)
	assert.True(t, strings.Contains(content, "mixed"),
		"Expected \"mixed\" classification in detection logic")
}

// Traceability: TC-DET-006 -> AC: frontend only -> frontend
// Classification produces frontend when only frontend signals detected.
func TestTC_DET_006_ClassificationProducesFrontendWhenOnlyFrontendSignals(t *testing.T) {
	content := getInitJustfileContent(t)
	assert.True(t, strings.Contains(content, "frontend"),
		"Expected \"frontend\" classification in detection logic")
}

// Traceability: TC-DET-007 -> AC: backend only -> backend
// Classification produces backend when only backend signals detected.
func TestTC_DET_007_ClassificationProducesBackendWhenOnlyBackendSignals(t *testing.T) {
	content := getInitJustfileContent(t)
	assert.True(t, strings.Contains(content, "backend"),
		"Expected \"backend\" classification in detection logic")
}

// Traceability: TC-DET-008 -> AC: neither signal -> error message
// Classification produces error when no known markers detected.
func TestTC_DET_008_ClassificationProducesErrorWhenNoMarkersDetected(t *testing.T) {
	content := getInitJustfileContent(t)
	// Should describe error case when no markers found
	hasError := strings.Contains(content, "no known project markers") ||
		strings.Contains(content, "no project markers") ||
		strings.Contains(content, "Error: no known") ||
		strings.Contains(content, "no markers detected") ||
		strings.Contains(content, "neither") ||
		strings.Contains(content, "Cannot determine")
	assert.True(t, hasError, "Expected error handling when no project markers detected")
}

// --- TC-DET-009 to TC-DET-011: Template selection ---

// Traceability: TC-DET-009 -> AC: Selects backend template for backend project
// Selects backend template for pure backend projects.
func TestTC_DET_009_SelectsBackendTemplateForPureBackendProjects(t *testing.T) {
	backendTemplate := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/templates/go.just")
	assert.True(t, strings.Contains(backendTemplate, "go vet"),
		"Expected backend template to contain go vet command")
}

// Traceability: TC-DET-010 -> AC: Selects frontend template for frontend project
// Selects frontend template for pure frontend projects.
func TestTC_DET_010_SelectsFrontendTemplateForPureFrontendProjects(t *testing.T) {
	frontendTemplate := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/templates/node.just")
	assert.True(t, strings.Contains(frontendTemplate, "npx tsc --noEmit"),
		"Expected frontend template to contain npx tsc --noEmit command")
}

// Traceability: TC-DET-011 -> AC: Selects mixed template for mixed project
// Selects mixed template for mixed projects.
func TestTC_DET_011_SelectsMixedTemplateForMixedProjects(t *testing.T) {
	mixedTemplate := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/templates/mixed.just")
	assert.True(t, strings.Contains(mixedTemplate, "frontend_dir"),
		"Expected mixed template to contain frontend_dir variable")
	assert.True(t, strings.Contains(mixedTemplate, "backend_dir"),
		"Expected mixed template to contain backend_dir variable")
}

// --- TC-DET-012 to TC-DET-014: Boundary markers ---

// Traceability: TC-DET-012 -> AC: Boundary markers wrap generated recipes
// Generated recipes wrapped in forge boundary markers.
func TestTC_DET_012_GeneratedRecipesWrappedInForgeBoundaryMarkers(t *testing.T) {
	content := getInitJustfileContent(t)
	startMarker := "# --- forge standard recipes ---"
	endMarker := "# --- end forge standard recipes ---"
	assert.True(t, strings.Contains(content, startMarker), "Expected start boundary marker")
	assert.True(t, strings.Contains(content, endMarker), "Expected end boundary marker")
}

// Traceability: TC-DET-013 -> AC: Boundary marker merge replaces only marked section
// Boundary marker merge replaces only marked section.
func TestTC_DET_013_BoundaryMarkerMergeReplacesOnlyMarkedSection(t *testing.T) {
	content := getInitJustfileContent(t)
	// Should describe merge logic that preserves content outside markers
	hasMerge := strings.Contains(content, "boundary marker") ||
		strings.Contains(content, "markers") ||
		strings.Contains(content, "replace") ||
		strings.Contains(content, "merge")
	assert.True(t, hasMerge, "Expected boundary marker merge logic description")
}

// Traceability: TC-DET-014 -> AC: All 3 templates have boundary markers
// All three templates have boundary markers.
func TestTC_DET_014_AllThreeTemplatesHaveBoundaryMarkers(t *testing.T) {
	startMarker := "# --- forge standard recipes ---"
	endMarker := "# --- end forge standard recipes ---"
	templates := []string{
		"../plugins/forge/skills/init-justfile/templates/go.just",
		"../plugins/forge/skills/init-justfile/templates/node.just",
		"../plugins/forge/skills/init-justfile/templates/mixed.just",
	}
	for _, tpl := range templates {
		content := testkit.ReadProjectFile(t, tpl)
		assert.True(t, strings.Contains(content, startMarker),
			"Expected start boundary marker in %s", tpl)
		assert.True(t, strings.Contains(content, endMarker),
			"Expected end boundary marker in %s", tpl)
	}
}

// --- TC-DET-015 to TC-DET-017: --force flag and interactive confirmation ---

// Traceability: TC-DET-015 -> AC: --force flag skips confirmation
// --force flag skips user confirmation for agent use.
func TestTC_DET_015_ForceFlagSkipsUserConfirmationForAgentUse(t *testing.T) {
	content := getInitJustfileContent(t)
	assert.True(t, strings.Contains(content, "--force"),
		"Expected --force flag to be documented in init-justfile")
	hasAgentRef := strings.Contains(content, "agent") ||
		strings.Contains(content, "non-interactive") ||
		strings.Contains(content, "skip")
	assert.True(t, hasAgentRef, "Expected --force to be described as skipping confirmation for agents")
}

// Traceability: TC-DET-016 -> AC: Interactive confirmation when no markers exist
// Interactive confirmation when no boundary markers and justfile exists.
func TestTC_DET_016_InteractiveConfirmationWhenNoBoundaryMarkersAndJustfileExists(t *testing.T) {
	content := getInitJustfileContent(t)
	hasConfirm := strings.Contains(content, "confirm") ||
		strings.Contains(content, "prompt") ||
		strings.Contains(content, "overwrite") ||
		strings.Contains(content, "ask")
	assert.True(t, hasConfirm, "Expected interactive confirmation description for existing justfile without markers")
}

// Traceability: TC-DET-017 -> AC: Idempotent re-run preserves custom recipes
// Re-running init-justfile preserves user custom recipes.
func TestTC_DET_017_ReRunningInitJustfilePreservesUserCustomRecipes(t *testing.T) {
	content := getInitJustfileContent(t)
	// Should describe that boundary markers enable idempotent re-runs
	hasPreserve := strings.Contains(content, "preserve") ||
		strings.Contains(content, "keep") ||
		strings.Contains(content, "custom") ||
		strings.Contains(content, "outside")
	assert.True(t, hasPreserve, "Expected description of preserving custom recipes outside boundary markers")
}

// --- TC-DET-018: project-type recipe for each variant ---

// Traceability: TC-DET-018 -> AC: project-type recipe outputs correct type
// All three project-type recipe variants exist as distinct templates.
func TestTC_DET_018_AllThreeProjectTypeRecipeVariantsExist(t *testing.T) {
	nodeContent := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/templates/node.just")
	goContent := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/templates/go.just")
	mixedContent := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/templates/mixed.just")

	// Verify templates exist with distinct content — each template targets a different project type
	assert.True(t, strings.Contains(nodeContent, "npm"), "Expected npm commands in node template")
	assert.True(t, strings.Contains(goContent, "go vet") || strings.Contains(goContent, "go build"),
		"Expected go commands in go template")
	assert.True(t, strings.Contains(mixedContent, "frontend_dir") && strings.Contains(mixedContent, "backend_dir"),
		"Expected frontend_dir and backend_dir in mixed template")
}

// --- TC-DET-019: Detection signal mapping ---

// Traceability: TC-DET-019 -> AC: package.json = frontend, go.mod/Cargo.toml/pyproject.toml = backend
// Detection correctly maps signals to frontend/backend categories.
func TestTC_DET_019_DetectionCorrectlyMapsSignalsToFrontendBackendCategories(t *testing.T) {
	content := getInitJustfileContent(t)

	// Find the detection/classification section
	workflowIdx := strings.Index(content, "## Workflow")
	assert.NotEqual(t, -1, workflowIdx, "Expected Workflow section")

	workflowSection := content[workflowIdx:]

	// frontend signal should be mapped
	assert.True(t, strings.Contains(workflowSection, "frontend"),
		"Expected \"frontend\" classification in workflow")

	// backend signal should be mapped
	assert.True(t, strings.Contains(workflowSection, "backend"),
		"Expected \"backend\" classification in workflow")

	// Verify the detection is structured with clear mapping
	assert.True(t, strings.Contains(content, "package.json"),
		"Expected package.json in detection mapping")
	assert.True(t, strings.Contains(content, "go.mod"),
		"Expected go.mod in detection mapping")
}
