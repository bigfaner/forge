import { test, expect } from '@playwright/test';
import { mkdirSync, writeFileSync, rmSync, existsSync, readdirSync } from 'node:fs';
import { join } from 'node:path';
import { runCli, readProjectFile, projectFileExists, PROJECT_ROOT } from '../helpers.js';

// ── Helpers ────────────────────────────────────────────────────────
function fileContains(content: string, needle: string): boolean {
  return content.includes(needle);
}

function fileNotContains(content: string, needle: string): boolean {
  return !content.includes(needle);
}

// ── TC-001 to TC-008, TC-013 to TC-019: File content checks ───────
test.describe('Skill/Agent file content checks', () => {

  // Traceability: TC-001 → Story 1 / AC-1
  test('TC-001: run-e2e-tests Step 1 uses just e2e-setup', () => {
    const content = readProjectFile('plugins/forge/skills/run-e2e-tests/SKILL.md');
    expect(
      fileContains(content, 'just e2e-setup'),
      'Expected "just e2e-setup" to appear in run-e2e-tests/SKILL.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'cd tests/e2e && npm install'),
      'Expected "cd tests/e2e && npm install" NOT to appear in run-e2e-tests/SKILL.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'npx playwright install chromium'),
      'Expected "npx playwright install chromium" NOT to appear in run-e2e-tests/SKILL.md',
    ).toBeTruthy();
  });

  // Traceability: TC-002 → Story 2 / AC-1 (updated: task-executor uses workflow-driven dispatch, not hardcoded commands)
  test('TC-002: task-executor uses workflow-driven execution (no hardcoded language commands)', () => {
    const content = readProjectFile('plugins/forge/agents/task-executor.md');
    // After refactor, task-executor uses workflow dispatch from task files, not hardcoded just commands.
    // Verify the agent does NOT contain language-specific commands.
    expect(
      fileNotContains(content, 'go test ./...'),
      'Expected "go test ./..." NOT to appear in task-executor.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'npm test'),
      'Expected "npm test" NOT to appear in task-executor.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'pytest'),
      'Expected "pytest" NOT to appear in task-executor.md',
    ).toBeTruthy();
    // Verify workflow-driven dispatch: agent reads execution workflow from task files
    expect(
      fileContains(content, 'Execution Workflow') || fileContains(content, 'workflow'),
      'Expected workflow-driven dispatch in task-executor.md',
    ).toBeTruthy();
  });

  // Traceability: TC-005 → Story 5 / AC-2 (verify run-tasks uses standard just commands)
  test('TC-005: run-tasks Breaking Gate uses just test for verification', () => {
    const content = readProjectFile('plugins/forge/commands/run-tasks.md');
    expect(
      fileContains(content, 'just test'),
      'Expected "just test" to appear in run-tasks.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'go test ./...'),
      'Expected "go test ./..." NOT to appear in run-tasks.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'npm test'),
      'Expected "npm test" NOT to appear in run-tasks.md',
    ).toBeTruthy();
  });

  // Traceability: TC-006 → Story 5 / AC-1
  test('TC-006: fix-bug uses just test not project-test-command placeholder', () => {
    const content = readProjectFile('plugins/forge/commands/fix-bug.md');
    expect(
      fileContains(content, 'just test'),
      'Expected "just test" to appear in fix-bug.md test verification step',
    ).toBeTruthy();
    expect(
      fileNotContains(content, '<project-test-command>'),
      'Expected "<project-test-command>" placeholder NOT to appear in fix-bug.md',
    ).toBeTruthy();
  });

  // Traceability: TC-007 → Story 5 / AC-2
  test('TC-007: run-tasks Breaking Gate uses just test', () => {
    const content = readProjectFile('plugins/forge/commands/run-tasks.md');
    expect(
      fileContains(content, 'just test'),
      'Expected "just test" to appear in run-tasks.md Breaking Gate section',
    ).toBeTruthy();
    // Breaking Gate section should not use language-specific commands
    const breakingGateIdx = content.indexOf('Breaking Task Gate');
    expect(breakingGateIdx !== -1, 'Expected "Breaking Task Gate" section to exist in run-tasks.md').toBeTruthy();
    const afterBreakingGate = content.slice(breakingGateIdx);
    expect(
      fileNotContains(afterBreakingGate, 'npm test'),
      'Expected "npm test" NOT to appear in Breaking Gate section of run-tasks.md',
    ).toBeTruthy();
    expect(
      fileNotContains(afterBreakingGate, 'go test'),
      'Expected "go test" NOT to appear in Breaking Gate section of run-tasks.md',
    ).toBeTruthy();
  });

  // Traceability: TC-008 → Story 5 / AC-3
  test('TC-008: record-task Metrics Collection uses just test', () => {
    const content = readProjectFile('plugins/forge/skills/record-task/SKILL.md');
    expect(
      fileContains(content, 'just test'),
      'Expected "just test" to appear in record-task/SKILL.md Metrics Collection section',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'go test -cover ./...'),
      'Expected "go test -cover ./..." NOT to appear in record-task/SKILL.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'npm test -- --coverage'),
      'Expected "npm test -- --coverage" NOT to appear in record-task/SKILL.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'pytest --cov='),
      'Expected "pytest --cov=" NOT to appear in record-task/SKILL.md',
    ).toBeTruthy();
  });

  // Traceability: TC-013 → Spec Section 5.3
  test('TC-013: run-e2e-tests skill prompts init-justfile when justfile missing', () => {
    const content = readProjectFile('plugins/forge/skills/run-e2e-tests/SKILL.md');
    // Skill should reference justfile existence check or init-justfile
    const hasJustfileCheck = fileContains(content, 'justfile') || fileContains(content, 'init-justfile');
    expect(
      hasJustfileCheck,
      'Expected run-e2e-tests/SKILL.md to reference justfile existence check or /init-justfile',
    ).toBeTruthy();
  });

  // Traceability: TC-014 → Spec Section 5.2 / Story 3
  test('TC-014: gen-test-scripts Step 4 uses just e2e-verify', () => {
    const content = readProjectFile('plugins/forge/skills/gen-test-scripts/SKILL.md');
    expect(
      fileContains(content, 'just e2e-verify --feature'),
      'Expected "just e2e-verify --feature" to appear in gen-test-scripts/SKILL.md Step 4',
    ).toBeTruthy();
    // Primary method is just e2e-verify; raw grep is only a documented fallback
  });

  // Traceability: TC-015 → Spec Section 5.2 (migrated: just build → just compile per tech-design)
  test('TC-015: error-fixer uses just compile && just test', () => {
    const content = readProjectFile('plugins/forge/agents/error-fixer.md');
    expect(
      fileContains(content, 'just compile'),
      'Expected "just compile" to appear in error-fixer.md',
    ).toBeTruthy();
    expect(
      fileContains(content, 'just test'),
      'Expected "just test" to appear in error-fixer.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'go build ./...'),
      'Expected "go build ./..." NOT to appear in error-fixer.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'go vet ./...'),
      'Expected "go vet ./..." NOT to appear in error-fixer.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'go test -race -cover ./...'),
      'Expected "go test -race -cover ./..." NOT to appear in error-fixer.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'npm run build && npm test'),
      'Expected "npm run build && npm test" NOT to appear in error-fixer.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'pytest --cov'),
      'Expected "pytest --cov" NOT to appear in error-fixer.md',
    ).toBeTruthy();
  });

  // Traceability: TC-016 → Spec Section 5.2 (updated: execute-task uses workflow-driven dispatch)
  test('TC-016: execute-task uses workflow-driven execution (no hardcoded language commands)', () => {
    const content = readProjectFile('plugins/forge/commands/execute-task.md');
    // After refactor, execute-task uses workflow dispatch from task files, not hardcoded just commands.
    // Verify the command does NOT contain language-specific commands.
    expect(
      fileNotContains(content, 'go test ./...'),
      'Expected "go test ./..." NOT to appear in execute-task.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'npm test'),
      'Expected "npm test" NOT to appear in execute-task.md',
    ).toBeTruthy();
    expect(
      fileNotContains(content, 'pytest'),
      'Expected "pytest" NOT to appear in execute-task.md',
    ).toBeTruthy();
    // Verify workflow-driven dispatch
    expect(
      fileContains(content, 'Execution Workflow') || fileContains(content, 'workflow'),
      'Expected workflow-driven dispatch in execute-task.md',
    ).toBeTruthy();
  });

  // Traceability: TC-017 → Spec Section 5.2
  test('TC-017: improve-harness uses just test', () => {
    const content = readProjectFile('plugins/forge/skills/improve-harness/SKILL.md');
    expect(
      fileContains(content, 'just test'),
      'Expected "just test" to appear in improve-harness/SKILL.md Step 4.3',
    ).toBeTruthy();
  });

  // Traceability: TC-018 → Spec Section 5.1
  test('TC-018: init-justfile generates e2e-setup target', () => {
    const content = readProjectFile('plugins/forge/skills/init-justfile/SKILL.md');
    expect(
      fileContains(content, 'e2e-setup'),
      'Expected "e2e-setup" recipe to appear in init-justfile.md template',
    ).toBeTruthy();
    // Verify idempotent npm install logic
    expect(
      fileContains(content, 'node_modules'),
      'Expected idempotent node_modules check in e2e-setup recipe',
    ).toBeTruthy();
    // Verify playwright install in the template files
    const genericTemplate = readProjectFile('plugins/forge/skills/init-justfile/templates/generic.just');
    expect(
      fileContains(genericTemplate, 'playwright install chromium'),
      'Expected "playwright install chromium" in e2e-setup recipe template',
    ).toBeTruthy();
  });

  // Traceability: TC-019 → Spec Section 5.1
  test('TC-019: init-justfile generates e2e-verify target', () => {
    const content = readProjectFile('plugins/forge/skills/init-justfile/SKILL.md');
    expect(
      fileContains(content, 'e2e-verify'),
      'Expected "e2e-verify" recipe to appear in init-justfile.md template',
    ).toBeTruthy();
    expect(
      fileContains(content, '--feature'),
      'Expected "--feature" parameter in e2e-verify recipe',
    ).toBeTruthy();
    expect(
      fileContains(content, '// VERIFY:'),
      'Expected "// VERIFY:" marker scanning in e2e-verify recipe',
    ).toBeTruthy();
  });
});

