---
id: "5"
title: "迁移 skill-ops journey"
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

# 5: 迁移 skill-ops journey

## Description
将 `forge-cli/tests/skill-ops/` 迁移到顶层 `tests/skill-ops/`。该 journey 包含 4 个测试文件（clean_code_skill、forensic、plugin_content、prompt）+ 1 个 main_test.go + 3 个 contracts。

## Reference Files
- `docs/proposals/consolidate-advanced-tests/proposal.md` — Scope > In Scope: 5 个独有 journey 迁移
- `forge-cli/tests/skill-ops/` — 源 journey（4 test files + main_test.go + 3 contracts）
- `tests/task-lifecycle/main_test.go` — 目标初始化模式参考

## Acceptance Criteria
- [ ] `tests/skill-ops/` 包含全部 4 个迁移的测试文件，import 使用 `testkit "forge-tests/testkit"`
- [ ] `main_test.go` 使用 ForgeBinary init 模式
- [ ] `contracts/` 目录的 3 个合约文件正确迁移
- [ ] `just test` 包含此 journey 且通过

## Implementation Notes
- `plugin_content_test.go` 和 `prompt_test.go` 可能读取 plugin 目录下的文件，需确认 `ProjectRoot` 返回的路径是否正确解析到项目根目录（而非 tests/ 目录）
- `forensic_test.go` 可能依赖特定的 session transcript 文件，需确认测试 fixture 路径
