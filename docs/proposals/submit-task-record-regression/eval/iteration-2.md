# Evaluation Report: Iteration 2

**Evaluator:** CTO Persona (Adversarial Scoring)
**Date:** 2026-05-24
**Iteration:** 2
**Document:** `docs/proposals/submit-task-record-regression/proposal.md`
**Previous:** Iteration 1 — 770/1000

---

## Total Score: 835/1000

---

## Dimension Scores

### 1. Problem Definition: 100/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 38/40 | "Go CLI 的 record 渲染管线（validateRecordData + markdown 模板渲染）接受 RecordData JSON 输入并输出 records/*.md。但当前缺乏系统性验证手段——无法确认每种 task type 的历史记录能否通过当前 Go CLI 的校验与渲染，输出与已提交 markdown 一致的结果" (line 13)。范围明确，指向具体函数和具体输出格式。标题也已重命名为 "Go CLI Record 渲染管线的回归验证（按任务类型）"，与内容完全匹配。唯一残留微小模糊：标题用 "渲染管线" 而非精确的函数名，但读者从第一段即可获得精确范围。 |
| Evidence provided | 36/40 | 四条 evidence 均为可验证事实。关键改进：evidence #3 从 "可能已与当前管线产生偏差" 修订为 "已与当前管线产生偏差——例如 forge-cli-clean-code 的 `coding.fix` 记录使用旧的 `summary` 字段格式，而当前 Go struct 已重命名为 `resolution`" (line 20)。这是具体的、可验证的偏差实例，不再是推测。扣分项：此偏差是前向不兼容（历史有当前无），但 evidence 未明确区分前向/后向两种偏差类型，读者可能误以为 golden dataset 能检测全部偏差。 |
| Urgency justified | 26/30 | "随着 task type 持续增加（当前 12 种活跃类型），Go CLI 渲染管线对各类 record 的兼容性风险逐步累积。现在用已有的丰富历史数据建立回归测试，成本最低——数据已存在，只需提取为 fixture" (line 26)。逻辑完整，时间窗口（数据已存在）论证有力。改进空间：未量化不做会怎样——一个 task type 渲染失败的实际用户影响是什么？修复成本多高？缺乏量化使紧迫感依赖定性判断。 |

**Attacks:**
1. [Problem Definition]: Evidence #3 给出了 `summary` → `resolution` 的具体实例，但未区分这是前向不兼容（历史有当前无）还是后向遗漏（当前有历史无）——读者可能从 evidence 的措辞推断 golden dataset 能检测所有 schema 偏差，但 Requirements line 46 已承认后向遗漏无法自动检测。Evidence 与 Requirements 之间有轻微的不对称。

---

### 2. Solution Clarity: 108/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 40/40 | "从已完成的 feature 提取历史 record 作为 golden dataset，按 task type 分组建立 Go CLI 层的回归测试：validateRecordData 校验 + markdown 模板渲染对比。测试文件位于 `internal/record/golden_test.go`，fixture 以 `testdata/{taskType}/{featureName}_{index}.json` + `testdata/{taskType}/{featureName}_{index}.golden.md` 的命名约定组织。table-driven 测试使用 `testCase` struct，字段包括 `name string`、`taskType string`、`inputJSON string`（fixture 路径）、`goldenMD string`（golden 路径）" (line 30)。读者可以完整复述将建什么——测试文件位置、fixture 命名、数据结构、测试模式全部明确。 |
| User-facing behavior described | 38/45 | Scope note (line 82) 明确排除了 submit-task skill 层，使测试边界清晰。开发者体验显著改善——fixture 命名约定、struct 定义、目录结构全部给出。扣分项：`RenderRecord()` 的调用方式仍未说明——是直接 Go 函数调用还是 CLI subprocess？这影响 fixture 格式设计（函数调用需要 JSON 输入，subprocess 需要考虑 CLI 参数和 stderr 输出）。 |
| Technical direction clear | 30/35 | `validateRecordData()`、`RenderRecord()`、table-driven by task type、`-update` flag、fixture 文件格式。方向清晰。微小不足：未说明 `-update` flag 的实现机制——是 Go 标准库的 `flag.Bool("update", ...)` 还是自定义逻辑？这影响 CI diff gating 的实现复杂度。 |

