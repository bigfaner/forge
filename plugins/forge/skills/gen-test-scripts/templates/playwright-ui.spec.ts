import { test, expect } from '@playwright/test';
import { screenshot, baseUrl, waitForApiAction
  // CONDITIONAL: Uncomment imports below only if auth-required-test exists
  // , loginViaUI, ensureAuthState, getApiToken, createAuthCurl, CurlResponse
  // CONDITIONAL: Uncomment import below only if login-test exists (for login form filling)
  // , defaultCreds, clearAuthState, clearCachedToken
  // CONDITIONAL: Uncomment imports below only if serial suite with shared data exists
  // , createTestResource, cleanupTestResources, withRetry
} from '../../helpers.js';

// ── FORBIDDEN PATTERNS (do NOT use in generated code) ───────────
// ❌ page.waitForTimeout(N)          — use waitForApiAction / expect().toBeVisible
// ❌ await new Promise(r => setTimeout(r, N)) — use withRetry for polling
// ❌ .locator('..')                  — DOM parent traversal breaks on any UI change
// ❌ .ant-*, .btn-* CSS classes      — use role/name/testid selectors only
// ─────────────────────────────────────────────────────────────────

test.describe('UI E2E Tests', () => {
  // CONDITIONAL: Uncomment the block below only if auth-required-test exists
  // and the playwright.config.ts setup project is NOT configured.
  // If the setup project IS configured, this block is unnecessary —
  // storageState is injected automatically by Playwright.
  //
  // test.beforeAll(async ({ browser }) => {
  //   const ctx = await browser.newContext();
  //   const page = await ctx.newPage();
  //   await ensureAuthState(page);
  //   await ctx.close();
  // });

  // ── Login Tests (no shared auth) ────────────────────────────────
  // CONDITIONAL: Uncomment this block only if login-test exists
  //
  // test.describe('Login', () => {
  //   test.afterAll(() => {
  //     clearCachedToken();
  //     clearAuthState();
  //   });
  //
  //   test('TC-001: Login with valid credentials', async ({ page }) => {
  //     await page.goto(`${baseUrl()}/login`); // VERIFY: login route from router files
  //     await page.getByRole('textbox', { name: 'Username' }).fill(defaultCreds.username); // VERIFY: username field locator
  //     await page.getByRole('textbox', { name: 'Password' }).fill(defaultCreds.password); // VERIFY: password field locator
  //     await page.getByRole('button', { name: 'Login' }).click(); // VERIFY: submit button locator
  //     await page.waitForURL('**/dashboard'); // VERIFY: post-login redirect route
  //     await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible(); // VERIFY: dashboard heading
  //     await screenshot(page, 'TC-001');
  //   });
  // });

  // ── Pattern References (replace with actual test cases from test-cases.md) ──

  // PATTERN REFERENCE: page render with heading assertion
  // test('TC-002: Page renders with expected heading', async ({ page }) => {
  //   await page.goto(`${baseUrl()}/`); // VERIFY: target route
  //   await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible(); // VERIFY: heading text
  //   await screenshot(page, 'TC-002');
  // });

  // PATTERN REFERENCE: form submission — wait for API response instead of timeout
  // test('TC-003: Form submission succeeds with valid data', async ({ page }) => {
  //   await page.goto(`${baseUrl()}/form`); // VERIFY: form route
  //   await page.getByRole('textbox', { name: 'Name' }).fill('Test User'); // VERIFY: field locator
  //   const res = await waitForApiAction(page,
  //     () => page.getByRole('button', { name: 'Submit' }).click(),
  //     '/api/resources', // VERIFY: API endpoint triggered by submit
  //   );
  //   expect(res.status).toBe(200); // VERIFY: expected status code
  //   await expect(page.getByText('Success')).toBeVisible(); // VERIFY: success message
  //   await screenshot(page, 'TC-003');
  // });

  // PATTERN REFERENCE: error message assertion for invalid input
  // test('TC-004: Error message shown for invalid input', async ({ page }) => {
  //   await page.goto(`${baseUrl()}/form`); // VERIFY: form route
  //   await page.getByRole('button', { name: 'Submit' }).click();
  //   await expect(page.getByText(/error|required|invalid/i)).toBeVisible(); // VERIFY: error message
  //   await screenshot(page, 'TC-004');
  // });

  // PATTERN REFERENCE: navigation between pages
  // test('TC-005: Navigation between pages works', async ({ page }) => {
  //   await page.goto(`${baseUrl()}/`); // VERIFY: start route
  //   await page.getByRole('link', { name: 'About' }).click(); // VERIFY: nav link
  //   await expect(page).toHaveURL(/\/about/); // VERIFY: target route pattern
  //   await screenshot(page, 'TC-005');
  // });

  // PATTERN REFERENCE: serial suite with shared data + cleanup
  // Limit serial suites to 15 tests max. Each suite must have afterAll cleanup.
  //
  // test.describe.serial('Resource lifecycle', () => {
  //   let resourceId: string;
  //   const _cleanupIds: string[] = [];
  //
  //   test.afterAll(async () => {
  //     // Cleanup: delete all created resources
  //     const token = await getApiToken(apiBaseUrl(), '/v1/auth/login'); // VERIFY: auth endpoint
  //     const del = createAuthCurl(apiBaseUrl(), token);
  //     for (const id of _cleanupIds) {
  //       await del('DELETE', `/v1/resources/${id}`).catch(() => {}); // best-effort
  //     }
  //   });
  //
  //   test('TC-006: Create resource via UI', async ({ page }) => {
  //     await page.goto(`${baseUrl()}/resources/new`); // VERIFY: create page route
  //     await page.getByRole('textbox', { name: 'Title' }).fill('E2E Test'); // VERIFY: field locator
  //     const res = await waitForApiAction(page,
  //       () => page.getByRole('button', { name: 'Create' }).click(),
  //       '/api/resources', // VERIFY: API endpoint
  //     );
  //     expect(res.status).toBe(201); // VERIFY: expected status
  //     const data = JSON.parse(res.body);
  //     resourceId = data.id ?? data.data?.id; // VERIFY: response ID field
  //     if (!resourceId) throw new Error('resourceId is undefined after create');
  //     _cleanupIds.push(resourceId);
  //   });
  //
  //   test('TC-007: Created resource is visible in list', async ({ page }) => {
  //     await page.goto(`${baseUrl()}/resources`); // VERIFY: list route
  //     await expect(page.getByText('E2E Test')).toBeVisible();
  //     await screenshot(page, 'TC-007');
  //   });
  // });
});
