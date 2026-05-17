---
scale: 1000
target: 900
iterations: 6
type: test-cases
context:
  conventions: [testing-isolation]
  business-rules: auto
---

# Test Cases Evaluation Rubric

**Total: 1000 points**

## Required Sections

The test-cases.md must contain these sections:

- [ ] Frontmatter with `feature`, `sources`, `generated`
- [ ] Summary table (counts per detected interface type + total)
- [ ] Grouped test case sections for each detected interface type (UI/TUI/Mobile/API/CLI Test Cases)
- [ ] Traceability table (TC ID → Source → Type → Target → Priority)
- [ ] Route Validation table (required when profile has `web-ui` or `api` capability AND route files can be discovered; omit without penalty otherwise)

**Missing section**: 0 pts for every dimension that depends on the missing section.

## Dimensions

### 1. PRD Traceability (170 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| TC-to-AC mapping exists | 0-60 | Every TC has a `Source` field pointing to a specific PRD acceptance criterion. Not just "PRD" but "PRD AC-3.1" level specificity |
| Traceability table complete | 0-60 | Traceability table lists every TC with its PRD source, type, target, and priority. No TCs missing from the table |
| Reverse coverage | 0-50 | Every PRD acceptance criterion has at least one TC. No AC is orphaned — check against prd-user-stories.md and prd-spec.md |

### 2. Step Actionability (220 pts)

**Blocking threshold**: If this dimension scores < 180, downstream gen-test-scripts is blocked.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Steps are concrete actions | 0-80 | Each step describes a single, unambiguous user action. "Click the Submit button" not "Submit the form". "GET /api/users?page=1" not "Fetch users" |
| Expected results are verifiable | 0-80 | Every expected result can be objectively verified: specific text, HTTP status, element state, data value. Not "should work" or "displays correctly" |
| Preconditions are explicit | 0-60 | TCs with dependencies (logged-in user, existing data, specific state) declare them in Pre-conditions. No implicit assumptions |

### 3. Interface Accuracy (130 pts)

This dimension adapts based on the project's active test profile capabilities. Read the active test profile manifest (resolved via `forge profile`) to determine the profile's `capabilities` field. These are interface-type capabilities (`web-ui`, `tui`, `mobile-ui`, `api`, `cli`), not build capabilities (`compile`, `test`, `lint`). Then select the matching evaluation criteria:

| Capability | Dimension name | Evaluation focus |
|-----------|---------------|------------------|
| `web-ui` | Route Accuracy | Routes are valid paths matching sitemap.json. UI TCs have Route fields with concrete paths. |
| `tui` | Output Assertion Accuracy | Expected outputs have specific text/snapshot comparison points. Terminal rendering assertions are concrete (exact strings, regex patterns, golden file refs). |
| `mobile-ui` | Interaction Accuracy | Touch/gesture/navigation flows are specific (tap coordinates, swipe directions, screen transitions). Element identification uses accessibility labels or resource IDs. |
| `api` | Contract Accuracy | Request/response structures match actual API schemas. Status codes, headers, body fields are explicit. Error response contracts are covered. |
| `cli` | Command Coverage | Flags, subcommands, arguments are explicitly tested. Output format assertions are concrete (exit codes, stdout/stderr content, error messages). |

**Active capability filtering**: Before dividing points, exclude capabilities that have zero test cases of the matching type in test-cases.md. Only capabilities with at least one corresponding TC participate in scoring.

When multiple capabilities are active, divide 130 pts equally across the remaining (non-excluded) active capabilities. Each capability's sub-criteria use percentage-based weights (not fixed points) to allow clean rescaling:

| Capability-specific criteria | Weight | What to check |
|---|--------|---------------|
| **web-ui**: Routes are valid and specific | 60% | Every Route field contains a real path (e.g., `/users/123/edit`), not vague descriptions. Matches sitemap.json routes where applicable. No placeholder paths |
| **web-ui**: Route consistency | 40% | UI, TUI, and Mobile TCs have Route fields with concrete paths/screen identifiers. API and CLI TCs omit Route — they use Target and command patterns instead. No Route field contains implementation details (testid, selector, CSS) |
| **tui**: Output assertions are concrete | 50% | Expected results specify exact text, snapshot comparison points, or regex patterns for terminal output |
| **tui**: Keyboard interaction coverage | 50% | Keyboard inputs, key sequences, and terminal state transitions are explicitly described |
| **mobile-ui**: Interaction specificity | 50% | Touch targets, gesture types, screen transitions are explicitly described with accessibility labels or resource IDs |
| **mobile-ui**: Navigation flow coverage | 50% | Screen transitions, back navigation, deep links are covered for all navigation paths |
| **api**: Contract accuracy | 50% | Request/response schemas match actual API. Status codes, headers, body fields are explicit |
| **api**: Error contract coverage | 50% | Error responses (4xx, 5xx) are covered with specific error body assertions |
| **cli**: Command coverage | 50% | All flags, subcommands, and argument combinations are tested |
| **cli**: Output assertion specificity | 50% | Exit codes, stdout/stderr content, and error messages are explicitly asserted |

