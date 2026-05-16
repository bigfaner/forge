---
id: "1"
title: "Fix checkDependenciesMet: add SourceTaskID == selfID check"
priority: "P0"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Fix checkDependenciesMet: add SourceTaskID == selfID check

## Description

`checkDependenciesMet` currently checks for active fix-tasks whose `SourceTaskID` matches any of the task's **dependencies**. It does NOT check for active fix-tasks targeting the task **itself** (`SourceTaskID == selfID`). This gap means `--block-source` scenarios can slip through: the source task appears eligible even while its fix-task is still active.

This is a prerequisite for the lazy unblock scan (task 2), which relies on `checkDependenciesMet` being complete.

## Reference Files
- `docs/proposals/task-lifecycle-hardening/proposal.md` — Source proposal
- `forge-cli/internal/cmd/claim.go` — Target file (checkDependenciesMet at ~line 238)

## Acceptance Criteria

- [ ] `checkDependenciesMet` returns false when an active fix-task has `SourceTaskID == selfID` (regardless of the task's own dependencies)
- [ ] Existing behavior unchanged for tasks without active fix-tasks targeting them
- [ ] Existing tests continue to pass
- [ ] New test cases added for the self-block scenario

## Hard Rules

- No changes to submit.go
- Only adds blocking conditions — never relaxes existing checks

## Implementation Notes

In the fix-task check loop (~line 271), add a second pass: iterate all tasks, if any active fix-task has `SourceTaskID == selfID`, add it to `unmet`. The proposal's mermaid diagram shows this as a new check node (A4: `SourceTaskID == T.ID`).

The function signature `checkDependenciesMet(index, selfID string, t task.Task)` already receives `selfID` — use it directly.
