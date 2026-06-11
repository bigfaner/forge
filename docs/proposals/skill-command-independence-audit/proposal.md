---
created: "2026-06-03"
author: "faner"
status: Draft
intent: "cleanup"
---

# Proposal: Skill & Command Independence Audit

## Problem

Forge plugin 中的 21 个 skill 和 16 个 command 存在三类文档质量问题：(1) 5 个 skill 存在跨 skill 内部文件引用，破坏了 skill 独立性（涉及 6 处引用关系）；(2) 多个 skill/command 包含冗余描述信息（total ~6000 行中有约 30% 可精简）；(3) 9 个 skill 的 Related Skills/Integration/References 章节仅罗列 pipeline 上下游，不提供附加信息量。

### Evidence

- **跨 skill 引用**: gen-journeys 引用 gen-contracts/rules/journey-contract-model.md（3次）、gen-test-scripts 引用 run-tests/rules/test-isolation.md、extract-design-md 引用 ui-design/templates/styles/（运行时数据读取，见 Scope 豁免说明）、init-justfile 引用 test-guide/references/test-type-model.md、gen-contracts 引用 gen-journeys/SKILL.md "Surface Detection" section（反向引用，与 gen-journeys→gen-contracts 形成双向耦合）
- **跨 skill 到 command**: fix-bug command 引用 learn/templates/ 和 consolidate-specs/rules/
- **冗余**: quick-tasks 与 breakdown-tasks 共享 ~150 行几乎相同内容；write-prd 与 tech-design 重复 Override Signals 表；execute-task 与 run-tasks command 共享约 20-30 行接口契约（claim 格式 + fix-type 表），核心逻辑各自独立；tech-design 4 种 intent 变体膨胀至 445 行；eval 334 行中约 100 行为 proposal-only 特性
- **Related 无用信息**: 9 个 skill 的 Related Skills/Integration/References 章节内容均可从正文中隐含推断

### Urgency

v3.0.0 开发阶段是清理文档债务的窗口期。随着 skill 数量增长，跨 skill 耦合会导致修改一处必须同步检查其他 skill，维护成本将持续上升。

## Proposed Solution

对有跨引用、冗余或低价值章节的 13 个 skill 和 3 个 command 执行三维度清理：
1. **消除跨引用**: 将引用的外部知识内联到引用方，使每个 skill/command 完全自洽（forensic 的动态 SKILL.md 加载除外——这是设计意图）
2. **精简描述**: 在维持自洽前提下压缩冗余展开，用简明描述替代冗长表格/多行说明
3. **删除 Related 章节**: 移除所有 Related Skills / Integration / References 章节，因为 pipeline 上下游关系已在正文流程中体现

### Innovation Highlights

无创新，标准文档清理。核心原则是"每个 skill 文件是一个独立的知识单元"。

## Requirements Analysis

### Key Scenarios

- AI agent 加载单个 skill 时，无需读取其他 skill 的内部文件即可完整理解并执行
- 修改一个 skill 时，不需要同步修改其他 skill
- 新开发者阅读某个 skill 时，该文件自洽、无悬挂引用

### Non-Functional Requirements

- 清理后所有 skill/command 仍能正确指导 AI agent 行为（功能等价）
- 总行数减少但不丢失关键决策信息

### Constraints & Dependencies

- forensic skill 的动态 SKILL.md 加载机制保留不动
- 清理仅涉及文档（.md 文件），不涉及代码逻辑

## Alternatives & Industry Benchmarking

### Industry Solutions

标准做法是"模块自包含"——每个模块的文档携带所需的全部上下文。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 耦合持续增长，维护成本上升 | Rejected: 债务会累积 |
| 提升到共享层 | 模块化标准 | 单份权威 | 引入新路径约定，增加分发复杂度 | Rejected: 过度设计 |
| 内容内联+精简 | 自包含模块 | 彻底独立，分发简单 | 知识多份存在，可能漂移 | **Selected: 符合 Forge 分发模型** |

## Feasibility Assessment

### Technical Feasibility

纯文档编辑，无技术风险。

### Resource & Timeline

预计 1 个 session 可完成。工作量主要集中在 5 个有跨引用的 skill + 1 个有跨引用的 command（fix-bug）+ 9 个有 Related 章节的 skill 的编辑。

### Dependency Readiness

无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| Related Skills 章节帮助 AI agent 理解 pipeline 上下文 | Occam's Razor | Refined: 正文流程已包含上下游关系，Related 章节是冗余信息 |
| 跨 skill 引用是"复用"优于"重复" | Assumption Flip | Overturned: 对 AI agent 而言，独立加载更可靠；多份拷贝的漂移风险低于跨 skill 耦合的维护负担 |
| forensic 需要特殊处理 | 5 Whys | Confirmed: forensic 的核心功能就是对比其他 skill 的定义与实际行为，动态引用是设计意图 |

