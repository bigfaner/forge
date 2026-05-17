# Test Cases Evaluation Rubric

**Total: 100 points**
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

### 1. PRD Traceability (25 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| TC-to-AC mapping exists | 0-9 | Every TC has a `Source` field pointing to a specific PRD acceptance criterion. Not just "PRD" but "PRD AC-3.1" level specificity |
| Traceability table complete | 0-8 | Traceability table lists every TC with its PRD source, type, target, and priority. No TCs missing from the table |
| Reverse coverage | 0-8 | Every PRD acceptance criterion has at least one TC. No AC is orphaned — check against prd-user-stories.md and prd-spec.md |

### 2. Step Actionability (25 pts)

**Blocking threshold**: If this dimension scores < 20, downstream gen-test-scripts is blocked.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Steps are concrete actions | 0-9 | Each step describes a single, unambiguous user action. "Click the Submit button" not "Submit the form". "GET /api/users?page=1" not "Fetch users" |
| Expected results are verifiable | 0-9 | Every expected result can be objectively verified: specific text, HTTP status, element state, data value. Not "should work" or "displays correctly" |
| Preconditions are explicit | 0-7 | TCs with dependencies (logged-in user, existing data, specific state) declare them in Pre-conditions. No implicit assumptions |

### 3. Route & Element Accuracy (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Routes are valid and specific | 0-7 | Every Route field contains a real path (e.g., `/users/123/edit`), not vague descriptions. Matches sitemap.json routes where applicable |
| Elements are identifiable | 0-7 | Every Element field uses a selector strategy: `data-testid`, `aria-label`, or semantic locator. Not "the button" or "the form" |
| Route/Element consistency | 0-6 | UI TCs have both Route and Element. API TCs have Route but no Element. CLI TCs have neither but have command patterns. No mismatches between TC type and fields |

### 4. Completeness (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Type coverage | 0-7 | All interface types present in the PRD (UI, API, CLI) have corresponding TCs. If PRD has API endpoints, there are API TCs |
| Boundary and edge cases | 0-7 | At least one TC per category covers: empty state, error handling, invalid input, or boundary condition. Not just happy paths |
| Integration scenarios | 0-6 | TCs cover cross-feature or cross-interface scenarios (e.g., UI action triggers API call, CLI command affects UI state) where applicable |

### 5. Structure & ID Integrity (10 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| TC IDs are sequential and unique | 0-4 | IDs follow the pattern (e.g., TC-001, TC-002...). No gaps, no duplicates, no re-used IDs |
| Classification is correct | 0-3 | Each TC is classified under the correct type (UI/API/CLI). No UI TCs in the API section or vice versa |
| Summary table matches actual | 0-3 | Counts in the summary table match the actual number of TCs in each section |

## Deduction Rules

- **Missing required section**: 0 pts for that dimension
- **Vague language without specificity**: -2 pts per instance ("click button" without identifying which, "check result" without expected value)
- **Cross-section inconsistency**: -3 pts per conflict (e.g., traceability table says TC-005 is API but it's listed under UI)
- **Placeholder text ("TBD", "TODO", "N/A")**: -2 pts per instance
