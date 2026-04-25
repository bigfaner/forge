import { describe, test, before, after } from 'node:test';
import assert from 'node:assert/strict';
import type { Page } from 'playwright';
import { setupBrowser, teardownBrowser, screenshot, baseUrl
  // CONDITIONAL: Uncomment imports below only if auth-required-test exists
  // , loginViaUI
  // , defaultCreds
} from './helpers.js';

describe('UI E2E Tests', () => {
  let page: Page;

  before(async () => {
    page = await setupBrowser();
    // CONDITIONAL: Uncomment the line below only if auth-required-test exists
    // await loginViaUI(page);
  });

  after(async () => {
    await teardownBrowser();
  });

  // ── Login Tests (no shared auth) ────────────────────────────────
  // CONDITIONAL: Uncomment this block only if login-test exists
  //
  // describe('Login', () => {
  //   let loginPage: Page;
  //
  //   before(async () => {
  //     loginPage = await page.context().newPage();
  //   });
  //
  //   after(async () => {
  //     await loginPage.close();
  //   });
  //
  //   // Traceability: TC-001 → Story 1 / AC-1
  //   test('TC-001: Login with valid credentials', async () => {
  //     await loginPage.goto(`${baseUrl}/login`);
  //     await loginPage.waitForLoadState('networkidle');
  //     await loginPage.getByRole('textbox', { name: 'Username' }).fill(defaultCreds.username);
  //     await loginPage.getByRole('textbox', { name: 'Password' }).fill(defaultCreds.password);
  //     await loginPage.getByRole('button', { name: 'Login' }).click();
  //     // TEMPLATE: Replace with actual post-login route
  //     await loginPage.waitForURL('**/dashboard');
  //     assert.match(loginPage.url(), /dashboard/, 'Redirected to dashboard after login');
  //     await screenshot(loginPage, 'TC-001');
  //   });
  // });

  // ── Authenticated Tests (use shared auth) ───────────────────────
  // Traceability: TC-002 → Story 2 / AC-1
  test('TC-002: Page renders with expected heading', async () => {
    await page.goto(`${baseUrl}/`);
    await page.waitForLoadState('networkidle');
    await page.getByRole('heading', { name: 'Dashboard' }).waitFor({ state: 'visible' });
    await screenshot(page, 'TC-002');
  });
});
