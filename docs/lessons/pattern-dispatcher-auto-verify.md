# `.forge/state.json` Not Created Due to Subagent CWD Mismatch

## Problem
After all tasks completed via `/run-tasks`, the Stop hook fires `task all-completed` but it silently skips — no e2e tests run, no verification happens.

## Root Cause
`task all-completed` has a dual-gate check (all_completed.go:70-74): all tasks must be done AND `.forge/state.json` must exist with `allCompleted=true`. The state file is never created because of a CWD mismatch.

Causal chain:
1. Symptom: Stop hook fires `task all-completed`, it exits 0 with no output, no tests run
2. Direct cause: `.forge/state.json` doesn't exist, so `checkAllCompleted` returns nil
3. Root cause: `task record` (called by subagents from `backend/` CWD) calls `saveIndexAndSignalCompletion` which calls `feature.WriteForgeState`. But `project.FindProjectRoot()` detects `backend/` as project root (finds `go.mod`), so `.forge/state.json` is written to the wrong location or fails silently
4. Evidence: `task all-completed -v` run from `backend/` shows `project root: Z:\...\backend` — wrong directory
5. Trigger: subagents inherit `backend/` as CWD because Go test/compile commands run from there

The flow that should work:
```
task record → saveIndexAndSignalCompletion → all tasks done? → WriteForgeState → .forge/state.json created
                                                                                        ↓
Session stops → Stop hook → task all-completed → checkAllCompleted finds state.json → runs e2e tests
```

The flow that actually happens:
```
task record (from backend/) → FindProjectRoot() returns backend/ → .forge/ written to wrong dir or fails
                                                                                          ↓
Session stops → Stop hook → task all-completed → no state.json → silent skip
```

## Solution
The subagent must run `task record` from the **project root** (not `backend/`). Options:
1. Subagent should `cd` to project root before running `task record`
2. Or `task record` should accept `--project-root` flag
3. Or `FindProjectRoot` should look for project-level markers (`.claude/`, `CLAUDE.md`) rather than `go.mod`

## Key Takeaway
When subagents run task-cli commands, their CWD affects `FindProjectRoot()`. The `task record` command creates `.forge/state.json` only when it finds the correct project root — which subagents running from `backend/` cannot locate. This breaks the entire automated verification chain: no state file → `task all-completed` skips → no e2e tests → no test graduation.
