---
created: 2026-05-31
author: faner
status: Draft
intent: refactor
---

# Proposal: Intent Enriched Enum

## Problem

Proposal intent 枚举只有 3 个值（`new-feature`、`refactor`、`cleanup`），与 task type 的 8 个值映射断裂，导致下游 pipeline 分支过粗——`refactor`、`cleanup`、`fix` 被等同对待，不同场景走了不合适的流程。具体而言：write-prd 和 tech-design 只有 2 条分支（`new-feature` → Full PRD；其余 → Spec-only PRD），8 个 task type 中有 5 个被压缩进同一条分支。

### Evidence

- `coding.fix`、`coding.enhancement`、`doc` 作为 task type 存在，但 proposal intent 没有对应值
- `fix` 在 brainstorm 中被启发式拆分为 `new-feature` 或 `refactor`，没有独立身份
- `refactor` 和 `cleanup` 在 write-prd/tech-design 中被完全等同对待（spec-only PRD），但两者的变更范围和风险等级可能很不同
- 一个改变外部 API 的 refactor 仍跳过 API handbook，因为 pipeline 分支只看 intent 不看内容

### Urgency

当前 Forge 已有 8 个 task type 中 5 个（`coding.fix`、`coding.enhancement`、`doc`、`doc.consolidate`、`doc.drift`）无对应 intent 值。brainstorm/SKILL.md 中 `coding.fix` 的启发式推断逻辑（line 96）每次执行都需要 LLM 做二元判断（是否有新用户可见行为），这是唯一需要运行时推断的 type→intent 映射。

**关于频率数据的说明**：本提案缺乏量化的频率数据（如"最近 N 次 brainstorm 中 heuristic miss 的次数"），因为 Forge 未内置 pipeline 执行日志，无法统计实际错配频率。但严重性论证如下：heuristic miss 的后果不是"稍不理想"而是"产出错误的产物集合"——bug fix 被当作 new-feature 会生成不必要的 User Stories 和 Full PRD（浪费 token + 用户困惑），涉及 API 变更的 refactor 跳过 API handbook 会遗漏关键审查步骤（质量风险）。即使错配频率不高，每次错配的修复成本（手动纠正已生成的 PRD/tech-design）和遗漏成本（未执行的检查步骤）都显著高于一次性修复映射的成本。

## Proposed Solution

1. **扩充 intent 为 6 值枚举**：`new-feature`、`enhancement`、`refactor`、`cleanup`、`fix`、`doc`，覆盖 8 个 task type 中的 6 个。映射关系：`new-feature` → `coding.feature`、`enhancement` → `coding.enhancement`、`refactor` → `coding.refactor`、`cleanup` → `coding.cleanup`、`fix` → `coding.fix`、`doc` → `doc`。`doc.consolidate` 和 `doc.drift` 属于低频内部任务，由 skill 自动生成，无需用户可见 intent，统一归入 `doc` umbrella（breakdown-tasks 的 Intent Propagation 将 `doc` intent 解析为 `doc` task type，不区分子类型）。
2. **混合模式 pipeline 分支**：intent 控制默认 pipeline 配置（一张表），PRD 内容中的明确信号可以覆盖默认值。Override Signals 检测条件如下：

   | 信号类型 | 关键词/模式 | 覆盖动作 |
   |---------|------------|---------|
   | API 变更 | "API"、"endpoint"、"命令重命名"、"接口变更"、"breaking change" | 开启 API handbook |
   | 用户可见行为 | "用户可见"、"UI 变更"、"CLI 输出"、"新选项" | 开启 User Stories |
   | 安全相关 | "认证"、"授权"、"权限"、"加密"、"token" | 开启 Security Review |
   | 性能相关 | "性能"、"延迟"、"吞吐量"、"缓存" | 开启 Performance Baseline |
   | 数据迁移 | "迁移"、"schema 变更"、"数据格式" | 开启 Migration Plan |

   检测机制：LLM 在生成 PRD/tech-design 内容的过程中，同步（而非先后）完成信号检测——内容生成与信号匹配是同一次 LLM 调用中的并行推理，不存在"先生成再扫描"的时序关系。命中任意一个信号即触发对应覆盖。多信号同时命中时各覆盖动作独立叠加（如同时出现"API"和"性能"则同时开启 API handbook 和 Performance Baseline）。Negation handling：信号检测仅匹配肯定性陈述语境（如"涉及 API 变更"），对于否定语境（如"不涉及 API 变更"、"不改变接口"），LLM 应跳过该信号。这依赖 LLM 的上下文理解能力而非纯关键词匹配——这是合理的，因为 Pipeline Configuration 步骤本身就是 LLM 执行的。当 override 触发时，write-prd/tech-design 在生成的文档中添加一个注释行标注被触发的覆盖信号（如 `<!-- Override: API handbook enabled by signal "接口变更" -->`），供用户 review 时确认。
