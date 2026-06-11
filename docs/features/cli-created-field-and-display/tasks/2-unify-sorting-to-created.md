---
id: "2"
title: "统一 feature/lesson 排序为 created 降序"
priority: "P1"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 2: 统一 feature/lesson 排序为 created 降序

## Description

将 `forge feature list` 和 `forge lesson` 的排序逻辑统一为 frontmatter `created` 字段降序，mtime 作为 fallback。`forge proposal` 已使用 `created` 排序，无需修改。

当前问题：
- Feature 按 manifest 文件 mtime 排序
- Lesson 按文件 mtime 排序（`Created` 字段已解析但未用于排序）
- Proposal 已按 `created` 排序

## Reference Files
- `docs/proposals/cli-created-field-and-display/proposal.md` — Source proposal

## Acceptance Criteria
- [ ] `forge feature list` 按 manifest frontmatter `created` 字段降序排列（缺少 `created` 时 fallback 到 mtime）
- [ ] `forge lesson` 按 frontmatter `created` 字段降序排列（缺少 `created` 时 fallback 到 mtime）
- [ ] `forge proposal` 排序逻辑不变（验证仍为 `created` 降序）
- [ ] 缺少 `created` 的旧文档 fallback 到 mtime，不报错、不跳过

## Hard Rules
- `created` 字段格式为 `YYYY-MM-DD`，可直接按字符串字典序比较
- mtime fallback 必须在 `created` 为空或不存在时触发，不能静默忽略有效值

## Implementation Notes
- **Feature**: `forge-cli/internal/cmd/feature.go:194-196` — 当前按 `ManifestMtime` 排序。需要解析 manifest.md 的 frontmatter 获取 `created` 字段，加入 Feature 结构体
- **Lesson**: `forge-cli/pkg/lesson/lesson.go:111-121` — 当前按 file mtime 排序。`Created` 字段已在 `Lesson` 结构体中存在（从 frontmatter 解析），但排序未使用。改为按 `Created` 排序
- **Proposal**: `forge-cli/internal/cmd/proposal.go:52-54` — 已按 `Created` 降序。无需修改，仅验证
- Key Risk: 存量 feature manifest 无 `created` 字段 → mtime fallback 已覆盖（低影响）
