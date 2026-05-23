package task

import (
	"fmt"

	"forge-cli/pkg/infocmd"

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
	Type          string   `yaml:"type"`
	Coverage      *int     `yaml:"coverage"`
}

// ParseFrontmatter extracts YAML frontmatter from a markdown file.
// Returns the parsed frontmatter and the remaining body content.
// If no frontmatter is found, returns a zero-value FrontmatterData with the
// original content as the body, and no error.
func ParseFrontmatter(content []byte) (FrontmatterData, []byte, error) {
	rawYAML, body, err := infocmd.ExtractFrontmatter(content)
	if err != nil {
		return FrontmatterData{}, nil, fmt.Errorf("parse frontmatter: %w", err)
	}
	if rawYAML == nil {
		return FrontmatterData{}, content, nil
	}

	var fm FrontmatterData
	if err := yaml.Unmarshal(rawYAML, &fm); err != nil {
		return FrontmatterData{}, nil, fmt.Errorf("parse frontmatter: %w", err)
	}

	return fm, body, nil
}
