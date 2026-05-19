package lesson

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscover_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, LessonsDir)
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	lessons, err := Discover(dir)
	assert.NoError(t, err)
	assert.Empty(t, lessons)
}

func TestDiscover_NoDir(t *testing.T) {
	dir := t.TempDir()

	lessons, err := Discover(dir)
	assert.NoError(t, err)
	assert.Empty(t, lessons)
}

func TestDiscover_SingleLesson(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, LessonsDir)
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

	lessons, err := Discover(dir)
	assert.NoError(t, err)
	require.Len(t, lessons, 1)

	l := lessons[0]
	assert.Equal(t, "lesson-test", l.Name)
	assert.Equal(t, "2026-05-10", l.Created)
	assert.Equal(t, "Test lesson", l.Title)
	assert.Equal(t, []string{"testing", "e2e"}, l.Tags)
	assert.Equal(t, "lesson", l.Category)
}

func TestDiscover_CategoryInference(t *testing.T) {
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
			lessonsDir := filepath.Join(dir, LessonsDir)
			require.NoError(t, os.MkdirAll(lessonsDir, 0755))

			content := "---\ndate: 2026-01-01\n---\n"
			require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, tt.filename), []byte(content), 0644))

			lessons, err := Discover(dir)
			assert.NoError(t, err)
			require.Len(t, lessons, 1)
			assert.Equal(t, tt.wantCategory, lessons[0].Category, "category for %s", tt.filename)
		})
	}
}

func TestDiscover_NoFrontmatterFallsBackToModTime(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, LessonsDir)
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	content := []byte("# No Frontmatter\n\nJust content.")
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "gotcha-no-fm.md"), content, 0644))

	lessons, err := Discover(dir)
	assert.NoError(t, err)
	require.Len(t, lessons, 1)
	assert.NotEmpty(t, lessons[0].Created)
	assert.Equal(t, "gotcha", lessons[0].Category)
}

func TestDiscover_SkipsDirectories(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, LessonsDir)
	require.NoError(t, os.MkdirAll(filepath.Join(lessonsDir, "subdir"), 0755))

	lessons, err := Discover(dir)
	assert.NoError(t, err)
	assert.Empty(t, lessons)
}

func TestDiscover_SkipsNonMdFiles(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, LessonsDir)
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "notes.txt"), []byte("text"), 0644))

	lessons, err := Discover(dir)
	assert.NoError(t, err)
	assert.Empty(t, lessons)
}

func TestDiscover_MultipleLessons(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, LessonsDir)
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	for _, name := range []string{"gotcha-a.md", "pattern-b.md", "lesson-c.md"} {
		content := "---\ndate: 2026-01-01\ntags: [test]\n---\n"
		require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, name), []byte(content), 0644))
	}

	lessons, err := Discover(dir)
	assert.NoError(t, err)
	assert.Len(t, lessons, 3)
}

func TestFindByName_Found(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, LessonsDir)
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	content := "---\ndate: 2026-03-15\ntitle: Found lesson\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "pattern-found.md"), []byte(content), 0644))

	l, err := FindByName(dir, "pattern-found")
	assert.NoError(t, err)
	require.NotNil(t, l)
	assert.Equal(t, "pattern-found", l.Name)
	assert.Equal(t, "pattern", l.Category)
}

func TestFindByName_NotFound(t *testing.T) {
	dir := t.TempDir()

	_, err := FindByName(dir, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lesson not found")
}

func TestDiscover_SortedByModTimeDescending(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, LessonsDir)
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	// Create three lesson files with specific modification times.
	fm := "---\ndate: 2026-01-01\ntitle: \"%s\"\n---\n"

	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-old.md"), []byte(fm), 0644))
	require.NoError(t, os.Chtimes(filepath.Join(lessonsDir, "lesson-old.md"), time.Time{}, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))

	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-newest.md"), []byte(fm), 0644))
	require.NoError(t, os.Chtimes(filepath.Join(lessonsDir, "lesson-newest.md"), time.Time{}, time.Date(2026, 5, 19, 0, 0, 0, 0, time.UTC)))

	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-mid.md"), []byte(fm), 0644))
	require.NoError(t, os.Chtimes(filepath.Join(lessonsDir, "lesson-mid.md"), time.Time{}, time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)))

	lessons, err := Discover(dir)
	require.NoError(t, err)
	require.Len(t, lessons, 3)

	// Sorted newest first: newest -> mid -> old.
	assert.Equal(t, "lesson-newest", lessons[0].Name, "first lesson should be newest")
	assert.Equal(t, "lesson-mid", lessons[1].Name, "second lesson should be mid")
	assert.Equal(t, "lesson-old", lessons[2].Name, "third lesson should be oldest")
}

func TestDiscover_OldestSortsLast(t *testing.T) {
	dir := t.TempDir()
	lessonsDir := filepath.Join(dir, LessonsDir)
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	fm := "---\ndate: 2026-01-01\ntitle: \"%s\"\n---\n"

	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-old.md"), []byte(fm), 0644))
	require.NoError(t, os.Chtimes(filepath.Join(lessonsDir, "lesson-old.md"), time.Time{}, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))

	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "lesson-recent.md"), []byte(fm), 0644))
	require.NoError(t, os.Chtimes(filepath.Join(lessonsDir, "lesson-recent.md"), time.Time{}, time.Date(2026, 5, 19, 0, 0, 0, 0, time.UTC)))

	lessons, err := Discover(dir)
	require.NoError(t, err)
	require.Len(t, lessons, 2)

	assert.Equal(t, "lesson-recent", lessons[0].Name, "most recently modified lesson should come first")
	assert.Equal(t, "lesson-old", lessons[1].Name, "oldest lesson should come last")
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
