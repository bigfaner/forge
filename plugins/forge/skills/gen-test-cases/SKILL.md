---
name: gen-test-cases
description: Generate structured test cases from PRD acceptance criteria. Classifies by type (UI/API/CLI) with full traceability to PRD sections.
---

# Gen Test Cases

Generate structured test cases from PRD acceptance criteria.

**Core principle**: The PRD is the sole input source. Every test case must be traceable to a specific acceptance criterion in the PRD. Do not invent acceptance criteria not present in the PRD.

<HARD-GATE>
This skill only generates test case documents (testing/test-cases.md), not executable test scripts.
Test script generation is handled by the `/gen-test-scripts` skill.
</HARD-GATE>

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

| Artifact | Missing prompt |
|----------|----------------|
| `docs/features/<slug>/prd/prd-user-stories.md` | Run `/write-prd` first |
| `docs/features/<slug>/prd/prd-spec.md` | Run `/write-prd` first |
| `docs/sitemap/sitemap.json` (optional, UI tests only) | Run `/gen-sitemap` — Element field defaults to `sitemap-missing` when absent |

This skill can be invoked manually or as the standard task T-test-1 appended by `/breakdown-tasks`.

## When to Use

**Trigger:**
- User asks to "generate test cases" or "create test cases"
- User provides `/gen-test-cases` command
- After PRD is finalized, before or after implementation

## Workflow

```
1. Read PRD sources → 2. Extract AC → 2.5. Detect interfaces → 3. Classify & generate → 3.5. Validate routes → 4. Write test-cases.md
```

### Step 1: Read PRD Sources

Read all available PRD documents:

1. `docs/features/<slug>/prd/prd-user-stories.md` — primary source for acceptance criteria (Given/When/Then format)
2. `docs/features/<slug>/prd/prd-spec.md` — functional specs, scope, quality checks at the end
3. `docs/features/<slug>/prd/prd-ui-functions.md` — UI-specific criteria (if exists)

Also read `docs/features/<slug>/ui/ui-design.md` if it exists — provides component-level verification points for UI tests.

### Step 2: Extract Acceptance Criteria

From each source, extract every verifiable criterion:

**From user stories** (`prd/prd-user-stories.md`):
- Each `Given/When/Then` block is one acceptance criterion
- Each story may have multiple AC blocks
- Preserve the story reference (e.g., "Story 1 / AC-1")

**From PRD spec** (`prd/prd-spec.md`):
- Quality check items at the end (checkboxes)
- Functional requirements in Functional Specs (prd-ui-functions.md UI specs and Related Changes)
- Performance/security requirements if testable

**From UI functions** (`prd/prd-ui-functions.md`):
- Each UI function's behavior description
- Interaction requirements
- State requirements (loading, empty, error)

<EXTREMELY-IMPORTANT>
Only extract acceptance criteria that **explicitly exist** in the PRD. Forbidden:
- Inventing test scenarios not mentioned in the PRD
- Interpreting vague descriptions as specific acceptance criteria
- Omitting any explicit Given/When/Then conditions
</EXTREMELY-IMPORTANT>

### Step 2.5: Detect Project Interfaces

Before classification, determine which interface types the project actually exposes. Use two signals:

1. **PRD signal** (primary): The PRD read in Step 1 describes the product's nature. A "web application" has UI+API but no CLI; a "CLI tool" has CLI but no UI.
2. **Codebase signal** (secondary): Scan the project for evidence of each interface type. Look broadly — don't hardcode specific framework patterns. Ask: "does this codebase contain code whose *purpose* is to serve this interface?"

**Interface concepts**:

| Interface | Meaning | NOT this |
|-----------|---------|----------|
| **UI** | The project renders pages/views that users interact with in a browser | — |
| **API** | The project exposes HTTP endpoints that clients consume | — |
| **CLI** | The project provides a **user-facing command-line binary** — a product feature the end user invokes from a terminal | Build commands (`go build`, `npm run build`), lint/test tools (`grep`, `eslint`), CI scripts — these are developer tooling, not product interfaces |

**Method**: Based on both signals, decide which interfaces the project exposes. Record as a set (e.g. `{UI, API}`).

<HARD-RULE>
If an interface type is absent from the detected set, **do not generate test cases for that type**. Criteria that would have matched an absent type should be:
1. Reclassified to a present type if they relate to product behavior under that interface
2. Omitted if they are purely build/tooling checks unrelated to any product interface
</HARD-RULE>

### Step 3: Classify & Generate Test Cases

For each extracted criterion, classify by type and generate a test case.

<HARD-RULE>
Every test case must include `Target` and `Test ID` fields:
- **Target**: `<type>/<page-or-resource>` (e.g. `ui/login`, `api/auth`, `cli/deploy`)
- **Test ID**: `<target>/<title-slug>` where title-slug = lowercase title + spaces to hyphens + remove punctuation
</HARD-RULE>

**Type classification rules:**

Only classify into types present in the detected set from Step 2.5. Skip absent types entirely.

