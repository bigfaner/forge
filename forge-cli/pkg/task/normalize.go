package task

import (
	"bytes"
	"strings"
)

// NormalizeTaskMD removes empty ## Hard Rules sections from task markdown content.
// A section is considered empty if it contains no non-whitespace content between
// the ## Hard Rules heading and the next ## heading (or end of content).
// Returns the normalized content. If no changes are needed, returns the input unchanged.
func NormalizeTaskMD(content []byte) []byte {
	_, bodyStart := splitFrontmatter(content)
	if bodyStart < 0 {
		return content
	}

	frontmatter := content[:bodyStart]
	body := string(content[bodyStart:])

	normalized := removeEmptyHardRulesSection(body)
	if normalized == body {
		return content
	}

	return append(frontmatter, []byte(normalized)...)
}

// removeEmptyHardRulesSection removes empty ## Hard Rules sections from the body.
func removeEmptyHardRulesSection(body string) string {
	const marker = "## Hard Rules"

	idx := strings.Index(body, marker)
	if idx < 0 {
		return body
	}

	// Find end of heading line
	headingEnd := idx + len(marker)
	for headingEnd < len(body) && body[headingEnd] != '\n' {
		headingEnd++
	}
	if headingEnd < len(body) {
		headingEnd++ // skip \n
	}

	// Find next ## heading or EOF
	nextSectionStart := findNextHeading(body, headingEnd)

	// Check if content between heading and next section is all whitespace
	sectionContent := body[headingEnd:nextSectionStart]
	if strings.TrimSpace(sectionContent) == "" {
		before := body[:idx]
		var after string
		if nextSectionStart < len(body) {
			after = body[nextSectionStart:]
		}

		before = strings.TrimRight(before, "\r\n") + "\n"
		if strings.HasPrefix(after, "##") {
			before += "\n"
		}
		return before + after
	}

	return body
}

// findNextHeading returns the index of the next ## heading starting from pos,
// or len(s) if no heading is found.
func findNextHeading(s string, pos int) int {
	for i := pos; i < len(s); {
		nl := strings.Index(s[i:], "\n")
		if nl < 0 {
			return len(s)
		}
		lineStart := i + nl + 1
		if lineStart < len(s) && s[lineStart] == '#' && lineStart+1 < len(s) && s[lineStart+1] == '#' {
			return lineStart
		}
		i = lineStart
	}
	return len(s)
}

// splitFrontmatter returns the index where frontmatter ends (after the closing ---).
// Returns (0, -1) if no valid frontmatter is found.
func splitFrontmatter(content []byte) (int, int) { //nolint:unparam
	if !bytes.HasPrefix(content, []byte("---")) {
		return 0, -1
	}
	closing := bytes.Index(content[3:], []byte("\n---"))
	if closing < 0 {
		return 0, -1
	}
	end := closing + 3 + 4
	for end < len(content) && (content[end] == '\n' || content[end] == '\r') {
		end++
	}
	return 0, end
}
