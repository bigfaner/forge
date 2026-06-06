---
id: "7"
title: "Update gen-journeys and gen-contracts references"
priority: "P1"
estimated_time: "1h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 7: Update gen-journeys and gen-contracts references

## Description
gen-journeys 和 gen-contracts 的 SKILL.md 及相关文件中引用了 `tests/<journey>/` 路径，需更新为 surface-key 自适应规则。

## Reference Files
- `docs/proposals/surface-key-test-alignment/proposal.md` — Proposed Solution, Scope, Key Scenarios

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/gen-journeys/SKILL.md | 检查 surface-key/type 在输出路径中的使用 |
| plugins/forge/skills/gen-journeys/templates/journey.md | 检查 journey frontmatter 的 surface 信息 |
| plugins/forge/skills/gen-contracts/SKILL.md | 更新 `tests/<journey>/` 引用 |
| plugins/forge/skills/gen-contracts/rules/journey-contract-model.md | 更新 contract-to-directory 映射 |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] gen-journeys SKILL.md 和 journey.md 模板中 surface 信息不依赖 surface-type 作为目录分区依据
- [ ] gen-contracts SKILL.md 中 `tests/<journey>/` 引用更新为自适应规则
- [ ] journey-contract-model.md 的 contract-to-directory 映射反映 surface-key 目录结构

## Implementation Notes
gen-journeys 生成 journey 文档，其 frontmatter 可能包含 `surface_types` 字段。需确认该字段是否影响测试目录的生成。gen-contracts 的 contract 模型中可能有 journey → test file 的路径映射需要更新。
