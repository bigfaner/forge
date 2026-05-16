---
name: gen-test-cases
description: Generate structured test cases from PRD acceptance criteria. Classifies by type (UI/TUI/Mobile/API/CLI) with full traceability to PRD sections.
---

# Gen Test Cases

Generate structured test cases from PRD acceptance criteria.

**Core principle**: The PRD is the sole input source. Every test case must be traceable to a specific acceptance criterion in the PRD. Do not invent acceptance criteria not present in the PRD.

<HARD-GATE>
This skill only generates test case documents (testing/test-cases.md), not executable test scripts.
Test script generation is handled by the `/gen-test-scripts` skill.
</HARD-GATE>

## Step 0: Resolve Profile

1. **Resolve profile**: Run `forge profile` to get the active test profile(s). This reads `.forge/config.yaml`, falls back to project structure detection.
2. **On failure** (output shows `PROFILE: (none)`): ask the user to choose from known profiles (`web-playwright`, `go-test`, `maestro`, `java-junit`, `rust-test`, `pytest`). Run `forge profile set <name>` to persist their choice.
3. **Load profile manifest**: Run `forge profile get <profile-name> --manifest`.

Use the loaded profile manifest for all subsequent steps.

<HARD-RULE>
Do NOT silently default to any profile. If `forge profile` returns no result and the user cannot decide, abort the skill.
</HARD-RULE>

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

| Artifact | Missing prompt |
|----------|----------------|
| `docs/features/<slug>/prd/prd-user-stories.md` | Run `/write-prd` first |
| `docs/features/<slug>/prd/prd-spec.md` | Run `/write-prd` first |
| `docs/sitemap/sitemap.json` (optional, only when profile has `web-ui` capability) | Run `/gen-sitemap`. Skip sitemap check entirely for non-web-ui profiles. |

**sitemap-missing fallback**: When `docs/sitemap/sitemap.json` is absent and the profile has `web-ui` capability, emit a warning in the output (`"sitemap-missing: route validation skipped. Run /gen-sitemap to enable route validation."`) and proceed without route verification. Do not abort the skill.

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

Before classification, determine which interface types the project actually exposes. Use the active test profile's capabilities as the primary signal:

1. **Profile capabilities** (primary): Read the profile manifest resolved in Step 0 → read `capabilities` field. Each capability maps to an interface type:

| Capability | Interface type |
|-----------|---------------|
| `web-ui` | UI (browser DOM interaction) |
| `tui` | TUI (terminal text rendering, keyboard) |
| `mobile-ui` | Mobile UI (touch, gestures) |
| `api` | API (HTTP/network) |
| `cli` | CLI (command-line) |

2. **PRD signal** (secondary): The PRD describes the product's nature. A "web application" has web-ui+api; a "CLI tool" has cli.
3. **Codebase signal** (tertiary): Scan the project for evidence of each interface type:
   - **web-ui**: presence of `package.json` with react/vue/angular dependency, or HTML files with DOM structure
   - **tui**: presence of terminal rendering libraries (tview, bubbletea, ncurses, ratatui)
   - **api**: presence of route handler files (`http.HandleFunc`, `express()`, `@app.route`, `@router.get`)
   - **cli**: presence of `cmd/` directory, `cobra.Command`, `argparse`, `main.go` with `os.Args`
   - **mobile-ui**: presence of `android/` or `ios/` directories, or mobile framework config (Expo, Flutter)

**Interface concepts**:

| Interface | Meaning | NOT this |
|-----------|---------|----------|
| **UI** | The project renders pages/views that users interact with in a browser | — |
| **TUI** | The project renders terminal-based interactive interfaces (text UI, keyboard-driven) | Raw CLI output (flags, exit codes) — TUI has full-screen rendering and keyboard navigation |
| **Mobile** | The project renders screens/views on a mobile device (native or cross-platform) that users interact with via touch and gestures | — |
| **API** | The project exposes HTTP endpoints that clients consume | — |
| **CLI** | The project provides a **user-facing command-line binary** — a product feature the end user invokes from a terminal | Build commands (`go build`, `npm run build`), lint/test tools (`grep`, `eslint`), CI scripts — these are developer tooling, not product interfaces |