**Attacks:**
2. [Solution Clarity]: `RenderRecord()` 的调用方式未说明——直接函数调用 vs CLI subprocess 的选择会影响 fixture 设计和 CI 集成方式，但文档未做出明确选择 (line 79 仅说 "可直接在测试中调用" 但未排除 subprocess 方式)

---

### 3. Industry Benchmarking: 72/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 24/40 | 仍只引用 Go `testing` 包的 `-update` flag 和 `jsonschema`。修订后 scope note 使方案聚焦于 Go CLI 层 golden testing，但 benchmarking 深度未改善。对于以 golden/snapshot testing 为核心的方案，缺少具体项目引用（Hugo、Terraform、Protoc 的 golden test 实现）、文章链接、或 Go 社区最佳实践文档。这是全文档最大短板，且 iteration 2 未改善。 |
| At least 3 meaningful alternatives | 20/30 | 三个方案（Do nothing、JSON Schema、Golden dataset）。修订后 JSON Schema 的 dismiss 从 "偏离 golden dataset 对比的目标" 改为 "不覆盖渲染管线一致性" (line 72)，不再是循环论证，但仍不够正面——未承认 JSON Schema 在 schema drift 检测上的优势（恰好是 line 46 承认的 golden dataset 的盲区）。缺少 property-based testing、template unit testing 等替代方案。 |
| Honest trade-off comparison | 14/25 | 修订后 cons 显著改善："fixture 需随格式演进更新；无法检测 record-format 模板文档与 Go struct 的偏差（只能测已有记录的回归，不能发现文档描述与实现的不一致）；历史 markdown 与当前渲染不一致时需人为判断以哪方为准" (line 73)。三个真实代价全部承认。扣分：JSON Schema 的 cons "需维护两套 schema" 是准确的，但 "不覆盖渲染管线一致性" 是用 solution 的优势去评判替代方案——这是隐含的循环论证。 |
| Chosen approach justified against benchmarks | 14/25 | "与目标最匹配" (line 73) 仍是论断。修订后 scope note 限定了目标为 Go CLI 渲染管线，golden dataset 确实匹配此目标，但缺少与 JSON Schema 的正面对比。文档本可论证："Golden dataset 覆盖了 validateRecordData + RenderRecord 的端到端一致性，而 JSON Schema 只能覆盖 validateRecordData 的结构合规——因此 golden dataset 在渲染一致性上更直接，但 JSON Schema 在 schema drift 检测上更强"。这种分层对比缺失。 |

**Attacks:**
3. [Industry Benchmarking]: 无具体项目引用（Hugo/Terraform/Protoc 的 golden test 实现）、无文章链接、无 Go 社区最佳实践文档——对于以 golden testing 为核心方法的方案，benchmarking 深度不足 (section "Industry Solutions" line 65 仅一段泛泛描述)
4. [Industry Benchmarking]: JSON Schema 替代方案的 dismiss 仍使用目标导向论证——"不覆盖渲染管线一致性" (line 72)——未正面承认 JSON Schema 恰好能覆盖 golden dataset 无法覆盖的盲区（schema drift 检测，如 line 46 所承认），两者应为互补关系而非互斥关系

---

### 4. Requirements Completeness: 94/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 36/40 | 修订后场景精准。关键改进："Schema 偏差（前向不兼容）" 限定为 "历史 RecordData JSON 包含当前 Go RecordData struct 已删除或重命名的字段，导致校验或渲染失败" (line 44)——可测试。"缺失字段（后向遗漏）" 限定为 "当前 Go 端新增的 required 字段在历史 record 中不存在——golden dataset 无法自动检测此类偏差，需通过 schema 对比人工补充" (line 46)——明确承认了覆盖盲区。扣分：缺少数值边界场景（如 RecordData 中字段值极长、特殊字符、空值）的考虑——这些不是 schema 偏差而是数据质量边界。 |
| Non-functional requirements | 32/40 | 修订后 CI 性能预算有详细分解："~30 fixtures × 2 次函数调用/fixture ≈ ~60 次函数调用，无网络 I/O，文件 I/O 极小（仅读取 fixture JSON 和 golden MD），预计 <10s，预留 buffer 至 30s" (line 52)。修正了 "无 I/O" 为 "无网络 I/O，文件 I/O 极小"。扣分：仍缺 fixture 存储大小估算、Go 版本要求。 |
| Constraints & dependencies | 26/30 | 两条约束清晰且无歧义。"不修改 Go 端校验逻辑（只测试，不改行为）" (line 59) 是关键约束。Scope note (line 82) 补充了隐含约束：submit-task skill 层不在自动化测试范围内。 |

