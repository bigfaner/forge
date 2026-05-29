---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["2", "3"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the derive-fix-task-type feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 2-add-doc-fix-template

- [ ] `doc-fix.md` template exists at `forge-cli/pkg/task/templates/doc-fix.md`
- [ ] Template contains fix instructions scoped to doc-type failures: no code quality gates, no test execution, only markdown/content fixes
- [ ] `GetTaskTemplate("doc.fix")` returns the template content without error


### 3-update-skill-files-derivation

- [ ] Error-handling instructions in task-executor.md, execute-task.md, run-tasks.md, submit-task/SKILL.md use derivation rule: extract `TASK_CATEGORY` from claim output, map doc/eval → `doc.fix`, coding/test/validation/gate → `coding.fix`
- [ ] Derivation rule table documented in at least one canonical location (run-tasks.md or execute-task.md) for agent reference
- [ ] `TYPE` and `TASK_CATEGORY` documented as extractable fields from `forge task claim` output in skill files
- [ ] `grep -rn "type coding\.fix" plugins/forge/ --include="*.md"` returns zero matches in error-handling contexts (informational mentions of `coding.fix` as a valid type or in reclassification examples are acceptable)


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/derive-fix-task-type/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/derive-fix-task-type/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
