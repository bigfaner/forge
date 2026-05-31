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

// helper to write a research report file in the test's docs/research directory.
func writeReport(t *testing.T, dir, filename, content string) {
	t.Helper()
	reportsDir := filepath.Join(dir, "docs", "research")
	require.NoError(t, os.MkdirAll(reportsDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(reportsDir, filename), []byte(content), 0644))
}

func TestDiscoverReports_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	reportsDir := filepath.Join(dir, "docs", "research")
	require.NoError(t, os.MkdirAll(reportsDir, 0755))

	reports, err := DiscoverReports(dir)
	assert.NoError(t, err)
	assert.Empty(t, reports)
}

func TestDiscoverReports_NoDir(t *testing.T) {
	dir := t.TempDir()

	reports, err := DiscoverReports(dir)
	assert.NoError(t, err)
	assert.Empty(t, reports)
}

func TestDiscoverReports_SingleReport(t *testing.T) {
	dir := t.TempDir()

	content := `---
created: "2026-05-10"
topic: "Claude Code Plugins"
mode: "deep-dive"
dimensions:
  - extensibility
  - security
candidates:
  - claude-code
---

# Research Report

Some content here.
`
	writeReport(t, dir, "claude-plugins.md", content)

	reports, err := DiscoverReports(dir)
	assert.NoError(t, err)
	require.Len(t, reports, 1)

	r := reports[0]
	assert.Equal(t, "claude-plugins", r.Slug)
	assert.Equal(t, "2026-05-10", r.Created)
	assert.Equal(t, "Claude Code Plugins", r.Topic)
	assert.Equal(t, "deep-dive", r.Mode)
	assert.Equal(t, []string{"extensibility", "security"}, r.Dimensions)
	assert.Equal(t, []string{"claude-code"}, r.Candidates)
	assert.Contains(t, r.FilePath, filepath.Join("docs", "research", "claude-plugins.md"))
}

func TestDiscoverReports_ComparisonReport(t *testing.T) {
	dir := t.TempDir()

	content := `---
created: "2026-03-15"
topic: "Build Tools Comparison"
mode: "comparison"
dimensions:
  - speed
  - ecosystem
candidates:
  - webpack
  - vite
  - esbuild
---

# Comparison
`
	writeReport(t, dir, "build-tools.md", content)

	reports, err := DiscoverReports(dir)
	assert.NoError(t, err)
	require.Len(t, reports, 1)

	r := reports[0]
	assert.Equal(t, "comparison", r.Mode)
	assert.Equal(t, []string{"webpack", "vite", "esbuild"}, r.Candidates)
	assert.Equal(t, []string{"speed", "ecosystem"}, r.Dimensions)
}

func TestDiscoverReports_MultipleReports(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"alpha.md", "beta.md", "gamma.md"} {
		content := fmt.Sprintf("---\ncreated: \"2026-01-01\"\ntopic: \"%s\"\nmode: \"deep-dive\"\n---\n", name)
		writeReport(t, dir, name, content)
	}

	reports, err := DiscoverReports(dir)
	assert.NoError(t, err)
	assert.Len(t, reports, 3)
}

func TestDiscoverReports_SkipsDirectories(t *testing.T) {
	dir := t.TempDir()
	reportsDir := filepath.Join(dir, "docs", "research")
	require.NoError(t, os.MkdirAll(filepath.Join(reportsDir, "subdir"), 0755))

	reports, err := DiscoverReports(dir)
	assert.NoError(t, err)
	assert.Empty(t, reports)
}

func TestDiscoverReports_SkipsNonMdFiles(t *testing.T) {
	dir := t.TempDir()
	reportsDir := filepath.Join(dir, "docs", "research")
	require.NoError(t, os.MkdirAll(reportsDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(reportsDir, "notes.txt"), []byte("text"), 0644))

	reports, err := DiscoverReports(dir)
	assert.NoError(t, err)
	assert.Empty(t, reports)
}

