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
2. **特征意图**（Feature Intent）：基于工作内容——`new-feature` / `refactor` / `cleanup`（task type → intent 映射规则：`coding.feature` 和 `coding.enhancement` 默认归为 `new-feature`（引入新的用户可观测行为）；`coding.cleanup` 归为 `cleanup`；`coding.refactor` 归为 `refactor`；`coding.fix` 按"是否引入新的用户可观测行为"判断——引入新行为归为 `new-feature`，仅内部调整归为 `refactor`。此判断由 brainstorm skill 的 AI 推断完成，用户确认后写入 proposal frontmatter）

Intent 驱动的是**测试管道段**的选择，不是管道模式的选择。完整矩阵：

| | Quick 模式 | Breakdown 模式 |
|---|---|---|
| **new-feature** | proposal → quick-tasks → gen-journeys → run-tests → validate-code → clean-code → doc-drift | proposal → write-prd → tech-design → breakdown-tasks → gen-journeys → eval-journey → gen-contracts → eval-contract → gen-scripts → run-tests → validate-code → clean-code → consolidate-specs |
| **refactor** | proposal → quick-tasks → clean-code → doc-drift | proposal → write-prd(spec-only) → tech-design(internal) → breakdown-tasks → validate-code → clean-code → consolidate-specs |
| **cleanup** | proposal → quick-tasks → clean-code → doc-drift | *(不适用——cleanup 始终走 Quick 模式，`build.go` 在 intent=`cleanup` 时强制 `mode=Quick`，忽略文档存在性)* |

关键区别：
- **`new-feature`**：现有完整管道不变
- **`refactor`**：PRD 内部分支（跳过 user stories，只生成 spec），测试管道跳过 Journey/Contract/Script，验证依赖已有的 quality-gate hook（compile + fmt + lint + test）。spec-only PRD 必须包含以下字段以满足 tech-design 输入需求：变更范围（affected modules/files）、约束条件（behavioral invariants to preserve）、验证标准（regression acceptance criteria）。tech-design 的 `prd-user-stories.md` 文件在 refactor 下不生成
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
6. **Intent 推断边界**：proposal 内容模糊（既有新行为又有重构）→ AI 按"是否引入新的用户可观测行为"判断主要 intent：若 proposal 的核心目标包含任何新增的外部可观测行为（新 API、新 CLI 命令、新输出格式），则主要 intent 为 `new-feature`，走完整测试管道以确保新行为有测试覆盖；若核心目标仅为重组内部实现，则主要 intent 为 `refactor`。用户可在确认阶段覆盖此判断。

### Non-Functional Requirements

- **向后兼容**：缺少 intent 字段的已有 proposal 默认为 `new-feature`，不改变现有行为
- **最小侵入**：修改集中在 `build.go`（`needsTestPipeline()` 加 intent 短路）、`autogen.go`（生成函数 + 接线函数感知 intent）和 skill 文档（brainstorm template, write-prd, tech-design），不改变 task 类型定义和 status 状态机。`IsTestableType()` 不修改，保持纯类型判断职责

### Constraints & Dependencies

- `intent` 字段需持久化在 proposal.md frontmatter 中，供 `forge task index` 读取。若 proposal.md 不存在（用户在未完成 brainstorm 的特征目录上运行 `forge task index`），`proposal.FindBySlug()` 返回空 Proposal，此时 CLI handler 将 `opts.Intent` 设为默认值 `new-feature`，行为与当前一致。不需要额外的错误处理或中断
- quality-gate hook 已存在且覆盖 compile+fmt+lint+test，无需新建
- breakdown-tasks 和 quick-tasks 已有 Intent Propagation 逻辑（从 proposal.md 的 intent 字段传播默认 task type），可复用

## Alternatives & Industry Benchmarking

### Industry Solutions

CI/CD 和构建系统普遍根据变更类型或目标选择管道拓扑：

- **GitHub Actions** 用 `paths` 过滤器决定是否触发工作流——只有 `src/` 变更触发 test job，`docs/` 变更触发 deploy job。本质是**基于变更范围的条件路由**，与 Forge 的 intent 驱动类似，但 GitHub Actions 的粒度是文件路径而非语义意图。Forge 无法用文件路径区分 refactor 和 new-feature（两者都改代码），因此需要语义层（intent）做路由决策。
- **Bazel** 用 `test_suite` 和 `tag` 过滤区分单元/集成/端到端测试，`query` 命令可按依赖图裁剪执行范围。Bazel 的模型是**基于依赖图的选择性执行**——只运行受变更影响的测试。Forge 的 refactor/cleanup 场景类似"无依赖变更"（无新的用户可观测行为），因此对应"跳过端到端测试，只跑回归验证"的策略。
- **Maven/Gradle** 有独立的 `compile`/`test`/`verify` 生命周期阶段，每个 phase 可绑定不同的 plugin goal。本质是**声明式阶段管道**——不同 goal 可以选择跳过某些 phase（如 `-DskipTests`）。Forge 的 intent 类似一个自动化的 `-DskipTests` 标志，但由语义推断而非手动指定。

