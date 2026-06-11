---
created: "2026-05-26"
author: "faner"
status: Draft
---

# Proposal: 清理无效测试用例（激进清理）

## Problem

测试套件中积累了大量无效测试：永久跳过的死测试、引用已删除功能的残留测试、不验证运行时行为的文本检查测试。这些测试不提供任何质量保障价值，反而制造噪音、增加维护成本、误导覆盖率指标。

### Evidence

对 `tests/`（e2e）和 `forge-cli/tests/`（集成）两层测试的全面审计：

| 反模式 | 数量 | 来源 | 具体问题 |
|--------|------|------|---------|
| 无条件 `t.Skip`（死测试） | 44 | `forge-cli/tests/` 中 10 个文件 | `t.Skip("requires manual setup")` 等，永远不执行 |
| 读静态文件 grep 文本 | 23 | `tests/` 中 `extract_design_md` (18), `quick_test_slim` (5) | 读 `.md`/`types.go` 检查字符串，不运行任何命令 |
| 重复测试（root 副本） | 12 | `tests/` 中 `cli_list_reverse_chronological` (6), `fix_task_claim_priority` (6) | graduate 后未删除源文件 |
| 递归 `go test` | 2 | `tests/` 中 `simplify_e2e_tests` TC-003/TC-004 | Windows 上 126+ 孤儿进程，6GB+ RAM |
| 永远 skip | 2 | `tests/` 中 `feature_set_command` TC-016/TC-017 | `t.Skip("requires real git worktree")` |
| 空洞断言 | 8 | `tests/` 中 `cli_lean_output` TC-006~011, TC-013, TC-018 | `if cond { assert }` 永不触发 |
| 引用已删除功能 | 3 | `forge-cli/tests/test-generation/` | 引用已退役的模板/skills |
| 需交互输入无法自动化 | 4 | `forge-cli/tests/forge-commands/forge_info_commands_test.go` | requires interactive stdin |

**总计约 98 个无效测试**，分布在约 20 个文件中。

### Urgency

v3.0.0 分支正在开发中，是 breaking change 窗口。递归测试在 Windows 上可导致机器卡死。无效测试污染覆盖率信号，误导开发决策。

## Proposed Solution

**激进清理**：删除所有识别到的无效测试，清理空文件、空目录、死辅助函数。不新增能力，不做重构，只做减法。

### Innovation Highlights

无创新——标准的技术债清理。洞察是：删除测试比修复测试更安全，因为无效测试提供零保障价值，保留在 git history 中随时可恢复。

## Requirements Analysis

### Key Scenarios

1. 开发者运行 `just test <journey>` — 只执行有实际价值的测试
2. 开发者运行 `just unit-test` — 不受影响（unit test 层未改）
3. CI 运行测试套件 — 更快、更可靠的信号
4. 新开发者阅读测试文件 — 只看到有意义的测试

### Non-Functional Requirements

| NFR | Requirement | Verification |
|-----|-------------|-------------|
| 兼容性 | `just test` 和 `just unit-test` 行为不变 | 两条命令编译通过且 pass |
| 安全性 | 不删除有实际断言价值的测试 | 逐文件审查删除清单 |
| 可维护性 | 清理后的空文件/目录一并删除 | 无空测试文件或空 suite 目录 |

### Constraints & Dependencies

- 必须在 `remove-forge-test-commands` 提案（已 Approved）执行后验证，因为它会删除 `pkg/contract/` 和部分测试文件
- `tests/test-suite-health/` 中的 TC-004 规则（零无条件 skip）目前只覆盖 `tests/`，清理后应扩展到 `forge-cli/tests/`

## Alternatives & Industry Benchmarking

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零风险 | 98 个无效测试持续存在；Windows 递归问题未解决 | Rejected: 成本远大于收益 |
| 保守清理（只删明确死测试） | — | 最小变更 | 保留文本验证和空洞断言测试 | Rejected: 用户选择激进方案 |
| **激进清理（删全部 + 清基础设施）** | 标准技术债清理 | 彻底清理，不留死角 | 可能删除"将来有用"的测试 | **Selected: 用户明确选择** |

