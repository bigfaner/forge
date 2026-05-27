---
created: "2026-05-26"
author: "fanhuifeng"
status: Draft
---

# Proposal: Run-Tests Journey Isolation & Pipeline Cleanup

## Problem

测试 pipeline 的 test-run 任务执行全量测试（`just test`），而非当前 feature 的 journey 测试，导致跨 feature 测试泄露。同时 quality-gate 缺乏 surface 编排能力，无法正确管理 dev/probe/test/teardown 生命周期。

### Evidence

1. **P0 阻断性 bug**：`forge-cli/pkg/prompt/data/test-run.md` 引用 `Skill(skill="forge:run-e2e-tests")`，但该 skill 不存在。正确名称为 `forge:run-tests`。
2. **P1 范围泄露（核心问题）**：test-run 任务模板声称执行 "staged test scripts for the {{FEATURE_SLUG}} feature"，但 run-tests skill 执行 `just test`（无 journey 参数），实际运行全部 surface 测试，非当前 feature 的 journey。
3. **P2 冗余任务**：verify-regression 紧跟 test-run，但 quality-gate（Stop hook）在所有任务完成后已执行 `just test`。verify-regression 无增量价值。
4. **P3 quality-gate 缺乏 surface 编排**：quality-gate 的 Phase 3 直接调用 `just test`，不感知 surface type。对于 web/api surface，不管理 dev/probe/teardown 生命周期，server 未启动时静默跳过测试。

### Urgency

P0 导致 test-run 任务必败。P1 导致测试结果不精确。两者使自动化测试 pipeline 不可用。

## Proposed Solution

分四个层次修复：

1. **修复 skill 引用（P0）**：prompt 模板中 `forge:run-e2e-tests` → `forge:run-tests`
2. **Journey 级隔离（P1）**：run-tests skill 通过 `ls docs/features/<slug>/testing/` 发现 journey 列表，对每个 journey 执行 `just test <journey>`。Surface 编排中 dev/probe 执行一次，per-journey 循环 test，teardown 执行一次。
3. **移除 verify-regression（P2）**：从 `GetBreakdownTestTasks()` 和 `GetQuickTestTasks()` 中移除，quality-gate 已提供全量回归保护。
4. **quality-gate surface 编排（P3）**：quality-gate 读取 `.forge/config.yaml` 的 surface 配置，对每个 surface 执行对应的编排序列（web/api: dev→probe→test→teardown；cli/tui: test→teardown）。

### Innovation Highlights

这不是创新，是修 bug + 消除技术债 + 架构补全。唯一的设计决策是 journey 发现机制：选择 `ls docs/features/<slug>/testing/` 目录扫描而非新增 CLI 命令，因为 journey 就是子目录，不需要额外抽象。

## Requirements Analysis

### Key Scenarios

- **Happy path**：feature 有 2 个 journey，run-tests 分别执行 `just test <journey1>` 和 `just test <journey2>`，合并结果生成报告
- **Edge case：无 journey**：`docs/features/<slug>/testing/` 不存在或为空，run-tests 报错提示先运行 gen-journeys
- **Edge case：单 journey**：只有一个 journey，等价于 `just test <journey>`
- **Surface 编排组合**：web/api 的 dev→probe 生命周期包裹所有 journey 测试（启动一次服务，依次跑每个 journey，最后 teardown）
- **quality-gate surface 编排**：多 surface 项目（web + cli），quality-gate 对每个 surface 执行对应编排序列
- **quality-gate 单 surface**：cli-only 项目，quality-gate 执行简化序列（test→teardown）

### Non-Functional Requirements

- Journey 发现必须零依赖（纯 `ls`），不需要额外工具
- Surface 编排序列行为不变，仅 test 步骤从单次 `just test` 改为循环 `just test <journey>`
- quality-gate 的 surface 编排在 Go 代码中实现，直接调用 just recipe

### Constraints & Dependencies

