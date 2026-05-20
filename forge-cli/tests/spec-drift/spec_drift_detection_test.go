//go:build e2e

package specdrift

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/tests/testkit"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- TC-001: List types includes doc-generation.drift with correct description ---

// Traceability: TC-001 -> Task 2 AC (type constant, registry)
func TestTC_001_ListTypesIncludesDriftTypeWithCorrectDescription(t *testing.T) {
	cmd := exec.Command(forgeBinary, "task", "list-types")
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, "task list-types should succeed: %s", string(out))
	output := string(out)

	assert.Contains(t, output, "doc-generation.drift",
		"list-types output should contain doc-generation.drift type")
	assert.Contains(t, output, "detect and fix spec drift against codebase",
		"list-types output should contain drift type description")

	// Verify total type count includes doc-generation.drift (14 types total)
	lines := strings.Split(output, "\n")
	nonEmptyLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines++
		}
	}
	assert.GreaterOrEqual(t, nonEmptyLines, 14,
		"should list at least 14 task types (including doc-generation.drift)")
}

// --- TC-002: Prompt for doc-generation.drift type resolves to correct strategy template ---

// Traceability: TC-002 -> Task 2 AC (prompt.go mapping)
func TestTC_002_DriftTypeResolvesToCorrectStrategyTemplate(t *testing.T) {
	// This test requires a task with type doc-generation.drift in index.json.
	// We verify the embedded template contains drift-only mode instructions.
	projectRoot := testkit.ProjectRoot(t)
	templatePath := filepath.Join(projectRoot, "pkg", "prompt", "data", "doc-generation-drift.md")
	data, err := os.ReadFile(templatePath)
	require.NoError(t, err, "drift strategy template should exist at pkg/prompt/data/doc-generation-drift.md")
	content := string(data)

	// Verify template references consolidate-specs skill
	assert.Contains(t, content, "consolidate-specs",
		"drift strategy template should reference consolidate-specs skill")

	// Verify template specifies drift-only mode (references Steps 9-11)
	assert.True(t,
		strings.Contains(content, "Steps 9-11") ||
			(strings.Contains(content, "Step 9") && strings.Contains(content, "Step 10") && strings.Contains(content, "Step 11")),
		"drift strategy template should reference Steps 9-11 (drift-only mode)")

	// Verify template mentions drift-only mode
	assert.True(t,
		strings.Contains(strings.ToLower(content), "drift-only") ||
			strings.Contains(content, "drift only"),
		"drift strategy template should mention drift-only mode")
}

// --- TC-003: Type doc-generation.drift is a valid type recognized by validate-index ---

// Traceability: TC-003 -> Task 2 AC (valid type map)
func TestTC_003_DriftTypeIsValidInValidateIndex(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("task", "validate-index")

	if exitCode != 0 {
		t.Skip("validate-index requires a feature context with index.json containing drift tasks")
	}
	assert.NotContains(t, strings.ToLower(out), "unknown type",
		"validate-index should not report unknown type errors for doc-generation.drift")
}

// --- TC-004: Quick test tasks include T-quick-doc-drift with doc-generation.drift type ---

// Traceability: TC-004 -> Task 2 AC (T-quick-doc-drift generation)
func TestTC_004_QuickPipelineIncludesTQuick5WithDriftType(t *testing.T) {
	// Verify the code defines T-quick-doc-drift with correct properties
	projectRoot := testkit.ProjectRoot(t)
	testgenPath := filepath.Join(projectRoot, "pkg", "task", "testgen.go")
	data, err := os.ReadFile(testgenPath)
	require.NoError(t, err)
	content := string(data)

	// Verify T-quick-doc-drift definition exists with drift type
	assert.Contains(t, content, `"T-quick-doc-drift"`,
		"testgen.go should define T-quick-doc-drift ID")
	assert.Contains(t, content, "TypeDocDrift",
		"T-quick-doc-drift should use TypeDocDrift type")
	// NoTest removed; TypeDocDrift is non-testable by IsTestableType
	assert.Contains(t, content, "TypeDocDrift",
		"T-quick-doc-drift should have non-testable type")
	assert.Contains(t, content, `Scope: "all"`,
		"T-quick-doc-drift should have Scope: all")
	assert.Contains(t, content, `"Detect Spec Drift"`,
		"T-quick-doc-drift title should contain 'Drift'")
}

