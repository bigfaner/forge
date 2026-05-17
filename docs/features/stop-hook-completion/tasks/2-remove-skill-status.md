---
id: "2"
title: "Remove status transition from quick.md and auto-push from run-tasks.md"
priority: "P1"
estimated_time: "30m"
dependencies: ["1"]
scope: "all"
breaking: false
type: "refactor"
mainSession: false
---

# 2: Remove status transition from quick.md and auto-push from run-tasks.md

## Description

After `forge feature complete --if-done` is verified working (Task 1), remove the now-redundant post-completion responsibilities from skill documentation:

1. **quick.md**: Remove the "Status Transition: Approved → Completed" section (lines ~122-136). The Stop hook now handles this automatically after quality-gate passes.
2. **run-tasks.md**: Remove "Step 5: Auto Git Push" section (lines ~120-153). The Stop hook now handles auto-push after status commit.

Both changes are safe because:
- The Stop hook fires after quality-gate passes, which is the same condition the skill code relied on
- The hook handles both quick and full pipeline modes
- The hook is idempotent — running after skill code would produce a harmless no-op

## Reference Files

- `docs/proposals/stop-hook-completion/proposal.md` — Source proposal
- `plugins/forge/commands/quick.md` — Contains "Status Transition: Approved → Completed" section to remove
- `plugins/forge/commands/run-tasks.md` — Contains "Step 5: Auto Git Push" section to remove

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/quick.md` | Remove "Status Transition: Approved → Completed" section (lines ~122-136) |
| `plugins/forge/commands/run-tasks.md` | Remove "Step 5: Auto Git Push" section (lines ~120-153) |

## Acceptance Criteria

- [ ] `quick.md` contains no "Status Transition: Approved → Completed" section
- [ ] `quick.md` still contains the "Status Transition: Draft → Approved" section (only the post-completion transition is removed)
- [ ] `run-tasks.md` contains no "Auto Git Push" or "gitPush" references
- [ ] `run-tasks.md` Step 4 (Post-Completion summary) remains intact — only Step 5 is removed
- [ ] Verify: `grep -c "status.*Completed" plugins/forge/commands/quick.md` returns 0
- [ ] Verify: `grep -c "gitPush\|git push" plugins/forge/commands/run-tasks.md` returns 0

## Hard Rules

- Do NOT modify any other section of quick.md or run-tasks.md
- Do NOT remove the "Draft → Approved" transition in quick.md — that is still needed at Step 2
- Do NOT remove Step 4 (Post-Completion summary) in run-tasks.md — only Step 5 (auto-push)
- Match existing formatting and indentation in both files

## Implementation Notes

1. **quick.md**: The "Status Transition: Approved → Completed" section is between Step 4 (Execute Tasks) and the Error Handling table. After removal, the Error Handling section should directly follow Step 4.

2. **run-tasks.md**: Step 5 is the last numbered step. After removal, the document ends at the end of Step 4 / Post-Completion section. Check if the line `Do NOT run e2e tests outside Step 3` at line 154 is part of Step 5 or a standalone note.

3. **Risk**: Removing these sections means the Stop hook is now the SOLE mechanism for post-completion. If Task 1 implementation has bugs, status won't be committed. Mitigation: Task 1 acceptance criteria must pass before this task starts (dependency enforced).
