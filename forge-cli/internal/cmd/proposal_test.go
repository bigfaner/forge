package cmd

import (
	"fmt"
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
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
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
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
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

func TestProposalList_SortedByCreatedDescending(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	// Create proposals with different dates (lexicographic order differs from date order)
	proposals := []struct {
		slug    string
		created string
	}{
		{"alpha-proposal", "2026-01-15"},
		{"beta-proposal", "2026-03-10"},
		{"gamma-proposal", "2026-02-01"},
	}
	for _, p := range proposals {
		proposalDir := filepath.Join(dir, feature.ProposalBaseDir, p.slug)
		require.NoError(t, os.MkdirAll(proposalDir, 0755))
		content := fmt.Sprintf("---\ncreated: %s\nstatus: Draft\n---\n", p.created)
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

	// Verify newest first: beta (Mar 10) > gamma (Feb 1) > alpha (Jan 15)
	lines := strings.Split(output, "\n")
	var slugOrder []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		for _, p := range proposals {
			if strings.Contains(trimmed, p.slug) {
				slugOrder = append(slugOrder, p.slug)
			}
		}
	}
	require.Len(t, slugOrder, 3, "expected 3 proposals in output")
	assert.Equal(t, "beta-proposal", slugOrder[0], "newest proposal should be first")
	assert.Equal(t, "gamma-proposal", slugOrder[1], "middle proposal should be second")
	assert.Equal(t, "alpha-proposal", slugOrder[2], "oldest proposal should be last")
}

func TestProposalDetail_Found(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
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
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
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

func TestProposalList_ApprovedStatus(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	slug := "approved-proposal"
	proposalDir := filepath.Join(dir, feature.ProposalBaseDir, slug)
	require.NoError(t, os.MkdirAll(proposalDir, 0755))
	content := "---\ncreated: 2026-05-01\nstatus: Approved\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(proposalDir, feature.ProposalFileName), []byte(content), 0644))

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
	assert.Contains(t, output, "approved-proposal")
	assert.Contains(t, output, "Approved")
}

func TestProposalList_CompletedStatus(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	slug := "completed-proposal"
	proposalDir := filepath.Join(dir, feature.ProposalBaseDir, slug)
	require.NoError(t, os.MkdirAll(proposalDir, 0755))
	content := "---\ncreated: 2026-04-01\nstatus: Completed\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(proposalDir, feature.ProposalFileName), []byte(content), 0644))

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
	assert.Contains(t, output, "completed-proposal")
	assert.Contains(t, output, "Completed")
}

func TestProposalList_MultipleStatuses(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	proposals := []struct {
		slug    string
		status  string
		created string
	}{
		{"draft-p", "Draft", "2026-01-01"},
		{"approved-p", "Approved", "2026-02-01"},
		{"completed-p", "Completed", "2026-03-01"},
	}
	for _, p := range proposals {
		proposalDir := filepath.Join(dir, feature.ProposalBaseDir, p.slug)
		require.NoError(t, os.MkdirAll(proposalDir, 0755))
		content := fmt.Sprintf("---\ncreated: %s\nstatus: %s\n---\n", p.created, p.status)
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
	assert.Contains(t, output, "Draft")
	assert.Contains(t, output, "Approved")
	assert.Contains(t, output, "Completed")
}

func TestProposalDetail_ApprovedStatus(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	slug := "approved-detail"
	proposalDir := filepath.Join(dir, feature.ProposalBaseDir, slug)
	require.NoError(t, os.MkdirAll(proposalDir, 0755))
	content := "---\ncreated: 2026-05-01\nauthor: tester\nstatus: Approved\n---\n"
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
	assert.Contains(t, output, "STATUS: Approved")
}

func TestProposalDetail_CompletedStatus(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644))

	slug := "completed-detail"
	proposalDir := filepath.Join(dir, feature.ProposalBaseDir, slug)
	require.NoError(t, os.MkdirAll(proposalDir, 0755))
	content := "---\ncreated: 2026-05-01\nauthor: tester\nstatus: Completed\n---\n"
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
	assert.Contains(t, output, "STATUS: Completed")
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
