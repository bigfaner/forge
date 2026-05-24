---
status: "completed"
started: "2026-05-24 15:54"
completed: "2026-05-24 16:23"
time_spent: "~29m"
---

# Task Record: 6 forge init: TUI confirmation for surfaces

## Summary
Added TUI confirmation flow for surfaces in forge init: scalar display shows single type, map form shows path->surface rows with conflict annotations, edit/add/delete operations, and confirm-to-write flow.

## Changes

### Files Created
- forge-cli/internal/cmd/init_surfaces.go
- forge-cli/internal/cmd/init_surfaces_test.go

### Files Modified
- forge-cli/pkg/forgeconfig/detect_surface.go
- forge-cli/internal/cmd/init.go
- forge-cli/internal/cmd/init_test.go

### Key Decisions
- Refactored detect_surface.go to expose per-manifest signal lists (detectPackageJSONSignals etc.) instead of just resolved types, enabling conflict metadata at the TUI layer
- Added DetectResult type with IsScalar flag and PathConflict slice to carry conflict info from detection to TUI
- TUI uses separate code paths for scalar (askScalarConfirmation) and map (askMapConfirmation) forms per Hard Rule that single-type should not show path column
- surfaceConfigFunc variable for testability mirrors existing configInitFunc pattern

## Test Results
- **Tests Executed**: Yes
- **Passed**: 13
- **Failed**: 0
- **Coverage**: 62.0%

## Acceptance Criteria
- [x] Detected surfaces displayed in TUI: scalar form shows single type, map form shows path->surface rows
- [x] Confirm button (or Enter) writes surfaces to config and proceeds
- [x] Edit entry: each row can enter edit mode to modify path or surface type
- [x] Conflict annotation: conflicting signals shown as path: surface (冲突信号: web + api，已按优先级选择 web)
- [x] Add: blank row input to add new mapping
- [x] Delete: select row + press d to remove

## Notes
Hard Rules verified: TUI works in both scalar and map display modes; single-type detection does NOT show path column. Existing 28 detection tests + all init tests pass unchanged.