**Scoring example**: If the active capabilities are `api` and `cli` (both with TCs present), each gets 65 pts. For `api`: Contract accuracy scores X/100 × 50% = X × 0.325 of the 65-pt allocation; Error contract coverage scores Y/100 × 50% = Y × 0.325. Total for `api` = X × 0.325 + Y × 0.325 (max 65).

### 4. Completeness (170 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Type coverage | 0-50 | All interface types present in the PRD (UI, TUI, Mobile, API, CLI) have corresponding TCs. If PRD has API endpoints, there are API TCs |
| Boundary and edge cases | 0-60 | Where the PRD explicitly mentions error states, empty states, or boundary conditions, at least one TC covers each. Do not invent scenarios not present in the PRD |
| Integration scenarios | 0-60 | TCs cover cross-feature or cross-interface scenarios (e.g., UI action triggers API call, CLI command affects UI state) where applicable and mentioned in the PRD |

### 5. Structure & ID Integrity (90 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| TC IDs are sequential and unique | 0-35 | IDs follow the pattern (e.g., TC-001, TC-002...). No gaps, no duplicates, no re-used IDs |
| Classification is correct | 0-30 | Each TC is classified under the correct type (UI/TUI/Mobile/API/CLI). No UI TCs in the API section or vice versa |
| Summary table matches actual | 0-25 | Counts in the summary table match the actual number of TCs in each section |

### 6. Antipattern Prevention (90 pts)

This dimension evaluates whether test cases are designed to avoid common downstream antipatterns in `/gen-test-scripts`. Well-designed test cases prevent these issues upstream, making script generation more reliable.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Pre-conditions are concrete and creatable | 0-25 | Score only pre-conditions that exist in the document — missing pre-conditions are penalized under D2 (Step Actionability). For pre-conditions that ARE listed, every one must describe HOW to create the required state (e.g., "a temp project with 3 pending tasks"), not just assert it exists (e.g., "pending tasks exist"). If a pre-condition cannot be created using an isolated test fixture (temp directory, test container, mock server, etc.), the downstream script will generate a conditional skip without fixture. Deduct 10 pts per non-creatable pre-condition |
| Steps describe runtime behavior | 0-20 | No step describes reading source files (`.md`, `.go`, `.json`), checking documentation content, or verifying file existence. Every step interacts with the running product (click, API call, CLI invocation). Deduct 10 pts per static-file-check step |
| No duplicate scenarios | 0-20 | No two TCs test the same scenario with identical inputs and conditions. Duplicate TCs generate duplicate test functions that double CI time. Deduct 10 pts per duplicate pair |
| No meta-testing | 0-15 | No TC verifies test infrastructure ("all tests pass", "test suite compiles", "config is valid"). Every TC must test product behavior. Meta-tests cause recursive test invocation. Deduct 15 pts per meta-test TC |
| Every TC is implementable | 0-10 | No TC describes a scenario that requires unavailable infrastructure (e.g., "real production database") without marking itself as non-implementable. Dead TCs generate unconditional skips — worse than no test at all. Deduct 5 pts per non-implementable TC without annotation |

### 7. Convention Compliance (130 pts)

This dimension uses injected context (project testing conventions) to verify test cases comply with established project patterns. The scorer should reference injected conventions to detect violations.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Test isolation compliance | 0-50 | Do test pre-conditions and setup steps follow the project's declared testing-isolation conventions? If conventions require temp directories, are TCs using them? If conventions require mock servers, do TCs set them up correctly? Deduct 10 pts per convention violation |
| Convention-aware assertions | 0-40 | Do expected results align with project conventions for output format, error format, and response structure? If conventions specify error message format, do TCs assert against that format? Deduct 10 pts per assertion that contradicts project conventions |
| Fixture strategy consistency | 0-40 | Do TCs use fixture strategies consistent with project conventions (e.g., inline fixtures vs factory functions vs seed scripts)? If conventions declare a specific fixture approach, TCs should follow it. Deduct 10 pts per TC using a non-standard fixture approach without justification |

## Deduction Rules

- **Score floor**: No criterion or dimension score may fall below 0. Clamp to 0 after applying all deductions.
- **Missing required section**: 0 pts for that dimension
- **Vague language without specificity**: -20 pts per instance ("click button" without identifying which, "check result" without expected value)
- **Cross-section inconsistency**: -30 pts per conflict (e.g., traceability table says TC-005 is API but it's listed under UI)
- **Placeholder text ("TBD", "TODO", "N/A")**: -20 pts per instance
