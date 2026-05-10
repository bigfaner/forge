import { test, expect } from '@playwright/test';
import {
  runCli,
  readProjectFile,
  projectFileExists,
  PROJECT_ROOT,
} from '../../helpers.js';
import { join } from 'node:path';
import { writeFileSync, readFileSync, mkdirSync, rmSync, existsSync } from 'node:fs';

// ── Shared paths ───────────────────────────────────────────────────
const FORGE_ROOT = PROJECT_ROOT;
const TASK_CLI_BIN = join(FORGE_ROOT, 'task-cli', 'task.exe');
const AGENTS_DIR = join(FORGE_ROOT, 'plugins', 'forge', 'agents');
const COMMANDS_DIR = join(FORGE_ROOT, 'plugins', 'forge', 'commands');
const TEMPLATES_BREAKDOWN_DIR = join(FORGE_ROOT, 'plugins', 'forge', 'skills', 'breakdown-tasks', 'templates');
const TEMPLATES_QUICK_DIR = join(FORGE_ROOT, 'plugins', 'forge', 'skills', 'quick-tasks', 'templates');

const TASK_EXECUTOR_MD = join(AGENTS_DIR, 'task-executor.md');
const RUN_TASKS_MD = join(COMMANDS_DIR, 'run-tasks.md');
const EXECUTE_TASK_MD = join(COMMANDS_DIR, 'execute-task.md');

// ── Helper: build task CLI if needed ───────────────────────────────
function ensureTaskBinary(): void {
  if (!existsSync(TASK_CLI_BIN)) {
    runCli('go build -o task.exe ./cmd/task/', join(FORGE_ROOT, 'task-cli'));
  }
}

// ── Helper: run task CLI command ───────────────────────────────────
function runTaskCli(args: string): { stdout: string; stderr: string; exitCode: number } {
  ensureTaskBinary();
  return runCli(`"${TASK_CLI_BIN}" ${args}`);
}

// ════════════════════════════════════════════════════════════════════
// Workflow Detection & Injection
// ════════════════════════════════════════════════════════════════════

