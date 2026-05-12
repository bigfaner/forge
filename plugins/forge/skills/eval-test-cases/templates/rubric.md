# Test Cases Evaluation Rubric

**Total: 1000 points**
**Report template:** `plugins/forge/skills/eval-test-cases/templates/report.md`

## Required Sections

The test-cases.md must contain these sections:

- [ ] Frontmatter with `feature`, `sources`, `generated`
- [ ] Summary table (UI/API/CLI counts + total)
- [ ] Grouped test case sections (UI Test Cases, API Test Cases, CLI Test Cases)
- [ ] Traceability table (TC ID → Source → Type → Target → Priority)
- [ ] Route Validation table

**Missing section**: 0 pts for every dimension that depends on the missing section.

## Dimensions

### 1. PRD Traceability (250 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| TC-to-AC mapping exists | 0-90 | Every TC has a `Source` field pointing to a specific PRD acceptance criterion. Not just "PRD" but "PRD AC-3.1" level specificity |
| Traceability table complete | 0-80 | Traceability table lists every TC with its PRD source, type, target, and priority. No TCs missing from the table |
| Reverse coverage | 0-80 | Every PRD acceptance criterion has at least one TC. No AC is orphaned — check against prd-user-stories.md and prd-spec.md |

### 2. Step Actionability (250 pts)

**Blocking threshold**: If this dimension scores < 200, downstream gen-test-scripts is blocked.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Steps are concrete actions | 0-90 | Each step describes a single, unambiguous user action. "Click the Submit button" not "Submit the form". "GET /api/users?page=1" not "Fetch users" |
| Expected results are verifiable | 0-90 | Every expected result can be objectively verified: specific text, HTTP status, element state, data value. Not "should work" or "displays correctly" |
| Preconditions are explicit | 0-70 | TCs with dependencies (logged-in user, existing data, specific state) declare them in Pre-conditions. No implicit assumptions |

### 3. Interface Accuracy (200 pts)

This dimension adapts based on the project's active test profile capabilities. Read `.forge/config.yaml` to determine capabilities, then select the matching evaluation criteria:

| Capability | Dimension name | Evaluation focus |
|-----------|---------------|------------------|
| `web-ui` | Route & Element Accuracy | Routes are valid paths matching sitemap.json. Elements use selector strategies (data-testid, aria-label, semantic locators). UI TCs have both Route and Element. |
| `tui` | Output Assertion Accuracy | Expected outputs have specific text/snapshot comparison points. Terminal rendering assertions are concrete (exact strings, regex patterns, golden file refs). |
| `mobile-ui` | Interaction Accuracy | Touch/gesture/navigation flows are specific (tap coordinates, swipe directions, screen transitions). Element identification uses accessibility labels or resource IDs. |
| `api` | Contract Accuracy | Request/response structures match actual API schemas. Status codes, headers, body fields are explicit. Error response contracts are covered. |
| `cli` | Command Coverage | Flags, subcommands, arguments are explicitly tested. Output format assertions are concrete (exit codes, stdout/stderr content, error messages). |

When multiple capabilities are active, evaluate each relevant section proportionally (divide 200 pts equally across active capability dimensions).

| Capability-specific criteria | Points | What to check |
|---|--------|---------------|
| **web-ui**: Routes are valid and specific | 0-70 | Every Route field contains a real path (e.g., `/users/123/edit`), not vague descriptions. Matches sitemap.json routes where applicable |
| **web-ui**: Elements are identifiable | 0-70 | Every Element field uses a selector strategy: `data-testid`, `aria-label`, or semantic locator. Not "the button" or "the form" |
| **web-ui**: Route/Element consistency | 0-60 | UI TCs have both Route and Element. API TCs have Route but no Element. CLI TCs have neither but have command patterns |
| **tui**: Output assertions are concrete | 0-100 | Expected results specify exact text, snapshot comparison points, or regex patterns for terminal output |
| **tui**: Keyboard interaction coverage | 0-100 | Keyboard inputs, key sequences, and terminal state transitions are explicitly described |
| **mobile-ui**: Interaction specificity | 0-100 | Touch targets, gesture types, screen transitions are explicitly described with accessibility labels or resource IDs |
| **mobile-ui**: Navigation flow coverage | 0-100 | Screen transitions, back navigation, deep links are covered for all navigation paths |
| **api**: Contract accuracy | 0-100 | Request/response schemas match actual API. Status codes, headers, body fields are explicit |
| **api**: Error contract coverage | 0-100 | Error responses (4xx, 5xx) are covered with specific error body assertions |
| **cli**: Command coverage | 0-100 | All flags, subcommands, and argument combinations are tested |
| **cli**: Output assertion specificity | 0-100 | Exit codes, stdout/stderr content, and error messages are explicitly asserted |

### 4. Completeness (200 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Type coverage | 0-70 | All interface types present in the PRD (UI, API, CLI) have corresponding TCs. If PRD has API endpoints, there are API TCs |
| Boundary and edge cases | 0-70 | At least one TC per category covers: empty state, error handling, invalid input, or boundary condition. Not just happy paths |
| Integration scenarios | 0-60 | TCs cover cross-feature or cross-interface scenarios (e.g., UI action triggers API call, CLI command affects UI state) where applicable |

### 5. Structure & ID Integrity (100 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| TC IDs are sequential and unique | 0-40 | IDs follow the pattern (e.g., TC-001, TC-002...). No gaps, no duplicates, no re-used IDs |
| Classification is correct | 0-30 | Each TC is classified under the correct type (UI/API/CLI). No UI TCs in the API section or vice versa |
| Summary table matches actual | 0-30 | Counts in the summary table match the actual number of TCs in each section |

## Deduction Rules

- **Missing required section**: 0 pts for that dimension
- **Vague language without specificity**: -20 pts per instance ("click button" without identifying which, "check result" without expected value)
- **Cross-section inconsistency**: -30 pts per conflict (e.g., traceability table says TC-005 is API but it's listed under UI)
- **Placeholder text ("TBD", "TODO", "N/A")**: -20 pts per instance
