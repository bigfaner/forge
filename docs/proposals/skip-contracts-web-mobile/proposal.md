---
created: "2026-06-09"
author: fanhuifeng
status: Draft
intent: "enhancement"
---

# Proposal: Skip gen-contracts for Interaction-Only Features

## Problem

Forge 测试流水线（gen-journeys → gen-contracts → gen-test-scripts）对纯 Web/Mobile 交互场景存在结构性断层：gen-contracts 只能为协议级 surface（API/CLI/TUI）生成有意义的契约，纯交互级 journey 无法产出契约文件，导致 gen-test-scripts 无输入可用，最终静默跳过这些 journey 的测试生成。

### Evidence

milestone-map feature 定义了 7 个旅程，其中 5 个配置为 `surface_types: ["web"]`（纯 Web）。经过完整流水线后，只有 2 个双 surface（web+api）旅程获得了测试脚本，5 个纯 Web 旅程（71%）零测试覆盖——且流水线未报告任何异常。

### Urgency

每新增一个 Web-only feature 都会重复此问题。测试覆盖率缺口对质量门禁不可见，生产代码可能上线而没有自动化测试保护。延迟修复的成本随 Web feature 数量线性增长。

## Proposed Solution

两层修复，决策优先级明确如下：**Skill 层（gen-test-scripts）是权威决策者**，负责判断 journey 的 surface 类型并路由到正确的生成路径；Pipeline 层的跳过仅为效率优化（避免执行注定失败的 gen-contracts）。即使 Pipeline 层判断出错（如 surface-type 字段缺失），Skill 层仍能正确兜底——不会因缺少 contract 而静默跳过，而是走直达路径或报告失败。

**1. 流水线层（forge-cli PipelineRegistry）**：效率优化层。检查 feature 所有业务任务的 `surface-type` 字段，若不存在 tui/cli/api 类型任务，跳过 `T-test-gen-contracts` 和 `T-eval-contract`，调整依赖链使 gen-scripts 直接依赖 gen-journeys。

**2. Skill 层（gen-test-scripts）**：权威决策层。为混合 surface 场景增加直达路径——当 journey 的 `surface_types` 仅含 web/mobile 且无对应 contract 文件时，直接从 journey.md + types/web.md 生成测试脚本，跳过 contract 前置检查。生成完成后执行覆盖率完整性自检（按 surface type 分别验证），缺口硬失败。

### Developer-Facing Observability

- **任务列表**：纯 Web/Mobile feature 不再出现 `T-test-gen-contracts`/`T-eval-contract`，`T-test-gen-scripts` 直接依赖 `T-gen-journeys`
- **跳过日志**：PipelineRegistry 输出 INFO 日志：`[skip] T-test-gen-contracts: no protocol-level surface tasks found (all tasks: web)`
- **直达路径标识**：`Generating test scripts for <journey> via direct path (surface: web, no contract required)`
- **覆盖率报告**：`Coverage: web 3/3 journeys, api 2/2 journeys`；缺口时 FAIL 并列缺失列表
- **失败回退**：直达路径生成失败时输出 `Direct path generation failed for <journey>: <reason>`，而非静默跳过

### Innovation Highlights

不是重新设计 Contract 格式来适配 Web，而是承认执行模型的本质差异：

- **协议级（API/CLI）**：输入输出是结构化数据，Contract 是自然抽象，保留 `journey → contract → test-script` 路径
- **交互级（Web/Mobile）**：输入是用户动作、输出是视觉状态，走 `journey → test-script` 直达路径

> **经验验证承诺**："journey.md 包含生成所需全部信息"是合理推断但未经实证。实施前需对 3 个现有 web-only journey 验证映射完整性（步骤描述 → step-action/fixture_spec/Outcome）。失败则回退到补充结构化字段方案。

## Requirements Analysis

### Key Scenarios

- **S1: 纯 Web feature**：所有业务任务 `surface-type: web` → 流水线跳过 gen-contracts，gen-scripts 走直达路径
- **S2: 纯 Mobile feature**：同 S1
- **S3: 混合 surface feature**（如 `backend=api, frontend=web`）：部分任务 web、部分 api → gen-contracts 正常生成（api journey 需要），web surface 的 gen-scripts 走直达路径
- **S4: 多 surface 项目只改前端**：业务任务全部 `surface-type: web` → 跳过 gen-contracts

### Constraints & Dependencies

- 业务任务必须有 `surface-type` 前置字段（breakdown-tasks / quick-tasks 已支持）
- gen-test-scripts 的 types/web.md 和 types/mobile.md 需要补充直达生成规则

### Non-Functional Requirements

- **可观测性**：Pipeline 跳过和直达路径决策通过日志和输出可追溯（详见 Developer-Facing Observability）
- **性能**：surface-type 条件检查 O(n) 遍历业务任务列表，不引入显著开销
- **向后兼容性**：已有 API/CLI/TUI 流水线不变（SC-7）。直达路径是条件分支，缺失时保守回退
- **可回退性**：feature flag `forge.skip_contracts.enabled`（默认 true）控制直达路径，详见 Rollback Mechanism

## Alternatives & Industry Benchmarking

### Industry Solutions

