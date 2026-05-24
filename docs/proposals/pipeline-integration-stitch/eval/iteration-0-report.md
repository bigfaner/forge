# Eval Report: Pre-Revision (Freeform Findings)

**Iteration**: 0
**Title**: Pre-Revision (Freeform Findings)

## ATTACK_POINTS

- **[high]** gen-and-run 移除范围不完整，生产代码文件数被低估 | quote: "The proposal should provide a complete file-by-file checklist with line numbers for the gen-and-run removal. The current '5 files ~15 places' for production code undercounts" | improvement: provide complete file-by-file checklist with line numbers for gen-and-run removal
- **[high]** gen-and-run 部分移除会导致编译失败（引用已删除常量） | quote: "if infer.go:32-33 (the T-quick-gen-and-run case) is removed but types.go:55 (TypeTestGenAndRun) is not, the code will not compile because the case references a deleted constant" | improvement: specify removal order that ensures incremental compilation
- **[medium]** 提案未指定 4 个新 prompt 模板的具体内容，eval 模板与现有模式根本不同 | quote: "the proposal does not specify the actual content of the 4 new templates" | improvement: specify key content patterns for each template type (test-gen vs eval)
- **[medium]** 提案建议参考的 code-quality-simplify.md 与 eval 任务结构无关 | quote: "code-quality-simplify.md is a coding-category task with coverage injection and quality gate workflow -- structurally irrelevant to eval tasks" | improvement: update risk mitigation to reference eval-relevant template pattern
- **[medium]** RecordData eval 字段命名使用 eval 前缀与现有无前缀惯例不一致 | quote: "The field names use an eval prefix (evalScore, evalFindings, evalSeverity, evalPassed). But the existing category-specific fields in RecordData do NOT use a category prefix" | improvement: use unprefixed field names consistent with existing convention
- **[medium]** RecordData 新增 4 个 eval 字段未指定 JSON tag 和 omitempty 行为 | quote: "The proposal introduces 4 new struct fields without specifying their JSON tag names, omitempty behavior, or how they interact with existing fields like Summary and Notes" | improvement: specify field details or note they follow existing pattern
- **[medium]** findFirstTestTaskIdx 修复不完整，应与 ResolveFirstTestDep 使用相同发现机制 | quote: "The proposal should ensure both functions use the same discovery mechanism" | improvement: specify using findTaskIndexByPrefix(tasks, "T-test-gen-journeys") for consistency
- **[medium]** RenderRecord 范围应同时涵盖 RecordTemplateData 和 NewRecordTemplateData 更新 | quote: "the scope entry should also mention updating RecordTemplateData to include eval-specific formatted fields and NewRecordTemplateData to populate them" | improvement: expand scope entry to include RecordTemplateData
- **[medium]** 依赖注入合并函数未指定签名，needsEval 参数处理不明 | quote: "The merged function must handle the case where needsEval is false. The proposal does not specify the function signature" | improvement: specify merged function signature with needsEval parameter
- **[medium]** record-format-test.md 列出从未注册的幻影类型 | quote: "test.gen-cases and test.eval-cases are not even registered in ValidTypes -- they appear to be phantom types that were never implemented" | improvement: specify exact correct type list for replacement
- **[medium]** CategoryForType default 分支会导致未来新类型被静音误分类 | quote: "any future type that lacks a matching prefix will still silently fall into CategoryCoding. This is how the eval types originally got misclassified" | improvement: address default branch hazard with logging or sentinel value
- **[medium]** grep 成功标准应排除历史文档目录 | quote: "The grep command should specify file type exclusions (e.g., --exclude-dir=docs/proposals)" | improvement: qualify grep success criterion to exclude historical docs

## BORDERLINE_FINDINGS

- **[medium]** 向 RecordData 添加 eval 特有字段会造成字段膨胀 | concern is legitimate but solution requires architectural discussion beyond proposal scope
- **[low]** validate_index.go 迁移错误信息位置未指定 | implementation detail, not proposal-level concern

## SKIPPED_FINDINGS

- (none classified as subjective preference — all findings reference verifiable gaps)

## Classification Audit

- Factual correction: 9 findings
- Structural/architectural suggestion: 3 findings
- Subjective preference: 0 findings

## Rubric

(Non-participatory for pre-revision cycle)
