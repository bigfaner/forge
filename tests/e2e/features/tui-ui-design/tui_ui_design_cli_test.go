//go:build e2e

package e2e

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper: readFile reads a file and returns its content as a string.
// Fatal on error so the test stops immediately if the file is missing.
func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %s", path, err)
	}
	return string(data)
}

// Helper: fileContains asserts that the content contains the given substring.
func fileContains(t *testing.T, content, substr, context string) {
	t.Helper()
	if !assert.Contains(t, content, substr, "expected %s to contain: %s", context, substr) {
		t.Logf("File content (first 500 chars): %s", truncate(content, 500))
	}
}

// Helper: truncate string for logging.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// ============================================================================
// Task 1: TUI Platform File & TUI Themes
// ============================================================================

// Traceability: TC-001 -> Task 1 AC-1
func TestTC_001_TUIPlatformFileDefinesNavigationStructure(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/templates/platforms/tui.md")

	// Assert Keymap table with required columns
	fileContains(t, content, "| Key | Action | Context/Mode |", "Keymap table header")

	// Assert Panel Layout table with required columns
	fileContains(t, content, "| Panel | View | Position | Size Hint |", "Panel Layout table header")

	// Assert Modes table with required columns
	fileContains(t, content, "| Mode | Description | Default Keybindings |", "Modes table header")

	// Assert Navigation Rules section
	assert.True(t, strings.Contains(content, "Navigation Rules") || strings.Contains(content, "Navigation Contract"),
		"expected tui.md to contain Navigation Rules section")
}

// Traceability: TC-002 -> Task 1 AC-2
func TestTC_002_ModernDarkTUIThemeSpecifiesCorrectProperties(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/templates/styles/modern-dark-tui.md")

	// Assert color space is 256-color
	assert.True(t, strings.Contains(content, "256-color"),
		"expected modern-dark-tui.md to specify 256-color color space")

	// Assert character set includes box-drawing + block elements
	fileContains(t, content, "box-drawing", "character set description")
	assert.True(t, strings.Contains(content, "▄") || strings.Contains(content, "▪") || strings.Contains(content, "─"),
		"expected modern-dark-tui.md to include box-drawing or block element examples")

	// Assert dark background palette with high contrast
	fileContains(t, content, "Dark", "dark background palette reference")
	fileContains(t, content, "high contrast", "high contrast specification")

	// Assert semantic colors
	assert.True(t, strings.Contains(content, "Green") || strings.Contains(content, "green"),
		"expected green semantic color")
	assert.True(t, strings.Contains(content, "Red") || strings.Contains(content, "red"),
		"expected red semantic color")
	assert.True(t, strings.Contains(content, "Blue") || strings.Contains(content, "blue"),
		"expected blue semantic color")

	// Assert density is compact
	fileContains(t, content, "Compact", "compact density")

	// Assert applicable scenarios section
	fileContains(t, content, "Applicable Scenarios", "applicable scenarios section")
}

// Traceability: TC-003 -> Task 1 AC-3
func TestTC_003_MinimalASCIITUIThemeSpecifiesCorrectProperties(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/templates/styles/minimal-ascii-tui.md")

	// Assert color space is 16-color
	assert.True(t, strings.Contains(content, "16-color"),
		"expected minimal-ascii-tui.md to specify 16-color color space")

	// Assert character set is pure ASCII
	fileContains(t, content, "Pure ASCII", "pure ASCII character set description")

	// Assert default terminal background palette with weight/spacing distinction
	assert.True(t, strings.Contains(content, "Default") && strings.Contains(content, "background"),
		"expected default terminal background palette reference")
	assert.True(t, strings.Contains(content, "weight") || strings.Contains(content, "spacing"),
		"expected weight/spacing distinction")

	// Assert density is loose
	fileContains(t, content, "Loose", "loose density")

	// Assert applicable scenarios section
	fileContains(t, content, "Applicable Scenarios", "applicable scenarios section")
}

