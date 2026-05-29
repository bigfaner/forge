---
feature: "unify-enum-constants"
journey: "validation-map-consolidation"
risk_level: "Medium"
surface_types: ["cli"]
sources:
  - docs/features/unify-enum-constants/prd/prd-user-stories.md
  - docs/features/unify-enum-constants/prd/prd-spec.md
generated: "2026-05-29"
---

# Journey: validation-map-consolidation

**Risk Level**: Medium

<!-- Risk Classification Criteria:
  High   = Workflow involves state mutation, data loss risk, or irreversible operations
  Medium = Workflow involves multi-step interaction without irreversible side effects
  Low    = Workflow is read-only or purely observational
-->

## Overview

Developer verifies that validation maps in `validate_index.go` correctly use `types.AllStatuses()` and `types.AllPriorities()` helper functions, ensuring validation logic stays synchronized with constant definitions (Story 4). This journey validates already-completed migration work.

## Setup

- `pkg/types/` package exists with Status, SurfaceType, Priority types and helper functions (`AllStatuses()`, `AllPriorities()`)
- `internal/cmd/task/validate_index.go` uses `buildValidStatusMap()` and `buildValidPriorityMap()` which call `types.AllStatuses()`/`types.AllPriorities()`
- All Status and Priority magic values already migrated to typed constants (prerequisite: status-migration and surface-type-migration journeys)

## Happy Path

### Step 1: Verify Validation Maps Use Typed Helper Functions

**User Action**: Developer inspects `internal/cmd/task/validate_index.go` to confirm `buildValidStatusMap()` calls `types.AllStatuses()` and `buildValidPriorityMap()` calls `types.AllPriorities()`

**Expected Result**: No hardcoded `validStatus`/`validPriority` maps exist. Validation maps are built from `types.AllStatuses()`/`types.AllPriorities()`. Adding a new enum constant automatically includes it in validation.

### Step 2: Verify AllStatuses() Coverage

**User Action**: Developer runs `go test ./pkg/types/... -run TestAllStatuses` to confirm `AllStatuses()` returns all 7 status constants

**Expected Result**: Test passes. `AllStatuses()` returns exactly `[StatusPending, StatusInProgress, StatusCompleted, StatusBlocked, StatusSuspended, StatusSkipped, StatusRejected]`.

### Step 3: Verify AllPriorities() Coverage

**User Action**: Developer runs `go test ./pkg/types/... -run TestAllPriorities` to confirm `AllPriorities()` returns all 3 priority constants

**Expected Result**: Test passes. `AllPriorities()` returns exactly `[PriorityP0, PriorityP1, PriorityP2]`.

### Step 4: Verify Validation Behavior Unchanged

**User Action**: Developer runs `go test ./internal/cmd/...` to verify validation still accepts valid values and rejects invalid ones

**Expected Result**: All tests pass. Valid Status/Priority values accepted, invalid ones rejected. Behavior identical to pre-migration.

## Edge Cases

### Step 2b: Empty AllStatuses() Return

**Precondition**: `AllStatuses()` returns an empty slice due to a bug in the implementation

**User Action**: Developer runs validation on a task with valid status

**Expected Result**: All statuses rejected. This should be caught in unit tests for `AllStatuses()`.

### Step 3b: Duplicate Values in Helper Output

**Precondition**: `AllPriorities()` accidentally returns duplicates

**User Action**: Developer builds validation map from `AllPriorities()`

**Expected Result**: Map construction deduplicates automatically (Go map key semantics). No functional impact, but signals a bug in the helper function.

### Step 4b: Validation of Case Sensitivity

**Precondition**: Config file contains `"Pending"` (capitalized) instead of `"pending"`

**User Action**: Developer runs `forge task validate` on the config

**Expected Result**: Validation rejects the incorrect casing. Typed constants enforce exact string matching. Error message clearly indicates the invalid value.

### Step 5: CLI Error — Task Status Not Found

**Precondition**: `forge task validate` is run on an index with a task whose status is not in `AllStatuses()` output (e.g., a typo like `"pendng"`)

**User Action**: Developer runs `forge task validate` on the malformed index

**Expected Result**: Exit code non-zero. stderr lists the invalid status value with guidance on valid values. No crash.

### Step 6: CLI Error — Priority Already Exists

**Precondition**: Task index contains two tasks with the same ID but different priorities (conflicting `validPriority` lookup)

**User Action**: Developer runs `forge task validate` on the conflicting index

**Expected Result**: Validation reports the duplicate/conflicting entry. Error message identifies the task ID and both priority values.

## Journey Invariants

- Validation maps must be derived from `types.AllStatuses()`/`types.AllPriorities()` — no hardcoded enum values in validation logic
- Adding a new enum constant must not require changes to validation code
- Validation behavior must be identical before and after consolidation
- All validation errors must produce clear, actionable error messages on stderr
