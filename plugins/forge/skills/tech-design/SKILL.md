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

1. **Detect language**: Run `forge test detect` to auto-detect the project's test language(s) from file signals.
2. **On failure** (no language detected): ask the user to add `languages` to `.forge/config.yaml` (e.g., `languages: [go]`).

<HARD-RULE>
Do NOT silently default to any language. If `forge test detect` returns no result and the user cannot configure `languages`, abort the skill.
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

### 5.1 PRD Coverage Verification

After drafting each section, verify every PRD acceptance criterion is addressed:

1. For each AC from `prd-user-stories.md`, identify which interface, model, or component handles it
2. If an AC has no corresponding design element, add one
3. Document the mapping in the "PRD Coverage Map" section of the template

### 5.2 Breakdown-Readiness Check

Before seeking approval, verify the design can be directly decomposed into implementation tasks:

| Check | Requirement |
|-------|-------------|
| Components enumerable | Can you list and count all components/modules by name? |
| Interfaces → tasks | Does each interface map to at least one implementation task? |
| Models → tasks | Does each data model map to at least one schema/migration task? |
| PRD AC coverage | Are all acceptance criteria from user stories addressed? |
| Cross-layer consistency | If feature spans layers, does the Data Map cover every field that crosses boundaries? |

If any check fails, add the missing detail before presenting to the user.

### 5.3 Cross-Layer Data Map

If the feature touches more than one architectural layer (database, API, UI, CLI, etc.):
- Complete the "Cross-Layer Data Map" table in the template
- Every field that appears in multiple layers must have a row showing its type/shape at each layer
- This becomes the Ground Truth for type decisions during task execution

If the feature is single-layer (e.g., only affects CLI output formatting):
- Write "Single-layer feature. Cross-Layer Data Map not applicable." in the section

### 5.4 Integration Specs

For each UI Function with `placement: existing-page:<route>`, generate an Integration Spec in the tech design document. Read the UI Design's Placement section for context.

The Integration Spec declares what file to modify and where:
- Do NOT specify implementation details (import statements, prop interfaces)
- Do specify: target file path, insertion point description, data source

This spec is consumed by breakdown-tasks to generate separate integration tasks.

If no UI Function has `placement: existing-page`, write "No existing-page integrations — not applicable."

### 5.5 DB Schema Branch (conditional)

**When `db-schema: "yes"`**:
1. Generate `design/er-diagram.md` using `templates/er-diagram.md` — Mermaid erDiagram + entity detail tables + index design + relationship descriptions
2. Generate `design/schema.sql` using `templates/schema.sql` — CREATE TABLE / ALTER TABLE with inline COMMENT syntax
3. Replace Data Models section in tech-design.md with cross-reference summary + Field Quick Reference table

**When `db-schema: "no"`**:
Data Models stays inline. After drafting, scan content for keywords: `TABLE`, `REFERENCES`, `FOREIGN KEY`, `CREATE TABLE`, `ALTER TABLE`, `migration`, `schema`. If found, prompt: "PRD marked db-schema 'no' but design references database tables. Generate er-diagram.md and schema.sql?" — Yes → proceed with db-schema "yes" path. No → keep inline.

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

Follow the tech-design archiving flow below. Use `templates/decision-entry.md` for the decision row format.

### tech-design Archiving Steps

Triggered after the user approves the tech-design document.

#### Step 7.1 — Check for candidates

Scan the approved tech-design document for entries marked as key decisions. If none exist, skip to Step 7.5.

#### Step 7.2 — Display candidate list

Show the numbered list of key decisions with their type in parentheses:

```
The following decisions are marked as key decisions and recommended for archiving:

  [1] Adopt event-driven architecture (Architecture)
  [2] Use SQLite as local cache storage (Data Model)
  [3] Choose Vitest over Jest as test framework (Dependencies)

Enter numbers to archive (comma-separated), or all / none:
```

#### Step 7.3 — Handle user input

- `none` → skip to Step 7.5
- `all` → archive every candidate
- comma-separated numbers (e.g. `1,3`) → archive only those entries
- `edit:<number>` → enter the edit sub-flow for that entry, then re-display the prompt

Invalid input (number not in candidate list): re-prompt with "Number X is not in the candidate list. Please re-enter."

#### Step 7.4 — Write and update

For each selected entry:
1. Append a decision row to `docs/decisions/<type>.md` (see decision entry row format below).
2. Update `docs/decisions/manifest.md` (see manifest update protocol below).

##### Decision Entry Row Format

Append to the end of `docs/decisions/<type>.md`:

```
| YYYY-MM-DD | <feature-slug> | <Decision, one sentence> | <Rationale, one sentence> | <feature-slug>/design/tech-design.md §<Section> |
```

