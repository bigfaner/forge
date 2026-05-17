---
name: consolidate-specs
description: Extract business rules and tech specs from feature docs into preview files, detect overlaps with existing knowledge, user confirms before integrating to project-level dirs. Also detects and fixes spec drift against the current codebase.
---

# /consolidate-specs

Extract key business rules and technical specifications from feature documents into structured spec files. User reviews and confirms before integration to project-level shared directories. After integration, verifies existing project-level specs still match the current codebase and fixes any drift found.

**Core principle**: Consolidation extracts cross-cutting knowledge that outlives a single feature. Feature-specific details stay in the feature; reusable rules and specs are promoted. Drift detection ensures project-level specs remain accurate as the codebase evolves.

<HARD-GATE>
- Do NOT integrate specs without explicit user confirmation
- Do NOT overwrite existing project-level spec files — append only, unless drift is detected in Step 9
- Do NOT run if `specs/.integrated` marker exists (idempotent)
- Do NOT infer rules not explicitly stated in source documents
</HARD-GATE>

## When to Use

- As T-specs-1 after all e2e tests graduate (standard task pipeline)
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
3. **Non-interactive session**: Running under `/run-tasks` dispatcher with no user present and CROSS items exist. Write preview files, mark task as `blocked`, and note "User review required for integration."

## Workflow

```
Step 1: Check idempotency marker
Step 2: Read feature documents
Step 3: Extract biz-specs from PRD
Step 4: Extract tech-specs from design
Step 5: Generate preview files + detect overlaps
Step 6: Present to user for review → write review-choices.md
Step 7: Integrate approved specs to project-level dirs
Step 8: Write integration marker + update manifest
Step 9: Detect drift in project-level specs
Step 10: Auto-fix drifted specs
Step 11: Commit spec changes
Step 12: Record task
```

## Step 1: Check Idempotency Marker

```bash
ls docs/features/<slug>/specs/.integrated
```

If the marker exists, this feature's specs have already been integrated. Read the marker to confirm, then skip with status `completed`.

## Step 2: Read Feature Documents

Read all available documents from the current feature:

- `docs/features/<slug>/prd/prd-spec.md` — business rules, constraints, acceptance criteria
- `docs/features/<slug>/prd/prd-user-stories.md` — user scenarios with business context
- `docs/features/<slug>/design/tech-design.md` — interfaces, data models, architecture decisions
- `docs/features/<slug>/design/api-handbook.md` — API contracts (if exists)
- `docs/features/<slug>/manifest.md` — feature overview

## Step 3: Extract Business Rules (biz-specs.md)

Scan PRD documents for:

1. **Business constraints** — validation rules, data integrity requirements, domain invariants
2. **Authorization rules** — who can do what, role-based access patterns
3. **State transitions** — allowed state flows, lifecycle rules
4. **Calculation rules** — formulas, thresholds, business logic
5. **Data ownership** — which entity owns what data, scoping rules

For each extracted rule:
- Preserve the original context (why the rule exists)
- Note which feature(s) it applies to
- Classify as `[CROSS]` or `[LOCAL]` using these criteria:
  - `[CROSS]`: Referenced by 2+ features, or expresses a domain invariant (not feature behavior), or establishes a naming/error-handling convention
  - `[LOCAL]`: Only meaningful within this feature's scope

## Step 4: Extract Technical Specifications (tech-specs.md)

Scan design documents for:

1. **Interface contracts** — API shapes, data structures, serialization rules
2. **Naming conventions** — field naming, API path conventions, coding standards
3. **Error handling patterns** — error types, status codes, error propagation
4. **Performance requirements** — latency targets, throughput, caching rules
5. **Security patterns** — authentication, authorization, data protection
6. **Data model rules** — soft-delete patterns, audit fields, indexing conventions

For each extracted spec, apply the same `[CROSS]`/`[LOCAL]` classification from Step 3.

## Domain Frontmatter

Convention and business-rule files carry a `domains` field in their YAML frontmatter that enables lightweight discovery by consumers (prompt templates, commands, agents).

```yaml
---
title: "Error Handling Conventions"  # existing field, unchanged
domains: [error, status, response, stderr]  # keywords this file covers
---
```

### Domain Derivation Rules

Domains are **derived programmatically from spec content** — never invented by the agent. The derivation algorithm:

1. **Spec ID keywords**: Extract tokens from project-global IDs in the file (e.g., `BIZ-auth-001` contributes `auth`, `TECH-api-003` contributes `api`)
2. **Source keywords**: Extract recurring domain-specific nouns from rule titles, requirement statements, and source references (e.g., a rule about "token validation" contributes `token`, `validation`)
3. **Deduplicate and normalize**: Lowercase, remove duplicates, keep only specific terms (not generic words like "rule", "spec", "requirement")
4. **Cardinality**: Each file gets **3-7 specific keywords**

