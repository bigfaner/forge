---
created: 2026-05-15
author: faner
status: Draft
---

# Proposal: Fix-Task IDs Lose Claim Priority to Numeric Business Tasks

## Problem

When a fix task (`fix-1`) and a numeric business task (`4`) are both P0 with met dependencies, `forge task claim` returns the business task instead of the fix task.

### Evidence

Observed during task execution on the `test-scripts-per-type` branch. After a compile error from task 3 triggered a fix task (`forge task add --template fix-task`), the dispatcher expected `fix-1` to be claimed next. Instead, `forge task claim` returned business task 4. Both were P0 with met dependencies.

Root cause: `checkDependenciesMet()` only checks whether dependency tasks are "completed" or "skipped". It does NOT check whether pending fix tasks exist for those dependencies. When task 3 completes but spawns fix-1, task 4 (depending on task 3) appears eligible because task 3 is "completed". The sort tiebreaker in `compareVersionIDs` then picks numeric `"4"` over alphabetic `"fix-1"`.

### Urgency

This breaks the dispatcher's core assumption (documented in `run-tasks.md`): "fix task (P0) will be claimed on next iteration." When fix tasks aren't claimed immediately, the dispatcher proceeds with the wrong business task while compile errors remain unresolved, leading to cascading failures.

## Proposed Solution

Enhance `checkDependenciesMet` (claim.go:238) to also check for pending fix tasks on upstream dependencies. When task 4 depends on task 3, and fix-1 has `sourceTaskID: "3"` with status pending/in_progress, task 4 is NOT eligible — regardless of sort ordering.

This uses the existing `Task.SourceTaskID` field, which the dispatcher already populates via `--source-task-id` when creating fix tasks.

### Innovation Highlights

The current dependency check is too literal: it checks if task 3 is "completed" but ignores the semantic state — task 3 may have spawned a fix task that hasn't run yet. The fix recognizes that a dependency is truly "met" only when no pending fix tasks exist for it. This is a standard pattern in build systems: a target is not "done" if its post-build validation has outstanding issues.

## Requirements Analysis

### Key Scenarios

- **Primary**: Task 4 depends on task 3. Fix-1 (sourceTaskID: "3") is pending → task 4 not eligible
- **Fix chain**: Fix-1 (sourceTaskID: "3") spawns fix-2 (sourceTaskID resolves to "3"). Both pending → task 4 not eligible until all fix tasks complete
- **Fix completed**: Fix-1 (sourceTaskID: "3") is completed → task 4 eligible (existing behavior)
- **No fix tasks**: No fix tasks for dependency → existing behavior unchanged
- **Fix for unrelated task**: Fix-1 (sourceTaskID: "2") pending, task 4 depends on task 3 → task 4 eligible (sourceTaskID doesn't match)
- **Multiple fix tasks**: Fix-1 and fix-2 both have sourceTaskID "3". Task 4 blocked until both complete

### Constraints & Dependencies

- `Task.SourceTaskID` is already populated by `forge task add --source-task-id` (dispatcher uses this)
- `Task.Type` field (`"fix"`) is available to identify fix tasks
- Must not change behavior when no fix tasks exist for a dependency
- `checkDependenciesMet` return type includes unmet dependency list — fix task blockages should appear there for diagnostics

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | No code changes | Dispatcher is broken; cascading failures | Rejected |
| Type-aware sort tier | Previous proposal | Sort fix tasks first | Only fixes ordering, not eligibility; downstream still appears eligible | Rejected: symptom fix, not root cause |
| compareVersionIDs hack | Initial proposal | Minimal code | Couples ID comparison to task semantics | Rejected: fragile |
| **Dependency check enhancement** | User insight | Fixes root cause; uses existing SourceTaskID; semantically correct | Slightly larger change; depends on SourceTaskID being populated | **Selected: most principled** |

## Feasibility Assessment

### Technical Feasibility

Straightforward. Add a loop in `checkDependenciesMet` after the existing dependency check: scan all tasks for fix tasks whose `SourceTaskID` matches any of the current task's dependencies and whose status is pending or in_progress. ~10 lines of code.

### Resource & Timeline

Single function change + tests. Estimated 30 minutes.

## Scope

### In Scope

- Enhance `checkDependenciesMet` in claim.go to check for pending fix tasks on dependencies
- Add unit tests for the new check behavior
- Add integration test to `TestClaimNextTask` for fix-task blocking downstream tasks
- Bump version (patch)

### Out of Scope

- Changes to `compareVersionIDs` or `parseSegment`
- Changes to fix-task ID format or template
- Changes to dispatcher protocol (`run-tasks.md`)
- Type-aware sort tier (dependency check makes this unnecessary)
- Changes to `SourceTaskID` resolution logic

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Fix tasks created without `--source-task-id` | L | H | Dispatcher always passes `--source-task-id`; if manually created, sort ordering is unchanged |
| Performance: scanning all tasks for each claim | L | L | Task lists are small (<100); linear scan is fine |
| Circular dependency with fix chains | L | M | `--source-task-id` resolves to root ancestor; fix chains don't create cycles |
| Existing claim tests break | L | H | New check only activates when fix tasks exist for dependencies |

## Success Criteria

- [ ] `checkDependenciesMet` returns false when a dependency has a pending fix task with matching sourceTaskID
- [ ] `checkDependenciesMet` returns true when all dependencies' fix tasks are completed
- [ ] `claimNextTask` returns fix-task when it coexists with business tasks that depend on the fix's source task
- [ ] All existing `TestCompareVersionIDs` and `TestClaimNextTask` tests pass unchanged
- [ ] Fix chain scenario: task blocked until all fix tasks for a dependency complete
- [ ] Unrelated fix tasks don't block tasks with different dependencies
