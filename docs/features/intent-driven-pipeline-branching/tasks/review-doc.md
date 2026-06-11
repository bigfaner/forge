---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["1", "2", "3"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the intent-driven-pipeline-branching feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-add-intent-to-proposal-template
- [ ] `plugins/forge/skills/brainstorm/templates/proposal.md` frontmatter 包含 `intent: "new-feature"` 字段，位于 `status` 之后
- [ ] brainstorm SKILL.md 包含 intent 推断步骤：AI 根据 proposal 内容和 task type → intent 映射规则推断 intent（`coding.feature`/`coding.enhancement` → `new-feature`，`coding.cleanup` → `cleanup`，`coding.refactor` → `refactor`，`coding.fix` 按"是否引入新的用户可观测行为"判断）
- [ ] brainstorm 使用 AskUserQuestion 展示推断的 intent，用户可覆盖确认后写入 proposal.md frontmatter
- [ ] 对 `coding.fix` 类型 proposal，intent 推断逻辑正确：引入新用户可观测行为 → `new-feature`，仅内部调整 → `refactor`
- [ ] 混合内容 proposal（既有新行为又有重构）按"是否引入新的用户可观测行为"判断主要 intent，用户可在确认阶段覆盖


### 2-update-write-prd-refactor
- [ ] write-prd SKILL.md 包含 intent 检测逻辑：当 proposal.md frontmatter 的 `intent` 为 `refactor` 时，执行 spec-only PRD 分支
- [ ] spec-only PRD 格式包含三个必需字段：变更范围（affected modules/files）、约束条件（behavioral invariants to preserve）、验证标准（regression acceptance criteria）
- [ ] refactor 分支下不生成 `prd-user-stories.md` 文件


### 3-update-tech-design-refactor
- [ ] tech-design SKILL.md 包含 intent 检测逻辑：当 proposal.md frontmatter 的 `intent` 为 `refactor` 时，执行内部架构侧重分支
- [ ] refactor 分支下不生成 API handbook 文件和 ER 图文件
- [ ] refactor 分支下不生成 `prd-user-stories.md` 文件


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/intent-driven-pipeline-branching/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/intent-driven-pipeline-branching/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
