---
id: "3"
title: "Update gen-test-scripts rule files"
priority: "P1"
estimated_time: "1h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 3: Update gen-test-scripts rule files

## Description
gen-test-scripts 的 4 个 rule 文件引用了 `tests/<journey>/` 目录路径，需要更新为 surface-key 自适应规则。

## Reference Files
- `docs/proposals/surface-key-test-alignment/proposal.md` — Proposed Solution, Requirements Analysis, Key Scenarios, Scope

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/gen-test-scripts/rules/step-0.5-validation.md | 目录引用更新 |
| plugins/forge/skills/gen-test-scripts/rules/step-1-contract-loading.md | contract-to-directory 映射对齐 |
| plugins/forge/skills/gen-test-scripts/rules/quality-gates.md | 测试文件路径验证对齐 |
| plugins/forge/skills/gen-test-scripts/rules/run-to-learn.md | 学习引用路径更新 |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] step-0.5-validation.md 目录引用与 SKILL.md 输出规则一致
- [ ] step-1-contract-loading.md 的 contract-to-directory 映射反映 surface-key 目录结构
- [ ] quality-gates.md 测试文件路径验证规则对齐
- [ ] run-to-learn.md 学习引用使用正确路径

## Implementation Notes
需检查 convention-guide.md 是否也引用了 `tests/<journey>/` 路径，如有则一并更新。
