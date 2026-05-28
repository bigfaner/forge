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

Developer consolidates hardcoded validation maps in `validate_index.go` to use `types.AllStatuses()` and `types.AllPriorities()` helper functions, ensuring validation logic stays synchronized with constant definitions (Story 4).

## Setup

- `pkg/types/` package exists with Status, SurfaceType, Priority types and helper functions (`AllStatuses()`, `AllPriorities()`)
- `internal/cmd/validate_index.go` contains hardcoded `validStatus` and `validPriority` maps
- All Status and Priority magic values already migrated to typed constants (prerequisite: journeys 1 and status-migration)

## Happy Path

### Step 1: Identify Hardcoded Validation Maps

**User Action**: Developer locates `validStatus map[string]bool{"pending": true, "in-progress": true, ...}` and similar `validPriority` map in `validate_index.go`

**Expected Result**: All hardcoded validation maps identified. Count of hardcoded entries matches expected enum count (7 Status, 3 Priority).

### Step 2: Replace Status Validation with AllStatuses()

**User Action**: Developer replaces `validStatus` hardcoded map with a map built from `types.AllStatuses()`, e.g., `buildValidMap(types.AllStatuses())`

**Expected Result**: Validation logic uses `types.AllStatuses()` as single source of truth. Any future Status constant added to `pkg/types/` is automatically included in validation.

### Step 3: Replace Priority Validation with AllPriorities()

**User Action**: Developer replaces `validPriority` hardcoded map with a map built from `types.AllPriorities()`

**Expected Result**: Validation logic uses `types.AllPriorities()` as single source of truth. Same auto-sync benefit as Status.

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

## Journey Invariants

- Validation maps must be derived from `types.AllStatuses()`/`types.AllPriorities()` — no hardcoded enum values in validation logic
- Adding a new enum constant must not require changes to validation code
- Validation behavior must be identical before and after consolidation
- All validation errors must produce clear, actionable error messages on stderr