**关键差异**：上述系统都是基于文件路径或依赖图的**结构化信号**做路由，而 Forge 引入了**语义信号**（intent）做路由。这是因为 Forge 的任务类型（如 `coding.refactor` vs `coding.feature`）在文件结构上不可区分，只有语义意图不同。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 每个重构/清理特征都浪费 eval 迭代 | Rejected: 问题必然复现 |
| gen-journeys 内部检测 task 类型跳过 | 本地 | 改动小 | 不解决 PRD 语义为空的问题，只隐藏下游症状 | Rejected: 治标不治本 |
| 修复 gen-journeys 的代码验证 + eval 循环 | lesson 建议 | 对所有特征都有用 | 不解决纯重构不应走 journey 的根本问题 | Rejected: 方向错误，重构根本不需要 journey |
| **Intent-Driven Pipeline Branching** | 本 proposal | 从源头解决：不同意图走不同管道 | 需要修改 forge-cli + 多个 skill；autogen.go 需覆盖 5 种有效接线场景（new-feature Breakdown/Quick、refactor Breakdown/Quick、cleanup Quick）+ 零 business task 边界分支 | **Selected: 最彻底，且复杂度可控** |

## Feasibility Assessment

### Technical Feasibility

完全可行。按改动位置分两层：

#### Skill 层（AI 行为引导）

1. **brainstorm template**：添加 `intent` frontmatter 字段，默认值为 `new-feature`
2. **brainstorm SKILL.md**：添加 AI 推断 intent 的步骤，用 AskUserQuestion 确认
3. **write-prd SKILL.md**：添加内部分支逻辑——refactor 跳过 user stories
4. **tech-design SKILL.md**：添加内部分支逻辑——refactor 侧重内部架构，跳过 API handbook / ER 图

#### CLI 层（forge-cli Go 代码）

5. **build.go — 读取 intent**：将 `intent` 作为 `BuildIndexOpts` 结构体的显式字段传入（与 `Mode` 字段并列）。完整数据流如下：
   - **CLI handler**（`cmd/task.go`）调用 `proposal.FindBySlug(slug)` 获取 `Proposal` 结构体（`proposal.go` 已解析 frontmatter 中的 `intent` 字段）
   - CLI handler 将 `Proposal.Intent` 赋值给 `BuildIndexOpts.Intent`，传入 `BuildIndex(opts)`
   - `BuildIndex()` 内部不再重复解析 frontmatter，直接使用 `opts.Intent`；若 `opts.Intent` 为空则默认 `new-feature`
   - `BuildIndex()` 将 `opts.Intent` 传递给 `needsTestPipeline(tasks, intent)` 和 `autogen.go` 的生成函数及接线函数
   - 这避免了 `BuildIndex()` 内部重复解析 frontmatter，且与 `detectMode()` 的时序解耦——`detectMode()` 只负责 Quick/Breakdown 判断，intent 独立通过 opts 传入