- **Playwright Codegen**：Browser Context Recording 捕获用户交互，直接生成 `page.click()`/`page.fill()` 等脚本，不经过契约层
- **Maestro**：声明式 YAML（`flows/` 目录）定义测试，`- tapOn: "Login"` 即测试步骤，核心理念 "flow as spec"
- **Cypress Studio**：从交互录制生成 `cy.get().click()` 代码，无接口契约中间产物

共同模式：**Action-Driven Testing**——用户动作直接映射测试步骤。Forge 的 gen-contracts 为协议级设计，对交互级属过度设计。

### Comparison Table

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 零改动 | 71% 覆盖缺口不可见 | Rejected |
| 泛化 Contract 模型 | 统一流水线 | Web Contract 退化为 DOM 操作+视觉断言列表，与测试脚本 1:1 重合，维护等价文档无抽象价值 | Rejected: 仪式而非实质 |
| 仅加覆盖率审计 | 最小改动 | 发现缺口但无法修复，开发者仍需手动写测试 | Rejected: 只诊断不治疗 |
| 仅流水线跳过 | 减少无效执行 | 混合场景 web-only journey 仍无测试 | Rejected: 不完整 |
| **流水线跳过 + Skill 直达路径** | 完整覆盖全场景 | 双路径维护 | **Selected** |

## Feasibility Assessment

- PipelineRegistry 已有 `UISurfaceOnly` 先例可复用
- 业务任务 `surface-type` 字段已存在
- types/web.md 已有 Web 规则，需扩展直达映射模板

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| Contract 是所有 surface 类型的必要中间层 | Assumption Flip | Overturned: 交互级 surface 输入是用户动作而非结构化数据 |
| gen-contracts 应始终包含在测试流水线中 | XY Detection | Refined: 仅 protocol-level surface 需要 |
| 流水线跳过应基于项目级 surface 配置 | 5 Whys | Overturned: 应基于业务任务的 surface-type |

## Scope

### In Scope

- PipelineRegistry 新增 `CondHasProtocolSurfaceTask` 条件：检查 surface-type 是否存在 tui/cli/api。缺失、为空、或未知值（如 "desktop"）采用保守策略视为"可能有协议级"，不触发跳过，并输出 WARN 日志
- 依赖链调整：gen-contracts 跳过时 gen-scripts 依赖 gen-journeys；混合 feature 保持原链路
- gen-test-scripts SKILL.md 路由修改：Step 2 检查 `surface_types`，仅含 web/mobile 时跳过 contract。直达映射：步骤段落 → `step-action`（click/type/navigate），前置条件描述 → `fixture_spec`，预期结果描述 → `Outcome`（assertVisible/assertText）。types/web.md 定义映射模板
- 覆盖率自检：`count(journeys_of_type) == count(test-scripts_of_type)`。Surface → Test Type 映射定义在 SKILL.md 中：`web→Web E2E Test`、`mobile→Mobile E2E Test`、`api→API Functional Test`、`cli→CLI Functional Test`、`tui→Terminal Functional Test`。缺口或类型不匹配硬失败
- types/web.md / types/mobile.md：补充直达生成规则

### Out of Scope

- gen-contracts / gen-journeys skill 修改
- Contract 六维度格式修改
- run-tests 修改
- TUI surface 处理方式变更

## Rollback Mechanism

Feature flag `forge.skip_contracts.enabled`（默认 true）控制直达路径：

1. **即时回退**：设 false 强制走 contract 路径
2. **Pipeline 回退**：移除 `CondHasProtocolSurfaceTask` 条件
3. **数据安全**：新增文件不覆盖现有文件
4. **监控**：跟踪直达路径脚本通过率，显著低于 contract 路径时触发评估

## Key Risks

| Risk | L | I | Mitigation |
|------|---|---|------------|
| surface-type 未填充/填写错误 | H | H | 双重兜底：(1) 缺失/空/未知值保守不跳过+WARN 日志；(2) Skill 层覆盖率+类型匹配自检。评 H：字段新引入，历史 feature 可能未标注 |
| journey.md 结构化信息不足（核心假设） | M | H | 实施前验证 3 个 web-only journey 映射完整性；SC-8 验证断言质量；失败则补充结构化字段 |
| 直达路径测试质量低于 Contract 驱动 | M | M | types/web.md 补充充分生成规则 |
| 双路径维护成本 | L | M | 分支逻辑按执行模型分流，仅 web/mobile 用直达路径 |

## Success Criteria

- [ ] SC-1: 纯 Web feature 流水线不生成 T-test-gen-contracts 和 T-eval-contract
- [ ] SC-2: 纯 Mobile feature 同 SC-1
- [ ] SC-3: 混合 surface feature 正常生成 T-test-gen-contracts 和 T-eval-contract
- [ ] SC-4: 同一项目中，前端-only feature 跳过 gen-contracts，后端 feature 保留
- [ ] SC-5: 直达路径生成的脚本包含与 journey 步骤对应的用户动作调用和至少一个可视化断言（非空、非骨架）
- [ ] SC-6: 按 surface type 自检覆盖率，缺口或类型不匹配时 FAIL 并输出缺口列表
- [ ] SC-7: 已有 API/CLI/TUI 流水线行为不变
- [ ] SC-8: types/web.md 和 types/mobile.md 直达规则产出包含有意义断言的测试脚本

consistency_check_result:
  status: pass
  pairs_checked: 21
  conflicts_found: 0

## Next Steps

- Proceed to `/write-prd` to formalize requirements
