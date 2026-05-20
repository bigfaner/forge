package docsync

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// extractDesignMdPath returns the path to the extract-design-md skill file.
func extractDesignMdPath(t *testing.T) string {
	t.Helper()
	root := projectRoot(t)
	p := filepath.Join(root, "..", "plugins", "forge", "skills", "extract-design-md", "SKILL.md")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		t.Fatalf("skill file not found: %s", p)
	}
	return p
}

// readExtractDesignMd reads the full command file content.
func readExtractDesignMd(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(extractDesignMdPath(t))
	if err != nil {
		t.Fatalf("cannot read extract-design-md.md: %v", err)
	}
	return string(data)
}

// extractYamlFrontmatter extracts the YAML frontmatter block from markdown content.
func extractYamlFrontmatter(content string) string {
	re := regexp.MustCompile("(?s)^---\n(.*?)\n---")
	m := re.FindStringSubmatch(content)
	if m == nil {
		return ""
	}
	return m[1]
}

// --- Test: Command frontmatter description mentions all three platforms ---

func TestExtractDesignMd_DescriptionMentionsAllPlatforms(t *testing.T) {
	content := readExtractDesignMd(t)
	fm := extractYamlFrontmatter(content)

	// The description field must mention all three platforms
	requiredKeywords := []string{"web", "mobile", "tui"}
	for _, kw := range requiredKeywords {
		if !strings.Contains(strings.ToLower(fm), kw) {
			t.Errorf("frontmatter description missing platform keyword %q", kw)
		}
	}
}

// --- Test: allowed-tools includes image analysis capability ---

func TestExtractDesignMd_AllowedToolsIncludesImageCapability(t *testing.T) {
	content := readExtractDesignMd(t)
	fm := extractYamlFrontmatter(content)

	// Find allowed-tools line
	if !strings.Contains(fm, "allowed-tools:") {
		t.Fatal("frontmatter missing allowed-tools field")
	}

	// Must include a tool capable of reading images (Read or an image analysis tool)
	hasImageCapability := strings.Contains(fm, "Read") ||
		strings.Contains(fm, "mcp__*analyze_image*") ||
		strings.Contains(fm, "mcp__*")

	if !hasImageCapability {
		// At minimum, Read must be present (it supports image files)
		t.Error("allowed-tools must include Read or image analysis tool for TUI screenshot support")
	}

	// Must also include WebFetch (web extraction needs it)
	if !strings.Contains(fm, "WebFetch") {
		t.Error("allowed-tools must include WebFetch for web extraction")
	}
}

// --- Test: argument-hints includes --platform with valid values ---

func TestExtractDesignMd_ArgumentHintsIncludesPlatform(t *testing.T) {
	content := readExtractDesignMd(t)
	fm := extractYamlFrontmatter(content)

	// Must have argument-hint with a platform entry
	if !strings.Contains(fm, "argument-hint") && !strings.Contains(fm, "argument_hint") {
		t.Fatal("frontmatter missing argument-hint")
	}

	// Must have a platform argument hint
	if !strings.Contains(fm, "platform") {
		t.Error("argument-hints missing platform argument")
	}

	// Must mention valid values: web, mobile, tui
	requiredValues := []string{"web", "mobile", "tui"}
	for _, v := range requiredValues {
		if !strings.Contains(strings.ToLower(fm), v) {
			t.Errorf("argument-hints for platform missing value %q", v)
		}
	}
}

// --- Test: Platform routing logic exists in command body ---

func TestExtractDesignMd_PlatformRoutingLogicExists(t *testing.T) {
	content := readExtractDesignMd(t)

	// After frontmatter, the command body must have routing logic
	// that handles at least "web", "mobile", and "tui" platforms
	bodyAfterFM := content
	if idx := strings.Index(content, "---\n"); idx >= 0 {
		secondDash := strings.Index(content[idx+4:], "---\n")
		if secondDash >= 0 {
			bodyAfterFM = content[idx+4+secondDash+4:]
		}
	}

	bodyLower := strings.ToLower(bodyAfterFM)

	// Must mention platform routing or conditional
	hasRouting := strings.Contains(bodyLower, "platform") &&
		(strings.Contains(bodyLower, "mobile") || strings.Contains(bodyLower, "tui"))

	if !hasRouting {
		t.Error("command body missing platform routing logic")
	}
}