3. **简化 brainstorm 推断**：`fix` 始终为 `fix`，移除启发式判断；每个 task type 直接映射到对应 intent
4. **Pipeline Configuration 表**（统一 write-prd 和 tech-design）：

   | Intent | PRD Format | User Stories | API Handbook | Test Pipeline | Security Review |
   |--------|-----------|-------------|-------------|--------------|----------------|
   | new-feature | Full | Yes | Yes | Yes | If signal |
   | enhancement | Simplified (保留 Background/Goals/Test Pipeline，跳过 User Stories) | No | If signal | Yes | If signal |
   | refactor | Spec-only | No | If signal | Yes | If signal |
   | cleanup | Spec-only | No | No | Yes | No |
   | fix | Spec-only | No | If signal | Yes (reproduce → fix → verify) | No |
   | doc | Minimal（标题：一句话描述文档变更对象和目的；目标：列出要更新/新增的文档文件和预期变更点；scope：界定变更涉及的文档范围和不涉及的边界） | No | No | No | No |

   enhancement 的简化 PRD：保留 Background（说明增强什么）、Goals（增强目标）、Test Pipeline（确保增强有测试覆盖），跳过 User Stories（已有用户群体，不涉及新用户流程）。

### Innovation Highlights

无特别创新。对标 task type 的现有分类体系，消除 intent 与 type 之间的映射鸿沟。

## Requirements Analysis

### Key Scenarios

1. **Bug fix 提案**：brainstorm 直接推断 `fix`，pipeline 默认 spec-only，跳过 user stories 和 API handbook
2. **Enhancement 提案**：brainstorm 推断 `enhancement`，pipeline 默认跳过 user stories 但保留 test pipeline（改善现有行为需要测试覆盖）
3. **改变外部 API 的 refactor**：brainstorm 推断 `refactor`（默认跳过 API handbook），但 PRD 内容包含"CLI 命令重命名"信号 → 覆盖开启 API handbook
4. **纯文档提案**：brainstorm 推断 `doc`，pipeline 全部跳过，直接进入 task 生成
5. **混合内容提案**：brainstorm 按核心目标推断主 intent，个别 task 通过 per-task type override 覆盖
6. **多信号触发**：PRD 同时提及"API 变更"和"性能"关键词 → 两个覆盖信号独立叠加，同时开启 API handbook 和 Performance Baseline。覆盖信号之间无互斥关系，全部采用"只加不减"策略
7. **Intent-content 不匹配**：intent 为 `doc` 但 PRD 意外包含"API"关键词。由于 intent 基线已决定 pipeline 路径（doc → Minimal PRD），override 信号在 Minimal PRD 格式下无产物可覆盖（doc intent 的 pipeline 没有可被"开启"的检查项）。此时 override 信号为 no-op，不改变 pipeline 行为。**设计约束**：这是显式的设计决策而非偶然属性——doc pipeline 的 Minimal 格式刻意不包含任何可被 override 开启的检查项，保证 doc intent 始终走最轻路径。如果未来需要为 doc pipeline 添加可被 override 的检查项，需重新评估此约束
8. **Invalid intent fallback**：brainstorm 输出不在 6 值枚举中的 intent（如 `bug-fix`、`documentation`、`hotfix`）。AskUserQuestion 中提供的是结构化选项列表，LLM 应从列表中选择。如果 LLM 仍输出非法值，下游 skill 的 Pipeline Configuration 表不匹配该 intent → LLM 将回退到最接近的合法值（如 `bug-fix` → `fix`，`documentation` → `doc`）。这是 best-effort 容错，不保证正确性——根本解决方案依赖 brainstorm SKILL.md 中对输出格式的明确约束

### Architecture Decision: fix 的映射策略

**问题**：`coding.fix` 在当前 breakdown-tasks/quick-tasks 中被标注为"Auto-generated for test failures via forge task add; do not assign manually"。引入 `fix` intent 后，用户可在 brainstorm 中选择 fix，breakdown-tasks 需要将其映射到 `coding.fix`，这与"do not assign manually"规则冲突。

