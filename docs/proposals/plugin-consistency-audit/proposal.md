---
created: "2026-05-29"
author: "faner"
status: Approved
intent: "refactor"
---

# Proposal: Forge Plugin 内部一致性审计

## Problem

v3.0.0 经历 test profile system 和 intent-driven pipeline branching 等大规模重构后，21 个 skill、18 个 command、1 个 agent 的 SKILL.md 与其各自的 templates/rules/data 文件之间存在指令矛盾、冗余信息或时序问题，且缺乏系统性验证。

### Evidence

- test profile 系统将 Playwright 硬编码替换为可插拔 profile，影响了 gen-test-scripts、run-tests、init-justfile 等多个 skill 的 rules 和 templates
- intent-driven branching 引入 `new-feature/refactor/cleanup` 三路分支，影响 breakdown-tasks、quick-tasks、run-tasks 等核心编排 skill
- **已确认的不一致实例**: `run-tests/SKILL.md` 已完全迁移到可插拔 test profile 机制（全文无 Playwright 引用），但其 `rules/env-check.md` 第 49 行仍硬编码 `npx playwright install`，与 SKILL.md 的 profile-agnostic 设计直接矛盾
- 208 个 .md 文件（skills 目录 21 个 SKILL.md + 49 个 templates + 75 个 rules + 6 个 data + 5 个 examples + 6 个 types + 26 个其他 md = 188，加上 18 个 command + 1 个 agent + 1 个 hooks/guide.md = 共 208 个可审计文件）通过手动维护交叉引用，重构过程中依赖局部修改而非全局验证

### Urgency

v3.0.0 计划于 2026 年 6 月 15 日前发版（当前处于 RC 阶段，版本号 3.0.0-rc.35）。内部不一致会导致运行时行为异常（流程卡死、模板字段缺失、步骤时序错乱），发版后修复需走 hotfix 流程，预计修复单个 P0 问题的周期（发现 → 定位 → 修复 → 验证）为 2-4 小时，且会中断用户工作流；而审计阶段预防性发现同类问题的边际成本约为 5-10 分钟/组件。审计须在 2 周内（2026 年 6 月 1 日前）完成，为后续修复留出缓冲。

## Proposed Solution

对 forge plugin 所有组件进行**内部逻辑自洽性审计**：逐一检查每个 skill 的 SKILL.md 与其 templates/rules/data/examples/types 之间、每个 command 的内部流程、以及 agent 的指令之间是否存在矛盾、冗余或时序问题。输出结构化问题报告（含文件路径、问题描述、严重等级、修复建议），不做实际修复。

### 审计方法论

每个组件的审计遵循三层协议，AI 在每一层扮演不同角色：

**Layer 1 — 结构完整性（文件存在性校验）**
- 目标: 确认 SKILL.md 中引用的所有文件路径实际存在，且无孤立文件
- AI 角色: 自动提取 SKILL.md 中的路径引用，与文件系统清单交叉比对
- 失败类型: REFERENCE（引用不存在的路径）
- 比对方法: 正则提取 `templates/xxx.md`、`rules/xxx.md` 等路径模式 → `ls` 验证

**Layer 2 — 指令一致性（语义比对）**
- 目标: 确认 SKILL.md 中的指令与 rules/templates/data 文件中的描述不矛盾
- AI 角色: 逐一将 SKILL.md 的每个步骤/约束与对应 rules 文件中的条款对比，识别"必须"vs"可选"等关键词冲突、字段名不一致、条件分支缺失
- 失败类型: CONFLICT（矛盾）、REDUNDANT（冗余）、INCOMPLETE（缺失）
- 比对方法: AI 读 SKILL.md 全文 → 读每个关联文件 → 逐条检查：关键词强度是否匹配（参见下方关键词映射表）、字段名是否一致、步骤时序是否对齐、是否有 SKILL.md 未提及但 rules 中存在的约束
- 冗余检测启发式: REDUNDANT 的判定不同于 CONFLICT——后者是"同一约束给出矛盾描述"，前者是"同一信息在多处重复但无矛盾"。具体检测方法：(1) 提取 SKILL.md 中的每条约束/步骤，记录其语义摘要；(2) 在 rules/templates 中搜索语义等价的描述；(3) 若同一语义在 ≥ 2 个文件中出现且无信息增量（即后文未添加 SKILL.md 中缺少的细节），则标记为 REDUNDANT。注意：rules 中对 SKILL.md 的合理展开（如添加具体示例、边界条件）不算冗余——只有纯重复才计入

