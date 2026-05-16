---
status: "blocked"
started: "2026-05-16 00:42"
completed: "N/A"
time_spent: ""
---

# Task Record: 1 Create generic eval skill and extract rubric files

## Summary
Created generic eval skill (SKILL.md) with scorer-gate-revise loop parameterized by rubric files. Extracted 9 rubric files from existing eval skill templates into skills/eval/rubrics/ with frontmatter declaring scale, target, iterations, and type. Generic skill handles both 100-point (harness) and 1000-point scales, detects UI platform for eval-ui, and skips reviser when iterations <= 1.

## Changes

### Files Created
- plugins/forge/skills/eval/SKILL.md
- plugins/forge/skills/eval/rubrics/proposal.md
- plugins/forge/skills/eval/rubrics/prd.md
- plugins/forge/skills/eval/rubrics/design.md
- plugins/forge/skills/eval/rubrics/ui-web.md
- plugins/forge/skills/eval/rubrics/ui-mobile.md
- plugins/forge/skills/eval/rubrics/ui-tui.md
- plugins/forge/skills/eval/rubrics/test-cases.md
- plugins/forge/skills/eval/rubrics/consistency.md
- plugins/forge/skills/eval/rubrics/harness.md

### Files Modified
无

### Key Decisions
- Rubric frontmatter carries scale/target/iterations/type so the generic skill reads them dynamically rather than hardcoding per-type defaults
- Single-pass path (Step 3a) for iterations <= 1 avoids the gate+revise loop entirely, used by eval-harness which has no reviser
- UI platform detection logic from eval-ui moved into the generic skill's Step 1.3, resolving --type ui to ui-web/ui-mobile/ui-tui

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 4
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] skills/eval/SKILL.md contains a single generic scorer-gate-revise loop that reads rubric from rubrics/<type>.md
- [x] Each rubric file is self-contained with frontmatter declaring scale, target, iterations, and type
- [x] eval-harness rubric declares scale: 100 (not 1000) and generic skill handles both
- [x] eval-ui rubric resolution: generic skill detects UI platform from manifest/config and selects ui-web, ui-mobile, or ui-tui
- [x] No eval-specific orchestration logic outside the generic skills/eval/SKILL.md

## Notes
4 pre-existing test failures in forge-cli/pkg/task (TestGetQuickTestTasks_PerType_*) unrelated to this documentation task. These failures exist on the clean branch state as well.
