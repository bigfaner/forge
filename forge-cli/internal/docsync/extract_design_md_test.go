package docsync

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// extractDesignMdPath returns the path to the extract-design-md command file.
func extractDesignMdPath(t *testing.T) string {
	t.Helper()
	root := projectRoot(t)
	p := filepath.Join(root, "..", "plugins", "forge", "commands", "extract-design-md.md")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		t.Fatalf("command file not found: %s", p)
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

// --- Test: allowed_tools includes image analysis capability ---

func TestExtractDesignMd_AllowedToolsIncludesImageCapability(t *testing.T) {
	content := readExtractDesignMd(t)
	fm := extractYamlFrontmatter(content)

	// Find allowed_tools line
	if !strings.Contains(fm, "allowed_tools") {
		t.Fatal("frontmatter missing allowed_tools field")
	}

	// Must include a tool capable of reading images (Read or an image analysis tool)
	hasImageCapability := strings.Contains(fm, `"Read"`) ||
		strings.Contains(fm, `"mcp__*analyze_image*"`) ||
		strings.Contains(fm, `"mcp__*`)

	if !hasImageCapability {
		// At minimum, Read must be present (it supports image files)
		t.Error("allowed_tools must include Read or image analysis tool for TUI screenshot support")
	}

	// Must also include WebFetch (web extraction needs it)
	if !strings.Contains(fm, `"WebFetch"`) {
		t.Error("allowed_tools must include WebFetch for web extraction")
	}
}

// --- Test: argument-hints includes --platform with valid values ---

func TestExtractDesignMd_ArgumentHintsIncludesPlatform(t *testing.T) {
	content := readExtractDesignMd(t)
	fm := extractYamlFrontmatter(content)

	// Must have argument-hints with a platform entry
	if !strings.Contains(fm, "argument-hints") && !strings.Contains(fm, "argument_hints") {
		t.Fatal("frontmatter missing argument-hints")
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
	content := readExtractDesignMd(t)

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

// --- Test: Command remains a single file ---

func TestExtractDesignMd_CommandIsSingleFile(t *testing.T) {
	root := projectRoot(t)

	// Check that extract-design-md exists as a single .md file (not a directory)
	mdPath := filepath.Join(root, "..", "plugins", "forge", "commands", "extract-design-md.md")
	dirPath := filepath.Join(root, "..", "plugins", "forge", "commands", "extract-design-md")

	mdExists := false
	if _, err := os.Stat(mdPath); err == nil {
		mdExists = true
	}

	dirExists := false
	if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
		dirExists = true
	}

	if !mdExists {
		t.Error("extract-design-md.md file does not exist")
	}
	if dirExists {
		t.Error("extract-design-md exists as a directory — command must remain a single file")
	}
}
