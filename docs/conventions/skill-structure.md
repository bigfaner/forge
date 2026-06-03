---
title: SKILL.md Structure Convention
domains: [skill, structure, splitting, rules, templates]
---

# SKILL.md Structure Convention

## SKILL.md 拆分启发式规则

SKILL.md 文件采用三层结构：流程骨架（SKILL.md）+ 规则细节（rules/）+ 模板资源（templates/）。

### 留在 SKILL.md 的内容（流程骨架层）

1. 所有步骤编号及其描述
2. 条件分支逻辑（"如果 X 则 A，否则 B"）
3. 输入/输出契约定义
4. 对 rules/ 和 templates/ 的引用指令

### 移至 rules/ 的内容（规则细节层）

1. 超过 5 行的规则定义和解释性文本
2. 术语定义和消歧文档
3. 命名约定、路径规范等参考性内容

### 移至 templates/ 的内容（模板资源层）

1. 超过 10 行的输出模板
2. 可复用的代码片段或配置模板
3. 示例输入/输出

### 边界规则

当一段内容同时包含流程指令和规则细节时：流程指令保留在 SKILL.md，规则细节移至 rules/ 并在原位置添加引用路径。

### 约束

- 每个 SKILL.md 行数不超过 500 行（放宽自 350 行，因复杂 skill 的流程骨架天然较长；3 个已知超标文件待拆分：gen-journeys 454 行, init-justfile 451 行, gen-test-scripts 373 行）
- 辅助文件位于 skill 目录内的以下子目录类型中：`rules/`（规则细节）、`templates/`（模板资源）、`experts/`（领域专家定义）、`rubrics/`（评分标准）、`references/`（参考文档）、`data/`（静态数据）、`types/`（类型定义）、`examples/`（示例文件）
- SKILL.md 必须包含完整流程步骤（遵守 skill-self-containment 原则）
- 不改变 skill 的输入/输出契约
