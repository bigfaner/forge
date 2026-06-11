---
id: "8"
title: "Fix: run-tests env-check.md Playwright hardcodes"
priority: "P1"
estimated_time: "30min"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 8: Fix: run-tests env-check.md Playwright hardcodes

## Description
`plugins/forge/skills/run-tests/rules/env-check.md` 第 49 行硬编码 `npx playwright install`，与 SKILL.md 的 profile-agnostic 设计直接矛盾。v3.0.0 test profile 系统已将 Playwright 替换为可插拔 Convention-based frameworks。需替换为 Convention-derived framework commands。(Source: RT-01, Report 02)

## Reference Files
- `docs/features/plugin-consistency-audit/reports/02-skills-batch-a.md#RT-01`: P1 级 CONFLICT，env-check.md 硬编码 Playwright (source: Report 02)
- `plugins/forge/skills/run-tests/rules/env-check.md`: 需修复的文件，Web surface check #3 的 Verify/Repair 列 (source: audit finding)
- `plugins/forge/skills/run-tests/SKILL.md`: SKILL.md 的 profile-agnostic 设计作为权威参考 (source: audit finding)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/run-tests/rules/env-check.md` | 替换 Web surface check 中的 Playwright 硬编码为 Convention-derived 命令 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Web surface environment check 中不再硬编码 `npx playwright install`
- [ ] 替换为 Convention-derived 的通用描述，如 "Run the browser automation framework install command (per Convention file)"
- [ ] Repair Suggestion 列同样替换为通用描述
- [ ] 文件中无其他 Playwright 硬编码残留（全局搜索 `playwright` 验证）

## Hard Rules
- 仅修改 `plugins/forge/skills/run-tests/rules/env-check.md`

## Implementation Notes
- RT-02 (env-check.md L13) 同时发现 env-check 引用了 gen-journeys 的 surface rules 而非 run-tests 自身的，此问题优先级较低(P2)，不在本任务范围内
