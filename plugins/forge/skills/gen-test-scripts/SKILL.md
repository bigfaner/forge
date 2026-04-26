---
name: gen-test-scripts
description: Generate executable TypeScript e2e test scripts from test cases. Uses Playwright for UI (semantic locators from sitemap.json), fetch for API, child_process for CLI. Zero external test frameworks (node:test + node:assert).
---

# Gen Test Scripts

Generate executable TypeScript e2e test scripts from test cases.

**Core principle**: Every generated `test()` must be independently runnable, repeatable, and have explicit assertions. Test cases are input; scripts are output.

<HARD-GATE>
This skill only generates test scripts (testing/scripts/), it does not execute tests.
Test execution is handled by the `/run-e2e-tests` skill.
</HARD-GATE>

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

| Artifact | Missing prompt |
|----------|----------------|
| `testing/test-cases.md` | Run `/gen-test-cases` first |
| `docs/sitemap/sitemap.json` (UI tests only) | Run `/gen-sitemap` first |
| `tests/e2e/config.yaml` | Run `/gen-sitemap` or create manually |

**Note**: This skill can be invoked manually or as the standard task T-test-2 appended by `/breakdown-tasks`.

```bash
ls docs/features/<slug>/testing/test-cases.md
ls docs/sitemap/sitemap.json  # UI tests only (project-level file, not per-feature)
```

**Note**: `<slug>` is the current feature name, obtained via `task feature` command. `docs/sitemap/sitemap.json` is a project-level file (one per application), not isolated per feature.

### Sitemap

Locator source for UI tests comes from `docs/sitemap/sitemap.json`. Generated and maintained by `/gen-sitemap` command.

**Location**: `docs/sitemap/sitemap.json`
**Full example**: `plugins/forge/references/shared/sitemap.json`

**Key locator fields**:

| Field | Usage |
|-------|-------|
| `layout.elements[].role` + `name` | Layout-level shared elements (navbar, etc.), available for all pages wrapped by `layout.wraps` |
| `elements[].role` + `name` | `page.getByRole(role, { name })` (priority 0) |
| `elements[].level` | `getByRole('heading', { level })` precise targeting |
| `elements[].label` | `page.getByLabel(label)` (priority 1) |
| `elements[].placeholder` | `page.getByPlaceholder(placeholder)` (priority 2) |
| `states[].trigger` | Click trigger element to enter dynamic state |
| `states[].elements` | Locators for elements within dynamic states |

Test cases match sitemap pages via the `Route` field. The `Element` field (referencing element IDs like E-NNN / L-NNN) is optional — when present, provides precise mapping; when absent, all page elements are used.

## When to Use

**Trigger:**

- User asks to "generate test scripts" or "create e2e scripts"
- User provides `/gen-test-scripts` command
- After `/gen-test-cases` has produced `testing/test-cases.md`

**Skip:**

- `testing/test-cases.md` doesn't exist (run `/gen-test-cases` first)

## Workflow

```
1. Read test cases → 2. Resolve sitemap → 3. Map locators → 4. Generate spec files → 5. Generate config files
```

### Step 1: Read Test Cases

Read `docs/features/<slug>/testing/test-cases.md`.

Parse each test case:

- Extract TC ID, title, type, route, feature, pre-conditions, steps, expected result, priority
- Group by type: UI / API / CLI
- Count each group

#### Auth Classification

For each test case, classify by authentication requirements:

| Category | Detection Rules | Generation Strategy |
|----------|----------------|---------------------|
| **login-test** | `Target` matches `ui/login`, `ui/auth`, `ui/signin`, `api/auth`, `api/login`, `api/token` | No shared auth. UI uses independent `loginPage`, API uses raw `curl()` |
| **auth-required-test** | `Pre-conditions` contain "authenticated", "logged in", or `Target` implies protected resource | Use shared auth |
| **public-test** | No auth-related pre-conditions, target is a public resource | No auth needed |

Count each category to decide whether to enable shared auth (enabled when `auth-required-test` exists).

<HARD-RULE>
Login tests and authenticated tests must not be mixed in the same `describe` block.
</HARD-RULE>

### Step 2: Resolve Sitemap

**Only execute when UI-type test cases exist.**

Read `docs/sitemap/sitemap.json`. For each route referenced in test cases:

