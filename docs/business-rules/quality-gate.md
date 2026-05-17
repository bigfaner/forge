---
title: "Quality Gate Rules"
domains: [quality-gate, pipeline, fix-task, compile, lint, retry]
---

# Quality Gate Rules

_Source: feature/forge-cli-v3_

## Pipeline

### BIZ-quality-gate-001: Quality Gate Multi-Phase Pipeline

**Rule**: `forge quality-gate` executes a multi-phase pipeline after all tasks complete: (1) compile -> fmt -> lint gate, (2) project-wide unit/integration tests (with retry-once policy for transient failures — if first attempt fails, tests are re-run once; only if both attempts fail is a fix task created), (3) e2e regression. Each failing step can auto-create a fix task (P0, breaking) up to a cap of 3 fix tasks per step; when the cap is reached, no new fix tasks are created and manual intervention is required. Each failure outputs a hook JSON block reason. Exit code 0 in all cases (hook JSON signals the actual decision). Docs-only features (no implementation or fix tasks) skip the quality gate entirely. The submit command uses a separate inline gate: compile -> fmt -> lint -> test, where failure exits with code 1.
**Context**: Provides project health enforcement after all feature tasks complete. The quality-gate hook evolved from a single pipeline into a multi-phase gate with fix-task auto-creation, retry-once tolerance for transient test failures, a per-step fix-task cap to prevent unbounded loop creation, and e2e regression support.
**Source**: feature/forge-cli-v3 BIZ-004