**决策**：允许 fix intent 映射到 `coding.fix`，更新 Type Assignment 表中的约束从"do not assign manually"改为"可由 fix intent 自动映射，但不可通过 `forge task add` CLI 手动创建"。理由：
1. fix intent 是用户在 brainstorm 中对 bug fix 场景的明确声明，与 CLI 手动创建有本质区别——前者经过意图确认，后者绕过 pipeline
2. 测试失败自动生成的 coding.fix 仍然保留，不受影响
3. 这样 fix intent 获得独立身份，不再被启发式拆分为 refactor/new-feature

**决策**：`enhancement` intent 映射到 `coding.enhancement` task type（不再是 `coding.feature`）。当前 brainstorm 将 `coding.feature` 和 `coding.enhancement` 都映射到 `new-feature` intent，引入 `enhancement` 后两者分道：`new-feature` → `coding.feature`，`enhancement` → `coding.enhancement`。这恢复了 Type Assignment 表中两个 task type 原有的语义区分（"adds new runtime behavior" vs "improves existing behavior without adding new capabilities"）。

### Non-Functional Requirements

- **Pipeline 向后兼容**：现有的 `new-feature`、`refactor`、`cleanup` 值的 pipeline 产物不变
- **交互向前演进**：brainstorm AskUserQuestion 从 3 选项扩展为 6 选项，是用户可见的交互变更（非破坏性——旧 3 值仍在选项中）
- 一致性：intent-to-type 映射覆盖 6/8 task type（`doc.consolidate`、`doc.drift` 为 skill 自动生成的低频类型，归入 `doc` umbrella），消除歧义
- 可维护性：Override Signals 关键词表和 Pipeline Configuration 表的变更有明确流程——修改 SKILL.md 中的结构化表格即可，无需改 Go 代码。write-prd 和 tech-design 各持有 Pipeline Configuration 表的一份副本，变更时需同步更新两处（已知风险，见 Risk Assessment）
- 可测试性：每个 override signal 需至少一个 PRD 输入→pipeline 输出的测试用例（如 PRD 包含"CLI 命令重命名" → 输出包含 API handbook section）。见 Success Criteria 中的具体验证条件

### Constraints & Dependencies

- Intent 分支逻辑全部在 skill markdown 中，无 Go 代码依赖
- 变更限于 plugins/forge/ 目录下的 8 个文件（grep `intent` 匹配的 skill 文件，排除 task-doc.md 的误匹配——该文件中 "intentionally" 包含 "intent" 子串但无 intent 字段）

## Alternatives & Industry Benchmarking

### Industry Solutions

分类系统的精确度随场景增长自然演进是常见模式。GitHub Issues 的 label 体系从最初的 `bug`/`enhancement` 二分法（GitHub 2013 年默认 label set）逐步演进——2016 年引入 `good first issue`（源于开源社区贡献者引导需求）、2020 年 GitHub 官方增加 `breaking change` 标签以支持 Semantic Versioning pipeline 自动化。TypeScript 在 2.0 版本（2016）引入 `diagnosticCategories` 枚举，将早期笼统的 error/warning 拆分为 `Suggestion`、`Message`、`Error`、`Warning` 等精确类别，每个 category 对应不同的编译器处理流程（Suggestion 不阻断编译，Error 阻断）。两者共同模式是：当使用者需要对不同类别施加不同行为（不同通知策略/不同编译流程）时，分类粒度必须匹配行为差异。Forge 的 intent→pipeline 分支是同类问题：不同 intent 需要不同的产物集合，3 值枚举无法区分需要不同处理的场景。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动 | fix/enhancement/doc 无对应 intent，pipeline 过粗 | Rejected: 覆盖缺口随场景增长扩大 |
| 只扩枚举 | — | 最小改动，解决 brainstorm 推断问题 | pipeline 分支仍只有 2 条，5 个 intent 挤在 Spec-only 分支里 | Rejected: 枚举对齐了但 pipeline 没跟上，type→pipeline 映射仍然过粗 |
| **扩枚举 + 混合 pipeline** | CI lint gate 模式（见下） | 完整解决两个动机 | 8 个文件变更；pipeline table 在 write-prd 和 tech-design 各存一份（见 Risk） | **Selected: 双重改进** |
| 完全内容驱动 pipeline | — | 最精准 | 无基线：每次 pipeline 配置完全依赖 LLM 对 PRD 内容的解读，同一 PRD 在不同 session 可能产出不同配置，缺少可复现的默认行为 | Rejected: 被拒绝的核心原因不是"LLM 判断不可靠"（本方案的 override 同样依赖 LLM），而是"没有基线"——用户无法预期默认产物集合，也无法判断 override 是否合理。本方案通过 intent 基线解决了这个问题 |

