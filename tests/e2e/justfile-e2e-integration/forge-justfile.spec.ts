import { describe, test } from 'node:test';
import assert from 'node:assert/strict';
import { readProjectFile, runCli } from '../helpers.js';

// ── Helpers ────────────────────────────────────────────────────────
function fileContains(content: string, needle: string): boolean {
  return content.includes(needle);
}

function getJustfile(): string {
  return readProjectFile('justfile');
}

// Extract the forge standard recipes section
function getStandardSection(): string {
  const content = getJustfile();
  const startMarker = '# --- forge standard recipes ---';
  const endMarker = '# --- end forge standard recipes ---';
  const startIdx = content.indexOf(startMarker);
  const endIdx = content.indexOf(endMarker);
  assert.notEqual(startIdx, -1, 'Expected start boundary marker in justfile');
  assert.notEqual(endIdx, -1, 'Expected end boundary marker in justfile');
  return content.slice(startIdx, endIdx + endMarker.length);
}

// ── TC-FJ-001 to TC-FJ-010: Standard recipes presence ─────────────
describe('Forge justfile: all 15 standard recipes present', () => {

  // Traceability: TC-FJ-001 -> AC: project-type outputs "mixed"
  test('TC-FJ-001: project-type recipe outputs "mixed"', () => {
    const result = runCli('just project-type');
    assert.equal(result.exitCode, 0, 'Expected exit code 0');
    assert.equal(result.stdout.trim(), 'mixed', 'Expected "mixed" output');
  });

  // Traceability: TC-FJ-002 -> AC: 10 scoped recipes present
  test('TC-FJ-002: 10 scoped recipes use bash case dispatch', () => {
    const section = getStandardSection();
    const scopedRecipes = ['compile', 'build', 'run', 'dev', 'test', 'lint', 'fmt', 'check', 'clean', 'install'];
    for (const recipe of scopedRecipes) {
      const pattern = `${recipe} scope=""`;
      assert.ok(
        fileContains(section, pattern),
        `Expected scoped recipe "${pattern}" in forge standard section`,
      );
    }
  });

  // Traceability: TC-FJ-003 -> AC: 5 unscoped recipes present
  test('TC-FJ-003: 5 unscoped recipes present (no scope parameter)', () => {
    const section = getStandardSection();
    const unscopedRecipes = ['project-type', 'test-e2e', 'ci', 'e2e-setup', 'e2e-verify'];
    for (const recipe of unscopedRecipes) {
      assert.ok(
        fileContains(section, recipe),
        `Expected recipe "${recipe}" in forge standard section`,
      );
    }
    // Verify these do NOT have scope=""
    assert.ok(!section.includes('project-type scope=""'), 'project-type should NOT have scope');
    assert.ok(!section.includes('ci scope=""'), 'ci should NOT have scope');
    assert.ok(!section.includes('e2e-setup scope=""'), 'e2e-setup should NOT have scope');
    assert.ok(!section.includes('e2e-verify scope=""'), 'e2e-verify should NOT have scope');
  });

  // Traceability: TC-FJ-004 -> AC: Boundary markers present
  test('TC-FJ-004: boundary markers present', () => {
    const content = getJustfile();
    assert.ok(
      fileContains(content, '# --- forge standard recipes ---'),
      'Expected start boundary marker',
    );
    assert.ok(
      fileContains(content, '# --- end forge standard recipes ---'),
      'Expected end boundary marker',
    );
  });

  // Traceability: TC-FJ-005 -> AC: compile frontend/backend/all work
  test('TC-FJ-005: compile recipe has frontend/backend/empty branches', () => {
    const section = getStandardSection();
    assert.ok(fileContains(section, 'compile scope=""'), 'Expected compile with scope');
    // Frontend: tsc --noEmit, Backend: go vet
    const compileMatch = section.match(/compile scope=""[\s\S]*?esac/);
    assert.ok(compileMatch, 'Expected compile recipe with bash case');
    assert.ok(compileMatch[0].includes('npx tsc --noEmit'), 'Expected frontend compile: npx tsc --noEmit');
    assert.ok(compileMatch[0].includes('go vet ./...'), 'Expected backend compile: go vet ./...');
  });

  // Traceability: TC-FJ-006 -> AC: build has correct frontend/backend branches
  test('TC-FJ-006: build recipe has frontend npm and backend go branches', () => {
    const section = getStandardSection();
    const buildMatch = section.match(/build scope=""[\s\S]*?esac/);
    assert.ok(buildMatch, 'Expected build recipe with bash case');
    assert.ok(buildMatch[0].includes('npm run build'), 'Expected frontend: npm run build');
    assert.ok(buildMatch[0].includes('go build ./...'), 'Expected backend: go build ./...');
  });

  // Traceability: TC-FJ-007 -> AC: test recipe has correct branches
  test('TC-FJ-007: test recipe has frontend npm and backend go branches', () => {
    const section = getStandardSection();
    const testMatch = section.match(/test scope=""[\s\S]*?esac/);
    assert.ok(testMatch, 'Expected test recipe with bash case');
    assert.ok(testMatch[0].includes('npm test'), 'Expected frontend: npm test');
    assert.ok(testMatch[0].includes('go test -race ./...'), 'Expected backend: go test -race ./...');
  });

  // Traceability: TC-FJ-008 -> AC: *) error branch in scoped recipes
  test('TC-FJ-008: all scoped recipes have *) error branch with stderr', () => {
    const section = getStandardSection();
    const scopedRecipes = ['compile', 'build', 'run', 'dev', 'test', 'lint', 'fmt', 'check', 'clean', 'install'];
    for (const recipe of scopedRecipes) {
      const recipeMatch = section.match(new RegExp(`${recipe} scope=""[\\s\\S]*?esac`));
      assert.ok(recipeMatch, `Expected ${recipe} recipe`);
      assert.ok(
        recipeMatch[0].includes("echo \"[forge] invalid scope"),
        `Expected *) error branch in ${recipe} recipe`,
      );
      assert.ok(
        recipeMatch[0].includes('>&2'),
        `Expected stderr redirect in ${recipe} recipe error branch`,
      );
    }
  });

  // Traceability: TC-FJ-009 -> AC: ci chains standard commands
  test('TC-FJ-009: ci recipe chains install, compile, build, test, lint', () => {
    const section = getStandardSection();
    assert.ok(fileContains(section, 'just install'), 'Expected "just install" in ci');
    assert.ok(fileContains(section, 'just compile'), 'Expected "just compile" in ci');
    assert.ok(fileContains(section, 'just build'), 'Expected "just build" in ci');
    assert.ok(fileContains(section, 'just test'), 'Expected "just test" in ci');
    assert.ok(fileContains(section, 'just lint'), 'Expected "just lint" in ci');
  });

  // Traceability: TC-FJ-010 -> AC: Custom recipes preserved
  test('TC-FJ-010: custom recipes (claude, claude-c) preserved outside boundary markers', () => {
    const content = getJustfile();
    assert.ok(fileContains(content, 'claude:'), 'Expected "claude:" recipe preserved');
    assert.ok(fileContains(content, 'claude-c:'), 'Expected "claude-c:" recipe preserved');
  });
});

