import { describe, test, before, after } from 'node:test';
import assert from 'node:assert/strict';
import { ab, abJson, snapshotContains, findElement, screenshot, baseUrl } from './helpers.js';

describe('UI E2E Tests', () => {
  before(() => {
    ab(`open ${baseUrl}`);
    ab('wait --load networkidle');
  });

  after(() => {
    ab('close');
  });

  // Traceability: TC-001 → Story 1 / AC-1
  test('TC-001: Page renders with expected heading', () => {
    ab(`open ${baseUrl}/index.html`);
    ab('wait --load networkidle');
    assert.ok(snapshotContains('Dashboard'), 'Expected heading found');
    screenshot('TC-001');
  });
});
