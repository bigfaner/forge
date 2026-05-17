// W1 + W2 violations: large serial suite without afterAll
import { test, expect } from '@playwright/test';

test.describe.serial('Large suite without cleanup', () => {
  // Traceability: TC-040 → Spec Section 1
  test('TC-040: test 1', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-041 → Spec Section 2
  test('TC-041: test 2', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-042 → Spec Section 3
  test('TC-042: test 3', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-043 → Spec Section 4
  test('TC-043: test 4', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-044 → Spec Section 5
  test('TC-044: test 5', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-045 → Spec Section 6
  test('TC-045: test 6', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-046 → Spec Section 7
  test('TC-046: test 7', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-047 → Spec Section 8
  test('TC-047: test 8', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-048 → Spec Section 9
  test('TC-048: test 9', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-049 → Spec Section 10
  test('TC-049: test 10', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-050 → Spec Section 11
  test('TC-050: test 11', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-051 → Spec Section 12
  test('TC-051: test 12', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-052 → Spec Section 13
  test('TC-052: test 13', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-053 → Spec Section 14
  test('TC-053: test 14', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-054 → Spec Section 15
  test('TC-054: test 15', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
  // Traceability: TC-055 → Spec Section 16
  test('TC-055: test 16', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('ok')).toBeVisible();
  });
});
