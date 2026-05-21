//go:build e2e

package skillops

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Clean Code Skill Tests ---
// Validates that the clean-code skill and command files exist and contain
// the required content per task 1 acceptance criteria.

// cleanCodeRepoRoot resolves the repository root.
func cleanCodeRepoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "plugins")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("cannot find repo root")
		}
		dir = parent
	}
}

// Traceability: Task 1 / AC-1 — SKILL.md exists with complete skill definition
func TestCleanCode_SkillFile_ExistsAndHasRequiredStructure(t *testing.T) {
	root := cleanCodeRepoRoot(t)
	skillPath := filepath.Join(root, "plugins", "forge", "skills", "clean-code", "SKILL.md")
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err, "plugins/forge/skills/clean-code/SKILL.md must exist")

	content := string(data)

	// Frontmatter must have name and description
	assert.True(t, strings.Contains(content, "name:"), "SKILL.md must have frontmatter 'name' field")
	assert.True(t, strings.Contains(content, "description:"), "SKILL.md must have frontmatter 'description' field")
}

// Traceability: Task 1 / AC-2 — Skill workflow: scope detection -> cleanup -> quality gate -> summary
func TestCleanCode_SkillFile_ContainsWorkflow(t *testing.T) {
	root := cleanCodeRepoRoot(t)
	skillPath := filepath.Join(root, "plugins", "forge", "skills", "clean-code", "SKILL.md")
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)

	content := string(data)

	// Must describe scope detection via git diff
	assert.True(t, strings.Contains(content, "git diff"), "SKILL.md must describe git diff scope detection")

	// Must mention quality gate
	assert.True(t, strings.Contains(content, "quality gate") || strings.Contains(content, "Quality Gate") || strings.Contains(content, "just test"),
		"SKILL.md must mention quality gate or just test")

	// Must mention cleanup summary
	assert.True(t, strings.Contains(content, "summary") || strings.Contains(content, "Summary"),
		"SKILL.md must mention cleanup summary")
}

// Traceability: Task 1 / AC-3 — Cleanup logic follows code-simplifier 5 principles
func TestCleanCode_SkillFile_ContainsFivePrinciples(t *testing.T) {
	root := cleanCodeRepoRoot(t)
	skillPath := filepath.Join(root, "plugins", "forge", "skills", "clean-code", "SKILL.md")
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)

	content := string(data)

	principles := []string{
		"Preserve Functionality",
		"Apply Project Standards",
		"Enhance Clarity",
		"Maintain Balance",
		"Focus Scope",
	}
	for _, p := range principles {
		assert.True(t, strings.Contains(content, p),
			"SKILL.md must mention principle: %s", p)
	}
}

// Traceability: Task 1 / AC-5 — Command file exists as slash command entry point
func TestCleanCode_CommandFile_ExistsAndReferencesSkill(t *testing.T) {
	root := cleanCodeRepoRoot(t)
	cmdPath := filepath.Join(root, "plugins", "forge", "commands", "clean-code.md")
	data, err := os.ReadFile(cmdPath)
	require.NoError(t, err, "plugins/forge/commands/clean-code.md must exist")

	content := string(data)

	// Command must reference the skill invocation
	assert.True(t, strings.Contains(content, "forge:clean-code"),
		"command file must reference 'forge:clean-code' skill invocation")

	// Frontmatter must have name and description
	assert.True(t, strings.Contains(content, "name:"), "command must have frontmatter 'name' field")
	assert.True(t, strings.Contains(content, "description:"), "command must have frontmatter 'description' field")
}

// Traceability: Task 1 / Hard Rule — Skill only modifies git diff scope
func TestCleanCode_SkillFile_EnforcesScopeConstraint(t *testing.T) {
	root := cleanCodeRepoRoot(t)
	skillPath := filepath.Join(root, "plugins", "forge", "skills", "clean-code", "SKILL.md")
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)

	content := string(data)

	// Must enforce scope constraint (only modify files in diff scope)
	lower := strings.ToLower(content)
	assert.True(t,
		strings.Contains(lower, "scope") || strings.Contains(lower, "only modify") || strings.Contains(lower, "only files"),
		"SKILL.md must enforce scope constraint")
}