Field constraints:
- `Date`: ISO 8601 (YYYY-MM-DD)
- `Feature`: feature slug, e.g. `feat-log-decisions`; use `-` if unknown
- `Decision`: single sentence, max 80 characters
- `Rationale`: single sentence, max 80 characters
- `Source`: `<feature-slug>/<file>.md §<Section>` or `manual`

##### Manifest Update Protocol

Target file: `docs/decisions/manifest.md`

**Operation A — Categories table**

Find the row matching the decision type. Increment the `Decisions` count by 1. Set `Last Updated` to today's date (YYYY-MM-DD).

**Operation B — Recent Decisions table**

Insert a new row immediately below the table header (newest first). Keep a maximum of 10 rows; remove the oldest row if the count exceeds 10.

Row format:

```
| YYYY-MM-DD | <feature-slug> | <Type Name> | <Decision, one sentence> | <source> |
```

#### Step 7.5 — Skip logic

If no key decisions exist in the tech-design document, silently skip the archiving step and proceed with the rest of the tech-design flow.

### edit Sub-flow

Triggered when the user inputs `edit:<number>` during the candidate selection prompt.

1. Validate that `<number>` exists in the current candidate list. If not, re-prompt: "Number X is not in the candidate list. Please re-enter."
2. Display the current Decision and Rationale for that entry.
3. Ask: "Enter new Decision (press Enter to keep current):"
4. Ask: "Enter new Rationale (press Enter to keep current):"
5. Update the in-memory candidate entry with the new values.
6. Return to the candidate selection prompt (Step 7.3).

See `examples/ask-question.md` for question formatting and `examples/exploration.md` for context exploration commands.

- If the approved document contains key decisions, display the candidate list and prompt the user to select which to archive.
- User may enter `none` to skip archiving entirely.
- If no key decisions exist in the document, skip this step silently.

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

After committing, use `AskUserQuestion` to ask:

> Run `/eval-design` for adversarial evaluation? (default: 900 points / 3 rounds)

- **Yes** → invoke `/eval-design` via `Skill` tool
- **Custom** → invoke `/eval-design --target X --iterations Y` via `Skill` tool
- **No** → proceed to `/breakdown-tasks`

## Step 11: Auto-Extract Knowledge

After writing design documents and updating the manifest, run the knowledge extraction routine to capture knowledge that the decision archiving in Step 7 may have missed.

### Parameters

| Parameter | Value |
|-----------|-------|
| `trigger` | `tech-design` |
| `artifacts` | `["docs/features/<slug>/design/tech-design.md"]` |

### Artifact Scanning Scope

Focus on architecture decisions, dependency choices, data model decisions in the design document.

### Knowledge Types

The extraction routine identifies four knowledge types:

| Type | Target | Format reference |
|------|--------|-----------------|
| Decision | `docs/decisions/<type>.md` | `decision-logging.md` Section 6 (row format), Section 7 (manifest update) |
| Lesson | `docs/lessons/<slug>.md` | `learn/templates/lesson-entry.md` |
| Convention | `docs/conventions/<topic>.md` | `/consolidate-specs` tech-specs entry format, with project-global ID |
| Business Rule | `docs/business-rules/<domain>.md` | `/consolidate-specs` biz-specs entry format, with project-global ID |

### Extraction Flow

#### Step 1: Scan artifacts

Read all design artifacts specified above.

#### Step 2: Identify notable knowledge

Apply the "notable knowledge" heuristics below to determine if any notable knowledge exists in the scanned artifacts. Classify each candidate by knowledge type (Decision, Lesson, Convention, Business Rule).

#### Step 3: Vocabulary-assisted classification

If `/consolidate-specs` has previously generated vocabulary (from drift-detection runs), use the domain keywords from existing `docs/conventions/` and `docs/business-rules/` files to suggest which target file each extracted item belongs to. This is a suggestion — the agent makes the final classification decision based on content.

If no vocabulary exists (no prior `/consolidate-specs` run), classify unassisted using the domain-to-file mapping from `/consolidate-specs` skill Step 5.

#### Step 4: Silent exit if no notable knowledge

If no candidates pass the "notable" heuristics (below), **produce no output**. Do not ask the user anything. Return silently.

#### Step 5: Present for user confirmation

Use AskUserQuestion to present extracted candidates:

```
Knowledge extracted from tech-design:

  [1] <Decision> → docs/decisions/<type>.md
  [2] <Lesson> → docs/lessons/<slug>.md
  [3] <Convention> → docs/conventions/<topic>.md
  [4] <Business Rule> → docs/business-rules/<domain>.md

Enter numbers to save (comma-separated), or all / none:
```

User input handling:
- `none` → discard all candidates, no output
- `all` → save all candidates
- comma-separated numbers → save only selected candidates

