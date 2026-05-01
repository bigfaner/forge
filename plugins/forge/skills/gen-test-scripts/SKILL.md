---
name: gen-test-scripts
description: Generate executable TypeScript e2e test scripts from test cases. Uses @playwright/test for all tests (no node:test or node:assert). Playwright for UI, fetch for API, child_process for CLI.
---

# Gen Test Scripts

Generate executable TypeScript e2e test scripts from test cases.

**Core principle**: Every generated `test()` must be independently runnable, repeatable, and have explicit assertions. Test cases are input; scripts are output.

<HARD-GATE>
This skill only generates test scripts (`tests/e2e/<feature>/`), it does not execute tests.
Test execution is handled by the `/run-e2e-tests` skill.
</HARD-GATE>

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

| Artifact | Missing prompt |
|----------|----------------|
| `testing/test-cases.md` | Run `/gen-test-cases` first |
| `docs/sitemap/sitemap.json` (UI tests only) | Run `/gen-sitemap` first |
| `tests/e2e/config.yaml` | Run `/gen-sitemap` or create from template at `plugins/forge/skills/gen-test-scripts/templates/config.yaml` |

**Note**: This skill can be invoked manually or as the standard task T-test-2 appended by `/breakdown-tasks`.

```bash
ls docs/features/<slug>/testing/test-cases.md
ls docs/sitemap/sitemap.json  # UI tests only (project-level file, not per-feature)
ls tests/e2e/config.yaml      # Test environment config (project-level)
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
1. Read test cases → 1.5. Code Reconnaissance → 2. Resolve sitemap → 3. Map locators → 4. Generate spec files → 5. Ensure shared infrastructure
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
| **custom-auth-test** | `Pre-conditions` mention "API key", "X-API-Key", "OAuth", "session cookie", or other non-Bearer auth patterns | Detect auth mechanism from codebase during Step 1.5. Generate custom auth setup in `test.beforeAll` — e.g., API key via `curl()` with custom header, or cookie-based session |

Count each category to decide whether to enable shared auth (enabled when `auth-required-test` exists).

<HARD-RULE>
Login tests and authenticated tests must not be mixed in the same `describe` block.
</HARD-RULE>

### Step 1.5: Code Reconnaissance (Build Fact Table)

Read actual source code files to extract ground-truth values. **Never guess or assume values** — every value in a generated script must come from the test-cases.md input or the Fact Table built here.

**Check test-cases.md** for Route Validation results (⚠️ warnings from gen-test-cases Step 3.5). Use corrected routes where available.

**Required reads** (adapt to project structure):

| Source | What to extract | Typical file patterns |
|--------|----------------|----------------------|
| Router files | Route paths, path parameters, middleware bindings | `internal/handler/router.go`, `src/routes/**`, `src/app/**/page.tsx` |
| Config files | API port, base path prefix, auth credentials | `config.yaml`, `.env`, `.env.example`, `src/config.*` |
| API handlers | Request/response schemas, status codes, validation rules | `internal/handler/*.go`, `src/handlers/**` |
| Auth implementation | Login endpoint path, token field name, header format | Auth module/controller files |
| CLI entry points | Command names, flag names, output formats | `src/cli/**`, `cmd/**/*.go` |

**Build the Fact Table**: After reading, record verified facts with source citations:

```markdown
## Fact Table
| Key | Value | Source |
|-----|-------|--------|
| API_PORT | 8080 | backend/config.yaml:3 |
| API_PREFIX | /v1 | internal/handler/router.go:20 |
| AUTH_ENDPOINT | POST /v1/auth/login | internal/handler/router.go:42 |
| AUTH_TOKEN_FIELD | data.token | internal/handler/auth.go:58 |
| AUTH_HEADER | Authorization: Bearer <token> | internal/middleware/auth.go:15 |
```

<HARD-RULE>
- Every value in the Fact Table must cite the source file and line number.
- If a source file cannot be found, note it as `UNKNOWN`. Do not fabricate values.
- The agent must use Fact Table values when generating spec files in Step 4. When a Fact Table value contradicts a template placeholder, the Fact Table wins.
- All `// VERIFY:` markers in templates must be resolved using Fact Table values.
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

