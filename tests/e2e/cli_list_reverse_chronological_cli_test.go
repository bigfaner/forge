//go:build e2e

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// projectRoot returns the forge project root directory.
func projectRootCLI(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine test file location")
	}
	dir := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..")
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("cannot resolve project root: %s", err)
	}
	return abs
}

// createTempProject creates a temporary directory with go.mod to act as a forge project root.
func createTempProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n\ngo 1.26\n"), 0644)
	assert.NoError(t, err, "failed to create go.mod")
	return dir
}

// createProposal creates a proposal directory and proposal.md with given frontmatter.
func createProposal(t *testing.T, projectDir, slug, frontmatter string) {
	t.Helper()
	proposalDir := filepath.Join(projectDir, "docs", "proposals", slug)
	err := os.MkdirAll(proposalDir, 0755)
	assert.NoError(t, err, "failed to create proposal directory for %s", slug)

	content := "---\n" + frontmatter + "\n---\n\n# " + slug + "\n"
	err = os.WriteFile(filepath.Join(proposalDir, "proposal.md"), []byte(content), 0644)
	assert.NoError(t, err, "failed to create proposal.md for %s", slug)
}

// createProposalNoFrontmatter creates a proposal directory and proposal.md without created field.
func createProposalNoFrontmatter(t *testing.T, projectDir, slug string) {
	t.Helper()
	proposalDir := filepath.Join(projectDir, "docs", "proposals", slug)
	err := os.MkdirAll(proposalDir, 0755)
	assert.NoError(t, err, "failed to create proposal directory for %s", slug)

	content := "---\nstatus: draft\n---\n\n# " + slug + "\n"
	err = os.WriteFile(filepath.Join(proposalDir, "proposal.md"), []byte(content), 0644)
	assert.NoError(t, err, "failed to create proposal.md for %s", slug)
}

// createFeature creates a feature directory with manifest.md.
func createFeature(t *testing.T, projectDir, slug string, mtime time.Time, withManifest bool) {
	t.Helper()
	featureDir := filepath.Join(projectDir, "docs", "features", slug)
	err := os.MkdirAll(featureDir, 0755)
	assert.NoError(t, err, "failed to create feature directory for %s", slug)

	if withManifest {
		manifestContent := "---\nstatus: active\n---\n\n# " + slug + "\n"
		manifestPath := filepath.Join(featureDir, "manifest.md")
		err = os.WriteFile(manifestPath, []byte(manifestContent), 0644)
		assert.NoError(t, err, "failed to create manifest.md for %s", slug)

		err = os.Chtimes(manifestPath, mtime, mtime)
		assert.NoError(t, err, "failed to set mtime for manifest.md of %s", slug)
	}
}

// extractSlugsFromTable parses CLI table output and returns the slug column values in order.
func extractSlugsFromTable(t *testing.T, output string) []string {
	t.Helper()
	lines := strings.Split(output, "\n")
	var slugs []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip separator lines, header lines, and block markers
		if trimmed == "" || trimmed == "---" || strings.HasPrefix(trimmed, "---") && strings.Contains(trimmed, "-") && len(strings.ReplaceAll(trimmed, "-", "")) == 0 {
			continue
		}
		// Skip header lines (contain uppercase column names)
		if strings.Contains(trimmed, "SLUG") || strings.Contains(trimmed, "FEATURES") || strings.Contains(trimmed, "PROPOSALS") {
			continue
		}
		// Skip separator lines (lines made of dashes and spaces)
		noSpaces := strings.ReplaceAll(trimmed, " ", "")
		if len(noSpaces) > 0 && strings.Count(noSpaces, "-") == len(noSpaces) {
			continue
		}
		// Extract first non-space field as slug
		fields := strings.Fields(trimmed)
		if len(fields) > 0 {
			slug := fields[0]
			// Skip lines that look like key-value pairs
			if strings.Contains(trimmed, ":") && strings.Index(trimmed, ":") < strings.Index(trimmed, slug) {
				continue
			}
			slugs = append(slugs, slug)
		}
	}
	return slugs
}

