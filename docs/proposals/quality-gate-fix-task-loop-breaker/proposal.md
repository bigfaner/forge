---
created: "2026-05-16"
author: "faner"
status: Draft
---

# Proposal: Quality-Gate Comprehensive Hardening

## Problem

The quality-gate stop hook has two confirmed bugs that cause infinite fix-task loops, plus latent defects that weaken its reliability as a safety net.

### Evidence

Two forensic reports confirm the feedback loop in production:

1. `docs/forensics/hook-feedback-loop/report.md` (2026-05-14): `just test` with `CGO_ENABLED=0` + `-race` caused deterministic failures. Fix tasks piled up: fix-1, fix-2, fix-3...
2. `docs/forensics/fix-task-loop/report.md` (2026-05-16): Flaky integration tests triggered the same loop after all quick-test-slim tasks completed.

### Root Cause Analysis

The feedback loop has three layers:

1. **Loop driver**: `saveIndexAndSignalCompletion` in `submit.go` re-writes `state.json` with `allCompleted=true` when fix tasks complete, re-triggering the Stop hook.
2. **Cap bypass (Bug A)**: `addFixTask` sets `Vars["SOURCE_TASK_ID"]` (template) but never `opts.SourceTaskID` (struct field). The counter's first filter `t.SourceTaskID != ""` always excludes quality-gate fix tasks — cap permanently reads 0.
3. **Cap reset (Bug B)**: `countActiveFixTasks` excludes completed/skipped tasks. After fix-1 completes, count resets to 0, allowing fix-2 for the same step.

The proposed fixes target layers 2 and 3 (make the cap effective). Layer 1 is mitigated: once the cap works, the loop terminates at 3 per step (up to 12 total across 4 steps) rather than running infinitely.

### Full Feature Audit Findings

Beyond the known bugs, a comprehensive code review uncovered these latent issues:

| # | Finding | Severity | File | Description |
|---|---------|----------|------|-------------|
| F1 | SourceTaskID not set (Bug A) | P0 | quality_gate.go:357-369 | Cap never triggers for quality-gate fix tasks |
| F2 | Cap counts only active tasks (Bug B) | P0 | quality_gate.go:307-319 | Completed fix tasks reset the counter |
| F3 | No retry before fix task | P1 | quality_gate.go:162-173 | Transient failures immediately create fix tasks |
| F4 | addFixTask silently returns ("", nil) on failure | P1 | quality_gate.go:374-400 | Template not found, task add failure, markdown creation failure — all produce empty fixID with no error. User sees "Failed to add fix task automatically" without knowing why. |
| F5 | Required missing recipe silently passes | P2 | just.go:105-107 | `RunGate` prints WARNING but returns true. By design: `Optional=false` means "must pass if exists, skip if absent". Correct for interpreted-language projects without a compile recipe. |
| F6 | No timeout on recipe execution | P2 | just.go (RunCapture) | A hanging `just` recipe blocks the quality-gate hook indefinitely. |
| F7 | No file locking on state.json | P3 | forge_state.go | TOCTOU race under concurrent hook invocation. Low risk (sessions are sequential). |
| F8 | Eager state consumption | P3 | quality_gate.go:84 | `ClearForgeState` before tests. A crash loses the signal. |
| F9 | Cap is per-step, not global | P3 | quality_gate.go:309 | compile + lint + unit-test + test-e2e each get 3 = 12 total per cycle. By design. |

### Urgency

The P0 bugs make the safety net an amplifier. Every session hitting a flaky test produces infinite fix tasks — confirmed twice in production. P1 issues mask failures and reduce operator confidence.

## Proposed Solution

Fix P0 and P1 issues. P2 and below are deferred (see rationale in Scope).

### P0: Break the feedback loop

1. **Set step-scoped `SourceTaskID`**: Use `"quality-gate:<step>"` sentinel in `addFixTask` (e.g., `"quality-gate:compile"`, `"quality-gate:unit-test"`). This makes quality-gate fix tasks identifiable by the counter **per step**. Step-scoping avoids cross-step dedup conflicts — `HasActiveFixTasks(index, "quality-gate:compile")` only blocks compile fix tasks, not lint or unit-test ones.
2. **Count cumulatively**: Change `countActiveFixTasks` to count ALL fix tasks per step (including completed/skipped), not just active ones. Rename to `countFixTasks`. Both A and B must be applied together — the counter's `SourceTaskID != ""` filter only passes after Fix A populates the field.

### P1: Reduce false-positive fix tasks + fix silent failures

3. **Retry unit-test once before fix task**: When unit tests fail, re-run once. If retry passes, emit a warning instead of creating a fix task. Implementation in `quality_gate.go` (not in `testrunner` — retry is a gate policy, not a runner feature). On retry-also-fails: the fix-task description includes `"retried once, both attempts failed"` plus the retry-run output.
4. **Return explicit errors from addFixTask**: Replace silent `("", nil)` returns with proper error propagation on template-not-found, task-add-failure, and markdown-creation-failure paths. Callers log the specific error instead of discarding with `_`. The `handleGateFailure` fallback (manual add instruction) remains unchanged.

### Housekeeping

5. **Version bump**: Patch bump in `scripts/version.txt` per project CLAUDE.md requirements.

### Innovation Highlights

The retry-before-fix pattern is borrowed from CI systems (GitHub Actions retries, Buildkite `retry.automatic`). Transient failures are the norm in test suites — auto-fix tasks should be a last resort.