#### Step 6: Write confirmed knowledge

For each confirmed candidate, write to the target file using the format defined by the knowledge type. Create target files if they do not exist. When creating new convention/business-rule files, include YAML frontmatter with `title` and `domains` per `/consolidate-specs` Domain Derivation Rules.

Do NOT write to knowledge directories without explicit user confirmation from Step 5.

### Coordination with Step 7

Step 7 archives **key decisions** explicitly marked in the tech-design document. This step focuses on knowledge types that Step 7 does not cover:

- **Lessons**: Non-obvious insights discovered during design (e.g., constraints that required workarounds)
- **Conventions**: Patterns established in the design that should be repeated across the project
- **Business Rules**: Cross-feature constraints surfaced during technical analysis
- **Decisions**: Only those not already archived by Step 7 (deduplication handles this)

If Step 7 was skipped (no key decisions), this step still runs and may surface notable knowledge from other categories.

### Notable Knowledge Heuristics

The heuristics determine whether a piece of knowledge is "notable" (worth extracting) vs "routine" (skip silently). The goal is a false-positive rate below 30%.

**Decisions — NOT notable when:**

- The choice is the standard/default option in the ecosystem (e.g., "used standard library", "used ORM for database access")
- No meaningful alternatives existed (e.g., "used the only available API")
- The decision is purely cosmetic or stylistic with no architectural impact
- The decision replicates an existing entry in `docs/decisions/`

**Decisions — NOTABLE when:**

- Multiple viable alternatives existed and the choice has lasting impact (e.g., "chose event-driven over polling for state sync")
- The decision involves a non-obvious tradeoff (e.g., "sacrificed consistency for availability in the cache layer")
- A constraint forced an unconventional approach (e.g., "used file-based locking because the Redis dependency was disallowed")

**Lessons — NOT notable when:**

- The root cause is a trivial mistake (e.g., typo, missing import, wrong variable name)
- The issue is standard to the framework/language (e.g., "null pointer from uninitialized field")
- The fix was obvious from the error message
- The lesson replicates an existing entry in `docs/lessons/`

**Lessons — NOTABLE when:**

- The root cause was non-obvious (e.g., race condition from hidden shared state, ordering dependency across services)
- The debugging path was indirect (e.g., "symptom appeared in module A but root cause was in module B")
- The issue would recur in similar contexts and the pattern is worth documenting (e.g., "non-thread-safe map in concurrent handler")

**Conventions — NOT notable when:**

- The pattern is already documented in `docs/conventions/`
- The pattern is a one-off choice specific to this feature
- The pattern is standard practice in the ecosystem (e.g., "used REST for HTTP API")

**Conventions — NOTABLE when:**

- The pattern should be repeated across the project (e.g., "all CLI commands use cobra with this flag structure")
- A project-specific standard was established (e.g., "config files use YAML with this schema structure")
- The pattern emerged from implementation and was not pre-designed

**Business Rules — NOT notable when:**

- The rule is feature-specific logic (e.g., "this feature's form validates email format")
- The rule is a standard CRUD constraint (e.g., "required fields must be non-empty")
- The rule replicates an existing entry in `docs/business-rules/`

**Business Rules — NOTABLE when:**

- The rule applies across features (e.g., "all monetary values use integer cents, never float")
- The rule expresses a domain invariant (e.g., "order status can only advance, never regress")
- The rule constrains user-facing behavior across the system (e.g., "all user actions require authentication except health-check endpoints")

### Deduplication

Before presenting candidates in Step 5, check for duplicates:

1. **Decisions**: grep `docs/decisions/<type>.md` for similar Decision text
2. **Lessons**: grep `docs/lessons/*.md` for similar Root Cause content
3. **Conventions**: grep `docs/conventions/*.md` for similar rule descriptions
4. **Business Rules**: grep `docs/business-rules/*.md` for similar rule statements

If a duplicate is found, exclude the candidate and do not present it. The heuristic goal is: if it is already documented, do not re-extract it.

### Rules

- Extraction logic must be **conservative**: only extract genuinely non-obvious knowledge
- Must not write to knowledge directories without explicit user confirmation
- Silent when no notable knowledge is detected — no output, no prompts
- All output formats must be compatible with `/learn` skill and `/consolidate-specs` overlap detection
- Deduplication runs before presentation — never present a duplicate of existing knowledge

## Integration

Works well with skills:

- `/write-prd` - Creates PRD input and manifest
- `/ui-design` - Preceding skill for UI features; UI design informs technical decisions
- `/eval-design` - Evaluate tech-design.md quality before handing off to breakdown-tasks
- `/breakdown-tasks` - Uses tech-design.md to create tasks
