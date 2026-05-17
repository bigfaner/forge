---
status: "completed"
started: "2026-05-17 02:40"
completed: "2026-05-17 02:42"
time_spent: "~2m"
---

# Task Record: 4 Update skill docs for renamed task IDs and auto config

## Summary
Updated all skill docs and guide files to use renamed task IDs (T-test-5 → T-specs-1, T-quick-5 → T-quick-specs-1) and added Auto-Behavior Configuration section to guide.md documenting the new auto config block in .forge/config.yaml.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/consolidate-specs/SKILL.md
- plugins/forge/hooks/guide.md
- plugins/forge/commands/quick.md

### Key Decisions
- Added auto config section to guide.md after All-Completed Hook section (natural placement since auto.gitPush ties into post-completion flow)
- Kept T-quick-1~4 range (only renamed T-quick-5 to T-quick-specs-1) to minimize unnecessary churn

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] Zero remaining references to T-test-5 or T-quick-5 in plugins/forge/
- [x] All renamed IDs use consistent new names: T-specs-1, T-quick-specs-1
- [x] Auto-behavior config documented in guide.md

## Notes
无