// --- Test: Input validation rejects invalid platform values ---

func TestExtractDesignMd_InputValidationRejectsInvalidPlatforms(t *testing.T) {
	content := readExtractDesignMd(t)

	// The command must include validation logic for platform values
	// Look for mentions of error/validation related to platform
	bodyLower := strings.ToLower(content)

	// Must have some form of validation message for invalid platform
	hasValidation := strings.Contains(bodyLower, "invalid") ||
		strings.Contains(bodyLower, "unsupported") ||
		strings.Contains(bodyLower, "unknown platform") ||
		strings.Contains(bodyLower, "must be one of") ||
		(strings.Contains(bodyLower, "error") && strings.Contains(bodyLower, "platform"))

	if !hasValidation {
		t.Error("command body missing input validation for invalid platform values")
	}
}

// --- Test: Web extraction unchanged when --platform web or no flag ---

func TestExtractDesignMd_WebExtractionStepsPreserved(t *testing.T) {
	content := readExtractDesignMd(t)

	// All original web extraction steps must still be present
	// These are the core web extraction elements from the original command
	requiredWebElements := []string{
		"Color palette",
		"Typography",
		"Components",
		"Layer 1",
		"Layer 2",
		"CSS custom properties",
		"WebFetch",
		"DESIGN.md",
	}

	for _, elem := range requiredWebElements {
		if !strings.Contains(content, elem) {
			t.Errorf("web extraction element %q missing from command (web behavior must be preserved)", elem)
		}
	}
}

// --- Mobile Adapter Tests ---

// --- Test: Mobile adapter has explicit extraction steps ---

func TestExtractDesignMd_MobileAdapterHasExplicitExtractionSteps(t *testing.T) {
	content := readExtractDesignMd(t)
	bodyLower := strings.ToLower(content)

	// Mobile adapter must not be a "placeholder" anymore
	if strings.Contains(bodyLower, "mobile placeholder") {
		t.Error("mobile adapter still contains placeholder text — must be replaced with actual extraction logic")
	}

	// Must mention mobile-specific extraction concepts
	mobileConcepts := []string{
		"responsive breakpoint",
		"touch target",
		"safe area",
	}
	for _, concept := range mobileConcepts {
		if !strings.Contains(bodyLower, concept) {
			t.Errorf("mobile adapter missing required concept: %q", concept)
		}
	}
}

// --- Test: Mobile adapter reuses web extraction layers ---

func TestExtractDesignMd_MobileAdapterReusesWebExtraction(t *testing.T) {
	content := readExtractDesignMd(t)
	bodyLower := strings.ToLower(content)

	// Must reference reuse of web extraction pipeline (Layers 1-5)
	hasLayerReuse := strings.Contains(bodyLower, "reuse") ||
		(strings.Contains(bodyLower, "layer 1") && strings.Contains(bodyLower, "mobile"))
	if !hasLayerReuse {
		t.Error("mobile adapter does not explicitly reference reusing web extraction layers")
	}

	// Must reference mobile user-agent or mobile viewport context
	hasMobileContext := strings.Contains(bodyLower, "user-agent") ||
		strings.Contains(bodyLower, "viewport")
	if !hasMobileContext {
		t.Error("mobile adapter missing mobile User-Agent or viewport context reference")
	}
}

// --- Test: Mobile adapter documents responsive breakpoint analysis ---

func TestExtractDesignMd_MobileAdapterDocumentsBreakpoints(t *testing.T) {
	content := readExtractDesignMd(t)

	// Must reference common mobile breakpoints
	requiredBreakpoints := []string{"320", "375", "414", "768"}
	for _, bp := range requiredBreakpoints {
		if !strings.Contains(content, bp) {
			t.Errorf("mobile adapter missing common breakpoint %spx", bp)
		}
	}
}

// --- Test: Mobile adapter documents touch target guideline ---

func TestExtractDesignMd_MobileAdapterDocumentsTouchTargetGuideline(t *testing.T) {
	content := readExtractDesignMd(t)

	// Must reference the 44x44pt touch target minimum guideline
	if !strings.Contains(content, "44") {
		t.Error("mobile adapter missing 44x44pt touch target minimum guideline reference")
	}
}

// --- Test: Mobile adapter documents safe area handling ---

