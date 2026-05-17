// Clean spec file - should produce NO errors and NO warnings
import { test, expect } from '@playwright/test';

// Traceability: TC-001 → Spec Section 1
test('TC-001: Page renders with heading', async ({ page }) => {
  await page.goto('http://localhost:3456/');
  await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible();
});

// Traceability: TC-002 → Spec Section 2
test('TC-002: Form submission succeeds', async ({ page }) => {
  await page.goto('http://localhost:3456/form');
  await page.getByRole('textbox', { name: 'Name' }).fill('Test');
  await page.getByRole('button', { name: 'Submit' }).click();
  await expect(page.getByText('Success')).toBeVisible();
});