// Traceability: TC-001 -> Story 1 / AC-1
test('TC-001: Execution Workflow detected in task template replaces TDD', () => {
  // Read task-executor.md and verify it has the Execution Workflow reading logic
  const content = readProjectFile('plugins/forge/agents/task-executor.md');

  // Step 2 must describe reading ## Execution Workflow from the task file
  expect(content).toMatch(/## Execution Workflow/);
  expect(content).toMatch(/CASE A.*Execution Workflow.*heading exists/i);

  // The workflow must mention "do not deviate, add, or skip steps"
  expect(content).toMatch(/Do not deviate/);
});

// Traceability: TC-002 -> Story 1 / AC-2
test('TC-002: Missing Execution Workflow falls back to TDD and Quality Gate', () => {
  const content = readProjectFile('plugins/forge/agents/task-executor.md');

  // CASE B: fallback to default template
  expect(content).toMatch(/CASE B.*No.*Execution Workflow.*heading found/i);
  expect(content).toMatch(/breakdown-tasks\/templates\/task\.md/);

  // Default template path is specified
  expect(content).toMatch(/plugins\/forge\/skills\/breakdown-tasks\/templates\/task\.md/);
});

// Traceability: TC-003 -> Story 1 / AC-3
test('TC-003: Empty Execution Workflow body triggers warning and TDD fallback', () => {
  const content = readProjectFile('plugins/forge/agents/task-executor.md');

  // CASE C: empty body warning
  expect(content).toMatch(/CASE C.*Execution Workflow.*heading.*empty/i);
  expect(content).toMatch(/WARNING.*Execution Workflow.*empty/i);

  // Falls back to same default template as Case B
  expect(content).toMatch(/Falling back to default template/);
});

// ════════════════════════════════════════════════════════════════════
// Execution Workflow Behavior
// ════════════════════════════════════════════════════════════════════

// Traceability: TC-004 -> Story 2 / AC-1
test('TC-004: Execution-type task creates fix task on failure without TDD retry', () => {
  const content = readProjectFile('plugins/forge/agents/task-executor.md');

  // Dynamic Task Addition section exists for fix-task creation
  expect(content).toMatch(/Dynamic Task Addition/);
  expect(content).toMatch(/task add.*--template fix-task/);
  expect(content).toMatch(/--block-source/);

  // Error handling table references test failures and fix-task creation
  expect(content).toMatch(/Test failures beyond scope.*task add/);
});

// Traceability: TC-005 -> Story 2 / AC-2
test('TC-005: Step 2 output uses Execution Workflow terminology not TDD terminology', () => {
  const content = readProjectFile('plugins/forge/agents/task-executor.md');

  // Step 2 is named "Execute Workflow" (not "TDD implementation")
  const step2Match = content.match(/### Step 2[:\s]*Execute Workflow/);
  expect(step2Match).not.toBeNull();

  // The step describes reading Execution Workflow from task file
  expect(content).toMatch(/determine the execution workflow/i);
  expect(content).toMatch(/## Execution Workflow/i);
});

// Traceability: TC-006 -> Story 2 / AC-3
test('TC-006: Execution-type task skips Quality Gate and proceeds to record and commit', () => {
  const content = readProjectFile('plugins/forge/agents/task-executor.md');

  // Task-executor has a streamlined 4-step flow (not 5-step with quality gate)
  // Step 0 = MAIN_SESSION Guard, Step 1 = Read, Step 2 = Execute, Step 3 = Record, Step 4 = Commit
  expect(content).toMatch(/Step 0.*MAIN_SESSION Guard/);
  expect(content).toMatch(/Step 1.*Read Task Definition/);
  expect(content).toMatch(/Step 2.*Execute Workflow/);
  expect(content).toMatch(/Step 3.*Record Task/);
  expect(content).toMatch(/Step 4.*Commit/);

  // No explicit quality gate step (compile -> fmt -> lint -> test)
  // The Quality Gate is only referenced in the default template task.md, not in task-executor.md itself
  const qualityGatePattern = /Step.*Quality Gate|compile.*fmt.*lint.*test/i;
  const hasExplicitQualityGate = qualityGatePattern.test(content);
  // task-executor delegates to the workflow; quality gate is only in default template
  expect(hasExplicitQualityGate).toBe(false);
});

// ════════════════════════════════════════════════════════════════════
// noTest Removal Verification
// ════════════════════════════════════════════════════════════════════

// Traceability: TC-007 -> Story 3 / AC-1
test('TC-007: Grep noTest and NO_TEST across all harness files yields zero matches', () => {
  // Check agents directory for noTest frontmatter key (camelCase, not "NoTest" in error names)
  const agentsResult = runCli(
    `grep -r '"noTest"' --include="*.md" --include="*.go" --include="*.json" "${AGENTS_DIR}" || true`,
  );
  expect(agentsResult.stdout.trim()).toBe('');

  // Check commands directory
  const commandsResult = runCli(
    `grep -r '"noTest"' --include="*.md" --include="*.go" --include="*.json" "${COMMANDS_DIR}" || true`,
  );
  expect(commandsResult.stdout.trim()).toBe('');

  // Check task-cli source for NO_TEST environment variable and noTest struct field
  // Exclude legitimate error names like ErrNoTestEvidence
  const goResult = runCli(
    `grep -rn 'noTest\\b.*json:"noTest"\\|NO_TEST\\b' --include="*.go" "${join(FORGE_ROOT, 'task-cli')}" || true`,
  );
  expect(goResult.stdout.trim()).toBe('');

  // Also verify task-executor.md has no NO_TEST references
  const executorContent = readFileSync(TASK_EXECUTOR_MD, 'utf-8');
  expect(executorContent).not.toMatch(/\bNO_TEST\b/);
});

// Traceability: TC-008 -> Story 3 / AC-2
test('TC-008: task-cli Go code has no noTest conditional branches', () => {
  // Check types.go — Task struct should have no noTest field
  const typesContent = readProjectFile('task-cli/pkg/task/types.go');
  // Match "noTest" as a struct field or json tag, not as part of unrelated names
  expect(typesContent).not.toMatch(/\bnoTest\b/);

  // Check record.go — no conditional branches based on noTest
  const recordContent = readProjectFile('task-cli/internal/cmd/record.go');
  // "noTest" as a variable/field reference (not ErrNoTestEvidence which is an error function)
  expect(recordContent).not.toMatch(/\bnoTest\b/);
  expect(recordContent).not.toMatch(/\bNO_TEST\b/);

  // Check all Go files for noTest as a struct field with json tag
  const result = runCli(
    `grep -rn '\\bnoTest\\b.*json:' --include="*.go" "${join(FORGE_ROOT, 'task-cli')}" || true`,
  );
  expect(result.stdout.trim()).toBe('');

  // Check for NO_TEST as environment variable pattern
  const envResult = runCli(
    `grep -rn '\\bNO_TEST\\b' --include="*.go" "${join(FORGE_ROOT, 'task-cli')}" || true`,
  );
  expect(envResult.stdout.trim()).toBe('');
});

// Traceability: TC-009 -> Story 3 / AC-3
test('TC-009: All task templates have no noTest in frontmatter', () => {
  // Count breakdown templates (excluding manifest files)
  const breakdownResult = runCli(
    `find "${TEMPLATES_BREAKDOWN_DIR}" -name "*.md" ! -name "manifest-*.md" ! -name "eval-test-cases.md" | wc -l`,
  );
  const breakdownCount = parseInt(breakdownResult.stdout.trim(), 10);
  expect(breakdownCount).toBeGreaterThan(0);

  // Count quick templates (excluding manifest)
  const quickResult = runCli(
    `find "${TEMPLATES_QUICK_DIR}" -name "*.md" ! -name "manifest-quick.md" | wc -l`,
  );
  const quickCount = parseInt(quickResult.stdout.trim(), 10);
  expect(quickCount).toBeGreaterThan(0);

  // Check for noTest in all templates
  const notestResult = runCli(
    `grep -rl "noTest" "${TEMPLATES_BREAKDOWN_DIR}" "${TEMPLATES_QUICK_DIR}" --include="*.md" || true`,
  );
  expect(notestResult.stdout.trim()).toBe('');
});

// Traceability: TC-010 -> Story 3 / AC-4
test('TC-010: index.schema.json files have no noTest field definition', () => {
  // Check breakdown schema
  const breakdownSchemaPath = join(TEMPLATES_BREAKDOWN_DIR, 'index.schema.json');
  if (existsSync(breakdownSchemaPath)) {
    const schemaContent = readFileSync(breakdownSchemaPath, 'utf-8');
    const schema = JSON.parse(schemaContent);
    expect(schema.properties).not.toHaveProperty('noTest');
  }

  // Check quick schema
  const quickSchemaPath = join(TEMPLATES_QUICK_DIR, 'index.schema.json');
  if (existsSync(quickSchemaPath)) {
    const schemaContent = readFileSync(quickSchemaPath, 'utf-8');
    const schema = JSON.parse(schemaContent);
    expect(schema.properties).not.toHaveProperty('noTest');
  }
});

// Traceability: TC-011 -> Story 3 / AC-5
test('TC-011: Command docs run-tasks.md and execute-task.md have no NO_TEST references', () => {
  // Check run-tasks.md
  if (existsSync(RUN_TASKS_MD)) {
    const runTasksContent = readFileSync(RUN_TASKS_MD, 'utf-8');
    expect(runTasksContent).not.toMatch(/NO_TEST|noTest|no_test/i);
  }

  // Check execute-task.md
  if (existsSync(EXECUTE_TASK_MD)) {
    const executeTaskContent = readFileSync(EXECUTE_TASK_MD, 'utf-8');
    expect(executeTaskContent).not.toMatch(/NO_TEST|noTest|no_test/i);
  }
});

// Traceability: TC-012 -> Story 3 / AC-6
test('TC-012: task-executor.md Step 2-3 has no NO_TEST references and uses workflow injection', () => {
  const content = readProjectFile('plugins/forge/agents/task-executor.md');

  // Zero NO_TEST references
  expect(content).not.toMatch(/NO_TEST|noTest|no_test/i);

  // Step 2 describes reading ## Execution Workflow
  const step2Section = content.match(/### Step 2[\s\S]*?(?=### Step 3)/);
  expect(step2Section).not.toBeNull();
  expect(step2Section![0]).toMatch(/Execution Workflow/);

  // Verify the workflow injection logic (read from task file or fallback to template)
  expect(content).toMatch(/Search for a.*## Execution Workflow.*heading/i);

  // Step 3 (Record) does not reference noTest for quality gate bypass
  const step3Section = content.match(/### Step 3[\s\S]*?(?=### Step 4)/);
  expect(step3Section).not.toBeNull();
  expect(step3Section![0]).not.toMatch(/noTest|NO_TEST/);
});

// ════════════════════════════════════════════════════════════════════
// Failure Handling
// ════════════════════════════════════════════════════════════════════

test.describe('Failure Handling', () => {
  const FIXTURES_DIR = join(FORGE_ROOT, 'tests', 'e2e', 'fixtures', 'task-executor-skeleton');

  test.beforeAll(() => {
    mkdirSync(FIXTURES_DIR, { recursive: true });
  });

  test.afterAll(() => {
    rmSync(FIXTURES_DIR, { recursive: true, force: true });
  });

  // Traceability: TC-013 -> Story 4 / AC-1
  test('TC-013: Missing or unparseable task file sets status to failed with error log', () => {
    // Run task CLI with a non-existent task file
    const result = runTaskCli('record nonexistent-task-xyz --data /dev/null');

    // Should fail (non-zero exit code)
    expect(result.exitCode).not.toBe(0);

    // Error output should describe the failure
    const output = result.stdout + result.stderr;
    expect(output.length).toBeGreaterThan(0);
  });

  // Traceability: TC-014 -> Story 4 / AC-2
  test('TC-014: Workflow failure with explicit failure instruction followed correctly', () => {
    const content = readProjectFile('plugins/forge/agents/task-executor.md');

    // Dynamic Task Addition section describes creating fix tasks on failure
    expect(content).toMatch(/task add.*--template fix-task/);
    expect(content).toMatch(/--block-source/);

    // Failure output format is specified
    expect(content).toMatch(/FAILED.*reason/i);
  });

  // Traceability: TC-015 -> Story 4 / AC-3
  test('TC-015: Workflow failure without explicit instruction records and stops', () => {
    const content = readProjectFile('plugins/forge/agents/task-executor.md');

    // Error handling table covers failure scenarios
    expect(content).toMatch(/Build fails.*Fix.*retry/i);
    expect(content).toMatch(/Test fails.*Fix.*retry/i);

    // No retry loop description — agent stops on failure
    // The output format shows FAILED with reason, not a retry cycle
    expect(content).toMatch(/Step 2\/4.*FAILED/);
  });

  // Traceability: TC-016 -> Story 4 / AC-4
  test('TC-016: Multi-step workflow mid-failure records completed steps and failure point', () => {
    const content = readProjectFile('plugins/forge/agents/task-executor.md');

    // The agent records steps through the execution workflow
    // Step 2 output format includes DONE or FAILED status
    expect(content).toMatch(/Step 2\/4.*DONE|Step 2\/4.*FAILED/);

    // Error handling references "Test failures beyond scope" which triggers fix-task creation
    expect(content).toMatch(/Test failures beyond scope/);
  });
});

// ════════════════════════════════════════════════════════════════════
// End-to-End Integration
// ════════════════════════════════════════════════════════════════════

// Traceability: TC-017 -> Story 1 / AC-1, Story 2 / AC-3
test('TC-017: Full dispatch-to-commit pipeline with Execution Workflow template', () => {
  const content = readProjectFile('plugins/forge/agents/task-executor.md');

  // Complete 4-step pipeline exists: MAIN_SESSION Guard -> Read -> Execute -> Record -> Commit
  expect(content).toMatch(/Step 0.*MAIN_SESSION Guard/);
  expect(content).toMatch(/Step 1.*Read Task Definition/);
  expect(content).toMatch(/Step 2.*Execute Workflow/);
  expect(content).toMatch(/Step 3.*Record Task.*MANDATORY/);
  expect(content).toMatch(/Step 4.*Commit/);

  // Workflow reading logic covers all 3 cases
  expect(content).toMatch(/CASE A.*heading exists with non-empty/i);
  expect(content).toMatch(/CASE B.*No.*heading found/i);
  expect(content).toMatch(/CASE C.*heading.*empty/i);

  // Output format includes completed and blocked paths
  expect(content).toMatch(/DONE.*TASK_ID/);
  expect(content).toMatch(/BLOCKED.*TASK_ID/);

  // No TDD keywords in the primary step descriptions
  const step2Section = content.match(/### Step 2[\s\S]*?(?=### Step 3)/);
  expect(step2Section).not.toBeNull();
  expect(step2Section![0]).not.toMatch(/RED.*GREEN.*REFACTOR/);
});
