---
status: "completed"
started: "2026-05-17 13:11"
completed: "2026-05-17 13:14"
time_spent: "~3m"
---

# Task Record: 1 Create types/cli.md instruction file

## Summary
Created types/cli.md instruction file for gen-test-scripts with CLI-specific reconnaissance strategy (grep commands for cobra/commander/click/argparse), Fact Table required keys (CLI_BINARY, CLI_COMMAND_*, CLI_FLAG_*), verification method (5 check commands), generation patterns (process execution, assertion patterns, argument/flag testing, environment isolation), and 3 CLI-specific antipattern guards (recursive test invocation, static file text grep, interactive prompts without automation).

## Changes

### Files Created
- plugins/forge/skills/gen-test-scripts/types/cli.md

### Files Modified
无

### Key Decisions
- Modeled structure after gen-test-cases/types/cli.md for cross-skill consistency (frontmatter format, section ordering, classification indicators)
- Added 'Interactive Prompts Without Automation' as a third CLI-specific antipattern beyond the two mentioned in the task (recursive invocation and static file grep), since interactive prompts are a common CLI-specific failure mode not covered by the generic 6 guards
- Included multi-language grep commands (Go/cobra, Node.js/commander, Python/click) in reconnaissance to cover all profiles, not just Go

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] plugins/forge/skills/gen-test-scripts/types/cli.md exists
- [x] Frontmatter declares type: cli and conventions: [testing-cli.md]
- [x] Contains dedicated Reconnaissance Strategy section with CLI-specific search patterns
- [x] Contains Fact Table Required Keys section listing minimum keys for CLI type
- [x] Contains Verification Method section describing how to confirm CLI interface
- [x] Contains Generation Patterns section describing CLI test case translation
- [x] Contains CLI Antipattern Guards section beyond the generic 6
- [x] At least 3 section headings are unique to this file
- [x] Content is grounded in current SKILL.md CLI-specific branches

## Notes
Documentation-only task. No code changes, no tests to run.
