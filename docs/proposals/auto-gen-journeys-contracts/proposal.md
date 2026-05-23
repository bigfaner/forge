---
created: 2026-05-23
author: "faner"
status: Draft
---

# Proposal: 将 gen-journeys/gen-contracts 升级为自动生成的测试流水线任务

## Problem

当前测试流水线存在自动化断口：`eval-journey` 和 `eval-contract` 已经是 `autogen.go` 中的自动生成任务，但它们依赖的 `gen-journeys` 和 `gen-contracts` 仍然需要手动调用。这意味着流水线无法从业务任务结束后自动流转到测试生成阶段。

### Evidence

- `autogen.go` 的 `GetBreakdownTestTasks()` 生成 eval-journey 和 eval-contract，但不生成 gen-journeys 和 gen-contracts
- Breakdown 模式下，用户需要手动执行 `/gen-journeys` 和 `/gen-contracts`，然后才能触发自动生成的 eval 任务
- Quick 模式使用 `gen-and-run` 组合任务跳过了 Journey-Contract 流程，直接从代码生成测试，导致测试质量无保证

### Urgency

v3.0.0 正在引入 test profile 系统（解耦 Playwright），测试流水线的完整性是 v3.0.0 发布的前提条件。当前的断口意味着流水线无法端到端自动化运行。

## Proposed Solution

将 `gen-journeys` 和 `gen-contracts` 从手动调用的 skills 升级为 `autogen.go` 中的自动生成测试任务，建立完整的自动化流水线。

**Breakdown 模式**：
```
gen-journeys → eval-journey → gen-contracts → eval-contract → gen-scripts → run → verify
```

**Quick 模式**（跳过 eval 质量关卡）：
```
gen-journeys → gen-contracts → gen-scripts → run → verify
```

Quick 模式的 `gen-and-run` 组合任务被替换为拆分的独立任务，与 Breakdown 模式共享相同的 gen-contracts/gen-scripts/run 任务定义。

**Quick 模式依赖拓扑**：采用 staged across types 策略（而非 sequential per-type）。即所有 profile type 的 gen-journeys 并行执行 → 汇聚到 gen-contracts → 所有 gen-scripts 并行 → run → verify。理由：(1) gen-journeys 各 type 间无依赖；(2) gen-contracts 需要所有 Journey 完成后才能执行代码侦察；(3) 与 Breakdown 模式的拓扑保持一致（仅少了 eval 中间关卡）。

**依赖解析重写**：当前 `resolveBreakdownDeps` 的 E2E 分支使用硬编码索引（`evalJourneyIdx := 0`, `evalContractIdx := 1`, `genStart := 2`）。在 eval-journey 前插入 gen-journeys 会导致所有索引偏移失效。需将整个 `resolveBreakdownDeps` 重写为基于 `findTaskIndexByPrefix` 的 ID 查找。同理，`resolveQuickDeps` 需完全重写以适配新的多任务拓扑。

### Innovation Highlights

无特殊创新。这是对已有自动化框架的自然延伸——将两个已成熟的 skill 纳入自动化任务生成系统，补全流水线断口。

## Requirements Analysis

### Key Scenarios

1. **Breakdown 模式完整流水线**：业务任务完成后，自动生成 gen-journeys 任务，依次执行完整链路（含 eval 质量关卡）
2. **Quick 模式精简流水线**：业务任务完成后，自动生成 gen-journeys 任务，跳过 eval 直接走 gen-contracts → gen-scripts → run
3. **Quick 模式 proposal 输入**：gen-journeys 需要支持 proposal.md 作为替代输入源（Quick 模式没有 PRD user stories）
4. **向后兼容**：已有 feature 的 index.json 不受影响
5. **零 Journey 中止**：gen-journeys 在 Quick 模式下若产出零 Journey 文件（proposal.md 信息不足），任务应 abort 并给出诊断信息，而非继续执行空链路
6. **索引匹配失败**：resolveBreakdownDeps 重写后若 findTaskIndex 返回 -1（任务未找到），应 panic 并输出明确错误信息，而非静默跳过

### Non-Functional Requirements

