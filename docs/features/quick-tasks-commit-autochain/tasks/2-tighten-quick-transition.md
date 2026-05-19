---
id: "2"
title: "Tighten /quick command Step 3→4 transition"
priority: "P1"
estimated_time: "15m"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 2: Tighten /quick command Step 3→4 transition

## Description
The /quick command already sequentially calls quick-tasks (Step 3) then run-tasks (Step 4), but lacks an explicit constraint preventing the agent from pausing or outputting intermediate summaries between steps. Add explicit transition instructions to ensure seamless handoff from planning to execution.

## Reference Files
- `docs/proposals/quick-tasks-commit-autochain/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/quick.md` | Add explicit Step 3→4 transition instructions |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] /quick command Step 3→4 section contains explicit "immediately proceed" instruction
- [ ] Transition uses `<EXTREMELY_IMPORTANT>` or equivalent high-visibility markup
- [ ] Agent is instructed NOT to output intermediate summary between quick-tasks completion and run-tasks start
- [ ] Instruction specifies that run-tasks should begin with no user confirmation needed after quick-tasks + commit succeeds

## Hard Rules
- Must read `docs/conventions/forge-distribution.md` before modifying plugin files

## Implementation Notes
- Key risk: agent may ignore transition instructions. Use strong markup tags (e.g., `<EXTREMELY_IMPORTANT>`) to maximize compliance
- The transition should reference that Step 8 (commit) must succeed before proceeding
