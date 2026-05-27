# Knowledge Extraction Rules

Detailed rules for auto-extracting knowledge from PRD artifacts.

## Parameters

| Parameter | Value |
|-----------|-------|
| `trigger` | `write-prd` |
| `artifacts` | PRD content (`docs/features/<slug>/prd/prd-spec.md`, `docs/features/<slug>/prd/prd-user-stories.md`) |

## Artifact Scanning Scope

Focus on new business rules and user-facing constraints in the PRD content.

## Knowledge Types

The extraction routine identifies four knowledge types:

| Type | Target | Format reference |
|------|--------|-----------------|
| Decision | `docs/decisions/<type>.md` | `learn/templates/decision-entry.md` Section 6 (row format), Section 7 (manifest update) |
| Lesson | `docs/lessons/<slug>.md` | `learn/templates/lesson-entry.md` |
| Convention | `docs/conventions/<topic>.md` | `/consolidate-specs` tech-specs entry format, with project-global ID |
| Business Rule | `docs/business-rules/<domain>.md` | `/consolidate-specs` biz-specs entry format, with project-global ID |

## Extraction Flow

### Step 1: Scan artifacts

Read all PRD artifacts specified above.

### Step 2: Identify notable knowledge

Apply the "notable knowledge" heuristics below to determine if any notable knowledge exists in the scanned artifacts. Classify each candidate by knowledge type (Decision, Lesson, Convention, Business Rule).

### Step 3: Vocabulary-assisted classification

If `/consolidate-specs` has previously generated vocabulary (from drift-detection runs), use the domain keywords from existing `docs/conventions/` and `docs/business-rules/` files to suggest which target file each extracted item belongs to. This is a suggestion — the agent makes the final classification decision based on content.

If no vocabulary exists (no prior `/consolidate-specs` run), classify unassisted using the domain-to-file mapping from `/consolidate-specs` skill Step 5.

### Step 4: Silent exit if no notable knowledge

If no candidates pass the "notable" heuristics (below), **produce no output**. Do not ask the user anything. Return silently.

### Step 4.5: Auto-save configuration check

Before presenting candidates, check the auto-save configuration:

```bash
forge config get auto.knowledgeSave
```

Capture stdout (trimmed) and exit code. Then:

| Exit Code | Mode match | Action |
|-----------|-----------|--------|
| 0 | Mode value is `true` | **Skip Step 5 entirely.** Treat all candidates as confirmed, proceed directly to Step 6. |
| 0 | Mode value is `false` | Present Step 5 confirmation (full flow below). |
| Non-zero (config missing/read error) | — | **Fallback: present Step 5 confirmation** (same as `false`). |

Mode context: `quick` when invoked via `/quick` pipeline, `full` when invoked via full pipeline. Parse the config output format `quick:<val> full:<val>` and select the value matching the current mode.

### Step 5: Present for user confirmation (skipped when auto-save is enabled)

This step is **only executed** when the auto-save configuration check (Step 4.5) returns `false` for the current mode.

Use AskUserQuestion to present extracted candidates:

```
Knowledge extracted from write-prd:

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

### Step 6: Write confirmed knowledge

For each confirmed candidate, write to the target file using the format defined by the knowledge type. Create target files if they do not exist. When creating new convention/business-rule files, include YAML frontmatter with `title` and `domains` per `/consolidate-specs` Domain Derivation Rules.

Do NOT write to knowledge directories without explicit user confirmation from Step 5.

## Notable Knowledge Heuristics

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

## Deduplication

Before presenting candidates in Step 5, check for duplicates:

1. **Decisions**: grep `docs/decisions/<type>.md` for similar Decision text
2. **Lessons**: grep `docs/lessons/*.md` for similar Root Cause content
3. **Conventions**: grep `docs/conventions/*.md` for similar rule descriptions
4. **Business Rules**: grep `docs/business-rules/*.md` for similar rule statements

If a duplicate is found, exclude the candidate and do not present it. The heuristic goal is: if it is already documented, do not re-extract it.

## Rules

- Extraction logic must be **conservative**: only extract genuinely non-obvious knowledge
- Must not write to knowledge directories without explicit user confirmation
- Silent when no notable knowledge is detected — no output, no prompts
- All output formats must be compatible with `/learn` skill and `/consolidate-specs` overlap detection
- Deduplication runs before presentation — never present a duplicate of existing knowledge
