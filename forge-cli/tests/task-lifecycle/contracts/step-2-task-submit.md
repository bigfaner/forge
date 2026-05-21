# Contract: task-lifecycle / Step 2: Task Submit

## Outcome "submit-success"
- Preconditions: "task exists in in_progress state with valid test data"
- Input: `forge task submit <task-id>` with required result data flags
- Output: "success confirmation with task status update, exit code 0"
- State: "task status changes to completed, record file created in tasks/records/"
- Side-effect: none

## Outcome "submit-terminal-state"
- Preconditions: "task exists in a terminal state (completed, blocked, rejected, or skipped)"
- Input: `forge task submit <task-id>` with result data flags
- Output: "error message indicating task is already in terminal state, exit code non-zero"
- State: "no state changes"
- Side-effect: none

## Outcome "submit-missing-data"
- Preconditions: "task exists but no result data provided via --data flag or stdin"
- Input: `forge task submit <task-id>` without data flags
- Output: "error message indicating missing input, exit code 1"
- State: "no state changes"
- Side-effect: none

## Outcome "concurrent-lock-contention"
- Preconditions: "task in in_progress state, concurrent submit process active"
- Input: `forge task submit <task-id>` while another submit is in progress
- Output: "lock contention handled gracefully, either queued success or error with retry suggestion"
- State: "first submit wins, second receives lock error or queued result"
- Side-effect: none

## Journey Invariants
- feature_slug consistent across all steps
- task_id stable once assigned
- index.json remains valid JSON throughout
