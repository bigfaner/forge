//go:build e2e

package skillops

import (
	"strings"
	"testing"

	"forge-cli/tests/testkit"

	"github.com/stretchr/testify/assert"
)

// --- Plugin Content Tests (TC-001) ---
// Converted from tests/e2e/plugin-content/skill-content.spec.ts (1 test).
// Validates that skill/agent/command files contain zero raw toolchain commands.

// Skill/agent/command files that should use `just <verb>` exclusively.
var skillFiles = []string{
	"plugins/forge/skills/run-tests/SKILL.md",
	"plugins/forge/agents/task-executor.md",
	"plugins/forge/commands/run-tasks.md",
	"plugins/forge/commands/fix-bug.md",
	"plugins/forge/skills/submit-task/SKILL.md",
	"plugins/forge/commands/execute-task.md",
}

// Raw toolchain commands that must NOT appear in skill/agent/command files.
var forbiddenCommands = []string{
	"go test ./...",
	"go build ./...",
	"go vet ./...",
	"npm run build",
	"npm test",
	"npm test -- --coverage",
	"npx serve",
	"cargo build",
	"pytest --cov=",
	"go test -cover ./...",
	"go test -race -cover ./...",
	"npm run build && npm test",
	"cd tests/e2e && npm install",
}

// Traceability: TC-001 -> Story 1 / AC-1
// Skill/agent/command files contain zero raw toolchain commands.
func TestTC_001_SkillAgentCommandFilesContainZeroRawToolchainCommands(t *testing.T) {
	var violations []string

	for _, relPath := range skillFiles {
		content := testkit.ReadProjectFile(t, "../"+relPath)
		for _, cmd := range forbiddenCommands {
			if strings.Contains(content, cmd) {
				violations = append(violations, relPath+` contains "`+cmd+`"`)
			}
		}
	}

	assert.Empty(t, violations,
		"Expected zero raw toolchain commands, but found:\n%s",
		strings.Join(violations, "\n"))
}
