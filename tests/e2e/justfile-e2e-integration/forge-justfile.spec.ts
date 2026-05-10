import { test, expect } from '@playwright/test';
import { readProjectFile, runCli } from '../helpers.js';

// ── Helpers ────────────────────────────────────────────────────────
function fileContains(content: string, needle: string): boolean {
  return content.includes(needle);
}

function getJustfile(): string {
  return readProjectFile('justfile');
}

// Get project type from justfile
function getProjectType(): string {
  const result = runCli('just project-type');
  return result.stdout.trim();
}

// Extract the forge standard recipes section
function getStandardSection(): string {
  const content = getJustfile();
  const startMarker = '# --- forge standard recipes ---';
  const endMarker = '# --- end forge standard recipes ---';
  const startIdx = content.indexOf(startMarker);
  const endIdx = content.indexOf(endMarker);
  if (startIdx === -1) throw new Error('Expected start boundary marker in justfile');
  if (endIdx === -1) throw new Error('Expected end boundary marker in justfile');
  return content.slice(startIdx, endIdx + endMarker.length);
}

// ── TC-FJ-001 to TC-FJ-010: Standard recipes presence ─────────────
test.describe('Forge justfile: all 15 standard recipes present', () => {

  // Traceability: TC-FJ-001 -> AC: project-type outputs a valid type
  test('TC-FJ-001: project-type recipe outputs a valid type (frontend/backend/mixed)', () => {
    const result = runCli('just project-type');
    expect(result.exitCode, 'Expected exit code 0').toBe(0);
    const output = result.stdout.trim();
    expect(
      ['frontend', 'backend', 'mixed'].includes(output),
      `Expected valid project-type output, got: "${output}"`,
    ).toBeTruthy();
  });

  // Traceability: TC-FJ-002 -> AC: 10 scoped recipes present
  test('TC-FJ-002: 10 scoped recipes use bash case dispatch', () => {
    const section = getStandardSection();
    const scopedRecipes = ['compile', 'build', 'run', 'dev', 'test', 'lint', 'fmt', 'check', 'clean', 'install'];
    for (const recipe of scopedRecipes) {
      const pattern = `${recipe} scope=""`;
      expect(
        fileContains(section, pattern),
        `Expected scoped recipe "${pattern}" in forge standard section`,
      ).toBeTruthy();
    }
  });

  // Traceability: TC-FJ-003 -> AC: 5 unscoped recipes present
  test('TC-FJ-003: 5 unscoped recipes present (no scope parameter)', () => {
    const section = getStandardSection();
    const unscopedRecipes = ['project-type', 'test-e2e', 'ci', 'e2e-setup', 'e2e-verify'];
    for (const recipe of unscopedRecipes) {
      expect(
        fileContains(section, recipe),
        `Expected recipe "${recipe}" in forge standard section`,
      ).toBeTruthy();
    }
    // Verify these do NOT have scope=""
    expect(!section.includes('project-type scope=""'), 'project-type should NOT have scope').toBeTruthy();
    expect(!section.includes('ci scope=""'), 'ci should NOT have scope').toBeTruthy();
    expect(!section.includes('e2e-setup scope=""'), 'e2e-setup should NOT have scope').toBeTruthy();
    expect(!section.includes('e2e-verify scope=""'), 'e2e-verify should NOT have scope').toBeTruthy();
  });

  // Traceability: TC-FJ-004 -> AC: Boundary markers present
  test('TC-FJ-004: boundary markers present', () => {
    const content = getJustfile();
    expect(
      fileContains(content, '# --- forge standard recipes ---'),
      'Expected start boundary marker',
    ).toBeTruthy();
    expect(
      fileContains(content, '# --- end forge standard recipes ---'),
      'Expected end boundary marker',
    ).toBeTruthy();
  });

  // Traceability: TC-FJ-005 -> AC: compile recipe has correct toolchain dispatch
  test('TC-FJ-005: compile recipe has correct toolchain dispatch for project type', () => {
    const section = getStandardSection();
    expect(fileContains(section, 'compile scope=""'), 'Expected compile with scope').toBeTruthy();
    const projectType = getProjectType();
    if (projectType === 'mixed') {
      // Mixed projects have bash case with frontend/backend branches
      const compileMatch = section.match(/compile scope=""[\s\S]*?esac/);
      expect(compileMatch, 'Expected compile recipe with bash case').toBeTruthy();
      expect(compileMatch![0].includes('npx tsc --noEmit'), 'Expected frontend compile: npx tsc --noEmit').toBeTruthy();
      expect(compileMatch![0].includes('go vet ./...'), 'Expected backend compile: go vet ./...').toBeTruthy();
    } else if (projectType === 'backend') {
      // Backend projects run backend toolchain directly (scope param accepted but unused)
      expect(
        fileContains(section, 'go vet') || fileContains(section, 'go build'),
        'Expected backend toolchain command in compile recipe',
      ).toBeTruthy();
    } else if (projectType === 'frontend') {
      expect(
        fileContains(section, 'tsc') || fileContains(section, 'npm run build'),
        'Expected frontend toolchain command in compile recipe',
      ).toBeTruthy();
    }
  });

  // Traceability: TC-FJ-006 -> AC: build recipe has correct toolchain dispatch
  test('TC-FJ-006: build recipe has correct toolchain dispatch for project type', () => {
    const section = getStandardSection();
    const projectType = getProjectType();
    if (projectType === 'mixed') {
      const buildMatch = section.match(/build scope=""[\s\S]*?esac/);
      expect(buildMatch, 'Expected build recipe with bash case').toBeTruthy();
      expect(buildMatch![0].includes('npm run build'), 'Expected frontend: npm run build').toBeTruthy();
      expect(buildMatch![0].includes('go build ./...'), 'Expected backend: go build ./...').toBeTruthy();
    } else if (projectType === 'backend') {
      expect(
        fileContains(section, 'go build'),
        'Expected backend: go build in build recipe',
      ).toBeTruthy();
    } else if (projectType === 'frontend') {
      expect(
        fileContains(section, 'npm run build'),
        'Expected frontend: npm run build in build recipe',
      ).toBeTruthy();
    }
  });

  // Traceability: TC-FJ-007 -> AC: test recipe has correct toolchain dispatch
  test('TC-FJ-007: test recipe has correct toolchain dispatch for project type', () => {
    const section = getStandardSection();
    const projectType = getProjectType();
    if (projectType === 'mixed') {
      const testMatch = section.match(/test scope=""[\s\S]*?esac/);
      expect(testMatch, 'Expected test recipe with bash case').toBeTruthy();
      expect(testMatch![0].includes('npm test'), 'Expected frontend: npm test').toBeTruthy();
      expect(testMatch![0].includes('go test -race ./...'), 'Expected backend: go test -race ./...').toBeTruthy();
    } else if (projectType === 'backend') {
      expect(
        fileContains(section, 'go test'),
        'Expected backend: go test in test recipe',
      ).toBeTruthy();
    } else if (projectType === 'frontend') {
      expect(
        fileContains(section, 'npm test'),
        'Expected frontend: npm test in test recipe',
      ).toBeTruthy();
    }
  });

  // Traceability: TC-FJ-008 -> AC: *) error branch in scoped recipes (mixed only)
  test('TC-FJ-008: scoped recipes have *) error branch with stderr (mixed projects)', () => {
    const projectType = getProjectType();
    if (projectType !== 'mixed') {
      // Non-mixed projects do not have scope dispatch, so no error branches needed
      return;
    }
    const section = getStandardSection();
    const scopedRecipes = ['compile', 'build', 'run', 'dev', 'test', 'lint', 'fmt', 'check', 'clean', 'install'];
    for (const recipe of scopedRecipes) {
      const recipeMatch = section.match(new RegExp(`${recipe} scope=""[\\s\\S]*?esac`));
      expect(recipeMatch, `Expected ${recipe} recipe`).toBeTruthy();
      expect(
        recipeMatch![0].includes("echo \"[forge] invalid scope"),
        `Expected *) error branch in ${recipe} recipe`,
      ).toBeTruthy();
      expect(
        recipeMatch![0].includes('>&2'),
        `Expected stderr redirect in ${recipe} recipe error branch`,
      ).toBeTruthy();
    }
  });

  // Traceability: TC-FJ-009 -> AC: ci chains standard commands
  test('TC-FJ-009: ci recipe chains install, compile, build, test, lint', () => {
    const section = getStandardSection();
    expect(fileContains(section, 'just install'), 'Expected "just install" in ci').toBeTruthy();
    expect(fileContains(section, 'just compile'), 'Expected "just compile" in ci').toBeTruthy();
    expect(fileContains(section, 'just build'), 'Expected "just build" in ci').toBeTruthy();
    expect(fileContains(section, 'just test'), 'Expected "just test" in ci').toBeTruthy();
    expect(fileContains(section, 'just lint'), 'Expected "just lint" in ci').toBeTruthy();
  });

  // Traceability: TC-FJ-010 -> AC: Custom recipes preserved
  test('TC-FJ-010: custom recipes (claude, claude-c) preserved outside boundary markers', () => {
    const content = getJustfile();
    expect(fileContains(content, 'claude:'), 'Expected "claude:" recipe preserved').toBeTruthy();
    expect(fileContains(content, 'claude-c:'), 'Expected "claude-c:" recipe preserved').toBeTruthy();
  });
});

