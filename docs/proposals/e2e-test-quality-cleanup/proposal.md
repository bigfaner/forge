---
created: "2026-05-16"
author: "faner"
status: Draft
---

# Proposal: E2E Test Quality Cleanup

## Problem

`tests/e2e/` 中大量测试没有实际验证价值：有的检查静态文件是否包含指定文本，有的永远 skip，有的环境不对就全部跳过，有的递归调用导致进程爆炸。

### Evidence

审计 101 个 e2e 测试，**43 个无意义**：

| 反模式 | 数量 | 文件 | 典型问题 |
|---|---|---|---|
| 读静态文件 grep 文本 | 23 | `extract_design_md` (18), `quick_test_slim` (5) | 读 `.md`/`types.go` 检查字符串，不运行任何命令 |
| 重复（root = features/） | 12 | `cli_list_reverse_chronological` (6), `fix_task_claim_priority` (6) | graduate 后未删除源文件 |
| 递归 `go test ./...` | 2 | `simplify_e2e_tests` TC-003/TC-004 | Windows 上 126+ 孤儿进程 |
| 永远 skip（死代码） | 2 | `feature_set_command` TC-016/TC-017 | `t.Skip("requires real git worktree")` |
| 条件 skip 无 fixture | 9 | `cli_lean_output` TC-001~019 | `claimTask()` 在无任务时全部 skip |
| 空洞断言 | 8 | `cli_lean_output` TC-006~011, TC-013, TC-018 | `if cond { assert }` 永不触发 |

**"读文件 grep 文本"的测试为什么没意义：** 它们读 `plugins/forge/commands/*.md` 或 `forge-cli/pkg/prompt/data/*.md` 或 `types.go`，用 `assert.Contains` 检查是否包含特定字符串。改一下源文件文案测试就断，但不验证任何运行时行为。

### Urgency

递归测试可在 Windows 上消耗 6GB+ RAM 导致机器卡死。其余无意义测试污染覆盖率信号，维护成本大于价值。

## Proposed Solution

**只做清理，不做功能增强。** 删除所有无意义的测试文件和测试函数，不新增能力。

## Scope

### In Scope — 删除

1. **删除 `extract_design_md_platform_adapters_cli_test.go`** — 整个文件（18 个测试全部是读静态命令文件 grep 文本）
2. **删除 `cli_list_reverse_chronological_cli_test.go`**（root 副本）— 与 `features/` 下的完全相同
3. **删除 `fix_task_claim_priority_cli_test.go`**（root 副本）— 同上
4. **删除 `cli_lean_output_cli_test.go`** — 整个文件（19 个测试：8 个空洞断言 + 11 个条件 skip 无 fixture）
5. **删除 `simplify_e2e_tests` TC-003/TC-004** — 递归调用 `go test`
6. **删除 `feature_set_command` TC-016/TC-017** — 永远 skip
7. **删除 `quick_test_slim` TC-003/TC-009/TC-010/TC-013/TC-016** — 读静态文件 grep 文本

### In Scope — 防止复发

8. **扩展 `rubrics/test-cases.md`** 增加反模式检测维度（递归调用、无条件 skip、空洞断言、无 fixture 的条件 skip、读静态文件 grep 文本）
9. **增强 `gen-test-scripts` SKILL.md** 在生成阶段明确禁止上述反模式

### Out of Scope

- 修复 `graduate-tests` workflow（将源文件从 features/ 删除）
- 重构条件 skip 测试为自包含 fixture

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 删除后 `just test-e2e` 编译失败 | L | H | 删除后立即 `just test-e2e` 验证 |
| 某些被删测试实际有价值 | L | L | 保留在 git history 中，随时可恢复 |

## Success Criteria

- [ ] `just test-e2e` 编译通过且全部 pass
- [ ] 零个 `t.Skip` 无条件跳过
- [ ] 零个递归 `exec.Command("go", "test"` 调用
- [ ] 零个读静态源文件做文本检查的测试
- [ ] `tests/e2e/` 中无与 `features/` 重复的文件
- [ ] `rubrics/test-cases.md` 包含反模式检测维度
- [ ] `gen-test-scripts` SKILL.md 明确禁止已知反模式
