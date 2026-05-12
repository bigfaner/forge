import { test, expect } from '@playwright/test';
import { runCli } from '../../helpers.js';

test.describe('CLI E2E Tests', () => {
  // Traceability: TC-020 → Spec Section X.Y
  test('TC-020: CLI command runs without error', () => {
    const result = runCli('echo hello'); // VERIFY: actual CLI command from CLI entry point
    expect(result.exitCode).toBe(0); // VERIFY: expected exit code
    expect(result.stdout).toMatch(/hello/); // VERIFY: expected stdout pattern
  });

  // Traceability: TC-021 -> Spec Section X.Y
  test('TC-021: CLI command with flags produces expected output', () => {
    const result = runCli('echo hello --flag value'); // VERIFY: actual CLI command with flags
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/hello/); // VERIFY: expected output pattern
  });

  // Traceability: TC-022 -> Spec Section X.Y
  test('TC-022: CLI command with invalid input shows error to stderr', () => {
    const result = runCli('echo --invalid-flag'); // VERIFY: command with invalid input
    expect(result.exitCode).not.toBe(0);
    expect(result.stderr).toMatch(/error|invalid|usage/i); // VERIFY: error pattern
  });
});
