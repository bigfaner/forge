# Task Executor Auto-Claims Next Task

## Problem

When dispatching `zcode:task-executor` via the `/run-tasks` dispatcher, the executor automatically calls `zcode:claim-task` after completing its assigned task and starts executing the next task — bypassing the dispatcher's control loop entirely.

## Root Cause

The `zcode:task-executor` skill (or its internal prompt) includes a step to claim the next task after recording the current one. This is designed for standalone use, but conflicts with the `/run-tasks` dispatcher architecture, which expects to be the sole entity claiming tasks.

## Solution

When dispatching `zcode:task-executor` from `/run-tasks`, explicitly instruct it **not** to claim or start any subsequent tasks:

```
Agent(
  subagent_type="task-executor",
  prompt="TASK_KEY: {{KEY}}
TASK_ID: {{ID}}
TASK_FILE: {{FILE}}

IMPORTANT: Do NOT claim or start any other tasks after completing this one. Stop after recording the task result."
)
```

## Key Takeaway

`zcode:task-executor` has autonomous task-chaining behavior by default. The `/run-tasks` dispatcher must explicitly suppress this with a "stop after one task" instruction, otherwise the dispatcher loses control of task ordering, parallelism, and error handling.
