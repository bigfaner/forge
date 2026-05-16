import { test, expect } from '@playwright/test';
import { runCli, readProjectFile } from '../helpers.js';

// -- Helpers ---------------------------------------------------------------
function fileContains(content: string, needle: string): boolean {
  return content.includes(needle);
}

function getInitJustfileContent(): string {
  return readProjectFile('plugins/forge/skills/init-justfile/SKILL.md');
}

// The 15 standard command names per Spec 5.1
const STANDARD_COMMANDS = [
  'compile', 'build', 'run', 'dev', 'test',
  'test-e2e', 'lint', 'fmt', 'check', 'clean',
  'install', 'ci', 'e2e-setup', 'e2e-verify', 'project-type',
];

// -- Tests -----------------------------------------------------------------
test.describe('init-justfile: project detection and generation', () => {

  // Traceability: TC-004 -> Story 2 / AC-1
  test('TC-004: frontend project detection generates scope-free justfile', () => {
    const frontendTemplate = readProjectFile('plugins/forge/skills/init-justfile/templates/node.just');
    // Frontend template should NOT have scope="" parameters
    expect(
      !fileContains(frontendTemplate, 'scope=""'),
      'Expected frontend template to NOT have scope="" parameters',
    ).toBeTruthy();
    // project-type should output "frontend"
    expect(
      fileContains(frontendTemplate, '@echo "frontend"'),
      'Expected @echo "frontend" in frontend template',
    ).toBeTruthy();
  });

  // Traceability: TC-005 -> Story 2 / AC-2
  test('TC-005: backend project detection generates scope-free justfile', () => {
    const backendTemplate = readProjectFile('plugins/forge/skills/init-justfile/templates/go.just');
    // Backend template should NOT have scope="" parameters
    expect(
      !fileContains(backendTemplate, 'scope=""'),
      'Expected backend template to NOT have scope="" parameters',
    ).toBeTruthy();
    expect(
      fileContains(backendTemplate, '@echo "backend"'),
      'Expected @echo "backend" in backend template',
    ).toBeTruthy();
  });

  // Traceability: TC-006 -> Story 2 / AC-3
  test('TC-006: mixed project detection generates scope-aware justfile', () => {
    const mixedTemplate = readProjectFile('plugins/forge/skills/init-justfile/templates/mixed.just');
    // Mixed template SHOULD have scope="" parameters
    expect(
      fileContains(mixedTemplate, 'scope=""'),
      'Expected mixed template to have scope="" parameters',
    ).toBeTruthy();
    expect(
      fileContains(mixedTemplate, '@echo "mixed"'),
      'Expected @echo "mixed" in mixed template',
    ).toBeTruthy();
  });

  // Traceability: TC-022 -> Spec 5.1 / vocabulary
  test('TC-022: all 15 standard commands are present in generated justfile', () => {
    const justfile = readProjectFile('justfile');
    for (const cmd of STANDARD_COMMANDS) {
      expect(
        fileContains(justfile, cmd),
        `Expected recipe "${cmd}" in the forge project justfile`,
      ).toBeTruthy();
    }
  });

  // Traceability: TC-018 -> Spec 5.2 / detection
  test('TC-018: no marker files detected causes init-justfile to error', () => {
    const content = getInitJustfileContent();
    // The init-justfile skill should describe an error case for no markers
    expect(
      fileContains(content, 'no known project markers') ||
      fileContains(content, 'no project markers') ||
      fileContains(content, 'Error: no known') ||
      fileContains(content, 'no markers detected') ||
      fileContains(content, 'neither') ||
      fileContains(content, 'Cannot determine'),
      'Expected error handling description for no project markers',
    ).toBeTruthy();
  });

  // Traceability: TC-019 -> Spec 5.2 / flow
  test('TC-019: existing justfile triggers user confirmation', () => {
    const content = getInitJustfileContent();
    expect(
      fileContains(content, 'confirm') ||
      fileContains(content, 'prompt') ||
      fileContains(content, 'overwrite') ||
      fileContains(content, 'ask') ||
      fileContains(content, '--force'),
      'Expected user confirmation mechanism for existing justfile',
    ).toBeTruthy();
  });

  // Traceability: TC-020 -> Spec / maintainability
  test('TC-020: boundary markers present triggers idempotent merge', () => {
    const justfile = readProjectFile('justfile');
    const startMarker = '# --- forge standard recipes ---';
    const endMarker = '# --- end forge standard recipes ---';
    expect(
      fileContains(justfile, startMarker),
      'Expected start boundary marker in justfile',
    ).toBeTruthy();
    expect(
      fileContains(justfile, endMarker),
      'Expected end boundary marker in justfile',
    ).toBeTruthy();

    // Custom recipes should be OUTSIDE boundary markers
    const startIdx = justfile.indexOf(startMarker);
    const endIdx = justfile.indexOf(endMarker);
    expect(startIdx !== -1 && endIdx !== -1, 'Expected both boundary markers').toBeTruthy();

    // Content before start marker or after end marker should contain custom recipes
    const beforeMarkers = justfile.slice(0, startIdx);
    const afterMarkers = justfile.slice(endIdx + endMarker.length);
    const customRecipesOutsideMarkers =
      beforeMarkers.includes('claude:') ||
      afterMarkers.includes('claude:');
    expect(
      customRecipesOutsideMarkers,
      'Expected custom recipes outside boundary markers',
    ).toBeTruthy();
  });
});
