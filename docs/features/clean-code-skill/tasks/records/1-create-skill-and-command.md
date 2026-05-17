---
status: "completed"
started: "2026-05-17 17:17"
completed: "2026-05-17 17:33"
time_spent: "~16m"
---

# Task Record: 1 Create clean-code skill and command

## Summary
Created clean-code skill and slash command. SKILL.md defines a 4-step workflow (scope detection via git diff, code cleanup with 5 principles from Anthropic code-simplifier, optional quality gate, cleanup summary). Command file provides /forge:clean-code entry point. Added 5 e2e tests validating file existence, workflow steps, 5 principles, command-skill wiring, and scope constraint.

## Changes

### Files Created
- plugins/forge/skills/clean-code/SKILL.md
- plugins/forge/commands/clean-code.md
- forge-cli/tests/e2e/clean_code_skill_cli_test.go

### Files Modified
无

### Key Decisions
- Modeled skill after Anthropic code-simplifier agent with 5 principles (Preserve Functionality, Apply Project Standards, Enhance Clarity, Maintain Balance, Focus Scope)
- Quality gate is optional: runs just test if available, skips if not, per hard rule in task definition
- Scope is strictly limited to git diff output, enforced by HARD-RULE blocks in SKILL.md
- Large diffs (50+ files) processed in batches of 10-15 to avoid context overflow

## Test Results
- **Tests Executed**: Yes
- **Passed**: 973
- **Failed**: 0
- **Coverage**: 81.0%

## Acceptance Criteria
- [x] plugins/forge/skills/clean-code/SKILL.md exists with complete skill definition
- [x] Skill workflow: scope detection (git diff) -> code cleanup (5 principles) -> quality gate (just test, optional) -> cleanup summary
- [x] Cleanup logic follows code-simplifier 5 principles: Preserve Functionality, Apply Project Standards, Enhance Clarity, Maintain Balance, Focus Scope
- [x] Skill can be invoked standalone via /forge:clean-code
- [x] plugins/forge/commands/clean-code.md exists as slash command entry point

## Notes
Task 1 scope was only skill + command creation. CLI wiring (typeToTemplate entry, prompt template) belongs to task 2.
