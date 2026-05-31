---
id: "10"
title: "CI gate integration (goconst + gofmt + go vet)"
priority: "P1"
estimated_time: "1.5h"
complexity: "medium"
dependencies: [8]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 10: CI gate integration (goconst + gofmt + go vet)

## Description
Phase 2b 收尾：将 `goconst`、`gofmt`、`go vet` 集入 `make lint`，确保 `golangci-lint` 配置中启用 `goconst` linter。修复 `goconst` 报告的所有现有违规（与 Task 8 魔法值提取同步完成后的残余检查）。

## Reference Files
- forge-cli/.golangci.yml: golangci-lint 配置文件，需启用 goconst linter (source: proposal.md#Scope item 13)
- forge-cli/Makefile: 需确保 lint target 包含 goconst、gofmt、go vet (source: proposal.md#Scope item 13)

## Acceptance Criteria
- [ ] `.golangci.yml` 中 `goconst` linter 已启用
- [ ] `make lint` 在 `forge-cli/` 目录下执行通过，包含 `goconst`、`gofmt`、`go vet` 检查
- [ ] `goconst` 报告的所有现有违规已修复（与 Task 8 同步完成后残余的违规）
- [ ] CI 通过（SC-8 的 lint 部分）

## Implementation Notes
- 如果 `.golangci.yml` 不存在，需创建；如果已有，需检查 `linters.enable` 列表
- `goconst` 的最小重复次数阈值建议设为 3（默认值）
