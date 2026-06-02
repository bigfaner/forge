---
id: "6"
title: "更新 init-justfile Test recipe"
priority: "P1"
estimated_time: "1h"
dependencies: [2]
type: "doc"
mainSession: false
---

# 6: 更新 init-justfile Test recipe

## Description
更新 init-justfile 的 Test recipe 生成逻辑，使其根据新的 surface-first 目录结构生成 recipe。Test recipe 命名和路径需遵循 Surface type 约定。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution (下游 skill 适配)
- `plugins/forge/skills/init-justfile/SKILL.md`: Test recipe 生成逻辑需修改 (ref: Proposed Solution)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/init-justfile/SKILL.md` | Test recipe 根据新目录结构生成 |

## Acceptance Criteria
- [ ] SKILL.md Test recipe 根据新目录结构 `testing/{surface}/` 生成
- [ ] Test recipe 命名遵循 Surface type 约定（如 `test-cli`、`test-api`）

## Hard Rules
- 必须先加载 `docs/conventions/forge-distribution.md` 了解 plugin 分发模型

## Implementation Notes
- init-justfile 的 test recipe 生成需要感知 Surface 类型，与 `.forge/config.yaml` 的 surfaces 配置对齐
