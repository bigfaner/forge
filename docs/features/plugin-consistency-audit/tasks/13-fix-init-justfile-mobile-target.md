---
id: "13"
title: "Fix: add init-justfile mobile test-setup target"
priority: "P1"
estimated_time: "15min"
dependencies: []
type: "doc"
complexity: "low"
mainSession: false
---

# 13: Fix: add init-justfile mobile test-setup target

## Description
init-justfile SKILL.md 的 Surface-Level Targets 表未包含 mobile surface 的 `<key>-test-setup` target。但 `rules/surfaces/mobile.md` 的编排序列包含 `test-setup -> dev -> probe -> test -> teardown`，有额外的 test-setup 步骤。SKILL.md 表格遗漏了此 target。(Source: C-26, Report 04)

## Reference Files
- `docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md#C-26`: P1 CONFLICT, mobile test-setup target 缺失 (source: Report 04)
- `plugins/forge/skills/init-justfile/SKILL.md`: 需修改的 Surface-Level Targets 表 (source: audit finding)
- `plugins/forge/skills/init-justfile/rules/surfaces/mobile.md`: mobile 编排序列参考 (source: audit finding)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/init-justfile/SKILL.md` | 在 Surface-Level Targets 表中添加 mobile 的 `<key>-test-setup` 行 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Surface-Level Targets 表包含 mobile-specific 的 `<key>-test-setup` 行
- [ ] 行描述注明 "Mobile only: prepare emulator and test environment"
- [ ] 与 `rules/surfaces/mobile.md` 的编排序列一致

## Hard Rules
- 仅修改 `plugins/forge/skills/init-justfile/SKILL.md`