**CI lint gate 模式**：CI 系统中的 lint gate（如 ESLint 的 `overrides` 配置、GitHub Actions 的 path-based trigger）采用"基线规则 + 条件覆盖"模式——默认对所有文件应用一组规则，然后通过 glob pattern 或条件表达式为特定路径/场景覆盖规则。本方案的 intent = 基线规则（对应 Pipeline Configuration 表的默认列），PRD 内容信号 = 条件覆盖（对应 Override Signals 表）。两者的共同优势是：基线保证最低行为一致，覆盖只做加法（开启额外检查），不做减法。

**与 CI lint gate 的关键差异**：ESLint overrides 使用机器可解析的 glob pattern（确定性、可在 CI 中自动化测试、版本可控），而本方案的 override 信号依赖 LLM 对自然语言的解读（概率性、无法在 CI 中自动化测试、解读结果不受版本控制）。这一差异是可接受的，原因在于：(1) Forge 的整个 skill 执行链本身依赖 LLM 指令遵循，override 信号的 LLM 解读不是引入新风险而是复用已有能力；(2) override 只做加法，最坏情况是多生成不必要的产物（可被用户 review 时发现并丢弃），不会遗漏必要步骤；(3) 结构化条件表降低了 LLM 解读的歧义空间。

## Feasibility Assessment

### Technical Feasibility

完全可行。所有变更是 markdown 编辑——更新推断表、分支表、覆盖规则。无 Go 代码变更。brainstorm 的 AskUserQuestion 从 3 选项变为 6 选项是用户可见的交互变化，但对已有提案的行为无影响（旧提案已生成，不回溯）。

### Resource & Timeline

中型变更：8 个 markdown 文件。write-prd 和 tech-design 变更量最大（需要重写 pipeline 分支逻辑），其余文件是小幅更新。预计 2-3 个任务可完成。

### Dependency Readiness

无外部依赖。所有 skill 文件已存在且结构清晰。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "3 值够用" | 5 Whys | Overturned: task type 有 8 个值，3 值 intent 无法覆盖。缺口通过启发式弥补，增加了不一致风险 |
| "fix 需要启发式区分" | Occam's Razor | Overturned: fix 作为独立 intent 更简单。是否有新用户可见行为由 pipeline 覆盖规则处理，不需要在 intent 层面区分 |
| "pipeline 分支只能靠 intent" | Assumption Flip | Refined: intent 提供稳定基线，PRD 内容提供覆盖信号。两者结合比任何单一维度都可靠 |
| "refactor 和 cleanup 应该等同对待" | Stress Test | Refined: refactor 和 cleanup 的 Pipeline Configuration 默认产物确实相同（Spec-only PRD），但两者的 override 概率不同——refactor 更可能触发 API/performance 信号（涉及结构重组），cleanup 几乎不会。因此分开两个 intent 的价值不在默认 pipeline 差异，而在于：(1) intent-to-type 语义对齐（`refactor` → `coding.refactor` vs `cleanup` → `coding.cleanup`）；(2) override 信号对不同 intent 的实际触发概率不同，区分 intent 让后续度量成为可能 |

## Scope

### In Scope

- **brainstorm/SKILL.md**：更新 Step 4.5 intent mapping 表（6 值），移除 fix 启发式，更新 AskUserQuestion 选项
- **brainstorm/templates/proposal.md**：更新 intent 有效值注释
- **write-prd/SKILL.md**：将二元分支（new-feature vs refactor/cleanup）替换为 Pipeline Configuration 表 + Override Signals；实现 `<!-- Override: ... -->` 注释行的生成逻辑
- **write-prd/rules/self-check.md**：更新 intent-gated 检查为 6 值
- **tech-design/SKILL.md**：将二元分支替换为 Pipeline Configuration 表 + Override Signals；实现 `<!-- Override: ... -->` 注释行的生成逻辑
- **tech-design/rules/design-quality-checks.md**：更新 intent-gated 检查为 6 值
- **breakdown-tasks/SKILL.md**：更新 Intent Propagation 为严格 1:1 映射（6 值），更新 Type Assignment 表中 `coding.fix` 的约束描述
- **quick-tasks/SKILL.md**：更新 Intent Propagation 为严格 1:1 映射（6 值）

