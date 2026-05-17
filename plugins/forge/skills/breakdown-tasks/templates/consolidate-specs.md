---
id: "T-test-5"
title: "Consolidate Specs"
priority: "P2"
estimated_time: "20min"
dependencies: ["T-test-4.5"]
status: pending
noTest: true
mainSession: false
---

# Consolidate Specs

## Description

Call `/consolidate-specs` skill to extract business rules from PRD and technical specifications from design into `specs/` directory. Present preview to user for review before integrating to project-level shared directories.

## Reference Files

- `docs/features/<slug>/prd/prd-spec.md` — Source for business rules
- `docs/features/<slug>/prd/prd-user-stories.md` — Source for business context
- `docs/features/<slug>/design/tech-design.md` — Source for technical specs
- `docs/features/<slug>/design/api-handbook.md` — Source for API contracts (if exists)

## Acceptance Criteria

- [ ] `docs/features/<slug>/specs/biz-specs.md` exists with extracted business rules
- [ ] `docs/features/<slug>/specs/tech-specs.md` exists with extracted technical specs
- [ ] If any `[CROSS]` items exist: `docs/features/<slug>/specs/review-choices.md` exists with user's approved/rejected items
- [ ] If integration occurred: only items marked "approved" in review-choices.md were integrated to project-level dirs
- [ ] `docs/features/<slug>/specs/.integrated` marker exists (with counts if integration occurred, or "skipped: all local" if early exit)

## Skip Conditions

If ALL extracted items are `[LOCAL]` (no cross-cutting candidates), generate preview files only and mark task completed. No integration step needed.

If no extractable rules found in PRD/design, mark task completed.

If running under `/run-tasks` (non-interactive session) and CROSS items exist, write preview files and mark task as `blocked` with note "User review required for integration." Do NOT auto-integrate.

## User Stories

No direct user story mapping. This is a standard knowledge consolidation task.

## Implementation Notes

**Step 1: Verify prerequisites**

Confirm feature documents exist:
- `docs/features/<slug>/prd/prd-spec.md`
- `docs/features/<slug>/design/tech-design.md`

If missing, mark task `blocked` and stop.

Check idempotency: if `docs/features/<slug>/specs/.integrated` exists, skip (already done).

**Step 2: Extract and classify**

Run `/consolidate-specs` skill. The skill will:
- Read all feature documents
- Extract business rules into `docs/features/<slug>/specs/biz-specs.md`
- Extract technical specs into `docs/features/<slug>/specs/tech-specs.md`
- Classify each item as `[CROSS]` (cross-cutting) or `[LOCAL]` (feature-specific)
  - `[CROSS]`: Referenced by 2+ features, or domain invariant, or naming/convention rule
  - `[LOCAL]`: Only meaningful within this feature
- Detect overlaps with existing entries in `docs/decisions/` and `docs/lessons/`

**Step 3: Early exit or user review**

If ALL items are `[LOCAL]`, skip to Step 5 (record as completed, no integration).

Otherwise, present the preview files to the user:
- List all `[CROSS]` items with suggested target files (`docs/business-rules/<domain>.md`, `docs/conventions/<topic>.md`)
- User decides which to integrate and which to skip
- Write choices to `docs/features/<slug>/specs/review-choices.md`

**Step 4: Integrate approved items**

For each item approved in `review-choices.md`:
- Append to the appropriate project-level file
- Create the file if it doesn't exist
- Add source reference back to the feature

Write `docs/features/<slug>/specs/.integrated` marker with integrated counts.

**Step 5: Record**

Record task via `/record-task` skill. The preview files remain in `specs/` for traceability.
