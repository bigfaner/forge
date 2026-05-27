package task

import (
	"fmt"
	"os"

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
	SurfaceKey    string   `yaml:"surface-key,omitempty"`
	SurfaceType   string   `yaml:"surface-type,omitempty"`
	Breaking      bool     `yaml:"breaking"`
	MainSession   bool     `yaml:"mainSession"`
	Type          string   `yaml:"type"`
	Coverage      *int     `yaml:"coverage"`
	// Complexity is the task complexity level: "low", "medium", or "high".
	// Empty value defaults to "medium" at runtime.
	Complexity string `yaml:"complexity,omitempty"`
	Scope      string `yaml:"scope"`
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

// WriteFrontmatter writes a task .md file with the given frontmatter data and body.
// It serializes the frontmatter as YAML and prepends it to the body with --- delimiters.
func WriteFrontmatter(path string, fm FrontmatterData, body []byte) error {
	var yamlBuf []byte
	yamlBuf = append(yamlBuf, "---\n"...)

	yamlData, err := yaml.Marshal(fm)
	if err != nil {
		return fmt.Errorf("marshal frontmatter: %w", err)
	}
	yamlBuf = append(yamlBuf, yamlData...)
	yamlBuf = append(yamlBuf, "---\n"...)
	yamlBuf = append(yamlBuf, body...)

	return os.WriteFile(path, yamlBuf, 0644)
}
