// Package research provides discovery and parsing for research reports.
package research

import (
	"fmt"
	"path/filepath"
	"time"

	"forge-cli/pkg/infocmd"
)

// researchDir is the base directory for research reports.
const researchDir = "docs/research"

// metadata holds the parsed frontmatter fields from a research report .md file.
type metadata struct {
	Created    string   `yaml:"created"`
	Topic      string   `yaml:"topic"`
	Mode       string   `yaml:"mode"`
	Dimensions []string `yaml:"dimensions"`
	Candidates []string `yaml:"candidates"`
}

// Report represents a discovered research report with metadata and derived info.
type Report struct {
	Slug       string
	Created    string
	Topic      string
	Mode       string
	Dimensions []string
	Candidates []string
	FilePath   string
}

// scanConfig defines how research reports are discovered using the infocmd framework.
var scanConfig = infocmd.ScanConfig[Report]{
	BaseDir:  researchDir,
	IsSubdir: false,
	IDKey:    func(r Report) string { return r.Slug },
	CreatedKey: func(r Report) string {
		return r.Created
	},
	ParseEntry: func(name, _ string, content []byte, _ time.Time) (Report, error) {
		var meta metadata
		if err := infocmd.ParseFrontmatter(content, &meta); err != nil {
			return Report{}, err
		}

		// Skip files with no meaningful frontmatter content.
		if meta.Topic == "" && meta.Mode == "" {
			return Report{}, fmt.Errorf("no topic or mode")
		}

		return Report{
			Slug:       name,
			Created:    meta.Created,
			Topic:      meta.Topic,
			Mode:       meta.Mode,
			Dimensions: meta.Dimensions,
			Candidates: meta.Candidates,
			FilePath:   filepath.Join(researchDir, name+".md"),
		}, nil
	},
}

// Discover walks docs/research/*.md and returns all reports sorted by
// frontmatter created field descending (newest first), with mtime as fallback.
// Files without valid frontmatter are skipped.
// A missing docs/research/ directory returns an empty slice with no error.
func Discover(projectRoot string) ([]Report, error) {
	return infocmd.Discover(projectRoot, scanConfig)
}

// FindBySlug returns a single report by slug (filename without .md), or an error if not found.
func FindBySlug(projectRoot, slug string) (*Report, error) {
	r, err := infocmd.FindByID(projectRoot, slug, scanConfig)
	if err != nil {
		return nil, fmt.Errorf("research report not found: %s", slug)
	}
	return r, nil
}
