---
id: "3"
title: "P1: Explicit error propagation from addFixTask"
priority: "P1"
estimated_time: "30m"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 3: P1: Explicit error propagation from addFixTask

## Description

`addFixTask` silently returns `("", nil)` on three failure paths: template not found, task add failure, and markdown creation failure. The caller discards the empty fixID with `_`, and the user sees only "Failed to add fix task automatically" without knowing why. Replace these silent returns with proper error propagation so callers can log the specific reason.

## Reference Files
- `docs/proposals/quality-gate-fix-task-loop-breaker/proposal.md` — Source proposal
- `forge-cli/internal/cmd/quality_gate.go` — Lines 374-400 (addFixTask error paths)
- `forge-cli/internal/cmd/quality_gate_test.go` — Existing tests

## Acceptance Criteria
- [ ] `addFixTask` returns explicit errors on template-not-found
- [ ] `addFixTask` returns explicit errors on task-add-failure
- [ ] `addFixTask` returns explicit errors on markdown-creation-failure
- [ ] Callers log the specific error instead of discarding with `_`
- [ ] `handleGateFailure` fallback (manual add instruction) remains unchanged
- [ ] New tests added for: error propagation on each failure path

## Hard Rules
- The `handleGateFailure` fallback behavior (showing manual add instruction) must remain unchanged — it's the user-facing safety net.

## Implementation Notes
- The three silent return points are in `addFixTask` at lines 374, 385, 397 (approximate — check current line numbers).
- Callers currently use `fixID, _ := addFixTask(...)` — change to capture and log the error.
