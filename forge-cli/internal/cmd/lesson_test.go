package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"forge-cli/pkg/lesson"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLessonList_Empty(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"lesson"})
		return rootCmd.Execute()
	})
	// Empty list prints to stderr
	_ = output
	_ = err
}

func TestLessonList_WithLessons(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	lessonsDir := filepath.Join(dir, lesson.LessonsDir)
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	// Create lessons
	content1 := "---\ndate: 2026-01-01\ntags: [testing]\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "gotcha-bad.md"), []byte(content1), 0644))

	content2 := "---\ndate: 2026-02-01\ntags: [arch, design]\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "pattern-reuse.md"), []byte(content2), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"lesson"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "LESSONS")
	assert.Contains(t, output, "gotcha-bad")
	assert.Contains(t, output, "pattern-reuse")
}

func TestLessonDetail_Found(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	lessonsDir := filepath.Join(dir, lesson.LessonsDir)
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	content := "---\ndate: 2026-03-15\ntags: [testing, e2e]\ntitle: My Lesson\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "hook-lifecycle.md"), []byte(content), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"lesson", "hook-lifecycle"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "NAME: hook-lifecycle")
	assert.Contains(t, output, "CREATED: 2026-03-15")
	assert.Contains(t, output, "CATEGORY: hook")
	assert.Contains(t, output, "TAGS: testing, e2e")
	assert.Contains(t, output, "TITLE: My Lesson")
	assert.Contains(t, output, "FILE:")
}

func TestLessonDetail_NoTags(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	lessonsDir := filepath.Join(dir, lesson.LessonsDir)
	require.NoError(t, os.MkdirAll(lessonsDir, 0755))

	content := "---\ndate: 2026-01-01\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(lessonsDir, "tool-makefile.md"), []byte(content), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"lesson", "tool-makefile"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "NAME: tool-makefile")
	assert.Contains(t, output, "CATEGORY: tool")
	assert.NotContains(t, output, "TAGS:")
}