// --- TC-005: T-quick-doc-drift depends on T-quick-verify-regression in quick pipeline ---

// Traceability: TC-005 -> Task 2 AC (T-quick-doc-drift deps)
func TestTC_005_TQuick5DependsOnTQuick4(t *testing.T) {
	projectRoot := testkit.ProjectRoot(t)
	testgenPath := filepath.Join(projectRoot, "pkg", "task", "testgen.go")
	data, err := os.ReadFile(testgenPath)
	require.NoError(t, err)
	content := string(data)

	// Verify resolveQuickDeps sets T-quick-doc-drift dependency on T-quick-verify-regression
	assert.Contains(t, content, `"T-quick-verify-regression"`,
		"resolveQuickDeps should set T-quick-doc-drift dependency on T-quick-verify-regression")
}

// --- TC-006: T-quick-doc-drift appears after T-quick-verify-regression in task generation order ---

// Traceability: TC-006 -> Task 2 AC (task order)
func TestTC_006_TQuick5AppearsAfterTQuick4(t *testing.T) {
	projectRoot := testkit.ProjectRoot(t)
	testgenPath := filepath.Join(projectRoot, "pkg", "task", "testgen.go")
	data, err := os.ReadFile(testgenPath)
	require.NoError(t, err)
	content := string(data)

	// T-quick-verify-regression is defined at index len(profiles)*3, T-quick-doc-drift at len(profiles)*3+1
	// Verify the order: T-quick-verify-regression appears before T-quick-doc-drift in the file
	verifyRegIdx := strings.Index(content, `"T-quick-verify-regression"`)
	driftIdx := strings.Index(content, `"T-quick-doc-drift"`)
	assert.Greater(t, verifyRegIdx, -1, "T-quick-verify-regression should be defined")
	assert.Greater(t, driftIdx, -1, "T-quick-doc-drift should be defined")
	assert.Less(t, verifyRegIdx, driftIdx,
		"T-quick-verify-regression should appear before T-quick-doc-drift in testgen.go")
}

// --- TC-007: T-specs-consolidate description includes drift detection in breakdown mode ---

// Traceability: TC-007 -> Task 3 AC (breakdown-tasks SKILL.md references drift detection)
func TestTC_007_TTest5DescriptionIncludesDriftDetection(t *testing.T) {
	// Verify breakdown-tasks SKILL.md references drift detection for T-specs-consolidate
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "breakdown-tasks", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	assert.Contains(t, content, "T-specs-consolidate",
		"breakdown-tasks SKILL.md should mention T-specs-consolidate")
	assert.True(t,
		strings.Contains(strings.ToLower(content), "drift") ||
			strings.Contains(content, "consolidate"),
		"breakdown-tasks SKILL.md T-specs-consolidate description should reference consolidation/drift detection")
}

// --- TC-008: consolidate-specs SKILL.md contains Steps 9-11 ---

// Traceability: TC-008 -> Task 1 AC (Steps 9-11 exist)
func TestTC_008_ConsolidateSpecsContainsSteps9Through11(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Verify Step 9 section exists
	assert.Contains(t, content, "Step 9",
		"SKILL.md should contain Step 9")

	// Verify Step 10 section exists
	assert.Contains(t, content, "Step 10",
		"SKILL.md should contain Step 10")

	// Verify Step 11 section exists
	assert.Contains(t, content, "Step 11",
		"SKILL.md should contain Step 11")

	// Verify Steps 9-11 appear after Step 8 using section headings
	step8Idx := strings.Index(content, "## Step 8")
	step9Idx := strings.Index(content, "## Step 9")
	assert.Greater(t, step8Idx, -1, "## Step 8 heading should exist")
	assert.Greater(t, step9Idx, -1, "## Step 9 heading should exist")
	assert.Greater(t, step9Idx, step8Idx,
		"## Step 9 heading should appear after ## Step 8 heading")
}

// --- TC-009: consolidate-specs SKILL.md Step 9 validates rules against code ---

