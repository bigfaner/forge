# E2E Test Report: typed-task-dispatch

**Date**: 2026-05-12 10:04:58
**Duration**: 1.0s

## Summary

| Type  | Total | Pass | Fail | Skip |
|-------|-------|------|------|------|
| UI    | 0     | 0    | 0    | 0    |
| API   | 0     | 0    | 0    | 0    |
| CLI   | 20    | 19   | 0    | 1    |
| **All** | **20** | **19** | **0** | **1** |

**Result**: ✅ ALL TESTS PASSED

---

## Results by Test Case

### [TC-001] ✓ PASS — doc-generation.summary task prompt contains no TDD steps
**Duration**: 18ms

### [TC-002] ✓ PASS — fix task prompt contains five-step diagnostic flow
**Duration**: 17ms

### [TC-003] ✓ PASS — new type template generates correct prompt output
**Duration**: 223ms

### [TC-004] ✓ PASS — unregistered type causes non-zero exit with error
**Duration**: 14ms

### [TC-005] ✓ PASS — task prompt outputs complete synthesized prompt within 500ms
**Duration**: 14ms

### [TC-006] ✓ PASS — missing type causes non-zero exit with error
**Duration**: 14ms

### [TC-007] ⊘ SKIP — task migrate is idempotent on already-typed index
**Duration**: 0ms

### [TC-008] ✓ PASS — task migrate rejects when tasks are in_progress
**Duration**: 13ms

### [TC-009] ✓ PASS — breakdown-tasks skill generates type fields for tasks
**Duration**: 1ms

### [TC-010] ✓ PASS — breakdown-tasks falls back to implementation for unrecognized descriptions
**Duration**: 0ms

### [TC-011] ✓ PASS — execute-task and run-tasks produce identical task prompt output
**Duration**: 26ms

### [TC-012] ✓ PASS — execute-task marks task blocked when task prompt fails
**Duration**: 13ms

### [TC-013] ✓ PASS — run-tasks dispatches fix task via task prompt with five-step prompt
**Duration**: 14ms

### [TC-014] ✓ PASS — task prompt --fix-record-missed outputs record-recovery prompt
**Duration**: 14ms

### [TC-015] ✓ PASS — task validate accepts valid type enum values and rejects invalid ones
**Duration**: 26ms

### [TC-016] ✓ PASS — task prompt injects phase summary path for first task of new phase
**Duration**: 103ms

### [TC-017] ✓ PASS — run-tasks routes eval-cases task to main session, not subagent
**Duration**: 0ms

### [TC-018] ✓ PASS — task prompt --fix-record-missed outputs record-recovery prompt
**Duration**: 13ms

### [TC-019] ✓ PASS — quick-tasks skill includes type assignment rules
**Duration**: 0ms

### [TC-020] ✓ PASS — git branch fallback provides feature when state.json is missing
**Duration**: 13ms

---

## Failed Tests Detail

No failed tests.

---

## Screenshots

No UI tests in this suite (CLI-only).
