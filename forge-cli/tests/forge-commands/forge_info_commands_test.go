//go:build e2e

package forgecommands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/tests/testkit"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Config Commands (Task 1) — TC-001 to TC-007
// =============================================================================

// Traceability: TC-001 -> Task 1 / AC-1, AC-2
func TestTC_001_ConfigInitInteractiveSetup(t *testing.T) {
	t.Skip("requires interactive stdin: multi-step prompt with project-type, profiles, capabilities selections")
}

// Traceability: TC-002 -> Task 1 / AC-3
func TestTC_002_ConfigInitReconfigurePrompt(t *testing.T) {
	t.Skip("requires interactive stdin and pre-existing .forge/config.yaml to trigger reconfigure prompt")
}

// Traceability: TC-003 -> Task 1 / AC-3
func TestTC_003_ConfigInitReconfigureAccepted(t *testing.T) {
	t.Skip("requires interactive stdin and pre-existing .forge/config.yaml with 'y' confirmation input")
}

// Traceability: TC-004 -> Task 1 / AC-4
func TestTC_004_ConfigGetProjectType(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("config", "get", "project-type")

	assert.Equal(t, 0, exitCode, "config get project-type should exit 0")
	assert.Equal(t, "backend", strings.TrimSpace(out),
		"config get project-type should output plain text 'backend' without formatting")
}

// Traceability: TC-005 -> Task 1 / AC-5
func TestTC_005_ConfigGetInterfacesArrayOutput(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("config", "get", "interfaces")

	assert.Equal(t, 0, exitCode, "config get interfaces should exit 0")
	lines := strings.Split(strings.TrimSpace(out), "\n")
	assert.True(t, len(lines) >= 3,
		"interfaces should output at least 3 lines (one per item), got %d: %q", len(lines), out)
	// Verify no formatting blocks or quotes — each line should be a plain interface name
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		assert.NotContains(t, trimmed, `"`, "each line should not contain quotes: %q", trimmed)
		assert.NotContains(t, trimmed, "`", "each line should not contain backticks: %q", trimmed)
	}
}

// Traceability: TC-006 -> Task 1 / AC-6
func TestTC_006_ConfigGetMissingKey(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("config", "get", "nonexistent-key")

	assert.Equal(t, 1, exitCode, "config get with missing key should exit 1")
	assert.Equal(t, "", strings.TrimSpace(out),
		"config get with missing key should produce no stdout output")
}

// Traceability: TC-007 -> Task 1 / AC-7
func TestTC_007_ForgeConfigStructFields(t *testing.T) {
	t.Skip("struct field validation requires Go source inspection, not CLI invocation; covered by unit tests")
}

// =============================================================================
// Proposal Commands (Task 2) — TC-008 to TC-012
// =============================================================================

// Traceability: TC-008 -> Task 2 / AC-1
func TestTC_008_ProposalListAllProposals(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("proposal")

	assert.Equal(t, 0, exitCode, "proposal list should exit 0")
	assert.True(t, len(strings.TrimSpace(out)) > 0,
		"proposal list output should not be empty")

	// Verify table columns present
	upper := strings.ToUpper(out)
	assert.True(t, strings.Contains(upper, "SLUG"),
		"proposal list output should contain SLUG column header: %s", out)
	assert.True(t, strings.Contains(upper, "CREATED"),
		"proposal list output should contain CREATED column header: %s", out)
	assert.True(t, strings.Contains(upper, "STATUS"),
		"proposal list output should contain STATUS column header: %s", out)
	assert.True(t, strings.Contains(upper, "PRD"),
		"proposal list output should contain PRD column header: %s", out)
	assert.True(t, strings.Contains(upper, "FEATURE"),
		"proposal list output should contain FEATURE column header: %s", out)
}

