---
status: "completed"
started: "2026-05-26 01:43"
completed: "2026-05-26 01:49"
time_spent: "~6m"
---

# Task Record: T-clean-code Simplify and Clean Code

## Summary
Simplified and cleaned code across surface-aware-justfile feature: eliminated duplicate KnownSurfaceTypes, merged three scan functions into one, simplified detectMobile redundant conditions, replaced manual string concatenation with strings.Join, and fixed base.Debugf variadic expansion bug.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/surfaces.go
- forge-cli/internal/cmd/output.go
- forge-cli/internal/cmd/base/output.go
- forge-cli/internal/cmd/init_surfaces.go
- forge-cli/pkg/forgeconfig/detect_surface.go
- forge-cli/pkg/forgeconfig/config.go

### Key Decisions
- Reused pkg/forgeconfig.KnownSurfaceTypes in cmd package instead of maintaining duplicate map
- Merged scanSubdirs + scanSubdirsWithConflicts + scanSubdirsWithSources into single scanSubdirsWithSources with nil-sources opt-out
- Removed dead detectSurfaceAtDirWithConflicts wrapper that was only called by removed scan functions
- Extracted selectSurfaceEntry helper to eliminate duplicated huh option construction in editMapEntry/deleteMapEntry

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] go build ./... passes
- [x] go vet ./... passes
- [x] All tests in modified packages pass with race detection
- [x] No duplicate KnownSurfaceTypes definitions
- [x] Dead scan/detect functions removed

## Notes
code-quality.simplify maps to coding category. Five cleanups applied: (1) duplicate KnownSurfaceTypes, (2) three scan functions merged, (3) detectMobile xcodeproj redundant branches, (4) joinSlice string concat, (5) base.Debugf variadic bug.
