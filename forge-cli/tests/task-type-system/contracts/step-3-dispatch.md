# Contract: task-type-system / Step 3: Type Dispatch

## Outcome "dispatch-doc-generation"
- Preconditions: "task with type doc-generation exists in feature index"
- Input: `forge prompt get-by-task-id <doc-task-id>`
- Output: "prompt without TDD steps (no RED/GREEN/REFACTOR, no 'just test')"
- State: "no state changes"
- Side-effect: none

## Outcome "dispatch-fix-type"
- Preconditions: "task with type fix exists in feature index"
- Input: `forge prompt get-by-task-id <fix-task-id>`
- Output: "prompt containing five-step diagnostic flow: diagnose, locate, fix, verify, commit"
- State: "no state changes"
- Side-effect: none

## Outcome "dispatch-unregistered-type"
- Preconditions: "index.json with task having unregistered type"
- Input: `forge prompt get-by-task-id <invalid-task-id>`
- Output: "non-zero exit code, error mentioning unknown/invalid type"
- State: "no state changes"
- Side-effect: none

## Outcome "dispatch-prompt-idempotent"
- Preconditions: "valid task exists in feature index"
- Input: `forge prompt get-by-task-id <task-id>` called twice
- Output: "both calls produce identical output, each completes within 500ms"
- State: "no state changes"
- Side-effect: none

## Outcome "dispatch-nonexistent-task"
- Preconditions: "task ID does not exist in any feature index"
- Input: `forge prompt get-by-task-id <nonexistent-id>`
- Output: "non-zero exit code, error mentioning 'not found' or 'error'"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- prompt output is deterministic for same task ID
- type system dispatches to correct template based on task type field
