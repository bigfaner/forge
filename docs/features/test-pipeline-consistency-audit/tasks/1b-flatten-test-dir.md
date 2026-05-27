---
id: "1b"
title: "扁平化 tests/ 物理目录结构"
priority: "P0"
estimated_time: "30m"
dependencies: ["1a"]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 1b: 扁平化 tests/ 物理目录结构

## Description
将物理目录从 `tests/e2e/` 扁平化为 `tests/`：移动 `tests/e2e/config.yaml` → `tests/config.yaml`，移动 `tests/e2e/results/` → `tests/results/`，删除 `tests/e2e/features/`（staging）和 `tests/e2e/.graduated/`（graduation）目录，确保 `tests/e2e/` 目录不再存在。

## Reference Files
- `proposal.md#新目录结构` — 物理目录扁平化的具体前后对比
- `proposal.md#Scope` — In Scope 第 2 项覆盖物理路径扁平化

## Acceptance Criteria
- [ ] `tests/config.yaml` 存在（探针配置文件）
- [ ] `tests/results/` 目录存在（测试结果目录）
- [ ] `tests/e2e/features/` 和 `tests/e2e/.graduated/` 目录已删除
- [ ] `tests/e2e/` 目录不再存在
- [ ] `go build ./...` 通过

## Hard Rules

## Implementation Notes
- commit 3f5f08f2 已删除 `tests/e2e/` 但未创建 `tests/config.yaml` 和 `tests/results/`
- `tests/config.yaml` 是 e2eprobe 探针配置文件，需要确认内容格式
- `tests/results/` 是测试结果输出目录，`.gitkeep` 或空目录即可

### Integration Test Impact
- Affected test suite(s): `forge-cli/pkg/e2eprobe/`
- Expected fixture changes: 探针配置路径测试
- Risk level: low
