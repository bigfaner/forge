---
created: "2026-06-08"
author: fanhuifeng
status: Draft
intent: "enhancement"
---

# Proposal: Behavioral Test Accuracy

## Problem

Forge 的测试管线生成的测试是**结构性的**（验证 API 不崩溃、数据格式正确），而非**行为性的**（验证功能真正可用）。**行为性测试**的定义：验证用户可观测的业务结果——实体间关系正确、状态转换符合业务规则、核心工作流端到端可达。与结构性测试（验证传输层/格式层正确性）的关键区别在于：行为性测试断言业务语义，结构性测试断言协议合规。当 seed data 为空容器、测试只覆盖孤立 CRUD 操作、断言仅检查 HTTP 200 时，全部测试通过但核心功能可能完全不可用。

### Evidence

pm-work-tracker 项目的里程碑地图功能：完整管线（gen-contracts → gen-test-scripts → run-test-backend → run-test-frontend）全部通过，但手动 QA 发现**所有里程碑地图内没有任何里程碑**。测试验证了空容器上的 CRUD 操作，从未验证核心用户工作流（在地图中创建里程碑）。

根因分析三层：
- **L1**: 测试仅验证单步 CRUD（create/edit/delete map），不验证多步工作流（创建地图 → 添加里程碑 → 转换状态）
- **L2**: Seed data 创建空容器，断言仅检查 HTTP 200 + response shape
- **L3**: Journey eval 787/1000、Contract eval 665/1000，均低于目标但被绕过。注：L3 的管线可靠性改进已由 eval-diagnostic-mode 提案独立覆盖，本提案聚焦 L1 和 L2 的测试有效性改进

### Urgency

这是管线架构的系统性缺陷，不是偶发问题。任何涉及父子实体关系的 feature 都会遇到同样问题。pm-work-tracker 项目已因此产生了虚假的"全部通过"结果，导致功能缺陷未被及时发现。

## Proposed Solution

从源头到终端全链路增强测试的行为性质量，覆盖三个管线阶段：

1. **gen-journeys 强制 Golden Path**：每个 feature 必须至少包含一个 Golden Path Journey。Golden Path 必须同时满足两个约束：(a) 跨越 3+ 步骤操作，(b) 覆盖 PRD/Design 文档中 primary user story 的核心领域动作序列（即语义完整性——不是任意 3 步拼凑，而是用户真实工作流中的关键步骤链）。gen-journeys 规则必须声明"Golden Path 的步骤序列必须从 PRD/Design 的 primary user story 或核心工作流描述中提取，禁止用不相关操作凑数"
2. **gen-contracts 新增 Fixture Specification 维度**：每个 Contract 的 Preconditions 必须声明前置数据状态（需要哪些实体、实体间关系、最小数据量）
3. **gen-test-scripts 断言深度 + seed data 丰富度规则**：从 Contract 读取 fixture spec，生成丰富 fixture；断言必须验证业务结果而非仅 HTTP 状态码。80% 阈值通过 gen-test-scripts 的规则文件在生成时强制执行——规则要求 agent 统计行为性断言占比，若 <80% 则自动补充

#### 断言分类判据

| 类别 | 定义 | 示例 | 计入 ≥80% 分母 |
|------|------|------|----------------|
| **行为性断言** | 验证业务语义：实体存在、状态正确、关系完整、业务规则满足、副作用可见 | `assert milestone.map_id == map.id`、`assert response.data.status == "completed"`、`assert created_count == 3` | 是 |
| **结构性断言** | 验证传输/格式层：HTTP 状态码、响应 schema、字段类型 | `assert response.status == 200`、`assert typeof response.data.id == "string"` | 否（但允许存在） |

边界案例明确：
- `assert response.status == 201 AND response.data.name == "milestone-1"`——混合断言中若包含至少一个业务字段验证，计为行为性
- `assert response.body contains "id"`——仅验证字段存在，为结构性
- `assert list.length > 0`——验证非空集合，为行为性（证明数据创建成功）
- health check / readiness 端点测试——结构性，合理存在，不受 80% 约束

