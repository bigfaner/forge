# Contract: error-handling / Step 1: Task Errors

## Outcome "task-status-not-found"
- Preconditions: "task ID does not exist in any feature index"
- Input: `forge task status <nonexistent-id>`
- Output: "exit code 1, output containing 'not found'"
- State: "no state changes"
- Side-effect: none

## Outcome "task-claim-no-available"
- Preconditions: "feature exists with no available pending tasks"
- Input: `forge task claim`
- Output: "error indicating no tasks available"
- State: "no state changes"
- Side-effect: none

## Outcome "task-claim-corrupted-index"
- Preconditions: "corrupted or missing index.json"
- Input: `forge task claim`
- Output: "error indicating index load failure"
- State: "no state changes"
- Side-effect: none

## Outcome "task-validate-invalid-schema"
- Preconditions: "index.json with schema validation errors"
- Input: `forge task validate-index <path>`
- Output: "non-zero exit code with validation error details"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- error messages are user-readable and actionable
- non-zero exit codes for all error conditions