- gen-journeys 作为自动任务时，MainSession=false（不涉及子代理编排）
- gen-contracts 作为自动任务时，MainSession=false
- 任务模板使用 embed.FS，与现有模板机制一致
- **执行时间影响**：gen-journeys 通常需 10-30 分钟，作为自动任务会阻塞流水线。应在任务模板的 estimatedTime 中准确标注
- **Token 消耗**：Quick 模式从 gen-and-run（1 个任务）变为 3+ 个任务，总体 Token 消耗约增加 50-80%
- **错误处理**：自动任务执行失败时，run-tasks 已有的 fix-task 机制自动生效

### Constraints & Dependencies

- gen-journeys 依赖 PRD user stories（Breakdown）或 proposal.md（Quick）作为输入
- gen-contracts 依赖 gen-journeys 输出的 Journey 文档
- Quick 模式的 gen-journeys 需要修改 SKILL.md 以支持 proposal.md 输入
- 依赖 test profile 系统（v3.0.0）已合入
- **gen-contracts SKILL.md 前置条件冲突**：gen-contracts 的 Prerequisites 要求 eval-journey 报告（Blocker），但 Quick 模式跳过 eval。解决方案：在 Quick 模式的任务模板中注入 `SKIP_EVAL_GATE=true` 标记，gen-contracts 据此跳过 eval 前置检查
- **gen-journeys SKILL.md 交互审批冲突**：gen-journeys 的 Step 6 HARD-RULE 要求用户审批后才能提交。自动任务模式下无法等待交互。解决方案：在自动任务模板中注入 `AUTO_COMMIT=true` 标记，gen-journeys 据此跳过人工审批直接提交生成的 Journey 文件

**模板注入执行模型**：AUTO_COMMIT 和 SKIP_EVAL_GATE 不是代码级变量解析——它们通过 embed 模板渲染写入任务 .md 文件。任务执行器读取任务 .md 内容作为 prompt，SKILL.md 中的执行步骤根据这些标记调整行为。具体来说，任务模板的 body 中包含条件指令（如 "若 AUTO_COMMIT=true，跳过 Step 6 的用户审批，直接 git add + commit"），SKILL.md 的执行步骤据此在 prompt 层面实现行为切换。

## Alternatives & Industry Benchmarking

### Industry Solutions

测试生成工具通常分为三类：

**录制回放模式**：Cypress Studio、Playwright Codegen、Selenium IDE。用户手动操作 UI，工具录制操作并生成测试脚本。局限性：依赖人工操作，无法从规格文档自动生成。

**AI 辅助生成**：Testim、Mabl、Reflection。利用 AI 分析应用状态并生成/维护测试。部分工具（如 Testim）提供端到端流水线，但测试规格仍需人工定义。

**声明式规格驱动**：Cucumber/SpecFlow（BDD 模式）通过 feature → scenario → step 的分层结构组织测试规格。Graphwalker（基于图的测试路径生成）通过模型定义状态转移和测试路径。

### Industry Pattern Adoption

本提案采用的 stage-per-artifact 分层模式与 BDD 工具链的分层结构对应：

| BDD 概念 | Forge 概念 | 对应关系 |
|----------|-----------|---------|
| Feature（用户行为规格） | Journey（用户工作流叙事） | 从高层行为描述中提取 |
| Scenario（验收条件） | Contract（6 维度可验证约束） | 细化为可验证的规格 |
| Step definition（执行代码） | Test script（具体测试代码） | 从规格自动生成执行代码 |

关键差异：Cucumber 的分层是为人工协作设计的（人写 feature，人写 step），Forge 的分层是为 AI 生成设计的（AI 生成全部三层）。这意味着 Forge 需要额外的 eval 质量关卡来验证 AI 生成的准确性。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零开发成本 | 流水线无法端到端自动化，用户需手动衔接 | Rejected: v3.0.0 需要完整流水线 |
| 合并到 gen-and-run 内部（类似 Testim 单一管道） | Testim | Quick 模式改动最小 | Journey/Contract 质量无法保证，架构耦合 | Rejected: 用户明确要求拆分 |
| **独立自动任务（stage-per-artifact）** | BDD 分层模式 | 架构清晰，两种模式共享任务定义，可独立调试，与行业 BDD 模式对应 | 开发量增加约 30%（依赖解析重写） | **Selected: 最符合 Forge 设计哲学，且与 BDD 行业模式一致** |

