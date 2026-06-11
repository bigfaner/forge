---
status: "completed"
started: "2026-05-17 00:28"
completed: "2026-05-17 00:33"
time_spent: "~5m"
---

# Task Record: 2 Create per-type instruction files (types/*.md)

## Summary
Created 5 per-type instruction files under plugins/forge/skills/gen-test-cases/types/ (ui.md, tui.md, mobile.md, api.md, cli.md), each containing self-contained Steps 3-4 instructions with classification indicators, TC format, target derivation rules, route validation, and quality rules extracted from the monolithic SKILL.md.

## Changes

### Files Created
- plugins/forge/skills/gen-test-cases/types/ui.md
- plugins/forge/skills/gen-test-cases/types/tui.md
- plugins/forge/skills/gen-test-cases/types/mobile.md
- plugins/forge/skills/gen-test-cases/types/api.md
- plugins/forge/skills/gen-test-cases/types/cli.md

### Files Modified
无

### Key Decisions
- UI and Mobile both include Integration TC generation for existing-page placements since both types can host integrated components
- Antipattern Prevention rules are referenced (not duplicated) in each type file, linking back to the dispatcher's shared rules
- TUI output assertion specificity section added to enforce concrete assertions (exact text, regex, snapshots) per the rubric's Output Assertion Accuracy criteria
- CLI includes explicit disambiguation rules (CLI vs TUI vs developer tooling) carried over from the monolithic SKILL.md

## Test Results
- **Tests Executed**: No (noTest task)
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] 5 instruction files created: types/ui.md, types/tui.md, types/mobile.md, types/api.md, types/cli.md
- [x] Each file has YAML frontmatter with conventions field listing type-specific convention dependencies
- [x] Each file covers type-specific Steps 3-4 completely (classification rules, TC format, target derivation, quality rules)
- [x] No content is lost from the current monolithic SKILL.md
- [x] UI instruction file includes Integration TC generation for existing-page placements
- [x] Convention frontmatter examples: UI -> [testing-ui.md, frontend.md], CLI -> [testing-cli.md]

## Notes
无
