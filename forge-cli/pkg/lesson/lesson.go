// Package lesson provides discovery and parsing for project lessons.
package lesson

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// LessonsDir is the base directory for lessons.
const LessonsDir = "docs/lessons"

// Metadata holds the parsed frontmatter fields from a lesson .md file.
type Metadata struct {
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

// Discover walks docs/lessons/*.md and returns all lessons sorted by file
// modification time in reverse chronological order (newest first).
// Lessons with zero or missing modification times sort to the end.
func Discover(projectRoot string) ([]Lesson, error) {
	lessonsDir := filepath.Join(projectRoot, LessonsDir)
	entries, err := os.ReadDir(lessonsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read lessons directory: %w", err)
	}

	type lessonWithTime struct {
		lesson  Lesson
		modTime time.Time
	}

	var items []lessonWithTime
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(lessonsDir, entry.Name())
		info, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var meta Metadata
		if err := parseFrontmatter(data, &meta); err != nil {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")
		created := meta.Date
		if created == "" {
			created = info.ModTime().Format("2006-01-02")
		}

		category := inferCategory(name)

		items = append(items, lessonWithTime{
			lesson: Lesson{
				Name:     name,
				Title:    meta.Title,
				Created:  created,
				Tags:     meta.Tags,
				Category: category,
				FilePath: filepath.Join(LessonsDir, entry.Name()),
			},
			modTime: info.ModTime(),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		mi, mj := items[i].modTime, items[j].modTime
		// Zero mod times sort to the end.
		if mi.IsZero() {
			return false
		}
		if mj.IsZero() {
			return true
		}
		return mi.After(mj)
	})

	lessons := make([]Lesson, len(items))
	for i, it := range items {
		lessons[i] = it.lesson
	}

	return lessons, nil
}

// FindByName returns a single lesson by name (without .md extension), or an error if not found.
func FindByName(projectRoot, name string) (*Lesson, error) {
	lessons, err := Discover(projectRoot)
	if err != nil {
		return nil, err
	}
	for _, l := range lessons {
		if l.Name == name {
			return &l, nil
		}
	}
	return nil, fmt.Errorf("lesson not found: %s", name)
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
