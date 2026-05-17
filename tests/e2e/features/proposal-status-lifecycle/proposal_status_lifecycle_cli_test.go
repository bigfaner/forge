//go:build e2e

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createTempProjectPSL creates a temporary directory with go.mod to act as a forge project root.
func createTempProjectPSL(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n\ngo 1.26\n"), 0644)
	assert.NoError(t, err, "failed to create go.mod")
	return dir
}

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
	cmd := exec.Command("forge", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("forge command failed in %s: %s: %s", dir, err, out)
	}
	return string(out)
}

// parseBlockPSL extracts lines between "---" separators from raw CLI output.
func parseBlockPSL(t *testing.T, raw string) []string {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) < 2 || strings.TrimSpace(lines[0]) != "---" || strings.TrimSpace(lines[len(lines)-1]) != "---" {
		t.Fatalf("output must be wrapped in --- separators, got:\n%s", raw)
	}
	inner := lines[1 : len(lines)-1]
	result := make([]string, 0, len(inner))
	for _, l := range inner {
		result = append(result, strings.TrimSpace(l))
	}
	return result
}

// hasFieldPSL checks that a parsed block contains a "KEY: value" line.
func hasFieldPSL(lines []string, key, value string) bool {
	prefix := key + ": "
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			if value == "" {
				return true
			}
			return l == prefix+value
		}
	}
	return false
}

// TC-004: Forge proposal list displays Approved status correctly
// Traceability: TC-004 -> Proposal Success Criterion 4 "forge proposal list displays Approved and Completed status values correctly"
func TestTC_004_ProposalListDisplaysApprovedStatus(t *testing.T) {
	// Step 1: Create a temporary project
	projectDir := createTempProjectPSL(t)

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

// TC-004 extended: Forge proposal list displays Completed status correctly
// Traceability: TC-004 -> Proposal Success Criterion 4 "forge proposal list displays Approved and Completed status values correctly"
func TestTC_004b_ProposalListDisplaysCompletedStatus(t *testing.T) {
	// Step 1: Create a temporary project
	projectDir := createTempProjectPSL(t)

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

// TC-005: Forge feature status displays Completed status correctly
// Traceability: TC-005 -> Proposal Success Criterion 5 "forge feature status <slug> correctly reflects when a feature's manifest status is completed"
func TestTC_005_FeatureStatusDisplaysCompleted(t *testing.T) {
	// Step 1: Create a temporary project
	projectDir := createTempProjectPSL(t)

	// Step 2: Create a feature with manifest status: completed
	createFeaturePSL(t, projectDir, "test-feature-completed", "completed")

	// Step 3: Run forge feature status test-feature-completed
	output := runForgePSL(t, projectDir, "feature", "status", "test-feature-completed")

	// Step 4: Verify STATUS field shows "completed"
	lines := parseBlockPSL(t, output)
	assert.True(t, hasFieldPSL(lines, "STATUS", "completed"), "STATUS field should show completed, got lines: %v", lines)
	assert.True(t, hasFieldPSL(lines, "SLUG", "test-feature-completed"), "SLUG field should show test-feature-completed")
}