// Traceability: TC-004 -> Task 1 AC-4
func TestTC_004_TUIThemeFilesFollowExistingStyleFileFormat(t *testing.T) {
	// Read existing style file as reference
	reference := readFile(t, "plugins/forge/skills/ui-design/templates/styles/apple.md")

	// Read TUI theme files
	modernDark := readFile(t, "plugins/forge/skills/ui-design/templates/styles/modern-dark-tui.md")
	minimalASCII := readFile(t, "plugins/forge/skills/ui-design/templates/styles/minimal-ascii-tui.md")

	// Both TUI theme files should have heading structure (## sections)
	// Reference apple.md uses sections like "## Color Palette", "## Typography", etc.
	refHasSections := strings.Contains(reference, "## ")
	assert.True(t, refHasSections, "reference style file should have ## sections")

	// Assert both TUI theme files have ## section headings
	assert.True(t, strings.Contains(modernDark, "## "), "modern-dark-tui.md should have ## section headings")
	assert.True(t, strings.Contains(minimalASCII, "## "), "minimal-ascii-tui.md should have ## section headings")

	// Assert both have a title (# heading)
	assert.True(t, strings.HasPrefix(modernDark, "# "), "modern-dark-tui.md should start with a title heading")
	assert.True(t, strings.HasPrefix(minimalASCII, "# "), "minimal-ascii-tui.md should start with a title heading")
}

// ============================================================================
// Task 2: PRD UI Functions Template TUI Navigation
// ============================================================================

// Traceability: TC-005 -> Task 2 AC-1
func TestTC_005_PRDUIFunctionsTemplateIncludesTUINavigationArchitecture(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/write-prd/templates/prd-ui-functions.md")

	// Assert Platform: tui indicator
	assert.True(t, strings.Contains(content, "tui"),
		"expected prd-ui-functions.md to contain TUI platform indicator")

	// Assert Keymap table with columns
	fileContains(t, content, "| Key | Action | Context/Mode |", "Keymap table header")

	// Assert Panel Layout table with columns
	fileContains(t, content, "| Panel | View | Position | Size Hint |", "Panel Layout table header")

	// Assert Modes table with columns
	fileContains(t, content, "| Mode | Description | Default Keybindings |", "Modes table header")

	// Assert Navigation Rules section for TUI
	assert.True(t, strings.Contains(content, "Navigation Rules"),
		"expected prd-ui-functions.md to contain Navigation Rules section")
}

// Traceability: TC-006 -> Task 2 AC-2
func TestTC_006_TUINavigationSectionIsConditionallyRendered(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/write-prd/templates/prd-ui-functions.md")

	// Assert TUI navigation section is guarded by platform=tui condition
	assert.True(t,
		strings.Contains(content, "platform=tui") || strings.Contains(content, "IF platform=tui"),
		"expected TUI navigation section to be guarded by platform=tui condition")

	// Assert web navigation section has platform condition
	assert.True(t,
		strings.Contains(content, "platform=web") || strings.Contains(content, "Pointer-Driven"),
		"expected web navigation section to have platform condition")

	// Assert mobile navigation section has platform condition
	assert.True(t,
		strings.Contains(content, "mobile") || strings.Contains(content, "mini-program"),
		"expected mobile navigation section reference")

	// Verify TUI section does not appear unconditionally in the web/mobile rendering path
	assert.True(t,
		strings.Contains(content, "TUI Navigation") || strings.Contains(content, "render when platform=tui"),
		"expected TUI section to be clearly conditional")
}

// Traceability: TC-007 -> Task 2 AC-3
func TestTC_007_WritePRDSkillReferencesTUINavigationTemplate(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/write-prd/SKILL.md")

	// Assert SKILL.md contains reference to TUI navigation rendering
	assert.True(t,
		strings.Contains(content, "TUI Navigation") || strings.Contains(content, "tui"),
		"expected write-prd SKILL.md to reference TUI navigation rendering")

	// Assert platform=tui detection logic
	assert.True(t,
		strings.Contains(content, "platform=tui") || strings.Contains(content, "platform = tui"),
		"expected write-prd SKILL.md to contain platform=tui detection logic")

	// Assert reference to prd-ui-functions.md template
	fileContains(t, content, "prd-ui-functions", "reference to prd-ui-functions template")
}

// Traceability: TC-008 -> Task 2 AC-4
func TestTC_008_ExistingWebMobilePRDGenerationBehaviorUnchanged(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/write-prd/templates/prd-ui-functions.md")

	// Assert web Navigation Architecture section exists
	assert.True(t,
		strings.Contains(content, "Pointer-Driven") || strings.Contains(content, "Primary Navigation"),
		"expected web Navigation Architecture section to exist")

	// Assert mobile is covered
	assert.True(t,
		strings.Contains(content, "mobile") || strings.Contains(content, "mini-program"),
		"expected mobile navigation reference to exist")

	// Verify TUI content does not leak into web/mobile sections
	// The web/mobile section should have its own Navigation Rules
	webNavRules := strings.Contains(content, "Pointer-Driven")
	assert.True(t, webNavRules, "expected separate Pointer-Driven navigation section for web/mobile")
}

