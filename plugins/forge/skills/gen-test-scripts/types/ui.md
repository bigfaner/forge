---
type: ui
conventions:
  - testing-ui.md
  - frontend.md
---

# UI Type Instructions

Type-specific generation instructions for **UI** (browser DOM interaction) test scripts. Loaded by the dispatcher after interface detection.

## Reconnaissance Strategy

UI-specific source code patterns to search during Step 1.5 Code Reconnaissance:

| Pattern | What to Extract | Discovery Method |
|---------|----------------|------------------|
| `data-testid` attributes | Static testid values, dynamic testid patterns | Grep for `data-testid` in component files (`.tsx`, `.vue`, `.svelte`). Note static vs dynamic patterns (e.g., `map-card-${id}` vs `map-card`). |
| Component files | Component structure, props, conditional rendering | Search for component definitions (React: `export default function`, `export const`; Vue: `<template>`; Svelte: `<script>`) |
| Route configs | Page routes, path parameters, nested routes | Search for route registration: React SPA `path="` or `path='` in route config, `<Route` component; Next.js `app/` directory structure, `page.tsx` files |
| Dynamic testid patterns | Template literals constructing testids | Grep for backtick patterns containing `data-testid` or `testid` with `${` interpolation |

## Fact Table Required Keys

Minimum Fact Table entries required for UI type (all must be non-UNKNOWN):

| Required Key Pattern | Example |
|---------------------|---------|
| `FRONTEND_BASE` | `http://localhost:3000` |
| `TESTID_*` entries | `TESTID_BTN_CONFIRM: btn-confirm`, `TESTID_MAP_CARD: map-card-${bizKey}` |

**If all required keys for UI are UNKNOWN**: skip UI test generation, emit a WARNING explaining why, and suggest the user verify frontend source files exist before re-running. Do NOT generate UI tests with UNKNOWN values.

**If some keys are non-UNKNOWN**: proceed. Individual UNKNOWN keys are acceptable (e.g., `TESTID_MAP_CARD` unknown but `TESTID_BTN_CONFIRM` known).

## Sitemap Resolution

**Only execute when the active profile has `web-ui` capability AND UI-type test cases exist.**

Read `docs/sitemap/sitemap.json`. For each route referenced in test cases:

1. Match `Route` field to `sitemap.json` `pages[].route` (dynamic routes like `/tasks/:id` match specific paths like `/tasks/123`)
2. Collect page structure data for matched pages (base elements + related states) as supplementary reference
3. **If sitemap has `layout` field**: for each route wrapped by `layout.wraps`, merge `layout.elements` as available page structure elements

The sitemap serves as a secondary reference for understanding page structure. Locators are NOT derived from sitemap element data keyed by test-case Element IDs -- all locator derivation comes from the Fact Table (Step 1.5).

If a route referenced in test cases does not exist in sitemap, **emit a WARNING** listing the missing routes and suggest re-running `/gen-sitemap`. Proceed using the Fact Table from Step 1.5 for the missing routes -- do not abort.

## Locator Mapping

Locator strategy is owned by this step -- `generate.md` provides only the framework-specific syntax (how to write a Playwright/Go/pytest locator), not the decision of *which* locator to use.

**Locator priority** (all locator types, including integration tests):

1. **Fact Table data** (from Step 1.5 Code Reconnaissance): `TESTID_*` entries, semantic roles, labels extracted from source code -> translate via `generate.md`. For dynamic testids (e.g., `map-card-${bizKey}`), use prefix selector (`page.locator('[data-testid^="map-card-"]')`)
2. **Sitemap page structure** (from Sitemap Resolution): role+name, label, placeholder -- used as supplementary confirmation of Fact Table data only, NOT as primary locator source -> translate via `generate.md`
3. **Semantic inference**: heading/section title, aria-label from source code -> translate via `generate.md`

<HARD-RULE>
Never guess `data-testid` values. Every testid locator must come from a Fact Table `TESTID_*` entry. If the Fact Table has no testid entry for the needed element, use semantic locators (role, label, text) instead.
</HARD-RULE>

