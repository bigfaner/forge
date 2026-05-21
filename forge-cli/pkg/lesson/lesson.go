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

// Discover walks docs/lessons/*.md and returns all lessons sorted by
// frontmatter created field descending (newest first), with mtime as fallback.
// Lessons without a created field fall back to file modification time.
func Discover(projectRoot string) ([]Lesson, error) {
	lessonsDir := filepath.Join(projectRoot, lessonsDir)
	entries, err := os.ReadDir(lessonsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read lessons directory: %w", err)
	}

	type lessonWithMeta struct {
		lesson  Lesson
		modTime time.Time
	}

	var items []lessonWithMeta
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

		var meta metadata
		if err := parseFrontmatter(data, &meta); err != nil {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")
		created := meta.Created
		if created == "" {
			created = meta.Date
		}

		category := inferCategory(name)

		items = append(items, lessonWithMeta{
			lesson: Lesson{
				Name:     name,
				Title:    meta.Title,
				Created:  created,
				Tags:     meta.Tags,
				Category: category,
				FilePath: filepath.Join(lessonsDir, entry.Name()),
			},
			modTime: info.ModTime(),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		ci, cj := items[i].lesson.Created, items[j].lesson.Created
		// Items with created field sort before items without (descending).
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
