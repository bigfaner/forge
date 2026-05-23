---
id: "1"
title: "Add git status to run-tasks Post-Completion"
priority: "P1"
estimated_time: "15m"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Add git status to run-tasks Post-Completion

## Description
After run-tasks completes all tasks, the user has no visibility into the current git state. Add a concise git summary to the Post-Completion section of `run-tasks.md` showing branch info and working tree changes.

## Reference Files
- `docs/proposals/run-tasks-git-status/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/run-tasks.md` | Add git status display instructions to Post-Completion section |

## Acceptance Criteria
- [ ] Post-Completion section instructs the dispatcher to run git commands showing: current branch name, ahead/behind relative to main, and changed/untracked files (`git status --short`)
- [ ] Git commands wrapped in error handling — skip silently on failure
- [ ] Existing Post-Completion content preserved (test task message, artifact commit prohibition)

## Hard Rules
- Do NOT modify the existing prohibition on committing post-loop artifacts
- Do NOT add diff statistics — keep output concise

## Implementation Notes
- Risk: git commands may fail in non-git contexts. Mitigation: wrap in conditional, skip on failure.
