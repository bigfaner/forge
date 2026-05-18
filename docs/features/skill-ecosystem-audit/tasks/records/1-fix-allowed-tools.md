---
status: "completed"
started: "2026-05-17 23:56"
completed: "2026-05-18 00:11"
time_spent: "~15m"
---

# Task Record: 1 Fix allowed-tools field name and format across all skills/commands

## Summary
Fixed allowed-tools field name (underscore to hyphen) and format (JSON array to space-separated string) across all 13 skill/command files. Also updated a docsync test that asserted the old format.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/clean-code.md
- plugins/forge/commands/execute-task.md
- plugins/forge/commands/extract-design-md.md
- plugins/forge/commands/fix-bug.md
- plugins/forge/commands/gen-sitemap.md
- plugins/forge/commands/git-checkout.md
- plugins/forge/commands/git-commit.md
- plugins/forge/commands/init-forge.md
- plugins/forge/commands/quick.md
- plugins/forge/commands/run-tasks.md
- plugins/forge/commands/simplify-skill.md
- plugins/forge/skills/clean-code/SKILL.md
- plugins/forge/skills/init-justfile/SKILL.md
- forge-cli/internal/docsync/extract_design_md_test.go

### Key Decisions
- Updated docsync test to match corrected frontmatter format (allowed-tools: with unquoted tool names) since the test validates frontmatter content

## Test Results
- **Tests Executed**: Yes
- **Passed**: 1061
- **Failed**: 0
- **Coverage**: 90.1%

## Acceptance Criteria
- [x] Zero files contain allowed_tools (underscore)
- [x] Zero files use JSON array format for allowed-tools
- [x] All allowed-tools values use space-separated format

## Notes
无