| Type | Indicators |
|------|-----------|
| **UI** | Page rendering, navigation, visual state, interactions, responsive behavior, component visibility, form input, modals, tabs, dropdowns |
| **API** | Endpoints, request/response, status codes, data contracts, HTTP methods, authentication headers |
| **CLI** | Commands, flags, output format, exit codes, arguments, stdin/stdout |

**Priority assignment:**
- **P0**: Criteria tied to core user stories or critical path
- **P1**: Criteria tied to secondary features or edge cases in core flow
- **P2**: Nice-to-have verifications, performance checks, edge cases

For each criterion, generate:

```markdown
## TC-{NNN}: {Title}
- **Source**: {Story N / AC-N} or {Spec Section X.Y} or {UI Function Name}
- **Type**: UI | API | CLI
- **Target**: <type>/<page-or-resource>
- **Test ID**: <target>/<title-slug>
- **Pre-conditions**: {What must be true before testing}
- **Route**: {Page route for UI tests}
- **Element**: {Required: sitemap element IDs (e.g., E-001) or `sitemap-missing`}
- **Steps**:
  1. {Step 1}
  2. {Step 2}
- **Expected**: {What the correct result looks like}
- **Priority**: P0 | P1 | P2
```

<HARD-RULE>
**Element field is required** for every test case. There are no exceptions — every TC must have an Element value.

**Sitemap presence check** (during Step 2, after reading PRD sources):
1. Check if `docs/sitemap/sitemap.json` exists.
2. **If sitemap.json does NOT exist**: Set Element to `sitemap-missing` for all test cases. Add a `> **WARNING**: sitemap.json not found — Element set to `sitemap-missing`. Run \`/gen-sitemap\` for precise element references.` note at the top of the test cases output.
3. **If sitemap.json exists**: For each test case, resolve the Route to sitemap page elements. If the sitemap page has element data, use the relevant element IDs. If the sitemap page has no element data for a specific route, set Element to `sitemap-missing` and add a note: `> ⚠️ Route {route} has no element data in sitemap — run \`/gen-sitemap\` to update.`

**Route Validation enhancement** (Step 3.5): When sitemap exists but lacks element data for a test case's route, the Route Validation step must report this gap alongside the route match result. The suggested remediation is running `/gen-sitemap` to regenerate or update the sitemap for that route.
</HARD-RULE>

### Integration Test Case Generation

For each UI Function with `placement: existing-page:<route>`, generate a dedicated integration verification test case:

```markdown
## TC-{NNN}: Integration — {{Component}} visible on {{Page}}
- **Source**: PRD UI Function "{{Function Name}}" Placement + Integration Spec
- **Type**: UI
- **Target**: ui/<page-name>
- **Test ID**: ui/<page-name>/integration-<component-slug>
- **Pre-conditions**: Component build complete, integration task complete
- **Route**: <route>
- **Element**: {{sitemap element IDs for the insertion point area}} (required: use sitemap IDs or `sitemap-missing`)
- **Steps**:
  1. Navigate to <route>
  2. Verify {{Component}} is visible at {{Position}}
  3. Verify {{Component}} renders with expected data
- **Expected**: Component appears at the specified position and displays data correctly
- **Priority**: P0
```

This test case MUST exist for every existing-page integration. It serves as a safety net: if the integration task is skipped, this test will fail.

<HARD-RULE>
**Numbering**: Start from TC-001, sequential. Group by type (UI first, then API, then CLI).

**Traceability**: Every test case's `Source` field must point to a specific location in the PRD. The file must end with a complete traceability table (TC ID → Source → Type → Target → Priority).

**Target derivation rules**:
- UI tests: `ui/<page-name>` (derived from URL or component name)
- API tests: `api/<resource>` (derived from endpoint)
- CLI tests: `cli/<command>` (derived from command name)

**Test ID generation rule**: `<target>/<title-slug>` where title-slug = lowercase title + spaces to hyphens + remove punctuation.
</HARD-RULE>

### Step 3.5: Validate Routes

Cross-reference each test case's `Route` and `Target` fields against actual project route files.

**Discovery**: Scan the project for route registration patterns using Grep (e.g., `r.Get(`, `router.get(`, `app.get(`).

**Validation**: For each test case with a `Route` field:
- Exact or prefix match against discovered routes → `✅ Matched (source:line)`
- No match → `⚠️ Route /path not found — verify path`

Write results as annotations on each test case's `Route` line, and add a summary section after the Traceability table.

<HARD-RULE>
If no route files can be discovered, skip this step entirely and omit the Route Validation section from the output. Do not fabricate validation results.
</HARD-RULE>

### Step 4: Write Output

Fill template at `plugins/forge/skills/gen-test-cases/templates/test-cases.md` and write to `docs/features/<slug>/testing/test-cases.md`. Create the `testing/` directory if it doesn't exist. Overwrite existing `test-cases.md` — the PRD is the source of truth.

## Related Skills

| Skill | Usage |
|-------|-------|
| `/write-prd` | Create PRD with acceptance criteria |
| `/gen-test-scripts` | Generate executable scripts from test cases |
| `/run-e2e-tests` | Execute test scripts and report results |