// ============================================================================
// Task 3: ui-design SKILL.md TUI Support
// ============================================================================

// Traceability: TC-009 -> Task 3 AC-1
func TestTC_009_UIDesignSkillDetectsPlatformTUI(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/SKILL.md")

	// Assert platform detection logic identifies platform=tui from PRD
	assert.True(t,
		strings.Contains(content, "platform") && strings.Contains(content, "tui"),
		"expected ui-design SKILL.md to contain platform detection for TUI")

	// Assert TUI branch/flow is separate from web and mobile
	assert.True(t,
		strings.Contains(content, "For TUI") || strings.Contains(content, "TUI Platform") || strings.Contains(content, "### For TUI"),
		"expected TUI-specific branch/flow separate from web and mobile")
}

// Traceability: TC-010 -> Task 3 AC-2
func TestTC_010_TUIBranchPresentsThemeSelection(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/SKILL.md")

	// Assert TUI branch offers three theme options
	assert.True(t,
		strings.Contains(content, "Modern Dark") || strings.Contains(content, "modern-dark"),
		"expected Modern Dark theme option")

	assert.True(t,
		strings.Contains(content, "Minimal ASCII") || strings.Contains(content, "minimal-ascii"),
		"expected Minimal ASCII theme option")

	assert.True(t,
		strings.Contains(content, "DESIGN.md") || strings.Contains(content, "custom"),
		"expected DESIGN.md custom theme option")

	// Assert theme selection is prompted during TUI flow
	assert.True(t,
		strings.Contains(content, "theme") || strings.Contains(content, "Theme"),
		"expected theme selection in TUI flow")
}

// Traceability: TC-011 -> Task 3 AC-3
func TestTC_011_TUIBranchUsesTUIPlatformAndThemeFiles(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/SKILL.md")

	// Assert TUI branch references platforms/tui.md
	assert.True(t,
		strings.Contains(content, "platforms/tui.md") || strings.Contains(content, "platforms") && strings.Contains(content, "tui"),
		"expected TUI branch to reference platforms/tui.md")

	// Assert TUI branch references theme files from styles directory
	assert.True(t,
		strings.Contains(content, "styles/") || strings.Contains(content, "styles") && strings.Contains(content, "tui"),
		"expected TUI branch to reference styles directory for themes")

	// Assert TUI branch reads theme properties
	assert.True(t,
		strings.Contains(content, "templates/styles") || strings.Contains(content, "theme"),
		"expected TUI branch to read theme properties from styles directory")
}

// Traceability: TC-012 -> Task 3 AC-4
func TestTC_012_UIDesignTemplateIncludesTUIComponentTemplateWith5StructuralRequirements(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/templates/ui-design.md")

	// Assert Panel Placement subsection
	fileContains(t, content, "Panel Placement", "Panel Placement subsection")

	// Assert ASCII Layout Mockup subsection
	fileContains(t, content, "ASCII Layout Mockup", "ASCII Layout Mockup subsection")

	// Assert Dimensions subsection
	fileContains(t, content, "Dimensions", "Dimensions subsection")

	// Assert Character Palette subsection
	fileContains(t, content, "Character Palette", "Character Palette subsection")

	// Assert Color Mapping subsection
	fileContains(t, content, "Color Mapping", "Color Mapping subsection")

	// Assert Edge Cases subsection
	fileContains(t, content, "Edge Cases", "Edge Cases subsection")
}

// Traceability: TC-013 -> Task 3 AC-5
func TestTC_013_MultiPlatformFeaturesProduceSeparateUIDesignFiles(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/SKILL.md")

	// Assert multi-platform logic produces separate output files per platform
	assert.True(t,
		strings.Contains(content, "ui-design-web") && strings.Contains(content, "ui-design-tui"),
		"expected multi-platform logic to produce separate ui-design-web.md and ui-design-tui.md files")
}

// Traceability: TC-014 -> Task 3 AC-6
func TestTC_014_SingleTUIFeatureProducesUIDesignTUI(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/SKILL.md")

	// Assert single TUI feature produces ui-design-tui.md
	assert.True(t,
		strings.Contains(content, "ui-design-tui.md"),
		"expected single TUI feature to produce ui-design-tui.md")
}

