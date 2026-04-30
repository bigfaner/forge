import { describe, test } from 'node:test';
import assert from 'node:assert/strict';
import { readProjectFile } from '../helpers.js';

// ── Helpers ────────────────────────────────────────────────────────
function fileContains(content: string, needle: string): boolean {
  return content.includes(needle);
}

function fileNotContains(content: string, needle: string): boolean {
  return !content.includes(needle);
}

// Extract the mixed template section from init-justfile.md
function getMixedTemplate(): string {
  const content = readProjectFile('plugins/forge/commands/init-justfile.md');
  const startMarker = '### Mixed Template';
  const startIdx = content.indexOf(startMarker);
  assert.notEqual(startIdx, -1, 'Expected "### Mixed Template" section in init-justfile.md');
  // Return content from the section header to the next ## or end
  const afterStart = content.slice(startIdx);
  const nextSection = afterStart.indexOf('\n## ', 1);
  return nextSection !== -1 ? afterStart.slice(0, nextSection) : afterStart;
}

// ── TC-MIX-001 to TC-MIX-015: Mixed template content checks ───────
describe('Mixed template content checks', () => {

  // Traceability: TC-MIX-001 -> AC: project-type outputs @echo "mixed"
  test('TC-MIX-001: project-type outputs "mixed"', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, '@echo "mixed"'),
      'Expected \'@echo "mixed"\' in mixed template project-type recipe',
    );
  });

  // Traceability: TC-MIX-002 -> AC: compile has scope with bash case
  test('TC-MIX-002: compile recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, 'compile scope=""'),
      'Expected \'compile scope=""\' in mixed template',
    );
    // Check for frontend/backend branches
    assert.ok(
      fileContains(template, 'frontend)') && fileContains(template, 'backend)'),
      'Expected "frontend)" and "backend)" case branches in compile recipe',
    );
  });

  // Traceability: TC-MIX-003 -> AC: build has scope with bash case
  test('TC-MIX-003: build recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, 'build scope=""'),
      'Expected \'build scope=""\' in mixed template',
    );
  });

  // Traceability: TC-MIX-004 -> AC: run has scope with bash case
  test('TC-MIX-004: run recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, 'run scope=""'),
      'Expected \'run scope=""\' in mixed template',
    );
  });

  // Traceability: TC-MIX-005 -> AC: dev has scope with bash case
  test('TC-MIX-005: dev recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, 'dev scope=""'),
      'Expected \'dev scope=""\' in mixed template',
    );
  });

  // Traceability: TC-MIX-006 -> AC: test has scope with bash case
  test('TC-MIX-006: test recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, 'test scope=""'),
      'Expected \'test scope=""\' in mixed template',
    );
  });

  // Traceability: TC-MIX-007 -> AC: lint has scope with bash case
  test('TC-MIX-007: lint recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, 'lint scope=""'),
      'Expected \'lint scope=""\' in mixed template',
    );
  });

  // Traceability: TC-MIX-008 -> AC: fmt has scope with bash case
  test('TC-MIX-008: fmt recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, 'fmt scope=""'),
      'Expected \'fmt scope=""\' in mixed template',
    );
  });

  // Traceability: TC-MIX-009 -> AC: check has scope with bash case
  test('TC-MIX-009: check recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, 'check scope=""'),
      'Expected \'check scope=""\' in mixed template',
    );
  });

  // Traceability: TC-MIX-010 -> AC: clean has scope with bash case
  test('TC-MIX-010: clean recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, 'clean scope=""'),
      'Expected \'clean scope=""\' in mixed template',
    );
  });

  // Traceability: TC-MIX-011 -> AC: install has scope with bash case
  test('TC-MIX-011: install recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, 'install scope=""'),
      'Expected \'install scope=""\' in mixed template',
    );
  });
});

