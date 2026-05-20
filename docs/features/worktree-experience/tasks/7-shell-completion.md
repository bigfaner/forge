---
id: "7"
title: "Shell completion for start/remove/resume subcommands"
priority: "P2"
estimated_time: "2h"
dependencies: ["1", "3", "6"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 7: Shell completion for start/remove/resume subcommands

## Description

Add Cobra `ValidArgsFunction` dynamic shell completion to the worktree subcommands. `start` completes with unfinished proposal/feature slugs; `remove` and `resume` complete with existing worktree slugs. This eliminates manual slug typing and context-switching to look up names.

## Reference Files
- `docs/proposals/worktree-experience/proposal.md` — Source proposal
- `forge-cli/internal/cmd/worktree.go` — Worktree commands
- `forge-cli/internal/cmd/worktree_test.go` — Existing tests
- `forge-cli/pkg/git/worktree_list.go` — Worktree listing

## Acceptance Criteria
- [ ] `forge worktree start <TAB>` shows unfinished proposal and feature slugs
- [ ] `forge worktree remove <TAB>` shows existing worktree slugs
- [ ] `forge worktree resume <TAB>` shows existing worktree slugs
- [ ] Completion response time < 200ms
- [ ] Completion works for bash, zsh, and fish shells (Cobra handles this automatically)
- [ ] No completion on `list`, `push`, `status` subcommands (they don't take slug args, or use different args)

## Hard Rules
- Use Cobra's `ValidArgsFunction` pattern — do NOT use static `ValidArgs`
- Completion functions must handle errors gracefully (return empty list, not error to shell)

## Implementation Notes
- Start completion: scan `docs/proposals/*/proposal.md` and `docs/features/*/manifest.md` for unfinished items. Reuse the listing logic from Task 2 (interactive mode).
- Remove/Resume completion: use `git.ListWorktrees()` to get forge-managed worktrees.
- Cobra ValidArgsFunction signature: `func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)`
- Registration: `cmd.ValidArgsFunction = myCompletionFunc` in `init()`.
- No external dependencies needed — Cobra handles shell-specific completion script generation.
