---
id: "5"
title: "Update init-justfile SKILL.md and surface rules"
priority: "P1"
estimated_time: "1h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 5: Update init-justfile SKILL.md and surface rules

## Description
init-justfile 的 SKILL.md 和 5 个 surface rule 文件定义了 justfile recipe 的生成逻辑。需确认 recipe 前缀逻辑和路径定义与 surface-key 目录结构对齐。

## Reference Files
- `docs/proposals/surface-key-test-alignment/proposal.md` — Proposed Solution, Scope, Key Risks

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/init-justfile/SKILL.md | 确认 recipe 前缀逻辑与新目录对齐 |
| plugins/forge/skills/init-justfile/rules/surfaces/api.md | 确认 recipe 定义兼容 |
| plugins/forge/skills/init-justfile/rules/surfaces/cli.md | 确认 recipe 定义兼容 |
| plugins/forge/skills/init-justfile/rules/surfaces/mobile.md | 确认 recipe 定义兼容 |
| plugins/forge/skills/init-justfile/rules/surfaces/tui.md | 确认 recipe 定义兼容 |
| plugins/forge/skills/init-justfile/rules/surfaces/web.md | 确认 recipe 定义兼容 |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] init-justfile SKILL.md recipe 前缀逻辑与 surface-key 目录结构对齐
- [ ] 5 个 surface rule 文件的 recipe 定义在 `tests/<surfaceKey>/<journey>/` 路径下正确工作

## Implementation Notes
init-justfile 的 recipe 前缀已使用 surface-key（如 `backend-test`），但 recipe 内部的 `tests/` 路径可能仍假设扁平结构。需检查 recipe 中 journey filter 路径是否需要包含 surfaceKey 层级。
