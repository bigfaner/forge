//go:build cli_functional

package featuremanagement

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Proposal/Feature status lifecycle tests — Journey: feature-management
// Tests verify proposal list status display and feature status command.
// ==============================================================================

// createProposalPSL creates a proposal directory and proposal.md with given frontmatter.
func createProposalPSL(t *testing.T, projectDir, slug, frontmatter string) {
	t.Helper()
	proposalDir := filepath.Join(projectDir, "docs", "proposals", slug)
	err := os.MkdirAll(proposalDir, 0755)
	assert.NoError(t, err, "failed to create proposal directory for %s", slug)

	content := "---\n" + frontmatter + "\n---\n\n# " + slug + "\n"
	err = os.WriteFile(filepath.Join(proposalDir, "proposal.md"), []byte(content), 0644)
	assert.NoError(t, err, "failed to create proposal.md for %s", slug)
}

// createFeaturePSL creates a feature directory with manifest.md containing given status.
func createFeaturePSL(t *testing.T, projectDir, slug, manifestStatus string) {
	t.Helper()
	featureDir := filepath.Join(projectDir, "docs", "features", slug)
	err := os.MkdirAll(featureDir, 0755)
	assert.NoError(t, err, "failed to create feature directory for %s", slug)

	manifestContent := "---\nstatus: " + manifestStatus + "\n---\n\n# " + slug + "\n"
	manifestPath := filepath.Join(featureDir, "manifest.md")
	err = os.WriteFile(manifestPath, []byte(manifestContent), 0644)
	assert.NoError(t, err, "failed to create manifest.md for %s", slug)
}

// runForgePSL runs the forge CLI in a given working directory.
func runForgePSL(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command(testkit.ForgeBinary, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("forge command failed in %s: %s: %s", dir, err, out)
	}
	return string(out)
}

// parseBlockPSL extracts lines between "---" separators from raw CLI output.
func parseBlockPSL(t *testing.T, raw string) []string {
	t.Helper()
	return testkit.ParseBlock(t, raw)
}

// hasFieldPSL checks that a parsed block contains a "KEY: value" line.
func hasFieldPSL(lines []string, key, value string) bool {
	return testkit.HasField(lines, key, value)
}

// TC-025: Forge proposal list displays Approved status correctly
// Traceability: TC-004 -> Proposal Success Criterion 4 "forge proposal list displays Approved and Completed status values correctly"
func TestTC_025_ProposalListDisplaysApprovedStatus(t *testing.T) {
	// Step 1: Create a temporary project
	projectDir := createTempProject(t)

	// Step 2: Create a proposal with status: Approved
	createProposalPSL(t, projectDir, "test-proposal-approved", "created: 2026-05-17\nstatus: Approved")

	// Step 3: Run forge proposal
	output := runForgePSL(t, projectDir, "proposal")

	// Step 4: Verify STATUS column shows "Approved"
	lines := parseBlockPSL(t, output)
	assert.True(t, hasFieldPSL(lines, "PROPOSALS", ""), "output should have PROPOSALS header")

	// Find the row containing our proposal slug and check STATUS column
	found := false
	for _, line := range lines {
		if strings.Contains(line, "test-proposal-approved") {
			assert.Contains(t, line, "Approved", "STATUS column should show Approved for proposal with status: Approved in frontmatter")
			found = true
		}
	}
	assert.True(t, found, "proposal test-proposal-approved should appear in output")
}

// TC-026: Forge proposal list displays Completed status correctly
// Traceability: TC-004 -> Proposal Success Criterion 4 "forge proposal list displays Approved and Completed status values correctly"
func TestTC_026_ProposalListDisplaysCompletedStatus(t *testing.T) {
	// Step 1: Create a temporary project
	projectDir := createTempProject(t)

	// Step 2: Create a proposal with status: Completed
	createProposalPSL(t, projectDir, "test-proposal-completed", "created: 2026-05-16\nstatus: Completed")

	// Step 3: Run forge proposal
	output := runForgePSL(t, projectDir, "proposal")

	// Step 4: Verify STATUS column shows "Completed"
	lines := parseBlockPSL(t, output)
	found := false
	for _, line := range lines {
		if strings.Contains(line, "test-proposal-completed") {
			assert.Contains(t, line, "Completed", "STATUS column should show Completed for proposal with status: Completed in frontmatter")
			found = true
		}
	}
	assert.True(t, found, "proposal test-proposal-completed should appear in output")
}

// TC-027: Forge feature status displays Completed status correctly
// Traceability: TC-005 -> Proposal Success Criterion 5 "forge feature status <slug> correctly reflects when a feature's manifest status is completed"
func TestTC_027_FeatureStatusDisplaysCompleted(t *testing.T) {
	// Step 1: Create a temporary project
	projectDir := createTempProject(t)

	// Step 2: Create a feature with manifest status: completed
	createFeaturePSL(t, projectDir, "test-feature-completed", "completed")

	// Step 3: Run forge feature status test-feature-completed
	output := runForgePSL(t, projectDir, "feature", "status", "test-feature-completed")

	// Step 4: Verify STATUS field shows "completed"
	lines := parseBlockPSL(t, output)
	assert.True(t, hasFieldPSL(lines, "STATUS", "completed"), "STATUS field should show completed, got lines: %v", lines)
	assert.True(t, hasFieldPSL(lines, "SLUG", "test-feature-completed"), "SLUG field should show test-feature-completed")
}
