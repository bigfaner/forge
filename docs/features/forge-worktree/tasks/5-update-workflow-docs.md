---
id: "5"
title: "Update WORKFLOW.md with forge worktree commands"
priority: "P2"
estimated_time: "30m"
dependencies: ["4"]
type: "documentation"
mainSession: false
---

# 5: Update WORKFLOW.md with forge worktree commands

## Description

Update the existing "Using Git Worktree" section in `forge-cli/docs/WORKFLOW.md` to reference the new `forge worktree` commands instead of manual git steps.

## Reference Files
- `docs/proposals/forge-worktree/proposal.md` — Source proposal
- `forge-cli/docs/WORKFLOW.md` — Existing workflow documentation (Section 9, Option 2)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/docs/WORKFLOW.md` | Replace manual worktree steps with `forge worktree start/list/resume/remove` commands |

## Acceptance Criteria

- [ ] Section 9 "Using Git Worktree" updated to show `forge worktree start <slug>` as the primary workflow
- [ ] All four commands documented with examples
- [ ] Manual git worktree steps kept as fallback reference

## Hard Rules

- Keep existing document structure and formatting style
- Both WORKFLOW.md and WORKFLOW.zh.md must be updated

## Implementation Notes

- The existing section shows: `git worktree add ../auth-login feature/auth-login` → `cd ../auth-login && forge task claim`
- Replace with: `forge worktree start auth-login` → feature auto-detected → `forge worktree remove auth-login` when done