func TestDiscoverReports_SkipsNoFrontmatter(t *testing.T) {
	dir := t.TempDir()

	content := "# No Frontmatter\n\nJust content."
	writeReport(t, dir, "no-frontmatter.md", content)

	reports, err := DiscoverReports(dir)
	assert.NoError(t, err)
	assert.Empty(t, reports, "files without frontmatter should be skipped")
}

func TestDiscoverReports_SkipsMalformedFrontmatter(t *testing.T) {
	dir := t.TempDir()

	// Only opening ---, no closing ---
	content := "---\ncreated: 2026-01-01\ntopic: broken\n"
	writeReport(t, dir, "malformed.md", content)

	reports, err := DiscoverReports(dir)
	assert.NoError(t, err)
	assert.Empty(t, reports, "malformed frontmatter should be skipped")
}

func TestDiscoverReports_OptionalDimensionsAndCandidates(t *testing.T) {
	dir := t.TempDir()

	content := `---
created: "2026-04-01"
topic: "Minimal Report"
mode: "deep-dive"
---

# Minimal
`
	writeReport(t, dir, "minimal.md", content)

	reports, err := DiscoverReports(dir)
	assert.NoError(t, err)
	require.Len(t, reports, 1)

	r := reports[0]
	assert.Nil(t, r.Dimensions, "missing dimensions should be nil")
	assert.Nil(t, r.Candidates, "missing candidates should be nil")
}

func TestDiscoverReports_CreatedFallsBackToMtime(t *testing.T) {
	dir := t.TempDir()

	content := `---
topic: "No Created Field"
mode: "deep-dive"
---

# No date
`
	writeReport(t, dir, "no-created.md", content)

	// Set mtime explicitly
	filePath := filepath.Join(dir, "docs", "research", "no-created.md")
	fixedTime := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	require.NoError(t, os.Chtimes(filePath, time.Time{}, fixedTime))

	reports, err := DiscoverReports(dir)
	assert.NoError(t, err)
	require.Len(t, reports, 1)

	// Created field stays empty when not in frontmatter; mtime is used for sorting only.
	assert.Equal(t, "", reports[0].Created, "created should be empty when not in frontmatter")
}

func TestDiscoverReports_SortedByCreatedDescending(t *testing.T) {
	dir := t.TempDir()

	reports := []struct {
		filename string
		created  string
	}{
		{"alpha.md", "2026-01-15"},
		{"beta.md", "2026-03-10"},
		{"gamma.md", "2026-02-01"},
	}

	for _, r := range reports {
		content := fmt.Sprintf("---\ncreated: \"%s\"\ntopic: \"%s\"\nmode: \"deep-dive\"\n---\n", r.created, r.filename)
		writeReport(t, dir, r.filename, content)
	}

	result, err := DiscoverReports(dir)
	require.NoError(t, err)
	require.Len(t, result, 3)

	// Sorted newest first: beta (Mar 10) > gamma (Feb 1) > alpha (Jan 15)
	assert.Equal(t, "beta.md", result[0].Topic, "first report should be newest by created")
	assert.Equal(t, "gamma.md", result[1].Topic, "second report should be mid by created")
	assert.Equal(t, "alpha.md", result[2].Topic, "third report should be oldest by created")
}

