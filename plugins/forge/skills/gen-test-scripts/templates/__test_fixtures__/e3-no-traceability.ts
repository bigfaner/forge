// E3 violation: test() without Traceability comment
import { test, expect } from '@playwright/test';

test('TC-020: Missing traceability', async ({ page }) => {
  await page.goto('http://localhost:3456/');
  await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible();
});
