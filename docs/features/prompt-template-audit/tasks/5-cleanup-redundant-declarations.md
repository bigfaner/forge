---
id: "5"
title: "P1: 清理未使用占位符声明 + coding-refactor 格式统一"
priority: "P1"
estimated_time: "20m"
dependencies: []
type: "doc"
mainSession: false
---

# 5: P1: 清理未使用占位符声明 + coding-refactor 格式统一

## Description
清理模板中声明但未使用的占位符变量，减少维护混淆：(1) doc.md 和 doc-eval.md 声明了 FEATURE_SLUG 但从未在正文中使用；(2) doc-consolidate.md 和 doc-drift.md 声明了 SCOPE 但未消费（consolidate 不跑 just 命令）。同时将 coding-refactor.md 的 CODING_PRINCIPLES 从 ### 标题格式统一为 - 列表格式。

## Reference Files
- `docs/proposals/prompt-template-audit/proposal.md` — Source proposal (Sections 1.2, 1.5, P1 #7/#8/#9)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/doc.md` | 移除 header 中 FEATURE_SLUG 声明行 |
| `forge-cli/pkg/prompt/data/doc-eval.md` | 移除 header 中 FEATURE_SLUG 声明行 |
| `forge-cli/pkg/prompt/data/doc-consolidate.md` | 移除 header 中 SCOPE 声明行 |
| `forge-cli/pkg/prompt/data/doc-drift.md` | 移除 header 中 SCOPE 声明行 |
| `forge-cli/pkg/prompt/data/coding-refactor.md` | CODING_PRINCIPLES 条目从 ### 标题格式改为 - 列表格式 |

## Acceptance Criteria
- [ ] doc.md 和 doc-eval.md 的 header 不再包含 FEATURE_SLUG 行
- [ ] doc-consolidate.md 和 doc-drift.md 的 header 不再包含 SCOPE 行
- [ ] coding-refactor.md 的 CODING_PRINCIPLES 使用 - 列表格式（与 coding-feature 一致）
- [ ] 所有修改不影响模板的实际输出（移除的是声明但未使用的变量）

## Implementation Notes
- FEATURE_SLUG 在 doc-summary.md 中有实际使用（构建 records 路径），不做修改
- SCOPE 在 coding.* 和 validation.* 模板中有实际使用，不做修改