func TestDiscoverReports_CreatedTakesPriorityOverMtime(t *testing.T) {
	dir := t.TempDir()

	// Report with older mtime but newer created date
	newCreatedContent := "---\ncreated: \"2026-05-01\"\ntopic: \"new-created\"\nmode: \"deep-dive\"\n---\n"
	writeReport(t, dir, "new-created.md", newCreatedContent)
	require.NoError(t, os.Chtimes(
		filepath.Join(dir, "docs", "research", "new-created.md"),
		time.Time{}, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	))

	// Report with newer mtime but older created date
	oldCreatedContent := "---\ncreated: \"2026-01-01\"\ntopic: \"old-created\"\nmode: \"deep-dive\"\n---\n"
	writeReport(t, dir, "old-created.md", oldCreatedContent)
	require.NoError(t, os.Chtimes(
		filepath.Join(dir, "docs", "research", "old-created.md"),
		time.Time{}, time.Date(2026, 5, 19, 0, 0, 0, 0, time.UTC),
	))

	reports, err := DiscoverReports(dir)
	require.NoError(t, err)
	require.Len(t, reports, 2)

	assert.Equal(t, "new-created", reports[0].Slug, "report with newer created date should sort first")
	assert.Equal(t, "old-created", reports[1].Slug, "report with older created date should sort after")
}

func TestDiscoverReports_MtimeFallbackSorting(t *testing.T) {
	dir := t.TempDir()

	// Two reports without created date -- sorting falls back to mtime
	content := "---\ntopic: \"no-date-a\"\nmode: \"deep-dive\"\n---\n"
	writeReport(t, dir, "old-fallback.md", content)
	require.NoError(t, os.Chtimes(
		filepath.Join(dir, "docs", "research", "old-fallback.md"),
		time.Time{}, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	))

	writeReport(t, dir, "new-fallback.md", content)
	require.NoError(t, os.Chtimes(
		filepath.Join(dir, "docs", "research", "new-fallback.md"),
		time.Time{}, time.Date(2026, 5, 19, 0, 0, 0, 0, time.UTC),
	))

	reports, err := DiscoverReports(dir)
	require.NoError(t, err)
	require.Len(t, reports, 2)

	assert.Equal(t, "new-fallback", reports[0].Slug, "report with newer mtime should come first (fallback)")
	assert.Equal(t, "old-fallback", reports[1].Slug, "report with older mtime should come last (fallback)")
}

func TestFindReportBySlug_Found(t *testing.T) {
	dir := t.TempDir()

	content := "---\ncreated: \"2026-05-01\"\ntopic: \"Find Me\"\nmode: \"deep-dive\"\n---\n"
	writeReport(t, dir, "find-me.md", content)

	r, err := FindReportBySlug(dir, "find-me")
	assert.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "find-me", r.Slug)
	assert.Equal(t, "Find Me", r.Topic)
}

func TestFindReportBySlug_NotFound(t *testing.T) {
	dir := t.TempDir()

	_, err := FindReportBySlug(dir, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "research report not found")
}

func TestDiscoverReports_InvalidYamlInFrontmatter(t *testing.T) {
	dir := t.TempDir()

	// Valid frontmatter delimiters but invalid YAML content
	content := "---\ncreated: [invalid: yaml\n---\n"
	writeReport(t, dir, "bad-yaml.md", content)

	reports, err := DiscoverReports(dir)
	assert.NoError(t, err)
	assert.Empty(t, reports, "invalid YAML within frontmatter should be skipped")
}

func TestDiscoverReports_ReportWithOnlyMode(t *testing.T) {
	dir := t.TempDir()

	// Has mode but no topic -- should still be included since mode is set
	content := "---\nmode: \"deep-dive\"\n---\n"
	writeReport(t, dir, "mode-only.md", content)

	reports, err := DiscoverReports(dir)
	assert.NoError(t, err)
	require.Len(t, reports, 1)
	assert.Equal(t, "deep-dive", reports[0].Mode)
	assert.Equal(t, "", reports[0].Topic)
}

func TestDiscoverReports_ReportWithOnlyTopic(t *testing.T) {
	dir := t.TempDir()

	// Has topic but no mode -- should still be included since topic is set
	content := "---\ntopic: \"Just Topic\"\n---\n"
	writeReport(t, dir, "topic-only.md", content)

	reports, err := DiscoverReports(dir)
	assert.NoError(t, err)
	require.Len(t, reports, 1)
	assert.Equal(t, "Just Topic", reports[0].Topic)
	assert.Equal(t, "", reports[0].Mode)
}
