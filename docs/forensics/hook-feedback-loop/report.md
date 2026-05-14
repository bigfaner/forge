# Forensic Report: Hook Feedback Loop (fix-3 through fix-10)

**Date**: 2026-05-14
**Session**: e36e4498-08de-496c-aef2-9d623246b70f
**Feature**: forge-cli-v3
**Investigator**: /forensic

## Executive Summary

The `task all-completed` Stop hook created a self-perpetuating feedback loop of false-positive fix tasks (fix-3 through fix-10). Every fix-task completion re-triggered the hook, which detected all tasks as completed, ran tests, and created another fix task on failure.

The test failure was **not intermittent** — it was **deterministic under specific conditions**. The root cause is a logic bug in the justfile's `test` recipe.

## Root Cause

**File**: `justfile:64`

```bash
# Line 60-64 (the buggy version)
race_flag=""
if command -v gcc &>/dev/null; then
    race_flag="-race"
fi
cd forge-cli && CGO_ENABLED=0 go test $race_flag ./...
```

**The bug**: The recipe checks for `gcc` to decide whether to enable `-race`, but then unconditionally sets `CGO_ENABLED=0`. When both conditions are true (gcc found + CGO disabled), `go test` fails with:

```
go: -race requires cgo; enable cgo by setting CGO_ENABLED=1
```

Exit code 2 → `set -euo pipefail` → recipe exits 1 → `just test` fails.

**Why it's environment-dependent**: On the developer's primary Windows environment, `gcc` is NOT in PATH, so `race_flag=""` and tests pass. But in certain hook execution contexts (where MSYS2/Git Bash toolchain adds `gcc` to PATH), `race_flag="-race"` is set, triggering the conflict.

**Reproduction**:

```bash
CGO_ENABLED=0 go test -race ./...
# go: -race requires cgo; enable cgo by setting CGO_ENABLED=1
# exit code: 2
```

## Feedback Loop Mechanism

```
1. All tasks complete → task record sets .forge/state.json allCompleted=true
2. Stop hook fires → task all-completed (or forge quality-gate) runs
3. Quality gate: just test → FAIL (gcc in PATH + CGO_ENABLED=0 + -race)
4. handleGateFailure → addFixTask → creates fix-N
5. task record fix-N → sets allCompleted=true again
6. Go to step 2 (infinite loop)
```

The loop was only broken by manually setting `allCompleted=false` after recording each fix task.

## Causal Chain (3 levels)

| Level | Finding |
|-------|---------|
| **Symptom** | 8 consecutive false-positive fix tasks (fix-3 through fix-10) created by the Stop hook |
| **Direct cause** | `just test` exits 1 in the hook context but exits 0 when run directly from the terminal |
| **Root cause** | justfile `test` recipe sets `CGO_ENABLED=0` unconditionally while conditionally enabling `-race` based on `gcc` availability — these two flags are mutually exclusive |

## Fix

The justfile recipe should conditionally set `CGO_ENABLED=0` only when NOT using `-race`:

```bash
test scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    if command -v gcc &>/dev/null; then
        cd forge-cli && go test -race ./...
    else
        cd forge-cli && CGO_ENABLED=0 go test ./...
    fi
```

This ensures `-race` and `CGO_ENABLED=0` are never combined.

## Secondary Finding: Hook Loop Prevention

The quality-gate hook lacks a cooldown or deduplication mechanism. Each `task record` call re-sets `allCompleted=true`, allowing the hook to fire again immediately. Suggested mitigations:

1. **Cooldown**: Skip quality-gate if last run was < 60 seconds ago
2. **State consumption**: The hook already clears state before running tests, but `task record` re-sets it. Consider not setting `allCompleted` for fix tasks (they are remedial, not progression)
3. **Fix-task cap**: The `maxFixTasksPerStep` (3) cap exists but the step name alternates between "login bug" and "unit-test failure" — each with its own cap counter

## Deviation Classification

| Category | Description |
|----------|-------------|
| `pipeline-gap` | No enforcement between test recipe flags — `-race` and `CGO_ENABLED=0` can conflict silently |
| `instruction-gap` | The run-tasks dispatcher had no protocol for handling hook feedback loops — it kept creating fix tasks reactively |
