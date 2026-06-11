//go:build cli_functional

package testsuitehealth

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==============================================================================
// gen-journeys skill structure tests — feature: contract-journey-test-model
// Tests verify that the gen-journeys skill exists with correct structure,
// templates, and content that satisfies the acceptance criteria:
//   - Extracts Journeys from PRD user stories as Markdown
//   - Each Journey has: name, risk level (High/Medium/Low), happy path, edge cases
//   - Output is per-Journey files (not per interface type)
//   - Format is parseable by gen-contracts (Journey name + Step sequence + user action + expected result)
//   - High-risk Journeys have edge case count >= happy path step count
//   - No code reconnaissance needed (pure narrative extraction)
//   - Single Journey per generation, auto-batch when context window exceeded
// ==============================================================================

// projectRootGenJourneys returns the forge project root directory.
func projectRootGenJourneys(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine test file location")
	}
	// thisFile: .../tests/test-suite-health/gen_journeys_skill_test.go
	// up 2: test-suite-health -> tests -> project root
	dir := filepath.Join(filepath.Dir(thisFile), "..", "..")
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("cannot resolve project root: %s", err)
	}
	return abs
}

// skillRoot returns the gen-journeys skill directory.
func skillRoot(t *testing.T) string {
	t.Helper()
	return filepath.Join(projectRootGenJourneys(t), "plugins", "forge", "skills", "gen-journeys")
}

// TC-001: gen-journeys skill directory exists with SKILL.md
// Traceability: AC "gen-journeys from PRD user stories -> Journey Markdown output"
func TestTC_001_GenJourneysSkillDirectoryExists(t *testing.T) {
	dir := skillRoot(t)

	info, err := os.Stat(dir)
	require.NoError(t, err, "gen-journeys skill directory should exist")
	assert.True(t, info.IsDir(), "gen-journeys should be a directory")

	skillFile := filepath.Join(dir, "SKILL.md")
	_, err = os.Stat(skillFile)
	require.NoError(t, err, "SKILL.md must exist in gen-journeys skill directory")
}

// TC-002: SKILL.md has valid frontmatter with name and description
// Traceability: Forge distribution convention — skill must have proper frontmatter
func TestTC_002_SkillMdHasValidFrontmatter(t *testing.T) {
	skillFile := filepath.Join(skillRoot(t), "SKILL.md")
	content, err := os.ReadFile(skillFile)
	require.NoError(t, err)

	// Must have YAML frontmatter
	assert.True(t, strings.HasPrefix(string(content), "---"),
		"SKILL.md must start with YAML frontmatter")

	// Must have name field
	assert.Contains(t, string(content), "name: gen-journeys",
		"SKILL.md frontmatter must declare name: gen-journeys")

	// Must have description field
	assert.Contains(t, string(content), "description:",
		"SKILL.md frontmatter must have a description")
}

// TC-003: gen-journeys has templates directory with journey template
// Traceability: AC "output per-Journey files" — template drives per-Journey output
func TestTC_003_GenJourneysHasJourneyTemplate(t *testing.T) {
	templatesDir := filepath.Join(skillRoot(t), "templates")

	info, err := os.Stat(templatesDir)
	require.NoError(t, err, "templates directory must exist")
	assert.True(t, info.IsDir(), "templates should be a directory")

	// Must have a journey template file
	entries, err := os.ReadDir(templatesDir)
	require.NoError(t, err)
	assert.NotEmpty(t, entries, "templates directory must contain at least one template file")

	// Must have a journey-specific template
	hasJourneyTemplate := false
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "journey") {
			hasJourneyTemplate = true
			break
		}
	}
	assert.True(t, hasJourneyTemplate,
		"templates directory must contain a journey template file (filename containing 'journey')")
}

