---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["3"]
type: "doc.review"
scope: "all"
---

Review documentation quality for the unify-task-add-type-param feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 2-update-plugin-components
- [ ] No plugin markdown file contains `--template fix-task` or `--template cleanup-task`
- [ ] All 6 files use `--type coding.fix` (or `--type coding.cleanup` where applicable)
- [ ] No changes to template variable flags (`--var`, `SOURCE_FILES`, `TEST_SCRIPT`, `TEST_RESULTS`)


### 3-update-cli-documentation
- [ ] No documentation file contains `--template` in the context of `forge task add`
- [ ] `WORKFLOW.md` flag table lists `--type` with description "Task type, auto-discovers matching template"
- [ ] All example commands use `--type coding.fix` instead of `--template fix-task`
- [ ] `README.md` CLI parameter table reflects new `--type` flag


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/unify-task-add-type-param/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/unify-task-add-type-param/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
