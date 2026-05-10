// W4 violation: CSS class selectors
import { test, expect } from '@playwright/test';

// Traceability: TC-070 → Spec Section 1
test('TC-070: CSS class selector', async ({ page }) => {
  await page.goto('/');
  const btn = page.locator('.ant-btn-primary');
  await expect(btn).toBeVisible();
});

// Traceability: TC-071 → Spec Section 2
test('TC-071: Another CSS selector', async ({ page }) => {
  await page.goto('/');
  const item = page.locator('.list-item');
  await expect(item).toBeVisible();
});
