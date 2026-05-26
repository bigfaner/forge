---
status: "completed"
started: "2026-05-27 01:29"
completed: "2026-05-27 01:37"
time_spent: "~8m"
---

# Task Record: T-review-doc Review Documentation Quality

## Summary
Reviewed all 53 AC items across 8 doc task groups. Fixed 9 non-conformances: unified manifest slug placeholder, removed stale test-cases reference, added template placeholder mappings, corrected output checklist, added commit step to breakdown-tasks, fixed gen-journeys error message, updated forge-distribution to reference /learn, deduplicated journey-contract-model.md, replaced all WebUI references with canonical Web.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/quick-tasks/templates/manifest-quick.md
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/gen-journeys/SKILL.md
- docs/conventions/forge-distribution.md
- plugins/forge/skills/gen-contracts/SKILL.md
- plugins/forge/skills/gen-contracts/rules/risk-density.md
- plugins/forge/skills/gen-test-scripts/SKILL.md
- plugins/forge/skills/gen-test-scripts/rules/run-to-learn.md
- plugins/forge/skills/gen-test-scripts/types/ui.md
- plugins/forge/skills/run-tests/rules/env-check.md
- plugins/forge/skills/eval/rubrics/contract.md
- plugins/forge/skills/eval/rubrics/journey.md

### Key Decisions
无

## Document Metrics
53 AC items checked: 44 passed, 9 fixed

## Referenced Documents
- docs/proposals/pipeline-spec-code-alignment/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
fixes-applied

## Acceptance Criteria
- [x] 10-submit-task-path: record format paths resolve correctly
- [x] 11-architecture-conventions: all 8 sub-items verified
- [x] 12-templates-guidance: manifest slug unified, test-cases removed, placeholder mappings added, commit step added, error message fixed, /learn reference updated, journey-contract-model deduplicated
- [x] 6-ghost-fields-cleanup: all 7 sub-items verified (no stale references)
- [x] 7-surface-type-consistency: two-layer strategy documented, task-doc explains absence, WebUI replaced with Web, doc.review listed correctly
- [x] 8-dispatcher-pipeline-logic: all 10 sub-items verified
- [x] 9-fix-task-template-vars: all --var params present, IT impact assessment in both skill docs, grouping by directory

## Notes
9 fixes applied across 13 files. No AC items remain failing. prompt-template-hierarchy.md does not exist (AC-20 noted but no deliverable to fix).
