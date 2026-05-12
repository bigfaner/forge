# Playwright Generate Strategy

Profile-specific test generation rules for the `gen-test-scripts` skill.

## Test Runner & Imports

| Test type | Runner | Assertion | HTTP | Process |
|-----------|--------|-----------|------|---------|
| UI | `@playwright/test` (full browser lifecycle) | `expect` from `@playwright/test` | — | — |
| API | `@playwright/test` (test runner only) | `expect` from `@playwright/test` | Node.js built-in `fetch` | — |
| CLI | `@playwright/test` (test runner only) | `expect` from `@playwright/test` | — | `child_process.execSync` via `runCli()` |

All tests import `test`, `expect`, `test.describe` from `@playwright/test`. The `node:test` and `node:assert` modules are **forbidden**. Using `agent-browser` in generated spec files is **forbidden**.

## Spec Template Mapping

| Test type | Template file | Output filename |
|-----------|--------------|-----------------|
| UI | `templates/playwright-ui.spec.ts` | `ui.spec.ts` |
| API | `templates/api.spec.ts` | `api.spec.ts` |
| CLI | `templates/cli.spec.ts` | `cli.spec.ts` |

Templates live at `plugins/forge/profiles/web-playwright/templates/`.

## Locator Mapping

Map sitemap element data to Playwright locators by priority:

| Priority | Condition | Generated code |
|----------|-----------|----------------|
| 0 | `role` + `name` | `page.getByRole('button', { name: 'Submit' })` |
| 0+ | heading + `level` | `page.getByRole('heading', { name: 'Dashboard', level: 1 })` |
| 1 | `label` | `page.getByLabel('Email address')` |
| 2 | `placeholder` (label empty) | `page.getByPlaceholder('Search...')` |
| 3 | Static text node | `page.getByText('No results found')` |
| 4 | `data-testid` visible | `page.getByTestId('user-avatar')` |
| 5 (fallback) | None of the above | `page.locator('.btn') // UNSTABLE: no semantic anchor` |

For dynamic states: click the trigger element's locator first, then map in-state elements.

## Auth Classification

| Category | Detection | Strategy |
|----------|-----------|----------|
| **login-test** | Target matches `ui/login`, `ui/auth`, `ui/signin`, `api/auth`, `api/login`, `api/token` | No shared auth. UI: independent `loginPage`. API: raw `curl()`. Must call `clearCachedToken()` / `clearAuthState()` after. |
| **auth-required-test** | Pre-conditions contain "authenticated", "logged in", or target implies protected resource | **Cached shared auth** — credentials acquired once, reused across all tests. |
| **public-test** | No auth pre-conditions, public target | No auth needed. |
| **custom-auth-test** | Pre-conditions mention "API key", "X-API-Key", "OAuth", "session cookie" | Custom auth setup in `test.beforeAll`. |

Classification priority (first match wins): custom-auth-test → login-test → auth-required-test → public-test.

Login tests and authenticated tests must NOT be mixed in the same `describe` block.

## Credential Caching

**API tests**: `getApiToken()` in `helpers.ts` caches token at module level. First call authenticates; subsequent calls return cached token. Transparent to spec files.

**UI tests**: Playwright `storageState` mechanism:
1. Generate `auth-setup.ts` (from `templates/auth-setup.ts`) into `tests/e2e/`
2. Uncomment the `projects` section in `playwright.config.ts`
3. Playwright injects saved cookies/localStorage into each test's browser context

**Fallback** (no projects configured): `ensureAuthState(page)` in `test.beforeAll` with fresh browser context.

**Login tests**: After login/logout tests, call both `clearCachedToken()` and `clearAuthState()` to invalidate stale credentials.

## Import Conventions

| Context | Import path |
|---------|-------------|
| Staging (`tests/e2e/features/<slug>/`) | `'../../helpers.js'` |
| Regression (`tests/e2e/<module>/`) | `'../helpers.js'` |
| Deep regression (`tests/e2e/<module>/sub/`) | `'../../helpers.js'` |

