import { test, expect } from '@playwright/test';
import { screenshot, baseUrl
  // CONDITIONAL: Uncomment imports below only if auth-required-test exists
  // , loginViaUI
  // , defaultCreds
} from '../helpers.js';

test.describe('UI E2E Tests', () => {
  // CONDITIONAL: Uncomment the block below only if auth-required-test exists
  // test.beforeEach(async ({ page }) => {
  //   await loginViaUI(page);
  // });

  // ── Login Tests (no shared auth) ────────────────────────────────
  // CONDITIONAL: Uncomment this block only if login-test exists
  //
  // test.describe('Login', () => {
  //   test('TC-001: Login with valid credentials', async ({ page }) => {
  //     await page.goto(`${baseUrl}/login`); // VERIFY: login route from router files
  //     await page.getByRole('textbox', { name: 'Username' }).fill(defaultCreds.username); // VERIFY: username field locator from sitemap/code
  //     await page.getByRole('textbox', { name: 'Password' }).fill(defaultCreds.password); // VERIFY: password field locator from sitemap/code
  //     await page.getByRole('button', { name: 'Login' }).click(); // VERIFY: submit button locator from sitemap/code
  //     // VERIFY: post-login redirect route from router files
  //     await page.waitForURL('**/dashboard');
  //     await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible(); // VERIFY: dashboard heading after login
  //     await screenshot(page, 'TC-001');
  //   });
  // });

  // ── Authenticated Tests (use shared auth) ───────────────────────
  // Traceability: TC-002 → Story 2 / AC-1
  test('TC-002: Page renders with expected heading', async ({ page }) => {
    await page.goto(`${baseUrl}/`); // VERIFY: target route from test case
    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible(); // VERIFY: heading text from sitemap/code
    await screenshot(page, 'TC-002');
  });
});
