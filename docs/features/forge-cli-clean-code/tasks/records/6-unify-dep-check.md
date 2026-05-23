---
status: "completed"
started: "2026-05-24 02:05"
completed: "2026-05-24 02:16"
time_spent: "~11m"
---

# Task Record: 6 Unify dependency check logic

## Summary
Consolidated duplicate dependency check logic (wildcard matching and unmet dep resolution) from 7 locations into 3 unified exported functions in pkg/task/deps.go: ResolveWildcardDep, GetUnmetDeps, and IsDepSatisfied. Net reduction of 114 lines of code.

## Changes

### Files Created
- forge-cli/pkg/task/deps.go
- forge-cli/pkg/task/deps_test.go

### Files Modified
- forge-cli/pkg/task/add.go
- forge-cli/pkg/task/statemachine.go
- forge-cli/internal/cmd/task/check_deps.go
- forge-cli/internal/cmd/task/claim.go
- forge-cli/internal/cmd/task/status.go
- forge-cli/internal/cmd/task/validate_index.go
- forge-cli/internal/cmd/task/validate_index_test.go
- forge-cli/internal/cmd/task/check_deps_test.go

### Key Decisions
- Created three complementary functions rather than one monolithic function: ResolveWildcardDep (pure matching), GetUnmetDeps (status-aware unmet resolution), IsDepSatisfied (status predicate)
- Preserved claim.go's vacuously-satisfied semantics for unknown deps by filtering GetUnmetDeps results rather than adding options/flags to the unified function
- CheckTransitionDeps uses IsDepSatisfied but not GetUnmetDeps because it intentionally does not expand wildcards (wildcards are treated as exact dep lookups that fail)
- validateLiveness uses ResolveWildcardDep + IsDepSatisfied rather than GetUnmetDeps because it needs finer-grained status classification (active vs completed vs blocked)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 500
- **Failed**: 0
- **Coverage**: 92.5%

## Acceptance Criteria
- [x] Single dependency check function created in pkg/task/
- [x] All duplicate implementations replaced with calls to the unified function
- [x] .x wildcard handling preserved in the unified implementation
- [x] go build ./... passes
- [x] go test ./... passes

## Notes
9 files changed, 51 insertions, 165 deletions (-114 net). pkg/task coverage: 92.5%, internal/cmd/task coverage: 70.1%.
