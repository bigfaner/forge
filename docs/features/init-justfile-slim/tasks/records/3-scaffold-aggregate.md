---
status: "completed"
started: "2026-06-09 22:43"
completed: "2026-06-09 22:51"
time_spent: "~8m"
---

# Task Record: 3 实现 scaffold aggregate 聚合模式

## Summary
Implemented scaffold --aggregate mode: GenerateAggregate function generating install/ci/clean aggregate recipes and test-setup multi-service orchestration recipe. Wired --aggregate flag into CLI register.go with forgeconfig.ReadSurfaces integration.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/scaffold/generate.go
- forge-cli/internal/cmd/scaffold/generate_test.go
- forge-cli/internal/cmd/scaffold/register.go
- forge-cli/internal/cmd/scaffold/types.go

### Key Decisions
- ReadSurfacesFunc exported as variable for testability (tests can mock forgeconfig.ReadSurfaces)
- serviceOrder map defines dependency-based startup order (api=0, web=1, mobile=2)
- Aggregate recipes use empty-body fallback 'echo No surfaces configured' instead of erroring
- surfacesToEntries converts scalar '.' key to empty string for consistent recipeName usage

## Test Results
- **Tests Executed**: Yes
- **Passed**: 53
- **Failed**: 0
- **Coverage**: 83.9%

## Acceptance Criteria
- [x] --aggregate generates install, ci, clean aggregate recipes without # user-customized marker
- [x] ci recipe aggregates all surface lint + compile + unit-test, excludes surface-level test
- [x] Multiple service-type surfaces generate test-setup with dependency-ordered startup and reverse teardown
- [x] Pure cli/tui combination does not generate test-setup aggregate recipe

## Notes
19 new tests added for aggregate mode. Total scaffold tests: 53 (48 generate + 5 register/helpers). Coverage 83.9% exceeds 80% target.
