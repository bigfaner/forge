---
id: "3"
title: "P1: 修复 gate 模板缺失 conventions 和 PHASE_SUMMARY"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 3: P1: 修复 gate 模板缺失 conventions 和 PHASE_SUMMARY

## Description
gate.md 缺少两个重要功能：(1) conventions 加载步骤——gate 检查应了解项目标准才能验证交付物；(2) PHASE_SUMMARY——gate 检查可能需要了解前序阶段的决策上下文。

## Reference Files
- `docs/proposals/prompt-template-audit/proposal.md` — Source proposal (Sections 1.3, 2.20)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/gate.md` | 添加 PHASE_SUMMARY header 声明和 Step 1 条件加载；添加 conventions 加载步骤 |

## Acceptance Criteria
- [ ] gate.md header 包含 PHASE_SUMMARY 声明
- [ ] Step 1 包含 PHASE_SUMMARY 条件加载语句
- [ ] Step 1 包含 `docs/conventions/` 和 `docs/business-rules/` 的加载指令（按 domains 字段过滤相关性）
- [ ] conventions 加载方式与其他模板一致（读 frontmatter domains 字段判断相关性）

## Implementation Notes
- 参考 coding-feature.md 或 validation-code.md 的 conventions 加载写法
- PHASE_SUMMARY 为空时由 cleanTemplateOutput 自动清理
