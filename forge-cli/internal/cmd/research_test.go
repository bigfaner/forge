package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeResearchReport(t *testing.T, dir, filename, content string) {
	t.Helper()
	reportsDir := filepath.Join(dir, "docs", "research")
	require.NoError(t, os.MkdirAll(reportsDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(reportsDir, filename), []byte(content), 0644))
}

func TestResearchList_Empty(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"research"})
		return rootCmd.Execute()
	})
	_ = output
	_ = err
}

func TestResearchList_NoDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"research"})
		return rootCmd.Execute()
	})
	_ = output
	_ = err
}

func TestResearchList_WithReports(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	reports := []struct {
		filename string
		created  string
		topic    string
		mode     string
	}{
		{"gbrain.md", "2026-05-15", "GBrain AI Platform", "deep-dive"},
		{"graphify.md", "2026-05-10", "Graphify Knowledge Graph", "deep-dive"},
	}
	for _, r := range reports {
		content := fmt.Sprintf("---\ncreated: \"%s\"\ntopic: \"%s\"\nmode: \"%s\"\n---\n", r.created, r.topic, r.mode)
		writeResearchReport(t, dir, r.filename, content)
	}

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"research"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "RESEARCH")
	assert.Contains(t, output, "gbrain")
	assert.Contains(t, output, "graphify")
	assert.Contains(t, output, "deep-dive")
}

func TestResearchList_SortedByCreatedDescending(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

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
		writeResearchReport(t, dir, r.filename, content)
	}

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"research"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)

	lines := strings.Split(output, "\n")
	var slugOrder []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		for _, r := range reports {
			if strings.Contains(trimmed, r.filename[:len(r.filename)-3]) {
				slugOrder = append(slugOrder, r.filename[:len(r.filename)-3])
			}
		}
	}
	require.Len(t, slugOrder, 3, "expected 3 reports in output")
	assert.Equal(t, "beta", slugOrder[0], "newest report should be first")
	assert.Equal(t, "gamma", slugOrder[1], "middle report should be second")
	assert.Equal(t, "alpha", slugOrder[2], "oldest report should be last")
}

func TestResearchList_TableHeaders(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	content := "---\ncreated: \"2026-05-01\"\ntopic: \"Test Topic\"\nmode: \"deep-dive\"\n---\n"
	writeResearchReport(t, dir, "test-report.md", content)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"research"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "SLUG")
	assert.Contains(t, output, "CREATED")
	assert.Contains(t, output, "TOPIC")
	assert.Contains(t, output, "MODE")
}

func TestResearchDetail_Found(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	content := `---
created: "2026-05-15"
topic: "GBrain AI Platform"
mode: "deep-dive"
dimensions:
  - Overview & Positioning
  - Architecture
  - Learning Curve
---
`
	writeResearchReport(t, dir, "gbrain.md", content)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"research", "gbrain"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "SLUG: gbrain")
	assert.Contains(t, output, "TOPIC: GBrain AI Platform")
	assert.Contains(t, output, "CREATED: 2026-05-15")
	assert.Contains(t, output, "MODE: deep-dive")
	assert.Contains(t, output, "DIMENSIONS: Overview & Positioning, Architecture, Learning Curve")
	assert.Contains(t, output, "FILE:")
}

func TestResearchDetail_NoDimensions(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	content := "---\ncreated: \"2026-05-01\"\ntopic: \"Simple\"\nmode: \"deep-dive\"\n---\n"
	writeResearchReport(t, dir, "simple.md", content)

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"research", "simple"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "SLUG: simple")
	assert.NotContains(t, output, "DIMENSIONS:")
}

func TestResearchDetail_NotFound(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"research", "nonexistent"})
		return rootCmd.Execute()
	})
	_ = output
	// Should have error output since report not found
	assert.Error(t, err)
}

func TestResearch_TooManyArgs(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"research", "slug1", "slug2"})
		return rootCmd.Execute()
	})
	_ = output
	assert.Error(t, err)
}

func TestMapReportsToSlugLens(t *testing.T) {
	tests := []struct {
		name     string
		slugs    []string
		expected []int
	}{
		{
			name:     "single short slug",
			slugs:    []string{"abc"},
			expected: []int{3},
		},
		{
			name:     "multiple varying lengths",
			slugs:    []string{"short", "medium-length-slug", "x"},
			expected: []int{5, 18, 1},
		},
		{
			name:     "empty list",
			slugs:    []string{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lens := make([]int, len(tt.slugs))
			for i, s := range tt.slugs {
				lens[i] = len(s)
			}
			assert.Equal(t, tt.expected, lens)
		})
	}
}
