import { test, expect } from '@playwright/test';
import { runCli } from '../helpers.js';

test.describe('CLI E2E Tests', () => {
  // Traceability: TC-020 → Spec Section X.Y
  test('TC-020: CLI command runs without error', () => {
    const result = runCli('echo hello'); // VERIFY: actual CLI command from CLI entry point
    expect(result.exitCode).toBe(0); // VERIFY: expected exit code
    expect(result.stdout).toMatch(/hello/); // VERIFY: expected stdout pattern
  });
});
