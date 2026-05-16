//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// extract-design-md platform adapters tests -- feature: extract-design-md-platform-adapters
// Tests verify that the extract-design-md command specification supports
// --platform flag with web, mobile, and tui adapters.
//
// Note: extract-design-md is a Claude Code slash command (not a forge CLI
// subcommand). These e2e tests verify the command specification file's
// structure and content to ensure correct behavior at runtime.
// ==============================================================================

// edmProjectRoot returns the forge project root directory.
func edmProjectRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine test file location")
	}
	dir := filepath.Join(filepath.Dir(thisFile), "..", "..")
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("cannot resolve project root: %s", err)
	}
	return abs
}

// edmCommandFilePath returns the path to the extract-design-md command file.
func edmCommandFilePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(edmProjectRoot(t), "plugins", "forge", "commands", "extract-design-md.md")
}

// edmReadCommandFile reads the full command file content.
func edmReadCommandFile(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(edmCommandFilePath(t))
	if err != nil {
		t.Fatalf("cannot read extract-design-md.md: %v", err)
	}
	return string(data)
}

// edmExtractFrontmatter extracts the YAML frontmatter block from markdown content.
func edmExtractFrontmatter(content string) string {
	re := regexp.MustCompile("(?s)^---\n(.*?)\n---")
	m := re.FindStringSubmatch(content)
	if m == nil {
		return ""
	}
	return m[1]
}

// edmCommandBody returns the content after the YAML frontmatter.
func edmCommandBody(content string) string {
	re := regexp.MustCompile("(?s)^---\n.*?\n---\n?")
	loc := re.FindStringIndex(content)
	if loc == nil {
		return content
	}
	return content[loc[1]:]
}

// --- Platform Flag & Scaffolding ---

// Traceability: TC-001 -> Task 1 / AC-2, Proposal Success Criteria item 3
func TestTC_001_DefaultPlatformProducesWebOutput(t *testing.T) {
	content := edmReadCommandFile(t)
	body := edmCommandBody(content)
	bodyLower := strings.ToLower(body)

	// Default platform must be "web" and produce standard web output
	assert.Contains(t, bodyLower, "default to `web`",
		"command must specify default platform as web")

	// Web extraction must produce standard web design tokens
	webSections := []string{
		"Color Palette",
		"Typography",
		"Components",
		"Layout",
		"Depth & Elevation",
	}
	for _, section := range webSections {
		assert.Contains(t, content, section,
			"web output must contain standard section: %s", section)
	}
}

// Traceability: TC-002 -> Task 1 / AC-2
func TestTC_002_ExplicitWebPlatformIdenticalOutput(t *testing.T) {
	content := edmReadCommandFile(t)

	// The command must route "web" platform through the same extraction pipeline
	// as the default (no flag) behavior. Verify web extraction is a single code path.
	assert.Contains(t, content, "web (default)",
		"command must show web as default in platform routing table")

	// Web extraction must route default and explicit --platform web to the same code path
	// The command file references Layer 1 in the web section and again in the mobile
	// section (which says "Reuse all web extraction layers (Layer 1-5)").
	// Both occurrences are expected — mobile reuses web layers, not duplicates them.
	assert.Contains(t, content, "web (default)",
		"routing table must show web as default platform")
	assert.Contains(t, strings.ToLower(content), "reuse all web extraction layers",
		"mobile must explicitly reuse web extraction layers, not duplicate them")
}

// Traceability: TC-003 -> Task 1 / AC-4
func TestTC_003_InvalidPlatformValueRejectedWithClearError(t *testing.T) {
	content := edmReadCommandFile(t)

	// Must have validation logic that rejects invalid platform values
	assert.Contains(t, content, "ERROR: unsupported platform",
		"command must include error message for unsupported platform")

	assert.Contains(t, content, "Must be one of: web, mobile, tui",
		"error message must list valid platforms: web, mobile, tui")

	// Verify no file is written on error - the command must stop immediately
	assert.Contains(t, strings.ToLower(content), "stop immediately",
		"command must stop immediately on invalid platform without writing files")
}