**关键词强度映射表**:

| 强度等级 | 中文关键词 | 英文关键词 | 隐含关键词 | 匹配规则 |
|---------|-----------|-----------|-----------|---------|
| 强制 | 必须、务必、一定、不可 | must、required、mandatory、always | 确保、保证、务必（语境强制） | SKILL.md 标"必须"的约束，rules 中不得标为"可选"；反之亦然 |
| 推荐 | 应该、建议、推荐 | should、recommended、prefer | 最好、建议（语境推荐） | SKILL.md 与 rules 中对同一约束的强度应一致（同为推荐或同为强制） |
| 可选 | 可以、可选、视情况 | optional、may、can、if needed | — | 不与强制/推荐冲突即可 |
| 禁止 | 不可、禁止、不要、切勿 | must not、never、avoid、do not | 避免（语境禁止） | SKILL.md 禁止的行为，rules 中不得要求执行 |

匹配规则：(1) 提取 SKILL.md 每条指令的关键词，查表确定强度等级；(2) 在 rules/templates 中找到对应约束，同样提取关键词并确定等级；(3) 比较两端等级是否一致——不一致则报告 CONFLICT。隐含关键词需结合语境判断，若无法确定则标记为"待确认"并在报告中标注低置信度。

**Layer 3 — 时序与流程（逻辑流校验）**
- 目标: 确认步骤顺序合理、前置条件已满足、无循环依赖
- AI 角色: 将步骤流程图式化，检查步骤 N 是否依赖步骤 N+1 的输出
- 失败类型: TIMING（时序错误）
- 比对方法: 提取 SKILL.md 中"先 X 再 Y"的时序约束 → 验证 template 中的字段使用顺序是否匹配

**审计协议步骤**:
1. 枚举组件清单: 列出全部 21 skills + 18 commands + 1 agent + hooks/guide.md
2. 对每个组件，列出其文件清单（SKILL.md + templates/ + rules/ + data/ + examples/ + types/）
3. 执行 Layer 1 结构检查，记录 REFERENCE 类问题
4. 执行 Layer 2 语义比对，记录 CONFLICT/REDUNDANT/INCOMPLETE 类问题
5. 执行 Layer 3 时序检查，记录 TIMING 类问题
6. 为每个问题标注严重等级（P0-P3）并给出修复建议
7. 输出结构化报告

**报告 schema**: 每条问题包含 `{component, file_path, layer, category, severity, description, fix_suggestion}`

**Prompt 策略**: 审计采用**逐组件多轮对话**模式，而非单条巨型 prompt。对每个组件执行以下对话序列：
1. 第一轮：读取 SKILL.md 全文，要求 AI 提取所有步骤、约束、引用路径和字段名，输出结构化摘要
2. 第二轮起（每个关联文件一轮）：读取一个 rules/templates/data 文件，要求 AI 将文件内容与第一轮的 SKILL.md 摘要逐条比对，按关键词映射表检查强度一致性，输出差异列表
3. 汇总轮：将所有差异列表汇总，要求 AI 去重、分级（P0-P3）、生成修复建议

此策略的优势：(1) 每轮对话上下文可控，避免单次塞入 36+ 文件导致截断；(2) 每轮有明确输入输出，便于人工复核中间结果；(3) 若某轮发现异常可针对性追问，而非重跑整个组件。

