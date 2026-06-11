package infocmd

import (
	"fmt"
	"path/filepath"
	"time"
)

// researchDir is the base directory for research reports.
const researchDir = "docs/research"

// researchMeta holds the parsed frontmatter fields from a research report .md file.
type researchMeta struct {
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

// researchScanConfig defines how research reports are discovered.
var researchScanConfig = ScanConfig[Report]{
	BaseDir:  researchDir,
	IsSubdir: false,
	IDKey:    func(r Report) string { return r.Slug },
	CreatedKey: func(r Report) string {
		return r.Created
	},
	ParseEntry: func(name, _ string, content []byte, _ time.Time) (Report, error) {
		var meta researchMeta
		if err := ParseFrontmatter(content, &meta); err != nil {
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

// DiscoverReports walks docs/research/*.md and returns all reports sorted by
// frontmatter created field descending (newest first), with mtime as fallback.
// Files without valid frontmatter are skipped.
// A missing docs/research/ directory returns an empty slice with no error.
func DiscoverReports(projectRoot string) ([]Report, error) {
	return Discover(projectRoot, researchScanConfig)
}

// FindReportBySlug returns a single report by slug (filename without .md), or an error if not found.
func FindReportBySlug(projectRoot, slug string) (*Report, error) {
	r, err := FindByID(projectRoot, slug, researchScanConfig)
	if err != nil {
		return nil, fmt.Errorf("research report not found: %s", slug)
	}
	return r, nil
}
