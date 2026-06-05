---
name: tech-design
description: Use after PRD (and UI design if applicable) is finalized to create technical design with architecture and implementation details.
effort: high
---

# Tech Design

## Overview

Produce technical design from PRD (and UI design if applicable), making technology decisions informed by the current project state.

**Core principle**: Resolve technical uncertainty during the design phase, avoiding rework during implementation.

<HARD-GATE>
Do NOT write any implementation code until tech-design.md is approved. The output of this skill is a design document, not code.
</HARD-GATE>

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

```bash
ls docs/features/<slug>/prd/prd-spec.md
```

| Artifact | Missing prompt |
|----------|----------------|
| `prd/prd-spec.md` | Run `/write-prd` first, then `/eval-prd`, then `/ui-design` (if UI features) |

## When to Use

**Trigger conditions:**

- Manifest exists at `docs/features/<slug>/manifest.md` with status `prd` or `design`
- PRD Spec exists at `prd/prd-spec.md`
- If feature has UI: `ui/ui-design.md` should exist (run `/ui-design` first)
- PRD is approved and ready for technical design

**Skip when:**

- No PRD exists (use `/write-prd` first)
- Design already exists for the feature

## Intent Detection

Read the `intent` field from `docs/features/<slug>/proposal.md` (or `docs/proposals/<slug>/proposal.md`) frontmatter before starting the process. This determines the design focus via the Pipeline Configuration table below.

### Pipeline Configuration

| Intent | Design Focus | API Handbook | CLI Handbook | Page Map | Screen Map | ER Diagram | User Stories | Security Review |
|--------|-------------|-------------|-------------|----------|------------|------------|-------------|----------------|
| `new-feature` | Full design | Generated | If CLI/TUI surface | If Web surface | If Mobile surface | Conditional on `db-schema` | Read and used | If signal |
| `enhancement` | Internal architecture (improvement to existing behavior) | If signal | If signal | If signal | If signal | If signal | **Skipped** | If signal |
| `refactor` | Internal architecture | If signal | If signal | If signal | If signal | **Skipped** | **Skipped** | If signal |
| `cleanup` | Minimal (typically skipped — cleanup uses Quick mode) | **Skipped** | **Skipped** | **Skipped** | **Skipped** | **Skipped** | **Skipped** | No |
| `fix` | Targeted fix design | If signal | If signal | If signal | If signal | **Skipped** | **Skipped** | No |
| `doc` | Minimal (title + goals + scope) | **Skipped** | **Skipped** | **Skipped** | **Skipped** | **Skipped** | **Skipped** | No |

**Default**: If `intent` is missing or empty, treat as `new-feature` — full design pipeline unchanged.

**Detection**:
```bash
# Read intent from proposal frontmatter
head -20 docs/proposals/<slug>/proposal.md | grep "^intent:"
```

### Override Signals

Pipeline defaults are determined by intent, but PRD content signals can enable additional pipeline steps. During content generation, detect the following signals within the same LLM call:

| Signal Type | Keywords / Patterns | Override Action |
|-------------|-------------------|-----------------|
| API 变更 | "API", "endpoint", "命令重命名", "接口变更", "breaking change" | Enable API Handbook |
| 用户可见行为 | "用户可见", "UI 变更", "CLI 输出", "新选项" | Enable User Stories |
| 安全相关 | "认证", "授权", "权限", "加密", "token" | Enable Security Review |
| 性能相关 | "性能", "延迟", "吞吐量", "缓存" | Enable Performance Baseline |
| 数据迁移 | "迁移", "schema 变更", "数据格式" | Enable Migration Plan |

**Detection rules**:
- Content generation and signal matching happen in parallel inference (same LLM call), not sequentially
- Negation handling: skip signals in negative context (e.g., "不涉及 API 变更"). Relies on LLM context understanding, not keyword matching
- Override only adds steps (开启), never removes. Worst case: unnecessary artifact generated, caught in user review
- Multiple signals trigger independently and stack (e.g., both "API" and "性能" → enable both API Handbook and Performance Baseline)
- When an override triggers, generate a comment in the tech-design output documenting it (e.g., `<!-- Override: API handbook enabled by signal "接口变更" -->`)
- For `doc` intent: override signals are no-op — minimal design format has no pipeline steps that can be overridden

## Process Flow

```
0. Detect test language → 1. Read PRD → 2. Explore context → 3. Identify decisions → 4. Ask questions → 5. Draft design → 6. Review → 7. Archive decisions (optional) → 8. Finalize → 9. Update Manifest → 10. Adversarial Eval Prompt → 11. Auto-extract Knowledge
```

**Intent-specific adjustments**:

| Step | `new-feature` | `enhancement` | `refactor` / `cleanup` / `fix` | `doc` |
|------|--------------|---------------|-------------------------------|-------|
| 1. Read PRD | Full | Simplified | Spec-only | Minimal |
| 3. Decisions focus | All types | Improvement-focused | Architecture-focused | Scope-focused |
| 5. Draft focus | Full design | Internal architecture | Internal architecture | Minimal: title + goals + scope |
| 8. Finalize | All artifacts | Tech-design only | Tech-design only (no API handbook, no ER diagram) | Tech-design only |

## Step 0: Detect Test Language

1. Read `docs/conventions/testing/index.md` to discover available Conventions. Based on the project's language/framework context, select the matching Convention and load it. Extract language from `Framework` section.
2. Fallback: scan existing source/test files (`go.mod`, `package.json`, `*_test.go`, etc.). Also check subdirectories for monorepo.
3. On failure: ask user.

<HARD-RULE>
Do NOT silently default to any language. Do NOT use `domains` frontmatter filtering — use index.md-based discovery.
</HARD-RULE>

## Step 1: Read Manifest → PRD

1. Read `manifest.md` to locate documents
2. Read `prd/prd-spec.md`:
   - Understand requirements
   - Note non-functional requirements — these are the **technical constraints** that drive your decisions
   - Identify acceptance criteria
3. **Intent gate for user stories**:
   - **`new-feature` intent (or missing)**: Read `prd/prd-user-stories.md` — extract all Given/When/Then acceptance criteria into a checklist. Keep this AC list visible throughout the design process — every AC must map to a design element
   - **`enhancement` intent**: **Skip `prd/prd-user-stories.md`** — this file is not generated for enhancement (write-prd uses simplified format). Instead, extract improvement goals from the PRD spec's Goals section. These goals serve as the AC checklist for the design
   - **`refactor` / `cleanup` / `fix` intent**: **Skip `prd/prd-user-stories.md`** — this file is not generated for these intents (write-prd uses spec-only format). Instead, extract acceptance criteria from the PRD spec's "验证标准 (Verification Criteria)" section. These regression criteria serve as the AC checklist for the design
   - **`doc` intent**: **Skip `prd/prd-user-stories.md`** — documentation changes have no user stories. The PRD is minimal (title + goals + scope)
4. Read `prd/prd-spec.md` frontmatter → extract `db-schema` value. Store for conditional branching in Step 5.

> **Note**: The PRD intentionally excludes technology selection (brainstorm and write-prd phases forbid it). All technology decisions start from this phase. Use non-functional constraints from the PRD as input conditions for technology selection.

## Step 2: Explore Context

| Source                 | What to Look For                                  |
| ---------------------- | ------------------------------------------------- |
| `docs/ARCHITECTURE.md` | Layer constraints                                 |
| `docs/decisions/`      | Existing decisions (category-based directory)     |
| `docs/business-rules/` | Cross-feature business rules from prior features  |
| `docs/conventions/`    | Technical conventions from prior features         |
| Package manager files  | Current dependencies (package.json, go.mod, etc.) |
| Source directories     | Existing patterns (src/, internal/, lib/, etc.)   |

## Step 3: Identify Decisions

### new-feature intent (default)

| Decision Type          | Example Questions        |
| ---------------------- | ------------------------ |
| Architecture           | Where does this fit?     |
| Interface              | What interfaces needed?  |
| Data Model             | What structures needed?  |
| Dependencies           | New dependencies?        |
| Error Handling         | How to handle errors?    |
| Testing                | Test strategy?           |
| Security               | Security considerations? |
| Local Dev & Deployment | Dev environment setup?   |

### enhancement intent

Focus on **improvement to existing behavior**. Decision types:

| Decision Type          | Example Questions                                            |
| ---------------------- | ------------------------------------------------------------ |
| Behavior Improvement   | What existing behavior is improved? How is it measured?      |
| Interface              | Does the enhancement change any interfaces?                  |
| Data Model             | Does the enhancement change data structures?                 |
| Dependencies           | New dependencies needed for the improvement?                 |
| Testing                | How to verify the improvement? Regression tests for unchanged behavior? |
| Security               | Does the improvement affect security surface?                |

### refactor / cleanup / fix intent

Focus on **internal architecture** decisions. Prioritize these decision types:

| Decision Type          | Example Questions                                            |
| ---------------------- | ------------------------------------------------------------ |
| Module Reorganization  | Which modules move where? What is the new dependency graph?  |
| Dependency Adjustment  | Which internal dependencies change? Circular deps to break?  |
| Code Structure         | How to restructure for clarity/maintainability?              |
| Behavioral Invariants  | What external behaviors must be preserved?                   |
| Regression Risk        | What could break? How to detect regressions early?           |
| Testing                | Do existing tests cover the refactored code? New tests needed? |

