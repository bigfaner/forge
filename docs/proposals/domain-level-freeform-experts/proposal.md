---
created: "2026-05-25"
author: "faner"
status: Draft
---

# Proposal: 领域级自由专家生成机制

## Problem

自由专家（freeform expert）的动态生成以单一 proposal 为锚点，生成的专家领域范围极窄（如 "Build Orchestration & Test Infrastructure Expert" 仅服务于 surface-aware-justfile proposal），导致跨 proposal 复用率近乎为零——12 个已生成的专家之间 Jaccard 相似度无法达到 0.3 的复用阈值。

### Evidence

- `docs/experts/` 下 12 个专家文件，对应 11 个唯一 proposal（build-orchestration-test-infra 和 surface-aware-dispatcher-orchestrator 共享同一个 proposal: surface-aware-justfile）
- `domain` 关键词高度特化：如 "pipeline-integration, type-system-categorization, dead-code-removal, go-backend, prompt-template-architecture"——5 个关键词中任何一个出现在其他 proposal 的概率极低
- 实际复用匹配从未成功过——P0.1 复用检查对全部 12 个专家计算 Jaccard 相似度，最高仅 0.095（prompt-compliance-architect vs 当前 proposal），远低于 0.3 阈值。每次评估都触发了全新专家生成

### Urgency

每次评估都经历完整的专家生成→确认循环（3 轮修改/拒绝上限），增加评审耗时。领域级专家一次生成即可服务同一领域的多个 proposal，显著降低评审启动成本。

## Proposed Solution

改造 `expert-inference.md` 的专家生成流程，引入预定义领域分类表实现分层识别：

1. **大领域匹配**：从 proposal 提取特征后，查分类表确定所属大领域（如"构建与测试基础设施"）
2. **领域内专家生成**：LLM 在该大领域范围内生成专家，关键词和背景覆盖整个领域而非单一 proposal

分类表以 Markdown 有序列表嵌入 `expert-inference.md` 的 prompt 正文中（紧接两步推理指令之后），格式为编号列表，每项包含领域名称和一句话描述（如 `1. **构建与测试基础设施** — task pipeline、justfile 构建、测试编排`）。预期大小：初始 6 个条目，上限 10 个条目，每个条目约 20-30 tokens，总嵌入量 < 300 tokens。该格式选择基于：(1) Markdown 列表与 prompt 文件格式一致，无额外解析依赖；(2) 有序列表便于 LLM 按编号引用，降低选择歧义。

两步推理的输出格式：LLM 在单次调用中输出结构化响应，使用 XML 标签分隔两个阶段——`<domain-match>` 输出匹配的领域编号（如 `3` 对应分类表第 3 项），`<expert-profile>` 输出完整的专家定义（按 expert-template.md 格式）。主控代码通过正则提取 `<domain-match>` 内容确定领域归属，`<expert-profile>` 内容作为专家文件写入。

同时更新 `expert-template.md` 增加 `scope` 字段，更新 `freeform-expert-persistence.md` 的复用匹配逻辑以适配新格式。

### 用户可观测行为

评估运行者在以下环节会体验到与当前系统的差异：

1. **专家生成阶段**：当 LLM 确定大领域后，确认提示中新增一行 `scope: domain-level` 及对应的大领域名称（如 `scope: domain-level [构建与测试基础设施]`），取代原先仅显示窄领域关键词的模式
2. **复用命中时**：当已有 domain-level 专家被复用，评估日志输出 `Reusing existing domain-level expert: <expert-name> (scope: domain-level)`，取代当前的 `Generating new expert...` 流程，用户无需经历生成→确认循环
3. **跨领域 proposal**：当 proposal 匹配到多个大领域时，确认提示默认推荐最高匹配领域，用户可通过现有的 Modify 选项切换到其他匹配领域——不引入新的交互步骤
4. **Findings 报告**：由 domain-level 专家产出的 findings 会标注适用范围（如 `[domain-level: 构建与测试基础设施]`），区分于 proposal-specific findings
### 跨领域策略

当 proposal 匹配到多个大领域时，默认使用匹配得分最高的单一领域专家进行评审。用户可在确认循环中通过 Modify 指定其他匹配领域，以覆盖面换取成本控制。不引入多专家串行评审——多专家协调涉及 pipeline 编排变更，超出本 proposal 范围。

对现有 12 个 proposal-specific 专家的 `domain` 字段做聚类，可归并为以下大领域：