## Scope

### In Scope

**Skill 跨引用修复（5 个 skill）:**
- gen-contracts: 内联 gen-journeys/SKILL.md "Surface Detection" section 所需知识（INJECT: surface 检测规则段落 ~20 行；SKIP: surface 生成步骤）；此为 gen-journeys→gen-contracts 反向引用的对应端，修复后消除双向耦合
- gen-journeys: 内联 gen-contracts/rules/journey-contract-model.md 所需内容（INJECT: Contract 结构定义 + Outcome 语义 ~60 行；SKIP: 代码示例和实现细节）
- gen-test-scripts: 内联 run-tests/rules/test-isolation.md 所需内容（INJECT: 隔离策略决策表 ~40 行；SKIP: run-tests 调度逻辑）
- extract-design-md: 对 ui-design/templates/styles/ 的引用是运行时数据读取而非知识引用，不适用内联。处理策略：在 extract-design-md 内部创建 rules/style-matching.md 包含匹配特征摘要，风格文件保留在 ui-design 中（设计意图豁免）
- init-justfile: 内联 test-guide/references/test-type-model.md 所需内容（INJECT: test-type 映射表 ~30 行；SKIP: 详细说明和示例）

**Command 跨引用修复（1 个 command）:**
- fix-bug: 内联 learn/templates 和 consolidate-specs/rules 所需内容（INJECT: 模板决策点 + spec 提取规则 ~40 行；SKIP: 其他 command 的上下文）

**Related Skills / Integration / References 章节删除（9 个 skill）:**
- consolidate-specs, gen-test-scripts, run-tests, tech-design, ui-design, write-prd：删除 Related Skills / Integration / References 章节（内容可从正文隐含推断）
- gen-contracts：删除 Related Skills / Integration 章节；**References 段落合并到内联知识中作为定义段落**（包含 Contract、Outcome、Semantic Descriptors 等 6 个概念定义，正文中被引用但未集中定义）
- gen-journeys：删除 Related Skills / Integration 章节；**References 段落合并到内联知识中作为定义段落**（包含正文中引用的概念定义）
- quick-tasks：仅删除 ## Integration 段落；**保留 ## Reference Files**（第 122 行，是模板占位符 {{REFERENCE_FILES}} 的替换规则说明，非 pipeline 上下游信息）

**冗余精简（重点 skill）:**
- quick-tasks: 精简与 breakdown-tasks 共享的内容描述
- breakdown-tasks: 精简共享内容
- write-prd: 精简 4 种 intent 变体的重复展开
- tech-design: 精简 4 种 intent 变体、Override Signals 表
- gen-journeys: 精简 5 个 per-surface 内联摘要
- eval: 精简 proposal-only 特性描述
- init-justfile: 精简 justfile 示例

**冗余精简（重点 command）:**
- execute-task: 精简与 run-tasks 重叠的逻辑
- run-tasks: 精简重叠逻辑
- fix-bug: 精简 Knowledge Review 段落

### Out of Scope

- forensic skill 的动态 SKILL.md 加载机制（设计意图，不是耦合问题）
- 功能行为变更（纯文档清理，不改变任何运行时行为）
- eval-* command stubs（已是纯委托，无清理必要）
- git-checkout, git-commit, clean-code, extract-design-md commands（已经足够精简）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 内联知识时遗漏关键信息导致 skill 执行不完整 | M | M | 内联后对比原文确保无遗漏 |
| 精简过度导致 AI agent 行为偏差 | L | M | 精简时保留所有硬规则和决策表，只压缩描述性文字 |
| 多份拷贝的知识在未来修改时未同步更新 | M | L | 可接受——独立性带来的维护简化大于同步成本；对内联段落使用 `<!-- INLINE:origin=<skill>/<file>#<section> -->` 标记提供可追溯性 |

## Success Criteria

- [ ] 0 处跨 skill 内部文件引用（forensic 的动态加载除外）
- [ ] 0 处 command 引用 skill 内部文件
- [ ] 0 个纯粹的 pipeline 上下游 Related/Integration/References 章节（gen-contracts/gen-journeys 的概念定义段落已合并到内联知识，quick-tasks 的 Reference Files 使用说明已保留）
- [ ] 总行数净减少 ≥ 10%
- [ ] 所有 skill/command 修改后功能等价，验证清单：
  - 所有 HARD-RULE / HARD-GATE / EXTREMELY-IMPORTANT / PROHIBITIONS 块计数不变
  - 所有决策表完整保留
  - 所有 Step 序号和流程步骤完整

## Next Steps

- Proceed to `/quick` for streamlined implementation (no PRD/design needed for doc cleanup)
