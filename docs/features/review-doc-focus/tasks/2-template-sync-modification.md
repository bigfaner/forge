---
id: "2"
title: "Template sync modification (autogen + prompt templates)"
priority: "P0"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 2: Template Sync Modification

## Description

同步修改 autogen 模板（`task/data/doc-review.md`）和 agent prompt 模板（`prompt/data/doc-review.md`）。两个模板存在强耦合——autogen 模板定义 AC 汇总区域的结构和位置，agent prompt 依赖该区域定位和读取 AC 数据。必须在同一个 task 中同步修改。

autogen 模板：新增 AC 汇总区域（使用 `{{DOC_TASK_AC}}` 占位符），更新 Discovery Strategy 为 allowlist。
prompt 模板：重构为三步（Load Pre-extracted AC → Discover Target Documents via allowlist → Review & Fix with constraints）。

## Reference Files
- `proposal.md#AC-汇总区域格式` — AC 汇总区域的完整 markdown schema 定义
- `proposal.md#Prompt-模板重构规格` — prompt 模板的三步重构规格和重构步骤
- `proposal.md#Coupling-Constraints` — 两个模板强耦合约束，必须同步修改
- `proposal.md#Key-Risks` — 过滤规则风险和 allowlist 策略说明

## Acceptance Criteria

- [ ] autogen 模板包含 `{{DOC_TASK_AC}}` 占位符位置（由 Task 1 的 renderBody 渲染）
- [ ] autogen 模板 Discovery Strategy 改为 allowlist：仅扫描 `docs/` 子树下 `.md` 文件
- [ ] autogen 模板排除 `tasks/`、`records/`、`manifest.md`、`index.json`
- [ ] prompt 模板 Step 1 改为 "Load Pre-extracted AC"（无 "scan tasks directory" 指令）
- [ ] prompt 模板 Step 2 使用 docs/ allowlist 文档发现策略
- [ ] prompt 模板 Step 3 包含 "仅修改 docs/ 下文件" 的显式禁止约束
- [ ] prompt 模板无任何 "scan tasks directory" 或 "Read the task's acceptance criteria from its .md file" 指令

## Hard Rules

- 两个模板必须在同一个 task 中修改——不能分开提交
- prompt 模板重构后保持 5-step workflow 总体结构（步骤数可能变化但流程逻辑不变）

## Implementation Notes

- 参考 proposal 的 AC 汇总区域格式 subsection 生成 autogen 模板的 AC 区域结构
- prompt 模板重构步骤：删除 Step 1 扫描指令 → 新增 AC 加载指令 → 重写 Step 2 为 allowlist → Step 3 增加禁止修改约束
- allowlist 范围需与当前模板的 `docs/features/{{FEATURE_SLUG}}/` 特性目录一致，不能扩大到全局 docs/
