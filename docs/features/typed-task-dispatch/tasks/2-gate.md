---
id: "2.gate"
title: "Phase 2 Gate: Schema 与模板验证"
priority: "P0"
estimated_time: "30min"
dependencies: ["2.summary"]
status: pending
breaking: true
type: "gate"
---

# 2.gate: Phase 2 Gate — Schema 与模板验证

## Description

Exit verification gate for Phase 2. Confirms that schema, templates, and skill rules are consistent before proceeding to Phase 3 (Agent & Commands).

## Verification Checklist

- [ ] `index.schema.json` 包含 type 枚举（11 个值）和 blockedReason 字段
- [ ] 所有任务模板 frontmatter 包含 type 字段，noTest 字段已移除
- [ ] breakdown-tasks SKILL.md 包含完整 Type Assignment 规则表
- [ ] quick-tasks SKILL.md 包含相同 Type Assignment 规则表
- [ ] 用 breakdown-tasks 生成一个测试 feature 的 index.json，`task validate` 无报错
- [ ] 生成的 index.json 中所有任务均含 type 字段，值与任务类型一致

## Acceptance Criteria

- [ ] 所有上述检查项通过
- [ ] `task validate docs/features/typed-task-dispatch/tasks/index.json` 无报错
