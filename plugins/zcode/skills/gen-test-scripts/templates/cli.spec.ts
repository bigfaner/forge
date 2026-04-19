import { describe, test } from 'node:test';
import assert from 'node:assert/strict';
import { runCli } from './helpers.js';

describe('CLI E2E Tests', () => {
  // Traceability: TC-020 → Spec Section X.Y
  test('TC-020: CLI command runs without error', () => {
    const result = runCli('echo hello');
    assert.equal(result.exitCode, 0);
    assert.match(result.stdout, /hello/);
  });
});
