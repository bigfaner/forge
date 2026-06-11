---
id: "2"
title: "迁移 error-handling 和 scope-resolution journeys"
priority: "P1"
estimated_time: "1h"
complexity: "low"
dependencies: [1]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 2: 迁移 error-handling 和 scope-resolution journeys

## Description
将 `forge-cli/tests/` 中的 2 个轻量级 journey（error-handling、scope-resolution）迁移到顶层 `tests/`。每个 journey 包含 1 个测试文件 + 1 个 main_test.go + contracts 目录。迁移涉及：更新 import 路径（`forge-cli/tests/testkit` → `forge-tests/testkit`）、重写 main_test.go 使用 ForgeBinary init 模式。

## Reference Files
- `docs/proposals/consolidate-advanced-tests/proposal.md` — Scope > In Scope: 5 个独有 journey 迁移
- `forge-cli/tests/error-handling/` — 源 journey（1 test file + main_test.go + 3 contracts）
- `forge-cli/tests/scope-resolution/` — 源 journey（1 test file + main_test.go + 2 contracts）
- `tests/task-lifecycle/main_test.go` — 目标初始化模式参考（ForgeBinary init）

## Acceptance Criteria
- [ ] `tests/error-handling/` 包含迁移的测试文件，import 使用 `testkit "forge-tests/testkit"`
- [ ] `tests/scope-resolution/` 包含迁移的测试文件，import 使用 `testkit "forge-tests/testkit"`
- [ ] 两个 journey 的 `main_test.go` 统一为 ForgeBinary init 模式（`_ = testkit.ForgeBinary`）
- [ ] `just test` 包含这两个 journey 且全部通过

## Implementation Notes
- 目标 main_test.go 模式：`import testkit "forge-tests/testkit"`，`TestMain` 中 `_ = testkit.ForgeBinary`，无需手动构建二进制
- 源测试文件中引用 `forgeBinary` 变量的地方需改为 `testkit.ForgeBinary`，引用 `testkit.RunCLIExitCode` 等函数的地方需确认新 testkit 中对应函数签名一致
- contracts 是 markdown 文档，直接复制即可
