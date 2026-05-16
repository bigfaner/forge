---
id: "3"
title: "Update tests for lazy unblock and self-block scenarios"
priority: "P0"
estimated_time: "1.5h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 3: Update tests for lazy unblock and self-block scenarios

## Description

Update `claim_test.go` to cover the new behaviors: `checkDependenciesMet` self-block check and `claimNextTask` lazy unblock scan. Existing tests that rely on the old "blocked stays blocked until manual intervention" behavior may need adjustment.

## Reference Files
- `docs/proposals/task-lifecycle-hardening/proposal.md` — Source proposal
- `forge-cli/internal/cmd/claim_test.go` — Target test file

## Acceptance Criteria

- [ ] Test: `checkDependenciesMet` returns false when active fix-task has `SourceTaskID == selfID`
- [ ] Test: `checkDependenciesMet` returns true when fix-task targeting self is completed
- [ ] Test: `claimNextTask` auto-unblocks blocked task whose dependency just completed
- [ ] Test: blocked task stays blocked when fix-task targeting it is still active
- [ ] Test: `--block-source` scenario (fix done → claim → source auto-unblocked)
- [ ] Test: auto-downgrade scenario (task blocked → submit passes → claim auto-unblocks)
- [ ] All existing tests pass with updated expectations

## Hard Rules

- Use table-driven tests following existing patterns in claim_test.go
- TDD: write failing tests first, then verify implementation from tasks 1-2 makes them pass

## Implementation Notes

Key scenarios from proposal:
1. Task A completed → next claim auto-unblocks blocked task B (depends on A) and claims it
2. Fix-task still active → source task stays blocked despite deps met
3. `--block-source` creates fix-task → source goes blocked → fix completes → next claim auto-unblocks source
4. Auto-downgrade (test fail) → task goes blocked → re-submit → next claim auto-unblocks
