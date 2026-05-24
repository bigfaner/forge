---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision (Freeform Findings)

## ATTACK_POINTS

### Factual Corrections

- **[high]** Proposal误标识修改面：`GetReviewDocTask()`不是实际修改目标 | quote: "修改 `autogen.go`：review-doc 任务定义支持嵌入 AC 数据" | improvement: 修正修改面描述为 BodyContext/renderBody 及模板文件

- **[high]** Prompt模板修改缺乏规格说明：最高风险变更仅有单行描述 | quote: "修改 `prompt/data/doc-review.md`：简化 agent prompt——移除任务扫描指令，改用内嵌 AC，加入文档过滤规则，禁止修改任务文件和记录" | improvement: 提供prompt模板mockup和重构步骤

- **[high]** AC汇总区域格式未定义，阻塞实现 | quote: "autogen 模板增加 AC 汇总区域，更新 Discovery Strategy 排除 tasks/ 和 records/" | improvement: 定义AC聚合的markdown格式schema（含task名称映射）

### Structural/Architectural Suggestions

- **[high]** AC提取设计与现有BodyContext schema不匹配，缺少per-task AC字段 | quote: "修改 `build.go`：`BuildIndex()` 从 doc 任务文件提取 AC 并传入 review-doc 生成" | improvement: 描述BodyContext需要的schema变更（新字段或新template placeholder）

- **[medium]** 两个行为锚点（autogen模板与prompt模板）的耦合未被承认 | quote: "agent prompt 的文档发现从'扫描全部'改为'只扫描内容文档'，明确排除 `tasks/`、`tasks/records/`、`manifest.md`、`index.json`" | improvement: 明确声明两个模板必须同步修改的耦合关系

- **[medium]** Proposal自身逻辑矛盾：否定prompt约束却依赖prompt过滤 | quote: "纯 prompt 约束被否决。Reason: LLM 经常忽略 prompt 约束，问题根源需从信息流层面解决" | improvement: 承认过滤策略仍属prompt约束，改为allowlist增强可靠性，或描述为何此处prompt约束可接受

- **[medium]** 实现工作量被低估，实际至少需4-5个任务而非2-3个 | quote: "预估 2-3 个 coding 任务" | improvement: 修正工作量估算为4-5个任务（extraction pipeline, BodyContext schema, autogen template, prompt template restructuring, testing）

- **[medium]** BuildIndex与T-review-doc执行间存在时间耦合，AC可能过期 | quote: "review-doc 依赖最后业务任务完成后才执行，手动修改发生在执行前，构建时状态即为最新" | improvement: 承认时间耦合风险，添加重新运行forge task index的约束说明

- **[medium]** AC为空时缺少行为定义，导致静默降级路径 | quote: "若 section 不存在则标记为空" | improvement: 定义AC为空时的具体行为（报错、跳过、还是自由评审）

- **[medium]** AC提取对标题格式严格依赖，格式偏差导致静默失败 | quote: "AC 覆盖率不降低（所有 doc 任务的 AC 均被检查）" | improvement: 添加构建时验证警告，定义标题匹配的容错策略

- **[medium]** 过滤策略的匹配方式未定义，存在边缘情况正确性风险 | quote: "过滤基于明确路径模式（tasks/、records/、manifest.md、index.json），不影响 docs/ 下的内容文件" | improvement: 明确指定使用allowlist而非denylist，定义路径匹配方式

## BORDERLINE_FINDINGS

- **[high]** Prompt模板修改是最高风险变更却获得最少规格说明 — 此发现介于factual correction（文档结构事实）和structural suggestion之间，归类为borderline，由Reviser决定是否处理

## SKIPPED_FINDINGS (Subjective Preferences)

- 建议：定义具体的AC聚合schema再实现 (severity: low)
- 建议：构建时增加AC提取验证步骤 (severity: low)
- 建议：承认时间耦合，改为执行时提取AC (severity: low)
- 建议：将prompt模板修改拆为独立实现任务 (severity: low)
- 建议：用allowlist替代denylist进行文档发现过滤 (severity: low)
- 建议：为空或异常AC定义降级行为路径 (severity: low)

## Classification Audit

- Total findings by triage layer: 3 factual correction / 8 structural suggestion / 6 subjective preference / 1 borderline
- 11 findings triaged (3 accepted + 8 structural) / 6 skipped / 1 borderline

## rubric

(all dimensions): N/A
