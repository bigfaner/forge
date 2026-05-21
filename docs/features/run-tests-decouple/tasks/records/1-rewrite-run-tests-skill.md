---
status: "completed"
started: "2026-05-21 23:43"
completed: "2026-05-21 23:48"
time_spent: "~5m"
---

# Task Record: 1 Rewrite run-e2e-tests → run-tests as pure executor

## Summary
Renamed run-e2e-tests skill to run-tests and rewrote SKILL.md as a pure executor. The new skill reads execution commands from .forge/config.yaml test.execution node instead of hardcoded just e2e-* commands. Workflow: Load Convention -> Load Config -> Validate Output Flags -> Setup -> Pre-check -> Run -> Parse -> Report -> Teardown. All commands come from config via template variables. Created config-schema reference, renamed template to test-report.md (removed E2E branding), preserved result-parsing.md and failure-diagnosis.md unchanged.

## Changes

### Files Created
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/skills/run-tests/templates/test-report.md
- plugins/forge/skills/run-tests/references/config-schema.md
- plugins/forge/skills/run-tests/rules/result-parsing.md
- plugins/forge/skills/run-tests/rules/failure-diagnosis.md

### Files Modified
无

### Key Decisions
- Config schema placed in skill references/ directory instead of non-existent commands/forge/lib/ path - follows forge-distribution model where skill reference files live alongside the skill
- Template renamed from e2e-report.md to test-report.md with Screenshots section removed (conditional rendering handled by SKILL.md Step 7 instructions)
- Output messages changed from 'E2E Test Results' to 'Test Results' for framework-agnostic branding

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Skill directory renamed: run-e2e-tests -> run-tests
- [x] SKILL.md frontmatter name changed to run-tests
- [x] SKILL.md contains zero hardcoded just or e2e commands
- [x] SKILL.md reads test.execution from .forge/config.yaml via forge config get test.execution
- [x] Template variables defined: {slug} required, {journey}, {test-dir}, {results-dir} optional with defaults
- [x] Escape rule: {{var}} -> literal {var}
- [x] Output-flags consistency validation step exists
- [x] Missing test.execution.run config produces clear error with config example
- [x] Missing {slug} variable produces clear error prompting forge feature <slug>
- [x] Workflow: Load Convention -> Load Config -> Validate -> Setup -> Pre-check -> Run -> Parse -> Report -> Teardown
- [x] Teardown uses state file .forge/test-state.json for cross-session reliability
- [x] result-parsing.md and failure-diagnosis.md preserved unchanged
- [x] test.execution schema added with all fields documented

## Notes
Config schema file placed at plugins/forge/skills/run-tests/references/config-schema.md instead of the task-specified path plugins/forge/commands/forge/lib/config-schema.yaml which does not exist in the current project structure. The references/ location follows the forge-distribution model where skill reference files are co-located with the skill for distribution to user environments.
