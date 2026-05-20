---
id: "6"
title: "Resume with claude -c session restore"
priority: "P2"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 6: Resume with claude -c session restore

## Description

Enhance `forge worktree resume` to use `claude -c <slug>` instead of just `claude --dangerously-skip-permissions`. This restores the previous Claude Code session context for the given slug, allowing users to continue where they left off.

## Reference Files
- `docs/proposals/worktree-experience/proposal.md` — Source proposal
- `forge-cli/internal/cmd/worktree.go` — Resume command implementation (runWorktreeResume)
- `forge-cli/internal/cmd/claude.go` — Claude launch utilities

## Acceptance Criteria
- [ ] `forge worktree resume <slug>` launches claude with `-c <slug>` flag for session restore
- [ ] `--dangerously-skip-permissions` is still passed (maintain current auto-approval behavior)
- [ ] If `claude -c` is not supported (old claude version), falls back to current behavior (launch without -c)
- [ ] The slug is used as the session name for `-c`

## Hard Rules
- Must validate that `claude -c` is supported before using it. If claude CLI doesn't support `-c`, fall back gracefully.
- Do NOT change the directory resolution logic — resume still navigates to the worktree directory

## Implementation Notes
- Current resume launches: `claude --dangerously-skip-permissions` in the worktree directory.
- New launch: `claude -c <slug> --dangerously-skip-permissions` in the worktree directory.
- Key risk: Claude Code's `-c` session name format may have restrictions. The slug format (e.g., "worktree-experience") should be safe. Need to verify this works.
- Fallback: if `claude -c` fails (unrecognized flag), catch the error and re-launch without `-c`.
