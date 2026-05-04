import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: '.',
  testIgnore: /^features\//,
  timeout: 30_000,
  expect: { timeout: 10_000 },
  globalTimeout: 300_000,
  // Increase retries for CI: set E2E_RETRIES env var (e.g. E2E_RETRIES=2)
  retries: Number(process.env.E2E_RETRIES ?? '0'),
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
