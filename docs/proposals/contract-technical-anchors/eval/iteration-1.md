# Proposal Evaluation Report: Contract Technical Anchors

**Evaluator**: CTO Adversarial Review
**Date**: 2026-06-05
**Iteration**: 1
**Document**: `docs/proposals/contract-technical-anchors/proposal.md`

---

## Phase 1: Reasoning Audit

### 1. Problem → Solution: Does the proposed solution address the stated problem?

**Alignment**: Good, but with a gap.

The problem states: "Contract 规格缺少技术锚点...导致 gen-test-scripts 只能依赖 LLM 推断技术细节...当推断错误时测试脚本与实际代码不匹配"

The solution provides: handbook 生成 → Contract 锚点填充 → 交叉验证 → 建议修复。这确实建立了一条信息链。

**Gap**: 提案标题说"Contract Technical Anchors"，但解决方案的核心工作量集中在"tech-design 生成 handbook"和"gen-test-scripts 交叉验证"。Contract 的锚点字段本身只是一个中间数据结构。提案没有明确说明当 handbook 不存在（且用户选择不生成）时的完整降级行为——虽然 Key Scenarios 提到了"缺少 handbook"场景，但 Success Criteria 对此只要求"管道正常运行"，缺乏对降级后测试准确率的任何承诺。

### 2. Solution → Evidence: Does evidence support the solution?

**Weakness**: 证据只来自一个项目（pm-work-tracker）的一个场景（Move sub-item）。

- 单一 case study，无法证明问题在其他 surface 类型（CLI、Web、Mobile）上是否存在
- "此类问题会在每个有 API 或 CLI surface 的项目中重复出现"是一个未经验证的泛化声明
- 没有数据说明当前 gen-test-scripts 推断错误的发生频率

### 3. Evidence → Success Criteria: Do SC test what matters?

**Partially**: Success Criteria 测试的是锚点覆盖率（100% coverage）和交叉验证功能，但没有测试提案声称要解决的核心问题——"测试脚本使用正确的技术细节"。

- SC 测的是"Contract 有 endpoint 字段"和"交叉验证能捕获 POST vs PUT"，但没有 SC 测量"测试脚本与实际代码的匹配率提升了多少"
- 缺少对 handbook 生成质量（格式正确性、完整性）的验证标准

### 4. Self-contradiction check

- "全 surface 覆盖" vs "分批实现（Phase 1: API → Phase 2: CLI → Phase 3: Web/Mobile）"——不矛盾，分批实现是合理的 scope 管理
- `consistency_check_result: status: pass, pairs_checked: 15, conflicts_found: 0`——这个嵌入在 SC 部分的 YAML 块看起来像是自动检查工具的输出，但文档没有解释它的含义或来源。这引入了理解歧义

---

## Phase 2: Rubric Scoring

### 1. Problem Definition: 78/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 32/40 | 核心问题清晰："Contract 缺少技术锚点 → LLM 推断 → 推断错误 → 测试不匹配"。但"三层测试均无法捕获"这一表述将问题范围扩大到整个测试体系，而实际证据仅涉及 API surface 的一个操作 |
| Evidence provided | 22/40 | 提供了 pm-work-tracker 的具体案例，有明确的技术细节（POST vs PUT）。但只有一个项目的单一场景，无频率数据，无其他项目/其他 surface 类型的佐证。"三层测试全部漏掉"的描述虽然具体，但仅反映该案例而非普遍情况 |
| Urgency justified | 24/30 | "会在每个有 API 或 CLI surface 的项目中重复出现"提供了理由，但缺乏量化依据（多少项目受影响？多久发生一次？）。"当前 Fact Table 代码侦察...不会与 Contract 做交叉比对"说明了为什么现有机制不够 |