// Traceability: TC-009 -> Task 2 / AC-2
func TestTC_009_ProposalSlugDetailView(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("proposal", "forge-info-commands")

	assert.Equal(t, 0, exitCode, "proposal detail should exit 0")
	assert.True(t, len(strings.TrimSpace(out)) > 0,
		"proposal detail output should not be empty")

	// Verify detail fields present
	upper := strings.ToUpper(out)
	assert.True(t, strings.Contains(upper, "SLUG"),
		"proposal detail should show SLUG field: %s", out)
	assert.True(t, strings.Contains(upper, "CREATED"),
		"proposal detail should show CREATED field: %s", out)
	assert.True(t, strings.Contains(upper, "STATUS"),
		"proposal detail should show STATUS field: %s", out)
	assert.True(t, strings.Contains(upper, "FILE"),
		"proposal detail should show FILE path: %s", out)
}

// Traceability: TC-010 -> Task 2 / AC-3
func TestTC_010_ProposalCreatedDateFromFrontmatter(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("proposal")

	assert.Equal(t, 0, exitCode, "proposal list should exit 0")
	// Verify that the Created column shows a date from frontmatter, not file system time
	assert.True(t, strings.Contains(out, "2026-05-14"),
		"proposal list should show created date '2026-05-14' from frontmatter: %s", out)
}

// Traceability: TC-011 -> Task 2 / AC-4
func TestTC_011_ProposalPRDColumnChecksPrdSpec(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("proposal")

	assert.Equal(t, 0, exitCode, "proposal list should exit 0")
	// For forge-info-commands which has no prd-spec.md, PRD column should show absence indicator
	lines := strings.Split(out, "\n")
	found := false
	for _, line := range lines {
		if strings.Contains(line, "forge-info-commands") {
			// PRD column should show "-" or "no" for proposals without prd-spec.md
			assert.True(t,
				strings.Contains(line, " - ") || strings.Contains(line, "No"),
				"PRD column for forge-info-commands should show absence indicator: %s", line)
			found = true
			break
		}
	}
	assert.True(t, found, "should find forge-info-commands row in proposal list: %s", out)
}

// Traceability: TC-012 -> Task 2 / AC-5
func TestTC_012_ProposalFeatureColumnReadsManifestStatus(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("proposal")

	assert.Equal(t, 0, exitCode, "proposal list should exit 0")
	// Feature column should show manifest status (e.g. "tasks") for proposals with features
	lines := strings.Split(out, "\n")
	found := false
	for _, line := range lines {
		if strings.Contains(line, "forge-info-commands") {
			assert.True(t,
				strings.Contains(line, "tasks") || strings.Contains(line, "pending") || strings.Contains(line, "completed"),
				"Feature column for forge-info-commands should show manifest status: %s", line)
			found = true
			break
		}
	}
	assert.True(t, found, "should find forge-info-commands row in proposal list: %s", out)
}

// =============================================================================
// Feature Commands (Task 2) — TC-013 to TC-017
// =============================================================================

// Traceability: TC-013 -> Task 2 / AC-6
func TestTC_013_FeatureListAllFeatures(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("feature", "list")

	assert.Equal(t, 0, exitCode, "feature list should exit 0")
	assert.True(t, len(strings.TrimSpace(out)) > 0,
		"feature list output should not be empty")

	// Verify table columns present
	upper := strings.ToUpper(out)
	assert.True(t, strings.Contains(upper, "SLUG"),
		"feature list should contain SLUG column header: %s", out)
	assert.True(t, strings.Contains(upper, "STATUS"),
		"feature list should contain STATUS column header: %s", out)
	assert.True(t, strings.Contains(upper, "PROGRESS"),
		"feature list should contain PROGRESS column header: %s", out)
}

// Traceability: TC-014 -> Task 2 / AC-7
func TestTC_014_FeatureListProgressFromIndexJSON(t *testing.T) {
	t.Skip("requires manual setup: feature with known task counts in index.json for precise progress assertion")
}

// Traceability: TC-015 -> Task 2 / AC-8
func TestTC_015_FeatureListScoresFromFrontmatter(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("feature", "list")

	assert.Equal(t, 0, exitCode, "feature list should exit 0")
	// Score columns should show em-dash when no score field exists in frontmatter
	// The em-dash character is used as placeholder
	assert.True(t, strings.Contains(out, "—"),
		"feature list should show em-dash for missing scores: %s", out)
}