### Design Rationale

审计按"单一组件自洽"而非"跨组件协调"组织——这不是创新，而是务实的范围裁剪：组件内部的不一致（重构残留）是最高概率的问题来源，且审计边界清晰、可增量执行。跨组件冗余可能是设计层面的合理重复，不在此次审计范围内。价值主张是**彻底性**（100% 覆盖而非抽样），而非方法创新。

## Requirements Analysis

### Key Scenarios

- SKILL.md 描述的步骤流程与 template 中假设的字段/结构不一致
- rules 文件中的约束条件与 SKILL.md 中的指令矛盾（已确认实例: `run-tests/SKILL.md` 使用可插拔 profile，但 `rules/env-check.md` 硬编码 `npx playwright install`）
- SKILL.md 引用的 template/rule/data/examples/types 文件路径不存在或已过时
- 同一 skill 内重复描述同一行为（SKILL.md 和 rules 各说一遍）
- Command 内部流程步骤时序错误（如先读后写、先验证后检查）
- 支持 SKILL.md 未提及的约束存在于 rules 文件中（INCOMPLETE）
- 支持 rules/templates 中存在的文件未被 SKILL.md 引用（孤立文件）
- SKILL.md 与 rules 对同一约束描述一致，但该约束引用了已废弃的功能或与实际 codebase 行为不符（一致但不正确——此类问题的完整验证需对照运行时行为，超出本次审计范围，但 AI 在比对过程中若发现明显的过时引用如已删除的 CLI 参数，应作为 INCOMPLETE 标记并在报告中注明"疑似一致但不正确，需运行时验证"）

### Audit Process Failure Scenarios

- AI 产生误报（hallucination）：报告中的问题在实际文件中不存在或描述不准确
- AI 遗漏真实问题（false negative）：运行时出现的问题未被审计报告捕获
- 审计结果不可复现：同一 commit 基准上再次运行产出不同的问题集合
- 大组件审计截断：因上下文窗口限制，仅检查了部分文件导致遗漏

### Non-Functional Requirements

- 审计覆盖率: 100% 的 skill（21个）、command（18个）、agent（1个）
- 问题分类: 矛盾(CONFLICT)、冗余(REDUNDANT)、时序(TIMING)、引用(REFERENCE)、缺失(INCOMPLETE)

#### Severity Level Definitions

- **P0 (Critical)**: 会导致运行时错误或流程完全卡死（如引用的模板文件不存在、步骤时序颠倒导致必选字段缺失）
- **P1 (High)**: 会导致行为偏差但不会完全阻断（如 SKILL.md 与 rules 对同一约束给出"必须"和"可选"的矛盾描述）
- **P2 (Medium)**: 信息冗余或轻微不一致，不影响运行时行为但增加维护负担（如同一 skill 内 SKILL.md 和 rules 重复描述同一行为）
- **P3 (Low)**: 风格或措辞不一致，不影响功能（如命名风格不统一、注释过时）

### Constraints & Dependencies

- 审计基于当前 v3.0.0 分支代码的特定 commit hash（报告标注基准），不依赖运行时测试
- 不修改任何代码，仅输出报告
- 时间约束: 1 个工作日内完成（见 Feasibility 预估）
- 工具约束: 使用 AI（Claude）+ 脚本（Layer 1），所有文件均为 markdown 可完整读取。推荐模型版本: Claude Sonnet 4 (claude-sonnet-4-20250514) 或同等能力模型；推荐参数: temperature=0（最低非确定性）；报告须记录实际使用的模型版本和参数
- 大组件（文件数 > 20）需分批审计以适配 AI 上下文窗口。分批策略：按子目录分组（如 templates/、rules/、data/ 各为一批），每批独立完成 Layer 1-3。分批完成后执行一次**汇总轮**：将各批的 SKILL.md 摘要与差异列表合并，进行跨组比对——检查 rules/ 中的约束是否与 templates/ 中的字段名/结构一致、data/ 中的枚举值是否与 rules/ 中的条件分支对齐。此汇总轮以 SKILL.md 为锚点，避免分批引入的盲区。