### Domain Overlap Detection

When multiple files in the same directory (`docs/conventions/` or `docs/business-rules/`) have `domains` fields, compute keyword overlap:

- **Overlap ratio** = `|intersection(domains_A, domains_B)| / min(|domains_A|, |domains_B|)`
- **Threshold**: If overlap ratio > 50%, flag as a potential duplicate/merge candidate during the user confirmation step (Step 6)
- **Action**: Display the warning; the user decides whether to merge or keep separate

## Step 5: Generate Preview Files + Detect Overlaps

Write preview files to `docs/features/<slug>/specs/`:

```
docs/features/<slug>/specs/
  biz-specs.md     — Extracted business rules
  tech-specs.md    — Extracted technical specifications
```

### Preview ID Numbering

Feature-local IDs in preview files use sequential 3-digit numbering starting at 001, independent per file:
- `biz-specs.md`: BIZ-001, BIZ-002, ...
- `tech-specs.md`: TECH-001, TECH-002, ...

### biz-specs.md Format

```markdown
---
feature: "<slug>"
generated: "<date>"
status: draft
---

# Business Rules: <Feature Name>

## <Rule Category>

### BIZ-NNN: <Rule Title>

**Rule**: <concise rule statement>
**Context**: <why this rule exists>
**Scope**: [CROSS|LOCAL]
**Source**: <prd section reference>

<Additional details, examples, or edge cases>
```

### tech-specs.md Format

```markdown
---
feature: "<slug>"
generated: "<date>"
status: draft
---

# Technical Specifications: <Feature Name>

## <Spec Category>

### TECH-NNN: <Spec Title>

**Requirement**: <concise requirement>
**Scope**: [CROSS|LOCAL]
**Source**: <design section reference>

<Implementation details, examples>
```

### Overlap Detection

Before presenting to user (Step 6), scan for related existing entries. This applies to BOTH business rules (biz-specs) and technical specs (tech-specs):

1. **Decisions**: match by filename → for each extracted entry's domain (biz or tech), map it to the corresponding decision type file using the table below. Then check `docs/decisions/<type>.md` for rows where the Decision column text matches the entry's topic.
2. **Lessons**: match by tags → for each extracted entry, infer which of the 8 tag vocabulary items best matches its domain (e.g., auth rules → `security`, error patterns → `error-handling`). Then grep `tags:` frontmatter in `docs/lessons/*.md` for exact tag value matches from the 8-item vocabulary.

**Domain-to-decision-file mapping** (from `plugins/forge/references/shared/decision-logging.md` Section 1):

| Spec domain keywords | Decision file |
|---------------------|---------------|
| system structure, layering, modules, architecture | `architecture.md` |
| API contracts, data shapes, serialization, interface | `interface.md` |
| schema, indexing, soft-delete, data model, data ownership | `data-model.md` |
| libraries, versions, packages, dependencies | `dependencies.md` |
| error types, status codes, error propagation, error handling | `error-handling.md` |
| test patterns, coverage, mocking, testing | `testing.md` |
| auth, permissions, data protection, security, access control | `security.md` |
| dev environment, tooling, deployment, local setup | `local-dev-deployment.md` |
| naming, conventions, coding standards | `architecture.md` |
| validation, state transitions, calculation rules | closest match or `architecture.md` |
| performance, latency, caching, throughput | `architecture.md` |

If a domain does not clearly map to any file, skip the decisions overlap check for that entry.

Collect matches as "Related existing entries" for display in Step 6.

### Early exit

If ALL extracted items are `[LOCAL]`, write the preview files but do NOT proceed to Step 6.

Write `docs/features/<slug>/specs/.integrated` for the early-exit case:
```yaml
feature: "<slug>"
integrated: "<date>"
status: "skipped: all local"
biz_count: 0
tech_count: 0
```

Then proceed directly to **Step 12** (Record Task). The preview files remain for traceability.

## Step 6: Present to User for Review

Display the preview with categorization and detected overlaps:

```
Extracted specs for <slug>:

Business Rules:
  [CROSS] BIZ-001: <rule> → docs/business-rules/<domain>.md (e.g., auth rules → auth.md)
  [LOCAL] BIZ-002: <rule> → stays in feature

Technical Specs:
  [CROSS] TECH-001: <spec> → docs/conventions/<topic>.md (e.g., API patterns → api.md)
  [LOCAL] TECH-002: <spec> → stays in feature

Related existing entries (may overlap):
  ⚠ decisions/error-handling.md: "Adopt AIError struct" — appears to overlap with TECH-001
  ⚠ lessons/gotcha-error-handling.md [error-handling] — may overlap with TECH-002

Domain overlap warnings:
  ⚠ docs/conventions/error-handling.md (domains: [error, status, response]) and docs/conventions/error-reporting.md (domains: [error, status, log]) share 66% of keywords — consider merging

For each overlap, choose: [skip] keep both | [replace] delete old + write new
```

Ask the user:
1. Which `[CROSS]` items should be integrated?
2. Which domain/topic file should each be merged into?
3. Any items to skip?
4. For each overlap: `[skip]` keep both, or `[replace]` delete old entry and write new?
5. For each domain overlap warning (>50% shared keywords): merge the files, keep separate, or adjust target file?

Write the user's choices to `docs/features/<slug>/specs/review-choices.md`:

```markdown
---
feature: "<slug>"
reviewed: "<date>"
---

# Review Choices

## Approved for Integration

- BIZ-001 → docs/business-rules/<domain>.md
- TECH-001 → docs/conventions/<topic>.md

## Skipped

- (any items the user chose to skip)

## Related Existing Entries

- decisions/error-handling.md row "Adopt AIError struct" → replaced by BIZ-auth-003
- lessons/gotcha-error-handling.md → deleted, superseded by TECH-error-005
```

## Step 7: Integrate Approved Specs

For each item listed as "Approved" in `review-choices.md`:

**Business rules** → append to `docs/business-rules/<domain>.md`:
- Create the file if it doesn't exist
- When creating a new file, write YAML frontmatter with `title` (derived from the domain name) and `domains` (derived per the Domain Derivation Rules above)
- Add a source reference linking back to the feature
- Group by rule category within the file

**Technical specs** → append to `docs/conventions/<topic>.md`:
- Create the file if it doesn't exist
- When creating a new file, write YAML frontmatter with `title` (derived from the topic name) and `domains` (derived per the Domain Derivation Rules above)
- Add a source reference linking back to the feature
- Group by spec category within the file

### New File Frontmatter

When creating a new project-level spec file, include this frontmatter:

```yaml
---
title: "<Descriptive Title>"
domains: [<keyword1>, <keyword2>, ..., <keywordN>]
---
```

- `title`: Human-readable title derived from the domain/topic name (existing behavior, unchanged)
- `domains`: 3-7 specific keywords derived from the spec content being written into the file, per the Domain Derivation Rules

For **existing files** that lack a `domains` field, derive and add it during integration (do not modify existing `title`).

### Project-Global ID Encoding

Each entry gets a project-global ID (not the feature-local BIZ-NNN/TECH-NNN):

- **Prefix** derived from target filename: `business-rules/auth.md` → prefix `auth`, `conventions/api.md` → prefix `api`
- **Sequence**: file-internal — find max existing NNN in the target file + 1
- **Format**: `BIZ-<domain>-<NNN>` for business rules, `TECH-<topic>-<NNN>` for tech specs
- **Examples**: `BIZ-auth-001`, `TECH-api-003`
- **Source traceability**: `Source: feature/<slug> BIZ-001` (links back to feature-local preview ID)

### Overlap Resolution

For each user-approved "replace" choice from Step 6:

1. **Decisions**: delete the matching table row from `docs/decisions/<type>.md` — match by the Decision column text (substring match on the text recorded in review-choices.md)
2. **Lessons**: delete `docs/lessons/<filename>.md`
3. Update `docs/decisions/manifest.md` counts if any decision rows were deleted
4. Write new entry to target file with project-global ID

## Step 8: Write Integration Marker + Update Manifest

Write `docs/features/<slug>/specs/.integrated`:

```yaml
feature: "<slug>"
integrated: "<date>"
biz_count: <N>
tech_count: <M>
replaced:
  - decisions/error-handling.md row "Adopt AIError..." → BIZ-auth-003
  - lessons/gotcha-error-handling.md → TECH-error-005
```

The `replaced` field is omitted if no overlaps were resolved.

Update `docs/features/<slug>/manifest.md` to reference the integrated specs.

## Step 9: Detect Drift in Project-Level Specs

Read all project-level spec files and validate each rule against the current codebase:

1. **Read all spec files**:
   - `docs/business-rules/*.md` — all business rule files
   - `docs/conventions/*.md` — all technical convention files

2. **Validate each rule against code**: For each rule in every spec file, search the codebase for the keywords and patterns described in the rule. Compare the rule's stated behavior against the actual code implementation:
   - Extract key domain terms, function names, file paths, or behavior descriptions from each rule
   - Search the relevant source files for those terms
   - Determine whether the rule's description still matches the code's current behavior

