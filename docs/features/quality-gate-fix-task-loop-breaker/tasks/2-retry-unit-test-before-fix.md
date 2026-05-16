---
id: "2"
title: "P1: Retry unit-test once before creating fix task"
priority: "P1"
estimated_time: "45m"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 2: P1: Retry unit-test once before creating fix task

## Description

Currently, transient test failures immediately create fix tasks. This is noisy — flaky tests are the norm, and auto-fix tasks should be a last resort. Add a retry-once policy for the unit-test step: if tests fail, re-run once. If retry passes, emit a warning and continue instead of creating a fix task.

The retry is a gate policy implemented in `quality_gate.go`, not a runner feature in `testrunner`.

## Reference Files
- `docs/proposals/quality-gate-fix-task-loop-breaker/proposal.md` — Source proposal
- `forge-cli/internal/cmd/quality_gate.go` — Lines 160-173 (unit-test step in runQualityGate)
- `forge-cli/internal/cmd/quality_gate_test.go` — Existing tests

## Acceptance Criteria
- [ ] When unit tests fail, they are retried once before creating a fix task
- [ ] If retry passes, a warning `WARNING: unit tests passed on retry (transient failure)` is logged to stderr, no fix task is created, gate continues to e2e step
- [ ] If retry also fails, a fix task is created with description including `"retried once, both attempts failed"` plus the retry-run output
- [ ] Retry logic only applies to the unit-test step, not compile/fmt/lint (deterministic failures don't benefit)
- [ ] New tests added for: retry pass (no fix task), retry fail (fix task with retry mention)

## Hard Rules
- Retry must not apply to compile, fmt, or lint steps — only unit-test.
- The retry-run output must be captured and included in the fix-task description on double-failure.

## Implementation Notes
- Pattern borrowed from CI systems (GitHub Actions retries, Buildkite `retry.automatic`).
- On retry-pass, the gate continues normally to the e2e step — the warning is informational only.