## Shared Infrastructure Files

Written to `tests/e2e/` (not per-feature):

| File | Source template | Policy |
|------|----------------|--------|
| `helpers.ts` | `templates/helpers.ts` | Create if missing; merge missing exports if exists. Strip `// VERIFY:` and `// TEMPLATE:` comments on copy. |
| `package.json` | `templates/package.json` | Create if missing; do NOT modify if exists. |
| `tsconfig.json` | `templates/tsconfig.json` | Create if missing; do NOT modify if exists. |
| `playwright.config.ts` | `templates/playwright.config.ts` | Create if missing; webServer audit only modification allowed. |
| `config.yaml` | `plugins/forge/references/shared/config.yaml` | Create if missing; do NOT modify if exists. |
| `auth-setup.ts` | `templates/auth-setup.ts` | Create if auth-required-test cases exist AND UI tests generated. |

## Anti-Patterns (Forbidden)

**Fixed delays**: No `waitForTimeout`, `sleep`, `setTimeout`. Use event-driven waits:

| Forbidden | Replacement |
|-----------|-------------|
| `page.waitForTimeout(2000)` | `await expect(locator).toBeVisible()` or `await page.waitForResponse('**/api/...')` |
| `await new Promise(r => setTimeout(r, N))` | `await withRetry(() => ..., { delayMs: N })` |
| `waitForTimeout` after navigation | `await page.waitForLoadState('networkidle')` |
| `waitForTimeout` after click/submit | `await page.waitForResponse('**/api/...')` wrapping the click |

**Serial suite limit**: 15 tests max. Split larger suites into multiple `test.describe.serial` blocks with independent setup/teardown. `afterAll` cleanup is mandatory for suites creating shared data.

**beforeEach login**: Forbidden for auth-required tests. Each `login(page)` costs 1-3 seconds. Use `storageState` (preferred) or `test.beforeAll` (fallback). Exception: tests verifying login behavior or needing clean session.

**CSS class selectors**: Never use `.ant-*`, `.btn-*`, etc. or DOM traversal (`locator('..')`). Use semantic locators or `data-testid` as fallback.

**`@eN` refs**: Must not appear in any generated spec file.

## beforeAll Safety

1. Wrap each operation in try/catch with `console.error` identifying which step failed
2. Explicit undefined check after assignment: `if (!var) throw new Error('var is undefined after <operation>')`
3. Prefer reusing existing resources over creating new ones
4. Use `test.describe.serial` for sequential dependencies
5. Use `withRetry()` for backend-dependent operations (token acquisition, resource creation)

## Integration Tests

Integration tests are a subcategory of UI tests (component on existing page, via `placement: existing-page:<route>`). Same framework, same spec file (`ui.spec.ts`).

Script generation strategy:
1. Locate target page by route (sitemap resolution)
2. Locate embedded component: heading/section title → `aria-label`/`data-testid` → sitemap element
3. Assert visibility: `expect(locator).toBeVisible()`
4. Assert position (if feasible): verify adjacent sibling elements
5. Verify data rendering: `await expect(locator).toContainText('expected data')`

## Traceability

Each `test()` must include a traceability comment:
```typescript
// Traceability: TC-NNN → {PRD Source}
```

## Compilation Check

After generating all spec files, run the project's compilation recipe (delegates to `npx tsc --noEmit` under the hood):
```bash
just e2e-compile
```

## VERIFY Markers

Resolve all `// VERIFY:` comments using Fact Table values. Post-generation check:
```bash
just e2e-verify --feature <slug>
```
or: `grep -rn '// VERIFY:' tests/e2e/features/<slug>/ --include='*.spec.ts'`

## Playwright webServer Audit

When processing `playwright.config.ts` (new or existing):

1. Check for `webServer` key: `grep -n "webServer" tests/e2e/playwright.config.ts`
2. If found: remove the entire `webServer: { ... }` block, insert comment: `// Server lifecycle managed by justfile (embedded in test-e2e recipe)`
3. This is the ONLY permitted modification to an existing `playwright.config.ts`
