---
status: "completed"
started: "2026-05-18 01:25"
completed: "2026-05-18 01:25"
time_spent: ""
---

# Task Record: 5 Extract shared logic from breakdown-tasks ↔ quick-tasks

## Summary
Extracted shared Type Assignment table, Intent Propagation logic, and Step 0 profile resolution from breakdown-tasks and quick-tasks SKILL.md into three shared reference files under plugins/forge/references/shared/. Both SKILL.md files now reference the shared files using ${CLAUDE_SKILL_DIR}/../../references/shared/ paths. Also fixed pre-existing argument-hints frontmatter in extract-design-md.md (singular to plural) that was blocking test suite.

## Changes

### Files Created
- plugins/forge/references/shared/type-assignment.md
- plugins/forge/references/shared/intent-propagation.md
- plugins/forge/references/shared/step0-profile-resolution.md

### Files Modified
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/commands/extract-design-md.md

### Key Decisions
- Extracted only the common core of Step 0 (2 numbered steps + HARD-RULE) into step0-profile-resolution.md; each skill retains its unique contextual paragraphs (breakdown-tasks keeps 'resolved languages drive' line, quick-tasks keeps 'Language resolution outcome' block)
- Used ${CLAUDE_SKILL_DIR}/../../references/shared/<file>.md paths per Hard Rules, consistent with forge-distribution.md conventions
- Fixed argument-hint to argument-hints in extract-design-md.md to unblock pre-existing test failure from task 1

## Test Results
- **Tests Executed**: Yes
- **Passed**: 1061
- **Failed**: 0
- **Coverage**: 83.5%

## Acceptance Criteria
- [x] Three shared reference files exist in plugins/forge/references/shared/
- [x] Both breakdown-tasks and quick-tasks reference the shared files
- [x] No content is lost - the extracted content matches the original inline versions

## Notes
Also fixed pre-existing argument-hint frontmatter in extract-design-md.md (task 1 scope) to clear the test suite blockage.
