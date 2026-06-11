// Package infocmd provides shared generic utilities for info-commands
// (research, proposal, lesson, etc.). It encapsulates the common patterns
// of directory scanning, frontmatter parsing, sorting, and lookup.
package infocmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ScanConfig defines how Discover[T] should scan a directory for entries.
// Each info-command provides its own ScanConfig to customize scanning behavior.
type ScanConfig[T any] struct {
	// BaseDir is the relative directory to scan (e.g. "docs/research", "docs/proposals").
	BaseDir string

	// IsSubdir controls the scanning mode:
	//   false (flat mode): scan BaseDir/*.md files, slug = filename minus .md
	//   true (subdir mode): scan BaseDir/*/FileName files, slug = subdirectory name
	IsSubdir bool

	// FileName is the file to look for inside each subdirectory.
	// Only used when IsSubdir is true (e.g. "proposal.md").
	FileName string

	// IDKey extracts the identifier from an entry (Slug or Name).
	IDKey func(T) string

	// CreatedKey extracts the created-date string from an entry for sorting.
	// Return empty string to fall back to file mtime.
	CreatedKey func(T) string

	// ParseEntry constructs a typed entry from a discovered file.
	// name is the entry identifier (slug or name).
	// path is the full filesystem path to the file.
	// content is the raw file bytes.
	// modTime is the file modification time.
	ParseEntry func(name, path string, content []byte, modTime time.Time) (T, error)
}

// ExtractFrontmatter extracts raw YAML bytes and body from markdown content.
// Returns the raw YAML bytes between the --- delimiters and the remaining body.
// If no valid frontmatter is found, returns nil bytes with no error.
func ExtractFrontmatter(content []byte) (rawYAML []byte, body []byte, err error) {
	text := string(content)

	if !strings.HasPrefix(text, "---") {
		return nil, content, nil
	}
	text = text[3:]

	// Require newline or EOF after opening ---
	if len(text) > 0 && text[0] != '\n' {
		return nil, content, nil
	}
	if len(text) > 0 {
		text = text[1:]
	}

	closeIdx := strings.Index(text, "\n---")
	if closeIdx < 0 {
		// Check if content ends with --- (no trailing newline)
		if strings.HasSuffix(strings.TrimSpace(text), "---") {
			lastIdx := strings.LastIndex(text, "---")
			raw := []byte(strings.TrimSpace(text[:lastIdx]))
			return raw, nil, nil
		}
		return nil, content, nil
	}

	rawYAML = []byte(text[:closeIdx])
	body = []byte(text[closeIdx+4:]) // skip \n---
	return rawYAML, body, nil
}

// ParseFrontmatter extracts YAML frontmatter from markdown content
// and unmarshals it into target.
// Returns nil with no error if no valid frontmatter is found.
// Returns an error if YAML parsing fails.
func ParseFrontmatter(content []byte, target any) error {
	rawYAML, _, err := ExtractFrontmatter(content)
	if err != nil {
		return err
	}
	if rawYAML == nil {
		return nil
	}
	return yaml.Unmarshal(rawYAML, target)
}

// Discover walks the directory defined by ScanConfig and returns all entries
// sorted by the created date descending (newest first), with mtime as fallback.
// Files that fail parsing are skipped.
// A missing directory returns an empty slice with no error.
func Discover[T any](projectRoot string, cfg ScanConfig[T]) ([]T, error) {
	scanDir := filepath.Join(projectRoot, cfg.BaseDir)
	entries, err := os.ReadDir(scanDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read directory %s: %w", cfg.BaseDir, err)
	}

	type itemWithMeta struct {
		item    T
		modTime time.Time
	}

	var items []itemWithMeta

	for _, entry := range entries {
		var name, filePath string

		if cfg.IsSubdir {
			// Subdir mode: entry must be a directory
			if !entry.IsDir() {
				continue
			}
			name = entry.Name()
			filePath = filepath.Join(scanDir, name, cfg.FileName)
		} else {
			// Flat mode: entry must be a .md file
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			name = strings.TrimSuffix(entry.Name(), ".md")
			filePath = filepath.Join(scanDir, entry.Name())
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		info, err := os.Stat(filePath)
		if err != nil {
			continue
		}
		modTime := info.ModTime()

		item, err := cfg.ParseEntry(name, filePath, data, modTime)
		if err != nil {
			continue
		}

		items = append(items, itemWithMeta{
			item:    item,
			modTime: modTime,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		ci, cj := cfg.CreatedKey(items[i].item), cfg.CreatedKey(items[j].item)

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

	result := make([]T, len(items))
	for i, it := range items {
		result[i] = it.item
	}

	return result, nil
}

// FindByID searches the results of Discover for an entry whose identifier
// (as defined by ScanConfig.IDKey) matches the given id.
// Returns an error if no matching entry is found.
func FindByID[T any](projectRoot, id string, cfg ScanConfig[T]) (*T, error) {
	items, err := Discover(projectRoot, cfg)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if cfg.IDKey(item) == id {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("item not found: %s", id)
}
