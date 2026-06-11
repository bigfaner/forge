# Pre-Revision Eval Report (Iteration 0)

**Title**: Pre-Revision (Freeform Findings)
**Iteration**: 0
**Source**: Freeform expert review by Prompt Compliance Architect

## ATTACK_POINTS

### Factual Corrections

- **[high]** complexity 字段注入需改动完整数据管道，提案严重低估改动面 | quote: "`renderTemplate()` 已有 7 个占位符，加一个 `{{COMPLEXITY}}` 是最小扩展。" | improvement: 更新 Feasibility 和 Constraints 部分，明确列出完整数据管道改动链：FrontmatterData → Task struct → index.json schema → renderTemplate

- **[high]** strings.ReplaceAll 模板系统无法实现条件性跳过段落，提案未讨论此架构约束 | quote: "low 跳过 Step 1.5 spec-code scan、简化探索" | improvement: 在 Constraints 中说明模板引擎限制，并给出条件段落实现方案（如 cleanTemplateOutput() 扩展或注释标记机制）

- **[high]** `<CRITICAL>` 与 `CODING_PRINCIPLES` 优先级矛盾未被直接解决，仅通过信息隔离间接规避 | quote: "CRITICAL 优先级高于一般原则，executor 遵循了冲突扫描结果而非 scope 边界。" | improvement: 在 Solution 中增加显式 scope boundary 声明机制，直接覆盖 spec-code scan 的越界倾向，而非仅靠信息隔离

- **[medium]** "搜索策略引导"概念模糊，在 5 个 template 中的具体位置和行为未定义 | quote: "加'先收集后修改'搜索策略引导" 以及 "搜索策略引导出现在所有 5 个 coding template 的 implementation 步骤前" | improvement: 给出搜索策略引导的具体指令内容模板，说明与 Step 1.5 的层叠关系

- **[medium]** "探索阶段 < 30s" 不可靠验证，forge 当前无 step-by-step 执行时间记录 | quote: "complexity: low 的任务执行时跳过 Step 1.5 spec-code scan，探索阶段 < 30s" | improvement: 将此 SC 改为可通过 `forge prompt get-by-task-id` 输出验证的确定性条件（如"prompt 输出不包含 Step 1.5 段落"），将 < 30s 降级为 NFR

- **[medium]** coding.fix 类型的 5-step 流程与 enhancement 不同，应用相同 complexity 判定逻辑语义不匹配 | quote: "5 个 coding prompt templates 加复杂度分支" | improvement: 明确 fix 类型不纳入 complexity routing 或做差异化处理，说明 fix-task 由 dispatcher 自动生成时不带 complexity 字段

### Structural Suggestions

- **[medium]** complexity 判定启发式依赖静态数量指标，忽略任务实际认知复杂度 | quote: "AC 数量 + Hard Rules + Reference Files 数量作为判定依据" | improvement: 考虑将硬编码阈值改为 SKILL.md 中的 LLM 判断指引，让任务生成阶段的 LLM 做认知判断

- **[medium]** Reference Files 内联化引入 stale reference 风险，未评估与 scope creep 的权衡 | quote: "proposal 变更时 task doc 不会自动同步。但对于 quick mode（≤15 任务），这个代价可接受。" | improvement: 增加 stale reference 缓解措施（如两层结构：内联 + 溯源链接），或在 Risks 中明确评估权衡

- **[medium]** quick-tasks 和 breakdown-tasks 同步修改缺少精确对应关系 | quote: "两个 SKILL.md 使用相同的判定规则描述，确保逻辑一致" | improvement: 在 In Scope 中分别列出两个文件各自需要修改的段落

- **[medium]** 移除 15 coding task 上限依据不充分，是与精度控制无关的独立架构决策 | quote: "移除 quick-tasks 15 coding task 上限" | improvement: 考虑将此项移出 scope 或补充充分的依据说明

## BORDERLINE_FINDINGS

无

## SKIPPED_FINDINGS (Subjective Preferences)

- 建议用 cleanTemplateOutput() 条件段落机制替代 {{COMPLEXITY}} 简单占位符 — 具体实现建议，留给 Scorer 评估
- 建议将 complexity 判定从硬编码阈值改为 SKILL.md 中的 LLM 判断指引 — 替代方案建议，已部分体现在 structural suggestion 中
- 建议 Reference Files 采用内联+溯源链接双层结构并加 scope boundary 声明 — 具体实现建议，已部分体现在 structural suggestion 中
- 建议将移除上限从本提案 scope 移出 — 范围调整建议，已体现在 structural suggestion 中
- 建议 Step 1.5 被跳过替代 < 30s 作为 SC — 已体现在 factual correction 中
- 建议 coding.fix 差异化处理 — 已体现在 factual correction 中

## Classification Audit

Total findings: 16
- Factual corrections: 6
- Structural/architectural suggestions: 4
- Subjective preferences (skipped): 6

## Rubric

All dimensions: N/A (pre-revision phase, rubric scoring deferred to Scorer cycle)
