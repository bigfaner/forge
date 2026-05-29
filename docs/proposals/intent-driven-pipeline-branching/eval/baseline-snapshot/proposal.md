---
created: "2026-05-29"
author: fanhuifeng
status: Draft
intent: new-feature
---

# Proposal: Intent-Driven Pipeline Branching

## Problem

Forge 的测试管道对**所有** coding 任务无差别生成 Journey → Contract → Test Script 全链路，但纯重构（`coding.refactor`）和代码清理（`coding.cleanup`）的特征不产生新的用户可观测行为，强行走 Journey 管道导致：生成器凭文档推测事实（常量名、CLI 命令未验证）、eval 循环"用新猜测修旧猜测"评分反而下降（630→585）、结构性问题无法通过迭代修复。

### Evidence

- `unify-enum-constants` 特征的 eval-journey 评分持续不达标（466→630→585，目标 850），3 轮迭代后反而下降
- `gen-journeys` SKILL.md 硬性规则禁止读取代码（"Do not read source code, test files, or implementation files"），所有事实声明纯靠 PRD 推测
- `build.go` 中 `IsTestableType()` 对所有 `coding.*` 类型返回 `true`，`coding.refactor` 和 `coding.cleanup` 触发完整测试管道
- PRD 的 user stories 格式（"As a user / I want / So that"）对纯重构语义为空

### Urgency

该问题在 `unify-enum-constants` 特征上已造成 3 轮无效迭代，浪费了 eval 计算。更根本的是：只要 Forge 处理纯重构或清理任务，此问题**必然复现**。每多一个此类特征，就多一轮无效的 journey 生成 + eval 循环。

## Proposed Solution

引入两个正交维度驱动 Pipeline 选择：

1. **管道模式**（Pipeline Mode）：基于文档存在性——Breakdown（有 PRD）或 Quick（仅有 proposal）
2. **特征意图**（Feature Intent）：基于工作内容——`new-feature` / `refactor` / `cleanup`

Intent 驱动的是**测试管道段**的选择，不是管道模式的选择。完整矩阵：

| | Quick 模式 | Breakdown 模式 |
|---|---|---|
| **new-feature** | proposal → quick-tasks → gen-journeys → run-tests → validate → clean-code → doc-drift | proposal → write-prd → tech-design → breakdown-tasks → full test pipeline → validate → clean-code → consolidate |
| **refactor** | proposal → quick-tasks → quality-gate → done | proposal → write-prd(spec-only) → tech-design(internal) → breakdown-tasks → quality-gate → consolidate |
| **cleanup** | proposal → quick-tasks → quality-gate → done | *(不适用，cleanup 不走 Breakdown 模式)* |

关键区别：
- **`new-feature`**：现有完整管道不变
- **`refactor`**：PRD 内部分支（跳过 user stories，只生成 spec），测试管道跳过 Journey/Contract/Script，验证依赖已有的 quality-gate hook（compile + fmt + lint + test）
- **`cleanup`**：始终走 Quick 模式，测试管道跳过，验证依赖 quality-gate hook

### Innovation Highlights

- **Intent 作为 Pipeline 一等公民**：不是在单个 skill 内部做 if/else，而是 intent 在 proposal 阶段确定后，驱动整个 pipeline 拓扑选择。`forge task index` 根据特征目录中的 intent 生成不同的测试任务链。
- **回归验证替代 Journey 验证**：对行为保持型任务，验证目标是"无回归"而非"新行为正确"。已有 quality-gate hook 天然承担此职责，无需新建验证基础设施。
- **PRD 自适应格式**：write-prd 内部分支——重构场景不生成 user stories，改为生成"变更范围 + 约束条件 + 验证标准"格式的 spec。

## Requirements Analysis

### Key Scenarios

