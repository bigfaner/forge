---
status: "completed"
started: "2026-05-24 14:59"
completed: "2026-05-24 15:12"
time_spent: "~13m"
---

# Task Record: 3 forge surfaces CLI command

## Summary
Added independent `forge surfaces` CLI command with three sub-invocations: full listing (scalar outputs single type, map outputs path=surface per line), path query (segment prefix matching, exit 0 on match, exit 1 with stderr on no match), and --types (space-separated deduplicated known types). Command registered at top level in root.go, not under forge config. interfaces field completely ignored.

## Changes

### Files Created
- forge-cli/internal/cmd/surfaces.go
- forge-cli/internal/cmd/surfaces_test.go

### Files Modified
- forge-cli/internal/cmd/root.go

### Key Decisions
- Used surfacesPathError sentinel type to write raw error to stderr without cobra's 'Error: ' prefix, satisfying the Hard Rules output format requirement
- Used the existing write() helper for stdout/stderr output to avoid golangci-lint errcheck warnings
- Reset surfacesTypesFlag global in each test to prevent state leakage between test functions
- Filtered unknown surface types from --types output using KnownSurfaceTypes map, consistent with proposal's unknown type handling strategy

## Test Results
- **Tests Executed**: Yes
- **Passed**: 16
- **Failed**: 0
- **Coverage**: 88.9%

## Acceptance Criteria
- [x] forge surfaces -- scalar form outputs single type (exit 0); map form outputs path=surface per line (exit 0)
- [x] forge surfaces <path> -- returns surface type string (exit 0) or stderr error + exit 1 with manual config hint
- [x] forge surfaces --types -- space-separated deduplicated type list (exit 0)
- [x] Scalar form: forge surfaces <any-path> always returns the single value (exit 0)
- [x] interfaces field completely ignored -- config with both surfaces and interfaces outputs only surfaces data
- [x] Command registered at top level in root.go (not under forge config)

## Notes
无
