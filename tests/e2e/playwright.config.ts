import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: '.',
  testIgnore: process.env.E2E_FEATURE ? undefined : /features\//,
  timeout: 30_000,
  expect: { timeout: 10_000 },
  globalTimeout: 300_000,
  retries: 0,
  workers: 1,
  reporter: [
    ['list'],
    ['json', { outputFile: 'results/test-results.json' }],
  ],
  use: {
    headless: true,
    screenshot: 'only-on-failure',
  },
  outputDir: 'results/',
});
