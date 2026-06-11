---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision Eval Report (Iteration 0)

## ATTACK_POINTS

- **[high]** PascalCase vs camelCase 命名矛盾 | quote: "YAML 使用 camelCase 但与 Go struct 的 PascalCase 字段名不同" | improvement: 统一分组内字段名为 PascalCase，与 Go struct FieldByName 大小写敏感匹配保持一致，消除 Key Risks 和 SC-FM-5 中的表述矛盾
- **[high]** PhaseSummary 格式与 Out of Scope 自相矛盾 | quote: "不修改 `prompt.go` 渲染逻辑（`Synthesize()` 等）" 为 Out of Scope，但 PhaseSummary section 需要 `phaseSummaryLine` 格式变更 | improvement: 将 `phaseSummaryLine` 的受控格式修改纳入 In Scope，明确修改范围为仅去除 `PHASE_SUMMARY:` 前缀
- **[high]** validateMetadataVariables 校验未覆盖 task/record 模板 | quote: "校验所有分组中的字段名集合均存在于对应的 Go struct 中" 但未定义 task/record 使用哪个 struct | improvement: 增加"分组规则表"定义三类模板各自的校验 struct 和分组分配
- **[high]** 向后兼容性的语义映射未定义 | quote: "TemplateMetadata.Variables 列表应该只包含 variables 分组下的字段，还是应该包含所有分组字段的并集？" | improvement: 定义新格式 TemplateMetadata 结构体和 AllFields() 语义等价规范
- **[high]** TASK_FILE 行格式稳定性未列为约束 | quote: "strings.Replace 将静默失败——不会报错" | improvement: 在 Constraints 中增加 TASK_FILE/TASK_ID/SURFACE_KEY 行格式不变性约束
- **[high]** 功能快照清单粒度标准未定义 | quote: "清单中每个节点的粒度是什么" | improvement: 在 Constraints 中增加快照清单的创建标准和分类字典
- **[high]** 行级 YAML 解析器扩展复杂度被低估 | quote: "扩展分组支持复杂度低" 但实际需要状态跟踪 | improvement: 在 Feasibility 中如实评估解析器扩展复杂度，建议引入 gopkg.in/yaml.v3
- **[medium]** 分组层级判定规则缺失 | quote: "为什么 task 模板不需要 conditional 分组？" | improvement: 增加分组的判定规则定义：conditional = 在正文中以 {{if .X}} 控制段落显示的字段
- **[medium]** SC-FM-1 迁移检测逻辑矛盾 | quote: "包括那些确实无需 identity 字段的模板" | improvement: 按模板类型分别定义迁移检测标准
- **[medium]** task-executor 步骤合并操作定义不充分 | quote: "从 11 步到 <=8 步需要减少 3 步" | improvement: 补充具体步骤合并方案

## BORDERLINE_FINDINGS

- PhaseSummary 语义冗余（标题行和旧前缀重复）→ 标记为 borderline，随 PhaseSummary 格式修改一并解决

## SKIPPED_FINDINGS

- 无（所有建议性 finding 已纳入上述 ATTACK_POINTS）