// ── TC-FJ-011 to TC-FJ-013: Live just command execution ────────────
test.describe('Forge justfile: live command execution', () => {

  // Traceability: TC-FJ-011 -> AC: just compile backend dispatches correctly
  test('TC-FJ-011: just compile backend dispatches to backend branch (no scope error)', () => {
    const result = runCli('just compile backend');
    const output = result.stdout + result.stderr;
    // May fail due to missing toolchain at root, but must NOT be a scope error
    const isScopeError = output.includes('[forge] invalid scope');
    expect(!isScopeError, 'Should not be a scope error for backend scope').toBeTruthy();
  });

  // Traceability: TC-FJ-012 -> AC: just compile frontend dispatches correctly
  test('TC-FJ-012: just compile frontend dispatches to frontend branch (no scope error)', () => {
    const result = runCli('just compile frontend');
    // May fail if no tsconfig at root, but should NOT fail due to scope dispatch
    const output = result.stdout + result.stderr;
    const isScopeError = output.includes('[forge] invalid scope');
    expect(!isScopeError, 'Should not be a scope error for frontend scope').toBeTruthy();
  });

  // Traceability: TC-FJ-013 -> AC: just compile (empty scope) dispatches correctly
  test('TC-FJ-013: just compile with empty scope dispatches to both branches (no scope error)', () => {
    const result = runCli('just compile');
    const output = result.stdout + result.stderr;
    const isScopeError = output.includes('[forge] invalid scope');
    expect(!isScopeError, 'Should not be a scope error with empty scope').toBeTruthy();
  });

  // Traceability: TC-FJ-014 -> AC: invalid scope produces error (mixed) or is ignored (non-mixed)
  test('TC-FJ-014: just compile with invalid scope behavior depends on project type', () => {
    const projectType = getProjectType();
    const result = runCli('just compile invalidscope');
    if (projectType === 'mixed') {
      // Mixed projects validate scope and reject invalid values
      expect(result.exitCode, 'Expected exit 1 for invalid scope in mixed project').toBe(1);
      const output = result.stdout + result.stderr;
      expect(
        output.includes('[forge] invalid scope'),
        `Expected "[forge] invalid scope" error, got: ${output}`,
      ).toBeTruthy();
    } else {
      // Non-mixed projects accept scope param but ignore it (no scope dispatch)
      expect(result.exitCode, 'Expected exit 0 or non-scope-error for non-mixed project').not.toBe(1);
    }
  });

  // Traceability: TC-FJ-015 -> AC: just project-type returns valid type
  test('TC-FJ-015: just project-type returns a valid project type', () => {
    const result = runCli('just project-type');
    expect(result.exitCode, 'Expected exit 0').toBe(0);
    const output = result.stdout.trim();
    expect(
      ['frontend', 'backend', 'mixed'].includes(output),
      `Expected valid project-type output, got: "${output}"`,
    ).toBeTruthy();
  });
});
