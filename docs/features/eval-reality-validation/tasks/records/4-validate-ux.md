---
status: "completed"
started: "2026-05-17 01:31"
completed: "2026-05-17 01:35"
time_spent: "~4m"
---

# Task Record: 4 validate-ux Implementation

## Summary
Created validate-ux eval type: rubric with 10 dimensions (1000 pts, 2 blocks: UX Rules + PRD Flow), two-phase pre-processing logic in eval SKILL.md, CLI/Web/TUI project type detection via forge profile capabilities, ux-snapshot.md format spec, 7 operation impact verification types, and task templates for both breakdown-tasks and quick-tasks.

## Changes

### Files Created
- plugins/forge/skills/eval/rubrics/validate-ux.md
- plugins/forge/skills/breakdown-tasks/templates/validate-ux-task.md
- plugins/forge/skills/quick-tasks/templates/validate-ux-task.md

### Files Modified
- plugins/forge/skills/eval/SKILL.md

### Key Decisions
- Two-phase model: Phase 1 (main session) compiles, runs, collects ux-snapshot.md; Phase 2 (doc-scorer) evaluates the snapshot against the rubric. Matches the harness type pattern.
- Rubric structured as 2 blocks: Block A (UX Rules, 400 pts) for standalone quality checks and Block B (PRD Flow Validation, 600 pts) for flow correctness and impact verification.
- TUI limited to non-interactive scenarios only (initial render, help, invalid input) per hard rule.
- Task templates positioned after T-test-5/T-quick-5 with mainSession: true since Phase 1 requires main session execution.

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] validate-ux.md rubric exists with scale: 1000, target: 700, iterations: 1, type: validate-ux, context frontmatter
- [x] Rubric defines 10 dimensions with correct point values totaling 1000
- [x] Rubric context declares conventions: [ux, cli, api] and business-rules: auto
- [x] Eval SKILL.md includes validate-ux pre-processing as two-phase
- [x] Project type detection uses forge profile capabilities: cli -> CLI, web-ui -> Web, tui -> TUI
- [x] Eval SKILL.md defines ux-snapshot.md format with Flow steps, Standalone Checks, and Effect Verification sections
- [x] Eval SKILL.md defines PRD-to-operation translation strategies per project type
- [x] Eval SKILL.md defines 7 operation impact verification types
- [x] Task templates exist in both breakdown-tasks and quick-tasks positioned after T-test/T-quick steps
- [x] validate-ux uses iterations=1, never triggers revise loop

## Notes
Documentation-only task. No tests to run. Hard rules verified: did not modify doc-scorer.md or doc-reviser.md; rubric totals exactly 1000 points; TUI limited to non-interactive; worktree/ temp dir requirement specified; iterations: 1 mandatory.