func TestExtractDesignMd_MobileAdapterDocumentsSafeAreaHandling(t *testing.T) {
	content := readExtractDesignMd(t)
	bodyLower := strings.ToLower(content)

	// Must reference safe-area-inset CSS env or viewport meta
	hasSafeAreaRef := strings.Contains(content, "safe-area-inset") ||
		strings.Contains(bodyLower, "viewport meta")
	if !hasSafeAreaRef {
		t.Error("mobile adapter missing safe-area-inset CSS env or viewport meta reference")
	}
}

// --- Test: Mobile DESIGN.md template extends web template ---

func TestExtractDesignMd_MobileTemplateExtendsWebTemplate(t *testing.T) {
	// Read the mobile template file directly (the skill references it via template path)
	root := projectRoot(t)
	mobileTemplatePath := filepath.Join(root, "..", "plugins", "forge", "skills", "extract-design-md", "templates", "design-mobile.md")
	data, err := os.ReadFile(mobileTemplatePath)
	if err != nil {
		t.Fatalf("cannot read mobile template: %v", err)
	}
	content := string(data)

	// Mobile DESIGN.md must extend the web template with mobile-specific sections
	// Must have mobile-specific sections in the output template
	requiredMobileSections := []string{
		"Touch Target",
		"Safe Area",
	}
	for _, section := range requiredMobileSections {
		if !strings.Contains(content, section) {
			t.Errorf("mobile DESIGN.md template missing required section: %q", section)
		}
	}
}

// --- Test: Mobile adapter marks estimated values ---

func TestExtractDesignMd_MobileAdapterMarksEstimatedValues(t *testing.T) {
	content := readExtractDesignMd(t)

	// Must have a mechanism for marking estimated values
	if !strings.Contains(content, "(estimated)") {
		t.Error("mobile adapter missing (estimated) marker for uncertain values")
	}
}

// --- Test: Mobile adapter documents limitation about responsive CSS ---

func TestExtractDesignMd_MobileAdapterDocumentsResponsiveCSSLimitation(t *testing.T) {
	content := readExtractDesignMd(t)
	bodyLower := strings.ToLower(content)

	// Must document that mobile extraction depends on target URL serving responsive CSS
	hasLimitation := strings.Contains(bodyLower, "responsive css") ||
		strings.Contains(bodyLower, "responsive stylesheet") ||
		(strings.Contains(bodyLower, "responsive") && strings.Contains(bodyLower, "limitation"))
	if !hasLimitation {
		t.Error("mobile adapter missing documentation about responsive CSS dependency limitation")
	}
}

// --- Test: Mobile extraction table entry updated from placeholder ---

func TestExtractDesignMd_MobileExtractionTableNotPlaceholder(t *testing.T) {
	content := readExtractDesignMd(t)

	// The platform routing table should NOT have "placeholder" in the mobile row
	routingTableRegex := regexp.MustCompile(`(?i)\|.*mobile.*\|.*\|.*\|`)
	matches := routingTableRegex.FindAllString(content, -1)
	for _, match := range matches {
		if strings.Contains(strings.ToLower(match), "placeholder") {
			t.Errorf("mobile routing table entry still contains 'placeholder': %s", strings.TrimSpace(match))
		}
	}
}

// --- TUI Adapter Tests ---

// --- Test: TUI adapter is not a placeholder ---

func TestExtractDesignMd_TuiAdapterNotPlaceholder(t *testing.T) {
	content := readExtractDesignMd(t)
	bodyLower := strings.ToLower(content)

	// TUI adapter must not be a "placeholder" or "not yet implemented" anymore
	if strings.Contains(bodyLower, "tui placeholder") {
		t.Error("TUI adapter still contains 'TUI placeholder' text — must be replaced with actual extraction logic")
	}
	if strings.Contains(bodyLower, "not yet implemented") && strings.Contains(bodyLower, "tui") {
		t.Error("TUI adapter still contains 'not yet implemented' message — must be replaced with actual extraction logic")
	}
}

// --- Test: TUI adapter requires local file path (not URL) ---

