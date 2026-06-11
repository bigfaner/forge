---
status: "completed"
started: "2026-05-27 00:29"
completed: "2026-05-27 00:42"
time_spent: "~13m"
---

# Task Record: 3 Add surface-aware orchestration to quality-gate

## Summary
Added surface-aware orchestration to quality-gate: reads surface types from .forge/config.yaml and executes per-surface lifecycle (dev→probe→test→teardown for web/api/mobile, test→teardown for cli/tui). Falls back to legacy behavior when no surfaces configured.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/quality_gate_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Split runTestRegression into surface-aware path and legacy path for backward compatibility
- Used resolveRecipe() to prefer surface-specific recipes (e.g. web-dev) over generic ones (dev)
- Probe uses just recipe (not e2eprobe.ProbeServers) for surface-agnostic health checking
- Teardown is always executed via explicit calls after each failure point, not defer

## Test Results
- **Tests Executed**: Yes
- **Passed**: 15
- **Failed**: 0
- **Coverage**: 64.8%

## Acceptance Criteria
- [x] quality-gate reads surface types from .forge/config.yaml via forgeconfig
- [x] For web/api surfaces: dev → probe → test → teardown sequence
- [x] For cli/tui surfaces: test → teardown sequence
- [x] No surfaces configured: falls back to current just test behavior
- [x] Probe retry logic (3 retries, 5s intervals) matches run-tests skill
- [x] Teardown is mandatory (executed even on prior step failure)
- [x] Existing quality_gate_test.go tests pass
- [x] New tests cover: single surface (cli), multi-surface, no surfaces, probe failure, dev failure
- [x] Version bump in scripts/version.txt (minor: new feature)

## Notes
Probe failure test takes ~10s due to 3 retries with 5s intervals. Tests use bash scripts for marker files to avoid Windows path issues in justfile recipes.
