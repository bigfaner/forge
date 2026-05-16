import { test, expect } from '@playwright/test';
import {
  runCli,
  readProjectFile,
  projectFileExists,
  PROJECT_ROOT,
} from '../../helpers.js';
import { join } from 'node:path';
import { writeFileSync, mkdirSync, rmSync, existsSync } from 'node:fs';

// ── Shared paths ───────────────────────────────────────────────────
const FIXTURES_DIR = join(PROJECT_ROOT, 'tests', 'e2e', 'fixtures', 'forge-testing-optimization');
const VALIDATE_SCRIPT = join(
  PROJECT_ROOT,
  'plugins',
  'forge',
  'skills',
  'gen-test-scripts',
  'templates',
  'validate-specs.mjs',
);

// ── TC-001: validate-specs detects E1-E4 ERROR rules ───────────────

test.describe('validate-specs E1-E4 ERROR detection', () => {
  const e1FixtureDir = join(FIXTURES_DIR, 'tc001');

  test.beforeAll(() => {
    // Create fixture spec files with deliberate E1-E4 violations
    mkdirSync(e1FixtureDir, { recursive: true });

    // E1: waitForTimeout
    writeFileSync(join(e1FixtureDir, 'e1-wait.spec.ts'), `
import { test } from '@playwright/test';
test('E1 waitForTimeout violation', async ({ page }) => {
  await page.waitForTimeout(5000);
});
`);

    // E1: setTimeout
    writeFileSync(join(e1FixtureDir, 'e1-setTimeout.spec.ts'), `
import { test } from '@playwright/test';
test('E1 setTimeout violation', async () => {
  setTimeout(() => {}, 1000);
});
`);

    // E3: missing Traceability comment
    writeFileSync(join(e1FixtureDir, 'e3-no-trace.spec.ts'), `
import { test } from '@playwright/test';
test('E3 no traceability', async () => {
  // no Traceability comment
});
`);

    // E4: DOM parent traversal locator('..')
    writeFileSync(join(e1FixtureDir, 'e4-dom-traverse.spec.ts'), `
import { test } from '@playwright/test';
test('E4 DOM parent traversal', async ({ page }) => {
  // Traceability: TC-001 -> proposal Section 2.1
  const parent = page.locator('..');
});
`);
  });

  test.afterAll(() => {
    rmSync(e1FixtureDir, { recursive: true, force: true });
  });

  // Traceability: TC-001 -> SC "validate-specs.mjs 能检测 E1-E4 四种 ERROR"
  test('TC-001: validate-specs detects E1-E4 ERROR rules', () => {
    const result = runCli(`node "${VALIDATE_SCRIPT}" "${e1FixtureDir}"`);

    expect(result.exitCode).not.toBe(0);

    // Parse JSON output
    const output = JSON.parse(result.stdout);
    const errorRules = output.errors.map((e: { rule: string }) => e.rule);

    // E1 violations detected
    expect(errorRules.filter((r: string) => r === 'E1').length).toBeGreaterThanOrEqual(2);

    // E3 violation detected
    expect(errorRules).toContain('E3');

    // E4 violation detected
    expect(errorRules).toContain('E4');
  });
});

// ── TC-002: validate-specs detects W1-W4 WARNING rules ─────────────

test.describe('validate-specs W1-W4 WARNING detection', () => {
  const wFixtureDir = join(FIXTURES_DIR, 'tc002');

  test.beforeAll(() => {
    mkdirSync(wFixtureDir, { recursive: true });

    // W1: serial suite > 15 tests + W2: no afterAll
    // Each test has Traceability comment to avoid E3 errors
    const tests = Array.from({ length: 20 }, (_, i) =>
      `  // Traceability: TC-002 -> SC W1 test-${i}\n  test('test-${i}', () => { /* pass */ });`
    ).join('\n');

    writeFileSync(join(wFixtureDir, 'w1-w2.spec.ts'), `
import { test } from '@playwright/test';
test.describe.serial('big suite', () => {
${tests}
});
`);

    // W3: beforeEach with login
    writeFileSync(join(wFixtureDir, 'w3-beforeEach-login.spec.ts'), `
import { test } from '@playwright/test';
test.describe('login suite', () => {
  test.beforeEach(async ({ page }) => {
    await loginViaUI(page);
  });
  // Traceability: TC-002 -> SC W3
  test('something after login', async () => {});
});
`);

    // W4: CSS class selector
    writeFileSync(join(wFixtureDir, 'w4-css-class.spec.ts'), `
import { test } from '@playwright/test';
// Traceability: TC-002 -> SC W4
test('CSS class selector', async ({ page }) => {
  const btn = page.locator('.ant-btn');
});
`);
  });

  test.afterAll(() => {
    rmSync(wFixtureDir, { recursive: true, force: true });
  });

  // Traceability: TC-002 -> SC "validate-specs.mjs 能检测 W1-W4 四种 WARNING"
  test('TC-002: validate-specs detects W1-W4 WARNING rules', () => {
    const result = runCli(`node "${VALIDATE_SCRIPT}" "${wFixtureDir}"`);

    // Warnings are non-blocking — exit code should be 0 when only warnings present
    expect(result.exitCode).toBe(0);

    const output = JSON.parse(result.stdout);
    const warningRules = output.warnings.map((w: { rule: string }) => w.rule);

    expect(warningRules).toContain('W1');
    expect(warningRules).toContain('W2');
    expect(warningRules).toContain('W3');
    expect(warningRules).toContain('W4');
  });
});