## Alternatives & Industry Benchmarking

### Industry Solutions

大型 prompt-based system 的一致性验证有成熟工具链：

- **promptfoo** (github.com/promptfoo/promptfoo, v0.110+ 截至 2026 年 5 月): 对 LLM prompt 进行自动化测试和回归检测，通过 assertion 框架验证 prompt 输出一致性。适用于 prompt 输出质量验证，但不检查 prompt 文件间的内部结构一致性。
- **guardrails-ai** (github.com/guardrails-ai/guardrails, v0.5+ 截至 2026 年 5 月): 通过 schema (Pydantic) 定义 LLM 输出的结构和验证规则，实现输出层面的结构化约束。适用于运行时 guard，不适用于静态 markdown 审计。
- **LangSmith** (smith.langchain.com, 截至 2026 年 5 月): 提供 prompt 版本管理和 trace 可视化，支持 prompt 模板的 diff 比对。最接近我们的需求，但它是 SaaS 平台，且关注 prompt 版本管理而非组件内部自洽性。
- **markdownlint + 自定义规则** (markdownlint v0.16+, 截至 2026 年 5 月): 可检查 markdown 格式一致性（标题层级、链接有效性），但无法检测语义层面的指令矛盾。

这些工具解决的是 prompt 输出质量、运行时 guard 或格式规范，而 forge plugin 需要的是**自由格式 markdown 文件间的语义一致性校验**——一个尚未有成熟工具覆盖的空白地带。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| 人工逐文件审读 | 传统代码审查方法 | 最细致，能捕获隐含语义问题 | 208 个可审计文件量巨大（预计 15-20 人时），易疲劳遗漏，不可复现 | Rejected: 效率太低，时间窗口不够 |
| 自动化 schema 验证 | promptfoo/guardrails 的思路 | 可重复执行，可 CI 集成 | 需先为每个 skill 定义 JSON Schema（预计 5-8 人时），后续维护 schema 成本持续存在 | Rejected: ROI 不够，schema 定义本身的正确性也需要验证 |
| markdownlint + 自定义规则 | markdownlint 社区 | 零成本 CI 集成，可检测链接失效和格式不一致 | 纯语法层面，无法检测"必须 vs 可选"等语义矛盾；自定义规则编写复杂度接近定义 schema | Rejected: 覆盖面太窄，只能做 Layer 1 |
| **AI 辅助分层审计** | 本次方案 | 覆盖全面（三层协议）、可理解上下文语义、无需预先定义 schema | AI 可能产生误报（hallucination）；非确定性——两次运行可能产出不同结果；大组件（如 eval 有 36+ 文件）可能超出 AI 上下文窗口；无法保证语义等价检测的召回率 | **Selected: 最适合当前规模（208 文件）和时间约束（RC 阶段）** |

**成本对比摘要**: schema 验证需 5-8 人时定义 JSON Schema（不含后续维护成本），实际执行约 1-2 人时；AI 辅助审计前期 0 人时定义，执行阶段 4-6 小时 + 约 $10-20 token 成本（200K-400K input tokens）。两者总耗时相近，但 AI 方案无需维护 schema 且可检测语义矛盾（schema 只能做结构校验），在"一次性审计"场景下 ROI 更优。

## Feasibility Assessment

### Technical Feasibility

所有文件均为 markdown，可完整读取和分析。Layer 1（结构检查）可通过脚本自动化，无技术障碍。Layer 2-3（语义比对和时序检查）依赖 AI 的语义理解能力，存在两个技术局限：(1) AI 无法保证 100% 召回率——隐含的逻辑矛盾（如"仅在 X 条件下执行"vs rules 中隐含假设 X 永远为真）可能被遗漏；(2) 大型组件（如 eval skill 含 36+ 文件）可能超出单次上下文窗口，需要分批处理。这两个局限是已知的精度上限，不构成阻断性障碍，但需要在报告中标注审计置信度。

