---
id: "5"
title: "重命名 Go build tag 为 surface-type-specific 标签"
priority: "P1"
estimated_time: "2h"
dependencies: [3, 4]
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 5: 重命名 Go build tag 为 surface-type-specific 标签

## Description
将 `//go:build e2e` 重命名为 surface-type-specific 标签（如 CLI 项目用 `cli-functional`），并同步更新 Convention 文件、init-justfile surface rule recipe 模板、生成的 justfile 和 test-guide 规则。变更传播链：Convention 文件（定义 build tag 规范）→ init-justfile recipe 模板 → 生成的 justfile → run-tests 调用。同时删除所有 deprecated alias（`test-e2e → <surface>-test`）。

## Reference Files
- `proposal.md#Layer-1-Go-代码层术语路径统一` — 第 7 项定义了 build tag 重命名规则和传播链
- `proposal.md#Scope` — In Scope 第 8、11 项覆盖 build tag 重命名和 alias 删除
- `proposal.md#Success-Criteria` — 验证条件：`//go:build e2e` 和 `-tags=e2e` grep 返回 0

## Acceptance Criteria
- [ ] 所有 `//go:build e2e` 替换为 `//go:build cli-functional`（CLI 项目）
- [ ] Convention 文件（`go.md`、`ginkgo.md`、`vitest.md`、`pytest.md`、`junit.md`、`rust.md`、`index.md` 共 7 个）中 `tags=e2e` 替换为 surface-specific 值
- [ ] `init-justfile/templates/` 下 6 个 justfile 模板（`python.just`、`rust.just`、`node.just`、`go.just`、`mixed.just`、`generic.just`）中 build tag 引用更新
- [ ] `test-guide/rules/` 中 build tag 表格更新
- [ ] 所有 deprecated alias（`test-e2e → <surface>-test`）已删除
- [ ] `grep -rn "//go:build e2e" tests/ forge-cli/` 返回 0 结果
- [ ] `grep -rn '\-tags=e2e' justfile plugins/forge/` 返回 0 结果
- [ ] `grep -rn "alias test-e2e" plugins/forge/` 返回 0 结果
- [ ] `go build ./...` 通过

## Hard Rules
- Convention 文件 + init-justfile + justfile 三处必须同步更新，不可遗漏任何一处

## Implementation Notes
- 每种 surface type 的 build tag 命名参考 `test-type-model.md`
- alias 直接删除（v3.0.0 大版本允许破坏性变更），不保留过渡期

### Integration Test Impact
- Affected test suite(s): `tests/` 目录下所有测试文件
- Expected fixture changes: 测试文件 build tag 行更新
- Risk level: medium
