// Auth Setup — run once before all authenticated tests.
// This file is used by the 'setup' project in playwright.config.ts.
// It logs in via UI and saves browser storage state for reuse by authenticated tests.
//
// This file is ONLY generated when auth-required-test cases exist.
// If all tests are public or login-only, this file is not needed.

import { test as setup, expect } from '@playwright/test';
import { baseUrl, ensureAuthState } from './helpers.js';

setup('authenticate', async ({ page }) => {
  await ensureAuthState(page);
  // Verify login succeeded — the saved storageState will be reused by authenticated tests
  await page.goto(`${baseUrl()}/`);
  await expect(page).toHaveURL(/(?!.*login)/);
});
