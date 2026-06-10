---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["1", "3", "4", "7", "8", "2", "5", "6", "9"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the forge-skill-audit feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-sync-eval-data
- [ ] rubric-reference.md journey row: scale=1150, target=975 (matching rubric frontmatter)
- [ ] rubric-reference.md contract row: scale=1100, target=935 (matching rubric frontmatter)
- [ ] rubric-reference.md header has maintenance comment: this file is a secondary cache of rubric scale/target
- [ ] eval-journey.md argument-hint shows --target 975; description lists all 7 dimensions with 1150-point scale
- [ ] eval-contract.md argument-hint shows --target 935; description lists all 8 dimensions with 1100-point scale
- [ ] eval/SKILL.md description reflects actual supported scales (no "100-point" claim)


### 2-remove-dead-path
- [ ] tech-design SKILL.md 仅引用 `docs/proposals/<slug>/proposal.md`，无 `docs/features/<slug>/proposal.md` 路径
- [ ] 搜索所有 skill 确认无其他 skill 引用 `docs/features/<slug>/proposal.md` 死路径


### 3-fix-task-template
- [ ] breakdown-tasks/templates/task.md 使用 `complexity: "{{COMPLEXITY}}"` 和 `type: "{{TYPE}}"` 占位符
- [ ] 模板包含与 quick-tasks/templates/task.md 一致的注释块，列出 COMPLEXITY 和 TYPE 的可选值及默认值


### 4-fix-record-format
- [ ] record-format-coding.md 不再列出 `doc.fix`
- [ ] record-format-doc.md 包含 `doc.fix` 覆盖


### 5-add-inline-markers
- [ ] 4 处 INLINE 引用均标注源文件路径和版本号标记（格式：`<!-- INLINE from <source-path> @ <version> -->`）
- [ ] `grep -r "INLINE" plugins/forge/skills/` 显示所有内联引用均有版本号标记


### 6-unify-config-keys
- [ ] 所有 auto.eval 配置键在 skill markdown 中统一为 kebab-case（`auto.eval.ui-design`, `auto.eval.tech-design`）
- [ ] Implementation notes 明确标注 Go config reader alias 兼容需求为后续任务


### 7-unify-ui-design-eval
- [ ] ui-design SKILL.md auto.eval 部分使用 bash script 模板（三路分支：disabled/skip/run），与 brainstorm/write-prd/tech-design 一致


### 8-cleanup-orphan-files
- [ ] `draft-generation.md` 和 `pattern-extraction.md` 已移至 `skills/test-guide/rules/_deprecated/` 目录
- [ ] SKILL.md 和其他 rules 文件无引用断裂（本就无引用）


### 9-fix-doc-completeness
- [ ] breakdown-tasks SKILL.md intent 读取部分指定完整路径 `docs/proposals/<slug>/proposal.md`
- [ ] test-isolation.md 头部有 `<!-- OWNER: run-tests | CONSUMERS: gen-test-scripts (INLINE) -->` 注释
- [ ] brainstorm SKILL.md Step 5 包含 `{{AUTHOR}}` 赋值指导（git config user.name 或询问用户）
- [ ] write-prd/templates/manifest.md 使用 `{{SLUG}}` 而非 `{{FEATURE_SLUG}}`


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/forge-skill-audit/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/forge-skill-audit/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.

## Acceptance Criteria

- [ ] All acceptance criteria met
