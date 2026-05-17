---
type: mobile
conventions:
  - testing-mobile.md
---

# Mobile Test Case Generation Instructions

Type-specific Steps 3-4 for **Mobile** (touch, gestures, screen transitions) test cases. Loaded by the dispatcher after Step 2.5 interface detection.

## Classification Indicators

Classify a PRD criterion as **Mobile** when it involves any of:

- Touch interactions (tap, double-tap, long-press)
- Gestures (swipe, pinch, drag, rotate)
- Screen transitions and navigation flows
- Accessibility labels and resource IDs
- App lifecycle events (background, foreground, terminate)
- Platform-specific UI components (bottom sheets, native dialogs, permissions)
- Push notifications, deep links

## Target Derivation

- **Target format**: `mobile/<screen-name>`
- Derive `<screen-name>` from the screen or navigation target name (e.g., `mobile/login`, `mobile/settings`, `mobile/product-detail`)

## Test ID Format

- **Test ID**: `<target>/<title-slug>`
- `title-slug` = lowercase title, spaces to hyphens, remove punctuation
- Example: `mobile/login/valid-credentials-navigate-to-home`

## Priority Assignment

1. Criterion tied to a core/critical Given/When/Then in the PRD → **P0**
2. Criterion tied to a secondary story, or an explicit error/boundary case for a core story → **P1**
3. Nice-to-have verifications, minor edge cases → **P2**

If the PRD has no explicit priority marking, default P0 for the first story's ACs and P1 for all others.

## TC Format

```markdown
## TC-{NNN}: {Title}
- **Source**: {Story N / AC-N} or {Spec Section X.Y}
- **Type**: Mobile
- **Target**: mobile/<screen-name>
- **Test ID**: mobile/<screen-name>/<title-slug>
- **Pre-conditions**: {What must be true before testing}
- **Route**: {Screen identifier or navigation path — required for Mobile tests}
- **Steps**:
  1. {Step 1}
  2. {Step 2}
- **Expected**: {What the correct result looks like}
- **Priority**: P0 | P1 | P2
```

- `Route` field is required for Mobile test cases — must contain a concrete screen identifier or navigation path.
- Touch targets and gesture types must be explicitly described in Steps (e.g., "Tap the Submit button", "Swipe left on the card list").
- Element identification should reference accessibility labels or resource IDs in Expected results where the PRD specifies them.

## Integration TC Generation (Existing-Page Placements)

When `prd/prd-ui-functions.md` exists and contains UI Functions with `placement: existing-page:<route>`, generate a dedicated integration test case for each:

```markdown
## TC-{NNN}: Integration — {Component} visible on {Screen}
- **Source**: PRD UI Function "{Function Name}" Placement + Integration Spec
- **Type**: Mobile
- **Target**: mobile/<screen-name>
- **Test ID**: mobile/<screen-name>/integration-<component-slug>
- **Pre-conditions**: Component build complete, integration task complete
- **Route**: <route>
- **Steps**:
  1. Navigate to <route>
  2. Verify {Component} is visible at {Position}
  3. Verify {Component} renders with expected data
- **Expected**: Component appears at the specified position and displays data correctly
- **Priority**: P0
```

This test case MUST exist for every existing-page integration. Skip when no `prd/prd-ui-functions.md` exists or no UI Functions have `existing-page` placements.

## Interaction Specificity

Mobile test cases require concrete interaction descriptions. Each step involving user input must specify:

- **Touch action**: tap, double-tap, long-press, swipe (direction), pinch (in/out), drag (source → target)
- **Target element**: identified by accessibility label, resource ID, or on-screen position
- **Screen context**: which screen the interaction occurs on

Vague descriptions like "interact with the list" or "navigate away" are not acceptable.

## Navigation Flow Coverage

Cover all navigation paths mentioned in the PRD:

- Forward navigation (screen A → screen B)
- Back navigation (system back button, up button, gesture)
- Deep links (opening the app at a specific screen)
- Tab/bar switching
- Modal presentations and dismissals

## Route Validation

Cross-reference each Mobile test case's `Route` field against actual screen/route definitions.

**Discovery patterns** (framework-specific):
- React Native / Expo: `Stack.Screen`, `navigation.navigate`, screen name strings
- Flutter: `Navigator.push`, `GetPage`, route definitions in `MaterialApp`
- iOS (SwiftUI): `NavigationLink`, `NavigationStack`, view identifiers
- Android (Jetpack Compose): `NavHost`, `composable()`, route strings

**Validation**: For each test case with a `Route` field:
- Match against discovered screen definitions → annotate `Matched (source:line)`
- No match → annotate `Route not found -- verify path`

If no route/screen definitions can be discovered, skip this step entirely. Do not fabricate validation results.

## Quality Rules

Apply the 6 Antipattern Prevention rules from the dispatcher's shared rules to every Mobile test case. Key Mobile-specific reminders:

- **Pre-conditions must be concrete and creatable**: Specify app state requirements (e.g., "user logged in", "app launched with test account").
- **Expected results must be specific and verifiable**: State exact screen content, element visibility, or navigation outcome. Not "works as expected".
- **Steps describe runtime behavior**: Interact with the running app (tap, swipe, navigate), not read source files or inspect layouts.

## Output

Write to `docs/features/<slug>/testing/mobile-test-cases.md`. Number test cases from TC-001 sequential. End the file with a traceability table:

```markdown
## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | Mobile | mobile/login | P0 |
```