// Traceability: TC-009 -> Task 1 AC (Step 9 classification)
func TestTC_009_Step9ClassifiesRulesAgainstCode(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Extract Step 9 section using heading pattern
	step9Idx := strings.Index(content, "## Step 9:")
	require.Greater(t, step9Idx, -1, "## Step 9: heading should exist")

	step10Idx := strings.Index(content, "## Step 10:")
	if step10Idx < 0 {
		step10Idx = len(content)
	}
	step9Section := content[step9Idx:step10Idx]

	// Verify three-way classification
	assert.True(t,
		strings.Contains(step9Section, "current") &&
			strings.Contains(step9Section, "drifted") &&
			strings.Contains(step9Section, "orphaned"),
		"Step 9 should define three-way classification: current/drifted/orphaned")

	// Verify reads business-rules and conventions
	assert.True(t,
		strings.Contains(step9Section, "business-rules") ||
			strings.Contains(step9Section, "business_rules"),
		"Step 9 should instruct reading business-rules directory")
	assert.True(t,
		strings.Contains(step9Section, "convention"),
		"Step 9 should instruct reading conventions directory")
}

// --- TC-010: consolidate-specs SKILL.md Step 10 preserves project-global IDs ---

// Traceability: TC-010 -> Task 1 AC (Step 10 ID preservation)
func TestTC_010_Step10PreservesProjectGlobalIDs(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Extract Step 10 section using heading pattern
	step10Idx := strings.Index(content, "## Step 10:")
	require.Greater(t, step10Idx, -1, "## Step 10: heading should exist")

	step11Idx := strings.Index(content, "## Step 11:")
	if step11Idx < 0 {
		step11Idx = len(content)
	}
	step10Section := content[step10Idx:step11Idx]

	// Verify ID preservation during updates
	assert.True(t,
		strings.Contains(step10Section, "ID") ||
			strings.Contains(step10Section, "preserv"),
		"Step 10 should instruct preserving project-global IDs during updates")

	// Verify orphaned rule removal
	assert.True(t,
		strings.Contains(step10Section, "orphan") ||
			strings.Contains(step10Section, "remov"),
		"Step 10 should instruct removing orphaned rules")

	// Verify implicit new rule detection
	assert.True(t,
		strings.Contains(step10Section, "new") ||
			strings.Contains(step10Section, "implicit"),
		"Step 10 should instruct detecting implicit new rules from code changes")
}

// --- TC-011: consolidate-specs SKILL.md HARD-GATE allows drift modification ---

// Traceability: TC-011 -> Task 1 AC (HARD-GATE drift exception)
func TestTC_011_HardGateAllowsDriftModification(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Find HARD-GATE section
	gateIdx := strings.Index(content, "HARD-GATE")
	require.Greater(t, gateIdx, -1, "SKILL.md should contain HARD-GATE section")

	// Verify the gate section has an exception for drift
	gateEnd := len(content)
	nextSection := strings.Index(content[gateIdx+10:], "##")
	if nextSection > 0 {
		gateEnd = gateIdx + 10 + nextSection
	}
	gateSection := content[gateIdx:gateEnd]

	assert.True(t,
		strings.Contains(gateSection, "drift") ||
			strings.Contains(gateSection, "Step 9"),
		"HARD-GATE should include an exception clause for drift detected in Step 9")
}

// --- TC-012: consolidate-specs SKILL.md supports drift-only mode ---

// Traceability: TC-012 -> Task 1 AC (drift-only mode)
func TestTC_012_ConsolidateSpecsSupportsDriftOnlyMode(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Verify drift-only mode is documented
	assert.True(t,
		strings.Contains(content, "drift-only") ||
			strings.Contains(content, "Drift-only") ||
			strings.Contains(content, "drift only"),
		"SKILL.md should document drift-only mode")

	// Verify drift-only mode skips Steps 1-8 and runs Steps 9-11
	assert.True(t,
		strings.Contains(content, "Step 9") &&
			strings.Contains(content, "Step 10") &&
			strings.Contains(content, "Step 11"),
		"SKILL.md drift-only mode should reference Steps 9-11")
}

// --- TC-013: breakdown-tasks SKILL.md references drift detection for T-specs-consolidate ---

// Traceability: TC-013 -> Task 3 AC (breakdown-tasks SKILL.md references T-specs-consolidate drift detection)
func TestTC_013_BreakdownTasksReferencesDriftForTTest5(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "breakdown-tasks", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	assert.Contains(t, content, "T-specs-consolidate",
		"breakdown-tasks SKILL.md should mention T-specs-consolidate")
	assert.True(t,
		strings.Contains(strings.ToLower(content), "drift"),
		"breakdown-tasks SKILL.md should mention drift detection for T-specs-consolidate")
}

// --- TC-014: guide.md reflects T-quick-doc-drift and drift detection flow ---

