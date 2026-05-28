package task

import (
	"os"
	"path/filepath"
	"strings"
)

// extractDocTaskCriteria scans a tasks directory for doc-category task .md files,
// extracts the "## Acceptance Criteria" section from each, and returns a map
// of task name (filename without .md) to raw AC markdown content.
// Only doc-category tasks are included; non-doc tasks are skipped.
// Tasks without an AC section are included with empty string content so that
// callers can detect missing AC and emit appropriate warnings.
func extractDocTaskCriteria(taskDir string) map[string]string {
	entries, err := os.ReadDir(taskDir)
	if err != nil {
		return nil
	}

	result := make(map[string]string)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		if shouldSkipFile(entry.Name()) {
			continue
		}

		// Read file to determine type
		filePath := filepath.Join(taskDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		// Parse frontmatter to get type
		fm, _, err := ParseFrontmatter(data)
		if err != nil || fm.ID == "" {
			continue
		}

		// Only process doc-category business tasks
		taskType := fm.Type
		if taskType == "" {
			taskType = InferType(fm.ID, nil)
		}
		if CategoryForType(taskType) != CategoryDoc {
			continue
		}
		// Skip system types (doc.review, doc.summary, etc.)
		if IsSystemType(taskType) {
			continue
		}

		// Extract AC section (empty string if section not found)
		content := string(data)
		acContent, _ := extractACSection(content)

		taskName := strings.TrimSuffix(entry.Name(), ".md")
		result[taskName] = acContent
	}

	return result
}

// isACHeading returns true if the trimmed line matches a recognized
// Acceptance Criteria heading. Supports:
//   - exact: "## Acceptance Criteria"
//   - case-insensitive: "## Acceptance criteria"
//   - Chinese alias: "## 验收标准"
func isACHeading(trimmed string) bool {
	if strings.HasPrefix(trimmed, "## ") {
		rest := strings.TrimSpace(trimmed[3:])
		if strings.EqualFold(rest, "Acceptance Criteria") {
			return true
		}
		if rest == "验收标准" {
			return true
		}
	}
	return false
}

// extractACSection extracts the content between "## Acceptance Criteria" and
// the next "## " heading (or end of file). Returns the content (everything
// after the heading line, including newlines) and true if found.
// Respects fenced code blocks: ## inside ``` blocks are not treated as section boundaries.
// Title matching is tolerant: supports case-insensitive "Acceptance Criteria" and Chinese alias "验收标准".
func extractACSection(content string) (string, bool) {
	lines := strings.Split(content, "\n")

	// Find the AC heading line (case-insensitive + Chinese alias)
	startIdx := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if isACHeading(trimmed) {
			startIdx = i + 1
			break
		}
	}
	if startIdx < 0 {
		return "", false
	}

	// Collect lines until next ## heading (respecting fenced code blocks)
	var collected []string
	inCodeBlock := false
	for i := startIdx; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Track fenced code blocks
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			collected = append(collected, line)
			continue
		}

		// Stop at next ## heading only if not inside a code block
		if !inCodeBlock && strings.HasPrefix(trimmed, "## ") {
			break
		}

		collected = append(collected, line)
	}

	result := strings.Join(collected, "\n")
	return result, true
}

// extractBodyContext builds a BodyContext by reading planning-time data from
// the proposal or PRD file. Missing files produce empty fields — this is valid.
func extractBodyContext(projectRoot, slug, mode string, surfaceTypes []string) BodyContext {
	ctx := BodyContext{
		FeatureSlug:  slug,
		Mode:         mode,
		SurfaceTypes: surfaceTypes,
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
		ctx.SuccessCriteria = extractSuccessCriteria(content)

	case "breakdown":
		prdPath := filepath.Join(projectRoot, "docs", "features", slug, "prd", "prd-spec.md")
		data, err := os.ReadFile(prdPath)
		if err != nil {
			return ctx // missing PRD is valid
		}
		content = string(data)
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