**Method**: Based on both signals, decide which interfaces the project exposes. Record as a set (e.g. `{UI, API}`).

**TUI vs CLI disambiguation**: TUI clears the terminal and redraws (full-screen rendering, e.g., `vim`, `htop`, `lazygit`). CLI produces line-oriented sequential output (e.g., `git`, `docker`, `npm`). Interactive prompts (inquirer, cobra — line-by-line Q&A) are CLI, not TUI.

<HARD-RULE>
If an interface type is absent from the detected set, **do not generate test cases for that type**. Criteria that would have matched an absent type should be:
1. Reclassified to a present type if they relate to product behavior under that interface
2. Omitted if they are purely build/tooling checks unrelated to any product interface
</HARD-RULE>

### Step 3: Classify & Generate Test Cases

For each extracted criterion, classify by type and generate a test case.

<HARD-RULE>
Every test case must include `Target` and `Test ID` fields:
- **Target**: `<type>/<page-or-resource>` (e.g. `ui/login`, `tui/dashboard`, `api/auth`, `cli/deploy`)
- **Test ID**: `<target>/<title-slug>` where title-slug = lowercase title + spaces to hyphens + remove punctuation
</HARD-RULE>

**Type classification rules:**

Only classify into types present in the detected set from Step 2.5. Skip absent types entirely.

| Type | Indicators |
|------|-----------|
| **UI** | Page rendering, navigation, visual state, interactions, responsive behavior, component visibility, form input, modals, tabs, dropdowns |
| **TUI** | Terminal screen rendering, keyboard navigation, text output assertions, screen transitions, cursor movement, key bindings, terminal state changes |
| **Mobile** | Touch interactions, gestures (swipe, pinch, long-press), screen transitions, accessibility labels, app lifecycle events, platform-specific UI components |
| **API** | Endpoints, request/response, status codes, data contracts, HTTP methods, authentication headers |
| **CLI** | Commands, flags, output format, exit codes, arguments, stdin/stdout |

**Priority assignment** (decision tree):

1. Is the criterion tied to a Given/When/Then in a user story marked as core/critical in the PRD? → **P0**
2. Is the criterion tied to a secondary story, or an error/boundary case explicitly mentioned for a core story? → **P1**
3. Otherwise (nice-to-have verifications, minor edge cases) → **P2**

If the PRD has no explicit priority marking, default P0 for the first story's ACs and P1 for all others.

For each criterion, generate:

```markdown
## TC-{NNN}: {Title}
- **Source**: {Story N / AC-N} or {Spec Section X.Y} or {UI Function Name}
- **Type**: UI | TUI | Mobile | API | CLI
- **Target**: <type>/<page-or-resource>
- **Test ID**: <target>/<title-slug>
- **Pre-conditions**: {What must be true before testing}
- **Route**: {Page route or screen path — only for UI, TUI, and Mobile tests; omit for API and CLI tests}
- **Steps**:
  1. {Step 1}
  2. {Step 2}
- **Expected**: {What the correct result looks like}
- **Priority**: P0 | P1 | P2
```

Technical implementation details (locators, selectors, testids) are the responsibility of `/gen-test-scripts`, which extracts them directly from source code. Keep test case descriptions in natural language only.

#### Test Case Quality Rules (Antipattern Prevention)

Well-designed test cases prevent downstream antipatterns in `/gen-test-scripts`. Apply these rules to every test case:

| # | Rule | Prevents downstream antipattern | How to apply |
|---|------|--------------------------------|--------------|
| 1 | **Pre-conditions must be concrete and creatable** | Conditional skip without self-contained fixture | Every pre-condition must specify HOW to create the required state (e.g., "a project with 3 pending tasks in temp dir" not "pending tasks exist"). If a pre-condition cannot be created in an isolated environment, rewrite it so it can |
| 2 | **Expected results must be specific and verifiable** | Vacuous assertions | Every expected result must be objectively checkable: exact text, specific status code, element state, data value. Not "works correctly" or "displays as expected" |
| 3 | **Steps describe runtime behavior, not file content** | Static-file text grep tests | Steps must describe interacting with the running product (click button, call API, run command), not reading source files or documentation |
| 4 | **No duplicate scenarios** | Duplicate test functions across packages | Each test case must test a distinct scenario. If two TCs test the same condition with the same inputs, merge them |
| 5 | **Test the product, not the test suite** | Recursive test invocation | Test cases must verify product behavior, not meta-properties like "all tests pass" or "the test suite compiles" |
| 6 | **Every test case must be implementable** | Unconditional t.Skip (dead tests) | If a test case describes a scenario that cannot be implemented without unavailable infrastructure (e.g., requires a physical device), note it after the Priority field — e.g., `Priority: P2 (manual-only: requires physical device)`. Do not leave it as a normal TC that will generate a dead skip |

These rules are derived from common e2e test quality antipatterns observed in automated test generation pipelines.

### Integration Test Case Generation

This section only applies when `prd/prd-ui-functions.md` exists and contains UI Functions with `placement: existing-page:<route>`. Skip if absent.

For each matching UI Function, generate a dedicated integration verification test case:

```markdown
## TC-{NNN}: Integration — {{Component}} visible on {{Page}}
- **Source**: PRD UI Function "{{Function Name}}" Placement + Integration Spec
- **Type**: {UI or Mobile, matching the detected interface type from Step 2.5}
- **Target**: {ui or mobile}/<page-name>
- **Test ID**: {ui or mobile}/<page-name>/integration-<component-slug>
- **Pre-conditions**: Component build complete, integration task complete
- **Route**: <route>
- **Steps**:
  1. Navigate to <route>
  2. Verify {{Component}} is visible at {{Position}}
  3. Verify {{Component}} renders with expected data
- **Expected**: Component appears at the specified position and displays data correctly
- **Priority**: P0
```

This test case MUST exist for every existing-page integration. It serves as a safety net: if the integration task is skipped, this test will fail.

<HARD-RULE>
**Numbering**: Start from TC-001, sequential. Group by type (UI first, then TUI, then Mobile, then API, then CLI).

**Traceability**: Every test case's `Source` field must point to a specific location in the PRD. The file must end with a complete traceability table (TC ID → Source → Type → Target → Priority).

**Target derivation rules**:
- UI tests: `ui/<page-name>` (derived from URL or component name)
- TUI tests: `tui/<screen-name>` (derived from screen/view name)
- Mobile tests: `mobile/<screen-name>` (derived from screen or navigation target)
- API tests: `api/<resource>` (derived from endpoint)
- CLI tests: `cli/<command>` (derived from command name)

**Test ID generation rule**: `<target>/<title-slug>` where title-slug = lowercase title + spaces to hyphens + remove punctuation.
</HARD-RULE>

### Step 3.5: Validate Routes

Cross-reference each test case's `Route` and `Target` fields against actual project route files.

**Discovery**: Scan the project for route registration patterns using Grep. Common patterns by framework:
- Go (chi/stdlib): `r.Get(`, `mux.HandleFunc`, `http.Handle`
- Express/Node: `router.get(`, `app.get(`, `app.post(`
- React SPA: `path="` or `path='` in route config, `<Route` component
- Next.js: `app/` directory structure, `page.tsx` files
- FastAPI/Flask: `@app.get(`, `@router.post(`, `@app.route`
- Spring: `@GetMapping`, `@PostMapping`, `@RequestMapping`
- React Native / Expo: `Stack.Screen`, `navigation.navigate`, screen name strings
- Flutter: `Navigator.push`, `GetPage`, route definitions in `MaterialApp`

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
