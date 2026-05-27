---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["5"]
type: "doc.review"
scope: "all"
---

Review documentation quality for the auto-eval-config feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 5-skill-config-check
- [ ] brainstorm 在 `auto.eval.proposal` 为 true 时跳过 AskUserQuestion，直接运行 eval-proposal
- [ ] brainstorm 在 `auto.eval.proposal` 为 false 时保持 AskUserQuestion
- [ ] write-prd 在 `auto.eval.prd` 为 true 时跳过 AskUserQuestion，直接运行 eval-prd
- [ ] write-prd 在 `auto.eval.prd` 为 false 时保持 AskUserQuestion
- [ ] ui-design 在 `auto.eval.uiDesign` 为 true 时跳过 AskUserQuestion，直接运行 eval-ui
- [ ] ui-design 在 `auto.eval.uiDesign` 为 false 时保持 AskUserQuestion
- [ ] tech-design 在 `auto.eval.techDesign` 为 true 时跳过 AskUserQuestion，直接运行 eval-design
- [ ] tech-design 在 `auto.eval.techDesign` 为 false 时保持 AskUserQuestion
- [ ] 4 个 skill 使用相同的 config check 模板（代码审查验证一致性）
- [ ] CLI 不可用时（退出码非零）回退到 AskUserQuestion


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/auto-eval-config/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/auto-eval-config/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
