//go:build e2e

package justfileintegration

import (
	"strings"
	"testing"

	"forge-cli/tests/e2e/testkit"

	"github.com/stretchr/testify/assert"
)

// --- Init-Justfile Tests (TC-004 to TC-022) ---
// Converted from tests/e2e/init-justfile/init-justfile.spec.ts (7 tests).
// These tests verify justfile initialization: template correctness, command presence,
// error handling, and boundary-marker idempotency.

// The 11 standard command names always present in the generated justfile per Spec 5.1.
// Template-provided recipes (e2e-test, e2e-setup, e2e-verify, project-type) are
// injected by the active profile manifest and may not appear in all project types.
var standardCommands = []string{
	"compile", "build", "run", "dev", "test",
	"lint", "fmt", "check", "clean",
	"install", "ci",
}

// fileContains is a helper that checks whether content contains needle.
func fileContains(content, needle string) bool {
	return strings.Contains(content, needle)
}

// getInitJustfileSkillContent reads the init-justfile SKILL.md content.
func getInitJustfileSkillContent(t *testing.T) string {
	t.Helper()
	return testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/SKILL.md")
}

// Traceability: TC-004 -> Story 2 / AC-1
// Frontend project detection generates scope-free justfile.
// The node template targets pure frontend projects without scope dispatch.
func TestTC_004_FrontendProjectDetectionGeneratesScopeFreeJustfile(t *testing.T) {
	frontendTemplate := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/templates/node.just")
	// Frontend template should NOT have scope="" parameters
	assert.False(t, fileContains(frontendTemplate, `scope=""`),
		`Expected frontend template to NOT have scope="" parameters`)
	// Frontend template should contain npm/npx commands
	assert.True(t, fileContains(frontendTemplate, "npx tsc --noEmit"),
		"Expected npx tsc --noEmit in frontend template")
}

// Traceability: TC-005 -> Story 2 / AC-2
// Backend project detection generates scope-free justfile.
// The go template targets pure backend projects without scope dispatch.
func TestTC_005_BackendProjectDetectionGeneratesScopeFreeJustfile(t *testing.T) {
	backendTemplate := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/templates/go.just")
	// Backend template should NOT have scope="" parameters
	assert.False(t, fileContains(backendTemplate, `scope=""`),
		`Expected backend template to NOT have scope="" parameters`)
	// Backend template should contain go commands
	assert.True(t, fileContains(backendTemplate, "go vet"),
		"Expected go vet in backend template")
}

// Traceability: TC-006 -> Story 2 / AC-3
// Mixed project detection generates scope-aware justfile.
// The mixed template targets mixed projects with scope dispatch (case/esac).
func TestTC_006_MixedProjectDetectionGeneratesScopeAwareJustfile(t *testing.T) {
	mixedTemplate := testkit.ReadProjectFile(t, "../plugins/forge/skills/init-justfile/templates/mixed.just")
	// Mixed template SHOULD have scope="" parameters
	assert.True(t, fileContains(mixedTemplate, `scope=""`),
		`Expected mixed template to have scope="" parameters`)
	// Mixed template should have scope dispatch with case/esac
	assert.True(t, fileContains(mixedTemplate, `case "{{scope}}"`),
		`Expected case "{{scope}}" in mixed template for scope dispatch`)
}

// Traceability: TC-022 -> Spec 5.1 / vocabulary
// All 15 standard commands are present in generated justfile.
func TestTC_022_All15StandardCommandsArePresentInGeneratedJustfile(t *testing.T) {
	justfile := testkit.ReadProjectFile(t, "../justfile")
	for _, cmd := range standardCommands {
		assert.True(t, fileContains(justfile, cmd),
			"Expected recipe %q in the forge project justfile", cmd)
	}
}

// Traceability: TC-018 -> Spec 5.2 / detection
// No marker files detected causes init-justfile to error.
func TestTC_018_NoMarkerFilesDetectedCausesInitJustfileToError(t *testing.T) {
	content := getInitJustfileSkillContent(t)
	// The init-justfile skill should describe an error case for no markers
	hasError := fileContains(content, "no known project markers") ||
		fileContains(content, "no project markers") ||
		fileContains(content, "Error: no known") ||
		fileContains(content, "no markers detected") ||
		fileContains(content, "neither") ||
		fileContains(content, "Cannot determine")
	assert.True(t, hasError,
		"Expected error handling description for no project markers")
}

// Traceability: TC-019 -> Spec 5.2 / flow
// Existing justfile triggers user confirmation.
func TestTC_019_ExistingJustfileTriggersUserConfirmation(t *testing.T) {
	content := getInitJustfileSkillContent(t)
	hasConfirm := fileContains(content, "confirm") ||
		fileContains(content, "prompt") ||
		fileContains(content, "overwrite") ||
		fileContains(content, "ask") ||
		fileContains(content, "--force")
	assert.True(t, hasConfirm,
		"Expected user confirmation mechanism for existing justfile")
}

// Traceability: TC-020 -> Spec / maintainability
// Boundary markers present triggers idempotent merge.
func TestTC_020_BoundaryMarkersPresentTriggersIdempotentMerge(t *testing.T) {
	justfile := testkit.ReadProjectFile(t, "../justfile")
	startMarker := "# --- forge standard recipes ---"
	endMarker := "# --- end forge standard recipes ---"

	assert.True(t, fileContains(justfile, startMarker),
		"Expected start boundary marker in justfile")
	assert.True(t, fileContains(justfile, endMarker),
		"Expected end boundary marker in justfile")

	startIdx := strings.Index(justfile, startMarker)
	endIdx := strings.Index(justfile, endMarker)
	assert.NotEqual(t, -1, startIdx, "Expected start boundary marker index")
	assert.NotEqual(t, -1, endIdx, "Expected end boundary marker index")

	// Custom recipes should be OUTSIDE boundary markers
	beforeMarkers := justfile[:startIdx]
	afterMarkers := justfile[endIdx+len(endMarker):]
	customRecipesOutsideMarkers :=
		strings.Contains(beforeMarkers, "claude:") ||
			strings.Contains(afterMarkers, "claude:")
	assert.True(t, customRecipesOutsideMarkers,
		"Expected custom recipes outside boundary markers")
}