func TestExtractDesignMd_TuiRequiresLocalFilePath(t *testing.T) {
	content := readExtractDesignMd(t)
	bodyLower := strings.ToLower(content)

	// Must explicitly state that TUI input must be a local file path
	hasLocalPath := strings.Contains(bodyLower, "local file path") ||
		strings.Contains(bodyLower, "local file") ||
		(strings.Contains(bodyLower, "file path") && !strings.Contains(bodyLower, "url"))
	if !hasLocalPath {
		t.Error("TUI adapter missing explicit requirement that input must be a local file path (not URL)")
	}

	// Must NOT support URL input for TUI
	// Find the TUI section and check it does NOT reference URL fetching
	tuiSectionRegex := regexp.MustCompile(`(?i)(tui|terminal).*?(?:mobile|$)`)
	tuiMatch := tuiSectionRegex.FindString(content)
	if tuiMatch != "" {
		// Within the TUI context, it should NOT mention URL as the primary input
		// (it's okay if the general command mentions URL for web/mobile)
		if strings.Contains(strings.ToLower(tuiMatch), "fetch") && strings.Contains(strings.ToLower(tuiMatch), "url") {
			t.Error("TUI adapter incorrectly suggests URL fetching — TUI must use local file path only")
		}
	}
}

// --- Test: TUI adapter uses AI vision for analysis ---

func TestExtractDesignMd_TuiUsesAIVision(t *testing.T) {
	content := readExtractDesignMd(t)
	bodyLower := strings.ToLower(content)

	// Must reference AI vision or image analysis for TUI screenshot
	hasVisionRef := strings.Contains(bodyLower, "vision") ||
		strings.Contains(bodyLower, "image analysis") ||
		strings.Contains(bodyLower, "screenshot analysis") ||
		strings.Contains(bodyLower, "analyze the screenshot") ||
		strings.Contains(bodyLower, "read the image") ||
		strings.Contains(bodyLower, "visual analysis")
	if !hasVisionRef {
		t.Error("TUI adapter missing reference to AI vision / image analysis for screenshot processing")
	}
}

// --- Test: TUI adapter extracts required token categories ---

func TestExtractDesignMd_TuiExtractsRequiredTokenCategories(t *testing.T) {
	content := readExtractDesignMd(t)
	bodyLower := strings.ToLower(content)

	// TUI must extract: ANSI color palette, character set, panel layout dimensions, key bindings
	requiredCategories := []struct {
		keyword  string
		category string
	}{
		{"ansi", "ANSI color palette"},
		{"color palette", "color palette"},
		{"character set", "character set"},
		{"panel layout", "panel layout dimensions"},
		{"key binding", "key bindings"},
	}
	for _, rc := range requiredCategories {
		if !strings.Contains(bodyLower, rc.keyword) {
			t.Errorf("TUI adapter missing required extraction category: %q (keyword: %q)", rc.category, rc.keyword)
		}
	}
}

// --- Test: TUI adapter references xterm-256 color numbers ---

func TestExtractDesignMd_TuiReferencesXterm256(t *testing.T) {
	content := readExtractDesignMd(t)

	// Must reference xterm-256 color numbers
	if !strings.Contains(strings.ToLower(content), "xterm-256") &&
		!strings.Contains(strings.ToLower(content), "xterm 256") &&
		!strings.Contains(content, "xterm-256") {
		t.Error("TUI adapter missing reference to xterm-256 color numbers for ANSI color palette")
	}
}

// --- Test: TUI adapter marks all values as estimated ---

func TestExtractDesignMd_TuiMarksAllValuesAsEstimated(t *testing.T) {
	content := readExtractDesignMd(t)

	// All TUI values MUST be marked (estimated)
	// Find TUI-specific sections and verify (estimated) marking instruction exists
	if !strings.Contains(content, "(estimated)") {
		t.Error("TUI adapter missing (estimated) marker requirement — all values must be marked estimated")
	}

	// There must be an explicit instruction about marking ALL TUI values as estimated
	bodyLower := strings.ToLower(content)
	hasEstimateInstruction := strings.Contains(bodyLower, "all") && strings.Contains(content, "(estimated)")
	if !hasEstimateInstruction {
		t.Error("TUI adapter missing instruction that ALL values must be marked (estimated)")
	}
}

// --- Test: TUI output structure aligns with modern-dark-tui sections ---