### Resource & Timeline

- 预计总工作量: 4-6 小时（AI 辅助执行 + 人工复核抽样）
- 其中 Layer 1 可脚本化（约 30 分钟编写 + 2 分钟运行），Layer 2-3 需 AI 逐组件执行（平均 10-15 分钟/组件，21 个 skill ≈ 3-4 小时，18 个 command + 1 agent + hooks ≈ 1 小时）
- 预计 AI token 消耗: 约 200K-400K input tokens（208 个 md 文件，平均 2-3K tokens/文件，含多轮比对）
- 墙钟时间: 1 个工作日内完成
- 锚定 commit hash 以避免审计期间文件变更（参见 B3 风险）

### Dependency Readiness

所有文件均在当前仓库中，无需外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 重构后的文件已经内部一致 | Assumption Flip: 假设不一致，逐一验证 | Refined: 重构是局部修改，残留不一致很可能存在（已通过 run-tests/env-check.md 确认至少 1 例） |
| 跨 skill 冗余是问题 | Occam's Razor: 跨 skill 重复可能是设计意图 | Confirmed: 用户确认忽略跨 skill 冗余，聚焦组件内部自洽 |
| 审计产出仅是"检查行为" | 价值审视: 审计的价值不仅是发现现有问题 | Refined: 审计的产出是 (1) 优先级排序的问题清单，(2) 严重性分级的热力图，(3) 面向未来维护的一致性基线（后续重构可 diff 此报告判断退化）。有效性通过已知问题的复现验证（见 SC）|

## Scope

### 跨组件引用边界规则

当 SKILL.md 引用另一个 skill 目录下的文件时（如共享 template 或 rule），按以下规则处理：
- **在审计范围内**: 验证该引用路径是否存在且文件可读（REFERENCE 类检查）
- **在审计范围内**: 验证 SKILL.md 中对该文件用途的描述与文件实际内容不矛盾（CONFLICT 类检查）
- **不在审计范围内**: 该被引用文件与其所属 skill 的 SKILL.md 之间的一致性（由被引用 skill 的审计覆盖）
- 简言之：审计检查的是"SKILL.md 视角下的引用正确性"，而非"被引用文件的内部一致性"

### In Scope

- 21 个 skill 的 SKILL.md 与其各自的 templates/rules/data/examples/types 之间的逻辑自洽性（具体定义：关键词强度一致、字段名匹配、步骤时序对齐、引用路径有效、无未提及的约束——参见审计方法论 Layer 1-3）
- 18 个 command 的内部流程一致性
- 1 个 agent (task-executor) 的内部指令一致性
- hooks/guide.md 的内部一致性（即 guide.md 中描述的 hook 行为与其引用的 hook 脚本文件路径和参数是否匹配，内部步骤之间是否存在矛盾）。**例外说明**: guide.md 是"单一组件自洽"原则的明确例外——它本身是跨 hook 脚本的索引文档，其职责正是描述多个 hook 脚本的行为和参数。审计 scope 限于：(1) guide.md 中引用的脚本路径是否存在（REFERENCE）；(2) guide.md 对脚本参数的描述是否与脚本文件中的参数声明一致（CONFLICT）；(3) guide.md 内部步骤之间是否矛盾。不深入验证脚本文件的逻辑正确性（那属于代码审查范畴）
- 问题分类: 矛盾(CONFLICT)、冗余(REDUNDANT)、时序(TIMING)、引用(REFERENCE)、缺失(INCOMPLETE)

### Out of Scope

