---
created: "2026-05-16"
sessions: [a71a854d]
skillsInvolved: [forge:quick, forge:quick-tasks, forge:run-tasks]
severity: medium
---

# Fix-Task Feedback Loop After All Tasks Completed

## Executive Summary

After the `/quick` pipeline completed all 7 tasks for `quick-test-slim`, the `Stop` hook triggered `forge quality-gate` which ran `just test`. Two flaky integration tests (`TestRunFeature_None`, `TestForgeStateLifecycle`) failed in the full suite despite passing individually. Each failure created a fix task, which when resolved triggered another `Stop` hook → quality-gate → same failure → new fix task. This is an infinite feedback loop caused by flaky tests in the quality-gate + Stop hook combination.

## Investigation Scope

| Dimension | Value |
|-----------|-------|
| Sessions analyzed | 1 (a71a854d) |
| Time range | 2026-05-16 13:58 to 16:04 |
| Skills involved | forge:quick, forge:quick-tasks, forge:run-tasks |
| Trigger | User asked why fix tasks kept appearing after all tasks completed |

## Timing Overview

| Session | Duration | Tool Time | Idle* | Top Bottleneck |
|---------|----------|-----------|-------|---------------|
| a71a854d | 2.1h | 6027.9s (100.5min) | ~25min | Agent (734.3s avg) |

*Idle = session duration minus total tool execution time.

| Tool | Calls | Total | Avg | Max |
|------|-------|-------|-----|-----|
| Agent | 8 | 5874.6s | 734.3s | 1870.1s |
| Bash | 51 | 149.8s | 2.9s | 68.6s |
| Read | 10 | 454ms | 45ms | 172ms |
| Write | 2 | 593ms | 296ms | 303ms |
| Skill | 2 | 88ms | 44ms | 52ms |
| AskUserQuestion | 1 | 2.4s | 2.4s | 2.4s |

## Findings

### Finding 1: Quality-Gate Stop Hook Creates Fix-Task Feedback Loop

**Category:** `pipeline-gap`

**Affected sessions:** a71a854d

**Symptom:**
After all 7 tasks completed successfully, the session's Stop hook ran `forge quality-gate`. Two integration tests failed (flaky in full suite). `quality_gate.go:156` (`addFixTask`) created fix-1. After fix-1 was resolved (tests passed on re-run), the session stopped again, triggering the Stop hook again, which failed on the same flaky tests, creating fix-2. This loop would continue indefinitely.

**Expected behavior:**
The all-completed hook (quality-gate) is a final safety net. When it detects a real regression, fix tasks should be created. But when the same tests fail repeatedly (flaky), the loop should break.

**Causal chain:**
1. **Symptom:** fix-1, fix-2 tasks created after all tasks completed
2. **Direct cause:** `Stop` hook runs `forge quality-gate` on every session stop, including after fix-task completion
3. **Root cause:** No mechanism to break the Stop → gate → fail → fix → Stop cycle. `addFixTask` has a per-step cap (`maxFixTasksPerStep = 3`) but it counts active fix tasks in the index. After fix-1 is marked completed, the cap resets, allowing fix-2 for the same step.

**Hook flow:**
```
Session ends → Stop hook → forge quality-gate
  → just test → FAIL (flaky TestRunFeature_None, TestForgeStateLifecycle)
  → addFixTask("unit-test", ...) → fix-1
  → session blocked (hook returns error)
  → dispatcher claims fix-1 → resolves → submits
  → session ends → Stop hook → forge quality-gate
  → just test → FAIL (same flaky tests)
  → addFixTask("unit-test", ...) → fix-2
  → ...
```

### Finding 2: Flaky Integration Tests in Full Suite

**Category:** `trust-without-verify`

**Affected sessions:** a71a854d (manifested in quality-gate runs)

**Symptom:**
`TestRunFeature_None` and `TestForgeStateLifecycle` in `forge-cli/internal/cmd/integration_test.go` fail when run in the full package suite but pass individually (verified 3x in this session). The tests create temp dirs and `os.Chdir()` into them, with cleanup via `t.Cleanup(func() { _ = os.Chdir(origWd) })`.