### 2. Solution Clarity: 88/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 32/40 | 5 步方案具体可解释：handbook 生成 → 锚点字段 → 填充 → 交叉验证 → 新鲜度检查。但第 4 步"交叉验证"的分类逻辑（高置信度/低置信度/无法验证）缺少分类标准——什么条件下判定为哪种置信度？ |
| User-facing behavior described | 38/45 | Key Scenarios 较好地描述了用户会经历什么：生成 handbook → Contract 有锚点 → 验证发现不匹配 → 用户确认修复 → 输出覆盖报告。但"用户确认后写入"的交互细节不足——是 inline 确认还是单独的审批步骤？ |
| Technical direction clear | 18/35 | "改动集中在三个现有 skill 内部"提供了方向，但缺乏关键细节：锚点字段在 Contract 文件中的具体位置（frontmatter 位置说明不够）、handbook 格式定义只是列出了字段名而无格式示例、交叉验证的具体比对算法未说明 |

### 3. Industry Benchmarking: 52/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 22/40 | 提到了 Pact 和 OpenAPI spec，但仅一段简短描述："行业常见的做法是 Contract Testing（Pact）或 OpenAPI spec 驱动的测试生成"。没有说明这些方案的具体做法、适用场景、与 Forge 的差异点。引用过于浅层 |
| At least 3 meaningful alternatives | 18/30 | 4 个选项（Do nothing、仅增强 Fact Table、OpenAPI spec 驱动、Contract Technical Anchors）。但"仅增强 Fact Table"是一个弱 straw man——它的 Cons 栏写"不解决设计-实现一致性问题"，但这是该方案本来就不打算解决的问题，不公平地否定了它。"OpenAPI spec 驱动"被标为"架构不匹配"但没解释为什么 Forge 的语义 Contract 与 OpenAPI 不兼容 |
| Honest trade-off comparison | 6/25 | Cons 栏存在 cherry-picking。对选定方案的 Cons 只写了"需要扩展 tech-design 和 gen-test-scripts"，这是一个低估了工作量的表述——还需要定义 4 种新的 handbook 格式、修改 eval-contract、增加交叉验证报告。对其他方案的 Cons 描述过于笼统 |
| Chosen approach justified against benchmarks | 6/25 | "最小改动覆盖最大范围"是一个断言而非论证。没有提供工作量估算来支持"最小改动"的说法。Phase 路线图暗示工作量不小（3 个 phase），与"最小改动"矛盾 |

### 4. Requirements Completeness: 72/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 28/40 | 5 个场景覆盖了主要路径：新功能、已有功能、不一致、缺少 handbook、handbook 过期。但缺少：handbook 存在但内容为空的情况、一个 Contract 对应多个 endpoint 的情况（批量操作）、锚点字段格式错误的情况 |
| Non-functional requirements | 28/40 | 向后兼容和性能影响都有提及。但"无显著额外开销（仅增加内存比对）"是主观判断，没有提供数据或分析支持。缺少对存储影响的说明（锚点字段增加的文件大小）。安全性方面未提及——锚点字段是否可能暴露内部 API 路径？ |
| Constraints & dependencies | 16/30 | 列出了 3 个依赖，但缺少：Claude API 的 context window 限制（大量 handbook 内容填充到 Contract 时的 token 限制）、各 surface handbook 格式定义的一致性约束、gen-test-scripts 当前架构是否支持新增交叉验证步骤 |

### 5. Solution Creativity: 55/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 18/40 | "设计文档驱动的锚点验证"核心思路并不新颖——这是 OpenAPI spec 驱动测试生成模式在 Forge 语义框架下的适配。Innovation Highlights 将其包装为"设计文档为 authority source"的洞察，但 OpenAPI 生态已有成熟的 schema-to-test 管道。与行业基线的区分度不高 |
| Cross-domain inspiration | 22/35 | 从 API handbook 扩展到 CLI/Web/Mobile 是合理的横向扩展。surface 统一验证的思路有一定创新性。但没有借鉴其他领域的类似问题解决方案（如编译器的类型检查、IDE 的 refactoring 安全检查） |
| Simplicity of insight | 15/25 | "Contract 需要技术锚点"的洞察直觉上合理，但实现方案引入了 4 种新 handbook 格式 + 3 phase 路线图 + 交叉验证分类体系，复杂度不低。"为什么之前没想到"的感觉更多是因为问题本身不难发现，而非因为解决方案精妙 |

