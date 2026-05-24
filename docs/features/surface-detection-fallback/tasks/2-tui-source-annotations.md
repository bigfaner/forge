---
id: "2"
title: "Update TUI confirmation flow with source annotations and re-run prompt"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: Update TUI confirmation flow with source annotations and re-run prompt

## Description
Update the TUI confirmation flow to display inference source annotations from the new `DetectResult.Sources` map. Show context alongside each surface type so users can judge whether to accept or override. Add re-run behavior: when `config.yaml` already has surfaces configured, prompt with Confirm/Re-detect/Edit options instead of silently re-running detection.

## Reference Files
- `proposal.md#Proposed-Solution` — TUI confirmation flow update spec, source display format, hint text, re-run prompt
- `proposal.md#Requirements-Analysis` — Key Scenarios 7-8 for re-run behavior and user override flow
- `proposal.md#Success-Criteria` — TUI display criteria (#8-9), re-run prompt behavior (#13-14)
- `proposal.md#Key-Risks` — re-run Edit flow divergence risk and shared function mitigation
- `forge-cli/internal/cmd/init_surfaces.go` — TUI functions: askSurfaceConfirmation (L25), askScalarConfirmation (L51), askMapConfirmation (L105), buildDisplayLines (L144), manualSurfaceEntry (L308)

## Acceptance Criteria
- [ ] `askScalarConfirmation` shows source annotation in description field: `"cli (inferred from cmd/ directory structure)"` vs `"cli (detected from cobra dependency)"`
- [ ] `askMapConfirmation` shows per-path source annotation in display lines
- [ ] TUI hint text: when an inferred surface is displayed, show `"This was inferred from project structure. Edit to correct if needed."`; hint absent for dependency-detected surfaces
- [ ] Re-run behavior: `config.yaml` with existing `surfaces: cli` → TUI shows `"Surfaces already configured: cli. Re-detect?"` with Confirm / Re-detect / Edit options
- [ ] Confirm → returns `SKIPPED surfaces (already configured)` action
- [ ] Re-detect → runs full detection + inference pipeline and presents standard TUI confirmation
- [ ] Edit → calls `manualSurfaceEntry` flow (same function as first-run manual entry)
- [ ] User override: after editing inferred `cli` to `api` in TUI → config receives `surfaces: api` with no source annotation; source is display-only
- [ ] Multi-surface TUI display shows map-form list with source annotations: `forge-cli/cli → cli (inferred from cmd/ directory structure)`

## Hard Rules
- Source annotation is display-only and must NOT be persisted to config via the SurfacesMap type
- Re-run Edit flow must call the same `manualSurfaceEntry` function used by first-run — no separate code path

## Implementation Notes
- `askSurfaceConfirmation` at L25 needs access to `DetectResult.Sources` — change signature to accept full `*DetectResult` instead of just `SurfacesMap`
- `buildDisplayLines` at L144 generates display entries — add source annotation formatting based on Sources map entries
- `formatConflictAnnotation` at L168 already formats conflict metadata — use similar pattern for source annotation
- Re-run detection in `runSurfaceConfig` at L483: check if config already has surfaces before calling detection
- Source annotation format: parse Sources map value — `inference:cmd-dir` → `"inferred from cmd/ directory structure"`, `dependency:cobra` → `"detected from cobra dependency"`
- `askSurfaceConfirmation` return type may need to include Sources info for downstream consumers (init summary, YAML comment)