1. **Breakdown + refactor**：用户发起"字符串字面量→类型常量"的重构 proposal → AI 推断 intent 为 `refactor` → 用户确认 → write-prd 生成 spec-only PRD（跳过 user stories）→ tech-design 侧重内部架构 → breakdown-tasks 使用 `coding.refactor` 类型 → 测试管道跳过 journey/contract/script → quality-gate 验证无回归
2. **Quick + refactor**：用户发起重构 proposal，选择 Quick 模式 → quick-tasks 使用 `coding.refactor` 类型 → 测试管道跳过 → quality-gate 验证
3. **Quick + cleanup**：用户发起"移除死代码"的 cleanup proposal → AI 推断 intent 为 `cleanup` → 用户确认 → quick-tasks 使用 `coding.cleanup` 类型 → 测试管道跳过 → quality-gate 验证
4. **Breakdown + new-feature**（默认）：用户发起新功能 proposal → intent 为 `new-feature`（默认）→ 现有完整管道不变
5. **Quick + new-feature**：用户发起新功能 proposal，选择 Quick 模式 → 现有 Quick 管道不变
6. **Intent 推断边界**：proposal 内容模糊（既有新行为又有重构）→ AI 推断为主 intent → 用户确认或覆盖

### Non-Functional Requirements

- **向后兼容**：缺少 intent 字段的已有 proposal 默认为 `new-feature`，不改变现有行为
- **最小侵入**：修改集中在 `build.go`（IsTestableType + autogen）和 skill 文档（brainstorm template, write-prd, tech-design），不改变 task 类型定义和 status 状态机

### Constraints & Dependencies

- `intent` 字段需持久化在 proposal.md frontmatter 中，供 `forge task index` 读取
- quality-gate hook 已存在且覆盖 compile+fmt+lint+test，无需新建
- breakdown-tasks 和 quick-tasks 已有 Intent Propagation 逻辑（从 proposal.md 的 intent 字段传播默认 task type），可复用

## Alternatives & Industry Benchmarking

### Industry Solutions

CI/CD 系统普遍根据变更类型选择管道：GitHub Actions 用 `paths` 过滤触发条件，Bazel 用 `test_suite` 区分单元/集成测试，Maven/Gradle 有独立的 `compile`/`test`/`verify` 生命周期阶段。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 每个重构/清理特征都浪费 eval 迭代 | Rejected: 问题必然复现 |
| gen-journeys 内部检测 task 类型跳过 | 本地 | 改动小 | 不解决 PRD 语义为空的问题，只隐藏下游症状 | Rejected: 治标不治本 |
| 修复 gen-journeys 的代码验证 + eval 循环 | lesson 建议 | 对所有特征都有用 | 不解决纯重构不应走 journey 的根本问题 | Rejected: 方向错误，重构根本不需要 journey |
| **Intent-Driven Pipeline Branching** | 本 proposal | 从源头解决：不同意图走不同管道 | 需要修改 forge-cli + 多个 skill | **Selected: 最彻底，且复杂度可控** |

## Feasibility Assessment

### Technical Feasibility

完全可行。按改动位置分两层：

#### Skill 层（AI 行为引导）

1. **brainstorm template**：添加 `intent` frontmatter 字段
2. **brainstorm SKILL.md**：添加 AI 推断 intent 的步骤，用 AskUserQuestion 确认
3. **write-prd SKILL.md**：添加内部分支逻辑——refactor 跳过 user stories
4. **tech-design SKILL.md**：添加内部分支逻辑——refactor 侧重内部架构，跳过 API handbook / ER 图

#### CLI 层（forge-cli Go 代码）

5. **build.go — 读取 intent**：`BuildIndex()` 在 `detectMode()` 之后读取 `proposal.md` frontmatter 的 `intent` 字段，缺失默认 `"new-feature"`
6. **build.go — needsTestPipeline() 分支**：增加 `intent` 参数，`refactor`/`cleanup` 返回 `false`，跳过测试管道生成
7. **build.go — needsTestPipeline() 以下连带影响**：
   - `GenerateTestTasks()` 对 refactor/cleanup 返回空列表（不生成 gen-journeys/gen-contracts/gen-scripts/run-tests）
   - 下游任务（validate-code, clean-code, consolidate-specs）的依赖链需重新挂载
8. **autogen.go — 依赖接线**：`resolveBreakdownDeps()` 和 `resolveQuickDeps()` 感知 intent，refactor/cleanup 时跳过 run-tests 节点，把下游任务直接接到 business tasks 尾部：

```
# 当前 new-feature（Breakdown）：
business tasks → gen-journeys → eval-journey → gen-contracts → eval-contract → gen-scripts → run-tests → validate-code → clean-code → consolidate

# refactor（Breakdown）：
business tasks → validate-code → clean-code → consolidate

# cleanup（Quick）：
business tasks → clean-code → doc-drift
```