External-facing decisions (API interfaces, new data models, security surface) are typically **not applicable** for refactoring/cleanup/fix — skip unless the refactoring/fix explicitly changes external contracts.

### doc intent

Focus on **documentation scope** decisions. Prioritize:

| Decision Type          | Example Questions                                            |
| ---------------------- | ------------------------------------------------------------ |
| Document Scope         | Which documents change? What is the update purpose?          |
| Content Boundaries     | What information is in/out of scope for the update?          |
| Cross-references       | Which other documents need updating for consistency?         |
| Verification           | How to verify the documentation is accurate and complete?    |

## Step 4: Ask Questions

Use `AskUserQuestion` for ALL uncertain areas.

## Step 5: Draft Design

Present incrementally, section by section. Apply quality checks from `rules/design-quality-checks.md` after each section (see Quality Check Matrix below).

### Section Matrix

| Section | `new-feature` | `enhancement` | `refactor` / `cleanup` / `fix` | `doc` |
|---------|--------------|---------------|-------------------------------|-------|
| Overview | High-level approach | Improvement approach + goals | Refactoring approach + goals | What changes and why |
| Architecture | Component diagram | Before/after improvement diagram | Before/after structure diagram | — |
| Scope | — | — | — | Files to update/create, expected changes |
| Interfaces | Interface definitions | **If** changes interfaces | **If** changes internal interfaces | — |
| Data Models | If `db-schema: "yes"`: generate `er-diagram.md` + `schema.sql`; otherwise struct definitions | **If** changes data structures | **If** changes data structures | — |
| Error Handling | Error strategy | **If** changes error paths | **If** changes error paths | — |
| Integration Specs | Integration specifications | — | — | — |
| Testing | Test strategy | Regression test strategy | Regression test strategy | — |
| Security | Security considerations | **If** touches security-sensitive code | **If** touches security-sensitive code | — |

**Blank cells** mean the section is skipped for that intent. "If" means conditional — include only when applicable.

**Explicitly skipped (not overridable by signals)**:
- `enhancement` / `refactor` / `cleanup` / `fix` / `doc`: **User Stories** — not read (no new user-observable behavior)
- `refactor` / `cleanup` / `fix` / `doc`: **ER Diagram** — not generated (no database schema change)
- `doc`: **Architecture, Interfaces, Data Models, Error Handling, Integration Specs, Testing, Security, API Handbook, ER Diagram** — all not applicable

**Override signal note**: During content generation, detect override signals. If triggered, generate `<!-- Override: ... -->` comments. For `doc` intent, override signals are no-op (no pipeline steps to override).

### Quality Check Matrix

After drafting each applicable section, apply quality checks from `rules/design-quality-checks.md`:

| Check | `new-feature` | `enhancement` | `refactor` / `cleanup` / `fix` | `doc` |
|-------|--------------|---------------|-------------------------------|-------|
| 5.1 PRD Coverage | Map AC from user stories | Map improvement goals from PRD | Map regression criteria from spec | Map documentation goals from PRD |
| 5.2 Breakdown Readiness | All sections | All applicable sections | All applicable sections | All applicable sections |
| 5.3 Cross-Layer Data Map | Yes | If spans layers | If spans layers | — |
| 5.4 Integration Specs | Yes | If override signal | If override signal | — |
| 5.5 DB Schema Branch | Conditional on `db-schema` | If override signal | — | — |

## Step 6: Get Approval

For each section, wait for user approval.

**Intent gate for DB Schema Review**: If `intent` is `enhancement`, `refactor`, `cleanup`, `fix`, or `doc`, **skip Step 6.1 entirely** — no ER diagram or schema is generated for these intents (unless override signal triggers ER Diagram for `enhancement`).

### 6.1 DB Schema Review Gate (when `db-schema: "yes"` and `intent: new-feature`)

<HARD-GATE>
When the Data Models section is reached and `er-diagram.md` + `schema.sql` have been generated, present them as a standalone review unit. Do NOT proceed to remaining sections until the user explicitly approves the database schema.
</HARD-GATE>

Present `er-diagram.md` and `schema.sql` alongside the Data Models cross-reference, and use `AskUserQuestion`:

> Database schema generated. Review the ER diagram and CREATE TABLE statements. Approve the schema?

- **Approved** → proceed to remaining sections
- **Request changes** → revise schema based on feedback, then re-present for approval

## Step 7: Archive Decisions (Optional)

Triggered automatically after user approves the tech-design in Step 6.

Follow the archiving flow in `rules/decision-archiving.md`:
- Scan for key decisions marked in the tech-design document
- Display candidate list for user selection (archive all / specific / none / edit)
- Write decision entries to `docs/decisions/<type>.md` using `templates/decision-entry.md`
- Update `docs/decisions/manifest.md` per the manifest update protocol