### Out of Scope

- Go CLI 代码变更（CLI 不引用 intent 字段）
- 新 skill 或 command 创建
- 已有提案的迁移（旧 3 值行为不变）
- task-sizing-gate 提案（独立提案）
- eval 系列 skill（eval-prd、eval-design 等）的 rubric/contract 更新——这些 skill 评估的是产物质量而非 pipeline 行为，intent 枚举变更不影响其评估逻辑。如果未来 eval skill 需要引用 intent 值，可在后续迭代中处理

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 6 值枚举仍不够（doc.consolidate/doc.drift 已被排除） | M | L | 当前 6 值覆盖所有用户可触发的 task type。doc.consolidate/doc.drift 由 skill 自动生成不需要用户 intent。如未来出现新 task type，可追加 intent 值——Pipeline Configuration 表新增一行即可 |
| 混合模式的覆盖规则被 LLM 忽略 | M | M | 覆盖信号以结构化条件表形式定义（信号类型 → 关键词 → 覆盖动作），不是 prose 描述。LLM 对结构化表格的遵守度显著高于自然语言指令——表格每行是原子化的 if-then 规则，不存在解读歧义。且覆盖规则只"开启"默认关闭的产物，不存在"错误关闭"风险，最坏情况是多做了不必要的产物而非遗漏 |
| 关键词误触发（否定语境） | L | L | 依赖 LLM 上下文理解而非纯字符串匹配。Pipeline Configuration 步骤由 LLM 执行，LLM 可识别"不涉及 API 变更"中的否定——同一 LLM 在处理 Pipeline Configuration 时已具备足够的上下文推理能力（否则整个 pipeline 配置逻辑都不可靠）。最坏情况是开启了一个不必要的检查（多做了产物），不影响正确性。Override 触发时生成的注释行（`<!-- Override: ... -->`）让用户可 review 并纠正 |
| write-prd/tech-design Pipeline Configuration 表同步漂移 | M | M | 两处各持有一份表格，变更时需手动同步。Mitigation：scope 中两处变更是同一批次的并行修改，且两处表格结构完全相同（6 行 × 6 列 + 同一 Override Signals 表），diff 检查成本极低 |
| write-prd/tech-design 分支重写引入不一致 | M | M | Pipeline Configuration 表统一两处逻辑，减少不一致可能性。验证策略：对旧 3 个 intent（new-feature、refactor、cleanup）分别运行变更前的 write-prd/tech-design 测试用例，对比变更后输出中的 pipeline 相关产物是否与变更前一致 |
| 旧提案与新规则不兼容 | L | L | 旧 3 值在表中仍有对应行，行为不变 |

## Success Criteria

- [ ] brainstorm 推断结果为 6 值之一，用户可在 AskUserQuestion 中选择全部 6 个值
- [ ] `fix` 始终推断为 `fix`，不再使用启发式
- [ ] write-prd 和 tech-design 使用统一的 Pipeline Configuration 表（6 行对应 6 个 intent，产生 4 种功能不同的 pipeline 配置：Full、Simplified、Spec-only、Minimal）
- [ ] enhancement intent 生成 Simplified PRD 格式（Background + Goals + Test Pipeline），跳过 User Stories
- [ ] Override Signals 规则存在且可被 PRD 内容触发；触发时生成的 PRD/tech-design 文档包含 `<!-- Override: ... -->` 注释。具体验证：对每个 override signal 各提供一个包含对应关键词的 PRD 输入（如"CLI 命令重命名"触发 API handbook），验证输出包含对应产物
- [ ] breakdown-tasks 和 quick-tasks 的 Intent Propagation 为严格 1:1 映射
- [ ] breakdown-tasks 将 `fix` intent 映射到 `coding.fix` task type（验证 Type Assignment 表更新）；brainstorm 将 `coding.feature` → `new-feature` 和 `coding.enhancement` → `enhancement` 作为独立路径（验证 intent mapping 表 split）
- [ ] 现有 `new-feature`、`refactor`、`cleanup` 值的 pipeline 产物不变——对旧 3 个 intent 分别执行变更后的 write-prd 和 tech-design，验证生成的产物集合（PRD sections、checklist items）与变更前一致
- [ ] `grep -rn "intent" plugins/forge/skills/{brainstorm,write-prd,tech-design,breakdown-tasks,quick-tasks}/ --include="*.md"` 输出中所有 intent 引用均反映 6 值枚举，无残留的 3 值引用

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
