---
status: "completed"
started: "2026-05-24 15:23"
completed: "2026-05-24 15:52"
time_spent: "~29m"
---

# Task Record: 4 forge task index: migrate from Interfaces to Surfaces

## Summary
Migrated BuildIndex from Interfaces to Surfaces: SurfaceTypes() filters unknown types with log.Warn, empty surfaces returns error, removed uiInterfaces map, renamed BodyContext.Interfaces to BodyContext.SurfaceTypes, updated all tests

## Changes

### Files Created
- forge-cli/pkg/forgeconfig/surfaces_test.go

### Files Modified
- forge-cli/pkg/forgeconfig/config.go
- forge-cli/pkg/forgeconfig/detect.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/extract.go
- forge-cli/pkg/task/extract_test.go

### Key Decisions
- Used KnownSurfaceTypes map for O(1) known-type lookup instead of linear slice search
- Unknown surface types produce slog.Warn (not error) and are silently excluded from results
- Empty surfaces is a hard error with actionable message pointing to forge init
- Added early return guard in ResolveFirstTestDep to prevent panic when no gen-journeys tasks exist
- Replaced panic-with-log.Fatal in gen-journeys tests with graceful no-op behavior

## Test Results
- **Tests Executed**: Yes
- **Passed**: 93
- **Failed**: 0
- **Coverage**: 87.6%

## Acceptance Criteria
- [x] BuildIndex reads Surfaces (map[string]string) from config instead of Interfaces ([]string)
- [x] SurfaceTypes() extracts deduplicated surface types, filtering unknown with log.Warn
- [x] Empty surfaces produces explicit error (exit 1), not silent skip
- [x] uiInterfaces map removed from autogen.go
- [x] BodyContext.Interfaces renamed to BodyContext.SurfaceTypes
- [x] All existing tests updated and passing
- [x] New tests for unknown type filtering and validation

## Notes
Coverage: task package 87.6%, forgeconfig package 87.3%. ResolveFirstTestDep now handles missing gen-journeys gracefully instead of panicking.