- justfile recipe 已支持 `just test <journey>` 参数（Go `-run` 正则过滤），无需修改 justfile
- quality-gate 是编译的 Go 二进制，run-tests 是 AI skill——两者独立实现相同的编排逻辑
- Surface 编排序列定义在 surface 规则文件（markdown）中，quality-gate Go 代码需要硬编码或从配置读取对应序列

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 无改动成本 | P0 持续阻断，pipeline 不可用 | Rejected: pipeline 完全失效 |
| 仅修 P0+P1 | — | 最小改动范围 | P2 冗余持续，P3 quality-gate 仍不感知 surface | Rejected: 遗留已知缺陷 |
| 新增 `forge test list-journeys` CLI 命令 | 内部 | 结构化输出 | 需改 Go 代码，ls 已足够 | Rejected: 过度工程 |
| **P0+P1+P2+P3 全部修复** | — | 一步到位，pipeline 完整可用 | 改动范围较大 | **Selected: 审计发现的全部问题应一并解决** |

## Feasibility Assessment

### Technical Feasibility

- P0：单行文本替换（prompt 模板）
- P1：修改 run-tests SKILL.md + surface 规则文件（skill 文档），无需改 Go 代码
- P2：修改 `GetBreakdownTestTasks()` 和 `GetQuickTestTasks()`（Go 代码），移除 verify-regression 任务定义和依赖
- P3：修改 quality-gate Go 代码，增加 surface 检测和编排逻辑

### Resource & Timeline

- P0：10min（文本替换）
- P1：1-2h（skill 文档修改 + surface 规则调整）
- P2：1-2h（Go 代码 + 测试更新）
- P3：2-3h（Go 代码实现 surface 编排 + 测试）
- 总计：4-7h

### Dependency Readiness

- justfile recipe 已支持 journey 参数
- Surface 规则文件已有 "Journey 过滤" 章节（待消费）
- quality-gate 已有 just recipe 调用能力（`just.RunCapture`）

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| verify-regression 在 pipeline 内有价值 | XY Detection | Y=确认不破坏已有测试，但 quality-gate 已覆盖全量回归。Overturned: 移除 |
| quality-gate 跑 `just test` 就够了 | Stress Test | web/api surface 需要 dev/probe/teardown 生命周期，直接 `just test` 在 server 未启动时失败或被跳过。Overturned: 需要 surface 编排 |
| journey 发现应该从 `tests/` 目录 | Assumption Flip | `docs/features/<slug>/testing/` 是 journey 定义（source of truth），`tests/` 是生成产物。选择从 docs/ 发现 |

## Scope

### In Scope

- 修复 `forge:run-e2e-tests` → `forge:run-tests` prompt 模板引用
- run-tests skill 增加 journey 发现步骤（ls docs/features/<slug>/testing/）
- run-tests skill 改为 per-journey 执行 `just test <journey>`
- run-tests skill 增加 surface 编排与 per-journey 执行的组合逻辑（dev/probe 启动一次，循环跑 journey，teardown 一次）
- 从 pipeline 移除 verify-regression 任务定义和依赖链
- quality-gate 增加 surface 感知和编排逻辑（Go 代码）
- 更新相关测试

### Out of Scope

- `forge test run-journey` CLI 命令的重构
- 新增 `forge test list-journeys` CLI 命令
- Surface 规则文件的 "Journey 过滤" 标签消费
- justfile recipe 变更（已支持 journey 参数）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| ls 目录扫描在非标准目录结构下失败 | L | M | 前置检查：testing/ 不存在时报明确错误 |
| quality-gate surface 编排逻辑与 run-tests skill 行为不一致 | M | M | 以 surface 规则文件为权威定义，Go 代码镜像其序列 |
| 移除 verify-regression 后缺少 pipeline 内回归保护 | L | L | quality-gate 在 Stop hook 中已提供全量回归 |
| quality-gate 改动引入回归 | M | M | 现有 quality-gate 测试覆盖 + 手动验证 |

## Success Criteria

- [ ] T-test-run 任务执行时正确调用 `forge:run-tests` skill（不再报 skill 不存在）
- [ ] run-tests skill 只执行当前 feature 的 journey 测试（`just test <journey>` per journey）
- [ ] Surface 编排中 dev/probe 执行一次，per-journey 循环 test，teardown 执行一次
- [ ] `GetBreakdownTestTasks()` 不再生成 `T-test-verify-regression` 任务
- [ ] `GetQuickTestTasks()` 不再生成 `T-test-verify-regression` 任务
- [ ] quality-gate 根据 surface type 执行对应的编排序列（非盲跑 `just test`）
- [ ] 现有 `TestGetQuickTestTasks_*` 和 `TestGetBreakdownTestTasks_*` 测试全部通过

## Next Steps

- Proceed to `/quick` to generate and execute tasks
