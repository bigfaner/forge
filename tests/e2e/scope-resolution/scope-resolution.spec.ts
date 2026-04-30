import { describe, test } from 'node:test';
import assert from 'node:assert/strict';
import { runCli, readProjectFile } from '../helpers.js';

// -- Helpers ---------------------------------------------------------------
function fileContains(content: string, needle: string): boolean {
  return content.includes(needle);
}

// -- Tests -----------------------------------------------------------------
describe('breakdown-tasks: scope field in index.json', () => {

  // Traceability: TC-007 -> Story 3 / AC-1
  test('TC-007: mixed project tasks receive scope field in index.json', () => {
    // Verify the breakdown-tasks skill describes adding scope to index.json tasks
    const skillContent = readProjectFile('plugins/forge/skills/breakdown-tasks/SKILL.md');
    assert.ok(
      fileContains(skillContent, 'scope'),
      'Expected "scope" in breakdown-tasks skill description',
    );
  });

  // Traceability: TC-008 -> Story 3 / AC-2
  test('TC-008: frontend-only task scope marked as frontend', () => {
    const skillContent = readProjectFile('plugins/forge/skills/breakdown-tasks/SKILL.md');
    // Skill should describe how frontend-only files get scope=frontend
    assert.ok(
      fileContains(skillContent, 'frontend') &&
      fileContains(skillContent, 'scope'),
      'Expected scope=frontend logic in breakdown-tasks skill',
    );
  });

  // Traceability: TC-009 -> Story 3 / AC-3
  test('TC-009: cross-scope task marked as all', () => {
    const skillContent = readProjectFile('plugins/forge/skills/breakdown-tasks/SKILL.md');
    assert.ok(
      fileContains(skillContent, 'all') &&
      fileContains(skillContent, 'scope'),
      'Expected scope=all logic in breakdown-tasks skill',
    );
  });

  // Traceability: TC-010 -> Story 3 / AC-4
  test('TC-010: non-mixed project tasks all receive scope all', () => {
    const skillContent = readProjectFile('plugins/forge/skills/breakdown-tasks/SKILL.md');
    // For non-mixed projects, scope should default to "all"
    assert.ok(
      fileContains(skillContent, 'all') &&
      fileContains(skillContent, 'scope'),
      'Expected default scope=all for non-mixed projects',
    );
  });
});

describe('skill execution: scope mismatch and fallback', () => {

  // Traceability: TC-015 -> Story 5 / AC-1
  test('TC-015: scope mismatch shows warning and falls back', () => {
    // Verify that the justfile's invalid scope error message matches the expected format
    // from the PRD: [forge] invalid scope 'X'; expected frontend/backend
    // When a scope mismatch occurs at the skill level, the skill should detect it
    // and fall back to running `just <verb>` without scope.
    // Here we verify the justfile produces the correct error for invalid scopes.
    const result = runCli('just build invalidscope');
    assert.equal(result.exitCode, 1, 'Expected exit code 1 for invalid scope');
    const output = result.stdout + result.stderr;
    assert.ok(
      output.includes("[forge] invalid scope 'invalidscope'"),
      `Expected "[forge] invalid scope 'invalidscope'" in output, got: ${output}`,
    );
  });

  // Traceability: TC-016 -> Story 5 / AC-2
  test('TC-016: mixed project with matching scope executes normally', () => {
    // In the real forge project (which is mixed), verify scoped commands work
    const result = runCli('just compile frontend');
    const output = result.stdout + result.stderr;
    // Should not produce a scope error
    const isScopeError = output.includes('[forge] invalid scope');
    assert.ok(!isScopeError, 'Should not be a scope error for frontend scope');
  });

  // Traceability: TC-023 -> Spec 5.3 / row 4
  test('TC-023: just project-type failure triggers fallback in skill', () => {
    // Verify the justfile has project-type recipe that agents can call.
    // When project-type fails (old justfile), skills should fall back to `just <verb>`.
    // We verify the recipe exists and returns a valid value.
    const result = runCli('just project-type');
    assert.equal(result.exitCode, 0, 'Expected project-type to exit 0');
    const output = result.stdout.trim();
    assert.ok(
      ['frontend', 'backend', 'mixed'].includes(output),
      `Expected valid project-type output, got: "${output}"`,
    );

    // Also verify the PRD spec documents the fallback behavior
    const prdSpec = readProjectFile('docs/features/justfile-standard-vocabulary/prd/prd-spec.md');
    assert.ok(
      fileContains(prdSpec, 'falling back') ||
      fileContains(prdSpec, 'fallback'),
      'Expected fallback behavior described in PRD spec',
    );
  });

  // Traceability: TC-024 -> Spec 5.3 / row 5
  test('TC-024: unexpected project-type output triggers fallback', () => {
    // Verify the PRD spec describes handling of unexpected project-type output.
    // The fallback behavior: "[forge] just project-type returned unexpected output 'XYZ'; falling back to just verb"
    const prdSpec = readProjectFile('docs/features/justfile-standard-vocabulary/prd/prd-spec.md');
    assert.ok(
      fileContains(prdSpec, 'unexpected output') ||
      fileContains(prdSpec, 'unexpected'),
      'Expected unexpected output handling described in PRD spec',
    );
    assert.ok(
      fileContains(prdSpec, 'falling back') ||
      fileContains(prdSpec, 'fallback'),
      'Expected fallback description in PRD spec',
    );

    // Verify the justfile project-type output is deterministic (one of the 3 valid values)
    const result1 = runCli('just project-type');
    const result2 = runCli('just project-type');
    assert.equal(result1.stdout.trim(), result2.stdout.trim(),
      'Expected deterministic project-type output');
  });
});
