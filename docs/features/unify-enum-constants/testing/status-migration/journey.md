---
feature: "unify-enum-constants"
journey: "status-migration"
risk_level: "High"
surface_types: ["cli"]
sources:
  - docs/features/unify-enum-constants/prd/prd-user-stories.md
  - docs/features/unify-enum-constants/prd/prd-spec.md
generated: "2026-05-29"
---

# Journey: status-migration

**Risk Level**: High

<!-- Risk Classification Criteria:
  High   = Workflow involves state mutation, data loss risk, or irreversible operations
  Medium = Workflow involves multi-step interaction without irreversible side effects
  Low    = Workflow is read-only or purely observational
-->

## Overview

Developer migrates all Status string literals across the codebase to typed constants in `pkg/types/`, ensuring compile-time type safety for task status values (Story 1 + Story 3 partial).

## Setup

- `pkg/types/status.go` exists with `type Status string` and 7 typed constants defined
- `pkg/feature/constants.go` re-exports Status constants from `pkg/types/` for backward compatibility
- A working forge-cli codebase with 117 Status string literals across 22 files

## Happy Path

### Step 1: Define Status Type and Constants

**User Action**: Developer creates `pkg/types/status.go` with `type Status string` and 7 constants (`StatusPending`, `StatusInProgress`, `StatusCompleted`, `StatusBlocked`, `StatusCancelled`, `StatusFailed`, `StatusReview`)

**Expected Result**: File compiles, `go build ./pkg/types/` succeeds. Constants are of type `Status` (not untyped string).

### Step 2: Migrate Existing Status Constants

**User Action**: Developer moves Status constant definitions from `pkg/feature/constants.go` to `pkg/types/status.go`, and adds re-export in the original file (`StatusPending = types.StatusPending`)

**Expected Result**: `pkg/feature/constants.go` still exports the same constants via re-export. Existing consumers compile without changes. `go build ./...` passes.

### Step 3: Replace Status Magic Values in State Machine

**User Action**: Developer replaces all Status string literals in `statemachine.go` with `types.StatusXxx` constants (transition table `From`/`To` fields, state checks)

**Expected Result**: State machine type signatures use `types.Status`. Compile-time verification catches any invalid status values. `go build ./pkg/task/` passes.

### Step 4: Replace Status Magic Values Across Remaining Files

**User Action**: Developer replaces Status string literals in all remaining files (`add.go`, `state.go`, `deps.go`, `build.go`, `autogen.go`, `types.go`, etc.)

**Expected Result**: All 117 Status string literals replaced with typed constants. Function signatures upgraded to use `types.Status`. `go build ./...` passes with zero errors.

### Step 5: Add Boundary Conversions at CLI Flags

**User Action**: Developer adds `types.Status(stringVal)` conversions at CLI flag parsing boundaries and config deserialization entry points

**Expected Result**: External string input (CLI args, config files) is explicitly converted to `types.Status` at the boundary. Internal code uses typed constants exclusively.

## Edge Cases

### Step 1b: Missing Status Constant

**Precondition**: A Status value exists in the codebase but is not defined as a constant in `pkg/types/status.go`

**User Action**: Developer runs `go build ./...` after partial migration

**Expected Result**: Compiler reports type mismatch error for the missing constant. Exit code non-zero. Error message clearly identifies the file and line.

### Step 2b: Circular Import via pkg/types

**Precondition**: `pkg/types/` accidentally imports an internal forge-cli package

**User Action**: Developer runs `go build ./pkg/types/`

**Expected Result**: Go compiler reports circular import error. `pkg/types/` must be a leaf package with no internal imports.

### Step 3b: Incorrect State Transition After Migration

**Precondition**: A string literal was incorrectly replaced with the wrong typed constant (e.g., `"completed"` replaced with `types.StatusFailed`)

**User Action**: Developer runs `go test ./pkg/task/...`

**Expected Result**: Test failure in state machine tests. The typed constant makes the error visible in code review even before test execution.

### Step 4b: Re-export Compatibility Break

**Precondition**: `pkg/feature/constants.go` re-export is missing or uses wrong type

**User Action**: Developer runs `go build ./...` to check for downstream breakage

**Expected Result**: External consumers of `pkg/feature/constants.go` fail to compile. Error identifies the missing/broken re-export.

### Step 5b: JSON Serialization Compatibility

**Precondition**: Struct fields changed from `string` to `types.Status`

**User Action**: Developer serializes/deserializes a task with Status field via JSON

**Expected Result**: JSON output unchanged — `type Status string` marshals/unmarshals identically to plain `string`. No behavior change.

### Step 5c: Status Comparison with Untyped String

**Precondition**: Developer writes `task.Status == "pending"` (untyped string comparison)

**User Action**: Developer runs `go build`

**Expected Result**: Compiler allows the comparison (untyped string constant is assignable to `types.Status`), but linter/review should flag it. The typed constant form `task.Status == types.StatusPending` is preferred.

## Journey Invariants

- `pkg/types/` must never import any forge-cli internal package — it is a leaf package
- All typed constant values must exactly match the original string literals (zero behavior change)
- Every Status field assignment in production code must use `types.StatusXxx` constants, never raw string literals
- `go build ./...` must pass after each migration step (incremental correctness)
- JSON serialization/deserialization of Status fields must produce identical output before and after migration