// Traceability: TC-004 -> Task 1 / AC-1
func TestTC_004_CommandFrontmatterDescribesAllThreePlatforms(t *testing.T) {
	content := edmReadCommandFile(t)
	fm := edmExtractFrontmatter(content)

	// Description must mention all three platforms
	platforms := []string{"web", "mobile", "tui"}
	for _, p := range platforms {
		assert.Contains(t, strings.ToLower(fm), p,
			"frontmatter description must mention platform: %s", p)
	}

	// argument-hints must include --platform
	assert.Contains(t, fm, "platform",
		"frontmatter must include platform in argument-hints")

	// Valid platform values must be documented
	for _, p := range platforms {
		assert.Contains(t, strings.ToLower(fm), p,
			"argument-hints must reference platform value: %s", p)
	}

	// allowed_tools must include Read for image analysis (TUI screenshots)
	assert.Contains(t, fm, `"Read"`,
		"allowed_tools must include Read for TUI screenshot analysis")

	// allowed_tools must include WebFetch for web/mobile extraction
	assert.Contains(t, fm, `"WebFetch"`,
		"allowed_tools must include WebFetch for web/mobile extraction")
}

// --- Mobile Adapter ---

// Traceability: TC-005 -> Task 2 / AC-1
func TestTC_005_MobilePlatformFetchesWithMobileContext(t *testing.T) {
	content := edmReadCommandFile(t)
	bodyLower := strings.ToLower(content)

	// Mobile must fetch with mobile User-Agent and viewport context
	assert.Contains(t, bodyLower, "user-agent",
		"mobile extraction must reference mobile User-Agent")

	// Must specify mobile viewport width
	assert.Contains(t, bodyLower, "viewport",
		"mobile extraction must reference viewport configuration")

	// Must reference mobile context in the routing table
	assert.Contains(t, bodyLower, "mobile",
		"command must reference mobile platform")
}

// Traceability: TC-006 -> Task 2 / AC-2
func TestTC_006_ResponsiveBreakpointAnalysisExtractsCommonBreakpoints(t *testing.T) {
	content := edmReadCommandFile(t)

	// Must include responsive breakpoint analysis section
	assert.Contains(t, strings.ToLower(content), "responsive breakpoint",
		"mobile extraction must include responsive breakpoint analysis")

	// Must list common breakpoints: 320px, 375px, 414px, 768px
	breakpoints := []string{"320", "375", "414", "768"}
	for _, bp := range breakpoints {
		assert.Contains(t, content, bp,
			"mobile extraction must reference breakpoint %spx", bp)
	}
}

// Traceability: TC-007 -> Task 2 / AC-3
func TestTC_007_TouchTargetEstimationAnalyzesInteractiveElements(t *testing.T) {
	content := edmReadCommandFile(t)
	bodyLower := strings.ToLower(content)

	// Must include touch target estimation
	assert.Contains(t, bodyLower, "touch target",
		"mobile extraction must include touch target estimation")

	// Must reference 44x44pt minimum guideline
	assert.Contains(t, content, "44",
		"mobile extraction must reference 44x44pt touch target minimum")

	// Must mention interactive elements analysis
	assert.Contains(t, bodyLower, "interactive element",
		"mobile extraction must analyze interactive elements")
}

// Traceability: TC-008 -> Task 2 / AC-4
func TestTC_008_SafeAreaHandlingInference(t *testing.T) {
	content := edmReadCommandFile(t)
	bodyLower := strings.ToLower(content)

	// Must include safe area handling section
	assert.Contains(t, bodyLower, "safe area",
		"mobile extraction must include safe area handling")

	// Must reference CSS env() for safe-area-inset
	assert.Contains(t, content, "safe-area-inset",
		"mobile extraction must reference CSS env(safe-area-inset-*) values")

	// Must mark estimated values when not detected from CSS
	assert.Contains(t, content, "(estimated)",
		"mobile extraction must mark undetected values as (estimated)")
}

// Traceability: TC-009 -> Task 2 / AC-5
func TestTC_009_MobileOutputExtendsWebTemplateWithMobileSections(t *testing.T) {
	content := edmReadCommandFile(t)

	// Must contain all standard web sections
	webSections := []string{
		"Color Palette",
		"Typography",
		"Components",
		"Layout",
		"Depth & Elevation",
		"Responsive Behavior",
		"Signature Patterns",
	}
	for _, section := range webSections {
		assert.Contains(t, content, section,
			"mobile output must preserve web section: %s", section)
	}

	// Must contain mobile-specific sections
	mobileSections := []string{
		"Touch Target",
		"Safe Area",
		"Responsive Breakpoint",
	}
	for _, section := range mobileSections {
		assert.Contains(t, content, section,
			"mobile output must include mobile-specific section: %s", section)
	}

	// The command file has mobile sections in a separate template block with an
	// instruction to "Insert them after 'Responsive Behavior' and before 'Signature Patterns'"
	assert.Contains(t, strings.ToLower(content), "after \"responsive behavior\" and before \"signature patterns\"",
		"command must instruct that mobile sections are inserted between Responsive Behavior and Signature Patterns")
}