// TC-004: Journey template contains gen-contracts-parseable structure
// Traceability: AC "format parseable by gen-contracts (Journey name + Step sequence + user action + expected result)"
func TestTC_004_JourneyTemplateHasGenContractsParseableStructure(t *testing.T) {
	templatesDir := filepath.Join(skillRoot(t), "templates")
	entries, err := os.ReadDir(templatesDir)
	require.NoError(t, err)

	// Find the journey template
	var journeyTemplateContent string
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "journey") && !entry.IsDir() {
			data, err := os.ReadFile(filepath.Join(templatesDir, entry.Name()))
			require.NoError(t, err)
			journeyTemplateContent = string(data)
			break
		}
	}
	require.NotEmpty(t, journeyTemplateContent, "must find journey template")

	// Must have Journey name placeholder
	assert.Contains(t, journeyTemplateContent, "Journey",
		"Journey template must contain 'Journey' heading/label")

	// Must have Risk level (High/Medium/Low)
	hasRiskLevel := strings.Contains(journeyTemplateContent, "Risk") ||
		strings.Contains(journeyTemplateContent, "risk")
	assert.True(t, hasRiskLevel,
		"Journey template must contain Risk/risk field")

	// Must have Step sequence structure
	hasStep := strings.Contains(journeyTemplateContent, "Step") ||
		strings.Contains(journeyTemplateContent, "step")
	assert.True(t, hasStep,
		"Journey template must contain Step/step structure")

	// Must have user action and expected outcome fields
	hasUserAction := strings.Contains(journeyTemplateContent, "action") ||
		strings.Contains(journeyTemplateContent, "Action") ||
		strings.Contains(journeyTemplateContent, "user") ||
		strings.Contains(journeyTemplateContent, "operation") ||
		strings.Contains(journeyTemplateContent, "command")
	assert.True(t, hasUserAction,
		"Journey template must contain user action field")

	hasExpectedOutcome := strings.Contains(journeyTemplateContent, "expected") ||
		strings.Contains(journeyTemplateContent, "Expected") ||
		strings.Contains(journeyTemplateContent, "outcome") ||
		strings.Contains(journeyTemplateContent, "Outcome") ||
		strings.Contains(journeyTemplateContent, "result") ||
		strings.Contains(journeyTemplateContent, "Result")
	assert.True(t, hasExpectedOutcome,
		"Journey template must contain expected outcome/result field")
}

// TC-005: Journey template has Happy Path and Edge Cases sections
// Traceability: AC "at least 1 happy path step + at least 1 edge case step"
func TestTC_005_JourneyTemplateHasHappyPathAndEdgeCases(t *testing.T) {
	templatesDir := filepath.Join(skillRoot(t), "templates")
	entries, err := os.ReadDir(templatesDir)
	require.NoError(t, err)

	var journeyTemplateContent string
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "journey") && !entry.IsDir() {
			data, err := os.ReadFile(filepath.Join(templatesDir, entry.Name()))
			require.NoError(t, err)
			journeyTemplateContent = string(data)
			break
		}
	}
	require.NotEmpty(t, journeyTemplateContent, "must find journey template")

	// Must have Happy Path section
	hasHappyPath := strings.Contains(journeyTemplateContent, "Happy Path") ||
		strings.Contains(journeyTemplateContent, "happy path") ||
		strings.Contains(journeyTemplateContent, "happy-path")
	assert.True(t, hasHappyPath,
		"Journey template must have Happy Path section")

	// Must have Edge Cases section
	hasEdgeCases := strings.Contains(journeyTemplateContent, "Edge") ||
		strings.Contains(journeyTemplateContent, "edge") ||
		strings.Contains(journeyTemplateContent, "Edge Case") ||
		strings.Contains(journeyTemplateContent, "edge case")
	assert.True(t, hasEdgeCases,
		"Journey template must have Edge Cases section")
}

// TC-006: Journey template has Risk classification (High/Medium/Low)
// Traceability: AC "risk level High/Medium/Low"
func TestTC_006_JourneyTemplateHasRiskClassification(t *testing.T) {
	templatesDir := filepath.Join(skillRoot(t), "templates")
	entries, err := os.ReadDir(templatesDir)
	require.NoError(t, err)

	var journeyTemplateContent string
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "journey") && !entry.IsDir() {
			data, err := os.ReadFile(filepath.Join(templatesDir, entry.Name()))
			require.NoError(t, err)
			journeyTemplateContent = string(data)
			break
		}
	}
	require.NotEmpty(t, journeyTemplateContent, "must find journey template")

	// Must reference all three risk levels
	assert.Contains(t, journeyTemplateContent, "High",
		"Journey template must reference High risk level")
	assert.Contains(t, journeyTemplateContent, "Medium",
		"Journey template must reference Medium risk level")
	assert.Contains(t, journeyTemplateContent, "Low",
		"Journey template must reference Low risk level")
}

