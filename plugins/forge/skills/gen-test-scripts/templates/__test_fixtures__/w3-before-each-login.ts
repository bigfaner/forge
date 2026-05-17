// W3 violation: beforeEach with login
import { test, expect } from '@playwright/test';
import { loginViaUI } from '../../helpers.js';

test.describe('Tests with beforeEach login', () => {
  test.beforeEach(async ({ page }) => {
    await loginViaUI(page);
  });

  // Traceability: TC-060 → Spec Section 1
  test('TC-060: Some test', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
});