6. **build.go — needsTestPipeline() 短路**：增加 `intent` 参数（签名变为 `needsTestPipeline(tasks map[string]Task, intent string) bool`）。当 intent 为 `refactor`/`cleanup` 时直接返回 `false`，不遍历 tasks 也不调用 `IsTestableType()`。`IsTestableType()` 不修改，保持纯类型判断职责（对所有 `coding.*` 返回 `true`）
7. **autogen.go — 生成函数内部跳过测试任务**：`GetBreakdownTestTasks()` 和 `GetQuickTestTasks()` 接收 intent 参数，当 intent 为 `refactor`/`cleanup` 时跳过测试任务（gen-journeys/gen-contracts/gen-scripts/run-tests）的生成，但仍生成 validate-code/clean-code 等验证任务。这确保 refactor/cleanup 下 autogen.go 仍被调用，validate-code/clean-code 通过现有生成路径产出
8. **autogen.go — 依赖接线**：`resolveBreakdownDeps()` 和 `resolveQuickDeps()` 感知 intent，refactor/cleanup 时跳过 run-tests 节点，把下游任务直接接到 business tasks 尾部。具体接线逻辑：
   - **零 business task 保护**：若 intent 为 refactor/cleanup 但 business task 列表为空（例如纯文档类型特征），则不生成 validate-code/clean-code 等下游任务——因为没有上游锚点，生成下游任务会产生悬空的 `depends_on` 引用。此场景下仅保留 quality-gate hook 作为验证手段。对 `new-feature` intent，现有逻辑已保证 business task 列表不为空（新功能必然有 coding task），因此零 business task 保护不影响 new-feature 行为
   - **refactor（Breakdown）**：最后一个 business task 的 taskID 作为 `validate-code` 的 `depends_on`；`clean-code` 依赖 `validate-code`；`consolidate-specs` 依赖 `clean-code`。不需要查找 `lastRunID`（因为 run-tests 节点不生成），直接取最后一个 business task 的 ID
   - **cleanup（Quick）**：最后一个 business task 的 taskID 作为 `clean-code` 的 `depends_on`；`doc-drift` 依赖 `clean-code`
   - **refactor（Quick）**：与 cleanup（Quick）相同的接线逻辑——最后一个 business task 的 taskID 作为 `clean-code` 的 `depends_on`；`doc-drift` 依赖 `clean-code`
   - **new-feature**：保持现有逻辑不变，`validate-code` 仍依赖 `run-tests` 的输出 taskID

```
# 当前 new-feature（Breakdown）：
business tasks → gen-journeys → eval-journey → gen-contracts → eval-contract → gen-scripts → run-tests → validate-code → clean-code → consolidate-specs

# refactor（Breakdown）：
business tasks → validate-code → clean-code → consolidate-specs
（validate-code.depends_on = 最后一个 business task 的 ID）

# cleanup（Quick）：
business tasks → clean-code → doc-drift
（clean-code.depends_on = 最后一个 business task 的 ID）
```

无需新增 task 模板——refactor/cleanup 不生成测试任务，business tasks 完成后直接由 quality-gate hook 验证。cleanup 的 Breakdown 路径不会到达 autogen.go，因为 build.go 在 intent=`cleanup` 时已将 mode 强制为 Quick。autogen.go 无需处理此组合。

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
- forge-cli `build.go`：`needsTestPipeline()` 增加 `intent` 参数做短路——当 intent 为 `refactor`/`cleanup` 时直接返回 `false`。`IsTestableType()` 不修改，保持纯类型判断职责
- forge-cli `build.go`：`detectMode()` 在 intent=`cleanup` 时强制返回 Quick 模式，忽略文档存在性
- forge-cli `autogen.go`：`GetBreakdownTestTasks()` 和 `GetQuickTestTasks()` 接收 intent 参数——当 intent 为 `refactor`/`cleanup` 时跳过测试任务（gen-journeys/gen-contracts/gen-scripts/run-tests）的生成，但仍生成 validate-code/clean-code 等验证任务
- forge-cli `autogen.go`：`resolveBreakdownDeps()` 和 `resolveQuickDeps()` 根据 intent 重新接线依赖链——refactor/cleanup 时下游任务（validate-code, clean-code）直接依赖最后一个 business task，跳过 run-tests 节点

### Out of Scope