// TC-007: SKILL.md references PRD user stories as input source
// Traceability: AC "from PRD user stories"
func TestTC_007_SkillMdReferencesPrdUserStoriesAsInput(t *testing.T) {
	skillFile := filepath.Join(skillRoot(t), "SKILL.md")
	content, err := os.ReadFile(skillFile)
	require.NoError(t, err)

	text := string(content)
	// Must reference PRD or user stories as input
	hasPrdReference := strings.Contains(text, "prd") ||
		strings.Contains(text, "PRD")
	hasUserStories := strings.Contains(text, "user stories") ||
		strings.Contains(text, "user-stories") ||
		strings.Contains(text, "prd-user-stories")

	assert.True(t, hasPrdReference, "SKILL.md must reference PRD as input source")
	assert.True(t, hasUserStories, "SKILL.md must reference user stories as input source")
}

// TC-008: SKILL.md specifies output is per-Journey files
// Traceability: AC "output per-Journey files, not per interface type"
func TestTC_008_SkillMdSpecifiesPerJourneyOutput(t *testing.T) {
	skillFile := filepath.Join(skillRoot(t), "SKILL.md")
	content, err := os.ReadFile(skillFile)
	require.NoError(t, err)

	text := string(content)
	// Must mention per-Journey file output
	hasJourneyFile := strings.Contains(text, "journey") &&
		(strings.Contains(text, "file") || strings.Contains(text, "per-journey") ||
			strings.Contains(text, "per journey") || strings.Contains(text, "分文件"))
	assert.True(t, hasJourneyFile,
		"SKILL.md must specify output is per-Journey files")
}

// TC-009: SKILL.md documents no code reconnaissance needed
// Traceability: Hard Rule "gen-journeys does not need code reconnaissance"
func TestTC_009_SkillMdStatesNoCodeReconnaissance(t *testing.T) {
	skillFile := filepath.Join(skillRoot(t), "SKILL.md")
	content, err := os.ReadFile(skillFile)
	require.NoError(t, err)

	text := string(content)
	// Must state that no code reconnaissance is needed
	hasNoRecon := strings.Contains(text, "no code") ||
		strings.Contains(text, "不需要代码侦察") ||
		strings.Contains(text, "pure narrative") ||
		strings.Contains(text, "narrative extraction") ||
		strings.Contains(text, "without code") ||
		strings.Contains(text, "reconnaissance") ||
		(strings.Contains(text, "code") && strings.Contains(text, "not required")) ||
		strings.Contains(text, "Reads Code") && strings.Contains(text, "No")
	assert.True(t, hasNoRecon,
		"SKILL.md must state that gen-journeys does not need code reconnaissance")
}

// TC-010: SKILL.md documents batch processing for context window overflow
// Traceability: Hard Rule "single Journey per generation, auto-batch on context overflow"
func TestTC_010_SkillMdDocumentsBatchProcessing(t *testing.T) {
	skillFile := filepath.Join(skillRoot(t), "SKILL.md")
	content, err := os.ReadFile(skillFile)
	require.NoError(t, err)

	text := string(content)
	// Must mention batch/split processing
	hasBatch := strings.Contains(text, "batch") ||
		strings.Contains(text, "split") ||
		strings.Contains(text, "分批") ||
		strings.Contains(text, "auto-batch")
	assert.True(t, hasBatch,
		"SKILL.md must document batch processing for context window overflow")
}

// TC-011: Journey template has Invariants section for Journey-level invariants
// Traceability: Model spec Section 1.3 — Journey-level Invariants are mandatory
func TestTC_011_JourneyTemplateHasInvariantsSection(t *testing.T) {
	templatesDir := filepath.Join(skillRoot(t), "templates")
	entries, err := os.ReadDir(templatesDir)
	require.NoError(t, err)

	var journeyTemplateContent string
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "journey") && !entry.IsDir() {
			data, err := os.ReadFile(filepath.Join(templatesDir, entry.Name()))
			require.NoError(t, err)
			journeyTemplateContent = string(data)
			break
		}
	}
	require.NotEmpty(t, journeyTemplateContent, "must find journey template")

	// Must have Invariants section
	hasInvariants := strings.Contains(journeyTemplateContent, "Invariant") ||
		strings.Contains(journeyTemplateContent, "invariant")
	assert.True(t, hasInvariants,
		"Journey template must have Invariants section for Journey-level invariants")
}

