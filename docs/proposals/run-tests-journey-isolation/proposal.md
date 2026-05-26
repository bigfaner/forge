---
created: "2026-05-26"
author: "fanhuifeng"
status: Draft
---

# Proposal: Run-Tests Journey Isolation & Pipeline Cleanup

## Problem

测试 pipeline 存在一个阻断性 bug 和多个逻辑缺陷，导致 test.run 任务无法正确执行，且 pipeline 整体存在冗余。

### Evidence

1. **P0 阻断性 bug**：`forge-cli/pkg/prompt/data/test-run.md` 引用 `Skill(skill="forge:run-e2e-tests")`，但代码库中不存在该 skill。实际 skill 名为 `forge:run-tests`。T-test-run 任务执行时必然失败。
2. **P1 范围泄露**：test.run 任务模板声称 "Execute staged e2e test scripts for the {{FEATURE_SLUG}} feature"，但 run-tests skill 执行 `just test`（无 journey 参数），实际运行全部 e2e 测试，而非当前 feature 的 journey。
3. **P2 冗余任务**：verify-regression 紧跟 test.run 之后，但 test.run 已跑全部测试。`forge test promote` 不在 pipeline 内自动执行，verify-regression 无增量价值。
4. **P3 Quick 模式过重**：gen-contracts（30-45min）对 Quick 模式的小 feature 来说投入产出比低，可以省略。

### Urgency

P0 导致 test.run 任务必败，pipeline 形同虚设。任何开启 `auto.test.quick=true` 或 `auto.test.full=true` 的项目都会触发此问题。

## Proposed Solution

分四个层次修复，按优先级递进：

1. **修复 skill 引用**：prompt 模板中 `forge:run-e2e-tests` → `forge:run-tests`
2. **Journey 级隔离**：run-tests skill 通过 `ls docs/features/<slug>/testing/` 发现 journey 列表，对每个 journey 执行 `just test <journey>`（而非 `just test`）
3. **移除 verify-regression**：从自动 pipeline 中移除，改为 `forge test promote` 后的手动步骤
4. **Quick 模式跳过 gen-contracts**：pipeline 从 gen-journeys → gen-contracts → gen-scripts 简化为 gen-journeys → gen-scripts

### Innovation Highlights

这不是创新，是修 bug 和消除技术债。唯一的设计决策是 journey 发现机制：选择了 `ls` 目录扫描而非新增 CLI 命令，因为 journey 就是 `docs/features/<slug>/testing/` 下的子目录，不需要额外抽象。

## Requirements Analysis

### Key Scenarios

- **Happy path**：feature 有 2 个 journey（cli + api），run-tests 分别执行 `just test <journey1>` 和 `just test <journey2>`，合并结果生成报告
- **Edge case：无 journey**：`docs/features/<slug>/testing/` 不存在或为空，run-tests 报错提示先运行 gen-journeys
- **Edge case：单 journey**：只有一个 journey，等价于 `just test <journey>`
- **Surface 编排**：web/api/mobile 的 dev → probe 生命周期包裹所有 journey 测试（启动一次服务，依次跑每个 journey，最后 teardown）
- **Quick 模式简化**：gen-journeys 完成后直接生成测试脚本，跳过 contract 形式化

### Non-Functional Requirements

- Journey 发现必须零依赖（纯 `ls`），不需要 forge CLI 或额外工具
- Surface 编排序列（dev/probe/test/teardown）的行为不变，仅 test 步骤从单次 `just test` 改为循环 `just test <journey>`

### Constraints & Dependencies

- justfile recipe 已支持 `just test <journey>` 参数（Go: `-run` 正则，Node/Playwright: 路径过滤），无需修改 justfile
- Surface 规则文件的编排序列表不变
- Quick 模式去掉 gen-contracts 需要 `GetQuickTestTasks()` 和 `GetBreakdownTestTasks()` 中的依赖链调整

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 无改动成本 | P0 bug 持续阻断 pipeline | Rejected: pipeline 不可用 |
| 仅修 P0+P1 | — | 最小改动范围 | P2/P3 冗余持续存在 | 可接受：分步实施 |
| 新增 `forge test list-journeys` CLI 命令 | 内部 | 结构化输出，可扩展 | 需改 Go 代码，ls 已足够 | Rejected: 过度工程 |
| **P0+P1+P2+P3 全部修复** | — | 一次性清理所有问题 | 改动范围较大 | **Selected: 审计发现的全部问题应一并解决** |

