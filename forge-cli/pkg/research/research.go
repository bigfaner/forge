// Package research provides discovery and parsing for research reports.
package research

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
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

// Discover walks docs/research/*.md and returns all reports sorted by
// frontmatter created field descending (newest first), with mtime as fallback.
// Files without valid frontmatter are skipped.
// A missing docs/research/ directory returns an empty slice with no error.
func Discover(projectRoot string) ([]Report, error) {
	reportsDir := filepath.Join(projectRoot, researchDir)
	entries, err := os.ReadDir(reportsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read research directory: %w", err)
	}

	type reportWithMeta struct {
		report  Report
		modTime time.Time
	}

	var items []reportWithMeta
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(reportsDir, entry.Name())
		info, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var meta metadata
		if err := parseFrontmatter(data, &meta); err != nil {
			continue
		}

		// Skip files with no frontmatter (parseFrontmatter returns nil for no FM).
		if meta.Topic == "" && meta.Mode == "" {
			continue
		}

		slug := strings.TrimSuffix(entry.Name(), ".md")

		items = append(items, reportWithMeta{
			report: Report{
				Slug:       slug,
				Created:    meta.Created,
				Topic:      meta.Topic,
				Mode:       meta.Mode,
				Dimensions: meta.Dimensions,
				Candidates: meta.Candidates,
				FilePath:   filepath.Join(researchDir, entry.Name()),
			},
			modTime: info.ModTime(),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		ci, cj := items[i].report.Created, items[j].report.Created
		// Both have created: sort descending by date string.
		if ci != "" && cj != "" {
			return ci > cj
		}
		// If only one has created, it sorts first.
		if ci != "" {
			return true
		}
		if cj != "" {
			return false
		}
		// Both missing created: fall back to mtime descending.
		mi, mj := items[i].modTime, items[j].modTime
		if mi.IsZero() {
			return false
		}
		if mj.IsZero() {
			return true
		}
		return mi.After(mj)
	})

	reports := make([]Report, len(items))
	for i, it := range items {
		reports[i] = it.report
	}

	return reports, nil
}

// FindBySlug returns a single report by slug (filename without .md), or an error if not found.
func FindBySlug(projectRoot, slug string) (*Report, error) {
	reports, err := Discover(projectRoot)
	if err != nil {
		return nil, err
	}
	for _, r := range reports {
		if r.Slug == slug {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("research report not found: %s", slug)
}

// parseFrontmatter extracts YAML frontmatter from markdown content.
func parseFrontmatter(content []byte, target any) error {
	text := string(content)

	if !strings.HasPrefix(text, "---") {
		return nil
	}
	text = text[3:]

	closeIdx := strings.Index(text, "\n---")
	if closeIdx < 0 {
		return nil
	}

	yamlContent := text[:closeIdx]
	return yaml.Unmarshal([]byte(yamlContent), target)
}
