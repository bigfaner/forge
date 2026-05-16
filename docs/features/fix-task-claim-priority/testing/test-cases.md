---
feature: "fix-task-claim-priority"
sources:
  - docs/proposals/fix-task-claim-priority/proposal.md
generated: "2026-05-16"
---

# Test Cases: fix-task-claim-priority

## Summary

| Type | Count |
|------|-------|
| CLI  | 6  |
| **Total** | **6** |

---

## CLI Test Cases

## TC-001: Pending fix task blocks dependent business task
- **Source**: Proposal Key Scenario "Primary" + Success Criterion 1
- **Type**: CLI
- **Target**: cli/task-claim
- **Test ID**: cli/task-claim/pending-fix-task-blocks-dependent-business-task
- **Pre-conditions**: Task 3 is completed. Fix-1 exists with sourceTaskID "3" and status pending. Task 4 depends on task 3. Both fix-1 and task 4 are P0.
- **Steps**:
  1. Create task dependency graph: task 3 (completed), fix-1 (pending, sourceTaskID "3"), task 4 (pending, depends on task 3)
  2. Run `forge task claim`
  3. Verify the claimed task is fix-1, not task 4
- **Expected**: `checkDependenciesMet` returns false for task 4 because fix-1 (sourceTaskID "3") is pending. `claimNextTask` returns fix-1.
- **Priority**: P0

## TC-002: Completed fix task allows dependent business task
- **Source**: Proposal Key Scenario "Fix completed" + Success Criterion 2
- **Type**: CLI
- **Target**: cli/task-claim
- **Test ID**: cli/task-claim/completed-fix-task-allows-dependent-business-task
- **Pre-conditions**: Task 3 is completed. Fix-1 exists with sourceTaskID "3" and status completed. Task 4 depends on task 3.
- **Steps**:
  1. Create task dependency graph: task 3 (completed), fix-1 (completed, sourceTaskID "3"), task 4 (pending, depends on task 3)
  2. Run `forge task claim`
  3. Verify task 4 is eligible and can be claimed
- **Expected**: `checkDependenciesMet` returns true for task 4 because fix-1 is completed. Task 4 can be claimed.
- **Priority**: P0

## TC-003: Fix task claimed before business task when both eligible
- **Source**: Success Criterion 3
- **Type**: CLI
- **Target**: cli/task-claim
- **Test ID**: cli/task-claim/fix-task-claimed-before-business-task
- **Pre-conditions**: Task 3 is completed. Fix-1 exists with sourceTaskID "3" and status pending. Task 4 depends on task 3. Fix-1 and task 4 are both P0 with met dependencies.
- **Steps**:
  1. Create task dependency graph: task 3 (completed), fix-1 (pending, sourceTaskID "3", depends on task 3), task 4 (pending, depends on task 3)
  2. Run `forge task claim`
  3. Verify fix-1 is returned, not task 4
- **Expected**: `claimNextTask` returns fix-1 when it coexists with business task 4. Fix-1 has met dependencies (task 3 completed, no pending fix tasks for its own deps).
- **Priority**: P0

## TC-004: Fix chain blocks dependent task until all fix tasks complete
- **Source**: Proposal Key Scenario "Fix chain" + Success Criterion 5
- **Type**: CLI
- **Target**: cli/task-claim
- **Test ID**: cli/task-claim/fix-chain-blocks-dependent-task
- **Pre-conditions**: Task 3 is completed. Fix-1 (sourceTaskID "3") is completed. Fix-2 (sourceTaskID "3") is pending. Task 4 depends on task 3.
- **Steps**:
  1. Create task dependency graph: task 3 (completed), fix-1 (completed, sourceTaskID "3"), fix-2 (pending, sourceTaskID "3"), task 4 (pending, depends on task 3)
  2. Run `forge task claim`
  3. Verify task 4 is not eligible
  4. Complete fix-2
  5. Run `forge task claim` again
  6. Verify task 4 is now eligible
- **Expected**: Task 4 remains blocked while fix-2 is pending. After fix-2 completes, all fix tasks for task 3 are done and task 4 becomes eligible.
- **Priority**: P1

## TC-005: Unrelated fix task does not block task with different dependency
- **Source**: Proposal Key Scenario "Fix for unrelated task" + Success Criterion 6
- **Type**: CLI
- **Target**: cli/task-claim
- **Test ID**: cli/task-claim/unrelated-fix-task-does-not-block
- **Pre-conditions**: Task 2 is completed. Task 3 is completed. Fix-1 exists with sourceTaskID "2" and status pending. Task 4 depends on task 3 only.
- **Steps**:
  1. Create task dependency graph: task 2 (completed), task 3 (completed), fix-1 (pending, sourceTaskID "2"), task 4 (pending, depends on task 3)
  2. Run `forge task claim`
  3. Verify task 4 is eligible
- **Expected**: `checkDependenciesMet` returns true for task 4 because fix-1's sourceTaskID "2" does not match task 4's dependency (task 3). Task 4 can be claimed.
- **Priority**: P1

## TC-006: No fix tasks preserves existing claim behavior
- **Source**: Proposal Key Scenario "No fix tasks" + Success Criterion 4
- **Type**: CLI
- **Target**: cli/task-claim
- **Test ID**: cli/task-claim/no-fix-tasks-preserves-existing-behavior
- **Pre-conditions**: Task 3 is completed. No fix tasks exist. Task 4 depends on task 3.
- **Steps**:
  1. Create task dependency graph: task 3 (completed), task 4 (pending, depends on task 3)
  2. Run `forge task claim`
  3. Verify task 4 is claimed as before
- **Expected**: Existing claim behavior is unchanged. `checkDependenciesMet` returns true for task 4. Task 4 is claimed.
- **Priority**: P0

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Proposal Key Scenario "Primary" + Success Criterion 1 | CLI | cli/task-claim | P0 |
| TC-002 | Proposal Key Scenario "Fix completed" + Success Criterion 2 | CLI | cli/task-claim | P0 |
| TC-003 | Success Criterion 3 | CLI | cli/task-claim | P0 |
| TC-004 | Proposal Key Scenario "Fix chain" + Success Criterion 5 | CLI | cli/task-claim | P1 |
| TC-005 | Proposal Key Scenario "Fix for unrelated task" + Success Criterion 6 | CLI | cli/task-claim | P1 |
| TC-006 | Proposal Key Scenario "No fix tasks" + Success Criterion 4 | CLI | cli/task-claim | P0 |
