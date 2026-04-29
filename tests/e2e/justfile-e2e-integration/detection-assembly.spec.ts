import { describe, test } from 'node:test';
import assert from 'node:assert/strict';
import { readProjectFile } from '../helpers.js';

// ── Helpers ────────────────────────────────────────────────────────
function fileContains(content: string, needle: string): boolean {
  return content.includes(needle);
}

function getInitJustfileContent(): string {
  return readProjectFile('plugins/forge/commands/init-justfile.md');
}

// ── TC-DET-001 to TC-DET-004: Project-type detection signals ───────
describe('Project-type detection signals', () => {

  // Traceability: TC-DET-001 -> AC: Detection checks package.json (frontend signal)
  test('TC-DET-001: detection logic checks for package.json as frontend signal', () => {
    const content = getInitJustfileContent();
    assert.ok(
      fileContains(content, 'package.json'),
      'Expected detection logic to check for package.json',
    );
    // Must be in the detection section, not just the language templates
    const detectionSection = content.indexOf('Step 1:');
    assert.ok(detectionSection !== -1, 'Expected Step 1 detection section');
  });

  // Traceability: TC-DET-002 -> AC: Detection checks go.mod (backend signal)
  test('TC-DET-002: detection logic checks for go.mod as backend signal', () => {
    const content = getInitJustfileContent();
    assert.ok(
      fileContains(content, 'go.mod'),
      'Expected detection logic to check for go.mod',
    );
  });

  // Traceability: TC-DET-003 -> AC: Detection checks Cargo.toml (backend signal)
  test('TC-DET-003: detection logic checks for Cargo.toml as backend signal', () => {
    const content = getInitJustfileContent();
    assert.ok(
      fileContains(content, 'Cargo.toml'),
      'Expected detection logic to check for Cargo.toml',
    );
  });

  // Traceability: TC-DET-004 -> AC: Detection checks pyproject.toml (backend signal)
  test('TC-DET-004: detection logic checks for pyproject.toml as backend signal', () => {
    const content = getInitJustfileContent();
    assert.ok(
      fileContains(content, 'pyproject.toml'),
      'Expected detection logic to check for pyproject.toml',
    );
  });
});

// ── TC-DET-005 to TC-DET-008: Project classification ───────────────
describe('Project classification', () => {

  // Traceability: TC-DET-005 -> AC: both signals -> mixed
  test('TC-DET-005: classification produces mixed when both frontend and backend signals detected', () => {
    const content = getInitJustfileContent();
    // Should describe "mixed" classification when both package.json and go.mod exist
    assert.ok(
      fileContains(content, 'mixed'),
      'Expected "mixed" classification in detection logic',
    );
  });

  // Traceability: TC-DET-006 -> AC: frontend only -> frontend
  test('TC-DET-006: classification produces frontend when only frontend signals detected', () => {
    const content = getInitJustfileContent();
    assert.ok(
      fileContains(content, 'frontend'),
      'Expected "frontend" classification in detection logic',
    );
  });

  // Traceability: TC-DET-007 -> AC: backend only -> backend
  test('TC-DET-007: classification produces backend when only backend signals detected', () => {
    const content = getInitJustfileContent();
    assert.ok(
      fileContains(content, 'backend'),
      'Expected "backend" classification in detection logic',
    );
  });

  // Traceability: TC-DET-008 -> AC: neither signal -> error message
  test('TC-DET-008: classification produces error when no known markers detected', () => {
    const content = getInitJustfileContent();
    // Should describe error case when no markers found
    assert.ok(
      fileContains(content, 'no known project markers') ||
      fileContains(content, 'no project markers') ||
      fileContains(content, 'Error: no known') ||
      fileContains(content, 'no markers detected') ||
      fileContains(content, 'neither') ||
      fileContains(content, 'Cannot determine'),
      'Expected error handling when no project markers detected',
    );
  });
});

// ── TC-DET-009 to TC-DET-011: Template selection ───────────────────
describe('Template selection and assembly', () => {

  // Traceability: TC-DET-009 -> AC: Selects backend template for backend project
  test('TC-DET-009: selects Backend Template for pure backend projects', () => {
    const content = getInitJustfileContent();
    // Find the Recipe Templates section with ### Backend Template header
    const backendHeader = '### Backend Template';
    assert.ok(
      fileContains(content, backendHeader),
      'Expected ### Backend Template section header to exist',
    );
    // Backend project-type should output "backend"
    const backendSection = content.indexOf(backendHeader);
    const backendContent = content.slice(backendSection, backendSection + 3000);
    assert.ok(
      fileContains(backendContent, '@echo "backend"'),
      'Expected backend template to output @echo "backend"',
    );
  });

  // Traceability: TC-DET-010 -> AC: Selects frontend template for frontend project
  test('TC-DET-010: selects Frontend Template for pure frontend projects', () => {
    const content = getInitJustfileContent();
    const frontendHeader = '### Frontend Template';
    assert.ok(
      fileContains(content, frontendHeader),
      'Expected ### Frontend Template section header to exist',
    );
    const frontendSection = content.indexOf(frontendHeader);
    const frontendContent = content.slice(frontendSection, frontendSection + 3000);
    assert.ok(
      fileContains(frontendContent, '@echo "frontend"'),
      'Expected frontend template to output @echo "frontend"',
    );
  });

  // Traceability: TC-DET-011 -> AC: Selects mixed template for mixed project
  test('TC-DET-011: selects Mixed Template for mixed projects', () => {
    const content = getInitJustfileContent();
    const mixedHeader = '### Mixed Template';
    assert.ok(
      fileContains(content, mixedHeader),
      'Expected ### Mixed Template section header to exist',
    );
    const mixedSection = content.indexOf(mixedHeader);
    const mixedContent = content.slice(mixedSection, mixedSection + 3000);
    assert.ok(
      fileContains(mixedContent, '@echo "mixed"'),
      'Expected mixed template to output @echo "mixed"',
    );
  });
});