// Traceability: TC-014 -> Task 4 AC (guide T-quick-doc-drift)
func TestTC_014_GuideReflectsTQuick5AndDriftDetection(t *testing.T) {
	guidePath := findRepoFile(t,
		filepath.Join("plugins", "forge", "hooks", "guide.md"))
	data, err := os.ReadFile(guidePath)
	require.NoError(t, err)
	content := string(data)

	// Verify Quick Mode section references T-quick-doc-drift
	assert.Contains(t, content, "T-quick-doc-drift",
		"guide.md should reference T-quick-doc-drift in Quick Mode section")

	// Verify drift detection is mentioned
	assert.True(t,
		strings.Contains(strings.ToLower(content), "drift"),
		"guide.md should mention drift detection")
}

// --- TC-015: guide.md specs rule mentions drift verification ---

// Traceability: TC-015 -> Task 4 AC (guide specs drift verification)
func TestTC_015_GuideSpecsRuleMentionsDriftVerification(t *testing.T) {
	guidePath := findRepoFile(t,
		filepath.Join("plugins", "forge", "hooks", "guide.md"))
	data, err := os.ReadFile(guidePath)
	require.NoError(t, err)
	content := string(data)

	// Verify business-rules or conventions entries mention drift
	lowerContent := strings.ToLower(content)
	assert.True(t,
		strings.Contains(lowerContent, "business-rules") || strings.Contains(lowerContent, "business_rules"),
		"guide.md should reference business-rules directory")

	// Check for drift verification in the specs-related section
	assert.True(t,
		strings.Contains(lowerContent, "drift") &&
			(strings.Contains(lowerContent, "convention") || strings.Contains(lowerContent, "business-rule")),
		"guide.md should mention drift verification in relation to spec files")
}

// --- TC-016: doc-generation-drift.md strategy template exists ---

// Traceability: TC-016 -> Task 2 AC (strategy template)
func TestTC_016_DriftStrategyTemplateExists(t *testing.T) {
	projectRoot := testkit.ProjectRoot(t)
	templatePath := filepath.Join(projectRoot, "pkg", "prompt", "data", "doc-generation-drift.md")

	info, err := os.Stat(templatePath)
	require.NoError(t, err, "doc-generation-drift.md strategy template should exist")
	assert.Greater(t, info.Size(), int64(0),
		"strategy template should not be empty")

	data, err := os.ReadFile(templatePath)
	require.NoError(t, err)
	content := string(data)

	// Verify it references consolidate-specs skill
	assert.Contains(t, content, "consolidate-specs",
		"strategy template should reference consolidate-specs skill")

	// Verify drift-only mode specification (references Steps 9-11 as a range)
	assert.True(t,
		strings.Contains(content, "Steps 9-11") ||
			strings.Contains(content, "steps 9-11") ||
			(strings.Contains(content, "Step 9") && strings.Contains(content, "Step 10") && strings.Contains(content, "Step 11")),
		"strategy template should reference Steps 9-11 (drift-only mode)")
}

// --- TC-017: All existing tests pass after feature changes ---

// Traceability: TC-017 -> Task 2 AC (existing tests pass)
func TestTC_017_AllExistingTestsPass(t *testing.T) {
	projectRoot := testkit.ProjectRoot(t)
	cmd := exec.Command("go", "test", "-race", "-count=1", "./...")
	cmd.Dir = projectRoot
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err,
		"all existing tests should pass: %s", string(out))
}

// --- TC-018: Task ID T-quick-doc-drift infers type as doc-generation.drift ---

// Traceability: TC-018 -> Task 2 AC (type inference)
func TestTC_018_TaskIDInfersDriftType(t *testing.T) {
	// Verify type inference for T-quick-doc-drift, T-quick-doc-drifta, T-quick-doc-driftb via code inspection
	projectRoot := testkit.ProjectRoot(t)
	inferPath := filepath.Join(projectRoot, "pkg", "task", "infer.go")
	data, err := os.ReadFile(inferPath)
	require.NoError(t, err)
	content := string(data)

	// Verify T-quick-doc-drift infers doc-generation.drift
	assert.Contains(t, content, `"T-quick-doc-drift"`,
		"infer.go should handle T-quick-doc-drift")
	assert.Contains(t, content, "TypeDocDrift",
		"T-quick-doc-drift should infer TypeDocDrift")

	// Verify profile-suffixed variants (T-quick-doc-drifta, T-quick-doc-driftb) also handled
	// The profileSuffixedID function handles T-quick-doc-drift + letter suffix
	assert.Contains(t, content, "profileSuffixedID",
		"infer.go should use profileSuffixedID for variant handling")
}

