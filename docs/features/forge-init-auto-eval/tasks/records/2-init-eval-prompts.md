---
status: "completed"
started: "2026-05-28 00:43"
completed: "2026-05-28 00:51"
time_spent: "~8m"
---

# Task Record: 2 autoBehaviorPrompts 新增 4 个 eval 提示

## Summary
Added 4 eval prompts (proposal, prd, uiDesign, techDesign) to autoBehaviorPrompts() in init.go, before the gitPush prompt. Both forge init and forge config init paths share askAutoBehavior() so both are covered.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/init.go
- forge-cli/internal/cmd/init_test.go

### Key Decisions
- Placed eval prompts before gitPush to match task spec ordering
- Used simple bool prompts (no quick/full prefix) matching the single-toggle style of gitPush
- No changes needed for forge config init stdin path -- it already calls askAutoBehavior() which iterates autoBehaviorPrompts()

## Test Results
- **Tests Executed**: Yes
- **Passed**: 5
- **Failed**: 0
- **Coverage**: 66.6%

## Acceptance Criteria
- [x] forge init interaction flow includes 4 eval prompts (proposal/prd/uiDesign/techDesign)
- [x] Prompt defaults match AutoConfigDefaults(): proposal:true, prd:false, uiDesign:true, techDesign:false
- [x] 4 eval prompts are before gitPush prompt
- [x] forge config init (stdin version) covered via shared askAutoBehavior()

## Notes
无
