---
created: "2026-05-22"
tags: [testing, architecture]
---

# Task executor 对目标明确的任务做冗余代码搜索

## Problem

task-executor 子代理执行 "删除这些文件" 类任务时，即使 task 文件已明确列出每个目标文件路径，仍然进行大量 grep/glob/read 搜索代码库，导致执行时间显著增加。

## Root Cause

1. task-executor 的通用执行流程包含"理解上下文"步骤，要求先读 task 文件、读 proposal、探索相关代码
2. 该流程是面向所有任务类型的通用设计，没有根据任务信息的完备程度做短路判断
3. 当 task 文件中 `## Implementation Notes` 已精确列出每个目标文件路径时，探索步骤变为纯冗余开销

## Solution

在 task 的 `## Hard Rules` 中添加约束来限制搜索行为：

```markdown
- 目标文件已在 Implementation Notes 中完整列出，跳过代码库搜索步骤
```

## Reusable Pattern

当 task 文件的信息已足够执行时（目标文件明确、变更范围精确），在 Hard Rules 中声明"跳过搜索"约束，避免 task-executor 的通用探索流程浪费时间。

适用于以下任务特征：
- 文件级删除/移动/重命名（目标文件已列出）
- 模板化重构（模式已知，无需发现）
- 配置更新（路径和值已确定）

不适用于：
- bug 修复（需要定位根因）
- 新功能实现（需要理解上下文）
- 影响范围不确定的变更
