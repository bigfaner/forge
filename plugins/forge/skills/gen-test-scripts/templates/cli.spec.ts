import { describe, test } from 'node:test';
import assert from 'node:assert/strict';
import { runCli } from '../helpers.js';

describe('CLI E2E Tests', () => {
  // Traceability: TC-020 → Spec Section X.Y
  test('TC-020: CLI command runs without error', () => {
    const result = runCli('echo hello'); // VERIFY: actual CLI command from CLI entry point
    assert.equal(result.exitCode, 0); // VERIFY: expected exit code
    assert.match(result.stdout, /hello/); // VERIFY: expected stdout pattern
  });
});
