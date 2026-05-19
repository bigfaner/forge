---
id: "5"
title: "Remove dispatcher breaking gate from run-tasks.md and execute-task.md"
priority: "P1"
estimated_time: "30min"
dependencies: []
type: "doc"
mainSession: false
---

# 5: Remove Dispatcher Breaking Gate

## Description

Both dispatchers (`run-tasks.md` and `execute-task.md`) have a Step 3 "Breaking Task Gate" that runs `just test` for tasks with `BREAKING=true`. This is now redundant because:

1. The CLI submit gate (task 1) runs the appropriate quality gate at submission time
2. Breaking tasks get the full gate (compile→fmt→lint→test) at submit
3. The dispatcher's `just test` is an extra redundant layer

Remove Step 3 entirely from both dispatchers. Also remove BREAKING from the claim output parsing in both files (BREAKING field removed from claim output in task 2).

## Reference Files
- `docs/proposals/deduplicate-quality-gate/proposal.md` — Source proposal (item 4)
- `plugins/forge/commands/run-tasks.md` — Step 3 Breaking Task Gate
- `plugins/forge/commands/execute-task.md` — Step 3 Breaking Task Gate

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/run-tasks.md` | Remove Step 3 entirely. Remove BREAKING from Step 1 extract list. Update flowchart. Update error handling table. |
| `plugins/forge/commands/execute-task.md` | Remove Step 3 entirely. Remove BREAKING from Step 1 extract list. Update error handling table. |

## Acceptance Criteria

- [ ] `run-tasks.md` has no Step 3 Breaking Task Gate
- [ ] `execute-task.md` has no Step 3 Breaking Task Gate
- [ ] Neither dispatcher extracts BREAKING from claim output
- [ ] Neither dispatcher runs `just test` directly (only the CLI submit gate does this)
- [ ] Flowcharts updated: dispatch→verify→STOP (no breaking gate)
- [ ] Error handling tables updated: remove breaking-gate-failure rows

## Implementation Notes

- `run-tasks.md`: Step 3 is lines 81-98. Also update mermaid flowchart (line 21) to remove E→LOOP. Update Dispatcher Iron Laws (line 28) — remove "EXCEPT in Step 3". Remove from Step 1 extract list (line 55). Remove test-failure row from error handling table (line 112).
- `execute-task.md`: Step 3 is lines 95-184. Also update mermaid flowchart to remove E node. Remove BREAKING from Step 1 extract (line 27). Remove 3a/3b/fix-task rows from error handling table (lines 209-211).
- The dispatchers still extract SCOPE, FEATURE, MAIN_SESSION from claim output — these remain needed.