// runForge runs the forge CLI in a given working directory.
func runForgeInDir(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("forge", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("forge command failed in %s: %s: %s", dir, err, out)
	}
	return string(out)
}

// runForgeRaw runs the forge CLI in a given working directory, returning output and exit code.
func runForgeRawInDir(t *testing.T, dir string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command("forge", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return string(out), exitCode
}

// TC-001: forge proposal lists proposals sorted by created date descending
// Traceability: TC-001 -> Proposal Success Criterion "forge proposal lists proposals newest-first by created date" + Task 1 AC-1
func TestTC_001_ProposalListSortedByCreatedDescending(t *testing.T) {
	// Step 1: Create a temporary project with go.mod
	projectDir := createTempProject(t)

	// Step 2: Create proposals with created dates: 2026-01-15, 2026-03-10, 2026-02-01
	createProposal(t, projectDir, "proposal-alpha", "created: 2026-01-15\nstatus: draft")
	createProposal(t, projectDir, "proposal-beta", "created: 2026-03-10\nstatus: draft")
	createProposal(t, projectDir, "proposal-gamma", "created: 2026-02-01\nstatus: draft")

	// Step 3: Run forge proposal
	output := runForgeInDir(t, projectDir, "proposal")

	// Step 4: Parse output for slug order
	slugs := extractSlugsFromTable(t, output)

	// Expected: 2026-03-10 (newest), 2026-02-01, 2026-01-15 (oldest)
	assert.GreaterOrEqual(t, len(slugs), 3, "expected at least 3 slugs in output, got: %v", slugs)
	if len(slugs) >= 3 {
		assert.Equal(t, "proposal-beta", slugs[0], "newest proposal (2026-03-10) should be first")
		assert.Equal(t, "proposal-gamma", slugs[1], "middle proposal (2026-02-01) should be second")
		assert.Equal(t, "proposal-alpha", slugs[2], "oldest proposal (2026-01-15) should be third")
	}
}

// TC-002: forge proposal handles proposals without created frontmatter (mtime fallback)
// Traceability: TC-002 -> Task 1 AC-2 "Proposals without created frontmatter still sort correctly (fallback mtime)" + Proposal Key Scenario
func TestTC_002_ProposalListMtimeFallback(t *testing.T) {
	// Step 1: Create a temporary project with go.mod
	projectDir := createTempProject(t)

	// Step 2: Create one proposal with created: 2026-05-01 in frontmatter
	createProposal(t, projectDir, "proposal-with-date", "created: 2026-05-01\nstatus: draft")

	// Step 3: Create one proposal without created field
	createProposalNoFrontmatter(t, projectDir, "proposal-no-date")

	// Step 4: Run forge proposal
	output, exitCode := runForgeRawInDir(t, projectDir, "proposal")

	// Expected: Both proposals appear, command completes without error
	assert.Equal(t, 0, exitCode, "forge proposal should exit 0, got output:\n%s", output)
	assert.Contains(t, output, "proposal-with-date", "output should contain proposal-with-date")
	assert.Contains(t, output, "proposal-no-date", "output should contain proposal-no-date")
}

// TC-003: forge proposal with empty proposals directory
// Traceability: TC-003 -> Task 1 AC-3 (existing tests pass, empty case)
func TestTC_003_ProposalListEmptyDirectory(t *testing.T) {
	// Step 1: Create a temporary project with go.mod but no proposals
	projectDir := createTempProject(t)
	// No proposals directory created

	// Step 2: Run forge proposal
	output, exitCode := runForgeRawInDir(t, projectDir, "proposal")

	// Expected: Command outputs "no proposals found" to stderr without error exit
	assert.Equal(t, 0, exitCode, "forge proposal should exit 0 for empty proposals, got output:\n%s", output)
	assert.Contains(t, output, "no proposals found", "output should contain 'no proposals found'")
}

// TC-004: forge feature list sorts features by manifest mtime descending
// Traceability: TC-004 -> Proposal Success Criterion "forge feature list lists features newest-first by manifest mtime" + Task 2 AC-1
func TestTC_004_FeatureListSortedByMtimeDescending(t *testing.T) {
	// Step 1: Create a temporary project with go.mod
	projectDir := createTempProject(t)

	// Step 2: Create features with manifest.mtimes: 2026-01-01, 2026-03-15, 2026-05-16
	createFeature(t, projectDir, "feature-alpha", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), true)
	createFeature(t, projectDir, "feature-beta", time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), true)
	createFeature(t, projectDir, "feature-gamma", time.Date(2026, 5, 16, 0, 0, 0, 0, time.UTC), true)

	// Step 3: Run forge feature list
	output := runForgeInDir(t, projectDir, "feature", "list")

	// Step 4: Parse output for slug order
	slugs := extractSlugsFromTable(t, output)

	// Expected: 2026-05-16 (newest), 2026-03-15, 2026-01-01 (oldest)
	assert.GreaterOrEqual(t, len(slugs), 3, "expected at least 3 slugs in output, got: %v", slugs)
	if len(slugs) >= 3 {
		assert.Equal(t, "feature-gamma", slugs[0], "newest feature (2026-05-16) should be first")
		assert.Equal(t, "feature-beta", slugs[1], "middle feature (2026-03-15) should be second")
		assert.Equal(t, "feature-alpha", slugs[2], "oldest feature (2026-01-01) should be third")
	}
}

