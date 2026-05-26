---
created: "2026-05-26"
author: "faner"
status: Draft
---

# Proposal: Brainstorm 阶段增加受影响点分析

## Problem

当用户通过 brainstorm 创建涉及**改动现有功能**的提案时，当前流程没有机制识别和列举该改动对现有代码库的波及范围。影响分析被推迟到 PRD（Related Changes 表格）、tech-design（ER Diagram Change Impact Analysis）甚至任务执行阶段（`breaking: true` 标志、IMPACT_DECLARATION）才进行，导致方向性错误在面试阶段未能被及时发现和纠正。

### Evidence

- `docs/lessons/gotcha-breaking-change-integration-test-blast-radius.md`: breaking task 仅范围化了直接单元测试，遗漏了 60+ 集成测试 fixture
- `docs/lessons/gotcha-breaking-change-quality-gate-deadlock.md`: breaking task 在质量门禁处死锁，因为 `go build` 在调用方编译失败
- 这些问题的根源都是在提案阶段缺少影响分析

### Urgency

每次遗漏影响分析都导致下游返工（任务拆分不当、质量门禁死锁、测试失败）。越早发现影响，修正成本越低。

## Proposed Solution

在 brainstorm skill 中深度嵌入影响分析能力：Step 1 自动扫描用户项目代码库检测现有功能变动，Step 2 面试中与用户交互确认受影响点，最终在 proposal.md 模板中新增 `## Affected Points` 章节记录分析结果。

### Innovation Highlights

将"波及范围分析"前置到最早可能的阶段（提案），而非等到设计或执行阶段。受影响点按四个维度（代码模块、API/接口契约、测试覆盖、配置/数据）分类标注影响类型（新增/修改/删除/不变），为下游所有阶段提供一致的参考基线。

## Requirements Analysis

### Key Scenarios

- **改动现有功能**: 用户提案涉及修改已有模块/接口/配置 → 自动触发影响扫描，四维度列举受影响点
- **纯新增功能**: 用户提案不涉及现有代码 → 不触发影响分析，不产生空章节
- **混合场景**: 部分新增部分修改 → 仅对修改部分触发影响分析
- **用户补充**: 自动扫描可能遗漏间接引用 → 面试中用户可补充和修正受影响点列表

### Non-Functional Requirements

- 扫描策略须语言无关，适用于任何技术栈的项目
- 扫描不应显著延长 brainstorm 流程时间（纯新增功能的提案零额外开销）

### Constraints & Dependencies

- 依赖 agent 的 Grep/Glob/Read 工具能力进行代码扫描
- 依赖 agent 对用户项目的代码结构理解（Step 1 Analyze Context 已有此能力）
- 受限于 agent 上下文窗口，扫描深度需要合理控制

## Alternatives & Industry Benchmarking

### Industry Solutions

静态分析工具（SonarQube、CodeClimate）通过 AST 分析计算变更影响范围，但需要特定语言支持且集成成本高。IDE 的 "Find Usages" 功能提供即时引用查找，但需要人工逐个确认。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动成本 | 影响遗漏导致下游返工，已有 lessons 证明 | Rejected: 已有实际痛点 |
| 仅扩展模板 | 行业常见做法 | 最小改动 | 发现太晚，面试中已做方向决策 | Rejected: 错过最佳修正时机 |
| Step 1 扫描 + Step 2 确认 + 模板新章节 | 借鉴 refactor IMPACT_DECLARATION 的前置思路 | 最早发现影响，面试中可调整方向 | 增加 brainstorm 流程复杂度 | **Selected: 最早发现 = 最低修正成本** |

## Feasibility Assessment

### Technical Feasibility

brainstorm skill 已有 Step 1 Analyze Context 阶段搜索代码库的能力。增强为影响扫描只需扩展搜索策略和分类逻辑，不引入新的工具依赖。

### Resource & Timeline

涉及 3 个文件的改动（1 新增 rule + 2 修改现有文件），属于 skill 指令文档的修改，无需编写代码。

### Dependency Readiness

无外部依赖。Grep/Glob/Read 工具已可用。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 影响分析在 PRD/tech-design 阶段做就够了 | 5 Whys | Overturned: lessons 证明 PRD 阶段的 Related Changes 表格缺乏系统性发现机制，tech-design 阶段仅覆盖 DB schema |
| agent 能可靠检测"改动现有功能" vs "纯新增" | Assumption Flip | Refined: 检测依赖 Step 1 的代码库搜索结果和用户描述中的关键词。边界情况（如"重构"类提案）需要在 Step 2 确认 |

## Scope

### In Scope

- 新增 `rules/impact-analysis.md`：定义扫描触发条件、四维度扫描策略、影响类型分类标准、输出格式
- 修改 `SKILL.md`：Step 1 增加智能检测和扫描指令，Step 2 增加受影响点确认环节
- 修改 `templates/proposal.md`：新增 `## Affected Points` 条件性章节

### Out of Scope

- PRD 的 `Related Changes` 表格变更（已有机制，层级不同）
- Tech Design 的 `Change Impact Analysis` 变更（DB schema 专用）
- coding-refactor 的 `IMPACT_DECLARATION` 变更（执行阶段机制）
- `breaking: true` 任务标志机制变更

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 智能检测误判（纯新增功能被误判为改动现有功能） | M | L | Step 2 确认环节允许用户修正；误判仅导致额外扫描，不影响提案质量 |
| 扫描不完整，遗漏间接引用 | M | M | 面试中用户可补充；rule 中定义最大扫描深度限制 |
| 四维度分类在非典型项目中不适用 | L | L | 每个维度标注为可选，agent 根据项目实际情况跳过不相关维度 |

## Success Criteria

- [ ] 涉及改动现有功能的提案自动包含 `## Affected Points` 章节，纯新增功能的提案不包含该章节
- [ ] Affected Points 章节按四个维度分组（代码模块/API接口/测试覆盖/配置数据），每条标注影响类型（新增/修改/删除/不变）
- [ ] Step 2 面试中包含受影响点的交互确认环节，用户可补充和修正扫描结果
- [ ] 下游 skill（write-prd、tech-design、breakdown-tasks）可引用 Affected Points 章节的信息
- [ ] 新 rule 文件定义语言无关的扫描策略，适用于 Go/TypeScript/Python/Java 等项目

## Next Steps

- Proceed to `/write-prd` to formalize requirements
