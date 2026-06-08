---
created: "2026-06-08"
author: fanhuifeng
status: Draft
intent: "enhancement"
---

# Proposal: Behavioral Test Accuracy

## Problem

Forge 的测试管线生成的测试是**结构性的**（验证 API 不崩溃、数据格式正确），而非**行为性的**（验证功能真正可用）。当 seed data 为空容器、测试只覆盖孤立 CRUD 操作、断言仅检查 HTTP 200 时，全部测试通过但核心功能可能完全不可用。

### Evidence

pm-work-tracker 项目的里程碑地图功能：完整管线（gen-contracts → gen-test-scripts → run-test-backend → run-test-frontend）全部通过，但手动 QA 发现**所有里程碑地图内没有任何里程碑**。测试验证了空容器上的 CRUD 操作，从未验证核心用户工作流（在地图中创建里程碑）。

根因分析三层：
- **L1**: 测试仅验证单步 CRUD（create/edit/delete map），不验证多步工作流（创建地图 → 添加里程碑 → 转换状态）
- **L2**: Seed data 创建空容器，断言仅检查 HTTP 200 + response shape
- **L3**: Journey eval 787/1000、Contract eval 665/1000，均低于目标但被绕过

### Urgency

这是管线架构的系统性缺陷，不是偶发问题。任何涉及父子实体关系的 feature 都会遇到同样问题。pm-work-tracker 项目已因此产生了虚假的"全部通过"结果，导致功能缺陷未被及时发现。

## Proposed Solution

从源头到终端全链路增强测试的行为性质量，覆盖三个管线阶段：

1. **gen-journeys 强制 Golden Path**：每个 feature 必须至少包含一个跨越多步操作的 Golden Path Journey，验证完整用户工作流而非孤立 CRUD
2. **gen-contracts 新增 Fixture Specification 维度**：每个 Contract 的 Preconditions 必须声明前置数据状态（需要哪些实体、实体间关系、最小数据量）
3. **gen-test-scripts 断言深度 + seed data 丰富度规则**：从 Contract 读取 fixture spec，生成丰富 fixture；断言必须验证业务结果而非仅 HTTP 状态码

两端 eval rubric 同步新增评估维度，确保质量可度量。

### Innovation Highlights

这不是标准行业实践的直接采用。行业中的 E2E 测试通常由人类 QA 工程师手动编写，他们天然理解"golden path"和"丰富数据"的重要性。Forge 的挑战在于**让 AI 生成管线自行推断这些需求**——从 PRD/Design 文档中提取工作流语义，从领域模型中推断实体关系，从 Contract 中传递 fixture 需求。

关键创新点：**Contract 级别的声明式 Fixture Specification**。不是在测试脚本生成阶段猜测需要什么数据，而是在 Contract 阶段就明确声明，让测试脚本生成时直接消费。这消除了"推断错误"的可能性，也使 fixture 需求可审计、可评估。

## Requirements Analysis

### Key Scenarios

**Happy path**: 用户执行 `gen-journeys` → 管线自动识别 feature 的核心工作流 → 生成包含 Golden Path 的 Journey → Contract 包含完整 Fixture Specification → 测试脚本创建丰富 fixture 并验证业务结果 → 测试能发现真实 bug

**Edge cases**:
- 简单 feature（无父子实体关系）：Golden Path Journey 仍然适用，但 fixture specification 可以声明为最小数据集
- 单实体 CRUD feature：Golden Path 可能就是完整的 CRUD 循环，但仍然跨越多个步骤

**Error scenarios**:
- PRD 中缺少工作流描述：gen-journeys 应从用户故事中推断，或要求补充
- 实体关系无法从代码推断：Fixture Specification 标记为 "unverifiable"，测试生成时使用保守的最小数据集

### Constraints & Dependencies

- 依赖 gen-journeys、gen-contracts、gen-test-scripts 三个现有 skill 的修改
- 依赖 eval rubric 的更新
- 不依赖外部系统或 API
- 不改变管线阶段顺序或新增阶段

## Alternatives & Industry Benchmarking

### Industry Solutions

