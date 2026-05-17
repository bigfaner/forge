---
status: "completed"
started: "2026-05-17 11:05"
completed: "2026-05-17 11:06"
time_spent: "~1m"
---

# Task Record: 1 Fix SKILL.md hardcoded .ts references in Step 3.5 and post-generation

## Summary
Replaced 3 hardcoded .ts filename references in gen-test-scripts SKILL.md with profile-manifest-derived references (manifest.templates.helpers). Affected lines: Task Splitting Guard (line 180), Auth Infrastructure point 4 (line 363), and Post-generation helper merge (line 502).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md

### Key Decisions
- Used 'profile manifest templates.* fields' as the generic phrasing in Task Splitting Guard since it's descriptive text, not a direct instruction
- Used 'manifest.templates.helpers' in the two actionable instruction lines (auth verification and post-generation merge) since agents need the exact field path to resolve the correct filename

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] Line ~180 (Task Splitting Guard): Replace hardcoded filenames with generic description referencing profile manifest templates
- [x] Line ~363 (Auth Infrastructure point 4): Replace hardcoded helpers.ts with profile-manifest-derived reference
- [x] Line ~502 (Post-generation helper merge): Replace hardcoded helpers.ts with profile-manifest-derived reference
- [x] grep for remaining hardcoded .ts references in Step 3.5 auth section and post-generation sections returns zero results (excluding playwright-specific branch)
- [x] web-playwright profile behavior is unchanged — references still resolve to helpers.ts via manifest

## Notes
Documentation-only task. No code changes, no tests affected.
