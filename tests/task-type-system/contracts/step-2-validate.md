# Contract: task-type-system / Step 2: Validate Index

## Outcome "validate-valid-types"
- Preconditions: "index.json with recognized task types (feature, enhancement, cleanup, refactor)"
- Input: `forge task validate-index <path-to-valid-index>`
- Output: "exit code 0, no validation errors"
- State: "no state changes"
- Side-effect: none

## Outcome "validate-invalid-type"
- Preconditions: "index.json with unrecognized task type"
- Input: `forge task validate-index <path-to-invalid-index>`
- Output: "exit code non-zero, error mentioning unknown/invalid type"
- State: "no state changes"
- Side-effect: none

## Outcome "validate-missing-type"
- Preconditions: "index.json with task entry missing type field"
- Input: `forge task validate-index <path-to-missing-index>`
- Output: "exit code non-zero, error mentioning missing/required type"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- index.json format is valid JSON with tasks map
- type field is mandatory for all task entries