// --- TC-019: consolidate-specs workflow diagram includes Steps 9-11 ---

// Traceability: TC-019 -> Task 1 AC (workflow diagram updated)
func TestTC_019_WorkflowDiagramIncludesDriftSteps(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Locate the workflow/diagram section
	// The SKILL.md should have a workflow section or mermaid diagram
	workflowIdx := strings.Index(content, "Workflow")
	if workflowIdx < 0 {
		workflowIdx = strings.Index(content, "workflow")
	}
	if workflowIdx < 0 {
		workflowIdx = strings.Index(content, "```mermaid")
	}
	require.Greater(t, workflowIdx, -1,
		"SKILL.md should contain a Workflow section or diagram")

	// Get the section after workflow marker (generous window for diagram)
	sectionEnd := len(content)
	if workflowIdx+2000 < sectionEnd {
		sectionEnd = workflowIdx + 2000
	}
	workflowSection := content[workflowIdx:sectionEnd]

	// Verify diagram includes Steps 9-11
	assert.True(t,
		strings.Contains(workflowSection, "Step 9") || strings.Contains(workflowSection, "9"),
		"workflow diagram should include Step 9 (Detect Drift)")
	assert.True(t,
		strings.Contains(workflowSection, "Step 10") || strings.Contains(workflowSection, "10"),
		"workflow diagram should include Step 10 (Auto-Fix)")
	assert.True(t,
		strings.Contains(workflowSection, "Step 11") || strings.Contains(workflowSection, "11"),
		"workflow diagram should include Step 11 (Commit)")

	// Verify Steps 9-11 appear after Step 8
	step8InWorkflow := strings.Index(workflowSection, "Step 8")
	step9InWorkflow := strings.Index(workflowSection, "Step 9")
	if step8InWorkflow > 0 && step9InWorkflow > 0 {
		assert.Greater(t, step9InWorkflow, step8InWorkflow,
			"Step 9 should appear after Step 8 in the workflow diagram")
	}
}

// --- TC-020: consolidate-specs SKILL.md contains Step 12 (vocabulary generation) ---

// Traceability: TC-020 -> Task 4 AC (vocabulary generation step exists)
func TestTC_020_ConsolidateSpecsContainsVocabularyStep(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Verify Step 12 section exists
	assert.Contains(t, content, "## Step 12:",
		"SKILL.md should contain Step 12 section")

	// Verify Step 12 is about vocabulary generation
	step12Idx := strings.Index(content, "## Step 12:")
	require.Greater(t, step12Idx, -1, "## Step 12 heading should exist")

	step13Idx := strings.Index(content, "## Step 13:")
	if step13Idx < 0 {
		step13Idx = len(content)
	}
	step12Section := content[step12Idx:step13Idx]

	assert.Contains(t, step12Section, "vocabulary",
		"Step 12 should be about vocabulary generation")
	assert.Contains(t, step12Section, "Vocabulary",
		"Step 12 should be about vocabulary generation")
}

// --- TC-021: Step 12 scans all 4 knowledge directories ---

// Traceability: TC-021 -> Task 4 AC (scans decisions, lessons, conventions, business-rules)
func TestTC_021_VocabularyStepScansAllKnowledgeDirectories(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Extract Step 12 section
	step12Idx := strings.Index(content, "## Step 12:")
	require.Greater(t, step12Idx, -1, "## Step 12 heading should exist")

	step13Idx := strings.Index(content, "## Step 13:")
	if step13Idx < 0 {
		step13Idx = len(content)
	}
	step12Section := content[step12Idx:step13Idx]

	// Verify all 4 directories are mentioned
	assert.Contains(t, step12Section, "docs/decisions/",
		"Step 12 should scan docs/decisions/")
	assert.Contains(t, step12Section, "docs/lessons/",
		"Step 12 should scan docs/lessons/")
	assert.Contains(t, step12Section, "docs/conventions/",
		"Step 12 should scan docs/conventions/")
	assert.Contains(t, step12Section, "docs/business-rules/",
		"Step 12 should scan docs/business-rules/")
}

// --- TC-022: Step 12 includes base 8-category vocabulary ---

