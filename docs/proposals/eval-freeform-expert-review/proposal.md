---
created: 2026-05-23
author: faner
status: Draft
---

# Proposal: eval-proposal 自由专家评审前置阶段

## Problem

eval-proposal 的评审体验过于机械化——固定 rubric 维度强制评分、三阶段流水线协议（推理审计→打分→盲点搜索），导致评审产出像一份检查清单而非专家洞察。rubric 覆盖了已知失败模式，但无法发现文档特有的、超出预设维度的问题。具体而言：rubric scorer 的 prompt 缺少领域特有上下文，导致 scorer 在已知维度上评分准确、在文档特有风险上视而不见。本提案不改变 rubric 本身（维度和评分标准保持不变），而是通过前置自由评审阶段，为 rubric scorer 补充其缺失的领域上下文，使 scorer 在评分时能"看到"rubric 维度之外的问题。注入机制的工作方式：提取的 key findings 以显式 attack points 列表形式追加到 scorer prompt 末尾，scorer 被要求在评分时回应这些 attack points。需明确界定：注入扩展的是 scorer 在固定 rubric 维度内的注意力范围（例如，在「隐藏成本」维度下，scorer 原本不会关注分布式一致性的脑裂成本，但注入后发现该维度后可以将其纳入考量），而非扩展评分维度本身。对于确实无法映射到任何 rubric 维度的发现（如纯架构层面的新颖风险），scorer 会将其记录在 blindspot 搜索阶段产出的「rubric 外发现」区域（以 `[beyond-rubric]` 标签标注，附在 eval report 的 ATTACKS 列表末尾，与 `[blindspot]` 标签并列），而非强行归入某个维度打分。

### Evidence

- 现有 eval-proposal 使用 10 维度 / 1000 分 rubric，CTO 专家角色固定，评审流程完全模板化。在 Forge 内部 dogfooding 中，对 5 份不同领域的提案进行 eval-proposal，10 个 rubric 维度的评审产出在关注点上重叠度超过 70%（均聚焦于「隐藏成本」「回滚计划」等通用维度），而每份提案的独特风险（如：分布式一致性提案中的脑裂场景未被覆盖、UI 设计提案中的无障碍合规盲点未被识别）均未被 rubric 自身发现
- rubric 维度的设计基于通用失败模式（如「隐藏成本」「回滚计划缺失」），无法针对特定提案的独特风险点进行深度挖掘。具体案例：一份关于 test profile 插件化的提案，其核心风险是「profile 注册时机对测试发现的影响」，但 rubric 的 10 个维度中没有任何一个覆盖此领域特定问题

### Urgency

eval 是 Forge 质量保障体系的核心。当前 rubric 评审缺少领域特有上下文，scorer 无法"看到"rubric 维度之外的问题，意味着每一份经过 eval 的提案仍可能携带未识别的盲点进入实施阶段。本提案优先聚焦 `eval --type proposal`，因为提案是决策上游——提案阶段的盲点会向下游的 PRD 和设计阶段传播，修复成本随阶段递增。如果自由评审推迟一个季度上线：在当前 Forge 每周约 2-3 份提案的评审频率下，约 30 份提案将仅依赖 rubric 评审通过，其中根据上述 dogfooding 数据，每份提案平均 1-2 个领域特有盲点将不被发现。这些盲点将在实施阶段转化为返工成本。其他 eval 类型（PRD、design、UI）存在相同的 rubric 盲点问题，将在 proposal 类型验证后复用相同架构扩展。

## Proposed Solution

为 eval-proposal 增加 `--freeform-expert` 参数，启用后在 rubric 循环之前插入 **Phase 0 自由专家评审**阶段。不传该参数时行为与现有 eval-proposal 完全一致：

