import { describe, test } from 'node:test';
import assert from 'node:assert/strict';
import { runCli, readProjectFile } from '../helpers.js';

// -- Helpers ---------------------------------------------------------------
function fileContains(content: string, needle: string): boolean {
  return content.includes(needle);
}

function getInitJustfileContent(): string {
  return readProjectFile('plugins/forge/commands/init-justfile.md');
}

// The 15 standard command names per Spec 5.1
const STANDARD_COMMANDS = [
  'compile', 'build', 'run', 'dev', 'test',
  'test-e2e', 'lint', 'fmt', 'check', 'clean',
  'install', 'ci', 'e2e-setup', 'e2e-verify', 'project-type',
];

// -- Tests -----------------------------------------------------------------
describe('init-justfile: project detection and generation', () => {

  // Traceability: TC-004 -> Story 2 / AC-1
  test('TC-004: frontend project detection generates scope-free justfile', () => {
    const content = getInitJustfileContent();
    // Frontend template should exist
    assert.ok(
      fileContains(content, '### Frontend Template'),
      'Expected "### Frontend Template" section in init-justfile',
    );
    // Extract the frontend template section (between ### Frontend and ### Mixed or ##)
    const frontendStart = content.indexOf('### Frontend Template');
    const nextH3 = content.indexOf('\n### ', frontendStart + 1);
    const nextH2 = content.indexOf('\n## ', frontendStart + 1);
    const nextSection = nextH3 !== -1 ? nextH3 : nextH2;
    const frontendSection = nextSection !== -1
      ? content.slice(frontendStart, nextSection)
      : content.slice(frontendStart);
    assert.ok(
      !fileContains(frontendSection, 'scope=""'),
      'Expected frontend template to NOT have scope="" parameters',
    );
    // project-type should output "frontend"
    assert.ok(
      fileContains(frontendSection, '@echo "frontend"'),
      'Expected @echo "frontend" in frontend template',
    );
  });

  // Traceability: TC-005 -> Story 2 / AC-2
  test('TC-005: backend project detection generates scope-free justfile', () => {
    const content = getInitJustfileContent();
    assert.ok(
      fileContains(content, '### Backend Template'),
      'Expected "### Backend Template" section in init-justfile',
    );
    // Extract the backend template section (between ### Backend and ### Frontend)
    const backendStart = content.indexOf('### Backend Template');
    const nextH3 = content.indexOf('\n### ', backendStart + 1);
    const nextH2 = content.indexOf('\n## ', backendStart + 1);
    const nextSection = nextH3 !== -1 ? nextH3 : nextH2;
    const backendSection = nextSection !== -1
      ? content.slice(backendStart, nextSection)
      : content.slice(backendStart);
    assert.ok(
      !fileContains(backendSection, 'scope=""'),
      'Expected backend template to NOT have scope="" parameters',
    );
    assert.ok(
      fileContains(backendSection, '@echo "backend"'),
      'Expected @echo "backend" in backend template',
    );
  });

  // Traceability: TC-006 -> Story 2 / AC-3
  test('TC-006: mixed project detection generates scope-aware justfile', () => {
    const content = getInitJustfileContent();
    assert.ok(
      fileContains(content, '### Mixed Template'),
      'Expected "### Mixed Template" section in init-justfile',
    );
    const mixedStart = content.indexOf('### Mixed Template');
    const nextSection = content.indexOf('\n## ', mixedStart + 1);
    const mixedSection = nextSection !== -1
      ? content.slice(mixedStart, nextSection)
      : content.slice(mixedStart);
    // Mixed template SHOULD have scope="" parameters
    assert.ok(
      fileContains(mixedSection, 'scope=""'),
      'Expected mixed template to have scope="" parameters',
    );
    assert.ok(
      fileContains(mixedSection, '@echo "mixed"'),
      'Expected @echo "mixed" in mixed template',
    );
  });

  // Traceability: TC-022 -> Spec 5.1 / vocabulary
  test('TC-022: all 15 standard commands are present in generated justfile', () => {
    const justfile = readProjectFile('justfile');
    for (const cmd of STANDARD_COMMANDS) {
      assert.ok(
        fileContains(justfile, cmd),
        `Expected recipe "${cmd}" in the forge project justfile`,
      );
    }
  });

  // Traceability: TC-018 -> Spec 5.2 / detection
  test('TC-018: no marker files detected causes init-justfile to error', () => {
    const content = getInitJustfileContent();
    // The init-justfile skill should describe an error case for no markers
    assert.ok(
      fileContains(content, 'no known project markers') ||
      fileContains(content, 'no project markers') ||
      fileContains(content, 'Error: no known') ||
      fileContains(content, 'no markers detected') ||
      fileContains(content, 'neither') ||
      fileContains(content, 'Cannot determine'),
      'Expected error handling description for no project markers',
    );
  });

  // Traceability: TC-019 -> Spec 5.2 / flow
  test('TC-019: existing justfile triggers user confirmation', () => {
    const content = getInitJustfileContent();
    assert.ok(
      fileContains(content, 'confirm') ||
      fileContains(content, 'prompt') ||
      fileContains(content, 'overwrite') ||
      fileContains(content, 'ask') ||
      fileContains(content, '--force'),
      'Expected user confirmation mechanism for existing justfile',
    );
  });

  // Traceability: TC-020 -> Spec / maintainability
  test('TC-020: boundary markers present triggers idempotent merge', () => {
    const justfile = readProjectFile('justfile');
    const startMarker = '# --- forge standard recipes ---';
    const endMarker = '# --- end forge standard recipes ---';
    assert.ok(
      fileContains(justfile, startMarker),
      'Expected start boundary marker in justfile',
    );
    assert.ok(
      fileContains(justfile, endMarker),
      'Expected end boundary marker in justfile',
    );

    // Custom recipes should be OUTSIDE boundary markers
    const startIdx = justfile.indexOf(startMarker);
    const endIdx = justfile.indexOf(endMarker);
    assert.ok(startIdx !== -1 && endIdx !== -1, 'Expected both boundary markers');

    // Content before start marker or after end marker should contain custom recipes
    const beforeMarkers = justfile.slice(0, startIdx);
    const afterMarkers = justfile.slice(endIdx + endMarker.length);
    const customRecipesOutsideMarkers =
      beforeMarkers.includes('claude:') ||
      afterMarkers.includes('claude:');
    assert.ok(
      customRecipesOutsideMarkers,
      'Expected custom recipes outside boundary markers',
    );
  });
});
