---
status: "completed"
started: "2026-05-17 15:14"
completed: "2026-05-17 15:31"
time_spent: "~17m"
---

# Task Record: 6 Migrate all 10 skills to forge testing CLI

## Summary
Migrated all 10 skills from forge profile CLI to forge testing CLI. Replaced all forge profile command references with forge testing equivalents. Replaced capabilities terminology with interfaces. Updated eval test-cases rubric scoring dimension from capabilities to interfaces. Removed all v2 profile name references (go-test, web-playwright, rust-test, java-junit, maestro, pytest). Updated profile-detection.md reference doc to use language keys and interfaces terminology.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/references/shared/profile-detection.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/templates/validate-ux-task.md
- plugins/forge/skills/eval/SKILL.md
- plugins/forge/skills/eval/rubrics/test-cases.md
- plugins/forge/skills/gen-test-cases/SKILL.md
- plugins/forge/skills/gen-test-scripts/SKILL.md
- plugins/forge/skills/graduate-tests/SKILL.md
- plugins/forge/skills/init-justfile/SKILL.md
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/quick-tasks/templates/validate-ux-task.md
- plugins/forge/skills/run-e2e-tests/SKILL.md
- plugins/forge/skills/tech-design/SKILL.md

### Key Decisions
- Step 0 sections in all skills renamed from 'Resolve Profile' to 'Resolve Language and Strategy' (or 'Resolve Language') instead of keeping 'Profile' naming
- Profile manifest concept replaced with strategy files loaded via forge testing get commands
- Error handling changed from asking user to pick profile names to asking user to add languages config
- No backward-compat fallbacks added per hard rule -- clean migration only

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 83.4%

## Acceptance Criteria
- [x] grep -r 'forge profile' plugins/ returns zero matches (excluding test files)
- [x] grep -rE '\bcapabilities\b' plugins/ returns zero matches in skill instruction text (excluding test files and comments)
- [x] gen-test-scripts calls forge testing get generate
- [x] run-e2e-tests calls forge testing get run
- [x] graduate-tests calls forge testing get graduate
- [x] gen-test-cases uses interfaces terminology
- [x] breakdown-tasks uses per-language task expansion
- [x] quick-tasks uses per-language task expansion
- [x] init-justfile calls forge testing get justfile
- [x] tech-design has no profile selection step, uses auto-detect
- [x] No skill references v2 profile names
- [x] eval-test-cases dynamic scoring dimension changes from capabilities to interfaces

## Notes
All 10 skills from D6 table migrated plus eval SKILL.md and test-cases rubric. The /quick skill is not a standalone skill file (it is the forge:quick skill that orchestrates brainstorm + quick-tasks + run-tasks), so its internal references are covered by the quick-tasks migration. profile-detection.md reference doc fully rewritten with language keys and interfaces.
