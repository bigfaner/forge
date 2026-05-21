---
step: 1
title: Quality Gate Fix-Task Loop Breaker
journey: quality-gate
---

# Step 1: Quality Gate Fix-Task Loop Breaker

## Given
- A forge project with tasks in index.json where all tasks are completed
- A justfile with compile/fmt/lint/test recipes
- Some recipes may fail (compile, fmt, lint)

## When
- `forge quality-gate` is executed after all tasks are completed

## Then
- Fix tasks are created for failing steps with step-scoped SourceTaskID (quality-gate:<step>)
- Cumulative cap of 3 fix tasks per step is enforced
- Cross-step independence: cap on compile does not block lint fix tasks
- Docs-only features are skipped
- Fix task markdown files are created on disk with proper template sections

## Contract Dimensions
- **Actor**: CLI user running quality gate after task completion
- **Input**: index.json with completed tasks, justfile with compile/fmt/lint/test recipes
- **Output**: CLI stdout/stderr with check results, fix task creation, cap warnings
- **Side Effects**: index.json updated with fix tasks, fix task markdown files created
- **Error Cases**: cap exceeded -> warning message, no new fix task created
- **Invariants**: SourceTaskID uses step-scoped sentinel format; cumulative count across all statuses
