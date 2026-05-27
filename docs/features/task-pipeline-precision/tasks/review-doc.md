---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["3", "4"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the task-pipeline-precision feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 3-quick-tasks-rules

- [ ] Merge rule changed from "time estimation (<30min)" to "independently verifiable" standard
- [ ] AC max 6 rule added: "if a task has >6 AC, the scope is too large, split further"
- [ ] Multi-verb detection rule added: "task descriptions with connectors linking independent actions (rename + flatten + confirm) should be split by functional boundary"
- [ ] Complexity判定 logic: default heuristic (AC≤3 AND no Hard Rules AND Reference Files≤1 → low; AC>6 OR has Hard Rules → high; else → medium) + LLM judgment override guidance
- [ ] Reference Files generation changed from `proposal.md#Section-Title` pointers to inline precise info format (file path + specific change description)
- [ ] 15 coding task cap removed from SKILL.md HARD-GATE section
- [ ] `templates/task.md` frontmatter has `complexity: "{{COMPLEXITY}}"` field with default "medium"
- [ ] `quick.md` command's 15 task cap reference removed from HARD-GATE


### 4-breakdown-tasks-sync

- [ ] Task Splitting Rules paragraph uses "independently verifiable" as the merge standard (same as quick-tasks)
- [ ] AC max 6 rule added (same wording as quick-tasks)
- [ ] Multi-verb detection rule added (same wording as quick-tasks)
- [ ] Complexity判定 logic with LLM override matches quick-tasks wording
- [ ] Reference Files generation changed to inline precise info format (same as quick-tasks)
- [ ] `templates/task.md` frontmatter has `complexity: "{{COMPLEXITY}}"` field with default "medium"
- [ ] breakdown-tasks specific features (Phase & Gate Detection, PRD Coverage Verification) are NOT affected by these changes


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/task-pipeline-precision/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/task-pipeline-precision/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