// ── TC-MIX-012 to TC-MIX-016: Bash case pattern checks ────────────
describe('Mixed template bash case pattern checks', () => {

  // Traceability: TC-MIX-012 -> AC: *) branch error message
  test('TC-MIX-012: scoped recipes have *) branch with error message to stderr', () => {
    const template = getMixedTemplate();
    const errorMsg = 'echo "[forge] invalid scope \'{{scope}}\'; expected frontend/backend" >&2; exit 1';
    // Count occurrences - should appear in all 10 scoped recipes
    const matches = template.split(errorMsg).length - 1;
    assert.ok(
      matches >= 10,
      `Expected at least 10 occurrences of *) error branch, got ${matches}`,
    );
  });

  // Traceability: TC-MIX-013 -> AC: "") branch runs both frontend and backend
  test('TC-MIX-013: "" branch executes both frontend and backend commands', () => {
    const template = getMixedTemplate();
    // Empty scope should chain frontend && backend
    assert.ok(
      fileContains(template, 'npm run build && go build ./...'),
      'Expected "" branch to chain npm and go build commands',
    );
    assert.ok(
      fileContains(template, 'npm test && go test -race ./...'),
      'Expected "" branch to chain npm and go test commands',
    );
  });

  // Traceability: TC-MIX-014 -> AC: All bash recipes use set -euo pipefail
  test('TC-MIX-014: all bash recipes use set -euo pipefail', () => {
    const template = getMixedTemplate();
    const bashRecipes = template.split('#!/usr/bin/env bash').length - 1;
    const pipefailCount = template.split('set -euo pipefail').length - 1;
    assert.ok(
      pipefailCount >= bashRecipes,
      `Expected at least ${bashRecipes} "set -euo pipefail" for ${bashRecipes} bash recipes, got ${pipefailCount}`,
    );
  });

  // Traceability: TC-MIX-015 -> AC: Frontend uses npm, backend uses Go
  test('TC-MIX-015: frontend commands use npm toolchain, backend uses Go toolchain', () => {
    const template = getMixedTemplate();
    // Frontend branch commands
    assert.ok(fileContains(template, 'npm run build'), 'Expected "npm run build" for frontend');
    assert.ok(fileContains(template, 'npm test'), 'Expected "npm test" for frontend');
    assert.ok(fileContains(template, 'npm run lint'), 'Expected "npm run lint" for frontend');
    assert.ok(fileContains(template, 'npm start'), 'Expected "npm start" for frontend');
    assert.ok(fileContains(template, 'npm run dev'), 'Expected "npm run dev" for frontend');
    assert.ok(fileContains(template, 'npx prettier --write .'), 'Expected "npx prettier --write ." for frontend');
    assert.ok(fileContains(template, 'npx tsc --noEmit'), 'Expected "npx tsc --noEmit" for frontend compile');
    assert.ok(fileContains(template, 'npm install'), 'Expected "npm install" for frontend');

    // Backend branch commands
    assert.ok(fileContains(template, 'go build ./...'), 'Expected "go build ./..." for backend');
    assert.ok(fileContains(template, 'go test -race ./...'), 'Expected "go test -race ./..." for backend');
    assert.ok(fileContains(template, 'golangci-lint run ./...'), 'Expected "golangci-lint run ./..." for backend');
    assert.ok(fileContains(template, 'go run .'), 'Expected "go run ." for backend');
    assert.ok(fileContains(template, 'gofmt -w .'), 'Expected "gofmt -w ." for backend');
    assert.ok(fileContains(template, 'go vet ./...'), 'Expected "go vet ./..." for backend');
    assert.ok(fileContains(template, 'go clean ./...'), 'Expected "go clean ./..." for backend');
    assert.ok(fileContains(template, 'go mod download'), 'Expected "go mod download" for backend');
  });
});

