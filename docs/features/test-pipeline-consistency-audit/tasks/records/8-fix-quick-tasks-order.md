---
status: "completed"
started: "2026-05-27 15:33"
completed: "2026-05-27 15:37"
time_spent: "~4m"
---

# Task Record: 8 修正 quick-tasks 执行顺序和术语

## Summary
Fixed quick-tasks execution order from gen-journeys → gen-contracts → gen-test-scripts to gen-journeys → run-test, and replaced 'Integration Test Impact Assessment' with 'Test Impact Assessment' in both quick-tasks and breakdown-tasks SKILL.md files

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md

### Key Decisions
无

## Document Metrics
2 files modified, 3 changes: 1 order fix + 2 terminology replacements

## Referenced Documents
- docs/proposals/test-pipeline-consistency-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] quick-tasks/SKILL.md execution order corrected to gen-journeys → run-test
- [x] Both SKILL.md files: 'Integration Test Impact Assessment' → 'Test Impact Assessment'
- [x] breakdown-tasks/SKILL.md execution order unchanged (already correct)

## Notes
breakdown-tasks order was already correct per task spec, only terminology was updated. gen-test-scripts/types/ui.md 'Integration Test' preserved as UI-specific concept per Implementation Notes.
