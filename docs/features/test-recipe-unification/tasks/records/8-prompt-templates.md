---
status: "completed"
started: "2026-05-24 22:19"
completed: "2026-05-24 22:21"
time_spent: "~2m"
---

# Task Record: 8 Update prompt templates for new recipe names

## Summary
Updated 3 prompt templates (gate.md, fix-record-missed.md, validation-code.md) to use `just unit-test` instead of `just test` for per-task gate scenarios

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/gate.md
- forge-cli/pkg/prompt/data/fix-record-missed.md
- forge-cli/pkg/prompt/data/validation-code.md

### Key Decisions
无

## Document Metrics
3 files modified, 6 replacements total (command references + table rows + mermaid node)

## Referenced Documents
- docs/proposals/test-recipe-unification/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] All 3 prompt templates reference `just unit-test` instead of `just test` for per-task gate scenarios
- [x] No residual `just test` references in gate/fix/validation prompt contexts

## Notes
Also updated failure-step table rows in all 3 files and mermaid flowchart decision node in gate.md to match unit-test naming
