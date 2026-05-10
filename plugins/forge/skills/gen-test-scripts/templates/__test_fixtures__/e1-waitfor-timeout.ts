// E1 violation: waitForTimeout and setTimeout
import { test, expect } from '@playwright/test';

// Traceability: TC-010 → Spec Section 1
test('TC-010: Bad test with waitForTimeout', async ({ page }) => {
  await page.goto('http://localhost:3456/');
  await page.waitForTimeout(5000);
  await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible();
});

// Traceability: TC-011 → Spec Section 2
test('TC-011: Bad test with setTimeout', async ({ page }) => {
  await page.goto('http://localhost:3456/');
  await new Promise(r => setTimeout(r, 3000));
  await expect(page.getByText('Done')).toBeVisible();
});
