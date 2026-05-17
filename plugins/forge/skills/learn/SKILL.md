---
name: learn
description: Unified knowledge accumulation skill. Capture decisions, lessons, conventions, and business rules from a single entry point. Absorbs /record-decision and /learn-lesson.
---

# /learn

Unified knowledge accumulation skill. Captures decisions, lessons, conventions, and business rules from a single entry point.

**Core principle**: Write knowledge immediately, report for review after. No confirmation gate before writing -- the user reviews what was written in the final report and corrects if needed.

## When to Use

**Trigger conditions:**
- User explicitly requests `/learn` (interactive) or `/learn "text"` (direct)
- User asks to record knowledge, a decision, a lesson, a convention, or a business rule
- Mid-task discovery worth preserving

**Delegate to `/consolidate-specs` when:**
- The input describes bulk extraction from feature documents (PRD, tech-design)
- The user wants to extract multiple rules/specs from structured documents
- The input references a feature slug and asks to scan its documents

## Input Modes

### Mode 1: Interactive (`/learn`)

Agent prompts:

```
What did you learn or decide?
```

User responds with free-form text. Agent proceeds to classification.

### Mode 2: Direct (`/learn "text"`)

Skip the prompt. Use the provided text directly. Proceed to classification.

## Workflow

```
1. Identify knowledge type(s) -> 2. Classify -> 3. Write -> 4. Report
```

## Step 1: Identify Knowledge Type(s)

Analyze the input text. Determine which knowledge type(s) apply. A single input can produce entries of multiple types.

| Type | Signals | Output Directory |
|------|---------|-----------------|
| **decision** | "decided to", "chose", "will use", "opted for", architectural choices, dependency selections | `docs/decisions/` |
| **lesson** | "found that", "discovered", "root cause was", "bug caused by", debugging insights, gotchas, patterns | `docs/lessons/` |
| **convention** | "we should always", "standard is", "convention:", coding standards, naming rules, API patterns | `docs/conventions/` |
| **business-rule** | "requirement:", "must", "business rule:", validation rules, authorization rules, state transitions | `docs/business-rules/` |

**Multi-type detection**: If the input contains signals for multiple types (e.g., "found bug caused by X, decided to use Y"), identify all applicable types.

**Bulk detection**: If the input references feature documents or requests scanning multiple sources, stop and suggest:

```
This looks like a bulk extraction from feature documents. Use /consolidate-specs for that.
```

## Step 2: Classify

For each identified type, classify the entry using the shared 8-category vocabulary.

### 8-Category Vocabulary

| Category | Tag | Decision Type File |
|----------|-----|-------------------|
| Architecture | `architecture` | `architecture.md` |
| Interface | `interface` | `interface.md` |
| Data Model | `data-model` | `data-model.md` |
| Dependencies | `dependencies` | `dependencies.md` |
| Error Handling | `error-handling` | `error-handling.md` |
| Testing | `testing` | `testing.md` |
| Security | `security` | `security.md` |
| Local Dev & Deployment | `local-dev-deployment` | `local-dev-deployment.md` |

**Classification behavior:**
- Present the vocabulary as **suggestions**, not enforced values
- Accept custom domain/type values without error
- When auto-generated vocabulary is available (from `/consolidate-specs`), use it to refine suggestions
- When vocabulary is not available, classify based on the 8-category defaults above

### Per-Type Classification

**Decision**: Map to one of the 8 type files in `docs/decisions/`. Select the best-fit type number (1-8).

**Lesson**: Select 1-4 tags from the vocabulary + one category prefix for the filename.

**Convention**: Determine the target topic file in `docs/conventions/` (existing or new). Derive domain keywords from content.

**Business-Rule**: Determine the target domain file in `docs/business-rules/` (existing or new). Derive domain keywords from content.

## Step 3: Write

Write entries immediately, one per identified type. Do not ask for confirmation before writing.

### Decision Entry

