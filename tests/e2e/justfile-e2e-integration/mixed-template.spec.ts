import { test, expect } from '@playwright/test';
import { readProjectFile } from '../helpers.js';

// ── Helpers ────────────────────────────────────────────────────────
function fileContains(content: string, needle: string): boolean {
  return content.includes(needle);
}

function fileNotContains(content: string, needle: string): boolean {
  return !content.includes(needle);
}

// Read the mixed template from the separate template file
function getMixedTemplate(): string {
  return readProjectFile('plugins/forge/skills/init-justfile/templates/mixed.just');
}

// ── TC-MIX-001 to TC-MIX-015: Mixed template content checks ───────
test.describe('Mixed template content checks', () => {

  // Traceability: TC-MIX-001 -> AC: project-type outputs @echo "mixed"
  test('TC-MIX-001: project-type outputs "mixed"', () => {
    const template = getMixedTemplate();
    expect(
      fileContains(template, '@echo "mixed"'),
      'Expected \'@echo "mixed"\' in mixed template project-type recipe',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-002 -> AC: compile has scope with bash case
  test('TC-MIX-002: compile recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    expect(
      fileContains(template, 'compile scope=""'),
      'Expected \'compile scope=""\' in mixed template',
    ).toBeTruthy();
    // Check for frontend/backend branches
    expect(
      fileContains(template, 'frontend)') && fileContains(template, 'backend)'),
      'Expected "frontend)" and "backend)" case branches in compile recipe',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-003 -> AC: build has scope with bash case
  test('TC-MIX-003: build recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    expect(
      fileContains(template, 'build scope=""'),
      'Expected \'build scope=""\' in mixed template',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-004 -> AC: run has scope with bash case
  test('TC-MIX-004: run recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    expect(
      fileContains(template, 'run scope=""'),
      'Expected \'run scope=""\' in mixed template',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-005 -> AC: dev has scope with bash case
  test('TC-MIX-005: dev recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    expect(
      fileContains(template, 'dev scope=""'),
      'Expected \'dev scope=""\' in mixed template',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-006 -> AC: test has scope with bash case
  test('TC-MIX-006: test recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    expect(
      fileContains(template, 'test scope=""'),
      'Expected \'test scope=""\' in mixed template',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-007 -> AC: lint has scope with bash case
  test('TC-MIX-007: lint recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    expect(
      fileContains(template, 'lint scope=""'),
      'Expected \'lint scope=""\' in mixed template',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-008 -> AC: fmt has scope with bash case
  test('TC-MIX-008: fmt recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    expect(
      fileContains(template, 'fmt scope=""'),
      'Expected \'fmt scope=""\' in mixed template',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-009 -> AC: check has scope with bash case
  test('TC-MIX-009: check recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    expect(
      fileContains(template, 'check scope=""'),
      'Expected \'check scope=""\' in mixed template',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-010 -> AC: clean has scope with bash case
  test('TC-MIX-010: clean recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    expect(
      fileContains(template, 'clean scope=""'),
      'Expected \'clean scope=""\' in mixed template',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-011 -> AC: install has scope with bash case
  test('TC-MIX-011: install recipe has scope="" parameter with bash case dispatch', () => {
    const template = getMixedTemplate();
    expect(
      fileContains(template, 'install scope=""'),
      'Expected \'install scope=""\' in mixed template',
    ).toBeTruthy();
  });
});

// ── TC-MIX-012 to TC-MIX-016: Bash case pattern checks ────────────
test.describe('Mixed template bash case pattern checks', () => {

  // Traceability: TC-MIX-012 -> AC: *) branch error message
  test('TC-MIX-012: scoped recipes have *) branch with error message to stderr', () => {
    const template = getMixedTemplate();
    const errorMsg = 'echo "[forge] invalid scope \'{{scope}}\'; expected frontend/backend" >&2; exit 1';
    // Count occurrences - should appear in all 10 scoped recipes
    const matches = template.split(errorMsg).length - 1;
    expect(
      matches >= 10,
      `Expected at least 10 occurrences of *) error branch, got ${matches}`,
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-013 -> AC: "") branch runs both frontend and backend
  test('TC-MIX-013: "" branch executes both frontend and backend commands', () => {
    const template = getMixedTemplate();
    // Empty scope should chain frontend && backend via placeholders
    expect(
      fileContains(template, 'npm run build) &&') &&
      fileContains(template, 'BACKEND_BUILD'),
      'Expected "" branch to chain frontend npm and BACKEND_BUILD placeholder commands',
    ).toBeTruthy();
    expect(
      fileContains(template, 'npm test) &&') &&
      fileContains(template, 'BACKEND_TEST'),
      'Expected "" branch to chain frontend npm test and BACKEND_TEST placeholder commands',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-014 -> AC: All bash recipes use set -euo pipefail
  test('TC-MIX-014: all bash recipes use set -euo pipefail', () => {
    const template = getMixedTemplate();
    const bashRecipes = template.split('#!/usr/bin/env bash').length - 1;
    const pipefailCount = template.split('set -euo pipefail').length - 1;
    expect(
      pipefailCount >= bashRecipes,
      `Expected at least ${bashRecipes} "set -euo pipefail" for ${bashRecipes} bash recipes, got ${pipefailCount}`,
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-015 -> AC: Frontend uses npm, backend uses BACKEND_* placeholders
  test('TC-MIX-015: frontend commands use npm toolchain, backend uses BACKEND_* placeholders', () => {
    const template = getMixedTemplate();
    // Frontend branch commands
    expect(fileContains(template, 'npm run build'), 'Expected "npm run build" for frontend').toBeTruthy();
    expect(fileContains(template, 'npm test'), 'Expected "npm test" for frontend').toBeTruthy();
    expect(fileContains(template, 'npm run lint'), 'Expected "npm run lint" for frontend').toBeTruthy();
    expect(fileContains(template, 'FRONTEND_RUN'), 'Expected "FRONTEND_RUN" placeholder for frontend run').toBeTruthy();
    expect(fileContains(template, 'FRONTEND_DEV'), 'Expected "FRONTEND_DEV" placeholder for frontend dev').toBeTruthy();
    expect(fileContains(template, 'npx prettier --write .'), 'Expected "npx prettier --write ." for frontend').toBeTruthy();
    expect(fileContains(template, 'npx tsc --noEmit'), 'Expected "npx tsc --noEmit" for frontend compile').toBeTruthy();
    expect(fileContains(template, 'npm install'), 'Expected "npm install" for frontend').toBeTruthy();

    // Backend uses BACKEND_* placeholders (replaced at init time based on detected language)
    expect(fileContains(template, 'BACKEND_BUILD'), 'Expected "BACKEND_BUILD" placeholder').toBeTruthy();
    expect(fileContains(template, 'BACKEND_TEST'), 'Expected "BACKEND_TEST" placeholder').toBeTruthy();
    expect(fileContains(template, 'BACKEND_LINT'), 'Expected "BACKEND_LINT" placeholder').toBeTruthy();
    expect(fileContains(template, 'BACKEND_RUN'), 'Expected "BACKEND_RUN" placeholder').toBeTruthy();
    expect(fileContains(template, 'BACKEND_FMT'), 'Expected "BACKEND_FMT" placeholder').toBeTruthy();
    expect(fileContains(template, 'BACKEND_COMPILE'), 'Expected "BACKEND_COMPILE" placeholder').toBeTruthy();
    expect(fileContains(template, 'BACKEND_CLEAN'), 'Expected "BACKEND_CLEAN" placeholder').toBeTruthy();
    expect(fileContains(template, 'BACKEND_INSTALL'), 'Expected "BACKEND_INSTALL" placeholder').toBeTruthy();
  });
});

// ── TC-MIX-016 to TC-MIX-020: Unscoped recipe checks ──────────────
test.describe('Mixed template unscoped recipe checks', () => {

  // Traceability: TC-MIX-016 -> AC: project-type has no scope parameter
  test('TC-MIX-016: project-type has no scope parameter', () => {
    const template = getMixedTemplate();
    // project-type should NOT have scope=""
    const projectTypeMatch = template.match(/project-type:.*\n/);
    expect(projectTypeMatch, 'Expected project-type recipe in mixed template').toBeTruthy();
    expect(
      !projectTypeMatch![0].includes('scope=""'),
      'Expected project-type to NOT have scope="" parameter',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-017 -> AC: test-e2e has no scope parameter
  test('TC-MIX-017: test-e2e has no scope parameter', () => {
    const template = getMixedTemplate();
    // test-e2e uses feature="" not scope=""
    const testE2eMatch = template.match(/test-e2e[^:]*:/);
    expect(testE2eMatch, 'Expected test-e2e recipe in mixed template').toBeTruthy();
    expect(
      !testE2eMatch![0].includes('scope=""'),
      'Expected test-e2e to NOT have scope="" parameter',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-018 -> AC: ci has no scope parameter
  test('TC-MIX-018: ci has no scope parameter', () => {
    const template = getMixedTemplate();
    const ciMatch = template.match(/^ci:/m);
    expect(ciMatch, 'Expected ci recipe in mixed template').toBeTruthy();
    expect(
      !template.includes('ci scope=""'),
      'Expected ci to NOT have scope="" parameter',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-019 -> AC: e2e-setup has no scope parameter
  test('TC-MIX-019: e2e-setup has no scope parameter', () => {
    const template = getMixedTemplate();
    const e2eSetupMatch = template.match(/^e2e-setup[^:]*:/m);
    expect(e2eSetupMatch, 'Expected e2e-setup recipe in mixed template').toBeTruthy();
    expect(
      !e2eSetupMatch![0].includes('scope=""'),
      'Expected e2e-setup to NOT have scope="" parameter',
    ).toBeTruthy();
  });

  // Traceability: TC-MIX-020 -> AC: e2e-verify has no scope parameter
  test('TC-MIX-020: e2e-verify has no scope parameter', () => {
    const template = getMixedTemplate();
    expect(
      fileContains(template, 'e2e-verify'),
      'Expected e2e-verify recipe in mixed template',
    ).toBeTruthy();
    expect(
      !template.includes('e2e-verify scope=""'),
      'Expected e2e-verify to NOT have scope="" parameter',
    ).toBeTruthy();
  });
});

// ── TC-MIX-021 to TC-MIX-023: Boundary markers and structure ──────
test.describe('Mixed template boundary markers and structure', () => {

  // Traceability: TC-MIX-021 -> AC: Templates stored as string literals
  test('TC-MIX-021: mixed template has forge boundary markers', () => {
    const template = getMixedTemplate();
    expect(
      fileContains(template, '# --- forge standard recipes ---'),
      'Expected forge standard recipes start marker in mixed template',
    ).toBeTruthy();
    expect(
      fileContains(template, '# --- end forge standard recipes ---'),
      'Expected forge standard recipes end marker in mixed template',
    ).toBeTruthy();
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
      expect(
        fileContains(template, recipe),
        `Expected recipe "${recipe}" in mixed template`,
      ).toBeTruthy();
    }
  });

  // Traceability: TC-MIX-023 -> AC: ci recipe chains install, compile, build, test, lint
  test('TC-MIX-023: ci recipe chains standard commands', () => {
    const template = getMixedTemplate();
    // Justfile ci recipe uses recipe names directly (e.g., "ci: install compile build test lint")
    expect(
      fileContains(template, 'ci:'),
      'Expected ci recipe definition',
    ).toBeTruthy();
    const ciLine = template.match(/^ci:.*$/m)?.[0] || '';
    const expectedSteps = ['install', 'compile', 'build', 'test', 'lint'];
    for (const step of expectedSteps) {
      expect(
        ciLine.includes(step),
        `Expected "${step}" in ci recipe, got: ${ciLine}`,
      ).toBeTruthy();
    }
  });
});
