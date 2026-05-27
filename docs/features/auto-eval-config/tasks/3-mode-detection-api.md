---
id: "3"
title: "Add mode detection API"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
surface-key: ""
surface-type: ""
breaking: false
type: "coding.feature"
mainSession: false
---

# 3: Add mode detection API

## Description
新增 `forge config get mode` 命令，返回当前管道模式（`quick`/`full`/`none`）。提取 `feature_complete.go` 中的 quick mode 判定逻辑为独立函数。CLI 通过解析当前工作目录路径推断 feature slug（匹配 `.forge/features/<slug>` 模式），然后检查 feature 目录下是否存在 `proposal.md`。

## Reference Files
- `docs/proposals/auto-eval-config/proposal.md#Constraints-Dependencies` — mode 检测 API 规格：3 种返回值、pwd 路径解析、feature slug 提取
- `docs/proposals/auto-eval-config/proposal.md#Key-Scenarios` — mode 检测场景（feature 内、feature 外）
- `docs/proposals/auto-eval-config/proposal.md#Success-Criteria` — PR-2 中 mode detection 的 3 条 SC

## Acceptance Criteria
- [ ] `forge config get mode` 在 feature 目录内 + proposal.md 存在时返回 `"quick"`
- [ ] `forge config get mode` 在 feature 目录内无 proposal.md 时返回 `"full"`
- [ ] `forge config get mode` 在非 feature 目录时返回 `"none"`
- [ ] 路径解析处理 Windows 和 Unix 路径分隔符
- [ ] 路径解析处理 symlink 场景（通过 `filepath.EvalSymlinks`）

## Hard Rules
- mode 检测通过 `forge config get mode` CLI 路径访问，不暴露为 Go 包级别 API
- 返回值严格为 `"quick"` / `"full"` / `"none"` 三选一

## Implementation Notes
- 从 `feature_complete.go` 提取 quick mode 判定逻辑为 `DetectPipelineMode(projectRoot string) string` 函数
- CLI 层在 `config.go` 的 `getConfigHandler` 中添加 `"mode"` 特殊 key 处理
- 路径解析：`filepath.EvalSymlinks` 解析 symlink → `strings.Contains` 匹配 `.forge/features/` → 提取 slug → 检查 `proposal.md`
