---
status: "completed"
started: "2026-05-27 01:03"
completed: "2026-05-27 01:07"
time_spent: "~4m"
---

# Task Record: 4 Add surface-suffixed type variants to ValidTypes

## Summary
Add IsValidType function with pattern-based validation for surface-suffixed task types (e.g. test.gen-scripts.cli). Updated 3 call sites (prompt.go, validate_index.go, add.go) to use IsValidType instead of ValidTypes map lookup.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/types.go
- forge-cli/pkg/task/types_test.go
- forge-cli/pkg/prompt/prompt.go
- forge-cli/internal/cmd/task/validate_index.go
- forge-cli/internal/cmd/task/add.go

### Key Decisions
- Used pattern matching (strip last dot-segment, check if base is in SystemTypes) rather than hardcoding surface keys into ValidTypes — per task Hard Rules
- Only system types (auto-generated) accept surface suffixes; non-generated types like coding.* remain strictly validated

## Test Results
- **Tests Executed**: Yes
- **Passed**: 5
- **Failed**: 0
- **Coverage**: 87.7%

## Acceptance Criteria
- [x] Surface-suffixed types (e.g. test.gen-scripts.cli) pass Synthesize validation
- [x] All existing ValidTypes entries still work
- [x] Existing tests pass (go test ./...)

## Notes
IsValidType replaces direct ValidTypes map lookup at all 3 validation call sites. ValidTypes map itself remains unchanged as the source of truth for base types.