---

### 5. Solution Creativity: 40/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 15/40 | 文档坦诚声明 "Golden dataset 回归测试是标准的软件工程实践，无特别创新" (line 34)。修订后亮点改为 "按 task type 分组组织 golden dataset，直接验证 Go CLI 渲染管线的端到端一致性" (line 35)。按 task type 分组是 table-driven testing 的标准组织方式。端到端声明已被 scope note 限定为 Go CLI 内部三步管道，不再是全链路。无超越行业基线的创新。 |
| Cross-domain inspiration | 5/35 | 无。纯标准 Go testing 实践，未借鉴其他领域的思路。 |
| Simplicity of insight | 20/25 | 核心洞察优雅："数据已存在，只需提取为 fixture"。修订后 scope note 使洞察更加聚焦——不试图解决无法自动化测试的 submit-task skill 层问题，而是将已有的历史数据直接转化为测试资产。这是 "利用现有资产" 的简洁洞察。 |

---

### 6. Feasibility: 95/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 40/40 | 高度可行。函数已存在、数据已存在、Go testing 框架成熟。Scope note 明确排除了 submit-task skill 层的复杂性。fixture 命名约定和 struct 定义已给出，技术路径无歧义。 |
| Resource & timeline feasibility | 27/30 | "预计 5-8 个 coding task" (line 86) 合理但偏宽。Scope note 将范围缩小后，实际工作量可能更接近 4-6 tasks。不过 5-8 包含了偏差分析和 fixture 调整的工作量，属于合理预留。 |
| Dependency readiness | 28/30 | "无外部依赖。所有数据在本地仓库" (line 90)。准确。微小风险：历史 record 的 JSON 格式是否完整——有些 feature 的 record.json 可能不包含完整的 RecordData 字段，但这属于数据质量问题而非依赖问题。 |

---

### 7. Scope Definition: 76/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 30/30 | 四项具体可交付物。修订后第 3 条从 "发现 record-format 模板中的问题并修复" 改为 "发现 Go 端 RecordData schema 与历史格式的偏差并更新 golden fixture 以匹配当前 Go 行为（若发现 Go struct 需兼容调整则单独提 issue，不在本 feature scope 内）" (line 143)。这消除了 scope creep 风险。"覆盖 11 种独立 task type + 1 个 alias（`fix` 为 `coding.fix` 的 legacy alias，共享同一模板路径）" (line 145) 精确描述了覆盖范围。 |
| Out-of-scope explicitly listed | 22/25 | 四项明确排除。Scope note (line 82) 补充排除了 submit-task skill 层。扣分：未明确排除 "Go 端代码修改"——虽然 line 59 说 "不修改 Go 端校验逻辑"，但 In Scope 第 3 条的 "偏差" 暗示可能修改 Go struct（虽然后面括号排除了），表述可以更清晰。 |
| Scope is bounded | 24/25 | 量化边界：11 种独立 type + 1 alias，10 个 feature source，5-8 coding tasks。Scope note 使范围与实现一致。微小问题：30 个 fixture vs "每种选 2-3 条" × 12 种 = 24-36，范围合理但跨度较大。 |

---

