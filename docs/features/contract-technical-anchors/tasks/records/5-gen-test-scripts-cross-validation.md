---
status: "completed"
started: "2026-06-05 17:42"
completed: "2026-06-05 17:48"
time_spent: "~6m"
---

# Task Record: 5 gen-test-scripts 增加交叉验证和 surface 覆盖报告

## Summary
为 gen-test-scripts SKILL.md 增加 Step 1.5 交叉验证步骤和 surface 覆盖报告，创建 rules/step-1.5-cross-validation.md 详细规则文件，更新 quality-gates.md 错误处理表

## Changes

### Files Created
- plugins/forge/skills/gen-test-scripts/rules/step-1.5-cross-validation.md

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md
- plugins/forge/skills/gen-test-scripts/rules/quality-gates.md

### Key Decisions
无

## Document Metrics
SKILL.md +120 lines (Step 1.5 section), step-1.5-cross-validation.md ~170 lines, quality-gates.md +6 error entries. All 6 AC items PASS.

## Referenced Documents
- docs/proposals/contract-technical-anchors/proposal.md
- plugins/forge/skills/gen-test-scripts/rules/step-0.5-validation.md

## Review Status
final

## Acceptance Criteria
- [x] 交叉验证比对 Fact Table 与 Contract frontmatter 锚点，结果分类为高置信度/低置信度/无法验证
- [x] 不匹配时以 handbook 为权威源生成建议修复，展示 diff 供用户确认后写入 Contract
- [x] 设计文档（handbook）与代码实现不一致时，生成明确的代码 bug 标记报告
- [x] 输出 surface 覆盖报告，明确列出已验证和未验证的 surface 类型
- [x] 缺少 handbook 或锚点字段时，降级为 Fact Table 推断（向后兼容），并提示用户
- [x] 能捕获 lesson 场景（POST vs PUT 不匹配），建议修复为 handbook 定义的 PUT

## Notes
交叉验证步骤插入在 Step 1（代码侦察）和 Step 2（读取 Contract）之间，以 handbook 为 authority source，设计-代码不一致时标记为 code bug。Degradation mode 确保向后兼容。Surface coverage report 输出完整的验证结果分类和未覆盖 surface 列表。
