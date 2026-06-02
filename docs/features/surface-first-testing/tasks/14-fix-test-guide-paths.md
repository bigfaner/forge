---
id: "14"
title: "Fix test-guide test directory paths in SKILL.md + surface templates"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 14: Fix test-guide test directory paths in SKILL.md + surface templates

## Description
test-guide 中测试目录路径存在两处不一致：(1) SKILL.md Step 4a 速查表中 web/mobile 行的文件位置仍为 `tests/e2e/`，但实际 gen-test-scripts 输出到 `tests/<journey>/`；(2) CLI/API/TUI 的 surface template 文件位置段落保留了 `tests/<surface>/` 作为备选路径，与 pipeline 实际输出 `tests/<journey>/` 矛盾。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution
- `plugins/forge/skills/test-guide/SKILL.md:169,171`: 速查表 web/mobile 行 `tests/e2e/` (ref: Proposed Solution)
- `plugins/forge/skills/test-guide/templates/surfaces/cli.md`: 文件位置双路径 (ref: Proposed Solution)
- `plugins/forge/skills/test-guide/templates/surfaces/api.md`: 文件位置双路径
- `plugins/forge/skills/test-guide/templates/surfaces/tui.md`: 文件位置双路径

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/test-guide/SKILL.md` | 速查表 web/mobile 行文件位置从 `tests/e2e/` 改为 `tests/<journey>/`；CLI/API/TUI 行从 `tests/<surface>/` 改为 `tests/<journey>/` |
| `plugins/forge/skills/test-guide/templates/surfaces/cli.md` | 文件位置移除 `tests/<surface>/` 选项，统一为 `tests/<journey>/` |
| `plugins/forge/skills/test-guide/templates/surfaces/api.md` | 同上 |
| `plugins/forge/skills/test-guide/templates/surfaces/tui.md` | 同上 |

## Acceptance Criteria
- [ ] SKILL.md 速查表中所有 5 个 surface 行的文件位置统一为 `tests/<journey>/`
- [ ] CLI template 文件位置移除 `tests/<surface>/` 备选，统一为 `tests/<journey>/`
- [ ] API template 文件位置移除 `tests/<surface>/` 备选，统一为 `tests/<journey>/`
- [ ] TUI template 文件位置移除 `tests/<surface>/` 备选，统一为 `tests/<journey>/`
- [ ] 无残留的 `tests/e2e/` 或 `tests/<surface>/` 引用（web/mobile template 已确认正确）

## Hard Rules
- 必须先加载 `docs/conventions/forge-distribution.md`
- 测试目录统一规则：所有 surface 的测试代码都生成到 `tests/<journey>/`

## Implementation Notes
- journey 名称由 gen-journeys 生成，不是 surface type 名称
- web/mobile template 文件已经是 `tests/<journey>/`，只需确认无需修改
