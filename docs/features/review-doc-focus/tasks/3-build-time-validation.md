---
id: "3"
title: "Build-time AC validation (warnings + empty AC + title tolerance)"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 3: Build-time AC Validation

## Description

在 AC 提取管线中添加构建时验证：标题匹配容错（大小写 + 中文别名）、缺失 AC 时输出 warning 日志、汇总区域显示占位文本、零 AC 特性输出特别警告。

## Reference Files
- `proposal.md#Constraints-&-Dependencies` — 标题匹配容错策略定义
- `proposal.md#Key-Risks` — AC 提取解析不稳定风险的缓解策略和 free-review 模式
- `proposal.md#Success-Criteria` — AC 覆盖率验证和 warning 输出要求

## Acceptance Criteria

- [ ] 标题匹配支持 `## Acceptance Criteria`、`## Acceptance criteria`（大小写差异）、`## 验收标准`（中文别名）
- [ ] section 不存在时输出 warning 日志 `[WARN] task <name> has no Acceptance Criteria section`
- [ ] AC 为空时汇总区域显示 `> No acceptance criteria defined.`
- [ ] 所有 doc 任务均缺 AC 时输出 `[WARN] feature has no AC for any doc task`
- [ ] `BuildIndex()` 生成后验证 `DocTaskCriteria` 的 key 集合与 doc 任务列表完全匹配

## Hard Rules

- 使用 soft-fail（warning + 占位文本），不中断构建
- 匹配失败时输出 warning，不静默跳过

## Implementation Notes

- 标题匹配使用 `strings.EqualFold` 或 `strings.ToLower` 比较
- 中文别名 `## 验收标准` 需硬编码支持
- warning 输出使用现有的 `result.Warnings` 累积机制