1. Read `references/shared/decision-logging.md` for the authoritative format.
2. Read `templates/decision-entry.md` for the row template.
3. Determine: date (today), feature slug (current feature or `-`), decision text, rationale, source (`/learn` or `manual`).
4. If `docs/decisions/` does not exist, auto-create the directory plus all 8 type files and `manifest.md` following decision-logging.md Section 8.
5. Append a decision row to `docs/decisions/<type>.md` (Section 6 row format).
6. Update `docs/decisions/manifest.md` (Section 7 manifest update protocol).

### Lesson Entry

1. Read `templates/lesson-entry.md` for the file template.
2. Generate filename: `<category-prefix><slug>.md` in `docs/lessons/`.
3. Fill the template sections: Problem, Root Cause (trace causal chain at least 3 levels deep), Solution, Reusable Pattern.
4. Set frontmatter: `created` (today's date), `tags` (from vocabulary).
5. Write the file to `docs/lessons/`.

### Convention Entry

1. Read `templates/convention-entry.md` for the entry format.
2. Determine target file: existing `docs/conventions/<topic>.md` or create new.
3. Assign project-global ID: `TECH-<topic>-<NNN>` (find max existing NNN + 1).
4. If creating new file: write frontmatter with `title` and `domains` (3-7 keywords derived from content).
5. Append entry to the target file.

### Business-Rule Entry

1. Read `templates/convention-entry.md` for the entry format.
2. Determine target file: existing `docs/business-rules/<domain>.md` or create new.
3. Assign project-global ID: `BIZ-<domain>-<NNN>` (find max existing NNN + 1).
4. If creating new file: write frontmatter with `title` and `domains` (3-7 keywords derived from content).
5. Append entry to the target file.

## Step 4: Report

After all entries are written, display a final report:

```
Learned knowledge recorded:

[1] Decision -> docs/decisions/<type>.md
    "<Decision text>"
    Rationale: <Rationale text>

[2] Lesson -> docs/lessons/<category-prefix><slug>.md
    "<Title>"
    Tags: <tags>

Review the entries above. Corrections can be made by editing the files directly.
```

For each entry written, include:
- Entry type
- Target file path
- Key content (decision text, lesson title, rule statement)
- Any relevant metadata (tags, rationale)

**The user reviews and corrects after the fact.** No pre-write confirmation.

## Auto-Generated Vocabulary

When available, `/learn` reads the vocabulary index generated by `/consolidate-specs` to refine classification suggestions. When the vocabulary is not available (first run, not yet generated), fall back to the 8-category defaults listed in Step 2.

The vocabulary is suggestive, not required. Custom values are always accepted.

## Directory Bootstrap

If any target directory does not exist:

| Directory | Bootstrap behavior |
|-----------|-------------------|
| `docs/decisions/` | Auto-create all 8 type files + `manifest.md` per decision-logging.md Section 8 |
| `docs/lessons/` | Auto-create the directory |
| `docs/conventions/` | Auto-create the directory |
| `docs/business-rules/` | Auto-create the directory |

## Compatibility

All file formats remain compatible with `/consolidate-specs` overlap detection:
- Decision row format matches decision-logging.md Section 6
- Lesson frontmatter `tags` use the 8-category vocabulary
- Convention/business-rule entries use project-global IDs (`TECH-<topic>-<NNN>`, `BIZ-<domain>-<NNN>`)
- Domain frontmatter follows the derivation rules from `/consolidate-specs`

## Common Mistakes

| Mistake | Correction |
|---------|------------|
| Asking for confirmation before writing | Write immediately, report for review after |
| Rejecting custom domain/type values | Accept any value; vocabulary is suggestive only |
| Writing lesson as "what I did" | Focus on "what to do next time" |
| Stopping at surface symptoms for lessons | Trace the causal chain at least 3 levels deep |
| Writing to wrong knowledge type | Use signal words in Step 1 to classify correctly |
| Inventing domain keywords | Derive from spec content, never fabricate |
