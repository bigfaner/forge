---
status: "completed"
started: "2026-05-19 14:45"
completed: "2026-05-19 14:47"
time_spent: "~2m"
---

# Task Record: 1 Restructure guide.md into 3 thematic sections

## Summary
Restructured guide.md from 241 lines to 110 lines by removing 6 sections of duplicated/reference content (Skill Workflow diagrams, Quick Mode details, Testing Lifecycle, Evaluation Parameter Exceptions, Knowledge Accumulation, Auxiliary Skills table) and reorganizing remaining content into 3 thematic sections: Directory Conventions, Execution Rules, and Automation Config. All factual rules preserved verbatim.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/hooks/guide.md

### Key Decisions
- Grouped Manifest subsection under Directory Conventions rather than standalone, since it describes what a feature entry point contains (directory-level concern)
- Moved Task-CLI into Execution Rules section since it describes the execution lifecycle flow

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] guide.md is ~100-120 lines (from 241)
- [x] 6 sections removed: Skill Workflow mermaid diagrams, Quick Mode details, Testing Lifecycle, Evaluation Parameter Exceptions, Knowledge Accumulation details, Auxiliary Skills table
- [x] 3 thematic sections present: Directory Conventions, Execution Rules, Automation Config
- [x] All remaining rules preserved accurately — quality gate, scope resolution, auto-config, all-completed hook, task-CLI flow
- [x] No functional behavior change

## Notes
Documentation-only task. No test metrics applicable.