无需新增 task 模板——refactor/cleanup 不生成测试任务，business tasks 完成后直接由 quality-gate hook 验证。

### Resource & Timeline

中等规模改动：forge-cli Go 代码修改（build.go + autogen.go）+ 4 个 skill 文档更新。无外部依赖。

### Dependency Readiness

quality-gate hook 已就绪。Intent Propagation 机制已在 breakdown-tasks / quick-tasks 中存在，只需扩展到 pipeline 级别。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "所有 coding 任务都需要 Journey 测试" | Assumption Flip：如果不是呢？ | **Overturned**：纯重构/清理无新用户可观测行为，Journey 语义为空 |
| "PRD 的 user stories 格式对所有特征都适用" | Stress Test：用 unify-enum-constants 验证 | **Overturned**："As a developer I want type-safe constants" 不是有效的 user story |
| "多轮 eval 迭代总能提高质量" | 5 Whys：为什么评分下降？ | **Overturned**：修订引入新的未验证声明，结构性问题无法通过迭代修复 |
| "Intent 只影响 task type，不影响 pipeline topology" | Assumption Flip：如果 intent 也驱动 pipeline 呢？ | **Confirmed**：intent 应该是一等公民，驱动整个 pipeline 选择 |

## Scope

### In Scope

- proposal.md template 添加 `intent` 字段
- brainstorm SKILL.md 添加 intent 推断 + 用户确认步骤
- write-prd SKILL.md 添加 refactor 内部分支（跳过 user stories，生成 spec-only PRD）
- tech-design SKILL.md 添加 refactor 内部分支（侧重内部架构）
- forge-cli `build.go`：`IsTestableType()` 区分行为变更 vs 行为保持
- forge-cli `autogen.go`：`GetBreakdownTestTasks()` 和 `GetQuickTestTasks()` 根据 intent 跳过 journey/contract/script 管道段
- forge-cli `build.go`：`needsTestPipeline()` 读取 intent，refactor/cleanup 跳过完整测试管道生成

### Out of Scope

- gen-journeys 的代码验证机制（Level 1 修复）—— 后续迭代
- eval-journey 修订循环的代码回查（Level 2 修复）—— 后续迭代
- eval-journey 结构性问题检测（Level 3 修复）—— 后续迭代
- 混合 intent 支持（一个 proposal 包含多种 intent）—— 不支持，按主要意图归类
- 意图自动推断的精确度优化—— 首版基于关键词匹配 + AI 推断，后续迭代优化

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Intent 推断错误，导致重构走了完整管道 | M | L — 只是多跑了不必要的步骤 | AI 推断后用户确认，可覆盖；默认 new-feature 不会丢失覆盖 |
| Intent 推断错误，导致新功能跳过了 journey | L | H — 测试覆盖不足 | 默认 new-feature，只有显式 refactor/cleanup 才跳过；用户必须确认 intent |
| refactor PRD spec-only 格式与下游 skill 不兼容 | M | M — tech-design / breakdown-tasks 可能期望 user stories | write-prd 分支确保 spec 格式包含 tech-design 需要的字段 |
| 已有 proposal 缺少 intent 字段 | H | L — 默认 new-feature | 向后兼容：缺失 intent = new-feature，行为不变 |

## Success Criteria

- [ ] `intent: refactor` 的 proposal（Breakdown 模式）生成 spec-only PRD（无 user stories），测试管道跳过 journey/contract/script
- [ ] `intent: refactor` 的 proposal（Quick 模式）不生成 gen-journeys 任务，quality-gate 验证无回归
- [ ] `intent: cleanup` 的 proposal 不生成 gen-journeys 任务，quality-gate 验证无回归
- [ ] `intent: new-feature` 的 proposal（Breakdown 和 Quick 模式）行为与当前完全一致（回归验证）
- [ ] `forge task index` 对 refactor/cleanup feature 不生成 test pipeline 任务（gen-journeys/gen-contracts/gen-scripts/run-tests）
- [ ] write-prd 对 `intent: refactor` 生成 spec-only PRD（无 prd-user-stories.md 文件）
- [ ] 已有 proposal（无 intent 字段）默认走 new-feature 管道，行为不变
- [ ] refactor/cleanup 特征的 business tasks 完成后，quality-gate hook 执行 compile+fmt+lint+test 验证无回归

## Next Steps

- Proceed to `/write-prd` to formalize requirements