## Feasibility Assessment

### Technical Feasibility

可行但需注意复杂度。主要技术挑战是 `resolveBreakdownDeps` 和 `resolveQuickDeps` 的完整重写——前者需要将硬编码索引（L425-449 的 evalJourneyIdx=0, evalContractIdx=1, genStart=2）全部替换为 `findTaskIndexByPrefix` 调用；后者需要从"gen-and-run 并行 → verify"拓扑重写为"staged across types"拓扑。

### Resource & Timeline

预计 6-7 个等效 coding 任务（8 项交付物，部分合并）：

1. types.go + autogenTypeToFile：新增 2 个类型常量、注册表条目、SystemTypes、模板映射（小）
2. autogen.go Breakdown 模式：插入 gen-journeys/gen-contracts 任务定义，重写 resolveBreakdownDeps（大）
3. autogen.go Quick 模式：替换 gen-and-run，重写 resolveQuickDeps，引入共享依赖链逻辑（大）
4. ResolveFirstTestDep + infer.go：更新首任务查找逻辑和类型推断（中）
5. embed 模板文件：test-gen-journeys.md, test-gen-contracts.md，含 AUTO_COMMIT/SKIP_EVAL_GATE 条件指令（中）
6. gen-journeys SKILL.md：支持 proposal.md 输入 + AUTO_COMMIT 非交互模式 + proposal.md 最低信息量检查（中）
7. gen-contracts SKILL.md：SKIP_EVAL_GATE 路径（小）
8. 文档更新：ARCHITECTURE.md、相关文档（小）

### Dependency Readiness

gen-journeys 和 gen-contracts 的 SKILL.md 已稳定，test profile 系统已合入。AUTO_COMMIT/SKIP_EVAL_GATE 的 prompt-level 行为切换不需要新机制——任务模板已有条件渲染能力（`{{MODE}}` 占位符）。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| gen-and-run 类型应保留以兼容旧 feature | Occam's Razor | Refined: gen-and-run 仅用于 Quick 模式，替换后无调用者。可安全删除类型定义，但保留 infer.go 中的识别逻辑以防有历史 index.json 引用 |
| gen-journeys 需要 MainSession=true | 5 Whys | Overturned: gen-journeys 不涉及子代理编排（不像 eval 需要 spawn scorer/reviser），MainSession=false 即可 |
| Quick 模式无法生成高质量 Journeys | Assumption Flip | Refined: proposal.md 的 scope + success criteria 为硬性前提（缺失时 abort），key scenarios 缺失时降级为 smoke-level Journey（仅覆盖 happy path）。代码侦察无法发现尚未实现的场景——这已作为 Quick 模式的已知质量差距被接受，而非依赖下游补偿 |
| resolveBreakdownDeps 已使用 findTaskIndex | 代码审计 | Overturned: 实际代码（autogen.go L425-449）E2E 分支使用硬编码索引（evalJourneyIdx=0, evalContractIdx=1），只有后续 validate/specs/clean-code 才用 findTaskIndex。必须重写 |
| gen-journeys/gen-contracts SKILL.md 无需修改 | 5 Whys | Overturned: gen-journeys 的交互审批 HARD-RULE 和 gen-contracts 的 eval-journey Blocker 前置条件都需要适配自动任务模式。通过任务模板注入条件指令解决，不修改核心生成逻辑 |

## Scope

### In Scope

- 新增 `test.gen-journeys` 和 `test.gen-contracts` 任务类型到 types.go
- 在 autogen.go 中为 Breakdown 和 Quick 模式生成 gen-journeys/gen-contracts 任务
- 重写 `resolveBreakdownDeps` 为基于 findTaskIndexByPrefix 的 ID 查找（消除硬编码索引）
- 重写 `resolveQuickDeps`（适配新的多任务拓扑：staged across types）
- 建立完整的依赖链（gen-journeys → eval-journey → gen-contracts → eval-contract → ...）
- 新增 embed 模板文件（test-gen-journeys.md, test-gen-contracts.md），含 AUTO_COMMIT/SKIP_EVAL_GATE 条件指令
- 更新 infer.go 识别新模式
- 更新 ResolveFirstTestDep：Breakdown 分支查找 T-test-gen-journeys，Quick 分支查找 T-test-gen-journeys（替代 T-quick-gen-and-run）
- 修改 gen-journeys SKILL.md 支持 proposal.md 作为替代输入源 + AUTO_COMMIT 非交互模式 + 最低信息量检查
- 修改 gen-contracts SKILL.md 添加 SKIP_EVAL_GATE 路径（Quick 模式跳过 eval-journey 前置检查）
- Quick 模式用拆分的独立任务替换 gen-and-run 组合任务
- 更新 ARCHITECTURE.md 和相关文档