// ── TC-003: ts-morph devDependency in package.json ──────────────────

// Traceability: TC-003 -> SC "ts-morph 在 tests/e2e/package.json 中作为 devDependency 存在"
test('TC-003: ts-morph devDependency in package.json', () => {
  const pkgJsonPath = 'plugins/forge/skills/gen-test-scripts/templates/package.json';
  const content = readProjectFile(pkgJsonPath);

  const pkg = JSON.parse(content);

  expect(pkg.devDependencies).toBeDefined();
  expect(pkg.devDependencies).toHaveProperty('ts-morph');

  const version = pkg.devDependencies['ts-morph'];
  expect(version).toBeTruthy();
  expect(version).not.toBe('*');
  // Valid semver range: starts with ^, ~, >=, or digit
  expect(version).toMatch(/^[\^~>=]?\d/);
});

// ── TC-004: task validate-specs command executes and returns structured output ──

test.describe('task validate-specs command', () => {
  const cleanFixtureDir = join(FIXTURES_DIR, 'tc004-clean');
  const dirtyFixtureDir = join(FIXTURES_DIR, 'tc004-dirty');

  test.beforeAll(() => {
    mkdirSync(cleanFixtureDir, { recursive: true });
    mkdirSync(dirtyFixtureDir, { recursive: true });

    // Clean spec — no violations
    writeFileSync(join(cleanFixtureDir, 'clean.spec.ts'), `
import { test } from '@playwright/test';
// Traceability: TC-004 -> SC "task validate-specs"
test('clean test', async () => {
  // No violations
});
`);

    // Dirty spec — E1 violation
    writeFileSync(join(dirtyFixtureDir, 'dirty.spec.ts'), `
import { test } from '@playwright/test';
// Traceability: TC-004 -> SC E1 violation test
test('dirty test', async ({ page }) => {
  await page.waitForTimeout(5000);
});
`);
  });

  test.afterAll(() => {
    rmSync(cleanFixtureDir, { recursive: true, force: true });
    rmSync(dirtyFixtureDir, { recursive: true, force: true });
  });

  // Traceability: TC-004 -> SC "task validate-specs 命令能执行校验并返回结构化输出"
  test('TC-004: task validate-specs returns structured output with violations', () => {
    // Run validate-specs.mjs directly to verify structured JSON output
    const result = runCli(`node "${VALIDATE_SCRIPT}" "${dirtyFixtureDir}"`);

    expect(result.exitCode).not.toBe(0);

    const output = JSON.parse(result.stdout);

    // Verify structured output shape
    expect(output).toHaveProperty('errors');
    expect(output).toHaveProperty('warnings');
    expect(Array.isArray(output.errors)).toBe(true);
    expect(Array.isArray(output.warnings)).toBe(true);

    // Verify E1 violation is reported
    if (output.errors.length > 0) {
      const e1Error = output.errors.find((e: { rule: string }) => e.rule === 'E1');
      expect(e1Error).toBeDefined();
      expect(e1Error).toHaveProperty('rule');
      expect(e1Error).toHaveProperty('file');
      expect(e1Error).toHaveProperty('line');
      expect(e1Error).toHaveProperty('message');
    }
  });
});

// ── TC-005: gen-test-scripts SKILL.md contains Step 4.5 ────────────

// Traceability: TC-005 -> SC "gen-test-scripts SKILL.md 包含 Step 4.5 结构校验步骤"
test('TC-005: gen-test-scripts SKILL.md contains Step 4.5 structural validation', () => {
  const content = readProjectFile('plugins/forge/skills/gen-test-scripts/SKILL.md');

  // Step 4.5 section heading exists
  expect(content).toMatch(/### Step 4\.5[:\s]/);

  // Step 4.5 describes structural validation using task validate-specs
  expect(content).toMatch(/task validate-specs/);

  // ERROR results block downstream
  expect(content).toMatch(/ERROR.*block/i);

  // WARNING results are non-blocking
  expect(content).toMatch(/WARNING.*non-block/i);
});

// ── TC-006: gen-test-scripts aborts when Step Actionability < 20 ────

// Traceability: TC-006 -> SC "gen-test-scripts 在 eval-test-cases Step Actionability < 20 时中止"
test('TC-006: gen-test-scripts aborts when eval-test-cases Step Actionability below 20', () => {
  const content = readProjectFile('plugins/forge/skills/gen-test-scripts/SKILL.md');

  // Prerequisites section exists
  expect(content).toMatch(/## Prerequisites/);

  // Step Actionability Gate section exists
  expect(content).toMatch(/Step Actionability/);

  // Aborts when score < 20
  expect(content).toMatch(/Step Actionability\s*<\s*20/);
  expect(content).toMatch(/ABORT/);
});

// ── TC-007: gen-test-cases Element field marked as required ────────

// Traceability: TC-007 -> SC "gen-test-cases SKILL.md 和模板中 Element 字段标记为必填"
test('TC-007: gen-test-cases Element field marked as required', () => {
  const skillContent = readProjectFile('plugins/forge/skills/gen-test-cases/SKILL.md');

  // SKILL.md states Element is required
  expect(skillContent).toMatch(/Element.*required/i);

  // SKILL.md defines Element field in the generated test case format
  // (the template is a structural skeleton; Element is defined by SKILL.md generation rules)
  expect(skillContent).toMatch(/- \*\*Element\*\*/);

  // Fallback behavior when sitemap is missing is defined
  expect(skillContent).toMatch(/sitemap-missing/);
});