// Traceability: TC-010 -> Task 2 / AC-6
func TestTC_010_MobileMatchStrategyWorksWithClosestBuiltInStyle(t *testing.T) {
	content := edmReadCommandFile(t)
	bodyLower := strings.ToLower(content)

	// Must include match strategy for mobile
	assert.Contains(t, bodyLower, "match closest built-in style",
		"mobile must support matching closest built-in style")

	// Must reference built-in styles for matching
	builtinStyles := []string{"vercel", "shadcn", "tailwind ui", "stripe", "apple"}
	foundStyle := false
	for _, style := range builtinStyles {
		if strings.Contains(bodyLower, style) {
			foundStyle = true
			break
		}
	}
	assert.True(t, foundStyle,
		"mobile match strategy must reference at least one built-in style")
}

// Traceability: TC-011 -> Task 2 / AC-7, Proposal Success Criteria item 4
func TestTC_011_MobileOutputConsumableByUiDesign(t *testing.T) {
	content := edmReadCommandFile(t)

	// Mobile DESIGN.md must follow the same structure convention as web DESIGN.md
	// so that ui-design skill can consume it without modification.
	// Verify the output template contains the DESIGN.md header structure.
	assert.Contains(t, content, "# Design System:",
		"mobile output must use standard DESIGN.md header for ui-design compatibility")

	assert.Contains(t, content, "Extracted from:",
		"mobile output must include 'Extracted from' metadata line")

	assert.Contains(t, content, "Based on:",
		"mobile output must include 'Based on' metadata for match strategy")

	// Must contain all token categories that ui-design expects
	uiDesignSections := []string{
		"Color Palette",
		"Typography",
		"Components",
		"Layout",
	}
	for _, section := range uiDesignSections {
		assert.Contains(t, content, section,
			"mobile output must include ui-design-compatible section: %s", section)
	}
}

// --- TUI Adapter ---

// Traceability: TC-012 -> Task 3 / AC-1
func TestTC_012_TuiPlatformAcceptsLocalScreenshotPath(t *testing.T) {
	content := edmReadCommandFile(t)
	bodyLower := strings.ToLower(content)

	// TUI must accept local file path as input
	assert.Contains(t, bodyLower, "local file path",
		"TUI extraction must accept local file path as input")

	// Must reference using Read tool for screenshot loading (backtick-quoted in source)
	assert.True(t, strings.Contains(bodyLower, "read tool") || strings.Contains(bodyLower, "`read` tool"),
		"TUI extraction must reference using Read tool to load screenshots")

	// Must reference screenshot file extension
	assert.Contains(t, bodyLower, ".png",
		"TUI extraction must reference screenshot file format (.png)")
}

// Traceability: TC-013 -> Task 3 / Hard Rule (TUI input must be local file path)
func TestTC_013_TuiRejectsUrlInputWithClearError(t *testing.T) {
	content := edmReadCommandFile(t)

	// Must include error message for URL input on TUI platform
	assert.Contains(t, content, "ERROR: TUI platform requires a local screenshot file path, not a URL",
		"command must include error message for TUI URL input rejection")

	assert.Contains(t, content, "./screenshot.png",
		"error message must provide example local file path")
}

// Traceability: TC-014 -> Task 3 / AC-2
func TestTC_014_TuiAiVisionExtractsAnsiColorPaletteAndCharacterSet(t *testing.T) {
	content := edmReadCommandFile(t)
	bodyLower := strings.ToLower(content)

	// Must reference ANSI color palette extraction
	assert.Contains(t, bodyLower, "ansi color palette",
		"TUI must extract ANSI color palette")

	// Must reference xterm-256 color numbers
	assert.Contains(t, bodyLower, "xterm-256",
		"TUI must reference xterm-256 color numbers")

	// Must reference character set identification
	assert.Contains(t, bodyLower, "character set",
		"TUI must identify character set type")

	// Must mention box-drawing, block elements, or ASCII
	charTypes := []string{"box-drawing", "block element", "ascii"}
	foundCharType := false
	for _, ct := range charTypes {
		if strings.Contains(bodyLower, ct) {
			foundCharType = true
			break
		}
	}
	assert.True(t, foundCharType,
		"TUI character set must mention box-drawing, block elements, or ASCII")

	// Must reference character palette reference table
	assert.Contains(t, content, "Character Palette Reference",
		"TUI output must include Character Palette Reference table")
}

