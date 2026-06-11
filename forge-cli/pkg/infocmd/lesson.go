package infocmd

import (
	"fmt"
	"strings"
	"time"
)

// lessonsDir is the base directory for lessons.
const lessonsDir = "docs/lessons"

// lessonMeta holds the parsed frontmatter fields from a lesson .md file.
type lessonMeta struct {
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

// lessonScanConfig is the ScanConfig for lesson discovery.
var lessonScanConfig = ScanConfig[Lesson]{
	BaseDir:  lessonsDir,
	IsSubdir: false,
	IDKey:    func(l Lesson) string { return l.Name },
	CreatedKey: func(l Lesson) string {
		return l.Created
	},
	ParseEntry: func(name, path string, content []byte, _ time.Time) (Lesson, error) {
		var meta lessonMeta
		if err := ParseFrontmatter(content, &meta); err != nil {
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

// DiscoverLessons walks docs/lessons/*.md and returns all lessons sorted by
// frontmatter created field descending (newest first), with mtime as fallback.
// Lessons without a created field fall back to file modification time.
func DiscoverLessons(projectRoot string) ([]Lesson, error) {
	return Discover(projectRoot, lessonScanConfig)
}

// FindLessonByName returns a single lesson by name (without .md extension), or an error if not found.
func FindLessonByName(projectRoot, name string) (*Lesson, error) {
	result, err := FindByID(projectRoot, name, lessonScanConfig)
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
