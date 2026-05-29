---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["6", "1", "2", "3", "4", "5"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the plugin-consistency-audit feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-inventory-structural-scan
- [ ] 全部 21 个 skill 枚举完成，每个 skill 列出完整文件清单（SKILL.md + templates/ + rules/ + data/ + examples/ + types/）
- [ ] 全部 18 个 command 枚举完成，列出其内部文件引用
- [ ] 1 个 agent (task-executor) 枚举完成，列出引用文件
- [ ] hooks/guide.md 枚举完成，列出引用的脚本路径
- [ ] Layer 1 完成：SKILL.md 中引用的每个路径已与文件系统交叉验证，REFERENCE 类问题已记录
- [ ] 孤立文件（存在于目录中但未被 SKILL.md 引用）已识别并记录
- [ ] 报告包含基准 commit hash


### 2-skills-audit-batch-a
- [ ] 7 个 skill 各自的 SKILL.md 已全文读取并提取结构化摘要（步骤、约束、引用路径、字段名）
- [ ] 每个 skill 的所有关联文件（templates/rules/data/examples/types）已逐一与 SKILL.md 摘要比对
- [ ] 使用关键词强度映射表检查"必须/应该/可选/禁止"的一致性，记录 CONFLICT 类问题
- [ ] 多步骤组件的步骤时序已验证（Layer 3），TIMING 类问题已记录
- [ ] **有效性验证**: run-tests skill 的 `rules/env-check.md` Playwright 硬编码已被识别为 P1 级 CONFLICT
- [ ] 每条问题按报告 schema 记录：`{component, file_path, layer, category, severity, description, fix_suggestion, confidence}`


### 3-skills-audit-batch-b
- [ ] 7 个 skill 各自的 SKILL.md 已全文读取并提取结构化摘要（步骤、约束、引用路径、字段名）
- [ ] 每个 skill 的所有关联文件（templates/rules/data）已逐一与 SKILL.md 摘要比对
- [ ] 使用关键词强度映射表检查"必须/应该/可选/禁止"的一致性，记录 CONFLICT 类问题
- [ ] 多步骤组件的步骤时序已验证（Layer 3），TIMING 类问题已记录
- [ ] 每条问题按报告 schema 记录：`{component, file_path, layer, category, severity, description, fix_suggestion, confidence}`


### 4-skills-audit-batch-c
- [ ] 7 个 skill 各自的 SKILL.md 已全文读取并提取结构化摘要（步骤、约束、引用路径、字段名）
- [ ] 每个 skill 的所有关联文件（templates/rules/data）已逐一与 SKILL.md 摘要比对
- [ ] 使用关键词强度映射表检查"必须/应该/可选/禁止"的一致性，记录 CONFLICT 类问题
- [ ] 多步骤组件的步骤时序已验证（Layer 3），TIMING 类问题已记录
- [ ] 每条问题按报告 schema 记录：`{component, file_path, layer, category, severity, description, fix_suggestion, confidence}`


### 5-commands-agent-hooks-audit
- [ ] 全部 18 个 command 文件已读取，内部流程步骤的时序和引用一致性已验证
- [ ] task-executor agent 文件已读取，指令之间的矛盾、冗余、时序问题已检查
- [ ] hooks/guide.md 已读取，引用的脚本路径存在性已验证（Layer 1 REFERENCE）
- [ ] hooks/guide.md 中脚本参数描述与实际脚本声明的一致性已验证（Layer 2 CONFLICT）
- [ ] 每条问题按报告 schema 记录：`{component, file_path, layer, category, severity, description, fix_suggestion, confidence}`


### 6-consolidated-report
- [ ] 所有审计发现已合并去重，按 P0→P3 排序输出
- [ ] **有效性验证通过**: run-tests 的 `rules/env-check.md` Playwright 硬编码问题在最终报告中作为 P1 级 CONFLICT 出现
- [ ] 五类分类（CONFLICT/REDUNDANT/TIMING/REFERENCE/INCOMPLETE）均至少有 1 个实例；若某类为 0，列出所有含多步骤流程的组件清单并确认已逐一验证
- [ ] 报告包含基准 commit hash、AI 模型版本、审计参数（temperature 等）
- [ ] 误报率抽检方案已定义：随机抽取 ≥20% 的 P0/P1 问题清单，标注待人工验证


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/plugin-consistency-audit/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/plugin-consistency-audit/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
