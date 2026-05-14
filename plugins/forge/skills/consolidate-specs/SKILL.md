---
name: consolidate-specs
description: Extract business rules and tech specs from feature docs into preview files, detect overlaps with existing knowledge, user confirms before integrating to project-level dirs.
---

# /consolidate-specs

Extract key business rules and technical specifications from feature documents into structured spec files. User reviews and confirms before integration to project-level shared directories.

**Core principle**: Consolidation extracts cross-cutting knowledge that outlives a single feature. Feature-specific details stay in the feature; reusable rules and specs are promoted.

<HARD-GATE>
- Do NOT integrate specs without explicit user confirmation
- Do NOT overwrite existing project-level spec files — append only
- Do NOT run if `specs/.integrated` marker exists (idempotent)
- Do NOT infer rules not explicitly stated in source documents
</HARD-GATE>

## When to Use

- As T-test-5 after all e2e tests graduate (standard task pipeline)
- User invokes `/consolidate-specs` manually at any time

## Prerequisites

Check before running. Abort and prompt user if missing:

| Artifact | Condition | Action if not met |
|----------|-----------|-------------------|
| Feature context | `forge feature` must be set | Run `forge feature <slug>` first |
| `prd/prd-spec.md` | Must exist | Run `/write-prd` first |
| `design/tech-design.md` | Must exist | Run `/tech-design` first |

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
Step 9: Record task
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

Then proceed directly to **Step 9** (Record Task). The preview files remain for traceability.

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

For each overlap, choose: [skip] keep both | [replace] delete old + write new
```

Ask the user:
1. Which `[CROSS]` items should be integrated?
2. Which domain/topic file should each be merged into?
3. Any items to skip?
4. For each overlap: `[skip]` keep both, or `[replace]` delete old entry and write new?

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
- Add a source reference linking back to the feature
- Group by rule category within the file

**Technical specs** → append to `docs/conventions/<topic>.md`:
- Create the file if it doesn't exist
- Add a source reference linking back to the feature
- Group by spec category within the file

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

## Step 9: Record Task

Invoke the skill:

```
Skill(skill="submit-task")
```

Omit `coverage` from record.json — the noTest flag in index.json auto-sets it.

## Rules

- Only extract rules that are **explicitly stated** in source documents — do not infer
- Feature-specific implementation details stay in the feature, not in specs
- Never overwrite existing project-level spec files — append and merge
- Overlap detection uses tag matching for lessons and filename matching for decisions
- Project-global IDs use filename-derived prefix + file-internal sequence

## Related Skills

| Skill | Relationship |
|-------|-------------|
| `/write-prd` | Upstream: source of business rules |
| `/tech-design` | Upstream: source of technical specs |
| `/graduate-tests` | Predecessor: T-test-4 before this T-test-5 |
| `submit-task` | Downstream: records task completion |
