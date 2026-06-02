---
id: "7"
title: "清理旧文件并重新生成 Forge 项目 conventions"
priority: "P2"
estimated_time: "1h"
dependencies: [2, 3, 4, 5, 6]
type: "doc"
mainSession: false
---

# 7: 清理旧文件并重新生成 Forge 项目 conventions

## Description
删除 Forge 项目自身的旧 convention 文件（6 个框架文件 + index.md）和旧 test-type-model（内容已迁入 plugin）。然后用新的 test-guide skill 重新生成 Forge 项目自身的 `docs/conventions/testing/cli/`。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Scope (Forge 项目层), Success Criteria

## Affected Files

### Create
| File | Description |
|------|-------------|
| `docs/conventions/testing/cli/index.md` | CLI convention 索引（由 test-guide 生成） |
| `docs/conventions/testing/cli/core.md` | CLI 测试策略（由 test-guide 生成） |
| `docs/conventions/testing/index.md` | 顶层速查表（由 test-guide 生成） |

### Delete
| File | Reason |
|------|--------|
| `docs/conventions/testing/ginkgo.md` | 旧框架文件，已由 surface-first 结构替代 |
| `docs/conventions/testing/go.md` | 旧框架文件，已由 surface-first 结构替代 |
| `docs/conventions/testing/junit.md` | 旧框架文件，已由 surface-first 结构替代 |
| `docs/conventions/testing/pytest.md` | 旧框架文件，已由 surface-first 结构替代 |
| `docs/conventions/testing/rust.md` | 旧框架文件，已由 surface-first 结构替代 |
| `docs/conventions/testing/vitest.md` | 旧框架文件，已由 surface-first 结构替代 |
| `docs/conventions/testing/index.md` | 旧速查表，将重新生成 |
| `docs/reference/test-type-model.md` | 内容已迁入 plugin 层 |

## Acceptance Criteria
- [ ] `docs/conventions/testing/` 下 6 个旧框架文件（ginkgo/go/junit/pytest/rust/vitest）已删除
- [ ] `docs/reference/test-type-model.md` 已删除
- [ ] Forge 项目自身 `docs/conventions/testing/cli/` 已用新 test-guide 重新生成（含 index.md + core.md）
- [ ] 顶层 `docs/conventions/testing/index.md` 已重新生成

## Implementation Notes
- 这是 Forge 项目自身的清理（不分发到用户项目），确保自洽性
- 删除前确认 test-type-model.md 的内容已完整迁入 `plugins/forge/skills/test-guide/references/test-type-model.md`
- 重新生成使用新的 test-guide skill，而非手动创建
