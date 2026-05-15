---
status: "blocked"
started: "2026-05-15 01:15"
completed: "N/A"
time_spent: ""
---

# Task Record: 5 eval-ui multi-platform rubrics and selection logic

## Summary
Created 3 independent platform rubrics (rubric-web.md renamed from rubric.md, rubric-mobile.md new, rubric-tui.md new) and modified eval-ui SKILL.md with platform detection and rubric path selection logic. Each rubric has 4 dimensions x 250 points = 1000 total with platform-specific scoring criteria and deduction rules per proposal D8/D9/D10.

## Changes

### Files Created
- plugins/forge/skills/eval-ui/templates/rubric-web.md
- plugins/forge/skills/eval-ui/templates/rubric-tui.md
- plugins/forge/skills/eval-ui/templates/rubric-mobile.md

### Files Modified
- plugins/forge/skills/eval-ui/SKILL.md

### Key Decisions
- rubric.md renamed to rubric-web.md with no content changes -- existing web rubric preserved exactly
- TUI rubric has Visual Specification dimension that enforces ASCII mockup completeness, character palette precision, and color mapping compliance -- directly maps to lesson's 5 structural requirements
- Mobile rubric has Touch Experience and Adaptive Layout dimensions with mobile-specific deduction rules (touch target sizing, landscape/portrait, safe areas)
- Platform detection in SKILL.md uses explicit platform field first, then structural inference from document content, defaulting to web
- Multi-platform features run independent score-revise loops per platform with respective rubrics

## Test Results
- **Tests Executed**: No
- **Passed**: 18
- **Failed**: 1
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] rubric-web.md contains the existing web rubric content (renamed from rubric.md)
- [x] rubric-tui.md has 4 dimensions: Requirement Coverage (250), Terminal Experience (250), Visual Specification (250), Implementability (250) per proposal D9
- [x] rubric-tui.md deduction rules: missing ASCII mockup, pending characters, missing edge cases, vague dimensions
- [x] rubric-mobile.md has 4 dimensions: Requirement Coverage (250), Touch Experience (250), Adaptive Layout (250), Implementability (250) per proposal D10
- [x] rubric-mobile.md deduction rules: touch targets without size, missing landscape/portrait, missing safe area
- [x] eval-ui/SKILL.md detects platform from ui-design document and selects matching rubric file
- [x] Multi-platform features evaluate each platform's ui-design file with its respective rubric

## Notes
Documentation/template-only task (markdown files, no Go code). The 1 test failure in forge-cli/internal/cmd is pre-existing (verified by running tests on clean stash) and unrelated to these changes.
