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

## Process Flow

```
0. Detect test language → 1. Read PRD → 2. Explore context → 3. Identify decisions → 4. Ask questions → 5. Draft design → 6. Review → 7. Archive decisions (optional) → 8. Finalize → 11. Auto-extract knowledge
```

## Step 0: Detect Test Language

1. Load Convention files from `docs/conventions/` by `domains` frontmatter (match `testing`, `go`, `typescript`, etc.). Extract language from `Framework` section.
2. Fallback: scan existing source/test files (`go.mod`, `package.json`, `*_test.go`, etc.). Also check subdirectories for monorepo.
3. On failure: ask user.

<HARD-RULE>
Do NOT silently default to any language.
</HARD-RULE>

## Step 1: Read Manifest → PRD

1. Read `manifest.md` to locate documents
2. Read `prd/prd-spec.md`:
   - Understand requirements
   - Note non-functional requirements — these are the **technical constraints** that drive your decisions
   - Identify acceptance criteria
3. Read `prd/prd-user-stories.md` — extract all Given/When/Then acceptance criteria into a checklist
   - Keep this AC list visible throughout the design process — every AC must map to a design element
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

## Step 4: Ask Questions

Use `AskUserQuestion` for ALL uncertain areas.

## Step 5: Draft Design

Present incrementally, section by section:

| Section        | Content                 |
| -------------- | ----------------------- |
| Overview       | High-level approach     |
| Architecture   | Component diagram       |
| Interfaces     | Interface definitions   |
| Data Models    | If `db-schema: "yes"`: generate `er-diagram.md` + `schema.sql`; inline becomes cross-reference. If `db-schema: "no"`: struct definitions as before. |
| Error Handling | Error strategy          |
| Integration Specs | Integration specifications for existing-page components |
| Testing        | Test strategy           |
| Security       | Security considerations |

After drafting each section, apply the quality checks from `rules/design-quality-checks.md`:
- PRD Coverage Verification (5.1)
- Breakdown-Readiness Check (5.2)
- Cross-Layer Data Map (5.3)
- Integration Specs (5.4)
- DB Schema Branch (5.5) — conditional on `db-schema` value

## Step 6: Get Approval

For each section, wait for user approval.

### 6.1 DB Schema Review Gate (when `db-schema: "yes"`)

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

Save to:
- `docs/features/<slug>/design/tech-design.md` — using `templates/tech-design.md`
- `docs/features/<slug>/design/api-handbook.md` — using `templates/api-handbook.md` (if feature has API surface)
- `docs/features/<slug>/design/er-diagram.md` — using `templates/er-diagram.md` (if `db-schema: "yes"`)
- `docs/features/<slug>/design/schema.sql` — using `templates/schema.sql` (if `db-schema: "yes"`)

## Step 9: Update Manifest

Update `manifest.md` using `templates/manifest-update-design.md`:
- Add Tech Design and API Handbook rows to Documents table
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

## Integration

Works well with skills:

- `/write-prd` - Creates PRD input and manifest
- `/ui-design` - Preceding skill for UI features; UI design informs technical decisions
- `/eval-design` - Evaluate tech-design.md quality before handing off to breakdown-tasks
- `/breakdown-tasks` - Uses tech-design.md to create tasks
