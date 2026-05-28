---
feature: "unify-enum-constants"
journey: "full-verification"
risk_level: "Low"
surface_types: ["cli"]
sources:
  - docs/features/unify-enum-constants/prd/prd-user-stories.md
  - docs/features/unify-enum-constants/prd/prd-spec.md
generated: "2026-05-29"
---

# Journey: full-verification

**Risk Level**: Low

<!-- Risk Classification Criteria:
  High   = Workflow involves state mutation, data loss risk, or irreversible operations
  Medium = Workflow involves multi-step interaction without irreversible side effects
  Low    = Workflow is read-only or purely observational
-->

## Overview

Maintainer verifies that the complete enum constants migration is correct: zero magic values remain, all code compiles, all tests pass, and `pkg/types/` is a proper leaf package (Story 3).

## Setup

- All enum migration journeys (status-migration, surface-type-migration, validation-map-consolidation) are completed
- `pkg/types/` package contains all typed constant definitions
- All production code uses typed constants instead of string literals

## Happy Path

### Step 1: Verify Zero Magic Values in Production Code

**User Action**: Maintainer runs `grep -r '"pending"\|"in-progress"\|"completed"\|"blocked"\|"cancelled"\|"failed"\|"review"' --include='*.go'` excluding test files

**Expected Result**: Zero matches in production code. All Status values are typed constants. Exit code 1 (grep found nothing).

### Step 2: Verify Leaf Package Property

**User Action**: Maintainer inspects `pkg/types/` imports — no forge-cli internal packages imported

**Expected Result**: `pkg/types/` only imports standard library packages. It is a true leaf package with no circular dependency risk.

### Step 3: Run Full Build

**User Action**: Maintainer runs `go build ./...`

**Expected Result**: Zero compilation errors. All type signatures are consistent across packages. Exit code 0.

### Step 4: Run Full Test Suite

**User Action**: Maintainer runs `go test ./...`

**Expected Result**: All tests pass. Zero behavior change from migration. Test output shows no failures.

### Step 5: Verify CLI Behavior Unchanged

**User Action**: Maintainer runs sample CLI commands (`forge task list`, `forge task status <id>`, etc.) and compares output to pre-migration baseline

**Expected Result**: CLI output identical to pre-migration. No user-visible changes.

## Edge Cases

### Step 1b: Residual Magic Value Found

**Precondition**: A grep match is found in a production file

**User Action**: Maintainer investigates the match location

**Expected Result**: The match is either in a test file (expected) or is a genuine missed migration that needs to be fixed before release.

### Step 2b: Indirect Circular Dependency

**Precondition**: `pkg/types/` imports a package that transitively imports `pkg/types/`

**User Action**: Maintainer runs `go build ./pkg/types/`

**Expected Result**: Go compiler detects and reports the circular dependency. Must be resolved before merge.

### Step 4b: Test Failure Due to Type Assertion

**Precondition**: A test uses `assert.Equal(t, "pending", task.Status)` with untyped string

**User Action**: Maintainer runs the failing test

**Expected Result**: Test fails because `types.StatusPending` is not equal to untyped `"pending"` in strict equality. Test must be updated to use typed constant or the comparison must account for the string underlying type.

## Journey Invariants

- Zero Status/SurfaceType/Priority string literals in production code
- `pkg/types/` is a leaf package with no forge-cli internal imports
- `go build ./...` passes with zero errors
- `go test ./...` passes with zero failures
- CLI output is identical before and after migration (zero behavior change)