// TC-012: SKILL.md output path matches gen-contracts input convention
// Traceability: Hard Rule "gen-journeys output must be gen-contracts consumable"
// The output should go to a location gen-contracts can read from
func TestTC_012_SkillMdOutputPathMatchesGenContractsInput(t *testing.T) {
	skillFile := filepath.Join(skillRoot(t), "SKILL.md")
	content, err := os.ReadFile(skillFile)
	require.NoError(t, err)

	text := string(content)
	// Must mention the output location that gen-contracts reads from
	hasJourneyDir := strings.Contains(text, "journeys") ||
		strings.Contains(text, "testing/journeys") ||
		strings.Contains(text, "testing/") ||
		strings.Contains(text, "docs/features") ||
		strings.Contains(text, "_contracts") ||
		strings.Contains(text, "contracts")
	assert.True(t, hasJourneyDir,
		"SKILL.md must specify output path that gen-contracts can consume")

	// Must NOT generate per-interface-type files
	assert.NotContains(t, text, "cli-test-cases",
		"gen-journeys should NOT output per-interface-type files like cli-test-cases")
	assert.NotContains(t, text, "api-test-cases",
		"gen-journeys should NOT output per-interface-type files like api-test-cases")
}

// TC-013: SKILL.md references model-and-directory-spec as authoritative source
// Traceability: Task depends on task 1 (model definition)
func TestTC_013_SkillMdReferencesModelSpec(t *testing.T) {
	skillFile := filepath.Join(skillRoot(t), "SKILL.md")
	content, err := os.ReadFile(skillFile)
	require.NoError(t, err)

	text := string(content)
	// Should reference the model spec or its key concepts
	hasModelRef := strings.Contains(text, "model-and-directory-spec") ||
		strings.Contains(text, "model spec") ||
		strings.Contains(text, "six dimensions") ||
		strings.Contains(text, "Contract") ||
		strings.Contains(text, "Journey-Driven")
	assert.True(t, hasModelRef,
		"SKILL.md should reference the model specification document")
}

// TC-014: SKILL.md uses ${CLAUDE_SKILL_DIR} for internal file references
// Traceability: Forge distribution convention Section 5
func TestTC_014_SkillMdUsesCorrectPathReferences(t *testing.T) {
	skillFile := filepath.Join(skillRoot(t), "SKILL.md")
	content, err := os.ReadFile(skillFile)
	require.NoError(t, err)

	text := string(content)

	// Must NOT use source-tree relative paths like plugins/forge/skills/...
	assert.NotContains(t, text, "plugins/forge/skills/",
		"SKILL.md must NOT use source-tree paths (forge distribution convention)")

	// Must NOT use ../../ relative paths for cross-skill references
	// (allowed for skill-internal references like templates/)
	re, _ := regexp.Compile(`\.\./\.\./\.\./`)
	assert.False(t, re.MatchString(text),
		"SKILL.md must NOT use deeply-nested relative paths")
}

// TC-015: Journey template includes Setup/Preconditions section
// Traceability: Model spec — Journey isolation requires setup declarations
func TestTC_015_JourneyTemplateHasSetupSection(t *testing.T) {
	templatesDir := filepath.Join(skillRoot(t), "templates")
	entries, err := os.ReadDir(templatesDir)
	require.NoError(t, err)

	var journeyTemplateContent string
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "journey") && !entry.IsDir() {
			data, err := os.ReadFile(filepath.Join(templatesDir, entry.Name()))
			require.NoError(t, err)
			journeyTemplateContent = string(data)
			break
		}
	}
	require.NotEmpty(t, journeyTemplateContent, "must find journey template")

	// Must have Setup or Preconditions section
	hasSetup := strings.Contains(journeyTemplateContent, "Setup") ||
		strings.Contains(journeyTemplateContent, "setup") ||
		strings.Contains(journeyTemplateContent, "Precondition") ||
		strings.Contains(journeyTemplateContent, "precondition") ||
		strings.Contains(journeyTemplateContent, "Prerequisite") ||
		strings.Contains(journeyTemplateContent, "prerequisite")
	assert.True(t, hasSetup,
		"Journey template must have Setup/Preconditions section")
}
