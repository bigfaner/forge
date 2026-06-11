---
id: "4"
title: "迁移 justfile-integration journey"
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

# 4: 迁移 justfile-integration journey

## Description
将 `forge-cli/tests/justfile-integration/` 迁移到顶层 `tests/justfile-integration/`。该 journey 包含 4 个测试文件（execution、forge_detection、init_justfile、mixed_cli）+ 1 个 main_test.go + 4 个 contracts。

## Reference Files
- `docs/proposals/consolidate-advanced-tests/proposal.md` — Scope > In Scope: 5 个独有 journey 迁移
- `forge-cli/tests/justfile-integration/` — 源 journey（4 test files + main_test.go + 4 contracts）
- `tests/task-lifecycle/main_test.go` — 目标初始化模式参考

## Acceptance Criteria
- [ ] `tests/justfile-integration/` 包含全部 4 个迁移的测试文件，import 使用 `testkit "forge-tests/testkit"`
- [ ] `main_test.go` 使用 ForgeBinary init 模式
- [ ] `contracts/` 目录的 4 个合约文件正确迁移
- [ ] `just test` 包含此 journey 且通过

## Implementation Notes
- `mixed_cli_test.go` 可能涉及跨 scope 的 CLI 调用，迁移后需确认 testkit 中 `RunCLIExitCode` 和 `ForgeBinary` 的行为与源 testkit 一致
- 该 journey 的 contracts 数量最多（4 个），直接复制
