---
step: 1
title: Task Claim
journey: task-lifecycle
---

# Step 1: Task Claim

## Given
- A forge project with tasks in index.json
- Tasks have various statuses (pending, blocked, in_progress, completed)
- Some tasks have fix-task dependencies (sourceTaskID)

## When
- `forge task claim` is executed

## Then
- Highest-priority eligible task is claimed
- Self-block: tasks with active fix-tasks (sourceTaskID == selfID) are blocked
- Lazy unblock: blocked tasks auto-transition to pending when dependencies met
- Fix-task priority: fix-tasks are claimed before business tasks
- Auto-unblock is logged to stdout

## Contract Dimensions
- **Actor**: CLI user executing `forge task claim`
- **Input**: index.json with task definitions, state.json for current feature
- **Output**: CLI stdout with ACTION, TASK_ID, STATUS fields in --- block format
- **Side Effects**: index.json status updates, state.json creation
- **Error Cases**: no eligible tasks -> non-zero exit code
- **Invariants**: completed fix-tasks do not block; multiple fix-tasks require all completed
