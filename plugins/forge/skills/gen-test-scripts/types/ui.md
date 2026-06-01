---
type: ui
conventions:
  - testing-ui.md
  - frontend.md
---

# UI Type Instructions

Type-specific generation instructions for **UI** (browser DOM interaction) test scripts. Loaded by the dispatcher after interface detection.

**Test type**: Web 端到端测试 (Web E2E Test). See `docs/reference/test-type-model.md` for the authoritative definition. Generated test code MUST use `@web-e2e` tags. This is one of the two surfaces where "e2e" terminology is correct (the other being Mobile).

This file defines two zones:

- **Golden Rules**: Framework-agnostic constraints that govern all UI test generation. These rules are mandatory and cannot be overridden by Convention files.
- **Reconnaissance Hints**: Discovery-only guidance for extracting UI-specific information from source code. Hints inform the Fact Table, not the generated code directly.

## Golden Rules

UI-specific constraints beyond the universal principles in `_shared.md` (Isolation, Determinism, Timeout Protection, Idempotency, Resource Cleanup).

### Session Reuse

Authentication scenarios must login once and reuse the session context across all tests within the same Journey. Per-test login is forbidden.

**Constraint**: Establish authentication state once at Journey setup. All subsequent tests inherit the authenticated session. No test may independently perform login unless the test case explicitly tests the authentication flow itself.

**Rationale**: Per-test login multiplies execution time by the number of tests, inflates CI minutes, and introduces flakiness from repeated auth service calls. Session reuse reduces a 50-test suite from 50 logins to 1.

**Antipattern guard**: Tests that independently navigate to a login page and fill credentials when they are not testing login behavior.

### Network Interception

For tests that depend on external services (third-party APIs, payment gateways, email services), intercept network requests and return fixed responses. Tests must not require network access to pass.

**Constraint**: External service calls must be replaced with deterministic responses via the framework's interception or mocking mechanism. The intercepted response data comes from the Fact Table or test case fixture data, never from live service calls.

**Rationale**: External service dependencies cause non-deterministic test failures. When the service is down or returns unexpected data, the test fails for reasons unrelated to the code under test.

**Antipattern guard**: Tests that call real third-party APIs, even in "read-only" or "safe" modes. Network access must not be required for tests to pass. See also `_shared.md` Principle: Determinism sub-dimension (b).

### Viewport Management

All UI tests must define an explicit default viewport size. Responsive testing uses explicit viewport switching as a pre-step, not as a default.

**Constraint**: The Convention defines the default viewport dimensions. Tests that verify responsive behavior must explicitly set the target viewport as a setup step before assertions. No test may depend on the runner's default viewport.

**Rationale**: Tests that pass at one viewport silently fail at another. CI runners often use non-standard viewports. Explicit viewport declaration makes the test's visual assumptions clear and reproducible.

**Antipattern guard**: Tests that assume a specific screen width without declaring it, or tests that verify responsive breakpoints without explicitly switching the viewport first.

### Element Location Strategy

Locator selection follows an abstract priority chain, not a framework-specific selector syntax:

1. **Test ID**: Use data-testid (or equivalent stable test attribute) from Fact Table entries
2. **Accessible role + name**: Use semantic role and accessible name (button named "Submit", heading "Dashboard")
3. **Label**: Use form field labels, aria-label, placeholder text
4. **Text content**: Use visible text as a last resort for interactive elements

**Constraint**: CSS class selectors are forbidden. Locators derived from CSS classes couple tests to styling implementation details that change during refactors without functional regressions. See also Antipattern Guards below.

**Hard Rule**: Never guess test ID values. Every test ID locator must come from a Fact Table entry. If the Fact Table has no entry for the needed element, use semantic locators instead.

### Locator Source Authority

Derive all locators from source code via the Fact Table (Step 1 Code Reconnaissance). The sitemap serves as supplementary confirmation only, not a primary locator source.

**Constraint**: Locator derivation priority:
1. Fact Table data (source code reconnaissance) -- primary source
2. Sitemap page structure -- supplementary confirmation only
3. Semantic inference from source code -- fallback

Test case files provide scenario context only, never locators or selectors.

## Reconnaissance Hints

<!-- Discovery-only: information extracted here populates the Fact Table. These hints do not directly guide code generation. -->

UI-specific source code patterns to search during Step 1 Code Reconnaissance:

| Pattern | What to Extract | Discovery Method |
|---------|----------------|------------------|
| `data-testid` attributes | Static testid values, dynamic testid patterns | Grep for `data-testid` in component files (`.tsx`, `.vue`, `.svelte`). Note static vs dynamic patterns (e.g., `map-card-${id}` vs `map-card`). |
| Component files | Component structure, props, conditional rendering | Search for component definitions (React: `export default function`, `export const`; Vue: `<template>`; Svelte: `<script>`) |
| Route configs | Page routes, path parameters, nested routes | Search for route registration: React SPA `path="` or `path='` in route config, `<Route` component; Next.js `app/` directory structure, `page.tsx` files |
| Dynamic testid patterns | Template literals constructing testids | Grep for backtick patterns containing `data-testid` or `testid` with `${` interpolation |
| Auth implementation | Login endpoint, session mechanism, token handling | Search for auth-related patterns in API handlers and frontend auth modules |

## Classification Indicators

Classify test cases as **UI** when they involve any of:

- Browser DOM interaction (click, type, select, hover)
- Page navigation and routing
- Form submission and validation
- Element visibility and content assertions
- Session and authentication state
- Responsive layout behavior

## Fact Table Required Keys

Minimum Fact Table entries required for UI type (all must be non-UNKNOWN):

| Required Key Pattern | Example |
|---------------------|---------|
| `FRONTEND_BASE` | `http://localhost:3000` |
| `TESTID_*` entries | `TESTID_BTN_CONFIRM: btn-confirm`, `TESTID_MAP_CARD: map-card-${bizKey}` |

**If all required keys for UI are UNKNOWN**: skip UI test generation, emit a WARNING explaining why, and suggest the user verify frontend source files exist before re-running. Do NOT generate UI tests with UNKNOWN values.

**If some keys are non-UNKNOWN**: proceed. Individual UNKNOWN keys are acceptable (e.g., `TESTID_MAP_CARD` unknown but `TESTID_BTN_CONFIRM` known).

## Sitemap Resolution

**Only execute when the project has `web-ui` interface AND UI-type test cases exist.**

Read `docs/sitemap/sitemap.json`. For each route referenced in test cases:

1. Match `Route` field to `sitemap.json` `pages[].route` (dynamic routes like `/tasks/:id` match specific paths like `/tasks/123`)
2. Collect page structure data for matched pages (base elements + related states) as supplementary reference
3. **If sitemap has `layout` field**: for each route wrapped by `layout.wraps`, merge `layout.elements` as available page structure elements

The sitemap serves as a secondary reference for understanding page structure. Locators are NOT derived from sitemap element data keyed by test-case Element IDs -- all locator derivation comes from the Fact Table (Step 1).

If a route referenced in test cases does not exist in sitemap, **emit a WARNING** listing the missing routes and suggest re-running `/gen-web-sitemap`. Proceed using the Fact Table from Step 1 for the missing routes -- do not abort.

## Verification Method

Confirm the project exposes a UI interface before generating test scripts:

| Check | Method |
|-------|--------|
| `data-testid` in frontend source | `grep -rn "data-testid" --include='*.tsx' --include='*.vue' --include='*.svelte' .` |
| `TESTID_*` or `FRONTEND_*` Fact Table entries | Verify Fact Table from Step 1 contains at least one non-UNKNOWN entry |

**Evidence of product interface**: Frontend source code was found and read during reconnaissance.

## Integration Test Scripts

Integration test cases verify that a component embedded on an existing page is visible and renders correctly. Follow the active Convention for framework-specific patterns.

**Identification**: A test case is an integration test when:
- `Source` field contains "Placement + Integration Spec", **or**
- `Test ID` field contains "integration-"

Integration tests use the same framework and file structure as standard tests. They are grouped with UI test cases and emitted into the same test file.

**Script generation strategy** -- for each integration test case:

1. **Locate the target page** by route (using sitemap resolution from Sitemap Resolution step)
2. **Locate the embedded component** using the locator strategy from the active Convention
3. **Assert visibility**: use framework-specific visibility assertion
4. **Assert position** (if feasible): verify the component appears at expected position
5. **Verify data rendering**: assert expected text content is present within the component locator

## UI Antipattern Guards

Beyond the shared antipattern guards in `_shared.md` (Sleep-Based Waits, Hardcoded Configuration, Vacuous Assertions, Source-Code-Level Testing), UI-specific forbidden patterns:

| Pattern | Why Forbidden | What To Do Instead |
|---------|--------------|-------------------|
| **CSS class selectors** | CSS classes are implementation details that change during styling refactors, breaking tests without functional regressions | Use `data-testid` locators from Fact Table, or semantic locators (role, label, text) |
| **Debug output** (`console.log`, `print`) in generated code | Bloats CI output, hides real failures in noise | Only assertions in generated test code |
| **Missing testid fallback to CSS selectors** | When Fact Table has no testid, falling back to CSS selectors creates fragile tests | Use semantic locators (role, label, text) as the fallback, never CSS selectors |
| **Screenshot-only assertions** | Screenshots are non-deterministic across environments and provide weak signal | Use structural assertions (visibility, text content) as primary; screenshots as supplementary evidence only |
| **Per-test authentication** | Multiplies execution time, introduces flakiness from repeated auth service calls | Login once at Journey setup, reuse session across tests (see Golden Rule: Session Reuse) |

## Test Ratio Constraint

Web surface targets a **balanced 50/50** ratio between Contract tests and Journey smoke tests.

- **Formula**: `Contract test functions / (Contract test functions + Journey smoke test functions) × 100%`
- **Target**: Approximately 50% Contract tests, 50% Journey smoke tests
- **Implementation**: Generate Contract tests for each Outcome, PLUS generate enriched Journey smoke tests that cover both happy path AND at least 1 failure path through the Journey
- **Minimum**: Every Journey MUST have at least 1 smoke test (happy path). The "balanced" target means the smoke test suite should be substantive — not just a trivial single-assertion test

**Ratio guidance**: If the generated plan skews heavily toward Contract tests (>70%), add additional Journey smoke test scenarios (e.g., alternative failure paths, state transition tests). If the plan skews toward Journey smoke tests (>70%), ensure Contract tests cover individual Outcomes adequately.

## Output

UI test scripts are written to `tests/<journey>/` following the active Convention's file naming and structure. Each test file includes a traceability comment linking back to the source Contract step.