// Traceability: TC-016 -> Task 2 / AC-9
func TestTC_016_FeatureStatusDetailView(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("feature", "status", "forge-info-commands")

	assert.Equal(t, 0, exitCode, "feature status should exit 0")
	assert.True(t, len(strings.TrimSpace(out)) > 0,
		"feature status output should not be empty")

	upper := strings.ToUpper(out)
	assert.True(t, strings.Contains(upper, "SLUG"),
		"feature status should show SLUG field: %s", out)
	assert.True(t, strings.Contains(upper, "STATUS"),
		"feature status should show STATUS field: %s", out)
	assert.True(t, strings.Contains(upper, "TASKS"),
		"feature status should show TASKS section: %s", out)
}

// Traceability: TC-017 -> Task 2 / Hard Rules
func TestTC_017_FeatureNoArgsKeepsExistingBehavior(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("feature")

	assert.Equal(t, 0, exitCode, "feature with no args should exit 0")
	assert.True(t, len(strings.TrimSpace(out)) > 0,
		"feature with no args should display current feature")
}

// =============================================================================
// Lesson Commands (Task 2) — TC-018 to TC-020
// =============================================================================

// Traceability: TC-018 -> Task 2 / AC-10
func TestTC_018_LessonListAllLessons(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("lesson")

	assert.Equal(t, 0, exitCode, "lesson list should exit 0")
	assert.True(t, len(strings.TrimSpace(out)) > 0,
		"lesson list output should not be empty")

	// Verify table columns present
	upper := strings.ToUpper(out)
	assert.True(t, strings.Contains(upper, "NAME"),
		"lesson list should contain NAME column header: %s", out)
	assert.True(t, strings.Contains(upper, "CREATED"),
		"lesson list should contain CREATED column header: %s", out)
	assert.True(t, strings.Contains(upper, "CATEGORY"),
		"lesson list should contain CATEGORY column header: %s", out)
	assert.True(t, strings.Contains(upper, "TAGS"),
		"lesson list should contain TAGS column header: %s", out)
}

// Traceability: TC-019 -> Task 2 / AC-11
func TestTC_019_LessonCategoryFromFilenamePrefix(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("lesson")

	assert.Equal(t, 0, exitCode, "lesson list should exit 0")
	// Category should be derived from filename prefix (e.g. "gotcha-" -> "gotcha")
	lines := strings.Split(out, "\n")
	foundGotcha := false
	for _, line := range lines {
		if strings.Contains(line, "gotcha") {
			foundGotcha = true
			break
		}
	}
	// If gotcha lessons exist, they should show "gotcha" as category
	assert.True(t, foundGotcha,
		"lesson list should show 'gotcha' category for gotcha-* lesson files: %s", out)
}

// Traceability: TC-020 -> Task 2 / AC-12
func TestTC_020_LessonNameDetailView(t *testing.T) {
	// First list lessons to find a valid name
	_, listOut := testkit.RunCLIExitCode("lesson")
	lines := strings.Split(listOut, "\n")

	// Find a lesson name from the list output (skip header/separators)
	var lessonName string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 1 {
			candidate := strings.TrimSpace(fields[0])
			// Skip header lines and separator lines
			if candidate != "" && candidate != "NAME" && !strings.HasPrefix(candidate, "-") &&
				!strings.HasPrefix(candidate, "---") && !strings.Contains(candidate, "LESSONS") {
				lessonName = candidate
				break
			}
		}
	}

	if lessonName == "" {
		t.Skip("no lessons found in docs/lessons/ to test detail view")
	}

	exitCode, out := testkit.RunCLIExitCode("lesson", lessonName)

	assert.Equal(t, 0, exitCode, "lesson detail should exit 0")
	assert.True(t, len(strings.TrimSpace(out)) > 0,
		"lesson detail output should not be empty")

	upper := strings.ToUpper(out)
	assert.True(t, strings.Contains(upper, "NAME"),
		"lesson detail should show NAME field: %s", out)
	assert.True(t, strings.Contains(upper, "FILE"),
		"lesson detail should show FILE path: %s", out)
	// Full content should NOT be printed
	assert.False(t, strings.Contains(out, "## "),
		"lesson detail should NOT print full markdown content (heading found): %s", out)
}

