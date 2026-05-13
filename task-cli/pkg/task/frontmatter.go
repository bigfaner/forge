package task

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// FrontmatterData holds the parsed YAML frontmatter fields from a task .md file.
type FrontmatterData struct {
	ID            string   `yaml:"id"`
	Title         string   `yaml:"title"`
	Priority      string   `yaml:"priority"`
	EstimatedTime string   `yaml:"estimated_time"`
	Dependencies  []string `yaml:"dependencies"`
	Scope         string   `yaml:"scope"`
	Breaking      bool     `yaml:"breaking"`
	MainSession   bool     `yaml:"mainSession"`
	NoTest        bool     `yaml:"noTest"`
	Type          string   `yaml:"type"`
	Profile       string   `yaml:"profile"`
}

// ParseFrontmatter extracts YAML frontmatter from a markdown file.
// Returns the parsed frontmatter and the remaining body content.
// If no frontmatter is found, returns a zero-value FrontmatterData with the
// original content as the body, and no error.
func ParseFrontmatter(content []byte) (FrontmatterData, []byte, error) {
	var fm FrontmatterData

	// Find opening ---
	line, rest, found := cutLine(content)
	if !found || !bytes.Equal(bytes.TrimSpace(line), []byte("---")) {
		return fm, content, nil
	}

	// Find closing ---
	closeIdx := bytes.Index(rest, []byte("\n---"))
	if closeIdx < 0 {
		// Check if it ends with ---\n at the very end
		if bytes.HasSuffix(bytes.TrimSpace(rest), []byte("---")) {
			closeIdx = bytes.LastIndex(rest, []byte("---"))
			yamlContent := rest[:closeIdx]
			if err := yaml.Unmarshal(bytes.TrimSpace(yamlContent), &fm); err != nil {
				return fm, nil, fmt.Errorf("parse frontmatter: %w", err)
			}
			return fm, nil, nil
		}
		return fm, content, nil
	}

	yamlContent := rest[:closeIdx]
	body := rest[closeIdx+4:] // skip \n---

	if err := yaml.Unmarshal(yamlContent, &fm); err != nil {
		return fm, nil, fmt.Errorf("parse frontmatter: %w", err)
	}

	return fm, body, nil
}

// cutLine splits content at the first newline.
func cutLine(content []byte) (line, rest []byte, found bool) {
	idx := bytes.IndexByte(content, '\n')
	if idx < 0 {
		return content, nil, false
	}
	return content[:idx], content[idx+1:], true
}
