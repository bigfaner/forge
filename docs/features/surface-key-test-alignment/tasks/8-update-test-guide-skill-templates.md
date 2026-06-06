---
id: "8"
title: "Update test-guide SKILL.md and surface templates"
priority: "P1"
estimated_time: "1h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 8: Update test-guide SKILL.md and surface templates

## Description
test-guide 是生成测试约定文档的 skill。其 SKILL.md、5 个 surface 模板和 references 文件中的测试目录约定需要更新为 surface-key 自适应规则。

## Reference Files
- `docs/proposals/surface-key-test-alignment/proposal.md` — Proposed Solution, Scope, Success Criteria

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/test-guide/SKILL.md | 目录约定部分更新 |
| plugins/forge/skills/test-guide/templates/surfaces/api.md | 测试目录路径对齐 |
| plugins/forge/skills/test-guide/templates/surfaces/cli.md | 测试目录路径对齐 |
| plugins/forge/skills/test-guide/templates/surfaces/mobile.md | 测试目录路径对齐 |
| plugins/forge/skills/test-guide/templates/surfaces/tui.md | 测试目录路径对齐 |
| plugins/forge/skills/test-guide/templates/surfaces/web.md | 测试目录路径对齐 |
| plugins/forge/skills/test-guide/references/test-type-model.md | 目录结构模型更新 |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] test-guide SKILL.md 目录约定部分包含多 surface 规则：`tests/<surfaceKey>/<journey>/`
- [ ] 5 个 surface 模板文件测试目录路径与新规则一致
- [ ] test-type-model.md 目录结构模型反映 surface-key 分区

## Implementation Notes
test-guide 生成的约定文件（如 `docs/conventions/testing/index.md`）被其他 skill 引用。确保模板中的目录描述准确，这样生成的约定文件才能正确指导后续任务。
