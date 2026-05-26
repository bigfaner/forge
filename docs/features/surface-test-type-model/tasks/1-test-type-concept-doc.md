---
id: "1"
title: "编写测试类型概念参考文档"
priority: "P0"
estimated_time: "2h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: 编写测试类型概念参考文档

## Description

创建 `docs/reference/test-type-model.md`，定义 Surface → Test Type 映射模型。这是整个功能的基石文档，后续所有术语更新任务都以本文档为权威参考。

文档内容基于 proposal.md 中定义的映射表、分类标准、语义定义和验证维度，但以面向 agent 的参考文档格式编写，便于 skill 文件引用。

## Reference Files
- `docs/proposals/surface-test-type-model/proposal.md#Proposed-Solution` — 定义了 Surface → Test Type 映射表、分类标准和语义定义
- `docs/proposals/surface-test-type-model/proposal.md#Technical-Direction` — task type 三段式命名规则和 justfile alias 方案
- `docs/proposals/surface-test-type-model/proposal.md#Requirements-Analysis` — Key Scenarios 和 Edge Cases 定义用户侧体验

## Affected Files

### Create
| File | Description |
|------|-------------|
| `docs/reference/test-type-model.md` | 测试类型概念参考文档 |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] 文档包含 5 种 surface 的 Test Type 名称（EN + CN）、验证维度和执行模型
- [ ] 文档包含分类标准声明：一级分类键（Surface）和二级属性（测试范围：功能测试/端到端测试）
- [ ] 文档包含每种测试类型的语义定义（如 CLI 功能测试 vs Web 端到端测试的区别）
- [ ] 文档包含 "e2e" 术语的使用约束：仅用于 Web/Mobile surface 的端到端测试上下文
- [ ] 文档 frontmatter 的 `domains` 字段包含 `testing`、`surface`、`test-type` 关键词

## Hard Rules

- 文档以中文撰写，技术术语附英文原文
- 篇幅控制在一页纸以内（映射表 + 分类标准 + 语义定义 + 使用约束）

## Implementation Notes

参考 proposal.md 的 "Proposed Solution" 部分，但重新组织为 agent 可快速查阅的参考格式（非 proposal 叙述格式）。映射表是文档核心，应放在最前面。
