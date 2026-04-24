import { describe, test, before, after } from 'node:test';
import assert from 'node:assert/strict';
import type { Page } from 'playwright';
import { setupBrowser, teardownBrowser, screenshot, loginViaUI, baseUrl } from './helpers.js';

describe('UI E2E Tests', () => {
  let page: Page;

  before(async () => {
    page = await setupBrowser();
    // Shared auth: login once for all non-login tests
    await loginViaUI(page);
  });

  after(async () => {
    await teardownBrowser();
  });

  // ── Login Tests (no shared auth) ────────────────────────────────
  describe('Login', () => {
    let loginPage: Page;

    before(async () => {
      // Fresh page without auth cookies
      loginPage = await page.context().newPage();
    });

    after(async () => {
      await loginPage.close();
    });

    // Traceability: TC-001 → Story 1 / AC-1
    test('TC-001: Login with valid credentials', async () => {
      await loginPage.goto(`${baseUrl}/login`);
      await loginPage.waitForLoadState('networkidle');
      await loginPage.getByRole('textbox', { name: 'Username' }).fill('admin');
      await loginPage.getByRole('textbox', { name: 'Password' }).fill('password');
      await loginPage.getByRole('button', { name: 'Login' }).click();
      await loginPage.waitForURL('**/dashboard');
      assert.ok(true, 'Redirected to dashboard after login');
      await screenshot(loginPage, 'TC-001');
    });
  });

  // ── Authenticated Tests (use shared auth) ───────────────────────
  // Traceability: TC-002 → Story 2 / AC-1
  test('TC-002: Page renders with expected heading', async () => {
    await page.goto(`${baseUrl}/index.html`);
    await page.waitForLoadState('networkidle');
    await page.getByRole('heading', { name: 'Dashboard' }).waitFor({ state: 'visible' });
    assert.ok(true, 'Heading visible');
    await screenshot(page, 'TC-002');
  });
});
