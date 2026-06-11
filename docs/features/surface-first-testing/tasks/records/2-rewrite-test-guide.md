---
status: "completed"
started: "2026-06-02 21:54"
completed: "2026-06-02 21:56"
time_spent: "~2m"
---

# Task Record: 2 重写 test-guide skill

## Summary
全量重写 test-guide SKILL.md：将主流程从 framework-first 改为 surface-first，新增 Surface 配置读取、per-surface 模板渲染、顶层速查表生成步骤；框架检测降格为辅助步骤仅填充断言偏好表。同步更新 convention-structure.md 为 surface-first 目录结构。

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/test-guide/SKILL.md
- plugins/forge/skills/test-guide/rules/convention-structure.md

### Key Decisions
无

## Document Metrics
SKILL.md ~200 lines, convention-structure.md ~120 lines, 5 AC all met

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- docs/conventions/forge-distribution.md
- plugins/forge/skills/test-guide/rules/signal-detection.md
- plugins/forge/skills/test-guide/templates/surfaces/cli.md
- plugins/forge/skills/test-guide/templates/surfaces/api.md
- plugins/forge/skills/test-guide/templates/surfaces/web.md
- plugins/forge/skills/test-guide/references/test-type-model.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md 包含读取 .forge/config.yaml surfaces 配置的步骤
- [x] SKILL.md 包含从 templates/surfaces/*.md 生成 per-surface convention 文件（index.md + core.md）的步骤
- [x] SKILL.md 包含生成顶层 docs/conventions/testing/index.md 速查表的步骤
- [x] 框架检测重构为辅助步骤，结果仅用于填充 core.md 断言偏好表
- [x] 旧的 framework-first 流程已移除

## Notes
SKILL.md 路径使用相对路径引用 rules/ 和 templates/ 目录，符合 forge-distribution.md 分发路径规范。旧 convention-template.md 保留但不再被新流程引用（非本次任务范围）。