3. **Classify each rule**:
   - `current` — rule description matches current code behavior
   - `drifted` — rule description is partially or fully inconsistent with current code (e.g., renamed function, changed threshold, modified behavior)
   - `orphaned` — the code the rule describes no longer exists (e.g., deleted module, removed feature)

4. **Output drift report**: Write a summary of all classifications. If no drift is found (all `current`), skip Steps 10-11 and proceed to Step 12.

### Drift-Only Mode Entry

If running in drift-only mode (no PRD/design files exist), start here at Step 9. Skip Steps 1-8 entirely.

## Step 10: Auto-Fix Drifted Specs

For each rule classified as `drifted` or `orphaned` in Step 9:

1. **Drifted rules**: Update the rule's description/behavior text in-place to match the current code. Preserve the project-global ID (e.g., `BIZ-auth-001`) — only update the descriptive text, not the ID or structural format. Mark the updated date in the entry's frontmatter or metadata.

2. **Orphaned rules**: Remove the rule entry from the spec file. Record the deletion for the commit message in Step 11:
   - Rule ID (e.g., `TECH-api-002`)
   - Reason for deletion (e.g., "corresponding module `X` removed in commit abc1234")

3. **Detect implicit new rules**: While scanning the code for drift, if you discover new patterns, conventions, or business logic that should be documented at the project level but are not in any spec file:
   - Extract the candidate rule with `[CROSS]` classification
   - Present to user for confirmation before appending
   - Append confirmed rules to the appropriate spec file with a new project-global ID

4. **Re-derive `domains` frontmatter**: When a file's content changes substantially (rules updated, added, or removed), re-derive the `domains` field per the Domain Derivation Rules. Compare the new domain set against the existing one:
   - If domains have changed, update the frontmatter in-place
   - If the updated `domains` cause a new >50% overlap with another file's domains, flag in the commit message and notify the user

### Preservation Rules

- Project-global IDs must never change during auto-fix — only description and behavior text updates
- File structure and formatting must remain consistent with the existing spec file conventions
- Deleted rules must be recorded with ID and reason for traceability in git history

## Step 11: Commit Spec Changes

If any spec files were modified in Step 10:

1. Stage all changed files under `docs/business-rules/` and `docs/conventions/`
2. Commit with a descriptive message listing:
   - Which rule IDs were updated (drift fix)
   - Which rule IDs were removed (orphaned) and why
   - Which rule IDs were added (implicit new rules)

Example commit message:

```
chore(specs): drift auto-fix — 2 updated, 1 removed, 1 added

Updated:
  - BIZ-auth-001: align with renamed validateToken → verifySession
  - TECH-api-003: reflect new rate limit threshold (100 → 200)

Removed:
  - TECH-api-002: corresponding legacy proxy module removed

Added:
  - TECH-error-006: implicit error wrapping convention (user-confirmed)
```

If no changes were made (all rules `current`), skip this step.

## Step 12: Record Task

Invoke the skill:

```
Skill(skill="submit-task")
```

Omit `coverage` from record.json — the noTest flag in index.json auto-sets it.

## Rules

- Only extract rules that are **explicitly stated** in source documents — do not infer
- Feature-specific implementation details stay in the feature, not in specs
- Never overwrite existing project-level spec files — append and merge, unless drift is detected in Step 9
- Overlap detection uses tag matching for lessons and filename matching for decisions
- Project-global IDs use filename-derived prefix + file-internal sequence
- Drift detection compares rule keywords against actual code — not simple text matching (mitigates false positives)
- Deleted rules must be recorded in commit message with ID and deletion reason
- Project-global IDs must be preserved during auto-fix (only update description/behavior text)
- New implicit rules from code are extracted with `[CROSS]` classification and presented to user before appending
- `domains` frontmatter is derived from spec ID keywords and source keywords — never invented by the agent
- Each file gets 3-7 specific domain keywords (not generic terms like "rule", "spec", "requirement")
- The existing `title` frontmatter behavior is unchanged — `domains` is an additive field
- Domain overlap >50% between files triggers a warning during the user confirmation step (Step 6)
- During drift detection (Steps 9-10), `domains` are re-derived when file content changes substantially

## Related Skills

| Skill | Relationship |
|-------|-------------|
| `/write-prd` | Upstream: source of business rules |
| `/tech-design` | Upstream: source of technical specs |
| `/graduate-tests` | Predecessor: T-test-4 before this T-specs-1 |
| `/submit-task` | Downstream: records task completion |
