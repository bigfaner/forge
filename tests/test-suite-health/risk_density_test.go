//go:build cli_functional

package testsuitehealth

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==============================================================================
// gen-contracts risk-driven density tests — feature: test-capability-v2
// Tests verify that the gen-contracts skill was enhanced with:
//   - risk-density.md rule file with correct density targets (High/Medium/Low)
//   - SKILL.md reads Journey risk_level and applies density rules
//   - Surface-required Outcomes are referenced from surface rules
//   - Schema validation with 1-retry logic is documented
//   - Static Fact Table output path (.forge/fact-table.json) is specified
//   - Inferred boundary Outcomes have source:inferred annotation rule
// ==============================================================================

// projectRootGenContracts returns the forge project root directory.
func projectRootGenContracts(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine test file location")
	}
	// thisFile: .../tests/test-suite-health/risk_density_test.go
	// up 2: test-suite-health -> tests -> project root
	dir := filepath.Join(filepath.Dir(thisFile), "..", "..")
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("cannot resolve project root: %s", err)
	}
	return abs
}

// genContractsSkillRoot returns the gen-contracts skill directory.
func genContractsSkillRoot(t *testing.T) string {
	t.Helper()
	return filepath.Join(projectRootGenContracts(t), "plugins", "forge", "skills", "gen-contracts")
}

// readSkillFile reads the gen-contracts SKILL.md content.
func readSkillFile(t *testing.T) string {
	t.Helper()
	skillFile := filepath.Join(genContractsSkillRoot(t), "SKILL.md")
	content, err := os.ReadFile(skillFile)
	require.NoError(t, err)
	return string(content)
}

// TC-RD-001: risk-density.md rule file exists
// Traceability: AC "新增规则文件 plugins/forge/skills/gen-contracts/rules/risk-density.md"
func TestTC_RD_001_RiskDensityRuleFileExists(t *testing.T) {
	ruleFile := filepath.Join(genContractsSkillRoot(t), "rules", "risk-density.md")
	info, err := os.Stat(ruleFile)
	require.NoError(t, err, "risk-density.md rule file must exist")
	assert.False(t, info.IsDir(), "risk-density.md must be a file, not a directory")
}

// TC-RD-002: risk-density.md defines three risk levels with density targets
// Traceability: AC "High 3-5/13-20, Medium 2-3/8-12, Low 1-2/4-7"
func TestTC_RD_002_RiskDensityDefinesThreeLevels(t *testing.T) {
	ruleFile := filepath.Join(genContractsSkillRoot(t), "rules", "risk-density.md")
	content, err := os.ReadFile(ruleFile)
	require.NoError(t, err)

	text := string(content)

	// Must reference all three risk levels
	assert.Contains(t, text, "High", "risk-density.md must define High risk level")
	assert.Contains(t, text, "Medium", "risk-density.md must define Medium risk level")
	assert.Contains(t, text, "Low", "risk-density.md must define Low risk level")

	// Must have Outcome per Step density targets
	assert.Contains(t, text, "3-5", "High risk must specify 3-5 Outcomes per Step")
	assert.Contains(t, text, "2-3", "Medium risk must specify 2-3 Outcomes per Step")
	assert.Contains(t, text, "1-2", "Low risk must specify 1-2 Outcomes per Step")

	// Must have total test count targets
	assert.Contains(t, text, "13-20", "High risk must specify 13-20 total tests")
	assert.Contains(t, text, "8-12", "Medium risk must specify 8-12 total tests")
	assert.Contains(t, text, "4-7", "Low risk must specify 4-7 total tests")
}

