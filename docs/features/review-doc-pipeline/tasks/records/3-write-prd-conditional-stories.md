---
status: "completed"
started: "2026-05-24 00:12"
completed: "2026-05-24 00:12"
time_spent: ""
---

# Task Record: 3 Add conditional user stories to write-prd Step 7

## Summary
Added conditional gate to write-prd Step 7: skip user story generation for doc-only features, generate normally for code features

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/write-prd/SKILL.md

### Key Decisions
无

## Document Metrics
1 file modified, 2 paragraphs added to Step 7

## Referenced Documents
- docs/proposals/review-doc-pipeline/proposal.md
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md

## Review Status
final

## Acceptance Criteria
- [x] Step 7 includes condition check for non-compilable/non-runnable In Scope items
- [x] Step 7 generates user stories normally for compilable/runnable file paths
- [x] Condition logic consistent with task type assignment rules (same heuristic)
- [x] When skipping, Step 7 outputs a brief note explaining why

## Notes
Gate prepended to existing Step 7 content without restructuring. Detection heuristic mirrors quick-tasks/breakdown-tasks type assignment: non-compilable/non-runnable = doc-only.
