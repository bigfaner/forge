---
id: "2"
title: "Remove claude/claude-c/claude-w recipes from project justfile"
priority: "P1"
estimated_time: "15min"
dependencies: ["1"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Remove claude/claude-c/claude-w recipes from project justfile

## Description

Now that `forge claude` provides a unified entry point, the redundant justfile recipes (`claude`, `claude-c`, `claude-w`) should be removed. The `claude-p` recipe stays (plugin-dir specific).

## Reference Files
- `docs/proposals/migrate-claude-commands/proposal.md` — Source proposal
- `justfile` — Project justfile (lines 1-13)

## Acceptance Criteria

- [ ] `claude`, `claude-c`, `claude-w` recipes removed from project justfile
- [ ] `claude-p` recipe preserved
- [ ] No other justfile content changed

## Hard Rules

- Do NOT remove `claude-p` recipe — it is plugin-dir specific and stays
- Do NOT modify forge standard recipes section

## Implementation Notes

- Simple deletion of lines 3-10 in the justfile
- Verify `claude-p` recipe still works after removal