## Feasibility Assessment

### Technical Feasibility

- P0：单行文本替换（prompt 模板）
- P1：修改 run-tests SKILL.md（skill 文档），无需改 Go 代码
- P2：修改 `GetBreakdownTestTasks()` 和 `GetQuickTestTasks()`（Go 代码），移除 verify-regression 任务定义和依赖
- P3：修改 `GetQuickTestTasks()`（Go 代码），移除 gen-contracts 任务，调整依赖链从 gen-journeys 直接连到 gen-scripts

### Resource & Timeline

- P0+P1：1-2h（skill 文档 + prompt 模板修改）
- P2+P3：2-3h（Go 代码 + 测试更新）
- 总计：3-5h

### Dependency Readiness

- justfile recipe 已支持 journey 参数（`surface-aware-justfile` feature 已完成）
- Surface 规则文件已有 "Journey 过滤" 章节（待消费）
- `forge test run-journey` CLI 命令已存在（本次不使用，保留为手动调试工具）

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| verify-regression 在 pipeline 内有价值 | XY Detection | 实际上是 Y=确认不破坏已有测试，但 X=verify-regression 在 promote 之前无增量。Overturned: 移除 |
| Quick 模式需要完整的 journey-contract-script 三层 | Occam's Razor | Quick 模式定位是小 feature，contract 层可省略。Refined: Quick 跳过 contract |
| `forge test run-journey` 适合 pipeline 使用 | Stress Test | 每次调用创建未使用的临时目录，N 个 journey 需要 N 次 CLI 进程。Overturned: skill 直接调用 `just test <journey>` |

## Scope

### In Scope

- 修复 `forge:run-e2e-tests` → `forge:run-tests` prompt 模板引用
- run-tests skill 增加 journey 发现步骤（ls 目录扫描）
- run-tests skill 改为 per-journey 执行 `just test <journey>`
- run-tests skill 增加 surface 编排与 per-journey 执行的组合逻辑（dev/probe 启动一次，循环跑 journey，teardown 一次）
- 从 pipeline 移除 `test.verify-regression` 任务定义和依赖链
- Quick 模式移除 `test.gen-contracts` 任务，调整依赖链
- 更新相关测试

### Out of Scope

- `forge test run-journey` CLI 命令的重构（保留为手动调试工具）
- 新增 `forge test list-journeys` CLI 命令
- Surface 规则文件的 "Journey 过滤" 标签消费（后续可做）
- `gen-contracts` 任务模板和 prompt 的清理（标记为 Quick 模式不使用即可）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| ls 目录扫描在非标准目录结构下失败 | L | M | 前置检查：`docs/features/<slug>/testing/` 不存在时报明确错误 |
| Quick 模式去掉 gen-contracts 后测试脚本质量下降 | M | M | gen-scripts skill 可从 journey 步骤描述直接提取测试维度，contract 知识内嵌到 gen-scripts 规则中 |
| 移除 verify-regression 后缺少回归保护 | L | L | promote 后用户手动跑 `just test-e2e`；all-completed hook 的 quality-gate 仍提供项目级安全网 |
| per-journey 执行增加了总耗时（N 次 just 调用 vs 1 次） | M | L | journey 数量通常 1-3 个，每次 just test <journey> 只跑子集，总耗时持平或更短 |

## Success Criteria

- [ ] T-test-run 任务执行时正确调用 `forge:run-tests` skill（不再报 skill 不存在）
- [ ] run-tests skill 只执行当前 feature 的 journey 测试，不运行其他 feature 的测试
- [ ] Surface 编排中 dev/probe 执行一次，per-journey 执行 test，teardown 执行一次
- [ ] `GetBreakdownTestTasks()` 不再生成 `T-test-verify-regression` 任务
- [ ] `GetQuickTestTasks()` 不再生成 `T-test-gen-contracts` 任务
- [ ] `GetQuickTestTasks()` 的依赖链为 gen-journeys → gen-scripts（跳过 gen-contracts）
- [ ] 现有测试 `TestGetQuickTestTasks_*` 和 `TestGetBreakdownTestTasks_*` 全部通过

## Next Steps

- Proceed to `/quick` to generate and execute tasks