- **Playwright/Cypress 最佳实践**: 强调 "test user workflows, not implementation details" 和 "use realistic test data"。这是手动测试的行业标准，但 Forge 需要将这些原则编码为生成规则。
- **TestNG/JUnit fixture patterns**: 使用 @BeforeMethod / beforeEach 声明 fixture 需求。Forge 的 Contract-level Fixture Specification 是将此概念提升到 spec 层面。
- **Contract Testing (Pact)**: 消费者驱动的契约测试关注 API 兼容性，但不关注功能完整性。Forge 的 challenge 是功能完整性。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 无成本 | 测试继续产生虚假通过；pm-work-tracker 问题重复发生 | Rejected: 已有真实项目失败证据 |
| 仅增强 eval rubric | eval-test-cases 提案 | 改动小 | 不解决 L2（seed data）根本问题 | Rejected: 不充分 |
| 轻量级规则增强 | — | 仅改 gen-test-scripts | 上游仍是 CRUD 描述，推断可能不准确 | Rejected: 上下游信息链断裂 |
| **全链路行为性增强** | Playwright 最佳实践 + Contract Testing 概念 | 从源头保证行为性描述；fixture 需求可审计 | 改动涉及 3 skill + 2 rubric | **Selected: 唯一能系统性解决三层根因的方案** |

## Feasibility Assessment

### Technical Feasibility

完全可行。所有改动都在现有 skill 的规则文件、模板文件和 rubric 文件中，不涉及架构变更。

### Resource & Timeline

修改量：3 个 skill 的规则/模板 + 2 个 eval rubric。每个 skill 的修改是独立的，可以并行。

### Dependency Readiness

无外部依赖。所有需要修改的文件都在 Forge plugin 内部。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 测试全部通过 = 功能可用 | 5 Whys | Overturned: pm-work-tracker 证明全部通过可能只是验证了空容器上的 CRUD |
| eval gate 阻断就能保证质量 | XY Detection | Refined: eval gate 是必要条件但不是充分条件；即使 eval 通过，如果 Journey/Contract 本身只描述 CRUD，下游测试仍然无效 |
| gen-test-scripts 可以自行推断 fixture 需求 | Assumption Flip | Overturned: 测试脚本生成阶段缺少领域上下文来推断"需要什么样的数据"，必须从上游 Contract 传递 |

## Scope

### In Scope

1. gen-journeys SKILL.md 和规则：新增 Golden Path Journey 强制要求
2. gen-contracts 模板和规则：新增 Fixture Specification 维度到 Contract Preconditions
3. gen-test-scripts 规则：断言深度规则 + seed data 丰富度规则（从 Contract 读取 fixture spec）
4. eval rubrics (journey.md, contract.md)：新增对应评估维度

### Out of Scope

- Pipeline 可靠性 / eval gate 行为（已由 eval-diagnostic-mode 覆盖）
- run-tests 执行阶段改动
- Mobile 测试特殊处理
- 共享 fixture 库 / 集中式 fixture 管理
- 新增 skill 或 pipeline 阶段

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Golden Path 强制要求导致简单 feature 的 Journey 过度膨胀 | M | L | Golden Path 可以是简洁的多步 CRUD 循环；规则应区分"简单 feature"和"复杂 feature"的期望 |
| Contract Fixture Specification 增加 Contract 复杂度，降低 eval 得分 | M | M | Fixture Specification 作为 Preconditions 的子维度，增加权重但独立评分；新 rubric 维度应设合理最低阈值 |
| 断言深度规则过于严格，导致某些合法测试被误判 | L | M | 80% 阈值允许 20% 的"结构性"断言存在，覆盖 logging、health check 等合理场景 |

## Success Criteria

consistency_check_result:
  status: pass
  pairs_checked: 6
  conflicts_found: 0

- [ ] SC-1: 每个包含父子实体关系的 feature，gen-journeys 至少生成 1 个 Golden Path Journey（跨越 3+ 步骤的完整工作流）
- [ ] SC-2: 每个 Contract 的 Preconditions 包含 Fixture Specification（声明前置实体类型、关系类型、最小数量）
- [ ] SC-3: gen-test-scripts 生成的测试中，≥80% 的断言验证业务结果（实体存在、状态正确、关系完整），而非仅 HTTP 状态码
- [ ] SC-4: 当 Fixture Specification 声明需要 N 个子实体时，生成的测试 fixture 必须创建 ≥N 个子实体
- [ ] SC-5: Journey eval rubric 新增 "Workflow Coverage" 维度（150 分），评分标准包含 Golden Path 存在性和多步覆盖度
- [ ] SC-6: Contract eval rubric 新增 "Fixture Specification" 维度（100 分），评分标准包含前置数据声明完整性和实体关系覆盖度

## Next Steps

- Proceed to `/write-prd` to formalize requirements
