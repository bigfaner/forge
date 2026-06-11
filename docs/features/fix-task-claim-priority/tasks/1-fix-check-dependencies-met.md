---
id: "1"
title: "Fix checkDependenciesMet to block on pending fix tasks"
priority: "P0"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Fix checkDependenciesMet to block on pending fix tasks

## Description
When a fix task (`fix-1`) and a business task (`4`) are both P0 with met dependencies, `forge task claim` returns the business task instead of the fix task. Root cause: `checkDependenciesMet()` only checks whether dependency tasks are "completed" or "skipped" — it does NOT check whether pending fix tasks exist for those dependencies.

Enhance `checkDependenciesMet` to also scan for pending fix tasks whose `SourceTaskID` matches any of the current task's dependencies. This uses the existing `Task.SourceTaskID` field.

## Reference Files
- `docs/proposals/fix-task-claim-priority/proposal.md` — Source proposal
- `forge-cli/internal/cmd/claim.go` — Contains `checkDependenciesMet` (line ~238)
- `forge-cli/internal/cmd/claim_test.go` — Existing `TestClaimNextTask` tests
- `forge-cli/pkg/task/types.go` — `Task.SourceTaskID` and `Task.Type` fields

## Acceptance Criteria
- [ ] `checkDependenciesMet` returns false when a dependency has a pending fix task with matching sourceTaskID
- [ ] `checkDependenciesMet` returns true when all dependencies' fix tasks are completed
- [ ] `claimNextTask` returns fix-task when it coexists with business tasks that depend on the fix's source task
- [ ] All existing `TestCompareVersionIDs` and `TestClaimNextTask` tests pass unchanged
- [ ] Fix chain scenario: task blocked until all fix tasks for a dependency complete
- [ ] Unrelated fix tasks (different sourceTaskID) don't block tasks with different dependencies

## Hard Rules
- Must NOT change `compareVersionIDs` or `parseSegment`
- Must NOT change fix-task ID format or template
- New check only activates when fix tasks exist for dependencies — no behavior change otherwise

## Implementation Notes
- Add a loop in `checkDependenciesMet` after the existing dependency check: scan all tasks for fix tasks whose `SourceTaskID` matches any of the current task's dependencies and whose status is pending or in_progress
- ~10 lines of code in claim.go
- Add unit tests for each scenario in the Key Scenarios section of the proposal
- Add integration test to `TestClaimNextTask` for the fix-task blocking downstream tasks scenario
- Fix tasks created without `--source-task-id` should not affect behavior (sort ordering unchanged)
- `--source-task-id` resolves to root ancestor, so fix chains don't create circular dependencies