### 6. Feasibility: 68/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 28/40 | 三个 skill 的改动方向可行。但"cli-handbook / page-map / screen-map 是新增文档类型"意味着需要定义 4 种新格式（含 API 的 reference 模式），这部分工作量被低估了。交叉验证的分类逻辑（高/低/无法验证）需要更精确的定义才能实现 |
| Resource & timeline | 22/30 | "单项 enhancement，改动点明确"——但 Scope 列出了 6 个 In Scope 项，涉及 4 个 skill 改动 + 4 种新文档格式。Phase 路线图暗示至少 3 个迭代。这与"单项 enhancement"的定位不匹配 |
| Dependency readiness | 18/30 | "api-handbook 已稳定运行"是对的，但 cli-handbook / page-map / screen-map 格式"无前置依赖"≠"已就绪"。这些格式需要从零设计，且格式设计质量直接影响整个方案的可行性。格式设计被标记为风险（handbook 格式设计不当）但没有给出具体的格式设计方案 |

### 7. Scope Definition: 52/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 20/30 | 6 项都是可交付物，但部分描述仍偏模糊："交叉验证输出 surface 覆盖报告"——报告格式？"明确列出已验证和未验证的 surface"——什么样的列表？ |
| Out-of-scope explicitly listed | 16/25 | 5 项 Out of Scope 清晰列出了推迟的内容。"Contract 手动编辑后的锚点漂移检测"被推迟但 `last_anchor_sync` 时间戳仍在 In Scope 中——这是部分实现 |
| Scope is bounded | 16/25 | Phase 路线图提供了分批策略，但 Phase 1 的具体范围未明确——是只做 API endpoint + method，还是也包含 eval-contract 的锚点完整性检查？Phase 之间的依赖关系未说明。6 个 In Scope 项跨 3 个 Phase，但没标注每项属于哪个 Phase |

### 8. Risk Assessment: 62/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 22/30 | 6 个风险，覆盖了主要方面。但缺少：handbook 格式碎片化风险（4 种格式各自演化导致维护困难）、交叉验证结果过多"低置信度"时的用户体验退化风险、Phase 1 完成后 Phase 2/3 可能因格式不兼容需要返工的风险 |
| Likelihood + impact rated | 18/30 | 评级大致合理，但"建议修复仍有误（设计文档本身有误）"标记为 L likelihood——设计文档本身有误的频率有多低？没有数据支持。pm-work-tracker 案例中 api-handbook 是正确的，但这是否是常态？ |
| Mitigations are actionable | 22/30 | 多数缓解措施可执行："复用 api-handbook 的成熟格式模式"、"用户确认环节"、"分批实现"、"handbook 新鲜度检查"。但"低置信度不自动处理"——那具体怎么处理？降级提示的交互设计是什么？"交叉验证报告明确列出已验证和未验证的 surface"——列出后用户需要做什么？ |

### 9. Success Criteria: 48/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 16/30 | "100% 包含 endpoint 字段"——可测量。"交叉验证能捕获 lesson 场景"——可测试但仅覆盖单一场景。"管道正常运行"——模糊，什么是"正常"？"生成明确的代码 bug 标记报告"——"明确"不可测量。嵌入在 SC 开头的 YAML 块（`consistency_check_result`）来源不明，不应作为 SC 的一部分 |
| Coverage is complete | 12/25 | 缺少对以下 In Scope 项的 SC：tech-design handbook 生成、`last_anchor_sync` 时间戳、handbook 新鲜度检查、eval-contract 锚点完整性检查。6 个 In Scope 项中只有约 4 个被 SC 覆盖 |
| SC internal consistency | 20/25 | SC 之间大致一致。"100% 包含 endpoint"与"缺少 handbook 时管道正常运行"通过"当 api-handbook 存在时"的限定条件解决了潜在冲突。`consistency_check_result: status: pass` 的 YAML 块含义不清，可能造成理解歧义 |

