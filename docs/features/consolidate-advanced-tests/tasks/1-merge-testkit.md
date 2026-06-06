---
id: "1"
title: "合并 testkit：补充 3 个缺失函数到顶层 tests/testkit"
priority: "P0"
estimated_time: "1h"
complexity: "low"
dependencies: []
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 1: 合并 testkit：补充 3 个缺失函数到顶层 tests/testkit

## Description
顶层 `tests/testkit/` 缺少 `forge-cli/tests/testkit/` 中的 3 个函数（`RunCLIExitCode`、`ProjectRoot`、`ReadProjectFile`）。这些函数被即将迁移的测试文件引用，必须先补齐以确保迁移后的测试能编译通过。

## Reference Files
- `docs/proposals/consolidate-advanced-tests/proposal.md` — Scope > In Scope: testkit 合并
- `tests/testkit/helpers.go` — 目标文件，新增 3 个函数
- `forge-cli/tests/testkit/helpers.go` — 源文件，待迁移函数的签名和实现参考

## Acceptance Criteria
- [ ] `tests/testkit/helpers.go` 包含 `RunCLIExitCode(args ...string) (int, string)` 函数
- [ ] `tests/testkit/helpers.go` 包含 `ProjectRoot(t *testing.T) string` 函数
- [ ] `tests/testkit/helpers.go` 包含 `ReadProjectFile(t *testing.T, relPath string) string` 函数
- [ ] 原有 `tests/` 下的测试全部通过（`just test`）

## Implementation Notes
- `RunCLIExitCode` 应使用 `ForgeBinary`（而非 `forgeBinaryPath`），与顶层 testkit 的 `RunCLI`/`RunCLIRaw` 保持一致
- `ProjectRoot` 需要使用 `runtime.Caller` 向上查找 `go.mod`，注意顶层 `tests/` 本身就是独立 Go module（`forge-tests`），所以应查找 `tests/go.mod`
- 源文件中 `SetForgeBinary` 不需要迁移——顶层 testkit 通过 `init()` 自动构建二进制，不需要外部设置路径
