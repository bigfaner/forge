---
title: "Quality Gate Rules"
domains: [quality-gate, pipeline, fix-task, compile, lint, retry]
---

# Quality Gate Rules

_Source: feature/forge-cli-v3_

## Pipeline

### BIZ-quality-gate-001: Quality Gate Multi-Phase Pipeline

**Rule**: `forge quality-gate` executes a multi-phase pipeline after all tasks complete using FullGateSequence (compile -> fmt -> lint -> unit-test -> test -> probe). Each failing step can auto-create a fix task (P0, breaking) up to a cap of 3 fix tasks per step; when the cap is reached, no new fix tasks are created and manual intervention is required. Each failure outputs a hook JSON block reason. Exit code 0 in all cases (hook JSON signals the actual decision). Docs-only features (no implementation or fix tasks) skip the quality gate entirely. The submit command uses a tiered gate: breaking tasks run UnitGateSequence (compile -> fmt -> lint -> unit-test), non-breaking coding tasks run NonBreakingGateSequence (compile -> fmt -> lint). Failure exits with code 1. The two-layer recipe model decouples `just unit-test` (language-level, fast per-task feedback) from `just test` (surface-level advanced tests, all-completed verification).
**Context**: Provides project health enforcement after all feature tasks complete. The quality-gate hook evolved from a single pipeline into a multi-phase gate with fix-task auto-creation, retry-once tolerance for transient test failures, a per-step fix-task cap to prevent unbounded loop creation, and two-layer recipe model (unit-test / test) replacing the former single-test approach.
**Source**: feature/forge-cli-v3 BIZ-004
