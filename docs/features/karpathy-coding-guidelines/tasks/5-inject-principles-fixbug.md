---
id: "5"
title: "Inject Karpathy principles into fix-bug command"
priority: "P1"
estimated_time: "45m"
dependencies: []
type: "doc"
mainSession: false
---

# 5: Inject Karpathy principles into fix-bug command

## Description

Add all 4 Karpathy principles to the `fix-bug` command file. This file is more complex than the coding templates — it has 6 workflow steps, `<EXTREMELY-IMPORTANT>` rules, and a knowledge review section. Principles must merge with existing rules rather than coexist.

## Reference Files
- `docs/proposals/karpathy-coding-guidelines/proposal.md` — Source proposal
- `docs/conventions/forge-distribution.md` — MUST read before modifying plugins/forge/ files

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/fix-bug.md` | Add `<CODING_PRINCIPLES>` block with 4 principles; merge overlapping rules in `<EXTREMELY-IMPORTANT>` |

## Acceptance Criteria
- Contains `<CODING_PRINCIPLES>` block (uppercase XML tags) with Think Before Coding + Simplicity First + Surgical Changes + Goal-Driven Execution
- `<CODING_PRINCIPLES>` positioned after the "Core principle" line, before `## Parameters`
- `<EXTREMELY-IMPORTANT>` rules reviewed: "Fix only what the failing tests require" is semantically preserved by Simplicity First + Surgical Changes — if exact overlap, merge into principles; if complementary, keep both
- "no scope creep, no refactoring, no improvements" rule absorbed by Simplicity First + Surgical Changes
- Knowledge Review section (Steps 1-6) untouched
- Workflow step numbering (Steps 1-6) unchanged
- Common Pitfalls table preserved

## Hard Rules
- MUST read `docs/conventions/forge-distribution.md` before modifying files under `plugins/forge/`
- MUST NOT remove the atomic commit requirement from `<EXTREMELY-IMPORTANT>`
- MUST NOT remove the HARD-GATE at Step 2 (reproduce before test)

## Implementation Notes
- The fix-bug command already has strong guardrails (`<EXTREMELY-IMPORTANT>`, `<HARD-GATE>`). The principles should reinforce these, not create instruction conflict.
- "Think Before Coding" maps to Step 1 (Understand) — principle should guide deeper investigation, not add a step.
- This is a plugin file distributed to user environments via cache. Path references use `${CLAUDE_SKILL_DIR}` — but this file doesn't use any path references, so no impact.
