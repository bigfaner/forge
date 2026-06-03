---
title: "Quality Gate Rules"
domains: [quality-gate, pipeline, fix-task, compile, lint, retry, unit-test, regression, NonBreakingGateSequence, UnitGateSequence, coding.cleanup]
---

# Quality Gate Rules

_Source: feature/forge-cli-v3_

## Pipeline

### BIZ-quality-gate-001: Quality Gate Multi-Phase Pipeline

**Rule**: `forge quality-gate` executes a three-phase pipeline after all tasks complete:

1. **Phase 1 — Compile/Fmt/Lint gate**: `NonBreakingGateSequence` (compile -> fmt -> lint). Each failing step can auto-create a fix task (P0, breaking for compile; P0, non-breaking for fmt/lint) up to a cap of 3 fix tasks per step. Failure outputs a hook JSON block reason and exits with code 0.
2. **Phase 2 — Unit tests**: `just unit-test [scope]` with retry-once policy (transient failure tolerance). If both attempts fail, a fix task is created. Failure outputs hook JSON and exits with code 0. The `scope` parameter uses `ResolveScope` probe-chain fallback: `just --dry-run compile <scope>` is tested first; if the justfile does not support scope, `just unit-test` runs without scope.
3. **Phase 3 — Test regression (surface-aware)**: When surfaces are configured in `.forge/config.yaml`, Phase 3 orchestrates per-surface-type lifecycle sequences via `runTestRegressionSurface`. For each unique surface type detected by `forgeconfig.SurfaceTypes()`:
   - **web/api surfaces**: dev -> probe -> test -> teardown (full lifecycle). Dev server started via `just <surface>-dev` or `just dev`. Probe checks server health with retry (3 attempts) via `just <surface>-probe` or `just probe`. Test runs via `just <surface>-test` or `just test`. Teardown always runs via `just <surface>-teardown` or `just teardown`.
   - **mobile surfaces**: dev -> probe -> test-setup -> test -> teardown (full lifecycle with mobile-specific setup step). Mobile test-setup runs between probe and test to handle emulator/simulator initialization.
   - **cli/tui surfaces**: test -> teardown (simplified lifecycle — no dev server or probe needed).
   - Surfaces of the same type share a single lifecycle run (dev/probe execute once per type, not per surface instance).
   - **On probe failure**: lifecycle aborts (returns error), teardown still runs. Error output saved, fix task created.
   - **On test failure**: teardown still runs (best-effort cleanup). Error output saved to `raw-output.txt`, fix task auto-created via `addRegressionFixTasks`.
   - **On teardown failure**: logged as warning, does not fail the lifecycle.
   - When no surfaces are configured, falls back to legacy behavior (`runTestRegressionLegacy`): optional `just test-setup`, `serverprobe.ProbeServers` health check, then `just test`.
   Failure creates fix task and outputs hook JSON. Exit code 0.

Exit code 0 in all cases (hook JSON signals the actual decision). Docs-only features (no implementation or fix tasks) skip the quality gate entirely. Fix tasks are grouped by source directory for parallel execution; compile/test failures use `coding.fix` type (breaking=true), fmt/lint failures use `coding.cleanup` type (breaking=false). The submit command uses a separate tiered gate: breaking tasks run `UnitGateSequence` (compile -> fmt -> lint -> unit-test), non-breaking coding tasks run `NonBreakingGateSequence` (compile -> fmt -> lint). Submit failure exits with code 1. The two-layer recipe model decouples `just unit-test` (language-level, fast per-task feedback) from `just test` (surface-level advanced tests, all-completed verification). Surface-aware orchestration details are defined in [surface-orchestration.md](./surface-orchestration.md).
**Context**: Provides project health enforcement after all feature tasks complete. The quality-gate hook evolved from a single FullGateSequence pipeline into a three-phase gate with fix-task auto-creation, retry-once tolerance for transient test failures, a per-step fix-task cap to prevent unbounded loop creation, directory-grouped parallel fix tasks, two-layer recipe model (unit-test / test) replacing the former single-test approach, and surface-aware Phase 3 that orchestrates per-surface-type lifecycle sequences (dev/probe/test/teardown for web/api/mobile, test/teardown for cli/tui).
**Source**: feature/forge-cli-v3 BIZ-004

## Surface-Specific Test Recipes

Phase 3 invokes surface-specific just recipes. Each surface type has a distinct lifecycle and recipe resolution chain. Recipes use a fallback pattern: `just <surface>-<verb>` is attempted first; if not defined, `just <verb>` is used as fallback.

| Surface Type | Lifecycle | Recipe Pattern (primary / fallback) |
|---|---|---|
| **web** | dev -> probe -> test -> teardown | `just <surface>-dev` / `just dev`; `just <surface>-probe` / `just probe`; `just <surface>-test` / `just test`; `just <surface>-teardown` / `just teardown` |
| **api** | dev -> probe -> test -> teardown | Same recipe pattern as web |
| **mobile** | dev -> probe -> test-setup -> test -> teardown | Same as web/api, plus `just <surface>-test-setup` / `just test-setup` between probe and test |
| **cli** | test -> teardown | `just <surface>-test` / `just test`; `just <surface>-teardown` / `just teardown` |
| **tui** | test -> teardown | Same recipe pattern as cli |