// ── TC-MIX-016 to TC-MIX-020: Unscoped recipe checks ──────────────
describe('Mixed template unscoped recipe checks', () => {

  // Traceability: TC-MIX-016 -> AC: project-type has no scope parameter
  test('TC-MIX-016: project-type has no scope parameter', () => {
    const template = getMixedTemplate();
    // project-type should NOT have scope=""
    const projectTypeMatch = template.match(/project-type:.*\n/);
    assert.ok(projectTypeMatch, 'Expected project-type recipe in mixed template');
    assert.ok(
      !projectTypeMatch[0].includes('scope=""'),
      'Expected project-type to NOT have scope="" parameter',
    );
  });

  // Traceability: TC-MIX-017 -> AC: test-e2e has no scope parameter
  test('TC-MIX-017: test-e2e has no scope parameter', () => {
    const template = getMixedTemplate();
    // test-e2e uses feature="" not scope=""
    const testE2eMatch = template.match(/test-e2e[^:]*:/);
    assert.ok(testE2eMatch, 'Expected test-e2e recipe in mixed template');
    assert.ok(
      !testE2eMatch[0].includes('scope=""'),
      'Expected test-e2e to NOT have scope="" parameter',
    );
  });

  // Traceability: TC-MIX-018 -> AC: ci has no scope parameter
  test('TC-MIX-018: ci has no scope parameter', () => {
    const template = getMixedTemplate();
    const ciMatch = template.match(/^ci:/m);
    assert.ok(ciMatch, 'Expected ci recipe in mixed template');
    assert.ok(
      !template.includes('ci scope=""'),
      'Expected ci to NOT have scope="" parameter',
    );
  });

  // Traceability: TC-MIX-019 -> AC: e2e-setup has no scope parameter
  test('TC-MIX-019: e2e-setup has no scope parameter', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, 'e2e-setup:'),
      'Expected e2e-setup recipe in mixed template',
    );
    assert.ok(
      !template.includes('e2e-setup scope=""'),
      'Expected e2e-setup to NOT have scope="" parameter',
    );
  });

  // Traceability: TC-MIX-020 -> AC: e2e-verify has no scope parameter
  test('TC-MIX-020: e2e-verify has no scope parameter', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, 'e2e-verify'),
      'Expected e2e-verify recipe in mixed template',
    );
    assert.ok(
      !template.includes('e2e-verify scope=""'),
      'Expected e2e-verify to NOT have scope="" parameter',
    );
  });
});

// ── TC-MIX-021 to TC-MIX-023: Boundary markers and structure ──────
describe('Mixed template boundary markers and structure', () => {

  // Traceability: TC-MIX-021 -> AC: Templates stored as string literals
  test('TC-MIX-021: mixed template has forge boundary markers', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, '# --- forge standard recipes ---'),
      'Expected forge standard recipes start marker in mixed template',
    );
    assert.ok(
      fileContains(template, '# --- end forge standard recipes ---'),
      'Expected forge standard recipes end marker in mixed template',
    );
  });

  // Traceability: TC-MIX-022 -> AC: All 15 recipes present
  test('TC-MIX-022: all 15 recipes are present in mixed template', () => {
    const template = getMixedTemplate();
    const expectedRecipes = [
      'project-type', 'compile', 'build', 'run', 'dev',
      'test', 'test-e2e', 'lint', 'fmt', 'check',
      'clean', 'install', 'ci', 'e2e-setup', 'e2e-verify',
    ];
    for (const recipe of expectedRecipes) {
      assert.ok(
        fileContains(template, recipe),
        `Expected recipe "${recipe}" in mixed template`,
      );
    }
  });

  // Traceability: TC-MIX-023 -> AC: ci recipe chains install, compile, build, test, lint
  test('TC-MIX-023: ci recipe chains standard commands', () => {
    const template = getMixedTemplate();
    assert.ok(
      fileContains(template, 'just install'),
      'Expected "just install" in ci recipe',
    );
    assert.ok(
      fileContains(template, 'just compile'),
      'Expected "just compile" in ci recipe',
    );
    assert.ok(
      fileContains(template, 'just build'),
      'Expected "just build" in ci recipe',
    );
    assert.ok(
      fileContains(template, 'just test'),
      'Expected "just test" in ci recipe',
    );
    assert.ok(
      fileContains(template, 'just lint'),
      'Expected "just lint" in ci recipe',
    );
  });
});