**断言深度质量指标**：除 80% 比例外，gen-test-scripts 规则要求行为性断言中至少 30% 必须是"深度断言"（验证实体间关系或状态转换），而非仅验证字段值匹配。防止通过大量浅层断言（如 `assert name == input`）凑比。

两端 eval rubric 同步新增评估维度，确保质量可度量。

#### Fixture Specification Schema

每个 Contract 的 Preconditions 中必须包含 `fixture_spec` 字段，schema 定义如下：

```yaml
fixture_spec:
  entities:                          # 必需，至少 1 个实体声明
    - entity_type: string            # 必需，实体类型名（如 "Project", "Milestone"）
      min_count: integer             # 必需，最小创建数量（≥1）
      relationship_type: string      # 可选，与父实体的关系（如 "belongs_to", "has_many"）
      parent_entity: string          # 可选，父实体类型名（建立实体间关系时必需）
      field_constraints:             # 可选，特定字段值约束
        - field: string              # 字段名
          value: any                 # 期望值或约束描述
  state_requirements:                # 可选，前置系统状态
    - description: string            # 状态描述
      prerequisite_entity: string    # 依赖的实体类型
```

**最小合法示例**（单实体 CRUD）：
```yaml
fixture_spec:
  entities:
    - entity_type: "Project"
      min_count: 1
```

**完整示例**（父子实体关系）：
```yaml
fixture_spec:
  entities:
    - entity_type: "Map"
      min_count: 1
    - entity_type: "Milestone"
      min_count: 3
      relationship_type: "belongs_to"
      parent_entity: "Map"
      field_constraints:
        - field: "status"
          value: "pending"
```

#### Feature 复杂度分类启发式规则

gen-journeys 规则必须按以下启发式判定 feature 复杂度，并应用差异化期望：

| 判据 | 简单 Feature | 复杂 Feature |
|------|-------------|-------------|
| 实体类型数量 | 1 个实体类型 | ≥2 个实体类型，且存在父子/关联关系 |
| PRD/Design 中的工作流描述 | 单一 CRUD 操作或线性流程 | 多步骤工作流，包含状态转换或实体间交互 |
| Golden Path 期望 | 完整 CRUD 循环（create → read → update → delete），3-5 步 | 覆盖 primary user story 的核心领域动作序列，5+ 步 |
| Fixture Specification 期望 | 声明单一实体 + 最小数量 | 声明实体关系 + 子实体最小数量 + 状态约束 |

判定优先级：实体关系 > 工作流描述。存在父子实体关系即判定为复杂 feature，无论步骤数。

Golden Path 语义完整性的自动化验证代理指标：gen-journeys 规则要求 Golden Path 的每个步骤必须引用 PRD/Design 中明确命名的用户操作（如"创建里程碑"），而非仅引用 HTTP 方法（如"POST /milestones"）。评审者可通过检查步骤描述是否包含领域术语（而非 API 术语）来判断语义完整性。

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
- **Backward compatibility**：已有项目的 Contract 可能不含 `fixture_spec` 字段。gen-test-scripts 规则必须处理缺失情况：当 `fixture_spec` 不存在时，回退到当前的隐式推断模式并输出 warning，不阻断管线。新项目和使用 `gen-contracts` 重新生成的 Contract 将自动包含 `fixture_spec`。

### Non-Functional Requirements

- **可维护性**：Fixture Specification Schema 新增的字段必须通过 gen-contracts 规则文件自动生成，不依赖人工填写
- **兼容性**：三个 skill 的修改通过共享的 `fixture_spec` schema 定义协调，schema 变更时需同步更新三个 skill 的规则文件

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
| 后处理校验（如 test-quality checker） | CI/CD lint 模式 | 不改管线 | 管线已经产出了劣质测试，事后检查为时已晚 | Rejected: 发现问题太迟 |
| **全链路行为性增强** | Playwright 最佳实践 + Contract Testing 概念 | 从源头保证行为性描述；fixture 需求可审计 | 改动涉及 3 skill + 2 rubric | **Selected: 唯一能系统性解决三层根因的方案** |

