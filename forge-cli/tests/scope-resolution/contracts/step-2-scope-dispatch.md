# Contract: scope-resolution / Step 2: Scope Dispatch

## Outcome "scope-mismatch-warning"
- Preconditions: "justfile with scope-aware build recipe exists"
- Input: `just build invalidscope`
- Output: "error for invalid scope, or success for scope-agnostic justfile"
- State: "no state changes"
- Side-effect: none

## Outcome "matching-scope-executes"
- Preconditions: "forge project justfile with compile recipe"
- Input: `just compile frontend`
- Output: "no scope error, may fail due to missing toolchain (acceptable)"
- State: "no state changes"
- Side-effect: none

## Outcome "project-type-fallback"
- Preconditions: "forge probe command or just project-type recipe available"
- Input: "forge probe or just project-type"
- Output: "valid project type (frontend/backend/mixed/go), or skip if unavailable"
- State: "no state changes"
- Side-effect: none

## Outcome "project-type-deterministic"
- Preconditions: "forge probe or just project-type available"
- Input: "forge probe called twice"
- Output: "identical output on both calls"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- PRD spec at docs/features/justfile-standard-vocabulary/prd/prd-spec.md documents fallback behavior
- project-type output is deterministic across runs
