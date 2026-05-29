---
feature: "unify-enum-constants"
journey: "surface-type-migration"
risk_level: "High"
surface_types: ["cli"]
sources:
  - docs/features/unify-enum-constants/prd/prd-user-stories.md
  - docs/features/unify-enum-constants/prd/prd-spec.md
generated: "2026-05-29"
---

# Journey: surface-type-migration

**Risk Level**: High

<!-- Risk Classification Criteria:
  High   = Workflow involves state mutation, data loss risk, or irreversible operations
  Medium = Workflow involves multi-step interaction without irreversible side effects
  Low    = Workflow is read-only or purely observational
-->

## Overview

Developer migrates all Surface Type string literals (~97 occurrences across 6 files) to typed constants in `pkg/types/`, ensuring compile-time type safety and single-point-of-definition for surface types (Story 2).

## Setup

- `pkg/types/surface.go` exists with `type SurfaceType string` and 5 typed constants (`SurfaceWeb`, `SurfaceAPI`, `SurfaceCLI`, `SurfaceTUI`, `SurfaceMobile`)
- `AllSurfaceTypes()` helper function returns all 5 constants
- A working forge-cli codebase with ~97 Surface Type string literals across 6 files

## Happy Path

### Step 1: Define SurfaceType and Constants

**User Action**: Developer creates `pkg/types/surface.go` with `type SurfaceType string` and 5 constants plus `AllSurfaceTypes()` helper

**Expected Result**: File compiles. Constants are of type `SurfaceType`. `AllSurfaceTypes()` returns a slice containing all 5 values.

### Step 2: Replace Surface Type Magic Values in Detection Module

**User Action**: Developer replaces Surface Type string literals in `pkg/forgeconfig/detect_surface.go` with `types.SurfaceXxx` constants (mapping tables, comparison expressions)

**Expected Result**: Surface detection logic uses typed constants. Function signatures upgraded to use `types.SurfaceType`. `go build ./pkg/forgeconfig/` passes.

### Step 3: Replace Surface Type Magic Values in Remaining Files

**User Action**: Developer replaces Surface Type string literals in remaining files (`execution_order.go`, `detect.go`, and other forgeconfig files)

**Expected Result**: All ~97 Surface Type string literals replaced. Struct fields use `types.SurfaceType`. `go build ./...` passes.

### Step 4: Add Boundary Conversions

**User Action**: Developer adds `types.SurfaceType(stringVal)` conversions at config parsing and surface detection output boundaries

**Expected Result**: External string input is explicitly converted to `types.SurfaceType` at boundaries. Internal code uses typed constants exclusively.

## Edge Cases

### Step 1b: Adding a New Surface Type

**Precondition**: A new Surface Type (e.g., "desktop") needs to be added to the system

**User Action**: Developer adds `SurfaceDesktop SurfaceType = "desktop"` to `surface.go` and updates `AllSurfaceTypes()`

**Expected Result**: Only `surface.go` modified. All references automatically pick up the new type through `AllSurfaceTypes()`. Detection maps can be updated in one place.

### Step 2b: Surface Detection Returns Unknown Value

**Precondition**: `forge surfaces <path>` encounters a project with no recognizable surface signals

**User Action**: Developer runs `forge surfaces ./unknown-project/`

**Expected Result**: Exit code 1, stderr contains error message with configuration guidance. The error path does not crash — it handles the unknown surface gracefully.

### Step 3b: Config File Contains Invalid Surface Type

**Precondition**: A config file specifies `surface: "desktop"` which is not a defined constant

**User Action**: Developer parses the config file during surface detection

**Expected Result**: Config parsing either rejects the invalid value (if strict validation) or the unknown value is handled as an error. The typed constant system ensures this is caught at the validation boundary.

### Step 4b: Cross-Package Surface Type Comparison

**Precondition**: Two packages compare Surface Type values using different representations (one typed, one string)

**User Action**: Developer runs `go build ./...` after partial migration

**Expected Result**: Type mismatch error identifies the inconsistent comparison. Forces completion of migration in the affected package.

### Step 4c: JSON/Config Serialization Compatibility

**Precondition**: SurfaceType struct fields changed from `string` to `types.SurfaceType`

**User Action**: Developer reads/writes a config file containing surface type values

**Expected Result**: Config serialization/deserialization produces identical output. `type SurfaceType string` is transparent in JSON/YAML.

### Step 5: CLI Error — Surface Not Found

**Precondition**: CLI command references a surface type that does not match any typed constant (e.g., `forge surfaces ./project/` returns no detectable surface)

**User Action**: Developer runs `forge surfaces ./unknown-project/`

**Expected Result**: Exit code 1. stderr contains error message indicating no surface detected, with configuration guidance. No crash.

### Step 6: CLI Error — Duplicate Surface Type in Config

**Precondition**: `.forge/config.yaml` defines a surface key that resolves to a type already configured

**User Action**: Developer runs `forge surfaces detect`

**Expected Result**: The already-existing surface type is reported. No duplicate entries created. If validation enforces uniqueness, error is reported on stderr.

## Journey Invariants

- `pkg/types/` must never import any forge-cli internal package — it is a leaf package
- All typed constant values must exactly match the original string literals (zero behavior change)
- Every SurfaceType field assignment in production code must use `types.SurfaceXxx` constants
- `AllSurfaceTypes()` must return exactly the set of defined constants — no more, no less
- `go build ./...` must pass after each migration step
