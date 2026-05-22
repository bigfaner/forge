---
name: gen-test-cases
description: Generate structured test cases from PRD acceptance criteria. Classifies by type (UI/TUI/Mobile/API/CLI) with full traceability to PRD sections.
conventions:
  - testing-isolation.md
---

# Gen Test Cases

Generate structured test cases from PRD acceptance criteria.

**Core principle**: The PRD is the sole input source. Every test case must be traceable to a specific acceptance criterion in the PRD. Do not invent acceptance criteria not present in the PRD.

<HARD-GATE>
This skill only generates test case documents (testing/{type}-test-cases.md), not executable test scripts.
Test script generation is handled by the `/gen-test-scripts` skill.
</HARD-GATE>

## Step 0: Resolve Language

1. Load Convention files from `docs/conventions/` by `domains` frontmatter (match `testing`, `go`, `typescript`, etc.). Extract language from `Framework` section.
2. Fallback: scan existing source/test files (`go.mod`, `package.json`, `*_test.go`, etc.). Also check subdirectories for monorepo.
3. On failure: ask user.

<HARD-RULE>
Do NOT silently default to any language.
</HARD-RULE>

## Prerequisites

| Artifact | Missing prompt |
|----------|----------------|
| `docs/features/<slug>/prd/prd-user-stories.md` | Run `/write-prd` first |
| `docs/features/<slug>/prd/prd-spec.md` | Run `/write-prd` first |
| `docs/sitemap/sitemap.json` (optional, only for `web-ui` interface) | Run `/gen-sitemap`. Skip for non-web-ui interfaces. |

**sitemap-missing fallback**: If sitemap absent with `web-ui` interface, emit warning and proceed without route verification. Do not abort.

This skill can be invoked manually or as the standard task T-test-1 appended by `/breakdown-tasks`.

## When to Use

- User asks to "generate test cases" or "create test cases"
- User provides `/gen-test-cases` command
- After PRD is finalized, before or after implementation

## Process Flow

```
1. Read PRD → 2. Extract AC → 2.5. Detect interfaces → 2.6. Load conventions → 3. Per-type loop → 3.5. Validate routes → 4. Generate manifest
```

### Step 1: Read PRD Sources

Read all available PRD documents:

1. `docs/features/<slug>/prd/prd-user-stories.md` -- primary source for acceptance criteria (Given/When/Then)
2. `docs/features/<slug>/prd/prd-spec.md` -- functional specs, scope, quality checks
3. `docs/features/<slug>/prd/prd-ui-functions.md` -- UI-specific criteria (if exists)
4. `docs/features/<slug>/ui/ui-design.md` -- component-level verification points (if exists)

### Step 2: Extract Acceptance Criteria

Extract every verifiable criterion from each source:

- **User stories**: Each `Given/When/Then` block is one AC. Preserve story reference (e.g., "Story 1 / AC-1").
- **PRD spec**: Quality check items, functional requirements, testable performance/security requirements.
- **UI functions**: Behavior descriptions, interaction requirements, state requirements (loading, empty, error).

<EXTREMELY-IMPORTANT>
Only extract acceptance criteria that **explicitly exist** in the PRD. Forbidden: inventing scenarios not in the PRD, interpreting vague descriptions as specific ACs, omitting explicit Given/When/Then conditions.
</EXTREMELY-IMPORTANT>

### Step 2.5: Detect Project Interfaces

Determine which interface types the project exposes:

1. **Project interfaces** (primary): Examine the project structure and configuration to determine the active interface types:
   - Check `docs/conventions/` for interface type configuration
   - Check project directory structure: `pages/` or `src/components/` → web-ui, `cmd/` with cobra/spf13 imports → cli, route handlers (`api/`, `routes/`) → api, terminal rendering libs (bubbletea, tview) → tui, `android/`/`ios/` → mobile-ui
   - Check `.forge/config.yaml` for a `project-type` field
   - Check `package.json` dependencies: react/vue/next → web-ui, express/fastify → api
   Mapping: `web-ui`->UI, `tui`->TUI, `mobile-ui`->Mobile, `api`->API, `cli`->CLI.
2. **PRD signal** (secondary): A "web application" has web-ui+api; a "CLI tool" has cli.
3. **Codebase signal** (tertiary): Scan for evidence (e.g., `package.json` with react for web-ui, `cmd/` + cobra for cli, route handlers for api, terminal rendering libs for tui, `android/`/`ios/` for mobile-ui).

Record as a set (e.g. `{UI, API}`). TUI clears terminal and redraws (full-screen). CLI produces line-oriented output. Interactive prompts (line-by-line Q&A) are CLI, not TUI.

<HARD-RULE>
If an interface type is absent from the detected set, do not generate test cases for that type. Reclassify to a present type if applicable, or omit.
</HARD-RULE>

### Step 2.6: Load Conventions

Load conventions into context before the per-type loop:

1. **Project-wide**: Read this SKILL.md's `conventions` frontmatter -> for each filename, check `docs/conventions/{filename}` exists -> load via Read tool or skip silently.
2. **Per-type**: For each active type, read `types/{type}.md` frontmatter `conventions` field -> for each filename, check `docs/conventions/{filename}` exists -> load via Read tool or skip silently.

Do not warn or abort on missing convention files.

### Step 3: Per-Type Dispatch Loop

For each active type, sequentially:

1. Read `types/{type}.md` via Read tool.
2. Execute the per-type instructions (type-specific Steps 3-4: classify criteria, generate test cases, write per-type output file to `docs/features/<slug>/testing/{type}-test-cases.md`).

<HARD-RULE>
Every test case must include `Target` (`<type>/<page-or-resource>`) and `Test ID` (`<target>/<title-slug>`, lowercase, hyphens, no punctuation). Number from TC-001 sequential, grouped by type. Every `Source` field must reference a specific PRD location. Each per-type file ends with a traceability table.
</HARD-RULE>

### Step 3.5: Validate Routes

Cross-reference all per-type test case `Route` and `Target` fields against actual project route files. Scan for route registration patterns (framework-specific: Go chi/stdlib, Express, React SPA, Next.js, FastAPI/Flask, Spring, React Native, Flutter). For each test case with a `Route` field: exact/prefix match -> annotate `Matched (source:line)`; no match -> annotate `Route not found -- verify path`.

<HARD-RULE>
If no route files can be discovered, skip this step. Do not fabricate validation results.
</HARD-RULE>

### Step 4: Generate Manifest

After all per-type files are written, generate `docs/features/<slug>/testing/manifest.md`:

```yaml
---
feature: "{{FEATURE_SLUG}}"
types: [{{ACTIVE_TYPES}}]
generated: "{{DATE}}"
---
```

Include a **Summary** table (Type | File | Count | Total) and a **Cross-Type Traceability** table (TC ID | Source | Type | Target | Priority | File). This file is the single entry point for downstream skills to discover all per-type test case files.

## Related Skills

| Skill | Usage |
|-------|-------|
| `/write-prd` | Create PRD with acceptance criteria |
| `/gen-test-scripts` | Generate executable scripts from test cases |
| `/run-e2e-tests` | Execute test scripts and report results |