## Feasibility Assessment

### Technical Feasibility

文件删除 + 测试函数删除。所有涉及的测试文件都可以通过 `go test -tags=e2e` 编译验证。无外部依赖。

### Resource & Timeline

小范围：约 20 个文件的删除/修改，可在一次会话完成。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "无效测试将来可能修复" | Occam's Razor | 44 个 skip 中大部分标注 "requires manual setup" 已超数月，无人修复，是死代码而非 TODO |
| "文本验证测试提供了某种覆盖" | Assumption Flip | 如果删掉这些测试，不会有任何 bug 遗漏，因为它们不验证运行时行为 |
| "空洞断言测试有隐含价值" | Stress Test | `if cond { assert }` 中 cond 永远为 false——这些测试从未真正执行断言 |

## Scope

### In Scope

#### `forge-cli/tests/` 集成测试清理

1. 删除或修复 44 个无条件 `t.Skip` 测试（10 个文件）
2. 删除引用已删除功能/退役 skills 的测试（`test-generation/gen_test_scripts_test.go`）
3. 删除需要交互输入无法自动化的测试（`forge_info_commands_test.go` 中 4 个 skip）
4. 删除清空后的空测试文件
5. 删除空 suite 目录
6. 清理 testkit 中仅被删除测试引用的辅助函数

#### `tests/` e2e 测试清理

7. 删除读静态文件 grep 文本的测试（23 个：`extract_design_md` 全部 18 个 + `quick_test_slim` 5 个）
8. 删除重复的 root 副本测试文件（12 个：`cli_list_reverse_chronological` 6 个 + `fix_task_claim_priority` 6 个）
9. 删除递归 `go test` 测试（`simplify_e2e_tests` TC-003/TC-004）
10. 删除永久 skip 测试（`feature_set_command` TC-016/TC-017）
11. 删除空洞断言测试（`cli_lean_output` 8 个）
12. 删除整个 `cli_lean_output_cli_test.go`（19 个测试：8 个空洞 + 11 个条件 skip 无 fixture）
13. 删除整个 `tui-ui-design/` 目录（31 个文本验证测试）

#### 防止复发

14. 扩展 `test-suite-health` TC-004 规则覆盖 `forge-cli/tests/`

### Out of Scope

- Co-located 单元测试（`forge-cli/internal/` 和 `forge-cli/pkg/` 中的 `*_test.go`）
- `remove-forge-test-commands` 提案覆盖的源代码删除
- 编写新测试替代被删除的测试
- 修改 CI pipeline 配置
- `tests/test-suite-health/` 元测试文件本身的修改

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 删除了实际有价值的测试 | L | M | 逐文件审查；git history 保留完整历史 |
| 清理后编译失败 | L | H | 每删一批文件后运行 `go test` 验证 |
| testkit 辅助函数被其他测试间接引用 | M | M | 删除前 `grep` 验证无其他引用 |
| 空目录删除遗漏 | L | L | 最终 `find` 验证无空目录 |

## Success Criteria

- [ ] `forge-cli/tests/` 中零个无条件 `t.Skip`
- [ ] `tests/` 中零个读静态源文件做文本检查的测试
- [ ] `tests/` 中零个递归 `exec.Command("go", "test"` 调用
- [ ] `tests/` 中零个与 `features/` 重复的 root 副本文件
- [ ] `tests/` 中零个空洞断言（`if cond { assert }` 且 cond 永远为 false）
- [ ] `just unit-test` 编译通过且全部 pass
- [ ] `just test` 编译通过（e2e 测试）
- [ ] 清理后无空测试文件
- [ ] 清理后无空 suite 目录（排除 `.graduated/` 元数据）
- [ ] testkit 辅助函数无死引用

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
