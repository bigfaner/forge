package task

import (
	"os"
	"path/filepath"
	"strings"
)

// extractBodyContext builds a BodyContext by reading planning-time data from
// the proposal or PRD file. Missing files produce empty fields — this is valid.
func extractBodyContext(projectRoot, slug, mode string, interfaces []string) BodyContext {
	ctx := BodyContext{
		FeatureSlug: slug,
		Mode:        mode,
		Interfaces:  interfaces,
	}

	if mode == "" {
		return ctx
	}

	// Read the source file (proposal for quick, PRD for breakdown)
	var content string
	switch mode {
	case "quick":
		propPath := filepath.Join(projectRoot, "docs", "proposals", slug, "proposal.md")
		data, err := os.ReadFile(propPath)
		if err != nil {
			return ctx // missing proposal is valid
		}
		content = string(data)
		ctx.Scope = extractScope(content)
		ctx.SuccessCriteria = extractSuccessCriteria(content)

	case "breakdown":
		prdPath := filepath.Join(projectRoot, "docs", "features", slug, "prd", "prd-spec.md")
		data, err := os.ReadFile(prdPath)
		if err != nil {
			return ctx // missing PRD is valid
		}
		content = string(data)
		ctx.Scope = extractScope(content)
		ctx.SuccessCriteria = extractSuccessCriteria(content)
		ctx.AcceptanceCriteria = extractAcceptanceCriteria(content)
	}

	return ctx
}

// extractScope parses "## Scope > ### In Scope" section and returns bullet items.
// Returns nil if the section is not found or empty.
func extractScope(content string) []string {
	return extractBulletItems(content, "### In Scope")
}

// extractSuccessCriteria parses "## Success Criteria" section and returns checkbox items.
// Returns nil if the section is not found or empty.
func extractSuccessCriteria(content string) []string {
	return extractCheckboxItems(content, "## Success Criteria")
}

// extractAcceptanceCriteria parses "## Acceptance Criteria" section and returns checkbox items.
// Returns nil if the section is not found or empty.
func extractAcceptanceCriteria(content string) []string {
	return extractCheckboxItems(content, "## Acceptance Criteria")
}

// extractBulletItems extracts "- " bullet list items under a subheading.
// The subheading is the exact line to match (e.g., "### In Scope").
// Collection stops at the next heading (## or ###) or end of content.
func extractBulletItems(content, subheading string) []string {
	lines := strings.Split(content, "\n")

	// Find the subheading line
	startIdx := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == subheading {
			startIdx = i + 1
			break
		}
	}
	if startIdx < 0 {
		return nil
	}

	// Collect bullet items until we hit another heading
	var items []string
	for i := startIdx; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Stop at next heading
		if strings.HasPrefix(trimmed, "##") {
			break
		}

		// Skip blank lines
		if trimmed == "" {
			continue
		}

		// Collect "- " items (strip checkbox prefixes if present)
		if strings.HasPrefix(trimmed, "- ") {
			item := strings.TrimPrefix(trimmed, "- ")
			item = strings.TrimSpace(item)
			// Strip checkbox prefixes: "[ ] " or "[x] " or "[X] "
			switch {
			case strings.HasPrefix(item, "[ ] "):
				item = strings.TrimPrefix(item, "[ ] ")
			case strings.HasPrefix(item, "[x] "):
				item = strings.TrimPrefix(item, "[x] ")
			case strings.HasPrefix(item, "[X] "):
				item = strings.TrimPrefix(item, "[X] ")
			}
			item = strings.TrimSpace(item)
			if item != "" {
				items = append(items, item)
			}
		}
	}

	return items
}

// extractCheckboxItems extracts "- [ ] " or "- [x] " checkbox items under a heading.
// The heading is the exact "## " line to match.
// Collection stops at the next ## heading or end of content.
// Sub-items (indented bullets) are skipped.
func extractCheckboxItems(content, heading string) []string {
	lines := strings.Split(content, "\n")

	// Find the heading line
	startIdx := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == heading {
			startIdx = i + 1
			break
		}
	}
	if startIdx < 0 {
		return nil
	}

	// Collect checkbox items until we hit another ## heading
	var items []string
	for i := startIdx; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Stop at next ## heading
		if strings.HasPrefix(trimmed, "## ") {
			break
		}

		// Skip blank lines and non-checkbox items
		if trimmed == "" {
			continue
		}

		// Only top-level checkboxes (no leading whitespace beyond trim)
		if strings.HasPrefix(trimmed, "- [ ] ") {
			item := strings.TrimPrefix(trimmed, "- [ ] ")
			item = strings.TrimSpace(item)
			if item != "" {
				items = append(items, item)
			}
		} else if strings.HasPrefix(trimmed, "- [x] ") {
			item := strings.TrimPrefix(trimmed, "- [x] ")
			item = strings.TrimSpace(item)
			if item != "" {
				items = append(items, item)
			}
		}
	}

	return items
}
