---
name: claim-task
description: Use when you need to claim and start working on the next available task from the project task list. Claims the highest priority task with all dependencies met.
---

# Claim Task

Claim the next available task from `docs/features/<slug>/tasks/index.json`.

## Usage

```bash
task claim
```

## Output

```
---
KEY: <task-key>
ID: <task-id>
FILE: <task-file>
---
```

## After Claiming

1. Read task file: `docs/features/<slug>/tasks/<FILE>`
2. Implement following TDD (RED → GREEN → REFACTOR)
3. Update record: `/record-task <TASK_ID>`
4. Mark complete: `/set-task-status <ID> completed`

## Related

- `/set-task-status` - Update task status
- `/run-tasks` - Auto-execute all tasks
- `/execute-task` - Manual single task workflow