### 8. Risk Assessment: 78/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 26/30 | 五项风险。修订后新增 "fixture 代表性不足，遗漏渲染边界情况" (line 162)，补充了 iteration-1 指出的缺失风险。legacy `fix` risk 修正为 "共享同一模板路径；fixture 中标注 alias 关系即可" (line 163)。扣分：仍缺少历史记录数据质量问题风险（JSON 格式不完整、字段缺失等非 schema 偏差问题）——这不是 fixture 代表性问题，而是原始数据质量问题。 |
| Likelihood + impact rated | 26/30 | 使用 M/L 矩阵一致。fix/coding.fix 混淆风险从 M/L 降至 L/L (line 163) 合理。扣分："fixture 数量大导致测试维护成本高" L/L (line 161) 可能低估——当 Go 模板频繁迭代时（iteration-1 也指出），每次迭代需审查全部 ~30 个 fixture 的 diff，实际维护成本取决于模板迭代频率。 |
| Mitigations are actionable | 26/30 | 大幅改善。Diff gating 机制描述具体："CI 中 diff gating 通过 Go test 的 `-update` flag 实现：测试失败时输出 diff，开发者传入 `-update` 重新生成 golden 文件，CI 检测 `.golden` 文件 git diff 非空即标记失败" (line 159)。这是可操作的具体机制。fixture 代表性的 mitigation "优先选取不同 feature 来源的记录以增加多样性；发现的边界 case 可增量补充 fixture" (line 162) 也具操作性。扣分：diff gating 的具体 CI 实现方式未说明——是 GitHub Actions 的 `git diff --exit-code`？还是自定义脚本？ |

---

### 9. Success Criteria: 70/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 48/55 | 修订后 criteria 有实质性改善。"≥11 种独立 task type + fix alias 的 golden dataset fixture 建立完成" (line 168)——计数与 scope 一致。"Go 端 RecordData schema 与历史记录的偏差全部识别并记录，golden fixture 更新至与当前 Go 行为一致（发现的偏差逐一记录在测试输出的诊断报告中）" (line 170)——"全部识别" 仍不可测量（不知总数），但 "逐一记录在诊断报告中" 使过程可审计。"新增 fixture 流程可操作：从已通过校验的历史 record 提取 RecordData JSON + 对应 markdown 为 golden pair，在 table-driven 测试中新增一条用例" (line 171)——可验证。扣分："全部识别" 在问题数量未知时仍不可测量——无法定义 "全部" 的终止条件。 |
| Coverage is complete | 22/25 | 五条 criteria 覆盖了 fixture 建立、测试执行、偏差识别、可扩展性、类型覆盖。与 in-scope items 对齐。 |

---

### 10. Logical Consistency: 82/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 32/35 | 修订后 Problem 与 Solution 精准匹配。Problem 说 Go CLI 渲染管线缺乏验证 (line 13)，Solution 是 Go CLI 回归测试 (line 30)。Scope note (line 82) 明确排除了 submit-task skill 层。逻辑闭环完整。唯一残留："端到端一致性" (line 35) 描述三步管道（validate → render → compare）仍稍显过度——更准确的措辞是 "渲染管线一致性" 或 "管线回归"。 |
| Scope <-> Solution <-> Success Criteria aligned | 28/30 | 大幅改善。Scope 排除 submit-task skill → Solution 只测 Go CLI → Criteria 不再声称验证 submit-task 模板选择。"≥11 种独立 task type + fix alias" (line 168) 与 Scope 的 "11 种独立 task type + 1 个 alias" (line 145) 一致。扣分：Innovation Highlights 的 "端到端一致性" (line 35) 与 Scope note 限定的 Go CLI 内部管道之间有轻微的修辞过度——不影响执行但影响沟通精确度。 |
| Requirements <-> Solution coherent | 22/25 | 修订后场景与 Solution 匹配度良好。"Schema 偏差" 限定为前向不兼容 (line 44)——golden dataset 可检测。"缺失字段" 限定为后向遗漏并承认 "golden dataset 无法自动检测此类偏差，需通过 schema 对比人工补充" (line 46)——不再声称自动化。扣分：line 46 承认后向遗漏无法自动检测，但 In Scope 列表未包含 "schema 对比人工补充" 这项工作——它被隐含在 In Scope 第 3 条的 "偏差" 中，但 "人工补充" 不如其他 in-scope 项那样具体。 |

---

## Blindspots

1. **[blindspot]** Coverage Matrix 中 `doc.eval` 和 `doc.drift` 引用的 Feature Sources 不在 Feature Sources 表中。`doc.eval` 列出 "eval-freeform-expert, run-tasks-git-status, enforce-forge-task-add, run-tests-decouple" (line 132)，`doc.drift` 列出 "worktree-unpushed, auto-task-main, forge-research, list-tasks, cli-restructure, refactor-impact" (line 134)。这些名称与 Feature Sources 表的 10 个 feature 不匹配。这是自 baseline 起连续三个版本未解决的问题，使这两种 task type 的 coverage 数据无法验证。考虑到 doc.eval 有 5 条记录、doc.drift 有 6 条记录，这些 fixture 的数据来源不明确会直接影响测试实施。

