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

两层修复：

**1. 流水线层（forge-cli PipelineRegistry）**：检查 feature 所有业务任务的 `surface-type` 字段，若不存在 tui/cli/api 类型任务，跳过 `T-test-gen-contracts` 和 `T-eval-contract`，调整依赖链使 gen-scripts 直接依赖 gen-journeys。

**2. Skill 层（gen-test-scripts）**：为混合 surface 场景增加直达路径——当 journey 的 `surface_types` 仅含 web/mobile 且无对应 contract 文件时，直接从 journey.md + types/web.md 生成测试脚本，跳过 contract 前置检查。生成完成后执行覆盖率完整性自检，缺口硬失败。

### Innovation Highlights

不是重新设计 Contract 格式来适配 Web（这等于在中间层重复写测试脚本），而是承认执行模型的本质差异：

- **协议级（API/CLI）**：输入输出是结构化数据，Contract 是自然抽象，保留 `journey → contract → test-script` 路径
- **交互级（Web/Mobile）**：输入是用户动作、输出是视觉状态，journey.md 已包含生成所需全部信息，中间层无信息增益，走 `journey → test-script` 直达路径

## Requirements Analysis

### Key Scenarios

- **S1: 纯 Web feature**：所有业务任务 `surface-type: web` → 流水线跳过 gen-contracts，gen-scripts 走直达路径
- **S2: 纯 Mobile feature**：同 S1
- **S3: 混合 surface feature**（如 `backend=api, frontend=web`）：部分任务 web、部分 api → gen-contracts 正常生成（api journey 需要），web surface 的 gen-scripts 任务对无 contract 的 web-only journey 走直达路径
- **S4: 多 surface 项目只改前端**：业务任务全部 `surface-type: web` → 跳过 gen-contracts，适配"同一项目不同 feature 涉及不同 surface"的场景

### Constraints & Dependencies

- 业务任务必须有 `surface-type` 前置字段（breakdown-tasks / quick-tasks 已支持）
- gen-test-scripts 的 types/web.md 和 types/mobile.md 需要补充直达生成规则

## Alternatives & Industry Benchmarking

### Industry Solutions

多数测试框架（Playwright、Cypress、Maestro）直接从用户场景描述生成测试脚本，不经过中间契约层。Forge 的 gen-contracts 是为 API 测试设计的特殊抽象。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动 | Web journey 零覆盖，缺口不可见 | Rejected: 71% 覆盖缺口不可接受 |
| 泛化 Contract 模型 | — | 统一流水线 | Web Contract = 重写测试脚本，中间层无信息增益 | Rejected: 仪式而非实质 |
| 仅加覆盖率审计 | — | 最小改动 | 只诊断不治疗 | Rejected: 不修复生成能力 |
| 仅流水线跳过 | — | 减少 gen-contracts 执行 | 混合场景下 web-only journey 仍无测试 | Rejected: 不完整 |
| **流水线跳过 + Skill 直达路径** | — | 完整覆盖两种场景 | 双路径维护 | **Selected: 唯一覆盖全场景的方案** |

## Feasibility Assessment

### Technical Feasibility

- PipelineRegistry 已有 surface 类型过滤先例（`UISurfaceOnly` for T-validate-ux）
- 业务任务 `surface-type` 字段已存在且由 skill 填充
- gen-test-scripts 的 `types/web.md` 已有 Web 特定规则，只需扩展

### Dependency Readiness

- 业务任务 surface-type 字段：已就绪
- PipelineRegistry 扩展机制：已就绪（可复用 `UISurfaceOnly` 模式）

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| Contract 是所有 surface 类型的必要中间层 | Assumption Flip | Overturned: 交互级 surface 输入是用户动作而非结构化数据，Contract 无信息增益 |
| gen-contracts 应始终包含在测试流水线中 | XY Detection | Refined: gen-contracts 是 protocol-level surface 的必要步骤，非交互级的必要步骤 |
| 流水线跳过应基于项目级 surface 配置 | 5 Whys | Overturned: 应基于业务任务的 surface-type，同一项目不同 feature 可能涉及不同 surface |

## Scope

### In Scope

- PipelineRegistry 新增 `CondHasProtocolSurfaceTask` 条件：检查业务任务 surface-type 是否存在 tui/cli/api，不存在则跳过 `T-test-gen-contracts` 和 `T-eval-contract`
- 依赖链调整：gen-contracts 跳过时，gen-scripts 依赖 gen-journeys / eval-journey
- gen-test-scripts SKILL.md：前置条件路由（protocol-level 需 Contract，interaction-level 跳过）
- gen-test-scripts 覆盖率自检：生成完成后 `count(journeys) == count(test-script-sets)`，缺口硬失败
- types/web.md / types/mobile.md：补充 journey.md 直达生成规则

### Out of Scope

- gen-contracts skill 修改
- Contract 六维度格式修改
- gen-journeys skill 修改
- run-tests 修改
- TUI surface 处理方式变更（保持 Contract 路径）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 业务任务 surface-type 未填充或填写错误 | M | H | 覆盖率自检硬失败兜底：即使流水线判断错误，自检会发现缺失的测试脚本 |
| 直达路径生成的测试脚本质量低于 Contract 驱动 | M | M | types/web.md 补充充分的生成规则；journey.md 本身已包含足够信息 |
| 双路径长期维护成本 | L | M | 两种路径的分支逻辑清晰（按执行模型分流），且仅 web/mobile 使用直达路径 |

## Success Criteria

- [ ] SC-1: 纯 Web feature（所有业务任务 surface-type=web）流水线不生成 T-test-gen-contracts 和 T-eval-contract 任务
- [ ] SC-2: 纯 Mobile feature（所有业务任务 surface-type=mobile）流水线不生成 T-test-gen-contracts 和 T-eval-contract 任务
- [ ] SC-3: 混合 surface feature（存在 surface-type 为 api/cli/tui 的业务任务）流水线正常生成 T-test-gen-contracts 和 T-eval-contract
- [ ] SC-4: 同一多 surface 项目中，只改前端的 feature 跳过 gen-contracts，同时改后端的 feature 保留 gen-contracts
- [ ] SC-5: gen-test-scripts 对无 contract 的 web/mobile journey 从 journey.md 直接生成测试脚本（不报错、不跳过）
- [ ] SC-6: gen-test-scripts 完成后自检覆盖率，journey 无对应测试脚本时任务 FAIL 并输出缺口列表
- [ ] SC-7: 已有 API/CLI/TUI 流水线行为不变（回归验证通过）

consistency_check_result:
  status: pass
  pairs_checked: 21
  conflicts_found: 0

## Next Steps

- Proceed to `/write-prd` to formalize requirements
