---
id: "3"
title: "Add surface-aware orchestration to quality-gate"
priority: "P1"
estimated_time: "2-3h"
dependencies: ["1"]
surface-key: ""
surface-type: ""
breaking: false
type: "coding.feature"
mainSession: false
---

# 3: Add surface-aware orchestration to quality-gate

## Description

The quality-gate (`forge quality-gate`) currently runs `just test` blindly in Phase 3 (test regression), without surface awareness. For web/api surfaces, this means no dev server lifecycle management — if the server isn't running, tests are silently skipped.

Add surface-aware orchestration: quality-gate reads `.forge/config.yaml` to detect surface types, then executes the appropriate orchestration sequence per surface:
- **web/api**: dev → probe → test → teardown (full lifecycle)
- **cli/tui**: test → teardown (simplified)
- **No surfaces configured**: fall back to current behavior (`just test`)

The orchestration sequences mirror the run-tests skill's surface rules (Task 1 defines the canonical behavior in SKILL.md).

## Reference Files
- `proposal.md#Proposed-Solution` — P3: quality-gate surface orchestration using surface config
- `proposal.md#Requirements-Analysis` — Key Scenarios: quality-gate surface编排 (multi-surface and single-surface)
- `proposal.md#Key-Risks` — risk of inconsistency between quality-gate Go code and run-tests skill
- `proposal.md#Constraints-&-Dependencies` — quality-gate is Go binary, run-tests is AI skill, independent implementations
- `proposal.md#Feasibility-Assessment` — quality-gate already has `just.RunCapture` capability

## Acceptance Criteria

- [ ] quality-gate reads surface types from `.forge/config.yaml` (via existing `forgeconfig` package)
- [ ] For web/api surfaces: quality-gate executes dev → probe → test → teardown sequence (full lifecycle)
- [ ] For cli/tui surfaces: quality-gate executes test → teardown sequence (simplified)
- [ ] When no surfaces configured: quality-gate falls back to current `just test` behavior
- [ ] Probe retry logic (3 retries, 5s intervals) matches run-tests skill behavior
- [ ] Teardown is mandatory (executed even on prior step failure)
- [ ] Existing `forge-cli/internal/cmd/quality_gate_test.go` tests pass
- [ ] New tests cover: single surface (cli), multi-surface, no surfaces, probe failure, dev failure
- [ ] Version bump in `scripts/version.txt` (minor: new feature)

## Hard Rules

- Surface orchestration sequences MUST mirror the run-tests skill's surface rules (Task 1 defines the canonical behavior)
- Teardown is mandatory regardless of prior step success/failure
- Probe failure after all retries → execute teardown → exit with error
- Dev failure → execute teardown → exit with dev's exit code
- Do NOT introduce new justfile recipes — use existing ones (`just dev`, `just probe`, `just test`, `just teardown` or surface-specific variants)

## Implementation Notes

- The surface config is in `.forge/config.yaml` under the `interfaces` key. Use `forgeconfig.LoadConfig()` to read surfaces.
- Surface types are available via `forgeconfig.SurfaceTypes(surfaces)`. The orchestration mapping:
  - `web`, `api`, `mobile`: full lifecycle (dev → probe → test → teardown)
  - `cli`, `tui`: simplified (test → teardown)
- The existing `runTestRegression()` function in `quality_gate.go` should be refactored to accept surface info and orchestrate per-surface.
- Use `just.HasRecipe()` to check if surface-specific recipes exist (e.g., `web-dev`, `web-probe`). Fall back to generic recipes (`dev`, `probe`, `test`, `teardown`) if surface-specific ones don't exist.
- The probe retry logic (3 retries, 5s intervals) should be extracted into a reusable function.
- The existing `e2eprobe.ProbeServers()` is specific to e2e/HTTP health checks. The new surface-aware probe should use `just probe` recipe instead, which is surface-agnostic.
- Consider creating a `SurfaceOrchestrator` struct or similar to encapsulate the per-surface lifecycle, keeping `runTestRegression()` clean.
