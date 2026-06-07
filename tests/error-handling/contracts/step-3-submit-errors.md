# Contract: error-handling / Step 3: Submit Errors

## Outcome "submit-concurrent-conflict"
- Preconditions: "concurrent lock contention on index.json"
- Input: `forge task submit <task-id>` during concurrent write
- Output: "retry error indicating write conflict"
- State: "no state changes (original index preserved)"
- Side-effect: none

## Outcome "submit-missing-index"
- Preconditions: "feature directory without index.json"
- Input: `forge task submit <task-id>`
- Output: "error indicating index not found"
- State: "no state changes"
- Side-effect: none

## Outcome "verify-incomplete-tasks"
- Preconditions: "feature with incomplete tasks and active state.json"
- Input: `forge task verify-done`
- Output: "error indicating not all tasks completed"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- submit is atomic — either fully succeeds or fully rolls back
- error messages guide the user toward resolution
