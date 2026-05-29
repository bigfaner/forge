---
status: "completed"
started: "2026-05-29 14:26"
completed: "2026-05-29 14:31"
time_spent: "~5m"
---

# Task Record: 1 Phase 0: 改进 fix task description 信息完整性

## Summary
Replaced ExtractConciseError tail-10 with --- FAIL: line extraction in addSingleFixTask description. New conciseError function extracts all --- FAIL: lines from test output; falls back to ExtractConciseError for compile/fmt/lint steps. Added ExtractFailLines to pkg/just and conciseError to quality_gate.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/just/just.go
- forge-cli/pkg/just/just_test.go
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/quality_gate_test.go

### Key Decisions
- Added ExtractFailLines as a standalone function in pkg/just for reusability
- Created conciseError helper in quality_gate.go to encapsulate FAIL-first, tail-fallback logic
- TrimSpace before prefix check to handle indented --- FAIL: lines from Go test output

## Test Results
- **Tests Executed**: Yes
- **Passed**: 11
- **Failed**: 0
- **Coverage**: 68.1%

## Acceptance Criteria
- [x] addSingleFixTask description replaces tail-10 with all --- FAIL: lines from output
- [x] Fallback to ExtractConciseError tail behavior when no --- FAIL: lines (compile/fmt/lint)
- [x] Unit tests cover: output with FAIL lines, output without FAIL lines, empty output

## Notes
pkg/just coverage: 84.9%. All existing tests remain green. No version bump needed as this is internal behavior improvement, not a CLI-visible change.
