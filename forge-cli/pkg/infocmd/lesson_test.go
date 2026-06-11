package infocmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoverLessons_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	lessons, err := DiscoverLessons(dir)
	assert.NoError(t, err)
	assert.Empty(t, lessons)
}

func TestDiscoverLessons_NoDir(t *testing.T) {
	dir := t.TempDir()

	lessons, err := DiscoverLessons(dir)
	assert.NoError(t, err)
	assert.Empty(t, lessons)
}

func TestDiscoverLessons_SingleLesson(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	content := `---
date: 2026-05-10
tags:
  - testing
  - e2e
title: "Test lesson"
---

# Test Lesson

Some content.
`
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-test.md"), []byte(content), 0644))

	lessons, err := DiscoverLessons(dir)
	assert.NoError(t, err)
	require.Len(t, lessons, 1)

	l := lessons[0]
	assert.Equal(t, "lesson-test", l.Name)
	assert.Equal(t, "2026-05-10", l.Created)
	assert.Equal(t, "Test lesson", l.Title)
	assert.Equal(t, []string{"testing", "e2e"}, l.Tags)
	assert.Equal(t, "lesson", l.Category)
}

func TestDiscoverLessons_CategoryInference(t *testing.T) {
	tests := []struct {
		filename     string
		wantCategory string
	}{
		{"gotcha-bad-thing.md", "gotcha"},
		{"arch-system-design.md", "architecture"},
		{"pattern-reuse.md", "pattern"},
		{"tool-makefile.md", "tool"},
		{"lesson-learned.md", "lesson"},
		{"hook-lifecycle.md", "hook"},
		{"other-topic.md", ""},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			dir := t.TempDir()
			lessonsDir := filepath.Join(dir, "docs/lessons")
			require.NoError(t, os.MkdirAll(lessonsDir, 0755))

			content := "---\ndate: 2026-01-01\n---\n"
			require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, tt.filename), []byte(content), 0644))

			lessons, err := DiscoverLessons(dir)
			assert.NoError(t, err)
			require.Len(t, lessons, 1)
			assert.Equal(t, tt.wantCategory, lessons[0].Category, "category for %s", tt.filename)
		})
	}
}

func TestDiscoverLessons_NoFrontmatterCreatedIsEmpty(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	content := []byte("# No Frontmatter\n\nJust content.")
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "gotcha-no-fm.md"), content, 0644))

	lessons, err := DiscoverLessons(dir)
	assert.NoError(t, err)
	require.Len(t, lessons, 1)
	// Without frontmatter, Created is empty -- sorting falls back to mtime.
	assert.Empty(t, lessons[0].Created)
	assert.Equal(t, "gotcha", lessons[0].Category)
}

func TestDiscoverLessons_SkipsDirectories(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(filepath.Join(lessonsDir, "subdir"), 0755))

	lessons, err := DiscoverLessons(dir)
	assert.NoError(t, err)
	assert.Empty(t, lessons)
}

func TestDiscoverLessons_SkipsNonMdFiles(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "notes.txt"), []byte("text"), 0644))

	lessons, err := DiscoverLessons(dir)
	assert.NoError(t, err)
	assert.Empty(t, lessons)
}

func TestDiscoverLessons_MultipleLessons(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	for _, name := range []string{"gotcha-a.md", "pattern-b.md", "lesson-c.md"} {
		content := "---\ndate: 2026-01-01\ntags: [test]\n---\n"
		require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, name), []byte(content), 0644))
	}

	lessons, err := DiscoverLessons(dir)
	assert.NoError(t, err)
	assert.Len(t, lessons, 3)
}

func TestFindLessonByName_Found(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	content := "---\ndate: 2026-03-15\ntitle: Found lesson\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "pattern-found.md"), []byte(content), 0644))

	l, err := FindLessonByName(dir, "pattern-found")
	assert.NoError(t, err)
	require.NotNil(t, l)
	assert.Equal(t, "pattern-found", l.Name)
	assert.Equal(t, "pattern", l.Category)
}

func TestFindLessonByName_NotFound(t *testing.T) {
	dir := t.TempDir()

	_, err := FindLessonByName(dir, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lesson not found")
}

func TestDiscoverLessons_SortedByCreatedDescending(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	// Create three lesson files with different created dates.
	// Use lexicographically different order from date order to verify proper sorting.
	lessons := []struct {
		filename string
		created  string
	}{
		{"lesson-alpha", "2026-01-15"},
		{"lesson-beta", "2026-03-10"},
		{"lesson-gamma", "2026-02-01"},
	}

	for _, l := range lessons {
		content := fmt.Sprintf("---\ncreated: %s\ntitle: \"%s\"\n---\n", l.created, l.filename)
		require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, l.filename+".md"), []byte(content), 0644))
	}

	result, err := DiscoverLessons(dir)
	require.NoError(t, err)
	require.Len(t, result, 3)

	// Sorted newest first by created: beta (Mar 10) > gamma (Feb 1) > alpha (Jan 15).
	assert.Equal(t, "lesson-beta", result[0].Name, "first lesson should be newest by created")
	assert.Equal(t, "lesson-gamma", result[1].Name, "second lesson should be mid by created")
	assert.Equal(t, "lesson-alpha", result[2].Name, "third lesson should be oldest by created")
}

