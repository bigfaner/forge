---
id: "14"
title: "Fix: gen-journeys test level emphasis mismatch"
priority: "P1"
estimated_time: "30min"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 14: Fix: gen-journeys test level emphasis mismatch

## Description
gen-journeys SKILL.md "Per-Surface Rule Application" 节中，5 个 surface type 的 test level emphasis 描述与对应 rules/surface-*.md 不一致。具体：(1) API: SKILL.md 说 "integration-heavy ratio"，rules 说 "Balanced 50/50"；(2) Web: SKILL.md 说 "e2e-heavy ratio"，rules 说 "Balanced 50/50"；(3) TUI: SKILL.md 说 "integration-heavy ratio"，rules 说 "Contract 80% / Journey smoke 20%"。以 rules 文件为准。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md#C-15`: P1 CONFLICT, 3/5 surfaces 不一致 (source: Report 04)
- `plugins/forge/skills/gen-journeys/SKILL.md`: 需修改的 "Per-Surface Rule Application" 节 (source: audit finding)
- `plugins/forge/skills/gen-journeys/rules/surface-api.md`: API surface 权威定义 (source: rules file)
- `plugins/forge/skills/gen-journeys/rules/surface-web.md`: Web surface 权威定义 (source: rules file)
- `plugins/forge/skills/gen-journeys/rules/surface-tui.md`: TUI surface 权威定义 (source: rules file)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-journeys/SKILL.md` | 更新 Per-Surface Rule Application 节中 API/Web/TUI 的 test level emphasis 描述 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] API surface 的 test level emphasis 与 rules/surface-api.md 一致（"Balanced 50/50"）
- [ ] Web surface 的 test level emphasis 与 rules/surface-web.md 一致（"Balanced 50/50"）
- [ ] TUI surface 的 test level emphasis 与 rules/surface-tui.md 一致（"Contract 80% / Journey smoke 20%"）
- [ ] CLI 和 Mobile surface 无变化（已一致）

## Hard Rules
- 仅修改 `plugins/forge/skills/gen-journeys/SKILL.md`
- 以 rules 文件为准，不要修改 rules 文件

## Implementation Notes
- 逐一对比 5 个 surface rules 文件中的 test level/ratio 描述，确保 SKILL.md 的摘要准确反映
