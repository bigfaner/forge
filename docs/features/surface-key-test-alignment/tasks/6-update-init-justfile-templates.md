---
id: "6"
title: "Update init-justfile justfile templates"
priority: "P1"
estimated_time: "1h"
dependencies: [5]
type: "doc"
mainSession: false
---

# 6: Update init-justfile justfile templates

## Description
init-justfile 的 6 个 justfile 模板文件（.just）中的 `tests/` 路径需要适配多 surface 目录结构。

## Reference Files
- `docs/proposals/surface-key-test-alignment/proposal.md` — Proposed Solution, Key Scenarios, Scope

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/init-justfile/templates/generic.just | tests/ 路径适配多 surface |
| plugins/forge/skills/init-justfile/templates/go.just | tests/ 路径适配多 surface |
| plugins/forge/skills/init-justfile/templates/mixed.just | tests/ 路径适配多 surface |
| plugins/forge/skills/init-justfile/templates/node.just | tests/ 路径适配多 surface |
| plugins/forge/skills/init-justfile/templates/python.just | tests/ 路径适配多 surface |
| plugins/forge/skills/init-justfile/templates/rust.just | tests/ 路径适配多 surface |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] 所有 6 个 justfile 模板的测试路径在多 surface 场景下使用 `tests/<surfaceKey>/<journey>/`
- [ ] 单 surface 场景下仍使用 `tests/<journey>/`（无额外目录层）
- [ ] 模板中的 recipe 参数传递与新目录结构一致

## Implementation Notes
justfile 模板是纯文本模板（非 Markdown），使用 just 语法。需注意模板中是否支持单/多 surface 的路径区分。如果模板不支持条件逻辑，可能需要在 SKILL.md 中通过指令控制生成哪种路径格式。