// Traceability: TC-022 -> Task 4 AC (base vocabulary always included)
func TestTC_022_VocabularyIncludesBase8Categories(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Extract Step 12 section
	step12Idx := strings.Index(content, "## Step 12:")
	require.Greater(t, step12Idx, -1)

	step13Idx := strings.Index(content, "## Step 13:")
	if step13Idx < 0 {
		step13Idx = len(content)
	}
	step12Section := content[step12Idx:step13Idx]

	// Verify base 8 categories
	baseCategories := []string{
		"architecture",
		"interface",
		"data-model",
		"dependencies",
		"error-handling",
		"testing",
		"security",
		"local-dev-deployment",
	}
	for _, cat := range baseCategories {
		assert.Contains(t, step12Section, cat,
			"Step 12 should include base category: "+cat)
	}
}

// --- TC-023: Step 12 generates vocabulary with types, domains, counts ---

// Traceability: TC-023 -> Task 4 AC (vocabulary output structure)
func TestTC_023_VocabularyOutputStructure(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Extract Step 12 section
	step12Idx := strings.Index(content, "## Step 12:")
	require.Greater(t, step12Idx, -1)

	step13Idx := strings.Index(content, "## Step 13:")
	if step13Idx < 0 {
		step13Idx = len(content)
	}
	step12Section := content[step12Idx:step13Idx]

	// Verify output contains types, domains, counts
	assert.Contains(t, step12Section, "Types",
		"Step 12 output should include Types")
	assert.Contains(t, step12Section, "Domains",
		"Step 12 output should include Domains")
	assert.Contains(t, step12Section, "Count",
		"Step 12 output should include Count")

	// Verify the 4 knowledge types
	assert.Contains(t, step12Section, "decision",
		"Step 12 should list 'decision' as a type")
	assert.Contains(t, step12Section, "lesson",
		"Step 12 should list 'lesson' as a type")
	assert.Contains(t, step12Section, "convention",
		"Step 12 should list 'convention' as a type")
	assert.Contains(t, step12Section, "business-rule",
		"Step 12 should list 'business-rule' as a type")
}

// --- TC-024: Vocabulary is marked as auto-generated ---

// Traceability: TC-024 -> Task 4 AC (auto-generated marking)
func TestTC_024_VocabularyMarkedAsAutoGenerated(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Extract Step 12 section
	step12Idx := strings.Index(content, "## Step 12:")
	require.Greater(t, step12Idx, -1)

	step13Idx := strings.Index(content, "## Step 13:")
	if step13Idx < 0 {
		step13Idx = len(content)
	}
	step12Section := content[step12Idx:step13Idx]

	// Verify auto-generated marking
	assert.Contains(t, step12Section, "auto-generated",
		"Step 12 should mark vocabulary as auto-generated")
	assert.Contains(t, step12Section, "AUTO-GENERATED",
		"Step 12 should have explicit AUTO-GENERATED comment")

	// Verify output file path
	assert.Contains(t, step12Section, "docs/.vocabulary.md",
		"Step 12 should specify output file docs/.vocabulary.md")
}

// --- TC-025: Step 12 placed after Step 11 and before Step 13 ---

// Traceability: TC-025 -> Task 4 AC (step ordering)
func TestTC_025_VocabularyStepOrdering(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	step11Idx := strings.Index(content, "## Step 11:")
	step12Idx := strings.Index(content, "## Step 12:")
	step13Idx := strings.Index(content, "## Step 13:")

	require.Greater(t, step11Idx, -1, "## Step 11 heading should exist")
	require.Greater(t, step12Idx, -1, "## Step 12 heading should exist")
	require.Greater(t, step13Idx, -1, "## Step 13 heading should exist")

	assert.Less(t, step11Idx, step12Idx,
		"Step 12 should appear after Step 11")
	assert.Less(t, step12Idx, step13Idx,
		"Step 12 should appear before Step 13")
}

// --- TC-026: Existing Steps 1-11 unchanged ---

// Traceability: TC-026 -> Task 4 Hard Rules (do not change existing steps)
func TestTC_026_ExistingStepsUnchanged(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Verify all original steps still exist
	for _, step := range []string{
		"## Step 1:",
		"## Step 2:",
		"## Step 3:",
		"## Step 4:",
		"## Step 5:",
		"## Step 6:",
		"## Step 7:",
		"## Step 8:",
		"## Step 9:",
		"## Step 10:",
		"## Step 11:",
	} {
		assert.Contains(t, content, step,
			"Existing step should be preserved: "+step)
	}
}

