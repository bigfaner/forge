---
id: "7"
title: "Overhaul justfile templates for two-layer test model"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
scope: "all"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 7: Overhaul justfile templates for two-layer test model

## Description

重写 6 个 `templates/*.just` 模板和项目根 `justfile`：将 `test` recipe 重命名为 `unit-test`，淘汰 `e2e-test`/`e2e-setup`/`e2e-verify`，新增 `test`（高级测试，接受可选 journey 参数）和 `test-setup`。为 Go/Node/Python/Rust 各生成语言对应的 `unit-test` recipe。

## Reference Files
- `proposal.md#Proposed-Solution` — defines recipe structure: unit-test (language-level), test (surface-level with optional journey), test-setup, probe
- `proposal.md#Requirements-Analysis` — Recipe 参数签名约定: `test journey=''` signature, init-justfile per-language/per-surface generation
- `proposal.md#Constraints-&-Dependencies` — language-specific unit-test generation for Go/Node/Python/Rust
- `proposal.md#Key-Risks` — template comment strategy to build new mental model: `# unit-test: language-level unit tests`, `# test: surface-level advanced tests`

## Acceptance Criteria
- All 6 templates (`generic.just`, `go.just`, `node.just`, `python.just`, `rust.just`, `mixed.just`) have:
  - `unit-test` recipe with language-appropriate command (Go: `go test ./...`, Node: `npm test`, Python: `pytest`, Rust: `cargo test`)
  - `test` recipe with optional `journey` parameter (`test journey=''`), dispatching by surface type
  - `test-setup` recipe (optional, for Playwright install, DB seed, etc.)
  - No `e2e-test`, `e2e-setup`, or `e2e-verify` recipes
  - Comments: `# unit-test: language-level unit tests` and `# test: surface-level advanced tests`
- Root `justfile`: `test` → `unit-test`; new `test` recipe; `ci` recipe updated
- `test` recipe passes `journey` parameter to underlying test framework when non-empty

## Hard Rules
- `test` recipe 必须接受可选第一个参数 `journey`（just 语法：`test journey=''`）
- 不硬编码 Playwright——`test` recipe 按 surface 分发

## Implementation Notes
- Go 模板 `unit-test` 不加 `-race` flag（race detection 由 CI 环境控制，不嵌入 per-task gate）
- Node 模板需要考虑 `npm test` vs `yarn test` 等变体，保持与现有模板风格一致
- Root justfile 的 `test` recipe 可以调用 `just unit-test && just test` 或分开定义