| 大领域 | 覆盖的现有专家 | 自然聚类依据 |
|--------|--------------|-------------|
| 构建与测试基础设施 | build-orchestration-test-infra, test-pipeline-architect, surface-aware-dispatcher-orchestrator | 均涉及 task pipeline、justfile 构建、测试编排 |
| Go 代码库健康 | go-codebase-health-engineer, go-pipeline-integration-type-system-engineer | 均涉及 Go 重构、代码质量、类型系统 |
| 评估管线架构 | eval-pipeline-infoflow-architect, requirements-consistency-forge-eval-specialist, expert-system-design-prompt-architecture | 均涉及 eval pipeline 设计、评分一致性、专家系统 prompt 架构 |
| Prompt 与 Agent 协议 | prompt-compliance-architect | prompt 模板、agent 行为控制 |
| Schema 与数据回归 | config-schema-surface-detection, golden-dataset-regression-architect | 均涉及 schema 迁移、快照测试 |
| 文档审计 | documentation-implementation-drift-auditor | 文档-代码一致性 |

12 个专家自然归并为 **6 个大领域**，加上未覆盖的潜力领域（如 UI/设计、知识管理、CLI 交互），上限为 **10 个大领域**（初始 6 个 + 预留 4 个增长空间）。

### Innovation Highlights
**分层领域识别**：通过预定义分类表缩小 LLM 的领域推断搜索空间。分类表将"从无限可能的领域名中选一个"简化为"从有限候选项中选一个"，降低了不一致的概率（从 "test infrastructure" vs "testing pipeline" 这类自由发散，收敛为在预定义的"构建与测试基础设施"条目中确认）。但这并非绝对保证——LLM 仍可能将边界模糊的 proposal 映射到不同条目。大多数专家系统要么用固定专家库（无灵活性），要么完全依赖 LLM 自由推断（无一致性）——分层方案在两者之间取得平衡，本质是用有限的选择集换取更高的匹配可靠性。

## Requirements Analysis

### Key Scenarios

- **场景 1（新领域首次评估）**：评估一个属于"构建与测试基础设施"领域的 proposal → 分类表匹配成功 → 生成该领域专家 → 保存到 `docs/experts/`
- **场景 2（同领域再次评估）**：评估另一个同领域 proposal → 复用匹配命中已有专家 → 跳过生成 → 直接用于评审- **场景 3（跨领域 proposal）**：proposal 涉及多个领域（如"Agent架构" + "配置Schema"）→ 匹配所有相关大领域 → 默认使用得分最高的单一领域专家评审 → 用户可通过 Modify 切换到其他匹配领域。**已知限制**：单一领域专家可能遗漏其他领域的评审视角，用户需手动切换领域以获取全面覆盖
- **场景 4（分类表未覆盖）**：proposal 属于分类表外的领域 → `<domain-match>` 输出 `0`（特殊值表示无匹配）→ prompt 中的 fallback 指令触发，LLM 切换为自由推断模式生成领域名和专家定义（不受分类表约束），专家的 `scope` 字段仍标记为 `domain-level`。用户确认时提示"该领域未在分类表中，已自动推断"

### Non-Functional Requirements

- **可扩展性**：新增领域仅需在 `expert-inference.md` 的分类表中追加一行（约 30 tokens），无需修改其他文件。扩展操作不涉及 schema 变更或代码改动
- **性能**：两步推理（领域匹配 + 专家生成）在单次 LLM 调用内完成（分类表嵌入 prompt，LLM 先输出领域编号再输出专家定义），不增加 API 调用次数。分类表嵌入增加约 300 tokens 输入成本，相对于现有 expert-inference prompt（约 2000 tokens）增幅 < 15%。推理延迟不增加，因为分类表查询是在同一推理过程中完成的
- **旧专家隔离**：现有 12 个 proposal-specific 专家文件保留在 `docs/experts/`，通过 `scope` 字段过滤使其不再参与新系统的复用匹配——本质是隐性废弃：文件仍存在但永远不会被新的 domain-level 评估选中。这是刻意选择：旧专家的窄领域关键词会干扰新系统的匹配准确性
- **分类准确率**：LLM 将 proposal 映射到分类表的准确率应 >= 80%（验证：将 11 个唯一 proposal 的领域归属与人工标注对比，LLM 映射结果与人工标注一致的为正确分类）

### Constraints & Dependencies

- 改动仅限 `experts/freeform/` 目录下的 prompt 文件（含 `extraction-prompt.md`）和 `rules/freeform-expert-persistence.md`
- 不影响 freeform-review-protocol、scorer-composition、reviser-composition
- 用户确认循环（Accept / Modify / Regenerate）保持不变

## Alternatives & Industry Benchmarking

### Industry Solutions