func TestExtractDesignMd_TuiOutputAlignsWithModernDarkTui(t *testing.T) {
	// Read the TUI template file directly (the skill references it via template path)
	root := projectRoot(t)
	tuiTemplatePath := filepath.Join(root, "..", "plugins", "forge", "skills", "extract-design-md", "templates", "design-tui.md")
	data, err := os.ReadFile(tuiTemplatePath)
	if err != nil {
		t.Fatalf("cannot read TUI template: %v", err)
	}
	content := string(data)

	// TUI output must have sections matching modern-dark-tui.md structure:
	// Color Space, Character Set, Character Palette Reference, Color Palette, Panel Layout
	requiredTuiSections := []string{
		"Color Space",
		"Character Set",
		"Character Palette Reference",
		"Color Palette",
		"Panel Layout",
	}
	for _, section := range requiredTuiSections {
		if !strings.Contains(content, section) {
			t.Errorf("TUI output template missing section matching modern-dark-tui: %q", section)
		}
	}
}

// --- Test: TUI match strategy supports built-in themes ---

func TestExtractDesignMd_TuiMatchStrategySupportsBuiltInThemes(t *testing.T) {
	content := readExtractDesignMd(t)
	bodyLower := strings.ToLower(content)

	// TUI must support matching against built-in TUI themes
	hasModernDark := strings.Contains(bodyLower, "modern-dark-tui")
	hasMinimalASCII := strings.Contains(bodyLower, "minimal-ascii-tui") ||
		strings.Contains(bodyLower, "minimal-ascii") ||
		strings.Contains(bodyLower, "minimal ascii")

	if !hasModernDark {
		t.Error("TUI match strategy missing reference to modern-dark-tui built-in theme")
	}
	if !hasMinimalASCII {
		t.Error("TUI match strategy missing reference to minimal-ascii-tui built-in theme")
	}
}

// --- Test: TUI adapter rejects blurry or low-resolution screenshots ---

func TestExtractDesignMd_TuiRejectsBlurryScreenshots(t *testing.T) {
	content := readExtractDesignMd(t)
	bodyLower := strings.ToLower(content)

	// Must have guidance or validation for screenshot quality
	hasQualityCheck := strings.Contains(bodyLower, "blurry") ||
		strings.Contains(bodyLower, "low-resolution") ||
		strings.Contains(bodyLower, "low resolution") ||
		strings.Contains(bodyLower, "unreadable") ||
		strings.Contains(bodyLower, "screenshot quality")

	if !hasQualityCheck {
		t.Error("TUI adapter missing screenshot quality validation — must reject blurry or low-resolution screenshots")
	}
}

// --- Test: TUI routing table entry is not placeholder ---

func TestExtractDesignMd_TuiExtractionTableNotPlaceholder(t *testing.T) {
	content := readExtractDesignMd(t)

	// The platform routing table should NOT have "placeholder" in the TUI row
	routingTableRegex := regexp.MustCompile(`(?i)\|.*tui.*\|.*\|.*\|`)
	matches := routingTableRegex.FindAllString(content, -1)
	for _, match := range matches {
		if strings.Contains(strings.ToLower(match), "placeholder") {
			t.Errorf("TUI routing table entry still contains 'placeholder': %s", strings.TrimSpace(match))
		}
	}
}

// --- Test: TUI adapter mentions box-drawing and block elements ---

func TestExtractDesignMd_TuiMentionsBoxDrawingAndBlockElements(t *testing.T) {
	content := readExtractDesignMd(t)
	bodyLower := strings.ToLower(content)

	// TUI must reference box-drawing and block elements as part of character set extraction
	if !strings.Contains(bodyLower, "box-drawing") && !strings.Contains(bodyLower, "box drawing") {
		t.Error("TUI adapter missing reference to box-drawing characters")
	}
	if !strings.Contains(bodyLower, "block element") {
		t.Error("TUI adapter missing reference to block elements")
	}
}

// --- Test: Skill file exists as SKILL.md in skill directory ---

func TestExtractDesignMd_SkillFileExists(t *testing.T) {
	root := projectRoot(t)

	// Check that extract-design-md exists as a skill directory with SKILL.md
	skillDir := filepath.Join(root, "..", "plugins", "forge", "skills", "extract-design-md")
	skillMD := filepath.Join(skillDir, "SKILL.md")

	dirExists := false
	if info, err := os.Stat(skillDir); err == nil && info.IsDir() {
		dirExists = true
	}

	mdExists := false
	if _, err := os.Stat(skillMD); err == nil {
		mdExists = true
	}

	if !dirExists {
		t.Error("extract-design-md skill directory does not exist")
	}
	if !mdExists {
		t.Error("extract-design-md SKILL.md file does not exist in skill directory")
	}
}