- 跨 skill 之间的冗余内容（设计层面的合理重复）
- rubrics/experts 的 prompt engineering 质量审查（注：eval skill 内部文件路径的交叉引用校验仍在审计范围内）
- Forge CLI Go 源码
- 用户项目目录结构
- 实际代码修复（仅产出报告）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 语义层面的隐含矛盾难以通过文本对比发现 | M（AI 可识别显式矛盾，但隐含假设需要推理链） | M（遗漏问题将在使用中暴露） | 重点检查关键词不一致（如"必须"vs"可选"）、条件分支覆盖 |
| AI 误报——生成不存在的问题（false positive） | H（基于通用 LLM 文献估算，结构化比对任务的典型误报率约 10-20%，非实测数据） | M（浪费开发者时间验证假问题） | 人工抽样复核：随机抽取 20% 的 P0/P1 问题独立验证真实性；报告标注置信度 |
| 报告问题过多导致修复优先级不清 | L（P0-P3 分级足以区分） | M | 每个问题标注严重等级（P0-P3），P0 问题附修复优先时间线 |
| 审计过程中遗漏某些文件 | L（Layer 1 脚本可自动枚举） | L | 使用 `find` 脚本自动生成文件清单，逐一勾选 |
| 审计期间文件变更导致报告过时 | M（v3.0.0 活跃开发中） | M（报告与实际状态不一致） | 锁定 commit hash，报告标注基准 commit；若审计期间有重大 merge 则重跑 |
| AI 非确定性——两次运行产出不同结果 | M（LLM temperature > 0 时必然存在） | L（问题集合差异通常 < 5%） | 记录 AI model 版本和参数；核心结论取交集（两次运行均报告的问题优先处理）。注：4-6 小时为单次运行基准；若需双运行验证，额外开销约 50-80%（非 100%，因 Layer 1 脚本化结果可缓存复用），总计约 6-11 小时，仍可在 2 个工作日内完成 |
| 审计疲劳——连续审查 21 个 skill 后注意力下降 | M（质量保证领域的已知现象：inspection fatigue，4-6 小时持续审查后遗漏率上升） | M（后期组件的问题捕获率降低） | 随机化审计顺序（非按字母序），避免系统性偏差；每个 skill 审计后休息间隔；P0 问题在最终汇总时做跨组件二次检查 |

## Success Criteria

- [ ] 21 个 skill 100% 覆盖审计，每个 skill 的 SKILL.md 与其 templates/rules/data/examples/types 逐一对比
- [ ] 18 个 command 100% 覆盖审计
- [ ] 1 个 agent (task-executor) 完成审计
- [ ] hooks/guide.md 完成审计（作为单一组件原则的例外，验证 guide.md 中引用的脚本路径存在性、参数描述与脚本声明的一致性、内部步骤无矛盾——不深入验证脚本逻辑）
- [ ] 输出结构化问题报告，每个问题包含: 文件路径、问题描述、严重等级(P0-P3)、修复建议、置信度(high/medium/low)
- [ ] 问题按 CONFLICT/REDUNDANT/TIMING/REFERENCE/INCOMPLETE 五类分类，且每类至少有 1 个实例（用于验证分类标准可操作）；若某类为 0 则需在报告中说明为何该类问题不存在。具体要求：对 TIMING 类为 0 的情况，报告须列出所有含多步骤流程的组件清单并确认每个组件的步骤排序已验证一致（而非泛泛声明"未观察到时序问题"）
- [ ] **有效性验证**: 对已知存在问题的 run-tests skill（`rules/env-check.md` 中残留 Playwright 硬编码）进行重审计，确认报告能复现该 P1 级矛盾——作为审计有效性的基线验证
- [ ] **误报率抽检**: 随机抽取 ≥ 20% 的 P0/P1 问题进行人工独立验证，确认 ≥ 80% 为真实问题（非 AI 幻觉）。若真阳性率 < 80%，触发升级流程：(1) 将全部 P0/P1 问题提交人工复核而非抽样；(2) 扩大抽样至所有类别的 50%；(3) 在报告中标注"审计置信度不足"并附人工复核结果
- [ ] 报告标注基准 commit hash，确保可复现