业内常见的专家/评审者选择机制：
1. **固定专家库**（如学术同行评审的 TPC 成员列表）——预先定义角色和领域覆盖范围，论文提交时由程序自动分配审稿人（如 OpenReview 的 affinity score 匹配）
2. **语义向量匹配**（如 embedding-based reviewer assignment）——将文档和审稿人 profile 分别编码为向量，通过 cosine similarity 匹配。NIPS/ICML 等 AI 会议已广泛采用
3. **LLM 自由推断**（如 ChatGPT 的 Custom Instructions persona）——完全依赖模型判断
4. **混合方案**（如 Claude Code 的 multi-expert parallel scoring）——部分固定 + 部分动态

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动成本 | 每次评估重复造轮子，复用率 0%。实际成本量化：每次评估的专家生成→确认循环平均 3 轮交互（含用户修改/拒绝），按每轮 2 分钟人工 + 30 秒 LLM 推理估算，单次评估浪费约 7.5 分钟在可复用的专家生成上。12 个已有 proposal 已累计触发 12 次完整生成，总浪费约 90 分钟。随着 proposal 数量增长，该成本线性累积 | Rejected: 可复用成本的持续累积超过改造的 2-3 小时投入 |
| 纯 Prompt 重写 | LLM self-guided | 改动最小 | 与选定方案的本质区别：纯重写仅改善单个 prompt 的指令措辞（如"生成更广领域的专家"），但 LLM 对"更广"的理解仍因 proposal 而异，无法保证两次独立推理产生相同的领域标签。选定方案通过嵌入预定义分类表将领域选择从"开放式生成"约束为"从有限候选项中选取"，从根本上改变了推理任务的性质 | Rejected: 指令措辞无法约束 LLM 的标签一致性 |
| 固定专家库 | 学术 TPC 模式（OpenReview） | 最强一致性；每个领域一个精心调校的专家 profile | Forge 的领域空间虽有限（当前 6 个大领域），但每个大领域内的 proposal 主题差异大（如"构建与测试基础设施"下的 justfile 编排 vs Go 代码库健康），固定专家无法覆盖领域内多样性。且 Forge 新增领域（如未来的 UI/设计、CLI 交互）需要新增专家 profile，每次都需要手动撰写和调校 | Rejected: 无法覆盖领域内子主题多样性 |
| 语义向量匹配 | NIPS/ICML embedding matching | 无需人工维护分类表；语义匹配精度高 | 需要 embedding 模型（额外 API 调用）+ 将专家 profile 编码为向量（需要存储和索引基础设施）。对 Forge 的规模（当前 12 个专家）而言，引入向量检索的工程复杂度远超收益——12 个候选的暴力 Jaccard 匹配足够快 | Rejected: 工程复杂度与规模不匹配（12 个候选不需要向量检索） |
| **分类表引导** | 分层识别 | 一致性 + 灵活性兼顾；零额外基础设施依赖 | 分类表需维护（约每 3-5 个新领域新增一行） | **Selected: 在一致性、灵活性和工程简洁性之间取得最优平衡** |

## Feasibility Assessment
### Technical Feasibility

可行，但需注意 prompt 链联动复杂度。改动涉及 4 个 Markdown prompt 文件的协调修改：

1. **expert-inference.md**：新增两步生成流程（大领域匹配 → 领域内专家生成），需嵌入分类表并与后续匹配逻辑对齐
2. **expert-template.md**：新增 `scope` 字段，影响所有下游消费专家文件的 prompt
3. **freeform-expert-persistence.md**：Jaccard 匹配逻辑需区分 `domain-level` 与 `proposal-specific` 专家，避免不同 scope 级别的专家误匹配
4. **extraction-prompt.md**：从 freeform review 提取 findings 时需感知专家 scope 级别，在提取的 JSON 中为每个 finding 标注 `scope` 字段

四个文件的变更存在依赖关系：template 的 schema 变更是 inference 和 persistence 的前提，extraction 依赖 template 的 scope 字段。

### Resource & Timeline

预估工作量 2-3 小时，基于类比：`expert-inference.md` 的两步流程改造类似 `freeform-expert-persistence.md` 曾做的 Jaccard 匹配逻辑重构（单文件约 1 小时），`expert-template.md` 的字段新增类似 `extraction-prompt.md` 的格式变更（单文件约 30 分钟）。交叉验证（用 2-3 个已有 proposal 走通新旧流程）约 1 小时。四个文件的协调修改需按依赖顺序逐个完成，无法并行。

### Dependency Readiness

无外部依赖。所有改动在 Forge plugin 内部完成。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 分类表会限制专家的领域覆盖 | Assumption Flip | Refined: 分类表只定义大领域方向，LLM 在领域内仍有充分细化空间。未覆盖的领域有降级路径 |
| Jaccard 匹配足以区分新旧专家 | Occam's Razor | Confirmed: 新专家关键词更广，自然匹配度更高。旧专家会在自然竞争中逐渐被淘汰 |
| 每个 proposal 只需一个领域专家 | Stress Test | Refined: 跨领域 proposal 匹配得分最高的单一领域专家。若评审视角不足，用户可通过 Modify 切换到其他匹配领域 |

## Scope

