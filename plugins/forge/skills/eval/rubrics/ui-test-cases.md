---
scale: 1000
target: 900
iterations: 6
type: ui-test-cases
context:
  conventions: [ux, frontend, testing-isolation]
  business-rules: auto
---

# UI Test Cases Evaluation Rubric

**Total: 1000 points**

## Required Sections

The ui-test-cases.md must contain these sections:

- [ ] Frontmatter with `feature`, `sources`, `generated`
- [ ] UI Test Cases section with individual test cases
- [ ] Traceability table (TC ID → Source → Type → Target → Priority)
- [ ] Route Validation table (required when route files can be discovered; omit without penalty otherwise)

**Missing section**: 0 pts for every dimension that depends on the missing section.

## Dimensions

### 1. PRD Traceability (200 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| TC-to-AC mapping exists | 0-70 | Every TC has a `Source` field pointing to a specific PRD acceptance criterion. Not just "PRD" but "PRD AC-3.1" level specificity |
| Traceability table complete | 0-70 | Traceability table lists every TC with its PRD source, type, target, and priority. No TCs missing from the table |
| Reverse coverage | 0-60 | Every PRD acceptance criterion has at least one TC. No AC is orphaned — check against prd-user-stories.md and prd-spec.md |

### 2. Step Actionability (250 pts)

**Blocking threshold**: If this dimension scores < 200, downstream gen-test-scripts is blocked.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Steps are concrete actions | 0-90 | Each step describes a single, unambiguous user action. "Click the Submit button" not "Submit the form". "GET /api/users?page=1" not "Fetch users" |
| Expected results are verifiable | 0-90 | Every expected result can be objectively verified: specific text, HTTP status, element state, data value. Not "should work" or "displays correctly" |
| Preconditions are explicit | 0-70 | TCs with dependencies (logged-in user, existing data, specific state) declare them in Pre-conditions. No implicit assumptions |

### 3. Visual State Accuracy (150 pts)

This dimension also checks compliance with injected project conventions for UI testing. The scorer should reference injected conventions to detect violations in route usage, element identification, and assertion patterns.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Route Accuracy | 0-60 | Every Route field contains a real path (e.g., `/users/123/edit`), not vague descriptions. Matches sitemap.json routes where applicable. No placeholder paths |
| Route Consistency | 0-40 | UI TCs have Route fields with concrete paths. No Route field contains implementation details (testid, selector, CSS). Route values correspond to actual application routes |
| Convention compliance | 0-50 | Do test steps and assertions comply with project conventions for UI testing? If injected conventions specify element identification patterns (e.g., accessibility labels vs CSS selectors), do TCs use the declared pattern? Are route assertions consistent with the project's routing conventions? Deduct 10 pts per convention violation |

### 4. Completeness (200 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Boundary and edge cases | 0-70 | Where the PRD explicitly mentions error states, empty states, or boundary conditions, at least one TC covers each. Do not invent scenarios not present in the PRD |
| Integration scenarios | 0-70 | TCs cover cross-feature or cross-interface scenarios (e.g., UI action triggers API call) where applicable and mentioned in the PRD |
| UI coverage breadth | 0-60 | All UI-related features described in the PRD have corresponding test cases. Every screen, page, or view mentioned in the PRD is tested |

### 5. Structure & ID Integrity (100 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| TC IDs are sequential and unique | 0-40 | IDs follow the pattern (e.g., TC-001, TC-002...). No gaps, no duplicates, no re-used IDs |
| Classification is correct | 0-30 | Each TC is classified as UI type. No API, CLI, TUI, or Mobile TCs in this file |
| Summary table matches actual | 0-30 | Counts in the summary table match the actual number of TCs in the section |

### 6. Antipattern Prevention (100 pts)

This dimension evaluates whether test cases are designed to avoid common downstream antipatterns in `/gen-test-scripts`. Well-designed test cases prevent these issues upstream, making script generation more reliable.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Pre-conditions are concrete and creatable | 0-30 | Score only pre-conditions that exist in the document — missing pre-conditions are penalized under D2 (Step Actionability). For pre-conditions that ARE listed, every one must describe HOW to create the required state (e.g., "a temp project with 3 pending tasks"), not just assert it exists (e.g., "pending tasks exist"). If a pre-condition cannot be created using an isolated test fixture (temp directory, test container, mock server, etc.), the downstream script will generate a conditional skip without fixture. Deduct 10 pts per non-creatable pre-condition |
| Steps describe runtime behavior | 0-25 | No step describes reading source files (`.md`, `.go`, `.json`), checking documentation content, or verifying file existence. Every step interacts with the running product (click, API call, CLI invocation). Deduct 10 pts per static-file-check step |
| No duplicate scenarios | 0-20 | No two TCs test the same scenario with identical inputs and conditions. Duplicate TCs generate duplicate test functions that double CI time. Deduct 10 pts per duplicate pair |
| No meta-testing | 0-15 | No TC verifies test infrastructure ("all tests pass", "test suite compiles", "config is valid"). Every TC must test product behavior. Meta-tests cause recursive test invocation. Deduct 15 pts per meta-test TC |
| Every TC is implementable | 0-10 | No TC describes a scenario that requires unavailable infrastructure (e.g., "real production database") without marking itself as non-implementable. Dead TCs generate unconditional skips — worse than no test at all. Deduct 5 pts per non-implementable TC without annotation |

## Deduction Rules

- **Score floor**: No criterion or dimension score may fall below 0. Clamp to 0 after applying all deductions.
- **Missing required section**: 0 pts for that dimension
- **Vague language without specificity**: -20 pts per instance ("click button" without identifying which, "check result" without expected value)
- **Cross-section inconsistency**: -30 pts per conflict (e.g., traceability table says TC-005 is API but it's listed under UI)
- **Placeholder text ("TBD", "TODO", "N/A")**: -20 pts per instance