- gen-journeys 的代码验证机制（Level 1 修复）—— 后续迭代
- eval-journey 修订循环的代码回查（Level 2 修复）—— 后续迭代
- eval-journey 结构性问题检测（Level 3 修复）—— 后续迭代
- 混合 intent 支持（一个 proposal 包含多种 intent）—— 不支持，按"是否引入新的用户可观测行为"判断主要意图归类（见 Key Scenarios #6）
- 意图自动推断的精确度优化—— 首版基于关键词匹配 + AI 推断，后续迭代优化
- gen-journeys 在 new-feature 场景下的事实幻觉问题（生成器凭文档推测常量名、CLI 命令等）—— 本 proposal 通过跳过 refactor/cleanup 的 journey 管道绕过此问题，但 new-feature 的 journey 仍可能产生未验证的事实声明。此问题需在后续迭代中通过 gen-journeys 代码验证机制（Level 1 修复）解决

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Intent 推断错误，导致重构走了完整管道 | M | L — 只是多跑了不必要的步骤 | AI 推断后用户确认，可覆盖；默认 new-feature 不会丢失覆盖 |
| Intent 推断错误，导致新功能跳过了 journey | L | H — 测试覆盖不足 | 默认 new-feature，只有显式 refactor/cleanup 才跳过；用户必须确认 intent |
| refactor PRD spec-only 格式与下游 skill 不兼容 | M | M — tech-design / breakdown-tasks 可能期望 user stories | write-prd 分支确保 spec 格式包含 tech-design 需要的三个字段（变更范围、约束条件、验证标准）；tech-design SKILL.md 添加 refactor 分支跳过 user-stories 引用 |
| 已有 proposal 缺少 intent 字段 | H | L — 默认 new-feature | 向后兼容：缺失 intent = new-feature，行为不变 |
| autogen.go 生成函数按 intent 跳过测试任务时 validate-code/clean-code 丢失 | M | M — refactor/cleanup 可能丢失 validate-code 任务 | 生成函数（`GetBreakdownTestTasks()`/`GetQuickTestTasks()`）内部按 intent 跳过测试任务时，需确保 validate-code/clean-code 仍正确生成。通过回归测试验证：对 `intent: new-feature` 生成完整管道（含 validate-code 依赖 run-tests），对 `intent: refactor` 跳过测试但保留 validate-code/clean-code 且依赖链重新接线到 business task |
| autogen.go 重构引入依赖接线 bug | M | H — 下游任务悬空或循环依赖 | 修改前为 autogen.go 添加单元测试覆盖所有 5 种接线场景（new-feature Breakdown/Quick, refactor Breakdown/Quick, cleanup Quick）+ 零 business task 边界情况 |

## Success Criteria

### Skill 层

- [ ] brainstorm skill 对新 proposal 能推断 intent（`new-feature`/`refactor`/`cleanup`）并通过 AskUserQuestion 让用户确认，确认后写入 proposal.md frontmatter 的 `intent` 字段。对包含混合内容（既有新行为又有重构）的 proposal，按"是否引入新的用户可观测行为"判断主要 intent；用户可在确认阶段覆盖
- [ ] brainstorm 对 `coding.fix` 类型 proposal 能根据"是否引入新的用户可观测行为"正确归类 intent（引入新行为 → `new-feature`，仅内部调整 → `refactor`）
- [ ] write-prd 对 `intent: refactor` 生成 spec-only PRD（无 `prd-user-stories.md` 文件），spec 包含变更范围、约束条件、验证标准三个必需字段
- [ ] `intent: refactor` 的 Breakdown 模式下，tech-design 输出不包含 API handbook 文件和 ER 图文件，不生成 `prd-user-stories.md`

### CLI 层

- [ ] `intent: new-feature` 的 Breakdown 模式 proposal，`forge task index` 输出的任务列表与未引入 intent 机制的基准输出完全相同（以当前 `main` 分支的输出为 golden file 基准，比对任务名称和依赖链结构，含 gen-journeys/gen-contracts/gen-scripts/run-tests 全链路及 validate-code/clean-code 依赖 run-tests）
- [ ] `intent: new-feature` 的 Quick 模式 proposal，`forge task index` 输出的任务列表与基准输出完全相同（含 gen-journeys/run-tests 及后续 validate-code/clean-code/doc-drift 依赖链）
- [ ] `forge task index` 对 `intent: refactor`/`intent: cleanup` 的 feature 不生成测试管道任务（gen-journeys/gen-contracts/gen-scripts/run-tests）
- [ ] `intent: cleanup` 的 proposal 即使存在 PRD 文档，`build.go` 的 `detectMode()` 仍强制返回 Quick 模式，不触发 Breakdown 管道
- [ ] refactor/cleanup 特征的 validate-code/clean-code 任务正确生成且依赖链完整。Breakdown+refactor：business tasks → validate-code → clean-code → consolidate-specs（validate-code 依赖最后一个 business task）；Quick+refactor/cleanup：business tasks → clean-code → doc-drift（clean-code 依赖最后一个 business task）
- [ ] `intent: refactor`/`intent: cleanup` 且 business task 列表为空时（如纯文档类型特征），不生成 validate-code/clean-code 等下游任务，避免悬空的 `depends_on` 引用，仅保留 quality-gate hook 验证
- [ ] 已有 proposal（无 intent 字段）默认走 `new-feature` 管道，行为不变
- [ ] autogen.go 的生成函数（`GetBreakdownTestTasks()`/`GetQuickTestTasks()`）和接线函数（`resolveBreakdownDeps()`/`resolveQuickDeps()`）存在单元测试，覆盖全部 5 种接线场景（new-feature Breakdown/Quick、refactor Breakdown/Quick、cleanup Quick）及零 business task 边界情况，测试全部通过

## Next Steps

- Proceed to `/write-prd` to formalize requirements