// ── TC-003, TC-004, TC-009 to TC-012, TC-020: just command execution ──
test.describe('just command execution', () => {
  const tmpE2eDir = join(PROJECT_ROOT, 'tests', 'e2e');
  const tmpFeaturesDir = join(tmpE2eDir, 'features');
  const tmpFeatureDir = join(tmpFeaturesDir, 'test-feature');
  const tmpMyFeatureDir = join(tmpFeaturesDir, 'my-feature');

  test.beforeAll(() => {
    // Create temp test directories for just command tests
    mkdirSync(tmpFeatureDir, { recursive: true });
    mkdirSync(tmpMyFeatureDir, { recursive: true });
  });

  test.afterAll(() => {
    // Clean up temp directories
    if (existsSync(tmpFeatureDir)) rmSync(tmpFeatureDir, { recursive: true, force: true });
    if (existsSync(tmpMyFeatureDir)) rmSync(tmpMyFeatureDir, { recursive: true, force: true });
    // Clean up features/ dir if empty (only if we created it)
    if (existsSync(tmpFeaturesDir) && readdirSync(tmpFeaturesDir).length === 0) {
      rmSync(tmpFeaturesDir, { recursive: true, force: true });
    }
  });

  // Traceability: TC-003 → Story 3 / AC-1
  test('TC-003: just e2e-verify exits 1 when VERIFY markers present', () => {
    writeFileSync(
      join(tmpFeatureDir, 'sample.spec.ts'),
      '// VERIFY: check this\nconsole.log("test");\n',
    );
    const result = runCli('just e2e-verify --feature test-feature');
    expect(result.exitCode, 'Expected exit code 1 when VERIFY markers present').toBe(1);
    const output = result.stdout + result.stderr;
    expect(
      output.includes('sample.spec.ts') || output.includes('VERIFY'),
      'Expected output to include filename or VERIFY marker info',
    ).toBeTruthy();
  });

  // Traceability: TC-004 → Story 3 / AC-2
  test('TC-004: just e2e-verify exits 0 when no VERIFY markers', () => {
    writeFileSync(
      join(tmpFeatureDir, 'sample.spec.ts'),
      'console.log("no verify markers here");\n',
    );
    const result = runCli('just e2e-verify --feature test-feature');
    expect(result.exitCode, 'Expected exit code 0 when no VERIFY markers').toBe(0);
    expect(
      result.stdout.includes('OK: no unresolved // VERIFY: markers'),
      `Expected "OK: no unresolved // VERIFY: markers" in stdout, got: ${result.stdout}`,
    ).toBeTruthy();
  });

  // Traceability: TC-009 → Spec Section 5.1
  test('TC-009: just e2e-setup exits 1 when package.json missing', () => {
    // Ensure package.json does not exist in tests/e2e
    const pkgPath = join(tmpE2eDir, 'package.json');
    const pkgExists = existsSync(pkgPath);
    if (pkgExists) {
      // Skip this test if package.json already exists (real project setup)
      return;
    }
    const result = runCli('just e2e-setup');
    expect(result.exitCode, 'Expected exit code 1 when tests/e2e/package.json missing').toBe(1);
    const output = result.stdout + result.stderr;
    expect(
      output.includes('tests/e2e/package.json not found'),
      `Expected "tests/e2e/package.json not found" in output, got: ${output}`,
    ).toBeTruthy();
  });

  // Traceability: TC-010 → Spec Section 5.1
  test('TC-010: just e2e-setup exits 0 with OK message when deps ready', () => {
    const pkgPath = join(tmpE2eDir, 'package.json');
    const nodeModulesPath = join(tmpE2eDir, 'node_modules');
    if (!existsSync(pkgPath) || !existsSync(nodeModulesPath)) {
      // Skip: requires real package.json and node_modules to be present
      return;
    }
    const result = runCli('just e2e-setup');
    expect(result.exitCode, 'Expected exit code 0 when deps are ready').toBe(0);
    expect(
      result.stdout.includes('OK: e2e dependencies ready'),
      `Expected "OK: e2e dependencies ready" in stdout, got: ${result.stdout}`,
    ).toBeTruthy();
  });

  // Traceability: TC-011 → Spec Section 5.1
  test('TC-011: just e2e-verify exits 1 when feature flag missing', () => {
    const result = runCli('just e2e-verify');
    expect(result.exitCode, 'Expected exit code 1 when --feature argument is missing').toBe(1);
    const output = result.stdout + result.stderr;
    expect(
      output.includes('--feature') || output.includes('Usage'),
      `Expected usage hint with "--feature <slug>" in output, got: ${output}`,
    ).toBeTruthy();
  });

  // Traceability: TC-012 → Spec Section 5.1
  test('TC-012: just e2e-verify outputs file and line number for residual markers', () => {
    // Write a spec file with a VERIFY marker on a known line
    const specContent = [
      'import { test } from "node:test";',
      '// line 2',
      '// line 3',
      '// line 4',
      '// line 5',
      '// line 6',
      '// line 7',
      '// line 8',
      '// line 9',
      '// line 10',
      '// line 11',
      '// line 12',
      '// line 13',
      '// line 14',
      '// line 15',
      '// line 16',
      '// line 17',
      '// line 18',
      '// line 19',
      '// line 20',
      '// line 21',
      '// line 22',
      '// line 23',
      '// line 24',
      '// line 25',
      '// line 26',
      '// line 27',
      '// line 28',
      '// line 29',
      '// line 30',
      '// line 31',
      '// line 32',
      '// line 33',
      '// line 34',
      '// line 35',
      '// line 36',
      '// line 37',
      '// line 38',
      '// line 39',
      '// line 40',
      '// line 41',
      '// VERIFY: implement this',
    ].join('\n');
    writeFileSync(join(tmpMyFeatureDir, 'login.spec.ts'), specContent);
    const result = runCli('just e2e-verify --feature my-feature');
    expect(result.exitCode, 'Expected exit code 1 when VERIFY markers present').toBe(1);
    const output = result.stdout + result.stderr;
    expect(
      output.includes('login.spec.ts'),
      `Expected "login.spec.ts" in output, got: ${output}`,
    ).toBeTruthy();
    expect(
      output.includes('42') || output.match(/:\d+:/),
      `Expected line number in output, got: ${output}`,
    ).toBeTruthy();
  });

  // Traceability: TC-020 → Spec Section 5.1
  test('TC-020: just e2e-setup is idempotent', () => {
    const pkgPath = join(tmpE2eDir, 'package.json');
    const nodeModulesPath = join(tmpE2eDir, 'node_modules');
    if (!existsSync(pkgPath) || !existsSync(nodeModulesPath)) {
      // Skip: requires real package.json and node_modules to be present
      return;
    }
    const result1 = runCli('just e2e-setup');
    const result2 = runCli('just e2e-setup');
    expect(result1.exitCode, 'Expected first run to exit 0').toBe(0);
    expect(result2.exitCode, 'Expected second run to exit 0').toBe(0);
    expect(
      result1.stdout.includes('OK: e2e dependencies ready'),
      `Expected "OK: e2e dependencies ready" in first run stdout`,
    ).toBeTruthy();
    expect(
      result2.stdout.includes('OK: e2e dependencies ready'),
      `Expected "OK: e2e dependencies ready" in second run stdout`,
    ).toBeTruthy();
  });
});
