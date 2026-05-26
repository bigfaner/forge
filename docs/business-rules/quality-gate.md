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
2. **Phase 2 — Unit tests**: `just unit-test` with retry-once policy (transient failure tolerance). If both attempts fail, a fix task is created. Failure outputs hook JSON and exits with code 0.
3. **Phase 3 — Test regression**: `just test` (full regression suite, requires probe health check first). Optional `just test-setup` step. Failure creates fix task and outputs hook JSON. Exit code 0.

Exit code 0 in all cases (hook JSON signals the actual decision). Docs-only features (no implementation or fix tasks) skip the quality gate entirely. Fix tasks are grouped by source directory for parallel execution; compile/test failures use `coding.fix` type (breaking=true), fmt/lint failures use `coding.cleanup` type (breaking=false). The submit command uses a separate tiered gate: breaking tasks run `UnitGateSequence` (compile -> fmt -> lint -> unit-test), non-breaking coding tasks run `NonBreakingGateSequence` (compile -> fmt -> lint). Submit failure exits with code 1. The two-layer recipe model decouples `just unit-test` (language-level, fast per-task feedback) from `just test` (surface-level advanced tests, all-completed verification).
**Context**: Provides project health enforcement after all feature tasks complete. The quality-gate hook evolved from a single FullGateSequence pipeline into a three-phase gate with fix-task auto-creation, retry-once tolerance for transient test failures, a per-step fix-task cap to prevent unbounded loop creation, directory-grouped parallel fix tasks, and two-layer recipe model (unit-test / test) replacing the former single-test approach.
**Source**: feature/forge-cli-v3 BIZ-004
