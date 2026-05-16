import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: '.',
  timeout: 30000,
  expect: { timeout: 10000 },
  retries: 0,
  workers: 1,
  reporter: [['list']],
  use: {
    headless: true,
  },
  outputDir: 'results/',
});
