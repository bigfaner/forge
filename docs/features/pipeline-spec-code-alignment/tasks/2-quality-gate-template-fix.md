---
id: "2"
title: "Fix quality_gate.go and template.go for cleanup-task and fix-task"
priority: "P0"
estimated_time: "2h"
dependencies: [1]
surface-key: "cli"
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 2: Fix quality_gate.go and template.go for cleanup-task and fix-task

## Description

Fix multiple issues in quality_gate.go and template.go:

1. **quality_gate.go inferSurface** (~line 352): Currently uses only the first extracted source file to infer surface. Multi-surface projects may assign tasks to wrong surface. Fix to use all source files.

2. **quality_gate.go fix-task splitting**: Fix-tasks are currently grouped by problem type. Change to group by test suite (same directory = same task), enabling parallel execution and bounded scope.

3. **quality_gate.go dual-source truth**: EstimatedTime and other fields exist in both Go code and task template frontmatter. Establish opts (Go code) as the authoritative source; frontmatter uses `{{...}}` template variables.

4. **template.go + coding.cleanup.md**: Cleanup-task (for lint/format failures) currently uses `Breaking: true`, triggering a full test gate. Fix to `Breaking: false` with `EstimatedTime: "15min"`. compile+lint still run.

## Reference Files
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Problem` — Evidence B4 (inferSurface single file), D2-D6 (fix-task quality issues)
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Proposed-Solution` — Clusters 2 and 5 descriptions
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Key-Risks` — Risk of cleanup-task Breaking:false missing regressions
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Success-Criteria` — SC for fix-task by suite, cleanup-task params, dual-source truth

## Acceptance Criteria
- [ ] `inferSurface` uses all source files, not just the first
- [ ] Fix-tasks are grouped by test suite (same directory), not problem type
- [ ] EstimatedTime comes from Go opts as authoritative source
- [ ] Cleanup-task uses `Breaking: false` and `EstimatedTime: "15min"`
- [ ] `coding.cleanup.md` template frontmatter has `breaking: false`
- [ ] Existing tests pass (`go test ./...`)

## Hard Rules
- compile and lint gates MUST still run for cleanup-tasks even with `Breaking: false`
- Do not change the quality gate interface — only internal behavior

## Implementation Notes
- For fix-task by suite: group failing tests by directory. Each group becomes one fix-task. Bottom-line rule: tests in same directory stay in one task.
- For dual-source truth: when template.go generates task .md from opts, use `{{estimated_time}}` from opts rather than hardcoding. The `coding.cleanup.md` template file should use a placeholder or match the opts value.
- Breaking:false means the task doesn't block downstream tasks on test failure, but compile+lint still enforce.
