---
id: "1a"
title: "重命名 E2E 常量和路径函数为 surface-neutral 名称"
priority: "P0"
estimated_time: "30m"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 1a: 重命名 E2E 常量和路径函数为 surface-neutral 名称

## Description
将 `constants.go` 中 `E2ETestsBaseDir` → `TestBaseDir`、删除 `E2EStagingDir`/`E2EGraduatedDir`；将 `paths.go` 中删除 `GetE2EStagingDir()`/`GetE2EGraduatedMarker()`/`GetE2ETargetDir()`，新增 `GetTestResultsDir()`/`GetTestConfigPath()`。更新所有引用点（quality_gate、autogen、testrunner、docsync_test）。

## Reference Files
- `proposal.md#Layer-1-Go-代码层术语路径统一` — 第 1 项定义常量重命名映射
- `proposal.md#新目录结构` — 常量值映射

## Acceptance Criteria
- [x] `E2ETestsBaseDir` 重命名为 `TestBaseDir`，值为 `"tests"`
- [x] `E2EStagingDir` / `E2EGraduatedDir` 常量已删除
- [x] `GetE2EStagingDir()` / `GetE2EGraduatedMarker()` / `GetE2ETargetDir()` 函数已删除
- [x] 新增 `GetTestResultsDir()` 和 `GetTestConfigPath()` 函数
- [x] 所有引用点已更新
- [x] `go build ./...` 通过

## Hard Rules

## Implementation Notes
- 已在 commit 3f5f08f2 中完成