1. **参数控制**：`forge eval --type proposal --freeform-expert` 启用自由专家阶段；不传参数时走标准 rubric 流程
2. **动态专家生成**：分析提案内容（domain、技术栈、复杂度、关键决策），推断最适合评审的专家档案（背景、专业领域、评审风格）。用户通过 AskUserQuestion 确认：可选择「接受此专家」「修改专家描述」或「重新生成」。选择修改时，用户直接输入文字描述期望的专家方向，系统基于修改后的描述重新生成档案并再次确认。最多允许 3 轮修改，超过后系统提示用户接受当前最佳版本或跳过自由评审
3. **自由叙事评审**：该专家以纯叙事形式对提案进行深度评审——无 rubric、无评分、无预设维度，完全由专家自主决定关注什么
4. **发现提取与注入**：通过 LLM 结构化输出提取自由评审叙事中的 key findings。提取 prompt 完整模板如下：

   > **System**: 你是一个分析助手。任务是从自由评审叙事中提取结构化的风险发现。
   >
   > **User**: 从以下自由评审叙事中提取所有显式提出的风险点和改进建议。
   >
   > 输出格式：JSON 数组，每个元素包含：
   > - `summary`: 一句话概括（不超过 50 字）
   > - `severity`: high / medium / low
   > - `quote`: 叙事中的原文引用（精确到句）
   >
   > 规则：
   > 1. 仅提取评审者明确表述的风险，不要推断隐含风险
   > 2. 每个风险点独立成条，不合并
   > 3. severity 基于评审者的语气强度判断（如使用「严重」「必须」为 high）
   > 4. quote 必须是原文逐字引用
   >
   > 叙事内容：
   > {{FREEFORM_REVIEW}}

   提取流程：自由评审叙事 + 上述提取 prompt → LLM 提取 → 验证提取产出非空且 JSON 格式合法 → 验证每个元素 summary/severity/quote 三个字段均非空 → 注入 rubric scorer prompt。若提取产出为空或格式校验失败，降级为标准 rubric 流程。需坦诚承认：从自由叙事中提取结构化发现本身是一个有损的机械过程——叙事中的推理链条、语气强度、上下文关联在提取为 bullet points 时不可避免地丢失。本提案接受这一损失，理由是：(1) 提取的 target 是 rubric scorer 而非人类读者，scorer 需要的是明确的 attack points 而非细腻的叙事；(2) 完整的自由评审叙事会同步保存在评审产出中供用户查阅，提取仅作为注入 rubric 的中间产物。提取完整性通过以下方式保障：(1) 提取 prompt 明确要求覆盖所有显式风险，(2) 提取产出与完整叙事一同保存供人工抽检，(3) 低于 50% 命中率时系统自动告警（见「部分提取失败」场景）
5. **专家持久化与复用**：动态专家档案保存到 `docs/experts/` 全局目录，后续评审可复用已有专家。专家档案包含质量追踪机制：每次使用该专家后，记录「注入后 rubric 评分是否有实质变化」（与无注入基线对比），「实质变化」定义为 rubric 评分差异 ≥ 15 分（即 1000 分制的 1.5%，高于 LLM 重复运行的典型评分方差）或 attack points 列表发生变动（新增/删除/修改至少 1 条）。若连续 3 次使用某专家均无实质变化，系统自动将该专家标记为 `deprecated: true`，后续匹配时跳过已弃用专家。用户也可手动标记专家为弃用状态

### Innovation Highlights

**动态专家生成**区别于行业常见的静态角色定义（如「你是一个 CTO」）。系统根据文档内容推断专家背景，使评审视角与文档特性匹配。例如：一个后端性能优化提案可能生成「分布式系统架构师，专注高并发场景」，而一个用户体验提案可能生成「产品心理学家，擅长行为设计分析」。

**自由叙事 → 结构化注入**的管道设计，兼顾了自由度和系统性：自由评审阶段不受 rubric 约束，但产出经过提取后进入 rubric 循环，确保后续评分能覆盖自由评审发现的盲点。

**专家库积累**：随使用积累的 `docs/experts/` 目录形成可复用的专家库，越用越丰富。

## Requirements Analysis

### Key Scenarios