// TC-RD-003: risk-density.md references surface-required Outcome derivation
// Traceability: AC "基于 surface 类型的 required_outcomes + 项目 Fact Table"
func TestTC_RD_003_RiskDensityReferencesSurfaceRequiredOutcomes(t *testing.T) {
	ruleFile := filepath.Join(genContractsSkillRoot(t), "rules", "risk-density.md")
	content, err := os.ReadFile(ruleFile)
	require.NoError(t, err)

	text := string(content)

	// Must reference surface types and their required outcomes
	assert.Contains(t, text, "CLI", "must reference CLI surface")
	assert.Contains(t, text, "API", "must reference API surface")
	assert.Contains(t, text, "TUI", "must reference TUI surface")
	assert.Contains(t, text, "Web", "must reference Web surface")

	// Must reference required outcome names from surface rules
	assert.Contains(t, text, "not-found", "must reference CLI required outcome: not-found")
	assert.Contains(t, text, "already-exists", "must reference CLI required outcome: already-exists")
	assert.Contains(t, text, "unauthorized", "must reference API required outcome: unauthorized")
	assert.Contains(t, text, "timeout", "must reference TUI required outcome: timeout")
	assert.Contains(t, text, "validation-error", "must reference Web required outcome: validation-error")
	assert.Contains(t, text, "session-expired", "must reference Web required outcome: session-expired")
}

// TC-RD-004: SKILL.md references risk-density rule
// Traceability: AC "gen-contracts SKILL.md 增强为：读取 Journey risk_level -> 按密度规则衍生 Outcome"
func TestTC_RD_004_SkillMdReferencesRiskDensityRule(t *testing.T) {
	text := readSkillFile(t)

	assert.Contains(t, text, "risk-density.md",
		"SKILL.md must reference risk-density.md rule file")
	assert.Contains(t, text, "risk_level",
		"SKILL.md must reference Journey risk_level field")
}

// TC-RD-005: SKILL.md defines Outcome density targets matching risk-density.md
// Traceability: AC "High 3-5/13-20, Medium 2-3/8-12, Low 1-2/4-7"
func TestTC_RD_005_SkillMdHasDensityTargets(t *testing.T) {
	text := readSkillFile(t)

	// SKILL.md must include density targets table or inline reference
	assert.Contains(t, text, "3-5", "SKILL.md must reference High density 3-5 Outcomes per Step")
	assert.Contains(t, text, "13-20", "SKILL.md must reference High density 13-20 total")
	assert.Contains(t, text, "2-3", "SKILL.md must reference Medium density 2-3 Outcomes per Step")
	assert.Contains(t, text, "8-12", "SKILL.md must reference Medium density 8-12 total")
	assert.Contains(t, text, "1-2", "SKILL.md must reference Low density 1-2 Outcomes per Step")
	assert.Contains(t, text, "4-7", "SKILL.md must reference Low density 4-7 total")
}

// TC-RD-006: SKILL.md documents schema validation with retry
// Traceability: AC "合约生成后执行 schema 验证... Schema 验证失败时自动重试 1 次"
func TestTC_RD_006_SkillMdHasSchemaValidationWithRetry(t *testing.T) {
	text := readSkillFile(t)

	// Must mention schema validation
	assert.True(t, strings.Contains(text, "Schema validation") || strings.Contains(text, "schema validation") || strings.Contains(text, "Validate Contracts"),
		"SKILL.md must mention schema validation")

	// Must mention retry logic
	assert.True(t, strings.Contains(text, "retry") || strings.Contains(text, "Retry"),
		"SKILL.md must document retry logic for validation failures")

	// Must mention exactly 1 retry
	assert.True(t, strings.Contains(text, "1 automatic retry") || strings.Contains(text, "retry once") || strings.Contains(text, "retries") || strings.Contains(text, "重试"),
		"SKILL.md must specify exactly 1 automatic retry")

	// Must mention pipeline pause on retry failure
	assert.True(t, strings.Contains(text, "pause") || strings.Contains(text, "Pause"),
		"SKILL.md must specify pipeline pause when retry fails")
}

