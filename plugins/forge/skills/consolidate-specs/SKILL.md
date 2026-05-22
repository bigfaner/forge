---
name: consolidate-specs
description: Extract business rules and tech specs from feature docs into preview files, detect overlaps with existing knowledge, user confirms before integrating to project-level dirs. Also detects and fixes spec drift against the current codebase.
argument-hint: "[--slug <feature-slug>]"
---

# /consolidate-specs

Extract key business rules and technical specifications from feature documents into structured spec files. User reviews and confirms before integration to project-level shared directories. After integration, verifies existing project-level specs still match the current codebase and fixes any drift found.

**Core principle**: Consolidation extracts cross-cutting knowledge that outlives a single feature. Feature-specific details stay in the feature; reusable rules and specs are promoted. Drift detection ensures project-level specs remain accurate as the codebase evolves.

<HARD-GATE>
- Do NOT integrate specs without explicit user confirmation (exception: non-interactive mode auto-integrates per Step 6)
- Do NOT overwrite existing project-level spec files -- append only, unless drift is detected in Step 9
- Do NOT run if `specs/.integrated` marker exists (idempotent)
- Do NOT infer rules not explicitly stated in source documents
</HARD-GATE>

## When to Use

- As T-specs-1 after all e2e tests are promoted (standard task pipeline)
- User invokes `/consolidate-specs` manually at any time

## Prerequisites

Check before running. Abort and prompt user if missing:

| Artifact | Condition | Action if not met |
|----------|-----------|-------------------|
| Feature context | `forge feature` must be set | Run `forge feature <slug>` first |
| `prd/prd-spec.md` | Must exist (unless drift-only mode) | Run `/write-prd` first, or run drift-only mode |
| `design/tech-design.md` | Must exist (unless drift-only mode) | Run `/tech-design` first, or run drift-only mode |

### Drift-Only Mode

If `prd/prd-spec.md` and `design/tech-design.md` do not exist, the skill runs in **drift-only mode**: skip Steps 1-8 and run only Steps 9-11. This enables spec maintenance in quick-mode workflows that lack PRD/design documents.

## Skip Conditions

1. **No extractable rules**: The PRD/design contains no explicit business rules or technical conventions (e.g., pure CRUD with no domain logic). Mark task as completed.
2. **All items are LOCAL**: After extraction, every item is feature-specific with no cross-cutting candidates. Generate preview files but skip integration. Mark task as completed.
3. **Non-interactive session**: Running under `/run-tasks` dispatcher with no user present and CROSS items exist. Auto-integrate all CROSS items without blocking. Commit with `[auto-specs]` tag for traceability (see Step 6 and Step 11 for non-interactive mode details).

## Process Flow

```
Step 1: Check idempotency marker
Step 2: Read feature documents
Step 3: Extract biz-specs from PRD
Step 4: Extract tech-specs from design
Step 5: Generate preview files + detect overlaps
Step 6: Present to user for review -> write review-choices.md
Step 7: Integrate approved specs to project-level dirs
Step 8: Write integration marker + update manifest
Step 9: Detect drift in project-level specs
Step 10: Auto-fix drifted specs
Step 11: Commit spec changes
Step 12: Generate vocabulary index
Step 13: Record task
```

## Step 1: Check Idempotency Marker

```bash
ls docs/features/<slug>/specs/.integrated
```

If the marker exists, this feature's specs have already been integrated. Read the marker to confirm, then skip with status `completed`.

## Step 2: Read Feature Documents

Read all available documents from the current feature:

- `docs/features/<slug>/prd/prd-spec.md` -- business rules, constraints, acceptance criteria
- `docs/features/<slug>/prd/prd-user-stories.md` -- user scenarios with business context
- `docs/features/<slug>/design/tech-design.md` -- interfaces, data models, architecture decisions
- `docs/features/<slug>/design/api-handbook.md` -- API contracts (if exists)
- `docs/features/<slug>/manifest.md` -- feature overview

## Step 3: Extract Business Rules (biz-specs.md)

Scan PRD documents for:

1. **Business constraints** -- validation rules, data integrity requirements, domain invariants
2. **Authorization rules** -- who can do what, role-based access patterns
3. **State transitions** -- allowed state flows, lifecycle rules
4. **Calculation rules** -- formulas, thresholds, business logic
5. **Data ownership** -- which entity owns what data, scoping rules

For each extracted rule:
- Preserve the original context (why the rule exists)
- Note which feature(s) it applies to
- Classify as `[CROSS]` or `[LOCAL]` per `rules/spec-classification.md`

## Step 4: Extract Technical Specifications (tech-specs.md)

Scan design documents for:

1. **Interface contracts** -- API shapes, data structures, serialization rules
2. **Naming conventions** -- field naming, API path conventions, coding standards
3. **Error handling patterns** -- error types, status codes, error propagation
4. **Performance requirements** -- latency targets, throughput, caching rules
5. **Security patterns** -- authentication, authorization, data protection
6. **Data model rules** -- soft-delete patterns, audit fields, indexing conventions

For each extracted spec, apply the same `[CROSS]`/`[LOCAL]` classification from `rules/spec-classification.md`.

## Step 5: Generate Preview Files + Detect Overlaps

Write preview files to `docs/features/<slug>/specs/`:

```
docs/features/<slug>/specs/
  biz-specs.md     -- Extracted business rules
  tech-specs.md    -- Extracted technical specifications
```

Use output formats from `templates/biz-specs.md` and `templates/tech-specs.md`.

Before presenting to user (Step 6), scan for related existing entries per `rules/overlap-detection.md`.

### Early exit

If ALL extracted items are `[LOCAL]`, write the preview files but do NOT proceed to Step 6.

Write `docs/features/<slug>/specs/.integrated` for the early-exit case using the early-exit marker template from `templates/markers.md`.

Then proceed directly to **Step 13** (Record Task). The preview files remain for traceability.

## Step 6: Present to User for Review

### Non-Interactive Mode

When running under `/run-tasks` dispatcher (no user present), skip interactive review and auto-integrate all `[CROSS]` items:

1. Auto-approve all `[CROSS]` items for integration into their detected target files
2. For overlaps with existing entries: use `[skip]` (keep both) by default -- safer than auto-replacing
3. For domain overlap warnings (>50% shared keywords): keep files separate, but include the warning note in the commit message
4. Auto-write `review-choices.md` with all CROSS items listed as "Approved for Integration"
5. Proceed directly to Step 7

This mode is safe because spec errors have no runtime risk and git revert provides perfect rollback. The `[auto-specs]` commit tag (Step 11) enables easy identification and batch revert.

### Interactive Mode

Display the preview with categorization and detected overlaps:

```
Extracted specs for <slug>:

Business Rules:
  [CROSS] BIZ-001: <rule> -> docs/business-rules/<domain>.md (e.g., auth rules -> auth.md)
  [LOCAL] BIZ-002: <rule> -> stays in feature

Technical Specs:
  [CROSS] TECH-001: <spec> -> docs/conventions/<topic>.md (e.g., API patterns -> api.md)
  [LOCAL] TECH-002: <spec> -> stays in feature

Related existing entries (may overlap):
  decisions/error-handling.md: "Adopt AIError struct" -- appears to overlap with TECH-001
  lessons/gotcha-error-handling.md [error-handling] -- may overlap with TECH-002

Domain overlap warnings:
  docs/conventions/error-handling.md (domains: [error, status, response]) and docs/conventions/error-reporting.md (domains: [error, status, log]) share 66% of keywords -- consider merging

For each overlap, choose: [skip] keep both | [replace] delete old + write new
```

Ask the user:
1. Which `[CROSS]` items should be integrated?
2. Which domain/topic file should each be merged into?
3. Any items to skip?
4. For each overlap: `[skip]` keep both, or `[replace]` delete old entry and write new?
5. For each domain overlap warning (>50% shared keywords): merge the files, keep separate, or adjust target file?

Write the user's choices to `docs/features/<slug>/specs/review-choices.md` using the template from `templates/review-choices.md`.