- **标准评审（无参数）**：`forge eval --type proposal` → 走现有 rubric 流程，行为完全不变
- **自由专家评审**：`forge eval --type proposal --freeform-expert` → 先走 Phase 0 自由评审，再走 rubric 流程
- **新专家生成**：评审一篇关于「测试框架插件化」的提案 → 系统推断需要一位「测试工具链架构师」→ 用户确认 → 专家进行自由评审
- **已有专家复用**：评审另一篇类似领域的提案 → 系统在 `docs/experts/` 中找到匹配的已有专家 → 用户确认复用
- **用户修改专家**：系统推断的专家不合适 → 用户修改专家档案 → 保存修改后的版本 → 继续评审
- **自由评审发现注入**：自由评审发现了「提案的扩展性假设缺乏验证」→ 该发现被提取为 attack point → 注入 rubric scorer → 后续评分覆盖此维度
- **专家生成失败**：LLM 生成的专家档案不连贯或与提案领域无关 → 系统检测到档案质量异常（如缺少领域关键词、背景描述空洞）→ 提示用户手动指定专家方向或降级为标准 rubric 流程
- **自由评审产出为空**：LLM 返回空内容或无实质发现的泛泛评价 → 提取阶段产出零条 key findings → 系统跳过注入步骤，直接进入标准 rubric 流程，并告知用户"自由评审未产出有效发现，已降级为标准流程"
- **用户反复拒绝专家**：用户连续 3 次拒绝系统生成的专家档案 → 系统提示用户手动输入专家描述，或选择跳过自由评审阶段
- **注入发现无效果**：自由评审产出的 key findings 被注入 rubric scorer，但 scorer 的评分和 attack points 与未注入时完全一致（通过 A/B 对比检测）→ 系统在评审报告中标注「自由评审发现未影响 rubric 评分」，提示用户该自由评审可能未产出与 rubric 维度相关的有效发现，由用户决定是否需要基于自由评审叙事手动补充评审
- **部分提取失败**：自由评审叙事包含 5 个风险点但提取仅成功 2 个（JSON 部分字段缺失或格式错误）→ 系统记录提取命中率（成功提取数 / 叙事中包含「风险」「问题」「隐患」「建议」等关键词的段落数），若命中率 < 50% 则在评审报告中标注「提取命中率低」，同时将完整自由评审叙事附加到评审报告供用户手动查阅

### Non-Functional Requirements

- **延迟容忍**：Phase 0 增加约一次 LLM 调用（专家生成 + 自由评审），eval 总时间增加约 30%——可接受
- **可审计性**：所有动态专家档案持久化到 `docs/experts/`，用户可事后审核
- **确定性**：同一提案 + 同一专家应产出方向一致的评审（非完全随机）。通过结构化 prompt 协议约束评审框架（固定评审段落结构：背景评估、关键风险识别、改进建议），并使用低 temperature（0.3）减少随机性。不要求逐字一致，但要求核心关注点列表的 Jaccard 相似度 ≥ 0.6

### Constraints & Dependencies

- 仅适用于 `eval --type proposal`
- 依赖现有 eval skill 的 scorer/reviser 子 agent 架构
- 专家档案格式需兼容现有 `experts/scorer/*.md` 的 prompt 格式
- 需遵守 `docs/conventions/forge-distribution.md` 的分发规范
- 要求底层 LLM 具备领域推理能力（如 Sonnet 级别或以上）；对高度专业化的领域（如密码学、编译器优化），专家生成质量可能下降，此时应降级为标准 rubric 流程

## Alternatives & Industry Benchmarking

### Industry Solutions

学术同行评审系统长期使用"领域匹配审稿人"机制——系统根据论文内容推荐具备相关背景的审稿人，审稿人以自由文本撰写评审意见，编辑从自由评审中提取决策依据。这一模式是学术出版的基础假设（尽管严格的 A/B 实证在学术界本身难以执行）：(1) 专家-内容匹配被普遍认为比固定角色更有可能覆盖领域特有问题，(2) 自由叙事评审能覆盖结构化检查表无法预设的问题。类比到提案评审：当前 eval-proposal 的固定 CTO 角色等同于所有论文都由同一位通用审稿人评审，缺少领域深度。