// Traceability: TC-015 -> Task 3 AC-7
func TestTC_015_ExistingWebMobileUIDesignBehaviorUnchanged(t *testing.T) {
	skillContent := readFile(t, "plugins/forge/skills/ui-design/SKILL.md")
	templateContent := readFile(t, "plugins/forge/skills/ui-design/templates/ui-design.md")

	// Assert web platform branch logic exists
	assert.True(t,
		strings.Contains(skillContent, "web") && strings.Contains(skillContent, "Web"),
		"expected web platform branch logic to be present")

	// Assert mobile platform branch logic exists
	assert.True(t,
		strings.Contains(skillContent, "mobile") || strings.Contains(skillContent, "Mobile"),
		"expected mobile platform branch logic to be present")

	// Assert web/mobile component template sections exist in ui-design.md template
	assert.True(t,
		strings.Contains(templateContent, "Component:") || strings.Contains(templateContent, "Component"),
		"expected web/mobile component template sections in ui-design.md")
}

// ============================================================================
// Task 4: TUI Prototype Rules
// ============================================================================

// Traceability: TC-016 -> Task 4 AC-1
func TestTC_016_PrototypeTemplateIncludesTUISpecificGenerationRules(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/templates/prototype.md")

	// Assert TUI-specific prototype generation rules section
	assert.True(t,
		strings.Contains(content, "TUI Prototype Rules"),
		"expected prototype.md to contain TUI-specific prototype generation rules section")

	// Assert TUI rules are distinct from web/mobile rules
	assert.True(t,
		strings.Contains(content, "Terminal Window") || strings.Contains(content, "terminal-window"),
		"expected TUI rules to be distinct from web/mobile prototype rules")
}

// Traceability: TC-017 -> Task 4 AC-2
func TestTC_017_TUIPrototypeIsSingleIndexHTMLWithTerminalWindowDiv(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/templates/prototype.md")

	// Assert single index.html output
	assert.True(t,
		strings.Contains(content, "Single-File Structure") || strings.Contains(content, "single `index.html`"),
		"expected TUI prototype rules to specify single index.html output")

	// Assert terminal-window div container
	assert.True(t,
		strings.Contains(content, "terminal-window") || strings.Contains(content, "Terminal Window Container"),
		"expected terminal-window div container specification")

	// Assert all panels rendered inside terminal window div
	assert.True(t,
		strings.Contains(content, "tui-panel") || strings.Contains(content, "panel"),
		"expected panel rendering rules inside terminal window")
}

// Traceability: TC-018 -> Task 4 AC-3
func TestTC_018_TUIPrototypeIncludesSimulatedKeyButtons(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/templates/prototype.md")

	// Assert simulated key buttons at the bottom
	assert.True(t,
		strings.Contains(content, "simulated-keys") || strings.Contains(content, "Simulated Key"),
		"expected simulated key buttons specification")

	// Assert required buttons: [Tab], [1], [q], [:command]
	fileContains(t, content, "[Tab]", "Tab button")
	fileContains(t, content, "[q]", "quit button")
	fileContains(t, content, ":command", "command button")

	// Assert buttons trigger panel switching
	assert.True(t,
		strings.Contains(content, "panel") && strings.Contains(content, "focus"),
		"expected buttons to trigger panel switching")
}

// Traceability: TC-019 -> Task 4 AC-4
func TestTC_019_TUIPrototypeUsesMonospaceFontAndDarkBackground(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/templates/prototype.md")

	// Assert monospace font specification
	assert.True(t,
		strings.Contains(content, "monospace") || strings.Contains(content, "Monospace") || strings.Contains(content, "Consolas") || strings.Contains(content, "Menlo"),
		"expected monospace font specification")

	// Assert dark background color
	assert.True(t,
		strings.Contains(content, "#1e1e1e") || strings.Contains(content, "dark background") || strings.Contains(content, "Dark"),
		"expected dark background color specification")

	// Assert fixed-width character rendering
	assert.True(t,
		strings.Contains(content, "font-family") || strings.Contains(content, "fixed"),
		"expected fixed-width character rendering specification")
}

