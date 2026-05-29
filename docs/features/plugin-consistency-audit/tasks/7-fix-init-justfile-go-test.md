---
id: "7"
title: "Fix: init-justfile go.just test recipe uses Node.js command"
priority: "P0"
estimated_time: "15min"
dependencies: []
type: "doc"
complexity: "low"
mainSession: false
---

# 7: Fix: init-justfile go.just test recipe uses Node.js command

## Description
`plugins/forge/skills/init-justfile/templates/go.just` 的 `test` recipe 使用 `npx playwright test`（Node.js/Playwright 命令），Go 项目运行时会导致测试失败。应替换为 Go 适用的 test 命令。(Source: C-24, Report 04)

## Reference Files
- `docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md#C-24`: P0 级发现，go.just test recipe 使用 Node.js 命令 (source: Report 04)
- `plugins/forge/skills/init-justfile/templates/go.just`: 需修复的文件，找到 `test` recipe 中的 `npx playwright test` (source: audit finding)
- `plugins/forge/skills/init-justfile/templates/node.just`: 参考 Node.js 模板的 test recipe 格式 (source: cross-reference)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/init-justfile/templates/go.just` | 将 `test` recipe 中的 `npx playwright test` 替换为 Go test 命令 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `go.just` 的 `test` recipe 不再包含 `npx playwright test` 或任何 Node.js 命令
- [ ] 替换后的命令使用 Go 生态的 test 命令（如 `go test ./...`），或引用 Convention 定义的 test runner
- [ ] `go.just` 中无其他残留的 Node.js/Playwright 引用

## Hard Rules
- 仅修改 `plugins/forge/skills/init-justfile/templates/go.just`

## Implementation Notes
- 参考该 skill 的 Convention 加载机制，Go 项目应使用 `go test` 或 Convention 中定义的 test runner
- 检查 go.just 中是否还有其他 Node.js 遗留命令（如 `npx` 开头的命令）