## Feasibility Assessment

### Technical Feasibility

完全可行。所有改动都在现有 skill 的规则文件、模板文件和 rubric 文件中，不涉及架构变更。

### Resource & Timeline

修改量：3 个 skill 的规则/模板 + 2 个 eval rubric。建议实施顺序：gen-journeys → gen-contracts → gen-test-scripts（因为 Fixture Specification schema 需先在 gen-contracts 中定义，gen-test-scripts 再消费），eval rubric 修改可与对应 skill 并行。

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
| Golden Path 强制要求导致简单 feature 的 Journey 过度膨胀 | M | L | Golden Path 可以是简洁的多步 CRUD 循环；按 Feature 复杂度分类启发式规则（见上文）区分简单/复杂 feature 的期望 |
| Contract Fixture Specification 增加 Contract 复杂度，降低 eval 得分 | M | M | Fixture Specification 作为 Preconditions 的子维度，增加权重但独立评分；新 rubric 维度设 60% 最低通过阈值，且关键子项一票否决 |
| 断言深度规则过于严格，导致某些合法测试被误判 | L | M | 80% 阈值允许 20% 的"结构性"断言存在，覆盖 logging、health check 等合理场景；断言分类判据表提供明确边界 |
| 三个 skill 共享 fixture_spec schema 的集成风险 | M | M | 先定义 schema 规范（gen-contracts），再实现消费端（gen-test-scripts）；schema 变更需同步更新三个 skill 的规则文件 |

## Success Criteria

consistency_check_result:
  status: pass
  pairs_checked: 7
  conflicts_found: 0

- [ ] SC-1: 每个包含父子实体关系的 feature，gen-journeys 至少生成 1 个 Golden Path Journey，该 Journey 必须覆盖 PRD/Design 中 primary user story 的核心领域动作序列（语义完整性约束），且跨越 3+ 步骤
- [ ] SC-2: 每个 Contract 的 Preconditions 包含 Fixture Specification（声明前置实体类型、关系类型、最小数量）
- [ ] SC-3: gen-test-scripts 生成的测试中，≥80% 的断言验证业务结果（按断言分类判据表定义：实体存在、状态正确、关系完整、业务规则满足），而非仅验证 HTTP 状态码或响应 schema
- [ ] SC-4: 当 Fixture Specification 声明需要 N 个子实体时，生成的测试 fixture 必须创建 ≥N 个子实体
- [ ] SC-5: Journey eval rubric 新增 "Workflow Coverage" 维度（150 分），评分标准包含 Golden Path 存在性和多步覆盖度。最低通过阈值：该维度得分 ≥90/150（60%），且 Golden Path 存在性子项不得为 0 分（一票否决）。防 checkbox-compliant 机制：eval prompt 要求评审者验证 Golden Path 的步骤序列是否对应 PRD/Design 中的具体用户故事，而非仅检查步骤数量
- [ ] SC-6: Contract eval rubric 新增 "Fixture Specification" 维度（100 分），评分标准包含前置数据声明完整性和实体关系覆盖度。最低通过阈值：该维度得分 ≥60/100（60%），且 fixture_spec.entities 必须包含 Contract 涉及的所有实体类型（完整性子项一票否决）。防 checkbox-compliant 机制：eval prompt 要求评审者验证 entity_type 是否与 Design 中的领域模型一致，而非仅检查字段是否存在
- [ ] SC-7: 以 pm-work-tracker 里程碑地图 feature 为回归基准，完成全管线（gen-journeys → gen-contracts → gen-test-scripts → run-tests）端到端验证：(a) 生成的 Journey 包含"创建地图 → 添加里程碑 → 验证地图中包含里程碑"工作流；(b) 生成的测试在里程碑地图为空时应失败（而非虚假通过）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