2. **[blindspot]** Feature Sources 表中 `fix` 和 `coding.fix` 的记录可能重复计数。forge-cli-clean-code 列出 "coding.fix (6), fix (6)" (line 108)——如果这些是同一批记录的两种计数（`fix` 是 `coding.fix` 的 alias），则 coding.fix 实际只有 6 条而非 Feature Sources 表暗示的 6+9=15 条。Coverage Matrix 中 coding.fix 列出 15 条 (line 127) 但 Feature Sources 仅引用 "forge-cli-clean-code, forge-arch"，而 forge-cli-clean-code 的 6 条 `coding.fix` 和 6 条 `fix` 是否计入 coding.fix 的 15 条未说明。自 baseline 起未解决。

3. **[blindspot]** Assumptions Challenged 表格第 2 行 "Go 端 RecordData schema 与历史记录格式兼容 | Assumption Flip | 待验证：这是本次 feature 的核心目标" (line 98)。修订后 scope note 将测试限定为 Go CLI 渲染管线回归测试，但此 assumption 仍声称 "验证" schema 兼容性。Golden dataset 测试通过只证明 "历史记录能通过当前的校验和渲染"，不证明 "当前 schema 与所有可能的合法输入兼容"——一个通过所有历史 fixture 的管线可能仍在新字段、新 edge case 上出错。Assumption 的措辞应从 "待验证" 改为 "部分验证（仅覆盖历史数据中出现的字段组合）"。

4. **[blindspot]** Industry Benchmarking section 在连续三个版本中几乎未改善，始终停留在 "Go `testing` 包内置 `-update` flag" 的泛泛提及。考虑到此方案以 golden testing 为核心方法，benchmarking 深度不足（缺少具体项目引用、文章链接、社区最佳实践）已成为系统性短板。这不是措辞问题而是调研深度问题——一个 CTO 审批时可能期望看到 "Hugo 用 golden tests 验证 template rendering output（链接），我们借鉴其 fixture 组织方式" 级别的引用。

---

## Bias Detection Report

Annotated regions (marked with `<!-- pre-revised: {severity} -->`): 13 paragraphs/regions
Unannotated regions: ~18 paragraphs/regions

- Annotated regions: 2 attack points / 13 paragraphs = density 0.15
  - Attack #1 (evidence 前向/后向不对称), Attack from blindspot #3 (assumption "待验证")
- Unannotated regions: 6 attack points / 18 paragraphs = density 0.33
  - Attack #2 (RenderRecord 调用方式), #3 (无项目引用), #4 (JSON Schema 互补关系), blindspot #1 (doc.eval/drift 数据来源), blindspot #2 (fix/coding.fix 重复计数), blindspot #4 (benchmarking 系统性短板)

Ratio (annotated/unannotated): 0.45

**Interpretation**: Annotated (pre-revised) regions receive fewer attacks than unannotated regions, with a ratio of 0.45. This is lower than iteration-1's ratio of 0.58, indicating a slightly stronger bias toward leniency on revised regions. Contributing factors: (1) revised regions genuinely addressed the most critical issues from baseline and iteration-1, leaving fewer substantive attacks; (2) unrevised regions (Industry Benchmarking, Coverage Matrix data) carry forward systemic weaknesses from baseline that accumulate attacks across iterations. The bias is not severe enough to invalidate the scoring — the unrevised sections objectively have more issues.

No `conflict-with-pre-revision` tags generated — no scorer judgment contradicts the pre-revision direction.

---

## Comparison to Previous Iterations

