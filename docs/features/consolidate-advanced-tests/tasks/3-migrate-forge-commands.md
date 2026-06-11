---
id: "3"
title: "迁移 forge-commands journey"
priority: "P1"
estimated_time: "1h"
complexity: "medium"
dependencies: [1]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 3: 迁移 forge-commands journey

## Description
将 `forge-cli/tests/forge-commands/` 迁移到顶层 `tests/forge-commands/`。该 journey 包含 4 个测试文件（discovery、e2e_commands、forge_info_commands、forge_init_install_just）+ 1 个 main_test.go + 3 个 contracts。

## Reference Files
- `docs/proposals/consolidate-advanced-tests/proposal.md` — Scope > In Scope: 5 个独有 journey 迁移
- `forge-cli/tests/forge-commands/` — 源 journey（4 test files + main_test.go + 3 contracts）
- `tests/task-lifecycle/main_test.go` — 目标初始化模式参考

## Acceptance Criteria
- [ ] `tests/forge-commands/` 包含全部 4 个迁移的测试文件，import 使用 `testkit "forge-tests/testkit"`
- [ ] `main_test.go` 使用 ForgeBinary init 模式
- [ ] `contracts/` 目录的 3 个合约文件正确迁移
- [ ] `just test` 包含此 journey 且通过

## Implementation Notes
- 该 journey 测试文件数量最多（4 个），需逐一检查 import 路径中是否有除 testkit 外的 forge-cli 内部包引用
- 注意 `forge_init_install_just_test.go` 可能依赖 justfile 相关的外部工具，测试失败时区分是迁移问题还是环境问题