// TC-RD-007: SKILL.md specifies surface-required Outcome derivation from Hard Rules
// Traceability: Hard Rule "必须 Outcome 衍生受 surface rule 的 required_outcomes 约束"
func TestTC_RD_007_SkillMdHasSurfaceRequiredOutcomeHardRule(t *testing.T) {
	text := readSkillFile(t)

	// Must have Hard Rule about surface-required Outcomes
	assert.Contains(t, text, "Surface-required Outcomes",
		"SKILL.md must mention Surface-required Outcomes")

	// Must reference specific surface required outcomes matching Hard Rules
	assert.Contains(t, text, "not-found", "must reference CLI: not-found")
	assert.Contains(t, text, "already-exists", "must reference CLI: already-exists")
	assert.Contains(t, text, "unauthorized", "must reference API: unauthorized")
	assert.Contains(t, text, "validation-error", "must reference WebUI: validation-error")
	assert.Contains(t, text, "session-expired", "must reference WebUI: session-expired")
}

// TC-RD-008: SKILL.md documents inferred Outcome annotation requirement
// Traceability: Hard Rule "LLM 衍生的边界 Outcome 需标注 source: inferred + 推理依据"
func TestTC_RD_008_SkillMdRequiresInferredAnnotation(t *testing.T) {
	text := readSkillFile(t)

	assert.Contains(t, text, "source: inferred",
		"SKILL.md must require source:inferred annotation for LLM-derived boundary Outcomes")
	assert.Contains(t, text, "reasoning",
		"SKILL.md must require reasoning explanation for inferred Outcomes")
}

// TC-RD-009: SKILL.md specifies static Fact Table output to .forge/fact-table.json
// Traceability: AC "静态 Fact Table 写入 .forge/fact-table.json（source: static）"
func TestTC_RD_009_SkillMdSpecifiesFactTableOutput(t *testing.T) {
	text := readSkillFile(t)

	assert.Contains(t, text, "fact-table.json",
		"SKILL.md must specify .forge/fact-table.json as Fact Table output path")
	assert.Contains(t, text, "static",
		"SKILL.md must specify source:static for reconnaissance facts")
}

// TC-RD-010: SKILL.md specifies surface type loading from surface rules
// Traceability: AC "基于 surface 类型的 required_outcomes"
func TestTC_RD_010_SkillMdLoadsSurfaceRules(t *testing.T) {
	text := readSkillFile(t)

	// Must reference surface rule loading via forge surfaces CLI
	assert.True(t,
		strings.Contains(text, "surface-") && strings.Contains(text, "forge surfaces"),
		"SKILL.md must reference surface rule loading via 'forge surfaces'")
}

// TC-RD-011: risk-density.md documents inferred Outcome annotation format
// Traceability: Hard Rule "source: inferred + 推理依据"
func TestTC_RD_011_RiskDensityDocumentsAnnotationFormat(t *testing.T) {
	ruleFile := filepath.Join(genContractsSkillRoot(t), "rules", "risk-density.md")
	content, err := os.ReadFile(ruleFile)
	require.NoError(t, err)

	text := string(content)

	assert.Contains(t, text, "source: inferred",
		"risk-density.md must document source:inferred annotation format")
	assert.Contains(t, text, "reasoning",
		"risk-density.md must require reasoning for inferred Outcomes")
}

// TC-RD-012: SKILL.md description mentions risk-driven density
// Traceability: Task description enhancement
func TestTC_RD_012_SkillMdDescriptionMentionsRiskDensity(t *testing.T) {
	text := readSkillFile(t)

	// Frontmatter description must mention risk-driven density
	assert.True(t,
		strings.Contains(text, "Risk-driven") || strings.Contains(text, "risk-driven") || strings.Contains(text, "risk_level"),
		"SKILL.md description must mention risk-driven Outcome density")
}

// TC-RD-013: SKILL.md uses correct path references (forge distribution convention)
// Traceability: Forge distribution convention -- no source-tree paths
func TestTC_RD_013_SkillMdUsesCorrectPaths(t *testing.T) {
	text := readSkillFile(t)

	assert.NotContains(t, text, "plugins/forge/skills/",
		"SKILL.md must NOT use source-tree paths (forge distribution convention)")
}