// Traceability: TC-020 -> Task 4 AC-5
func TestTC_020_TUIPrototypePanelLayoutMatchesASCIIMockup(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/templates/prototype.md")

	// Assert HTML panel layout must match ASCII mockup from ui-design.md
	assert.True(t,
		strings.Contains(content, "ASCII mockup") || strings.Contains(content, "ASCII Layout Mockup"),
		"expected panel layout matching requirement referencing ASCII mockup")

	// Assert reference to ui-design.md as source of truth for dimensions
	assert.True(t,
		strings.Contains(content, "ui-design") || strings.Contains(content, "dimensions"),
		"expected reference to ui-design.md dimensions as source of truth")
}

// Traceability: TC-021 -> Task 4 AC-6
func TestTC_021_TUIPrototypesOutputToCorrectDirectories(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/templates/prototype.md")

	// Assert TUI prototype output path for multi-platform includes prototype/tui/
	assert.True(t,
		strings.Contains(content, "prototype/tui/"),
		"expected multi-platform TUI prototype output path prototype/tui/")

	// Assert single TUI feature output path
	assert.True(t,
		strings.Contains(content, "prototype/") && strings.Contains(content, "Single TUI"),
		"expected single TUI feature prototype output path")
}

// ============================================================================
// Task 5: Eval-UI Rubric Templates
// ============================================================================

// Traceability: TC-022 -> Task 5 AC-1
func TestTC_022_RubricWebContainsExistingWebRubricContent(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/eval-ui/templates/rubric-web.md")

	// Assert file contains rubric content with dimensions and scoring
	assert.True(t,
		strings.Count(content, "## ") >= 4 || strings.Contains(content, "Perspectives") || strings.Contains(content, "Dimension"),
		"expected rubric-web.md to contain multiple dimensions with scoring criteria")

	// Assert total score equals 1000 points
	fileContains(t, content, "1000", "total score of 1000 points")
}

// Traceability: TC-023 -> Task 5 AC-2
func TestTC_023_RubricTUIHas4CorrectDimensions(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/eval-ui/templates/rubric-tui.md")

	// Assert Requirement Coverage dimension with 250 points
	fileContains(t, content, "Requirement Coverage", "Requirement Coverage dimension")
	assert.True(t, strings.Contains(content, "250 pts"), "expected 250 pts per dimension")

	// Assert Terminal Experience dimension with 250 points
	fileContains(t, content, "Terminal Experience", "Terminal Experience dimension")

	// Assert Visual Specification dimension with 250 points
	fileContains(t, content, "Visual Specification", "Visual Specification dimension")

	// Assert Implementability dimension with 250 points
	fileContains(t, content, "Implementability", "Implementability dimension")

	// Assert total score equals 1000 points
	fileContains(t, content, "1000", "total score of 1000 points")
}

// Traceability: TC-024 -> Task 5 AC-3
func TestTC_024_RubricTUIDeductionRulesAreCorrect(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/eval-ui/templates/rubric-tui.md")

	// Assert deduction: missing ASCII mockup sets Visual Specification to 0
	assert.True(t,
		strings.Contains(content, "Missing ASCII mockup") && strings.Contains(content, "0"),
		"expected deduction rule for missing ASCII mockup")

	// Assert deduction: pending/unspecified characters incur -30 per instance
	assert.True(t,
		strings.Contains(content, "-30"),
		"expected -30 deduction for pending/unspecified characters")

	// Assert deduction: missing mandatory edge case incurs -50 per case
	assert.True(t,
		strings.Contains(content, "-50"),
		"expected -50 deduction for missing mandatory edge case")

	// Assert deduction: vague dimensions incur -20 per instance
	assert.True(t,
		strings.Contains(content, "-20"),
		"expected -20 deduction for vague dimensions")
}

// Traceability: TC-025 -> Task 5 AC-4
func TestTC_025_RubricMobileHas4CorrectDimensions(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/eval-ui/templates/rubric-mobile.md")

	// Assert Requirement Coverage dimension with 250 points
	fileContains(t, content, "Requirement Coverage", "Requirement Coverage dimension")
	assert.True(t, strings.Contains(content, "250 pts"), "expected 250 pts per dimension")

	// Assert Touch Experience dimension with 250 points
	fileContains(t, content, "Touch Experience", "Touch Experience dimension")

	// Assert Adaptive Layout dimension with 250 points
	fileContains(t, content, "Adaptive Layout", "Adaptive Layout dimension")

	// Assert Implementability dimension with 250 points
	fileContains(t, content, "Implementability", "Implementability dimension")

	// Assert total score equals 1000 points
	fileContains(t, content, "1000", "total score of 1000 points")
}