**Root cause hypothesis:**
These tests modify global state (working directory) and rely on `t.Cleanup` for restoration. Earlier tests in the package may also modify global state (`.forge/` config, feature state) that leaks through shared filesystem or environment. Under resource pressure (the `fork/exec link.exe: The paging file is too small` error confirms Windows memory pressure during parallel test execution), cleanup may not execute properly or tests may run in unexpected order.

**Evidence:**
- First quality-gate run: `fork/exec link.exe: The paging file is too small` + both tests failed
- Second quality-gate run: both tests failed (no paging error, but resource pressure likely)
- Individual runs: both tests pass consistently
- Both tests together: both pass
- Full suite: both pass (when not under resource pressure)

### Finding 3: Task Blocking Cascade (T-quick-3, T-quick-4, T-quick-6)

**Category:** `pipeline-gap`

**Affected sessions:** a71a854d

**Symptom:**
After T-quick-2, T-quick-3, and T-quick-5 completed, their downstream tasks (T-quick-3, T-quick-4, T-quick-6) were found in "blocked" status despite all dependencies being completed. The dispatcher had to manually unblock each one with `forge task status <id> pending`.

**Pattern observed:**
| After task completed | Next task found blocked |
|---------------------|----------------------|
| T-quick-2 (gen-scripts) | T-quick-3 (run-tests) |
| T-quick-3 (run-tests) | T-quick-4 (graduate) |
| T-quick-4 (graduate) | T-quick-5 (verify-regression) — NOT blocked |
| T-quick-5 (verify-regression) | T-quick-6 (drift-detection) |

**Root cause hypothesis:**
The `forge task submit` command (`submit.go:283-286`) auto-downgrades tasks with test failures to "blocked". The raw test output shows multiple `auto-downgrading status from 'completed' to 'blocked'` warnings. When a task-executor subagent submits a task, if the quality gate detects test failures (possibly from the same flaky suite), the NEXT task in the chain may be affected. Alternatively, the task-executor may be explicitly blocking downstream tasks when it detects issues with the merged gen+run code.

**Note:** Without extracting individual subagent transcripts, the exact mechanism (auto-downgrade of submitted task vs. explicit blocking of downstream task) cannot be confirmed. This would require further subagent forensic extraction.

## Recommendations

| Priority | Action | Target File | Finding |
|----------|--------|-------------|---------|
| P0 | **Break the fix-task feedback loop:** Track cumulative fix-task count per step across completed fix-tasks (not just active ones). After N completed fix-tasks for the same step with the same test failures, stop creating new fix-tasks and emit a warning instead. | `forge-cli/internal/cmd/quality_gate.go:304-327` (countActiveFixTasks + addFixTask) | Finding 1 |
| P0 | **Fix flaky integration tests:** Investigate test isolation in `TestRunFeature_None` and `TestForgeStateLifecycle`. Likely need to use `t.Parallel()` isolation or stub out feature/config resolution instead of relying on filesystem state. | `forge-cli/internal/cmd/integration_test.go:1230,1785` | Finding 2 |
| P1 | **Investigate task blocking cascade:** Extract subagent transcripts for T-quick-2, T-quick-3, T-quick-5 to determine if `forge task submit` auto-downgrade or explicit subagent action caused the downstream blocking. | Subagent transcripts in `~/.claude/projects/` | Finding 3 |
| P1 | **Add flaky-test resilience to quality-gate:** Before creating fix tasks, re-run the failing tests once. If they pass on retry, emit a warning instead of creating a fix task. | `forge-cli/internal/cmd/quality_gate.go:160-173` | Finding 1, 2 |

## Evidence

Evidence files at: `docs/forensics/fix-task-loop/evidence/`

| File | Source | Size |
|------|--------|------|
| evidence.json | Main session a71a854d | ~15 KB |
