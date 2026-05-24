---
status: "completed"
started: "2026-05-24 09:49"
completed: "2026-05-24 09:50"
time_spent: "~1m"
---

# Task Record: 3 Add Reference Files fallback rule to task-executor.md Hard Constraints

## Summary
Added rule 8 (SPEC AUTHORITY FALLBACK) to task-executor.md Hard Constraints block, requiring agents to proactively read task file Reference Files even when synthesized strategy omits the declaration.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/agents/task-executor.md

### Key Decisions
无

## Document Metrics
1 file modified, 4 lines added (rule 8 + 3 sub-bullets)

## Referenced Documents
- docs/proposals/spec-authority-enforcement/proposal.md#Proposed-Solution
- docs/proposals/spec-authority-enforcement/proposal.md#Priority-Rules
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] New rule added as item 8 in <EXTREMELY-IMPORTANT> Hard Constraints block after rule 7
- [x] Rule states agent must proactively read ## Reference Files from task file, even if synthesized strategy omits it
- [x] Rule references priority order: ## Hard Rules > ## Reference Files > existing code
- [x] Rule is a fallback (not a duplicate of template-layer IMPORTANT block)
- [x] Modification takes effect immediately (task-executor.md distributed as-is)
- [x] Did NOT modify existing rules 1-7
- [x] Did NOT create a second EXTREMELY-IMPORTANT block

## Notes
Rule wording follows the suggested template from the task Implementation Notes, with minor formatting consistency adjustments.