### 10. Logical Consistency: 65/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 25/35 | 方案确实解决了"Contract 缺少技术锚点"的问题。但问题仅基于 API surface 的证据，方案却扩展到全 surface——这个扩展缺乏等量的证据支撑 |
| Scope ↔ Solution ↔ Success Criteria aligned | 18/30 | In Scope 的"eval-contract 评分规则增加技术锚点完整性检查"在 Solution 部分未被提及。In Scope 的"handbook 新鲜度检查"在 Solution 第 5 点有提及但无对应 SC。`last_anchor_sync` 时间戳在 In Scope 中但 Solution 和 SC 都缺乏详细说明 |
| Requirements ↔ Solution coherent | 22/25 | Key Scenarios 与 Solution 的 5 步方案对应关系良好。Anchor Field Schema 提供了字段定义。但多 endpoint 的 Contract 场景（一个操作可能涉及多个 API 调用）在 Requirements 中未提及，Schema 只定义了单锚点 |

---

## Phase 3: Blindspot Hunt

### [blindspot] 1: 单锚点假设

文档中 Anchor Field Schema 为每种 surface 类型定义了单组字段（如 API 的 `endpoint` + `method`）。但实际场景中，一个 Contract 操作可能涉及多个 endpoint（例如：创建资源用 POST，确认用 PUT，回滚用 DELETE）。文档未讨论多锚点场景。

**Quote**: "API | `endpoint` (string)、`method` (string)" — Schema 定义 endpoint 为 string 而非 array，暗示单锚点假设。

### [blindspot] 2: handbook 内容 vs Contract 锚点的粒度不匹配

提案假设 handbook 中的技术细节可以一对一映射到 Contract 的锚点字段。但 api-handbook 按 route 组织（如 `PUT /teams/:teamId/sub-items/:subId/move`），而 Contract 按用户操作/业务场景组织（如"Move sub-item"）。两者的粒度可能不一致——一个业务操作可能跨越多个 route，或一个 route 服务多个操作。

**Quote**: "gen-contracts 从设计文档填充锚点：读取对应 handbook，提取技术细节写入 Contract frontmatter" — "对应" handbook 的"对应"关系如何建立？没有说明。

### [blindspot] 3: eval-contract 改动的测试覆盖

In Scope 包含"eval-contract 评分规则增加技术锚点完整性检查，包含 handbook 内部一致性检查"。但 eval-contract 是评估工具本身，增加新检查项意味着评估标准变更。这会影响所有现有 Contract 的评分——现有没有锚点的 Contract 会因此被扣分。提案未讨论这个向后兼容性问题。

**Quote**: "eval-contract 评分规则增加技术锚点完整性检查" — 新检查对现有无锚点 Contract 的处理方式未说明。

### [blindspot] 4: 交叉验证的分类标准缺失

Solution 和 Scope 反复提到"高置信度/低置信度/无法验证"三类结果，但从未定义分类标准。什么条件判定为高置信度？（Fact Table 和 handbook 完全一致？部分一致？）什么条件判定为低置信度？什么条件判定为无法验证？没有这些定义，实现者无法编写代码。

**Quote**: "将 Fact Table（代码侦察）与 Contract frontmatter 比对，不匹配时输出分类结果（高置信度/低置信度/无法验证）" — 分类标准缺失。

### [blindspot] 5: `consistency_check_result` 来源不明

Success Criteria 部分开头嵌入了一个 YAML 块：

