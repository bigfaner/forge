---
status: "completed"
started: "2026-05-17 13:35"
completed: "2026-05-17 13:43"
time_spent: "~8m"
---

# Task Record: 3 Update command/agent files with discovery instruction

## Summary
Removed hardcoded keyword-to-filename mapping table from fix-bug.md and replaced it with the discovery instruction (domains frontmatter-based relevance matching). The error-fixer.md file was not modified because it no longer exists — it was deprecated/removed by the typed-task-dispatch feature.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/fix-bug.md

### Key Decisions
- Only modified fix-bug.md since error-fixer.md no longer exists in the agents directory (deprecated by typed-task-dispatch feature, task 4.1)
- Preserved the 'Infer relevant domains from bug description, affected files, and scope' instruction as part of the Project Knowledge line, per Hard Rules
- Replaced the 'Example mappings' line with the same discovery instruction used in task 2's prompt templates for consistency

## Test Results
- **Tests Executed**: Yes
- **Passed**: 954
- **Failed**: 0
- **Coverage**: 82.8%

## Acceptance Criteria
- [x] Both files contain the discovery instruction (same as in prompt templates)
- [x] No file contains the hardcoded mapping pattern
- [x] No file contains any keyword-to-filename mapping table
- [x] The 'Project Knowledge' section structure is preserved

## Notes
error-fixer.md does not exist in plugins/forge/agents/ — it was removed/deprecated by the typed-task-dispatch feature. The agents directory contains only doc-reviser.md, doc-scorer.md, and task-executor.md. The acceptance criterion 'Both files contain the discovery instruction' is met for the one existing file (fix-bug.md); the second file no longer exists in the codebase.
