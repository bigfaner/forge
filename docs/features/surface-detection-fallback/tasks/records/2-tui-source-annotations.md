---
status: "completed"
started: "2026-05-24 20:55"
completed: "2026-05-24 21:13"
time_spent: "~18m"
---

# Task Record: 2 Update TUI confirmation flow with source annotations and re-run prompt

## Summary
Updated TUI confirmation flow with source annotations and re-run prompt. Added formatSourceAnnotation, isInferred, formatSurfacesSummary, askRerunPrompt. Threaded SourcesMap through buildDisplayLines, askScalarConfirmation, askMapConfirmation. Added re-run behavior in runSurfaceConfig with Confirm/Re-detect/Edit options using shared manualSurfaceEntry function.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/init_surfaces.go
- forge-cli/internal/cmd/init_surfaces_test.go
- forge-cli/internal/cmd/init.go

### Key Decisions
- Source annotation is display-only via formatSourceAnnotation -- not persisted to SurfacesMap (Hard Rule)
- Re-run Edit calls the same manualSurfaceEntry variable -- no separate code path (Hard Rule)
- Made manualSurfaceEntry, askRerunPrompt, runNewSurfaceDetection into function variables for testability
- askSurfaceConfirmation preserved as entry point for runNewSurfaceDetectionImpl (used by task 5 forge surfaces detect)
- TUI hint text shown only for inferred surfaces, not dependency-detected surfaces

## Test Results
- **Tests Executed**: Yes
- **Passed**: 47
- **Failed**: 0
- **Coverage**: 60.1%

## Acceptance Criteria
- [x] askScalarConfirmation shows source annotation in description field
- [x] askMapConfirmation shows per-path source annotation in display lines
- [x] TUI hint text for inferred surfaces, absent for dependency-detected
- [x] Re-run behavior: existing surfaces triggers Confirm/Re-detect/Edit prompt
- [x] Confirm returns SKIPPED surfaces (already configured)
- [x] Re-detect runs full detection + inference pipeline
- [x] Edit calls manualSurfaceEntry (same function as first-run)
- [x] User override: edited value has no source annotation persisted
- [x] Multi-surface TUI display shows map-form list with source annotations

## Notes
Core logic functions (formatSourceAnnotation, formatInferenceDetail, isInferred, formatSurfacesSummary, buildDisplayLines) at 100% coverage. TUI functions at 0% due to TTY requirement -- tested indirectly via mocked function variables. All 31 packages in forge-cli pass.