<HARD-RULE>
**Source-code-first locator derivation**: Derive all locators from source code (Fact Table). Do NOT reference any testid, selector, or locator from test-cases.md -- test-cases provides scenario context only. The Fact Table built in Step 1.5 is the sole primary locator source. The sitemap (Sitemap Resolution) is supplementary confirmation, not a locator source driven by test-case Element IDs.
</HARD-RULE>

For test steps within dynamic states: first click the trigger element's locator, then map locators for in-state elements. Build an in-memory mapping table for use in spec generation.

**Integration test component locators**: Use the same priority chain above. Locate the embedded component via Fact Table data (source code) first, then confirm with sitemap page structure, then semantic inference.

## Verification Method

Confirm the project exposes a UI interface before generating test scripts:

| Check | Method |
|-------|--------|
| `data-testid` in frontend source | `grep -rn "data-testid" --include='*.tsx' --include='*.vue' --include='*.svelte' .` |
| `TESTID_*` or `FRONTEND_*` Fact Table entries | Verify Fact Table from Step 1.5 contains at least one non-UNKNOWN entry |

**Evidence of product interface**: Frontend source code was found and read during reconnaissance.

## Generation Patterns

How UI test cases translate to executable scripts:

| Test Case Action | Script Pattern |
|-----------------|----------------|
| Navigate to page | Framework navigation to route URL (e.g., `page.goto('/login')`) |
| Click element | Locate via Fact Table testid, then click (e.g., `page.getByTestId('btn-submit').click()`) |
| Fill form field | Locate input, clear, type value (e.g., `page.getByTestId('input-email').fill('user@test.com')`) |
| Submit form | Click submit button, wait for navigation or response |
| Assert visibility | Framework visibility assertion (e.g., `expect(page.getByTestId('dashboard')).toBeVisible()`) |
| Assert text content | Framework text assertion (e.g., `expect(page.getByTestId('welcome-msg')).toHaveText('Hello')`) |
| Assert navigation | Verify URL changed to expected path |
| Screenshot capture | Framework screenshot method for visual regression or debug evidence |

All DOM interactions must use locators derived from the Locator Mapping priority chain (Fact Table > Sitemap > Semantic inference).

## Integration Test Scripts

Integration test cases verify that a component embedded on an existing page is visible and renders correctly. Follow the profile's `generate.md` for framework-specific patterns.

**Identification**: A test case is an integration test when:
- `Source` field contains "Placement + Integration Spec", **or**
- `Test ID` field contains "integration-"

Integration tests use the same framework and file structure as standard tests. They are grouped with UI test cases and emitted into the same test file.

**Script generation strategy** -- for each integration test case:

1. **Locate the target page** by route (using sitemap resolution from Sitemap Resolution step)
2. **Locate the embedded component** using the locator strategy from `generate.md`
3. **Assert visibility**: use framework-specific visibility assertion
4. **Assert position** (if feasible): verify the component appears at expected position
5. **Verify data rendering**: assert expected text content is present within the component locator

## UI Antipattern Guards

Beyond the generic 6 antipattern guards in the dispatcher, UI-specific forbidden patterns:

| Pattern | Why Forbidden | What To Do Instead |
|---------|--------------|-------------------|
| **CSS class selectors** | CSS classes are implementation details that change during styling refactors, breaking tests without functional regressions | Use `data-testid` locators from Fact Table, or semantic locators (role, label, text) |
| **Fixed delays/sleeps** | Mask real timing issues, waste time when the app is fast, fail when the app is slow | Use event-driven waits (wait for element visible, wait for navigation) or polling with timeout |
| **Debug output** (`console.log`, `print`) in generated code | Bloats CI output, hides real failures in noise | Only assertions in generated test code |
| **Missing testid fallback to CSS selectors** | When Fact Table has no testid, falling back to CSS selectors creates fragile tests | Use semantic locators (role, label, text) as the fallback, never CSS selectors |
| **Screenshot-only assertions** | Screenshots are non-deterministic across environments and provide weak signal | Use structural assertions (visibility, text content) as primary; screenshots as supplementary evidence only |
