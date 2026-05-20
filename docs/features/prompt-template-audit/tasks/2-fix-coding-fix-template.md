---
id: "2"
title: "P1: 修复 coding-fix 模板缺失和顺序问题"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 2: P1: 修复 coding-fix 模板缺失和顺序问题

## Description
coding-fix.md 存在两个问题：(1) 缺少 PHASE_SUMMARY header 声明和 Step 1 的条件加载，修复任务可能需要前序阶段上下文；(2) conventions 加载在 task 文件读取之后（与其他 coding 模板相反），应统一顺序。

## Reference Files
- `docs/proposals/prompt-template-audit/proposal.md` — Source proposal (Sections 1.2, 2.5)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/coding-fix.md` | 添加 PHASE_SUMMARY header 声明；将 conventions 加载移到 task 文件读取之前 |

## Acceptance Criteria
- [ ] coding-fix.md header 包含 PHASE_SUMMARY 声明
- [ ] Step 1 中包含 PHASE_SUMMARY 条件加载语句
- [ ] Step 1 中 conventions 加载在 task 文件读取之前（与其他 coding.* 模板一致）
- [ ] 模板结构与 coding-feature/enhancement/cleanup 保持一致的 conventions 加载模式

## Implementation Notes
- 参考 coding-feature.md 的 PHASE_SUMMARY 使用方式
- PHASE_SUMMARY 为空时由 cleanTemplateOutput 自动清理，无需额外处理
