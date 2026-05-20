---
feature: "prompt-template-audit"
status: completed
mode: quick
---

# Feature (Quick): prompt-template-audit

<!-- Status flow: tasks -> in-progress -> completed -->

## Documents

| Document | Path |
|----------|------|
| Proposal | ../../proposals/prompt-template-audit/proposal.md |
| Test Cases | testing/test-cases.md |

## Tasks

| ID | Title | Status | File |
|----|-------|--------|------|
| 1 | P0: 移除6个模板的显式submit步骤（双重提交） | pending | tasks/1-fix-double-submit.md |
| 2 | P1: 修复 coding-fix 模板缺失和顺序问题 | pending | tasks/2-fix-coding-fix-template.md |
| 3 | P1: 修复 gate 模板缺失 conventions 和 PHASE_SUMMARY | pending | tasks/3-fix-gate-template.md |
| 4 | P1: test.* 模板 HARD-RULE 标签重命名为 TASK-CONSTRAINTS | pending | tasks/4-rename-hard-rule-tags.md |
| 5 | P1: 清理未使用占位符声明 + coding-refactor 格式统一 | pending | tasks/5-cleanup-redundant-declarations.md |
| 6 | P2: 增强模板指令精确度 | pending | tasks/6-enhance-template-precision.md |
| 7 | P3: coding-refactor fmt处理简化 + coding-enhancement同步注释 | pending | tasks/7-p3-template-misc-optimization.md |
| 8 | P3: prompt.go 添加占位符 escaping 警告注释 | pending | tasks/8-prompt-go-comments.md |
| T-quick-doc-drift | Detect Spec Drift | pending | tasks/quick-drift-detection.md |
