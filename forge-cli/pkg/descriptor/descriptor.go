// Package descriptor converts semantic descriptors from Contract specifications
// into precise regex patterns using Fact Table entries as ground truth.
package descriptor

import (
	"fmt"
	"regexp"
	"strings"
)

// FactEntry represents a single entry from the Fact Table built during
// code reconnaissance. Each entry maps a conceptual key to a concrete
// value observed in the codebase, with a source citation.
type FactEntry struct {
	Key    string // e.g. "CLI_TASK_CLAIM_OUTPUT"
	Value  string // e.g. "claimed task <task_id>"
	Source string // e.g. "internal/cmd/claim.go:42"
}

// RegexBuildResult holds the output of a semantic descriptor to regex conversion.
type RegexBuildResult struct {
	// Pattern is the generated regex pattern string.
	Pattern string
	// Captures maps capture group names to their group indices.
	Captures map[string]int
	// Unresolved is true when the descriptor could not be fully resolved
	// from the Fact Table (e.g. no matching entry found).
	Unresolved bool
	// Placeholder keeps the original semantic descriptor when unresolved.
	Placeholder string
}

// BuildRegex converts a semantic descriptor string into a regex pattern
// based on Fact Table entries. The descriptor uses natural language patterns
// like "success confirmation containing feature-slug" and the Fact Table
// provides actual output examples from code reconnaissance.
//
// Conversion rules:
//   - If a FactEntry matches the descriptor's semantic intent, its Value is used
//     as the basis for the regex pattern.
//   - Placeholder tokens like <task_id>, <feature-slug> are converted to
//     named capture groups.
//   - If no matching FactEntry is found, the result is marked Unresolved.
func BuildRegex(descriptor string, facts []FactEntry) RegexBuildResult {
	// Try to find a matching fact entry by semantic similarity
	matched := findMatchingFact(descriptor, facts)
	if matched == nil {
		return RegexBuildResult{
			Unresolved:  true,
			Placeholder: descriptor,
		}
	}

	pattern := valueToRegex(matched.Value)
	captures := extractCaptures(matched.Value)

	return RegexBuildResult{
		Pattern:  pattern,
		Captures: captures,
	}
}

// findMatchingFact attempts to match a semantic descriptor to a Fact Table entry.
// It uses keyword matching: extracts significant words from the descriptor
// and looks for entries whose Key or Value contains those keywords.
func findMatchingFact(descriptor string, facts []FactEntry) *FactEntry {
	keywords := extractKeywords(descriptor)
	if len(keywords) == 0 {
		return nil
	}

	var bestMatch *FactEntry
	bestScore := 0

	for i := range facts {
		f := &facts[i]
		score := scoreMatch(keywords, f.Key, f.Value)
		if score > bestScore && score >= 2 { // minimum 2 keyword matches
			bestScore = score
			bestMatch = f
		}
	}

	return bestMatch
}

// extractKeywords extracts significant words from a semantic descriptor.
// Filters out common stop words and short words.
func extractKeywords(descriptor string) []string {
	stopWords := map[string]bool{
		"a": true, "an": true, "the": true, "is": true, "are": true,
		"was": true, "were": true, "be": true, "been": true,
		"being": true, "have": true, "has": true, "had": true,
		"do": true, "does": true, "did": true, "will": true,
		"would": true, "shall": true, "should": true, "may": true,
		"might": true, "can": true, "could": true, "of": true,
		"in": true, "to": true, "for": true, "with": true,
		"on": true, "at": true, "by": true, "from": true,
		"or": true, "and": true, "not": true, "no": true,
		"but": true, "it": true, "its": true, "this": true,
		"that": true, "than": true, "then": true, "so": true,
		"if": true, "as": true, "up": true, "out": true,
	}

	words := strings.Fields(strings.ToLower(descriptor))
	var keywords []string
	for _, w := range words {
		// Clean punctuation
		w = strings.Trim(w, ".,;:!?'\"()[]{}<>")
		if len(w) >= 3 && !stopWords[w] {
			keywords = append(keywords, w)
		}
	}
	return keywords
}

