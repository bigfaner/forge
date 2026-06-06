---
id: "2"
title: "Update gen-test-scripts SKILL.md and type templates"
priority: "P1"
estimated_time: "1h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 2: Update gen-test-scripts SKILL.md and type templates

## Description
gen-test-scripts 的 SKILL.md 和 6 个 type 模板文件中引用的输出目录规则需要从 `tests/<journey>/` 更新为自适应规则：多 surface → `tests/<surfaceKey>/<journey>/`，单 surface → `tests/<journey>/`。

## Reference Files
- `docs/proposals/surface-key-test-alignment/proposal.md` — Proposed Solution, Requirements Analysis, Key Scenarios, Scope, Success Criteria

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/gen-test-scripts/SKILL.md | 输出目录规则更新为 surface-key 自适应 |
| plugins/forge/skills/gen-test-scripts/types/_shared.md | 更新多 surface 输出目录指导 |
| plugins/forge/skills/gen-test-scripts/types/api.md | 更新输出目录路径描述 |
| plugins/forge/skills/gen-test-scripts/types/cli.md | 更新输出目录路径描述 |
| plugins/forge/skills/gen-test-scripts/types/mobile.md | 更新输出目录路径描述 |
| plugins/forge/skills/gen-test-scripts/types/tui.md | 更新输出目录路径描述 |
| plugins/forge/skills/gen-test-scripts/types/web.md | 更新输出目录路径描述 |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] SKILL.md 输出目录规则更新：多 surface → `tests/<surfaceKey>/<journey>/`，单 surface → `tests/<journey>/`
- [ ] 所有 5 个 surface type 文件（api.md, cli.md, mobile.md, tui.md, web.md）输出目录描述与 SKILL.md 一致
- [ ] _shared.md 包含多 surface 输出目录指导

## Implementation Notes
SKILL.md 是 agent 执行任务时加载的主要指令文件，输出目录描述的准确性直接影响生成文件的位置。需确保 agent 能根据 surface 数量（单/多）正确选择目录层级。