// Traceability: TC-015 -> Task 3 / AC-2
func TestTC_015_TuiExtractsPanelLayoutDimensionsAndKeyBindings(t *testing.T) {
	content := edmReadCommandFile(t)
	bodyLower := strings.ToLower(content)

	// Must reference panel layout dimensions
	assert.Contains(t, bodyLower, "panel layout",
		"TUI must extract panel layout dimensions")

	// Must reference terminal dimensions (rows x columns)
	assert.Contains(t, bodyLower, "rows",
		"TUI must reference terminal rows")
	assert.Contains(t, bodyLower, "column",
		"TUI must reference terminal columns")

	// Must reference key bindings extraction
	assert.Contains(t, bodyLower, "key binding",
		"TUI must extract key bindings from status bar or help panel")

	// Must reference status bar or help legend
	statusRef := strings.Contains(bodyLower, "status bar") ||
		strings.Contains(bodyLower, "help panel") ||
		strings.Contains(bodyLower, "help legend")
	assert.True(t, statusRef,
		"TUI key bindings must reference status bar, help panel, or help legend")
}

// Traceability: TC-016 -> Task 3 / AC-3, Proposal Success Criteria item 5
func TestTC_016_TuiOutputMatchesBuiltinThemeStructure(t *testing.T) {
	content := edmReadCommandFile(t)

	// TUI output sections must match modern-dark-tui structure
	requiredTuiSections := []string{
		"Color Space",
		"Character Set",
		"Character Palette Reference",
		"Color Palette",
		"Typography",
		"Panel Layout",
		"Key Bindings",
	}
	for _, section := range requiredTuiSections {
		assert.Contains(t, content, section,
			"TUI output must include section matching modern-dark-tui: %s", section)
	}

	// Must reference built-in TUI themes for match strategy
	assert.Contains(t, strings.ToLower(content), "modern-dark-tui",
		"TUI must reference modern-dark-tui built-in theme")
}

// Traceability: TC-017 -> Task 3 / AC-4, Task 3 / Hard Rule
func TestTC_017_TuiAllValuesMarkedAsEstimated(t *testing.T) {
	content := edmReadCommandFile(t)
	bodyLower := strings.ToLower(content)

	// Must have explicit instruction to mark ALL TUI values as (estimated)
	assert.Contains(t, content, "(estimated)",
		"TUI output must include (estimated) markers")

	// Must have explicit instruction about ALL values
	// Look for "all values" or "all extracted values" combined with "(estimated)"
	hasAllEstimateInstruction := (strings.Contains(bodyLower, "all values") ||
		strings.Contains(bodyLower, "all extracted values") ||
		strings.Contains(bodyLower, "must be marked")) &&
		strings.Contains(content, "(estimated)")
	assert.True(t, hasAllEstimateInstruction,
		"TUI must have explicit instruction that ALL values must be marked (estimated)")

	// Must not contain guarantees about precision
	// Check that estimated marker appears in TUI-specific context
	tuiSection := ""
	tuiIdx := strings.Index(strings.ToLower(content), "tui extraction")
	if tuiIdx >= 0 {
		// Take a reasonable section of content after TUI extraction heading
		end := tuiIdx + 3000
		if end > len(content) {
			end = len(content)
		}
		tuiSection = content[tuiIdx:end]
	}
	assert.Contains(t, tuiSection, "(estimated)",
		"TUI extraction section must contain (estimated) markers")
}

// Traceability: TC-018 -> Task 3 / AC-6
func TestTC_018_TuiRejectsLowQualityScreenshotWithClearError(t *testing.T) {
	content := edmReadCommandFile(t)
	bodyLower := strings.ToLower(content)

	// Must include screenshot quality validation
	assert.Contains(t, content, "ERROR: Screenshot quality is too low for reliable analysis",
		"command must include error message for low-quality screenshots")

	// Must reference quality assessment criteria
	qualityCriteria := []string{"blurry", "low-resolution", "unreadable"}
	foundCriteria := false
	for _, criterion := range qualityCriteria {
		if strings.Contains(bodyLower, criterion) {
			foundCriteria = true
			break
		}
	}
	assert.True(t, foundCriteria,
		"command must reference quality criteria (blurry, low-resolution, or unreadable)")

	// Must provide tips for better screenshots
	assert.Contains(t, bodyLower, "high-resolution",
		"error message must suggest providing high-resolution screenshot")

	// Must not write any file on quality rejection
	// Verify "stop" appears near the quality check instruction
	qualityCheckIdx := strings.Index(bodyLower, "screenshot quality")
	if qualityCheckIdx < 0 {
		qualityCheckIdx = strings.Index(bodyLower, "blurry")
	}
	if qualityCheckIdx >= 0 {
		nearbyContent := bodyLower[qualityCheckIdx:]
		if len(nearbyContent) > 300 {
			nearbyContent = nearbyContent[:300]
		}
		assert.Contains(t, nearbyContent, "stop",
			"command must stop and not write files when screenshot quality is too low")
	}
}
