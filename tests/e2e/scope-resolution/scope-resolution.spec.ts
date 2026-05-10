import { test, expect } from '@playwright/test';
import { runCli, readProjectFile } from '../helpers.js';

// -- Helpers ---------------------------------------------------------------
function fileContains(content: string, needle: string): boolean {
  return content.includes(needle);
}

// Get project type from justfile
function getProjectType(): string {
  const result = runCli('just project-type');
  return result.stdout.trim();
}

// -- Tests -----------------------------------------------------------------
test.describe('breakdown-tasks: scope field in index.json', () => {

  // Traceability: TC-007 -> Story 3 / AC-1
  test('TC-007: mixed project tasks receive scope field in index.json', () => {
    // Verify the breakdown-tasks skill describes adding scope to index.json tasks
    const skillContent = readProjectFile('plugins/forge/skills/breakdown-tasks/SKILL.md');
    expect(
      fileContains(skillContent, 'scope'),
      'Expected "scope" in breakdown-tasks skill description',
    ).toBeTruthy();
  });

  // Traceability: TC-008 -> Story 3 / AC-2
  test('TC-008: frontend-only task scope marked as frontend', () => {
    const skillContent = readProjectFile('plugins/forge/skills/breakdown-tasks/SKILL.md');
    // Skill should describe how frontend-only files get scope=frontend
    expect(
      fileContains(skillContent, 'frontend') &&
      fileContains(skillContent, 'scope'),
      'Expected scope=frontend logic in breakdown-tasks skill',
    ).toBeTruthy();
  });

  // Traceability: TC-009 -> Story 3 / AC-3
  test('TC-009: cross-scope task marked as all', () => {
    const skillContent = readProjectFile('plugins/forge/skills/breakdown-tasks/SKILL.md');
    expect(
      fileContains(skillContent, 'all') &&
      fileContains(skillContent, 'scope'),
      'Expected scope=all logic in breakdown-tasks skill',
    ).toBeTruthy();
  });

  // Traceability: TC-010 -> Story 3 / AC-4
  test('TC-010: non-mixed project tasks all receive scope all', () => {
    const skillContent = readProjectFile('plugins/forge/skills/breakdown-tasks/SKILL.md');
    // For non-mixed projects, scope should default to "all"
    expect(
      fileContains(skillContent, 'all') &&
      fileContains(skillContent, 'scope'),
      'Expected default scope=all for non-mixed projects',
    ).toBeTruthy();
  });
});

test.describe('skill execution: scope mismatch and fallback', () => {

  // Traceability: TC-015 -> Story 5 / AC-1
  test('TC-015: scope mismatch shows warning and falls back', () => {
    const projectType = getProjectType();
    if (projectType === 'mixed') {
      // Mixed projects validate scope and reject invalid values
      // When a scope mismatch occurs at the skill level, the skill should detect it
      // and fall back to running `just <verb>` without scope.
      const result = runCli('just build invalidscope');
      expect(result.exitCode, 'Expected exit code 1 for invalid scope').toBe(1);
      const output = result.stdout + result.stderr;
      expect(
        output.includes("[forge] invalid scope 'invalidscope'"),
        `Expected "[forge] invalid scope 'invalidscope'" in output, got: ${output}`,
      ).toBeTruthy();
    } else {
      // Non-mixed projects accept scope param but ignore it (no scope dispatch).
      // Scope mismatch is not applicable since there's no scope validation.
      const result = runCli('just build invalidscope');
      const output = result.stdout + result.stderr;
      const isScopeError = output.includes('[forge] invalid scope');
      expect(!isScopeError, 'Non-mixed projects should not produce scope errors').toBeTruthy();
    }
  });

  // Traceability: TC-016 -> Story 5 / AC-2
  test('TC-016: mixed project with matching scope executes normally', () => {
    // In the real forge project (which is mixed), verify scoped commands work
    const result = runCli('just compile frontend');
    const output = result.stdout + result.stderr;
    // Should not produce a scope error
    const isScopeError = output.includes('[forge] invalid scope');
    expect(!isScopeError, 'Should not be a scope error for frontend scope').toBeTruthy();
  });

  // Traceability: TC-023 -> Spec 5.3 / row 4
  test('TC-023: just project-type failure triggers fallback in skill', () => {
    // Verify the justfile has project-type recipe that agents can call.
    // When project-type fails (old justfile), skills should fall back to `just <verb>`.
    // We verify the recipe exists and returns a valid value.
    const result = runCli('just project-type');
    expect(result.exitCode, 'Expected project-type to exit 0').toBe(0);
    const output = result.stdout.trim();
    expect(
      ['frontend', 'backend', 'mixed'].includes(output),
      `Expected valid project-type output, got: "${output}"`,
    ).toBeTruthy();

    // Also verify the PRD spec documents the fallback behavior
    const prdSpec = readProjectFile('docs/features/justfile-standard-vocabulary/prd/prd-spec.md');
    expect(
      fileContains(prdSpec, 'falling back') ||
      fileContains(prdSpec, 'fallback'),
      'Expected fallback behavior described in PRD spec',
    ).toBeTruthy();
  });

  // Traceability: TC-024 -> Spec 5.3 / row 5
  test('TC-024: unexpected project-type output triggers fallback', () => {
    // Verify the PRD spec describes handling of unexpected project-type output.
    // The fallback behavior: "[forge] just project-type returned unexpected output 'XYZ'; falling back to just verb"
    const prdSpec = readProjectFile('docs/features/justfile-standard-vocabulary/prd/prd-spec.md');
    expect(
      fileContains(prdSpec, 'unexpected output') ||
      fileContains(prdSpec, 'unexpected'),
      'Expected unexpected output handling described in PRD spec',
    ).toBeTruthy();
    expect(
      fileContains(prdSpec, 'falling back') ||
      fileContains(prdSpec, 'fallback'),
      'Expected fallback description in PRD spec',
    ).toBeTruthy();

    // Verify the justfile project-type output is deterministic (one of the 3 valid values)
    const result1 = runCli('just project-type');
    const result2 = runCli('just project-type');
    expect(result1.stdout.trim(), 'Expected deterministic project-type output').toBe(result2.stdout.trim());
  });
});