// ── TC-FJ-011 to TC-FJ-013: Live just command execution ────────────
describe('Forge justfile: live command execution', () => {

  // Traceability: TC-FJ-011 -> AC: just compile backend dispatches correctly
  test('TC-FJ-011: just compile backend dispatches to backend branch (no scope error)', () => {
    const result = runCli('just compile backend');
    const output = result.stdout + result.stderr;
    // May fail due to missing toolchain at root, but must NOT be a scope error
    const isScopeError = output.includes('[forge] invalid scope');
    assert.ok(!isScopeError, 'Should not be a scope error for backend scope');
  });

  // Traceability: TC-FJ-012 -> AC: just compile frontend dispatches correctly
  test('TC-FJ-012: just compile frontend dispatches to frontend branch (no scope error)', () => {
    const result = runCli('just compile frontend');
    // May fail if no tsconfig at root, but should NOT fail due to scope dispatch
    const output = result.stdout + result.stderr;
    const isScopeError = output.includes('[forge] invalid scope');
    assert.ok(!isScopeError, 'Should not be a scope error for frontend scope');
  });

  // Traceability: TC-FJ-013 -> AC: just compile (empty scope) dispatches correctly
  test('TC-FJ-013: just compile with empty scope dispatches to both branches (no scope error)', () => {
    const result = runCli('just compile');
    const output = result.stdout + result.stderr;
    const isScopeError = output.includes('[forge] invalid scope');
    assert.ok(!isScopeError, 'Should not be a scope error with empty scope');
  });

  // Traceability: TC-FJ-014 -> AC: invalid scope produces error
  test('TC-FJ-014: just compile with invalid scope exits 1', () => {
    const result = runCli('just compile invalidscope');
    assert.equal(result.exitCode, 1, 'Expected exit 1 for invalid scope');
    const output = result.stdout + result.stderr;
    assert.ok(
      output.includes('[forge] invalid scope'),
      `Expected "[forge] invalid scope" error, got: ${output}`,
    );
  });

  // Traceability: TC-FJ-015 -> AC: just project-type returns mixed
  test('TC-FJ-015: just project-type returns exactly "mixed"', () => {
    const result = runCli('just project-type');
    assert.equal(result.exitCode, 0, 'Expected exit 0');
    assert.equal(result.stdout.trim(), 'mixed', 'Expected exactly "mixed"');
  });
});
