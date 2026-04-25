---
name: set-task-status
description: Use when you need to update the status of a task. Supports pending, in_progress, completed, blocked, skipped statuses.
---

# Set Task Status

## Overview

Update the status of a task in `docs/features/<slug>/tasks/index.json`.

## Usage

```bash
task status <task-id-or-key> <status>
```

## Valid Statuses

| Status        | When to Use                          |
| ------------- | ------------------------------------ |
| `pending`     | Task not started (default)           |
| `in_progress` | Currently working on it              |
| `completed`   | Task finished successfully           |
| `blocked`     | Cannot proceed due to external issue |
| `skipped`     | Task not needed                      |

## Workflow

```
pending → in_progress → completed
                 ↓
              blocked → in_progress
                 ↓
              skipped
```

## Related Skills

- `/claim-task` - Claim next available task
- `/run-tasks` - Auto-execute all tasks