| Dimension | Baseline (i0) | Iteration 1 | Iteration 2 | Delta (i1→i2) |
|-----------|--------------|-------------|-------------|---------------|
| Problem Definition | 84 | 92 | 100 | +8 |
| Solution Clarity | 92 | 98 | 108 | +10 |
| Industry Benchmarking | 64 | 68 | 72 | +4 |
| Requirements Completeness | 76 | 86 | 94 | +8 |
| Solution Creativity | 28 | 35 | 40 | +5 |
| Feasibility | 88 | 92 | 95 | +3 |
| Scope Definition | 64 | 72 | 76 | +4 |
| Risk Assessment | 58 | 72 | 78 | +6 |
| Success Criteria | 46 | 62 | 70 | +8 |
| Logical Consistency | 43 | 75 | 82 | +7 |
| **Total** | **643** | **770** | **835** | **+65** |

Key improvements from iteration 1 to 2:
1. **Solution Clarity +10**: 添加了测试文件路径、fixture 命名约定、testCase struct 定义，使开发者可操作性达到接近完整的水平
2. **Problem Definition +8**: Evidence #3 从 "可能" 改为具体偏差实例（summary→resolution），消除推测性
3. **Requirements Completeness +8**: 后向遗漏盲区明确承认、CI 预算详细分解、NFR I/O 措辞修正
4. **Logical Consistency +7**: "≥11 种独立 task type + fix alias" 与 scope 对齐、phantom scenarios 限定为可测试范围
5. **Success Criteria +8**: 计数与 scope 一致、fixture 流程描述更准确

Remaining gaps (prioritized):
1. **Industry Benchmarking (72/120)**: 最大短板且改善最慢（+4 per iteration）。缺少具体项目引用、文章链接、社区最佳实践。JSON Schema 替代方案的评估仍不够正面
2. **Solution Creativity (40/100)**: 诚实承认无创新，分数合理但低于大多数维度
3. **Coverage Matrix 数据不可验证**: doc.eval 和 doc.drift 引用不在 Feature Sources 表中的 feature 名称——连续三个版本未解决
4. **fix/coding.fix 记录可能重复计数**: Feature Sources 表和 Coverage Matrix 之间的一致性问题——连续三个版本未解决

---

## Attacks Summary

1. [Problem Definition]: Evidence #3 给出具体偏差实例但未区分前向/后向偏差类型——"历史记录格式已与当前管线产生偏差" (line 20)——读者可能误以为 golden dataset 能检测所有偏差，但 line 46 已承认后向遗漏无法自动检测，两者之间存在轻微不对称
2. [Solution Clarity]: `RenderRecord()` 的调用方式未说明——直接函数调用 vs CLI subprocess 影响 fixture 设计和 CI 集成——line 79 仅说 "可直接在测试中调用" 但未排除 subprocess 方式
3. [Industry Benchmarking]: 无具体项目引用（Hugo/Terraform/Protoc 的 golden test 实现）、无文章链接、无 Go 社区最佳实践——Industry Solutions section (line 65) 仅一段泛泛描述，对于以 golden testing 为核心方法的方案，调研深度不足
4. [Industry Benchmarking]: JSON Schema 替代方案与 golden dataset 应为互补关系而非互斥——"不覆盖渲染管线一致性" (line 72) 未承认 JSON Schema 恰好能覆盖 line 46 承认的 golden dataset 盲区
5. [Requirements Completeness]: 缺少数值边界场景（字段值极长、特殊字符、空值等数据质量边界）
6. [Risk Assessment]: 缺少历史记录数据质量问题风险——JSON 格式不完整、字段缺失等非 schema 偏差问题
7. [Success Criteria]: "偏差全部识别并记录" (line 170) 中 "全部" 在问题数量未知时不可测量——无法定义终止条件
8. [Logical Consistency]: "端到端一致性" (line 35) 对三步管道（validate → render → compare）仍属修辞过度，更准确的措辞是 "渲染管线一致性"
9. [blindspot]: Coverage Matrix 中 doc.eval 和 doc.drift 引用的 feature 不在 Feature Sources 表中 (line 132-134)——连续三个版本未解决
10. [blindspot]: Feature Sources 表 fix/coding.fix 记录可能重复计数 (line 108)——连续三个版本未解决
11. [blindspot]: Assumption "Go 端 RecordData schema 与历史记录格式兼容...待验证" (line 98) 措辞应限定为 "部分验证（仅覆盖历史数据中出现的字段组合）"
12. [blindspot]: Industry Benchmarking 在连续三个版本中几乎未改善，已成为系统性短板——缺少调研深度而非措辞问题
