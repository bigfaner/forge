---
id: "3"
title: "Update knowledge extraction in fix-bug, write-prd, tech-design"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 3: Update knowledge extraction in fix-bug, write-prd, tech-design

## Description

Update the knowledge extraction flow in 3 skill/command files to read `auto.knowledgeSave` config and skip user confirmation when the mode-specific value is `true`. When enabled, the extraction proceeds silently: Step 5 (AskUserQuestion) is replaced with auto-save-all, and Step 6 writes directly.

## Reference Files
- `docs/proposals/auto-knowledge-save/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/fix-bug.md` | Add config check before Step 5; auto-save when `true` |
| `plugins/forge/skills/write-prd/rules/knowledge-extraction.md` | Add config check before Step 5; auto-save when `true` |
| `plugins/forge/skills/tech-design/rules/knowledge-extraction.md` | Add config check before Step 5; auto-save when `true` |

## Acceptance Criteria
- [ ] All 3 files read `auto.knowledgeSave` via `forge config get auto.knowledgeSave` before Step 5
- [ ] When mode is `true`, Step 5 (AskUserQuestion) is skipped and all candidates auto-saved to Step 6
- [ ] When mode is `false`, existing AskUserQuestion flow is preserved unchanged
- [ ] The config read and branching logic is consistent across all 3 files

## Hard Rules
- The config read pattern: `forge config get auto.knowledgeSave` → parse `"quick:<val> full:<val>"` → check current mode value
- Must use the same branching pattern in all 3 files for consistency

## Implementation Notes
- Insert the config check as a new "Step 4.5" or modify Step 5 to branch on the config value
- The mode context (quick vs full) is available from the skill's execution context (quick-tasks vs full pipeline)
- Each file already has the same Extraction Flow structure (Steps 1-6), so the modification point is identical
- Pattern for auto-save: skip AskUserQuestion, treat as if user selected "all", proceed to Step 6 write