// TC-005: forge feature list sorts features with missing manifest to the end
// Traceability: TC-005 -> Task 2 AC-2 "Features with missing/unreadable manifest sort to the end"
func TestTC_005_FeatureListMissingManifestToEnd(t *testing.T) {
	// Step 1: Create a temporary project with go.mod
	projectDir := createTempProject(t)

	// Step 2: Create one feature with a valid manifest.md (recent mtime)
	createFeature(t, projectDir, "feature-with-manifest", time.Date(2026, 5, 16, 0, 0, 0, 0, time.UTC), true)

	// Step 3: Create one feature directory without manifest.md
	createFeature(t, projectDir, "feature-no-manifest", time.Time{}, false)

	// Step 4: Run forge feature list
	output := runForgeInDir(t, projectDir, "feature", "list")

	// Step 5: Compare output positions of both features
	slugs := extractSlugsFromTable(t, output)

	// Expected: Feature with manifest appears before feature without manifest
	assert.GreaterOrEqual(t, len(slugs), 2, "expected at least 2 slugs in output, got: %v", slugs)
	if len(slugs) >= 2 {
		withIdx := -1
		withoutIdx := -1
		for i, s := range slugs {
			if s == "feature-with-manifest" {
				withIdx = i
			}
			if s == "feature-no-manifest" {
				withoutIdx = i
			}
		}
		assert.NotEqual(t, -1, withIdx, "feature-with-manifest should appear in output")
		assert.NotEqual(t, -1, withoutIdx, "feature-no-manifest should appear in output")
		if withIdx != -1 && withoutIdx != -1 {
			assert.Less(t, withIdx, withoutIdx,
				"feature with manifest should appear before feature without manifest")
		}
	}
}

// TC-006: forge feature list with empty features directory
// Traceability: TC-006 -> Task 2 AC-3 (existing tests pass, empty case)
func TestTC_006_FeatureListEmptyDirectory(t *testing.T) {
	// Step 1: Create a temporary project with go.mod but no features
	projectDir := createTempProject(t)
	// No features directory created

	// Step 2: Run forge feature list
	output, exitCode := runForgeRawInDir(t, projectDir, "feature", "list")

	// Expected: Command outputs "no features found" to stderr without error exit
	assert.Equal(t, 0, exitCode, "forge feature list should exit 0 for empty features, got output:\n%s", output)
	assert.Contains(t, output, "no features found", "output should contain 'no features found'")
}
