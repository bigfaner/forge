---
id: "3"
title: "列表命令 slug 列动态宽度"
priority: "P2"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 3: 列表命令 slug 列动态宽度

## Description

将 `forge feature list`、`forge lesson`、`forge proposal` 三个命令的 slug/name 列从固定宽度改为动态宽度，根据实际数据计算（取最长 slug 长度），设置最小 30 字符、最大 60 字符边界。

当前状态：feature/proposal 为 30 字符固定宽度，lesson 为 35 字符。最长 slug 为 42 字符（`profile-aware-shared-infra-precise-staging`），导致截断。

## Reference Files
- `docs/proposals/cli-created-field-and-display/proposal.md` — Source proposal

## Acceptance Criteria
- [ ] `forge feature list` 输出中无 slug 被截断（42 字符 slug 完整显示）
- [ ] `forge lesson` 输出中无 name 被截断
- [ ] `forge proposal` 输出中无 slug 被截断
- [ ] 列宽 = clamp(max(30, maxSlugLen + 2), 60)，即最小 30、最大 60
- [ ] `truncateSlug` 函数根据动态宽度截断（超过动态宽度的罕见情况仍截断）

## Hard Rules
- 动态宽度下限 30 字符，上限 60 字符
- 表头、分隔线、数据行的列宽必须一致

## Implementation Notes
- **Feature**: `forge-cli/internal/cmd/feature.go:208-229` — SLUG 列 `%-30s`，`truncateSlug(f.Slug, 30)`
- **Lesson**: `forge-cli/internal/cmd/lesson.go:54-71` — NAME 列 `%-35s`，`truncateSlug(l.Name, 35)`
- **Proposal**: `forge-cli/internal/cmd/proposal.go:61-84` — SLUG 列 `%-30s`，`truncateSlug(p.Slug, 30)`
- `truncateSlug` 定义在 `proposal.go:117-122`，被 feature 复用
- 实现策略：在排序后、渲染前，遍历数据计算 maxSlugLen，然后 clamp 到 [30, 60]，传入渲染逻辑