### Out of Scope

- 修改 gen-journeys/gen-contracts 的核心生成逻辑
- 新增或修改 eval rubrics
- 修改 gen-test-scripts 或 run-tests skill
- 删除 `TypeTestGenAndRun` 类型（保留以防历史引用，标记 deprecated）
- `data/test-gen-and-run.md` 模板文件保留在 embed.FS 中（编译需要），但不再有新任务使用它。在文件头部添加 deprecated 注释说明

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| resolveBreakdownDeps 重写引入索引计算错误 | H | H | 重写为完全基于 findTaskIndex/ID 的依赖解析，消除所有算术索引。单测覆盖所有现有拓扑 + 新拓扑 |
| resolveQuickDeps 重写引入拓扑错误 | H | H | Quick 模式新拓扑与 Breakdown 共享核心逻辑（参数化 includeEval），减少两套独立实现的不一致风险。单测覆盖 staged-across-types 拓扑 |
| proposal.md 信息不足导致 Journey 质量低 | M | M | 定义硬性前提：scope + success criteria 必须存在（缺失则 abort）；key scenarios 缺失时降级为 smoke-level Journey 并标注 quality=low |
| ResolveFirstTestDep Quick 分支前缀不匹配 | H | H | Quick 分支从 findTaskIndexByPrefix(tasks, "T-quick-gen-and-run") 更新为 findTaskIndex(tasks, "T-test-gen-journeys")。单测覆盖 |
| Quick 模式替换 gen-and-run 导致已有测试失败 | L | M | 保留 TypeTestGenAndRun 类型定义和 infer 逻辑，仅停止生成新任务 |

## Success Criteria

- [ ] `forge task index --feature <slug>` 在 Breakdown 模式自动生成 T-test-gen-journeys 和 T-test-gen-contracts 任务
- [ ] Quick 模式自动生成拆分的 gen-journeys → gen-contracts → gen-scripts → run 链路（不再生成 gen-and-run）
- [ ] `forge task validate-index` 通过，包含新的任务类型
- [ ] gen-journeys 支持 proposal.md 作为输入：运行 `forge task claim T-test-gen-journeys` 后，生成的 Journey 文件数量 ≥ 1，内容与 proposal.md 的 scope/success criteria 对应（人工审验）
- [ ] 所有现有测试通过（`go test ./...`）
- [ ] `forge -h` 输出的类型列表包含 test.gen-journeys 和 test.gen-contracts
- [ ] resolveBreakdownDeps 重写后，Breakdown 模式的依赖链正确（eval-journey 依赖 gen-journeys，gen-contracts 依赖 eval-journey）
- [ ] resolveQuickDeps 重写后，Quick 模式生成正确的任务数量（staged across types 拓扑）
- [ ] AUTO_COMMIT=true 注入到 gen-journeys 任务模板时，自动任务模式下跳过人工审批；手动调用 /gen-journeys 时仍要求用户审批
- [ ] SKIP_EVAL_GATE=true 注入到 gen-contracts 任务模板时，Quick 模式下跳过 eval 前置检查
- [ ] infer.go 正确识别 T-test-gen-journeys 和 T-test-gen-contracts 的类型
- [ ] 重写后的 autogen.go 所有单测通过（包含新拓扑和旧回归）
- [ ] 已有 feature 的 index.json 不受影响（向后兼容：新增类型不影响旧类型的解析和执行）
- [ ] gen-journeys 在 proposal.md 缺少 scope 或 success criteria 时 abort 并输出诊断信息

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