// scoreMatch counts how many keywords appear in (or partially match) the key or value.
// A keyword partially matches if either side of a hyphen in the keyword matches.
func scoreMatch(keywords []string, factKey, factValue string) int {
	combined := strings.ToLower(factKey + " " + factValue)
	score := 0
	for _, kw := range keywords {
		if strings.Contains(combined, kw) {
			score++
			continue
		}
		// Partial match: check if any hyphen-separated segment of the keyword matches
		parts := strings.Split(kw, "-")
		matched := false
		for _, part := range parts {
			if len(part) >= 3 && strings.Contains(combined, part) {
				matched = true
				break
			}
		}
		if matched {
			score++
		}
	}
	return score
}

// valueToRegex converts a Fact Table value into a regex pattern.
// Placeholder tokens like <task_id> become named capture groups ([\w-]+).
// Literal text is regex-escaped.
func valueToRegex(value string) string {
	// Find all <placeholder> tokens and their positions
	placeholderRe := regexp.MustCompile(`<([\w-]+)>`)

	// Split the value by placeholders, escaping literal parts
	matches := placeholderRe.FindAllStringSubmatchIndex(value, -1)

	var sb strings.Builder
	lastEnd := 0
	for _, m := range matches {
		// Write escaped literal text before this placeholder
		literal := value[lastEnd:m[0]]
		sb.WriteString(regexp.QuoteMeta(literal))

		// Write named capture group for the placeholder
		name := value[m[2]:m[3]]
		// Sanitize name for regex group: replace hyphens with underscores
		cleanName := strings.ReplaceAll(name, "-", "_")
		fmt.Fprintf(&sb, `(?P<%s>[\w-]+)`, cleanName)

		lastEnd = m[1]
	}

	// Write remaining literal text
	if lastEnd < len(value) {
		sb.WriteString(regexp.QuoteMeta(value[lastEnd:]))
	}

	return sb.String()
}

// extractCaptures builds a map of capture group names to their indices
// from a value containing <placeholder> tokens.
func extractCaptures(value string) map[string]int {
	placeholderRe := regexp.MustCompile(`<([\w-]+)>`)
	matches := placeholderRe.FindAllStringSubmatch(value, -1)

	captures := make(map[string]int)
	for i, m := range matches {
		name := strings.ReplaceAll(m[1], "-", "_")
		captures[name] = i + 1 // 1-indexed
	}
	return captures
}

// BuildAssertionRegex is a convenience function that takes a semantic descriptor
// and Fact Table, and returns a complete regex pattern suitable for use in test
// assertions. Returns the original descriptor in quotes if unresolved.
func BuildAssertionRegex(descriptor string, facts []FactEntry) string {
	result := BuildRegex(descriptor, facts)
	if result.Unresolved {
		return regexp.QuoteMeta(result.Placeholder)
	}
	return result.Pattern
}

// BuildSensitiveFieldPlaceholder replaces known sensitive field names
// with environment variable placeholders. This implements the test data
// safety hard rule: no hardcoded secrets in generated test code.
func BuildSensitiveFieldPlaceholder(fieldName string) string {
	sensitiveFields := map[string]string{
		"token":        "<from-env>",
		"api_key":      "<from-env>",
		"secret":       "<from-env>",
		"password":     "<from-env>",
		"access_token": "<from-env>",
		"auth_token":   "<from-env>",
		"session_id":   "<from-env>",
	}

	lower := strings.ToLower(fieldName)
	if placeholder, ok := sensitiveFields[lower]; ok {
		return placeholder
	}
	return fieldName
}

// IsSensitiveField checks whether a field name is considered sensitive
// and should use placeholder values in generated test code.
func IsSensitiveField(fieldName string) bool {
	sensitivePatterns := []string{"token", "secret", "password", "key", "credential"}
	lower := strings.ToLower(fieldName)
	for _, p := range sensitivePatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}
