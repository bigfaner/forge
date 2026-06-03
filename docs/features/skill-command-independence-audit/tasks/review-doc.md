---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["1", "10", "11", "13", "2", "3", "4", "5", "12", "6", "7", "8", "9"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the skill-command-independence-audit feature (quick mode).

## Acceptance Criteria
- [ ] All doc task deliverables reviewed against their acceptance criteria
- [ ] Review findings documented (pass/fail per task)

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-inline-gen-contracts
- [ ] gen-contracts/SKILL.md 包含从 gen-journeys/SKILL.md 内联的 surface 检测规则段落（~20 行），带 `<!-- INLINE:origin=... -->` 标记
- [ ] Related Skills、Integration 章节已删除
- [ ] References 中的概念定义（Contract、Outcome、Semantic Descriptors 等 6 个概念）已合并到内联知识作为定义段落
- [ ] gen-contracts 不再包含对 gen-journeys 内部文件的跨 skill 引用


### 10-reduce-breakdown-tasks
- [ ] 与 quick-tasks 共享的 ~150 行内容已精简，保留 breakdown-tasks 专属逻辑
- [ ] breakdown-tasks 仍能完整指导 AI agent 执行 task breakdown


### 11-reduce-eval
- [ ] proposal-only 特性描述已精简，保留通用 eval 功能说明
- [ ] 所有 eval 子类型（17 种 rubric）的触发条件和执行流程完整保留


### 12-reduce-execute-run-tasks
- [ ] claim 格式描述在两个 command 中各自独立且简洁
- [ ] fix-type 表在两个 command 中各自独立且简洁
- [ ] 两个 command 各自完整可执行，无悬挂引用


### 13-delete-related-batch
- [ ] 三个 skill 的 Related Skills、Integration、References 章节已删除
- [ ] 删除的内容均可从正文中隐含推断，无独有信息丢失
- [ ] 仅修改以下文件：consolidate-specs/SKILL.md、run-tests/SKILL.md、ui-design/SKILL.md


### 2-inline-gen-journeys
- [ ] gen-journeys/SKILL.md 包含内联的 Contract 结构定义 + Outcome 语义（~60 行），带 `<!-- INLINE:origin=... -->` 标记
- [ ] Related Skills、Integration 章节已删除
- [ ] References 中的概念定义已合并到内联知识
- [ ] 5 个 per-surface 内联摘要已精简，保留关键差异信息
- [ ] gen-journeys 不再包含对 gen-contracts 内部文件的跨 skill 引用


### 3-inline-gen-test-scripts
- [ ] gen-test-scripts/SKILL.md 包含内联的隔离策略决策表（~40 行），带 `<!-- INLINE:origin=... -->` 标记
- [ ] Related Skills、Integration、References 章节已删除
- [ ] gen-test-scripts 不再包含对 run-tests 内部文件的跨 skill 引用


### 4-create-style-matching-rules
- [ ] `rules/style-matching.md` 已创建，包含各风格的匹配特征摘要
- [ ] SKILL.md 引用新 rules 文件替代对 ui-design/templates/styles/ 的路径说明


### 5-inline-init-justfile
- [ ] init-justfile/SKILL.md 包含内联的 test-type 映射表（~30 行），带 `<!-- INLINE:origin=... -->` 标记
- [ ] justfile 示例已精简，保留关键模式
- [ ] init-justfile 不再包含对 test-guide 内部文件的跨 skill 引用


### 6-inline-fix-bug
- [ ] fix-bug.md 包含内联的模板决策点和 spec 提取规则（~40 行），带 `<!-- INLINE:origin=... -->` 标记
- [ ] Knowledge Review 段落已精简
- [ ] fix-bug 不再包含对 learn 或 consolidate-specs 内部文件的跨引用


### 7-clean-tech-design
- [ ] Related Skills、Integration、References 章节已删除，且内容可从正文隐含推断
- [ ] 4 种 intent 变体的重复展开已精简，保留核心差异
- [ ] Override Signals 表已精简，去除与 write-prd 的重复内容


### 8-clean-write-prd
- [ ] Related Skills、Integration、References 章节已删除
- [ ] 4 种 intent 变体的重复展开已精简，保留核心差异
- [ ] 与 tech-design 共享的 Override Signals 表已精简为 write-prd 专属版本


### 9-clean-quick-tasks
- [ ] ## Integration 段落已删除
- [ ] ## Reference Files 段落已保留（模板占位符替换规则说明）
- [ ] 与 breakdown-tasks 共享的 ~150 行内容已精简，保留 quick-tasks 专属逻辑


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/skill-command-independence-audit/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/skill-command-independence-audit/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