代码评审领域的趋势同样从「检查清单」走向「上下文感知评审」。虽然 GitHub Copilot Code Review 和 Phabricator 的 Herald Rules 面向代码而非提案，但它们的核心思路——根据内容动态调整评审策略——与提案评审的需求一致。关键类比：代码评审解决"通用 linter 发现不了架构问题"，提案评审解决"通用 rubric 发现不了领域特有盲点"。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本，不增加复杂度 | 不解决「机械评审」痛点，dogfooding 数据显示每份提案 1-2 个领域盲点未覆盖 | Rejected: 数据表明 rubric 单独不足以保障质量 |
| 混合首轮（自由+rubric 同步） | 内部方案 | 延迟最优（零额外迭代，单轮完成）；自由评审可参考 rubric 维度避免遗漏已知陷阱 | rubric 维度会隐式约束「自由」评审的关注点，降低发现意外盲点的概率；两路评审同步进行增加单轮 prompt 复杂度 | Rejected: 混合方案在延迟上有明确优势（节省约 30% 总时间），但代价是锚定效应——自由评审的注意力会被 rubric 维度吸引，导致其产出趋向 rubric 已覆盖的领域。这一取舍的本质是「速度 vs. 独立性」：我们选择独立性，因为本提案的核心价值在于发现 rubric 之外的盲点，如果自由评审被 rubric 锚定，其增量价值将大幅缩水。若后续实践表明延迟不可接受，可回退到混合方案作为降级策略 |
| **Pre-scorer 前置阶段** | 本提案 | 自由评审完全独立于 rubric（无锚定效应）；纯增量修改，不触碰现有流程；发现注入 rubric 后可提升后续评分的覆盖面 | 多一轮 LLM 调用（约 30% 时间增加、约 2k-4k 额外 token）；自由评审质量高度依赖专家档案的匹配准确度 | **Selected: 自由评审的独立性最大化了发现未知盲点的概率，且增量式修改风险最低** |

## Feasibility Assessment

### Technical Feasibility

完全可行。核心改动点：
1. eval skill / eval-proposal command 增加 `--freeform-expert` 参数解析
2. eval skill 的 proposal 类型处理流程中增加 Phase 0（仅在参数启用时进入）
3. 新增一个「专家推断 + 自由评审」子 agent prompt
4. 新增发现提取逻辑
5. 修改 rubric scorer prompt 以接收注入的发现

所有改动均在现有 eval skill 架构内完成，无需新框架或外部依赖。

### Resource & Timeline

预计 4-6 个 coding task：
1. `--freeform-expert` 参数解析 + 条件分支（command + skill SKILL.md）
2. 专家推断逻辑 + 专家档案模板
3. 自由评审子 agent prompt + 协议
4. 发现提取 + 注入机制
5. eval skill 集成（Phase 0 编排）
6. 专家持久化与复用逻辑

加上 doc 类型 task（专家档案模板文档、协议文档），总量在 quick mode 范围内，预计 2-3 个工程日可完成实现。

### Dependency Readiness

无外部依赖。所有改动在 `plugins/forge/skills/eval/` 目录内完成。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 自由叙事评审比 rubric 打分更有洞察力 | Assumption Flip | Confirmed: rubric 覆盖已知模式但遗漏文档特有问题；自由评审捕获盲点但可能遗漏已知陷阱。两者互补而非替代 |
| 动态生成的专家足够可靠 | Stress Test | Refined: 需要用户确认环节作为安全网，避免不合适的专家产出低质量评审 |
| Phase 0 增加的延迟可接受 | Occam's Razor | Confirmed: 一次额外 LLM 调用（约 30% 时间增加）换来显著提升的评审质量，ROI 合理 |

## Scope

### In Scope

- `--freeform-expert` 参数解析与条件分支（启用 / 未启用两条路径）
- 动态专家档案推断机制（分析提案 → 推断专家 → 生成详细档案）
- 用户确认机制（接受 / 修改 / 重新生成）
- 专家档案持久化到 `docs/experts/` 全局目录
- 已有专家复用逻辑（匹配提案内容与已有专家档案）
- 自由叙事评审协议（纯叙事、无 rubric、无评分）
- 发现提取机制（从叙事中提取结构化 key findings）
- 注入机制（将发现作为额外上下文注入 rubric scorer）
- eval skill 的 proposal 类型集成（Phase 0 编排）

### Out of Scope

