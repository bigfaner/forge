---
status: "completed"
started: "2026-05-24 21:14"
completed: "2026-05-24 21:23"
time_spent: "~9m"
---

# Task Record: 3 Improve init summary to show actual detected surface types

## Summary
Changed forge init summary format from opaque 'N mappings' to actual detected surface types with compact source annotations: scalar form shows 'cli (inferred:cmd-dir)' or 'cli (from cobra)'; map form shows 'forge-cli/cli=cli (inferred:cmd-dir)' per entry.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/init.go
- forge-cli/internal/cmd/init_surfaces.go
- forge-cli/internal/cmd/init_surfaces_test.go

### Key Decisions
- Added formatCompactSourceAnnotation for init summary (short: '(inferred:cmd-dir)', '(from cobra)') separate from existing formatSourceAnnotation for TUI (long: '(inferred from cmd/ directory structure)')
- Made askSurfaceConfirmation a variable function returning (SurfacesMap, SourcesMap, bool) so Sources propagate from DetectResult through to the init summary detail string
- writeSurfacesToConfig now accepts SourcesMap parameter for the detail string only — sources are never persisted to config

## Test Results
- **Tests Executed**: Yes
- **Passed**: 18
- **Failed**: 0
- **Coverage**: 60.8%

## Acceptance Criteria
- [x] Scalar form: 'cli (inferred:cmd-dir)' when inferred, 'cli (from cobra)' when dependency, 'cli' bare
- [x] Map form: 'forge-cli=cli (inferred:cmd-dir)' per entry instead of (N mappings)
- [x] Subdir detection summary shows compact annotations
- [x] Existing signal detection: cobra -> 'cli (from cobra)'
- [x] Re-run skip: SKIPPED surfaces (already configured)

## Notes
Hard Rule verified: grep confirms no programmatic parser of init summary exists in codebase. All existing tests pass unchanged.
