---
id: "12"
title: "Fix gen-contracts: surface detection, Convention loading, surfaceType naming"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 12: Fix gen-contracts: surface detection, Convention loading, surfaceType naming

## Description
gen-contracts 是唯一不使用 `forge surfaces` CLI 进行 surface 检测的 skill（与 gen-journeys 的 HARD-RULE 矛盾），且 Convention 加载路径仍为旧 framework-first 扁平结构。同时 risk-density.md 使用 "WebUI" 应为 "Web"，journey-contract-model.md 使用 "CLI/API/TUI/UI/Mobile" 应为 "CLI/API/TUI/Web/Mobile"。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution, Out of Scope
- `plugins/forge/skills/gen-contracts/SKILL.md:56,59,80,165`: surface 检测和 Convention 加载逻辑 (ref: Proposed Solution)
- `plugins/forge/skills/gen-contracts/rules/risk-density.md:56`: "WebUI" → "Web" (ref: Proposed Solution)
- `plugins/forge/skills/gen-contracts/rules/journey-contract-model.md:155`: "UI" → "Web" in type list (ref: Proposed Solution)
- `plugins/forge/skills/gen-journeys/SKILL.md:26-57`: gen-journeys 的 `forge surfaces` 使用方式作为参考 (ref: Proposed Solution)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-contracts/SKILL.md` | surface 检测改用 `forge surfaces` CLI；Convention 加载改为 surface-first 路径 `testing/{surface}/core.md` |
| `plugins/forge/skills/gen-contracts/rules/risk-density.md` | "WebUI" → "Web" |
| `plugins/forge/skills/gen-contracts/rules/journey-contract-model.md` | "CLI/API/TUI/UI/Mobile" → "CLI/API/TUI/Web/Mobile" |

## Acceptance Criteria
- [ ] gen-contracts SKILL.md surface 检测改用 `forge surfaces <path>` CLI，不再自行从目录结构推断
- [ ] gen-contracts SKILL.md Convention 加载路径改为 `docs/conventions/testing/{surface}/core.md`，有 legacy fallback
- [ ] `risk-density.md` 中 "WebUI" 改为 "Web"
- [ ] `journey-contract-model.md` 中 interface type 列表 "UI" 改为 "Web"
- [ ] 所有 surfaceType 统一使用 web/api/cli/tui/mobile

## Hard Rules
- 必须先加载 `docs/conventions/forge-distribution.md`
- surface 检测必须使用 `forge surfaces` CLI，禁止自行从目录结构推断
- 不引用其它 skill 的内部文件

## Implementation Notes
- gen-contracts 的 surface 检测参考 gen-journeys SKILL.md:26-57 的实现方式：`forge surfaces <path>` + exit code 契约
- Convention 加载参考 gen-test-scripts SKILL.md 的 surface-first 路径逻辑
- 检测到旧 Convention 结构时输出迁移提示（与其它 skill 一致）