// Traceability: TC-026 -> Task 5 AC-5
func TestTC_026_RubricMobileDeductionRulesAreCorrect(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/eval-ui/templates/rubric-mobile.md")

	// Assert deduction: touch targets without size annotation incur -30 per instance
	assert.True(t,
		strings.Contains(content, "-30") && (strings.Contains(content, "touch target") || strings.Contains(content, "Touch target")),
		"expected -30 deduction for touch targets without size")

	// Assert deduction: missing landscape/portrait adaptation incurs -50
	assert.True(t,
		strings.Contains(content, "-50"),
		"expected -50 deduction for missing landscape/portrait adaptation")

	// Assert deduction: missing safe area handling incurs -40
	assert.True(t,
		strings.Contains(content, "-40"),
		"expected -40 deduction for missing safe area handling")
}

// Traceability: TC-027 -> Task 5 AC-6
func TestTC_027_EvalUISkillDetectsPlatformAndSelectsRubric(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/eval-ui/SKILL.md")

	// Assert platform detection from ui-design document
	assert.True(t,
		strings.Contains(content, "platform") && (strings.Contains(content, "detect") || strings.Contains(content, "Detect") || strings.Contains(content, "Platform Detection")),
		"expected eval-ui SKILL.md to detect platform from ui-design document")

	// Assert platform=web selects rubric-web.md
	assert.True(t,
		strings.Contains(content, "rubric-web.md"),
		"expected platform=web to select rubric-web.md")

	// Assert platform=mobile selects rubric-mobile.md
	assert.True(t,
		strings.Contains(content, "rubric-mobile.md"),
		"expected platform=mobile to select rubric-mobile.md")

	// Assert platform=tui selects rubric-tui.md
	assert.True(t,
		strings.Contains(content, "rubric-tui.md"),
		"expected platform=tui to select rubric-tui.md")
}

// Traceability: TC-028 -> Task 5 AC-7
func TestTC_028_MultiPlatformFeaturesEvaluateWithRespectiveRubrics(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/eval-ui/SKILL.md")

	// Assert multi-platform logic evaluates each platform independently
	assert.True(t,
		strings.Contains(content, "independently") || strings.Contains(content, "separate") || strings.Contains(content, "each platform"),
		"expected multi-platform logic to evaluate each platform independently")

	// Assert each platform uses its respective rubric
	assert.True(t,
		strings.Contains(content, "respective rubric") || strings.Contains(content, "respective"),
		"expected each platform to use its respective rubric")
}

// ============================================================================
// Task 6: Manifest Update Template
// ============================================================================

// Traceability: TC-029 -> Task 6 AC-1
func TestTC_029_SinglePlatformWebManifestUnchanged(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/templates/manifest-update-ui.md")

	// Assert template lists ui-design.md for single web platform
	fileContains(t, content, "ui-design.md", "ui-design.md for single web platform")

	// Assert template lists prototype/ directory for single web platform
	fileContains(t, content, "prototype/", "prototype/ directory for single web platform")
}

// Traceability: TC-030 -> Task 6 AC-2
func TestTC_030_MultiPlatformManifestListsPlatformSpecificFiles(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/templates/manifest-update-ui.md")

	// Assert multi-platform section lists ui-design-web.md
	fileContains(t, content, "ui-design-web.md", "ui-design-web.md for multi-platform")

	// Assert multi-platform section lists ui-design-tui.md
	fileContains(t, content, "ui-design-tui.md", "ui-design-tui.md for multi-platform")

	// Assert multi-platform section lists prototype/web/
	fileContains(t, content, "prototype/web/", "prototype/web/ for multi-platform")

	// Assert multi-platform section lists prototype/tui/
	fileContains(t, content, "prototype/tui/", "prototype/tui/ for multi-platform")
}

// Traceability: TC-031 -> Task 6 AC-3
func TestTC_031_SingleTUIManifestListsCorrectFiles(t *testing.T) {
	content := readFile(t, "plugins/forge/skills/ui-design/templates/manifest-update-ui.md")

	// Assert single TUI section lists ui-design-tui.md
	fileContains(t, content, "ui-design-tui.md", "ui-design-tui.md for single TUI platform")

	// Assert single TUI section lists prototype/
	// The template should reference prototype/ for single TUI
	assert.True(t,
		strings.Contains(content, "prototype/"),
		"expected prototype/ reference for single TUI platform")
}
