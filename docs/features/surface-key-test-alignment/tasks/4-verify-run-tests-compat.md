---
id: "4"
title: "Verify run-tests SKILL.md compatibility"
priority: "P1"
estimated_time: "30min"
dependencies: [1]
type: "doc"
mainSession: false
---

# 4: Verify run-tests SKILL.md compatibility

## Description
run-tests 已使用 per-surface-key expansion，需验证其 SKILL.md 中的 journey discovery 和 recipe 调用路径与新目录结构兼容。

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
| plugins/forge/skills/run-tests/SKILL.md | 如有不兼容则更新 journey discovery 逻辑 |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] run-tests SKILL.md 的 journey discovery 逻辑与 `tests/<surfaceKey>/<journey>/` 目录结构兼容
- [ ] run-tests recipe 调用路径在单 surface 和多 surface 项目中均正确

## Implementation Notes
run-tests 已使用 per-surface-key，大概率兼容。重点检查 journey discovery（如何找到 tests/ 下的测试文件）和 recipe 调用（just 命令如何传递路径参数）是否需要适配。如无需修改，仅需记录验证结果。