// =============================================================================
// Init Command (Task 3) — TC-021 to TC-029
// =============================================================================

// Traceability: TC-021 -> Task 3 / AC-1
func TestTC_021_InitCreatesForgeDir(t *testing.T) {
	t.Skip("requires clean project state (no .forge/ directory); destructive to run against real project")
}

// Traceability: TC-022 -> Task 3 / AC-2
func TestTC_022_InitGeneratesCLAUDEmd(t *testing.T) {
	t.Skip("requires clean project state (no CLAUDE.md); destructive to run against real project")
}

// Traceability: TC-023 -> Task 3 / AC-3
func TestTC_023_InitAppendsGitignoreWithDedup(t *testing.T) {
	t.Skip("requires isolated project state; modifies .gitignore which is destructive")
}

// Traceability: TC-024 -> Task 3 / AC-3, Hard Rules
func TestTC_024_InitGitignoreDedupSkipsExisting(t *testing.T) {
	t.Skip("requires .gitignore with pre-existing forge entries; modifies .gitignore")
}

// Traceability: TC-025 -> Task 3 / AC-4
func TestTC_025_InitAppendsJustfileRecipes(t *testing.T) {
	t.Skip("requires isolated project state; modifies justfile which is destructive")
}

// Traceability: TC-026 -> Task 3 / AC-4, Hard Rules
func TestTC_026_InitJustfileDedupSkipsExisting(t *testing.T) {
	t.Skip("requires justfile with pre-existing claude recipe; modifies justfile")
}

// Traceability: TC-027 -> Task 3 / AC-5
func TestTC_027_InitRunsConfigInitWhenNoConfig(t *testing.T) {
	t.Skip("requires clean project state (no .forge/config.yaml) and interactive stdin")
}

// Traceability: TC-028 -> Task 3 / AC-6, Hard Rules
func TestTC_028_InitSkipsExistingFiles(t *testing.T) {
	t.Skip("requires pre-existing .forge/, CLAUDE.md, .forge/config.yaml; complex setup")
}

// Traceability: TC-029 -> Task 3 / AC-7
func TestTC_029_InitResultReportFormat(t *testing.T) {
	t.Skip("requires clean project state; destructive to run against real project")
}

// =============================================================================
// Migration (Task 4) — TC-030 to TC-032
// =============================================================================

// Traceability: TC-030 -> Task 4 / AC-1, Hard Rules
func TestTC_030_ResolveScopeReadsConfigDirectly(t *testing.T) {
	t.Skip("ResolveScope() is an internal function, not a CLI command; covered by unit tests")
}

// Traceability: TC-031 -> Task 4 / AC-6, Hard Rules
func TestTC_031_ResolveScopeMissingConfigReturnsEmpty(t *testing.T) {
	t.Skip("ResolveScope() is an internal function, not a CLI command; covered by unit tests")
}

// Traceability: TC-032 -> Task 4 / AC-3
func TestTC_032_JustfileHasNoProjectTypeRecipe(t *testing.T) {
	// Read justfile and verify no project-type: recipe exists
	justfilePath := filepath.Join("..", "..", "..", "justfile")
	data, err := os.ReadFile(justfilePath)
	if err != nil {
		// Try relative to working directory
		data, err = os.ReadFile("justfile")
		if err != nil {
			t.Skip("cannot locate justfile for migration check")
		}
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		assert.False(t,
			strings.HasPrefix(trimmed, "project-type:") || strings.HasPrefix(trimmed, "project-type :"),
			"justfile should NOT contain a 'project-type:' recipe (migration task): found %q", trimmed)
	}
}
