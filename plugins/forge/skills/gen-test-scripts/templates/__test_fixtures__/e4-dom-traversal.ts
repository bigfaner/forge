// E4 violation: DOM parent traversal
import { test, expect } from '@playwright/test';

// Traceability: TC-030 → Spec Section 1
test('TC-030: DOM parent traversal', async ({ page }) => {
  await page.goto('http://localhost:3456/');
  const parent = page.locator('..');
  await expect(parent).toBeVisible();
});

// Traceability: TC-031 → Spec Section 2
test('TC-031: Another DOM traversal', async ({ page }) => {
  await page.goto('http://localhost:3456/');
  const el = page.locator('div/..');
  await expect(el).toBeVisible();
});
