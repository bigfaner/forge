# Contract: task-lifecycle / Step 1: Task Claim

## Outcome "claim-success"
- Preconditions: "feature exists with at least one pending task with met dependencies"
- Input: `forge task claim` with CLAUDE_PROJECT_DIR pointing to a valid project root
- Output: "block output with ACTION: CLAIMED, TASK_ID, and task metadata fields, exit code 0"
- State: "claimed task status changes to in_progress, state.json created with task assignment"
- Side-effect: none

## Outcome "fix-task-priority"
- Preconditions: "pending fix task with sourceTaskID blocking a dependent business task"
- Input: `forge task claim` with CLAUDE_PROJECT_DIR pointing to a project with fix tasks
- Output: "block output with ACTION: CLAIMED and TASK_ID set to the fix task, exit code 0"
- State: "fix task status changes to in_progress, dependent business task remains blocked"
- Side-effect: none

## Outcome "fix-chain"
- Preconditions: "multiple pending fix tasks with same sourceTaskID, dependent business task"
- Input: `forge task claim` in sequence as fix tasks are completed
- Output: "each claim returns the next pending fix task until all complete, then business task becomes available"
- State: "fix tasks complete one by one, business task unblocked after last fix task completes"
- Side-effect: none

## Outcome "no-tasks-available"
- Preconditions: "feature exists with no pending tasks or all tasks have unmet dependencies"
- Input: `forge task claim` with CLAUDE_PROJECT_DIR pointing to a project with no eligible tasks
- Output: "message indicating no tasks available, exit code non-zero or specific no-task message"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- feature_slug consistent across all steps
- task_id stable once assigned
- index.json remains valid JSON throughout
