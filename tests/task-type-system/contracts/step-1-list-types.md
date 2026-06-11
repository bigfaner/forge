# Contract: task-type-system / Step 1: List Types

## Outcome "list-types-output"
- Preconditions: "forge CLI binary built from current source tree"
- Input: `forge task list-types`
- Output: "all registered task types with descriptions, one per line, at least 5 types, each description <= 60 chars, exit code 0"
- State: "no state changes"
- Side-effect: none

## Outcome "list-types-dispatch-prompt"
- Preconditions: "feature index with typed tasks exists"
- Input: `forge prompt get-by-task-id <task-id>`
- Output: "synthesized prompt matching the task type template, exit code 0"
- State: "no state changes"
- Side-effect: none

## Outcome "list-types-validate"
- Preconditions: "index.json file with task entries"
- Input: `forge task validate-index <path>`
- Output: "valid types pass, invalid or missing types fail with type-related error message"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- forge binary built once per test session via TestMain
- testkit.SetForgeBinary propagated to all test helpers
