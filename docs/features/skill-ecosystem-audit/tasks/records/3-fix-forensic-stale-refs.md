---
status: "blocked"
started: "2026-05-18 00:25"
completed: "N/A"
time_spent: ""
---

# Task Record: 3 Fix forensic hardcoded paths and stale record-task reference

## Summary
Replaced 4 hardcoded developer-specific paths in forensic SKILL.md with generic <project-hash> patterns and ${CLAUDE_SESSION_ID} variable. Fixed stale /record-task reference to /submit-task in consolidate-specs template.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/forensic/SKILL.md
- plugins/forge/skills/breakdown-tasks/templates/consolidate-specs.md

### Key Decisions
- Used generic <project-hash> placeholder instead of literal hash, with comment explaining ${CLAUDE_SESSION_ID} usage for session-specific paths

## Test Results
- **Tests Executed**: Yes
- **Passed**: 2289
- **Failed**: 1
- **Coverage**: 90.5%

## Acceptance Criteria
- [x] Zero hardcoded user-specific paths in forensic SKILL.md
- [x] Zero record-task references (excluding submit-task)
- [x] Forensic SKILL.md uses generic path patterns or CLI commands that work for any user

## Notes
1 pre-existing test failure in TestExtractDesignMd_ArgumentHintsIncludesPlatform (unrelated to this task, fails on clean branch too). All compile/fmt/lint pass.
