---
status: "completed"
started: "2026-05-23 10:07"
completed: "2026-05-23 10:31"
time_spent: "~24m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Drift-only spec consolidation: detected and auto-fixed drift in docs/conventions/forge-cli-reference.md. Removed deleted commands (forge e2e group, forge probe), corrected forge task status to read-only, fixed function name casing and constant locations. All other 16 spec files validated as current.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/forge-cli-reference.md
- docs/.vocabulary.md

### Key Decisions
- Moved deleted commands (forge e2e group, forge probe) to '已移除的命令' section for traceability rather than silently removing them
- Corrected exported function names (CalcSlugColWidth, PadRight, TruncateSlug) and constant location (base/output.go) to match actual code

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
无

## Notes
Drift caused by commit 168deb5e which deleted the forge e2e command group and forge probe command. Vocabulary index regenerated with date bump only — no domain changes.