- 扩展到其他 eval 类型（prd、design、ui 等）——优先聚焦 proposal 是因为提案处于决策上游，其盲点向下游传播且修复成本递增；proposal 类型验证后将复用相同架构扩展至其他类型
- 修改现有 proposal rubric 本身
- 多专家并行自由评审
- 对自由评审产出进行评分
- 任何 UI/交互式组件

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| AI 推断的专家角色不合适 | M | M | 用户确认环节可修改或重新生成 |
| 自由叙事产出难以提取结构化发现 | M | H | 提取协议要求评审产出中标注关键段落；若提取产出为空，跳过注入步骤并降级为标准 rubric 流程，向用户报告降级原因 |
| 专家库膨胀导致匹配困难 | L | L | 命名约定 + domain 标签帮助检索；质量追踪机制自动弃用低效专家，控制有效专家数量 |
| LLM 生成虚假领域专家（hallucinated expertise） | M | H | 专家档案必须包含可验证的领域关键词和背景描述；用户确认环节作为最终校验。当用户无法判断专家的领域胜任度时，系统提供辅助校验：(1) 自动将专家档案中的领域关键词与提案中的技术术语做交叉引用，生成「关键词覆盖率」报告供用户参考；(2) 为专家档案附加 3-5 个自检问题（如「该专家是否覆盖了提案中提到的 X 技术？」），用户通过回答这些问题间接验证专家匹配度；发现注入时 rubric scorer 仍独立评分，不会被虚假专家意见主导 |
| 自由评审发现与 rubric 评分方向矛盾 | M | M | 矛盾本身是有价值的信息——rubric scorer 在注入发现时应同时接收矛盾标记，在评分报告中明确标注"自由评审与 rubric 存在分歧"，由用户最终判断 |
| 同行评审类比的实证基础薄弱 | L | M | 本提案的核心理念（动态专家 + 自由叙事）基于学术同行评审的类比推理，而非严格的 A/B 实证数据。若实践表明自由评审未产出 rubric 之外的增量发现，可通过 success criteria 中的 A/B 对比检测及时识别，并将此特性降级或移除 |

## Success Criteria

- [ ] eval-proposal 不传 `--freeform-expert` 时，对同一文档的评审输出与未引入该参数前相比，rubric 评分差异 ≤ 5%（以 5 份历史提案的评审结果作为基线对比）
- [ ] 传入 `--freeform-expert` 时进入 Phase 0，生成动态专家档案并经用户确认
- [ ] 自由评审产出为纯叙事格式，无 rubric 维度、无评分
- [ ] 自由评审的 key findings 成功提取并注入后续 rubric scorer
- [ ] 动态专家档案保存到 `docs/experts/`，每个档案必须包含以下必填字段：`domain`（适用领域）、`background`（专业背景）、`review_style`（评审风格描述）、`generated_for`（首次生成时的目标提案路径）、`created_at`（生成时间戳）、`review_history`（使用该专家的评审记录列表，含提案路径和评审日期），格式为 YAML front matter + Markdown 正文
- [ ] 后续评审可复用 `docs/experts/` 中的已有专家
- [ ] eval 总时间增加不超过 40%（基线：当前单轮 eval 时间）
- [ ] 专家废弃机制可被触发和验证：对同一专家连续 3 次评审均无「实质变化」（rubric 评分差异 < 3 分且 attack points 列表未变动）时，该专家的 YAML front matter 中 `deprecated` 字段自动置为 `true`
- [ ] 注入有效性达标：在 5 份历史提案的 A/B 对比中，注入自由评审发现后 rubric scorer 的 attack points 列表至少在 3 份（60%）中出现变动（新增或修改 attack points），否则系统在评审报告中标注「注入效果不显著」警告
- [ ] eval report 的 ATTACKS 列表包含 `[beyond-rubric]` 标签条目：对于确实无法映射到 rubric 维度的自由评审发现，scorer 必须将其以 `[beyond-rubric]: [finding]` 格式记录在 ATTACKS 列表末尾（与 `[blindspot]` 并列），格式校验：至少包含 summary 和 quote 两个子字段
- [ ] 对同一份提案，freeform+rubric 评审最终产出的 attack points 数量 ≥ 纯 rubric 评审的 attack points 数量 × 1.3（以 5 份历史提案作为对比基线），且自由评审阶段产出的 attack points 中至少 30% 不在纯 rubric 评审的关注点列表内；鉴于样本量有限，量化指标需辅以人工定性审核：由提案作者逐份判定自由评审发现的 attack points 是否揭示了 rubric 未覆盖的真实风险

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
