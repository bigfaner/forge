---
id: "2"
title: "配置驱动框架系统"
priority: "P0"
estimated_time: "4h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 2: 配置驱动框架系统

## Description

扩展 `.forge/config.yaml` 的配置系统，支持项目声明测试框架和执行命令。Forge 不再硬编码 6 种语言 profile，而是通过配置驱动框架选择。内置模板作为便利默认值，项目可覆盖。

来源：proposal 的"配置驱动框架选择"和 Scope 中的"配置驱动框架系统"。

## Reference Files
- `docs/proposals/contract-journey-test-model/proposal.md` — Source proposal
- `.forge/config.yaml` — 现有配置
- `forge-cli/internal/cmd/config*.go` — 配置加载逻辑
- `plugins/forge/skills/gen-test-scripts/references/` — 现有 language profile 模板

## Acceptance Criteria

- [ ] `.forge/config.yaml` 支持声明 `test-framework`（如 `mocha`、`pytest`、`go-testing`）和 `test-command`（如 `go test ./...`）
- [ ] 声明 `mocha` → gen-test-scripts 生成 `describe/it` 结构；声明 `pytest` → 生成 `def test_` 函数；声明 `go-testing` → 生成 `func Test*` 函数
- [ ] 零配置时使用内置模板默认值，生成结果与声明内置模板名完全一致（diff 为空）
- [ ] 现有 `.forge/config.yaml` 的 `languages` 和 `interfaces` 字段继续工作（向后兼容）

## Hard Rules

- 不硬编码语言名称到框架选择逻辑中
- 内置模板必须可被项目自定义模板完全覆盖
- 配置缺失时必须有合理的默认行为（内置模板）

## Implementation Notes

- 现有 6 个 language profile（`plugins/forge/skills/gen-test-scripts/references/` 下的 `go/`、`javascript/`、`python/`、`java/`、`rust/`、`mobile/`）将作为内置模板基础
- 参考现有 config 加载逻辑 `forge-cli/internal/cmd/config*.go`
