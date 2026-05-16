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
