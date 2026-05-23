---
status: "completed"
started: "2026-05-23 10:06"
completed: "2026-05-23 10:10"
time_spent: "~4m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Drift-only spec consolidation: detected and fixed 4 drifts in forge-cli-reference.md (forge task status signature, forge probe removed, forge research added, forge config set added). Regenerated vocabulary index. All other 16 spec files verified current.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/forge-cli-reference.md
- docs/.vocabulary.md

### Key Decisions
- forge probe marked as removed (not deleted from table) since it was referenced in spec but never existed as a CLI command in current codebase -- pkg/e2eprobe is a library only

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] All project-level spec files in docs/business-rules/ and docs/conventions/ validated against current codebase
- [x] Drifted specs auto-fixed and committed with [auto-specs] tag
- [x] Vocabulary index regenerated

## Notes
Doc-type task (doc.drift) -- no test metrics applicable. Drift summary: 1 drifted (task status signature), 1 orphaned (forge probe), 2 implicit new (forge research, forge config set). All 3 business-rules files and 13/14 convention files verified current.
