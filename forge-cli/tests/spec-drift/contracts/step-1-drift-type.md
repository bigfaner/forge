# Contract: spec-drift / Step 1: Drift Type

## Outcome "drift-type-in-list-types"
- Preconditions: "forge binary built from source"
- Input: "Run forge task list-types"
- Output: "Output contains doc-generation.drift type with correct description"
- State: "no state changes"
- Side-effect: none

## Outcome "drift-strategy-template-exists"
- Preconditions: "pkg/prompt/data/doc-generation-drift.md exists"
- Input: "Read template file"
- Output: "Template references consolidate-specs skill and Steps 9-11 (drift-only mode)"
- State: "no state changes"
- Side-effect: none

## Outcome "drift-type-valid-in-validate-index"
- Preconditions: "feature context with index.json containing drift tasks"
- Input: "Run forge task validate-index"
- Output: "No unknown type errors for doc-generation.drift"
- State: "no state changes"
- Side-effect: none

## Outcome "quick-pipeline-includes-drift-task"
- Preconditions: "pkg/task/testgen.go exists"
- Input: "Read testgen.go content"
- Output: "T-quick-doc-drift defined with TypeDocDrift, Scope all, correct title and dependency on T-quick-verify-regression"
- State: "no state changes"
- Side-effect: none

## Outcome "drift-task-ordering"
- Preconditions: "pkg/task/testgen.go exists"
- Input: "Read testgen.go content"
- Output: "T-quick-verify-regression appears before T-quick-doc-drift"
- State: "no state changes"
- Side-effect: none

## Outcome "task-id-infers-drift-type"
- Preconditions: "pkg/task/infer.go exists"
- Input: "Read infer.go content"
- Output: "T-quick-doc-drift infers TypeDocDrift via profileSuffixedID"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- doc-generation.drift is a valid task type in the forge CLI type registry
- T-quick-doc-drift is the quick-pipeline task for drift detection
- Drift type maps to doc-generation-drift.md strategy template