If no key decisions exist, silently skip this step.

## Step 8: Write Design Documents

All intents save to `docs/features/<slug>/design/tech-design.md` using `templates/tech-design.md`.

Additional artifacts by intent:

| Artifact | `new-feature` | `enhancement` | `refactor` / `cleanup` / `fix` | `doc` |
|----------|--------------|---------------|-------------------------------|-------|
| `api-handbook.md` | If API surface | If override signal | If override signal | — |
| `cli-handbook.md` | If CLI/TUI surface | If override signal | If override signal | — |
| `page-map.md` | If Web surface | If override signal | If override signal | — |
| `screen-map.md` | If Mobile surface | If override signal | If override signal | — |
| `er-diagram.md` | If `db-schema: "yes"` | If override signal | — | — |
| `schema.sql` | If `db-schema: "yes"` | If override signal | — | — |

### Surface Detection for Handbook Generation

Handbook template selection is based on the project's surface configuration:

```bash
# Detect surfaces from project configuration
forge surfaces
```

Surface-to-template mapping:

| Surface | Template | Output File |
|---------|----------|-------------|
| `api` | `templates/api-handbook.md` | `design/api-handbook.md` |
| `cli` | `templates/cli-handbook.md` | `design/cli-handbook.md` |
| `tui` | `templates/cli-handbook.md` | `design/cli-handbook.md` |
| `web` | `templates/page-map.md` | `design/page-map.md` |
| `mobile` | `templates/screen-map.md` | `design/screen-map.md` |

Each generated handbook includes a `created` frontmatter timestamp for freshness checks by downstream skills (gen-contracts, gen-test-scripts).

## Step 9: Update Manifest

Update `manifest.md` using `templates/manifest-update-design.md`:
- Add Tech Design row to Documents table
- **`new-feature` intent only**: Add API Handbook row to Documents table (if feature has API surface)
- **`new-feature` intent only**: Add CLI Handbook row to Documents table (if feature has CLI/TUI surface)
- **`new-feature` intent only**: Add Page Map row to Documents table (if feature has Web surface)
- **`new-feature` intent only**: Add Screen Map row to Documents table (if feature has Mobile surface)
- **`new-feature` intent only**: Add ER Diagram and Schema rows (if `db-schema: "yes"`)
- **`enhancement` intent**: Add API Handbook row only if override signal triggers it
- **`enhancement` intent**: Add CLI Handbook row only if override signal triggers it
- **`enhancement` intent**: Add Page Map row only if override signal triggers it
- **`enhancement` intent**: Add Screen Map row only if override signal triggers it
- **`refactor` / `cleanup` / `fix` intent**: Add API Handbook row only if override signal triggers it
- **`refactor` / `cleanup` / `fix` intent**: Add CLI Handbook / Page Map / Screen Map rows only if override signal triggers them
- **`doc` intent**: Only Tech Design row (no API Handbook, ER Diagram, or Schema)
- Add traceability links from PRD sections to design sections
- Advance status to `design` if `/ui-design` already completed or if UI is not applicable

## Step 10: Adversarial Eval Prompt

<EXTREMELY-IMPORTANT>
Eval auto-run check — do NOT use AskUserQuestion when config enables auto-run.

Run the following config check sequence via Bash tool:

```bash
# Eval auto-run check (techDesign)
EVAL_ENABLED=$(forge config get auto.eval.techDesign 2>/dev/null)
if [ "$EVAL_ENABLED" = "true" ]; then
  echo "AUTO_RUN"
elif [ "$EVAL_ENABLED" = "false" ]; then
  echo "SKIP"
else
  echo "FALLBACK_ASK"
fi
```

Based on the output:
- **AUTO_RUN** → invoke `/eval-design` via `Skill` tool (default: 900 points / 3 rounds)
- **SKIP** → skip eval, output "eval-design 已通过配置跳过", proceed to `/breakdown-tasks`
- **FALLBACK_ASK** → ask via `AskUserQuestion`: "Run `/eval-design` for adversarial evaluation? (default: 900 points / 3 rounds)"
  - **Yes** → invoke `/eval-design` via `Skill` tool
  - **Custom** → invoke `/eval-design --target X --iterations Y` via `Skill` tool
  - **No** → proceed to `/breakdown-tasks`
</EXTREMELY-IMPORTANT>

## Step 11: Auto-Extract Knowledge

After writing design documents and updating the manifest, run the knowledge extraction routine per `rules/knowledge-extraction.md` to capture knowledge that Step 7 may have missed.

Extraction covers four knowledge types: Decisions, Lessons, Conventions, Business Rules. Only genuinely non-obvious knowledge is extracted (conservative approach). User confirmation is required before writing to any knowledge directory.

