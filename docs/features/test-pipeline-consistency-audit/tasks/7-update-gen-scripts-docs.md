---
id: "7"
title: "更新 gen-contracts 和 gen-test-scripts Skill 文档术语"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 7: 更新 gen-contracts 和 gen-test-scripts Skill 文档术语

## Description
替换 `gen-contracts/` 和 `gen-test-scripts/` Skill 文档中的旧术语和路径引用："e2e 测试管道" → "Forge 测试管道"；`tests/e2e/` 旧路径 → `tests/<journey>/`；更新 `gen-test-scripts/rules/step-1-contract-loading.md` 中示例路径；更新 `gen-test-scripts/rules/convention-guide.md` 中 "e2e tests" 引用。注意 `gen-contracts/rules/journey-contract-model.md` 第 159 行 "language profile" 属于旧模型对比表，保留不改动。

## Reference Files
- `proposal.md#Layer-2-Skill-文档层术语统一` — 第 9 项定义了术语替换范围和保留例外
- `proposal.md#Scope` — In Scope 第 3 项覆盖 Skill 文档层术语统一

## Acceptance Criteria
- [ ] `gen-contracts/` 下所有 Skill 文档中 "e2e 测试管道" 替换为 "Forge 测试管道"
- [ ] `gen-test-scripts/` 下所有 Skill 文档中 `tests/e2e/` 旧路径替换为 `tests/<journey>/`
- [ ] `gen-test-scripts/rules/step-1-contract-loading.md` 中 `tests/e2e/step1_test.go` 示例路径已更新
- [ ] `gen-test-scripts/rules/convention-guide.md` 中 "e2e tests" 引用已替换
- [ ] `gen-contracts/rules/journey-contract-model.md` 第 159 行 "language profile" 未被修改

## Implementation Notes
- `journey-contract-model.md` 第 159 行的 "language profile" 是旧模型对比表的一部分，是解释性内容而非规范，保留

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-contracts/SKILL.md` | "e2e 测试管道" → "Forge 测试管道" |
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | 旧路径和术语替换 |
| `plugins/forge/skills/gen-test-scripts/rules/step-1-contract-loading.md` | 示例路径 `tests/e2e/step1_test.go` 更新 |
| `plugins/forge/skills/gen-test-scripts/rules/convention-guide.md` | "e2e tests" 引用替换 |

### Delete
| File | Reason |
|------|--------|