```
consistency_check_result:
  status: pass
  pairs_checked: 15
  conflicts_found: 0
```

这不是一个 Success Criterion，而像是某个自动检查工具的输出。它出现在 SC 列表中会造成混淆——读者不清楚这是提案的承诺还是已完成的验证。如果是为了证明 SC 内部一致性已通过检查，应该放在单独的验证报告而非 SC 部分。

**Quote**: SC 部分开头的 YAML 块。

### [blindspot] 6: "设计文档为准"的 authority assumption 未经审视

提案的核心假设是"设计文档是 authority source，设计-代码不一致时定位为代码 bug"。但在实际开发中，设计文档过时是非常常见的——代码可能已经有意偏离了设计文档（例如：紧急修复、性能优化、需求变更后未更新文档）。将所有不一致都标记为"代码 bug"可能产生大量误报。

**Quote**: "交叉验证以设计文档为 authority source，设计-实现不一致时定位为代码 bug 而非测试问题" — 将设计文档作为绝对权威的假设未充分论证。

### [blindspot] 7: Fact Table 的可靠性假设

提案承认"静态分析无法覆盖所有路由注册模式"，但交叉验证框架仍然依赖 Fact Table 作为验证源之一。当 Fact Table 不完整时（插件系统、动态加载），交叉验证退化为 handbook 单源验证——但提案没有说明这种退化场景下如何保证验证质量。

**Quote**: "代码侦察覆盖度：静态分析无法覆盖所有路由注册模式（插件系统、动态加载、反射机制等），侦察结果不完整时会标记为'低置信度'或'无法验证'" — 退化后的质量保障缺失。

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 78 | 110 |
| 2. Solution Clarity | 88 | 120 |
| 3. Industry Benchmarking | 52 | 120 |
| 4. Requirements Completeness | 72 | 110 |
| 5. Solution Creativity | 55 | 100 |
| 6. Feasibility | 68 | 100 |
| 7. Scope Definition | 52 | 80 |
| 8. Risk Assessment | 62 | 90 |
| 9. Success Criteria | 48 | 80 |
| 10. Logical Consistency | 65 | 90 |
| **Total** | **640** | **1000** |

---

## Bias Detection Report

Annotated regions: 6 attack points / 5 paragraphs = density 1.20
Unannotated regions: 9 attack points / ~20 paragraphs = density 0.45
Ratio (annotated/unannotated): 2.67

**Interpretation**: Annotated regions (pre-revised) show higher attack density, likely because the revisions introduced new content (Known Limitations, Anchor Field Schema, Phased Roadmap, expanded SC, expanded Risks) that expanded the document's surface area but did not always achieve full precision. The pre-revision process correctly identified these as areas needing attention. Several attacks on annotated regions focus on missing precision in newly added details (classification criteria, field schema assumptions), which is consistent with the pre-revision markers flagging these as high/medium severity areas.

**Conflict-with-pre-revision tags**: None. All attacks on annotated regions align with the pre-revision direction — the revised content improved the document but still has gaps that the pre-revision process correctly anticipated.

---

## Key Improvement Recommendations

1. **补充分类标准定义**：明确定义"高置信度/低置信度/无法验证"的判定条件
2. **增加多锚点场景讨论**：修改 Anchor Field Schema 以支持一个 Contract 对应多个技术锚点
3. **深化 Industry Benchmarking**：详细分析 Pact/OpenAPI 的具体做法和适用性，消除 straw-man alternative
4. **补充 SC 覆盖度**：为 eval-contract 改动、handbook 新鲜度检查、`last_anchor_sync` 添加对应 SC
5. **论证"设计文档为 authority"假设**：提供退化和冲突场景的处理策略，而非绝对标记为代码 bug
6. **量化证据**：补充问题发生频率数据、受影响项目范围、以支持 urgency 论证
7. **移除或解释 `consistency_check_result` YAML 块**：不应出现在 SC 列表中