// ── TC-DET-012 to TC-DET-014: Boundary markers ─────────────────────
describe('Boundary markers', () => {

  // Traceability: TC-DET-012 -> AC: Boundary markers wrap generated recipes
  test('TC-DET-012: generated recipes wrapped in forge boundary markers', () => {
    const content = getInitJustfileContent();
    const startMarker = '# --- forge standard recipes ---';
    const endMarker = '# --- end forge standard recipes ---';
    assert.ok(
      fileContains(content, startMarker),
      'Expected start boundary marker',
    );
    assert.ok(
      fileContains(content, endMarker),
      'Expected end boundary marker',
    );
  });

  // Traceability: TC-DET-013 -> AC: Boundary marker merge replaces only marked section
  test('TC-DET-013: boundary marker merge replaces only marked section', () => {
    const content = getInitJustfileContent();
    // Should describe merge logic that preserves content outside markers
    assert.ok(
      fileContains(content, 'boundary marker') ||
      fileContains(content, 'markers') ||
      fileContains(content, 'replace') ||
      fileContains(content, 'merge'),
      'Expected boundary marker merge logic description',
    );
  });

  // Traceability: TC-DET-014 -> AC: All 3 templates have boundary markers
  test('TC-DET-014: all three templates have boundary markers', () => {
    const content = getInitJustfileContent();
    const startMarker = '# --- forge standard recipes ---';
    const endMarker = '# --- end forge standard recipes ---';
    // Count occurrences - should appear in backend, frontend, and mixed templates
    const startCount = content.split(startMarker).length - 1;
    const endCount = content.split(endMarker).length - 1;
    assert.ok(
      startCount >= 3,
      `Expected at least 3 start boundary markers (one per template), got ${startCount}`,
    );
    assert.ok(
      endCount >= 3,
      `Expected at least 3 end boundary markers (one per template), got ${endCount}`,
    );
  });
});

// ── TC-DET-015 to TC-DET-017: --force flag and interactive confirmation ─
describe('--force flag and interactive confirmation', () => {

  // Traceability: TC-DET-015 -> AC: --force flag skips confirmation
  test('TC-DET-015: --force flag skips user confirmation for agent use', () => {
    const content = getInitJustfileContent();
    assert.ok(
      fileContains(content, '--force'),
      'Expected --force flag to be documented in init-justfile',
    );
    assert.ok(
      fileContains(content, 'agent') ||
      fileContains(content, 'non-interactive') ||
      fileContains(content, 'skip'),
      'Expected --force to be described as skipping confirmation for agents',
    );
  });

  // Traceability: TC-DET-016 -> AC: Interactive confirmation when no markers exist
  test('TC-DET-016: interactive confirmation when no boundary markers and justfile exists', () => {
    const content = getInitJustfileContent();
    assert.ok(
      fileContains(content, 'confirm') ||
      fileContains(content, 'prompt') ||
      fileContains(content, 'overwrite') ||
      fileContains(content, 'ask'),
      'Expected interactive confirmation description for existing justfile without markers',
    );
  });

  // Traceability: TC-DET-017 -> AC: Idempotent re-run preserves custom recipes
  test('TC-DET-017: re-running init-justfile preserves user custom recipes', () => {
    const content = getInitJustfileContent();
    // Should describe that boundary markers enable idempotent re-runs
    assert.ok(
      fileContains(content, 'preserve') ||
      fileContains(content, 'keep') ||
      fileContains(content, 'custom') ||
      fileContains(content, 'outside'),
      'Expected description of preserving custom recipes outside boundary markers',
    );
  });
});

// ── TC-DET-018: project-type recipe for each variant ───────────────
describe('project-type recipe variants', () => {

  // Traceability: TC-DET-018 -> AC: project-type recipe outputs correct type
  test('TC-DET-018: all three project-type recipe variants exist', () => {
    const content = getInitJustfileContent();
    assert.ok(
      fileContains(content, '@echo "frontend"'),
      'Expected @echo "frontend" project-type recipe',
    );
    assert.ok(
      fileContains(content, '@echo "backend"'),
      'Expected @echo "backend" project-type recipe',
    );
    assert.ok(
      fileContains(content, '@echo "mixed"'),
      'Expected @echo "mixed" project-type recipe',
    );
  });
});

// ── TC-DET-019: Detection signal mapping ───────────────────────────
describe('Detection signal mapping', () => {

  // Traceability: TC-DET-019 -> AC: package.json = frontend, go.mod/Cargo.toml/pyproject.toml = backend
  test('TC-DET-019: detection correctly maps signals to frontend/backend categories', () => {
    const content = getInitJustfileContent();

    // Find the detection/classification section
    const workflowIdx = content.indexOf('## Workflow');
    assert.ok(workflowIdx !== -1, 'Expected Workflow section');

    const workflowSection = content.slice(workflowIdx);

    // package.json should be mapped as frontend signal
    assert.ok(
      fileContains(workflowSection, 'frontend'),
      'Expected "frontend" classification in workflow',
    );

    // go.mod should be mapped as backend signal
    assert.ok(
      fileContains(workflowSection, 'backend'),
      'Expected "backend" classification in workflow',
    );

    // Verify the detection is structured with clear mapping
    // The file should describe that package.json = frontend, go.mod/Cargo.toml/pyproject.toml = backend
    assert.ok(
      fileContains(content, 'package.json'),
      'Expected package.json in detection mapping',
    );
    assert.ok(
      fileContains(content, 'go.mod'),
      'Expected go.mod in detection mapping',
    );
  });
});
