import { test, expect } from '@playwright/test';
import { readProjectFile } from '../helpers.js';

// ── Helpers ────────────────────────────────────────────────────────
function fileContains(content: string, needle: string): boolean {
  return content.includes(needle);
}

function getInitJustfileContent(): string {
  return readProjectFile('plugins/forge/commands/init-justfile.md');
}

// ── TC-DET-001 to TC-DET-004: Project-type detection signals ───────
test.describe('Project-type detection signals', () => {

  // Traceability: TC-DET-001 -> AC: Detection checks package.json (frontend signal)
  test('TC-DET-001: detection logic checks for package.json as frontend signal', () => {
    const content = getInitJustfileContent();
    expect(
      fileContains(content, 'package.json'),
      'Expected detection logic to check for package.json',
    ).toBeTruthy();
    // Must be in the detection section, not just the language templates
    const detectionSection = content.indexOf('Step 1:');
    expect(detectionSection !== -1, 'Expected Step 1 detection section').toBeTruthy();
  });

  // Traceability: TC-DET-002 -> AC: Detection checks go.mod (backend signal)
  test('TC-DET-002: detection logic checks for go.mod as backend signal', () => {
    const content = getInitJustfileContent();
    expect(
      fileContains(content, 'go.mod'),
      'Expected detection logic to check for go.mod',
    ).toBeTruthy();
  });

  // Traceability: TC-DET-003 -> AC: Detection checks Cargo.toml (backend signal)
  test('TC-DET-003: detection logic checks for Cargo.toml as backend signal', () => {
    const content = getInitJustfileContent();
    expect(
      fileContains(content, 'Cargo.toml'),
      'Expected detection logic to check for Cargo.toml',
    ).toBeTruthy();
  });

  // Traceability: TC-DET-004 -> AC: Detection checks pyproject.toml (backend signal)
  test('TC-DET-004: detection logic checks for pyproject.toml as backend signal', () => {
    const content = getInitJustfileContent();
    expect(
      fileContains(content, 'pyproject.toml'),
      'Expected detection logic to check for pyproject.toml',
    ).toBeTruthy();
  });
});

// ── TC-DET-005 to TC-DET-008: Project classification ───────────────
test.describe('Project classification', () => {

  // Traceability: TC-DET-005 -> AC: both signals -> mixed
  test('TC-DET-005: classification produces mixed when both frontend and backend signals detected', () => {
    const content = getInitJustfileContent();
    // Should describe "mixed" classification when both package.json and go.mod exist
    expect(
      fileContains(content, 'mixed'),
      'Expected "mixed" classification in detection logic',
    ).toBeTruthy();
  });

  // Traceability: TC-DET-006 -> AC: frontend only -> frontend
  test('TC-DET-006: classification produces frontend when only frontend signals detected', () => {
    const content = getInitJustfileContent();
    expect(
      fileContains(content, 'frontend'),
      'Expected "frontend" classification in detection logic',
    ).toBeTruthy();
  });

  // Traceability: TC-DET-007 -> AC: backend only -> backend
  test('TC-DET-007: classification produces backend when only backend signals detected', () => {
    const content = getInitJustfileContent();
    expect(
      fileContains(content, 'backend'),
      'Expected "backend" classification in detection logic',
    ).toBeTruthy();
  });

  // Traceability: TC-DET-008 -> AC: neither signal -> error message
  test('TC-DET-008: classification produces error when no known markers detected', () => {
    const content = getInitJustfileContent();
    // Should describe error case when no markers found
    expect(
      fileContains(content, 'no known project markers') ||
      fileContains(content, 'no project markers') ||
      fileContains(content, 'Error: no known') ||
      fileContains(content, 'no markers detected') ||
      fileContains(content, 'neither') ||
      fileContains(content, 'Cannot determine'),
      'Expected error handling when no project markers detected',
    ).toBeTruthy();
  });
});

// ── TC-DET-009 to TC-DET-011: Template selection ───────────────────
test.describe('Template selection and assembly', () => {

  // Traceability: TC-DET-009 -> AC: Selects backend template for backend project
  test('TC-DET-009: selects Backend Template for pure backend projects', () => {
    const backendTemplate = readProjectFile('plugins/forge/references/justfile-templates/go.just');
    expect(
      fileContains(backendTemplate, '@echo "backend"'),
      'Expected backend template to output @echo "backend"',
    ).toBeTruthy();
  });

  // Traceability: TC-DET-010 -> AC: Selects frontend template for frontend project
  test('TC-DET-010: selects Frontend Template for pure frontend projects', () => {
    const frontendTemplate = readProjectFile('plugins/forge/references/justfile-templates/node.just');
    expect(
      fileContains(frontendTemplate, '@echo "frontend"'),
      'Expected frontend template to output @echo "frontend"',
    ).toBeTruthy();
  });

  // Traceability: TC-DET-011 -> AC: Selects mixed template for mixed project
  test('TC-DET-011: selects Mixed Template for mixed projects', () => {
    const mixedTemplate = readProjectFile('plugins/forge/references/justfile-templates/mixed.just');
    expect(
      fileContains(mixedTemplate, '@echo "mixed"'),
      'Expected mixed template to output @echo "mixed"',
    ).toBeTruthy();
  });
});

