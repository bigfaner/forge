---
id: "6"
title: "合并 task-lifecycle 和 task-type-system 重叠 journeys"
priority: "P1"
estimated_time: "2h"
complexity: "high"
dependencies: [1]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 6: 合并 task-lifecycle 和 task-type-system 重叠 journeys

## Description
将 `forge-cli/tests/` 中的 2 个重叠 journey 的测试文件合并到顶层 `tests/` 对应的已有目录中。与独有 journey 不同，这些目标目录已存在测试文件和 contracts，需要处理文件名冲突和内容共存。

## Reference Files
- `docs/proposals/consolidate-advanced-tests/proposal.md` — Scope > In Scope: 2 个重叠 journey 合并, Key Risks
- `forge-cli/tests/task-lifecycle/` — 源：3 test files + main_test.go + 3 contracts
- `forge-cli/tests/task-type-system/` — 源：2 test files + main_test.go + 3 contracts
- `tests/task-lifecycle/` — 已有：2 test files + main_test.go + 2 contracts
- `tests/task-type-system/` — 已有：1 test file + main_test.go + 3 contracts

## Acceptance Criteria
- [ ] `tests/task-lifecycle/` 包含来自源端的 3 个新测试文件（fix_task_claim_priority_test、submit_test、task_stage_gates_test），import 使用 `forge-tests/testkit`
- [ ] `tests/task-type-system/` 包含来自源端的新测试文件（解决 `task_type_refinement_test.go` 同名冲突），import 使用 `forge-tests/testkit`
- [ ] 两端 contracts 合并后无内容丢失（处理 step-1、step-2 同名冲突）
- [ ] 所有新增测试文件 import `testkit "forge-tests/testkit"`
- [ ] `just test` 包含这两个 journey 且全部通过

## Hard Rules
- 不修改已有测试文件（`tests/task-lifecycle/task_lifecycle_test.go`、`tests/task-lifecycle/task_record_test.go`、`tests/task-type-system/task_type_refinement_test.go`）的内容

## Implementation Notes
- **文件名冲突**：`task-type-system` 两端都有 `task_type_refinement_test.go`，需对比内容。若测试函数名不冲突可合并到同一文件；若冲突则重命名源端文件（如 `task_type_refinement_v2_test.go`）
- **contracts 同名冲突**：task-lifecycle 两端都有 `step-1-task-claim.md` 和 `step-2-task-submit.md`，需对比内容后决定是否需要重命名或合并
- **main_test.go**：目标已存在，不需要源端的 main_test.go（目标已使用正确的 init 模式）
- 源端的 `main_test.go` 不迁移（目标已有正确版本）