**Verify project interfaces before generating**: For each type group from Step 1, confirm the project actually exposes that interface. If a type group has test cases but the project lacks that interface (no evidence in codebase that this is a product interface, not developer tooling): warn the user, then **skip that spec file entirely**. Key distinction: build/test/lint commands (`go build`, `grep`, `npm test`) are developer tooling, not a CLI product interface.

For each non-empty, verified type group, generate a spec file from the corresponding template.

**Template usage**: Templates contain `CONDITIONAL` comment blocks marking code segments to enable/disable based on auth classification. Based on Step 1 auth classification results, **uncomment** matching CONDITIONAL blocks, remove non-matching blocks, then fill in test data. Do not rewrite template structure from scratch.

**Import path**: All spec files must import from `'../helpers.js'` (one level up to shared helpers.ts at `tests/e2e/`).

**VERIFY marker resolution**: All `// VERIFY:` comments in templates must be resolved during generation:
1. Look up the corresponding value in the Fact Table (Step 1.5)
2. Replace the placeholder value with the Fact Table value
3. Remove the `// VERIFY:` comment
4. If no Fact Table value exists, keep the `// VERIFY:` comment as-is so it is visible during code review

**Post-generation check**: After generating all spec files, verify no unresolved `// VERIFY:` markers remain:
1. Run `just e2e-verify --feature <slug>` if the recipe exists in the Justfile
2. Otherwise, fall back to `grep -r '// VERIFY:' tests/e2e/<slug>/`
3. Any remaining `// VERIFY:` = skill incomplete — do not proceed to run-e2e-tests until all markers are resolved

<EXTREMELY-IMPORTANT>
**UI tests use `@playwright/test`** (full test runner with automatic browser lifecycle). Using agent-browser in generated spec files is forbidden.

**API tests use Node.js built-in `fetch`**. External HTTP libraries like axios, supertest are forbidden.

**CLI tests use `child_process.execSync`** via `runCli()` helper.

**All tests use `@playwright/test`** (`test`, `expect`, `test.describe`). The `node:test` and `node:assert` modules are not used.

**`@eN` refs must not appear in any generated spec file.**
</EXTREMELY-IMPORTANT>

**UI tests (`tests/e2e/<feature>/ui.spec.ts`)**:

- Read template: `plugins/forge/skills/gen-test-scripts/templates/playwright-ui.spec.ts`
- **Auth setup** (`test.beforeEach` hook):
  - If `auth-required-test` exists: call `await loginViaUI(page)` using the injected `page` fixture
  - If only `public-test` and/or `login-test`: do not call `loginViaUI()`
- **Login test cases**: place in nested `test.describe('Login', ...)` block
  - Each test gets its own `page` via Playwright fixture injection
  - No shared auth state between tests
- **Authenticated test cases** + **public-test cases**: at top level, use shared `page` from fixture
- Use locator mapping table from Step 3 to replace template example locators
- Use `await expect(locator).toBeVisible()` pattern instead of `assert.ok()`
- Keep `screenshot(page, 'TC-NNN')` calls
- Fallback locators annotated with `// UNSTABLE: no semantic anchor`
- Each test includes a comment: `// Traceability: TC-NNN → {Source}`

**API tests (`tests/e2e/<feature>/api.spec.ts`)**:

- Read template: `plugins/forge/skills/gen-test-scripts/templates/api.spec.ts`
- **Auth setup** (`test.beforeAll` block):
  - If `auth-required-test` exists: call `const token = await getApiToken(apiBaseUrl, authPath)` (authPath from Fact Table `AUTH_ENDPOINT` path) and create `authCurl = createAuthCurl(apiBaseUrl, token)`
  - If `custom-auth-test` exists: write custom auth setup from scratch in `test.beforeAll` using `curl()` directly — do not use `getApiToken` (it assumes username/password body schema). See CONDITIONAL scaffold in template.
  - If only `public-test` and/or `login-test`: no auth setup
- **Login/auth API test cases**: use raw `curl()` (no Bearer header)
- **Authenticated test cases**: use `authCurl(method, path)` (auto-inject Bearer token)
- **Custom-auth test cases**: use `curl()` with custom auth headers from `test.beforeAll`
- **Public test cases**: use `curl()` (no auth)
- For each API test case, generate a `test()` block:
  - Assert on status code, response body using `expect()`
  - Each test includes a traceability comment

**CLI tests (`tests/e2e/<feature>/cli.spec.ts`)**:

- Read template: `plugins/forge/skills/gen-test-scripts/templates/cli.spec.ts`
- For each CLI test case, generate a `test()` block:
  - Use `runCli()` helper to execute commands
  - Assert on stdout, stderr, exit code using `expect()`
  - Each test includes a traceability comment

<HARD-RULE>
**Traceability**: Each `test()` must include a traceability comment `// Traceability: TC-NNN → {PRD Source}` ensuring complete traceability chain from test script to PRD.

**Skip empty groups**: If no test cases of a given type exist, or the project lacks that interface (Step 4 verification), skip generating that spec file.

**Empty result guard**: If zero spec files were generated (all groups empty or all interfaces undetected), abort with a clear message: "No spec files generated — either test-cases.md has no testable cases, or all interface types were undetected. Re-run /gen-test-cases or verify project structure."
</HARD-RULE>

### Step 5: Ensure Shared Infrastructure

Check if shared infrastructure exists at `tests/e2e/`:

```bash
ls tests/e2e/helpers.ts tests/e2e/package.json tests/e2e/tsconfig.json tests/e2e/playwright.config.ts 2>/dev/null
```

**If any file is missing**, create it from the corresponding template:
- `helpers.ts`: copy from `plugins/forge/skills/gen-test-scripts/templates/helpers.ts` (only if tests/e2e/helpers.ts does not exist)
- `package.json`: copy from `plugins/forge/skills/gen-test-scripts/templates/package.json` (only if tests/e2e/package.json does not exist)
- `tsconfig.json`: copy from `plugins/forge/skills/gen-test-scripts/templates/tsconfig.json` (only if tests/e2e/tsconfig.json does not exist)
- `playwright.config.ts`: copy from `plugins/forge/skills/gen-test-scripts/templates/playwright.config.ts` (only if tests/e2e/playwright.config.ts does not exist)

**If all exist**: skip. Do not modify existing shared files.

Install dependencies if `node_modules` is missing:

```bash
just e2e-setup
```

## Output Files

All generated spec files go to `tests/e2e/<feature>/`:

```
tests/e2e/                        # Shared infrastructure (created once)
  helpers.ts                      # Shared utilities (curl, auth, runCli, screenshot)
  package.json                    # @playwright/test + yaml
  tsconfig.json                   # ES2022 + NodeNext
  config.yaml                     # Test environment config (already exists)
  <feature>/                      # Generated per-feature (only for types detected in project)
    ui.spec.ts                    # [if UI detected] UI tests via Playwright
    api.spec.ts                   # [if API detected] API tests via fetch
    cli.spec.ts                   # [if CLI detected] CLI tests via child_process
```

<HARD-RULE>
**Shared file policy**: `helpers.ts`, `package.json`, `tsconfig.json`, and `playwright.config.ts` at `tests/e2e/` are shared across all features.
- If they do not exist: create them from templates.
- If they already exist: DO NOT overwrite. The existing versions are canonical.
- Only spec files (`*.spec.ts`) are written per-feature.
</HARD-RULE>

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