### In Scope
- `expert-inference.md` 嵌入领域分类表，改造为两步生成流程（大领域匹配 → 领域内专家生成）
- `expert-template.md` 增加 `scope` 字段（`domain-level` / `proposal-specific`），同时在 `domain` 字段中存储大领域名称（如 `"构建与测试基础设施"`），取代当前的窄关键词列表- `freeform-expert-persistence.md` 更新复用匹配逻辑：区分 `scope` 级别（domain-level 专家仅与 domain-level 专家匹配），过滤旧 proposal-specific 专家的噪声
- `extraction-prompt.md` 更新 findings 提取逻辑：在提取的 JSON 数组中为每个 finding 增加 `scope` 字段（`domain-level` 或 `proposal-specific`），使下游消费者（Pre-Revision、Scorer）可按 scope 标注 findings 的适用范围。改动内容：extraction prompt 的输出 schema 新增 `scope` 字段，取值从专家文件的 `scope` front matter 字段读取

### Out of Scope
- 推广自由专家到 PRD / tech-design / ui-design 评估（todo #166）
- 修改 freeform-review-protocol、scorer-composition、reviser-composition
- 修改 freeform-pipeline 的跨步骤编排逻辑（如 scorer、reviser 的调用顺序）——注意：`expert-inference.md` 内部的两步推理（领域匹配 → 专家生成）是对单个 prompt 内部推理流程的改造，不涉及 pipeline 步骤间的编排变更
- 迁移或废弃现有 12 个专家文件（通过 scope 字段隔离，不主动迁移）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 分类表初期覆盖不全，部分 proposal 无法匹配 | M | L | 提供 LLM 自由推断降级路径，未匹配时自动降级 || 领域级专家的评审深度不如 proposal-specific 专家 | M | M | 承认 trade-off：领域级专家覆盖面更广但单 proposal 深度可能不足。缓解机制为：(1) 专家 prompt 中注入当前 proposal 全文作为评审焦点，使领域专家在宽泛背景下针对 proposal 的特定问题发力；(2) 用户确认循环（Accept / Modify / Regenerate）作为最终质量保障——若深度不足，用户可通过 Modify 要求专家聚焦特定子领域 |
| 分类表随时间膨胀难以维护 | L | L | 分类表控制在大领域粒度，初始 6 个上限 10 个，不细化到子领域 || 旧 proposal-specific 专家意外参与新系统匹配 | L | M | `scope` 字段过滤确保 domain-level 匹配仅匹配 domain-level 专家。若过滤逻辑有 bug，旧专家的窄关键词可能产生噪声 |

## Success Criteria
- [ ] 新生成的专家 `domain` 关键词覆盖范围 ≥ 2 个 proposal 的领域交集。计算方式：从同大领域内选取 2 个已有 proposal，提取其 `domain` 关键词集合 K₁ 和 K₂，新专家的关键词集合 K_new 须满足 `|K_new ∩ K₁| / |K₁| ≥ 0.5` 且 `|K_new ∩ K₂| / |K₂| ≥ 0.5`（即新专家至少覆盖每个 proposal 一半以上的领域关键词）。**已知 trade-off**：覆盖面扩大会稀释单 proposal 的评审深度，此 SC 仅衡量覆盖面，深度保障由用户确认循环兜底
- [ ] 同领域内的第二个 proposal 评估时，复用匹配成功（当前为 0 成功）
- [ ] `extraction-prompt.md` 正确区分 domain-level 与 proposal-specific findings：domain-level 专家产出的 findings 在 freeform review 中被标记为适用范围更广，不因单一 proposal 的上下文被过度窄化（验证：用同一 domain-level 专家评审两个同领域 proposal，两者共享的 findings 比例 ≥ 30%。"共享 findings"判定算法：对两次评审的 findings 列表 F₁ 和 F₂，逐对比较 summary 字段，若两 findings 的 summary 语义等价（由 LLM 判定：将两个 summary 拼接后询问"是否表达同一问题"，回答 yes 则视为共享），则计为 1 个 shared finding。比例 = |shared(F₁, F₂)| / max(|F₁|, |F₂|)）
- [ ] 专家生成后经用户确认的轮次 ≤ 2（当前经常需要修改以扩大领域范围）
- [ ] `scope` 字段在复用匹配中正确隔离：使用 12 个现有专家（均为 `proposal-specific`）测试，当评估一个生成 `domain-level` 专家的 proposal 时，复用匹配的候选列表中不包含任何 `proposal-specific` 专家（验证：检查 P0.1 候选列表的 scope 字段值）
- [ ] 分类表覆盖 ≥ 80% 的已有 proposal（分母为 `docs/experts/` 下已有专家文件对应的唯一 proposal 数量，即当前 11 个；验证：将每个已有 proposal 逐一匹配分类表，统计成功匹配的比例）

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
