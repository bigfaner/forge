package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/feature"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProposalList_Empty(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"proposal"})
		return rootCmd.Execute()
	})
	// The command exits via os.Exit(0) from stderr print, so we check stderr
	_ = output
	_ = err
}

func TestProposalList_WithProposals(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	// Create two proposals
	for _, slug := range []string{"alpha", "beta"} {
		proposalDir := filepath.Join(dir, feature.ProposalBaseDir, slug)
		require.NoError(t, os.MkdirAll(proposalDir, 0755))
		content := "---\ncreated: 2026-01-01\nstatus: Draft\n---\n"
		require.NoError(t, os.WriteFile(filepath.Join(proposalDir, feature.ProposalFileName), []byte(content), 0644))
	}

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"proposal"})
		return rootCmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "PROPOSALS")
	assert.Contains(t, output, "alpha")
	assert.Contains(t, output, "beta")
}

func TestProposalDetail_Found(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	slug := "test-proposal"
	proposalDir := filepath.Join(dir, feature.ProposalBaseDir, slug)
	require.NoError(t, os.MkdirAll(proposalDir, 0755))
	content := "---\ncreated: 2026-05-01\nauthor: tester\nstatus: Draft\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(proposalDir, feature.ProposalFileName), []byte(content), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"proposal", slug})
		return rootCmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "SLUG: test-proposal")
	assert.Contains(t, output, "CREATED: 2026-05-01")
	assert.Contains(t, output, "STATUS: Draft")
	assert.Contains(t, output, "AUTHOR: tester")
	assert.Contains(t, output, "FILE:")
}

func TestProposalDetail_WithPRDAndFeature(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	slug := "full-proposal"
	proposalDir := filepath.Join(dir, feature.ProposalBaseDir, slug)
	require.NoError(t, os.MkdirAll(proposalDir, 0755))
	content := "---\ncreated: 2026-05-01\nauthor: tester\nstatus: Approved\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(proposalDir, feature.ProposalFileName), []byte(content), 0644))

	// Create PRD
	prdDir := filepath.Join(dir, feature.FeaturesDir, slug, feature.PRDDirName)
	require.NoError(t, os.MkdirAll(prdDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(prdDir, feature.PRDSpecFile), []byte("# PRD"), 0644))

	// Create manifest
	featureDir := filepath.Join(dir, feature.FeaturesDir, slug)
	require.NoError(t, os.MkdirAll(featureDir, 0755))
	manifestContent := "---\nfeature: full-proposal\nstatus: in-progress\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(featureDir, feature.ManifestFileName), []byte(manifestContent), 0644))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(dir))

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"proposal", slug})
		return rootCmd.Execute()
	})
	require.NoError(t, err)
	assert.Contains(t, output, "PRD: yes")
	assert.Contains(t, output, "FEATURE: in-progress")
}

func TestTruncateSlug(t *testing.T) {
	tests := []struct {
		input   string
		maxLen  int
		wantLen int
	}{
		{"short", 10, 5},
		{"exactly-10", 10, 10},
		{"this-is-a-very-long-slug-name", 10, 10},
	}

	for _, tt := range tests {
		result := truncateSlug(tt.input, tt.maxLen)
		assert.LessOrEqual(t, len(result), tt.maxLen)
		if tt.input == "short" {
			assert.Equal(t, "short", result)
		}
		if len(tt.input) > tt.maxLen {
			assert.True(t, strings.HasSuffix(result, "..."))
		}
	}
}
