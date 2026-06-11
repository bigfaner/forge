---
id: "9"
title: "Add fix-task template variables to all creation points"
priority: "P1"
estimated_time: "1.5h"
dependencies: [8]
type: "doc"
mainSession: false
---

# 9: Add fix-task template variables to all creation points

## Description

All fix-task creation points (10+ locations across skill/command docs) are missing `--var SOURCE_FILES/TEST_SCRIPT/TEST_RESULTS` parameters. The `ApplyVars` function in Go code expects these, so fix-tasks fail at render time. Also add integration test impact assessment guidance for breaking tasks. Covers Cluster 4 (issues D1-D2):

1. **run-tasks.md**: 4 fix-task creation points need `--var` parameters added
2. **task-executor.md**: Fix-task creation needs `--var` parameters
3. **execute-task.md**: Fix-task creation needs `--var` parameters + `--description`
4. **submit-task/SKILL.md**: Recovery fix-task needs `--var` parameters
5. **quick-tasks/SKILL.md + breakdown-tasks/SKILL.md**: Add guidance for integration test impact assessment when `breaking: true`

The fix-task grouping should be by test suite (directory), not by problem type.

## Reference Files
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Problem` — Evidence D1 (all fix-task creation points missing --var), D2 (fix-task grouped by problem type not suite)
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Proposed-Solution` — Cluster 4 description
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Success-Criteria` — SC for all fix-task creation points passing --var

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/run-tasks.md` | 4 creation points: add `--var SOURCE_FILES/TEST_SCRIPT/TEST_RESULTS` |
| `plugins/forge/agents/task-executor.md` | Fix-task creation: add `--var` |
| `plugins/forge/commands/execute-task.md` | Fix-task creation: add `--var` + `--description` |
| `plugins/forge/skills/submit-task/SKILL.md` | Recovery fix-task: add `--var` |
| `plugins/forge/skills/quick-tasks/SKILL.md` | Breaking task IT impact assessment guidance |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Breaking task IT impact assessment guidance |

## Acceptance Criteria
- [ ] Every fix-task creation point in run-tasks.md includes `--var SOURCE_FILES=... --var TEST_SCRIPT=... --var TEST_RESULTS=...`
- [ ] task-executor.md fix-task creation includes all three `--var` parameters
- [ ] execute-task.md fix-task creation includes all three `--var` + `--description`
- [ ] submit-task/SKILL.md recovery includes all three `--var`
- [ ] quick-tasks and breakdown-tasks SKILL.md have breaking task IT impact assessment guidance
- [ ] Fix-task grouping guidance specifies by test suite (directory), not problem type

## Hard Rules
- All three template variables (`SOURCE_FILES`, `TEST_SCRIPT`, `TEST_RESULTS`) are mandatory at every creation point
- Fix-task grouping is by test suite directory — same directory = same task

## Implementation Notes
- The `--var` syntax follows the pattern: `forge task add --type coding.fix --var SOURCE_FILES="<paths>" --var TEST_SCRIPT="<test>" --var TEST_RESULTS="<results>"`
- For IT impact assessment guidance: when a task has `breaking: true`, the description must include which integration test fixtures are affected