func TestDiscoverLessons_OldestSortsLast(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	oldContent := "---\ncreated: 2020-01-01\ntitle: \"old\"\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-old.md"), []byte(oldContent), 0644))

	recentContent := "---\ncreated: 2026-05-19\ntitle: \"recent\"\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-recent.md"), []byte(recentContent), 0644))

	lessons, err := DiscoverLessons(dir)
	require.NoError(t, err)
	require.Len(t, lessons, 2)

	assert.Equal(t, "lesson-recent", lessons[0].Name, "lesson with newer created date should come first")
	assert.Equal(t, "lesson-old", lessons[1].Name, "oldest lesson should come last")
}

func TestDiscoverLessons_MtimeFallbackWhenNoCreated(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	// Lessons without created/date fields -- should fall back to mtime.
	content := "---\ntitle: \"no-date\"\n---\n"

	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-old-fallback.md"), []byte(content), 0644))
	require.NoError(t, os.Chtimes(filepath.Join(lessonsDir, "lesson-old-fallback.md"), time.Time{}, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))

	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-new-fallback.md"), []byte(content), 0644))
	require.NoError(t, os.Chtimes(filepath.Join(lessonsDir, "lesson-new-fallback.md"), time.Time{}, time.Date(2026, 5, 19, 0, 0, 0, 0, time.UTC)))

	lessons, err := DiscoverLessons(dir)
	require.NoError(t, err)
	require.Len(t, lessons, 2)

	assert.Equal(t, "lesson-new-fallback", lessons[0].Name, "lesson with newer mtime should come first (fallback)")
	assert.Equal(t, "lesson-old-fallback", lessons[1].Name, "lesson with older mtime should come last (fallback)")
}

func TestDiscoverLessons_CreatedTakesPriorityOverMtime(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	// Lesson with older mtime but newer created date should sort first.
	newCreatedContent := "---\ncreated: 2026-05-01\ntitle: \"new-created\"\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-new-created.md"), []byte(newCreatedContent), 0644))
	require.NoError(t, os.Chtimes(filepath.Join(lessonsDir, "lesson-new-created.md"), time.Time{}, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))

	// Lesson with newer mtime but older created date should sort second.
	oldCreatedContent := "---\ncreated: 2026-01-01\ntitle: \"old-created\"\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-old-created.md"), []byte(oldCreatedContent), 0644))
	require.NoError(t, os.Chtimes(filepath.Join(lessonsDir, "lesson-old-created.md"), time.Time{}, time.Date(2026, 5, 19, 0, 0, 0, 0, time.UTC)))

	lessons, err := DiscoverLessons(dir)
	require.NoError(t, err)
	require.Len(t, lessons, 2)

	assert.Equal(t, "lesson-new-created", lessons[0].Name, "lesson with newer created date should sort first regardless of mtime")
	assert.Equal(t, "lesson-old-created", lessons[1].Name, "lesson with older created date should sort after")
}

func TestDiscoverLessons_CreatedField(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	content := `---
created: 2026-05-10
tags:
  - testing
title: "Created field lesson"
---
`
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-created.md"), []byte(content), 0644))

	lessons, err := DiscoverLessons(dir)
	assert.NoError(t, err)
	require.Len(t, lessons, 1)

	l := lessons[0]
	assert.Equal(t, "2026-05-10", l.Created)
	assert.Equal(t, "Created field lesson", l.Title)
}

func TestDiscoverLessons_CreatedTakesPriorityOverDate(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	content := `---
created: 2026-04-15
date: 2026-01-01
title: "Both fields"
---
`
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-both.md"), []byte(content), 0644))

	lessons, err := DiscoverLessons(dir)
	assert.NoError(t, err)
	require.Len(t, lessons, 1)

	// created should win over date
	assert.Equal(t, "2026-04-15", lessons[0].Created)
}

func TestDiscoverLessons_DateFieldStillWorks(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, "docs/lessons")
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	content := `---
date: 2026-03-20
title: "Date field only"
---
`
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-date.md"), []byte(content), 0644))

	lessons, err := DiscoverLessons(dir)
	assert.NoError(t, err)
	require.Len(t, lessons, 1)

	assert.Equal(t, "2026-03-20", lessons[0].Created)
}

func TestInferCategory(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"gotcha-something", "gotcha"},
		{"arch-design", "architecture"},
		{"pattern-reuse", "pattern"},
		{"tool-makefile", "tool"},
		{"lesson-learned", "lesson"},
		{"hook-lifecycle", "hook"},
		{"random-name", ""},
		{"gotcha", ""},          // no hyphen after prefix, doesn't match
		{"gotchasomething", ""}, // no hyphen after prefix, doesn't match
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, inferCategory(tt.name))
		})
	}
}