Step-scoped sentinels (`"quality-gate:<step>"`) maintain the existing per-step cap granularity while working correctly with both `countFixTasks` (step-scoped by title prefix) and `HasActiveFixTasks` (step-scoped by sentinel value). This avoids the cross-step blocking that a flat `"quality-gate"` sentinel would cause.

## Requirements Analysis

### Key Scenarios

1. **Flaky test in full suite**: `just test` fails, retry passes, warning logged, no fix task created.
2. **Genuine test failure**: `just test` fails, retry also fails, fix task created with `SourceTaskID: "quality-gate:unit-test"`. Cap counts it cumulatively.
3. **Same step fails 3 times total**: After 3 cumulative fix tasks for a step, no more are created. Human intervention required.
4. **Cross-step independence**: Fix-1 for compile (pending) does NOT block creating fix-2 for unit-test. Each step has its own sentinel and its own cap.
5. **Template not found**: `addFixTask` returns an explicit error with reason. Caller logs it, `handleGateFailure` shows manual add instruction.
6. **Fix task completes, same failure recurs**: Cumulative count is 1 of 3. Next failure creates fix-2. After 3 total, stop.
7. **Retry passes**: Warning `WARNING: unit tests passed on retry (transient failure)` logged to stderr. No fix task. Gate continues to e2e step.

### Constraints & Dependencies

- Sentinel `"quality-gate:<step>"` must not collide with real task IDs. `FindTask` returns nil for it, so source resolution/blocking is correctly skipped for all steps.
- Existing test at `quality_gate_test.go:598` asserts `SourceTaskID == ""` — must be updated to assert the step-scoped sentinel.
- `Vars["SOURCE_TASK_ID"]` remains `"N/A (project-wide gate)"` for template rendering (intentionally diverges from the struct field value).

### P2/P3 Deferrals with Rationale

| Finding | Why Deferred |
|---------|-------------|
| F5 (required missing recipe) | Current "warn and skip" is correct: `Optional=false` means "must pass if exists, skip if absent". Failing would break Python/interpreted-language projects that lack a `compile` recipe. |
| F6 (recipe timeout) | Requires `RunCapture` API change (add `context.Context` / timeout). Separate feature. |
| F7 (state.json locking) | Sessions are sequential — practical risk near zero. |
| F8 (eager state consumption) | Crash edge case only. Adding deferred state management adds complexity for minimal gain. |

## Alternatives & Industry Benchmarking

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | No code change | Infinite loop continues | Rejected |
| Fix P0 only (A+B) | Minimal change, breaks loop | No retry resilience; silent failures remain | Viable but incomplete |
| P0+P1 with flat `"quality-gate"` sentinel | Simpler sentinel | Cross-step dedup: pending compile fix blocks lint fix creation | Rejected |
| **P0+P1 with step-scoped sentinel** | Breaks loop, retry resilience, correct per-step semantics | Slightly more complex sentinel format | **Selected** |

## Scope

### In Scope

- Set `opts.SourceTaskID: "quality-gate:" + step` in `addFixTask`
- Rename `countActiveFixTasks` → `countFixTasks`, count all statuses (not just active)
- Add retry-once logic for unit-test step in `quality_gate.go` (not `testrunner`)
- Replace silent `("", nil)` returns in `addFixTask` with proper errors; callers log errors
- Patch version bump in `scripts/version.txt`
- Update existing tests for all changed behavior
- Add new tests for: step-scoped sentinel, cumulative counting, retry pass, retry fail, error propagation

### Out of Scope

- Fixing specific flaky tests — separate concern
- Retry for other quality-gate steps (compile, fmt, lint) — deterministic failures don't benefit
- Changing the cap value (stays at 3)
- Failing on missing required recipe (F5) — breaks interpreted-language projects
- File locking on state.json (F7), eager state consumption (F8) — P3
- Timeout on recipe execution (F6) — requires `RunCapture` API redesign
- Scope resolution in quality-gate (hardcoded empty scope) — existing behavior, unrelated to loop fix

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Sentinel `"quality-gate:<step>"` collides with a real task ID | L | M | Task IDs are UUIDs — impossible in practice |
| Cumulative cap permanently blocks fix tasks after 3 per step | L | M | Retry-once absorbs transient failures. 3 persistent failures warrant human attention. Manual `forge task add` still available as escape hatch. |
| Retry masks a real regression (passes on retry) | L | M | Warning is logged. Failure visible in stderr. |
| `HasActiveFixTasks` still cross-blocks if step is misspelled in sentinel | L | L | Sentinel is constructed programmatically (`"quality-gate:" + step`) using the same `step` variable used for title prefix — no typo risk. |

## Success Criteria

- [ ] `addFixTask` creates tasks with `SourceTaskID: "quality-gate:<step>"` (verified by test per step)
- [ ] `countFixTasks` counts fix tasks regardless of status (completed + active + blocked)
- [ ] When unit tests fail, they are retried once before creating a fix task
- [ ] If retry passes, a warning is logged and no fix task is created
- [ ] If retry fails, fix-task description mentions "retried once, both attempts failed"
- [ ] After 3 cumulative fix tasks for a step, no more are created regardless of status
- [ ] Pending fix for step A does NOT block fix creation for step B (cross-step independence)
- [ ] `addFixTask` returns explicit errors on template-not-found, task-add-failure, and markdown-creation-failure
- [ ] Version bumped in `scripts/version.txt`
- [ ] All existing quality-gate tests pass with updated assertions
- [ ] No regression in non-quality-gate fix task flows (dispatcher-created fix tasks still work)

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
