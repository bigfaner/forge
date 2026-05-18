---
id: "2"
title: "Task templates — remove user confirmation, adapt for auto mode"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "enhancement"
mainSession: false
---

# 2: Task templates — remove user confirmation, adapt for auto mode

## Description

Update the task execution templates to remove user confirmation steps and adapt the workflow for fully automated pipeline execution. These templates control how the task agent runs consolidate-specs during `/run-tasks`.

## Reference Files
- `docs/proposals/auto-consolidate-specs/proposal.md` — Source proposal
- `plugins/forge/skills/breakdown-tasks/templates/consolidate-specs.md` — Full mode task template (T-specs-1)
- `forge-cli/pkg/prompt/data/doc-generation-consolidate.md` — Full mode prompt template
- `forge-cli/pkg/prompt/data/doc-generation-drift.md` — Quick mode prompt template (drift-only)

## Acceptance Criteria

- [ ] `consolidate-specs.md` task template: "early exit or user review" step skips confirmation in pipeline mode
- [ ] `consolidate-specs.md` task template: all-LOCAL items still auto-proceed without blocking
- [ ] `doc-generation-consolidate.md`: prompt template instructs agent to run in non-interactive mode
- [ ] `doc-generation-drift.md`: prompt template instructs agent to run in non-interactive mode
- [ ] Full pipeline `/run-tasks` no longer blocks on consolidate-specs waiting for user input
- [ ] Quick pipeline `/run-tasks` executes drift-only consolidate-specs without blocking

## Hard Rules

- Must preserve the existing behavior when invoked manually (outside pipeline)
- Templates must reference the updated SKILL.md auto-integration behavior from Task 1

## Implementation Notes

- The consolidate-specs task template has a "Skip conditions" section — verify these still work correctly in auto mode
- The prompt templates should pass a non-interactive flag or context to the skill invocation
- Risk: drift auto-fix deleting correct rules → mitigated by separate `[auto-specs]` commit
