---
type: ui
conventions:
  - testing-ui.md
  - frontend.md
---

# UI Test Case Generation Instructions

Type-specific Steps 3-4 for **UI** (browser DOM interaction) test cases. Loaded by the dispatcher after Step 2.5 interface detection.

## Classification Indicators

Classify a PRD criterion as **UI** when it involves any of:

- Page rendering and visual state
- Navigation between pages
- Form input, modals, tabs, dropdowns
- Component visibility and responsive behavior
- Interactions with browser DOM elements
- Screen layout and visual feedback

## Target Derivation

- **Target format**: `ui/<page-name>`
- Derive `<page-name>` from the URL path or component name (e.g., `ui/login`, `ui/user-profile`, `ui/dashboard`)

## Test ID Format

- **Test ID**: `<target>/<title-slug>`
- `title-slug` = lowercase title, spaces to hyphens, remove punctuation
- Example: `ui/login/valid-credentials-redirect-to-dashboard`

## Priority Assignment

1. Criterion tied to a core/critical Given/When/Then in the PRD → **P0**
2. Criterion tied to a secondary story, or an explicit error/boundary case for a core story → **P1**
3. Nice-to-have verifications, minor edge cases → **P2**

If the PRD has no explicit priority marking, default P0 for the first story's ACs and P1 for all others.

## TC Format

```markdown
## TC-{NNN}: {Title}
- **Source**: {Story N / AC-N} or {Spec Section X.Y} or {UI Function Name}
- **Type**: UI
- **Target**: ui/<page-name>
- **Test ID**: ui/<page-name>/<title-slug>
- **Pre-conditions**: {What must be true before testing}
- **Route**: {Page route path, e.g. /login}
- **Steps**:
  1. {Step 1}
  2. {Step 2}
- **Expected**: {What the correct result looks like}
- **Priority**: P0 | P1 | P2
```

- `Route` field is required for UI test cases — must contain a concrete path (e.g., `/users/123/edit`), not vague descriptions.
- Do not include implementation details (testid, selector, CSS) in Route or any field.
- Steps must describe natural-language user actions, not technical locators.

## Integration TC Generation (Existing-Page Placements)

When `prd/prd-ui-functions.md` exists and contains UI Functions with `placement: existing-page:<route>`, generate a dedicated integration test case for each:

```markdown
## TC-{NNN}: Integration — {Component} visible on {Page}
- **Source**: PRD UI Function "{Function Name}" Placement + Integration Spec
- **Type**: UI
- **Target**: ui/<page-name>
- **Test ID**: ui/<page-name>/integration-<component-slug>
- **Pre-conditions**: Component build complete, integration task complete
- **Route**: <route>
- **Steps**:
  1. Navigate to <route>
  2. Verify {Component} is visible at {Position}
  3. Verify {Component} renders with expected data
- **Expected**: Component appears at the specified position and displays data correctly
- **Priority**: P0
```

This test case MUST exist for every existing-page integration. It serves as a safety net: if the integration task is skipped, this test will fail.

Skip this section when no `prd/prd-ui-functions.md` exists or no UI Functions have `existing-page` placements.

## Route Validation

Cross-reference each UI test case's `Route` field against actual project route files.

**Discovery patterns** (framework-specific):
- Go (chi/stdlib): `r.Get(`, `mux.HandleFunc`, `http.Handle`
- Express/Node: `router.get(`, `app.get(`, `app.post(`
- React SPA: `path="` or `path='` in route config, `<Route` component
- Next.js: `app/` directory structure, `page.tsx` files
- Spring: `@GetMapping`, `@PostMapping`, `@RequestMapping`

**Validation**: For each test case with a `Route` field:
- Exact or prefix match → annotate `Matched (source:line)`
- No match → annotate `Route not found -- verify path`

If no route files can be discovered, skip this step entirely. Do not fabricate validation results.

## Quality Rules

Apply the 6 Antipattern Prevention rules from the dispatcher's shared rules to every UI test case. Key UI-specific reminders:

- **Pre-conditions must be concrete and creatable**: Specify how to create the required page state (e.g., "a user logged in with admin role" not "admin state exists").
- **Expected results must be specific and verifiable**: State exact text, element visibility, URL, or data value. Not "displays correctly".
- **Steps describe runtime behavior**: Interact with the running application (click, navigate, type), not read source files.

## Output

Write to `docs/features/<slug>/testing/ui-test-cases.md`. Number test cases from TC-001 sequential. End the file with a traceability table:

```markdown
## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | UI | ui/login | P0 |
```
