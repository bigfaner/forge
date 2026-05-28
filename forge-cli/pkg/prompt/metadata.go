package prompt

import (
	"fmt"
	"reflect"
	"strings"
)

// TemplateMetadata holds the parsed metadata frontmatter from a template file.
type TemplateMetadata struct {
	Type      string   // task type constant (e.g. "coding.feature")
	Category  string   // task category (coding/doc/test/eval/validation/gate/record)
	Variables []string // list of template variable names used for validation
}

// parseMetadataFrontmatter extracts metadata from between the first pair of ---
// markers in a template file. Returns the remaining content (after the closing ---)
// and the parsed metadata.
// If no frontmatter is found (no leading ---), returns the original content with
// a nil metadata (not an error — metadata is optional for backward compatibility).
func parseMetadataFrontmatter(content string) (body string, meta *TemplateMetadata) {
	trimmed := strings.TrimLeft(content, " \t\n")
	if !strings.HasPrefix(trimmed, "---") {
		// No frontmatter — return content as-is
		return content, nil
	}

	// Find the closing ---
	afterOpen := trimmed[3:] // skip opening ---
	// Skip optional newline after opening ---
	if len(afterOpen) > 0 && afterOpen[0] == '\n' {
		afterOpen = afterOpen[1:]
	} else if len(afterOpen) > 1 && afterOpen[0] == '\r' && afterOpen[1] == '\n' {
		afterOpen = afterOpen[2:]
	}

	closeIdx := strings.Index(afterOpen, "\n---")
	if closeIdx < 0 {
		// No closing --- found — not valid frontmatter, return as-is
		return content, nil
	}

	frontmatter := afterOpen[:closeIdx]
	// Skip past closing --- and newline
	remaining := afterOpen[closeIdx+4:] // skip "\n---"
	if len(remaining) > 0 && remaining[0] == '\n' {
		remaining = remaining[1:]
	} else if len(remaining) > 1 && remaining[0] == '\r' && remaining[1] == '\n' {
		remaining = remaining[2:]
	}

	// Parse YAML-like frontmatter (simple line-based parser)
	meta = &TemplateMetadata{}
	for _, line := range strings.Split(frontmatter, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		switch {
		case strings.HasPrefix(line, "type:"):
			meta.Type = strings.TrimSpace(strings.TrimPrefix(line, "type:"))
			// Remove surrounding quotes if present
			meta.Type = strings.Trim(meta.Type, "\"")
		case strings.HasPrefix(line, "category:"):
			meta.Category = strings.TrimSpace(strings.TrimPrefix(line, "category:"))
			meta.Category = strings.Trim(meta.Category, "\"")
		case strings.HasPrefix(line, "- ") && meta.Variables != nil:
			// Variable list item
			varName := strings.TrimSpace(strings.TrimPrefix(line, "- "))
			varName = strings.Trim(varName, "\"")
			meta.Variables = append(meta.Variables, varName)
		case strings.HasPrefix(line, "variables:"):
			// Initialize variables list (may be empty)
			meta.Variables = []string{}
		}
	}

	return remaining, meta
}

// stripMetadataFrontmatter removes metadata frontmatter from template content.
// Returns the body content (everything after the closing ---).
// If no metadata frontmatter is found, returns the original content unchanged.
func stripMetadataFrontmatter(content string) string {
	body, _ := parseMetadataFrontmatter(content)
	return body
}

// validateMetadataVariables checks that each variable declared in metadata
// exists as an exported field on the given struct type.
// Returns an error listing any variables that don't have matching struct fields.
func validateMetadataVariables(meta *TemplateMetadata, structType reflect.Type) error {
	if meta == nil || len(meta.Variables) == 0 {
		return nil
	}

	var mismatches []string
	for _, varName := range meta.Variables {
		if !structHasField(structType, varName) {
			mismatches = append(mismatches, varName)
		}
	}

	if len(mismatches) > 0 {
		return fmt.Errorf("metadata variables not found in %s struct: %s", structType.Name(), strings.Join(mismatches, ", "))
	}
	return nil
}

// structHasField checks if a struct type has an exported field with the given name.
func structHasField(t reflect.Type, fieldName string) bool {
	for t.Kind() == 4 { // reflect.Ptr
		t = t.Elem()
	}
	if t.Kind() != 25 { // reflect.Struct
		return false
	}
	_, ok := t.FieldByName(fieldName)
	return ok
}
