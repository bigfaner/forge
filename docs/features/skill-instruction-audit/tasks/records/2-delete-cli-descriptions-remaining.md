---
status: "completed"
started: "2026-05-28 22:53"
completed: "2026-05-28 22:55"
time_spent: "~2m"
---

# Task Record: 2 Delete CLI behavior descriptions from remaining skills/commands

## Summary
Deleted CLI behavior descriptions from 6 skill/command files (gen-journeys, run-tests, eval, forensic, ui-design, quick) per three-category boundary rule

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-journeys/SKILL.md
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/skills/eval/SKILL.md
- plugins/forge/skills/forensic/SKILL.md
- plugins/forge/skills/ui-design/SKILL.md
- plugins/forge/commands/quick.md

### Key Decisions
无

## Document Metrics
6 files modified, ~30 lines deleted/replaced, 6/6 AC passed

## Referenced Documents
- docs/proposals/skill-instruction-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] gen-journeys/SKILL.md has no example output blocks for forge surfaces; Exit Code table remains
- [x] run-tests/SKILL.md has no segment prefix matching; command remains
- [x] eval/SKILL.md has no repeated tool usage explanations
- [x] forensic/SKILL.md has no go build or ~/.zcode-forge-cli references
- [x] ui-design/SKILL.md config check is natural language, not bash script
- [x] quick.md has no behavioral descriptions of run-tasks or brainstorm internals

## Notes
All deletions followed the three-category boundary rule: removed behavior explanations (category 3), preserved instructions (category 1) and output contracts (category 2).
