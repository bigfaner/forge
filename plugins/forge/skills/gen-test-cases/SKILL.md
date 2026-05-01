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
| `prd/prd-user-stories.md` | Run `/write-prd` first |
| `prd/prd-spec.md` | Run `/write-prd` first |
| `docs/sitemap/sitemap.json` (optional, UI tests only) | Run `/gen-sitemap` for more precise element references |

**Note**: This skill can be invoked manually or as the standard task T-test-1 appended by `/breakdown-tasks`.

```bash
task feature
ls docs/features/<slug>/prd/prd-user-stories.md
ls docs/features/<slug>/prd/prd-spec.md
```

## When to Use

**Trigger:**
- User asks to "generate test cases" or "create test cases"
- User provides `/gen-test-cases` command
- After PRD is finalized, before or after implementation

**Skip:**
- No PRD exists yet (use `/write-prd` first)

## Workflow

```
1. Read PRD sources → 2. Extract AC → 2.5. Detect interfaces → 3. Classify & generate → 3.5. Validate routes → 4. Write test-cases.md
```

### Step 1: Read PRD Sources

Read all available PRD documents:

1. `prd/prd-user-stories.md` — primary source for acceptance criteria (Given/When/Then format)
2. `prd/prd-spec.md` — functional specs, scope, quality checks at the end
3. `prd/prd-ui-functions.md` — UI-specific criteria (if exists)

Also read `ui/ui-design.md` if it exists — provides component-level verification points for UI tests.

### Step 2: Extract Acceptance Criteria

From each source, extract every verifiable criterion:

**From user stories** (`prd-user-stories.md`):
- Each `Given/When/Then` block is one acceptance criterion
- Each story may have multiple AC blocks
- Preserve the story reference (e.g., "Story 1 / AC-1")

**From PRD spec** (`prd/prd-spec.md`):
- Quality check items at the end (checkboxes)
- Functional requirements in Section 5 (list pages, button operations, forms)
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
- **Target**: <type>/<page-or-resource>          ← e.g. ui/login, api/auth, cli/deploy
- **Test ID**: <target>/<title-slug>            ← e.g. ui/login/login-with-valid-credentials
- **Pre-conditions**: {What must be true before testing}
- **Route**: {Page route for UI tests}            ← e.g. /login, /settings
- **Element**: {Optional: sitemap element IDs}    ← e.g. E-001, L-003 (only if sitemap exists)
- **Steps**:
  1. {Step 1}
  2. {Step 2}
  ...
- **Expected**: {What the correct result looks like}
- **Priority**: P0 | P1 | P2
```

**Element field rules**:
- Only generate when `docs/sitemap/sitemap.json` exists
- Reference element IDs from sitemap (E-NNN for page elements, L-NNN for layout elements)
- List element IDs directly operated on in test steps, comma-separated for multiple
- Omit this field when no sitemap exists; gen-test-scripts will use all page elements

<HARD-RULE>
**Numbering**: Start from TC-001, sequential. Group by type (UI first, then API, then CLI).

**Traceability**: Every test case's `Source` field must point to a specific location in the PRD (Story number, Spec section number, UI Function name). The file must end with a complete traceability table (TC ID → Source → Type → Target → Priority).

**Target derivation rules**:
- UI tests: `ui/<page-name>` (derived from URL or component name, e.g. login page → `ui/login`)
- API tests: `api/<resource>` (derived from endpoint, e.g. `/api/auth` → `api/auth`)
- CLI tests: `cli/<command>` (derived from command name, e.g. `task claim` → `cli/claim`)

**Test ID generation rule**: `<target>/<title-slug>` where title-slug = lowercase title + spaces to hyphens + remove punctuation.
</HARD-RULE>

### Step 3.5: Validate Routes

Cross-reference each test case's `Route` and `Target` fields against actual project route files.

**Discovery**: Scan the project root for route definitions:

| Pattern | Framework |
|---------|-----------|
| `internal/handler/router.go`, `*routes*.go` | Go (Gin/Echo/Fiber) |
| `src/routes/**`, `src/app/**/route.*` | React/Next.js |
| `src/pages/**`, `src/app/**/page.*` | Next.js App Router |
| `config/routes.*`, `routes/**` | Rails/Phoenix |

Use `Grep` to search for route registration patterns (e.g., `r.Get(`, `r.Post(`, `router.get(`, `router.post(`, `app.get(`).

**Validation**: For each test case with a `Route` field:
- Exact or prefix match against discovered routes → `✅ Matched (source:line)`
- No match → `⚠️ Route /path not found — verify path`
- Suggest closest matching route if available

**Write results** as annotations on each test case's `Route` line, and add a summary section after the Traceability table (see template).

<HARD-RULE>
If no route files can be discovered, skip this step entirely and omit the Route Validation section from the output. Do not fabricate validation results.
</HARD-RULE>

### Step 4: Write Output

Read the template at `plugins/forge/skills/gen-test-cases/templates/test-cases.md`.

Fill in:
- Frontmatter with feature slug, source references, generation date
- All test cases
- Traceability table at the end

Write to: `docs/features/<slug>/testing/test-cases.md`

Create the `testing/` directory if it doesn't exist.

## Overwrite Policy

If `testing/test-cases.md` already exists:
- **Overwrite without asking** — this skill regenerates from current PRD state
- The old file is replaced; PRD is the source of truth
- If user wants to preserve, they should commit the previous version first

## Related Skills

| Skill | Usage |
|-------|-------|
| `/write-prd` | Create PRD with acceptance criteria |
| `/gen-test-scripts` | Generate executable scripts from test cases |
| `/run-e2e-tests` | Execute test scripts and report results |
