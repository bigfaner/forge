# Contract: skill-ops / Step 3: Prompt

## Outcome "prompt-get-by-task-id"
- Preconditions: "valid task ID exists in feature index"
- Input: `forge prompt get-by-task-id <task-id>`
- Output: "exit code 0, non-empty output with TASK_ID and TASK_FILE substituted (no template placeholders)"
- State: "no state changes"
- Side-effect: none

## Outcome "prompt-nonexistent-task-id"
- Preconditions: "no task with given ID"
- Input: `forge prompt get-by-task-id NONEXISTENT-999`
- Output: "exit code 1, output contains 'not found'"
- State: "no state changes"
- Side-effect: none

## Outcome "prompt-invalid-type"
- Preconditions: "task with missing/invalid type in index.json"
- Input: `forge prompt get-by-task-id <id>`
- Output: "appropriate error message"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- forge binary path consistent across all steps
- all commands use built binary, not system-installed
