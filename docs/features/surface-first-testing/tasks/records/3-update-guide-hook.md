---
status: "completed"
started: "2026-06-02 21:57"
completed: "2026-06-02 21:58"
time_spent: "~1m"
---

# Task Record: 3 更新 guide.md hook 测试速查表

## Summary
Added Testing section to guide.md with Surface -> Test Type mapping table, e2e terminology constraint, test file location rules, and /test-guide prompt. Increment: 14 lines (within 20-line budget).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/hooks/guide.md

### Key Decisions
无

## Document Metrics
14 lines added, 5 surfaces covered, all AC met

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- plugins/forge/skills/test-guide/references/test-type-model.md

## Review Status
final

## Acceptance Criteria
- [x] guide.md new Testing section, increment <= 20 lines
- [x] Contains Surface -> Test Type mapping table (5 surfaces)
- [x] Contains e2e terminology constraint
- [x] Contains test file location rules
- [x] Contains prompt to run /test-guide for full strategy

## Notes
Mapping table content verified against test-type-model.md. Execution column adds reasoning context per Implementation Notes guidance.