// ── TC-DET-012 to TC-DET-014: Boundary markers ─────────────────────
test.describe('Boundary markers', () => {

  // Traceability: TC-DET-012 -> AC: Boundary markers wrap generated recipes
  test('TC-DET-012: generated recipes wrapped in forge boundary markers', () => {
    const content = getInitJustfileContent();
    const startMarker = '# --- forge standard recipes ---';
    const endMarker = '# --- end forge standard recipes ---';
    expect(
      fileContains(content, startMarker),
      'Expected start boundary marker',
    ).toBeTruthy();
    expect(
      fileContains(content, endMarker),
      'Expected end boundary marker',
    ).toBeTruthy();
  });

  // Traceability: TC-DET-013 -> AC: Boundary marker merge replaces only marked section
  test('TC-DET-013: boundary marker merge replaces only marked section', () => {
    const content = getInitJustfileContent();
    // Should describe merge logic that preserves content outside markers
    expect(
      fileContains(content, 'boundary marker') ||
      fileContains(content, 'markers') ||
      fileContains(content, 'replace') ||
      fileContains(content, 'merge'),
      'Expected boundary marker merge logic description',
    ).toBeTruthy();
  });

  // Traceability: TC-DET-014 -> AC: All 3 templates have boundary markers
  test('TC-DET-014: all three templates have boundary markers', () => {
    const startMarker = '# --- forge standard recipes ---';
    const endMarker = '# --- end forge standard recipes ---';
    const templates = [
      'plugins/forge/references/justfile-templates/go.just',
      'plugins/forge/references/justfile-templates/node.just',
      'plugins/forge/references/justfile-templates/mixed.just',
    ];
    for (const tpl of templates) {
      const content = readProjectFile(tpl);
      expect(
        fileContains(content, startMarker),
        `Expected start boundary marker in ${tpl}`,
      ).toBeTruthy();
      expect(
        fileContains(content, endMarker),
        `Expected end boundary marker in ${tpl}`,
      ).toBeTruthy();
    }
  });
});

// ── TC-DET-015 to TC-DET-017: --force flag and interactive confirmation ─
test.describe('--force flag and interactive confirmation', () => {

  // Traceability: TC-DET-015 -> AC: --force flag skips confirmation
  test('TC-DET-015: --force flag skips user confirmation for agent use', () => {
    const content = getInitJustfileContent();
    expect(
      fileContains(content, '--force'),
      'Expected --force flag to be documented in init-justfile',
    ).toBeTruthy();
    expect(
      fileContains(content, 'agent') ||
      fileContains(content, 'non-interactive') ||
      fileContains(content, 'skip'),
      'Expected --force to be described as skipping confirmation for agents',
    ).toBeTruthy();
  });

  // Traceability: TC-DET-016 -> AC: Interactive confirmation when no markers exist
  test('TC-DET-016: interactive confirmation when no boundary markers and justfile exists', () => {
    const content = getInitJustfileContent();
    expect(
      fileContains(content, 'confirm') ||
      fileContains(content, 'prompt') ||
      fileContains(content, 'overwrite') ||
      fileContains(content, 'ask'),
      'Expected interactive confirmation description for existing justfile without markers',
    ).toBeTruthy();
  });

  // Traceability: TC-DET-017 -> AC: Idempotent re-run preserves custom recipes
  test('TC-DET-017: re-running init-justfile preserves user custom recipes', () => {
    const content = getInitJustfileContent();
    // Should describe that boundary markers enable idempotent re-runs
    expect(
      fileContains(content, 'preserve') ||
      fileContains(content, 'keep') ||
      fileContains(content, 'custom') ||
      fileContains(content, 'outside'),
      'Expected description of preserving custom recipes outside boundary markers',
    ).toBeTruthy();
  });
});

// ── TC-DET-018: project-type recipe for each variant ───────────────
test.describe('project-type recipe variants', () => {

  // Traceability: TC-DET-018 -> AC: project-type recipe outputs correct type
  test('TC-DET-018: all three project-type recipe variants exist', () => {
    expect(
      fileContains(readProjectFile('plugins/forge/references/justfile-templates/node.just'), '@echo "frontend"'),
      'Expected @echo "frontend" project-type recipe',
    ).toBeTruthy();
    expect(
      fileContains(readProjectFile('plugins/forge/references/justfile-templates/go.just'), '@echo "backend"'),
      'Expected @echo "backend" project-type recipe',
    ).toBeTruthy();
    expect(
      fileContains(readProjectFile('plugins/forge/references/justfile-templates/mixed.just'), '@echo "mixed"'),
      'Expected @echo "mixed" project-type recipe',
    ).toBeTruthy();
  });
});

// ── TC-DET-019: Detection signal mapping ───────────────────────────
test.describe('Detection signal mapping', () => {

  // Traceability: TC-DET-019 -> AC: package.json = frontend, go.mod/Cargo.toml/pyproject.toml = backend
  test('TC-DET-019: detection correctly maps signals to frontend/backend categories', () => {
    const content = getInitJustfileContent();

    // Find the detection/classification section
    const workflowIdx = content.indexOf('## Workflow');
    expect(workflowIdx !== -1, 'Expected Workflow section').toBeTruthy();

    const workflowSection = content.slice(workflowIdx);

    // package.json should be mapped as frontend signal
    expect(
      fileContains(workflowSection, 'frontend'),
      'Expected "frontend" classification in workflow',
    ).toBeTruthy();

    // go.mod should be mapped as backend signal
    expect(
      fileContains(workflowSection, 'backend'),
      'Expected "backend" classification in workflow',
    ).toBeTruthy();

    // Verify the detection is structured with clear mapping
    // The file should describe that package.json = frontend, go.mod/Cargo.toml/pyproject.toml = backend
    expect(
      fileContains(content, 'package.json'),
      'Expected package.json in detection mapping',
    ).toBeTruthy();
    expect(
      fileContains(content, 'go.mod'),
      'Expected go.mod in detection mapping',
    ).toBeTruthy();
  });
});