1. Match `Route` field to `sitemap.json` `pages[].route` (dynamic routes like `/tasks/:id` match specific paths like `/tasks/123`)
2. Collect all element data for matched pages (base elements + related states)
3. **If test case has `Element` field** (referencing sitemap element IDs like E-NNN / L-NNN): use only specified elements for precise test-step-to-sitemap-element mapping
4. **If test case has no `Element` field**: use all available elements for that page to build locator mapping; the agent matches test step descriptions to the most appropriate elements
5. **If sitemap has `layout` field**: for each route wrapped by `layout.wraps`, merge `layout.elements` as available elements (layout elements are also available in that page's locator mapping)

If a route referenced in test cases does not exist in sitemap, report the missing route and suggest re-running `/gen-sitemap`.

### Step 3: Map Locators

Based on sitemap element data collected in Step 2, generate Playwright locators for each element by priority:

| Priority | Condition | Generated Code |
|----------|-----------|----------------|
| 0 (most stable) | `role` ∈ {button,link,heading,...} + `name` non-empty | `page.getByRole('button', { name: 'Submit' })` |
| 0+ | heading + `level` non-empty | `page.getByRole('heading', { name: 'Dashboard', level: 1 })` |
| 1 | `label` non-empty | `page.getByLabel('Email address')` |
| 2 | `placeholder` non-empty, label empty | `page.getByPlaceholder('Search...')` |
| 3 | Static text node | `page.getByText('No results found')` |
| 4 | `data-testid` visible | `page.getByTestId('user-avatar')` |
| 5 (fallback) | None of the above | `page.locator('.btn') // UNSTABLE: no semantic anchor` |

For test steps within dynamic states: first click the trigger element's locator, then map locators for in-state elements.

Build an in-memory mapping table for use in Step 4.

### Step 4: Generate Spec Files

For each type group, generate a spec file from the corresponding template.

**Template usage**: Templates contain `CONDITIONAL` comment blocks marking code segments to enable/disable based on auth classification. Based on Step 1 auth classification results, **uncomment** matching CONDITIONAL blocks, remove non-matching blocks, then fill in test data. Do not rewrite template structure from scratch.

<EXTREMELY-IMPORTANT>
**UI tests use Playwright Locator API** (browser driver library only). Using agent-browser in generated spec files is forbidden.

**API tests use Node.js built-in `fetch`**. External HTTP libraries like axios, supertest are forbidden.

**CLI tests use `child_process.execSync`** via `runCli()` helper.

**All tests use `node:test` + `node:assert`**. External test frameworks like jest, mocha, vitest are forbidden.

**`@eN` refs must not appear in any generated spec file.**
</EXTREMELY-IMPORTANT>

**UI tests (`testing/scripts/ui.spec.ts`)**:

- Read template: `plugins/forge/skills/gen-test-scripts/templates/playwright-ui.spec.ts`
- **Auth setup** (top-level `before` hook):
  - If `auth-required-test` exists: call `await loginViaUI(page)` after `setupBrowser()`
  - If only `public-test` and/or `login-test`: do not call `loginViaUI()`
- **Login test cases**: place in nested `describe('Login', ...)` block
  - Create independent `loginPage` via `page.context().newPage()` (no auth cookie)
  - After block ends, `loginPage.close()`
  - No shared auth
- **Authenticated test cases** + **public-test cases**: at top level, use shared `page` (logged-in state)
- Use locator mapping table from Step 3 to replace template example locators
- Each `test()` uses `await page.xxx.waitFor()` instead of `assert.ok(snapshotContains())`
- Keep `screenshot(page, 'TC-NNN')` calls
- Fallback locators annotated with `// UNSTABLE: no semantic anchor`
- Each test includes a comment: `// Traceability: TC-NNN → {Source}`

**API tests (`testing/scripts/api.spec.ts`)**:

- Read template: `plugins/forge/skills/gen-test-scripts/templates/api.spec.ts`
- **Auth setup** (top-level `before` hook):
  - If `auth-required-test` exists: call `const token = await getApiToken(apiBaseUrl)` and create `authCurl = createAuthCurl(apiBaseUrl, token)`
  - If only `public-test` and/or `login-test`: no auth setup
- **Login/auth API test cases**: use raw `curl()` (no Bearer header)
- **Authenticated test cases**: use `authCurl(method, path)` (auto-inject Bearer token)
- **Public test cases**: use `curl()` (no auth)
- For each API test case, generate a `test()` block:
  - Assert on status code, response body, headers
  - Each test includes a traceability comment

**CLI tests (`testing/scripts/cli.spec.ts`)**:

- Read template: `plugins/forge/skills/gen-test-scripts/templates/cli.spec.ts`
- For each CLI test case, generate a `test()` block:
  - Use `runCli()` helper to execute commands
  - Assert on stdout, stderr, exit code
  - Each test includes a traceability comment

<HARD-RULE>
**Traceability**: Each `test()` must include a traceability comment `// Traceability: TC-NNN → {PRD Source}` ensuring complete traceability chain from test script to PRD.

**Skip empty groups**: If no test cases of a given type exist, skip generating that spec file.
</HARD-RULE>

### Step 5: Generate Config Files

Write `testing/scripts/package.json` from template:

- `tsx` + `playwright` as devDependencies
- Scripts: `"test:ui"`, `"test:api"`, `"test:cli"`, `"test:all"`

Write `testing/scripts/tsconfig.json` from template:

- Target ES2022, module NodeNext, strict mode

## Output Files

All files go to `docs/features/<slug>/testing/scripts/`:

```
scripts/
  helpers.ts       # Shared utilities (browser lifecycle, curl, auth helpers, runCli)
  ui.spec.ts       # UI tests via Playwright (shared auth via loginViaUI)
  api.spec.ts      # API tests via fetch (shared auth via authCurl)
  cli.spec.ts      # CLI tests via child_process
  package.json     # tsx + playwright
  tsconfig.json    # ES2022 + NodeNext
```

## Locator Priority Reference

When generating UI test code, use Playwright locators following this priority:

| Priority | Method                    | Example                                              |
| -------- | ------------------------- | ---------------------------------------------------- |
| 0        | `page.getByRole()`        | `page.getByRole('button', { name: 'Submit' })`       |
| 0+       | `page.getByRole()`        | `page.getByRole('heading', { name: 'Dashboard', level: 1 })` |
| 1        | `page.getByLabel()`       | `page.getByLabel('Email address')`                   |
| 2        | `page.getByPlaceholder()` | `page.getByPlaceholder('Search...')`                 |
| 3        | `page.getByText()`        | `page.getByText('No results found')`                 |
| 4        | `page.getByTestId()`      | `page.getByTestId('user-avatar')`                    |
| 5        | `page.locator()`          | `page.locator('.btn') // UNSTABLE`                   |

## Related Skills

| Skill             | Usage                                   |
| ----------------- | --------------------------------------- |
| `/gen-test-cases` | Generate test cases from PRD            |
| `/gen-sitemap`    | Generate and maintain sitemap.json      |
| `/run-e2e-tests`  | Execute test scripts and report results |