## Step 7: Integrate Approved Specs

For each item listed as "Approved" in `review-choices.md`:

**Business rules** -> append to `docs/business-rules/<domain>.md`:
- Create the file if it doesn't exist
- When creating a new file, write frontmatter per `rules/spec-classification.md` and `rules/domain-frontmatter.md`
- Add a source reference linking back to the feature
- Group by rule category within the file

**Technical specs** -> append to `docs/conventions/<topic>.md`:
- Create the file if it doesn't exist
- When creating a new file, write frontmatter per `rules/spec-classification.md` and `rules/domain-frontmatter.md`
- Add a source reference linking back to the feature
- Group by spec category within the file

Assign project-global IDs per `rules/spec-classification.md`.

### Overlap Resolution

For each user-approved "replace" choice from Step 6:

1. **Decisions**: delete the matching table row from `docs/decisions/<type>.md` -- match by the Decision column text (substring match on the text recorded in review-choices.md)
2. **Lessons**: delete `docs/lessons/<filename>.md`
3. Update `docs/decisions/manifest.md` counts if any decision rows were deleted
4. Write new entry to target file with project-global ID

## Step 8: Write Integration Marker + Update Manifest

Write `docs/features/<slug>/specs/.integrated` using the standard integration marker template from `templates/markers.md`.

Update `docs/features/<slug>/manifest.md` to reference the integrated specs.

## Step 9: Detect Drift in Project-Level Specs

Read all project-level spec files and validate each rule against the current codebase:

1. **Read all spec files**:
   - `docs/business-rules/*.md` -- all business rule files
   - `docs/conventions/*.md` -- all technical convention files

2. **Validate each rule against code**: For each rule in every spec file, search the codebase for the keywords and patterns described in the rule. Compare the rule's stated behavior against the actual code implementation:
   - Extract key domain terms, function names, file paths, or behavior descriptions from each rule
   - Search the relevant source files for those terms
   - Determine whether the rule's description still matches the code's current behavior

3. **Classify each rule**:
   - `current` -- rule description matches current code behavior
   - `drifted` -- rule description is partially or fully inconsistent with current code (e.g., renamed function, changed threshold, modified behavior)
   - `orphaned` -- the code the rule describes no longer exists (e.g., deleted module, removed feature)

4. **Output drift report**: Write a summary of all classifications. If no drift is found (all `current`), skip Steps 10-11 and proceed to Step 12 (vocabulary generation).

### Drift-Only Mode Entry

If running in drift-only mode (no PRD/design files exist), start here at Step 9. Skip Steps 1-8 entirely.

## Step 10: Auto-Fix Drifted Specs

For each rule classified as `drifted` or `orphaned` in Step 9:

1. **Drifted rules**: Update the rule's description/behavior text in-place to match the current code. Preserve the project-global ID (e.g., `BIZ-auth-001`) -- only update the descriptive text, not the ID or structural format.

2. **Orphaned rules**: Remove the rule entry from the spec file. Record the deletion for the commit message in Step 11:
   - Rule ID (e.g., `TECH-api-002`)
   - Reason for deletion (e.g., "corresponding module `X` removed in commit abc1234")

3. **Detect implicit new rules**: While scanning the code for drift, if you discover new patterns, conventions, or business logic that should be documented at the project level but are not in any spec file:
   - Extract the candidate rule with `[CROSS]` classification
   - **Interactive mode**: Present to user for confirmation before appending
   - **Non-interactive mode**: Auto-append with `[auto-specs]` tag -- include in commit message
   - Append confirmed rules to the appropriate spec file with a new project-global ID

4. **Re-derive `domains` frontmatter**: When a file's content changes substantially (rules updated, added, or removed), re-derive the `domains` field per `rules/domain-frontmatter.md`. Compare the new domain set against the existing one:
   - If domains have changed, update the frontmatter in-place
   - If the updated `domains` cause a new >50% overlap with another file's domains, flag in the commit message and notify the user

### Preservation Rules

