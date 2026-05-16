import { test, expect } from '@playwright/test';
import { readProjectFile } from '../helpers.js';

// -- Helpers ---------------------------------------------------------------
function fileContains(content: string, needle: string): boolean {
  return content.includes(needle);
}

// Skill/agent/command files that should use `just <verb>` exclusively.
const SKILL_FILES: string[] = [
  'plugins/forge/skills/run-e2e-tests/SKILL.md',
  'plugins/forge/agents/task-executor.md',
  'plugins/forge/commands/run-tasks.md',
  'plugins/forge/commands/fix-bug.md',
  'plugins/forge/skills/record-task/SKILL.md',
  'plugins/forge/agents/error-fixer.md',
  'plugins/forge/commands/execute-task.md',
  'plugins/forge/skills/improve-harness/SKILL.md',
];

// Raw toolchain commands that must NOT appear in skill/agent/command files.
const FORBIDDEN_COMMANDS: string[] = [
  'go test ./...',
  'go build ./...',
  'go vet ./...',
  'npm run build',
  'npm test',
  'npm test -- --coverage',
  'npx serve',
  'cargo build',
  'pytest --cov=',
  'go test -cover ./...',
  'go test -race -cover ./...',
  'npm run build && npm test',
  'cd tests/e2e && npm install',
];

// -- Tests -----------------------------------------------------------------
test.describe('Skill commands use standard just verbs', () => {

  // Traceability: TC-001 -> Story 1 / AC-1
  test('TC-001: skill/agent/command files contain zero raw toolchain commands', () => {
    const violations: string[] = [];

    for (const relPath of SKILL_FILES) {
      const content = readProjectFile(relPath);
      for (const cmd of FORBIDDEN_COMMANDS) {
        if (fileContains(content, cmd)) {
          violations.push(`${relPath} contains "${cmd}"`);
        }
      }
    }

    expect(
      violations.length,
      `Expected zero raw toolchain commands, but found:\n${violations.join('\n')}`,
    ).toBe(0);
  });
});
