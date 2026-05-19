---
id: "1"
title: "CLI submit gate — static gate for non-breaking tasks + zero metrics acceptance"
priority: "P1"
estimated_time: "1-2h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 1: CLI Submit Gate — Static Gate + Zero Metrics

## Description

The CLI submit gate currently runs the full quality gate (compile→fmt→lint→test) for all `coding.*` tasks via `validateQualityGate()`. This is redundant for non-breaking tasks — the dispatcher's breaking gate (being removed in task 5) and gate tasks provide `just test` coverage at the right checkpoints.

Replace with a tiered model:
- **breaking=true**: full gate (`DefaultGateSequence()`: compile→fmt→lint→test)
- **breaking=false** + `coding.*` type: static gate (`LintGateSequence()`: compile→fmt→lint)
- **non-`coding.*` types**: skip entirely (already handled by `IsTestableType()`)

Also update `validateRecordData()` to accept zero test metrics for non-breaking coding tasks (the agent ran targeted tests during development, not `just test`).

## Reference Files
- `docs/proposals/deduplicate-quality-gate/proposal.md` — Source proposal (Tiered Test Execution Model)
- `forge-cli/pkg/just/just.go` — `LintGateSequence()` already exists for compile→fmt→lint

## Acceptance Criteria

- [ ] `validateQualityGate()` reads `t.Breaking` to choose gate sequence: breaking → `DefaultGateSequence()`, non-breaking → `LintGateSequence()`
- [ ] Non-`coding.*` types skip the quality gate entirely (existing behavior, unchanged)
- [ ] `validateRecordData()` accepts `testsPassed=0` + `testsFailed=0` for non-breaking coding tasks (set coverage=-1.0)
- [ ] Breaking coding tasks still require test evidence in `validateRecordData()`
- [ ] `forge task submit` for a non-breaking coding task passes with only compile+fmt+lint
- [ ] `forge task submit` for a breaking coding task requires compile+fmt+lint+test
- [ ] Existing tests pass; new tests cover the tiered gate logic

## Hard Rules

- Read `t.Breaking` from the task struct (already loaded from index.json), not from claim output or state
- Use existing `just.LintGateSequence()` — do not create a new gate sequence function

## Implementation Notes

- `forge-cli/internal/cmd/submit.go` line 492: `validateQualityGate()` currently always calls `just.DefaultGateSequence()`. Change to conditionally use `just.LintGateSequence()` when `!t.Breaking`.
- `forge-cli/internal/cmd/submit.go` lines 129-133 and 318-323: `validateRecordData()` currently rejects zero test evidence. Add a condition to accept it for non-breaking tasks.
- `forge-cli/pkg/just/just.go` line 34: `LintGateSequence()` returns compile→fmt→lint. Already used by `quality_gate.go` line 152.
- TDD: write tests for both gate paths (breaking and non-breaking) before modifying production code.
