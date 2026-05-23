// Package lesson provides discovery and parsing for project lessons.
package lesson

import (
	"fmt"
	"strings"
	"time"

	"forge-cli/pkg/infocmd"
)

// lessonsDir is the base directory for lessons.
const lessonsDir = "docs/lessons"

// metadata holds the parsed frontmatter fields from a lesson .md file.
type metadata struct {
	Created  string   `yaml:"created"`
	Date     string   `yaml:"date"`
	Tags     []string `yaml:"tags"`
	Title    string   `yaml:"title"`
	Severity string   `yaml:"severity"`
}

// Lesson represents a discovered lesson with metadata and derived info.
type Lesson struct {
	Name     string
	Title    string
	Created  string
	Tags     []string
	Category string
	FilePath string
}

// Category prefixes map filename prefix to category name.
var categoryPrefixes = map[string]string{
	"gotcha-":  "gotcha",
	"arch-":    "architecture",
	"pattern-": "pattern",
	"tool-":    "tool",
	"lesson-":  "lesson",
	"hook-":    "hook",
}

// scanConfig is the infocmd.ScanConfig for lesson discovery.
var scanConfig = infocmd.ScanConfig[Lesson]{
	BaseDir:  lessonsDir,
	IsSubdir: false,
	IDKey:    func(l Lesson) string { return l.Name },
	CreatedKey: func(l Lesson) string {
		return l.Created
	},
	ParseEntry: func(name, path string, content []byte, _ time.Time) (Lesson, error) {
		var meta metadata
		if err := infocmd.ParseFrontmatter(content, &meta); err != nil {
			return Lesson{}, err
		}

		created := meta.Created
		if created == "" {
			created = meta.Date
		}

		return Lesson{
			Name:     name,
			Title:    meta.Title,
			Created:  created,
			Tags:     meta.Tags,
			Category: inferCategory(name),
			FilePath: path,
		}, nil
	},
}

// Discover walks docs/lessons/*.md and returns all lessons sorted by
// frontmatter created field descending (newest first), with mtime as fallback.
// Lessons without a created field fall back to file modification time.
func Discover(projectRoot string) ([]Lesson, error) {
	return infocmd.Discover(projectRoot, scanConfig)
}

// FindByName returns a single lesson by name (without .md extension), or an error if not found.
func FindByName(projectRoot, name string) (*Lesson, error) {
	result, err := infocmd.FindByID(projectRoot, name, scanConfig)
	if err != nil {
		return nil, fmt.Errorf("lesson not found: %s", name)
	}
	return result, nil
}

// inferCategory determines category from filename prefix.
func inferCategory(name string) string {
	for prefix, category := range categoryPrefixes {
		if strings.HasPrefix(name, prefix) {
			return category
		}
	}
	return ""
}