- Project-global IDs must never change during auto-fix -- only description and behavior text updates
- File structure and formatting must remain consistent with the existing spec file conventions
- Deleted rules must be recorded with ID and reason for traceability in git history

## Step 11: Commit Spec Changes

If any spec files were modified in Step 7 (integration) or Step 10 (drift fix):

1. Stage all changed files under `docs/business-rules/` and `docs/conventions/`
2. Commit with a descriptive message using templates from `templates/commit-messages.md`:
   - Which rule IDs were updated (drift fix)
   - Which rule IDs were removed (orphaned) and why
   - Which rule IDs were added (implicit new rules)
3. **Non-interactive mode**: Include `[auto-specs]` tag in the commit message for traceability. This enables `git log --grep="[auto-specs]"` to find all auto-integrated commits. Include overlap warnings (>50% domain overlap) as notes in the commit message body.

The `[auto-specs]` tag must always be present when running in non-interactive (pipeline) mode, including drift-only mode (Steps 9-11). This ensures all automated spec changes are traceable and reversible.

If no changes were made (all rules `current`), skip this step.

## Step 12: Generate Vocabulary Index

<!-- AUTO-GENERATED -- do not edit manually. Regenerated on every /consolidate-specs run. -->

Scan all four knowledge directories and produce a vocabulary index for use by `/learn` and auto-extract triggers. This step runs unconditionally -- even when knowledge directories are sparse or empty.

### Scan Targets

| Directory | What to extract | Source field |
|-----------|----------------|--------------|
| `docs/decisions/*.md` | Type names from decision row table, domain keywords from Decision column text | Table rows |
| `docs/lessons/*.md` | Tags from YAML frontmatter `tags:` field | `tags` frontmatter |
| `docs/conventions/*.md` | Domains from YAML frontmatter `domains:` field | `domains` frontmatter |
| `docs/business-rules/*.md` | Domains from YAML frontmatter `domains:` field | `domains` frontmatter |

### Base Vocabulary

The base 8-category vocabulary is always included, even when no knowledge files exist:

1. **architecture** -- system structure, layering, modules
2. **interface** -- API contracts, data shapes, serialization
3. **data-model** -- schema, indexing, soft-delete, data ownership
4. **dependencies** -- libraries, versions, packages
5. **error-handling** -- error types, status codes, error propagation
6. **testing** -- test patterns, coverage, mocking
7. **security** -- auth, permissions, data protection
8. **local-dev-deployment** -- dev environment, tooling, deployment

### Aggregation

1. **Types**: Collect unique knowledge types found: `decision`, `lesson`, `convention`, `business-rule`. Report which directories are non-empty vs empty.

2. **Domains**: Aggregate unique domain keywords from all scanned sources (tags from lessons + domains from conventions/business-rules + type-derived keywords from decisions). Merge with the base 8 categories. Deduplicate and normalize (lowercase, sorted).

3. **Counts**: For each type and domain keyword, count how many entries exist across all directories.

### Output Format

Use the output template from `templates/vocabulary-index.md`.

### Idempotency

The vocabulary file is fully regenerated on every `/consolidate-specs` run -- no incremental updates. Previous content is replaced entirely.

### Usage by Other Skills

- `/learn` reads `docs/.vocabulary.md` at runtime to suggest type and domain classifications for user input
- Auto-extract triggers (in `run-tasks`, `fix-bug`, `write-prd`, `tech-design`) read the vocabulary to classify extracted knowledge
- Both `/learn` and triggers accept values outside the vocabulary -- it is suggestive, not restrictive

## Step 13: Record Task

Invoke the skill:

```
Skill(skill="submit-task")
```

Omit `coverage` from record.json -- the noTest flag in index.json auto-sets it.

## Rules

See `rules/constraints.md` for the complete list of constraints and rules governing this skill.

## Related Skills

| Skill | Relationship |
|-------|-------------|
| `/write-prd` | Upstream: source of business rules |
| `/tech-design` | Upstream: source of technical specs |
| `forge test promote` | Predecessor: T-test-4 before this T-specs-1 |
| `/submit-task` | Downstream: records task completion |
