package proposal

import (
	"os"
	"path/filepath"
	"testing"

	"forge-cli/pkg/feature"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscover_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	proposalsDir := filepath.Join(dir, feature.ProposalBaseDir)
	require.NoError(t, os.MkdirAll(proposalsDir, 0755))

	proposals, err := Discover(dir)
	assert.NoError(t, err)
	assert.Empty(t, proposals)
}

func TestDiscover_NoDir(t *testing.T) {
	dir := t.TempDir()

	proposals, err := Discover(dir)
	assert.NoError(t, err)
	assert.Empty(t, proposals)
}

func TestDiscover_SingleProposal(t *testing.T) {
	dir := t.TempDir()

	// Create proposal
	slug := "test-proposal"
	proposalDir := filepath.Join(dir, feature.ProposalBaseDir, slug)
	require.NoError(t, os.MkdirAll(proposalDir, 0755))

	proposalContent := `---
created: 2026-01-15
author: tester
status: Draft
---

# Test Proposal

Some content here.
`
	require.NoError(t, os.WriteFile(filepath.Join(proposalDir, feature.ProposalFileName), []byte(proposalContent), 0644))

	proposals, err := Discover(dir)
	assert.NoError(t, err)
	require.Len(t, proposals, 1)

	p := proposals[0]
	assert.Equal(t, "test-proposal", p.Slug)
	assert.Equal(t, "2026-01-15", p.Created)
	assert.Equal(t, "Draft", p.Status)
	assert.Equal(t, "tester", p.Author)
	assert.False(t, p.HasPRD)
	assert.Equal(t, "", p.FeatureStatus)
}

func TestDiscover_ProposalWithPRD(t *testing.T) {
	dir := t.TempDir()
	slug := "my-feature"

	// Create proposal
	proposalDir := filepath.Join(dir, feature.ProposalBaseDir, slug)
	require.NoError(t, os.MkdirAll(proposalDir, 0755))
	proposalContent := `---
created: 2026-03-01
status: Approved
---
`
	require.NoError(t, os.WriteFile(filepath.Join(proposalDir, feature.ProposalFileName), []byte(proposalContent), 0644))

	// Create PRD
	prdDir := filepath.Join(dir, feature.FeaturesDir, slug, feature.PRDDirName)
	require.NoError(t, os.MkdirAll(prdDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(prdDir, feature.PRDSpecFile), []byte("# PRD"), 0644))

	proposals, err := Discover(dir)
	assert.NoError(t, err)
	require.Len(t, proposals, 1)
	assert.True(t, proposals[0].HasPRD)
}

func TestDiscover_ProposalWithFeatureStatus(t *testing.T) {
	dir := t.TempDir()
	slug := "active-feature"

	// Create proposal
	proposalDir := filepath.Join(dir, feature.ProposalBaseDir, slug)
	require.NoError(t, os.MkdirAll(proposalDir, 0755))
	proposalContent := `---
created: 2026-02-10
status: Approved
---
`
	require.NoError(t, os.WriteFile(filepath.Join(proposalDir, feature.ProposalFileName), []byte(proposalContent), 0644))

	// Create manifest
	featureDir := filepath.Join(dir, feature.FeaturesDir, slug)
	require.NoError(t, os.MkdirAll(featureDir, 0755))
	manifestContent := `---
feature: active-feature
status: in-progress
---
`
	require.NoError(t, os.WriteFile(filepath.Join(featureDir, feature.ManifestFileName), []byte(manifestContent), 0644))

	proposals, err := Discover(dir)
	assert.NoError(t, err)
	require.Len(t, proposals, 1)
	assert.Equal(t, "in-progress", proposals[0].FeatureStatus)
}

func TestDiscover_NoFrontmatterFallsBackToModTime(t *testing.T) {
	dir := t.TempDir()
	slug := "no-fm-proposal"

	proposalDir := filepath.Join(dir, feature.ProposalBaseDir, slug)
	require.NoError(t, os.MkdirAll(proposalDir, 0755))

	content := []byte("# No Frontmatter Proposal\n\nJust content.")
	require.NoError(t, os.WriteFile(filepath.Join(proposalDir, feature.ProposalFileName), content, 0644))

	proposals, err := Discover(dir)
	assert.NoError(t, err)
	require.Len(t, proposals, 1)

	// Should have a date from mod time
	assert.NotEmpty(t, proposals[0].Created)
	// Status should be empty since no frontmatter
	assert.Equal(t, "", proposals[0].Status)
}

func TestDiscover_MultipleProposals(t *testing.T) {
	dir := t.TempDir()

	for _, slug := range []string{"alpha", "beta", "gamma"} {
		proposalDir := filepath.Join(dir, feature.ProposalBaseDir, slug)
		require.NoError(t, os.MkdirAll(proposalDir, 0755))
		content := "---\ncreated: 2026-01-01\nstatus: Draft\n---\n"
		require.NoError(t, os.WriteFile(filepath.Join(proposalDir, feature.ProposalFileName), []byte(content), 0644))
	}

	proposals, err := Discover(dir)
	assert.NoError(t, err)
	assert.Len(t, proposals, 3)
}

func TestDiscover_SkipsNonDirectoryEntries(t *testing.T) {
	dir := t.TempDir()

	proposalsDir := filepath.Join(dir, feature.ProposalBaseDir)
	require.NoError(t, os.MkdirAll(proposalsDir, 0755))

	// Create a file (not directory) in proposals dir
	require.NoError(t, os.WriteFile(filepath.Join(proposalsDir, "README.md"), []byte("# Proposals"), 0644))

	proposals, err := Discover(dir)
	assert.NoError(t, err)
	assert.Empty(t, proposals)
}

func TestFindBySlug_Found(t *testing.T) {
	dir := t.TempDir()
	slug := "find-me"

	proposalDir := filepath.Join(dir, feature.ProposalBaseDir, slug)
	require.NoError(t, os.MkdirAll(proposalDir, 0755))
	content := "---\ncreated: 2026-05-01\nstatus: Draft\n---\n"
	require.NoError(t, os.WriteFile(filepath.Join(proposalDir, feature.ProposalFileName), []byte(content), 0644))

	p, err := FindBySlug(dir, "find-me")
	assert.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, "find-me", p.Slug)
}

func TestFindBySlug_NotFound(t *testing.T) {
	dir := t.TempDir()

	_, err := FindBySlug(dir, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "proposal not found")
}
