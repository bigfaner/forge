---
status: "completed"
started: "2026-05-20 13:39"
completed: "2026-05-20 13:44"
time_spent: "~5m"
---

# Task Record: 1 Slim consolidate-specs (607→≤350 lines)

## Summary
将 consolidate-specs/SKILL.md 从 607 行精简至 348 行（减少 43%）。按 Splitting Heuristic 规则拆分：步骤编号+条件分支+I/O 契约保留在 SKILL.md，规则细节移至 rules/，输出模板移至 templates/。

## Changes

### Files Created
- plugins/forge/skills/consolidate-specs/rules/constraints.md
- plugins/forge/skills/consolidate-specs/rules/domain-frontmatter.md
- plugins/forge/skills/consolidate-specs/rules/overlap-detection.md
- plugins/forge/skills/consolidate-specs/rules/spec-classification.md
- plugins/forge/skills/consolidate-specs/templates/biz-specs.md
- plugins/forge/skills/consolidate-specs/templates/tech-specs.md
- plugins/forge/skills/consolidate-specs/templates/review-choices.md
- plugins/forge/skills/consolidate-specs/templates/markers.md
- plugins/forge/skills/consolidate-specs/templates/vocabulary-index.md
- plugins/forge/skills/consolidate-specs/templates/commit-messages.md

### Files Modified
- plugins/forge/skills/consolidate-specs/SKILL.md

### Key Decisions
- Domain Frontmatter 规则（约25行）移至 rules/domain-frontmatter.md
- CROSS/LOCAL 分类标准 + Project-Global ID 编码移至 rules/spec-classification.md
- 重叠检测规则 + domain-to-decision-file 映射表移至 rules/overlap-detection.md
- Rules 部分约23条约束规则移至 rules/constraints.md
- 6 个输出模板（biz-specs, tech-specs, review-choices, markers, vocabulary-index, commit-messages）移至 templates/
- SKILL.md 中所有引用使用 ${CLAUDE_SKILL_DIR} 路径变量，确保分发后路径正确

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] SKILL.md 行数 <= 350 行
- [x] 所有步骤编号及其描述保留在 SKILL.md 中（13步）
- [x] 所有条件分支逻辑保留在 SKILL.md 中
- [x] 输入/输出契约定义保留在 SKILL.md 中
- [x] SKILL.md 中引用的所有 rules/templates 路径存在且文件可读
- [x] 无流程步骤遗漏

## Notes
无
