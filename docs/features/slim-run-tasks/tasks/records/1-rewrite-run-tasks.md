---
status: "completed"
started: "2026-05-16 10:08"
completed: "2026-05-16 10:10"
time_spent: "~2m"
---

# Task Record: 1 Rewrite run-tasks.md for token efficiency

## Summary
Compressed run-tasks.md from 256 lines to 146 lines while preserving all semantics. Three changes applied: (1) replaced forge task query with forge task status in Step 2b, (2) silent gate execution with output redirect to .forge/tmp/ files and tail on failure, (3) removed verbose explanations and condensed error handling table.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/run-tasks.md

### Key Decisions
- Merged 2a/2b/2c sub-sections into compact prose blocks instead of separate headed sections
- Combined pre-flight checks inline with gate execution code blocks to save lines
- Used single-line fix-task command format with combined --var flags to reduce code block size

## Test Results
- **Tests Executed**: No (noTest task)
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] run-tasks.md reduced from ~250 to ~150 lines (preserve mermaid diagram, EXTREMELY-IMPORTANT blocks verbatim)
- [x] Step 2b uses forge task status instead of forge task query
- [x] Breaking Gate (Step 3a): test output redirected to file, only exit code checked on success; on failure, tail last 20 lines
- [x] E2E Gate (Step 3b): same silent treatment — redirect output, check exit code, tail on failure
- [x] Error Handling table condensed to inline notes or compact format
- [x] All functionality preserved: claim, dispatch+verify, fix-task spawning, main-session routing, record-missing recovery, post-completion summary
- [x] Pre-flight checks for Breaking Gate preserved
- [x] Pre-flight checks for E2E Gate preserved

## Notes
noTest documentation task. Line count: 256 -> 146 (43% reduction). Mermaid diagram and all EXTREMELY-IMPORTANT blocks preserved verbatim.
