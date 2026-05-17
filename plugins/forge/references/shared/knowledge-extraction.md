# Knowledge Extraction — Shared Auto-Extract Routine

This file defines the shared knowledge extraction logic used by auto-extract trigger points (`run-tasks`, `fix-bug`, `write-prd`, `tech-design`). Trigger points **include** this file — do not copy-paste it.

---

## 1. Parameters

The calling trigger provides these parameters before executing this routine:

| Parameter | Description | Example |
|-----------|-------------|---------|
| `trigger` | Which trigger invoked the extraction | `run-tasks`, `fix-bug`, `write-prd`, `tech-design` |
| `artifacts` | List of artifact paths to scan (varies by trigger — see Section 2) | `["docs/features/x/tasks/", "docs/features/x/manifest.md"]` |

---

## 2. Artifact Scanning Scope per Trigger

| Trigger | Artifacts to scan | Notes |
|---------|-------------------|-------|
| `run-tasks` | Task outcomes (`docs/features/<slug>/tasks/*.md`), code changes (`git diff` against feature branch base), manifest (`docs/features/<slug>/manifest.md`) | Focus on outcomes and patterns that emerged during implementation |
| `fix-bug` | Root cause analysis, fix approach (from the fix session context) | Focus on non-obvious root causes and debugging patterns |
| `write-prd` | PRD content (`docs/features/<slug>/prd/prd-spec.md`, `prd-user-stories.md`) | Focus on new business rules and user-facing constraints |
| `tech-design` | Design document (`docs/features/<slug>/design/tech-design.md`) | Focus on architecture decisions, dependency choices, data model decisions |

---

## 3. Knowledge Types

The extraction routine identifies four knowledge types:

### 3.1 Decisions

Non-obvious choices where alternatives existed. Written to `docs/decisions/<type>.md` using the row format from `decision-logging.md` Section 6.

### 3.2 Lessons

Root causes that were not immediately apparent, debugging patterns worth remembering. Written to `docs/lessons/<slug>.md` using the learn skill's `lesson-entry.md` template.

### 3.3 Conventions

Patterns that should be repeated across the project — coding standards, naming conventions, structural patterns. Appended to `docs/conventions/<topic>.md` (or new file) using the format from `/consolidate-specs` tech-specs entries.

### 3.4 Business Rules

Constraints that apply across features — validation rules, domain invariants, authorization rules. Appended to `docs/business-rules/<domain>.md` (or new file) using the format from `/consolidate-specs` biz-specs entries.

---

## 4. Extraction Flow

### Step 1: Scan artifacts

Read all artifacts specified by the trigger's scanning scope (Section 2). For `run-tasks`, also review `git diff` output for the feature branch to capture code-level patterns.

### Step 2: Identify notable knowledge

Apply the heuristics in Section 5 to determine if any notable knowledge exists in the scanned artifacts. Classify each candidate by knowledge type (Section 3).

### Step 3: Vocabulary-assisted classification

If `/consolidate-specs` has previously generated vocabulary (from drift-detection runs), use the domain keywords from existing `docs/conventions/` and `docs/business-rules/` files to suggest which target file each extracted item belongs to. This is a suggestion — the agent makes the final classification decision based on content.

If no vocabulary exists (no prior `/consolidate-specs` run), classify unassisted using the domain-to-file mapping from `/consolidate-specs` skill Step 5.

### Step 4: Silent exit if no notable knowledge

If no candidates pass the "notable" heuristics (Section 5), **produce no output**. Do not ask the user anything. Return silently to the calling trigger.

### Step 5: Present for user confirmation

Use AskUserQuestion to present extracted candidates:

```
Knowledge extracted from <trigger>:

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

For each confirmed candidate, write to the target file using the format defined by the knowledge type:

| Type | Target | Format reference |
|------|--------|-----------------|
| Decision | `docs/decisions/<type>.md` | `decision-logging.md` Section 6 (row format), Section 7 (manifest update) |
| Lesson | `docs/lessons/<slug>.md` | `learn/templates/lesson-entry.md` |
| Convention | `docs/conventions/<topic>.md` | `/consolidate-specs` tech-specs entry format, with project-global ID |
| Business Rule | `docs/business-rules/<domain>.md` | `/consolidate-specs` biz-specs entry format, with project-global ID |

Create target files if they do not exist. When creating new convention/business-rule files, include YAML frontmatter with `title` and `domains` per `/consolidate-specs` Domain Derivation Rules.

Do NOT write to knowledge directories without explicit user confirmation from Step 5.

---

## 5. Notable Knowledge Heuristics

These heuristics determine whether a piece of knowledge is "notable" (worth extracting) vs "routine" (skip silently). The goal is a false-positive rate below 30%.

### 5.1 Decisions — NOT notable when:

- The choice is the standard/default option in the ecosystem (e.g., "used standard library", "used ORM for database access")
- No meaningful alternatives existed (e.g., "used the only available API")
- The decision is purely cosmetic or stylistic with no architectural impact
- The decision replicates an existing entry in `docs/decisions/`

### 5.2 Decisions — NOTABLE when:

- Multiple viable alternatives existed and the choice has lasting impact (e.g., "chose event-driven over polling for state sync")
- The decision involves a non-obvious tradeoff (e.g., "sacrificed consistency for availability in the cache layer")
- A constraint forced an unconventional approach (e.g., "used file-based locking because the Redis dependency was disallowed")

### 5.3 Lessons — NOT notable when:

- The root cause is a trivial mistake (e.g., typo, missing import, wrong variable name)
- The issue is standard to the framework/language (e.g., "null pointer from uninitialized field")
- The fix was obvious from the error message
- The lesson replicates an existing entry in `docs/lessons/`

### 5.4 Lessons — NOTABLE when:

- The root cause was non-obvious (e.g., race condition from hidden shared state, ordering dependency across services)
- The debugging path was indirect (e.g., "symptom appeared in module A but root cause was in module B")
- The issue would recur in similar contexts and the pattern is worth documenting (e.g., "non-thread-safe map in concurrent handler")

### 5.5 Conventions — NOT notable when:

- The pattern is already documented in `docs/conventions/`
- The pattern is a one-off choice specific to this feature
- The pattern is standard practice in the ecosystem (e.g., "used REST for HTTP API")

### 5.6 Conventions — NOTABLE when:

- The pattern should be repeated across the project (e.g., "all CLI commands use cobra with this flag structure")
- A project-specific standard was established (e.g., "config files use YAML with this schema structure")
- The pattern emerged from implementation and was not pre-designed

### 5.7 Business Rules — NOT notable when:

- The rule is feature-specific logic (e.g., "this feature's form validates email format")
- The rule is a standard CRUD constraint (e.g., "required fields must be non-empty")
- The rule replicates an existing entry in `docs/business-rules/`

### 5.8 Business Rules — NOTABLE when:

- The rule applies across features (e.g., "all monetary values use integer cents, never float")
- The rule expresses a domain invariant (e.g., "order status can only advance, never regress")
- The rule constrains user-facing behavior across the system (e.g., "all user actions require authentication except health-check endpoints")

---

## 6. Deduplication

Before presenting candidates in Step 5, check for duplicates:

1. **Decisions**: grep `docs/decisions/<type>.md` for similar Decision text
2. **Lessons**: grep `docs/lessons/*.md` for similar Root Cause content
3. **Conventions**: grep `docs/conventions/*.md` for similar rule descriptions
4. **Business Rules**: grep `docs/business-rules/*.md` for similar rule statements

If a duplicate is found, exclude the candidate and do not present it. The heuristic goal is: if it is already documented, do not re-extract it.

---

## 7. Calling Convention

Trigger points include this routine by reading this file and following the steps. The trigger provides `trigger` and `artifacts` parameters, then the extraction flow runs from Step 1 through Step 6 (or exits silently at Step 4).

Example inclusion in a trigger point's skill/command file:

```markdown
## Auto-Extract Knowledge

After completing the main workflow, read this extraction routine (`knowledge-extraction.md`)
and execute its extraction flow with:
- trigger: <trigger-name>
- artifacts: <list of artifact paths to scan>
```

---

## 8. Rules

- Extraction logic must be **conservative**: only extract genuinely non-obvious knowledge
- Must not write to knowledge directories without explicit user confirmation
- Silent when no notable knowledge is detected — no output, no prompts
- All output formats must be compatible with `/learn` skill and `/consolidate-specs` overlap detection
- This is a shared reference — trigger points include it, not copy-paste it
- Deduplication runs before presentation — never present a duplicate of existing knowledge
