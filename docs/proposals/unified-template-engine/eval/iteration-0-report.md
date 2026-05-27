---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision (Freeform Findings)

## ATTACK_POINTS

- **[high]** promptTemplateData struct 缺少 4 个必要字段 | quote: "type promptTemplateData struct { TaskID, TaskCategory, PhaseSummary, CoverageStrategy, SurfaceKey, SurfaceType, Complexity }" — 缺少 TaskFile、FeatureSlug、TestTypeArg、CoverageTarget。`renderTemplate()` 在 prompt.go:114-168 中消费 8 个独立替换操作，但 struct 只有 7 个字段。 | improvement: 补全 promptTemplateData 为 11 个字段（增加 TaskFile、FeatureSlug、TestTypeArg、CoverageTarget）

- **[high]** autogenTemplateData struct 严重不完整，缺失 5/9 数据维度 | quote: "type autogenTemplateData struct { TaskID, TaskType, SurfaceKey, SurfaceType, ScopeDisplay }" — renderBody() 实际消费 9 个数据维度（FeatureSlug、Mode、ScopeDisplay/SurfaceTypes、SurfaceKey、AcceptanceCriteria、DocTaskCriteria 等），但 struct 只定义了 5 个字段。 | improvement: 补全 autogenTemplateData 为完整字段集（增加 FeatureSlug、Mode、SurfaceTypes、AcceptanceCriteria、DocTaskCriteria 的预格式化字段）

- **[high]** PHASE_SUMMARY 需要两个独立的 {{if}} 块而非一个 | quote: "`{{if .PhaseSummary}}...{{end}}` 替换 `If {{PHASE_SUMMARY}} is non-empty` 模式" — 所有 18 个使用 PHASE_SUMMARY 的模板中，占位符出现在两个位置（标签行 + 条件指令行），相隔 10-30 行，不能用单个 {{if}} 包裹。 | improvement: 明确说明 PHASE_SUMMARY 迁移需要两个独立 {{if}} 块，并给出具体模板示例

- **[medium]** `just compile {{SURFACE_KEY}}` 尾部空格处理需要模板级方案 | quote: "`cleanTemplateOutput()` 仅保留空白行塌陷逻辑" — 当 SurfaceKey 为空时 `just compile {{.SurfaceKey}}` 产生尾部空格。当前靠 cleanTemplateOutput 第三条规则清理。提案说只保留空白行塌陷。 | improvement: 在 Scope 中明确使用 `just compile{{if .SurfaceKey}} {{.SurfaceKey}}{{end}}` 模板模式替代，应用于 6 个相关模板

- **[medium]** injectSurfaceFrontmatter 描述不准确 | quote: "injectSurfaceFrontmatter() 同时执行替换已有字段和插入缺失字段两种行为" — 实际代码（add.go:272-280）只做 strings.Replace 替换，无插入逻辑。描述基于函数注释而非实现。 | improvement: 修正描述为"替换 surface-key: "" 和 surface-type: "" 字面值为推断值"

- **[low]** doc-consolidate.md 的 SCOPE 分类错误 | quote: "`{{SCOPE}}` 两种使用模式（段落块 vs 行内值）...段落级（doc-consolidate, test-gen-contracts, test-gen-journeys, test-run）" — doc-consolidate.md 第 4 行是 `- Scope: {{SCOPE}}`，属行内值模式而非段落级。 | improvement: 将 doc-consolidate 从"段落级"移至"行内值"分类

## BORDERLINE_FINDINGS

(none)

## SKIPPED_FINDINGS

(none)

## Classification Audit

| Triage Layer | Count | Details |
|---|---|---|
| Factual correction | 3 | promptTemplateData 缺字段、injectSurfaceFrontmatter 描述不准确、doc-consolidate SCOPE 分类 |
| Structural/architectural suggestion | 3 | autogenTemplateData 不完整、PHASE_SUMMARY 双块模式、just 尾部空格 |
| Subjective preference | 0 | (none) |
