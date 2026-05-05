import { defineConfig } from '@playwright/test';

// E2E_FEATURE=1 disables testIgnore so `just test-e2e --feature <slug>` can run
// staging area tests. Without it, testIgnore excludes features/ from the regression suite.
const featureMode = !!process.env.E2E_FEATURE;

export default defineConfig({
  testDir: '.',
  testIgnore: featureMode ? [] : /^features\//,
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
  // CONDITIONAL: Uncomment the projects section below when auth-required UI tests exist.
  // This enables Playwright's storageState mechanism: login once in setup, reuse across all tests.
  // When projects are defined, testIgnore at top level is inherited by each project.
  // The setup project runs auth-setup.ts to save browser state via ensureAuthState().
  // The authenticated project depends on setup and reuses the saved storageState.
  //
  // projects: [
  //   {
  //     name: 'setup',
  //     testDir: '.',
  //     testMatch: /auth-setup\.ts/,
  //     testIgnore: [],
  //   },
  //   {
  //     name: 'authenticated',
  //     testDir: 'features',
  //     testMatch: /.*\.spec\.ts/,
  //     testIgnore: [],
  //     dependencies: ['setup'],
  //     use: {
  //       storageState: 'results/.auth/state.json',
  //     },
  //   },
  // ],
});
