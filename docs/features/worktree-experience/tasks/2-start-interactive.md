---
id: "2"
title: "Start interactive mode: -i flag for proposal/feature selection"
priority: "P1"
estimated_time: "2h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: Start interactive mode: -i flag for proposal/feature selection

## Description

Add `-i/--interactive` flag to `forge worktree start` that presents a selectable list of unfinished proposals and features. Users pick one, and the slug is auto-filled. This eliminates the need to manually look up and type slugs.

## Reference Files
- `docs/proposals/worktree-experience/proposal.md` — Source proposal
- `forge-cli/internal/cmd/worktree.go` — Start command implementation
- `forge-cli/internal/cmd/worktree_test.go` — Existing tests

## Acceptance Criteria
- [ ] `forge worktree start -i` lists all unfinished proposals (in `docs/proposals/*/`) and features (in `docs/features/*/`) with their status
- [ ] User can select one item from the list; the slug is extracted and used as the argument
- [ ] When `-i` is used, the `<slug>` positional arg becomes optional (one or the other)
- [ ] When both `-i` and `<slug>` are provided, `<slug>` takes precedence (ignore -i)
- [ ] Empty list (no proposals or features) prints a helpful message and exits
- [ ] Selection prompt works in a terminal context (TTY detection for non-interactive environments)

## Hard Rules
- Do NOT add external dependencies for TUI/selection — use fmt.Scanln or a simple numbered-list approach
- The slug derivation must match existing conventions (directory name under proposals/ or features/)

## Implementation Notes
- "Unfinished" proposals: any proposal directory that exists. Optionally filter by `status: Draft` in frontmatter.
- "Unfinished" features: features whose manifest status is not `completed`.
- Need a function to scan `docs/proposals/` and `docs/features/` directories and list slugs.
- Simple selection: print numbered list, read integer from stdin.