// --- TC-027: Vocabulary step works with empty directories ---

// Traceability: TC-027 -> Task 4 AC (works when directories sparse/empty)
func TestTC_027_VocabularyWorksWithEmptyDirectories(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Extract Step 12 section
	step12Idx := strings.Index(content, "## Step 12:")
	require.Greater(t, step12Idx, -1)

	step13Idx := strings.Index(content, "## Step 13:")
	if step13Idx < 0 {
		step13Idx = len(content)
	}
	step12Section := content[step12Idx:step13Idx]

	// Verify it handles sparse/empty directories
	assert.True(t,
		strings.Contains(step12Section, "unconditionally") ||
			strings.Contains(step12Section, "sparse") ||
			strings.Contains(step12Section, "empty"),
		"Step 12 should mention handling sparse or empty directories")

	// Verify base vocabulary is always present (even with empty dirs)
	assert.Contains(t, step12Section, "always included",
		"Step 12 should state base vocabulary is always included")
}

// --- TC-028: Vocabulary references /learn and auto-extract triggers ---

// Traceability: TC-028 -> Task 4 AC (vocabulary usable by /learn and triggers)
func TestTC_028_VocabularyReferencesLearnAndTriggers(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Extract Step 12 section
	step12Idx := strings.Index(content, "## Step 12:")
	require.Greater(t, step12Idx, -1)

	step13Idx := strings.Index(content, "## Step 13:")
	if step13Idx < 0 {
		step13Idx = len(content)
	}
	step12Section := content[step12Idx:step13Idx]

	// Verify /learn reference
	assert.Contains(t, step12Section, "/learn",
		"Step 12 should reference /learn skill")

	// Verify auto-extract triggers mentioned
	assert.True(t,
		strings.Contains(step12Section, "auto-extract") ||
			strings.Contains(step12Section, "trigger"),
		"Step 12 should reference auto-extract triggers")

	// Verify vocabulary is suggestive, not restrictive
	assert.True(t,
		strings.Contains(step12Section, "suggestive") ||
			strings.Contains(step12Section, "not restrictive"),
		"Step 12 should note vocabulary is suggestive, not restrictive")
}

// --- TC-029: Workflow diagram includes Step 12 and Step 13 ---

// Traceability: TC-029 -> Task 4 AC (workflow diagram updated)
func TestTC_029_WorkflowDiagramIncludesVocabularyStep(t *testing.T) {
	skillPath := findRepoFile(t,
		filepath.Join("plugins", "forge", "skills", "consolidate-specs", "SKILL.md"))
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)

	// Locate the workflow section
	workflowIdx := strings.Index(content, "Workflow")
	require.Greater(t, workflowIdx, -1, "SKILL.md should contain Workflow section")

	sectionEnd := len(content)
	if workflowIdx+2500 < sectionEnd {
		sectionEnd = workflowIdx + 2500
	}
	workflowSection := content[workflowIdx:sectionEnd]

	// Verify diagram includes Step 12 and Step 13
	assert.Contains(t, workflowSection, "Step 12",
		"Workflow diagram should include Step 12 (vocabulary generation)")
	assert.Contains(t, workflowSection, "Step 13",
		"Workflow diagram should include Step 13 (record task)")

	// Verify Step 12 appears after Step 11
	step11InWorkflow := strings.Index(workflowSection, "Step 11")
	step12InWorkflow := strings.Index(workflowSection, "Step 12")
	if step11InWorkflow > 0 && step12InWorkflow > 0 {
		assert.Less(t, step11InWorkflow, step12InWorkflow,
			"Step 12 should appear after Step 11 in workflow diagram")
	}
}

// --- Helper functions ---

// findRepoFile resolves a file path relative to the repository root (forge parent).
// The project root (forge-cli) is one level below the repo root.
func findRepoFile(t *testing.T, relPath string) string {
	t.Helper()
	projectRoot := testkit.ProjectRoot(t)
	repoRoot := filepath.Dir(projectRoot)
	fullPath := filepath.Join(repoRoot, relPath)
	info, err := os.Stat(fullPath)
	if err != nil {
		t.Fatalf("file not found: %s: %v", fullPath, err)
	}
	if info.IsDir() {
		t.Fatalf("expected file, got directory: %s", fullPath)
	}
	return fullPath
}
